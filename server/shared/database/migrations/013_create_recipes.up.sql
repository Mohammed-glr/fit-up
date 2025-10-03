
CREATE TABLE system_recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL CHECK (category IN ('breakfast', 'lunch', 'dinner', 'snack', 'dessert')),
    calories INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    fiber INTEGER NOT NULL,
    prep_time INTEGER NOT NULL,
    cook_time INTEGER,
    difficulty VARCHAR(20) CHECK (difficulty IN ('easy', 'medium', 'hard')),
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE system_recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES system_recipes(id) ON DELETE CASCADE,
    item VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    order_index INTEGER NOT NULL, 
    CONSTRAINT unique_ingredient_order UNIQUE(recipe_id, order_index)
);

CREATE TABLE system_recipe_instructions (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES system_recipes(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    instruction TEXT NOT NULL,
    CONSTRAINT unique_step UNIQUE(recipe_id, step_number)
);

CREATE TABLE system_recipe_tags (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES system_recipes(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    CONSTRAINT unique_recipe_tag UNIQUE(recipe_id, tag)
);

CREATE INDEX idx_system_recipes_category ON system_recipes(category);
CREATE INDEX idx_system_recipes_calories ON system_recipes(calories);
CREATE INDEX idx_system_recipe_tags_tag ON system_recipe_tags(tag);


CREATE TABLE user_recipes (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL CHECK (category IN ('breakfast', 'lunch', 'dinner', 'snack', 'dessert')),
    calories INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    fiber INTEGER,
    prep_time INTEGER, 
    cook_time INTEGER,
    difficulty VARCHAR(20) CHECK (difficulty IN ('easy', 'medium', 'hard')),
    image_url VARCHAR(500),
    is_favorite BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES user_recipes(id) ON DELETE CASCADE,
    item VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    order_index INTEGER NOT NULL,
    CONSTRAINT unique_user_ingredient_order UNIQUE(recipe_id, order_index)
);


CREATE TABLE user_recipe_instructions (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES user_recipes(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    instruction TEXT NOT NULL,
    CONSTRAINT unique_user_step UNIQUE(recipe_id, step_number)
);

CREATE TABLE user_recipe_tags (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES user_recipes(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    CONSTRAINT unique_user_recipe_tag UNIQUE(recipe_id, tag)
);

CREATE INDEX idx_user_recipes_user_id ON user_recipes(user_id);
CREATE INDEX idx_user_recipes_category ON user_recipes(category);
CREATE INDEX idx_user_recipes_is_favorite ON user_recipes(is_favorite);


CREATE TABLE user_favorite_recipes (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id INTEGER NOT NULL REFERENCES system_recipes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_favorite UNIQUE(user_id, recipe_id)
);


CREATE TABLE food_log_entries (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_date DATE NOT NULL,
    meal_type VARCHAR(50) CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    system_recipe_id INTEGER REFERENCES system_recipes(id) ON DELETE SET NULL,
    user_recipe_id INTEGER REFERENCES user_recipes(id) ON DELETE SET NULL,
    calories INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    fiber INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (system_recipe_id IS NOT NULL AND user_recipe_id IS NULL) OR
        (system_recipe_id IS NULL AND user_recipe_id IS NOT NULL) OR
        (system_recipe_id IS NULL AND user_recipe_id IS NULL)
    )
);

CREATE INDEX idx_food_log_user_date ON food_log_entries(user_id, log_date);
CREATE INDEX idx_food_log_date ON food_log_entries(log_date);
CREATE INDEX idx_food_log_system_recipe ON food_log_entries(system_recipe_id) WHERE system_recipe_id IS NOT NULL;
CREATE INDEX idx_food_log_user_recipe ON food_log_entries(user_recipe_id) WHERE user_recipe_id IS NOT NULL;





CREATE VIEW user_all_recipes AS
SELECT 
    'system' as source,
    sr.id,
    sr.name,
    sr.category,
    sr.calories,
    sr.protein,
    sr.carbs,
    sr.fat,
    sr.fiber,
    sr.prep_time,
    sr.image_url,
    NULL as user_id,
    EXISTS(SELECT 1 FROM user_favorite_recipes ufr WHERE ufr.recipe_id = sr.id) as is_favorite
FROM system_recipes sr
WHERE sr.is_active = true

UNION ALL

SELECT 
    'user' as source,
    ur.id,
    ur.name,
    ur.category,
    ur.calories,
    ur.protein,
    ur.carbs,
    ur.fat,
    ur.fiber,
    ur.prep_time,
    ur.image_url,
    ur.user_id,
    ur.is_favorite
FROM user_recipes ur;

CREATE VIEW daily_nutrition_summary AS
SELECT 
    user_id,
    log_date,
    SUM(calories) as total_calories,
    SUM(protein) as total_protein,
    SUM(carbs) as total_carbs,
    SUM(fat) as total_fat,
    SUM(fiber) as total_fiber,
    COUNT(*) as total_entries
FROM food_log_entries
GROUP BY user_id, log_date;

