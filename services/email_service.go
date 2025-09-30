package services

import (
	"bytes"
	"fmt"
	"html/template"
	"notification-server/config"
	"notification-server/models"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

// NewEmailService creates a new email service instance
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendUserProfileEmail sends an email with user profile details
func (s *EmailService) SendUserProfileEmail(user *models.User) error {
	// Create email template
	emailBody, err := s.generateUserProfileHTML(user)
	if err != nil {
		return fmt.Errorf("failed to generate email template: %w", err)
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.Email.FromName, s.config.Email.FromEmail))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Welcome! Your User Profile Details")
	m.SetBody("text/html", emailBody)

	// Send email
	port, err := strconv.Atoi(s.config.Email.SMTPPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	d := gomail.NewDialer(
		s.config.Email.SMTPHost,
		port,
		s.config.Email.Username,
		s.config.Email.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// generateUserProfileHTML generates HTML email template for user profile
func (s *EmailService) generateUserProfileHTML(user *models.User) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            text-align: center;
            margin: -30px -30px 30px -30px;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .profile-section {
            margin: 20px 0;
        }
        .profile-item {
            background-color: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid #667eea;
            border-radius: 5px;
        }
        .profile-item label {
            font-weight: bold;
            color: #667eea;
            display: block;
            margin-bottom: 5px;
        }
        .profile-item .value {
            color: #333;
            font-size: 16px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 2px solid #eee;
            text-align: center;
            color: #888;
            font-size: 14px;
        }
        .welcome-text {
            text-align: center;
            margin: 20px 0;
            color: #555;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to Our Platform!</h1>
        </div>
        
        <div class="welcome-text">
            <p>Hi <strong>{{.FirstName}}</strong>,</p>
            <p>Thank you for signing up! Here are your profile details:</p>
        </div>

        <div class="profile-section">
            <div class="profile-item">
                <label>Full Name:</label>
                <div class="value">{{.FirstName}} {{.LastName}}</div>
            </div>

            <div class="profile-item">
                <label>Email Address:</label>
                <div class="value">{{.Email}}</div>
            </div>

            {{if .Phone}}
            <div class="profile-item">
                <label>Phone Number:</label>
                <div class="value">{{.Phone}}</div>
            </div>
            {{end}}

            {{if .Address}}
            <div class="profile-item">
                <label>Address:</label>
                <div class="value">{{.Address}}</div>
            </div>
            {{end}}

            {{if .City}}
            <div class="profile-item">
                <label>City:</label>
                <div class="value">{{.City}}</div>
            </div>
            {{end}}

            {{if .Country}}
            <div class="profile-item">
                <label>Country:</label>
                <div class="value">{{.Country}}</div>
            </div>
            {{end}}

            <div class="profile-item">
                <label>Account Created:</label>
                <div class="value">{{.CreatedAt.Format "January 02, 2006 at 3:04 PM"}}</div>
            </div>
        </div>

        <div class="footer">
            <p>If you have any questions, feel free to contact our support team.</p>
            <p>&copy; 2025 Notification Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("userProfile").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, user); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SendOTPEmail sends an OTP verification email to the user
func (s *EmailService) SendOTPEmail(email, otp string) error {
	// Create email template
	emailBody, err := s.generateOTPHTML(otp)
	if err != nil {
		return fmt.Errorf("failed to generate OTP email template: %w", err)
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.Email.FromName, s.config.Email.FromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your OTP Verification Code")
	m.SetBody("text/html", emailBody)

	// Send email
	port, err := strconv.Atoi(s.config.Email.SMTPPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	d := gomail.NewDialer(
		s.config.Email.SMTPHost,
		port,
		s.config.Email.Username,
		s.config.Email.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

// generateOTPHTML generates HTML email template for OTP verification
func (s *EmailService) generateOTPHTML(otp string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            text-align: center;
            margin: -30px -30px 30px -30px;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .otp-container {
            text-align: center;
            margin: 30px 0;
        }
        .otp-code {
            font-size: 36px;
            font-weight: bold;
            color: #667eea;
            letter-spacing: 8px;
            padding: 20px;
            background-color: #f8f9fa;
            border-radius: 10px;
            border: 2px dashed #667eea;
            display: inline-block;
        }
        .message {
            text-align: center;
            color: #555;
            margin: 20px 0;
        }
        .warning {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 2px solid #eee;
            text-align: center;
            color: #888;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Email Verification</h1>
        </div>
        
        <div class="message">
            <p>Thank you for signing up! Please use the following OTP to verify your email address:</p>
        </div>

        <div class="otp-container">
            <div class="otp-code">{{.}}</div>
        </div>

        <div class="warning">
            <strong>⚠️ Important:</strong>
            <ul style="margin: 10px 0; padding-left: 20px;">
                <li>This OTP is valid for 10 minutes only</li>
                <li>Do not share this code with anyone</li>
                <li>If you didn't request this, please ignore this email</li>
            </ul>
        </div>

        <div class="footer">
            <p>If you have any questions, feel free to contact our support team.</p>
            <p>&copy; 2025 Notification Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("otp").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, otp); err != nil {
		return "", err
	}

	return buf.String(), nil
}
