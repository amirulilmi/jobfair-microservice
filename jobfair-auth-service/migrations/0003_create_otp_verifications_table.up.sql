-- Create otp_verifications table
CREATE TABLE IF NOT EXISTS otp_verifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    purpose VARCHAR(50) NOT NULL CHECK (purpose IN ('phone_verification', 'password_reset', 'login_verification')),
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Key
    CONSTRAINT fk_otp_verifications_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- Create indexes for fast OTP lookup
CREATE INDEX idx_otp_verifications_user_id ON otp_verifications(user_id);
CREATE INDEX idx_otp_verifications_phone_number ON otp_verifications(phone_number);
CREATE INDEX idx_otp_verifications_purpose ON otp_verifications(purpose);
CREATE INDEX idx_otp_verifications_expires_at ON otp_verifications(expires_at);
CREATE INDEX idx_otp_verifications_is_used ON otp_verifications(is_used);

-- Composite index for common query pattern
CREATE INDEX idx_otp_verifications_lookup 
ON otp_verifications(phone_number, purpose, is_used, expires_at);

-- Comments
COMMENT ON TABLE otp_verifications IS 'OTP codes for phone and security verification';
COMMENT ON COLUMN otp_verifications.purpose IS 'Purpose of OTP: phone_verification, password_reset, or login_verification';
COMMENT ON COLUMN otp_verifications.expires_at IS 'OTP expiration timestamp (usually 5 minutes)';
