-- Drop trigger
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop indexes
DROP INDEX IF EXISTS idx_companies_user_id;
DROP INDEX IF EXISTS idx_companies_name;
DROP INDEX IF EXISTS idx_companies_industry;
DROP INDEX IF EXISTS idx_companies_city;
DROP INDEX IF EXISTS idx_companies_country;
DROP INDEX IF EXISTS idx_companies_is_verified;
DROP INDEX IF EXISTS idx_companies_is_featured;
DROP INDEX IF EXISTS idx_companies_slug;
DROP INDEX IF EXISTS idx_companies_deleted_at;
DROP INDEX IF EXISTS idx_companies_created_at;
DROP INDEX IF EXISTS idx_companies_video_urls;
DROP INDEX IF EXISTS idx_companies_tags;
DROP INDEX IF EXISTS idx_companies_name_fulltext;
DROP INDEX IF EXISTS idx_companies_description_fulltext;

-- Drop table
DROP TABLE IF EXISTS companies;
