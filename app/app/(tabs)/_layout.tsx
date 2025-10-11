import { Redirect, Tabs } from 'expo-router';
import React from 'react';
import { Platform, StyleSheet } from 'react-native';

import { HapticTab } from '@/components/haptic-tab';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { COLORS } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { useAuth } from '@/context/auth-context';

export default function TabLayout() {
  const colorScheme = useColorScheme();
  const { isAuthenticated, isLoading } = useAuth();

  // if (!isAuthenticated && !isLoading) {
  //   return <Redirect href="../(auth)/login" />;
  // }

  const isDark = colorScheme === 'dark';

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: COLORS.primary,
        tabBarInactiveTintColor: isDark ? '#8E8E93' : '#999999',
        headerShown: false,
        tabBarButton: HapticTab,
        tabBarStyle: [
          styles.tabBar,
          {
            backgroundColor: isDark 
              ? 'rgba(28, 28, 30, 0.85)' 
              : 'rgba(255, 255, 255, 0.85)',
            borderTopColor: isDark 
              ? 'rgba(84, 84, 88, 0.3)' 
              : 'rgba(0, 0, 0, 0.08)',
          }
        ],
        tabBarLabelStyle: styles.tabBarLabel,
        tabBarItemStyle: styles.tabBarItem,
      }}>
      <Tabs.Screen
        name="dashboard"
        options={{
          title: 'Dashboard',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="chart.bar.fill" 
              color={color} 
            />
          ),
          tabBarLabel: 'Dashboard',
        }}
      />
      <Tabs.Screen
        name="schema"
        options={{
          title: 'Schema',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="calendar" 
              color={color} 
            />
          ),
          tabBarLabel: 'Schema',
        }}
      />
      <Tabs.Screen
        name="messages"
        options={{
          title: 'Messages',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="message.fill" 
              color={color} 
            />
          ),
          tabBarLabel: 'Messages',
          tabBarBadge: undefined, // You can add badge count here later
        }}
      />
      <Tabs.Screen
        name="mindfullness"
        options={{
          title: 'Mindfulness',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="brain.head.profile" 
              color={color} 
            />
          ),
          tabBarLabel: 'Mind',
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Profile',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="person.crop.circle.fill" 
              color={color} 
            />
          ),
          tabBarLabel: 'Profile',
        }}
      />
    </Tabs>
  );
}

const styles = StyleSheet.create({
  tabBar: {
    position: 'absolute',
    bottom: Platform.OS === 'ios' ? 20 : 16,
    left: 20,
    right: 20,
    height: Platform.OS === 'ios' ? 88 : 72,
    borderRadius: 30,
    borderTopWidth: 0,
    paddingBottom: Platform.OS === 'ios' ? 20 : 12,
    paddingTop: 12,
    paddingHorizontal: 8,
    elevation: 12,
    shadowColor: '#000',
    shadowOffset: { 
      width: 0, 
      height: 8 
    },
    shadowOpacity: 0.25,
    shadowRadius: 16,
    overflow: 'visible',
    backdropFilter: 'blur(20px)',
  },
  tabBarLabel: {
    fontSize: 11,
    fontWeight: '700',
    marginTop: 4,
    marginBottom: 0,
    letterSpacing: 0.3,
  },
  tabBarItem: {
    paddingVertical: 8,
    marginHorizontal: 2,
    borderRadius: 16,
  },
});
