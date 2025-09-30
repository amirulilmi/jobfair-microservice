DROP INDEX IF EXISTS idx_jobs_description_fulltext;
DROP INDEX IF EXISTS idx_jobs_title_fulltext;
DROP INDEX IF EXISTS idx_jobs_benefits;
DROP INDEX IF EXISTS idx_jobs_skills;
DROP INDEX IF EXISTS idx_jobs_created_at;
DROP INDEX IF EXISTS idx_jobs_deleted_at;
DROP INDEX IF EXISTS idx_jobs_published_at;
DROP INDEX IF EXISTS idx_jobs_slug;
DROP INDEX IF EXISTS idx_jobs_is_remote;
DROP INDEX IF EXISTS idx_jobs_country;
DROP INDEX IF EXISTS idx_jobs_city;
DROP INDEX IF EXISTS idx_jobs_job_type;
DROP INDEX IF EXISTS idx_jobs_status;
DROP INDEX IF EXISTS idx_jobs_company_id;

DROP TRIGGER IF EXISTS update_jobs_updated_at ON jobs;

DROP TABLE IF EXISTS jobs;