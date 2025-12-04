package weight

import (
	"testing"
	"time"

	"peso/internal/domain/user"
)

func TestWeight_NewWeight(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	weightValue, _ := NewWeightValue(70.5)
	unit, _ := NewWeightUnit("kg")
	measuredAt := time.Now().AddDate(0, 0, -1) // Yesterday

	tests := []struct {
		name       string
		id         string
		userID     user.UserID
		value      WeightValue
		unit       WeightUnit
		measuredAt time.Time
		notes      string
		wantErr    bool
	}{
		{
			name:       "valid weight",
			id:         "weight_123",
			userID:     userID,
			value:      weightValue,
			unit:       unit,
			measuredAt: measuredAt,
			notes:      "Morning weight",
			wantErr:    false,
		},
		{
			name:       "valid weight without notes",
			id:         "weight_124",
			userID:     userID,
			value:      weightValue,
			unit:       unit,
			measuredAt: measuredAt,
			notes:      "",
			wantErr:    false,
		},
		{
			name:       "invalid weight id",
			id:         "",
			userID:     userID,
			value:      weightValue,
			unit:       unit,
			measuredAt: measuredAt,
			notes:      "",
			wantErr:    true,
		},
		{
			name:       "empty user id",
			id:         "weight_125",
			userID:     user.UserID(""),
			value:      weightValue,
			unit:       unit,
			measuredAt: measuredAt,
			notes:      "",
			wantErr:    true,
		},
		{
			name:       "future measurement date should fail",
			id:         "weight_126",
			userID:     userID,
			value:      weightValue,
			unit:       unit,
			measuredAt: time.Now().AddDate(0, 0, 1), // Tomorrow
			notes:      "",
			wantErr:    true,
		},
		{
			name:       "zero weight value",
			id:         "weight_127",
			userID:     userID,
			value:      WeightValue(0),
			unit:       unit,
			measuredAt: measuredAt,
			notes:      "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight, err := NewWeight(tt.id, tt.userID, tt.value, tt.unit, tt.measuredAt, tt.notes)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if weight != nil {
					t.Errorf("expected nil weight but got %v", weight)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if weight == nil {
					t.Error("expected weight but got nil")
				} else {
					if weight.ID().String() != tt.id {
						t.Errorf("expected ID %s but got %s", tt.id, weight.ID().String())
					}
					if weight.UserID().String() != tt.userID.String() {
						t.Errorf("expected userID %s but got %s", tt.userID.String(), weight.UserID().String())
					}
					if weight.Value().Float64() != tt.value.Float64() {
						t.Errorf("expected value %f but got %f", tt.value.Float64(), weight.Value().Float64())
					}
					if weight.Unit().String() != tt.unit.String() {
						t.Errorf("expected unit %s but got %s", tt.unit.String(), weight.Unit().String())
					}
					if weight.MeasuredAt().Truncate(time.Second) != tt.measuredAt.Truncate(time.Second) {
						t.Errorf("expected measuredAt %v but got %v", tt.measuredAt, weight.MeasuredAt())
					}
					if weight.Notes() != tt.notes {
						t.Errorf("expected notes %s but got %s", tt.notes, weight.Notes())
					}
					if weight.CreatedAt().IsZero() {
						t.Error("expected non-zero CreatedAt")
					}
				}
			}
		})
	}
}
