package goal

import (
	"errors"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"
)

type Goal struct {
	id           GoalID
	userID       user.UserID
	targetWeight weight.WeightValue
	unit         weight.WeightUnit
	targetDate   TargetDate
	description  string
	active       bool
	createdAt    time.Time
	updatedAt    time.Time
}

var (
	ErrEmptyUserID      = errors.New("user ID cannot be empty")
	ErrZeroTargetWeight = errors.New("target weight cannot be zero")
	ErrZeroTargetDate   = errors.New("target date cannot be zero")
)

func NewGoal(id string, userID user.UserID, targetWeight weight.WeightValue, unit weight.WeightUnit, targetDate TargetDate, description string) (*Goal, error) {
	goalID, err := NewGoalID(id)
	if err != nil {
		return nil, err
	}

	if userID.IsEmpty() {
		return nil, ErrEmptyUserID
	}

	if targetWeight.IsZero() {
		return nil, ErrZeroTargetWeight
	}

	if targetDate.IsZero() {
		return nil, ErrZeroTargetDate
	}

	now := time.Now()

	return &Goal{
		id:           goalID,
		userID:       userID,
		targetWeight: targetWeight,
		unit:         unit,
		targetDate:   targetDate,
		description:  description,
		active:       true,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (g *Goal) ID() GoalID {
	return g.id
}

func (g *Goal) UserID() user.UserID {
	return g.userID
}

func (g *Goal) TargetWeight() weight.WeightValue {
	return g.targetWeight
}

func (g *Goal) Unit() weight.WeightUnit {
	return g.unit
}

func (g *Goal) TargetDate() TargetDate {
	return g.targetDate
}

func (g *Goal) Description() string {
	return g.description
}

func (g *Goal) IsActive() bool {
	return g.active
}

func (g *Goal) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Goal) UpdatedAt() time.Time {
	return g.updatedAt
}

func (g *Goal) Deactivate() {
	g.active = false
	g.updatedAt = time.Now()
}

func (g *Goal) Activate() {
	g.active = true
	g.updatedAt = time.Now()
}

func (g *Goal) UpdateDescription(description string) {
	g.description = description
	g.updatedAt = time.Now()
}

func (g *Goal) IsExpired() bool {
	return g.targetDate.IsPast()
}

func (g *Goal) DaysRemaining() int {
	return g.targetDate.DaysUntil()
}
