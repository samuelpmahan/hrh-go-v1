package teacherwishlist

import (
	"testing"
	"time"
)

func TestNewTeacher(t *testing.T) {

	tests := []struct {
		name        string
		email       string
		firstName   string
		lastName    string
		schoolID    string
		gradeLevel  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid teacher creation",
			email:      "teacher@example.com",
			firstName:  "John",
			lastName:   "Doe",
			schoolID:   "school-123",
			gradeLevel: "5th Grade",
			wantErr:    false,
		},
		{
			name:        "empty email",
			email:       "",
			firstName:   "John",
			lastName:    "Doe",
			schoolID:    "school-123",
			gradeLevel:  "5th Grade",
			wantErr:     true,
			errContains: "email is required",
		},
		{
			name:        "invalid email format",
			email:       "invalid-email",
			firstName:   "John",
			lastName:    "Doe",
			schoolID:    "school-123",
			gradeLevel:  "5th Grade",
			wantErr:     true,
			errContains: "email format is invalid",
		},
		{
			name:        "empty first name",
			email:       "teacher@example.com",
			firstName:   "",
			lastName:    "Doe",
			schoolID:    "school-123",
			gradeLevel:  "5th Grade",
			wantErr:     true,
			errContains: "first name is required",
		},
		{
			name:        "empty last name",
			email:       "teacher@example.com",
			firstName:   "John",
			lastName:    "",
			schoolID:    "school-123",
			gradeLevel:  "5th Grade",
			wantErr:     true,
			errContains: "last name is required",
		},
		{
			name:        "empty school ID",
			email:       "teacher@example.com",
			firstName:   "John",
			lastName:    "Doe",
			schoolID:    "",
			gradeLevel:  "5th Grade",
			wantErr:     true,
			errContains: "school ID is required",
		},
		{
			name:       "email normalization",
			email:      "  TEACHER@EXAMPLE.COM  ",
			firstName:  "John",
			lastName:   "Doe",
			schoolID:   "school-123",
			gradeLevel: "5th Grade",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teacher, err := NewTeacher(tt.email, tt.firstName, tt.lastName, tt.schoolID, tt.gradeLevel)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTeacher() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("NewTeacher() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTeacher() unexpected error = %v", err)
				return
			}

			// Verify teacher properties
			if teacher.ID == "" {
				t.Error("NewTeacher() should generate an ID")
			}
			if teacher.Status != TeacherStatusPending {
				t.Errorf("NewTeacher() status = %v, want %v", teacher.Status, TeacherStatusPending)
			}
			if teacher.CreatedAt.IsZero() {
				t.Error("NewTeacher() should set CreatedAt")
			}
			if teacher.UpdatedAt.IsZero() {
				t.Error("NewTeacher() should set UpdatedAt")
			}

			// Test email normalization
			if tt.name == "email normalization" {
				if teacher.Email != "teacher@example.com" {
					t.Errorf("NewTeacher() email = %v, want %v", teacher.Email, "teacher@example.com")
				}
			}
		})
	}
}

