package application

import (
	"errors"
	"testing"
	"time"

	"peso/internal/domain/goal"
	"peso/internal/domain/user"
	"peso/internal/domain/weight"
)

type MockGoalRepository struct {
	calls map[string][]interface{}
	data  map[string]interface{}
}

func NewMockGoalRepository() *MockGoalRepository {
	return &MockGoalRepository{
		calls: make(map[string][]interface{}),
		data:  make(map[string]interface{}),
	}
}

func (m *MockGoalRepository) Save(g *goal.Goal) error {
	m.calls["Save"] = append(m.calls["Save"], g)
	if err, ok := m.data["SaveError"]; ok {
		return err.(error)
	}
	return nil
}

func (m *MockGoalRepository) FindByID(id goal.GoalID) (*goal.Goal, error) {
	m.calls["FindByID"] = append(m.calls["FindByID"], id)
	if err, ok := m.data["FindByIDError"]; ok {
		return nil, err.(error)
	}
	if g, ok := m.data["FindByIDResult"]; ok {
		return g.(*goal.Goal), nil
	}
	return nil, errors.New("not found")
}

func (m *MockGoalRepository) FindActiveByUserID(userID user.UserID) (*goal.Goal, error) {
	m.calls["FindActiveByUserID"] = append(m.calls["FindActiveByUserID"], userID)
	if err, ok := m.data["FindActiveByUserIDError"]; ok {
		return nil, err.(error)
	}
	if g, ok := m.data["FindActiveByUserIDResult"]; ok {
		return g.(*goal.Goal), nil
	}
	return nil, errors.New("no active goal")
}

func (m *MockGoalRepository) FindByUserID(userID user.UserID) ([]*goal.Goal, error) {
	m.calls["FindByUserID"] = append(m.calls["FindByUserID"], userID)
	if err, ok := m.data["FindByUserIDError"]; ok {
		return nil, err.(error)
	}
	if goals, ok := m.data["FindByUserIDResult"]; ok {
		return goals.([]*goal.Goal), nil
	}
	return nil, errors.New("not found")
}

func (m *MockGoalRepository) DeactivateByUserID(userID user.UserID) error {
	m.calls["DeactivateByUserID"] = append(m.calls["DeactivateByUserID"], userID)
	if err, ok := m.data["DeactivateByUserIDError"]; ok {
		return err.(error)
	}
	return nil
}

func (m *MockGoalRepository) Delete(id goal.GoalID) error {
	m.calls["Delete"] = append(m.calls["Delete"], id)
	if err, ok := m.data["DeleteError"]; ok {
		return err.(error)
	}
	return nil
}

func TestGoalTracker_SetGoal(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockWeightRepo := NewMockWeightRepository()
	mockGoalRepo := NewMockGoalRepository()

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
				mockUserRepo.data["FindByIDResult"] = testUser
				mockWeightRepo.data["FindLatestByUserIDResult"] = currentWeight
				mockGoalRepo.data["FindActiveByUserIDError"] = errors.New("no active goal")
			},
			expectedErr: false,
		},
		{
			name: "user not found",
			setupMocks: func() {
				mockUserRepo.data["FindByIDError"] = errors.New("user not found")
			},
			expectedErr: true,
			errorMsg:    "user not found",
		},
		{
			name: "user inactive",
			setupMocks: func() {
				inactiveUser, _ := user.NewUser("giada", "Giada", "")
				inactiveUser.Deactivate()
				mockUserRepo.data["FindByIDResult"] = inactiveUser
			},
			expectedErr: true,
			errorMsg:    "user is not active",
		},
		{
			name: "no current weight found",
			setupMocks: func() {
				mockUserRepo.data["FindByIDResult"] = testUser
				mockWeightRepo.data["FindLatestByUserIDError"] = errors.New("no weight found")
			},
			expectedErr: true,
			errorMsg:    "no current weight found",
		},
		{
			name: "target weight same as current",
			setupMocks: func() {
				sameWeight, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(65.0)), unit, time.Now(), "")
				mockUserRepo.data["FindByIDResult"] = testUser
				mockWeightRepo.data["FindLatestByUserIDResult"] = sameWeight
			},
			expectedErr: true,
			errorMsg:    "target weight must be different from current weight",
		},
		{
			name: "already has active goal",
			setupMocks: func() {
				existingGoal, _ := goal.NewGoal("g1", userID, targetWeight, unit, targetDate, "existing")
				mockUserRepo.data["FindByIDResult"] = testUser
				mockWeightRepo.data["FindLatestByUserIDResult"] = currentWeight
				mockGoalRepo.data["FindActiveByUserIDResult"] = existingGoal
			},
			expectedErr: true,
			errorMsg:    "user already has an active goal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUserRepo.calls = make(map[string][]interface{})
			mockWeightRepo.calls = make(map[string][]interface{})
			mockGoalRepo.calls = make(map[string][]interface{})
			mockUserRepo.data = make(map[string]interface{})
			mockWeightRepo.data = make(map[string]interface{})
			mockGoalRepo.data = make(map[string]interface{})

			tt.setupMocks()

			tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)

			result, err := tracker.SetGoal(userID, targetWeight, unit, targetDate, description)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q but got %q", tt.errorMsg, err.Error())
				}
				if result != nil {
					t.Errorf("expected nil result but got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected result but got nil")
				}
				if result != nil && result.UserID().String() != userID.String() {
					t.Errorf("expected userID %s but got %s", userID.String(), result.UserID().String())
				}
				if result != nil && result.TargetWeight().Float64() != targetWeight.Float64() {
					t.Errorf("expected target weight %f but got %f", targetWeight.Float64(), result.TargetWeight().Float64())
				}
			}
		})
	}
}

