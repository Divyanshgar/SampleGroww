# TODO: Replace GORM with SQLC in Notification Server Project

## Overview
Replace GORM ORM with SQLC for type-safe database operations across the entire project.

## Steps to Complete

### 1. Update Dependencies
- [ ] Remove GORM dependencies from go.mod (gorm.io/gorm and gorm.io/driver/postgres)
- [ ] Keep SQLC and other dependencies
- [ ] Run `go mod tidy` to clean up

### 2. Update Models
- [ ] Remove GORM tags from models/user.go and models/stock.go
- [ ] Create conversion functions between SQLC models (with sql.NullString) and application models
- [ ] Update models to work with SQLC types

### 3. Update Database Initialization
- [ ] Remove GORM setup from database/database.go
- [ ] Keep only SQLC initialization (Queries)
- [ ] Remove GormDB variable and GetDB function
- [ ] Update AutoMigrate logic if needed (or remove if not using GORM migrations)

### 4. Replace Database Operations in Utils
- [ ] Update utils/otp_store.go to use SQLC queries instead of GORM
- [ ] Convert SaveOTP, SaveOrUpdateOTP, CheckOTP, VerifyAndMarkOTP, DeleteOTP functions

### 5. Replace Database Operations in Services
- [ ] Update services/excel_service.go to use SQLC for fetching stocks

### 6. Replace Database Operations in Handlers
- [ ] Update handlers/notification_handler.go for all user and stock operations
- [ ] Convert CRUD operations to SQLC queries

### 7. Update Example and Seeder
- [ ] Update example/simple_excel_example.go to use SQLC
- [ ] Update database/seeder.go to use SQLC

### 8. Handle Type Conversions
- [ ] Add helper functions for converting between SQLC models and app models
- [ ] Handle sql.NullString, sql.NullTime, etc. properly

### 9. Regenerate SQLC Code
- [ ] Run `sqlc generate` to ensure all queries are up to date

### 10. Test Thoroughly
- [ ] Build the project and check for compilation errors
- [ ] Run the application and test all endpoints
- [ ] Test user registration, OTP verification, profile generation
- [ ] Verify database operations (create, read, update, delete) work correctly
- [ ] Test Excel and PDF generation
- [ ] Check frontend integration if applicable

## Notes
- SQLC models use sql.NullString for nullable fields, need to handle conversions
- Ensure all existing functionality is preserved
- Pay special attention to OTP operations as they involve complex queries
- Test with real database to ensure no data corruption
