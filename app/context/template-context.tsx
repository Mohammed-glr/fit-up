import React, { createContext, useContext, useState, ReactNode } from 'react';

interface TemplateContextType {
  onCreateTemplate: (() => void) | null;
  setOnCreateTemplate: (callback: (() => void) | null) => void;
}

const TemplateContext = createContext<TemplateContextType | undefined>(undefined);

export const TemplateProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [onCreateTemplate, setOnCreateTemplate] = useState<(() => void) | null>(null);

  return (
    <TemplateContext.Provider value={{ onCreateTemplate, setOnCreateTemplate }}>
      {children}
    </TemplateContext.Provider>
  );
};

export const useTemplateContext = () => {
  const context = useContext(TemplateContext);
  if (context === undefined) {
    throw new Error('useTemplateContext must be used within a TemplateProvider');
  }
  return context;
};
