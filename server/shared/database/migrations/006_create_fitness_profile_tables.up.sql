CREATE TABLE workout_profiles (
    workout_profile_id SERIAL PRIMARY KEY,
    auth_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level VARCHAR(20) NOT NULL CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    goal VARCHAR(20) NOT NULL CHECK (goal IN ('strength', 'muscle_gain', 'fat_loss', 'endurance', 'general_fitness')),
    frequency INT NOT NULL CHECK (frequency >= 1 AND frequency <= 7),
    equipment JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(auth_user_id)
);

CREATE TABLE fitness_assessments (
    assessment_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assessment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    overall_level VARCHAR(20) NOT NULL CHECK (overall_level IN ('beginner', 'intermediate', 'advanced')),
    strength_level VARCHAR(20) NOT NULL CHECK (strength_level IN ('beginner', 'intermediate', 'advanced')),
    cardio_level VARCHAR(20) NOT NULL CHECK (cardio_level IN ('beginner', 'intermediate', 'advanced')),
    flexibility_level VARCHAR(20) NOT NULL CHECK (flexibility_level IN ('beginner', 'intermediate', 'advanced')),
    assessment_data JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE fitness_goal_targets (
    goal_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    goal_type VARCHAR(20) NOT NULL CHECK (goal_type IN ('strength', 'muscle_gain', 'fat_loss', 'endurance', 'general_fitness')),
    target_value FLOAT NOT NULL CHECK (target_value >= 0),
    current_value FLOAT NOT NULL DEFAULT 0 CHECK (current_value >= 0),
    target_date TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}'
);

CREATE TABLE movement_assessments (
    assessment_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assessment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    movement_data JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE movement_limitations (
    limitation_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movement_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('mild', 'moderate', 'severe')),
    description TEXT NOT NULL
);

CREATE TABLE one_rep_max_estimates (
    estimate_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    estimated_max FLOAT NOT NULL CHECK (estimated_max > 0),
    estimate_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    method VARCHAR(50) NOT NULL DEFAULT 'epley',
    confidence FLOAT NOT NULL DEFAULT 0.8 CHECK (confidence BETWEEN 0 AND 1),
    UNIQUE(user_id, exercise_id, estimate_date)
);

CREATE INDEX idx_workout_profiles_auth_user_id ON workout_profiles(auth_user_id);
CREATE INDEX idx_workout_profiles_level_goal ON workout_profiles(level, goal);

CREATE INDEX idx_fitness_assessments_user_id ON fitness_assessments(user_id);
CREATE INDEX idx_fitness_assessments_date ON fitness_assessments(assessment_date);

CREATE INDEX idx_fitness_goal_targets_user_id ON fitness_goal_targets(user_id);
CREATE INDEX idx_fitness_goal_targets_active ON fitness_goal_targets(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_fitness_goal_targets_type ON fitness_goal_targets(goal_type);
CREATE INDEX idx_fitness_goal_targets_target_date ON fitness_goal_targets(target_date);

CREATE INDEX idx_movement_assessments_user_id ON movement_assessments(user_id);
CREATE INDEX idx_movement_assessments_date ON movement_assessments(assessment_date);

CREATE INDEX idx_movement_limitations_user_id ON movement_limitations(user_id);
CREATE INDEX idx_movement_limitations_type ON movement_limitations(movement_type);

CREATE INDEX idx_one_rep_max_estimates_user_exercise ON one_rep_max_estimates(user_id, exercise_id);
CREATE INDEX idx_one_rep_max_estimates_date ON one_rep_max_estimates(estimate_date);
