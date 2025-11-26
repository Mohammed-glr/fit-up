import React from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useCoachDashboard, useCoachClients } from '@/hooks/schema/use-coach';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

export default function SchemaScreen() {
  const router = useRouter();
  const { data: dashboard, isLoading: loadingDashboard, refetch: refetchDashboard } = useCoachDashboard();
  const { data: clientsData, isLoading: loadingClients, refetch: refetchClients } = useCoachClients();

  const [refreshing, setRefreshing] = React.useState(false);

  const onRefresh = async () => {
    setRefreshing(true);
    await Promise.all([refetchDashboard(), refetchClients()]);
    setRefreshing(false);
  };

  if (loadingDashboard || loadingClients) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.contentContainer}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor={COLORS.primary} />
      }
    >
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Workout Schemas</Text>
        <Text style={styles.headerSubtitle}>Manage your clients' training programs</Text>
      </View>

      <View style={styles.statsContainer}>
        <View style={[styles.statCard, styles.statCardPrimary]}>
          <Ionicons name="people" size={32} color={COLORS.primary} />
          <Text style={styles.statNumber}>{dashboard?.total_clients || 0}</Text>
          <Text style={styles.statLabel}>Total Clients</Text>
        </View>

        <View style={styles.statCard}>
          <Ionicons name="calendar" size={32} color={COLORS.success} />
          <Text style={styles.statNumber}>{dashboard?.active_schemas || 0}</Text>
          <Text style={styles.statLabel}>Active Schemas</Text>
        </View>
      </View>

      <View style={styles.statsContainer}>
        <View style={styles.statCard}>
          <Ionicons name="fitness" size={32} color={COLORS.success} />
          <Text style={styles.statNumber}>{dashboard?.total_workouts || 0}</Text>
          <Text style={styles.statLabel}>Workouts This Month</Text>
        </View>

        <View style={styles.statCard}>
          <Ionicons name="trending-up" size={32} color={COLORS.info} />
          <Text style={styles.statNumber}>
            {dashboard?.average_completion ? `${Math.round(dashboard.average_completion * 100)}%` : '0%'}
          </Text>
          <Text style={styles.statLabel}>Avg Completion</Text>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Quick Actions</Text>
        <View style={styles.actionsGrid}>
          <TouchableOpacity
            style={styles.actionCard}
            onPress={() => {
              router.push('/(coach)/schema-create');
            }}
          >
            <View style={[styles.actionIcon, { backgroundColor: `${COLORS.primary}20` }]}>
              <Ionicons name="add-circle" size={32} color={COLORS.primary} />
            </View>
            <Text style={styles.actionTitle}>Create Schema</Text>
            <Text style={styles.actionSubtitle}>Build new program</Text>
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.actionCard}
            onPress={() => router.push('/(coach)/clients')}
          >
            <View style={[styles.actionIcon, { backgroundColor: `${COLORS.success}20` }]}>
              <Ionicons name="people" size={32} color={COLORS.success} />
            </View>
            <Text style={styles.actionTitle}>My Clients</Text>
            <Text style={styles.actionSubtitle}>Manage clients</Text>
          </TouchableOpacity>

        </View>
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Recent Clients</Text>
          <TouchableOpacity onPress={() => router.push('/(coach)/clients')}>
            <Text style={styles.viewAllText}>View All</Text>
          </TouchableOpacity>
        </View>

        {(clientsData?.clients ?? []) .slice(0, 5).map((client) => (
          <TouchableOpacity
            key={client.user_id}
            style={styles.clientCard}
            onPress={() => router.push({ pathname: '/(coach)/client-details', params: { userId: client.user_id.toString() } })}
          >
            <View style={styles.clientAvatar}>
              <Text style={styles.clientAvatarText}>
                {client.first_name[0]}{client.last_name[0]}
              </Text>
            </View>
            <View style={styles.clientInfo}>
              <Text style={styles.clientName}>
                {client.first_name} {client.last_name}
              </Text>
              <View style={styles.clientStats}>
                <View style={styles.clientStat}>
                  <Ionicons name="barbell" size={14} color={COLORS.text.tertiary} />
                  <Text style={styles.clientStatText}>{client.total_workouts} workouts</Text>
                </View>
                <View style={styles.clientStat}>
                  <Ionicons name="flame" size={14} color={COLORS.text.tertiary} />
                  <Text style={styles.clientStatText}>{client.current_streak} day streak</Text>
                </View>
              </View>
            </View>
            <View style={styles.clientCompletion}>
              <Text style={styles.completionText}>{Math.round(client.completion_rate)}%</Text>
              <Text style={styles.completionLabel}>Complete</Text>
            </View>
          </TouchableOpacity>
        ))}

        {(!clientsData?.clients || clientsData.clients.length === 0) && (
          <View style={styles.emptyState}>
            <Ionicons name="people-outline" size={64} color={COLORS.text.tertiary} />
            <Text style={styles.emptyTitle}>No clients yet</Text>
            <Text style={styles.emptySubtitle}>
              Start by assigning clients to create their workout programs
            </Text>
          </View>
        )}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  contentContainer: {
    paddingBottom: SPACING['6xl'],
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.background.auth,
  },
  header: {
    paddingHorizontal: SPACING.lg,
    paddingTop: SPACING.xl,
    paddingBottom: SPACING.lg,
  },
  headerTitle: {
    fontSize: FONT_SIZES['3xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.xs,
  },
  headerSubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
  },
  statsContainer: {
    flexDirection: 'row',
    paddingHorizontal: SPACING.lg,
    marginBottom: SPACING.md,
    gap: SPACING.md,
  },
  statCard: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    padding: SPACING.base,
    borderRadius: BORDER_RADIUS['2xl'],
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  statCardPrimary: {
    borderWidth: 2,
    borderColor: COLORS.primary,
  },
  statNumber: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.xs,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginTop: SPACING.xs,
  },
  section: {
    paddingHorizontal: SPACING.lg,
    marginTop: SPACING.xl,
    marginBottom: SPACING.xl,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  sectionTitle: {
    marginBottom: SPACING.md,
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  viewAllText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
  },
  actionsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  actionCard: {
    width: '48%',
    backgroundColor: COLORS.background.card,
    padding: SPACING.base,
    borderRadius: BORDER_RADIUS['2xl'],
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  actionIcon: {
    width: 64,
    height: 64,
    borderRadius: BORDER_RADIUS.full,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  actionTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.xs,
  },
  actionSubtitle: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
  clientCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    padding: SPACING.base,
    borderRadius: BORDER_RADIUS['2xl'],
    marginBottom: SPACING.md,
    ...SHADOWS.sm,
  },
  clientAvatar: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  clientAvatarText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  clientInfo: {
    flex: 1,
  },
  clientName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: SPACING.xs,
  },
  clientStats: {
    flexDirection: 'row',
    gap: SPACING.md,
  },
  clientStat: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  clientStatText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  clientCompletion: {
    alignItems: 'center',
  },
  completionText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  completionLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: SPACING['4xl'],
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.md,
    marginBottom: SPACING.xs,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    paddingHorizontal: SPACING.xl,
  },
});
