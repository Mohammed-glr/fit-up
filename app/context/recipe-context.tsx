import React from 'react';

interface RecipeContextValue {
  onCreateRecipe?: () => void;
  setOnCreateRecipe: (callback: (() => void) | undefined) => void;
}

const RecipeContext = React.createContext<RecipeContextValue>({
  onCreateRecipe: undefined,
  setOnCreateRecipe: () => {},
});

export const RecipeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [onCreateRecipe, setOnCreateRecipe] = React.useState<(() => void) | undefined>();

  const value = React.useMemo(
    () => ({
      onCreateRecipe,
      setOnCreateRecipe: (callback: (() => void) | undefined) => {
        setOnCreateRecipe(() => callback);
      },
    }),
    [onCreateRecipe]
  );

  return <RecipeContext.Provider value={value}>{children}</RecipeContext.Provider>;
};

export const useRecipeContext = () => {
  const context = React.useContext(RecipeContext);
  if (!context) {
    throw new Error('useRecipeContext must be used within RecipeProvider');
  }
  return context;
};
