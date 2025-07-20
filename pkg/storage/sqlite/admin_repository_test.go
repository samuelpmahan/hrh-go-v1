package sqlite

import (
	"context"
	"testing"

	"hrh-backend/internal/admin"
)

// createTestAdmin creates a test admin for use in tests
func createTestAdmin(t *testing.T, username, plainPassword string) *admin.Admin {
	adminUser, err := admin.NewAdmin(username, plainPassword)
	if err != nil {
		t.Fatalf("Failed to create test admin: %v", err)
	}
	return adminUser
}

// createTestAdminWithHash creates a test admin with an already hashed password
func createTestAdminWithHash(t *testing.T, username, passwordHash string) *admin.Admin {
	adminUser, err := admin.NewAdminWithHash(username, passwordHash)
	if err != nil {
		t.Fatalf("Failed to create test admin with hash: %v", err)
	}
	return adminUser
}

func TestAdminRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	adminUser := createTestAdmin(t, "testadmin", "hashedpassword123")

	// Test successful creation
	err := repo.Create(ctx, adminUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify admin was created by retrieving it
	retrieved, err := repo.GetByUsername(ctx, adminUser.Username)
	if err != nil {
		t.Errorf("Failed to retrieve created admin: %v", err)
	}

	if retrieved.Username != adminUser.Username {
		t.Errorf("Expected username %s, got %s", adminUser.Username, retrieved.Username)
	}
	if retrieved.PasswordHash != adminUser.PasswordHash {
		t.Errorf("Expected password hash %s, got %s", adminUser.PasswordHash, retrieved.PasswordHash)
	}
	if retrieved.ID != adminUser.ID {
		t.Errorf("Expected ID %s, got %s", adminUser.ID, retrieved.ID)
	}
}

func TestAdminRepository_Create_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	admin1 := createTestAdmin(t, "duplicateuser", "password1")

	// Create first admin
	err := repo.Create(ctx, admin1)
	if err != nil {
		t.Fatalf("Failed to create first admin: %v", err)
	}

	// Try to create second admin with same username
	admin2 := createTestAdmin(t, "duplicateuser", "password2")

	err = repo.Create(ctx, admin2)
	if err == nil {
		t.Error("Expected duplicate key error, got nil")
	}
	if !IsDuplicateKeyError(err) {
		t.Errorf("Expected duplicate key error, got %v", err)
	}
}

func TestAdminRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	adminUser := createTestAdmin(t, "getbyusername", "testpassword")

	// Create admin first
	err := repo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Test successful retrieval
	retrieved, err := repo.GetByUsername(ctx, adminUser.Username)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.Username != adminUser.Username {
		t.Errorf("Expected username %s, got %s", adminUser.Username, retrieved.Username)
	}
	if retrieved.ID != adminUser.ID {
		t.Errorf("Expected ID %s, got %s", adminUser.ID, retrieved.ID)
	}
	if retrieved.PasswordHash != adminUser.PasswordHash {
		t.Errorf("Expected password hash %s, got %s", adminUser.PasswordHash, retrieved.PasswordHash)
	}
}

func TestAdminRepository_GetByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	_, err := repo.GetByUsername(ctx, "nonexistentuser")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestAdminRepository_ValidateCredentials_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	username := "validuser"
	password := "correctpassword"
	adminUser := createTestAdmin(t, username, password)

	// Create admin first
	err := repo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Test successful credential validation
	validated, err := repo.ValidateCredentials(ctx, username, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if validated.Username != username {
		t.Errorf("Expected username %s, got %s", username, validated.Username)
	}
	if validated.ID != adminUser.ID {
		t.Errorf("Expected ID %s, got %s", adminUser.ID, validated.ID)
	}
}

