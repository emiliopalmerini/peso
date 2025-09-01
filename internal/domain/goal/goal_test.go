package goal

import (
	"testing"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"

	"github.com/stretchr/testify/assert"
)

func TestGoal_NewGoal(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := NewTargetDate(2030, 12, 31)

	tests := []struct {
		name         string
		id           string
		userID       user.UserID
		targetWeight weight.WeightValue
		unit         weight.WeightUnit
		targetDate   TargetDate
		description  string
		wantErr      bool
	}{
		{
			name:         "valid goal",
			id:           "goal_123",
			userID:       userID,
			targetWeight: targetWeight,
			unit:         unit,
			targetDate:   targetDate,
			description:  "Lose weight for summer",
			wantErr:      false,
		},
		{
			name:         "valid goal without description",
			id:           "goal_124",
			userID:       userID,
			targetWeight: targetWeight,
			unit:         unit,
			targetDate:   targetDate,
			description:  "",
			wantErr:      false,
		},
		{
			name:         "invalid goal id",
			id:           "",
			userID:       userID,
			targetWeight: targetWeight,
			unit:         unit,
			targetDate:   targetDate,
			description:  "",
			wantErr:      true,
		},
		{
			name:         "empty user id",
			id:           "goal_125",
			userID:       user.UserID(""),
			targetWeight: targetWeight,
			unit:         unit,
			targetDate:   targetDate,
			description:  "",
			wantErr:      true,
		},
		{
			name:         "zero target weight",
			id:           "goal_126",
			userID:       userID,
			targetWeight: weight.WeightValue(0),
			unit:         unit,
			targetDate:   targetDate,
			description:  "",
			wantErr:      true,
		},
		{
			name:         "zero target date",
			id:           "goal_127",
			userID:       userID,
			targetWeight: targetWeight,
			unit:         unit,
			targetDate:   TargetDate{},
			description:  "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goal, err := NewGoal(tt.id, tt.userID, tt.targetWeight, tt.unit, tt.targetDate, tt.description)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, goal)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, goal)
				assert.Equal(t, tt.id, goal.ID().String())
				assert.Equal(t, tt.userID.String(), goal.UserID().String())
				assert.Equal(t, tt.targetWeight.Float64(), goal.TargetWeight().Float64())
				assert.Equal(t, tt.unit.String(), goal.Unit().String())
				assert.Equal(t, tt.targetDate.String(), goal.TargetDate().String())
				assert.Equal(t, tt.description, goal.Description())
				assert.True(t, goal.IsActive())
				assert.False(t, goal.CreatedAt().IsZero())
				assert.False(t, goal.UpdatedAt().IsZero())
			}
		})
	}
}

func TestGoal_Deactivate(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := NewTargetDate(2030, 12, 31)

	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "Test goal")
	assert.NoError(t, err)
	
	assert.True(t, goal.IsActive())
	
	originalUpdatedAt := goal.UpdatedAt()
	time.Sleep(time.Millisecond)
	
	goal.Deactivate()
	assert.False(t, goal.IsActive())
	assert.True(t, goal.UpdatedAt().After(originalUpdatedAt))
}

func TestGoal_Activate(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := NewTargetDate(2030, 12, 31)

	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "Test goal")
	assert.NoError(t, err)
	
	goal.Deactivate()
	assert.False(t, goal.IsActive())
	
	originalUpdatedAt := goal.UpdatedAt()
	time.Sleep(time.Millisecond)
	
	goal.Activate()
	assert.True(t, goal.IsActive())
	assert.True(t, goal.UpdatedAt().After(originalUpdatedAt))
}

func TestGoal_UpdateDescription(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := NewTargetDate(2030, 12, 31)

	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "Original description")
	assert.NoError(t, err)
	
	originalUpdatedAt := goal.UpdatedAt()
	time.Sleep(time.Millisecond)
	
	goal.UpdateDescription("Updated description")
	assert.Equal(t, "Updated description", goal.Description())
	assert.True(t, goal.UpdatedAt().After(originalUpdatedAt))
}

func TestGoal_IsExpired(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	
	// Future goal (not expired)
	futureDate, _ := NewTargetDate(2030, 12, 31)
	futureGoal, err := NewGoal("goal_123", userID, targetWeight, unit, futureDate, "")
	assert.NoError(t, err)
	assert.False(t, futureGoal.IsExpired())
}

func TestGoal_DaysRemaining(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	
	// Goal in 30 days
	futureDate := time.Now().AddDate(0, 0, 30)
	targetDate, _ := NewTargetDate(futureDate.Year(), int(futureDate.Month()), futureDate.Day())
	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "")
	assert.NoError(t, err)
	
	days := goal.DaysRemaining()
	assert.Equal(t, 30, days)
}