import { Tabs } from 'expo-router';
import React from 'react';
import { Platform, StyleSheet } from 'react-native';

import { AnimatedTabButton } from '@/components/animated-tab-button';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { DynamicButton } from './dynamic-button';
import { BORDER_RADIUS, COLORS, SPACING } from '@/constants/theme';

export function UserTabNavigator() {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'light';

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: '#8FE507',
        tabBarInactiveTintColor: isDark ? '#8FE507' : '#8E8E93',
        headerShown: true,
        headerLeft: () => <DynamicButton />,
        headerStyle: {
          backgroundColor: isDark ? '#0A0A0A' : '#0A0A0A',
          borderBottomColor: 'transparent',
          shadowColor: 'transparent',
          height: 110,
        },
        headerTintColor: isDark ? '#000000ff' : '#ffffff',
        headerTitleStyle: {
          fontWeight: '600',
          fontSize: 18,
          padding: SPACING.md,
          backgroundColor: isDark ? COLORS.background.card : COLORS.background.accent,
          borderRadius: BORDER_RADIUS.full,
        },
        tabBarButton: AnimatedTabButton,
        tabBarStyle: [
          styles.tabBar,
          {
            backgroundColor: isDark 
              ? 'rgba(255, 255, 255, 0.95)'
              : 'rgba(28, 28, 30, 0.95)',
            borderWidth: isDark ? 0.5 : 1,
            borderColor: isDark 
              ? 'rgba(0, 0, 0, 0.06)'
              : 'rgba(255, 255, 255, 0.1)',
          }
        ],
        tabBarLabelStyle: styles.tabBarLabel,
        tabBarItemStyle: styles.tabBarItem,
      }}>

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
          tabBarLabel: 'Plan',
        }}
      />
      <Tabs.Screen
        name="conversations"
        options={{
          title: 'Conversations',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="message.fill" 
              color={color} 
            />
          ),
          tabBarLabel: 'Chat',
          tabBarBadge: undefined, 
        }}
      />
      <Tabs.Screen
        name="chat"
        options={{
          href: null,
          title: 'Chat',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="plan-generator"
        options={{
          href: null,
          title: 'Plan Generator',
          tabBarStyle: { display: 'none' },
          headerShown: false,
        }}
      />
      <Tabs.Screen
        name="index"
        options={{
          title: 'Dashboard',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24} 
              name="house.fill"
              color={color} 
            />
          ),
          tabBarLabel: 'Home',
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
          tabBarLabel: 'Me',
        }}
      />
    </Tabs>
  );
}
export const styles = StyleSheet.create({
  tabBar: {
    position: 'absolute',
    bottom: Platform.OS === 'ios' ? 20 : 16,
    marginHorizontal: 16,
    height: Platform.OS === 'ios' ? 76 : 76,
    borderRadius: BORDER_RADIUS.full,
    borderTopWidth: 0,
    paddingBottom: 12,
    paddingTop: 12,
    paddingHorizontal: 8,
    elevation: 20,
    shadowColor: '#000',
    shadowOffset: { 
      width: 0, 
      height: 12 
    },
    shadowOpacity: 0.35,
    shadowRadius: 24,
    overflow: 'visible',
    backdropFilter: 'blur(30px)',

  },
  tabBarLabel: {
    fontSize: 12,
    fontWeight: '800',
    marginTop: 6,
    marginBottom: 0,
    letterSpacing: 0.5,
    textAlign: 'center',
    display: Platform.OS === 'ios' ? 'flex' : 'none',
  },
  tabBarItem: {
    flex: 1,
    paddingVertical: 0,
    marginHorizontal: 4,
    borderRadius: 20,
    alignItems: 'center',
    justifyContent: 'center',
  },
});
