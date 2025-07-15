package schooldirectory

import (
	"strings"
	"testing"
	"time"

	"hrh-backend/internal/shared/domain"
)

func TestNewSchool(t *testing.T) {
	// Create a valid address for testing
	location, err := domain.NewLocation(40.7128, -74.0060, "New York", "NY")
	if err != nil {
		t.Fatalf("Failed to create location: %v", err)
	}

	validAddress, err := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)
	if err != nil {
		t.Fatalf("Failed to create address: %v", err)
	}

	tests := []struct {
		name        string
		schoolName  string
		address     domain.Address
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid school creation",
			schoolName: "Lincoln Elementary",
			address:    validAddress,
			wantErr:    false,
		},
		{
			name:        "empty school name",
			schoolName:  "",
			address:     validAddress,
			wantErr:     true,
			errContains: "school name cannot be empty",
		},
		{
			name:        "whitespace only school name",
			schoolName:  "   ",
			address:     validAddress,
			wantErr:     true,
			errContains: "school name cannot be empty",
		},
		{
			name:        "empty address",
			schoolName:  "Lincoln Elementary",
			address:     domain.Address{}, // empty address
			wantErr:     true,
			errContains: "school address cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			school, err := NewSchool(tt.schoolName, tt.address)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSchool() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewSchool() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewSchool() unexpected error = %v", err)
				return
			}

			// Validate the created school
			if school == nil {
				t.Error("NewSchool() returned nil school")
				return
			}

			if school.ID == "" {
				t.Error("NewSchool() school ID is empty")
			}

			if school.Name != strings.TrimSpace(tt.schoolName) {
				t.Errorf("NewSchool() school name = %v, want %v", school.Name, strings.TrimSpace(tt.schoolName))
			}

			if !school.Address.Equals(tt.address) {
				t.Errorf("NewSchool() school address = %v, want %v", school.Address, tt.address)
			}

			if school.CreatedAt.IsZero() {
				t.Error("NewSchool() CreatedAt is zero")
			}

			if school.UpdatedAt.IsZero() {
				t.Error("NewSchool() UpdatedAt is zero")
			}

			if !school.CreatedAt.Equal(school.UpdatedAt) {
				t.Error("NewSchool() CreatedAt and UpdatedAt should be equal for new school")
			}
		})
	}
}

func TestSchool_UpdateName(t *testing.T) {
	// Create a valid school for testing
	location, _ := domain.NewLocation(40.7128, -74.0060, "New York", "NY")
	address, _ := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)
	school, _ := NewSchool("Original Name", address)
	originalUpdatedAt := school.UpdatedAt

	// Wait a small amount to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name        string
		newName     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid name update",
			newName: "Updated School Name",
			wantErr: false,
		},
		{
			name:        "empty name",
			newName:     "",
			wantErr:     true,
			errContains: "school name cannot be empty",
		},
		{
			name:        "whitespace only name",
			newName:     "   ",
			wantErr:     true,
			errContains: "school name cannot be empty",
		},
		{
			name:    "name with leading/trailing whitespace",
			newName: "  Trimmed Name  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh school for each test
			testSchool, _ := NewSchool("Test School", address)
			testSchool.UpdatedAt = originalUpdatedAt // Reset to test timestamp update

			err := testSchool.UpdateName(tt.newName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateName() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateName() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateName() unexpected error = %v", err)
				return
			}

			expectedName := strings.TrimSpace(tt.newName)
			if testSchool.Name != expectedName {
				t.Errorf("UpdateName() school name = %v, want %v", testSchool.Name, expectedName)
			}

			if !testSchool.UpdatedAt.After(originalUpdatedAt) {
				t.Error("UpdateName() should update the UpdatedAt timestamp")
			}
		})
	}
}

func TestSchool_UpdateAddress(t *testing.T) {
	// Create a valid school for testing
	location, _ := domain.NewLocation(40.7128, -74.0060, "New York", "NY")
	originalAddress, _ := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)
	school, _ := NewSchool("Test School", originalAddress)
	originalUpdatedAt := school.UpdatedAt

	// Wait a small amount to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Create a new valid address
	newLocation, _ := domain.NewLocation(34.0522, -118.2437, "Los Angeles", "CA")
	newAddress, _ := domain.NewAddress("456 Oak Ave", "Los Angeles", "CA", "90210", newLocation)

	tests := []struct {
		name        string
		newAddress  domain.Address
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid address update",
			newAddress: newAddress,
			wantErr:    false,
		},
		{
			name:        "empty address",
			newAddress:  domain.Address{}, // empty address
			wantErr:     true,
			errContains: "school address cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh school for each test
			testSchool, _ := NewSchool("Test School", originalAddress)
			testSchool.UpdatedAt = originalUpdatedAt // Reset to test timestamp update

			err := testSchool.UpdateAddress(tt.newAddress)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateAddress() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateAddress() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateAddress() unexpected error = %v", err)
				return
			}

			if !testSchool.Address.Equals(tt.newAddress) {
				t.Errorf("UpdateAddress() school address = %v, want %v", testSchool.Address, tt.newAddress)
			}

			if !testSchool.UpdatedAt.After(originalUpdatedAt) {
				t.Error("UpdateAddress() should update the UpdatedAt timestamp")
			}
		})
	}
}

