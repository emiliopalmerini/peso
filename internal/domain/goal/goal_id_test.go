package goal

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.Error(t, err)
				assert.Empty(t, goalID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, goalID.String())
			}
		})
	}
}

func TestGoalID_String(t *testing.T) {
	goalID, err := NewGoalID("goal_123")
	assert.NoError(t, err)
	assert.Equal(t, "goal_123", goalID.String())
}

func TestGoalID_IsEmpty(t *testing.T) {
	var goalID GoalID
	assert.True(t, goalID.IsEmpty())

	goalID, err := NewGoalID("goal_123")
	assert.NoError(t, err)
	assert.False(t, goalID.IsEmpty())
}