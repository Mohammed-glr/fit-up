-- Drop workout_templates table and related objects
DROP TRIGGER IF EXISTS workout_template_updated_at_trigger ON workout_templates;
DROP FUNCTION IF EXISTS update_workout_template_updated_at();
DROP INDEX IF EXISTS idx_workout_templates_public;
DROP INDEX IF EXISTS idx_workout_templates_user_id;
DROP TABLE IF EXISTS workout_templates;
