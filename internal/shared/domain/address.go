package domain

import (
	"errors"
	"regexp"
	"strings"
)

// Address represents a physical address with standardized components
// This is a Value Object - immutable and defined by its attributes
type Address struct {
	Street   string   `json:"street"`
	City     string   `json:"city"`
	State    string   `json:"state"`
	ZipCode  string   `json:"zip_code"`
	Location Location `json:"location"`
}

// NewAddress creates a new Address with validation
func NewAddress(street, city, state, zipCode string, location Location) (Address, error) {
	addr := Address{
		Street:   strings.TrimSpace(street),
		City:     strings.TrimSpace(city),
		State:    strings.TrimSpace(state),
		ZipCode:  strings.TrimSpace(zipCode),
		Location: location,
	}

	if err := addr.validate(); err != nil {
		return Address{}, err
	}

	return addr, nil
}

// WithStreet returns a new Address with updated street
func (a Address) WithStreet(street string) (Address, error) {
	return NewAddress(street, a.City, a.State, a.ZipCode, a.Location)
}

// WithCity returns a new Address with updated city
func (a Address) WithCity(city string) (Address, error) {
	return NewAddress(a.Street, city, a.State, a.ZipCode, a.Location)
}

// WithState returns a new Address with updated state
func (a Address) WithState(state string) (Address, error) {
	return NewAddress(a.Street, a.City, state, a.ZipCode, a.Location)
}

// WithZipCode returns a new Address with updated zip code
func (a Address) WithZipCode(zipCode string) (Address, error) {
	return NewAddress(a.Street, a.City, a.State, zipCode, a.Location)
}

// WithLocation returns a new Address with updated location
func (a Address) WithLocation(location Location) (Address, error) {
	return NewAddress(a.Street, a.City, a.State, a.ZipCode, location)
}

// validate performs validation on address components
func (a Address) validate() error {
	if a.City == "" {
		return errors.New("city is required")
	}

	if a.State == "" {
		return errors.New("state is required")
	}

	// Validate state format (2-letter state code)
	if len(a.State) != 2 {
		return errors.New("state must be a 2-letter code")
	}

	// Validate zip code format (basic US zip code validation)
	if a.ZipCode != "" {
		zipRegex := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
		if !zipRegex.MatchString(a.ZipCode) {
			return errors.New("zip code must be in format 12345 or 12345-6789")
		}
	}

	return nil
}

// IsEmpty returns true if the address has no meaningful data
func (a Address) IsEmpty() bool {
	return a.Street == "" && a.City == "" && a.State == "" && a.ZipCode == ""
}

// String returns a formatted string representation of the address
func (a Address) String() string {
	parts := []string{}

	if a.Street != "" {
		parts = append(parts, a.Street)
	}

	cityState := strings.TrimSpace(a.City + " " + a.State)
	if cityState != "" {
		parts = append(parts, cityState)
	}

	if a.ZipCode != "" {
		parts = append(parts, a.ZipCode)
	}

	return strings.Join(parts, ", ")
}

// Equals compares two addresses for equality
func (a Address) Equals(other Address) bool {
	return a.Street == other.Street &&
		a.City == other.City &&
		a.State == other.State &&
		a.ZipCode == other.ZipCode &&
		a.Location.Equals(other.Location)
}