func TestGoalTracker_CalculateProgress(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockWeightRepo := NewMockWeightRepository()
	mockGoalRepo := NewMockGoalRepository()

	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 6, 15) // Future date

	testGoal, _ := goal.NewGoal("g1", userID, targetWeight, unit, targetDate, "test goal")
	currentWeight, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(68.0)), unit, time.Now(), "")

	tests := []struct {
		name                    string
		setupMocks              func()
		expectedErr             bool
		expectedWeightToLose    float64
		expectedIsOnTrack       bool
	}{
		{
			name: "successful progress calculation - losing weight",
			setupMocks: func() {
				mockGoalRepo.data["FindActiveByUserIDResult"] = testGoal
				mockWeightRepo.data["FindLatestByUserIDResult"] = currentWeight
			},
			expectedErr:          false,
			expectedWeightToLose: 3.0, // 68 - 65 = 3kg to lose
			expectedIsOnTrack:    true,
		},
		{
			name: "no active goal",
			setupMocks: func() {
				mockGoalRepo.data["FindActiveByUserIDError"] = errors.New("no active goal")
			},
			expectedErr: true,
		},
		{
			name: "no current weight",
			setupMocks: func() {
				mockGoalRepo.data["FindActiveByUserIDResult"] = testGoal
				mockWeightRepo.data["FindLatestByUserIDError"] = errors.New("no weight found")
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockGoalRepo.calls = make(map[string][]interface{})
			mockWeightRepo.calls = make(map[string][]interface{})
			mockGoalRepo.data = make(map[string]interface{})
			mockWeightRepo.data = make(map[string]interface{})

			tt.setupMocks()

			tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)

			progress, err := tracker.CalculateProgress(userID)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if progress.Goal == nil {
					t.Error("expected goal but got nil")
				}
				if progress.WeightToLose.Float64() != tt.expectedWeightToLose {
					t.Errorf("expected weight to lose %f but got %f", tt.expectedWeightToLose, progress.WeightToLose.Float64())
				}
				if progress.IsOnTrack != tt.expectedIsOnTrack {
					t.Errorf("expected IsOnTrack %v but got %v", tt.expectedIsOnTrack, progress.IsOnTrack)
				}
				if progress.DaysRemaining <= 0 {
					t.Errorf("expected days remaining > 0 but got %d", progress.DaysRemaining)
				}
			}
		})
	}
}

func TestGoalTracker_GetActiveGoal(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockWeightRepo := NewMockWeightRepository()
	mockGoalRepo := NewMockGoalRepository()

	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 6, 15)

	testGoal, _ := goal.NewGoal("g1", userID, targetWeight, unit, targetDate, "test goal")

	mockGoalRepo.data["FindActiveByUserIDResult"] = testGoal

	tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)

	foundGoal, err := tracker.GetActiveGoal(userID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if foundGoal == nil {
		t.Error("expected goal but got nil")
	}
	if foundGoal != nil && foundGoal.ID().String() != testGoal.ID().String() {
		t.Errorf("expected goal ID %s but got %s", testGoal.ID().String(), foundGoal.ID().String())
	}
}

func TestGoalTracker_DeactivateGoal(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockWeightRepo := NewMockWeightRepository()
	mockGoalRepo := NewMockGoalRepository()

	goalID, _ := goal.NewGoalID("g1")
	userID, _ := user.NewUserID("giada")
	targetWeight, _ := weight.NewWeightValue(65.0)
	unit, _ := weight.NewWeightUnit("kg")
	targetDate, _ := goal.NewTargetDate(2030, 6, 15)

	testGoal, _ := goal.NewGoal(goalID.String(), userID, targetWeight, unit, targetDate, "test goal")

	mockGoalRepo.data["FindByIDResult"] = testGoal

	tracker := NewGoalTracker(mockUserRepo, mockWeightRepo, mockGoalRepo)

	err := tracker.DeactivateGoal(goalID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
