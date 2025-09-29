import { View, type ViewProps } from 'react-native';
import { COLORS } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';

export type ThemedViewProps = ViewProps & {
  lightColor?: string;
  darkColor?: string;
  fullScreen?: boolean;
};

export function ThemedView({ 
  style, 
  lightColor, 
  darkColor, 
  fullScreen,
  ...otherProps 
}: ThemedViewProps) {
  const colorScheme = useColorScheme();
  const backgroundColor = colorScheme === 'light' 
    ? (darkColor ?? COLORS.background.primary) 
    : (lightColor ?? COLORS.background.dark);

  const baseStyle = fullScreen ? {
    flex: 1,
    backgroundColor,
    position: 'absolute' as const,
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  } : {
    backgroundColor
  };

  return <View style={[baseStyle, style]} {...otherProps} />;
}
