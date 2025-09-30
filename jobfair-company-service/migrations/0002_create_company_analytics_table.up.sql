-- Create company_analytics table (JF-101 analytics dashboard)
CREATE TABLE IF NOT EXISTS company_analytics (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    
    -- Virtual Booth Metrics
    booth_visits INTEGER DEFAULT 0,
    booth_unique_visits INTEGER DEFAULT 0,
    profile_views INTEGER DEFAULT 0,
    profile_unique_views INTEGER DEFAULT 0,
    
    -- Job Metrics
    job_views INTEGER DEFAULT 0,
    job_applications INTEGER DEFAULT 0,
    job_applications_qualified INTEGER DEFAULT 0,
    
    -- Engagement Metrics
    video_plays INTEGER DEFAULT 0,
    gallery_views INTEGER DEFAULT 0,
    website_clicks INTEGER DEFAULT 0,
    social_clicks INTEGER DEFAULT 0,
    contact_clicks INTEGER DEFAULT 0,
    
    -- Conversion Metrics
    saved_by_users INTEGER DEFAULT 0,
    shared_count INTEGER DEFAULT 0,
    average_time_on_booth INTEGER DEFAULT 0, -- in seconds
    
    -- Period Metrics
    analytics_date DATE NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Key
    CONSTRAINT fk_company_analytics_company_id 
        FOREIGN KEY (company_id) 
        REFERENCES companies(id) 
        ON DELETE CASCADE,
    
    -- Unique constraint (one record per company per day)
    CONSTRAINT company_analytics_company_date_unique UNIQUE (company_id, analytics_date)
);

-- Create indexes
CREATE INDEX idx_company_analytics_company_id ON company_analytics(company_id);
CREATE INDEX idx_company_analytics_date ON company_analytics(analytics_date);
CREATE INDEX idx_company_analytics_booth_visits ON company_analytics(booth_visits);
CREATE INDEX idx_company_analytics_created_at ON company_analytics(created_at);

-- Create trigger for updated_at
CREATE TRIGGER update_company_analytics_updated_at BEFORE UPDATE ON company_analytics
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE company_analytics IS 'Daily analytics metrics for virtual booth and company profile';
COMMENT ON COLUMN company_analytics.analytics_date IS 'Date for these analytics (one record per day)';
COMMENT ON COLUMN company_analytics.booth_unique_visits IS 'Unique visitors to virtual booth';
