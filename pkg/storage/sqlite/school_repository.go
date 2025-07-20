package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"hrh-backend/internal/schooldirectory"
	"hrh-backend/internal/shared/domain"
)

// SchoolRepository implements the schooldirectory.SchoolRepository interface for SQLite
type SchoolRepository struct {
	db DBInterface
}

// NewSchoolRepository creates a new SQLite school repository
func NewSchoolRepository(db DBInterface) *SchoolRepository {
	return &SchoolRepository{db: db}
}

// Create persists a new school aggregate with proper Value Object flattening
func (r *SchoolRepository) Create(ctx context.Context, school *schooldirectory.School) error {
	// Validate aggregate invariants before persistence
	if err := school.ValidateInvariants(); err != nil {
		return fmt.Errorf("school invariant validation failed: %w", err)
	}

	query := `
		INSERT INTO schools (
			id, name, 
			address_street, address_city, address_state, address_zip_code,
			location_latitude, location_longitude, location_county, location_region,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		school.ID,
		school.Name,
		// Address Value Object fields
		school.Address.Street,
		school.Address.City,
		school.Address.State,
		school.Address.ZipCode,
		// Location Value Object fields
		school.Address.Location.Latitude,
		school.Address.Location.Longitude,
		school.Address.Location.County,
		school.Address.Location.Region,
		school.CreatedAt,
		school.UpdatedAt,
	)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("school with name '%s' in %s, %s already exists: %w",
				school.Name, school.Address.City, school.Address.State, ErrDuplicateKey)
		}
		return fmt.Errorf("failed to create school: %w", err)
	}

	return nil
}

// GetByID retrieves a school by its unique identifier with Value Object reconstruction
func (r *SchoolRepository) GetByID(ctx context.Context, id string) (*schooldirectory.School, error) {
	query := `
		SELECT 
			id, name,
			address_street, address_city, address_state, address_zip_code,
			location_latitude, location_longitude, location_county, location_region,
			created_at, updated_at
		FROM schools
		WHERE id = ?`

	var school schooldirectory.School
	var addressStreet, addressCity, addressState, addressZipCode sql.NullString
	var locationLatitude, locationLongitude sql.NullFloat64
	var locationCounty, locationRegion sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&school.ID,
		&school.Name,
		// Address Value Object fields
		&addressStreet,
		&addressCity,
		&addressState,
		&addressZipCode,
		// Location Value Object fields
		&locationLatitude,
		&locationLongitude,
		&locationCounty,
		&locationRegion,
		&school.CreatedAt,
		&school.UpdatedAt,
	)

	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("school with ID %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get school by ID: %w", err)
	}

	// Reconstruct Address and Location Value Objects
	location, err := domain.NewLocation(
		locationLatitude.Float64,
		locationLongitude.Float64,
		locationCounty.String,
		locationRegion.String,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct location value object: %w", err)
	}

	address, err := domain.NewAddress(
		addressStreet.String,
		addressCity.String,
		addressState.String,
		addressZipCode.String,
		location,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct address value object: %w", err)
	}

	school.Address = address
	return &school, nil
}

// GetAll retrieves all schools from the repository with Value Object reconstruction
func (r *SchoolRepository) GetAll(ctx context.Context) ([]*schooldirectory.School, error) {
	query := `
		SELECT 
			id, name,
			address_street, address_city, address_state, address_zip_code,
			location_latitude, location_longitude, location_county, location_region,
			created_at, updated_at
		FROM schools
		ORDER BY name, address_city, address_state`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query schools: %w", err)
	}
	defer rows.Close()

	var schools []*schooldirectory.School

	for rows.Next() {
		var school schooldirectory.School
		var addressStreet, addressCity, addressState, addressZipCode sql.NullString
		var locationLatitude, locationLongitude sql.NullFloat64
		var locationCounty, locationRegion sql.NullString

		err := rows.Scan(
			&school.ID,
			&school.Name,
			// Address Value Object fields
			&addressStreet,
			&addressCity,
			&addressState,
			&addressZipCode,
			// Location Value Object fields
			&locationLatitude,
			&locationLongitude,
			&locationCounty,
			&locationRegion,
			&school.CreatedAt,
			&school.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan school row: %w", err)
		}

		// Reconstruct Address and Location Value Objects
		location, err := domain.NewLocation(
			locationLatitude.Float64,
			locationLongitude.Float64,
			locationCounty.String,
			locationRegion.String,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct location value object for school %s: %w", school.ID, err)
		}

		address, err := domain.NewAddress(
			addressStreet.String,
			addressCity.String,
			addressState.String,
			addressZipCode.String,
			location,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct address value object for school %s: %w", school.ID, err)
		}

		school.Address = address
		schools = append(schools, &school)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating school rows: %w", err)
	}

	return schools, nil
}

// Update persists changes to an existing school aggregate with Value Object flattening
func (r *SchoolRepository) Update(ctx context.Context, school *schooldirectory.School) error {
	// Validate aggregate invariants before persistence
	if err := school.ValidateInvariants(); err != nil {
		return fmt.Errorf("school invariant validation failed: %w", err)
	}

	query := `
		UPDATE schools 
		SET name = ?, 
		    address_street = ?, address_city = ?, address_state = ?, address_zip_code = ?,
		    location_latitude = ?, location_longitude = ?, location_county = ?, location_region = ?,
		    updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		school.Name,
		// Address Value Object fields
		school.Address.Street,
		school.Address.City,
		school.Address.State,
		school.Address.ZipCode,
		// Location Value Object fields
		school.Address.Location.Latitude,
		school.Address.Location.Longitude,
		school.Address.Location.County,
		school.Address.Location.Region,
		school.UpdatedAt,
		school.ID,
	)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("school with name '%s' in %s, %s already exists: %w",
				school.Name, school.Address.City, school.Address.State, ErrDuplicateKey)
		}
		return fmt.Errorf("failed to update school: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("school with ID %s not found: %w", school.ID, ErrNotFound)
	}

	return nil
}

