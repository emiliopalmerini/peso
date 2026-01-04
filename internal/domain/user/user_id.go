package user

import (
	"errors"
	"strings"
)

type UserID string

const maxUserIDLength = 50

var (
	ErrEmptyUserID   = errors.New("user ID cannot be empty")
	ErrUserIDTooLong = errors.New("user ID too long")
)

func NewUserID(value string) (UserID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyUserID
	}

	if len(trimmed) > maxUserIDLength {
		return "", ErrUserIDTooLong
	}

	return UserID(trimmed), nil
}

func (u UserID) String() string {
	return string(u)
}

func (u UserID) IsEmpty() bool {
	return string(u) == ""
}
