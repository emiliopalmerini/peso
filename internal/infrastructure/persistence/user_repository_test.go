package persistence

import (
	"testing"

	"peso/internal/domain/user"
)

func setupTestDB(t *testing.T) *DB {
	// Create a temporary database file
	dbFile := t.TempDir() + "/test.db"

	db, err := NewDB(dbFile)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Create tables manually for testing
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT DEFAULT '',
			password_hash TEXT DEFAULT '',
			active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	return db
}

func TestUserRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	testUser, err := user.NewUser("giada", "Giada", "giada@example.com")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test save
	err = repo.Save(testUser)
	if err != nil {
		t.Errorf("unexpected error saving user: %v", err)
	}

	// Test save again (update)
	testUser.UpdateEmail("giada.updated@example.com")
	err = repo.Save(testUser)
	if err != nil {
		t.Errorf("unexpected error updating user: %v", err)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create and save user
	originalUser, err := user.NewUser("giada", "Giada", "giada@example.com")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	err = repo.Save(originalUser)
	if err != nil {
		t.Fatalf("failed to save test user: %v", err)
	}

	// Find by ID
	userID, _ := user.NewUserID("giada")
	foundUser, err := repo.FindByID(userID)

	if err != nil {
		t.Errorf("unexpected error finding user: %v", err)
	}
	if foundUser == nil {
		t.Error("expected user but got nil")
	} else {
		if foundUser.ID().String() != originalUser.ID().String() {
			t.Errorf("expected ID %s but got %s", originalUser.ID().String(), foundUser.ID().String())
		}
		if foundUser.Name() != originalUser.Name() {
			t.Errorf("expected name %s but got %s", originalUser.Name(), foundUser.Name())
		}
		if foundUser.Email() != originalUser.Email() {
			t.Errorf("expected email %s but got %s", originalUser.Email(), foundUser.Email())
		}
		if foundUser.IsActive() != originalUser.IsActive() {
			t.Errorf("expected active %v but got %v", originalUser.IsActive(), foundUser.IsActive())
		}
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	userID, _ := user.NewUserID("nonexistent")
	foundUser, err := repo.FindByID(userID)

	if err == nil {
		t.Error("expected error but got nil")
	}
	if foundUser != nil {
		t.Errorf("expected nil user but got %v", foundUser)
	}
}

func TestUserRepository_FindByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create and save user
	originalUser, err := user.NewUser("giada", "Giada", "")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	err = repo.Save(originalUser)
	if err != nil {
		t.Fatalf("failed to save test user: %v", err)
	}

	// Find by name
	foundUser, err := repo.FindByName("Giada")

	if err != nil {
		t.Errorf("unexpected error finding user: %v", err)
	}
	if foundUser == nil {
		t.Error("expected user but got nil")
	} else {
		if foundUser.ID().String() != originalUser.ID().String() {
			t.Errorf("expected ID %s but got %s", originalUser.ID().String(), foundUser.ID().String())
		}
		if foundUser.Name() != originalUser.Name() {
			t.Errorf("expected name %s but got %s", originalUser.Name(), foundUser.Name())
		}
	}
}

func TestUserRepository_FindActive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create and save active user
	activeUser, err := user.NewUser("giada", "Giada", "")
	if err != nil {
		t.Fatalf("failed to create active user: %v", err)
	}
	err = repo.Save(activeUser)
	if err != nil {
		t.Fatalf("failed to save active user: %v", err)
	}

	// Create and save inactive user
	inactiveUser, err := user.NewUser("emilio", "Emilio", "")
	if err != nil {
		t.Fatalf("failed to create inactive user: %v", err)
	}
	inactiveUser.Deactivate()
	err = repo.Save(inactiveUser)
	if err != nil {
		t.Fatalf("failed to save inactive user: %v", err)
	}

	// Find active users
	activeUsers, err := repo.FindActive()

	if err != nil {
		t.Errorf("unexpected error finding active users: %v", err)
	}
	if len(activeUsers) != 1 {
		t.Errorf("expected 1 active user but got %d", len(activeUsers))
	} else if activeUsers[0].ID().String() != "giada" {
		t.Errorf("expected user giada but got %s", activeUsers[0].ID().String())
	}
}

func TestUserRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create and save user
	testUser, err := user.NewUser("giada", "Giada", "")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	err = repo.Save(testUser)
	if err != nil {
		t.Fatalf("failed to save test user: %v", err)
	}

	// Test exists
	exists, err := repo.Exists(testUser.ID())
	if err != nil {
		t.Errorf("unexpected error checking existence: %v", err)
	}
	if !exists {
		t.Error("expected user to exist")
	}

	// Test doesn't exist
	nonExistentID, _ := user.NewUserID("nonexistent")
	exists, err = repo.Exists(nonExistentID)
	if err != nil {
		t.Errorf("unexpected error checking existence: %v", err)
	}
	if exists {
		t.Error("expected user to not exist")
	}
}
