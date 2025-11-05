import React, { useState } from 'react';
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  TextInput,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  Modal,
  Alert,
} from 'react-native';
import { router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useAssignClient, useCoachClients } from '@/hooks/schema/use-coach';
import type { ClientSummary } from '@/types/schema';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { InputField, Button } from '@/components/forms';

export default function ClientListScreen() {
  const { data, isLoading, refetch, isRefetching } = useCoachClients();
  const assignClient = useAssignClient();
  const [searchQuery, setSearchQuery] = useState('');
  const [isAssignModalVisible, setAssignModalVisible] = useState(false);
  const [userIdInput, setUserIdInput] = useState('');
  const [notes, setNotes] = useState('');
  const [assignError, setAssignError] = useState<string | null>(null);

  const filteredClients = data?.clients.filter((client) => {
    const query = searchQuery.toLowerCase();
    return (
      client.first_name.toLowerCase().includes(query) ||
      client.last_name.toLowerCase().includes(query) ||
      client.email.toLowerCase().includes(query)
    );
  }) || [];

  const handleOpenAssignModal = () => {
    setAssignError(null);
    setUserIdInput('');
    setNotes('');
    assignClient.reset();
    setAssignModalVisible(true);
  };

  const handleCloseAssignModal = () => {
    setAssignModalVisible(false);
    setAssignError(null);
    assignClient.reset();
  };

  const handleAssignSubmit = () => {
    const trimmed = userIdInput.trim();
    if (!trimmed) {
      setAssignError('Enter the workout profile ID for the client.');
      return;
    }

    const parsedUserId = Number(trimmed);
    if (!Number.isInteger(parsedUserId) || parsedUserId <= 0) {
      setAssignError('Profile ID must be a positive number.');
      return;
    }

    setAssignError(null);

    assignClient.mutate(
      { user_id: parsedUserId, notes: notes.trim() || undefined },
      {
        onSuccess: () => {
          setAssignModalVisible(false);
          setUserIdInput('');
          setNotes('');
          Alert.alert('Client assigned', 'The client has been assigned successfully.');
        },
        onError: (error) => {
          setAssignError(error?.message ?? 'Unable to assign client.');
        },
      }
    );
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
    const hasActiveSchema = !!item.current_schema_id;

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
              <View style={[styles.badge, hasActiveSchema ? styles.activeBadge : styles.inactiveBadge]}>
                <Text style={[styles.badgeText, hasActiveSchema ? styles.activeBadgeText : styles.inactiveBadgeText]}>
                  {hasActiveSchema ? 'Active Schema' : 'No Schema'}
                </Text>
              </View>
              <View style={styles.badge}>
                <Text style={styles.badgeText}>{item.fitness_level}</Text>
              </View>
            </View>
          </View>
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

        {item.last_workout_date && (
          <Text style={styles.lastWorkout}>
            Last workout: {new Date(item.last_workout_date).toLocaleDateString()}
          </Text>
        )}

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
        <View style={styles.headerRow}>
          <Text style={styles.headerTitle}>Clients</Text>
          <TouchableOpacity
            style={styles.assignButton}
            onPress={handleOpenAssignModal}
            activeOpacity={0.8}
          >
            <Ionicons name="person-add" size={18} color={COLORS.text.primary} />
            <Text style={styles.assignButtonText}>Assign Client</Text>
          </TouchableOpacity>
        </View>
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
        <View style={styles.statsBar}>
          <Text style={styles.statsText}>
            {filteredClients.length} {filteredClients.length === 1 ? 'client' : 'clients'}
          </Text>
        </View>
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

      <Modal
        visible={isAssignModalVisible}
        transparent
        animationType="fade"
        onRequestClose={handleCloseAssignModal}
      >
        <View style={styles.modalBackdrop}>
          <View style={styles.modalContainer}>
            <Text style={styles.modalTitle}>Assign Client</Text>

            <InputField
              label="Workout profile ID"
              placeholder="e.g. 42"
              value={userIdInput}
              onChangeText={(text) => {
                setAssignError(null);
                setUserIdInput(text.replace(/[^0-9]/g, ''));
              }}
              keyboardType="numeric"
              leftIcon="person"
            />

            <Text style={styles.notesLabel}>Notes (optional)</Text>
            <TextInput
              value={notes}
              onChangeText={(text) => setNotes(text)}
              style={styles.notesInput}
              placeholder="Share any context you want to remember."
              placeholderTextColor={COLORS.text.placeholder}
              multiline
            />

            {assignError && <Text style={styles.assignError}>{assignError}</Text>}

            <View style={styles.modalActions}>
              <Button
                title="Cancel"
                variant="secondary"
                onPress={handleCloseAssignModal}
                style={styles.modalActionButton}
              />
              <Button
                title="Assign"
                onPress={handleAssignSubmit}
                loading={assignClient.isPending}
                style={styles.modalActionButton}
              />
            </View>
          </View>
        </View>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
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
    gap: SPACING.md,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: SPACING.lg,
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
  assignButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    ...SHADOWS.sm,
  },
  assignButtonText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
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
    alignItems: 'center',
  },
  statsText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  listContent: {
    padding: SPACING.base,
  },
  clientCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
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
    borderRadius: BORDER_RADIUS.md,
    gap: SPACING.xs,
  },
  primaryAction: {
    backgroundColor: COLORS.primary,
  },
  secondaryAction: {
    backgroundColor: 'transparent',
    borderWidth: 1,
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
  modalBackdrop: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.6)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.lg,
  },
  modalContainer: {
    width: '100%',
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.lg,
    gap: SPACING.base,
  },
  modalTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
  },
  notesLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.secondary,
    marginBottom: SPACING.xs,
  },
  notesInput: {
    minHeight: 96,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    padding: SPACING.sm,
    color: COLORS.text.auth.primary,
    backgroundColor: COLORS.background.surface || COLORS.background.card,
    textAlignVertical: 'top',
  },
  assignError: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.error,
  },
  modalActions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: SPACING.sm,
  },
  modalActionButton: {
    flex: 1,
  },
});
