package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a base user entity with common authentication properties
type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// generateID creates a simple UUID-like string using crypto/rand
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// NewUser creates a new User with validation and password hashing
func NewUser(username, plainPassword string) (*User, error) {
	id := generateID()
	now := time.Now()

	// Hash the password using bcrypt
	passwordHash, err := HashPassword(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		ID:           id,
		Username:     strings.TrimSpace(username),
		PasswordHash: passwordHash,
		CreatedAt:    now,
	}

	if err := user.validateInvariants(); err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserWithHash creates a new User with an already hashed password (for testing/migration)
func NewUserWithHash(username, passwordHash string) (*User, error) {
	id := generateID()
	now := time.Now()

	user := &User{
		ID:           id,
		Username:     strings.TrimSpace(username),
		PasswordHash: passwordHash,
		CreatedAt:    now,
	}

	if err := user.validateInvariants(); err != nil {
		return nil, err
	}

	return user, nil
}

// validateInvariants checks all user invariants
func (u *User) validateInvariants() error {
	// Username validation
	if u.Username == "" {
		return errors.New("username is required")
	}

	// Username length validation
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	// Password hash validation
	if u.PasswordHash == "" {
		return errors.New("password hash is required")
	}

	return nil
}

// UpdatePassword updates the user's password with proper hashing
func (u *User) UpdatePassword(plainPassword string) error {
	if plainPassword == "" {
		return errors.New("password cannot be empty")
	}

	passwordHash, err := HashPassword(plainPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = passwordHash
	return u.validateInvariants()
}

// ValidatePassword checks if the provided plain password matches the stored hash
func (u *User) ValidatePassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	return err == nil
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(plainPassword string) (string, error) {
	if plainPassword == "" {
		return "", errors.New("password cannot be empty")
	}

	// Use bcrypt with default cost (currently 10)
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}

	return string(hash), nil
}

// ValidatePasswordHash checks if a plain password matches a bcrypt hash
func ValidatePasswordHash(plainPassword, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword))
	return err == nil
}
