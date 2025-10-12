import { 
  COLORS, 
  SPACING, 
  FONT_SIZES, 
  FONT_WEIGHTS, 
  BORDER_RADIUS, 
  SHADOWS,
  getColorWithOpacity 
} from "@/constants/theme";
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
                        <Ionicons name="close" size={18} color={COLORS.text.inverse} />
                    </TouchableOpacity>
                )}
            </View>
        </Animated.View>
    );
}

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    left: SPACING.base,
    right: SPACING.base,
    marginHorizontal: 'auto',
    maxWidth: screenWidth - (SPACING.base * 2),
    minHeight: 56,
    borderRadius: BORDER_RADIUS.full,
    ...SHADOWS.lg,
    zIndex: 1000,
  },
  
  topPosition: {
    top: SPACING['5xl'],
  },
  bottomPosition: {
    bottom: SPACING['6xl'],
  },
  centerPosition: {
    top: '50%',
    marginTop: -28,
  },
  
  content: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.md,
    minHeight: 56,
  },
  
  icon: {
    marginRight: SPACING.md,
  },
  
  message: {
    flex: 1,
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    lineHeight: 20,
  },
  
  actionButton: {
    marginLeft: SPACING.md,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs + 2,
    borderRadius: BORDER_RADIUS.sm,
    backgroundColor: COLORS.background.secondary,
    borderWidth: 1,
    borderColor: COLORS.border.light,
  },
  
  actionButtonText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  
  closeButton: {
    marginLeft: SPACING.sm,
    padding: SPACING.xs,
  },
  
  errorContainer: {
    backgroundColor: COLORS.background.errorSoft,
    borderColor: COLORS.error,
    borderWidth: 0.5,
  },
  successContainer: {
    backgroundColor: COLORS.background.successSoft,

    borderColor: COLORS.success,
    borderWidth: 0.5,
  },
  infoContainer: {
    backgroundColor: COLORS.background.infoSoft,

    borderColor: COLORS.info,
    borderWidth: 0.5  ,
  },
  warningContainer: {
    backgroundColor: COLORS.background.warningSoft,

    borderColor: COLORS.warning,
    borderWidth: 0.5,
  },
  
  errorText: {
    color: COLORS.text.sc.error,
  },
  successText: {
    color: COLORS.text.sc.success,
  },
  infoText: {
    color: COLORS.text.sc.info ,
  },
  warningText: {
    color: COLORS.text.sc.warning,
  },
});
