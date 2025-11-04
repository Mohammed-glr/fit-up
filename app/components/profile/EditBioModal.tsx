import React, { useState } from 'react';
import {
  Modal,
  View,
  Text,
  TextInput,
  StyleSheet,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from 'react-native';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { Ionicons } from '@expo/vector-icons';
import { Button } from '@/components/forms';

interface EditBioModalProps {
  visible: boolean;
  currentBio: string;
  onClose: () => void;
  onSave: (bio: string) => void;
  isLoading?: boolean;
}

export const EditBioModal: React.FC<EditBioModalProps> = ({
  visible,
  currentBio,
  onClose,
  onSave,
  isLoading = false,
}) => {
  const [bio, setBio] = useState(currentBio);
  const maxLength = 500;

  const handleSave = () => {
    onSave(bio);
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={onClose}
    >
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.modalContainer}
      >
        <TouchableOpacity
          style={styles.modalOverlay}
          activeOpacity={1}
          onPress={onClose}
        />
        <View style={styles.modalContent}>
          <View style={styles.header}>
            <Text style={styles.title}>Edit Bio</Text>
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.primary} />
            </TouchableOpacity>
          </View>

          <ScrollView showsVerticalScrollIndicator={false}>
            <View>
              <TextInput
                style={styles.textInput}
                value={bio}
                onChangeText={setBio}
                placeholder="Tell us about yourself..."
                placeholderTextColor={COLORS.text.placeholder}
                multiline
                maxLength={maxLength}
                textAlignVertical="top"
                editable={!isLoading}
              />
              <Text style={styles.charCount}>
                {bio.length}/{maxLength}
              </Text>
            </View>
          </ScrollView>

          <View style={styles.actions}>
            <Button
              title="Cancel"
              onPress={onClose}
              variant="outline"
              disabled={isLoading}
            />
            <View style={styles.saveButtonContainer}>
              <Button
                title="Save"
                onPress={handleSave}
                loading={isLoading}
                disabled={isLoading || bio === currentBio}
              />
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const styles = StyleSheet.create({
  modalContainer: {
    flex: 1,
    justifyContent: 'flex-end',
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  modalContent: {
    maxHeight: '90%',
      backgroundColor: COLORS.background.auth,
    borderTopLeftRadius: BORDER_RADIUS['3xl'],
    borderTopRightRadius: BORDER_RADIUS['3xl'],
        paddingTop: SPACING.lg,
        paddingHorizontal: SPACING.lg,
        paddingBottom: SPACING['3xl'],
        ...SHADOWS.lg,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: SPACING.md,
  },
  closeButton: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    alignItems: 'center',
    justifyContent: 'center',
    },
  title: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },

  textInput: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
    minHeight: 150,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  charCount: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'right',
    marginTop: SPACING.sm,
  },
  actions: {
    flexDirection: 'row',
    
    paddingTop: SPACING.base,
    gap: SPACING.base,
  },
  cancelButton: {
    flex: 1,
    paddingVertical: SPACING.base,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: BORDER_RADIUS.base,
    borderWidth: 1,
    borderColor: COLORS.border.medium,
  },
  cancelButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.secondary,
  },
  saveButtonContainer: {
    flex: 1,
  },
});
