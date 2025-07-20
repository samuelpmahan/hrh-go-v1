package sqlite

import (
	"context"
	"testing"

	"hrh-backend/internal/schooldirectory"
	"hrh-backend/internal/shared/domain"
)

// createTestSchool creates a test school for use in tests
func createTestSchool(t *testing.T) *schooldirectory.School {
	location, err := domain.NewLocation(40.7128, -74.0060, "New York County", "Northeast")
	if err != nil {
		t.Fatalf("Failed to create test location: %v", err)
	}

	address, err := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)
	if err != nil {
		t.Fatalf("Failed to create test address: %v", err)
	}

	school, err := schooldirectory.NewSchool("Test Elementary School", address)
	if err != nil {
		t.Fatalf("Failed to create test school: %v", err)
	}

	return school
}

// createTestSchoolWithName creates a test school with a specific name
func createTestSchoolWithName(t *testing.T, name, city, state string) *schooldirectory.School {
	location, err := domain.NewLocation(40.7128, -74.0060, "Test County", "Northeast")
	if err != nil {
		t.Fatalf("Failed to create test location: %v", err)
	}

	address, err := domain.NewAddress("123 Test St", city, state, "12345", location)
	if err != nil {
		t.Fatalf("Failed to create test address: %v", err)
	}

	school, err := schooldirectory.NewSchool(name, address)
	if err != nil {
		t.Fatalf("Failed to create test school: %v", err)
	}

	return school
}

func TestSchoolRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school := createTestSchool(t)

	// Test successful creation
	err := repo.Create(ctx, school)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify school was created by retrieving it
	retrieved, err := repo.GetByID(ctx, school.ID)
	if err != nil {
		t.Errorf("Failed to retrieve created school: %v", err)
	}

	// Verify basic fields
	if retrieved.Name != school.Name {
		t.Errorf("Expected name %s, got %s", school.Name, retrieved.Name)
	}
	if retrieved.ID != school.ID {
		t.Errorf("Expected ID %s, got %s", school.ID, retrieved.ID)
	}

	// Verify Address Value Object reconstruction
	if retrieved.Address.Street != school.Address.Street {
		t.Errorf("Expected street %s, got %s", school.Address.Street, retrieved.Address.Street)
	}
	if retrieved.Address.City != school.Address.City {
		t.Errorf("Expected city %s, got %s", school.Address.City, retrieved.Address.City)
	}
	if retrieved.Address.State != school.Address.State {
		t.Errorf("Expected state %s, got %s", school.Address.State, retrieved.Address.State)
	}
	if retrieved.Address.ZipCode != school.Address.ZipCode {
		t.Errorf("Expected zip code %s, got %s", school.Address.ZipCode, retrieved.Address.ZipCode)
	}

	// Verify Location Value Object reconstruction
	if retrieved.Address.Location.Latitude != school.Address.Location.Latitude {
		t.Errorf("Expected latitude %f, got %f", school.Address.Location.Latitude, retrieved.Address.Location.Latitude)
	}
	if retrieved.Address.Location.Longitude != school.Address.Location.Longitude {
		t.Errorf("Expected longitude %f, got %f", school.Address.Location.Longitude, retrieved.Address.Location.Longitude)
	}
	if retrieved.Address.Location.County != school.Address.Location.County {
		t.Errorf("Expected county %s, got %s", school.Address.Location.County, retrieved.Address.Location.County)
	}
	if retrieved.Address.Location.Region != school.Address.Location.Region {
		t.Errorf("Expected region %s, got %s", school.Address.Location.Region, retrieved.Address.Location.Region)
	}
}

func TestSchoolRepository_Create_DuplicateSchool(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school1 := createTestSchoolWithName(t, "Duplicate School", "Test City", "NY")

	// Create first school
	err := repo.Create(ctx, school1)
	if err != nil {
		t.Fatalf("Failed to create first school: %v", err)
	}

	// Try to create second school with same name, city, state (should violate unique constraint)
	school2 := createTestSchoolWithName(t, "Duplicate School", "Test City", "NY")

	err = repo.Create(ctx, school2)
	if err == nil {
		t.Error("Expected duplicate key error, got nil")
	}
	if !IsDuplicateKeyError(err) {
		t.Errorf("Expected duplicate key error, got %v", err)
	}
}

