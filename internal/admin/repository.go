package admin

import (
	"context"
)

// AdminRepository defines the interface for admin data access operations
type AdminRepository interface {
	// GetByUsername retrieves an admin by their username
	GetByUsername(ctx context.Context, username string) (*Admin, error)

	// ValidateCredentials checks if the provided username and password are valid
	// Returns the admin if credentials are valid, error otherwise
	ValidateCredentials(ctx context.Context, username, password string) (*Admin, error)
}
