package user

import (
	"testing"
)

func TestUserID_NewUserID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid user id",
			value:   "giada",
			wantErr: false,
		},
		{
			name:    "valid user id emilio",
			value:   "emilio",
			wantErr: false,
		},
		{
			name:    "empty user id should fail",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only should fail",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "too long should fail",
			value:   "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := NewUserID(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if userID.String() != "" {
					t.Errorf("expected empty userID but got %s", userID.String())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if userID.String() != tt.value {
					t.Errorf("expected %s but got %s", tt.value, userID.String())
				}
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	userID, err := NewUserID("giada")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if userID.String() != "giada" {
		t.Errorf("expected giada but got %s", userID.String())
	}
}

func TestUserID_IsEmpty(t *testing.T) {
	var userID UserID
	if !userID.IsEmpty() {
		t.Error("expected empty userID")
	}

	userID, err := NewUserID("giada")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if userID.IsEmpty() {
		t.Error("expected non-empty userID")
	}
}
