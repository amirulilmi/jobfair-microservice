-- Drop trigger
DROP TRIGGER IF EXISTS update_company_media_updated_at ON company_media;

-- Drop indexes
DROP INDEX IF EXISTS idx_company_media_company_id;
DROP INDEX IF EXISTS idx_company_media_type;
DROP INDEX IF EXISTS idx_company_media_display_order;
DROP INDEX IF EXISTS idx_company_media_is_public;
DROP INDEX IF EXISTS idx_company_media_deleted_at;

-- Drop table
DROP TABLE IF EXISTS company_media;