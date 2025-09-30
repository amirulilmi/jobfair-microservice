-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('job_seeker', 'company', 'admin')),
    
    -- Basic Profile
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone_number VARCHAR(50),
    country_code VARCHAR(10) DEFAULT '+62',
    country VARCHAR(100),
    profile_photo VARCHAR(500),
    
    -- Verification Status
    is_email_verified BOOLEAN DEFAULT FALSE,
    is_phone_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP,
    phone_verified_at TIMESTAMP,
    
    -- Account Status
    is_active BOOLEAN DEFAULT TRUE,
    is_profile_complete BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT users_email_unique UNIQUE (email)
);

-- Create partial unique index for phone_number (only non-null values)
CREATE UNIQUE INDEX idx_users_phone_number_unique 
ON users (phone_number) 
WHERE phone_number IS NOT NULL AND phone_number != '';

-- Create indexes for better query performance
CREATE INDEX idx_users_user_type ON users(user_type);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE users IS 'Main users table for authentication';
COMMENT ON COLUMN users.user_type IS 'User type: job_seeker, company, or admin';
COMMENT ON COLUMN users.phone_number IS 'Unique phone number, nullable';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp';