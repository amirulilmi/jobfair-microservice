DROP INDEX IF EXISTS idx_job_applications_created_at;
DROP INDEX IF EXISTS idx_job_applications_applied_at;
DROP INDEX IF EXISTS idx_job_applications_status;
DROP INDEX IF EXISTS idx_job_applications_company_id;
DROP INDEX IF EXISTS idx_job_applications_user_id;
DROP INDEX IF EXISTS idx_job_applications_job_id;

DROP TRIGGER IF EXISTS update_job_applications_updated_at ON job_applications;

DROP TABLE IF EXISTS job_applications;