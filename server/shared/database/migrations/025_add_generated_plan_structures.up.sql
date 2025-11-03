CREATE TABLE IF NOT EXISTS generated_plan_days (
    plan_day_id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES generated_plans(plan_id) ON DELETE CASCADE,
    day_index INT NOT NULL CHECK (day_index >= 1),
    day_title TEXT NOT NULL,
    focus TEXT,
    is_rest BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(plan_id, day_index)
);

CREATE TABLE IF NOT EXISTS generated_plan_exercises (
    plan_exercise_id SERIAL PRIMARY KEY,
    plan_day_id INT NOT NULL REFERENCES generated_plan_days(plan_day_id) ON DELETE CASCADE,
    exercise_order INT NOT NULL CHECK (exercise_order >= 1),
    exercise_id INT,
    name TEXT NOT NULL,
    sets INT,
    reps TEXT,
    rest_seconds INT,
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_generated_plan_days_plan_id ON generated_plan_days(plan_id);
CREATE INDEX IF NOT EXISTS idx_generated_plan_exercises_day_id ON generated_plan_exercises(plan_day_id);
