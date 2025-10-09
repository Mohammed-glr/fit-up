


type FitnessLevel = 'beginner' | 'intermediate' | 'advanced';
type FitnessGoal = 'strength' | 'muscle_gain' | 'fat_loss' | 'endurance' | 'general_fitness';
type ExerciseType = 'strength' | 'cardio' | 'mobility' | 'hiit' | 'stretching';
type EquipmentType = 'barbell' | 'dumbbell' | 'bodyweight' | 'machine' | 'kettlebell' | 'resistance_band';



interface Exercise {
  exercise_id: number;
  name: string;
  muscle_groups: string;
  difficulty: FitnessLevel;
  equipment: EquipmentType;
  type: ExerciseType;
  default_sets: number;
  default_reps: string;
  rest_seconds: number;
}

interface ExerciseResponse {
  exercise_id: number;
  name: string;
  muscle_groups: string[];
  difficulty: FitnessLevel;
  equipment: EquipmentType;
  type: ExerciseType;
  default_sets: number;
  default_reps: string;
  rest_seconds: number;
}

interface WeeklySchema {
  schema_id: number;
  user_id: number;
  week_start: string;
  active: boolean;
}

interface Workout {
  workout_id: number;
  schema_id: number;
  day_of_week: number;
  focus: string;
}

interface WorkoutExercise {
  we_id: number;
  workout_id: number;
  exercise_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
}

interface WorkoutExerciseDetail {
  we_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
  exercise: ExerciseResponse;
}

interface WorkoutWithExercises {
  workout_id: number;
  schema_id: number;
  day_of_week: number;
  focus: string;
  exercises: WorkoutExerciseDetail[];
}

interface WeeklySchemaWithWorkouts {
  schema_id: number;
  user_id: number;
  week_start: string;
  active: boolean;
  workouts: WorkoutWithExercises[];
}

// ============= WORKOUT PROFILE =============
interface WorkoutProfile {
  workout_profile_id: number;
  auth_user_id: string;
  level: FitnessLevel;
  goal: FitnessGoal;
  frequency: number;
  equipment: string[];
  created_at: string;
}

interface WorkoutProfileRequest {
  level: FitnessLevel;
  goal: FitnessGoal;
  frequency: number;
  equipment: string[];
}

// ============= PLAN GENERATION =============
interface GeneratedPlan {
  plan_id: number;
  user_id: number;
  week_start: string;
  generated_at: string;
  algorithm: string;
  effectiveness: number;
  is_active: boolean;
  metadata: any;
}

interface PlanGenerationMetadata {
  user_goals: FitnessGoal[];
  available_equipment: EquipmentType[];
  fitness_level: FitnessLevel;
  weekly_frequency: number;
  time_per_workout: number;
  algorithm: string;
  parameters: Record<string, any>;
}

interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export type {
  FitnessLevel,
  FitnessGoal,
  ExerciseType,
  EquipmentType,
  
  Exercise,
  ExerciseResponse,
  WeeklySchema,
  Workout,
  WorkoutExercise,
  WorkoutExerciseDetail,
  WorkoutWithExercises,
  WeeklySchemaWithWorkouts,
  
  WorkoutProfile,
  WorkoutProfileRequest,
  
  GeneratedPlan,
  PlanGenerationMetadata,
  
  APIResponse,
};
