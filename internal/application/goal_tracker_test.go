package application

import (
	"errors"
	"testing"
	"time"

	"peso/internal/domain/goal"
	"peso/internal/domain/user"
	"peso/internal/domain/weight"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGoalRepository struct {
	mock.Mock
}

func (m *MockGoalRepository) Save(g *goal.Goal) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockGoalRepository) FindByID(id goal.GoalID) (*goal.Goal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goal.Goal), args.Error(1)
}

func (m *MockGoalRepository) FindActiveByUserID(userID user.UserID) (*goal.Goal, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goal.Goal), args.Error(1)
}

func (m *MockGoalRepository) FindByUserID(userID user.UserID) ([]*goal.Goal, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goal.Goal), args.Error(1)
}

func (m *MockGoalRepository) DeactivateByUserID(userID user.UserID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockGoalRepository) Delete(id goal.GoalID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGoalTracker_SetGoal(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWeightRepo := new(MockWeightRepository)
	mockGoalRepo := new(MockGoalRepository)
	
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 12, 31)
	description := "Summer goal"
	
	testUser, _ := user.NewUser("giada", "Giada", "")
	currentWeight, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(70.0)), unit, time.Now(), "")
	
	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr bool
		errorMsg    string
	}{
		{
			name: "successful goal setting",
			setupMocks: func() {
				mockUserRepo.On("FindByID", userID).Return(testUser, nil)
				mockWeightRepo.On("FindLatestByUserID", userID).Return(currentWeight, nil)
				mockGoalRepo.On("FindActiveByUserID", userID).Return(nil, errors.New("no active goal"))
				mockGoalRepo.On("Save", mock.AnythingOfType("*goal.Goal")).Return(nil)
			},
			expectedErr: false,
		},
		{
			name: "user not found",
			setupMocks: func() {
				mockUserRepo.On("FindByID", userID).Return(nil, errors.New("user not found"))
			},
			expectedErr: true,
			errorMsg:    "user not found",
		},
		{
			name: "user inactive",
			setupMocks: func() {
				inactiveUser, _ := user.NewUser("giada", "Giada", "")
				inactiveUser.Deactivate()
				mockUserRepo.On("FindByID", userID).Return(inactiveUser, nil)
			},
			expectedErr: true,
			errorMsg:    "user is not active",
		},
		{
			name: "no current weight found",
			setupMocks: func() {
				mockUserRepo.On("FindByID", userID).Return(testUser, nil)
				mockWeightRepo.On("FindLatestByUserID", userID).Return(nil, errors.New("no weight found"))
			},
			expectedErr: true,
			errorMsg:    "no current weight found",
		},
		{
			name: "target weight same as current",
			setupMocks: func() {
				sameWeight, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(65.0)), unit, time.Now(), "")
				mockUserRepo.On("FindByID", userID).Return(testUser, nil)
				mockWeightRepo.On("FindLatestByUserID", userID).Return(sameWeight, nil)
			},
			expectedErr: true,
			errorMsg:    "target weight must be different from current weight",
		},
		{
			name: "already has active goal",
			setupMocks: func() {
				existingGoal, _ := goal.NewGoal("g1", userID, targetWeight, unit, targetDate, "existing")
				mockUserRepo.On("FindByID", userID).Return(testUser, nil)
				mockWeightRepo.On("FindLatestByUserID", userID).Return(currentWeight, nil)
				mockGoalRepo.On("FindActiveByUserID", userID).Return(existingGoal, nil)
			},
			expectedErr: true,
			errorMsg:    "user already has an active goal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUserRepo.ExpectedCalls = nil
			mockWeightRepo.ExpectedCalls = nil
			mockGoalRepo.ExpectedCalls = nil
			
			tt.setupMocks()
			
			tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)
			
			result, err := tracker.SetGoal(userID, targetWeight, unit, targetDate, description)
			
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, userID.String(), result.UserID().String())
				assert.Equal(t, targetWeight.Float64(), result.TargetWeight().Float64())
			}
			
			mockUserRepo.AssertExpectations(t)
			mockWeightRepo.AssertExpectations(t)
			mockGoalRepo.AssertExpectations(t)
		})
	}
}

func TestGoalTracker_CalculateProgress(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWeightRepo := new(MockWeightRepository)
	mockGoalRepo := new(MockGoalRepository)
	
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 6, 15) // Future date
	
	testGoal, _ := goal.NewGoal("g1", userID, targetWeight, unit, targetDate, "test goal")
	currentWeight, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(68.0)), unit, time.Now(), "")

	tests := []struct {
		name                string
		setupMocks         func()
		expectedErr        bool
		expectedWeightToLose float64
		expectedIsOnTrack   bool
	}{
		{
			name: "successful progress calculation - losing weight",
			setupMocks: func() {
				mockGoalRepo.On("FindActiveByUserID", userID).Return(testGoal, nil)
				mockWeightRepo.On("FindLatestByUserID", userID).Return(currentWeight, nil)
			},
			expectedErr:         false,
			expectedWeightToLose: 3.0, // 68 - 65 = 3kg to lose
			expectedIsOnTrack:   true,
		},
		{
			name: "no active goal",
			setupMocks: func() {
				mockGoalRepo.On("FindActiveByUserID", userID).Return(nil, errors.New("no active goal"))
			},
			expectedErr: true,
		},
		{
			name: "no current weight",
			setupMocks: func() {
				mockGoalRepo.On("FindActiveByUserID", userID).Return(testGoal, nil)
				mockWeightRepo.On("FindLatestByUserID", userID).Return(nil, errors.New("no weight found"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockGoalRepo.ExpectedCalls = nil
			mockWeightRepo.ExpectedCalls = nil
			
			tt.setupMocks()
			
			tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)
			
			progress, err := tracker.CalculateProgress(userID)
			
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, progress.Goal)
				assert.Equal(t, tt.expectedWeightToLose, progress.WeightToLose.Float64())
				assert.Equal(t, tt.expectedIsOnTrack, progress.IsOnTrack)
				assert.True(t, progress.DaysRemaining > 0)
			}
			
			mockGoalRepo.AssertExpectations(t)
			mockWeightRepo.AssertExpectations(t)
		})
	}
}

func TestGoalTracker_GetActiveGoal(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWeightRepo := new(MockWeightRepository)
	mockGoalRepo := new(MockGoalRepository)
	
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 6, 15)
	
	testGoal, _ := goal.NewGoal("g1", userID, targetWeight, unit, targetDate, "test goal")
	
	mockGoalRepo.On("FindActiveByUserID", userID).Return(testGoal, nil)
	
	tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)
	
	foundGoal, err := tracker.GetActiveGoal(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, foundGoal)
	assert.Equal(t, testGoal.ID().String(), foundGoal.ID().String())
	
	mockGoalRepo.AssertExpectations(t)
}

func TestGoalTracker_DeactivateGoal(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWeightRepo := new(MockWeightRepository)
	mockGoalRepo := new(MockGoalRepository)
	
	goalID, _ := goal.NewGoalID("g1")
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 6, 15)
	
	testGoal, _ := goal.NewGoal(goalID.String(), userID, targetWeight, unit, targetDate, "test goal")
	
	mockGoalRepo.On("FindByID", goalID).Return(testGoal, nil)
	mockGoalRepo.On("Save", mock.AnythingOfType("*goal.Goal")).Return(nil)
	
	tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)
	
	err := tracker.DeactivateGoal(goalID)
	
	assert.NoError(t, err)
	mockGoalRepo.AssertExpectations(t)
}