-- Drop trigger
DROP TRIGGER IF EXISTS update_job_seeker_profiles_updated_at ON job_seeker_profiles;

-- Drop indexes
DROP INDEX IF EXISTS idx_job_seeker_profiles_user_id;
DROP INDEX IF EXISTS idx_job_seeker_profiles_employment_status;
DROP INDEX IF EXISTS idx_job_seeker_profiles_job_search_status;
DROP INDEX IF EXISTS idx_job_seeker_profiles_created_at;
DROP INDEX IF EXISTS idx_job_seeker_profiles_desired_positions;
DROP INDEX IF EXISTS idx_job_seeker_profiles_preferred_locations;
DROP INDEX IF EXISTS idx_job_seeker_profiles_job_types;

-- Drop table
DROP TABLE IF EXISTS job_seeker_profiles;