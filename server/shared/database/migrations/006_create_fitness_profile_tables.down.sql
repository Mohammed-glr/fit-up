
DROP INDEX IF EXISTS idx_one_rep_max_estimates_date;
DROP INDEX IF EXISTS idx_one_rep_max_estimates_user_exercise;
DROP INDEX IF EXISTS idx_movement_limitations_type;
DROP INDEX IF EXISTS idx_movement_limitations_user_id;
DROP INDEX IF EXISTS idx_movement_assessments_date;
DROP INDEX IF EXISTS idx_movement_assessments_user_id;
DROP INDEX IF EXISTS idx_fitness_goal_targets_target_date;
DROP INDEX IF EXISTS idx_fitness_goal_targets_type;
DROP INDEX IF EXISTS idx_fitness_goal_targets_active;
DROP INDEX IF EXISTS idx_fitness_goal_targets_user_id;
DROP INDEX IF EXISTS idx_fitness_assessments_date;
DROP INDEX IF EXISTS idx_fitness_assessments_user_id;
DROP INDEX IF EXISTS idx_workout_profiles_level_goal;
DROP INDEX IF EXISTS idx_workout_profiles_auth_user_id;

DROP TABLE IF EXISTS one_rep_max_estimates CASCADE;
DROP TABLE IF EXISTS movement_limitations CASCADE;
DROP TABLE IF EXISTS movement_assessments CASCADE;
DROP TABLE IF EXISTS fitness_goal_targets CASCADE;
DROP TABLE IF EXISTS fitness_assessments CASCADE;
DROP TABLE IF EXISTS workout_profiles CASCADE;
