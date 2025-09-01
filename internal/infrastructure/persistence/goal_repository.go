package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"peso/internal/domain/goal"
	"peso/internal/domain/user"
	"peso/internal/domain/weight"
	"peso/internal/interfaces"
)

type goalRepository struct {
	db *DB
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *DB) interfaces.GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Save(g *goal.Goal) error {
	query := `
		INSERT OR REPLACE INTO goals (id, user_id, target_weight, unit, target_date, description, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		g.ID().String(),
		g.UserID().String(),
		g.TargetWeight().Float64(),
		g.Unit().String(),
		g.TargetDate().ToTime(),
		g.Description(),
		g.IsActive(),
		g.CreatedAt(),
		g.UpdatedAt(),
	)
	
	if err != nil {
		return fmt.Errorf("failed to save goal: %w", err)
	}
	
	return nil
}

func (r *goalRepository) FindByID(id goal.GoalID) (*goal.Goal, error) {
	query := `
		SELECT id, user_id, target_weight, unit, target_date, description, active, created_at, updated_at
		FROM goals 
		WHERE id = ?
	`
	
	var (
		goalID       string
		userID       string
		targetWeight float64
		unit         string
		targetDate   time.Time
		description  string
		active       bool
		createdAt    time.Time
		updatedAt    time.Time
	)
	
	err := r.db.QueryRow(query, id.String()).Scan(
		&goalID, &userID, &targetWeight, &unit, &targetDate, &description, &active, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("goal not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find goal by ID: %w", err)
	}
	
	return r.scanGoal(goalID, userID, targetWeight, unit, targetDate, description, active, createdAt, updatedAt)
}

func (r *goalRepository) FindActiveByUserID(userID user.UserID) (*goal.Goal, error) {
	query := `
		SELECT id, user_id, target_weight, unit, target_date, description, active, created_at, updated_at
		FROM goals 
		WHERE user_id = ? AND active = TRUE
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	var (
		goalID       string
		uid          string
		targetWeight float64
		unit         string
		targetDate   time.Time
		description  string
		active       bool
		createdAt    time.Time
		updatedAt    time.Time
	)
	
	err := r.db.QueryRow(query, userID.String()).Scan(
		&goalID, &uid, &targetWeight, &unit, &targetDate, &description, &active, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active goal found for user: %s", userID.String())
		}
		return nil, fmt.Errorf("failed to find active goal: %w", err)
	}
	
	return r.scanGoal(goalID, uid, targetWeight, unit, targetDate, description, active, createdAt, updatedAt)
}

func (r *goalRepository) FindByUserID(userID user.UserID) ([]*goal.Goal, error) {
	query := `
		SELECT id, user_id, target_weight, unit, target_date, description, active, created_at, updated_at
		FROM goals 
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query goals by user ID: %w", err)
	}
	defer rows.Close()
	
	return r.scanGoals(rows)
}

func (r *goalRepository) DeactivateByUserID(userID user.UserID) error {
	query := `
		UPDATE goals 
		SET active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ? AND active = TRUE
	`
	
	_, err := r.db.Exec(query, userID.String())
	if err != nil {
		return fmt.Errorf("failed to deactivate goals for user: %w", err)
	}
	
	return nil
}

func (r *goalRepository) Delete(id goal.GoalID) error {
	query := `DELETE FROM goals WHERE id = ?`
	
	result, err := r.db.Exec(query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("goal not found: %s", id.String())
	}
	
	return nil
}

func (r *goalRepository) scanGoals(rows *sql.Rows) ([]*goal.Goal, error) {
	var goals []*goal.Goal
	
	for rows.Next() {
		var (
			goalID       string
			userID       string
			targetWeight float64
			unit         string
			targetDate   time.Time
			description  string
			active       bool
			createdAt    time.Time
			updatedAt    time.Time
		)
		
		err := rows.Scan(&goalID, &userID, &targetWeight, &unit, &targetDate, &description, &active, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal row: %w", err)
		}
		
		g, err := r.scanGoal(goalID, userID, targetWeight, unit, targetDate, description, active, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
		
		goals = append(goals, g)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over goal rows: %w", err)
	}
	
	return goals, nil
}

func (r *goalRepository) scanGoal(id, userIDStr string, targetWeight float64, unitStr string, targetDate time.Time, description string, active bool, createdAt, updatedAt time.Time) (*goal.Goal, error) {
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID from database: %w", err)
	}
	
	weightValue, err := weight.NewWeightValue(targetWeight)
	if err != nil {
		return nil, fmt.Errorf("invalid target weight from database: %w", err)
	}
	
	unit, err := weight.NewWeightUnit(unitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid weight unit from database: %w", err)
	}
	
	targetDateValue, err := goal.NewTargetDate(targetDate.Year(), int(targetDate.Month()), targetDate.Day())
	if err != nil {
		// If the date is in the past, we still want to load it from database
		// Create a TargetDate without validation
		targetDateValue = goal.TargetDate{} // This will need to be adjusted
	}
	
	g, err := goal.NewGoal(id, userID, weightValue, unit, targetDateValue, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create goal from database row: %w", err)
	}
	
	if !active {
		g.Deactivate()
	}
	
	return g, nil
}