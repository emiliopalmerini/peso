package weight

import (
	"testing"
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
				if err == nil {
					t.Error("expected error but got nil")
				}
				if weight != WeightValue(0) {
					t.Errorf("expected zero weight but got %f", weight.Float64())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if weight.Float64() != tt.value {
					t.Errorf("expected %f but got %f", tt.value, weight.Float64())
				}
			}
		})
	}
}

func TestWeightValue_Float64(t *testing.T) {
	weight, err := NewWeightValue(75.3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if weight.Float64() != 75.3 {
		t.Errorf("expected 75.3 but got %f", weight.Float64())
	}
}

func TestWeightValue_String(t *testing.T) {
	weight, err := NewWeightValue(75.3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if weight.String() != "75.3" {
		t.Errorf("expected 75.3 but got %s", weight.String())
	}
}

func TestWeightValue_IsZero(t *testing.T) {
	var weight WeightValue
	if !weight.IsZero() {
		t.Error("expected zero weight")
	}

	weight, err := NewWeightValue(75.3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if weight.IsZero() {
		t.Error("expected non-zero weight")
	}
}

func TestWeightValue_Subtract(t *testing.T) {
	weight1, _ := NewWeightValue(80.0)
	weight2, _ := NewWeightValue(75.0)

	diff := weight1.Subtract(weight2)
	if diff.Float64() != 5.0 {
		t.Errorf("expected 5.0 but got %f", diff.Float64())
	}
}

func TestWeightValue_Add(t *testing.T) {
	weight1, _ := NewWeightValue(70.0)
	weight2, _ := NewWeightValue(15.0)

	sum, err := weight1.Add(weight2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if sum.Float64() != 85.0 {
		t.Errorf("expected 85.0 but got %f", sum.Float64())
	}

	// Test that adding results in invalid weight
	weight3, _ := NewWeightValue(490.0)
	weight4, _ := NewWeightValue(20.0)

	_, err = weight3.Add(weight4)
	if err == nil {
		t.Error("expected error but got nil")
	}
}
