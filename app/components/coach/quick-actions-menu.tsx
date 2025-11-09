import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Modal, Pressable, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { router } from 'expo-router';
import type { ClientSummary } from '@/types/schema';
import { COLORS, SPACING, BORDER_RADIUS, FONT_SIZES, FONT_WEIGHTS, SHADOWS } from '@/constants/theme';

interface QuickActionsMenuProps {
  client: ClientSummary;
  onRemove?: () => void;
}

interface QuickAction {
  icon: string;
  label: string;
  action: () => void;
  color?: string;
  variant?: 'default' | 'danger';
}

export function QuickActionsMenu({ client, onRemove }: QuickActionsMenuProps) {
  const [isOpen, setIsOpen] = useState(false);

  const handleAssignSchema = () => {
    setIsOpen(false);
    router.push({
      pathname: '/(coach)/schema-create',
      params: { userId: client.user_id.toString() },
    });
  };

  const handleSendMessage = () => {
    setIsOpen(false);
    router.push({
      pathname: '/(coach)/chat',
      params: { clientId: client.auth_id },
    });
  };

  const handleViewProgress = () => {
    setIsOpen(false);
    router.push({
      pathname: '/(coach)/client-details',
      params: { userId: client.user_id.toString() },
    });
  };

  const handleViewWorkouts = () => {
    setIsOpen(false);
    router.push({
      pathname: '/(coach)/client-details',
      params: { userId: client.user_id.toString(), tab: 'workouts' },
    });
  };

  const handleRemoveClient = () => {
    setIsOpen(false);
    Alert.alert(
      'Remove Client',
      `Are you sure you want to remove ${client.first_name} ${client.last_name} as your client?`,
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Remove',
          style: 'destructive',
          onPress: () => onRemove?.(),
        },
      ]
    );
  };

  const actions: QuickAction[] = [
    {
      icon: 'calendar-outline',
      label: 'Assign Schema',
      action: handleAssignSchema,
      color: COLORS.primary,
    },
    {
      icon: 'chatbubble-outline',
      label: 'Send Message',
      action: handleSendMessage,
      color: COLORS.info,
    },
    {
      icon: 'stats-chart-outline',
      label: 'View Progress',
      action: handleViewProgress,
      color: COLORS.success,
    },
    {
      icon: 'barbell-outline',
      label: 'View Workouts',
      action: handleViewWorkouts,
      color: COLORS.warning,
    },
    {
      icon: 'remove-circle-outline',
      label: 'Remove Client',
      action: handleRemoveClient,
      color: COLORS.error,
      variant: 'danger',
    },
  ];

  return (
    <>
      <TouchableOpacity
        style={styles.trigger}
        onPress={() => setIsOpen(true)}
        activeOpacity={0.7}
      >
        <Ionicons name="ellipsis-vertical" size={20} color={COLORS.text.secondary} />
      </TouchableOpacity>

      <Modal
        visible={isOpen}
        transparent
        animationType="fade"
        onRequestClose={() => setIsOpen(false)}
      >
        <Pressable style={styles.overlay} onPress={() => setIsOpen(false)}>
          <View style={styles.menu}>
            <View style={styles.header}>
              <Text style={styles.headerText}>
                {client.first_name} {client.last_name}
              </Text>
              <TouchableOpacity onPress={() => setIsOpen(false)}>
                <Ionicons name="close" size={24} color={COLORS.text.secondary} />
              </TouchableOpacity>
            </View>

            {actions.map((action, index) => (
              <TouchableOpacity
                key={index}
                style={[
                  styles.action,
                  action.variant === 'danger' && styles.actionDanger,
                  index === actions.length - 1 && styles.lastAction,
                ]}
                onPress={action.action}
                activeOpacity={0.7}
              >
                <View
                  style={[
                    styles.iconContainer,
                    { backgroundColor: `${action.color}20` },
                  ]}
                >
                  <Ionicons
                    name={action.icon as any}
                    size={20}
                    color={action.color}
                  />
                </View>
                <Text
                  style={[
                    styles.actionText,
                    action.variant === 'danger' && styles.actionTextDanger,
                  ]}
                >
                  {action.label}
                </Text>
                <Ionicons
                  name="chevron-forward"
                  size={18}
                  color={COLORS.text.tertiary}
                />
              </TouchableOpacity>
            ))}
          </View>
        </Pressable>
      </Modal>
    </>
  );
}

const styles = StyleSheet.create({
  trigger: {
    width: 32,
    height: 32,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.dark,
    justifyContent: 'center',
    alignItems: 'center',
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  menu: {
    backgroundColor: COLORS.background.card,
    borderTopLeftRadius: BORDER_RADIUS.xl,
    borderTopRightRadius: BORDER_RADIUS.xl,
    paddingBottom: SPACING.xl,
    ...SHADOWS.lg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.dark,
  },
  headerText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
  },
  action: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
    gap: SPACING.md,
  },
  lastAction: {
    borderTopWidth: 1,
    borderTopColor: COLORS.border.dark,
    marginTop: SPACING.xs,
  },
  actionDanger: {
    backgroundColor: `${COLORS.error}05`,
  },
  iconContainer: {
    width: 36,
    height: 36,
    borderRadius: BORDER_RADIUS.md,
    justifyContent: 'center',
    alignItems: 'center',
  },
  actionText: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.inverse,
  },
  actionTextDanger: {
    color: COLORS.error,
  },
});
