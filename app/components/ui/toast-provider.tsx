import React, { createContext, useContext, useCallback, useState } from 'react';
import { Toast } from './toast';
import { ViewStyle, TextStyle } from 'react-native';
import { SPACING } from '@/constants/theme';

type ToastType = 'error' | 'success' | 'info' | 'warning';
type ToastPosition = 'top' | 'bottom' | 'center';

interface ToastConfig {
  id: string;
  message: string;
  type: ToastType;
  position?: ToastPosition;
  duration?: number;
  showCloseButton?: boolean;
  style?: ViewStyle;
  textStyle?: TextStyle;
  actionButton?: {
    text: string;
    onPress: () => void;
  };
}

interface ToastContextType {
  showToast: (config: Omit<ToastConfig, 'id'>) => string;
  hideToast: (id: string) => void;
  hideAllToasts: () => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

interface ToastProviderProps {
  children: React.ReactNode;
  maxToasts?: number;
}

export const ToastProvider: React.FC<ToastProviderProps> = ({ 
  children, 
  maxToasts = 3 
}) => {
  const [toasts, setToasts] = useState<ToastConfig[]>([]);

  const showToast = useCallback((config: Omit<ToastConfig, 'id'>) => {
    const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    const newToast: ToastConfig = {
      ...config,
      id,
      position: config.position || 'top',
      duration: config.duration ?? 3000,
      showCloseButton: config.showCloseButton ?? true,
    };

    setToasts(prevToasts => {
      const updatedToasts = [...prevToasts, newToast];
      if (updatedToasts.length > maxToasts) {
        return updatedToasts.slice(-maxToasts);
      }
      return updatedToasts;
    });

    return id;
  }, [maxToasts]);

  const hideToast = useCallback((id: string) => {
    setToasts(prevToasts => prevToasts.filter(toast => toast.id !== id));
  }, []);

  const hideAllToasts = useCallback(() => {
    setToasts([]);
  }, []);

  const contextValue = React.useMemo(() => ({
    showToast,
    hideToast,
    hideAllToasts,
  }), [showToast, hideToast, hideAllToasts]);

  return (
    <ToastContext.Provider value={contextValue}>
      {children}
      {toasts.map((toast, index) => (
        <Toast
          key={toast.id}
          message={toast.message}
          type={toast.type}
          position={toast.position}
          isVisible={true}
          onClose={() => hideToast(toast.id)}
          duration={toast.duration}
          showCloseButton={toast.showCloseButton}
          style={{
            ...toast.style,
            ...(toast.position === 'top' && index > 0 && {
              top: SPACING['5xl'] + (index * SPACING.md),
            }),
            ...(toast.position === 'bottom' && index > 0 && {
              bottom: SPACING['6xl'] + (index * SPACING.md),
            }),
          }}
          textStyle={toast.textStyle}
          actionButton={toast.actionButton}
        />
      ))}
    </ToastContext.Provider>  
  );
};

export const useToastMethods = () => {
  const { showToast } = useToast();

  return {
    showError: (message: string, options?: Partial<Omit<ToastConfig, 'id' | 'message' | 'type'>>) =>
      showToast({ message, type: 'error', ...options }),
    
    showSuccess: (message: string, options?: Partial<Omit<ToastConfig, 'id' | 'message' | 'type'>>) =>
      showToast({ message, type: 'success', ...options }),
    
    showInfo: (message: string, options?: Partial<Omit<ToastConfig, 'id' | 'message' | 'type'>>) =>
      showToast({ message, type: 'info', ...options }),
    
    showWarning: (message: string, options?: Partial<Omit<ToastConfig, 'id' | 'message' | 'type'>>) =>
      showToast({ message, type: 'warning', ...options }),
  };
};