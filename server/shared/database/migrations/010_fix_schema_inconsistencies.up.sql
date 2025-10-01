
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


DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_exercises_default_sets') THEN
    ALTER TABLE exercises ADD CONSTRAINT check_exercises_default_sets 
      CHECK (default_sets >= 1 AND default_sets <= 10);
  END IF;
END $$;

DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_exercises_rest_seconds') THEN
    ALTER TABLE exercises ADD CONSTRAINT check_exercises_rest_seconds 
      CHECK (rest_seconds >= 0 AND rest_seconds <= 600);
  END IF;
END $$;

DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_workout_exercises_sets') THEN
    ALTER TABLE workout_exercises ADD CONSTRAINT check_workout_exercises_sets 
      CHECK (sets >= 1 AND sets <= 10);
  END IF;
END $$;

DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_workout_exercises_rest_seconds') THEN
    ALTER TABLE workout_exercises ADD CONSTRAINT check_workout_exercises_rest_seconds 
      CHECK (rest_seconds >= 0 AND rest_seconds <= 600);
  END IF;
END $$;

DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_workout_templates_days_per_week') THEN
    ALTER TABLE workout_templates ADD CONSTRAINT check_workout_templates_days_per_week 
      CHECK (days_per_week >= 1 AND days_per_week <= 7);
  END IF;
END $$;


ALTER TABLE exercises ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE workout_templates ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

DROP TRIGGER IF EXISTS set_timestamp_exercises ON exercises;
CREATE TRIGGER set_timestamp_exercises
    BEFORE UPDATE ON exercises
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

DROP TRIGGER IF EXISTS set_timestamp_workout_templates ON workout_templates;
CREATE TRIGGER set_timestamp_workout_templates
    BEFORE UPDATE ON workout_templates
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

DROP TRIGGER IF EXISTS set_timestamp_workouts ON workouts;
CREATE TRIGGER set_timestamp_workouts
    BEFORE UPDATE ON workouts
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

DROP TRIGGER IF EXISTS set_timestamp_workout_exercises ON workout_exercises;
CREATE TRIGGER set_timestamp_workout_exercises
    BEFORE UPDATE ON workout_exercises
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();