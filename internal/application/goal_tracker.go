package application

import (
	"errors"
	"fmt"
	"math"
	"time"

	"peso/internal/domain/goal"
	"peso/internal/domain/user"
	"peso/internal/domain/weight"
	"peso/internal/interfaces"
)

// GoalProgress represents progress towards a goal
type GoalProgress struct {
	Goal            *goal.Goal
	CurrentWeight   weight.WeightValue
	WeightToLose    weight.WeightValue // Positive = lose weight, negative = gain weight
	DaysRemaining   int
	WeightPerDay    weight.WeightValue // Required weight change per day
	ProgressPercent float64
	IsOnTrack       bool
}

// GoalTracker implements goal tracking business logic
type GoalTracker struct {
	userRepo   interfaces.UserRepository
	weightRepo interfaces.WeightRepository
	goalRepo   interfaces.GoalRepository
}

var (
	ErrNoActiveGoal           = errors.New("no active goal found")
	ErrNoCurrentWeight        = errors.New("no current weight found")
	ErrSameWeight             = errors.New("target weight must be different from current weight")
	ErrActiveGoalExists       = errors.New("user already has an active goal")
	ErrUnrealisticGoal        = errors.New("goal is unrealistic (max 2kg per week)")
	ErrGoalNotFound          = errors.New("goal not found")
)

const (
	minWeightDifference = 0.1  // Minimum difference in kg
	maxWeightLossPerWeek = 2.0  // Maximum realistic weight loss per week in kg
)

// NewGoalTracker creates a new goal tracker service
func NewGoalTracker(userRepo interfaces.UserRepository, weightRepo interfaces.WeightRepository, goalRepo interfaces.GoalRepository) *GoalTracker {
	return &GoalTracker{
		userRepo:   userRepo,
		weightRepo: weightRepo,
		goalRepo:   goalRepo,
	}
}

// SetGoal sets a new goal for a user
func (gt *GoalTracker) SetGoal(userID user.UserID, targetWeight weight.WeightValue, unit weight.WeightUnit, targetDate goal.TargetDate, description string) (*goal.Goal, error) {
	// Verify user exists and is active
	u, err := gt.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUserNotFound, err.Error())
	}
	
	if !u.IsActive() {
		return nil, ErrUserNotActive
	}
	
	// Get current weight
	currentWeightRecord, err := gt.weightRepo.FindLatestByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNoCurrentWeight, err.Error())
	}
	
	// Validate target weight is different from current
	weightDiff := abs(targetWeight.Float64() - currentWeightRecord.Value().Float64())
	if weightDiff < minWeightDifference {
		return nil, ErrSameWeight
	}
	
	// Check if user already has an active goal
	existingGoal, err := gt.goalRepo.FindActiveByUserID(userID)
	if err == nil && existingGoal != nil {
		return nil, ErrActiveGoalExists
	}
	
	// Validate goal is realistic (max 2kg per week)
	daysUntilGoal := targetDate.DaysUntil()
	weeksUntilGoal := float64(daysUntilGoal) / 7.0
	requiredWeightChangePerWeek := abs(weightDiff / weeksUntilGoal)
	
	if requiredWeightChangePerWeek > maxWeightLossPerWeek {
		return nil, ErrUnrealisticGoal
	}
	
	// Generate unique goal ID
	goalID := fmt.Sprintf("goal_%s_%d", userID.String(), time.Now().UnixNano())
	
	// Create goal
	newGoal, err := goal.NewGoal(goalID, userID, targetWeight, unit, targetDate, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}
	
	// Save goal
	if err := gt.goalRepo.Save(newGoal); err != nil {
		return nil, fmt.Errorf("failed to save goal: %w", err)
	}
	
	return newGoal, nil
}

// GetActiveGoal gets the active goal for a user
func (gt *GoalTracker) GetActiveGoal(userID user.UserID) (*goal.Goal, error) {
	activeGoal, err := gt.goalRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNoActiveGoal, err.Error())
	}
	
	return activeGoal, nil
}

