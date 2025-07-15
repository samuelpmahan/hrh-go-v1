package domain

import (
	"testing"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name     string
		street   string
		city     string
		state    string
		zipCode  string
		location Location
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid address",
			street:   "123 Main St",
			city:     "Springfield",
			state:    "IL",
			zipCode:  "62701",
			location: Location{},
			wantErr:  false,
		},
		{
			name:     "valid address with extended zip",
			street:   "456 Oak Ave",
			city:     "Chicago",
			state:    "IL",
			zipCode:  "60601-1234",
			location: Location{},
			wantErr:  false,
		},
		{
			name:     "missing city",
			street:   "123 Main St",
			city:     "",
			state:    "IL",
			zipCode:  "62701",
			location: Location{},
			wantErr:  true,
			errMsg:   "city is required",
		},
		{
			name:     "missing state",
			street:   "123 Main St",
			city:     "Springfield",
			state:    "",
			zipCode:  "62701",
			location: Location{},
			wantErr:  true,
			errMsg:   "state is required",
		},
		{
			name:     "invalid state format",
			street:   "123 Main St",
			city:     "Springfield",
			state:    "Illinois",
			zipCode:  "62701",
			location: Location{},
			wantErr:  true,
			errMsg:   "state must be a 2-letter code",
		},
		{
			name:     "invalid zip code format",
			street:   "123 Main St",
			city:     "Springfield",
			state:    "IL",
			zipCode:  "1234",
			location: Location{},
			wantErr:  true,
			errMsg:   "zip code must be in format 12345 or 12345-6789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := NewAddress(tt.street, tt.city, tt.state, tt.zipCode, tt.location)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewAddress() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewAddress() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewAddress() unexpected error = %v", err)
				return
			}

			if addr.Street != tt.street {
				t.Errorf("NewAddress() street = %v, want %v", addr.Street, tt.street)
			}
			if addr.City != tt.city {
				t.Errorf("NewAddress() city = %v, want %v", addr.City, tt.city)
			}
			if addr.State != tt.state {
				t.Errorf("NewAddress() state = %v, want %v", addr.State, tt.state)
			}
			if addr.ZipCode != tt.zipCode {
				t.Errorf("NewAddress() zipCode = %v, want %v", addr.ZipCode, tt.zipCode)
			}
		})
	}
}

func TestAddress_ImmutabilityMethods(t *testing.T) {
	location, _ := NewLocation(40.7128, -74.0060, "New York", "NY")
	original, _ := NewAddress("123 Main St", "Springfield", "IL", "62701", location)

	// Test WithStreet
	newAddr, err := original.WithStreet("456 Oak Ave")
	if err != nil {
		t.Errorf("WithStreet() unexpected error = %v", err)
	}
	if newAddr.Street != "456 Oak Ave" {
		t.Errorf("WithStreet() street = %v, want %v", newAddr.Street, "456 Oak Ave")
	}
	if original.Street != "123 Main St" {
		t.Errorf("WithStreet() modified original address")
	}

	// Test WithCity
	newAddr, err = original.WithCity("Chicago")
	if err != nil {
		t.Errorf("WithCity() unexpected error = %v", err)
	}
	if newAddr.City != "Chicago" {
		t.Errorf("WithCity() city = %v, want %v", newAddr.City, "Chicago")
	}
	if original.City != "Springfield" {
		t.Errorf("WithCity() modified original address")
	}

	// Test WithState
	newAddr, err = original.WithState("NY")
	if err != nil {
		t.Errorf("WithState() unexpected error = %v", err)
	}
	if newAddr.State != "NY" {
		t.Errorf("WithState() state = %v, want %v", newAddr.State, "NY")
	}
	if original.State != "IL" {
		t.Errorf("WithState() modified original address")
	}

	// Test WithZipCode
	newAddr, err = original.WithZipCode("60601")
	if err != nil {
		t.Errorf("WithZipCode() unexpected error = %v", err)
	}
	if newAddr.ZipCode != "60601" {
		t.Errorf("WithZipCode() zipCode = %v, want %v", newAddr.ZipCode, "60601")
	}
	if original.ZipCode != "62701" {
		t.Errorf("WithZipCode() modified original address")
	}
}

func TestAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		address  Address
		expected bool
	}{
		{
			name:     "empty address",
			address:  Address{},
			expected: true,
		},
		{
			name: "address with only street",
			address: Address{
				Street: "123 Main St",
			},
			expected: false,
		},
		{
			name: "address with city and state",
			address: Address{
				City:  "Springfield",
				State: "IL",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.address.IsEmpty()
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		address  Address
		expected string
	}{
		{
			name:     "empty address",
			address:  Address{},
			expected: "",
		},
		{
			name: "full address",
			address: Address{
				Street:  "123 Main St",
				City:    "Springfield",
				State:   "IL",
				ZipCode: "62701",
			},
			expected: "123 Main St, Springfield IL, 62701",
		},
		{
			name: "address without street",
			address: Address{
				City:    "Springfield",
				State:   "IL",
				ZipCode: "62701",
			},
			expected: "Springfield IL, 62701",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.address.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAddress_Equals(t *testing.T) {
	location1, _ := NewLocation(40.7128, -74.0060, "New York", "NY")
	location2, _ := NewLocation(41.8781, -87.6298, "Cook", "IL")

	addr1, _ := NewAddress("123 Main St", "Springfield", "IL", "62701", location1)
	addr2, _ := NewAddress("123 Main St", "Springfield", "IL", "62701", location1)
	addr3, _ := NewAddress("456 Oak Ave", "Springfield", "IL", "62701", location1)
	addr4, _ := NewAddress("123 Main St", "Springfield", "IL", "62701", location2)

	tests := []struct {
		name     string
		addr1    Address
		addr2    Address
		expected bool
	}{
		{
			name:     "identical addresses",
			addr1:    addr1,
			addr2:    addr2,
			expected: true,
		},
		{
			name:     "different streets",
			addr1:    addr1,
			addr2:    addr3,
			expected: false,
		},
		{
			name:     "different locations",
			addr1:    addr1,
			addr2:    addr4,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.addr1.Equals(tt.addr2)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}
