package goal

import (
	"testing"
)

func TestGoalID_NewGoalID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid goal id",
			value:   "goal_123",
			wantErr: false,
		},
		{
			name:    "empty goal id should fail",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only should fail",
			value:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalID, err := NewGoalID(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if goalID.String() != "" {
					t.Errorf("expected empty goalID but got %s", goalID.String())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if goalID.String() != tt.value {
					t.Errorf("expected %s but got %s", tt.value, goalID.String())
				}
			}
		})
	}
}

func TestGoalID_String(t *testing.T) {
	goalID, err := NewGoalID("goal_123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if goalID.String() != "goal_123" {
		t.Errorf("expected goal_123 but got %s", goalID.String())
	}
}

func TestGoalID_IsEmpty(t *testing.T) {
	var goalID GoalID
	if !goalID.IsEmpty() {
		t.Error("expected empty goalID")
	}

	goalID, err := NewGoalID("goal_123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if goalID.IsEmpty() {
		t.Error("expected non-empty goalID")
	}
}
