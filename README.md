# Notification Server with OTP Verification

A secure notification server built with Go and Gin framework featuring email OTP verification, user registration, and a React frontend.

## Features

- 🔐 **OTP Email Verification**: Secure 6-digit OTP sent via email
- ✉️ **Email Notifications**: Beautiful HTML email templates
- 📝 **User Registration**: Multi-step registration with email verification
- 🗄️ **PostgreSQL Database**: User data persistence with GORM
- ⚛️ **React Frontend**: Modern, responsive UI with step-by-step flow
- 🚀 **RESTful API**: Clean API endpoints using Gin framework
- ⚙️ **Environment Configuration**: Easy configuration via `.env` file

## Tech Stack

### Backend
- **Go 1.21**
- **Gin Web Framework**
- **GORM** (PostgreSQL)
- **Gomail** (Email sending)

### Frontend
- **React 18**
- **Axios** (HTTP client)
- **Modern CSS** (Responsive design)

## Project Structure

```
notification-server/
├── config/                 # Configuration management
│   └── config.go
├── database/              # Database connection and migrations
│   └── database.go
├── handlers/              # HTTP request handlers
│   └── notification_handler.go
├── models/                # Data models
│   ├── user.go
│   └── otp.go
├── routes/                # Route definitions
│   └── routes.go
├── services/              # Business logic services
│   └── email_service.go
├── utils/                 # Utility functions
│   └── otp_store.go
├── frontend/              # React frontend
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── App.js
│   │   └── index.js
│   └── package.json
├── .env                   # Environment variables
├── .env.example          # Environment variables template
├── go.mod                # Go dependencies
└── main.go               # Application entry point
```

## Prerequisites

Before running the application, ensure you have:

1. **Go 1.21+** installed
2. **Node.js 16+** and npm installed
3. **PostgreSQL** database running
4. **SMTP credentials** (e.g., Gmail App Password)

## Setup Instructions

### 1. Clone or Navigate to Project

```bash
cd "d:/GO Learning/Notification"
```

### 2. Install Backend Dependencies

```bash
go mod download
go mod tidy
```

### 3. Install Frontend Dependencies

```bash
cd frontend
npm install
cd ..
```

### 4. Setup PostgreSQL Database

Create a database for the application:

```sql
CREATE DATABASE notification_db;
```

### 5. Configure Environment Variables

Copy the example environment file and update with your credentials:

```bash
cp .env.example .env
```

Edit `.env` file with your settings:

```env
# Server Configuration
SERVER_PORT=8443
SERVER_HOST=0.0.0.0

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=notification_db
DB_SSLMODE=disable

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM_EMAIL=your_email@gmail.com
SMTP_FROM_NAME=Notification Service
```

### 6. Setup Gmail SMTP (if using Gmail)

1. Enable 2-Factor Authentication in your Google Account
2. Generate an App Password:
   - Go to Google Account Settings
   - Security → 2-Step Verification → App passwords
   - Generate a new app password
3. Use this app password in your `.env` file as `SMTP_PASSWORD`

### 7. Build React Frontend

```bash
cd frontend
npm run build
cd ..
```

This creates a production build in `frontend/build/` that will be served by the Go backend.

## Running the Application

### Option 1: Development Mode (Frontend + Backend Separate)

**Terminal 1 - Start Backend:**
```bash
go run main.go
```

**Terminal 2 - Start Frontend Dev Server:**
```bash
cd frontend
npm start
```

Frontend will be available at `http://localhost:3000` with hot reload.

### Option 2: Production Mode (Integrated)

**Build Frontend First:**
```bash
cd frontend
npm run build
cd ..
```

**Start Backend:**
```bash
go run main.go
```

The complete application will be available at `http://localhost:8443`

### Build and Run Executable

```bash
# Build the Go application
go build -o notification-server

# Run the executable (Windows)
.\notification-server.exe

# Run the executable (Linux/Mac)
./notification-server
```

