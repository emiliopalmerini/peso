package goal

import (
	"errors"
	"fmt"
	"time"
)

type TargetDate struct {
	year  int
	month int
	day   int
}

var (
	ErrInvalidDate = errors.New("invalid date")
	ErrPastDate    = errors.New("target date cannot be in the past")
)

func NewTargetDate(year, month, day int) (TargetDate, error) {
	// Validate date using time.Date, it will normalize invalid dates
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Check if the normalized date matches what we provided
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return TargetDate{}, ErrInvalidDate
	}

	// Check if date is in the past
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if t.Before(today) {
		return TargetDate{}, ErrPastDate
	}

	return TargetDate{
		year:  year,
		month: month,
		day:   day,
	}, nil
}

func (td TargetDate) Year() int {
	return td.year
}

func (td TargetDate) Month() int {
	return td.month
}

func (td TargetDate) Day() int {
	return td.day
}

func (td TargetDate) IsValid() bool {
	if td.year == 0 && td.month == 0 && td.day == 0 {
		return false
	}

	t := time.Date(td.year, time.Month(td.month), td.day, 0, 0, 0, 0, time.UTC)
	return t.Year() == td.year && int(t.Month()) == td.month && t.Day() == td.day
}

func (td TargetDate) IsPast() bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	targetTime := td.ToTime()

	return targetTime.Before(today)
}

func (td TargetDate) DaysUntil() int {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	targetTime := td.ToTime()

	diff := targetTime.Sub(today)
	return int(diff.Hours() / 24)
}

func (td TargetDate) ToTime() time.Time {
	return time.Date(td.year, time.Month(td.month), td.day, 0, 0, 0, 0, time.UTC)
}

func (td TargetDate) String() string {
	return fmt.Sprintf("%02d/%02d/%04d", td.day, td.month, td.year)
}

func (td TargetDate) IsZero() bool {
	return td.year == 0 && td.month == 0 && td.day == 0
}
