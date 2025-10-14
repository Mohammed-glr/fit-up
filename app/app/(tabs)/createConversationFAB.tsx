// components/CreateConversationFAB.tsx
import { useCreateConversation } from '@/hooks/message/use-conversation';
import React, { useState } from 'react';
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  Modal,
  View,
  Platform,
} from 'react-native';
import Animated, { useAnimatedStyle, withSpring, useSharedValue } from 'react-native-reanimated';
import { useToastMethods } from '@/components/ui';
import { Button, InputField } from '@/components/forms';
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
      showError('Failed to create conversation. Please try again.');
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

            <InputField
              placeholder="Coach ID"
              value={coachId}
              onChangeText={setCoachId}
              autoCapitalize="none"
              disabled={createConversation.isPending}
            />

            <InputField
              placeholder="Client ID"
              value={clientId}
              onChangeText={setClientId}
              autoCapitalize="none"
              disabled={createConversation.isPending}
            />

            <View style={styles.buttonRow}>
              <Button
                title="Cancel"
                variant="secondary"
                onPress={() => {
                  setModalVisible(false);
                  setCoachId('');
                  setClientId('');
                }}
                disabled={createConversation.isPending}
                style={styles.modalButton}
              />

              <Button
                title="Create"
                variant="primary"
                onPress={handleCreate}
                disabled={createConversation.isPending}
                loading={createConversation.isPending}
                style={styles.modalButton}
              />
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
    right: 16,
    top: Platform.OS === 'ios' ? 20 : 20,
    zIndex: 800,
  },
  fab: {
width: 48,
        height: 48,
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
    color: COLORS.text.inverse,
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
    borderRadius: 60,
    padding: 20,
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
  buttonRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: 16,
    gap: 12,
  },
  modalButton: {
    flex: 1,
  },
});