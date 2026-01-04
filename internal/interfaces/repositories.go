package interfaces

import (
	"time"

	"peso/internal/domain/goal"
	"peso/internal/domain/session"
	"peso/internal/domain/user"
	"peso/internal/domain/weight"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Save(user *user.User) error
	FindByID(id user.UserID) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	FindByName(name string) (*user.User, error)
	FindActive() ([]*user.User, error)
	Exists(id user.UserID) (bool, error)
	EmailExists(email string) (bool, error)
}

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	Save(session *session.Session) error
	FindByToken(token string) (*session.Session, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID user.UserID) error
	DeleteExpired() error
}

// WeightRepository defines the interface for weight persistence
type WeightRepository interface {
	Save(weight *weight.Weight) error
	FindByID(id weight.WeightID) (*weight.Weight, error)
	FindByUserID(userID user.UserID, limit int) ([]*weight.Weight, error)
	FindByUserIDAndPeriod(userID user.UserID, from, to time.Time) ([]*weight.Weight, error)
	FindLatestByUserID(userID user.UserID) (*weight.Weight, error)
	CountByUserIDAndDate(userID user.UserID, date time.Time) (int, error)
	Delete(id weight.WeightID) error
}

// GoalRepository defines the interface for goal persistence
type GoalRepository interface {
	Save(goal *goal.Goal) error
	FindByID(id goal.GoalID) (*goal.Goal, error)
	FindActiveByUserID(userID user.UserID) (*goal.Goal, error)
	FindByUserID(userID user.UserID) ([]*goal.Goal, error)
	DeactivateByUserID(userID user.UserID) error
	Delete(id goal.GoalID) error
}