## User Flow

1. **Enter Email**: User enters their email address
2. **Receive OTP**: System sends a 6-digit OTP to the email
3. **Verify OTP**: User enters the OTP to verify their email
4. **Complete Registration**: User fills out profile details
5. **Success**: User is registered and receives a welcome email

## API Endpoints

### Authentication Endpoints

#### 1. Request OTP
**POST** `/api/v1/auth/request-otp`

Sends an OTP to the user's email for verification.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "message": "OTP sent successfully to your email",
  "email": "user@example.com"
}
```

#### 2. Verify OTP and Register
**POST** `/api/v1/auth/verify-otp`

Verifies the OTP and registers the user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main Street",
  "city": "New York",
  "country": "USA"
}
```

**Response (Success):**
```json
{
  "message": "User profile email sent successfully",
  "user_id": 1,
  "email": "john.doe@example.com"
}
```

**Response (Error):**
```json
{
  "error": "Failed to send email",
  "details": "error message"
}
```

### 3. Get User Profile
**GET** `/api/v1/users/:id`

Retrieve a user's profile by their ID.

**Response:**
```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main Street",
  "city": "New York",
  "country": "USA",
  "created_at": "2025-09-30T11:39:43Z",
  "updated_at": "2025-09-30T11:39:43Z"
}
```

## Testing with cURL

### Test Health Endpoint
```bash
curl -k https://localhost:8443/health
```

### Send User Profile Email
```bash
curl -k -X POST https://localhost:8443/api/v1/notifications/send-user-profile \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "address": "123 Main Street",
    "city": "New York",
    "country": "USA"
  }'
```

### Get User Profile
```bash
curl -k https://localhost:8443/api/v1/users/1
```

**Note:** The `-k` flag is used to skip SSL certificate verification for self-signed certificates.

## Database Schema

### Users Table

| Column     | Type      | Constraints          |
|------------|-----------|----------------------|
| id         | SERIAL    | PRIMARY KEY          |
| first_name | VARCHAR   | NOT NULL             |
| last_name  | VARCHAR   | NOT NULL             |
| email      | VARCHAR   | NOT NULL, UNIQUE     |
| phone      | VARCHAR   |                      |
| address    | VARCHAR   |                      |
| city       | VARCHAR   |                      |
| country    | VARCHAR   |                      |
| created_at | TIMESTAMP | AUTO                 |
| updated_at | TIMESTAMP | AUTO                 |
| deleted_at | TIMESTAMP | NULLABLE             |

## Email Template

The service sends beautifully formatted HTML emails with:
- Gradient header design
- User profile details in styled cards
- Responsive design
- Professional footer

## Security Considerations

### Development
- Self-signed certificates are used for development
- Certificate warnings in browsers are normal

### Production
- Use certificates from a trusted Certificate Authority (Let's Encrypt, etc.)
- Store sensitive credentials in secure vaults (not in `.env`)
- Enable database SSL mode
- Use strong passwords
- Implement rate limiting
- Add authentication middleware

## Troubleshooting

### SSL Certificate Issues
If you encounter SSL certificate errors:
```bash
# Regenerate certificates
cd scripts
bash generate_certs.sh
```

### Email Not Sending
1. Verify SMTP credentials in `.env`
2. Check if Gmail allows "Less secure apps" (if using Gmail)
3. Use App Password for Gmail with 2FA
4. Check firewall settings for port 587

### Database Connection Issues
1. Ensure PostgreSQL is running
2. Verify database credentials in `.env`
3. Check if the database exists
4. Verify network connectivity

### Port Already in Use
If port 8443 is in use, change `SERVER_PORT` in `.env`:
```env
SERVER_PORT=9443
```

## Development

### Add New Dependencies
```bash
go get <package-name>
go mod tidy
```

### Run Tests
```bash
go test ./...
```

### Format Code
```bash
go fmt ./...
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue in the repository.

---

**Happy Coding! 🚀**
