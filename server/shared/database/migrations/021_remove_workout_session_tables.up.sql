-- Remove unused workout session tracking tables and other unused tables

-- Drop workout session tables in order (child tables first due to foreign keys)
DROP TABLE IF EXISTS weekly_session_stats CASCADE;
DROP TABLE IF EXISTS session_metrics CASCADE;
DROP TABLE IF EXISTS set_performances CASCADE;
DROP TABLE IF EXISTS exercise_performances CASCADE;
DROP TABLE IF EXISTS skipped_workouts CASCADE;
DROP TABLE IF EXISTS workout_sessions CASCADE;

-- Drop unused coach activity log (not referenced in code)
DROP TABLE IF EXISTS coach_activity_log CASCADE;

-- Drop associated indexes (if not automatically dropped)
DROP INDEX IF EXISTS idx_workout_sessions_user_id;
DROP INDEX IF EXISTS idx_workout_sessions_status;
DROP INDEX IF EXISTS idx_workout_sessions_start_time;
DROP INDEX IF EXISTS idx_workout_sessions_user_date;
DROP INDEX IF EXISTS idx_skipped_workouts_user_id;
DROP INDEX IF EXISTS idx_skipped_workouts_date;
DROP INDEX IF EXISTS idx_exercise_performances_session_id;
DROP INDEX IF EXISTS idx_exercise_performances_exercise_id;
DROP INDEX IF EXISTS idx_set_performances_performance_id;
DROP INDEX IF EXISTS idx_set_performances_completed_at;
DROP INDEX IF EXISTS idx_session_metrics_session_id;
DROP INDEX IF EXISTS idx_weekly_session_stats_user_id;
DROP INDEX IF EXISTS idx_weekly_session_stats_week_start;

-- Drop coach activity log indexes
DROP INDEX IF EXISTS idx_coach_activity_coach_id;
DROP INDEX IF EXISTS idx_coach_activity_timestamp;
