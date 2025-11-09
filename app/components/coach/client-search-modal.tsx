import React, { useState, useEffect } from 'react';
import {
  Modal,
  View,
  Text,
  TextInput,
  TouchableOpacity,
  FlatList,
  StyleSheet,
  ActivityIndicator,
  Platform,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useSearchUsers, useAssignClient } from '@/hooks/schema/use-coach';
import type { UserSearchResult } from '@/types/schema';
import {
  COLORS,
  SPACING,
  FONT_SIZES,
  FONT_WEIGHTS,
  BORDER_RADIUS,
} from '@/constants/theme';

interface ClientSearchModalProps {
  visible: boolean;
  onClose: () => void;
  onAssigned?: () => void;
}

export const ClientSearchModal: React.FC<ClientSearchModalProps> = ({
  visible,
  onClose,
  onAssigned,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');
  const [selectedUser, setSelectedUser] = useState<UserSearchResult | null>(null);
  const [notes, setNotes] = useState('');

  const { data, isLoading } = useSearchUsers(debouncedQuery);
  const assignClient = useAssignClient();

  // Debounce search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(searchQuery.trim());
    }, 300);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const handleClose = () => {
    setSearchQuery('');
    setDebouncedQuery('');
    setSelectedUser(null);
    setNotes('');
    assignClient.reset();
    onClose();
  };

  const handleAssign = () => {
    if (!selectedUser) return;

    assignClient.mutate(
      { username: selectedUser.username, notes: notes.trim() || undefined },
      {
        onSuccess: () => {
          Alert.alert('Success', `${selectedUser.first_name} ${selectedUser.last_name} has been assigned as your client.`);
          handleClose();
          if (onAssigned) onAssigned();
        },
        onError: (error) => {
          Alert.alert('Error', error.message || 'Failed to assign client');
        },
      }
    );
  };

  const renderUserCard = ({ item }: { item: UserSearchResult }) => {
    const fullName = `${item.first_name} ${item.last_name}`;
    const initials = fullName.split(' ').map(n => n[0]).join('').toUpperCase();

    return (
      <TouchableOpacity
        style={[
          styles.userCard,
          selectedUser?.auth_user_id === item.auth_user_id && styles.userCardSelected,
        ]}
        onPress={() => setSelectedUser(item)}
        activeOpacity={0.7}
      >
        <View style={styles.avatar}>
            
            <Text style={styles.avatarText}>{initials}</Text>
        </View>
        <View style={styles.userInfo}>
          <Text style={styles.userName}>{fullName}</Text>
          <Text style={styles.userEmail}>@{item.username}</Text>
          <View style={styles.badgeRow}>
            {item.has_coach && (
              <View style={[styles.badge, styles.badgeWarning]}>
                <Text style={[styles.badgeText, styles.badgeTextWarning]}>Has Coach</Text>
              </View>
            )}
          </View>
        </View>
        <Ionicons
          name={selectedUser?.auth_user_id === item.auth_user_id ? 'checkmark-circle' : 'radio-button-off'}
          size={24}
          color={selectedUser?.auth_user_id === item.auth_user_id ? COLORS.primary : COLORS.text.tertiary}
        />
      </TouchableOpacity>
    );
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={handleClose}
    >
      <View style={styles.modalBackdrop}>
        <View style={styles.modalContainer}>
          <View style={styles.header}>
            <Text style={styles.title}>Assign Client</Text>
            <TouchableOpacity onPress={handleClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.inverse} />
            </TouchableOpacity>
          </View>

          <View style={styles.searchContainer}>
            <Ionicons name="search" size={20} color={COLORS.text.tertiary} />
            <TextInput
              style={styles.searchInput}
              placeholder="Search by name, username, or email..."
              placeholderTextColor={COLORS.text.placeholder}
              value={searchQuery}
              onChangeText={setSearchQuery}
              autoFocus
            />
            {isLoading && <ActivityIndicator size="small" color={COLORS.primary} />}
          </View>

          {debouncedQuery.length < 2 ? (
            <View style={styles.emptyState}>
              <Ionicons name="search-outline" size={48} color={COLORS.text.tertiary} />
              <Text style={styles.emptyText}>Start typing to search for users</Text>
            </View>
          ) : (
            <FlatList
              data={data?.users || []}
              renderItem={renderUserCard}
              keyExtractor={(item) => item.auth_user_id}
              style={styles.list}
              ListEmptyComponent={
                !isLoading ? (
                  <View style={styles.emptyState}>
                    <Ionicons name="people-outline" size={48} color={COLORS.text.tertiary} />
                    <Text style={styles.emptyText}>No users found</Text>
                  </View>
                ) : null
              }
            />
          )}

          {selectedUser && (
            <View style={styles.notesSection}>
              <Text style={styles.notesLabel}>Notes (optional)</Text>
              <TextInput
                style={styles.notesInput}
                placeholder="Add any notes about this client..."
                placeholderTextColor={COLORS.text.placeholder}
                value={notes}
                onChangeText={setNotes}
                multiline
                numberOfLines={3}
              />
            </View>
          )}

          <View style={styles.footer}>
            <TouchableOpacity
              style={[styles.button, styles.buttonSecondary]}
              onPress={handleClose}
              disabled={assignClient.isPending}
            >
              <Text style={styles.buttonTextSecondary}>Cancel</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={[
                styles.button,
                styles.buttonPrimary,
                (!selectedUser || assignClient.isPending) && styles.buttonDisabled,
              ]}
              onPress={handleAssign}
              disabled={!selectedUser || assignClient.isPending}
            >
              {assignClient.isPending ? (
                <ActivityIndicator size="small" color={COLORS.text.primary} />
              ) : (
                <Text style={styles.buttonTextPrimary}>Assign Client</Text>
              )}
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  modalBackdrop: {
    flex: 1,    
    backgroundColor: COLORS.surface.overlay,
    justifyContent: 'flex-start',
    alignItems: 'center',
    padding: SPACING.md,
    paddingTop: Platform.OS === 'android' ? SPACING['6xl'] : SPACING['5xl'],
  },
  modalContainer: {
    width: '100%',
    maxWidth: 600,
    maxHeight: '90%',
    backgroundColor: COLORS.background.auth,
    borderRadius: BORDER_RADIUS['3xl'],
    borderWidth: 1,
    borderColor: 'rgba(255,255,255,0.08)',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 12 },
    shadowOpacity: 0.3,
    shadowRadius: 18,
    elevation: 18,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: SPACING.lg,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.dark,
  },
  title: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.inverse,
  },
  closeButton: {
    padding: SPACING.xs,
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    margin: SPACING.lg,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    gap: SPACING.sm,
  },
  searchInput: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
    paddingVertical: SPACING.xs,
  },
  list: {
    maxHeight: 400,
    paddingHorizontal: SPACING.lg,
  },
  userCard: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: SPACING.md,
    marginBottom: SPACING.sm,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    gap: SPACING.md,
  },
  userCardSelected: {
    borderColor: COLORS.primary,
    borderWidth: 2,
  },
  avatar: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
  },
  avatarText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
  },
  userInfo: {
    flex: 1,
    gap: 4,
  },
  userName: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
  },
  userEmail: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  badgeRow: {
    flexDirection: 'row',
    gap: SPACING.xs,
    marginTop: 4,
  },
  badge: {
    paddingHorizontal: SPACING.sm,
    paddingVertical: 4,
    borderRadius: BORDER_RADIUS.sm,
    backgroundColor: COLORS.background.accent,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  badgeText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    textTransform: 'capitalize',
  },
  badgeWarning: {
    backgroundColor: 'rgba(255,152,0,0.1)',
    borderColor: 'rgba(255,152,0,0.3)',
  },
  badgeTextWarning: {
    color: COLORS.warning,
  },
  emptyState: {
    padding: SPACING['4xl'],
    alignItems: 'center',
    justifyContent: 'center',
  },
  emptyText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    marginTop: SPACING.md,
    textAlign: 'center',
  },
  notesSection: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.dark,
  },
  notesLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.auth.secondary,
    marginBottom: SPACING.xs,
  },
  notesInput: {
    minHeight: 80,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    padding: SPACING.md,
    color: COLORS.text.inverse,
    backgroundColor: COLORS.background.card,
    textAlignVertical: 'top',
    fontSize: FONT_SIZES.base,
  },
  footer: {
    flexDirection: 'row',
    padding: SPACING.lg,
    gap: SPACING.sm,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.dark,
  },
  button: {
    flex: 1,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttonPrimary: {
    backgroundColor: COLORS.primary,
  },
  buttonSecondary: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  buttonDisabled: {
    opacity: 0.5,
  },
  buttonTextPrimary: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.primary,
  },
  buttonTextSecondary: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.tertiary,
  },
});