func TestSchool_GetUniqueKey(t *testing.T) {
	tests := []struct {
		name        string
		schoolName  string
		city        string
		state       string
		expectedKey string
	}{
		{
			name:        "basic unique key",
			schoolName:  "Lincoln Elementary",
			city:        "New York",
			state:       "NY",
			expectedKey: "lincoln elementary|new york|ny",
		},
		{
			name:        "case insensitive key",
			schoolName:  "LINCOLN ELEMENTARY",
			city:        "NEW YORK",
			state:       "ny",
			expectedKey: "lincoln elementary|new york|ny",
		},
		{
			name:        "whitespace trimming",
			schoolName:  "  Lincoln Elementary  ",
			city:        "  New York  ",
			state:       "  NY  ",
			expectedKey: "lincoln elementary|new york|ny",
		},
		{
			name:        "mixed case and whitespace",
			schoolName:  "  LiNcOlN ElEmEnTaRy  ",
			city:        "  NeW yOrK  ",
			state:       "  Ny  ",
			expectedKey: "lincoln elementary|new york|ny",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, _ := domain.NewLocation(40.7128, -74.0060, "County", "Region")
			address, _ := domain.NewAddress("123 Main St", tt.city, tt.state, "10001", location)
			school, _ := NewSchool(tt.schoolName, address)

			key := school.GetUniqueKey()
			if key != tt.expectedKey {
				t.Errorf("GetUniqueKey() = %v, want %v", key, tt.expectedKey)
			}
		})
	}
}

func TestSchool_ValidateInvariants(t *testing.T) {
	// Create valid components
	location, _ := domain.NewLocation(40.7128, -74.0060, "New York", "NY")
	validAddress, _ := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)

	tests := []struct {
		name        string
		setupSchool func() *School
		wantErr     bool
		errContains string
	}{
		{
			name: "valid school",
			setupSchool: func() *School {
				school, _ := NewSchool("Lincoln Elementary", validAddress)
				return school
			},
			wantErr: false,
		},
		{
			name: "empty school name",
			setupSchool: func() *School {
				school, _ := NewSchool("Lincoln Elementary", validAddress)
				school.Name = ""
				return school
			},
			wantErr:     true,
			errContains: "school name cannot be empty",
		},
		{
			name: "whitespace only school name",
			setupSchool: func() *School {
				school, _ := NewSchool("Lincoln Elementary", validAddress)
				school.Name = "   "
				return school
			},
			wantErr:     true,
			errContains: "school name cannot be empty",
		},
		{
			name: "empty address",
			setupSchool: func() *School {
				school, _ := NewSchool("Lincoln Elementary", validAddress)
				school.Address = domain.Address{} // empty address
				return school
			},
			wantErr:     true,
			errContains: "school address cannot be empty",
		},
		{
			name: "address without city",
			setupSchool: func() *School {
				school, _ := NewSchool("Lincoln Elementary", validAddress)
				// Directly set an address with empty city to bypass domain validation
				school.Address = domain.Address{
					Street:   "123 Main St",
					City:     "",
					State:    "NY",
					ZipCode:  "10001",
					Location: location,
				}
				return school
			},
			wantErr:     true,
			errContains: "school address must have a city",
		},
		{
			name: "address without state",
			setupSchool: func() *School {
				school, _ := NewSchool("Lincoln Elementary", validAddress)
				// Directly set an address with empty state to bypass domain validation
				school.Address = domain.Address{
					Street:   "123 Main St",
					City:     "New York",
					State:    "",
					ZipCode:  "10001",
					Location: location,
				}
				return school
			},
			wantErr:     true,
			errContains: "school address must have a state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			school := tt.setupSchool()
			err := school.ValidateInvariants()

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateInvariants() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateInvariants() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateInvariants() unexpected error = %v", err)
			}
		})
	}
}

func TestSchool_String(t *testing.T) {
	location, _ := domain.NewLocation(40.7128, -74.0060, "New York", "NY")
	address, _ := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)
	school, _ := NewSchool("Lincoln Elementary", address)

	expected := "Lincoln Elementary (New York, NY)"
	result := school.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}

func TestGenerateID(t *testing.T) {
	// Test that generateID produces unique IDs
	ids := make(map[string]bool)

	for i := 0; i < 100; i++ {
		id, err := generateID()
		if err != nil {
			t.Errorf("generateID() unexpected error = %v", err)
		}

		if id == "" {
			t.Error("generateID() returned empty string")
		}

		if ids[id] {
			t.Errorf("generateID() produced duplicate ID: %v", id)
		}

		ids[id] = true

		// Verify ID format (should be hex with dashes)
		if !strings.Contains(id, "-") {
			t.Errorf("generateID() ID format incorrect: %v", id)
		}
	}
}

// Test that demonstrates the unique key invariant enforcement
func TestSchool_UniqueKeyInvariant(t *testing.T) {
	location, _ := domain.NewLocation(40.7128, -74.0060, "New York", "NY")
	address, _ := domain.NewAddress("123 Main St", "New York", "NY", "10001", location)

	// Create two schools with the same name, city, state
	school1, _ := NewSchool("Lincoln Elementary", address)
	school2, _ := NewSchool("Lincoln Elementary", address)

	// They should have the same unique key (this would be caught at the application layer)
	key1 := school1.GetUniqueKey()
	key2 := school2.GetUniqueKey()

	if key1 != key2 {
		t.Errorf("Schools with same name/city/state should have same unique key: %v != %v", key1, key2)
	}

	// But different IDs
	if school1.ID == school2.ID {
		t.Error("Schools should have different IDs even with same unique key")
	}
}
