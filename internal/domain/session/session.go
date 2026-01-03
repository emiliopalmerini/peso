package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"peso/internal/domain/user"
)

const (
	tokenLength    = 32
	defaultExpiry  = 30 * 24 * time.Hour // 30 days
)

var (
	ErrSessionExpired = errors.New("session has expired")
	ErrInvalidToken   = errors.New("invalid session token")
)

type Session struct {
	id        SessionID
	userID    user.UserID
	token     string
	expiresAt time.Time
	createdAt time.Time
}

func NewSession(userID user.UserID) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Session{
		id:        NewSessionID(),
		userID:    userID,
		token:     token,
		expiresAt: now.Add(defaultExpiry),
		createdAt: now,
	}, nil
}

func ReconstructSession(id SessionID, userID user.UserID, token string, expiresAt, createdAt time.Time) *Session {
	return &Session{
		id:        id,
		userID:    userID,
		token:     token,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}
}

func (s *Session) ID() SessionID {
	return s.id
}

func (s *Session) UserID() user.UserID {
	return s.userID
}

func (s *Session) Token() string {
	return s.token
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsExpired()
}

func generateToken() (string, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
