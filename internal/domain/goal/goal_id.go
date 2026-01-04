package goal

import (
	"errors"
	"strings"
)

type GoalID string

var (
	ErrEmptyGoalID = errors.New("goal ID cannot be empty")
)

func NewGoalID(value string) (GoalID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyGoalID
	}

	return GoalID(trimmed), nil
}

func (g GoalID) String() string {
	return string(g)
}

func (g GoalID) IsEmpty() bool {
	return string(g) == ""
}
