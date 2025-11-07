ALTER TABLE system_recipes
    ADD COLUMN IF NOT EXISTS servings INTEGER DEFAULT 1 CHECK (servings > 0);
