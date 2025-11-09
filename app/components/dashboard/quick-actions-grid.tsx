import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

interface QuickAction {
  icon: keyof typeof Ionicons.glyphMap;
  label: string;
  route: string;
  color: string;
  gradient?: [string, string];
}

const quickActions: QuickAction[] = [
  {
    icon: 'play-circle',
    label: 'Start Workout',
    route: '/(user)/plans',
    color: COLORS.primary,
  },
  {
    icon: 'calendar-outline',
    label: 'View Plans',
    route: '/(user)/plans',
    color: COLORS.info,
  },
  {
    icon: 'stats-chart',
    label: 'Track Progress',
    route: '/(user)/progress',
    color: COLORS.success,
  },
  {
    icon: 'restaurant-outline',
    label: 'Nutrition',
    route: '/(user)/nutrition',
    color: COLORS.warning,
  },
  {
    icon: 'chatbubble-outline',
    label: 'Message Coach',
    route: '/(user)/conversations',
    color: COLORS.error,
  },
  {
    icon: 'fitness-outline',
    label: 'Mindfulness',
    route: '/(user)/mindfullness',
    color: '#8B5CF6',
  },
];

interface QuickActionCardProps {
  action: QuickAction;
  index: number;
  onPress: () => void;
}

const QuickActionCard: React.FC<QuickActionCardProps> = ({ action, index, onPress }) => (
  <MotiView
    from={{ opacity: 0, scale: 0.8 }}
    animate={{ opacity: 1, scale: 1 }}
    transition={{ type: 'spring', delay: index * 50, damping: 15 }}
    style={styles.actionCardContainer}
  >
    <TouchableOpacity
      style={[styles.actionCard, { backgroundColor: `${action.color}10` }]}
      onPress={onPress}
      activeOpacity={0.7}
    >
      <View style={[styles.actionIconContainer, { backgroundColor: action.color }]}>
        <Ionicons name={action.icon} size={24} color={COLORS.white} />
      </View>
      <Text style={styles.actionLabel} numberOfLines={2}>
        {action.label}
      </Text>
    </TouchableOpacity>
  </MotiView>
);

export const QuickActionsGrid: React.FC = () => {
  const router = useRouter();

  const handleActionPress = (route: string) => {
    router.push(route as any);
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Quick Actions</Text>
      <View style={styles.grid}>
        {quickActions.map((action, index) => (
          <QuickActionCard
            key={action.label}
            action={action}
            index={index}
            onPress={() => handleActionPress(action.route)}
          />
        ))}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: SPACING.xl,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  actionCardContainer: {
    width: '31%', // 3 columns with gaps
    minWidth: 100,
  },
  actionCard: {
    aspectRatio: 1,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.sm,
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    ...SHADOWS.sm,
  },
  actionIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.xs,
  },
  actionLabel: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
    textAlign: 'center',
  },
});
