package teacherwishlist

import (
	"context"
)

// TeacherFilters represents filtering criteria for teacher queries
type TeacherFilters struct {
	Status     *TeacherStatus
	SchoolID   *string
	GradeLevel *string
	Email      *string
	Limit      int
	Offset     int
}

// TeacherRepository defines the interface for persisting and retrieving Teacher aggregates.
// This interface follows aggregate-focused methods and includes proper context handling.
type TeacherRepository interface {
	// Create persists a new teacher aggregate
	Create(ctx context.Context, teacher *Teacher) error

	// GetByID retrieves a teacher by their unique identifier
	GetByID(ctx context.Context, id string) (*Teacher, error)

	// GetByEmail retrieves a teacher by their email address (unique constraint)
	GetByEmail(ctx context.Context, email string) (*Teacher, error)

	// Update persists changes to an existing teacher aggregate
	Update(ctx context.Context, teacher *Teacher) error

	// List retrieves teachers based on filtering criteria with pagination support
	List(ctx context.Context, filters TeacherFilters) ([]*Teacher, error)

	// GetActiveTeachers retrieves all approved teachers with valid wishlist URLs for donor search
	GetActiveTeachers(ctx context.Context) ([]*Teacher, error)

	// GetPendingTeachers retrieves all teachers with pending status for admin review
	GetPendingTeachers(ctx context.Context) ([]*Teacher, error)

	// Delete removes a teacher aggregate (for admin management)
	Delete(ctx context.Context, id string) error

	// ExistsByEmail checks if a teacher with the given email already exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// CountBySchool returns the number of teachers associated with a specific school
	CountBySchool(ctx context.Context, schoolID string) (int, error)
}
