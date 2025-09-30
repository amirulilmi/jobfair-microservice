-- Drop trigger
DROP TRIGGER IF EXISTS update_company_analytics_updated_at ON company_analytics;

-- Drop indexes
DROP INDEX IF EXISTS idx_company_analytics_company_id;
DROP INDEX IF EXISTS idx_company_analytics_date;
DROP INDEX IF EXISTS idx_company_analytics_booth_visits;
DROP INDEX IF EXISTS idx_company_analytics_created_at;

-- Drop table
DROP TABLE IF EXISTS company_analytics;