import React, { createContext, useContext, useState, ReactNode } from 'react';

interface MindfulnessContextType {
  gratitudeMode: 'list' | 'create';
  setGratitudeMode: (mode: 'list' | 'create') => void;
  reflectionMode: 'main' | 'history';
  setReflectionMode: (mode: 'main' | 'history') => void;
  triggerGratitudeCreate: () => void;
  triggerReflectionHistory: () => void;
  isGratitudeWriting: boolean;
  setIsGratitudeWriting: (writing: boolean) => void;
  isReflectionResponding: boolean;
  setIsReflectionResponding: (responding: boolean) => void;
  isReflectionHistory: boolean;
  setIsReflectionHistory: (history: boolean) => void;
  onSaveGratitude?: () => void;
  setOnSaveGratitude: (callback: () => void) => void;
  onSaveReflection?: () => void;
  setOnSaveReflection: (callback: () => void) => void;
  isSavingGratitude: boolean;
  setIsSavingGratitude: (saving: boolean) => void;
  isSavingReflection: boolean;
  setIsSavingReflection: (saving: boolean) => void;
}

const MindfulnessContext = createContext<MindfulnessContextType | undefined>(undefined);

export const MindfulnessProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [gratitudeMode, setGratitudeMode] = useState<'list' | 'create'>('list');
  const [reflectionMode, setReflectionMode] = useState<'main' | 'history'>('main');
  const [isGratitudeWriting, setIsGratitudeWriting] = useState(false);
  const [isReflectionResponding, setIsReflectionResponding] = useState(false);
  const [isReflectionHistory, setIsReflectionHistory] = useState(false);
  const [onSaveGratitude, setOnSaveGratitude] = useState<(() => void) | undefined>();
  const [onSaveReflection, setOnSaveReflection] = useState<(() => void) | undefined>();
  const [isSavingGratitude, setIsSavingGratitude] = useState(false);
  const [isSavingReflection, setIsSavingReflection] = useState(false);

  const triggerGratitudeCreate = () => {
    setGratitudeMode('create');
  };

  const triggerReflectionHistory = () => {
    setReflectionMode('history');
  };

  return (
    <MindfulnessContext.Provider
      value={{
        gratitudeMode,
        setGratitudeMode,
        reflectionMode,
        setReflectionMode,
        triggerGratitudeCreate,
        triggerReflectionHistory,
        isGratitudeWriting,
        setIsGratitudeWriting,
        isReflectionResponding,
        setIsReflectionResponding,
        isReflectionHistory,
        setIsReflectionHistory,
        onSaveGratitude,
        setOnSaveGratitude,
        onSaveReflection,
        setOnSaveReflection,
        isSavingGratitude,
        setIsSavingGratitude,
        isSavingReflection,
        setIsSavingReflection,
      }}
    >
      {children}
    </MindfulnessContext.Provider>
  );
};

export const useMindfulnessContext = () => {
  const context = useContext(MindfulnessContext);
  if (!context) {
    throw new Error('useMindfulnessContext must be used within MindfulnessProvider');
  }
  return context;
};
