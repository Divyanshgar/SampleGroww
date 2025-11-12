package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"notification-server/db"
	"notification-server/models"
	"time"
)

// SeedStockData seeds the database with random stock data for sample users
func SeedStockData() error {
	log.Println("Seeding stock data...")

	// Ensure database queries are initialized
	if Queries == nil {
		return fmt.Errorf("database queries not initialized; run database.Initialize first")
	}

	// Check if users exist, if not, create sample users
	users, err := Queries.ListUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No users found, creating sample users...")
		sampleUsers := []models.User{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Phone:     "+1234567890",
				Address:   "123 Main St",
				City:      "New York",
				Country:   "USA",
			},
			{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
				Phone:     "+1234567891",
				Address:   "456 Elm St",
				City:      "Los Angeles",
				Country:   "USA",
			},
			{
				FirstName: "Bob",
				LastName:  "Johnson",
				Email:     "bob.johnson@example.com",
				Phone:     "+1234567892",
				Address:   "789 Oak St",
				City:      "Chicago",
				Country:   "USA",
			},
		}

		for _, user := range sampleUsers {
			_, err := Queries.CreateUser(context.Background(), db.CreateUserParams{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Phone:     sql.NullString{String: user.Phone, Valid: true},
				Address:   sql.NullString{String: user.Address, Valid: true},
				City:      sql.NullString{String: user.City, Valid: true},
				Country:   sql.NullString{String: user.Country, Valid: true},
			})
			if err != nil {
				return err
			}
		}
		log.Println("Sample users created successfully")
	}

	// Fetch all users again after potential creation
	users, err = Queries.ListUsers(context.Background())
	if err != nil {
		return err
	}

	// Stock names
	stocks := []string{"Apple", "Google", "Microsoft", "Tesla", "Amazon"}

	// Seed random stock transactions for each user
	rand.Seed(time.Now().UnixNano())

	for _, user := range users {
		// Check if user already has stock data
		userStocks, err := Queries.ListStocksByUser(context.Background(), user.ID)
		if err != nil {
			return err
		}
		if len(userStocks) > 0 {
			log.Printf("User %s already has stock data, skipping...", user.Email)
			continue
		}

		// Generate 5-15 random transactions per user
		numTransactions := rand.Intn(11) + 5 // 5 to 15

		for i := 0; i < numTransactions; i++ {
			stockName := stocks[rand.Intn(len(stocks))]
			action := "buy"
			if rand.Float32() < 0.3 { // 30% chance of sell
				action = "sell"
			}
			quantity := rand.Intn(100) + 1   // 1 to 100
			price := rand.Float64()*490 + 10 // 10 to 500
			originalValue := float64(quantity) * price

			// Random date in the last year
			daysAgo := rand.Intn(365)
			date := time.Now().AddDate(0, 0, -daysAgo)

			_, err := Queries.CreateStock(context.Background(), db.CreateStockParams{
				UserID:        user.ID,
				StockName:     stockName,
				Action:        action,
				Quantity:      int32(quantity),
				Price:         fmt.Sprintf("%.2f", price),
				OriginalValue: fmt.Sprintf("%.2f", originalValue),
				Date:          sql.NullTime{Time: date, Valid: true},
			})
			if err != nil {
				return err
			}
		}
		log.Printf("Seeded %d stock transactions for user %s", numTransactions, user.Email)
	}

	log.Println("Stock data seeding completed successfully")
	return nil
}
