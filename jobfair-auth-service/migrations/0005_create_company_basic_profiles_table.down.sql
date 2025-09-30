-- Drop trigger
DROP TRIGGER IF EXISTS update_company_basic_profiles_updated_at ON company_basic_profiles;

-- Drop indexes
DROP INDEX IF EXISTS idx_company_basic_profiles_user_id;
DROP INDEX IF EXISTS idx_company_basic_profiles_company_name;
DROP INDEX IF EXISTS idx_company_basic_profiles_industry;
DROP INDEX IF EXISTS idx_company_basic_profiles_created_at;

-- Drop table
DROP TABLE IF EXISTS company_basic_profiles;