import React, { useState } from 'react';
import {
  Modal,
  Platform,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
  Image,
} from 'react-native';
import Animated, {
  useAnimatedStyle,
  useSharedValue,
  withSpring,
} from 'react-native-reanimated';
import { useRouter } from 'expo-router';

import { Button, InputField } from '@/components/forms';
import { useToastMethods } from '@/components/ui';
import { COLORS } from '@/constants/theme';
import { useAuth } from '@/context/auth-context';
import {
  useCreateConversation,
  useUserLookup,
} from '@/hooks/message/use-conversation';
import type { PublicUserResponse } from '@/types/auth';
import {
  canMessage,
  formatRole,
  getInitial,
  roleRestrictionMessage,
} from '@/utils/conversation';

interface CreateConversationFABProps {
  onConversationCreated?: (conversationId: number) => void;
}

const FAB_SPRING_CONFIG = { damping: 14, stiffness: 260 } as const;

export const CreateConversationFAB: React.FC<CreateConversationFABProps> = ({
  onConversationCreated,
}) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedUser, setSelectedUser] =
    useState<PublicUserResponse | null>(null);

  const { user } = useAuth();
  const router = useRouter();
  const { showError, showInfo } = useToastMethods();

  const createConversation = useCreateConversation();
  const searchUser = useUserLookup();

  const scale = useSharedValue(1);

  const animatedFabStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
  }));

  const resetState = () => {
    setSearchTerm('');
    setSelectedUser(null);
    searchUser.reset();
  };

  const handleOpen = () => {
    scale.value = withSpring(0.92, FAB_SPRING_CONFIG, () => {
      scale.value = withSpring(1, FAB_SPRING_CONFIG);
    });
    setIsModalVisible(true);
  };

  const handleClose = () => {
    setIsModalVisible(false);
    resetState();
  };

  const handleSearch = async () => {
    if (!user?.id) {
      showError('You must be logged in to start a conversation.');
      return;
    }

    const username = searchTerm.trim();
    if (!username) {
      showInfo('Enter a username to search.');
      return;
    }

    try {
      setSelectedUser(null);
      searchUser.reset();
      const result = await searchUser.mutateAsync(username);

      if (result.id === user.id) {
        showInfo('You cannot start a conversation with yourself.');
        return;
      }

      if (!canMessage(user.role, result.role)) {
        showInfo(roleRestrictionMessage(user.role));
        return;
      }

      setSelectedUser(result);
    } catch (error) {
      console.error('Failed to lookup user:', error);
      showError('User not found. Double-check the username and try again.');
    }
  };

  const resolveParticipants = () => {
    if (!user || !selectedUser) {
      return null;
    }

    if (user.role === 'coach' && selectedUser.role === 'user') {
      return { coachId: user.id, clientId: selectedUser.id };
    }

    if (user.role === 'user' && selectedUser.role === 'coach') {
      return { coachId: selectedUser.id, clientId: user.id };
    }

    return null;
  };

  const handleCreateConversation = async () => {
    if (!user?.id) {
      showError('You must be logged in to start a conversation.');
      return;
    }

    if (!selectedUser) {
      showInfo('Search for a user and select them first.');
      return;
    }

    const participants = resolveParticipants();
    if (!participants) {
      showInfo(roleRestrictionMessage(user.role));
      return;
    }

    try {
      const result = await createConversation.mutateAsync({
        coach_id: participants.coachId,
        client_id: participants.clientId,
      });

      handleClose();

      if (onConversationCreated) {
        onConversationCreated(result.conversation.conversation_id);
      }
    } catch (error) {
      console.error('Failed to create conversation:', error);
      showError('Failed to create conversation. Please try again.');
    }
  };

  const searchHint = (() => {
    if (user?.role === 'coach') {
      return 'Search for a client by username to start chatting.';
    }
    if (user?.role === 'user') {
      return 'Search for a coach by username to start chatting.';
    }
    return 'Search for an account by username.';
  })();

  const disableActions = createConversation.isPending || searchUser.isPending;

  return (
    <>
      <Animated.View style={[styles.fabWrapper, animatedFabStyle]}>
        <TouchableOpacity
          accessibilityRole="button"
          accessibilityLabel="Start a new conversation"
          activeOpacity={0.85}
          onPress={handleOpen}
          style={styles.fab}
        >
          <Text style={styles.fabIcon}>+</Text>
        </TouchableOpacity>
      </Animated.View>

      <Modal
        visible={isModalVisible}
        animationType="fade"
        transparent
        onRequestClose={handleClose}
      >
        <View style={styles.modalOverlay}>
          <View style={styles.modalCard}>
            <Text style={styles.modalTitle}>New Conversation</Text>
            <Text style={styles.modalSubtitle}>{searchHint}</Text>

            <View style={styles.searchRow}>
              <View style={styles.searchInputWrapper}>
                <InputField
                  placeholder="Search by username"
                  value={searchTerm}
                  onChangeText={(value) => {
                    setSearchTerm(value);
                    setSelectedUser(null);
                  }}
                  autoCapitalize="none"
                  disabled={disableActions}
                />
              </View>
              <Button
                title="Search"
                variant="secondary"
                onPress={handleSearch}
                disabled={disableActions}
                loading={searchUser.isPending}
                style={styles.searchButton}
              />
            </View>

            {selectedUser && (
              <View style={styles.profileCard}>
                <View style={styles.profileHeader}>
                  {selectedUser.image ? (
                    <Image
                      source={{ uri: selectedUser.image }}
                      style={styles.avatarImage}
                    />
                  ) : (
                    <View style={styles.avatarFallback}>
                      <Text style={styles.avatarInitial}>
                        {getInitial(selectedUser.name || selectedUser.username)}
                      </Text>
                    </View>
                  )}

                  <View style={styles.profileMeta}>
                    <Text style={styles.profileName}>
                      {selectedUser.name || selectedUser.username}
                    </Text>
                    <Text style={styles.profileUsername}>
                      @{selectedUser.username} · {formatRole(selectedUser.role)}
                    </Text>
                  </View>
                </View>

                {selectedUser.bio ? (
                  <Text style={styles.profileBio}>{selectedUser.bio}</Text>
                ) : null}

                <Button
                  title="View Profile"
                  variant="outline"
                  onPress={() => {
                    handleClose();
                    router.push(`/profile/${selectedUser.username}` as never);
                  }}
                  style={styles.profileButton}
                />
              </View>
            )}

            {!selectedUser && searchUser.isSuccess && !searchUser.isPending ? (
              <Text style={styles.helperText}>
                No eligible account found. Try another username.
              </Text>
            ) : null}

            <View style={styles.actionsRow}>
              <Button
                title="Cancel"
                variant="secondary"
                onPress={handleClose}
                disabled={createConversation.isPending}
                style={styles.actionButton}
              />
              <Button
                title="Start Conversation"
                variant="primary"
                onPress={handleCreateConversation}
                disabled={!selectedUser || createConversation.isPending}
                loading={createConversation.isPending}
                style={styles.actionButton}
              />
            </View>
          </View>
        </View>
      </Modal>
    </>
  );
};

