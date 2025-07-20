package sqlite

import (
	"context"
	"fmt"
	"time"

	"hrh-backend/internal/admin"
)

// AdminRepository implements the admin.AdminRepository interface for SQLite
type AdminRepository struct {
	db DBInterface
}

// NewAdminRepository creates a new SQLite admin repository
func NewAdminRepository(db DBInterface) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetByUsername retrieves an admin by their username
func (r *AdminRepository) GetByUsername(ctx context.Context, username string) (*admin.Admin, error) {
	query := `
		SELECT id, username, password_hash, created_at
		FROM admins
		WHERE username = ?`

	// Create temporary variables to scan into
	var id, usernameDB, passwordHash string
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&id,
		&usernameDB,
		&passwordHash,
		&createdAt,
	)

	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("admin with username %s not found: %w", username, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}

	// Create the admin with the retrieved data
	adminUser, err := admin.NewAdminWithHash(usernameDB, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin from database data: %w", err)
	}

	// Set the ID and CreatedAt from database
	adminUser.ID = id
	adminUser.CreatedAt = createdAt

	return adminUser, nil
}

// ValidateCredentials checks if the provided username and password are valid
// Returns the admin if credentials are valid, error otherwise
func (r *AdminRepository) ValidateCredentials(ctx context.Context, username, password string) (*admin.Admin, error) {
	// First get the admin by username
	adminUser, err := r.GetByUsername(ctx, username)
	if err != nil {
		// Return a generic error to avoid username enumeration
		return nil, fmt.Errorf("invalid credentials")
	}

	// Use bcrypt to validate the password against the stored hash
	if !adminUser.ValidatePassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return adminUser, nil
}

// Create persists a new admin (helper method for testing and initial setup)
func (r *AdminRepository) Create(ctx context.Context, adminUser *admin.Admin) error {
	query := `
		INSERT INTO admins (id, username, password_hash, created_at)
		VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		adminUser.ID,
		adminUser.Username,
		adminUser.PasswordHash,
		adminUser.CreatedAt,
	)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("admin with username %s already exists: %w", adminUser.Username, ErrDuplicateKey)
		}
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}
