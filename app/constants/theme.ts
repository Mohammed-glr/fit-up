/**
 * FitUp App Theme Configuration
 * 
 * This file contains the complete design system for the FitUp app including:
 * - Colors (Primary green accent, Secondary black/white, Support states)
 * - Typography (Font families, sizes, weights)
 * - Spacing (Margins, paddings, gaps)
 * - Border radius, shadows, and other design tokens
 */

import { Platform } from 'react-native';

// ========================================
// COLOR PALETTE
// ========================================

// Primary Colors (Green Accent)
export const PRIMARY_COLORS = {
  primary: '#8FE507',        // Main green accent
  primaryDark: '#6AB000',    // Darker shade for hover/pressed states
  primaryLight: '#B5F542',   // Lighter tint for highlights
  primarySoft: '#E8FFD1',    // Soft tint for backgrounds
} as const;

// Secondary Colors (Modern Light Theme)
export const SECONDARY_COLORS = {
  black: '#1F2937',          // Softer dark gray instead of pure black
  darkGray: '#374151',       // Modern dark gray for secondary text
  mediumGray: '#6B7280',     // Lighter medium gray for borders/dividers
  lightGray: '#9CA3AF',      // Softer light gray for placeholders
  white: '#FFFFFF',          // Pure white
  offWhite: '#F9FAFB',       // Very light gray for subtle backgrounds
} as const;

// Support/State Colors
export const SUPPORT_COLORS = {
  success: '#4CAF50',        // Success messages (green family)
  warning: '#FFC107',        // Warning messages
  error: '#F44336',          // Error messages
  info: '#2196F3',           // Info messages
} as const;

// Combined Colors Object
export const COLORS = {
  ...PRIMARY_COLORS,
  ...SECONDARY_COLORS,
  ...SUPPORT_COLORS,
  
  // Semantic Color Mappings
  text: {
    primary: SECONDARY_COLORS.black,
    secondary: SECONDARY_COLORS.darkGray,
    tertiary: SECONDARY_COLORS.mediumGray,
    placeholder: SECONDARY_COLORS.lightGray,
    inverse: SECONDARY_COLORS.white,
  },
  background: {
    primary: SECONDARY_COLORS.white,
    secondary: SECONDARY_COLORS.offWhite,
    accent: PRIMARY_COLORS.primarySoft,
    dark: SECONDARY_COLORS.darkGray,
    surface: '#FFFFFF',
    card: '#FFFFFF',
    gradient: {
      start: '#F9FAFB',
      end: '#FFFFFF',
    },
  },
  border: {
    light: SECONDARY_COLORS.lightGray,
    medium: SECONDARY_COLORS.mediumGray,
    dark: SECONDARY_COLORS.darkGray,
    accent: PRIMARY_COLORS.primary,
    subtle: '#E5E7EB',
  },
  surface: {
    elevated: '#FFFFFF',
    overlay: 'rgba(0, 0, 0, 0.5)',
    backdrop: 'rgba(0, 0, 0, 0.3)',
  },
} as const;

// ========================================
// TYPOGRAPHY
// ========================================

export const FONT_FAMILIES = Platform.select({
  ios: {
    regular: 'System',
    medium: 'System-Medium',
    semiBold: 'System-SemiBold',
    bold: 'System-Bold',
    mono: 'Menlo-Regular',
  },
  android: {
    regular: 'Roboto-Regular',
    medium: 'Roboto-Medium',
    semiBold: 'Roboto-Medium',
    bold: 'Roboto-Bold',
    mono: 'RobotoMono-Regular',
  },
  web: {
    regular: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    medium: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    semiBold: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    bold: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    mono: "SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace",
  },
  default: {
    regular: 'normal',
    medium: 'normal',
    semiBold: 'normal',
    bold: 'bold',
    mono: 'monospace',
  },
}) || {
  regular: 'normal',
  medium: 'normal',
  semiBold: 'normal',
  bold: 'bold',
  mono: 'monospace',
};

export const FONT_SIZES = {
  xs: 12,
  sm: 14,
  base: 16,
  lg: 18,
  xl: 20,
  '2xl': 24,
  '3xl': 30,
  '4xl': 36,
  '5xl': 48,
  '6xl': 60,
} as const;

