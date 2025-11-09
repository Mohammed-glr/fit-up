import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { MotiView } from 'moti';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';

interface DashboardGreetingProps {
  name?: string;
}

export const DashboardGreeting: React.FC<DashboardGreetingProps> = ({ name }) => {
  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Good Morning';
    if (hour < 18) return 'Good Afternoon';
    return 'Good Evening';
  };

  const getMotivationalMessage = () => {
    const messages = [
      "Let's crush today's workout! 💪",
      "Your fitness journey starts now!",
      "Every rep counts! 🔥",
      "Stronger than yesterday!",
      "Make today count! ⚡",
    ];
    const randomIndex = Math.floor(Math.random() * messages.length);
    return messages[randomIndex];
  };

  return (
    <MotiView
      from={{ opacity: 0, translateY: -20 }}
      animate={{ opacity: 1, translateY: 0 }}
      transition={{ type: 'timing', duration: 600 }}
      style={styles.container}
    >
      <Text style={styles.greeting}>{getGreeting()},</Text>
      <Text style={styles.name}>{name || 'Champion'}! 👋</Text>
      <Text style={styles.motivation}>{getMotivationalMessage()}</Text>
    </MotiView>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: SPACING.lg,
  },
  greeting: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.inverse,
  },
  name: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
    marginBottom: SPACING.xs,
  },
  motivation: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.primary,
  },
});
