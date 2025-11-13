-- Rollback mindfulness tables

DROP TABLE IF EXISTS mindfulness_streaks CASCADE;
DROP TABLE IF EXISTS reflection_responses CASCADE;
DROP TABLE IF EXISTS reflection_prompts CASCADE;
DROP TABLE IF EXISTS gratitude_entries CASCADE;
DROP TABLE IF EXISTS breathing_exercises CASCADE;
DROP TABLE IF EXISTS mindfulness_sessions CASCADE;

DROP INDEX IF EXISTS idx_mindfulness_sessions_user_date;
DROP INDEX IF EXISTS idx_breathing_exercises_user_date;
DROP INDEX IF EXISTS idx_gratitude_entries_user_date;
DROP INDEX IF EXISTS idx_reflection_responses_user_date;
DROP INDEX IF EXISTS idx_mindfulness_streaks_user;
