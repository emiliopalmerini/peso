package goal

import (
	"testing"
	"time"
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
				if err == nil {
					t.Error("expected error but got nil")
				}
				if !targetDate.IsZero() {
					t.Error("expected zero targetDate")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if targetDate.Year() != tt.year {
					t.Errorf("expected year %d but got %d", tt.year, targetDate.Year())
				}
				if targetDate.Month() != tt.month {
					t.Errorf("expected month %d but got %d", tt.month, targetDate.Month())
				}
				if targetDate.Day() != tt.day {
					t.Errorf("expected day %d but got %d", tt.day, targetDate.Day())
				}
			}
		})
	}
}

func TestTargetDate_IsValid(t *testing.T) {
	validDate, err := NewTargetDate(2030, 6, 15)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if !validDate.IsValid() {
		t.Error("expected valid date")
	}

	var invalidDate TargetDate
	if invalidDate.IsValid() {
		t.Error("expected invalid date")
	}
}

func TestTargetDate_IsPast(t *testing.T) {
	// Test with future date
	futureDate, _ := NewTargetDate(2030, 12, 31)
	if futureDate.IsPast() {
		t.Error("expected future date to not be past")
	}
}

func TestTargetDate_DaysUntil(t *testing.T) {
	// Test with a future date
	tomorrow := time.Now().AddDate(0, 0, 1)
	targetDate, err := NewTargetDate(tomorrow.Year(), int(tomorrow.Month()), tomorrow.Day())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	days := targetDate.DaysUntil()
	if days != 1 {
		t.Errorf("expected 1 day until but got %d", days)
	}
}

func TestTargetDate_ToTime(t *testing.T) {
	targetDate, err := NewTargetDate(2030, 6, 15)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	timeVal := targetDate.ToTime()
	if timeVal.Year() != 2030 {
		t.Errorf("expected year 2030 but got %d", timeVal.Year())
	}
	if timeVal.Month() != time.Month(6) {
		t.Errorf("expected month 6 but got %d", timeVal.Month())
	}
	if timeVal.Day() != 15 {
		t.Errorf("expected day 15 but got %d", timeVal.Day())
	}
}

func TestTargetDate_String(t *testing.T) {
	targetDate, err := NewTargetDate(2030, 6, 15)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if targetDate.String() != "15/06/2030" {
		t.Errorf("expected 15/06/2030 but got %s", targetDate.String())
	}
}

func TestTargetDate_IsZero(t *testing.T) {
	var zeroDate TargetDate
	if !zeroDate.IsZero() {
		t.Error("expected zero date")
	}

	validDate, err := NewTargetDate(2030, 6, 15)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if validDate.IsZero() {
		t.Error("expected non-zero date")
	}
}
