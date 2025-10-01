
CREATE TABLE recovery_metrics (
    metric_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    sleep_hours FLOAT DEFAULT NULL CHECK (sleep_hours IS NULL OR (sleep_hours >= 0 AND sleep_hours <= 24)),
    sleep_quality FLOAT DEFAULT NULL CHECK (sleep_quality IS NULL OR (sleep_quality >= 0 AND sleep_quality <= 10)),
    stress_level FLOAT DEFAULT NULL CHECK (stress_level IS NULL OR (stress_level >= 0 AND stress_level <= 10)),
    energy_level FLOAT DEFAULT NULL CHECK (energy_level IS NULL OR (energy_level >= 0 AND energy_level <= 10)),
    soreness FLOAT DEFAULT NULL CHECK (soreness IS NULL OR (soreness >= 0 AND soreness <= 10)),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

CREATE TABLE strength_progressions (
    progression_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    starting_max FLOAT NOT NULL CHECK (starting_max > 0),
    current_max FLOAT NOT NULL CHECK (current_max > 0),
    progression_rate FLOAT NOT NULL DEFAULT 0,
    trend VARCHAR(20) NOT NULL DEFAULT 'stable' CHECK (trend IN ('increasing', 'decreasing', 'stable', 'plateau')),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    timeframe_days INT NOT NULL DEFAULT 30 CHECK (timeframe_days > 0),
    UNIQUE(user_id, exercise_id, timeframe_days)
);

CREATE TABLE plateau_detections (
    detection_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    plateau_detected BOOLEAN NOT NULL DEFAULT FALSE,
    plateau_duration_days INT DEFAULT NULL CHECK (plateau_duration_days IS NULL OR plateau_duration_days >= 0),
    recommendation TEXT DEFAULT NULL,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, exercise_id)
);

CREATE TABLE training_volumes (
    volume_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    total_sets INT NOT NULL DEFAULT 0,
    total_reps INT NOT NULL DEFAULT 0,
    total_weight FLOAT NOT NULL DEFAULT 0,
    volume_load FLOAT NOT NULL DEFAULT 0,
    intensity_average FLOAT DEFAULT NULL CHECK (intensity_average IS NULL OR (intensity_average >= 0 AND intensity_average <= 1)),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, week_start)
);

CREATE TABLE intensity_progressions (
    intensity_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    baseline_intensity FLOAT NOT NULL CHECK (baseline_intensity >= 0 AND baseline_intensity <= 1),
    current_intensity FLOAT NOT NULL CHECK (current_intensity >= 0 AND current_intensity <= 1),
    progression_rate FLOAT NOT NULL DEFAULT 0,
    recommended_next FLOAT DEFAULT NULL CHECK (recommended_next IS NULL OR (recommended_next >= 0 AND recommended_next <= 1)),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, exercise_id)
);

CREATE TABLE goal_progress (
    progress_id SERIAL PRIMARY KEY,
    goal_id INT NOT NULL REFERENCES fitness_goal_targets(goal_id) ON DELETE CASCADE,
    progress_percent FLOAT NOT NULL CHECK (progress_percent >= 0 AND progress_percent <= 100),
    on_track BOOLEAN NOT NULL DEFAULT TRUE,
    estimated_completion TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(goal_id, calculated_at::DATE)
);

CREATE TABLE goal_predictions (
    prediction_id SERIAL PRIMARY KEY,
    goal_id INT NOT NULL REFERENCES fitness_goal_targets(goal_id) ON DELETE CASCADE,
    probability_of_success FLOAT NOT NULL CHECK (probability_of_success >= 0 AND probability_of_success <= 1),
    estimated_days INT NOT NULL CHECK (estimated_days >= 0),
    confidence FLOAT NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    predicted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(goal_id)
);

CREATE TABLE goal_adjustments (
    adjustment_id SERIAL PRIMARY KEY,
    goal_id INT NOT NULL REFERENCES fitness_goal_targets(goal_id) ON DELETE CASCADE,
    recommendation_type VARCHAR(50) NOT NULL,
    adjustment TEXT NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE optimal_loads (
    load_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recommended_sets INT NOT NULL CHECK (recommended_sets > 0),
    recommended_reps INT NOT NULL CHECK (recommended_reps > 0),
    intensity_range VARCHAR(20) NOT NULL, -- e.g., "70-80%"
    volume_target FLOAT NOT NULL CHECK (volume_target > 0),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(user_id, calculated_at::DATE)
);

CREATE INDEX idx_recovery_metrics_user_id ON recovery_metrics(user_id);
CREATE INDEX idx_recovery_metrics_date ON recovery_metrics(date);
CREATE INDEX idx_recovery_metrics_user_date ON recovery_metrics(user_id, date);

CREATE INDEX idx_strength_progressions_user_id ON strength_progressions(user_id);
CREATE INDEX idx_strength_progressions_exercise_id ON strength_progressions(exercise_id);
CREATE INDEX idx_strength_progressions_user_exercise ON strength_progressions(user_id, exercise_id);

CREATE INDEX idx_plateau_detections_user_id ON plateau_detections(user_id);
CREATE INDEX idx_plateau_detections_detected ON plateau_detections(plateau_detected) WHERE plateau_detected = TRUE;

CREATE INDEX idx_training_volumes_user_id ON training_volumes(user_id);
CREATE INDEX idx_training_volumes_week_start ON training_volumes(week_start);

CREATE INDEX idx_intensity_progressions_user_id ON intensity_progressions(user_id);
CREATE INDEX idx_intensity_progressions_exercise_id ON intensity_progressions(exercise_id);

CREATE INDEX idx_goal_progress_goal_id ON goal_progress(goal_id);
CREATE INDEX idx_goal_progress_calculated_at ON goal_progress(calculated_at);

CREATE INDEX idx_goal_predictions_goal_id ON goal_predictions(goal_id);

CREATE INDEX idx_goal_adjustments_goal_id ON goal_adjustments(goal_id);
CREATE INDEX idx_goal_adjustments_created_at ON goal_adjustments(created_at);

CREATE INDEX idx_optimal_loads_user_id ON optimal_loads(user_id);
CREATE INDEX idx_optimal_loads_valid_until ON optimal_loads(valid_until);