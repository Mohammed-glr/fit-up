-- Create workout_templates table for saving workout configurations
CREATE TABLE IF NOT EXISTS workout_templates (
    template_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    exercises JSONB NOT NULL, -- Array of exercises with sets, reps, weight
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster user template lookups
CREATE INDEX idx_workout_templates_user_id ON workout_templates(user_id);

-- Create index for public templates discovery
CREATE INDEX idx_workout_templates_public ON workout_templates(is_public) WHERE is_public = TRUE;

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_workout_template_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER workout_template_updated_at_trigger
    BEFORE UPDATE ON workout_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_workout_template_updated_at();

-- Insert some default templates for reference
INSERT INTO workout_templates (user_id, name, description, is_public, exercises) VALUES
(
    '00000000-0000-0000-0000-000000000000', -- System user placeholder
    'Push Day (Chest, Shoulders, Triceps)',
    'A complete push workout focusing on chest, shoulders, and triceps',
    TRUE,
    '[
        {"exercise_name": "Bench Press", "sets": 4, "target_reps": 8, "target_weight": 0, "rest_seconds": 120},
        {"exercise_name": "Incline Dumbbell Press", "sets": 3, "target_reps": 10, "target_weight": 0, "rest_seconds": 90},
        {"exercise_name": "Shoulder Press", "sets": 4, "target_reps": 8, "target_weight": 0, "rest_seconds": 90},
        {"exercise_name": "Lateral Raises", "sets": 3, "target_reps": 12, "target_weight": 0, "rest_seconds": 60},
        {"exercise_name": "Tricep Dips", "sets": 3, "target_reps": 10, "target_weight": 0, "rest_seconds": 60}
    ]'::jsonb
),
(
    '00000000-0000-0000-0000-000000000000',
    'Pull Day (Back, Biceps)',
    'A complete pull workout focusing on back and biceps',
    TRUE,
    '[
        {"exercise_name": "Deadlift", "sets": 4, "target_reps": 6, "target_weight": 0, "rest_seconds": 180},
        {"exercise_name": "Pull-ups", "sets": 4, "target_reps": 8, "target_weight": 0, "rest_seconds": 120},
        {"exercise_name": "Barbell Rows", "sets": 4, "target_reps": 8, "target_weight": 0, "rest_seconds": 90},
        {"exercise_name": "Face Pulls", "sets": 3, "target_reps": 15, "target_weight": 0, "rest_seconds": 60},
        {"exercise_name": "Bicep Curls", "sets": 3, "target_reps": 10, "target_weight": 0, "rest_seconds": 60}
    ]'::jsonb
),
(
    '00000000-0000-0000-0000-000000000000',
    'Leg Day (Quads, Hamstrings, Glutes)',
    'A complete leg workout',
    TRUE,
    '[
        {"exercise_name": "Squats", "sets": 4, "target_reps": 8, "target_weight": 0, "rest_seconds": 180},
        {"exercise_name": "Romanian Deadlifts", "sets": 4, "target_reps": 8, "target_weight": 0, "rest_seconds": 120},
        {"exercise_name": "Leg Press", "sets": 3, "target_reps": 12, "target_weight": 0, "rest_seconds": 90},
        {"exercise_name": "Leg Curls", "sets": 3, "target_reps": 12, "target_weight": 0, "rest_seconds": 60},
        {"exercise_name": "Calf Raises", "sets": 4, "target_reps": 15, "target_weight": 0, "rest_seconds": 60}
    ]'::jsonb
),
(
    '00000000-0000-0000-0000-000000000000',
    'Full Body Beginner',
    'A beginner-friendly full body workout',
    TRUE,
    '[
        {"exercise_name": "Goblet Squats", "sets": 3, "target_reps": 10, "target_weight": 0, "rest_seconds": 90},
        {"exercise_name": "Push-ups", "sets": 3, "target_reps": 10, "target_weight": 0, "rest_seconds": 60},
        {"exercise_name": "Dumbbell Rows", "sets": 3, "target_reps": 10, "target_weight": 0, "rest_seconds": 60},
        {"exercise_name": "Overhead Press", "sets": 3, "target_reps": 8, "target_weight": 0, "rest_seconds": 90},
        {"exercise_name": "Plank", "sets": 3, "target_reps": 30, "target_weight": 0, "rest_seconds": 60}
    ]'::jsonb
);

COMMENT ON TABLE workout_templates IS 'Stores user-created workout templates for quick workout starts';
COMMENT ON COLUMN workout_templates.exercises IS 'JSONB array of exercise configurations with sets, reps, and weights';
COMMENT ON COLUMN workout_templates.is_public IS 'Whether the template can be discovered by other users';