export const FONT_WEIGHTS = {
  normal: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: '700' as const,
} as const;

export const LINE_HEIGHTS = {
  tight: 1.2,
  normal: 1.4,
  relaxed: 1.6,
  loose: 1.8,
} as const;

// ========================================
// SPACING
// ========================================

export const SPACING = {
  xs: 4,
  sm: 8,
  md: 12,
  base: 16,
  lg: 20,
  xl: 24,
  '2xl': 32,
  '3xl': 40,
  '4xl': 48,
  '5xl': 64,
  '6xl': 80,
} as const;

// ========================================
// BORDER RADIUS
// ========================================

export const BORDER_RADIUS = {
  none: 0,
  sm: 6,
  base: 12,
  md: 16,
  lg: 20,
  xl: 24,
  '2xl': 32,
  '3xl': 40,
  full: 999,
} as const;

// ========================================
// SHADOWS
// ========================================

export const SHADOWS = {
  none: {
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0,
    shadowRadius: 0,
    elevation: 0,
  },
  sm: {
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 3,
    elevation: 1,
  },
  base: {
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 6,
    elevation: 3,
  },
  lg: {
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 12,
    elevation: 6,
  },
  xl: {
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.12,
    shadowRadius: 20,
    elevation: 10,
  },
  modern: {
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 8,
    elevation: 2,
  },
} as const;

// ========================================
// COMPONENT SPECIFIC TOKENS
// ========================================

export const BUTTON = {
  height: {
    sm: 36,
    base: 48,
    lg: 56,
  },
  borderRadius: BORDER_RADIUS.md,
  fontSize: {
    sm: FONT_SIZES.sm,
    base: FONT_SIZES.base,
    lg: FONT_SIZES.lg,
  },
} as const;

export const INPUT = {
  height: 52,
  borderRadius: BORDER_RADIUS.md,
  fontSize: FONT_SIZES.base,
  borderWidth: 1.5,
} as const;

// ========================================
// THEME OBJECT
// ========================================

export const THEME = {
  colors: COLORS,
  fonts: FONT_FAMILIES,
  fontSizes: FONT_SIZES,
  fontWeights: FONT_WEIGHTS,
  lineHeights: LINE_HEIGHTS,
  spacing: SPACING,
  borderRadius: BORDER_RADIUS,
  shadows: SHADOWS,
  button: BUTTON,
  input: INPUT,
} as const;

// ========================================
// TYPE DEFINITIONS
// ========================================

export type ThemeColors = typeof COLORS;
export type FontSizes = typeof FONT_SIZES;
export type Spacing = typeof SPACING;
export type BorderRadius = typeof BORDER_RADIUS;
export type Theme = typeof THEME;

// ========================================
// USAGE EXAMPLES & HELPER FUNCTIONS
// ========================================

/**
 * Usage Guide:
 * 
 * Primary Green (#8FE507) → buttons, highlights, active states
 * Black & White → text, backgrounds, contrast
 * Grays → borders, placeholders, disabled states
 * Support Colors → feedback messages, alerts
 * 
 * Example usage in components:
 * 
 * import { COLORS, SPACING, FONT_SIZES } from '@/constants/theme';
 * 
 * const styles = StyleSheet.create({
 *   button: {
 *     backgroundColor: COLORS.primary,
 *     paddingHorizontal: SPACING.base,
 *     paddingVertical: SPACING.md,
 *     borderRadius: BORDER_RADIUS.base,
 *   },
 *   text: {
 *     color: COLORS.text.primary,
 *     fontSize: FONT_SIZES.base,
 *   },
 * });
 */

// Helper function to get color with opacity
export const getColorWithOpacity = (color: string, opacity: number): string => {
  // Remove # if present
  const hex = color.replace('#', '');
  
  // Parse r, g, b values
  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);
  
  return `rgba(${r}, ${g}, ${b}, ${opacity})`;
};

// Helper function to get shadow color
export const getShadowColor = (color: string = COLORS.black, opacity: number = 0.1): string => {
  return getColorWithOpacity(color, opacity);
};
