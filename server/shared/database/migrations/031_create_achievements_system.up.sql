-- Create achievements/badges table
CREATE TABLE IF NOT EXISTS achievements (
    achievement_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    badge_icon VARCHAR(50) NOT NULL,
    badge_color VARCHAR(20) NOT NULL,
    category VARCHAR(50) NOT NULL,
    requirement_type VARCHAR(50) NOT NULL,
    requirement_value INT NOT NULL,
    points INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS user_achievements (
    user_achievement_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id INT NOT NULL REFERENCES achievements(achievement_id) ON DELETE CASCADE,
    earned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    progress INT DEFAULT 0,
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS idx_user_achievements_user ON user_achievements(user_id, earned_at DESC);
CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);

INSERT INTO achievements (name, description, badge_icon, badge_color, category, requirement_type, requirement_value, points) VALUES
('Getting Started', 'Complete your first workout', 'rocket', 'primary', 'milestone', 'total_workouts', 1, 10),
('Consistent Beginner', 'Maintain a 7-day workout streak', 'flame', 'warning', 'streak', 'streak_days', 7, 50),
('Dedicated Athlete', 'Maintain a 30-day workout streak', 'flame', 'warning', 'streak', 'streak_days', 30, 200),
('Unstoppable', 'Maintain a 100-day workout streak', 'flame', 'error', 'streak', 'streak_days', 100, 500),

('Lightweight', 'Lift a total of 10,000 lbs', 'barbell', 'success', 'volume', 'total_volume_lbs', 10000, 30),
('Middleweight', 'Lift a total of 50,000 lbs', 'barbell', 'success', 'volume', 'total_volume_lbs', 50000, 100),
('Heavyweight', 'Lift a total of 100,000 lbs', 'barbell', 'success', 'volume', 'total_volume_lbs', 100000, 250),
('Titan', 'Lift a total of 500,000 lbs', 'barbell', 'error', 'volume', 'total_volume_lbs', 500000, 1000),

('First PR', 'Set your first personal record', 'trophy', 'warning', 'pr', 'pr_count', 1, 25),
('PR Collector', 'Set 10 personal records', 'trophy', 'warning', 'pr', 'pr_count', 10, 100),
('Record Breaker', 'Set 50 personal records', 'trophy', 'error', 'pr', 'pr_count', 50, 400),

('Ten Club', 'Complete 10 workouts', 'fitness', 'primary', 'milestone', 'total_workouts', 10, 30),
('Fifty Club', 'Complete 50 workouts', 'fitness', 'primary', 'milestone', 'total_workouts', 50, 150),
('Century Club', 'Complete 100 workouts', 'fitness', 'success', 'milestone', 'total_workouts', 100, 300),
('Elite 500', 'Complete 500 workouts', 'fitness', 'error', 'milestone', 'total_workouts', 500, 1500),

('Weekend Warrior', 'Complete 20 weekend workouts', 'calendar', 'info', 'consistency', 'weekend_workouts', 20, 75),
('Early Bird', 'Complete 30 morning workouts (before 9 AM)', 'sunny', 'warning', 'consistency', 'morning_workouts', 30, 100),
('Night Owl', 'Complete 30 evening workouts (after 6 PM)', 'moon', 'primary', 'consistency', 'evening_workouts', 30, 100)
ON CONFLICT (name) DO NOTHING;

COMMENT ON TABLE achievements IS 'Available achievement badges in the system';
COMMENT ON TABLE user_achievements IS 'User-specific achievement progress and earned badges';
