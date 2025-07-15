package schooldirectory

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"hrh-backend/internal/shared/domain"
)

// School is an Aggregate Root that manages school information and location data.
// Internal Invariants:
// - A School must have a unique combination of [Name, City, State]
// - A School must have a valid Address with Location information
// - A School's name cannot be empty
type School struct {
	ID        string         `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Address   domain.Address `json:"address"` // Address Value Object composition
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// NewSchool creates a new School aggregate root with validation
func NewSchool(name string, address domain.Address) (*School, error) {
	// Validate required fields
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("school name cannot be empty")
	}

	// Validate address
	if address.IsEmpty() {
		return nil, fmt.Errorf("school address cannot be empty")
	}

	// Generate unique ID
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate school ID: %w", err)
	}

	now := time.Now()

	school := &School{
		ID:        id,
		Name:      name,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return school, nil
}

// UpdateName updates the school's name with validation
func (s *School) UpdateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("school name cannot be empty")
	}

	s.Name = name
	s.UpdatedAt = time.Now()
	return nil
}

// UpdateAddress updates the school's address with validation
func (s *School) UpdateAddress(address domain.Address) error {
	if address.IsEmpty() {
		return fmt.Errorf("school address cannot be empty")
	}

	s.Address = address
	s.UpdatedAt = time.Now()
	return nil
}

// GetUniqueKey returns the unique key for this school (name+city+state)
// This is used to enforce the cross-aggregate invariant of school uniqueness
func (s *School) GetUniqueKey() string {
	return fmt.Sprintf("%s|%s|%s",
		strings.ToLower(strings.TrimSpace(s.Name)),
		strings.ToLower(strings.TrimSpace(s.Address.City)),
		strings.ToLower(strings.TrimSpace(s.Address.State)))
}

// ValidateInvariants validates all internal aggregate invariants
func (s *School) ValidateInvariants() error {
	// Validate name is not empty
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("school name cannot be empty")
	}

	// Validate address is not empty
	if s.Address.IsEmpty() {
		return fmt.Errorf("school address cannot be empty")
	}

	// Validate address has required components for uniqueness constraint
	if strings.TrimSpace(s.Address.City) == "" {
		return fmt.Errorf("school address must have a city")
	}

	if strings.TrimSpace(s.Address.State) == "" {
		return fmt.Errorf("school address must have a state")
	}

	return nil
}

// String returns a string representation of the school
func (s *School) String() string {
	return fmt.Sprintf("%s (%s, %s)", s.Name, s.Address.City, s.Address.State)
}

// generateID generates a unique identifier for the school
func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}
