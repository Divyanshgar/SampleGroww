# Quick Start Guide

Get up and running with the Notification Service in minutes!

## Prerequisites

- Go 1.21+
- Node.js 16+
- PostgreSQL
- SMTP credentials (e.g., Gmail App Password)

## 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd Notification

# Install backend dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install
cd ..
```

## 2. Configure Environment

```bash
# Copy the example .env file
cp .env.example .env
```

Edit `.env` with your configuration:
- Update database credentials
- Add your SMTP settings

## 3. Setup Database

```sql
CREATE DATABASE notification_db;
```

## 4. Run in Development Mode

### Terminal 1 - Backend
```bash
go run main.go
```

### Terminal 2 - Frontend
```bash
cd frontend
npm start
```

## 5. Access the Application

Open your browser to:
- Frontend: http://localhost:3000
- API: http://localhost:8443

## 6. Test the Flow

1. Enter your email to receive an OTP
2. Check your email for the 6-digit code
3. Enter the OTP to verify
4. Complete the registration form
5. Check your email for the welcome message

## Production Deployment

1. Build the frontend:
   ```bash
   cd frontend
   npm run build
   cd ..
   ```

2. Build the Go application:
   ```bash
   go build -o notification-server
   ```

3. Run the server:
   ```bash
   ./notification-server
   ```

The application will be available at `http://your-domain.com:8443`

## Troubleshooting

- **Emails not sending?** Verify your SMTP settings in `.env`
- **Database connection issues?** Check your PostgreSQL credentials
- **Frontend not loading?** Ensure you've built the React app

For more details, see the full [README.md](README.md)
