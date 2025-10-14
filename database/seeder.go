package database

import (
	"log"
	"math/rand"
	"notification-server/models"
	"time"
)

// SeedStockData seeds the database with random stock data for sample users
func SeedStockData() error {
	log.Println("Seeding stock data...")

	// Check if users exist, if not, create sample users
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
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
			if err := DB.Create(&user).Error; err != nil {
				return err
			}
		}
		log.Println("Sample users created successfully")
	}

	// Fetch all users
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		return err
	}

	// Stock names
	stocks := []string{"Apple", "Google", "Microsoft", "Tesla", "Amazon"}

	// Seed random stock transactions for each user
	rand.Seed(time.Now().UnixNano())

	for _, user := range users {
		// Check if user already has stock data
		var stockCount int64
		DB.Model(&models.Stock{}).Where("user_id = ?", user.ID).Count(&stockCount)
		if stockCount > 0 {
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
			quantity := rand.Intn(100) + 1 // 1 to 100
			price := rand.Float64()*490 + 10 // 10 to 500
			originalValue := float64(quantity) * price

			// Random date in the last year
			daysAgo := rand.Intn(365)
			date := time.Now().AddDate(0, 0, -daysAgo)

			stock := models.Stock{
				UserID:        user.ID,
				StockName:     stockName,
				Action:        action,
				Quantity:      quantity,
				Price:         price,
				OriginalValue: originalValue,
				Date:          date,
			}

			if err := DB.Create(&stock).Error; err != nil {
				return err
			}
		}
		log.Printf("Seeded %d stock transactions for user %s", numTransactions, user.Email)
	}

	log.Println("Stock data seeding completed successfully")
	return nil
}
