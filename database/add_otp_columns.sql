-- Add OTP columns to existing users table
-- Run this if the users table already exists without OTP columns

ALTER TABLE users
ADD COLUMN IF NOT EXISTS otp_code VARCHAR(6),
ADD COLUMN IF NOT EXISTS otp_expires TIMESTAMP,
ADD COLUMN IF NOT EXISTS otp_verified BOOLEAN DEFAULT FALSE;

-- Create index on otp_verified if not exists
CREATE INDEX IF NOT EXISTS idx_users_otp_verified ON users(otp_verified);
