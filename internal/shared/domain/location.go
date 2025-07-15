package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Location represents geographic coordinates and standardized address components
// This is a Value Object - immutable and defined by its attributes
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	County    string  `json:"county"`
	Region    string  `json:"region"`
}

// NewLocation creates a new Location with validation
func NewLocation(latitude, longitude float64, county, region string) (Location, error) {
	loc := Location{
		Latitude:  latitude,
		Longitude: longitude,
		County:    strings.TrimSpace(county),
		Region:    strings.TrimSpace(region),
	}

	if err := loc.validate(); err != nil {
		return Location{}, err
	}

	return loc, nil
}

// WithCoordinates returns a new Location with updated coordinates
func (l Location) WithCoordinates(latitude, longitude float64) (Location, error) {
	return NewLocation(latitude, longitude, l.County, l.Region)
}

// WithCounty returns a new Location with updated county
func (l Location) WithCounty(county string) (Location, error) {
	return NewLocation(l.Latitude, l.Longitude, county, l.Region)
}

// WithRegion returns a new Location with updated region
func (l Location) WithRegion(region string) (Location, error) {
	return NewLocation(l.Latitude, l.Longitude, l.County, region)
}

// validate performs validation on location components
func (l Location) validate() error {
	// Validate latitude range (-90 to 90)
	if l.Latitude < -90 || l.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90 degrees")
	}

	// Validate longitude range (-180 to 180)
	if l.Longitude < -180 || l.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180 degrees")
	}

	return nil
}

// IsEmpty returns true if the location has no meaningful coordinate data
func (l Location) IsEmpty() bool {
	return l.Latitude == 0 && l.Longitude == 0
}

// DistanceTo calculates the distance in kilometers to another location using Haversine formula
func (l Location) DistanceTo(other Location) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := l.Latitude * math.Pi / 180
	lat2Rad := other.Latitude * math.Pi / 180
	deltaLatRad := (other.Latitude - l.Latitude) * math.Pi / 180
	deltaLngRad := (other.Longitude - l.Longitude) * math.Pi / 180

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLngRad/2)*math.Sin(deltaLngRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// String returns a formatted string representation of the location
func (l Location) String() string {
	parts := []string{}

	if !l.IsEmpty() {
		parts = append(parts, fmt.Sprintf("%.6f, %.6f", l.Latitude, l.Longitude))
	}

	if l.County != "" {
		parts = append(parts, l.County+" County")
	}

	if l.Region != "" {
		parts = append(parts, l.Region)
	}

	return strings.Join(parts, " - ")
}

// Equals compares two locations for equality (with small tolerance for floating point comparison)
func (l Location) Equals(other Location) bool {
	const tolerance = 1e-6

	return math.Abs(l.Latitude-other.Latitude) < tolerance &&
		math.Abs(l.Longitude-other.Longitude) < tolerance &&
		l.County == other.County &&
		l.Region == other.Region
}
