package session

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrEmptySessionID   = errors.New("session ID cannot be empty")
	ErrInvalidSessionID = errors.New("invalid session ID format")
)

type SessionID struct {
	value string
}

func NewSessionID() SessionID {
	return SessionID{value: uuid.New().String()}
}

func ParseSessionID(id string) (SessionID, error) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return SessionID{}, ErrEmptySessionID
	}

	if _, err := uuid.Parse(trimmed); err != nil {
		return SessionID{}, ErrInvalidSessionID
	}

	return SessionID{value: trimmed}, nil
}

func (id SessionID) String() string {
	return id.value
}

func (id SessionID) IsEmpty() bool {
	return id.value == ""
}
