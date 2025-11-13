export interface MindfulnessSession {
  session_id: number;
  user_id: string;
  session_type: 'pre_workout' | 'post_workout' | 'breathing' | 'meditation' | 'gratitude';
  duration_seconds: number;
  completed_at: string;
  notes?: string;
  mood_before?: number; // 1-5
  mood_after?: number; // 1-5
}

export interface BreathingExercise {
  exercise_id: number;
  user_id: string;
  breathing_type: 'box' | '478' | 'energizing' | 'calming' | 'custom';
  duration_seconds: number;
  cycles_completed: number;
  completed_at: string;
  heart_rate_before?: number;
  heart_rate_after?: number;
}

export interface GratitudeEntry {
  entry_id: number;
  user_id: string;
  entry_text: string;
  tags?: string[];
  mood?: number; // 1-5
  created_at: string;
  workout_session_id?: number;
}

export interface ReflectionPrompt {
  prompt_id: number;
  prompt_text: string;
  category: string;
  is_active: boolean;
}

export interface ReflectionResponse {
  response_id: number;
  user_id: string;
  prompt_id?: number;
  response_text: string;
  created_at: string;
}

export interface MindfulnessStreak {
  streak_id: number;
  user_id: string;
  current_streak: number;
  longest_streak: number;
  last_activity_date: string;
  total_sessions: number;
}

// Request types
export interface CreateMindfulnessSessionRequest {
  session_type: string;
  duration_seconds: number;
  notes?: string;
  mood_before?: number;
  mood_after?: number;
}

export interface CreateBreathingExerciseRequest {
  breathing_type: string;
  duration_seconds: number;
  cycles_completed: number;
  heart_rate_before?: number;
  heart_rate_after?: number;
}

export interface CreateGratitudeEntryRequest {
  entry_text: string;
  tags?: string[];
  mood?: number;
  workout_session_id?: number;
}

export interface CreateReflectionResponseRequest {
  prompt_id?: number;
  response_text: string;
}

export interface MoodDataPoint {
  date: string;
  avg_before: number;
  avg_after: number;
}

export interface MindfulnessStats {
  total_sessions: number;
  total_minutes: number;
  current_streak: number;
  longest_streak: number;
  sessions_by_type: Record<string, number>;
  recent_sessions: MindfulnessSession[];
  mood_trend?: MoodDataPoint[];
}

export interface BreathingStats {
  total_exercises: number;
  total_minutes: number;
  total_cycles: number;
  exercises_by_type: Record<string, number>;
  recent_exercises: BreathingExercise[];
}

export interface BreathingPattern {
  name: string;
  type: 'box' | '478' | 'energizing' | 'calming' | 'custom';
  description: string;
  pattern: number[]; // [inhale, hold, exhale, hold]
  defaultCycles: number;
  duration: number; // Total duration in seconds
  benefits: string[];
}

export const BREATHING_PATTERNS: Record<string, BreathingPattern> = {
  box: {
    name: 'Box Breathing',
    type: 'box',
    description: 'Equal parts inhale, hold, exhale, hold',
    pattern: [4, 4, 4, 4],
    defaultCycles: 4,
    duration: 240, // 4 minutes
    benefits: ['Reduces stress', 'Improves focus', 'Calms nervous system'],
  },
  '478': {
    name: '4-7-8 Breathing',
    type: '478',
    description: 'Relaxing breath for sleep and anxiety',
    pattern: [4, 7, 8, 0],
    defaultCycles: 4,
    duration: 152, // ~2.5 minutes
    benefits: ['Promotes sleep', 'Reduces anxiety', 'Lowers heart rate'],
  },
  energizing: {
    name: 'Energizing Breath',
    type: 'energizing',
    description: 'Quick breaths to boost energy',
    pattern: [2, 0, 2, 0],
    defaultCycles: 10,
    duration: 120, // 2 minutes
    benefits: ['Increases energy', 'Boosts alertness', 'Improves circulation'],
  },
  calming: {
    name: 'Calming Breath',
    type: 'calming',
    description: 'Extended exhale for relaxation',
    pattern: [4, 0, 6, 0],
    defaultCycles: 6,
    duration: 180, // 3 minutes
    benefits: ['Deep relaxation', 'Reduces tension', 'Calms mind'],
  },
};

export const MOOD_SCALE = [
  { value: 1, label: '😞 Very Low', emoji: '😞' },
  { value: 2, label: '😕 Low', emoji: '😕' },
  { value: 3, label: '😐 Neutral', emoji: '😐' },
  { value: 4, label: '😊 Good', emoji: '😊' },
  { value: 5, label: '😄 Excellent', emoji: '😄' },
];

export const GRATITUDE_TAGS = [
  'health',
  'fitness',
  'progress',
  'relationships',
  'achievement',
  'nature',
  'self-care',
  'mindfulness',
  'recovery',
  'strength',
];