// GetStartingWeightForGoal gets the weight closest to when the goal was created
func (gt *GoalTracker) GetStartingWeightForGoal(userID user.UserID, goalCreatedAt time.Time) (*weight.Weight, error) {
	// Look for weights around the goal creation date (±7 days)
	from := goalCreatedAt.AddDate(0, 0, -7)
	to := goalCreatedAt.AddDate(0, 0, 7)
	
	weights, err := gt.weightRepo.FindByUserIDAndPeriod(userID, from, to)
	if err != nil {
		return nil, err
	}
	
	if len(weights) == 0 {
		// Fallback to latest weight before goal creation
		return gt.weightRepo.FindLatestByUserID(userID)
	}
	
	// Find the weight closest to goal creation date
	var closest *weight.Weight
	minDiff := time.Duration(math.MaxInt64)
	
	for _, w := range weights {
		diff := w.MeasuredAt().Sub(goalCreatedAt)
		if diff < 0 {
			diff = -diff
		}
		if diff < minDiff {
			minDiff = diff
			closest = w
		}
	}
	
	return closest, nil
}

// CalculateProgress calculates progress towards the user's active goal
func (gt *GoalTracker) CalculateProgress(userID user.UserID) (GoalProgress, error) {
	// Get active goal
	activeGoal, err := gt.GetActiveGoal(userID)
	if err != nil {
		return GoalProgress{}, err
	}
	
	// Get current weight
	currentWeightRecord, err := gt.weightRepo.FindLatestByUserID(userID)
	if err != nil {
		return GoalProgress{}, fmt.Errorf("%w: %s", ErrNoCurrentWeight, err.Error())
	}
	
	currentWeight := currentWeightRecord.Value()
	targetWeight := activeGoal.TargetWeight()
	daysRemaining := activeGoal.DaysRemaining()
	
	// Calculate weight to lose/gain (positive = need to lose, negative = need to gain)
	weightToLose := currentWeight.Subtract(targetWeight)
	
	// Calculate required weight change per day
	var weightPerDay weight.WeightValue
	if daysRemaining > 0 {
		dailyChange := weightToLose.Float64() / float64(daysRemaining)
		weightPerDay = weight.WeightValue(abs(dailyChange))
	}
	
	// Calculate progress percentage
	// If we need to lose weight, progress = (startWeight - currentWeight) / (startWeight - targetWeight)
	// For simplicity, we'll calculate based on remaining vs total
	var progressPercent float64
	totalWeightChange := abs(weightToLose.Float64())
	if totalWeightChange > 0 {
		// This is a simplified progress calculation
		// In a real app, you'd want to track the starting weight when goal was set
		progressPercent = math.Max(0, (1.0 - totalWeightChange/10.0) * 100) // Simplified calculation
	}
	
	// Determine if on track (simplified: if requiring less than 0.3kg per day)
	isOnTrack := weightPerDay.Float64() <= 0.3
	
	return GoalProgress{
		Goal:            activeGoal,
		CurrentWeight:   currentWeight,
		WeightToLose:    weight.WeightValue(abs(weightToLose.Float64())),
		DaysRemaining:   daysRemaining,
		WeightPerDay:    weightPerDay,
		ProgressPercent: progressPercent,
		IsOnTrack:       isOnTrack,
	}, nil
}

// DeactivateGoal deactivates a specific goal
func (gt *GoalTracker) DeactivateGoal(goalID goal.GoalID) error {
	// Find the goal
	g, err := gt.goalRepo.FindByID(goalID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrGoalNotFound, err.Error())
	}
	
	// Deactivate it
	g.Deactivate()
	
	// Save the updated goal
	if err := gt.goalRepo.Save(g); err != nil {
		return fmt.Errorf("failed to deactivate goal: %w", err)
	}
	
	return nil
}