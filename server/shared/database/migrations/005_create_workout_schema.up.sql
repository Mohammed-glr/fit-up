CREATE TABLE exercises (
    exercise_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    muscle_groups VARCHAR(100) NOT NULL, 
    difficulty VARCHAR(20) NOT NULL,     
    equipment VARCHAR(50) NOT NULL,     
    type VARCHAR(20) NOT NULL,          
    default_sets INT NOT NULL,
    default_reps VARCHAR(20) NOT NULL,   
    rest_seconds INT NOT NULL
);

CREATE TABLE workout_templates (
    template_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    min_level VARCHAR(20) NOT NULL,
    max_level VARCHAR(20) NOT NULL,
    suitable_goals VARCHAR(100) NOT NULL,
    days_per_week INT NOT NULL
);

CREATE TABLE weekly_schemas (
    schema_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE workouts (
    workout_id SERIAL PRIMARY KEY,
    schema_id INT REFERENCES weekly_schemas(schema_id) ON DELETE CASCADE,
    day_of_week INT NOT NULL,
    focus VARCHAR(50) NOT NULL
);

CREATE TABLE workout_exercises (
    we_id SERIAL PRIMARY KEY,
    workout_id INT REFERENCES workouts(workout_id) ON DELETE CASCADE,
    exercise_id INT REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    sets INT NOT NULL,
    reps VARCHAR(20) NOT NULL,
    rest_seconds INT NOT NULL
);

CREATE TABLE progress_logs (
    log_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INT REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    date DATE NOT NULL,
    sets_completed INT,
    reps_completed INT,
    weight_used FLOAT,
    duration_seconds INT
);

CREATE INDEX idx_exercise_muscle_groups ON exercises USING GIN (to_tsvector('english', muscle_groups));
CREATE INDEX idx_workout_schema_user ON weekly_schemas(user_id, week_start);
CREATE INDEX idx_progress_logs_user_date ON progress_logs(user_id, date);
CREATE INDEX idx_workout_exercises_workout ON workout_exercises(workout_id);
CREATE INDEX idx_workouts_schema_day ON workouts(schema_id, day_of_week);
