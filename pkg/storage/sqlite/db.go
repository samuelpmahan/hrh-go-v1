package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DBInterface defines the interface that our repositories expect
type DBInterface interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// DB wraps the sql.DB with additional functionality
type DB struct {
	*sql.DB
}

// Ensure DB implements DBInterface
var _ DBInterface = (*DB)(nil)

// Config holds database configuration
type Config struct {
	DatabasePath string
}

// NewDB creates a new SQLite database connection with proper configuration
func NewDB(config Config) (*DB, error) {
	// SQLite connection string with useful pragmas for performance and reliability
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=ON", config.DatabasePath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ensure foreign keys are enabled for this connection
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Configure connection pool for SQLite (single writer, multiple readers)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize schema if needed
	if err := initializeSchema(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &DB{DB: db}, nil
}

// initializeSchema creates the necessary tables if they don't exist
func initializeSchema(ctx context.Context, db *sql.DB) error {
	schema := `
	-- Schools table with Address and Location Value Objects flattened
	CREATE TABLE IF NOT EXISTS schools (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		-- Address Value Object fields
		address_street TEXT,
		address_city TEXT,
		address_state TEXT,
		address_zip_code TEXT,
		-- Location Value Object fields
		location_latitude REAL,
		location_longitude REAL,
		location_county TEXT,
		location_region TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		-- Unique constraint for school uniqueness invariant
		UNIQUE(name, address_city, address_state)
	);

	-- Teachers table
	CREATE TABLE IF NOT EXISTS teachers (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		school_id TEXT REFERENCES schools(id),
		grade_level TEXT,
		wishlist_url TEXT,
		status TEXT DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Admins table
	CREATE TABLE IF NOT EXISTS admins (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Indexes for performance and cross-aggregate invariant enforcement
	CREATE INDEX IF NOT EXISTS idx_teachers_school_id ON teachers(school_id);
	CREATE INDEX IF NOT EXISTS idx_teachers_status ON teachers(status);
	CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email);
	CREATE INDEX IF NOT EXISTS idx_schools_location ON schools(location_latitude, location_longitude);
	CREATE INDEX IF NOT EXISTS idx_schools_city_state ON schools(address_city, address_state);
	`

	_, err := db.ExecContext(ctx, schema)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// WithTx executes a function within a database transaction
func (db *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// Common error types for repository operations
var (
	ErrNotFound     = fmt.Errorf("record not found")
	ErrDuplicateKey = fmt.Errorf("duplicate key violation")
	ErrForeignKey   = fmt.Errorf("foreign key violation")
	ErrInvalidInput = fmt.Errorf("invalid input")
)

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound)
}

// IsDuplicateKeyError checks if an error is a duplicate key violation
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// SQLite constraint error messages
	return contains(err.Error(), "UNIQUE constraint failed") ||
		contains(err.Error(), "duplicate key")
}

// IsForeignKeyError checks if an error is a foreign key violation
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	// SQLite foreign key error messages
	return contains(err.Error(), "FOREIGN KEY constraint failed") ||
		contains(err.Error(), "foreign key")
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAt(s, substr, 1))))
}

func containsAt(s, substr string, start int) bool {
	if start >= len(s) {
		return false
	}
	if start+len(substr) > len(s) {
		return containsAt(s, substr, start+1)
	}
	if s[start:start+len(substr)] == substr {
		return true
	}
	return containsAt(s, substr, start+1)
}
