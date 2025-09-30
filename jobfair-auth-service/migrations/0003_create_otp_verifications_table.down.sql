-- Drop indexes
DROP INDEX IF EXISTS idx_otp_verifications_user_id;
DROP INDEX IF EXISTS idx_otp_verifications_phone_number;
DROP INDEX IF EXISTS idx_otp_verifications_purpose;
DROP INDEX IF EXISTS idx_otp_verifications_expires_at;
DROP INDEX IF EXISTS idx_otp_verifications_is_used;
DROP INDEX IF EXISTS idx_otp_verifications_lookup;

-- Drop table
DROP TABLE IF EXISTS otp_verifications;