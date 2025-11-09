import React, { useState } from 'react';
import {
  TouchableOpacity,
  StyleSheet,
  View,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import {
  COLORS,
  SPACING,
  BORDER_RADIUS,
  SHADOWS,
} from '@/constants/theme';
import { ClientSearchModal } from './client-search-modal';

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
  const [isModalVisible, setModalVisible] = useState(false);

  const handleOpen = () => {
    if (disabled) return;
    setModalVisible(true);
  };

  const handleClose = () => {
    setModalVisible(false);
  };

  const handleAssigned = () => {
    setModalVisible(false);
    if (onAssigned) {
      onAssigned();
    }
  };

  return (
    <View style={styles.container}>
      <TouchableOpacity
        style={[
          styles.button,
          variant === 'icon' && styles.iconButton,
          disabled && styles.disabled,
        ]}
        onPress={handleOpen}
        activeOpacity={0.8}
        disabled={disabled}
        accessibilityLabel={label}
      >
        <Ionicons
          name="person-add"
          size={variant === 'icon' ? 26 : 24}
          color={COLORS.text.inverse}
        />
      </TouchableOpacity>

      <ClientSearchModal
        visible={isModalVisible}
        onClose={handleClose}
        onAssigned={handleAssigned}
      />
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
  disabled: {
    opacity: 0.6,
  },
});
