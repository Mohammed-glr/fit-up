
DROP INDEX IF EXISTS idx_plan_generation_metadata_plan_id;
DROP INDEX IF EXISTS idx_plan_adaptations_trigger;
DROP INDEX IF EXISTS idx_plan_adaptations_date;
DROP INDEX IF EXISTS idx_plan_adaptations_plan_id;
DROP INDEX IF EXISTS idx_plan_performance_data_measured_at;
DROP INDEX IF EXISTS idx_plan_performance_data_plan_id;
DROP INDEX IF EXISTS idx_generated_plans_algorithm;
DROP INDEX IF EXISTS idx_generated_plans_active;
DROP INDEX IF EXISTS idx_generated_plans_week_start;
DROP INDEX IF EXISTS idx_generated_plans_user_id;

DROP TABLE IF EXISTS plan_generation_metadata CASCADE;
DROP TABLE IF EXISTS plan_adaptations CASCADE;
DROP TABLE IF EXISTS plan_performance_data CASCADE;
DROP TABLE IF EXISTS generated_plans CASCADE;
