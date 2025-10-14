package handlers

import (
	"fmt"
	"log"
	"net/http"
	"notification-server/database"
	"notification-server/models"
	"notification-server/services"
	"notification-server/utils"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	emailService *services.EmailService
	pdfService   *services.PDFService
	excelService *services.ExcelService
}

func NewNotificationHandler(emailService *services.EmailService, pdfService *services.PDFService, excelService *services.ExcelService) *NotificationHandler {
	return &NotificationHandler{
		emailService: emailService,
		pdfService:   pdfService,
		excelService: excelService,
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

func (h *NotificationHandler) GenerateUserProfileHTML(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("Generating HTML for user ID: %s", userID)

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

	// Fetch stocks for the user
	var stocks []models.Stock
	if err := database.GetDB().Where("user_id = ?", user.ID).Find(&stocks).Error; err != nil {
		log.Printf("Failed to fetch stocks: %v", err)
		stocks = []models.Stock{}
	}

	// Sort stocks by date
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Date.Before(stocks[j].Date)
	})

	// Calculate total value
	var totalValue float64
	for _, stock := range stocks {
		totalValue += stock.OriginalValue
	}

	// Generate HTML
	htmlContent := h.generateUserProfileHTML(&user, stocks, totalValue)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlContent)
}

// generateUserProfileHTML generates the HTML content for user profile
func (h *NotificationHandler) generateUserProfileHTML(user *models.User, stocks []models.Stock, totalValue float64) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>User Profile and Stock Data</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f9f9f9; }
        .header { text-align: center; margin-bottom: 20px; }
        .logo { max-width: 200px; height: auto; }
        .company-name { font-size: 18px; font-weight: bold; color: #000; margin: 10px 0; }
        .generated-date { font-size: 12px; color: #666; }
        .section-title { font-size: 16px; font-weight: bold; background-color: #667eea; color: white; padding: 10px; text-align: center; margin: 20px 0 10px 0; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; background-color: white; }
        th, td { border: 1px solid #000; padding: 8px; text-align: left; }
        th { background-color: #667eea; color: white; font-weight: bold; text-align: center; }
        .even-row { background-color: #f2f2f2; }
        .odd-row { background-color: white; }
        .buy { background-color: #d4edda; }
        .sell { background-color: #f8d7da; }
        .totals-row { background-color: #e9ecef; font-weight: bold; }
        .no-data { text-align: center; color: #666; font-style: italic; padding: 20px; }
    </style>
</head>
<body>
    <div class="header">
        <img src="./logo.jpg" alt="Centricity Logo" class="logo" onerror="this.style.display='none'">
        <div class="company-name">Centricity Financial Distribution Private Limited</div>
        <div class="generated-date">Generated on: ` + currentTime + `</div>
    </div>

    <div class="section-title">User Profile Summary</div>
    <table>
        <thead>
            <tr>
                <th>Field</th>
                <th>Value</th>
            </tr>
        </thead>
        <tbody>
            <tr class="even-row"><td>User ID</td><td>` + fmt.Sprintf("%d", user.ID) + `</td></tr>
            <tr class="odd-row"><td>First Name</td><td>` + user.FirstName + `</td></tr>
            <tr class="even-row"><td>Last Name</td><td>` + user.LastName + `</td></tr>
            <tr class="odd-row"><td>Email</td><td>` + user.Email + `</td></tr>
            <tr class="even-row"><td>Phone</td><td>` + user.Phone + `</td></tr>
            <tr class="odd-row"><td>Address</td><td>` + user.Address + `</td></tr>
            <tr class="even-row"><td>City</td><td>` + user.City + `</td></tr>
            <tr class="odd-row"><td>Country</td><td>` + user.Country + `</td></tr>
        </tbody>
    </table>

    <div class="section-title">Stock Transactions</div>
    <table>
        <thead>
            <tr>
                <th>Stock Name</th>
                <th>Action</th>
                <th>Quantity</th>
                <th>Price</th>
                <th>Original Value</th>
                <th>Date</th>
            </tr>
        </thead>
        <tbody>`

	if len(stocks) == 0 {
		htmlTemplate += `<tr><td colspan="6" class="no-data">No stock transactions available for this user.</td></tr>`
	} else {
		for i, stock := range stocks {
			rowClass := "even-row"
			if i%2 == 1 {
				rowClass = "odd-row"
			}
			actionClass := ""
			if stock.Action == "buy" {
				actionClass = " buy"
			} else if stock.Action == "sell" {
				actionClass = " sell"
			}
			htmlTemplate += fmt.Sprintf(`
            <tr class="%s%s">
                <td>%s</td>
                <td>%s</td>
                <td>%d</td>
                <td>%.2f</td>
                <td>%.2f</td>
                <td>%s</td>
            </tr>`, rowClass, actionClass, stock.StockName, stock.Action, stock.Quantity, stock.Price, stock.OriginalValue, stock.Date.Format("2006-01-02"))
		}
		htmlTemplate += fmt.Sprintf(`
            <tr class="totals-row">
                <td colspan="4">Total Original Value</td>
                <td>%.2f</td>
                <td></td>
            </tr>`, totalValue)
	}

	htmlTemplate += `
        </tbody>
    </table>
</body>
</html>`

	return htmlTemplate
}

// ====================== GENERATE USER PROFILE EXCEL ======================

func (h *NotificationHandler) GenerateUserProfileExcel(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("Generating Excel for user ID: %s", userID)

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

	excelBytes, err := h.excelService.GenerateUserProfileExcel(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate Excel",
			"details": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=user_profile.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}
