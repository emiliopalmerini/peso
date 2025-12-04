package weight

import (
	"testing"
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
				if err == nil {
					t.Error("expected error but got nil")
				}
				if weightID.String() != "" {
					t.Errorf("expected empty weightID but got %s", weightID.String())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if weightID.String() != tt.value {
					t.Errorf("expected %s but got %s", tt.value, weightID.String())
				}
			}
		})
	}
}

func TestWeightID_String(t *testing.T) {
	weightID, err := NewWeightID("weight_123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if weightID.String() != "weight_123" {
		t.Errorf("expected weight_123 but got %s", weightID.String())
	}
}

func TestWeightID_IsEmpty(t *testing.T) {
	var weightID WeightID
	if !weightID.IsEmpty() {
		t.Error("expected empty weightID")
	}

	weightID, err := NewWeightID("weight_123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if weightID.IsEmpty() {
		t.Error("expected non-empty weightID")
	}
}
