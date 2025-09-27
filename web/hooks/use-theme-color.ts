/**
 * Learn more about light and dark modes:
 * https://docs.expo.dev/guides/color-schemes/
 */

import { COLORS } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';

export function useThemeColor(
  props: { light?: string; dark?: string },
  colorName: keyof typeof COLORS.text | keyof typeof COLORS.background | 'text' | 'background' | 'tint' | 'icon' | 'tabIconDefault' | 'tabIconSelected'
) {
  const theme = useColorScheme() ?? 'light';
  const colorFromProps = props[theme];

  if (colorFromProps) {
    return colorFromProps;
  } else {
    // For backwards compatibility, map common color names to our new structure
    switch (colorName as string) {
      case 'text':
        return theme === 'dark' ? COLORS.text.inverse : COLORS.text.primary;
      case 'background':
        return theme === 'dark' ? COLORS.background.dark : COLORS.background.secondary;
      case 'primary':
        return COLORS.background.primary;
      case 'secondary':
        return COLORS.background.secondary;
      case 'tint':
        return COLORS.primary;
      case 'icon':
        return theme === 'dark' ? COLORS.text.secondary : COLORS.text.tertiary;
      case 'tabIconDefault':
        return COLORS.text.tertiary;
      case 'tabIconSelected':
        return COLORS.primary;
      default:
        return COLORS.text.primary;
    }
  }
}
