package database

import (
	"fmt"
	"log"
	"notification-server/config"
	"notification-server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize sets up the database connection
func Initialize(cfg *config.Config) error {
	// Step 1: Prepare DSN string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Step 2: Connect to DB
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("✅ Database connection established successfully")

	// ✅ Step 3: AutoMigrate all your models here (this creates tables)
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Stock{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %w", err)
	}
	log.Println("🚀 AutoMigration completed successfully")

	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
