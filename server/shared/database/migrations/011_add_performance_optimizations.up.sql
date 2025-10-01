
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_progress_logs_user_exercise_date') THEN
        CREATE INDEX idx_progress_logs_user_exercise_date 
            ON progress_logs(user_id, exercise_id, date DESC);
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_workout_exercises_workout_exercise') THEN
        CREATE INDEX idx_workout_exercises_workout_exercise 
            ON workout_exercises(workout_id, exercise_id);
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_exercises_equipment_difficulty') THEN
        CREATE INDEX idx_exercises_equipment_difficulty 
            ON exercises(equipment, difficulty);
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_exercises_type_muscle_groups') THEN
        CREATE INDEX idx_exercises_type_muscle_groups 
            ON exercises(type, muscle_groups);
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_weekly_schemas_active_user') THEN
        CREATE INDEX idx_weekly_schemas_active_user 
            ON weekly_schemas(user_id, week_start DESC) WHERE active = TRUE;
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_fitness_goal_targets_active_user') THEN
        CREATE INDEX idx_fitness_goal_targets_active_user 
            ON fitness_goal_targets(user_id, target_date) WHERE is_active = TRUE;
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_workout_sessions_active') THEN
        CREATE INDEX idx_workout_sessions_active 
            ON workout_sessions(user_id, start_time DESC) WHERE status = 'active';
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_exercises_name_search') THEN
        CREATE INDEX idx_exercises_name_search 
            ON exercises USING GIN (to_tsvector('english', name));
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_workout_templates_name_search') THEN
        CREATE INDEX idx_workout_templates_name_search 
            ON workout_templates USING GIN (to_tsvector('english', name || ' ' || COALESCE(description, '')));
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_progress_logs_date_weight') THEN
        CREATE INDEX idx_progress_logs_date_weight 
            ON progress_logs(date, weight_used) WHERE weight_used IS NOT NULL;
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_exercise_performances_session_volume') THEN
        CREATE INDEX idx_exercise_performances_session_volume 
            ON exercise_performances(session_id, total_volume DESC);
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_workout_profiles_auth_user') THEN
        ALTER TABLE workout_profiles 
            ADD CONSTRAINT fk_workout_profiles_auth_user 
            FOREIGN KEY (auth_user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'uk_workout_profiles_auth_user_id') THEN
        ALTER TABLE workout_profiles 
            ADD CONSTRAINT uk_workout_profiles_auth_user_id 
            UNIQUE (auth_user_id);
    END IF;
END $$;

ANALYZE exercises;
ANALYZE workout_templates;
ANALYZE workouts;
ANALYZE workout_exercises;
ANALYZE progress_logs;
ANALYZE weekly_schemas;

COMMENT ON TABLE exercises IS 'Master table of all available exercises';
COMMENT ON TABLE workout_templates IS 'Pre-defined workout program templates';
COMMENT ON TABLE weekly_schemas IS 'User weekly workout schedules';
COMMENT ON TABLE workouts IS 'Individual workout days within a weekly schema';
COMMENT ON TABLE workout_exercises IS 'Exercises assigned to specific workouts';
COMMENT ON TABLE progress_logs IS 'User exercise performance tracking';

COMMENT ON COLUMN exercises.muscle_groups IS 'Comma-separated list of primary muscle groups';
COMMENT ON COLUMN exercises.default_reps IS 'Recommended rep range (e.g., "8-12")';
COMMENT ON COLUMN workout_templates.suitable_goals IS 'Comma-separated list of fitness goals';
COMMENT ON COLUMN weekly_schemas.week_start IS 'Monday of the workout week';
COMMENT ON COLUMN workouts.day_of_week IS '1=Monday, 2=Tuesday, ..., 7=Sunday';

CREATE OR REPLACE VIEW v_exercise_popularity AS
SELECT 
    e.exercise_id,
    e.name,
    COUNT(we.we_id) as usage_count,
    COUNT(DISTINCT w.schema_id) as unique_schemas
FROM exercises e
LEFT JOIN workout_exercises we ON e.exercise_id = we.exercise_id
LEFT JOIN workouts w ON we.workout_id = w.workout_id
GROUP BY e.exercise_id, e.name
ORDER BY usage_count DESC;

CREATE OR REPLACE VIEW v_user_activity_summary AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(DISTINCT ws.session_id) as total_sessions,
    COUNT(DISTINCT DATE(ws.start_time)) as workout_days,
    AVG(sm.completion_rate) as avg_completion_rate,
    MAX(ws.start_time) as last_workout
FROM users u
LEFT JOIN workout_sessions ws ON u.id = ws.user_id
LEFT JOIN session_metrics sm ON ws.session_id = sm.session_id
WHERE ws.status = 'completed'
GROUP BY u.id, u.username;

CREATE OR REPLACE FUNCTION cleanup_old_data()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER := 0;
    temp_count INTEGER := 0;
BEGIN
    DELETE FROM workout_sessions 
    WHERE start_time < NOW() - INTERVAL '2 years' 
      AND status IN ('completed', 'skipped', 'abandoned');
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM recovery_metrics 
    WHERE date < NOW() - INTERVAL '1 year';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM plan_adaptations 
    WHERE adaptation_date < NOW() - INTERVAL '6 months';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;