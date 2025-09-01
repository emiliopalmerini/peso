package weight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightUnit_NewWeightUnit(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    WeightUnit
		wantErr bool
	}{
		{
			name:    "valid kg unit",
			value:   "kg",
			want:    WeightUnitKg,
			wantErr: false,
		},
		{
			name:    "valid lb unit",
			value:   "lb",
			want:    WeightUnitLb,
			wantErr: false,
		},
		{
			name:    "invalid unit should fail",
			value:   "invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty unit should fail",
			value:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unit, err := NewWeightUnit(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, WeightUnit(""), unit)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, unit)
			}
		})
	}
}

func TestWeightUnit_String(t *testing.T) {
	assert.Equal(t, "kg", WeightUnitKg.String())
	assert.Equal(t, "lb", WeightUnitLb.String())
}

func TestWeightUnit_IsValid(t *testing.T) {
	assert.True(t, WeightUnitKg.IsValid())
	assert.True(t, WeightUnitLb.IsValid())
	assert.False(t, WeightUnit("invalid").IsValid())
}