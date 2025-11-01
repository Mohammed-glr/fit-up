import { useQuery } from '@tanstack/react-query';
import { workoutService } from '@/api/services/schema-service';
import type { Workout, WorkoutExerciseDetail } from '@/types/schema';
import { APIError } from '@/api/client';

export const workoutKeys = {
  all: ['workouts'] as const,
  detail: (id: number) => [...workoutKeys.all, 'detail', id] as const,
  exercises: (id: number) => [...workoutKeys.all, 'exercises', id] as const,
};


export const useWorkout = (workoutID: number) => {
  return useQuery<Workout, APIError>({
    queryKey: workoutKeys.detail(workoutID),
    queryFn: () => workoutService.Retrieve(workoutID),
    enabled: !!workoutID && workoutID > 0,
    staleTime: 10 * 60 * 1000, 
  });
};


export const useWorkoutExercises = (workoutID: number) => {
  return useQuery<WorkoutExerciseDetail[], APIError>({
    queryKey: workoutKeys.exercises(workoutID),
    queryFn: () => workoutService.GetWorkoutExercises(workoutID),
    enabled: !!workoutID && workoutID > 0,
    staleTime: 10 * 60 * 1000, 
  });
};
