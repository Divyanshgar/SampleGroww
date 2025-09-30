package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"notification-server/models"
	"sync"
	"time"
)

// OTPStore handles in-memory storage of OTPs
type OTPStore struct {
	otps map[string]*models.OTP
	mu   sync.RWMutex
}

var store *OTPStore
var once sync.Once

// GetOTPStore returns a singleton instance of OTPStore
func GetOTPStore() *OTPStore {
	once.Do(func() {
		store = &OTPStore{
			otps: make(map[string]*models.OTP),
		}
		// Start cleanup goroutine
		go store.cleanupExpiredOTPs()
	})
	return store
}

// GenerateOTP generates a 6-digit OTP code
func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Save stores an OTP for a given email
func (s *OTPStore) Save(email, code string, expiryMinutes int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.otps[email] = &models.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Duration(expiryMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}
}

// Verify checks if the OTP is valid for the given email
func (s *OTPStore) Verify(email, code string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	otp, exists := s.otps[email]
	if !exists {
		return false
	}

	if otp.IsExpired() {
		return false
	}

	return otp.Code == code
}

// Delete removes an OTP for a given email
func (s *OTPStore) Delete(email string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.otps, email)
}

// cleanupExpiredOTPs periodically removes expired OTPs
func (s *OTPStore) cleanupExpiredOTPs() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for email, otp := range s.otps {
			if otp.IsExpired() {
				delete(s.otps, email)
			}
		}
		s.mu.Unlock()
	}
}
