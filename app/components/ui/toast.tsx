import { COLORS } from "@/constants/theme";
import Ionicons from "@expo/vector-icons/build/Ionicons";
import React from "react";

interface ToastProps {
  message: string;
  type?: 'error' | 'success' | 'info';
  isVisible: boolean;
  onClose?: () => void;
  duration?: number;
  showCloseButton?: boolean;
  style?: React.CSSProperties;
}

export const Toast: React.FC<ToastProps> = ({
    message,
    type = 'error',
    isVisible,
    onClose,
    duration = 2000,
    showCloseButton = true,
    style,
}) => {
    React.useEffect(() => {
        if (isVisible && duration > 0 && onClose) {
            const timer = setTimeout(() => {
                onClose();
            }, duration);

            return () => clearTimeout(timer);
        }
    }, [isVisible, duration, onClose]);

    if (!isVisible) return null;

    const renderIcon = () => {
        switch (type) {
            case 'error':
                return <Ionicons name="alert-circle" size={24} color={COLORS.error} />;
            case 'success':
                return <Ionicons name="checkmark-circle" size={24} color={COLORS.success} />;
            case 'info':
                return <Ionicons name="information-circle" size={24} color={COLORS.info} />;
            default:
                return null;
        }
    }; 

    return (
        <div style={style.toast}>
            {renderIcon()}
            <span>{message}</span>
            {showCloseButton && <button onClick={onClose}>Close</button>}
        </div>
    );
}

const styles = {
  toastContainer: {
    display: 'flex',
    alignItems: 'center',
    padding: '10px 20px',
    borderRadius: 5,
    boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
    maxWidth: 400,
    margin: '10px auto',
    zIndex: 1000,
  },
  toastContent: {
    flex: 1,
    marginLeft: 10,
  },
    toastIcon: {
    marginRight: 10,
  },
    toastMessage: {
    flex: 1,
    marginLeft: 10,
  },
    toastClose: {
    background: 'none',
    border: 'none',
    cursor: 'pointer',
    padding: 0,
    marginLeft: 10,
  },
    error: {
    color: COLORS.error,
  },
    success: {
    color: COLORS.success,
  },
    info: {
    color: COLORS.info,
  },
    fixed: {
    position: 'fixed',
    top: 20,
    left: '50%',
    transform: 'translateX(-50%)',
    width: '90%',
    maxWidth: 400,
  },
    relative: {
    position: 'relative',
    width: '100%',
  },
};
