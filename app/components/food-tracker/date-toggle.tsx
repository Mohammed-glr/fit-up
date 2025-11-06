import React from 'react';
import { StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import Ionicons from '@expo/vector-icons/Ionicons';
import { BORDER_RADIUS, COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING } from '@/constants/theme';

type DateToggleProps = {
  date: string;
  onChange: (date: string) => void;
};

const parseDate = (value: string) => {
  const [year, month, day] = value.split('-').map(Number);
  return new Date(year, (month ?? 1) - 1, day);
};

const dateToKey = (value: Date) => {
  const year = value.getFullYear();
  const month = `${value.getMonth() + 1}`.padStart(2, '0');
  const day = `${value.getDate()}`.padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const addDays = (date: string, amount: number) => {
  const baseDate = parseDate(date);
  baseDate.setDate(baseDate.getDate() + amount);
  return dateToKey(baseDate);
};

const formatDisplayDate = (date: string) => {
  const parsed = parseDate(date);
  return parsed.toLocaleDateString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  });
};

export function DateToggle({ date, onChange }: DateToggleProps) {
  const handleChange = React.useCallback(
    (amount: number) => {
      const next = addDays(date, amount);
      onChange(next);
    },
    [date, onChange],
  );

  const isToday = React.useMemo(() => {
    const today = dateToKey(new Date());
    return today === date;
  }, [date]);

  return (
    <View style={styles.container}>
      <TouchableOpacity
        accessibilityRole="button"
        accessibilityLabel="View previous day"
        onPress={() => handleChange(-1)}
        style={styles.iconButton}
      >
  <Ionicons name="chevron-back" size={18} color={COLORS.text.inverse} />
      </TouchableOpacity>

      <View style={styles.centerContent}>
        <Text style={styles.dateText}>{formatDisplayDate(date)}</Text>
        <Text style={styles.subText}>{isToday ? 'Today' : 'Daily summary'}</Text>
      </View>

      <TouchableOpacity
        accessibilityRole="button"
        accessibilityLabel="View next day"
        onPress={() => handleChange(1)}
        style={[styles.iconButton, isToday && styles.disabledButton]}
        disabled={isToday}
      >
  <Ionicons name="chevron-forward" size={18} color={isToday ? COLORS.text.placeholder : COLORS.text.tertiary} />
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.darkGray,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
  },
  iconButton: {
    width: 36,
    height: 36,
    borderRadius: 18,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.background.accent,
  },
  centerContent: {
    flex: 1,
    alignItems: 'center',
  },
  dateText: {
    color: COLORS.text.inverse,
  fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    letterSpacing: 0.2,
  },
  subText: {
    marginTop: 2,
  color: COLORS.text.tertiary,
    fontSize: FONT_SIZES.xs,
  },
  disabledButton: {
    backgroundColor: COLORS.background.secondary,
  },
});
