package utils

import (
	"fmt"
	"math/rand"
	"notification-server/database"
	"notification-server/models"
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
	user := &models.User{
		Email:       email,
		OTPCode:     otp,
		OTPExpires:  expiresAt,
		OTPVerified: false,
	}

	return database.GetDB().Create(user).Error
}

// SaveOrUpdateOTP stores or updates the OTP in the database with expiration (minutes)
func SaveOrUpdateOTP(email, otp string, expireMinutes int) error {
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)

	// Try to update existing user
	result := database.GetDB().Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"otp_code":     otp,
		"otp_expires":  expiresAt,
		"otp_verified": false,
	})

	if result.Error != nil {
		return result.Error
	}

	// If no rows updated, create new user
	if result.RowsAffected == 0 {
		user := &models.User{
			Email:       email,
			OTPCode:     otp,
			OTPExpires:  expiresAt,
			OTPVerified: false,
		}
		return database.GetDB().Create(user).Error
	}

	return nil
}

// CheckOTP checks the OTP for a given email without marking as verified
func CheckOTP(email, code string) bool {
	var user models.User
	err := database.GetDB().
		Where("email = ? AND otp_code = ?", email, code).
		Order("created_at desc").
		First(&user).Error

	if err != nil {
		// OTP not found
		return false
	}

	if user.OTPVerified {
		// Already verified
		return false
	}

	if time.Now().After(user.OTPExpires) {
		// Expired OTP
		return false
	}

	return true
}

// VerifyAndMarkOTP checks the OTP and marks it as verified if valid
func VerifyAndMarkOTP(email, code string) bool {
	if !CheckOTP(email, code) {
		return false
	}

	// Mark OTP as verified
	database.GetDB().Model(&models.User{}).Where("email = ?", email).Update("otp_verified", true)
	return true
}

// DeleteOTP clears OTP fields for the user (cleanup)
func DeleteOTP(email string) error {
	return database.GetDB().Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"otp_code":     "",
		"otp_expires":  time.Time{},
		"otp_verified": true,
	}).Error
}
