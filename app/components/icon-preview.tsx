import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { IconSymbol } from '@/components/ui/icon-symbol';

/**
 * Icon Preview Component - Use this to test different icons
 * Import and render this component anywhere to see icon options
 */
export function IconPreview() {
  const iconSets = {
    'Dashboard Icons': [
      { name: 'analytics', label: 'Analytics' },
      { name: 'pulse', label: 'Pulse' },
      { name: 'chart.bar.fill', label: 'Bar Chart' },
      { name: 'dashboard-alt', label: 'Pie Chart' },
    ],
    'Wellness Icons': [
      { name: 'brain.head.profile', label: 'Brain' },
      { name: 'meditation', label: 'Flower' },
      { name: 'spa', label: 'Spa' },
      { name: 'heart-pulse', label: 'Heart Pulse' },
    ],
    'Communication Icons': [
      { name: 'message.fill', label: 'Chat Bubbles' },
      { name: 'chat', label: 'Chat Box' },
      { name: 'paperplane.fill', label: 'Send' },
    ],
    'Profile Icons': [
      { name: 'person.crop.circle.fill', label: 'Person Circle' },
    ],
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>🎨 Available Icons</Text>
      {Object.entries(iconSets).map(([category, icons]) => (
        <View key={category} style={styles.category}>
          <Text style={styles.categoryTitle}>{category}</Text>
          <View style={styles.iconGrid}>
            {icons.map((icon) => (
              <View key={icon.name} style={styles.iconItem}>
                <View style={styles.iconCircle}>
                  <IconSymbol name={icon.name as any} size={28} color="#8FE507" />
                </View>
                <Text style={styles.iconLabel}>{icon.label}</Text>
                <Text style={styles.iconName}>{icon.name}</Text>
              </View>
            ))}
          </View>
        </View>
      ))}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0A0A0A',
    padding: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#FFFFFF',
    marginBottom: 24,
  },
  category: {
    marginBottom: 32,
  },
  categoryTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#8FE507',
    marginBottom: 16,
  },
  iconGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 16,
  },
  iconItem: {
    alignItems: 'center',
    width: 100,
  },
  iconCircle: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: 'rgba(143, 229, 7, 0.1)',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 8,
  },
  iconLabel: {
    fontSize: 12,
    fontWeight: '600',
    color: '#FFFFFF',
    textAlign: 'center',
  },
  iconName: {
    fontSize: 9,
    color: '#666',
    textAlign: 'center',
    marginTop: 4,
  },
});
