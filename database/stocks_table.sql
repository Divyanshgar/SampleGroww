-- Create stocks table for PostgreSQL
-- This table stores stock transactions for users

CREATE TABLE IF NOT EXISTS stocks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_name VARCHAR(100) NOT NULL,
    action VARCHAR(10) NOT NULL,  -- 'buy' or 'sell'
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    original_value DECIMAL(10,2) NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_stocks_user_id ON stocks(user_id);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_stocks_deleted_at ON stocks(deleted_at);
