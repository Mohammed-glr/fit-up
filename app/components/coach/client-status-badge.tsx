import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { ClientSummary } from '@/types/schema';
import { getClientStatus } from '@/utils/client-status';

interface ClientStatusBadgeProps {
  client: ClientSummary;
  size?: 'small' | 'medium';
  showLabel?: boolean;
}

export function ClientStatusBadge({ 
  client, 
  size = 'medium',
  showLabel = true 
}: ClientStatusBadgeProps) {
  const statusInfo = getClientStatus(client);
  
  const isSmall = size === 'small';
  
  return (
    <View 
      style={[
        styles.badge,
        isSmall && styles.badgeSmall,
        { backgroundColor: `${statusInfo.color}20` }
      ]}
    >
      <Ionicons
        name={statusInfo.icon as any}
        size={isSmall ? 12 : 14}
        color={statusInfo.color}
      />
      {showLabel && (
        <Text 
          style={[
            styles.label,
            isSmall && styles.labelSmall,
            { color: statusInfo.color }
          ]}
        >
          {statusInfo.label}
        </Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 10,
    paddingVertical: 5,
    borderRadius: 12,
    gap: 4,
  },
  badgeSmall: {
    paddingHorizontal: 6,
    paddingVertical: 3,
    borderRadius: 8,
    gap: 3,
  },
  label: {
    fontSize: 12,
    fontWeight: '600',
  },
  labelSmall: {
    fontSize: 10,
    fontWeight: '500',
  },
});
