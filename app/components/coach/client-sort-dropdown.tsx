import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Modal, Pressable } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, BORDER_RADIUS, FONT_SIZES, FONT_WEIGHTS, SHADOWS } from '@/constants/theme';
import { SORT_OPTIONS, type SortOption, type SortOrder } from '@/utils/client-sorting';

interface ClientSortDropdownProps {
  sortBy: SortOption;
  sortOrder: SortOrder;
  onSortChange: (sortBy: SortOption, sortOrder: SortOrder) => void;
}

export function ClientSortDropdown({ sortBy, sortOrder, onSortChange }: ClientSortDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);

  const currentSort = SORT_OPTIONS.find((opt) => opt.key === sortBy);

  const handleOptionSelect = (option: SortOption) => {
    // If selecting the same option, toggle order
    if (option === sortBy) {
      onSortChange(option, sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      // New option, default to descending (most relevant first)
      onSortChange(option, 'desc');
    }
    setIsOpen(false);
  };

  return (
    <>
      <TouchableOpacity
        style={styles.trigger}
        onPress={() => setIsOpen(true)}
        activeOpacity={0.7}
      >
        <Ionicons name="swap-vertical" size={18} color={COLORS.text.secondary} />
        <Text style={styles.triggerText}>
          {currentSort?.label || 'Sort'} {sortOrder === 'asc' ? '↑' : '↓'}
        </Text>
        <Ionicons name="chevron-down" size={16} color={COLORS.text.tertiary} />
      </TouchableOpacity>

      <Modal
        visible={isOpen}
        transparent
        animationType="fade"
        onRequestClose={() => setIsOpen(false)}
      >
        <Pressable style={styles.overlay} onPress={() => setIsOpen(false)}>
          <View style={styles.dropdown}>
            <View style={styles.header}>
              <Text style={styles.headerText}>Sort By</Text>
              <TouchableOpacity onPress={() => setIsOpen(false)}>
                <Ionicons name="close" size={24} color={COLORS.text.secondary} />
              </TouchableOpacity>
            </View>

            {SORT_OPTIONS.map((option) => {
              const isActive = option.key === sortBy;
              
              return (
                <TouchableOpacity
                  key={option.key}
                  style={[styles.option, isActive && styles.optionActive]}
                  onPress={() => handleOptionSelect(option.key)}
                  activeOpacity={0.7}
                >
                  <View style={styles.optionLeft}>
                    <Ionicons
                      name={option.icon as any}
                      size={20}
                      color={isActive ? COLORS.primary : COLORS.text.secondary}
                    />
                    <Text style={[styles.optionText, isActive && styles.optionTextActive]}>
                      {option.label}
                    </Text>
                  </View>
                  {isActive && (
                    <View style={styles.orderIndicator}>
                      <Ionicons
                        name={sortOrder === 'asc' ? 'arrow-up' : 'arrow-down'}
                        size={16}
                        color={COLORS.primary}
                      />
                    </View>
                  )}
                </TouchableOpacity>
              );
            })}
          </View>
        </Pressable>
      </Modal>
    </>
  );
}

const styles = StyleSheet.create({
  trigger: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    gap: SPACING.xs,
  },
  triggerText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  dropdown: {
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
  option: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
  },
  optionActive: {
    backgroundColor: `${COLORS.primary}15`,
  },
  optionLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.md,
  },
  optionText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
  },
  optionTextActive: {
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.semibold as any,
  },
  orderIndicator: {
    width: 28,
    height: 28,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: `${COLORS.primary}20`,
    justifyContent: 'center',
    alignItems: 'center',
  },
});
