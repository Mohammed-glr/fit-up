-- Remove unused analytics tables that are not referenced in code

-- Drop unused goal tracking analytics tables
DROP TABLE IF EXISTS goal_adjustments CASCADE;
DROP TABLE IF EXISTS goal_predictions CASCADE;
DROP TABLE IF EXISTS goal_progress CASCADE;

-- Drop unused load optimization table
DROP TABLE IF EXISTS optimal_loads CASCADE;

-- Drop associated indexes
DROP INDEX IF EXISTS idx_goal_progress_goal_id;
DROP INDEX IF EXISTS idx_goal_progress_calculated_at;
DROP INDEX IF EXISTS idx_goal_predictions_goal_id;
DROP INDEX IF EXISTS idx_goal_adjustments_goal_id;
DROP INDEX IF EXISTS idx_goal_adjustments_created_at;
DROP INDEX IF EXISTS idx_optimal_loads_user_id;
DROP INDEX IF EXISTS idx_optimal_loads_valid_until;
