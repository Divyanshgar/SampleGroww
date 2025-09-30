package handlers

import (
	"log"
	"net/http"
	"notification-server/database"
	"notification-server/models"
	"notification-server/services"
	"notification-server/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	emailService *services.EmailService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(emailService *services.EmailService) *NotificationHandler {
	return &NotificationHandler{
		emailService: emailService,
	}
}

// SendUserProfileRequest represents the request body for sending user profile email
type SendUserProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
}

// SendUserProfile handles sending user profile details via email
// @Summary Send user profile email
// @Description Sends an email to the user with their profile details after signup
// @Tags Notifications
// @Accept json
// @Produce json
// @Param user body SendUserProfileRequest true "User Profile Details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/notifications/send-user-profile [post]
func (h *NotificationHandler) SendUserProfile(c *gin.Context) {
	var req SendUserProfileRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create user in database
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		City:      req.City,
		Country:   req.Country,
	}

	// Save user to database
	if err := database.GetDB().Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save user to database",
			"details": err.Error(),
		})
		return
	}

	// Send email
	if err := h.emailService.SendUserProfileEmail(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send email",
			"details": err.Error(),
			"message": "User was saved to database but email sending failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile email sent successfully",
		"user_id": user.ID,
		"email":   user.Email,
	})
}

// GetUserProfile retrieves a user profile by ID
// @Summary Get user profile
// @Description Retrieves a user profile from the database
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [get]
func (h *NotificationHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// HealthCheck endpoint for server health monitoring
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "notification-server",
	})
}

// RequestOTPRequest represents the request body for OTP request
type RequestOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestOTP sends an OTP to the user's email
// @Summary Request OTP
// @Description Sends an OTP to the user's email for verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RequestOTPRequest true "Email for OTP"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/request-otp [post]
func (h *NotificationHandler) RequestOTP(c *gin.Context) {
	var req RequestOTPRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already registered",
		})
		return
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate OTP",
			"details": err.Error(),
		})
		return
	}

	// Store OTP (valid for 10 minutes)
	otpStore := utils.GetOTPStore()
	otpStore.Save(req.Email, otp, 10)

	// Send OTP email
	if err := h.emailService.SendOTPEmail(req.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send OTP email",
			"details": err.Error(),
		})
		return
	}

	log.Printf("OTP sent to %s: %s", req.Email, otp) // For development only

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully to your email",
		"email":   req.Email,
	})
}

// VerifyOTPAndRegisterRequest represents the request body for OTP verification and registration
type VerifyOTPAndRegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	OTP       string `json:"otp" binding:"required,len=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
}

// VerifyOTPAndRegister verifies the OTP and registers the user
// @Summary Verify OTP and Register User
// @Description Verifies OTP and registers the user in the database
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyOTPAndRegisterRequest true "User Registration with OTP"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/verify-otp [post]
func (h *NotificationHandler) VerifyOTPAndRegister(c *gin.Context) {
	var req VerifyOTPAndRegisterRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Verify OTP
	otpStore := utils.GetOTPStore()
	if !otpStore.Verify(req.Email, req.OTP) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired OTP",
		})
		return
	}

	// Create user in database
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		City:      req.City,
		Country:   req.Country,
	}

	// Save user to database
	if err := database.GetDB().Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save user to database",
			"details": err.Error(),
		})
		return
	}

	// Delete the used OTP
	otpStore.Delete(req.Email)

	// Send welcome email
	if err := h.emailService.SendUserProfileEmail(user); err != nil {
		log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
		// Don't fail the request if email sending fails
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
		"email":   user.Email,
	})
}
