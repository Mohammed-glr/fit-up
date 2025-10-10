
type FitnessLevel = 'beginner' | 'intermediate' | 'advanced';
type FitnessGoal = 'strength' | 'muscle_gain' | 'fat_loss' | 'endurance' | 'general_fitness';
type ExerciseType = 'strength' | 'cardio' | 'mobility' | 'hiit' | 'stretching';
type EquipmentType = 'barbell' | 'dumbbell' | 'bodyweight' | 'machine' | 'kettlebell' | 'resistance_band';

interface ExerciseSummary {
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

interface ExerciseDetail {
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

type Exercise = ExerciseSummary;

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

interface WorkoutExerciseDetail {
  we_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
  exercise: ExerciseDetail;
}

interface WorkoutWithExercises extends Workout {
  exercises: WorkoutExerciseDetail[];
}

interface WeeklySchemaWithWorkouts extends WeeklySchema {
  workouts: WorkoutWithExercises[];
}

interface PlanGenerationMetadata {
  user_goals: FitnessGoal[];
  available_equipment: EquipmentType[];
  fitness_level: FitnessLevel;
  weekly_frequency: number;
  time_per_workout: number;
  algorithm?: string;
  parameters?: Record<string, unknown>;
}

interface CreatePlanRequest {
  user_id: number;
  metadata: PlanGenerationMetadata;
}

interface GeneratedPlan {
  plan_id: number;
  user_id: number;
  week_start: string;
  generated_at: string;
  algorithm: string;
  effectiveness: number;
  is_active: boolean;
  metadata: unknown;
}

interface PlanPerformancePayload {
  completion_rate: number;
  average_rpe: number;
  progress_rate: number;
  user_satisfaction: number;
  injury_rate: number;
}

interface PaginatedResponse<T> {
  data: T[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export type {
  FitnessLevel,
  FitnessGoal,
  ExerciseType,
  EquipmentType,

  Exercise,
  ExerciseDetail,

  WeeklySchema,
  Workout,
  WorkoutExerciseDetail,
  WorkoutWithExercises,
  WeeklySchemaWithWorkouts,

  PlanGenerationMetadata,
  CreatePlanRequest,
  GeneratedPlan,
  PlanPerformancePayload,

  PaginatedResponse,
};
