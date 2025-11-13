
CREATE TABLE IF NOT EXISTS mindfulness_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_type VARCHAR(20) NOT NULL CHECK (session_type IN ('pre_workout', 'post_workout', 'breathing', 'meditation', 'gratitude')),
    duration_seconds INT NOT NULL CHECK (duration_seconds >= 0),
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    mood_before INT CHECK (mood_before >= 1 AND mood_before <= 5),
    mood_after INT CHECK (mood_after >= 1 AND mood_after <= 5)
);

CREATE TABLE IF NOT EXISTS breathing_exercises (
    exercise_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    breathing_type VARCHAR(30) NOT NULL CHECK (breathing_type IN ('box', '478', 'energizing', 'calming', 'custom')),
    duration_seconds INT NOT NULL CHECK (duration_seconds >= 0),
    cycles_completed INT NOT NULL DEFAULT 0,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    heart_rate_before INT,
    heart_rate_after INT
);

CREATE TABLE IF NOT EXISTS gratitude_entries (
    entry_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entry_text TEXT NOT NULL,
    tags TEXT[],
    mood INT CHECK (mood >= 1 AND mood <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    workout_session_id INT REFERENCES workout_sessions(session_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS reflection_prompts (
    prompt_id SERIAL PRIMARY KEY,
    prompt_text TEXT NOT NULL,
    category VARCHAR(30) NOT NULL CHECK (category IN ('fitness', 'gratitude', 'growth', 'recovery', 'motivation')),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS reflection_responses (
    response_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prompt_id INT REFERENCES reflection_prompts(prompt_id) ON DELETE SET NULL,
    response_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mindfulness_streaks (
    streak_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    last_activity_date DATE NOT NULL,
    total_sessions INT NOT NULL DEFAULT 0,
    UNIQUE(user_id)
);

CREATE INDEX IF NOT EXISTS idx_mindfulness_sessions_user_date ON mindfulness_sessions(user_id, completed_at DESC);
CREATE INDEX IF NOT EXISTS idx_breathing_exercises_user_date ON breathing_exercises(user_id, completed_at DESC);
CREATE INDEX IF NOT EXISTS idx_gratitude_entries_user_date ON gratitude_entries(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reflection_responses_user_date ON reflection_responses(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_mindfulness_streaks_user ON mindfulness_streaks(user_id);

INSERT INTO reflection_prompts (prompt_text, category) VALUES
('What did you accomplish today?', 'fitness'),
('How did your body feel during training?', 'fitness'),
('One thing you''re grateful for right now.', 'gratitude'),
('What energy level do you feel? (1-10)', 'recovery'),
('What was your biggest challenge today?', 'growth'),
('What made you smile today?', 'gratitude'),
('How can you improve tomorrow?', 'growth'),
('What are you proud of this week?', 'motivation'),
('How well did you sleep last night?', 'recovery'),
('What intention do you set for tomorrow?', 'motivation')
ON CONFLICT DO NOTHING;
