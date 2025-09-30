package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FirstName string         `gorm:"size:100;not null" json:"first_name"`
	LastName  string         `gorm:"size:100;not null" json:"last_name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"size:500" json:"address"`
	City      string         `gorm:"size:100" json:"city"`
	Country   string         `gorm:"size:100" json:"country"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}