func TestTeacher_UpdateWishlistURL(t *testing.T) {
	teacher, _ := NewTeacher("teacher@example.com", "John", "Doe", "school-123", "5th Grade")

	tests := []struct {
		name        string
		wishlistURL string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid Amazon URL",
			wishlistURL: "https://www.amazon.com/hz/wishlist/ls/ABC123",
			wantErr:     false,
		},
		{
			name:        "valid Amazon short URL",
			wishlistURL: "https://amzn.to/abc123",
			wantErr:     false,
		},
		{
			name:        "empty URL",
			wishlistURL: "",
			wantErr:     false,
		},
		{
			name:        "non-Amazon URL",
			wishlistURL: "https://www.example.com/wishlist",
			wantErr:     true,
			errContains: "wishlist URL must be from Amazon",
		},
		{
			name:        "HTTP instead of HTTPS",
			wishlistURL: "http://www.amazon.com/hz/wishlist/ls/ABC123",
			wantErr:     true,
			errContains: "wishlist URL must use HTTPS",
		},
		{
			name:        "invalid URL format",
			wishlistURL: "not-a-url",
			wantErr:     true,
			errContains: "invalid URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := teacher.UpdatedAt
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			err := teacher.UpdateWishlistURL(tt.wishlistURL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateWishlistURL() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("UpdateWishlistURL() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateWishlistURL() unexpected error = %v", err)
				return
			}

			if teacher.WishlistURL != tt.wishlistURL {
				t.Errorf("UpdateWishlistURL() wishlistURL = %v, want %v", teacher.WishlistURL, tt.wishlistURL)
			}

			if !teacher.UpdatedAt.After(originalUpdatedAt) {
				t.Error("UpdateWishlistURL() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestTeacher_Approve(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  TeacherStatus
		hasWishlistURL bool
		wantErr        bool
		errContains    string
	}{
		{
			name:           "approve pending teacher with wishlist",
			initialStatus:  TeacherStatusPending,
			hasWishlistURL: true,
			wantErr:        false,
		},
		{
			name:           "approve pending teacher without wishlist",
			initialStatus:  TeacherStatusPending,
			hasWishlistURL: false,
			wantErr:        true,
			errContains:    "teacher must have a wishlist URL to be approved",
		},
		{
			name:           "approve already approved teacher",
			initialStatus:  TeacherStatusApproved,
			hasWishlistURL: true,
			wantErr:        true,
			errContains:    "teacher is already approved",
		},
		{
			name:           "approve rejected teacher",
			initialStatus:  TeacherStatusRejected,
			hasWishlistURL: true,
			wantErr:        true,
			errContains:    "cannot approve a rejected teacher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teacher, _ := NewTeacher("teacher@example.com", "John", "Doe", "school-123", "5th Grade")
			teacher.Status = tt.initialStatus

			if tt.hasWishlistURL {
				teacher.WishlistURL = "https://www.amazon.com/hz/wishlist/ls/ABC123"
			}

			originalUpdatedAt := teacher.UpdatedAt
			time.Sleep(1 * time.Millisecond)

			err := teacher.Approve()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Approve() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("Approve() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Approve() unexpected error = %v", err)
				return
			}

			if teacher.Status != TeacherStatusApproved {
				t.Errorf("Approve() status = %v, want %v", teacher.Status, TeacherStatusApproved)
			}

			if !teacher.UpdatedAt.After(originalUpdatedAt) {
				t.Error("Approve() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestTeacher_Reject(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus TeacherStatus
		wantErr       bool
		errContains   string
	}{
		{
			name:          "reject pending teacher",
			initialStatus: TeacherStatusPending,
			wantErr:       false,
		},
		{
			name:          "reject already rejected teacher",
			initialStatus: TeacherStatusRejected,
			wantErr:       true,
			errContains:   "teacher is already rejected",
		},
		{
			name:          "reject approved teacher",
			initialStatus: TeacherStatusApproved,
			wantErr:       true,
			errContains:   "cannot reject an approved teacher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teacher, _ := NewTeacher("teacher@example.com", "John", "Doe", "school-123", "5th Grade")
			teacher.Status = tt.initialStatus

			if tt.initialStatus == TeacherStatusApproved {
				teacher.WishlistURL = "https://www.amazon.com/hz/wishlist/ls/ABC123"
			}

			originalUpdatedAt := teacher.UpdatedAt
			time.Sleep(1 * time.Millisecond)

			err := teacher.Reject()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Reject() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("Reject() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Reject() unexpected error = %v", err)
				return
			}

			if teacher.Status != TeacherStatusRejected {
				t.Errorf("Reject() status = %v, want %v", teacher.Status, TeacherStatusRejected)
			}

			if !teacher.UpdatedAt.After(originalUpdatedAt) {
				t.Error("Reject() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestTeacher_UpdateProfile(t *testing.T) {
	teacher, _ := NewTeacher("teacher@example.com", "John", "Doe", "school-123", "5th Grade")

	originalUpdatedAt := teacher.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	err := teacher.UpdateProfile("Jane", "Smith", "3rd Grade")

	if err != nil {
		t.Errorf("UpdateProfile() unexpected error = %v", err)
		return
	}

	if teacher.FirstName != "Jane" {
		t.Errorf("UpdateProfile() firstName = %v, want %v", teacher.FirstName, "Jane")
	}
	if teacher.LastName != "Smith" {
		t.Errorf("UpdateProfile() lastName = %v, want %v", teacher.LastName, "Smith")
	}
	if teacher.GradeLevel != "3rd Grade" {
		t.Errorf("UpdateProfile() gradeLevel = %v, want %v", teacher.GradeLevel, "3rd Grade")
	}
	if !teacher.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdateProfile() should update UpdatedAt timestamp")
	}
}

func TestTeacher_StatusMethods(t *testing.T) {
	teacher, _ := NewTeacher("teacher@example.com", "John", "Doe", "school-123", "5th Grade")

	// Test pending status
	if !teacher.IsPending() {
		t.Error("IsPending() should return true for new teacher")
	}
	if teacher.IsApproved() {
		t.Error("IsApproved() should return false for new teacher")
	}
	if teacher.IsRejected() {
		t.Error("IsRejected() should return false for new teacher")
	}

	// Test approved status
	teacher.WishlistURL = "https://www.amazon.com/hz/wishlist/ls/ABC123"
	teacher.Approve()

	if teacher.IsPending() {
		t.Error("IsPending() should return false for approved teacher")
	}
	if !teacher.IsApproved() {
		t.Error("IsApproved() should return true for approved teacher")
	}
	if teacher.IsRejected() {
		t.Error("IsRejected() should return false for approved teacher")
	}

	// Test rejected status (create new teacher since we can't reject approved)
	rejectedTeacher, _ := NewTeacher("rejected@example.com", "Jane", "Smith", "school-456", "2nd Grade")
	rejectedTeacher.Reject()

	if rejectedTeacher.IsPending() {
		t.Error("IsPending() should return false for rejected teacher")
	}
	if rejectedTeacher.IsApproved() {
		t.Error("IsApproved() should return false for rejected teacher")
	}
	if !rejectedTeacher.IsRejected() {
		t.Error("IsRejected() should return true for rejected teacher")
	}
}

func TestTeacherStatus_Methods(t *testing.T) {
	tests := []struct {
		status   TeacherStatus
		isValid  bool
		expected string
	}{
		{TeacherStatusPending, true, "pending"},
		{TeacherStatusApproved, true, "approved"},
		{TeacherStatusRejected, true, "rejected"},
		{TeacherStatus("invalid"), false, "invalid"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if tt.status.IsValid() != tt.isValid {
				t.Errorf("IsValid() = %v, want %v", tt.status.IsValid(), tt.isValid)
			}
			if tt.status.String() != tt.expected {
				t.Errorf("String() = %v, want %v", tt.status.String(), tt.expected)
			}
		})
	}
}

func TestTeacher_InvariantValidation(t *testing.T) {
	teacher, _ := NewTeacher("teacher@example.com", "John", "Doe", "school-123", "5th Grade")

	// Test approved teacher without wishlist URL invariant
	teacher.Status = TeacherStatusApproved
	teacher.WishlistURL = ""

	err := teacher.validateInvariants()
	if err == nil {
		t.Error("validateInvariants() should fail for approved teacher without wishlist URL")
	}
	if !containsString(err.Error(), "approved teachers must have a wishlist URL") {
		t.Errorf("validateInvariants() error = %v, want error containing 'approved teachers must have a wishlist URL'", err)
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
