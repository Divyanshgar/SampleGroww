package routes

import (
	"notification-server/handlers"
	"notification-server/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, emailService *services.EmailService) {
	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(emailService)

	// Health check endpoint
	router.GET("/health", notificationHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/request-otp", notificationHandler.RequestOTP)
			auth.POST("/verify-otp", notificationHandler.VerifyOTPAndRegister)
		}

		// Notification routes
		notifications := v1.Group("/notifications")
		{
			notifications.POST("/send-user-profile", notificationHandler.SendUserProfile)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/:id", notificationHandler.GetUserProfile)
		}
	}
}
