package domain

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	password := "testpassword123"

	user, err := NewUser(username, password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}

	if user.PasswordHash == password {
		t.Error("Password hash should not be plain text")
	}

	if user.ID == "" {
		t.Error("User ID should not be empty")
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	// Verify password can be validated
	if !user.ValidatePassword(password) {
		t.Error("Password validation should succeed with correct password")
	}

	if user.ValidatePassword("wrongpassword") {
		t.Error("Password validation should fail with wrong password")
	}
}

func TestNewUser_EmptyUsername(t *testing.T) {
	_, err := NewUser("", "password123")
	if err == nil {
		t.Error("Expected error for empty username, got nil")
	}

	expectedError := "username is required"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewUser_ShortUsername(t *testing.T) {
	_, err := NewUser("ab", "password123")
	if err == nil {
		t.Error("Expected error for short username, got nil")
	}

	expectedError := "username must be at least 3 characters long"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewUser_EmptyPassword(t *testing.T) {
	_, err := NewUser("testuser", "")
	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}

	expectedError := "password cannot be empty"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewUser_WhitespaceUsername(t *testing.T) {
	user, err := NewUser("  testuser  ", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username to be trimmed to 'testuser', got '%s'", user.Username)
	}
}

func TestNewUserWithHash(t *testing.T) {
	username := "testuser"
	passwordHash := "$2a$10$abcdefghijklmnopqrstuvwxyz1234567890"

	user, err := NewUserWithHash(username, passwordHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}

	if user.PasswordHash != passwordHash {
		t.Errorf("Expected password hash %s, got %s", passwordHash, user.PasswordHash)
	}

	if user.ID == "" {
		t.Error("User ID should not be empty")
	}
}

func TestNewUserWithHash_EmptyHash(t *testing.T) {
	_, err := NewUserWithHash("testuser", "")
	if err == nil {
		t.Error("Expected error for empty password hash, got nil")
	}

	expectedError := "password hash is required"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	user, err := NewUser("testuser", "oldpassword")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	oldHash := user.PasswordHash
	newPassword := "newpassword123"

	err = user.UpdatePassword(newPassword)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.PasswordHash == oldHash {
		t.Error("Password hash should have changed")
	}

	if user.PasswordHash == newPassword {
		t.Error("Password hash should not be plain text")
	}

	// Old password should no longer work
	if user.ValidatePassword("oldpassword") {
		t.Error("Old password should no longer be valid")
	}

	// New password should work
	if !user.ValidatePassword(newPassword) {
		t.Error("New password should be valid")
	}
}

func TestUser_UpdatePassword_Empty(t *testing.T) {
	user, err := NewUser("testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = user.UpdatePassword("")
	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}

	expectedError := "password cannot be empty"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	password := "testpassword123"
	user, err := NewUser("testuser", password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Correct password should validate
	if !user.ValidatePassword(password) {
		t.Error("Correct password should validate")
	}

	// Wrong password should not validate
	if user.ValidatePassword("wrongpassword") {
		t.Error("Wrong password should not validate")
	}

	// Empty password should not validate
	if user.ValidatePassword("") {
		t.Error("Empty password should not validate")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if hash == password {
		t.Error("Hash should not be plain text")
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	// Verify it's a valid bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Errorf("Hash should be valid bcrypt hash: %v", err)
	}
}

func TestHashPassword_Empty(t *testing.T) {
	_, err := HashPassword("")
	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}

	expectedError := "password cannot be empty"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidatePasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Correct password should validate
	if !ValidatePasswordHash(password, hash) {
		t.Error("Correct password should validate against hash")
	}

	// Wrong password should not validate
	if ValidatePasswordHash("wrongpassword", hash) {
		t.Error("Wrong password should not validate against hash")
	}

	// Empty password should not validate
	if ValidatePasswordHash("", hash) {
		t.Error("Empty password should not validate against hash")
	}

	// Invalid hash should not validate
	if ValidatePasswordHash(password, "invalidhash") {
		t.Error("Invalid hash should not validate")
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == "" {
		t.Error("Generated ID should not be empty")
	}

	if id2 == "" {
		t.Error("Generated ID should not be empty")
	}

	if id1 == id2 {
		t.Error("Generated IDs should be unique")
	}

	// Check format (should be like: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
	parts := strings.Split(id1, "-")
	if len(parts) != 5 {
		t.Errorf("Expected ID to have 5 parts separated by hyphens, got %d parts", len(parts))
	}
}
