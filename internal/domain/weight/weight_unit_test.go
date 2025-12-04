package weight

import (
	"testing"
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
				if err == nil {
					t.Error("expected error but got nil")
				}
				if unit != WeightUnit("") {
					t.Errorf("expected empty unit but got %s", unit)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if unit != tt.want {
					t.Errorf("expected %s but got %s", tt.want, unit)
				}
			}
		})
	}
}

func TestWeightUnit_String(t *testing.T) {
	if WeightUnitKg.String() != "kg" {
		t.Errorf("expected kg but got %s", WeightUnitKg.String())
	}
	if WeightUnitLb.String() != "lb" {
		t.Errorf("expected lb but got %s", WeightUnitLb.String())
	}
}

func TestWeightUnit_IsValid(t *testing.T) {
	if !WeightUnitKg.IsValid() {
		t.Error("expected kg to be valid")
	}
	if !WeightUnitLb.IsValid() {
		t.Error("expected lb to be valid")
	}
	if WeightUnit("invalid").IsValid() {
		t.Error("expected invalid to be invalid")
	}
}
