-- Rollback: Recreate unused analytics tables
-- Note: This rollback recreates the tables but data will be lost

CREATE TABLE goal_progress (
    progress_id SERIAL PRIMARY KEY,
    goal_id INT NOT NULL REFERENCES fitness_goal_targets(goal_id) ON DELETE CASCADE,
    progress_percent FLOAT NOT NULL CHECK (progress_percent >= 0 AND progress_percent <= 100),
    on_track BOOLEAN NOT NULL DEFAULT TRUE,
    estimated_completion TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
    intensity_range VARCHAR(20) NOT NULL,
    volume_target FLOAT NOT NULL CHECK (volume_target > 0),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_goal_progress_goal_id ON goal_progress(goal_id);
CREATE INDEX idx_goal_progress_calculated_at ON goal_progress(calculated_at);
CREATE INDEX idx_goal_predictions_goal_id ON goal_predictions(goal_id);
CREATE INDEX idx_goal_adjustments_goal_id ON goal_adjustments(goal_id);
CREATE INDEX idx_goal_adjustments_created_at ON goal_adjustments(created_at);
CREATE INDEX idx_optimal_loads_user_id ON optimal_loads(user_id);
CREATE INDEX idx_optimal_loads_valid_until ON optimal_loads(valid_until);
