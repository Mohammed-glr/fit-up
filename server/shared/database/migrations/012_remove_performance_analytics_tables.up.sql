
DROP INDEX IF EXISTS idx_optimal_loads_valid_until;
DROP INDEX IF EXISTS idx_optimal_loads_user_id;
DROP INDEX IF EXISTS idx_intensity_progressions_exercise_id;
DROP INDEX IF EXISTS idx_intensity_progressions_user_id;
DROP INDEX IF EXISTS idx_training_volumes_week_start;
DROP INDEX IF EXISTS idx_training_volumes_user_id;
DROP INDEX IF EXISTS idx_plateau_detections_detected;
DROP INDEX IF EXISTS idx_plateau_detections_user_id;
DROP INDEX IF EXISTS idx_strength_progressions_user_exercise;
DROP INDEX IF EXISTS idx_strength_progressions_exercise_id;
DROP INDEX IF EXISTS idx_strength_progressions_user_id;
DROP TABLE IF EXISTS optimal_loads CASCADE;
DROP TABLE IF EXISTS intensity_progressions CASCADE;
DROP TABLE IF EXISTS training_volumes CASCADE;
DROP TABLE IF EXISTS plateau_detections CASCADE;
DROP TABLE IF EXISTS strength_progressions CASCADE;


