package admin

import (
	"context"
)

// AdminFilters represents filtering criteria for admin queries
type AdminFilters struct {
	Username *string
	Limit    int
	Offset   int
}

// AdminRepository defines the interface for persisting and retrieving Admin entities.
// This interface follows aggregate-focused methods and includes proper context handling.
type AdminRepository interface {
	// Create persists a new admin entity
	Create(ctx context.Context, admin *Admin) error

	// GetByID retrieves an admin by their unique identifier
	GetByID(ctx context.Context, id string) (*Admin, error)

	// GetByUsername retrieves an admin by their username (unique constraint)
	GetByUsername(ctx context.Context, username string) (*Admin, error)

	// Update persists changes to an existing admin entity
	Update(ctx context.Context, admin *Admin) error

	// List retrieves admins based on filtering criteria with pagination support
	List(ctx context.Context, filters AdminFilters) ([]*Admin, error)

	// Delete removes an admin entity (for admin management)
	Delete(ctx context.Context, id string) error

	// ExistsByUsername checks if an admin with the given username already exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ValidateCredentials validates admin credentials and returns the admin if valid
	// This method handles the authentication logic at the repository level
	ValidateCredentials(ctx context.Context, username, passwordHash string) (*Admin, error)

	// Count returns the total number of admins in the system
	Count(ctx context.Context) (int, error)
}
