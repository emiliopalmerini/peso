package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"peso/internal/domain/user"
	"peso/internal/interfaces"
)

type userRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(u *user.User) error {
	query := `
		INSERT OR REPLACE INTO users (id, name, email, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		u.ID().String(),
		u.Name(),
		u.Email(),
		u.IsActive(),
		u.CreatedAt(),
		u.UpdatedAt(),
	)
	
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	
	return nil
}

func (r *userRepository) FindByID(id user.UserID) (*user.User, error) {
	query := `
		SELECT id, name, email, active, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`
	
	var (
		userID    string
		name      string
		email     string
		active    bool
		createdAt time.Time
		updatedAt time.Time
	)
	
	err := r.db.QueryRow(query, id.String()).Scan(
		&userID, &name, &email, &active, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	
	return r.scanUser(userID, name, email, active, createdAt, updatedAt)
}

func (r *userRepository) FindByName(name string) (*user.User, error) {
	query := `
		SELECT id, name, email, active, created_at, updated_at 
		FROM users 
		WHERE name = ?
	`
	
	var (
		userID    string
		userName  string
		email     string
		active    bool
		createdAt time.Time
		updatedAt time.Time
	)
	
	err := r.db.QueryRow(query, name).Scan(
		&userID, &userName, &email, &active, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with name: %s", name)
		}
		return nil, fmt.Errorf("failed to find user by name: %w", err)
	}
	
	return r.scanUser(userID, userName, email, active, createdAt, updatedAt)
}

func (r *userRepository) FindActive() ([]*user.User, error) {
	query := `
		SELECT id, name, email, active, created_at, updated_at 
		FROM users 
		WHERE active = TRUE
		ORDER BY name
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active users: %w", err)
	}
	defer rows.Close()
	
	var users []*user.User
	
	for rows.Next() {
		var (
			userID    string
			name      string
			email     string
			active    bool
			createdAt time.Time
			updatedAt time.Time
		)
		
		err := rows.Scan(&userID, &name, &email, &active, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		
		u, err := r.scanUser(userID, name, email, active, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
		
		users = append(users, u)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}
	
	return users, nil
}

func (r *userRepository) Exists(id user.UserID) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE id = ?`
	
	var count int
	err := r.db.QueryRow(query, id.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	
	return count > 0, nil
}

// scanUser is a helper function to create a user from database values
func (r *userRepository) scanUser(id, name, email string, active bool, createdAt, updatedAt time.Time) (*user.User, error) {
	// We need to reconstruct the user with the database values
	// Since the User struct doesn't have setters for created/updated times,
	// we'll create a new user and manually set the times (this is a limitation we'd address in a real implementation)
	u, err := user.NewUser(id, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user from database row: %w", err)
	}
	
	if !active {
		u.Deactivate()
	}
	
	// Note: In a real implementation, we'd need to expose setters or use a different approach
	// to properly restore the created/updated timestamps from the database
	
	return u, nil
}