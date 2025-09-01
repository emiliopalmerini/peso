package weight

import (
	"testing"
	"time"

	"peso/internal/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestWeight_NewWeight(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	weightValue, _ := NewWeightValue(70.5)
	unit, _ := NewWeightUnit("kg")
	measuredAt := time.Now().AddDate(0, 0, -1) // Yesterday

	tests := []struct {
		name        string
		id          string
		userID      user.UserID
		value       WeightValue
		unit        WeightUnit
		measuredAt  time.Time
		notes       string
		wantErr     bool
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
				assert.Error(t, err)
				assert.Nil(t, weight)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, weight)
				assert.Equal(t, tt.id, weight.ID().String())
				assert.Equal(t, tt.userID.String(), weight.UserID().String())
				assert.Equal(t, tt.value.Float64(), weight.Value().Float64())
				assert.Equal(t, tt.unit.String(), weight.Unit().String())
				assert.Equal(t, tt.measuredAt.Truncate(time.Second), weight.MeasuredAt().Truncate(time.Second))
				assert.Equal(t, tt.notes, weight.Notes())
				assert.False(t, weight.CreatedAt().IsZero())
			}
		})
	}
}

func TestWeight_IsRecent(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	weightValue, _ := NewWeightValue(70.5)
	unit, _ := NewWeightUnit("kg")
	
	// Recent weight (today)
	recentWeight, err := NewWeight("weight_123", userID, weightValue, unit, time.Now().AddDate(0, 0, 0), "")
	assert.NoError(t, err)
	assert.True(t, recentWeight.IsRecent())
	
	// Old weight (2 weeks ago)
	oldWeight, err := NewWeight("weight_124", userID, weightValue, unit, time.Now().AddDate(0, 0, -14), "")
	assert.NoError(t, err)
	assert.False(t, oldWeight.IsRecent())
}

func TestWeight_IsSameDay(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	weightValue, _ := NewWeightValue(70.5)
	unit, _ := NewWeightUnit("kg")
	
	today := time.Now()
	weight, err := NewWeight("weight_123", userID, weightValue, unit, today, "")
	assert.NoError(t, err)
	
	assert.True(t, weight.IsSameDay(today))
	assert.False(t, weight.IsSameDay(today.AddDate(0, 0, 1)))
}

func TestWeight_UpdateNotes(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	weightValue, _ := NewWeightValue(70.5)
	unit, _ := NewWeightUnit("kg")
	measuredAt := time.Now().AddDate(0, 0, -1)
	
	weight, err := NewWeight("weight_123", userID, weightValue, unit, measuredAt, "")
	assert.NoError(t, err)
	
	assert.Equal(t, "", weight.Notes())
	
	weight.UpdateNotes("Updated notes")
	assert.Equal(t, "Updated notes", weight.Notes())
}