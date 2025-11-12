package database

import (
	"database/sql"
	"fmt"
	"log"
	"notification-server/config"
	"notification-server/db"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *db.Queries

// Initialize sets up the database connection and runs migrations
func Initialize(cfg *config.Config) error {
	// Step 1: Prepare DSN string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Step 2: Connect to DB
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("✅ Database connection established successfully")

	// Step 3: Run migrations (execute schema.sql)
	if err = runMigrations(DB); err != nil {
		log.Printf("Migration failed (possibly tables already exist): %v", err)
		log.Println("Continuing with existing database schema...")
	} else {
		log.Println("🚀 Migrations completed successfully")
	}

	// Step 4: Initialize SQLC Queries
	Queries = db.New(DB)
	log.Println("🔧 SQLC Queries initialized successfully")

	return nil
}

// runMigrations executes the schema.sql file to create tables
func runMigrations(db *sql.DB) error {
	schemaPath := filepath.Join("database", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// GetDB returns the SQL database instance (for backward compatibility)
func GetDB() *sql.DB {
	return DB
}

// GetQueries returns the SQLC Queries instance
func GetQueries() *db.Queries {
	return Queries
}
