n
DROP FUNCTION IF EXISTS cleanup_old_data();

DROP VIEW IF EXISTS v_user_activity_summary;
DROP VIEW IF EXISTS v_exercise_popularity;

COMMENT ON COLUMN workouts.day_of_week IS NULL;
COMMENT ON COLUMN weekly_schemas.week_start IS NULL;
COMMENT ON COLUMN workout_templates.suitable_goals IS NULL;
COMMENT ON COLUMN exercises.default_reps IS NULL;
COMMENT ON COLUMN exercises.muscle_groups IS NULL;

COMMENT ON TABLE progress_logs IS NULL;
COMMENT ON TABLE workout_exercises IS NULL;
COMMENT ON TABLE workouts IS NULL;
COMMENT ON TABLE weekly_schemas IS NULL;
COMMENT ON TABLE workout_templates IS NULL;
COMMENT ON TABLE exercises IS NULL;

ALTER TABLE workout_profiles DROP CONSTRAINT IF EXISTS uk_workout_profiles_auth_user_id;

ALTER TABLE workout_profiles DROP CONSTRAINT IF EXISTS fk_workout_profiles_auth_user;

DROP INDEX CONCURRENTLY IF EXISTS idx_exercise_performances_session_volume;
DROP INDEX CONCURRENTLY IF EXISTS idx_progress_logs_date_weight;
DROP INDEX CONCURRENTLY IF EXISTS idx_workout_templates_name_search;
DROP INDEX CONCURRENTLY IF EXISTS idx_exercises_name_search;
DROP INDEX CONCURRENTLY IF EXISTS idx_workout_sessions_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_fitness_goal_targets_active_user;
DROP INDEX CONCURRENTLY IF EXISTS idx_weekly_schemas_active_user;
DROP INDEX CONCURRENTLY IF EXISTS idx_exercises_type_muscle_groups;
DROP INDEX CONCURRENTLY IF EXISTS idx_exercises_equipment_difficulty;
DROP INDEX CONCURRENTLY IF EXISTS idx_workout_exercises_workout_exercise;
DROP INDEX CONCURRENTLY IF EXISTS idx_progress_logs_user_exercise_date;
