package application

import (
	"errors"
	"testing"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id user.UserID) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByName(name string) (*user.User, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindActive() ([]*user.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserRepository) Exists(id user.UserID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

type MockWeightRepository struct {
	mock.Mock
}

func (m *MockWeightRepository) Save(w *weight.Weight) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *MockWeightRepository) FindByID(id weight.WeightID) (*weight.Weight, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*weight.Weight), args.Error(1)
}

func (m *MockWeightRepository) FindByUserID(userID user.UserID, limit int) ([]*weight.Weight, error) {
	args := m.Called(userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*weight.Weight), args.Error(1)
}

func (m *MockWeightRepository) FindByUserIDAndPeriod(userID user.UserID, from, to time.Time) ([]*weight.Weight, error) {
	args := m.Called(userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*weight.Weight), args.Error(1)
}

func (m *MockWeightRepository) FindLatestByUserID(userID user.UserID) (*weight.Weight, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*weight.Weight), args.Error(1)
}

func (m *MockWeightRepository) CountByUserIDAndDate(userID user.UserID, date time.Time) (int, error) {
	args := m.Called(userID, date)
	return args.Int(0), args.Error(1)
}

func (m *MockWeightRepository) Delete(id weight.WeightID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestWeightTracker_RecordWeight(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockWeightRepo := new(MockWeightRepository)
	
	userID, _ := user.NewUserID("giada")
	value, _ := weight.NewWeightValue(70.5)
	unit, _ := weight.NewWeightUnit("kg")
	measuredAt := time.Now().AddDate(0, 0, -1) // Yesterday
	
	testUser, _ := user.NewUser("giada", "Giada", "")
	
	tests := []struct {
		name         string
		setupMocks   func()
		userID       user.UserID
		value        weight.WeightValue
		unit         weight.WeightUnit
		measuredAt   time.Time
		expectedErr  bool
		errorMessage string
	}{
		{
			name: "successful weight recording",
			setupMocks: func() {
				mockUserRepo.On("FindByID", userID).Return(testUser, nil)
				mockWeightRepo.On("CountByUserIDAndDate", userID, mock.AnythingOfType("time.Time")).Return(2, nil)
				mockWeightRepo.On("Save", mock.AnythingOfType("*weight.Weight")).Return(nil)
			},
			userID:      userID,
			value:       value,
			unit:        unit,
			measuredAt:  measuredAt,
			expectedErr: false,
		},
		{
			name: "user not found",
			setupMocks: func() {
				mockUserRepo.On("FindByID", userID).Return(nil, errors.New("user not found"))
			},
			userID:       userID,
			value:        value,
			unit:         unit,
			measuredAt:   measuredAt,
			expectedErr:  true,
			errorMessage: "user not found",
		},
		{
			name: "user inactive",
			setupMocks: func() {
				inactiveUser, _ := user.NewUser("giada", "Giada", "")
				inactiveUser.Deactivate()
				mockUserRepo.On("FindByID", userID).Return(inactiveUser, nil)
			},
			userID:       userID,
			value:        value,
			unit:         unit,
			measuredAt:   measuredAt,
			expectedErr:  true,
			errorMessage: "user is not active",
		},
		{
			name: "too many weights per day",
			setupMocks: func() {
				mockUserRepo.On("FindByID", userID).Return(testUser, nil)
				mockWeightRepo.On("CountByUserIDAndDate", userID, mock.AnythingOfType("time.Time")).Return(10, nil)
			},
			userID:       userID,
			value:        value,
			unit:         unit,
			measuredAt:   measuredAt,
			expectedErr:  true,
			errorMessage: "maximum daily weight recordings exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUserRepo.ExpectedCalls = nil
			mockWeightRepo.ExpectedCalls = nil
			
			tt.setupMocks()
			
			tracker := NewWeightTracker(mockUserRepo, mockWeightRepo)
			
			result, err := tracker.RecordWeight(tt.userID, tt.value, tt.unit, tt.measuredAt, "")
			
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID.String(), result.UserID().String())
				assert.Equal(t, tt.value.Float64(), result.Value().Float64())
			}
			
			mockUserRepo.AssertExpectations(t)
			mockWeightRepo.AssertExpectations(t)
		})
	}
}

func TestWeightTracker_GetWeightHistory(t *testing.T) {
	mockWeightRepo := new(MockWeightRepository)
	mockUserRepo := new(MockUserRepository)
	
	userID, _ := user.NewUserID("giada")
	period := TimePeriodLastWeek
	
	// Create test weights
	weight1, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(70.0)), must(weight.NewWeightUnit("kg")), time.Now().AddDate(0, 0, -1), "")
	weight2, _ := weight.NewWeight("w2", userID, must(weight.NewWeightValue(69.5)), must(weight.NewWeightUnit("kg")), time.Now().AddDate(0, 0, -2), "")
	expectedWeights := []*weight.Weight{weight1, weight2}
	
	mockWeightRepo.On("FindByUserIDAndPeriod", userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(expectedWeights, nil)
	
	tracker := NewWeightTracker(mockUserRepo, mockWeightRepo)
	
	weights, err := tracker.GetWeightHistory(userID, period)
	
	assert.NoError(t, err)
	assert.Len(t, weights, 2)
	assert.Equal(t, expectedWeights, weights)
	
	mockWeightRepo.AssertExpectations(t)
}

func TestWeightTracker_CalculateWeightTrend(t *testing.T) {
	mockWeightRepo := new(MockWeightRepository)
	mockUserRepo := new(MockUserRepository)
	
	userID, _ := user.NewUserID("giada")
	period := TimePeriodLastMonth
	
	// Create test weights showing a downward trend
	now := time.Now()
	weight1, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(72.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -30), "")
	weight2, _ := weight.NewWeight("w2", userID, must(weight.NewWeightValue(71.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -20), "")
	weight3, _ := weight.NewWeight("w3", userID, must(weight.NewWeightValue(70.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -10), "")
	weight4, _ := weight.NewWeight("w4", userID, must(weight.NewWeightValue(69.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -1), "")
	
	weights := []*weight.Weight{weight1, weight2, weight3, weight4}
	
	mockWeightRepo.On("FindByUserIDAndPeriod", userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(weights, nil)
	
	tracker := NewWeightTracker(mockUserRepo, mockWeightRepo)
	
	trend, err := tracker.CalculateWeightTrend(userID, period)
	
	assert.NoError(t, err)
	assert.Equal(t, TrendDecreasing, trend.Direction)
	assert.Equal(t, 3.0, trend.TotalChange.Float64()) // 72 - 69 = 3kg lost
	assert.True(t, trend.AverageChangePerWeek < 0) // Negative because losing weight
}

// Helper function to avoid repetitive error handling in tests
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}