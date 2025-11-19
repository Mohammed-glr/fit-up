import { useQuery } from '@tanstack/react-query';
import { httpClient } from '@/api/client';

export interface SchemaWorkout {
  workout_id: number;
  schema_id: number;
  day_of_week: number;
  focus: string;
  exercises: WorkoutExerciseDetail[];
}

export interface WorkoutExerciseDetail {
  we_id: number;
  sets: number;
  reps: string;
  rest_seconds: number;
  exercise: {
    exercise_id: number;
    name: string;
    muscle_groups: string[];
    difficulty: string;
    equipment: string;
    type: string;
  };
}

export interface SchemaExercise {
  exercise_id: number;
  name: string;
  sets: number;
  reps: string;
  rest_seconds: number;
  notes?: string;
}

export interface UserSchema {
  schema_id: number;
  user_id: string;
  week_start: string;
  active: boolean;
  workouts?: SchemaWorkout[];
}

interface SchemaWithMetadata extends UserSchema {
  total_workouts?: number;
}


export const useUserSchemas = (userId: string) => {
  return useQuery<SchemaWithMetadata[]>({
    queryKey: ['schemas', 'user', userId],
    queryFn: async () => {
      const response = await httpClient.get(`/schemas/user/${userId}`);
      const data = response.data?.schemas || response.data;
      const schemas = Array.isArray(data) ? data : [];
      
      // Add total_workouts count if not present
      return schemas.map((schema: UserSchema) => ({
        ...schema,
        total_workouts: schema.workouts?.length || 0,
      }));
    },
    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
    gcTime: 15 * 60 * 1000,
  });
};

export const useSchemaWithWorkouts = (schemaId: number) => {
  return useQuery<{
    schema_id: number;
    user_id: string;
    week_start: string;
    active: boolean;
    workouts: SchemaWorkout[];
  }>({
    queryKey: ['schemas', schemaId, 'workouts'],
    queryFn: async () => {
      const response = await httpClient.get(`/schemas/${schemaId}/workouts`);
      return response.data;
    },
    enabled: !!schemaId,
    staleTime: 5 * 60 * 1000,
    gcTime: 15 * 60 * 1000,
  });
};


export const useActiveSchema = (userId: string) => {
  return useQuery<SchemaWithMetadata | null>({
    queryKey: ['schemas', 'user', userId, 'active'],
    queryFn: async () => {
      const response = await httpClient.get(`/schemas/user/${userId}`);
      const data = response.data?.schemas || response.data;
      const schemas = Array.isArray(data) ? data : [];
      const activeSchema = schemas.find((s: UserSchema) => s.active);
      
      if (!activeSchema) return null;
      
      return {
        ...activeSchema,
        total_workouts: activeSchema.workouts?.length || 0,
      };
    },
    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
    gcTime: 15 * 60 * 1000,
  });
};
