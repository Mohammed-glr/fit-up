DROP INDEX IF EXISTS idx_generated_plan_exercises_day_id;
DROP INDEX IF EXISTS idx_generated_plan_days_plan_id;

DROP TABLE IF EXISTS generated_plan_exercises CASCADE;
DROP TABLE IF EXISTS generated_plan_days CASCADE;
