package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_NewUser(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		userName string
		email   string
		wantErr bool
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
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.id, user.ID().String())
				assert.Equal(t, tt.userName, user.Name())
				assert.Equal(t, tt.email, user.Email())
				assert.True(t, user.IsActive())
				assert.False(t, user.CreatedAt().IsZero())
				assert.False(t, user.UpdatedAt().IsZero())
			}
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	user, err := NewUser("giada", "Giada", "")
	assert.NoError(t, err)
	
	assert.True(t, user.IsActive())
	
	user.Deactivate()
	assert.False(t, user.IsActive())
	assert.False(t, user.UpdatedAt().IsZero())
}

func TestUser_Activate(t *testing.T) {
	user, err := NewUser("giada", "Giada", "")
	assert.NoError(t, err)
	
	user.Deactivate()
	assert.False(t, user.IsActive())
	
	user.Activate()
	assert.True(t, user.IsActive())
	assert.False(t, user.UpdatedAt().IsZero())
}

func TestUser_UpdateEmail(t *testing.T) {
	user, err := NewUser("giada", "Giada", "")
	assert.NoError(t, err)
	
	originalUpdatedAt := user.UpdatedAt()
	time.Sleep(time.Millisecond) // Ensure time difference
	
	user.UpdateEmail("giada@example.com")
	assert.Equal(t, "giada@example.com", user.Email())
	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
}

func TestUser_UpdateName(t *testing.T) {
	tests := []struct {
		name     string
		newName  string
		wantErr  bool
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
			assert.NoError(t, err)
			
			originalUpdatedAt := user.UpdatedAt()
			time.Sleep(time.Millisecond)
			
			err = user.UpdateName(tt.newName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "Giada", user.Name()) // Should remain unchanged
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, user.Name())
				assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
			}
		})
	}
}