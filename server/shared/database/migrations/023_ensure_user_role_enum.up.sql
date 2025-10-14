-- Ensure the user_role enum has all required values
DO $$ 
BEGIN
    -- Check if 'coach' value exists in the enum, if not add it
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'coach' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'user_role')
    ) THEN
        ALTER TYPE user_role ADD VALUE 'coach';
    END IF;
    
    -- Check if 'admin' value exists in the enum, if not add it
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumlabel = 'admin' 
        AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'user_role')
    ) THEN
        ALTER TYPE user_role ADD VALUE 'admin';
    END IF;
END $$;
