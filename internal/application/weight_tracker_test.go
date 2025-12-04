package application

import (
	"errors"
	"testing"
	"time"

	"peso/internal/domain/user"
	"peso/internal/domain/weight"
)

// Mock repositories
type MockUserRepository struct {
	calls map[string][]interface{}
	data  map[string]interface{}
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		calls: make(map[string][]interface{}),
		data:  make(map[string]interface{}),
	}
}

func (m *MockUserRepository) Save(u *user.User) error {
	m.calls["Save"] = append(m.calls["Save"], u)
	if err, ok := m.data["SaveError"]; ok {
		return err.(error)
	}
	return nil
}

func (m *MockUserRepository) FindByID(id user.UserID) (*user.User, error) {
	m.calls["FindByID"] = append(m.calls["FindByID"], id)
	if err, ok := m.data["FindByIDError"]; ok {
		return nil, err.(error)
	}
	if u, ok := m.data["FindByIDResult"]; ok {
		return u.(*user.User), nil
	}
	return nil, errors.New("not found")
}

func (m *MockUserRepository) FindByName(name string) (*user.User, error) {
	m.calls["FindByName"] = append(m.calls["FindByName"], name)
	if err, ok := m.data["FindByNameError"]; ok {
		return nil, err.(error)
	}
	if u, ok := m.data["FindByNameResult"]; ok {
		return u.(*user.User), nil
	}
	return nil, errors.New("not found")
}

func (m *MockUserRepository) FindActive() ([]*user.User, error) {
	m.calls["FindActive"] = append(m.calls["FindActive"], nil)
	if err, ok := m.data["FindActiveError"]; ok {
		return nil, err.(error)
	}
	if users, ok := m.data["FindActiveResult"]; ok {
		return users.([]*user.User), nil
	}
	return nil, errors.New("not found")
}

func (m *MockUserRepository) Exists(id user.UserID) (bool, error) {
	m.calls["Exists"] = append(m.calls["Exists"], id)
	if err, ok := m.data["ExistsError"]; ok {
		return false, err.(error)
	}
	if exists, ok := m.data["ExistsResult"]; ok {
		return exists.(bool), nil
	}
	return false, nil
}

type MockWeightRepository struct {
	calls map[string][]interface{}
	data  map[string]interface{}
}

func NewMockWeightRepository() *MockWeightRepository {
	return &MockWeightRepository{
		calls: make(map[string][]interface{}),
		data:  make(map[string]interface{}),
	}
}

func (m *MockWeightRepository) Save(w *weight.Weight) error {
	m.calls["Save"] = append(m.calls["Save"], w)
	if err, ok := m.data["SaveError"]; ok {
		return err.(error)
	}
	return nil
}

func (m *MockWeightRepository) FindByID(id weight.WeightID) (*weight.Weight, error) {
	m.calls["FindByID"] = append(m.calls["FindByID"], id)
	if err, ok := m.data["FindByIDError"]; ok {
		return nil, err.(error)
	}
	if w, ok := m.data["FindByIDResult"]; ok {
		return w.(*weight.Weight), nil
	}
	return nil, errors.New("not found")
}

func (m *MockWeightRepository) FindByUserID(userID user.UserID, limit int) ([]*weight.Weight, error) {
	m.calls["FindByUserID"] = append(m.calls["FindByUserID"], userID, limit)
	if err, ok := m.data["FindByUserIDError"]; ok {
		return nil, err.(error)
	}
	if weights, ok := m.data["FindByUserIDResult"]; ok {
		return weights.([]*weight.Weight), nil
	}
	return nil, errors.New("not found")
}

func (m *MockWeightRepository) FindByUserIDAndPeriod(userID user.UserID, from, to time.Time) ([]*weight.Weight, error) {
	m.calls["FindByUserIDAndPeriod"] = append(m.calls["FindByUserIDAndPeriod"], userID, from, to)
	if err, ok := m.data["FindByUserIDAndPeriodError"]; ok {
		return nil, err.(error)
	}
	if weights, ok := m.data["FindByUserIDAndPeriodResult"]; ok {
		return weights.([]*weight.Weight), nil
	}
	return nil, errors.New("not found")
}

func (m *MockWeightRepository) FindLatestByUserID(userID user.UserID) (*weight.Weight, error) {
	m.calls["FindLatestByUserID"] = append(m.calls["FindLatestByUserID"], userID)
	if err, ok := m.data["FindLatestByUserIDError"]; ok {
		return nil, err.(error)
	}
	if w, ok := m.data["FindLatestByUserIDResult"]; ok {
		return w.(*weight.Weight), nil
	}
	return nil, errors.New("not found")
}

func (m *MockWeightRepository) CountByUserIDAndDate(userID user.UserID, date time.Time) (int, error) {
	m.calls["CountByUserIDAndDate"] = append(m.calls["CountByUserIDAndDate"], userID, date)
	if err, ok := m.data["CountByUserIDAndDateError"]; ok {
		return 0, err.(error)
	}
	if count, ok := m.data["CountByUserIDAndDateResult"]; ok {
		return count.(int), nil
	}
	return 0, nil
}

func (m *MockWeightRepository) Delete(id weight.WeightID) error {
	m.calls["Delete"] = append(m.calls["Delete"], id)
	if err, ok := m.data["DeleteError"]; ok {
		return err.(error)
	}
	return nil
}

