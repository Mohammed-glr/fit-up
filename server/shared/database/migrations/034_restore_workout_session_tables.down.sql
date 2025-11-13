
DROP TABLE IF EXISTS weekly_session_stats CASCADE;
DROP TABLE IF EXISTS session_metrics CASCADE;
DROP TABLE IF EXISTS set_performances CASCADE;
DROP TABLE IF EXISTS exercise_performances CASCADE;
DROP TABLE IF EXISTS skipped_workouts CASCADE;
DROP TABLE IF EXISTS workout_sessions CASCADE;

DROP INDEX IF EXISTS idx_workout_sessions_user_id;
DROP INDEX IF EXISTS idx_workout_sessions_status;
DROP INDEX IF EXISTS idx_workout_sessions_start_time;
DROP INDEX IF EXISTS idx_workout_sessions_user_date;
DROP INDEX IF EXISTS idx_workout_sessions_active;
DROP INDEX IF EXISTS idx_skipped_workouts_user_id;
DROP INDEX IF EXISTS idx_skipped_workouts_date;
DROP INDEX IF EXISTS idx_exercise_performances_session_id;
DROP INDEX IF EXISTS idx_exercise_performances_exercise_id;
DROP INDEX IF EXISTS idx_set_performances_performance_id;
DROP INDEX IF EXISTS idx_set_performances_completed_at;
DROP INDEX IF EXISTS idx_session_metrics_session_id;
DROP INDEX IF EXISTS idx_weekly_session_stats_user_id;
DROP INDEX IF EXISTS idx_weekly_session_stats_week_start;
