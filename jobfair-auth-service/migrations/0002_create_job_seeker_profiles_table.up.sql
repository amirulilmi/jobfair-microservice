-- Create job_seeker_profiles table
CREATE TABLE IF NOT EXISTS job_seeker_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    
    -- Employment Information
    current_job_title VARCHAR(255),
    current_company VARCHAR(255),
    employment_status VARCHAR(50) CHECK (employment_status IN ('student', 'employed', 'self_employed', 'unemployed')),
    job_search_status VARCHAR(50) CHECK (job_search_status IN ('actively_looking', 'open_to_opportunities', 'not_looking')),
    
    -- Job Preferences (PostgreSQL Arrays)
    desired_positions TEXT[],
    preferred_locations TEXT[],
    job_types TEXT[],
    
    -- Preferences & Settings
    notifications_enabled BOOLEAN DEFAULT TRUE,
    location_enabled BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Key
    CONSTRAINT fk_job_seeker_profiles_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    -- Unique constraint
    CONSTRAINT job_seeker_profiles_user_id_unique UNIQUE (user_id)
);

-- Create indexes
CREATE INDEX idx_job_seeker_profiles_user_id ON job_seeker_profiles(user_id);
CREATE INDEX idx_job_seeker_profiles_employment_status ON job_seeker_profiles(employment_status);
CREATE INDEX idx_job_seeker_profiles_job_search_status ON job_seeker_profiles(job_search_status);
CREATE INDEX idx_job_seeker_profiles_created_at ON job_seeker_profiles(created_at);

-- GIN indexes for array columns (for faster array searches)
CREATE INDEX idx_job_seeker_profiles_desired_positions ON job_seeker_profiles USING GIN (desired_positions);
CREATE INDEX idx_job_seeker_profiles_preferred_locations ON job_seeker_profiles USING GIN (preferred_locations);
CREATE INDEX idx_job_seeker_profiles_job_types ON job_seeker_profiles USING GIN (job_types);

-- Create trigger for updated_at
CREATE TRIGGER update_job_seeker_profiles_updated_at BEFORE UPDATE ON job_seeker_profiles
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE job_seeker_profiles IS 'Extended profile information for job seekers';
COMMENT ON COLUMN job_seeker_profiles.desired_positions IS 'Array of desired job positions';
COMMENT ON COLUMN job_seeker_profiles.preferred_locations IS 'Array of preferred work locations';
COMMENT ON COLUMN job_seeker_profiles.job_types IS 'Array of preferred job types (fulltime, parttime, etc)';
