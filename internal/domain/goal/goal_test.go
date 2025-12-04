package goal

import (
	"testing"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"
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
				if err == nil {
					t.Error("expected error but got nil")
				}
				if goal != nil {
					t.Errorf("expected nil goal but got %v", goal)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if goal == nil {
					t.Error("expected goal but got nil")
				} else {
					if goal.ID().String() != tt.id {
						t.Errorf("expected ID %s but got %s", tt.id, goal.ID().String())
					}
					if goal.UserID().String() != tt.userID.String() {
						t.Errorf("expected userID %s but got %s", tt.userID.String(), goal.UserID().String())
					}
					if goal.TargetWeight().Float64() != tt.targetWeight.Float64() {
						t.Errorf("expected target weight %f but got %f", tt.targetWeight.Float64(), goal.TargetWeight().Float64())
					}
					if goal.Unit().String() != tt.unit.String() {
						t.Errorf("expected unit %s but got %s", tt.unit.String(), goal.Unit().String())
					}
					if goal.TargetDate().String() != tt.targetDate.String() {
						t.Errorf("expected target date %s but got %s", tt.targetDate.String(), goal.TargetDate().String())
					}
					if goal.Description() != tt.description {
						t.Errorf("expected description %s but got %s", tt.description, goal.Description())
					}
					if !goal.IsActive() {
						t.Error("expected goal to be active")
					}
					if goal.CreatedAt().IsZero() {
						t.Error("expected non-zero CreatedAt")
					}
					if goal.UpdatedAt().IsZero() {
						t.Error("expected non-zero UpdatedAt")
					}
				}
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
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !goal.IsActive() {
		t.Error("expected goal to be active initially")
	}

	originalUpdatedAt := goal.UpdatedAt()
	time.Sleep(time.Millisecond)

	goal.Deactivate()
	if goal.IsActive() {
		t.Error("expected goal to be inactive after deactivate")
	}
	if !goal.UpdatedAt().After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be after original")
	}
}

func TestGoal_Activate(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := NewTargetDate(2030, 12, 31)

	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "Test goal")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	goal.Deactivate()
	if goal.IsActive() {
		t.Error("expected goal to be inactive after deactivate")
	}

	originalUpdatedAt := goal.UpdatedAt()
	time.Sleep(time.Millisecond)

	goal.Activate()
	if !goal.IsActive() {
		t.Error("expected goal to be active after activate")
	}
	if !goal.UpdatedAt().After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be after original")
	}
}

func TestGoal_UpdateDescription(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := NewTargetDate(2030, 12, 31)

	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "Original description")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	originalUpdatedAt := goal.UpdatedAt()
	time.Sleep(time.Millisecond)

	goal.UpdateDescription("Updated description")
	if goal.Description() != "Updated description" {
		t.Errorf("expected description Updated description but got %s", goal.Description())
	}
	if !goal.UpdatedAt().After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be after original")
	}
}

func TestGoal_IsExpired(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")

	// Future goal (not expired)
	futureDate, _ := NewTargetDate(2030, 12, 31)
	futureGoal, err := NewGoal("goal_123", userID, targetWeight, unit, futureDate, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if futureGoal.IsExpired() {
		t.Error("expected future goal to not be expired")
	}
}

func TestGoal_DaysRemaining(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")

	// Goal in 30 days
	futureDate := time.Now().AddDate(0, 0, 30)
	targetDate, _ := NewTargetDate(futureDate.Year(), int(futureDate.Month()), futureDate.Day())
	goal, err := NewGoal("goal_123", userID, targetWeight, unit, targetDate, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	days := goal.DaysRemaining()
	if days != 30 {
		t.Errorf("expected 30 days remaining but got %d", days)
	}
}
