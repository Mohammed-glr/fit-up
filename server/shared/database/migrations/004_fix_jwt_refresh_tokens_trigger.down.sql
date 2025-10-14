-- Drop the trigger
DROP TRIGGER IF EXISTS set_timestamp_jwt_refresh_tokens ON jwt_refresh_tokens;

-- Drop the function
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- Remove the updated_at column
ALTER TABLE jwt_refresh_tokens DROP COLUMN IF EXISTS updated_at;
