package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.Error(t, err)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, userID.String())
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	userID, err := NewUserID("giada")
	assert.NoError(t, err)
	assert.Equal(t, "giada", userID.String())
}

func TestUserID_IsEmpty(t *testing.T) {
	var userID UserID
	assert.True(t, userID.IsEmpty())

	userID, err := NewUserID("giada")
	assert.NoError(t, err)
	assert.False(t, userID.IsEmpty())
}