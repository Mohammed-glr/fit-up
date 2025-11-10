/**
 * Types for Workout Sharing feature
 * Allows users to share completed workouts with coach, friends, or social media
 */

export interface ShareWorkoutOptions {
  shareWithCoach?: boolean;
  exportAsImage?: boolean;
  copyAsText?: boolean;
  shareToSocial?: boolean;
}

export interface ShareWorkoutRequest {
  session_id: number;
  share_type: 'coach' | 'image' | 'text' | 'social';
  message?: string;
}

export interface WorkoutShareSummary {
  session_id: number;
  workout_title: string;
  completed_at: string;
  duration_minutes: number;
  total_exercises: number;
  total_sets: number;
  total_reps: number;
  total_volume_lbs: number;
  exercises: WorkoutShareExercise[];
  prs_achieved: number;
  user_name?: string;
  user_photo_url?: string;
}

export interface WorkoutShareExercise {
  exercise_name: string;
  sets_completed: number;
  total_reps: number;
  total_volume_lbs: number;
  pr_achieved: boolean;
  best_set?: {
    weight: number;
    reps: number;
  };
}

export interface ShareWorkoutResponse {
  success: boolean;
  share_url?: string;
  share_text?: string;
  message: string;
}

export interface ShareToCoachRequest {
  session_id: number;
  coach_id: string;
  message?: string;
  include_summary: boolean;
}

export interface ShareToCoachResponse {
  success: boolean;
  message_id?: number;
  message: string;
}
