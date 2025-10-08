-- Create users table for PostgreSQL
-- This table stores user information including OTP verification data

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    address VARCHAR(500),
    city VARCHAR(100),
    country VARCHAR(100),
    otp_code VARCHAR(6),  -- Stores the 6-digit OTP code
    otp_expires TIMESTAMP,  -- Expiration time for the OTP
    otp_verified BOOLEAN DEFAULT FALSE,  -- Whether OTP has been verified
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Optional: Create index on otp_verified for OTP-related queries
CREATE INDEX IF NOT EXISTS idx_users_otp_verified ON users(otp_verified);
