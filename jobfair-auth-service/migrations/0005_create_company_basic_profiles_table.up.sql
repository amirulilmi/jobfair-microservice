-- File: migrations/0005_create_company_basic_profiles_table.up.sql

-- Create company_basic_profiles table (stored in auth service)
CREATE TABLE IF NOT EXISTS company_basic_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    
    -- Basic Company Information
    company_name VARCHAR(255) NOT NULL,
    industry VARCHAR(100),
    phone_number VARCHAR(50),
    address TEXT,
    website VARCHAR(255),
    logo_url VARCHAR(500),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Key
    CONSTRAINT fk_company_basic_profiles_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    -- Unique constraint
    CONSTRAINT company_basic_profiles_user_id_unique UNIQUE (user_id)
);

-- Create indexes
CREATE INDEX idx_company_basic_profiles_user_id ON company_basic_profiles(user_id);
CREATE INDEX idx_company_basic_profiles_company_name ON company_basic_profiles(company_name);
CREATE INDEX idx_company_basic_profiles_industry ON company_basic_profiles(industry);
CREATE INDEX idx_company_basic_profiles_created_at ON company_basic_profiles(created_at);

-- Create trigger for updated_at
CREATE TRIGGER update_company_basic_profiles_updated_at BEFORE UPDATE ON company_basic_profiles
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE company_basic_profiles IS 'Basic company profile information stored in auth service';
COMMENT ON COLUMN company_basic_profiles.user_id IS 'Reference to users table (one-to-one)';
COMMENT ON COLUMN company_basic_profiles.logo_url IS 'Company logo URL';