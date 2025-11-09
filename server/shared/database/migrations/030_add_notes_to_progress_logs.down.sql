-- Remove notes column from progress_logs table
ALTER TABLE progress_logs 
DROP COLUMN IF EXISTS notes;
