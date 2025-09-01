package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"
	"peso/internal/interfaces"
)

type weightRepository struct {
	db *DB
}

// NewWeightRepository creates a new weight repository
func NewWeightRepository(db *DB) interfaces.WeightRepository {
	return &weightRepository{db: db}
}

func (r *weightRepository) Save(w *weight.Weight) error {
	query := `
		INSERT OR REPLACE INTO weights (id, user_id, value, unit, measured_at, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		w.ID().String(),
		w.UserID().String(),
		w.Value().Float64(),
		w.Unit().String(),
		w.MeasuredAt(),
		w.Notes(),
		w.CreatedAt(),
	)
	
	if err != nil {
		return fmt.Errorf("failed to save weight: %w", err)
	}
	
	return nil
}

func (r *weightRepository) FindByID(id weight.WeightID) (*weight.Weight, error) {
	query := `
		SELECT id, user_id, value, unit, measured_at, notes, created_at 
		FROM weights 
		WHERE id = ?
	`
	
	var (
		weightID   string
		userID     string
		value      float64
		unit       string
		measuredAt time.Time
		notes      string
		createdAt  time.Time
	)
	
	err := r.db.QueryRow(query, id.String()).Scan(
		&weightID, &userID, &value, &unit, &measuredAt, &notes, &createdAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("weight not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find weight by ID: %w", err)
	}
	
	return r.scanWeight(weightID, userID, value, unit, measuredAt, notes, createdAt)
}

func (r *weightRepository) FindByUserID(userID user.UserID, limit int) ([]*weight.Weight, error) {
	query := `
		SELECT id, user_id, value, unit, measured_at, notes, created_at 
		FROM weights 
		WHERE user_id = ?
		ORDER BY measured_at DESC
		LIMIT ?
	`
	
	rows, err := r.db.Query(query, userID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query weights by user ID: %w", err)
	}
	defer rows.Close()
	
	return r.scanWeights(rows)
}

func (r *weightRepository) FindByUserIDAndPeriod(userID user.UserID, from, to time.Time) ([]*weight.Weight, error) {
	query := `
		SELECT id, user_id, value, unit, measured_at, notes, created_at 
		FROM weights 
		WHERE user_id = ? AND measured_at >= ? AND measured_at <= ?
		ORDER BY measured_at ASC
	`
	
	rows, err := r.db.Query(query, userID.String(), from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query weights by user ID and period: %w", err)
	}
	defer rows.Close()
	
	return r.scanWeights(rows)
}

func (r *weightRepository) FindLatestByUserID(userID user.UserID) (*weight.Weight, error) {
	query := `
		SELECT id, user_id, value, unit, measured_at, notes, created_at 
		FROM weights 
		WHERE user_id = ?
		ORDER BY measured_at DESC
		LIMIT 1
	`
	
	var (
		weightID   string
		uid        string
		value      float64
		unit       string
		measuredAt time.Time
		notes      string
		createdAt  time.Time
	)
	
	err := r.db.QueryRow(query, userID.String()).Scan(
		&weightID, &uid, &value, &unit, &measuredAt, &notes, &createdAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no weight found for user: %s", userID.String())
		}
		return nil, fmt.Errorf("failed to find latest weight: %w", err)
	}
	
	return r.scanWeight(weightID, uid, value, unit, measuredAt, notes, createdAt)
}

func (r *weightRepository) CountByUserIDAndDate(userID user.UserID, date time.Time) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM weights 
		WHERE user_id = ? AND DATE(measured_at) = DATE(?)
	`
	
	var count int
	err := r.db.QueryRow(query, userID.String(), date).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count weights by user and date: %w", err)
	}
	
	return count, nil
}

func (r *weightRepository) Delete(id weight.WeightID) error {
	query := `DELETE FROM weights WHERE id = ?`
	
	result, err := r.db.Exec(query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete weight: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("weight not found: %s", id.String())
	}
	
	return nil
}

func (r *weightRepository) scanWeights(rows *sql.Rows) ([]*weight.Weight, error) {
	var weights []*weight.Weight
	
	for rows.Next() {
		var (
			weightID   string
			userID     string
			value      float64
			unit       string
			measuredAt time.Time
			notes      string
			createdAt  time.Time
		)
		
		err := rows.Scan(&weightID, &userID, &value, &unit, &measuredAt, &notes, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weight row: %w", err)
		}
		
		w, err := r.scanWeight(weightID, userID, value, unit, measuredAt, notes, createdAt)
		if err != nil {
			return nil, err
		}
		
		weights = append(weights, w)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over weight rows: %w", err)
	}
	
	return weights, nil
}

func (r *weightRepository) scanWeight(id, userIDStr string, value float64, unitStr string, measuredAt time.Time, notes string, createdAt time.Time) (*weight.Weight, error) {
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID from database: %w", err)
	}
	
	weightValue, err := weight.NewWeightValue(value)
	if err != nil {
		return nil, fmt.Errorf("invalid weight value from database: %w", err)
	}
	
	unit, err := weight.NewWeightUnit(unitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid weight unit from database: %w", err)
	}
	
	w, err := weight.NewWeight(id, userID, weightValue, unit, measuredAt, notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create weight from database row: %w", err)
	}
	
	return w, nil
}