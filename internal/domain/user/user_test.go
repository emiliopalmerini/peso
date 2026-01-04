package user

import (
	"testing"
	"time"
)

func TestUser_NewUser(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		userName string
		email    string
		wantErr  bool
	}{
		{
			name:     "valid user",
			id:       "giada",
			userName: "Giada",
			email:    "giada@example.com",
			wantErr:  false,
		},
		{
			name:     "valid user without email",
			id:       "emilio",
			userName: "Emilio",
			email:    "",
			wantErr:  false,
		},
		{
			name:     "invalid user id",
			id:       "",
			userName: "Test User",
			email:    "",
			wantErr:  true,
		},
		{
			name:     "empty name",
			id:       "test",
			userName: "",
			email:    "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.userName, tt.email)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if user != nil {
					t.Errorf("expected nil user but got %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Error("expected user but got nil")
				} else {
					if user.ID().String() != tt.id {
						t.Errorf("expected ID %s but got %s", tt.id, user.ID().String())
					}
					if user.Name() != tt.userName {
						t.Errorf("expected name %s but got %s", tt.userName, user.Name())
					}
					if user.Email() != tt.email {
						t.Errorf("expected email %s but got %s", tt.email, user.Email())
					}
					if !user.IsActive() {
						t.Error("expected user to be active")
					}
					if user.CreatedAt().IsZero() {
						t.Error("expected non-zero CreatedAt")
					}
					if user.UpdatedAt().IsZero() {
						t.Error("expected non-zero UpdatedAt")
					}
				}
			}
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	user, err := NewUser("giada", "Giada", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !user.IsActive() {
		t.Error("expected user to be active initially")
	}

	user.Deactivate()
	if user.IsActive() {
		t.Error("expected user to be inactive after deactivate")
	}
	if user.UpdatedAt().IsZero() {
		t.Error("expected non-zero UpdatedAt after deactivate")
	}
}

func TestUser_Activate(t *testing.T) {
	user, err := NewUser("giada", "Giada", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	user.Deactivate()
	if user.IsActive() {
		t.Error("expected user to be inactive after deactivate")
	}

	user.Activate()
	if !user.IsActive() {
		t.Error("expected user to be active after activate")
	}
	if user.UpdatedAt().IsZero() {
		t.Error("expected non-zero UpdatedAt after activate")
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	user, err := NewUser("giada", "Giada", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	originalUpdatedAt := user.UpdatedAt()
	time.Sleep(time.Millisecond) // Ensure time difference

	user.UpdateEmail("giada@example.com")
	if user.Email() != "giada@example.com" {
		t.Errorf("expected email giada@example.com but got %s", user.Email())
	}
	if !user.UpdatedAt().After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be after original")
	}
}

func TestUser_UpdateName(t *testing.T) {
	tests := []struct {
		name    string
		newName string
		wantErr bool
	}{
		{
			name:    "valid name update",
			newName: "Giada Updated",
			wantErr: false,
		},
		{
			name:    "empty name should fail",
			newName: "",
			wantErr: true,
		},
		{
			name:    "whitespace name should fail",
			newName: "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser("giada", "Giada", "")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			originalUpdatedAt := user.UpdatedAt()
			time.Sleep(time.Millisecond)

			err = user.UpdateName(tt.newName)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if user.Name() != "Giada" {
					t.Errorf("expected name to remain Giada but got %s", user.Name())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user.Name() != tt.newName {
					t.Errorf("expected name %s but got %s", tt.newName, user.Name())
				}
				if !user.UpdatedAt().After(originalUpdatedAt) {
					t.Error("expected UpdatedAt to be after original")
				}
			}
		})
	}
}
