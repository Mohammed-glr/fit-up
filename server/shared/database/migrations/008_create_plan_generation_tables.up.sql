

CREATE TABLE generated_plans (
    plan_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    algorithm VARCHAR(50) NOT NULL DEFAULT 'fitup_v1',
    effectiveness FLOAT DEFAULT NULL CHECK (effectiveness IS NULL OR (effectiveness >= 0 AND effectiveness <= 1)),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE plan_performance_data (
    performance_id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES generated_plans(plan_id) ON DELETE CASCADE,
    completion_rate FLOAT NOT NULL CHECK (completion_rate >= 0 AND completion_rate <= 1),
    average_rpe FLOAT DEFAULT NULL CHECK (average_rpe IS NULL OR (average_rpe >= 1 AND average_rpe <= 10)),
    progress_rate FLOAT NOT NULL DEFAULT 0,
    user_satisfaction FLOAT DEFAULT NULL CHECK (user_satisfaction IS NULL OR (user_satisfaction >= 0 AND user_satisfaction <= 1)),
    injury_rate FLOAT NOT NULL DEFAULT 0 CHECK (injury_rate >= 0 AND injury_rate <= 1),
    measured_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE plan_adaptations (
    adaptation_id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES generated_plans(plan_id) ON DELETE CASCADE,
    adaptation_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reason TEXT NOT NULL,
    changes JSONB NOT NULL DEFAULT '{}',
    trigger VARCHAR(50) NOT NULL DEFAULT 'manual'
);

CREATE TABLE plan_generation_metadata (
    metadata_id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES generated_plans(plan_id) ON DELETE CASCADE,
    user_goals JSONB NOT NULL DEFAULT '[]',
    available_equipment JSONB NOT NULL DEFAULT '[]',
    fitness_level VARCHAR(20) NOT NULL CHECK (fitness_level IN ('beginner', 'intermediate', 'advanced')),
    weekly_frequency INT NOT NULL CHECK (weekly_frequency >= 1 AND weekly_frequency <= 7),
    time_per_workout INT NOT NULL CHECK (time_per_workout > 0),
    algorithm_version VARCHAR(20) NOT NULL,
    algorithm_parameters JSONB NOT NULL DEFAULT '{}',
    UNIQUE(plan_id)
);

CREATE INDEX idx_generated_plans_user_id ON generated_plans(user_id);
CREATE INDEX idx_generated_plans_week_start ON generated_plans(week_start);
CREATE INDEX idx_generated_plans_active ON generated_plans(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_generated_plans_algorithm ON generated_plans(algorithm);

CREATE INDEX idx_plan_performance_data_plan_id ON plan_performance_data(plan_id);
CREATE INDEX idx_plan_performance_data_measured_at ON plan_performance_data(measured_at);

CREATE INDEX idx_plan_adaptations_plan_id ON plan_adaptations(plan_id);
CREATE INDEX idx_plan_adaptations_date ON plan_adaptations(adaptation_date);
CREATE INDEX idx_plan_adaptations_trigger ON plan_adaptations(trigger);

CREATE INDEX idx_plan_generation_metadata_plan_id ON plan_generation_metadata(plan_id);