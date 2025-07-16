package schooldirectory

import (
	"context"
)

// SchoolRepository defines the interface for school data access operations
type SchoolRepository interface {
	// Create stores a new school in the repository
	Create(ctx context.Context, school *School) error

	// GetByID retrieves a school by its unique identifier
	GetByID(ctx context.Context, id string) (*School, error)

	// GetAll retrieves all schools from the repository
	GetAll(ctx context.Context) ([]*School, error)

	// Update modifies an existing school's information
	Update(ctx context.Context, school *School) error

	// Delete removes a school from the repository
	Delete(ctx context.Context, id string) error
}
