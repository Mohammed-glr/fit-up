
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
  user_id: string; // Changed from number to string (auth_user_id)
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

interface GeneratedPlanExerciseDetail {
  plan_exercise_id: number;
  plan_day_id: number;
  exercise_order: number;
  exercise_id?: number | null;
  name: string;
  sets: number;
  reps: string;
  rest_seconds: number;
  notes?: string | null;
}

interface GeneratedPlanWorkout {
  workout_id: number;
  plan_id: number;
  day_index: number;
  day_title: string;
  focus: string;
  is_rest: boolean;
  exercises: GeneratedPlanExerciseDetail[];
}

interface GeneratedPlan {
  plan_id: number;
  user_id: number;
  week_start: string;
  generated_at: string;
  algorithm: string;
  effectiveness: number;
  is_active: boolean;
  metadata: Record<string, unknown> | null;
  workouts?: GeneratedPlanWorkout[];
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

interface PlanAdaptation {
  adaptation_id: number;
  plan_id: number;
  adaptation_date: string;
  reason: string;
  changes?: Record<string, unknown> | null;
  trigger: string;
}

interface PlanEffectivenessResponse {
  plan_id: number;
  effectiveness_score: number;
}

interface ClientSummary {
  user_id: number;
  auth_id: string;
  first_name: string;
  last_name: string;
  email: string;
  assigned_at: string;
  current_schema_id?: number | null;
  active_goals: number;
  completion_rate: number;
  last_workout_date?: string | null;
  total_workouts: number;
  current_streak: number;
  fitness_level: string;
}

interface UserSearchResult {
  workout_profile_id: number;
  auth_user_id: string;
  username: string;
  first_name: string;
  last_name: string;
  email: string;
  fitness_level: string;
  fitness_goal: string;
  has_coach: boolean;
  current_coach_id?: string | null;
  created_at: string;
}

interface CoachDashboard {
  coach_id: string;
  total_clients: number;
  active_clients: number;
  active_schemas: number;
  total_workouts: number;
  average_completion: number;
  clients: ClientSummary[];
  recent_activity: CoachActivity[];
}

interface CoachActivity {
  coach_id: string;
  activity_id: number;
  activity_type: string;
  user_id: number;
  user_name: string;
  description: string;
  timestamp: string;
}

interface SchemaMetadata {
  created_by: string;
  creator_id: string;
  is_custom: boolean;
  base_template_id?: number | null;
  last_modified_by: string;
  modified_at?: string | null;
  version: number;
  tags: string[];
  custom_data?: Record<string, unknown>;
}

interface WorkoutDetail extends Workout {
  exercises: WorkoutExerciseDetail[];
  estimated_minutes: number;
  notes: string;
}

interface WeeklySchemaExtended extends WeeklySchema {
  coach_id?: string | null;
  coach_name?: string;
  metadata: SchemaMetadata;
  workouts: WorkoutDetail[];
}

interface ManualExerciseRequest {
  exercise_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
  weight?: string;
  tempo?: string;
  notes?: string;
  order_index?: number;
  is_superset?: boolean;
  superset_group?: number;
}

interface ManualWorkoutRequest {
  day_of_week: number;
  workout_name: string;
  focus: string;
  notes?: string;
  estimated_minutes?: number;
  exercises: ManualExerciseRequest[];
}

interface ManualSchemaRequest {
  user_id: number;
  coach_id: string;
  name: string;
  description?: string;
  start_date: string;
  end_date?: string | null;
  is_template?: boolean;
  workouts: ManualWorkoutRequest[];
}

interface CoachAssignmentRequest {
  coach_id: string;
  username: string;
  notes?: string;
}

interface CoachAssignment {
  assignment_id: number;
  coach_id: string;
  user_id: number;
  assigned_at: string;
  assigned_by: string;
  is_active: boolean;
  deactivated_at?: string | null;
  notes?: string;
}

interface WorkoutTemplate {
  template_id: number;
  name: string;
  description: string;
  min_level: FitnessLevel;
  max_level: FitnessLevel;
  suitable_goals: string;
  days_per_week: number;
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

  ClientSummary,
  UserSearchResult,
  CoachDashboard,
  CoachActivity,
  SchemaMetadata,
  WorkoutDetail,
  WeeklySchemaExtended,
  ManualExerciseRequest,
  ManualWorkoutRequest,
  ManualSchemaRequest,
  CoachAssignmentRequest,
  CoachAssignment,
  WorkoutTemplate,
  PlanAdaptation,
  PlanEffectivenessResponse,
  GeneratedPlanWorkout,
  GeneratedPlanExerciseDetail,
};
