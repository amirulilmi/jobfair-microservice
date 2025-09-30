-- 1. Create function for auto updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,

    -- Basic Information
    name VARCHAR(255) NOT NULL,
    description TEXT,
    industry VARCHAR(100),
    company_size VARCHAR(50) CHECK (company_size IN ('1-10', '11-50', '51-200', '201-500', '501-1000', '1000+')),
    founded_year INTEGER,

    -- Contact Information
    email VARCHAR(255),
    phone VARCHAR(50),
    website VARCHAR(255),

    -- Location
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Virtual Booth Assets (JF-101 requirement)
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    video_urls TEXT[],
    gallery_urls TEXT[],

    -- Social Media
    linkedin_url VARCHAR(255),
    facebook_url VARCHAR(255),
    twitter_url VARCHAR(255),
    instagram_url VARCHAR(255),

    -- Verification & Status
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    verification_badge VARCHAR(50),
    is_featured BOOLEAN DEFAULT FALSE,
    is_premium BOOLEAN DEFAULT FALSE,
    subscription_tier VARCHAR(50) DEFAULT 'free',

    -- SEO & Metadata
    slug VARCHAR(255),
    meta_title VARCHAR(255),
    meta_description TEXT,
    tags TEXT[],

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT companies_user_id_unique UNIQUE (user_id),
    CONSTRAINT companies_slug_unique UNIQUE (slug)
);

-- 3. Create indexes
CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_industry ON companies(industry);
CREATE INDEX idx_companies_city ON companies(city);
CREATE INDEX idx_companies_country ON companies(country);
CREATE INDEX idx_companies_is_verified ON companies(is_verified);
CREATE INDEX idx_companies_is_featured ON companies(is_featured);
CREATE INDEX idx_companies_slug ON companies(slug);
CREATE INDEX idx_companies_deleted_at ON companies(deleted_at);
CREATE INDEX idx_companies_created_at ON companies(created_at);

-- 4. GIN indexes for array columns
CREATE INDEX idx_companies_video_urls ON companies USING GIN (video_urls);
CREATE INDEX idx_companies_tags ON companies USING GIN (tags);

-- 5. Full-text search index
CREATE INDEX idx_companies_name_fulltext ON companies USING GIN (to_tsvector('english', name));
CREATE INDEX idx_companies_description_fulltext ON companies USING GIN (to_tsvector('english', description));

-- 6. Create trigger for updated_at
CREATE TRIGGER update_companies_updated_at
BEFORE UPDATE ON companies
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- 7. Comments
COMMENT ON TABLE companies IS 'Company profiles and virtual booth information';
COMMENT ON COLUMN companies.user_id IS 'Reference to users table in auth service (loose coupling)';
COMMENT ON COLUMN companies.video_urls IS 'Array of video URLs for virtual booth';
COMMENT ON COLUMN companies.slug IS 'URL-friendly company identifier';