const styles = StyleSheet.create({
  fabWrapper: {
    position: 'absolute',
    right: 20,
    top: Platform.select({ ios: 28, android: 20, default: 20 }),
    zIndex: 10,
  },
  fab: {
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    elevation: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.3,
    shadowRadius: 10,
  },
  fabIcon: {
    color: COLORS.text.inverse,
    fontSize: 40,
    fontWeight: '700',
    marginTop: -4,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: COLORS.surface.overlay,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  modalCard: {
    width: '92%',
    maxWidth: 460,
    padding: 24,
    borderRadius: 28,
    backgroundColor: COLORS.background.surface,
    borderWidth: 1,
    borderColor: 'rgba(255,255,255,0.08)',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 12 },
    shadowOpacity: 0.3,
    shadowRadius: 18,
    elevation: 18,
  },
  modalTitle: {
    fontSize: 24,
    fontWeight: '700',
    color: COLORS.text.primary,
    textAlign: 'center',
    marginBottom: 12,
  },
  modalSubtitle: {
    fontSize: 14,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginBottom: 20,
  },
  searchRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  searchInputWrapper: {
    flex: 1,
    marginRight: 12,
  },
  searchButton: {
    minWidth: 110,
  },
  profileCard: {
    borderRadius: 24,
    padding: 16,
    backgroundColor: 'rgba(0,0,0,0.04)',
    borderWidth: 1,
    borderColor: 'rgba(255,255,255,0.08)',
    marginBottom: 12,
  },
  profileHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 10,
  },
  avatarImage: {
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: 'rgba(255,255,255,0.08)',
  },
  avatarFallback: {
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: 'rgba(255,255,255,0.08)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  avatarInitial: {
    fontSize: 22,
    fontWeight: '700',
    color: COLORS.text.primary,
  },
  profileMeta: {
    flex: 1,
    marginLeft: 12,
  },
  profileName: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text.primary,
  },
  profileUsername: {
    fontSize: 14,
    color: COLORS.text.tertiary,
    marginTop: 4,
  },
  profileBio: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  profileButton: {
    marginTop: 8,
  },
  helperText: {
    fontSize: 14,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginBottom: 12,
  },
  actionsRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginTop: 12,
  },
  actionButton: {
    flex: 1,
    marginHorizontal: 6,
  },
});

