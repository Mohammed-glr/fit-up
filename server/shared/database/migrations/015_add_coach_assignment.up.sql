CREATE TABLE IF NOT EXISTS coach_assignments (
    assignment_id SERIAL PRIMARY KEY,
    coach_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    deactivated_at TIMESTAMP,
    notes TEXT,
    
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES workout_profiles(workout_profile_id) ON DELETE CASCADE,
    CONSTRAINT unique_active_assignment UNIQUE (coach_id, user_id, is_active)
);

CREATE INDEX idx_coach_assignments_coach_id ON coach_assignments(coach_id);
CREATE INDEX idx_coach_assignments_user_id ON coach_assignments(user_id);
CREATE INDEX idx_coach_assignments_active ON coach_assignments(is_active);

CREATE TABLE IF NOT EXISTS coach_activity_log (
    activity_id SERIAL PRIMARY KEY,
    coach_id TEXT NOT NULL,
    user_id INTEGER,
    activity_type TEXT NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_activity_type CHECK (activity_type IN (
        'schema_created', 'schema_updated', 'schema_deleted',
        'client_assigned', 'client_removed', 'goal_created',
        'assessment_created', 'note_added'
    ))
);

CREATE INDEX idx_coach_activity_coach_id ON coach_activity_log(coach_id);
CREATE INDEX idx_coach_activity_timestamp ON coach_activity_log(created_at);
