-- Create job_applications table
CREATE TABLE IF NOT EXISTS job_applications (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    company_id INTEGER NOT NULL,

    -- Application Status
    status VARCHAR(50) DEFAULT 'applied' CHECK (status IN ('applied', 'shortlisted', 'interview', 'hired', 'rejected')),

    -- Application Details
    cover_letter TEXT,
    resume_url VARCHAR(500),

    -- Additional Information
    notes TEXT,

    -- Timestamps
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    viewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_job_applications_job_id
        FOREIGN KEY (job_id)
        REFERENCES jobs(id)
        ON DELETE CASCADE,

    -- Unique constraint: one user can only apply once to a job
    CONSTRAINT job_applications_user_job_unique UNIQUE (job_id, user_id)
);

-- Create indexes
CREATE INDEX idx_job_applications_job_id ON job_applications(job_id);
CREATE INDEX idx_job_applications_user_id ON job_applications(user_id);
CREATE INDEX idx_job_applications_company_id ON job_applications(company_id);
CREATE INDEX idx_job_applications_status ON job_applications(status);
CREATE INDEX idx_job_applications_applied_at ON job_applications(applied_at);
CREATE INDEX idx_job_applications_created_at ON job_applications(created_at);

-- Create trigger for updated_at
CREATE TRIGGER update_job_applications_updated_at
BEFORE UPDATE ON job_applications
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE job_applications IS 'Job applications from users to company job postings';
COMMENT ON COLUMN job_applications.status IS 'Application status: applied, shortlisted, interview, hired, rejected';
COMMENT ON COLUMN job_applications.user_id IS 'Reference to users table in auth service';
COMMENT ON COLUMN job_applications.viewed_at IS 'When the company first viewed this application';