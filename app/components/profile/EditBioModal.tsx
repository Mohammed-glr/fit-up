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
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
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
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.primary} />
            </TouchableOpacity>
            <Text style={styles.title}>Edit Bio</Text>
            <View style={styles.placeholder} />
          </View>

          <ScrollView showsVerticalScrollIndicator={false}>
            <View style={styles.inputContainer}>
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

            <View style={styles.tipsContainer}>
              <Text style={styles.tipsTitle}>💡 Tips for a great bio:</Text>
              <Text style={styles.tipText}>• Share your fitness goals and journey</Text>
              <Text style={styles.tipText}>• Mention your favorite workouts</Text>
              <Text style={styles.tipText}>• Include any certifications or expertise</Text>
            </View>
          </ScrollView>

          <View style={styles.actions}>
            <TouchableOpacity
              style={styles.cancelButton}
              onPress={onClose}
              disabled={isLoading}
            >
              <Text style={styles.cancelButtonText}>Cancel</Text>
            </TouchableOpacity>
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
    backgroundColor: COLORS.white,
    borderTopLeftRadius: BORDER_RADIUS.xl,
    borderTopRightRadius: BORDER_RADIUS.xl,
    maxHeight: '85%',
    paddingBottom: SPACING.xl,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.base,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.subtle,
  },
  closeButton: {
    padding: SPACING.sm,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  placeholder: {
    width: 40,
  },
  inputContainer: {
    margin: SPACING.lg,
  },
  textInput: {
    backgroundColor: COLORS.background.secondary,
    borderRadius: BORDER_RADIUS.base,
    padding: SPACING.base,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.primary,
    minHeight: 150,
    borderWidth: 1,
    borderColor: COLORS.border.subtle,
  },
  charCount: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'right',
    marginTop: SPACING.sm,
  },
  tipsContainer: {
    marginHorizontal: SPACING.lg,
    padding: SPACING.base,
    backgroundColor: COLORS.background.accent,
    borderRadius: BORDER_RADIUS.base,
    marginBottom: SPACING.lg,
  },
  tipsTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
    marginBottom: SPACING.sm,
  },
  tipText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.secondary,
    marginVertical: SPACING.xs / 2,
  },
  actions: {
    flexDirection: 'row',
    paddingHorizontal: SPACING.lg,
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
