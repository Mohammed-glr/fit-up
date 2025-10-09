

type FitnessLevel = 'beginner' | 'intermediate' | 'advanced';
type FitnessGoal = 'strength' | 'muscle_gain' | 'fat_loss' | 'endurance' | 'general_fitness';
type ExerciseType = 'strength' | 'cardio' | 'mobility' | 'hiit' | 'stretching';
type EquipmentType = 'barbell' | 'dumbbell' | 'bodyweight' | 'machine' | 'kettlebell' | 'resistance_band';
type SessionStatus = 'active' | 'completed' | 'skipped' | 'abandoned';


interface UserDisplayInfo {
  auth_user_id: string;
  first_name: string;
  last_name: string;
  email: string;
}

interface CoachAssignment {
  assignment_id: number;
  coach_id: string;
  user_id: string;
  assigned_at: string;
  assigned_by: string;
  is_active: boolean;
  deactivated_at?: string | null;
  notes: string;
}

interface CoachAssignmentRequest {
  coach_id: string;
  user_id: string;
  notes?: string;
}

interface ClientSummary {
  user_id: number;
  auth_id: string;
  first_name: string;
  last_name: string;
  email: string;
  assigned_at: string;
  current_schema_id?: number;
  active_goals: number;
  completion_rate: number;
  last_workout_date?: string | null;
  total_workouts: number;
  current_streak: number;
  fitness_level: string;
}

interface CoachDashboard {
  coach_id: string;
  total_clients: number;
  active_clients: number;
  active_schemas: number;
  total_workouts: number;
  average_completion: number;
  clients: ClientSummary[];
}


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

