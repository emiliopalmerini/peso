package persistence

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

// DB wraps sql.DB with additional functionality
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Migrate runs database migrations from the given filesystem
func (db *DB) Migrate(migrationsFS fs.FS) error {
	// Create migrations table if it doesn't exist
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := fs.Glob(migrationsFS, "*.sql")
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	// Sort files to ensure proper execution order
	sort.Strings(files)

	for _, file := range files {
		migrationName := strings.TrimSuffix(filepath.Base(file), ".sql")
		
		// Check if migration has already been applied
		applied, err := db.isMigrationApplied(migrationName)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}
		
		if applied {
			continue
		}

		// Read and execute migration
		content, err := fs.ReadFile(migrationsFS, file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if err := db.executeMigration(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		// Mark migration as applied
		if err := db.markMigrationApplied(migrationName); err != nil {
			return fmt.Errorf("failed to mark migration as applied: %w", err)
		}

		fmt.Printf("Applied migration: %s\n", migrationName)
	}

	return nil
}

func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func (db *DB) isMigrationApplied(name string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) executeMigration(content string) error {
	// Execute the entire content as one statement for SQLite
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clean up the content - remove comments and excessive whitespace
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}
	
	if len(cleanedLines) == 0 {
		return nil // No actual SQL to execute
	}
	
	cleanedContent := strings.Join(cleanedLines, " ")
	
	if _, err := tx.Exec(cleanedContent); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return tx.Commit()
}

func (db *DB) markMigrationApplied(name string) error {
	_, err := db.Exec("INSERT INTO migrations (name) VALUES (?)", name)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}