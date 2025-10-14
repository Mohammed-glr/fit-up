-- Cannot remove enum values in PostgreSQL without recreating the entire enum
-- This migration is effectively irreversible
-- If you need to roll back, you would need to:
-- 1. Change all users with 'coach' or 'admin' roles to 'user'
-- 2. Drop and recreate the enum
-- 3. Recreate the users table

-- For safety, this down migration does nothing
SELECT 1;
