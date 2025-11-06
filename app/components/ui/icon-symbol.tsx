// Beautiful modern icons using Ionicons, FontAwesome, and Material Icons

import Ionicons from '@expo/vector-icons/Ionicons';
import FontAwesome5 from '@expo/vector-icons/FontAwesome5';
import MaterialCommunityIcons from '@expo/vector-icons/MaterialCommunityIcons';
import { SymbolWeight } from 'expo-symbols';
import { ComponentProps } from 'react';
import { OpaqueColorValue, type StyleProp, type TextStyle } from 'react-native';

type IoniconsName = ComponentProps<typeof Ionicons>['name'];
type FontAwesome5Name = ComponentProps<typeof FontAwesome5>['name'];
type MaterialCommunityName = ComponentProps<typeof MaterialCommunityIcons>['name'];

type IconMapping = {
  name: IoniconsName | FontAwesome5Name | MaterialCommunityName;
  library: 'ionicons' | 'fontawesome' | 'material-community';
};

type IconSymbolName = keyof typeof MAPPING;

/**
 * Modern icon mappings using Ionicons (primary), FontAwesome5, and Material Community Icons
 * These icon libraries provide much more beautiful and varied icons than basic Material Icons
 */
const MAPPING = {
  // Original mappings
  'house.fill': { name: 'home', library: 'ionicons' },
  'paperplane.fill': { name: 'send', library: 'ionicons' },
  'chevron.left.forwardslash.chevron.right': { name: 'code-slash', library: 'ionicons' },
  'chevron.right': { name: 'chevron-forward', library: 'ionicons' },
  
  // Tab bar icons - Beautiful modern alternatives
  'chart.bar.fill': { name: 'stats-chart', library: 'ionicons' },
  'calendar': { name: 'calendar', library: 'ionicons' },
  'message.fill': { name: 'chatbubbles', library: 'ionicons' },
  'brain.head.profile': { name: 'brain', library: 'material-community' },
  'person.crop.circle.fill': { name: 'person-circle', library: 'ionicons' },
  
  // Alternative beautiful icons you can use
  'dashboard-alt': { name: 'pie-chart', library: 'ionicons' },
  'analytics': { name: 'analytics', library: 'ionicons' },
  'pulse': { name: 'pulse', library: 'ionicons' },
  'fitness': { name: 'fitness', library: 'ionicons' },
  'meditation': { name: 'flower-outline', library: 'ionicons' },
  'spa': { name: 'spa', library: 'material-community' },
  'heart-pulse': { name: 'heart-pulse', library: 'material-community' },
  'chat': { name: 'chatbox-ellipses', library: 'ionicons' },
  
  // Additional common icons
  'speedometer': { name: 'speedometer', library: 'ionicons' },
  'square.grid.2x2.fill': { name: 'grid', library: 'ionicons' },
  'bubble.left.and.bubble.right.fill': { name: 'chatbubbles', library: 'ionicons' },
  'wind': { name: 'leaf-outline', library: 'ionicons' },
  'dumbbell.fill': { name: 'barbell', library: 'ionicons' },
  'nutrition': { name: 'nutrition', library: 'ionicons' },
} as Record<string, IconMapping>;

/**
 * A versatile icon component that uses the best icon libraries available
 * Supports Ionicons, FontAwesome5, and Material Community Icons
 * This ensures beautiful, modern icons across all platforms
 */
export function IconSymbol({
  name,
  size = 24,
  color,
  style,
}: {
  name: IconSymbolName;
  size?: number;
  color: string | OpaqueColorValue;
  style?: StyleProp<TextStyle>;
  weight?: SymbolWeight;
}) {
  const iconConfig = MAPPING[name];
  
  if (!iconConfig) {
    return <Ionicons color={color} size={size} name="help-circle-outline" style={style} />;
  }

  switch (iconConfig.library) {
    case 'fontawesome':
      return <FontAwesome5 color={color} size={size} name={iconConfig.name as FontAwesome5Name} style={style} />;
    case 'material-community':
      return <MaterialCommunityIcons color={color} size={size} name={iconConfig.name as MaterialCommunityName} style={style} />;
    case 'ionicons':
    default:
      return <Ionicons color={color} size={size} name={iconConfig.name as IoniconsName} style={style} />;
  }
}
