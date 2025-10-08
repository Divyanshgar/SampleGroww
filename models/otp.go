package models

import (
	"time"
)

// OTP represents a one-time password for email verification
type OTP struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	Code       string    `json:"code"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// IsExpired checks if the OTP has expired
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}
