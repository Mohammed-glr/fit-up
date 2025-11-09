-- Add notes column to progress_logs table for workout tracking
ALTER TABLE progress_logs 
ADD COLUMN IF NOT EXISTS notes TEXT;

COMMENT ON COLUMN progress_logs.notes IS 'Optional notes about the set (form, feeling, observations)';
