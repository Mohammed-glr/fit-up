/**
 * Types for User Workout Templates feature
 * Allows users to save and reuse workout configurations (different from schema templates)
 */

export interface TemplateExercise {
  exercise_name: string;
  sets: number;
  target_reps: string;
  target_weight?: number;
  rest_seconds: number;
}

export interface UserWorkoutTemplate {
  template_id: string;
  user_id: string;
  name: string;
  description: string;
  is_public: boolean;
  exercises: TemplateExercise[];
  created_at: string;
  updated_at: string;
}

export interface CreateUserTemplateRequest {
  name: string;
  description?: string;
  is_public?: boolean;
  exercises: TemplateExercise[];
}

export interface UpdateUserTemplateRequest {
  name?: string;
  description?: string;
  is_public?: boolean;
  exercises?: TemplateExercise[];
}

export interface UserTemplateListResponse {
  templates: UserWorkoutTemplate[];
  total_count: number;
  has_more: boolean;
}

export interface UserTemplateListParams {
  page?: number;
  page_size?: number;
}

