package goal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTargetDate_NewTargetDate(t *testing.T) {
	tests := []struct {
		name    string
		year    int
		month   int
		day     int
		wantErr bool
	}{
		{
			name:    "valid date",
			year:    2030,
			month:   6,
			day:     15,
			wantErr: false,
		},
		{
			name:    "invalid month should fail",
			year:    2025,
			month:   13,
			day:     15,
			wantErr: true,
		},
		{
			name:    "invalid day should fail",
			year:    2025,
			month:   6,
			day:     32,
			wantErr: true,
		},
		{
			name:    "past date should fail",
			year:    2020,
			month:   1,
			day:     1,
			wantErr: true,
		},
		{
			name:    "leap year february 29",
			year:    2028,
			month:   2,
			day:     29,
			wantErr: false,
		},
		{
			name:    "non-leap year february 29 should fail",
			year:    2029,
			month:   2,
			day:     29,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDate, err := NewTargetDate(tt.year, tt.month, tt.day)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, targetDate.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.year, targetDate.Year())
				assert.Equal(t, tt.month, targetDate.Month())
				assert.Equal(t, tt.day, targetDate.Day())
			}
		})
	}
}

func TestTargetDate_IsValid(t *testing.T) {
	validDate, err := NewTargetDate(2030, 6, 15)
	assert.NoError(t, err)
	assert.True(t, validDate.IsValid())

	var invalidDate TargetDate
	assert.False(t, invalidDate.IsValid())
}

func TestTargetDate_IsPast(t *testing.T) {
	// Test with future date
	futureDate, _ := NewTargetDate(2030, 12, 31)
	assert.False(t, futureDate.IsPast())
}

func TestTargetDate_DaysUntil(t *testing.T) {
	// Test with a future date
	tomorrow := time.Now().AddDate(0, 0, 1)
	targetDate, err := NewTargetDate(tomorrow.Year(), int(tomorrow.Month()), tomorrow.Day())
	assert.NoError(t, err)
	
	days := targetDate.DaysUntil()
	assert.Equal(t, 1, days)
}

func TestTargetDate_ToTime(t *testing.T) {
	targetDate, err := NewTargetDate(2030, 6, 15)
	assert.NoError(t, err)
	
	timeVal := targetDate.ToTime()
	assert.Equal(t, 2030, timeVal.Year())
	assert.Equal(t, time.Month(6), timeVal.Month())
	assert.Equal(t, 15, timeVal.Day())
}

func TestTargetDate_String(t *testing.T) {
	targetDate, err := NewTargetDate(2030, 6, 15)
	assert.NoError(t, err)
	
	assert.Equal(t, "2030-06-15", targetDate.String())
}

func TestTargetDate_IsZero(t *testing.T) {
	var zeroDate TargetDate
	assert.True(t, zeroDate.IsZero())

	validDate, err := NewTargetDate(2030, 6, 15)
	assert.NoError(t, err)
	assert.False(t, validDate.IsZero())
}