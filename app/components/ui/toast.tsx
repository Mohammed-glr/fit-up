import { COLORS } from "@/constants/theme";
import Ionicons from "@expo/vector-icons/build/Ionicons";
import React, { useEffect, useRef } from "react";
import { 
  View, 
  Text, 
  TouchableOpacity, 
  StyleSheet, 
  Animated,
  Dimensions,
  ViewStyle,
  TextStyle,
} from "react-native";
import * as Haptics from 'expo-haptics';

export type ToastType = 'error' | 'success' | 'info' | 'warning';
export type ToastPosition = 'top' | 'bottom' | 'center';

interface ToastProps {
  message: string;
  type?: ToastType;
  position?: ToastPosition;
  isVisible: boolean;
  onClose?: () => void;
  duration?: number;
  showCloseButton?: boolean;
  style?: ViewStyle;
  textStyle?: TextStyle;
  actionButton?: {
    text: string;
    onPress: () => void;
  };
}

const { width: screenWidth } = Dimensions.get('window');

export const Toast: React.FC<ToastProps> = ({
    message,
    type = 'error',
    position = 'top',
    isVisible,
    onClose,
    duration = 3000,
    showCloseButton = true,
    style,
    textStyle,
    actionButton,
}) => {
    const slideAnim = useRef(new Animated.Value(0)).current;
    const opacityAnim = useRef(new Animated.Value(0)).current;

    useEffect(() => {
        if (isVisible) {
            Haptics.notificationAsync(
                type === 'success' 
                    ? Haptics.NotificationFeedbackType.Success
                    : type === 'error' 
                    ? Haptics.NotificationFeedbackType.Error
                    : Haptics.NotificationFeedbackType.Warning
            );

            Animated.parallel([
                Animated.timing(slideAnim, {
                    toValue: 1,
                    duration: 300,
                    useNativeDriver: true,
                }),
                Animated.timing(opacityAnim, {
                    toValue: 1,
                    duration: 300,
                    useNativeDriver: true,
                }),
            ]).start();

            if (duration > 0 && onClose) {
                const timer = setTimeout(() => {
                    handleClose();
                }, duration);
                return () => clearTimeout(timer);
            }
        }
    }, [isVisible, duration, onClose, type]);

    const handleClose = () => {
        Animated.parallel([
            Animated.timing(slideAnim, {
                toValue: 0,
                duration: 250,
                useNativeDriver: true,
            }),
            Animated.timing(opacityAnim, {
                toValue: 0,
                duration: 250,
                useNativeDriver: true,
            }),
        ]).start(() => {
            onClose?.();
        });
    };

    const renderIcon = () => {
        const iconProps = {
            size: 20,
            style: styles.icon,
        };

        switch (type) {
            case 'error':
                return <Ionicons name="alert-circle" color={COLORS.error} {...iconProps} />;
            case 'success':
                return <Ionicons name="checkmark-circle" color={COLORS.success} {...iconProps} />;
            case 'info':
                return <Ionicons name="information-circle" color={COLORS.info} {...iconProps} />;
            case 'warning':
                return <Ionicons name="warning" color={COLORS.warning} {...iconProps} />;
            default:
                return null;
        }
    };

    const getContainerStyle = () => {
        const baseTransform = [
            {
                translateY: slideAnim.interpolate({
                    inputRange: [0, 1],
                    outputRange: position === 'top' ? [-100, 0] : position === 'bottom' ? [100, 0] : [0, 0],
                }),
            },
            {
                scale: slideAnim.interpolate({
                    inputRange: [0, 1],
                    outputRange: [0.8, 1],
                }),
            },
        ];

        return [
            styles.container,
            styles[`${position}Position`],
            styles[`${type}Container`],
            {
                opacity: opacityAnim,
                transform: baseTransform,
            },
            style,
        ];
    };

    if (!isVisible) return null;

    return (
        <Animated.View style={getContainerStyle()}>
            <View style={styles.content}>
                {renderIcon()}
                <Text style={[styles.message, styles[`${type}Text`], textStyle]} numberOfLines={3}>
                    {message}
                </Text>
                
                {actionButton && (
                    <TouchableOpacity 
                        style={styles.actionButton}
                        onPress={actionButton.onPress}
                        activeOpacity={0.7}
                    >
                        <Text style={styles.actionButtonText}>{actionButton.text}</Text>
                    </TouchableOpacity>
                )}
                
                {showCloseButton && (
                    <TouchableOpacity 
                        style={styles.closeButton}
                        onPress={handleClose}
                        activeOpacity={0.7}
                        hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
                    >
                        <Ionicons name="close" size={18} color={COLORS.mediumGray} />
                    </TouchableOpacity>
                )}
            </View>
        </Animated.View>
    );
}

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    left: 16,
    right: 16,
    marginHorizontal: 'auto',
    maxWidth: screenWidth - 32,
    minHeight: 56,
    borderRadius: 12,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 4,
    },
    shadowOpacity: 0.15,
    shadowRadius: 12,
    elevation: 8,
    zIndex: 1000,
  },
  
  topPosition: {
    top: 60,
  },
  bottomPosition: {
    bottom: 100,
  },
  centerPosition: {
    top: '50%',
    marginTop: -28,
  },
  
  content: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    minHeight: 56,
  },
  
  icon: {
    marginRight: 12,
  },
  
  message: {
    flex: 1,
    fontSize: 14,
    fontWeight: '500',
    lineHeight: 20,
  },
  
  actionButton: {
    marginLeft: 12,
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 6,
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
  },
  
  actionButtonText: {
    fontSize: 12,
    fontWeight: '600',
    color: COLORS.white,
  },
  
  closeButton: {
    marginLeft: 8,
    padding: 4,
  },
  
  errorContainer: {
    backgroundColor: COLORS.error,
  },
  successContainer: {
    backgroundColor: COLORS.success,
  },
  infoContainer: {
    backgroundColor: COLORS.info,
  },
  warningContainer: {
    backgroundColor: COLORS.warning,
  },
  
  errorText: {
    color: COLORS.white,
  },
  successText: {
    color: COLORS.white,
  },
  infoText: {
    color: COLORS.white,
  },
  warningText: {
    color: COLORS.white,
  },
});
