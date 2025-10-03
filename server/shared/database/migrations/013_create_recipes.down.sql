
DROP VIEW IF EXISTS user_all_recipes;
DROP VIEW IF EXISTS daily_nutrition_summary;
DROP TABLE IF EXISTS user_favorite_recipes CASCADE;
DROP TABLE IF EXISTS user_recipes CASCADE;
DROP TABLE IF EXISTS system_recipe_instructions CASCADE;
DROP TABLE IF EXISTS system_recipe_ingredients CASCADE;
DROP TABLE IF EXISTS system_recipes CASCADE;
DROP TABLE IF EXISTS food_log_entries CASCADE;
DROP TYPE IF EXISTS recipe_category;
DROP TYPE IF EXISTS ingredient_unit;