import React, { useState } from 'react';
import {
  SafeAreaView,
  StyleSheet,
  Text,
  View,
  ScrollView,
  TouchableOpacity,
  Dimensions,
} from 'react-native';
import { router } from 'expo-router';
import { useMindfulnessStreak, useMindfulnessStats } from '../../hooks/mindfulness/use-mindfulness';
import { BREATHING_PATTERNS } from '../../types/mindfulness';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import { IconSymbol } from '@/components/ui/icon-symbol';

const { width } = Dimensions.get('window');

export default function MindfullnessScreen() {
  const { data: streak } = useMindfulnessStreak();
  const { data: stats } = useMindfulnessStats();

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          <View style={styles.header}>
            <View style={styles.titleContainer}>
              <IconSymbol name="figure.mind.and.body" size={32} color={COLORS.ms} />
              <Text style={styles.title}>Mindfulness</Text>
            </View>
            <Text style={styles.subtitle}>Find your inner calm</Text>
          </View>

          {streak && (
            <View style={styles.streakCard}>
              <View style={styles.streakRow}>
                <View style={styles.streakItem}>
                  <View style={styles.streakNumberContainer}>
                    <IconSymbol name="flame.fill" size={20} color={COLORS.ms} />
                    <Text style={styles.streakNumber}>{streak.current_streak}</Text>
                  </View>
                  <Text style={styles.streakLabel}>Day Streak</Text>
                </View>
                <View style={styles.streakDivider} />
                <View style={styles.streakItem}>
                  <View style={styles.streakNumberContainer}>
                    <IconSymbol name="trophy.fill" size={20} color={COLORS.ms} />
                    <Text style={styles.streakNumber}>{streak.longest_streak}</Text>
                  </View>
                  <Text style={styles.streakLabel}>Best Streak</Text>
                </View>
                <View style={styles.streakDivider} />
                <View style={styles.streakItem}>
                  <Text style={styles.streakNumber}>{streak.total_sessions}</Text>
                  <Text style={styles.streakLabel}>Total Sessions</Text>
                </View>
              </View>
            </View>
          )}

          {/* Quick Stats */}
          {stats && (
            <View style={styles.statsCard}>
              <Text style={styles.statsTitle}>This Week</Text>
              <View style={styles.statsRow}>
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>{stats.total_sessions}</Text>
                  <Text style={styles.statLabel}>Sessions</Text>
                </View>
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>{stats.total_minutes}</Text>
                  <Text style={styles.statLabel}>Minutes</Text>
                </View>
              </View>
            </View>
          )}

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Breathing Exercises</Text>
            <Text style={styles.sectionSubtitle}>
              Guided breathing patterns to calm your mind
            </Text>

            {Object.values(BREATHING_PATTERNS).map((pattern) => (
              <TouchableOpacity
                key={pattern.type}
                style={styles.card}
                onPress={() => router.push(`/(user)/breathing/${pattern.type}` as any)}
                activeOpacity={0.7}
              >
                <View style={styles.cardHeader}>
                  <Text style={styles.cardTitle}>{pattern.name}</Text>
                  <Text style={styles.cardDuration}>
                    {Math.floor(pattern.duration / 60)} min
                  </Text>
                </View>
                <Text style={styles.cardDescription}>
                  {pattern.description}
                </Text>
                <View style={styles.benefitsContainer}>
                  {pattern.benefits.slice(0, 2).map((benefit, index) => (
                    <View key={index} style={styles.benefitTag}>
                      <Text style={styles.benefitText}>{benefit}</Text>
                    </View>
                  ))}
                </View>
              </TouchableOpacity>
            ))}
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Gratitude Journal</Text>
            <Text style={styles.sectionSubtitle}>
              Reflect on what you're grateful for
            </Text>
            <TouchableOpacity
              style={styles.card}
              onPress={() => router.push('/(user)/gratitude')}
              activeOpacity={0.7}
            >
              <View style={styles.cardHeader}>
                <IconSymbol name="square.and.pencil" size={24} color={COLORS.ms} />
                <Text style={styles.cardTitle}>Write Entry</Text>
              </View>
              <Text style={styles.cardDescription}>
                Take a moment to express gratitude
              </Text>
            </TouchableOpacity>
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Daily Reflection</Text>
            <Text style={styles.sectionSubtitle}>
              Answer thoughtful prompts for self-discovery
            </Text>
            <TouchableOpacity
              style={styles.card}
              onPress={() => router.push('/(user)/reflection')}
              activeOpacity={0.7}
            >
              <View style={styles.cardHeader}>
                <IconSymbol name="text.bubble.fill" size={24} color={COLORS.ms} />
                <Text style={styles.cardTitle}>Today's Prompt</Text>
              </View>
              <Text style={styles.cardDescription}>
                Explore your thoughts and feelings
              </Text>
            </TouchableOpacity>
          </View>

          <View style={styles.bottomSpacer} />
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0A0A0A',
  },
  safeArea: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.base,
  },
  header: {
    marginBottom: 24,
  },
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 4,
  },
  title: {
    fontSize: FONT_SIZES['3xl'],
    fontWeight: '700',
    color: '#FFFFFF',
  },
  subtitle: {
    fontSize: FONT_SIZES.sm,
    color: '#888888',
  },
  streakCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: 16,
  },
  streakRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  streakItem: {
    flex: 1,
    alignItems: 'center',
  },
  streakNumberContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    marginBottom: SPACING.xs,
  },
  streakNumber: {
    fontSize: 24,
    fontWeight: '700',
    color: COLORS.text.inverse,
  },
  streakLabel: {
    fontSize: 12,
    color: COLORS.text.tertiary,
  },
  streakDivider: {
    width: 1,
    height: 40,
    backgroundColor: '#2A2A2A',
  },
  statsCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: 24,
  },
  statsTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.inverse,
    marginBottom: SPACING.md,
  },
  statsRow: {
    flexDirection: 'row',
    gap: 20,
  },
  statItem: {
    flex: 1,
  },
  statValue: {
    fontSize: 28,
    fontWeight: '700',
    color: '#6C63FF',
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 14,
    color: '#888888',
  },
  section: {
    marginBottom: 32,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFFFFF',
    marginBottom: 4,
  },
  sectionSubtitle: {
    fontSize: 14,
    color: '#888888',
    marginBottom: 16,
  },
  card: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: 12,
  },
  cardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 8,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  cardDuration: {
    fontSize: 14,
    fontWeight: '500',
    color: '#6C63FF',
  },
  cardDescription: {
    fontSize: 14,
    color: '#888888',
    marginBottom: 12,
  },
  benefitsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  benefitTag: {
    backgroundColor: '#2A2A2A',
    borderRadius: BORDER_RADIUS.full,
    paddingVertical: 6,
    paddingHorizontal: 12,
  },
  benefitText: {
    fontSize: 12,
    color: '#6C63FF',
    fontWeight: '500',
  },
  bottomSpacer: {
    height: 40,
  },
});

