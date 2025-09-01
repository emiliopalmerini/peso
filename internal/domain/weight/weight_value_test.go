package weight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightValue_NewWeightValue(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{
			name:    "valid weight",
			value:   70.5,
			wantErr: false,
		},
		{
			name:    "minimum valid weight",
			value:   10.0,
			wantErr: false,
		},
		{
			name:    "maximum valid weight",
			value:   500.0,
			wantErr: false,
		},
		{
			name:    "zero weight should fail",
			value:   0.0,
			wantErr: true,
		},
		{
			name:    "negative weight should fail",
			value:   -10.0,
			wantErr: true,
		},
		{
			name:    "too low weight should fail",
			value:   9.9,
			wantErr: true,
		},
		{
			name:    "too high weight should fail",
			value:   500.1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight, err := NewWeightValue(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, WeightValue(0), weight)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, weight.Float64())
			}
		})
	}
}

func TestWeightValue_Float64(t *testing.T) {
	weight, err := NewWeightValue(75.3)
	assert.NoError(t, err)
	assert.Equal(t, 75.3, weight.Float64())
}

func TestWeightValue_String(t *testing.T) {
	weight, err := NewWeightValue(75.3)
	assert.NoError(t, err)
	assert.Equal(t, "75.3", weight.String())
}

func TestWeightValue_IsZero(t *testing.T) {
	var weight WeightValue
	assert.True(t, weight.IsZero())

	weight, err := NewWeightValue(75.3)
	assert.NoError(t, err)
	assert.False(t, weight.IsZero())
}

func TestWeightValue_Subtract(t *testing.T) {
	weight1, _ := NewWeightValue(80.0)
	weight2, _ := NewWeightValue(75.0)
	
	diff := weight1.Subtract(weight2)
	assert.Equal(t, 5.0, diff.Float64())
}

func TestWeightValue_Add(t *testing.T) {
	weight1, _ := NewWeightValue(70.0)
	weight2, _ := NewWeightValue(15.0)
	
	sum, err := weight1.Add(weight2)
	assert.NoError(t, err)
	assert.Equal(t, 85.0, sum.Float64())
	
	// Test that adding results in invalid weight
	weight3, _ := NewWeightValue(490.0)
	weight4, _ := NewWeightValue(20.0)
	
	_, err = weight3.Add(weight4)
	assert.Error(t, err)
}