CREATE TABLE nutrition_goals (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calories_goal INTEGER NOT NULL,
    protein_goal INTEGER NOT NULL,
    carbs_goal INTEGER NOT NULL,
    fat_goal INTEGER NOT NULL,
    fiber_goal INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_goal UNIQUE(user_id)
);
CREATE INDEX idx_nutrition_goals_user_id ON nutrition_goals(user_id);