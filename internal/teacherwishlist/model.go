package teacherwishlist

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// TeacherStatus represents the possible states of a teacher registration
type TeacherStatus string

const (
	TeacherStatusPending  TeacherStatus = "pending"
	TeacherStatusApproved TeacherStatus = "approved"
	TeacherStatusRejected TeacherStatus = "rejected"
)

// String returns the string representation of TeacherStatus
func (ts TeacherStatus) String() string {
	return string(ts)
}

// IsValid checks if the teacher status is valid
func (ts TeacherStatus) IsValid() bool {
	switch ts {
	case TeacherStatusPending, TeacherStatusApproved, TeacherStatusRejected:
		return true
	default:
		return false
	}
}

// Teacher represents a teacher aggregate root in the domain
// This is an Aggregate Root that encapsulates teacher registration, profile management, and wishlist functionality
type Teacher struct {
	ID          string         `json:"id" db:"id"`
	Email       string         `json:"email" db:"email"`
	FirstName   string         `json:"first_name" db:"first_name"`
	LastName    string         `json:"last_name" db:"last_name"`
	SchoolID    string         `json:"school_id" db:"school_id"`
	GradeLevel  string         `json:"grade_level" db:"grade_level"`
	WishlistURL string         `json:"wishlist_url" db:"wishlist_url"`
	Status      TeacherStatus  `json:"status" db:"status"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// emailRegex for basic email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// generateID creates a simple UUID-like string using crypto/rand
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// NewTeacher creates a new Teacher with validation
func NewTeacher(email, firstName, lastName, schoolID, gradeLevel string) (*Teacher, error) {
	id := generateID()
	now := time.Now()

	teacher := &Teacher{
		ID:         id,
		Email:      strings.TrimSpace(strings.ToLower(email)),
		FirstName:  strings.TrimSpace(firstName),
		LastName:   strings.TrimSpace(lastName),
		SchoolID:   strings.TrimSpace(schoolID),
		GradeLevel: strings.TrimSpace(gradeLevel),
		Status:     TeacherStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := teacher.validateInvariants(); err != nil {
		return nil, err
	}

	return teacher, nil
}

// validateInvariants checks all aggregate invariants
func (t *Teacher) validateInvariants() error {
	// Email validation
	if t.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(t.Email) {
		return errors.New("email format is invalid")
	}

	// Name validation
	if t.FirstName == "" {
		return errors.New("first name is required")
	}
	if t.LastName == "" {
		return errors.New("last name is required")
	}

	// School association validation
	if t.SchoolID == "" {
		return errors.New("school ID is required")
	}

	// Status validation
	if !t.Status.IsValid() {
		return errors.New("invalid teacher status")
	}

	// Validate wishlist URL if provided
	if t.WishlistURL != "" {
		if err := t.validateWishlistURL(t.WishlistURL); err != nil {
			return err
		}
	}

	return nil
}

// validateWishlistURL validates Amazon wishlist URL format
func (t *Teacher) validateWishlistURL(wishlistURL string) error {
	if wishlistURL == "" {
		return nil // Empty URL is valid for non-approved teachers
	}

	// Parse URL
	parsedURL, err := url.Parse(wishlistURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Check if URL has a valid scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid URL format")
	}

	// Check if it's HTTPS
	if parsedURL.Scheme != "https" {
		return errors.New("wishlist URL must use HTTPS")
	}

	// Check if it's an Amazon domain
	if !strings.Contains(parsedURL.Host, "amazon.com") && !strings.Contains(parsedURL.Host, "amzn.to") {
		return errors.New("wishlist URL must be from Amazon")
	}

	return nil
}

// UpdateWishlistURL updates the teacher's wishlist URL with validation
func (t *Teacher) UpdateWishlistURL(wishlistURL string) error {
	wishlistURL = strings.TrimSpace(wishlistURL)

	if err := t.validateWishlistURL(wishlistURL); err != nil {
		return err
	}

	t.WishlistURL = wishlistURL
	t.UpdatedAt = time.Now()

	// Re-validate invariants after update
	return t.validateInvariants()
}

// Approve transitions the teacher status to approved
func (t *Teacher) Approve() error {
	if t.Status == TeacherStatusApproved {
		return errors.New("teacher is already approved")
	}

	if t.Status == TeacherStatusRejected {
		return errors.New("cannot approve a rejected teacher")
	}

	// Must have wishlist URL to be approved
	if t.WishlistURL == "" {
		return errors.New("teacher must have a wishlist URL to be approved")
	}

	t.Status = TeacherStatusApproved
	t.UpdatedAt = time.Now()

	return t.validateInvariants()
}

// Reject transitions the teacher status to rejected
func (t *Teacher) Reject() error {
	if t.Status == TeacherStatusRejected {
		return errors.New("teacher is already rejected")
	}

	if t.Status == TeacherStatusApproved {
		return errors.New("cannot reject an approved teacher")
	}

	t.Status = TeacherStatusRejected
	t.UpdatedAt = time.Now()

	return t.validateInvariants()
}

// UpdateProfile updates the teacher's profile information
func (t *Teacher) UpdateProfile(firstName, lastName, gradeLevel string) error {
	t.FirstName = strings.TrimSpace(firstName)
	t.LastName = strings.TrimSpace(lastName)
	t.GradeLevel = strings.TrimSpace(gradeLevel)
	t.UpdatedAt = time.Now()

	return t.validateInvariants()
}

// IsApproved returns true if the teacher is approved
func (t *Teacher) IsApproved() bool {
	return t.Status == TeacherStatusApproved
}

// IsPending returns true if the teacher is pending approval
func (t *Teacher) IsPending() bool {
	return t.Status == TeacherStatusPending
}

// IsRejected returns true if the teacher is rejected
func (t *Teacher) IsRejected() bool {
	return t.Status == TeacherStatusRejected
}

// FullName returns the teacher's full name
func (t *Teacher) FullName() string {
	return strings.TrimSpace(t.FirstName + " " + t.LastName)
}

// HasWishlist returns true if the teacher has a wishlist URL
func (t *Teacher) HasWishlist() bool {
	return t.WishlistURL != ""
}
