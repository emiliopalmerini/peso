package application

import (
	"errors"
	"fmt"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"
	"peso/internal/interfaces"
)

// TimePeriod represents different time periods for weight tracking
type TimePeriod int

const (
	TimePeriodLastWeek TimePeriod = iota
	TimePeriodLastMonth
	TimePeriodLast3Months
	TimePeriodLast6Months
	TimePeriodLastYear
	TimePeriodAll
)

// TrendDirection indicates the direction of weight change
type TrendDirection int

const (
	TrendIncreasing TrendDirection = iota
	TrendDecreasing
	TrendStable
	TrendNoData
)

// WeightTrend represents weight change over time
type WeightTrend struct {
	Direction            TrendDirection
	TotalChange          weight.WeightValue
	AverageChangePerWeek float64
	StartWeight          weight.WeightValue
	EndWeight            weight.WeightValue
	DataPoints           int
}

// WeightTracker implements weight tracking business logic
type WeightTracker struct {
    userRepo   interfaces.UserRepository
    weightRepo interfaces.WeightRepository
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user is not active")
	ErrMaxDailyRecordings = errors.New("maximum daily weight recordings exceeded")
)

const maxDailyWeightRecordings = 10

// NewWeightTracker creates a new weight tracker service
func NewWeightTracker(userRepo interfaces.UserRepository, weightRepo interfaces.WeightRepository) *WeightTracker {
	return &WeightTracker{
		userRepo:   userRepo,
		weightRepo: weightRepo,
	}
}

// RecordWeight records a new weight measurement for a user
func (wt *WeightTracker) RecordWeight(userID user.UserID, value weight.WeightValue, unit weight.WeightUnit, measuredAt time.Time, notes string) (*weight.Weight, error) {
	// Verify user exists and is active
	u, err := wt.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUserNotFound, err.Error())
	}
	
	if !u.IsActive() {
		return nil, ErrUserNotActive
	}
	
	// Check daily recording limit
	dayStart := time.Date(measuredAt.Year(), measuredAt.Month(), measuredAt.Day(), 0, 0, 0, 0, measuredAt.Location())
	dailyCount, err := wt.weightRepo.CountByUserIDAndDate(userID, dayStart)
	if err != nil {
		return nil, fmt.Errorf("failed to check daily recording count: %w", err)
	}
	
	if dailyCount >= maxDailyWeightRecordings {
		return nil, ErrMaxDailyRecordings
	}
	
	// Generate a unique ID for the weight record
	weightID := fmt.Sprintf("weight_%s_%d", userID.String(), time.Now().UnixNano())
	
	// Create weight record
	w, err := weight.NewWeight(weightID, userID, value, unit, measuredAt, notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create weight record: %w", err)
	}
	
	// Save to repository
	if err := wt.weightRepo.Save(w); err != nil {
		return nil, fmt.Errorf("failed to save weight record: %w", err)
	}
	
	return w, nil
}

// GetWeightHistory retrieves weight history for a user within a time period
func (wt *WeightTracker) GetWeightHistory(userID user.UserID, period TimePeriod) ([]*weight.Weight, error) {
	from, to := wt.getPeriodBounds(period)
	
	weights, err := wt.weightRepo.FindByUserIDAndPeriod(userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve weight history: %w", err)
	}
	
	return weights, nil
}

// GetRecentWeights retrieves the most recent N weights for a user (descending by date)
func (wt *WeightTracker) GetRecentWeights(userID user.UserID, limit int) ([]*weight.Weight, error) {
    // Verify user exists
    if _, err := wt.userRepo.FindByID(userID); err != nil {
        return nil, fmt.Errorf("%w: %s", ErrUserNotFound, err.Error())
    }
    if limit <= 0 {
        limit = 10
    }
    ws, err := wt.weightRepo.FindByUserID(userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve recent weights: %w", err)
    }
    return ws, nil
}

// GetLatestWeight returns the most recent weight for a user
func (wt *WeightTracker) GetLatestWeight(userID user.UserID) (*weight.Weight, error) {
    // Verify user exists
    if _, err := wt.userRepo.FindByID(userID); err != nil {
        return nil, fmt.Errorf("%w: %s", ErrUserNotFound, err.Error())
    }
    w, err := wt.weightRepo.FindLatestByUserID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve latest weight: %w", err)
    }
    return w, nil
}

// CalculateWeightTrend calculates weight trend over a time period
func (wt *WeightTracker) CalculateWeightTrend(userID user.UserID, period TimePeriod) (WeightTrend, error) {
	weights, err := wt.GetWeightHistory(userID, period)
	if err != nil {
		return WeightTrend{}, err
	}
	
	if len(weights) < 2 {
		return WeightTrend{
			Direction:  TrendNoData,
			DataPoints: len(weights),
		}, nil
	}
	
	// Sort by measurement date (assume repository returns sorted data)
	startWeight := weights[0].Value()
	endWeight := weights[len(weights)-1].Value()
	
	// Calculate total change
	totalChange := endWeight.Subtract(startWeight)
	
	// Calculate average change per week
	firstDate := weights[0].MeasuredAt()
	lastDate := weights[len(weights)-1].MeasuredAt()
	
	daysDiff := lastDate.Sub(firstDate).Hours() / 24
	weeksDiff := daysDiff / 7
	
	var avgChangePerWeek float64
	if weeksDiff > 0 {
		avgChangePerWeek = totalChange.Float64() / weeksDiff
	}
	
	// Determine trend direction
	var direction TrendDirection
	if totalChange.Float64() > 0.1 { // More than 0.1kg increase
		direction = TrendIncreasing
	} else if totalChange.Float64() < -0.1 { // More than 0.1kg decrease
		direction = TrendDecreasing
	} else {
		direction = TrendStable
	}
	
	return WeightTrend{
		Direction:            direction,
		TotalChange:          weight.WeightValue(abs(totalChange.Float64())), // Always positive for display
		AverageChangePerWeek: avgChangePerWeek,
		StartWeight:          startWeight,
		EndWeight:            endWeight,
		DataPoints:           len(weights),
	}, nil
}

// getPeriodBounds returns the time bounds for a given period
func (wt *WeightTracker) getPeriodBounds(period TimePeriod) (from, to time.Time) {
	now := time.Now()
	to = now
	
	switch period {
	case TimePeriodLastWeek:
		from = now.AddDate(0, 0, -7)
	case TimePeriodLastMonth:
		from = now.AddDate(0, -1, 0)
	case TimePeriodLast3Months:
		from = now.AddDate(0, -3, 0)
	case TimePeriodLast6Months:
		from = now.AddDate(0, -6, 0)
	case TimePeriodLastYear:
		from = now.AddDate(-1, 0, 0)
	case TimePeriodAll:
		from = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) // Arbitrary start date
	default:
		from = now.AddDate(0, -1, 0) // Default to last month
	}
	
	return from, to
}

// DeleteWeight removes a weight record
func (wt *WeightTracker) DeleteWeight(userID user.UserID, weightID weight.WeightID) error {
	// Verify user exists
	if _, err := wt.userRepo.FindByID(userID); err != nil {
		return fmt.Errorf("%w: %s", ErrUserNotFound, err.Error())
	}

	// Get weight to verify it belongs to the user
	w, err := wt.weightRepo.FindByID(weightID)
	if err != nil {
		return fmt.Errorf("weight not found: %w", err)
	}

	if w.UserID() != userID {
		return fmt.Errorf("weight does not belong to user")
	}

	if err := wt.weightRepo.Delete(weightID); err != nil {
		return fmt.Errorf("failed to delete weight: %w", err)
	}

	return nil
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