func TestSchoolRepository_Create_SameNameDifferentLocation(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	// Create two schools with same name but different cities (should be allowed)
	school1 := createTestSchoolWithName(t, "Lincoln Elementary", "New York", "NY")
	school2 := createTestSchoolWithName(t, "Lincoln Elementary", "Albany", "NY")

	err := repo.Create(ctx, school1)
	if err != nil {
		t.Fatalf("Failed to create first school: %v", err)
	}

	err = repo.Create(ctx, school2)
	if err != nil {
		t.Errorf("Expected no error for same name different city, got %v", err)
	}
}

func TestSchoolRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school := createTestSchool(t)

	// Create school first
	err := repo.Create(ctx, school)
	if err != nil {
		t.Fatalf("Failed to create school: %v", err)
	}

	// Test successful retrieval
	retrieved, err := repo.GetByID(ctx, school.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.ID != school.ID {
		t.Errorf("Expected ID %s, got %s", school.ID, retrieved.ID)
	}
	if retrieved.Name != school.Name {
		t.Errorf("Expected name %s, got %s", school.Name, retrieved.Name)
	}
}

func TestSchoolRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestSchoolRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	// Create multiple schools
	school1 := createTestSchoolWithName(t, "Alpha Elementary", "City A", "NY")
	school2 := createTestSchoolWithName(t, "Beta Elementary", "City B", "CA")
	school3 := createTestSchoolWithName(t, "Gamma Elementary", "City C", "TX")

	err := repo.Create(ctx, school1)
	if err != nil {
		t.Fatalf("Failed to create school1: %v", err)
	}

	err = repo.Create(ctx, school2)
	if err != nil {
		t.Fatalf("Failed to create school2: %v", err)
	}

	err = repo.Create(ctx, school3)
	if err != nil {
		t.Fatalf("Failed to create school3: %v", err)
	}

	// Test GetAll
	schools, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(schools) != 3 {
		t.Errorf("Expected 3 schools, got %d", len(schools))
	}

	// Verify schools are ordered by name, city, state
	if schools[0].Name != "Alpha Elementary" {
		t.Errorf("Expected first school to be Alpha Elementary, got %s", schools[0].Name)
	}
}

func TestSchoolRepository_GetAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	schools, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(schools) != 0 {
		t.Errorf("Expected 0 schools, got %d", len(schools))
	}
}

func TestSchoolRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school := createTestSchool(t)

	// Create school first
	err := repo.Create(ctx, school)
	if err != nil {
		t.Fatalf("Failed to create school: %v", err)
	}

	// Update school information
	err = school.UpdateName("Updated School Name")
	if err != nil {
		t.Fatalf("Failed to update school name: %v", err)
	}

	newLocation, err := domain.NewLocation(41.8781, -87.6298, "Cook County", "Midwest")
	if err != nil {
		t.Fatalf("Failed to create new location: %v", err)
	}

	newAddress, err := domain.NewAddress("456 Updated St", "Chicago", "IL", "60601", newLocation)
	if err != nil {
		t.Fatalf("Failed to create new address: %v", err)
	}

	err = school.UpdateAddress(newAddress)
	if err != nil {
		t.Fatalf("Failed to update school address: %v", err)
	}

	// Update in repository
	err = repo.Update(ctx, school)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify updates
	retrieved, err := repo.GetByID(ctx, school.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated school: %v", err)
	}

	if retrieved.Name != "Updated School Name" {
		t.Errorf("Expected name Updated School Name, got %s", retrieved.Name)
	}
	if retrieved.Address.Street != "456 Updated St" {
		t.Errorf("Expected street 456 Updated St, got %s", retrieved.Address.Street)
	}
	if retrieved.Address.City != "Chicago" {
		t.Errorf("Expected city Chicago, got %s", retrieved.Address.City)
	}
	if retrieved.Address.State != "IL" {
		t.Errorf("Expected state IL, got %s", retrieved.Address.State)
	}
	if retrieved.Address.Location.County != "Cook County" {
		t.Errorf("Expected county Cook County, got %s", retrieved.Address.Location.County)
	}
}

func TestSchoolRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school := createTestSchool(t)
	school.ID = "non-existent-id"

	err := repo.Update(ctx, school)
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestSchoolRepository_Update_DuplicateSchool(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	// Create two schools
	school1 := createTestSchoolWithName(t, "School One", "Test City", "NY")
	school2 := createTestSchoolWithName(t, "School Two", "Other City", "CA")

	err := repo.Create(ctx, school1)
	if err != nil {
		t.Fatalf("Failed to create first school: %v", err)
	}

	err = repo.Create(ctx, school2)
	if err != nil {
		t.Fatalf("Failed to create second school: %v", err)
	}

	// Try to update school2 to have same name/city/state as school1
	err = school2.UpdateName("School One")
	if err != nil {
		t.Fatalf("Failed to update school name: %v", err)
	}

	newLocation, err := domain.NewLocation(40.7128, -74.0060, "Test County", "Northeast")
	if err != nil {
		t.Fatalf("Failed to create new location: %v", err)
	}

	newAddress, err := domain.NewAddress("123 Test St", "Test City", "NY", "12345", newLocation)
	if err != nil {
		t.Fatalf("Failed to create new address: %v", err)
	}

	err = school2.UpdateAddress(newAddress)
	if err != nil {
		t.Fatalf("Failed to update school address: %v", err)
	}

	err = repo.Update(ctx, school2)
	if err == nil {
		t.Error("Expected duplicate key error, got nil")
	}
	if !IsDuplicateKeyError(err) {
		t.Errorf("Expected duplicate key error, got %v", err)
	}
}

func TestSchoolRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school := createTestSchool(t)

	// Create school first
	err := repo.Create(ctx, school)
	if err != nil {
		t.Fatalf("Failed to create school: %v", err)
	}

	// Delete school
	err = repo.Delete(ctx, school.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify school was deleted
	_, err = repo.GetByID(ctx, school.ID)
	if err == nil {
		t.Error("Expected not found error after deletion, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error after deletion, got %v", err)
	}
}

func TestSchoolRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestSchoolRepository_GetByNameCityState(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	school := createTestSchoolWithName(t, "Unique School", "Unique City", "NY")

	// Create school first
	err := repo.Create(ctx, school)
	if err != nil {
		t.Fatalf("Failed to create school: %v", err)
	}

	// Test successful retrieval
	retrieved, err := repo.GetByNameCityState(ctx, "Unique School", "Unique City", "NY")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.ID != school.ID {
		t.Errorf("Expected ID %s, got %s", school.ID, retrieved.ID)
	}
	if retrieved.Name != school.Name {
		t.Errorf("Expected name %s, got %s", school.Name, retrieved.Name)
	}

	// Test case insensitive matching
	retrieved2, err := repo.GetByNameCityState(ctx, "unique school", "unique city", "ny")
	if err != nil {
		t.Errorf("Expected no error for case insensitive match, got %v", err)
	}

	if retrieved2.ID != school.ID {
		t.Errorf("Expected case insensitive match to return same school")
	}
}

func TestSchoolRepository_GetByNameCityState_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	_, err := repo.GetByNameCityState(ctx, "Non Existent", "No City", "XX")
	if err == nil {
		t.Error("Expected not found error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error, got %v", err)
	}
}

func TestSchoolRepository_CRUD_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSchoolRepository(db)
	ctx := context.Background()

	// Create a school
	school := createTestSchool(t)

	// Test Create
	err := repo.Create(ctx, school)
	if err != nil {
		t.Fatalf("Failed to create school: %v", err)
	}

	// Test Read by ID
	retrieved, err := repo.GetByID(ctx, school.ID)
	if err != nil {
		t.Fatalf("Failed to get school by ID: %v", err)
	}
	if retrieved.Name != school.Name {
		t.Errorf("Expected name %s, got %s", school.Name, retrieved.Name)
	}

	// Test Read by Name/City/State
	retrievedByKey, err := repo.GetByNameCityState(ctx, school.Name, school.Address.City, school.Address.State)
	if err != nil {
		t.Fatalf("Failed to get school by name/city/state: %v", err)
	}
	if retrievedByKey.ID != school.ID {
		t.Errorf("Expected ID %s, got %s", school.ID, retrievedByKey.ID)
	}

	// Test GetAll
	schools, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get all schools: %v", err)
	}
	if len(schools) != 1 {
		t.Errorf("Expected 1 school, got %d", len(schools))
	}

	// Test Update
	err = school.UpdateName("Updated School")
	if err != nil {
		t.Fatalf("Failed to update school name: %v", err)
	}

	err = repo.Update(ctx, school)
	if err != nil {
		t.Fatalf("Failed to update school in repository: %v", err)
	}

	updated, err := repo.GetByID(ctx, school.ID)
	if err != nil {
		t.Fatalf("Failed to get updated school: %v", err)
	}
	if updated.Name != "Updated School" {
		t.Errorf("Expected name Updated School, got %s", updated.Name)
	}

	// Test Delete
	err = repo.Delete(ctx, school.ID)
	if err != nil {
		t.Fatalf("Failed to delete school: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, school.ID)
	if !IsNotFoundError(err) {
		t.Errorf("Expected not found error after deletion, got %v", err)
	}
}
