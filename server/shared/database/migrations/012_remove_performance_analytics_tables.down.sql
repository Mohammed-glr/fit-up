-- Migration Rollback: Recreate Performance Analytics Tables if needed
-- This is the down migration for 013_remove_performance_analytics_tables

-- Recreate strength_progressions table
CREATE TABLE IF NOT EXISTS strength_progressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL,
    week_start DATE NOT NULL,
    max_weight DECIMAL(6,2),
    total_volume DECIMAL(10,2),
    average_reps INTEGER,
    progression_rate DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_strength_progressions_user_id ON strength_progressions(user_id);
CREATE INDEX idx_strength_progressions_exercise_id ON strength_progressions(exercise_id);
CREATE INDEX idx_strength_progressions_user_exercise ON strength_progressions(user_id, exercise_id);

-- Recreate plateau_detections table
CREATE TABLE IF NOT EXISTS plateau_detections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    duration_weeks INTEGER,
    suggested_deload BOOLEAN DEFAULT false,
    alternative_exercises JSONB,
    resolved BOOLEAN DEFAULT false,
    resolved_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_plateau_detections_user_id ON plateau_detections(user_id);
CREATE INDEX idx_plateau_detections_detected ON plateau_detections(detected_at);

-- Recreate training_volumes table
CREATE TABLE IF NOT EXISTS training_volumes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    total_sets INTEGER,
    total_reps INTEGER,
    total_volume DECIMAL(12,2),
    training_frequency INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_training_volumes_user_id ON training_volumes(user_id);
CREATE INDEX idx_training_volumes_week_start ON training_volumes(week_start);

-- Recreate intensity_progressions table
CREATE TABLE IF NOT EXISTS intensity_progressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL,
    week_start DATE NOT NULL,
    average_intensity DECIMAL(5,2),
    max_intensity DECIMAL(5,2),
    volume_load DECIMAL(12,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_intensity_progressions_user_id ON intensity_progressions(user_id);
CREATE INDEX idx_intensity_progressions_exercise_id ON intensity_progressions(exercise_id);

-- Recreate optimal_loads table
CREATE TABLE IF NOT EXISTS optimal_loads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL,
    optimal_1rm DECIMAL(6,2),
    optimal_5rm DECIMAL(6,2),
    optimal_10rm DECIMAL(6,2),
    confidence_score DECIMAL(3,2),
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    valid_until TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_optimal_loads_user_id ON optimal_loads(user_id);
CREATE INDEX idx_optimal_loads_valid_until ON optimal_loads(valid_until);
