-- Create company_analytics table (JF-101 analytics dashboard)
CREATE TABLE IF NOT EXISTS company_analytics (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,

    -- Virtual Booth Metrics
    booth_visits INTEGER DEFAULT 0,
    profile_views INTEGER DEFAULT 0,

    -- Job Metrics
    job_views INTEGER DEFAULT 0,
    applications INTEGER DEFAULT 0,
    total_jobs_posted INTEGER DEFAULT 0,
    active_jobs INTEGER DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Key
    CONSTRAINT fk_company_analytics_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,

    -- Unique constraint (one record per company)
    CONSTRAINT company_analytics_company_unique UNIQUE (company_id)
);

-- Create indexes
CREATE INDEX idx_company_analytics_company_id ON company_analytics(company_id);
CREATE INDEX idx_company_analytics_created_at ON company_analytics(created_at);

-- Create trigger for updated_at
CREATE TRIGGER update_company_analytics_updated_at BEFORE UPDATE ON company_analytics
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE company_analytics IS 'Analytics metrics for virtual booth and company profile';
COMMENT ON COLUMN company_analytics.booth_visits IS 'Total visitors to virtual booth';
COMMENT ON COLUMN company_analytics.profile_views IS 'Total profile views';
