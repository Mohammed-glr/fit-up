
DROP TRIGGER IF EXISTS set_timestamp_workout_exercises ON workout_exercises;
DROP TRIGGER IF EXISTS set_timestamp_workouts ON workouts;
DROP TRIGGER IF EXISTS set_timestamp_workout_templates ON workout_templates;
DROP TRIGGER IF EXISTS set_timestamp_exercises ON exercises;

ALTER TABLE workout_exercises DROP COLUMN IF EXISTS updated_at;
ALTER TABLE workouts DROP COLUMN IF EXISTS updated_at;
ALTER TABLE workout_templates DROP COLUMN IF EXISTS updated_at;
ALTER TABLE exercises DROP COLUMN IF EXISTS updated_at;

ALTER TABLE workout_templates DROP CONSTRAINT IF EXISTS check_workout_templates_days_per_week;
ALTER TABLE workout_exercises DROP CONSTRAINT IF EXISTS check_workout_exercises_rest_seconds;
ALTER TABLE workout_exercises DROP CONSTRAINT IF EXISTS check_workout_exercises_sets;
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS check_exercises_rest_seconds;
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS check_exercises_default_sets;

DROP TYPE IF EXISTS session_status CASCADE;
DROP TYPE IF EXISTS equipment_type CASCADE;
DROP TYPE IF EXISTS exercise_type CASCADE;
DROP TYPE IF EXISTS fitness_goal CASCADE;
DROP TYPE IF EXISTS fitness_level CASCADE;
