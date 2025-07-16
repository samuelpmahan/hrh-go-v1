package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"hrh-backend/internal/teacherwishlist"
)

// setupTestDB creates a temporary SQLite database for testing
func setupTestDB(t *testing.T) *DB {
	// Create temporary database file
	tmpFile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database file: %v", err)
	}
	tmpFile.Close()

	// Clean up the file when test completes
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	config := Config{
		DatabasePath: tmpFile.Name(),
	}

	db, err := NewDB(config)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Clean up database connection when test completes
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// createTestTeacher creates a test teacher for use in tests
func createTestTeacher(t *testing.T) *teacherwishlist.Teacher {
	teacher, err := teacherwishlist.NewTeacher(
		"test@example.com",
		"John",
		"Doe",
		"school-123",
		"5th Grade",
	)
	if err != nil {
		t.Fatalf("Failed to create test teacher: %v", err)
	}
	return teacher
}

// insertTestSchool inserts a test school to satisfy foreign key constraints
func insertTestSchool(t *testing.T, db *DB, schoolID string) {
	query := `
		INSERT INTO schools (id, name, address_city, address_state, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.ExecContext(context.Background(), query,
		schoolID, "Test School", "Test City", "Test State", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert test school: %v", err)
	}
}

func TestTeacherRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher := createTestTeacher(t)
	insertTestSchool(t, db, teacher.SchoolID)

	// Test successful creation
	err := repo.Create(ctx, teacher)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify teacher was created by retrieving it
	retrieved, err := repo.GetByID(ctx, teacher.ID)
	if err != nil {
		t.Errorf("Failed to retrieve created teacher: %v", err)
	}

	if retrieved.Email != teacher.Email {
		t.Errorf("Expected email %s, got %s", teacher.Email, retrieved.Email)
	}
	if retrieved.FirstName != teacher.FirstName {
		t.Errorf("Expected first name %s, got %s", teacher.FirstName, retrieved.FirstName)
	}
	if retrieved.Status != teacher.Status {
		t.Errorf("Expected status %s, got %s", teacher.Status, retrieved.Status)
	}
}

func TestTeacherRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher1 := createTestTeacher(t)
	insertTestSchool(t, db, teacher1.SchoolID)

	// Create first teacher
	err := repo.Create(ctx, teacher1)
	if err != nil {
		t.Fatalf("Failed to create first teacher: %v", err)
	}

	// Try to create second teacher with same email
	teacher2, _ := teacherwishlist.NewTeacher(
		teacher1.Email, // Same email
		"Jane",
		"Smith",
		teacher1.SchoolID,
		"3rd Grade",
	)

	err = repo.Create(ctx, teacher2)
	if err == nil {
		t.Error("Expected duplicate key error, got nil")
	}
	if !IsDuplicateKeyError(err) {
		t.Errorf("Expected duplicate key error, got %v", err)
	}
}

func TestTeacherRepository_Create_InvalidSchoolID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	// Create a teacher with a unique email and non-existent school ID
	teacher, err := teacherwishlist.NewTeacher(
		"unique-fk-test@example.com",
		"John",
		"Doe",
		"non-existent-school-id",
		"5th Grade",
	)
	if err != nil {
		t.Fatalf("Failed to create test teacher: %v", err)
	}

	// Don't insert the school, so foreign key constraint will fail
	err = repo.Create(ctx, teacher)
	if err == nil {
		t.Error("Expected foreign key error, got nil")
	}
	if !IsForeignKeyError(err) {
		t.Errorf("Expected foreign key error, got %v", err)
	}
}

func TestTeacherRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher := createTestTeacher(t)
	insertTestSchool(t, db, teacher.SchoolID)

	// Create teacher first
	err := repo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	// Test successful retrieval
	retrieved, err := repo.GetByID(ctx, teacher.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.ID != teacher.ID {
		t.Errorf("Expected ID %s, got %s", teacher.ID, retrieved.ID)
	}
	if retrieved.Email != teacher.Email {
		t.Errorf("Expected email %s, got %s", teacher.Email, retrieved.Email)
	}
}

func TestTeacherRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestTeacherRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher := createTestTeacher(t)
	insertTestSchool(t, db, teacher.SchoolID)

	// Create teacher first
	err := repo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	// Test successful retrieval
	retrieved, err := repo.GetByEmail(ctx, teacher.Email)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.Email != teacher.Email {
		t.Errorf("Expected email %s, got %s", teacher.Email, retrieved.Email)
	}
	if retrieved.ID != teacher.ID {
		t.Errorf("Expected ID %s, got %s", teacher.ID, retrieved.ID)
	}
}

func TestTeacherRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestTeacherRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher := createTestTeacher(t)
	insertTestSchool(t, db, teacher.SchoolID)

	// Create teacher first
	err := repo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	// Update teacher information
	err = teacher.UpdateProfile("Jane", "Smith", "3rd Grade")
	if err != nil {
		t.Fatalf("Failed to update teacher profile: %v", err)
	}

	err = teacher.UpdateWishlistURL("https://amazon.com/wishlist/test")
	if err != nil {
		t.Fatalf("Failed to update wishlist URL: %v", err)
	}

	// Update in repository
	err = repo.Update(ctx, teacher)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify updates
	retrieved, err := repo.GetByID(ctx, teacher.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated teacher: %v", err)
	}

	if retrieved.FirstName != "Jane" {
		t.Errorf("Expected first name Jane, got %s", retrieved.FirstName)
	}
	if retrieved.LastName != "Smith" {
		t.Errorf("Expected last name Smith, got %s", retrieved.LastName)
	}
	if retrieved.GradeLevel != "3rd Grade" {
		t.Errorf("Expected grade level 3rd Grade, got %s", retrieved.GradeLevel)
	}
	if retrieved.WishlistURL != "https://amazon.com/wishlist/test" {
		t.Errorf("Expected wishlist URL https://amazon.com/wishlist/test, got %s", retrieved.WishlistURL)
	}
}

func TestTeacherRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher := createTestTeacher(t)
	teacher.ID = "non-existent-id"

	err := repo.Update(ctx, teacher)
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestTeacherRepository_Update_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	// Create two teachers
	teacher1 := createTestTeacher(t)
	insertTestSchool(t, db, teacher1.SchoolID)

	teacher2, _ := teacherwishlist.NewTeacher(
		"teacher2@example.com",
		"Jane",
		"Smith",
		teacher1.SchoolID,
		"3rd Grade",
	)

	err := repo.Create(ctx, teacher1)
	if err != nil {
		t.Fatalf("Failed to create first teacher: %v", err)
	}

	err = repo.Create(ctx, teacher2)
	if err != nil {
		t.Fatalf("Failed to create second teacher: %v", err)
	}

	// Try to update teacher2 with teacher1's email
	teacher2.Email = teacher1.Email
	err = repo.Update(ctx, teacher2)
	if err == nil {
		t.Error("Expected duplicate key error, got nil")
	}
	if !IsDuplicateKeyError(err) {
		t.Errorf("Expected duplicate key error, got %v", err)
	}
}

func TestTeacherRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	teacher := createTestTeacher(t)
	insertTestSchool(t, db, teacher.SchoolID)

	// Create teacher first
	err := repo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	// Delete teacher
	err = repo.Delete(ctx, teacher.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify teacher was deleted
	_, err = repo.GetByID(ctx, teacher.ID)
	if err == nil {
		t.Error("Expected not found error after deletion, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error after deletion, got %v", err)
	}
}

func TestTeacherRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestTeacherRepository_CRUD_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTeacherRepository(db)
	ctx := context.Background()

	// Create a teacher
	teacher := createTestTeacher(t)
	insertTestSchool(t, db, teacher.SchoolID)

	// Test Create
	err := repo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	// Test Read by ID
	retrieved, err := repo.GetByID(ctx, teacher.ID)
	if err != nil {
		t.Fatalf("Failed to get teacher by ID: %v", err)
	}
	if retrieved.Email != teacher.Email {
		t.Errorf("Expected email %s, got %s", teacher.Email, retrieved.Email)
	}

	// Test Read by Email
	retrievedByEmail, err := repo.GetByEmail(ctx, teacher.Email)
	if err != nil {
		t.Fatalf("Failed to get teacher by email: %v", err)
	}
	if retrievedByEmail.ID != teacher.ID {
		t.Errorf("Expected ID %s, got %s", teacher.ID, retrievedByEmail.ID)
	}

	// Test Update
	err = teacher.UpdateProfile("Updated", "Name", "Updated Grade")
	if err != nil {
		t.Fatalf("Failed to update teacher profile: %v", err)
	}

	err = repo.Update(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to update teacher in repository: %v", err)
	}

	updated, err := repo.GetByID(ctx, teacher.ID)
	if err != nil {
		t.Fatalf("Failed to get updated teacher: %v", err)
	}
	if updated.FirstName != "Updated" {
		t.Errorf("Expected first name Updated, got %s", updated.FirstName)
	}

	// Test Delete
	err = repo.Delete(ctx, teacher.ID)
	if err != nil {
		t.Fatalf("Failed to delete teacher: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, teacher.ID)
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error after deletion, got %v", err)
	}
}
