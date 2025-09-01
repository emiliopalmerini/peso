package weight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightID_NewWeightID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid weight id",
			value:   "weight_123",
			wantErr: false,
		},
		{
			name:    "empty weight id should fail",
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
			weightID, err := NewWeightID(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, weightID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, weightID.String())
			}
		})
	}
}

func TestWeightID_String(t *testing.T) {
	weightID, err := NewWeightID("weight_123")
	assert.NoError(t, err)
	assert.Equal(t, "weight_123", weightID.String())
}

func TestWeightID_IsEmpty(t *testing.T) {
	var weightID WeightID
	assert.True(t, weightID.IsEmpty())

	weightID, err := NewWeightID("weight_123")
	assert.NoError(t, err)
	assert.False(t, weightID.IsEmpty())
}