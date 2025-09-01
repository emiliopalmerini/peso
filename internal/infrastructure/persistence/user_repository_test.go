package persistence

import (
	"testing"

	"peso/internal/domain/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *DB {
	// Create a temporary database file
	dbFile := t.TempDir() + "/test.db"
	
	db, err := NewDB(dbFile)
	require.NoError(t, err)
	
	// Create tables manually for testing
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT DEFAULT '',
			active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)
	
	return db
}

func TestUserRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewUserRepository(db)
	
	user, err := user.NewUser("giada", "Giada", "giada@example.com")
	require.NoError(t, err)
	
	// Test save
	err = repo.Save(user)
	assert.NoError(t, err)
	
	// Test save again (update)
	user.UpdateEmail("giada.updated@example.com")
	err = repo.Save(user)
	assert.NoError(t, err)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewUserRepository(db)
	
	// Create and save user
	originalUser, err := user.NewUser("giada", "Giada", "giada@example.com")
	require.NoError(t, err)
	
	err = repo.Save(originalUser)
	require.NoError(t, err)
	
	// Find by ID
	userID, _ := user.NewUserID("giada")
	foundUser, err := repo.FindByID(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, originalUser.ID().String(), foundUser.ID().String())
	assert.Equal(t, originalUser.Name(), foundUser.Name())
	assert.Equal(t, originalUser.Email(), foundUser.Email())
	assert.Equal(t, originalUser.IsActive(), foundUser.IsActive())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewUserRepository(db)
	
	userID, _ := user.NewUserID("nonexistent")
	foundUser, err := repo.FindByID(userID)
	
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestUserRepository_FindByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewUserRepository(db)
	
	// Create and save user
	originalUser, err := user.NewUser("giada", "Giada", "")
	require.NoError(t, err)
	
	err = repo.Save(originalUser)
	require.NoError(t, err)
	
	// Find by name
	foundUser, err := repo.FindByName("Giada")
	
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, originalUser.ID().String(), foundUser.ID().String())
	assert.Equal(t, originalUser.Name(), foundUser.Name())
}

func TestUserRepository_FindActive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewUserRepository(db)
	
	// Create and save active user
	activeUser, err := user.NewUser("giada", "Giada", "")
	require.NoError(t, err)
	err = repo.Save(activeUser)
	require.NoError(t, err)
	
	// Create and save inactive user
	inactiveUser, err := user.NewUser("emilio", "Emilio", "")
	require.NoError(t, err)
	inactiveUser.Deactivate()
	err = repo.Save(inactiveUser)
	require.NoError(t, err)
	
	// Find active users
	activeUsers, err := repo.FindActive()
	
	assert.NoError(t, err)
	assert.Len(t, activeUsers, 1)
	assert.Equal(t, "giada", activeUsers[0].ID().String())
}

func TestUserRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewUserRepository(db)
	
	// Create and save user
	testUser, err := user.NewUser("giada", "Giada", "")
	require.NoError(t, err)
	err = repo.Save(testUser)
	require.NoError(t, err)
	
	// Test exists  
	exists, err := repo.Exists(testUser.ID())
	assert.NoError(t, err)
	assert.True(t, exists)
	
	// Test doesn't exist
	nonExistentID, _ := user.NewUserID("nonexistent")
	exists, err = repo.Exists(nonExistentID)
	assert.NoError(t, err)
	assert.False(t, exists)
}