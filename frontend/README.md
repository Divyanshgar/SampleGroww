# Notification Service - Frontend

React frontend for user registration with OTP email verification.

## Features

- ✉️ Email input and OTP request
- 🔐 6-digit OTP verification
- 📝 Multi-step registration form
- ✅ Success confirmation with user details
- 📱 Responsive design
- 🎨 Beautiful gradient UI

## Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Development Mode

Run the React dev server (with hot reload):

```bash
npm start
```

The app will open at `http://localhost:3000` and proxy API requests to `http://localhost:8443`.

### 3. Build for Production

```bash
npm run build
```

This creates an optimized production build in the `build/` directory.

## Project Structure

```
frontend/
├── public/
│   └── index.html          # HTML template
├── src/
│   ├── components/         # React components
│   │   ├── EmailForm.js
│   │   ├── OTPVerification.js
│   │   ├── RegistrationForm.js
│   │   └── SuccessMessage.js
│   ├── App.js             # Main application component
│   ├── App.css            # Application styles
│   ├── index.js           # Entry point
│   └── index.css          # Global styles
└── package.json           # Dependencies and scripts
```

## API Integration

The frontend communicates with the backend API:

- **POST** `/api/v1/auth/request-otp` - Request OTP
- **POST** `/api/v1/auth/verify-otp` - Verify OTP and register user

## Environment

The app uses a proxy configuration in `package.json` to forward API requests to the backend server during development.

## Deployment

After building, the `build/` directory can be served by:
- The Go backend (configured in `main.go`)
- Any static file server (Nginx, Apache, etc.)
- CDN or hosting platform (Netlify, Vercel, etc.)

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
