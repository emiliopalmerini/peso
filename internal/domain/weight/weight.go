package weight

import (
	"errors"
	"time"

	"peso/internal/domain/user"
)

type Weight struct {
	id         WeightID
	userID     user.UserID
	value      WeightValue
	unit       WeightUnit
	measuredAt time.Time
	notes      string
	createdAt  time.Time
}

var (
	ErrEmptyUserID      = errors.New("user ID cannot be empty")
	ErrFutureMeasurement = errors.New("measurement date cannot be in the future")
	ErrZeroWeight       = errors.New("weight value cannot be zero")
)

func NewWeight(id string, userID user.UserID, value WeightValue, unit WeightUnit, measuredAt time.Time, notes string) (*Weight, error) {
	weightID, err := NewWeightID(id)
	if err != nil {
		return nil, err
	}
	
	if userID.IsEmpty() {
		return nil, ErrEmptyUserID
	}
	
	if value.IsZero() {
		return nil, ErrZeroWeight
	}
	
	// Check if measurement is in the future (allow same day)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	if measuredAt.After(today) {
		return nil, ErrFutureMeasurement
	}
	
	return &Weight{
		id:         weightID,
		userID:     userID,
		value:      value,
		unit:       unit,
		measuredAt: measuredAt,
		notes:      notes,
		createdAt:  time.Now(),
	}, nil
}

func (w *Weight) ID() WeightID {
	return w.id
}

func (w *Weight) UserID() user.UserID {
	return w.userID
}

func (w *Weight) Value() WeightValue {
	return w.value
}

func (w *Weight) Unit() WeightUnit {
	return w.unit
}

func (w *Weight) MeasuredAt() time.Time {
	return w.measuredAt
}

func (w *Weight) Notes() string {
	return w.notes
}

func (w *Weight) CreatedAt() time.Time {
	return w.createdAt
}

func (w *Weight) IsRecent() bool {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	return w.measuredAt.After(weekAgo)
}

func (w *Weight) IsSameDay(date time.Time) bool {
	wYear, wMonth, wDay := w.measuredAt.Date()
	dYear, dMonth, dDay := date.Date()
	return wYear == dYear && wMonth == dMonth && wDay == dDay
}

func (w *Weight) UpdateNotes(notes string) {
	w.notes = notes
}