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
	pdfService   *services.PDFService
}

func NewNotificationHandler(emailService *services.EmailService, pdfService *services.PDFService) *NotificationHandler {
	return &NotificationHandler{
		emailService: emailService,
		pdfService:   pdfService,
	}
}

// ====================== HEALTH CHECK ======================
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "notification-server",
	})
}

// ====================== SEND USER PROFILE EMAIL ======================

type SendUserProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
}

func (h *NotificationHandler) SendUserProfile(c *gin.Context) {
	var req SendUserProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Save user to database
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		City:      req.City,
		Country:   req.Country,
	}

	if err := database.GetDB().Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save user to database",
			"details": err.Error(),
		})
		return
	}

	// Generate PDF for user profile
	pdfBytes, err := h.pdfService.GenerateUserProfilePDF(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate PDF",
			"details": err.Error(),
		})
		return
	}

	// Debug: Save PDF to file for verification
	err = h.pdfService.SavePDFToFile("debug_user_profile.pdf", pdfBytes)
	if err != nil {
		log.Printf("Failed to save debug PDF file: %v", err)
	}

	// Send email with PDF attachment
	if err := h.emailService.SendUserProfileEmailWithAttachment(user, pdfBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send email with PDF attachment",
			"details": err.Error(),
			"message": "User was saved to database but email sending failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile email with PDF sent successfully",
		"user_id": user.ID,
		"email":   user.Email,
	})
}

// ====================== GET USER PROFILE ======================

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

// ====================== REQUEST OTP ======================

type RequestOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *NotificationHandler) RequestOTP(c *gin.Context) {
	var req RequestOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// User exists, check if already verified/registered
		if existingUser.OTPVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
			return
		}
		// User exists but not verified, allow resend OTP
	} else {
		// User does not exist, create new
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP", "details": err.Error()})
		return
	}

	// Save or update OTP in DB (valid for 10 minutes)
	if err := utils.SaveOrUpdateOTP(req.Email, otp, 10); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP in database", "details": err.Error()})
		return
	}

	// Send OTP via email
	if err := h.emailService.SendOTPEmail(req.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully to your email",
		"email":   req.Email,
	})
}

// ====================== VERIFY OTP AND REGISTER ======================

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

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

func (h *NotificationHandler) VerifyOTPAndRegister(c *gin.Context) {
	var req VerifyOTPAndRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Verify OTP from DB
	if !utils.VerifyAndMarkOTP(req.Email, req.OTP) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Find the existing user
	var user models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found", "details": err.Error()})
		return
	}

	// Update user with full information
	updates := map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"phone":      req.Phone,
		"address":    req.Address,
		"city":       req.City,
		"country":    req.Country,
	}
	if err := database.GetDB().Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	// Optional: cleanup OTPs after successful registration
	utils.DeleteOTP(req.Email)

	// Generate PDF for user profile
	pdfBytes, err := h.pdfService.GenerateUserProfilePDF(&user)
	if err != nil {
		log.Printf("Failed to generate PDF for user %s: %v", user.Email, err)
		// Send email without attachment if PDF generation fails
		if err := h.emailService.SendUserProfileEmail(&user); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
		}
	} else {
		// Send welcome email with PDF attachment
		if err := h.emailService.SendUserProfileEmailWithAttachment(&user, pdfBytes); err != nil {
			log.Printf("Failed to send welcome email with attachment to %s: %v", user.Email, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
		"email":   user.Email,
	})
}

func (h *NotificationHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Verify OTP from DB
	if !utils.CheckOTP(req.Email, req.OTP) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"email":   req.Email,
	})
}

// ====================== GENERATE USER PROFILE PDF ======================

func (h *NotificationHandler) GenerateUserProfilePDF(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("Generating PDF for user ID: %s", userID)

	if database.GetDB() == nil {
		log.Printf("Database connection is nil")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		log.Printf("User not found error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	log.Printf("User found: %+v", user)

	pdfBytes, err := h.pdfService.GenerateUserProfilePDF(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate PDF",
			"details": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=user_profile.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
