package schooldirectory

import (
	"context"
	"fmt"
	"strings"
	"sync" // For thread-safe map access in the in-memory implementation
)

// SchoolFilters represents filtering criteria for school queries
type SchoolFilters struct {
	Name   *string
	City   *string
	State  *string
	Limit  int
	Offset int
}

// SchoolRepository defines the interface for persisting and retrieving School aggregates.
// This interface follows aggregate-focused methods and includes proper context handling.
type SchoolRepository interface {
	// Create persists a new school aggregate
	Create(ctx context.Context, school *School) error

	// GetByID retrieves a school by its unique identifier
	GetByID(ctx context.Context, id string) (*School, error)

	// GetByNameCityState retrieves a school by its unique combination of name, city, and state
	// This supports the cross-aggregate invariant: "[SchoolName, City, State] is a unique combination"
	GetByNameCityState(ctx context.Context, name, city, state string) (*School, error)

	// Update persists changes to an existing school aggregate
	Update(ctx context.Context, school *School) error

	// List retrieves schools based on filtering criteria with pagination support
	List(ctx context.Context, filters SchoolFilters) ([]*School, error)

	// GetAll retrieves all schools (for admin management and directory browsing)
	GetAll(ctx context.Context) ([]*School, error)

	// Delete removes a school aggregate (for admin management)
	Delete(ctx context.Context, id string) error

	// ExistsByNameCityState checks if a school with the given name, city, and state combination already exists
	ExistsByNameCityState(ctx context.Context, name, city, state string) (bool, error)

	// Search performs text-based search across school names and locations
	Search(ctx context.Context, query string, limit int) ([]*School, error)
}

// InMemorySchoolRepository is a concrete implementation of SchoolRepository
// that stores data in memory. Suitable for testing and development.
type InMemorySchoolRepository struct {
	mu      sync.RWMutex // Protects map access for concurrency
	schools map[string]*School
	// An additional map to quickly check for unique combinations for the invariant.
	// Key: "name|city|state"
	uniqueCombos map[string]*School
}

// NewInMemorySchoolRepository creates and returns a new InMemorySchoolRepository.
func NewInMemorySchoolRepository() *InMemorySchoolRepository {
	return &InMemorySchoolRepository{
		schools:      make(map[string]*School),
		uniqueCombos: make(map[string]*School),
	}
}

// Create persists a new school aggregate
func (r *InMemorySchoolRepository) Create(ctx context.Context, school *School) error {
	return r.Save(school)
}

// Update persists changes to an existing school aggregate
func (r *InMemorySchoolRepository) Update(ctx context.Context, school *School) error {
	return r.Save(school)
}

// GetByID retrieves a school by its unique identifier
func (r *InMemorySchoolRepository) GetByID(ctx context.Context, id string) (*School, error) {
	return r.GetByID_Legacy(id)
}

// GetByNameCityState retrieves a school by its unique combination of name, city, and state
func (r *InMemorySchoolRepository) GetByNameCityState(ctx context.Context, name, city, state string) (*School, error) {
	return r.FindByUniqueCombination(name, city, state)
}

// List retrieves schools based on filtering criteria with pagination support
func (r *InMemorySchoolRepository) List(ctx context.Context, filters SchoolFilters) ([]*School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*School
	for _, school := range r.schools {
		// Apply filters
		if filters.Name != nil && !strings.Contains(strings.ToLower(school.Name), strings.ToLower(*filters.Name)) {
			continue
		}
		if filters.City != nil && !strings.Contains(strings.ToLower(school.Address.City), strings.ToLower(*filters.City)) {
			continue
		}
		if filters.State != nil && !strings.Contains(strings.ToLower(school.Address.State), strings.ToLower(*filters.State)) {
			continue
		}
		results = append(results, school)
	}

	// Apply pagination
	if filters.Offset > 0 && filters.Offset < len(results) {
		results = results[filters.Offset:]
	}
	if filters.Limit > 0 && filters.Limit < len(results) {
		results = results[:filters.Limit]
	}

	return results, nil
}

// GetAll retrieves all schools
func (r *InMemorySchoolRepository) GetAll(ctx context.Context) ([]*School, error) {
	return r.GetAll_Legacy()
}

// Delete removes a school aggregate
func (r *InMemorySchoolRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	school, found := r.schools[id]
	if !found {
		return fmt.Errorf("school with ID %s not found", id)
	}

	// Remove from both maps
	delete(r.schools, id)
	delete(r.uniqueCombos, school.GetUniqueKey())
	return nil
}

// ExistsByNameCityState checks if a school with the given name, city, and state combination already exists
func (r *InMemorySchoolRepository) ExistsByNameCityState(ctx context.Context, name, city, state string) (bool, error) {
	_, err := r.FindByUniqueCombination(name, city, state)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Search performs text-based search across school names and locations
func (r *InMemorySchoolRepository) Search(ctx context.Context, query string, limit int) ([]*School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	var results []*School

	for _, school := range r.schools {
		// Search in name, city, and state
		if strings.Contains(strings.ToLower(school.Name), query) ||
			strings.Contains(strings.ToLower(school.Address.City), query) ||
			strings.Contains(strings.ToLower(school.Address.State), query) {
			results = append(results, school)
		}

		// Apply limit
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	return results, nil
}

// Legacy methods for backward compatibility

// Save adds or updates a school in the in-memory store.
// It also updates the unique combination index.
func (r *InMemorySchoolRepository) Save(school *School) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Before saving, if a school with this ID already exists,
	// remove its old unique combo entry if the name/location changed.
	if existingSchool, found := r.schools[school.ID]; found {
		oldComboKey := existingSchool.GetUniqueKey()
		delete(r.uniqueCombos, oldComboKey)
	}

	comboKey := school.GetUniqueKey()
	// Enforce local uniqueness on save - although a cross-aggregate invariant,
	// the repository can prevent *itself* from storing duplicates.
	// A higher-level application service would use FindByUniqueCombination *before* calling Save.
	if existingSchool, found := r.uniqueCombos[comboKey]; found && existingSchool.ID != school.ID {
		return fmt.Errorf("school with name '%s', city '%s', and state '%s' already exists (ID: %s)",
			school.Name, school.Address.City, school.Address.State, existingSchool.ID)
	}

	r.schools[school.ID] = school
	r.uniqueCombos[comboKey] = school
	return nil
}

// GetByID_Legacy retrieves a school by its ID.
func (r *InMemorySchoolRepository) GetByID_Legacy(id string) (*School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	school, found := r.schools[id]
	if !found {
		return nil, fmt.Errorf("school with ID %s not found", id)
	}
	return school, nil
}

// FindByUniqueCombination finds a school by its unique name, city, and state combination.
func (r *InMemorySchoolRepository) FindByUniqueCombination(name, city, state string) (*School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create the same key format as School.GetUniqueKey()
	comboKey := fmt.Sprintf("%s|%s|%s",
		strings.ToLower(strings.TrimSpace(name)),
		strings.ToLower(strings.TrimSpace(city)),
		strings.ToLower(strings.TrimSpace(state)))

	school, found := r.uniqueCombos[comboKey]
	if !found {
		return nil, fmt.Errorf("school with name '%s', city '%s', state '%s' not found", name, city, state)
	}
	return school, nil
}

// GetAll_Legacy returns all schools currently in the repository.
func (r *InMemorySchoolRepository) GetAll_Legacy() ([]*School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	schools := make([]*School, 0, len(r.schools))
	for _, school := range r.schools {
		schools = append(schools, school)
	}
	return schools, nil
}

// District-related repository implementation will be added in future tasks