func TestAdminRepository_ValidateCredentials_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	username := "testuser"
	correctPassword := "correctpassword"
	wrongPassword := "wrongpassword"
	adminUser := createTestAdmin(t, username, correctPassword)

	// Create admin first
	err := repo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Test credential validation with wrong password
	_, err = repo.ValidateCredentials(ctx, username, wrongPassword)
	if err == nil {
		t.Error("Expected invalid credentials error, got nil")
	}

	expectedError := "invalid credentials"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAdminRepository_ValidateCredentials_NonexistentUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	// Test credential validation with nonexistent user
	_, err := repo.ValidateCredentials(ctx, "nonexistentuser", "anypassword")
	if err == nil {
		t.Error("Expected invalid credentials error, got nil")
	}

	expectedError := "invalid credentials"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAdminRepository_ValidateCredentials_EmptyCredentials(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	// Test with empty username
	_, err := repo.ValidateCredentials(ctx, "", "password")
	if err == nil {
		t.Error("Expected invalid credentials error for empty username, got nil")
	}

	// Test with empty password
	_, err = repo.ValidateCredentials(ctx, "username", "")
	if err == nil {
		t.Error("Expected invalid credentials error for empty password, got nil")
	}

	// Test with both empty
	_, err = repo.ValidateCredentials(ctx, "", "")
	if err == nil {
		t.Error("Expected invalid credentials error for empty credentials, got nil")
	}
}

func TestAdminRepository_CRUD_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	// Create an admin
	username := "integrationtest"
	password := "testpassword123"
	adminUser := createTestAdmin(t, username, password)

	// Test Create
	err := repo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Test Read by Username
	retrieved, err := repo.GetByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get admin by username: %v", err)
	}
	if retrieved.Username != username {
		t.Errorf("Expected username %s, got %s", username, retrieved.Username)
	}
	// Password hash should be bcrypt hash, not plain text
	if retrieved.PasswordHash == password {
		t.Error("Password hash should not be plain text")
	}
	// Verify the password can be validated
	if !retrieved.ValidatePassword(password) {
		t.Error("Password validation should succeed with correct password")
	}

	// Test ValidateCredentials with correct credentials
	validated, err := repo.ValidateCredentials(ctx, username, password)
	if err != nil {
		t.Fatalf("Failed to validate correct credentials: %v", err)
	}
	if validated.Username != username {
		t.Errorf("Expected validated username %s, got %s", username, validated.Username)
	}

	// Test ValidateCredentials with incorrect password
	_, err = repo.ValidateCredentials(ctx, username, "wrongpassword")
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}

	// Test ValidateCredentials with nonexistent user
	_, err = repo.ValidateCredentials(ctx, "nonexistent", password)
	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
}

func TestAdminRepository_MultipleAdmins(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAdminRepository(db)
	ctx := context.Background()

	// Create multiple admins with their plain text passwords
	adminData := []struct {
		username string
		password string
		admin    *admin.Admin
	}{
		{"admin1", "password1", nil},
		{"admin2", "password2", nil},
		{"admin3", "password3", nil},
	}

	// Create admin objects
	for i := range adminData {
		adminData[i].admin = createTestAdmin(t, adminData[i].username, adminData[i].password)
	}

	// Create all admins in database
	for _, data := range adminData {
		err := repo.Create(ctx, data.admin)
		if err != nil {
			t.Fatalf("Failed to create admin %s: %v", data.username, err)
		}
	}

	// Verify each admin can be retrieved and validated
	for _, data := range adminData {
		// Test GetByUsername
		retrieved, err := repo.GetByUsername(ctx, data.username)
		if err != nil {
			t.Errorf("Failed to get admin %s: %v", data.username, err)
			continue
		}
		if retrieved.Username != data.username {
			t.Errorf("Expected username %s, got %s", data.username, retrieved.Username)
		}

		// Test ValidateCredentials with correct password
		validated, err := repo.ValidateCredentials(ctx, data.username, data.password)
		if err != nil {
			t.Errorf("Failed to validate credentials for admin %s: %v", data.username, err)
			continue
		}
		if validated.Username != data.username {
			t.Errorf("Expected validated username %s, got %s", data.username, validated.Username)
		}
	}

	// Test cross-validation (admin1's password shouldn't work for admin2)
	_, err := repo.ValidateCredentials(ctx, adminData[0].username, adminData[1].password)
	if err == nil {
		t.Error("Expected error when using admin2's password for admin1, got nil")
	}
}