// Delete removes a school from the repository
func (r *SchoolRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM schools WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		// Check if this is a foreign key constraint violation (teachers referencing this school)
		if IsForeignKeyError(err) {
			return fmt.Errorf("cannot delete school: teachers are still associated with this school: %w", ErrForeignKey)
		}
		return fmt.Errorf("failed to delete school: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("school with ID %s not found: %w", id, ErrNotFound)
	}

	return nil
}

// GetByNameCityState retrieves a school by its unique combination of name, city, and state
// This method supports the cross-aggregate invariant enforcement
func (r *SchoolRepository) GetByNameCityState(ctx context.Context, name, city, state string) (*schooldirectory.School, error) {
	query := `
		SELECT 
			id, name,
			address_street, address_city, address_state, address_zip_code,
			location_latitude, location_longitude, location_county, location_region,
			created_at, updated_at
		FROM schools
		WHERE LOWER(TRIM(name)) = LOWER(TRIM(?)) 
		  AND LOWER(TRIM(address_city)) = LOWER(TRIM(?)) 
		  AND LOWER(TRIM(address_state)) = LOWER(TRIM(?))`

	var school schooldirectory.School
	var addressStreet, addressCity, addressState, addressZipCode sql.NullString
	var locationLatitude, locationLongitude sql.NullFloat64
	var locationCounty, locationRegion sql.NullString

	err := r.db.QueryRowContext(ctx, query, name, city, state).Scan(
		&school.ID,
		&school.Name,
		// Address Value Object fields
		&addressStreet,
		&addressCity,
		&addressState,
		&addressZipCode,
		// Location Value Object fields
		&locationLatitude,
		&locationLongitude,
		&locationCounty,
		&locationRegion,
		&school.CreatedAt,
		&school.UpdatedAt,
	)

	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("school with name '%s' in %s, %s not found: %w", name, city, state, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get school by name/city/state: %w", err)
	}

	// Reconstruct Address and Location Value Objects
	location, err := domain.NewLocation(
		locationLatitude.Float64,
		locationLongitude.Float64,
		locationCounty.String,
		locationRegion.String,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct location value object: %w", err)
	}

	address, err := domain.NewAddress(
		addressStreet.String,
		addressCity.String,
		addressState.String,
		addressZipCode.String,
		location,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct address value object: %w", err)
	}

	school.Address = address
	return &school, nil
}