func TestWeightTracker_RecordWeight(t *testing.T) {
	userID, _ := user.NewUserID("giada")
	value, _ := weight.NewWeightValue(70.5)
	unit, _ := weight.NewWeightUnit("kg")
	measuredAt := time.Now().AddDate(0, 0, -1) // Yesterday

	testUser, _ := user.NewUser("giada", "Giada", "")

	tests := []struct {
		name         string
		setupMocks   func(*MockUserRepository, *MockWeightRepository)
		userID       user.UserID
		value        weight.WeightValue
		unit         weight.WeightUnit
		measuredAt   time.Time
		expectErr    bool
		errorMessage string
	}{
		{
			name: "successful weight recording",
			setupMocks: func(ur *MockUserRepository, wr *MockWeightRepository) {
				ur.data["FindByIDResult"] = testUser
				wr.data["CountByUserIDAndDateResult"] = 2
			},
			userID:     userID,
			value:      value,
			unit:       unit,
			measuredAt: measuredAt,
			expectErr:  false,
		},
		{
			name: "user not found",
			setupMocks: func(ur *MockUserRepository, wr *MockWeightRepository) {
				ur.data["FindByIDError"] = errors.New("user not found")
			},
			userID:       userID,
			value:        value,
			unit:         unit,
			measuredAt:   measuredAt,
			expectErr:    true,
			errorMessage: "user not found",
		},
		{
			name: "user inactive",
			setupMocks: func(ur *MockUserRepository, wr *MockWeightRepository) {
				inactiveUser, _ := user.NewUser("giada", "Giada", "")
				inactiveUser.Deactivate()
				ur.data["FindByIDResult"] = inactiveUser
			},
			userID:       userID,
			value:        value,
			unit:         unit,
			measuredAt:   measuredAt,
			expectErr:    true,
			errorMessage: "user is not active",
		},
		{
			name: "too many weights per day",
			setupMocks: func(ur *MockUserRepository, wr *MockWeightRepository) {
				ur.data["FindByIDResult"] = testUser
				wr.data["CountByUserIDAndDateResult"] = 10
			},
			userID:       userID,
			value:        value,
			unit:         unit,
			measuredAt:   measuredAt,
			expectErr:    true,
			errorMessage: "maximum daily weight recordings exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := NewMockUserRepository()
			mockWeightRepo := NewMockWeightRepository()

			tt.setupMocks(mockUserRepo, mockWeightRepo)

			tracker := NewWeightTracker(mockUserRepo, mockWeightRepo)

			result, err := tracker.RecordWeight(tt.userID, tt.value, tt.unit, tt.measuredAt, "")

			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorMessage != "" && !contains(err.Error(), tt.errorMessage) {
					t.Errorf("expected error containing %q but got %q", tt.errorMessage, err.Error())
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
				if result != nil && result.UserID().String() != tt.userID.String() {
					t.Errorf("expected userID %s but got %s", tt.userID.String(), result.UserID().String())
				}
				if result != nil && result.Value().Float64() != tt.value.Float64() {
					t.Errorf("expected value %f but got %f", tt.value.Float64(), result.Value().Float64())
				}
			}
		})
	}
}

func TestWeightTracker_GetWeightHistory(t *testing.T) {
	mockWeightRepo := NewMockWeightRepository()
	mockUserRepo := NewMockUserRepository()

	userID, _ := user.NewUserID("giada")
	period := TimePeriodLastWeek

	// Create test weights
	weight1, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(70.0)), must(weight.NewWeightUnit("kg")), time.Now().AddDate(0, 0, -1), "")
	weight2, _ := weight.NewWeight("w2", userID, must(weight.NewWeightValue(69.5)), must(weight.NewWeightUnit("kg")), time.Now().AddDate(0, 0, -2), "")
	expectedWeights := []*weight.Weight{weight1, weight2}

	mockWeightRepo.data["FindByUserIDAndPeriodResult"] = expectedWeights

	tracker := NewWeightTracker(mockUserRepo, mockWeightRepo)

	weights, err := tracker.GetWeightHistory(userID, period)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(weights) != 2 {
		t.Errorf("expected 2 weights but got %d", len(weights))
	}
	if len(weights) > 0 && !equals(weights, expectedWeights) {
		t.Error("weights do not match expected values")
	}
}

func TestWeightTracker_CalculateWeightTrend(t *testing.T) {
	mockWeightRepo := NewMockWeightRepository()
	mockUserRepo := NewMockUserRepository()

	userID, _ := user.NewUserID("giada")
	period := TimePeriodLastMonth

	// Create test weights showing a downward trend
	now := time.Now()
	weight1, _ := weight.NewWeight("w1", userID, must(weight.NewWeightValue(72.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -30), "")
	weight2, _ := weight.NewWeight("w2", userID, must(weight.NewWeightValue(71.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -20), "")
	weight3, _ := weight.NewWeight("w3", userID, must(weight.NewWeightValue(70.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -10), "")
	weight4, _ := weight.NewWeight("w4", userID, must(weight.NewWeightValue(69.0)), must(weight.NewWeightUnit("kg")), now.AddDate(0, 0, -1), "")

	weights := []*weight.Weight{weight1, weight2, weight3, weight4}

	mockWeightRepo.data["FindByUserIDAndPeriodResult"] = weights

	tracker := NewWeightTracker(mockUserRepo, mockWeightRepo)

	trend, err := tracker.CalculateWeightTrend(userID, period)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if trend.Direction != TrendDecreasing {
		t.Errorf("expected TrendDecreasing but got %v", trend.Direction)
	}
	if trend.TotalChange.Float64() != 3.0 {
		t.Errorf("expected total change 3.0 but got %f", trend.TotalChange.Float64())
	}
	if trend.AverageChangePerWeek >= 0 {
		t.Errorf("expected negative average change but got %f", trend.AverageChangePerWeek)
	}
}

// Helper function to avoid repetitive error handling in tests
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func equals(a, b []*weight.Weight) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ID().String() != b[i].ID().String() {
			return false
		}
	}
	return true
}