interface ExerciseRequest {
  name: string;
  muscle_groups: string[];
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

interface WorkoutTemplate {
  template_id: number;
  name: string;
  description: string;
  min_level: FitnessLevel;
  max_level: FitnessLevel;
  suitable_goals: string;
  days_per_week: number;
}

interface WorkoutTemplateRequest {
  name: string;
  description?: string;
  min_level: FitnessLevel;
  max_level: FitnessLevel;
  suitable_goals: FitnessGoal[];
  days_per_week: number;
}

interface WorkoutTemplateResponse {
  template_id: number;
  name: string;
  description: string;
  min_level: FitnessLevel;
  max_level: FitnessLevel;
  suitable_goals: FitnessGoal[];
  days_per_week: number;
}

interface WeeklySchema {
  schema_id: number;
  user_id: number;
  week_start: string;
  active: boolean;
}

interface WeeklySchemaRequest {
  user_id: number;
  week_start: string;
}

interface Workout {
  workout_id: number;
  schema_id: number;
  day_of_week: number;
  focus: string;
}

interface WorkoutRequest {
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

interface WorkoutExerciseRequest {
  workout_id: number;
  exercise_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
}

interface ProgressLog {
  log_id: number;
  user_id: number;
  exercise_id: number;
  date: string;
  sets_completed?: number | null;
  reps_completed?: number | null;
  weight_used?: number | null;
  duration_seconds?: number | null;
}

interface ProgressLogRequest {
  user_id: number;
  exercise_id: number;
  date: string;
  sets_completed?: number;
  reps_completed?: number;
  weight_used?: number;
  duration_seconds?: number;
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

interface PersonalBest {
  exercise_id: number;
  exercise_name: string;
  best_weight?: number | null;
  best_reps?: number | null;
  best_volume?: number | null;
  achieved_at: string;
}

interface UserProgressSummary {
  user_id: number;
  total_workouts: number;
  current_streak: number;
  last_workout?: string | null;
  personal_bests: PersonalBest[];
}

interface ExerciseFilter {
  muscle_groups?: string[];
  difficulty?: FitnessLevel;
  equipment?: EquipmentType[];
  type?: ExerciseType[];
  search?: string;
}

interface TemplateFilter {
  level?: FitnessLevel;
  goals?: FitnessGoal[];
  days_per_week?: number;
  search?: string;
}

interface ProgressFilter {
  user_id: number;
  exercise_id?: number;
  date_from?: string;
  date_to?: string;
}

interface PaginationSchemaParams {
  offset?: number;
  limit?: number;
  page?: number;
  page_size?: number;
}

interface PaginatedSchemaResponse<T> {
  data: T[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}



interface OneRepMaxEstimate {
  estimate_id: number;
  user_id: number;
  exercise_id: number;
  estimated_max: number;
  estimate_date: string;
  method: string;
  confidence: number;
}

interface PerformanceData {
  weight: number;
  reps: number;
  sets: number;
  rpe: number;
  duration_seconds: number;
}

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

interface PlanPerformanceData {
  completion_rate: number;
  average_rpe: number;
  progress_rate: number;
  user_satisfaction: number;
  injury_rate: number;
}

interface PlanAdaptation {
  adaptation_id: number;
  plan_id: number;
  adaptation_date: string;
  reason: string;
  changes: any;
  trigger: string;
}


interface WeeklySessionStats {
  week_start: string;
  sessions_planned: number;
  sessions_completed: number;
  total_volume: number;
  average_rpe: number;
  completion_rate: number;
}

interface GoalProgress {
  goal_id: number;
  progress_percent: number;
  on_track: boolean;
  estimated_completion?: string | null;
}

interface TimeToGoalEstimate {
  goal_id: number;
  estimated_days: number;
  confidence: number;
  assumptions: string[];
}

interface GoalAdjustment {
  goal_id: number;
  recommendation_type: string;
  adjustment: string;
  reason: string;
}

interface TrainingHistory {
  total_workouts: number;
  weeks_active: number;
  average_frequency: number;
  consistency_score: number;
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

interface SchemaMetadata {
  created_by: string;
  creator_id: string;
  is_custom: boolean;
  base_template_id?: number | null;
  last_modified_by: string;
  modified_at?: string | null;
  version: number;
  tags: string[];
  custom_data: Record<string, any>;
}

interface WorkoutDetail {
  workout_id: number;
  schema_id: number;
  day_of_week: number;
  focus: string;
  exercises: WorkoutExerciseDetail[];
  estimated_minutes: number;
  notes: string;
}

interface WeeklySchemaExtended extends WeeklySchema {
  coach_id?: string;
  coach_name?: string;
  metadata: SchemaMetadata;
  workouts: WorkoutDetail[];
}

export type {
  FitnessLevel,
  FitnessGoal,
  ExerciseType,
  EquipmentType,
  SessionStatus,
  UserDisplayInfo,
  CoachAssignment,
  CoachAssignmentRequest,
  ClientSummary,
  CoachDashboard,
  
  Exercise,
  ExerciseRequest,
  ExerciseResponse,
  
  WorkoutTemplate,
  WorkoutTemplateRequest,
  WorkoutTemplateResponse,
  
  WeeklySchema,
  WeeklySchemaRequest,
  Workout,
  WorkoutRequest,
  WorkoutExercise,
  WorkoutExerciseRequest,
  
  ProgressLog,
  ProgressLogRequest,
  
  WorkoutExerciseDetail,
  WorkoutWithExercises,
  WeeklySchemaWithWorkouts,
  PersonalBest,
  UserProgressSummary,
  
  ExerciseFilter,
  TemplateFilter,
  ProgressFilter,
  
  PaginationSchemaParams,
  PaginatedSchemaResponse,
  APIResponse,
  

  
  OneRepMaxEstimate,
  PerformanceData,

  
  GeneratedPlan,
  PlanGenerationMetadata,
  PlanPerformanceData,
  PlanAdaptation,
  WeeklySessionStats,
  
  GoalProgress,
  TimeToGoalEstimate,
  GoalAdjustment,
  TrainingHistory,
  
  ManualExerciseRequest,
  ManualWorkoutRequest,
  ManualSchemaRequest,
  SchemaMetadata,
  WorkoutDetail,
  WeeklySchemaExtended,
};
