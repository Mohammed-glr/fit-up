// components/CreateConversationFAB.tsx
import { useCreateConversation } from '@/hooks/message/use-conversation';
import React, { useState } from 'react';
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  Modal,
  View,
  TextInput,
  ActivityIndicator,
  Platform,
} from 'react-native';
import Animated, { useAnimatedStyle, withSpring, useSharedValue } from 'react-native-reanimated';
import { useToastMethods } from '@/components/ui';
import { COLORS } from '@/constants/theme';
interface CreateConversationFABProps {
  onConversationCreated?: (conversationId: number) => void;
}

export const CreateConversationFAB: React.FC<CreateConversationFABProps> = ({
  onConversationCreated,
}) => {
  const [modalVisible, setModalVisible] = useState(false);
  const [coachId, setCoachId] = useState('');
  const [clientId, setClientId] = useState('');
  const { showError, showSuccess, showInfo } = useToastMethods();
  const scale = useSharedValue(1);

  const createConversation = useCreateConversation();

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
  }));

  const handleCreate = async () => {
    if (!coachId.trim() || !clientId.trim()) {
        showError('Both Coach ID and Client ID are required.');
      return;
    }

    try {
      const result = await createConversation.mutateAsync({
        coach_id: coachId.trim(),
        client_id: clientId.trim(),
      });

      setModalVisible(false);
      setCoachId('');
      setClientId('');

      showSuccess('Conversation created successfully!');

      if (onConversationCreated) {
        onConversationCreated(result.conversation.conversation_id);
      }
    } catch (error) {
      console.error('Failed to create conversation:', error);
    }
  };

  return (
    <>
      <Animated.View style={[styles.fabContainer, animatedStyle]}>
        <TouchableOpacity
          style={styles.fab}
          onPress={() => {
            scale.value = withSpring(0.9, { damping: 10 }, () => {
              scale.value = withSpring(1, { damping: 10 });
            });
            setModalVisible(true);
          }}
          activeOpacity={0.9}
        >
          <Text style={styles.fabIcon}>+</Text>
        </TouchableOpacity>
      </Animated.View>

      <Modal
        visible={modalVisible}
        animationType="fade"
        transparent={true}
        onRequestClose={() => setModalVisible(false)}
      >
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>New Conversation</Text>

            <TextInput
              style={styles.input}
              placeholder="Coach ID"
              value={coachId}
              onChangeText={setCoachId}
              autoCapitalize="none"
            />

            <TextInput
              style={styles.input}
              placeholder="Client ID"
              value={clientId}
              onChangeText={setClientId}
              autoCapitalize="none"
            />

            <View style={styles.buttonRow}>
              <TouchableOpacity
                style={[styles.modalButton, styles.cancelButton]}
                onPress={() => {
                  setModalVisible(false);
                  setCoachId('');
                  setClientId('');
                }}
                disabled={createConversation.isPending}
              >
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[styles.modalButton, styles.createButton]}
                onPress={handleCreate}
                disabled={createConversation.isPending}
              >
                {createConversation.isPending ? (
                  <ActivityIndicator color="#fff" />
                ) : (
                  <Text style={styles.createButtonText}>Create</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </Modal>
    </>
  );
};

const styles = StyleSheet.create({
  fabContainer: {
    position: 'absolute',
    right: 20,
    top: Platform.OS === 'ios' ? 20 : 20,
    zIndex: 800,
  },
  fab: {
    width: 64,
    height: 64,
    borderRadius: 32,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    elevation: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.4,
    shadowRadius: 16,
  },
  fabIcon: {
    color: COLORS.text.primary,
    fontSize: 44,
    fontWeight: '700',
    marginTop: -8,
    lineHeight: 32,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: COLORS.surface.overlay,
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalContent: {
    backgroundColor: COLORS.darkGray,
    borderWidth: 0.5,
    borderColor: 'rgba(255, 255, 255, 0.1)',
    borderRadius: 63,
    padding: 24,
    width: '88%',
    maxWidth: 420,
    elevation: 24,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 12 },
    shadowOpacity: 0.58,
    shadowRadius: 16,
  },
  modalTitle: {
    fontSize: 24,
    fontWeight: '800',
    marginBottom: 24,
    textAlign: 'center',
    color: COLORS.text.inverse,
    letterSpacing: 0.5,
  },
  input: {
    borderWidth: 2,
    borderColor: COLORS.border.subtle,
    backgroundColor: COLORS.background.card,
    borderRadius: 32,
    padding: 16,
    marginBottom: 16,
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.primary,
  },
  buttonRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: 16,
    gap: 12,
  },
  modalButton: {
    flex: 1,
    paddingVertical: 16,
    borderRadius: 32,
    alignItems: 'center',
    justifyContent: 'center',
    elevation: 4,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.2,
    shadowRadius: 4,
  },
  cancelButton: {
    backgroundColor: COLORS.background.secondary,
  },
  cancelButtonText: {
    color: COLORS.text.primary,
    fontSize: 16,
    fontWeight: '700',
    letterSpacing: 0.3,
  },
  createButton: {
    backgroundColor: COLORS.primary,
  },
  createButtonText: {
    color: COLORS.text.primary,
    fontSize: 16,
    fontWeight: '700',
    letterSpacing: 0.3,
  },
});