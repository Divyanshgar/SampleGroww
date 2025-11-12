package utils

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"notification-server/database"
	"notification-server/db"
	"time"
)

// GenerateOTP creates a 6-digit random OTP
func GenerateOTP() (string, error) {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp, nil
}

// SaveOTP stores the OTP in the database with expiration (minutes)
func SaveOTP(email, otp string, expireMinutes int) error {
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)

	// Create user with OTP fields
	params := db.CreateUserParams{
		FirstName: "",
		LastName:  "",
		Email:     email,
		Phone:     sql.NullString{Valid: false},
		Address:   sql.NullString{Valid: false},
		City:      sql.NullString{Valid: false},
		Country:   sql.NullString{Valid: false},
	}

	user, err := database.GetQueries().CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	// Update with OTP fields
	updateParams := db.UpdateUserOTPParams{
		ID:          user.ID,
		OtpCode:     sql.NullString{String: otp, Valid: true},
		OtpExpires:  sql.NullTime{Time: expiresAt, Valid: true},
		OtpVerified: sql.NullBool{Bool: false, Valid: true},
	}

	_, err = database.GetQueries().UpdateUserOTP(context.Background(), updateParams)
	return err
}

// SaveOrUpdateOTP stores or updates the OTP in the database with expiration (minutes)
func SaveOrUpdateOTP(email, otp string, expireMinutes int) error {
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)

	// Try to find existing user
	existingUser, err := database.GetQueries().GetUserByEmail(context.Background(), email)
	if err != nil {
		// User doesn't exist, create new one
		params := db.CreateUserParams{
			FirstName: "",
			LastName:  "",
			Email:     email,
			Phone:     sql.NullString{Valid: false},
			Address:   sql.NullString{Valid: false},
			City:      sql.NullString{Valid: false},
			Country:   sql.NullString{Valid: false},
		}
		_, err = database.GetQueries().CreateUser(context.Background(), params)
		if err != nil {
			return err
		}
		// Get the newly created user to update OTP
		existingUser, err = database.GetQueries().GetUserByEmail(context.Background(), email)
		if err != nil {
			return err
		}
	}

	// Update OTP fields
	updateParams := db.UpdateUserOTPParams{
		ID:          existingUser.ID,
		OtpCode:     sql.NullString{String: otp, Valid: true},
		OtpExpires:  sql.NullTime{Time: expiresAt, Valid: true},
		OtpVerified: sql.NullBool{Bool: false, Valid: true},
	}

	_, err = database.GetQueries().UpdateUserOTP(context.Background(), updateParams)
	return err
}

// CheckOTP checks the OTP for a given email without marking as verified
func CheckOTP(email, code string) bool {
	user, err := database.GetQueries().GetUserByEmailAndOTP(context.Background(), db.GetUserByEmailAndOTPParams{
		Email:   email,
		OtpCode: sql.NullString{String: code, Valid: true},
	})

	if err != nil {
		// OTP not found
		return false
	}

	if user.OtpVerified.Valid && user.OtpVerified.Bool {
		// Already verified
		return false
	}

	if user.OtpExpires.Valid && time.Now().After(user.OtpExpires.Time) {
		// Expired OTP
		return false
	}

	return true
}

// VerifyAndMarkOTP checks the OTP and marks it as verified if valid
func VerifyAndMarkOTP(email, code string) bool {
	// Get user by email first
	user, err := database.GetQueries().GetUserByEmail(context.Background(), email)
	if err != nil {
		return false
	}

	// Verify OTP using the VerifyUserOTP query
	_, err = database.GetQueries().VerifyUserOTP(context.Background(), db.VerifyUserOTPParams{
		ID:      user.ID,
		OtpCode: sql.NullString{String: code, Valid: true},
	})

	return err == nil
}

// DeleteOTP clears OTP fields for the user (cleanup)
func DeleteOTP(email string) error {
	user, err := database.GetQueries().GetUserByEmail(context.Background(), email)
	if err != nil {
		return err
	}

	// Clear OTP fields
	updateParams := db.UpdateUserOTPParams{
		ID:          user.ID,
		OtpCode:     sql.NullString{Valid: false},
		OtpExpires:  sql.NullTime{Valid: false},
		OtpVerified: sql.NullBool{Bool: true, Valid: true},
	}

	_, err = database.GetQueries().UpdateUserOTP(context.Background(), updateParams)
	return err
}
