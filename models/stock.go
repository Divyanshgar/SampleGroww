package models

import (
	"time"
)

// Stock represents a stock transaction for a user
type Stock struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	StockName     string    `json:"stock_name"`
	Action        string    `json:"action"` // buy or sell
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	OriginalValue float64   `json:"original_value"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
