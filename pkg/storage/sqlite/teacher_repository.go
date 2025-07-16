package sqlite

import (
	"context"
	"fmt"

	"hrh-backend/internal/teacherwishlist"
)

// TeacherRepository implements the teacherwishlist.TeacherRepository interface for SQLite
type TeacherRepository struct {
	db DBInterface
}

// NewTeacherRepository creates a new SQLite teacher repository
func NewTeacherRepository(db DBInterface) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// Create persists a new teacher aggregate
func (r *TeacherRepository) Create(ctx context.Context, teacher *teacherwishlist.Teacher) error {
	query := `
		INSERT INTO teachers (id, email, first_name, last_name, school_id, grade_level, wishlist_url, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		teacher.ID,
		teacher.Email,
		teacher.FirstName,
		teacher.LastName,
		teacher.SchoolID,
		teacher.GradeLevel,
		teacher.WishlistURL,
		teacher.Status.String(),
		teacher.CreatedAt,
		teacher.UpdatedAt,
	)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("teacher with email %s already exists: %w", teacher.Email, ErrDuplicateKey)
		}
		if IsForeignKeyError(err) {
			return fmt.Errorf("school with ID %s does not exist: %w", teacher.SchoolID, ErrForeignKey)
		}
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	return nil
}

// GetByID retrieves a teacher by their unique identifier
func (r *TeacherRepository) GetByID(ctx context.Context, id string) (*teacherwishlist.Teacher, error) {
	query := `
		SELECT id, email, first_name, last_name, school_id, grade_level, wishlist_url, status, created_at, updated_at
		FROM teachers
		WHERE id = ?`

	teacher := &teacherwishlist.Teacher{}
	var statusStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&teacher.ID,
		&teacher.Email,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.SchoolID,
		&teacher.GradeLevel,
		&teacher.WishlistURL,
		&statusStr,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("teacher with ID %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get teacher by ID: %w", err)
	}

	teacher.Status = teacherwishlist.TeacherStatus(statusStr)
	return teacher, nil
}

// GetByEmail retrieves a teacher by their email address
func (r *TeacherRepository) GetByEmail(ctx context.Context, email string) (*teacherwishlist.Teacher, error) {
	query := `
		SELECT id, email, first_name, last_name, school_id, grade_level, wishlist_url, status, created_at, updated_at
		FROM teachers
		WHERE email = ?`

	teacher := &teacherwishlist.Teacher{}
	var statusStr string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&teacher.ID,
		&teacher.Email,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.SchoolID,
		&teacher.GradeLevel,
		&teacher.WishlistURL,
		&statusStr,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("teacher with email %s not found: %w", email, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get teacher by email: %w", err)
	}

	teacher.Status = teacherwishlist.TeacherStatus(statusStr)
	return teacher, nil
}

// Update persists changes to an existing teacher aggregate
func (r *TeacherRepository) Update(ctx context.Context, teacher *teacherwishlist.Teacher) error {
	query := `
		UPDATE teachers 
		SET email = ?, first_name = ?, last_name = ?, school_id = ?, grade_level = ?, 
		    wishlist_url = ?, status = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		teacher.Email,
		teacher.FirstName,
		teacher.LastName,
		teacher.SchoolID,
		teacher.GradeLevel,
		teacher.WishlistURL,
		teacher.Status.String(),
		teacher.UpdatedAt,
		teacher.ID,
	)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("teacher with email %s already exists: %w", teacher.Email, ErrDuplicateKey)
		}
		if IsForeignKeyError(err) {
			return fmt.Errorf("school with ID %s does not exist: %w", teacher.SchoolID, ErrForeignKey)
		}
		return fmt.Errorf("failed to update teacher: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("teacher with ID %s not found: %w", teacher.ID, ErrNotFound)
	}

	return nil
}

// Delete removes a teacher aggregate
func (r *TeacherRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM teachers WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("teacher with ID %s not found: %w", id, ErrNotFound)
	}

	return nil
}
