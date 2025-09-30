-- ============================
-- Schema Inconsistency Fixes
-- ============================

-- PROBLEM 1: Duplicate users table in workout schema
-- The workout schema (005_create_workout_schema.sql) creates a users table
-- but users already exist in migration 001. This creates conflicts.
-- 
-- SOLUTION: Drop the duplicate users table and fix references

-- First, update the weekly_schemas table to reference the correct users table
ALTER TABLE weekly_schemas DROP CONSTRAINT IF EXISTS weekly_schemas_user_id_fkey;
ALTER TABLE weekly_schemas ADD CONSTRAINT weekly_schemas_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Update the user_id column type to match the auth users table
ALTER TABLE weekly_schemas ALTER COLUMN user_id TYPE TEXT;

-- Update progress_logs to reference correct users table  
ALTER TABLE progress_logs DROP CONSTRAINT IF EXISTS progress_logs_user_id_fkey;
ALTER TABLE progress_logs ALTER COLUMN user_id TYPE TEXT;
ALTER TABLE progress_logs ADD CONSTRAINT progress_logs_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Drop the duplicate users table from workout schema if it exists
DROP TABLE IF EXISTS users CASCADE;

-- PROBLEM 2: Inconsistent user ID references
-- Some tables use INT user_id, others use TEXT auth_user_id
-- Let's standardize on TEXT user_id to match the auth system

-- Update any remaining INT user_id columns to TEXT
-- (This would need to be done table by table with data migration if there's existing data)

-- PROBLEM 3: Missing enum types for consistency
-- Create enum types to match Go constants for better data integrity

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'fitness_level') THEN
    CREATE TYPE fitness_level AS ENUM ('beginner', 'intermediate', 'advanced');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'fitness_goal') THEN
    CREATE TYPE fitness_goal AS ENUM ('strength', 'muscle_gain', 'fat_loss', 'endurance', 'general_fitness');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exercise_type') THEN
    CREATE TYPE exercise_type AS ENUM ('strength', 'cardio', 'mobility', 'hiit', 'stretching');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'equipment_type') THEN
    CREATE TYPE equipment_type AS ENUM ('barbell', 'dumbbell', 'bodyweight', 'machine', 'kettlebell', 'resistance_band');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'session_status') THEN
    CREATE TYPE session_status AS ENUM ('active', 'completed', 'skipped', 'abandoned');
  END IF;
END $$;

-- Update existing tables to use enum types (if no data exists yet)
-- These commands should be run carefully with existing data

-- ALTER TABLE exercises 
--   ALTER COLUMN difficulty TYPE fitness_level USING difficulty::fitness_level,
--   ALTER COLUMN equipment TYPE equipment_type USING equipment::equipment_type,
--   ALTER COLUMN type TYPE exercise_type USING type::exercise_type;

-- PROBLEM 4: Missing constraints and validations
-- Add missing check constraints that match Go validations

-- Add check constraints to existing tables
ALTER TABLE exercises 
  ADD CONSTRAINT check_exercises_default_sets 
    CHECK (default_sets >= 1 AND default_sets <= 10);

ALTER TABLE exercises 
  ADD CONSTRAINT check_exercises_rest_seconds 
    CHECK (rest_seconds >= 0 AND rest_seconds <= 600);

-- Add constraints to workout_exercises
ALTER TABLE workout_exercises 
  ADD CONSTRAINT check_workout_exercises_sets 
    CHECK (sets >= 1 AND sets <= 10);

ALTER TABLE workout_exercises 
  ADD CONSTRAINT check_workout_exercises_rest_seconds 
    CHECK (rest_seconds >= 0 AND rest_seconds <= 600);

-- Add constraints to workout_templates
ALTER TABLE workout_templates 
  ADD CONSTRAINT check_workout_templates_days_per_week 
    CHECK (days_per_week >= 1 AND days_per_week <= 7);

-- PROBLEM 5: Missing triggers for updated_at timestamps
-- Add updated_at columns and triggers where needed

ALTER TABLE exercises ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE workout_templates ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Create triggers for updated_at
CREATE TRIGGER set_timestamp_exercises
    BEFORE UPDATE ON exercises
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_workout_templates
    BEFORE UPDATE ON workout_templates
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_workouts
    BEFORE UPDATE ON workouts
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_workout_exercises
    BEFORE UPDATE ON workout_exercises
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();