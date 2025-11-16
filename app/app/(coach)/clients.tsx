import React, { useState, useMemo } from 'react';
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  TextInput,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
} from 'react-native';
import { router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useCoachClients } from '@/hooks/schema/use-coach';
import type { ClientSummary } from '@/types/schema';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { ClientStatusBadge } from '@/components/coach/client-status-badge';
import { StatusFilterChips } from '@/components/coach/status-filter-chips';
import { ClientSortDropdown } from '@/components/coach/client-sort-dropdown';
import { QuickActionsMenu } from '@/components/coach/quick-actions-menu';
import { 
  filterClientsByStatus, 
  getStatusCounts, 
  getLastWorkoutDescription,
  type ClientStatus 
} from '@/utils/client-status';
import { sortClients, type SortOption, type SortOrder } from '@/utils/client-sorting';

export default function ClientListScreen() {
  const { data, isLoading, refetch, isRefetching } = useCoachClients();
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<ClientStatus | 'all'>('all');
  const [sortBy, setSortBy] = useState<SortOption>('last_active');
  const [sortOrder, setSortOrder] = useState<SortOrder>('desc');

  const allClients = data?.clients || [];

  const statusCounts = useMemo(() => getStatusCounts(allClients), [allClients]);

  const filteredClients = useMemo(() => {
    let result = allClients;

    result = filterClientsByStatus(result, statusFilter);

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      result = result.filter((client) => {
        return (
          client.first_name.toLowerCase().includes(query) ||
          client.last_name.toLowerCase().includes(query) ||
          client.email.toLowerCase().includes(query)
        );
      });
    }

    result = sortClients(result, sortBy, sortOrder);

    return result;
  }, [allClients, statusFilter, searchQuery, sortBy, sortOrder]);

  const handleSortChange = (newSortBy: SortOption, newSortOrder: SortOrder) => {
    setSortBy(newSortBy);
    setSortOrder(newSortOrder);
  };

  if (isLoading) {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
        <Text style={styles.loadingText}>Loading your clients...</Text>
      </View>
    );
  }

  const handleClientPress = (client: ClientSummary) => {
    router.push({
      pathname: '/(coach)/client-details',
      params: { userId: client.user_id.toString() },
    });
  };

  const handleCreateSchema = (client: ClientSummary) => {
    router.push({
      pathname: '/(coach)/schema-create',
      params: { userId: client.user_id.toString() },
    });
  };

  const renderClientItem = ({ item }: { item: ClientSummary }) => {
    const fullName = `${item.first_name} ${item.last_name}`;
    const completionPercentage = Math.round(item.completion_rate * 100);
    const lastWorkoutText = getLastWorkoutDescription(item.last_workout_date || null);

    return (
      <TouchableOpacity
        style={styles.clientCard}
        onPress={() => handleClientPress(item)}
        activeOpacity={0.7}
      >
        <View style={styles.clientHeader}>
          <View style={styles.avatar}>
            <Text style={styles.avatarText}>
              {item.first_name[0]}{item.last_name[0]}
            </Text>
          </View>
          <View style={styles.clientInfo}>
            <Text style={styles.clientName}>{fullName}</Text>
            <Text style={styles.clientEmail}>{item.email}</Text>
            <View style={styles.badgeContainer}>
              <ClientStatusBadge client={item} size="small" />
              <View style={styles.badge}>
                <Text style={styles.badgeText}>{item.fitness_level}</Text>
              </View>
            </View>
          </View>
          <QuickActionsMenu client={item} onRemove={() => refetch()} />
        </View>

        <View style={styles.statsRow}>
          <View style={styles.statItem}>
            <Ionicons name="barbell" size={16} color={COLORS.primary} />
            <Text style={styles.statText}>{item.total_workouts} workouts</Text>
          </View>
          <View style={styles.statItem}>
            <Ionicons name="flame" size={16} color={COLORS.warning} />
            <Text style={styles.statText}>{item.current_streak} day streak</Text>
          </View>
          <View style={styles.statItem}>
            <Ionicons name="checkmark-circle" size={16} color={COLORS.success} />
            <Text style={styles.statText}>{completionPercentage}% complete</Text>
          </View>
        </View>

        <Text style={styles.lastWorkout}>
          Last workout: {lastWorkoutText}
        </Text>

        <View style={styles.actions}>
          <TouchableOpacity
            style={[styles.actionButton, styles.primaryAction]}
            onPress={() => handleCreateSchema(item)}
          >
            <Ionicons name="add-circle" size={20} color={COLORS.text.primary} />
            <Text style={styles.primaryActionText}>Create Schema</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.actionButton, styles.secondaryAction]}
            onPress={() => handleClientPress(item)}
          >
            <Text style={styles.secondaryActionText}>View Details</Text>
          </TouchableOpacity>
        </View>
      </TouchableOpacity>
    );
  };

  
  return (
    <View style={styles.container}>
      
      <View style={styles.header}>
        <View style={styles.searchContainer}>
          <Ionicons
            name="search"
            size={20}
            color={COLORS.text.tertiary}
            style={styles.searchIcon}
          />
          <TextInput
            style={styles.searchInput}
            placeholder="Search clients..."
            placeholderTextColor={COLORS.text.placeholder}
            value={searchQuery}
            onChangeText={setSearchQuery}
          />
          {searchQuery.length > 0 && (
            <TouchableOpacity onPress={() => setSearchQuery('')}>
              <Ionicons name="close-circle" size={20} color={COLORS.text.tertiary} />
            </TouchableOpacity>
          )}
        </View>
      </View>

      <StatusFilterChips
        activeFilter={statusFilter}
        onFilterChange={setStatusFilter}
        counts={statusCounts}
      />

      <View style={styles.statsBar}>
        <Text style={styles.statsText}>
          {filteredClients.length} {filteredClients.length === 1 ? 'client' : 'clients'}
        </Text>
        <ClientSortDropdown
          sortBy={sortBy}
          sortOrder={sortOrder}
          onSortChange={handleSortChange}
        />
      </View>

      <FlatList
        data={filteredClients}
        renderItem={renderClientItem}
        keyExtractor={(item) => item.user_id.toString()}
        contentContainerStyle={styles.listContent}
        refreshControl={
          <RefreshControl
            refreshing={isRefetching}
            onRefresh={refetch}
            tintColor={COLORS.primary}
            colors={[COLORS.primary]}
          />
        }
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Ionicons name="people-outline" size={64} color={COLORS.text.tertiary} />
            <Text style={styles.emptyTitle}>No clients found</Text>
            <Text style={styles.emptySubtitle}>
              {searchQuery ? 'Try a different search term' : 'Start by assigning clients'}
              
            </Text>
          </View>
        }
      />

    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    
    backgroundColor: COLORS.background.auth,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.background.auth,
  },
  loadingText: {
    marginTop: SPACING.sm,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
  },
  header: {
    padding: SPACING.sm,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
   searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    marginHorizontal: SPACING.lg,
    marginTop: SPACING.md,
    paddingHorizontal: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  searchIcon: {
    marginRight: SPACING.sm,
  },
  searchInput: {
    flex: 1,
    paddingVertical: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.primary,
  },
  statsBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: SPACING.lg,
    paddingBottom: SPACING.sm,
  },
  statsText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    flex: 1,
  },
  listContent: {
    padding: SPACING.base,
  },
  clientCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.base,
    marginBottom: SPACING.md,
    ...SHADOWS.sm,
  },
  clientHeader: {
    flexDirection: 'row',
    marginBottom: SPACING.md,
  },
  avatar: {
    width: 56,
    height: 56,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  avatarText: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  clientInfo: {
    flex: 1,
  },
  clientName: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: 4,
  },
  clientEmail: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.xs,
  },
  badgeContainer: {
    flexDirection: 'row',
    gap: SPACING.xs,
  },
  badge: {
    paddingHorizontal: SPACING.sm,
    paddingVertical: 4,
    borderRadius: BORDER_RADIUS.sm,
    backgroundColor: COLORS.background.primary,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  activeBadge: {
    backgroundColor: COLORS.success,
    borderColor: COLORS.success,
  },
  inactiveBadge: {
    backgroundColor: 'transparent',
    borderColor: COLORS.border.dark,
  },
  badgeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.tertiary,
    textTransform: 'capitalize',
  },
  activeBadgeText: {
    color: COLORS.text.primary,
  },
  inactiveBadgeText: {
    color: COLORS.text.tertiary,
  },
  statsRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: SPACING.sm,
  },
  statItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  statText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.auth.secondary,
  },
  lastWorkout: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginBottom: SPACING.md,
  },
  actions: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  actionButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    gap: SPACING.xs,
  },
  primaryAction: {
    backgroundColor: COLORS.primary,
  },
  secondaryAction: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
    borderColor: COLORS.border.dark,
  },
  primaryActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  secondaryActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.tertiary,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING['5xl'],
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.base,
    marginBottom: SPACING.xs,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
});
