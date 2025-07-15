package domain

import (
	"math"
	"testing"
)

func TestNewLocation(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		county    string
		region    string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid location",
			latitude:  40.7128,
			longitude: -74.0060,
			county:    "New York",
			region:    "NY",
			wantErr:   false,
		},
		{
			name:      "valid location at equator",
			latitude:  0,
			longitude: 0,
			county:    "",
			region:    "",
			wantErr:   false,
		},
		{
			name:      "latitude too high",
			latitude:  91,
			longitude: -74.0060,
			county:    "New York",
			region:    "NY",
			wantErr:   true,
			errMsg:    "latitude must be between -90 and 90 degrees",
		},
		{
			name:      "latitude too low",
			latitude:  -91,
			longitude: -74.0060,
			county:    "New York",
			region:    "NY",
			wantErr:   true,
			errMsg:    "latitude must be between -90 and 90 degrees",
		},
		{
			name:      "longitude too high",
			latitude:  40.7128,
			longitude: 181,
			county:    "New York",
			region:    "NY",
			wantErr:   true,
			errMsg:    "longitude must be between -180 and 180 degrees",
		},
		{
			name:      "longitude too low",
			latitude:  40.7128,
			longitude: -181,
			county:    "New York",
			region:    "NY",
			wantErr:   true,
			errMsg:    "longitude must be between -180 and 180 degrees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := NewLocation(tt.latitude, tt.longitude, tt.county, tt.region)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewLocation() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewLocation() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewLocation() unexpected error = %v", err)
				return
			}

			if loc.Latitude != tt.latitude {
				t.Errorf("NewLocation() latitude = %v, want %v", loc.Latitude, tt.latitude)
			}
			if loc.Longitude != tt.longitude {
				t.Errorf("NewLocation() longitude = %v, want %v", loc.Longitude, tt.longitude)
			}
			if loc.County != tt.county {
				t.Errorf("NewLocation() county = %v, want %v", loc.County, tt.county)
			}
			if loc.Region != tt.region {
				t.Errorf("NewLocation() region = %v, want %v", loc.Region, tt.region)
			}
		})
	}
}

func TestLocation_ImmutabilityMethods(t *testing.T) {
	original, _ := NewLocation(40.7128, -74.0060, "New York", "NY")

	// Test WithCoordinates
	newLoc, err := original.WithCoordinates(41.8781, -87.6298)
	if err != nil {
		t.Errorf("WithCoordinates() unexpected error = %v", err)
	}
	if newLoc.Latitude != 41.8781 {
		t.Errorf("WithCoordinates() latitude = %v, want %v", newLoc.Latitude, 41.8781)
	}
	if newLoc.Longitude != -87.6298 {
		t.Errorf("WithCoordinates() longitude = %v, want %v", newLoc.Longitude, -87.6298)
	}
	if original.Latitude != 40.7128 {
		t.Errorf("WithCoordinates() modified original location")
	}

	// Test WithCounty
	newLoc, err = original.WithCounty("Cook")
	if err != nil {
		t.Errorf("WithCounty() unexpected error = %v", err)
	}
	if newLoc.County != "Cook" {
		t.Errorf("WithCounty() county = %v, want %v", newLoc.County, "Cook")
	}
	if original.County != "New York" {
		t.Errorf("WithCounty() modified original location")
	}

	// Test WithRegion
	newLoc, err = original.WithRegion("IL")
	if err != nil {
		t.Errorf("WithRegion() unexpected error = %v", err)
	}
	if newLoc.Region != "IL" {
		t.Errorf("WithRegion() region = %v, want %v", newLoc.Region, "IL")
	}
	if original.Region != "NY" {
		t.Errorf("WithRegion() modified original location")
	}
}

func TestLocation_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		expected bool
	}{
		{
			name:     "empty location",
			location: Location{},
			expected: true,
		},
		{
			name: "location with coordinates only",
			location: Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expected: false,
		},
		{
			name: "location with county",
			location: Location{
				County: "New York",
			},
			expected: false,
		},
		{
			name: "location with region",
			location: Location{
				Region: "NY",
			},
			expected: false,
		},
		{
			name: "location with coordinates and county",
			location: Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
				County:    "New York",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.location.IsEmpty()
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLocation_DistanceTo(t *testing.T) {
	// New York City coordinates
	nyc, _ := NewLocation(40.7128, -74.0060, "New York", "NY")
	// Chicago coordinates
	chicago, _ := NewLocation(41.8781, -87.6298, "Cook", "IL")

	distance := nyc.DistanceTo(chicago)

	// Expected distance between NYC and Chicago is approximately 1145 km
	expectedDistance := 1145.0
	tolerance := 50.0 // Allow 50km tolerance

	if math.Abs(distance-expectedDistance) > tolerance {
		t.Errorf("DistanceTo() = %v, want approximately %v (±%v)", distance, expectedDistance, tolerance)
	}

	// Test distance to self should be 0
	selfDistance := nyc.DistanceTo(nyc)
	if selfDistance != 0 {
		t.Errorf("DistanceTo() self distance = %v, want 0", selfDistance)
	}
}

func TestLocation_String(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		expected string
	}{
		{
			name:     "empty location",
			location: Location{},
			expected: "",
		},
		{
			name: "location with coordinates only",
			location: Location{
				Latitude:  40.712800,
				Longitude: -74.006000,
			},
			expected: "40.712800, -74.006000",
		},
		{
			name: "full location",
			location: Location{
				Latitude:  40.712800,
				Longitude: -74.006000,
				County:    "New York",
				Region:    "NY",
			},
			expected: "40.712800, -74.006000 - New York County - NY",
		},
		{
			name: "location with county only",
			location: Location{
				Latitude:  40.712800,
				Longitude: -74.006000,
				County:    "New York",
			},
			expected: "40.712800, -74.006000 - New York County",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.location.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	loc1, _ := NewLocation(40.7128, -74.0060, "New York", "NY")
	loc2, _ := NewLocation(40.7128, -74.0060, "New York", "NY")
	loc3, _ := NewLocation(41.8781, -87.6298, "Cook", "IL")
	loc4, _ := NewLocation(40.7128, -74.0060, "Manhattan", "NY")

	tests := []struct {
		name     string
		loc1     Location
		loc2     Location
		expected bool
	}{
		{
			name:     "identical locations",
			loc1:     loc1,
			loc2:     loc2,
			expected: true,
		},
		{
			name:     "different coordinates",
			loc1:     loc1,
			loc2:     loc3,
			expected: false,
		},
		{
			name:     "different county",
			loc1:     loc1,
			loc2:     loc4,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.loc1.Equals(tt.loc2)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLocation_WithCoordinates_ValidationError(t *testing.T) {
	original, _ := NewLocation(40.7128, -74.0060, "New York", "NY")

	// Test invalid latitude
	_, err := original.WithCoordinates(91, -74.0060)
	if err == nil {
		t.Errorf("WithCoordinates() expected error for invalid latitude")
	}

	// Test invalid longitude
	_, err = original.WithCoordinates(40.7128, 181)
	if err == nil {
		t.Errorf("WithCoordinates() expected error for invalid longitude")
	}
}
