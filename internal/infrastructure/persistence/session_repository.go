package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"peso/internal/domain/session"
	"peso/internal/domain/user"
	"peso/internal/interfaces"
)

type sessionRepository struct {
	db *DB
}

func NewSessionRepository(db *DB) interfaces.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Save(s *session.Session) error {
	query := `
		INSERT OR REPLACE INTO sessions (id, user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		s.ID().String(),
		s.UserID().String(),
		s.Token(),
		s.ExpiresAt(),
		s.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (r *sessionRepository) FindByToken(token string) (*session.Session, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = ?
	`

	var (
		id        string
		userID    string
		tkn       string
		expiresAt time.Time
		createdAt time.Time
	)

	err := r.db.QueryRow(query, token).Scan(&id, &userID, &tkn, &expiresAt, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	sessionID, err := session.ParseSessionID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	uid, err := user.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return session.ReconstructSession(sessionID, uid, tkn, expiresAt, createdAt), nil
}

func (r *sessionRepository) DeleteByToken(token string) error {
	query := `DELETE FROM sessions WHERE token = ?`

	_, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *sessionRepository) DeleteByUserID(userID user.UserID) error {
	query := `DELETE FROM sessions WHERE user_id = ?`

	_, err := r.db.Exec(query, userID.String())
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

func (r *sessionRepository) DeleteExpired() error {
	query := `DELETE FROM sessions WHERE expires_at < ?`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
