-- Create jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,

    -- Basic Information
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    requirements TEXT,
    responsibilities TEXT,

    -- Job Details
    job_type VARCHAR(50) NOT NULL CHECK (job_type IN ('full_time', 'part_time', 'contract', 'internship', 'freelance')),
    job_level VARCHAR(50) CHECK (job_level IN ('entry', 'junior', 'mid', 'senior', 'lead', 'manager', 'director', 'executive')),
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'closed', 'expired')),

    -- Location
    location VARCHAR(255),
    city VARCHAR(100),
    country VARCHAR(100),
    is_remote BOOLEAN DEFAULT FALSE,

    -- Compensation
    salary_min INTEGER,
    salary_max INTEGER,
    salary_currency VARCHAR(10) DEFAULT 'IDR',
    show_salary BOOLEAN DEFAULT FALSE,

    -- Skills and Benefits
    skills TEXT[],
    benefits TEXT[],

    -- Positions
    positions INTEGER DEFAULT 1,

    -- Metrics
    view_count INTEGER DEFAULT 0,
    application_count INTEGER DEFAULT 0,

    -- Dates
    expires_at TIMESTAMP,
    published_at TIMESTAMP,

    -- SEO
    slug VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign Key
    CONSTRAINT fk_jobs_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_jobs_company_id ON jobs(company_id);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_job_type ON jobs(job_type);
CREATE INDEX idx_jobs_city ON jobs(city);
CREATE INDEX idx_jobs_country ON jobs(country);
CREATE INDEX idx_jobs_is_remote ON jobs(is_remote);
CREATE INDEX idx_jobs_slug ON jobs(slug);
CREATE INDEX idx_jobs_published_at ON jobs(published_at);
CREATE INDEX idx_jobs_deleted_at ON jobs(deleted_at);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);

-- GIN indexes for array columns
CREATE INDEX idx_jobs_skills ON jobs USING GIN (skills);
CREATE INDEX idx_jobs_benefits ON jobs USING GIN (benefits);

-- Full-text search index
CREATE INDEX idx_jobs_title_fulltext ON jobs USING GIN (to_tsvector('english', title));
CREATE INDEX idx_jobs_description_fulltext ON jobs USING GIN (to_tsvector('english', description));

-- Create trigger for updated_at
CREATE TRIGGER update_jobs_updated_at
BEFORE UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE jobs IS 'Job postings created by companies';
COMMENT ON COLUMN jobs.status IS 'draft: not published, active: live, paused: temporarily hidden, closed: no longer accepting applications, expired: past expiry date';
COMMENT ON COLUMN jobs.job_type IS 'Type of employment';
COMMENT ON COLUMN jobs.is_remote IS 'Whether the job can be done remotely';