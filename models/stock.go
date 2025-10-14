package models

import (
	"time"

	"gorm.io/gorm"
)

// Stock represents a stock transaction for a user
type Stock struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	StockName     string         `gorm:"size:100;not null" json:"stock_name"`
	Action        string         `gorm:"size:10;not null" json:"action"` // buy or sell
	Quantity      int            `gorm:"not null" json:"quantity"`
	Price         float64        `gorm:"not null" json:"price"`
	OriginalValue float64        `gorm:"not null" json:"original_value"`
	Date          time.Time      `json:"date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name
func (Stock) TableName() string {
	return "stocks"
}
