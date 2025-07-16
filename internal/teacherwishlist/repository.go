package teacherwishlist

import (
	"context"
)

// TeacherRepository defines the interface for teacher data access operations
type TeacherRepository interface {
	// Create stores a new teacher in the repository
	Create(ctx context.Context, teacher *Teacher) error

	// GetByID retrieves a teacher by their unique identifier
	GetByID(ctx context.Context, id string) (*Teacher, error)

	// GetByEmail retrieves a teacher by their email address
	GetByEmail(ctx context.Context, email string) (*Teacher, error)

	// Update modifies an existing teacher's information
	Update(ctx context.Context, teacher *Teacher) error

	// Delete removes a teacher from the repository
	Delete(ctx context.Context, id string) error
}
