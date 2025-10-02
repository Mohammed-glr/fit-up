
CREATE TABLE workout_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id INT NOT NULL REFERENCES workouts(workout_id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'skipped', 'abandoned')),
    total_exercises INT NOT NULL DEFAULT 0,
    completed_exercises INT NOT NULL DEFAULT 0,
    total_volume FLOAT NOT NULL DEFAULT 0,
    notes TEXT DEFAULT ''
);

CREATE TABLE skipped_workouts (
    skip_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id INT NOT NULL REFERENCES workouts(workout_id) ON DELETE CASCADE,
    skip_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reason TEXT NOT NULL
);

CREATE TABLE exercise_performances (
    performance_id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES workout_sessions(session_id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    sets_completed INT NOT NULL DEFAULT 0,
    total_volume FLOAT NOT NULL DEFAULT 0,
    rpe FLOAT DEFAULT NULL CHECK (rpe IS NULL OR (rpe >= 1 AND rpe <= 10)),
    notes TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE set_performances (
    set_id SERIAL PRIMARY KEY,
    performance_id INT NOT NULL REFERENCES exercise_performances(performance_id) ON DELETE CASCADE,
    set_number INT NOT NULL,
    reps INT NOT NULL CHECK (reps >= 0),
    weight FLOAT NOT NULL DEFAULT 0 CHECK (weight >= 0),
    rpe FLOAT DEFAULT NULL CHECK (rpe IS NULL OR (rpe >= 1 AND rpe <= 10)),
    rest_seconds INT DEFAULT NULL CHECK (rest_seconds IS NULL OR rest_seconds >= 0),
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(performance_id, set_number)
);

CREATE TABLE session_metrics (
    metric_id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES workout_sessions(session_id) ON DELETE CASCADE,
    duration_seconds INT NOT NULL CHECK (duration_seconds >= 0),
    total_volume FLOAT NOT NULL DEFAULT 0,
    average_intensity FLOAT DEFAULT NULL CHECK (average_intensity IS NULL OR (average_intensity >= 0 AND average_intensity <= 1)),
    completion_rate FLOAT NOT NULL CHECK (completion_rate >= 0 AND completion_rate <= 1),
    average_rpe FLOAT DEFAULT NULL CHECK (average_rpe IS NULL OR (average_rpe >= 1 AND average_rpe <= 10)),
    calories_burned INT DEFAULT NULL CHECK (calories_burned IS NULL OR calories_burned >= 0),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id)
);

CREATE TABLE weekly_session_stats (
    stat_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    sessions_planned INT NOT NULL DEFAULT 0,
    sessions_completed INT NOT NULL DEFAULT 0,
    total_volume FLOAT NOT NULL DEFAULT 0,
    average_rpe FLOAT DEFAULT NULL CHECK (average_rpe IS NULL OR (average_rpe >= 1 AND average_rpe <= 10)),
    completion_rate FLOAT NOT NULL CHECK (completion_rate >= 0 AND completion_rate <= 1),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, week_start)
);

CREATE INDEX idx_workout_sessions_user_id ON workout_sessions(user_id);
CREATE INDEX idx_workout_sessions_status ON workout_sessions(status);
CREATE INDEX idx_workout_sessions_start_time ON workout_sessions(start_time);
CREATE INDEX idx_workout_sessions_user_date ON workout_sessions(user_id, start_time);

CREATE INDEX idx_skipped_workouts_user_id ON skipped_workouts(user_id);
CREATE INDEX idx_skipped_workouts_date ON skipped_workouts(skip_date);

CREATE INDEX idx_exercise_performances_session_id ON exercise_performances(session_id);
CREATE INDEX idx_exercise_performances_exercise_id ON exercise_performances(exercise_id);

CREATE INDEX idx_set_performances_performance_id ON set_performances(performance_id);
CREATE INDEX idx_set_performances_completed_at ON set_performances(completed_at);

CREATE INDEX idx_session_metrics_session_id ON session_metrics(session_id);

CREATE INDEX idx_weekly_session_stats_user_id ON weekly_session_stats(user_id);
CREATE INDEX idx_weekly_session_stats_week_start ON weekly_session_stats(week_start);