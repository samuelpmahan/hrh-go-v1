package admin

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Admin represents an administrator entity with system management privileges
type Admin struct {
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

// NewAdmin creates a new Admin with validation
func NewAdmin(username, passwordHash string) (*Admin, error) {
	id := generateID()
	now := time.Now()

	admin := &Admin{
		ID:           id,
		Username:     strings.TrimSpace(username),
		PasswordHash: passwordHash,
		CreatedAt:    now,
	}

	if err := admin.validateInvariants(); err != nil {
		return nil, err
	}

	return admin, nil
}

// validateInvariants checks all admin invariants
func (a *Admin) validateInvariants() error {
	// Username validation
	if a.Username == "" {
		return errors.New("username is required")
	}

	// Username length validation
	if len(a.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	// Password hash validation
	if a.PasswordHash == "" {
		return errors.New("password hash is required")
	}

	return nil
}

// UpdatePassword updates the admin's password hash
func (a *Admin) UpdatePassword(passwordHash string) error {
	if passwordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	a.PasswordHash = passwordHash
	return a.validateInvariants()
}
