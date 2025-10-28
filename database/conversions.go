package database

import (
	"database/sql"
	"notification-server/db"
	"notification-server/models"
	"strconv"
	"time"
)

// ConvertDBUserToModel converts SQLC User to models.User
func ConvertDBUserToModel(dbUser db.User) models.User {
	user := models.User{
		ID:          uint(dbUser.ID),
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		Email:       dbUser.Email,
		Phone:       "",
		Address:     "",
		City:        "",
		Country:     "",
		OTPCode:     "",
		OTPExpires:  time.Time{},
		OTPVerified: false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	// Handle nullable fields
	if dbUser.Phone.Valid {
		user.Phone = dbUser.Phone.String
	}
	if dbUser.Address.Valid {
		user.Address = dbUser.Address.String
	}
	if dbUser.City.Valid {
		user.City = dbUser.City.String
	}
	if dbUser.Country.Valid {
		user.Country = dbUser.Country.String
	}
	if dbUser.OtpCode.Valid {
		user.OTPCode = dbUser.OtpCode.String
	}
	if dbUser.OtpExpires.Valid {
		user.OTPExpires = dbUser.OtpExpires.Time
	}
	if dbUser.OtpVerified.Valid {
		user.OTPVerified = dbUser.OtpVerified.Bool
	}
	if dbUser.CreatedAt.Valid {
		user.CreatedAt = dbUser.CreatedAt.Time
	}
	if dbUser.UpdatedAt.Valid {
		user.UpdatedAt = dbUser.UpdatedAt.Time
	}

	return user
}

// ConvertDBStockToModel converts SQLC Stock to models.Stock
func ConvertDBStockToModel(dbStock db.Stock) models.Stock {
	stock := models.Stock{
		ID:            uint(dbStock.ID),
		UserID:        uint(dbStock.UserID),
		StockName:     dbStock.StockName,
		Action:        dbStock.Action,
		Quantity:      int(dbStock.Quantity),
		Price:         0.0,
		OriginalValue: 0.0,
		Date:          time.Time{},
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	// Handle price and original_value (stored as strings in SQLC)
	if dbStock.Price != "" {
		if price, err := strconv.ParseFloat(dbStock.Price, 64); err == nil {
			stock.Price = price
		}
	}
	if dbStock.OriginalValue != "" {
		if originalValue, err := strconv.ParseFloat(dbStock.OriginalValue, 64); err == nil {
			stock.OriginalValue = originalValue
		}
	}
	if dbStock.Date.Valid {
		stock.Date = dbStock.Date.Time
	}
	if dbStock.CreatedAt.Valid {
		stock.CreatedAt = dbStock.CreatedAt.Time
	}
	if dbStock.UpdatedAt.Valid {
		stock.UpdatedAt = dbStock.UpdatedAt.Time
	}

	return stock
}

// ConvertModelUserToDBParams converts models.User to CreateUserParams
func ConvertModelUserToDBParams(user models.User) db.CreateUserParams {
	return db.CreateUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Address:   sql.NullString{String: user.Address, Valid: user.Address != ""},
		City:      sql.NullString{String: user.City, Valid: user.City != ""},
		Country:   sql.NullString{String: user.Country, Valid: user.Country != ""},
	}
}

// ConvertModelStockToDBParams converts models.Stock to CreateStockParams
func ConvertModelStockToDBParams(stock models.Stock) db.CreateStockParams {
	return db.CreateStockParams{
		UserID:        int32(stock.UserID),
		StockName:     stock.StockName,
		Action:        stock.Action,
		Quantity:      int32(stock.Quantity),
		Price:         strconv.FormatFloat(stock.Price, 'f', -1, 64),
		OriginalValue: strconv.FormatFloat(stock.OriginalValue, 'f', -1, 64),
		Date:          sql.NullTime{Time: stock.Date, Valid: !stock.Date.IsZero()},
	}
}
