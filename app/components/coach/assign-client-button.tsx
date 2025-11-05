import React, { useState } from 'react';
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  Modal,
  View,
  TextInput,
  Alert,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import {
  COLORS,
  SPACING,
  FONT_SIZES,
  FONT_WEIGHTS,
  BORDER_RADIUS,
  SHADOWS,
} from '@/constants/theme';
import { InputField, Button } from '@/components/forms';
import { useAssignClient } from '@/hooks/schema/use-coach';

type AssignClientVariant = 'default' | 'icon';

interface AssignClientButtonProps {
  label?: string;
  variant?: AssignClientVariant;
  disabled?: boolean;
  onAssigned?: () => void;
}

export const AssignClientButton: React.FC<AssignClientButtonProps> = ({
  label = 'Assign Client',
  variant = 'default',
  disabled = false,
  onAssigned,
}) => {
  const assignClient = useAssignClient();
  const [isModalVisible, setModalVisible] = useState(false);
  const [usernameInput, setUsernameInput] = useState('');
  const [notes, setNotes] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleOpen = () => {
    if (disabled) {
      return;
    }
    assignClient.reset();
    setError(null);
  setUsernameInput('');
    setNotes('');
    setModalVisible(true);
  };

  const handleClose = () => {
    setModalVisible(false);
    setError(null);
    assignClient.reset();
  };

  const handleSubmit = () => {
    const trimmed = usernameInput.trim();
    if (!trimmed) {
      setError('Enter the username for the client.');
      return;
    }

    if (/\s/.test(trimmed)) {
      setError('Username cannot contain spaces.');
      return;
    }

    setError(null);

    assignClient.mutate(
      { username: trimmed, notes: notes.trim() || undefined },
      {
        onSuccess: () => {
          setModalVisible(false);
          setUsernameInput('');
          setNotes('');
          if (onAssigned) {
            onAssigned();
          }
          Alert.alert('Client assigned', 'The client has been assigned successfully.');
        },
        onError: (mutationError) => {
          setError(mutationError?.message ?? 'Unable to assign client.');
        },
      }
    );
  };

  const isPending = assignClient.isPending;
  const triggerDisabled = disabled || isPending;

  return (
    <View style={styles.container}>
      <TouchableOpacity
        style={[
          styles.button,
          variant === 'icon' && styles.iconButton,
          triggerDisabled && styles.disabled,
        ]}
        onPress={handleOpen}
        activeOpacity={0.8}
        disabled={triggerDisabled}
        accessibilityLabel={label}
      >
        <Ionicons
          name="person-add"
          size={variant === 'icon' ? 26 : 24}
          color={COLORS.text.inverse}
        />
      </TouchableOpacity>

      <Modal
        visible={isModalVisible}
        transparent
        animationType="fade"
        onRequestClose={handleClose}
      >
        <View style={styles.modalBackdrop}>
          <View style={styles.modalContainer}>
            <Text style={styles.modalTitle}>{label}</Text>
            <Text style={styles.modalSubtitle}>
              Assign a client by entering their username below.
            </Text>
            <InputField
              label="Username"
              placeholder="e.g. janedoe"
              value={usernameInput}
              onChangeText={(text) => {
                setError(null);
                setUsernameInput(text);
              }}
              leftIcon="person"
            />

            <Text style={styles.notesLabel}>Notes (optional)</Text>
            <TextInput
              value={notes}
              onChangeText={setNotes}
              style={styles.notesInput}
              placeholder="Share any context you want to remember."
              placeholderTextColor={COLORS.text.placeholder}
              multiline
            />

            {error && <Text style={styles.errorText}>{error}</Text>}

            <View style={styles.modalActions}>
              <Button
                title="Cancel"
                variant="outline"
                onPress={handleClose}
              />
              <Button
                title={isPending ? 'Assigning...' : 'Assign'}
                variant="primary"
                onPress={handleSubmit}
                loading={isPending}
                style={styles.modalActionButton}
              />
            </View>
          </View>
        </View>
      </Modal>
    </View>
  );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: SPACING.md,
    },
  button: {
    width: 48,
    height: 48,
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.background.accent,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    ...SHADOWS.sm,
  },
  iconButton: {
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.sm,
  },
  label: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  disabled: {
    opacity: 0.6,
  },
  modalBackdrop: {
    flex: 1,
    backgroundColor: COLORS.surface.overlay,
    justifyContent: 'flex-start',
    alignItems: 'center',
    padding: SPACING.md,
    paddingTop: Platform.OS === 'android' ? SPACING['4xl'] : SPACING['5xl'],
  },
  modalContainer: {
    width: '92%',
    maxWidth: 550,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['3xl'],
    backgroundColor: COLORS.background.auth,
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
    color: COLORS.text.inverse,
    textAlign: 'left',
    marginBottom: 12,
    marginLeft: 5,
  },
   modalSubtitle: {
    fontSize: 14,
    color: COLORS.text.tertiary,
    textAlign: 'left',
    marginBottom: 20,
    marginLeft: 5,
  },
  notesLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.auth.secondary,
    marginBottom: SPACING.xs,
  },
  notesInput: {
    minHeight: 96,
    borderRadius: BORDER_RADIUS.xl,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    padding: SPACING.md,
    color: COLORS.text.primary,
    backgroundColor: COLORS.background.surface || COLORS.background.card,
    textAlignVertical: 'top',
    marginBottom: SPACING.md,
  },
  errorText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.error,
    marginBottom: SPACING.md,
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
