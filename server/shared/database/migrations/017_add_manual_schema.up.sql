ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS coach_id TEXT;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS created_by TEXT DEFAULT 'system';
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS is_custom BOOLEAN DEFAULT FALSE;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS base_template_id INTEGER;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS metadata JSONB;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;

ALTER TABLE workouts ADD COLUMN IF NOT EXISTS workout_name TEXT;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS estimated_minutes INTEGER;

ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS weight TEXT;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS tempo TEXT;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS order_index INTEGER DEFAULT 0;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS is_superset BOOLEAN DEFAULT FALSE;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS superset_group INTEGER;

CREATE INDEX idx_weekly_schemas_coach_id ON weekly_schemas(coach_id);
CREATE INDEX idx_weekly_schemas_created_by ON weekly_schemas(created_by);

