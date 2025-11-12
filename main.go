package main

import (
	"fmt"
	"log"
	"notification-server/config"
	"notification-server/database"
	"notification-server/routes"
	"notification-server/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Println("Configuration loaded successfully")

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Seed stock data
	if err := database.SeedStockData(); err != nil {
		log.Fatalf("Failed to seed stock data: %v", err)
	}
	log.Println("Stock data seeded successfully")

	// Initialize email service
	emailService := services.NewEmailService(cfg)
	log.Println("Email service initialized successfully")

	// Initialize PDF service
	pdfService := services.NewPDFService()
	log.Println("PDF service initialized successfully")

	// Initialize Excel service
	excelService := services.NewExcelService()
	log.Println("Excel service initialized successfully")

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	// Serve static files from React build
	router.Static("/static", "./frontend/build/static")
	router.StaticFile("/favicon.ico", "./frontend/build/favicon.ico")
	router.StaticFile("/manifest.json", "./frontend/build/manifest.json")
	router.StaticFile("/logo.jpg", "./frontend/build/logo.jpg")

	// Setup API routes
	routes.SetupRoutes(router, emailService, pdfService, excelService)

	// Serve React app for all other routes (SPA fallback)
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})

	// Start HTTP server
	serverAddr := fmt.Sprintf("127.0.0.1:%s", cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", serverAddr)
	log.Printf("Frontend available at: http://localhost:%s", cfg.Server.Port)

	// Run without TLS
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
