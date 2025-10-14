-- Drop the problematic trigger
DROP TRIGGER IF EXISTS set_timestamp_jwt_refresh_tokens ON jwt_refresh_tokens;

-- Add updated_at column to jwt_refresh_tokens if it doesn't exist
ALTER TABLE jwt_refresh_tokens ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Create the trigger function if it doesn't exist
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate the trigger
CREATE TRIGGER set_timestamp_jwt_refresh_tokens
    BEFORE UPDATE ON jwt_refresh_tokens
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
