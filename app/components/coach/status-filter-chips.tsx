import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ScrollView } from 'react-native';
import type { ClientStatus } from '@/utils/client-status';
import { COLORS, SPACING, BORDER_RADIUS, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';

interface StatusFilterChipsProps {
  activeFilter: ClientStatus | 'all';
  onFilterChange: (filter: ClientStatus | 'all') => void;
  counts: Record<ClientStatus | 'all', number>;
}

const FILTER_OPTIONS: Array<{
  key: ClientStatus | 'all';
  label: string;
  color: string;
}> = [
  { key: 'all', label: 'All', color: COLORS.text.secondary },
  { key: 'active', label: 'Active', color: '#10B981' },
  { key: 'needs_attention', label: 'Needs Attention', color: '#F59E0B' },
  { key: 'inactive', label: 'Inactive', color: '#EF4444' },
  { key: 'no_schema', label: 'No Schema', color: '#8B5CF6' },
];

export function StatusFilterChips({
  activeFilter,
  onFilterChange,
  counts,
}: StatusFilterChipsProps) {
  return (
    <ScrollView
      horizontal
      showsHorizontalScrollIndicator={false}
      contentContainerStyle={styles.container}
    >
      {FILTER_OPTIONS.map((option) => {
        const isActive = activeFilter === option.key;
        const count = counts[option.key] || 0;

        return (
          <View
            key={option.key}
            style={styles.chipContainer}
            >
          <TouchableOpacity
            key={option.key}
            style={[
              styles.chip,
              isActive && { 
                backgroundColor: option.color,
                borderColor: option.color,
              },
            ]}
            onPress={() => onFilterChange(option.key)}
            activeOpacity={0.7}
          >
            <Text
              style={[
                styles.chipText,
                isActive && styles.chipTextActive,
              ]}
            >
              {option.label}
            </Text>
            <View
              style={[
                styles.countBadge,
                isActive && styles.countBadgeActive,
              ]}
            >
              <Text
                style={[
                  styles.countText,
                  isActive && styles.countTextActive,
                ]}
              >
                {count}
              </Text>
            </View>
          </TouchableOpacity>
            </View>
        );
      })}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    gap: SPACING.sm,
  },
  chipContainer: {
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.darkGray,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
    gap: SPACING.xs,
  },
  chipText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.inverse,
  },
  chipTextActive: {
    color: COLORS.text.primary,
    fontWeight: FONT_WEIGHTS.semibold as any,
  },
  countBadge: {
    minWidth: 20,
    height: 20,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.dark,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 6,
  },
  countBadgeActive: {
    backgroundColor: 'rgba(255, 255, 255, 0.3)',
  },
  countText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.tertiary,
  },
  countTextActive: {
    color: COLORS.text.primary,
  },
});
