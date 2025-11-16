import React, { createContext, useContext, useState, ReactNode } from 'react';
import type { ManualWorkoutRequest } from '@/types/schema';

interface WorkoutEditorContextType {
  onSaveWorkout: (() => void) | undefined;
  setOnSaveWorkout: (callback: (() => void) | undefined) => void;
  isSavingWorkout: boolean;
  setIsSavingWorkout: (isSaving: boolean) => void;
  currentWorkout: ManualWorkoutRequest | null;
  setCurrentWorkout: (workout: ManualWorkoutRequest | null) => void;
}

const WorkoutEditorContext = createContext<WorkoutEditorContextType | undefined>(undefined);

export const WorkoutEditorProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [onSaveWorkout, setOnSaveWorkout] = useState<(() => void) | undefined>(undefined);
  const [isSavingWorkout, setIsSavingWorkout] = useState(false);
  const [currentWorkout, setCurrentWorkout] = useState<ManualWorkoutRequest | null>(null);

  return (
    <WorkoutEditorContext.Provider
      value={{
        onSaveWorkout,
        setOnSaveWorkout,
        isSavingWorkout,
        setIsSavingWorkout,
        currentWorkout,
        setCurrentWorkout,
      }}
    >
      {children}
    </WorkoutEditorContext.Provider>
  );
};

export const useWorkoutEditorContext = () => {
  const context = useContext(WorkoutEditorContext);
  if (!context) {
    throw new Error('useWorkoutEditorContext must be used within WorkoutEditorProvider');
  }
  return context;
};
