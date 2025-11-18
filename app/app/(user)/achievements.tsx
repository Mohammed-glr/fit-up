import React, { useState } from 'react';
import {
  SafeAreaView,
  StyleSheet,
  Text,
  View,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
  Dimensions,
} from 'react-native';
import { useAchievements, useAchievementStats, Achievement } from '@/hooks/user/use-achievements';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';

const { width } = Dimensions.get('window');
const CARD_WIDTH = (width - SPACING.lg * 3) / 2;

type CategoryFilter = 'all' | 'streak' | 'volume' | 'pr' | 'milestone' | 'consistency';

const CATEGORY_COLORS: Record<string, string[]> = {
  streak: ['#FF6B6B', '#FF8E53'],
  volume: ['#4ECDC4', '#44A08D'],
  pr: ['#A8E6CF', '#56AB91'],
  milestone: ['#FFD93D', '#F6C90E'],
  consistency: ['#B4A7D6', '#8E7CC3'],
};

const CATEGORY_ICONS: Record<string, keyof typeof Ionicons.glyphMap> = {
  streak: 'flame',
  volume: 'barbell',
  pr: 'trophy',
  milestone: 'star',
  consistency: 'calendar',
};

export default function AchievementsScreen() {
  const { data: achievements, isLoading } = useAchievements();
  const { data: stats, isLoading: isLoadingStats } = useAchievementStats();
  const [selectedCategory, setSelectedCategory] = useState<CategoryFilter>('all');

  const filteredAchievements = achievements?.filter(
    (a) => selectedCategory === 'all' || a.category === selectedCategory
  ) || [];

  const earnedCount = filteredAchievements.filter((a) => a.is_completed).length;
  const totalCount = filteredAchievements.length;

  if (isLoading || isLoadingStats) {
    return (
      <View style={styles.container}>
        <SafeAreaView style={styles.safeArea}>
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" color={COLORS.primary} />
            <Text style={styles.loadingText}>Loading achievements...</Text>
          </View>
        </SafeAreaView>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>

        <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false}>
        <View style={styles.header}>
                <View style={styles.titleContainer}>
                    <Ionicons name="ribbon" size={32} color={COLORS.primary} />
                <Text style={styles.headerTitle}>Achievement</Text>
                </View>
                <Text style={styles.headerSubtitle}>
                    Track your achievements and milestones!
                </Text>
          </View>
          {/* Header Stats */}
          <MotiView
            from={{ opacity: 0, translateY: -20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 400 }}
            style={styles.statsCard}
          >
            <Text style={styles.statsTitle}>Your Achievements</Text>
            <View style={styles.statsRow}>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{stats?.earned_achievements || 0}</Text>
                <Text style={styles.statLabel}>Earned</Text>
              </View>
              <View style={styles.statDivider} />
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{stats?.earned_points || 0}</Text>
                <Text style={styles.statLabel}>Points</Text>
              </View>
              <View style={styles.statDivider} />
              <View style={styles.statItem}>
                <Text style={styles.statValue}>
                  {stats?.completion_rate ? Math.round(stats.completion_rate) : 0}%
                </Text>
                <Text style={styles.statLabel}>Complete</Text>
              </View>
            </View>
          </MotiView>

          {/* Category Filters */}
          <ScrollView
            horizontal
            showsHorizontalScrollIndicator={false}
            contentContainerStyle={styles.filterContainer}
          >
            {(['all', 'streak', 'volume', 'pr', 'milestone', 'consistency'] as CategoryFilter[]).map(
              (category) => (
                <TouchableOpacity
                  key={category}
                  style={[
                    styles.filterButton,
                    selectedCategory === category && styles.filterButtonActive,
                  ]}
                  onPress={() => setSelectedCategory(category)}
                  activeOpacity={0.7}
                >
                  <Text
                    style={[
                      styles.filterText,
                      selectedCategory === category && styles.filterTextActive,
                    ]}
                  >
                    {category.charAt(0).toUpperCase() + category.slice(1)}
                  </Text>
                </TouchableOpacity>
              )
            )}
          </ScrollView>

          {/* Achievement Grid */}
          <View style={styles.gridContainer}>
            <Text style={styles.sectionTitle}>
              {selectedCategory === 'all' ? 'All Badges' : `${selectedCategory.charAt(0).toUpperCase() + selectedCategory.slice(1)} Badges`}
            </Text>
            <Text style={styles.sectionSubtitle}>
              {earnedCount} of {totalCount} earned
            </Text>
            
            <View style={styles.grid}>
              {filteredAchievements.map((achievement, index) => (
                <AchievementCard
                  key={achievement.achievement_id}
                  achievement={achievement}
                  index={index}
                />
              ))}
            </View>
          </View>
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

interface AchievementCardProps {
  achievement: Achievement;
  index: number;
}

function AchievementCard({ achievement, index }: AchievementCardProps) {
  const colors = CATEGORY_COLORS[achievement.category] || ['#999', '#666'];
  const icon = CATEGORY_ICONS[achievement.category] || 'trophy';
  const isEarned = achievement.is_completed;

  return (
    <MotiView
      from={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{
        type: 'timing',
        duration: 400,
        delay: index * 50,
      }}
      style={styles.achievementCard}
    >
      <View
        style={[
          styles.cardGradient,
          { backgroundColor: isEarned ? colors[0] : '#333' }
        ]}
      >
        <View style={[styles.badgeContainer, !isEarned && styles.badgeContainerLocked]}>
          <Ionicons
            name={icon}
            size={32}
            color={isEarned ? '#FFF' : '#666'}
          />
        </View>

        <Text style={[styles.achievementName, !isEarned && styles.achievementNameLocked]}>
          {achievement.name}
        </Text>
        <Text style={[styles.achievementDescription, !isEarned && styles.achievementDescriptionLocked]} numberOfLines={2}>
          {achievement.description}
        </Text>

        <View style={styles.progressContainer}>
          <View style={styles.progressBarBackground}>
            <View
              style={[
                styles.progressBarFill,
                {
                  width: `${achievement.completion_rate}%`,
                  backgroundColor: isEarned ? '#FFF' : COLORS.primary,
                },
              ]}
            />
          </View>
          <Text style={[styles.progressText, !isEarned && styles.progressTextLocked]}>
            {achievement.progress} / {achievement.requirement_value}
          </Text>
        </View>

        <View style={styles.pointsBadge}>
          <Text style={styles.pointsText}>{achievement.points} pts</Text>
        </View>

        {isEarned && achievement.earned_at && (
          <View style={styles.earnedContainer}>
            <Ionicons name="checkmark-circle" size={14} color="#FFF" />
            <Text style={styles.earnedText}>
              {new Date(achievement.earned_at).toLocaleDateString()}
            </Text>
          </View>
        )}

        {!isEarned && (
          <View style={styles.lockedOverlay}>
            <Ionicons name="lock-closed" size={24} color="#666" />
          </View>
        )}
      </View>
    </MotiView>
  );
}

const styles = {
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  safeArea: {
    flex: 1,
  },
  scrollContent: {
    padding: SPACING.base,
    paddingBottom: SPACING.xl * 2,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center' as const,
    alignItems: 'center' as const,
  },
  loadingText: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.base,
    marginTop: SPACING.base,
  },
  statsCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING.lg,
    ...SHADOWS.base,
  },
  statsTitle: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    marginBottom: SPACING.base,
    textAlign: 'center' as const,
  },
  statsRow: {
    flexDirection: 'row' as const,
    justifyContent: 'space-around' as const,
    alignItems: 'center' as const,
  },
  statItem: {
    alignItems: 'center' as const,
  },
  statValue: {
    color: COLORS.primary,
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
  },
  statLabel: {
    color: COLORS.text.tertiary,
fontSize: FONT_SIZES.sm,
    marginTop: SPACING.xs,
  },
  statDivider: {
    width: 1,
    height: 40,
    backgroundColor: COLORS.border.light,
  },
  filterContainer: {
    paddingVertical: SPACING.base,
    gap: SPACING.sm,
  },
  filterButton: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    marginRight: SPACING.sm,
  },
  filterButtonActive: {
    backgroundColor: COLORS.primary,
  },
  filterText: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
  },
  filterTextActive: {
    color: COLORS.white,
  },
  gridContainer: {
    marginTop: SPACING.base,
  },
  sectionTitle: {
    color: COLORS.text.inverse,
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    marginBottom: SPACING.xs,
  },
  sectionSubtitle: {
    color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.sm,
    marginBottom: SPACING.base,
  },
  grid: {
    flexDirection: 'row' as const,
    flexWrap: 'wrap' as const,
    gap: SPACING.base,
  },
  achievementCard: {
    width: CARD_WIDTH,
    marginBottom: SPACING.sm,
  },
  cardGradient: {
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    minHeight: 220,
    backgroundColor: COLORS.background.card,
    ...SHADOWS.base,
  },
  badgeContainer: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    justifyContent: 'center' as const,
    alignItems: 'center' as const,
    marginBottom: SPACING.base,
  },
  badgeContainerLocked: {
    backgroundColor: 'rgba(0, 0, 0, 0.3)',
  },
  achievementName: {
    color: '#FFF',
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    marginBottom: SPACING.xs,
  },
  achievementNameLocked: {
    color: '#888',
  },
  achievementDescription: {
    color: 'rgba(255, 255, 255, 0.9)',
    fontSize: FONT_SIZES.sm,
    marginBottom: SPACING.base,
    lineHeight: 18,
  },
  achievementDescriptionLocked: {
    color: '#666',
  },
  progressContainer: {
    marginTop: 'auto' as const,
  },
  progressBarBackground: {
    height: 6,
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    borderRadius: BORDER_RADIUS.sm,
    overflow: 'hidden' as const,
    marginBottom: SPACING.xs,
  },
  progressBarFill: {
    height: '100%' as const,
    borderRadius: BORDER_RADIUS.sm,
  },
  progressText: {
    color: '#FFF',
    fontSize: FONT_SIZES.xs,
    textAlign: 'right' as const,
  },
  progressTextLocked: {
    color: '#888',
  },
  pointsBadge: {
    position: 'absolute' as const,
    top: SPACING.sm,
    right: SPACING.sm,
    backgroundColor: 'rgba(0, 0, 0, 0.3)',
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.sm,
  },
  pointsText: {
    color: '#FFF',
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.bold,
  },
  earnedContainer: {
    position: 'absolute' as const,
    bottom: SPACING.sm,
    left: SPACING.sm,
    flexDirection: 'row' as const,
    alignItems: 'center' as const,
    gap: SPACING.xs,
  },
  earnedText: {
    color: '#FFF',
    fontSize: FONT_SIZES.xs,
  },
  lockedOverlay: {
    position: 'absolute' as const,
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    borderRadius: BORDER_RADIUS.lg,
    justifyContent: 'center' as const,
    alignItems: 'center' as const,
  },
   header: {
    marginBottom: 24,
  },
  titleContainer: {
    flexDirection: 'row' as const,
    alignItems: 'center' as const,
    gap: 12,
    marginBottom: 4,
  },
  title: {
    fontSize: 32,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  subtitle: {
    fontSize: 16,
    color: '#888888',
  },
  headerInfo: {
    flex: 1,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: '700' as const,
    color: '#FFFFFF',
  },
  headerSubtitle: {
    fontSize: 14,
    color: '#888888',
    marginTop: 2,
  },
};
