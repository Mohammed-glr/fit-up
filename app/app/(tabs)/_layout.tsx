import { Redirect, Tabs } from 'expo-router';
import React from 'react';
import { Platform, StyleSheet } from 'react-native';

import { AnimatedTabButton } from '@/components/animated-tab-button';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { COLORS } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { useAuth } from '@/context/auth-context';

export default function TabLayout() {
  const colorScheme = useColorScheme();
  const { isAuthenticated, isLoading } = useAuth();

  if (!isAuthenticated && !isLoading) {
    return <Redirect href="../(auth)/login" />;
  }

  const isDark = colorScheme === 'dark';

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: '#8FE507',
        tabBarInactiveTintColor: isDark ? '#8E8E93' : '#8FE507',
        headerShown: false,
        tabBarButton: AnimatedTabButton,
        tabBarStyle: [
          styles.tabBar,
          {
            backgroundColor: isDark 
              ? 'rgba(28, 28, 30, 0.95)' 
              : 'rgba(255, 255, 255, 0.95)',
            borderTopColor: isDark 
              ? 'rgba(84, 84, 88, 0.3)' 
              : 'rgba(0, 0, 0, 0.08)',
            borderWidth: isDark ? 0.5 : 1,
            borderColor: isDark 
              ? 'rgba(255, 255, 255, 0.1)' 
              : 'rgba(0, 0, 0, 0.06)',
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
          name="index"
          options={{
            title: 'Dashboard',
            tabBarIcon: ({ color, focused }) => (
              <IconSymbol 
                size={focused ? 28 : 24} 
                name="analytics"
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
              name={focused ? "meditation" : "brain.head.profile"} 
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

const styles = StyleSheet.create({
  tabBar: {
    position: 'absolute',
    bottom: Platform.OS === 'ios' ? 20 : 16,
    marginHorizontal: 16,
    height: Platform.OS === 'ios' ? 86 : 76,
    borderRadius: 32,
    borderTopWidth: 0,
    paddingBottom: Platform.OS === 'ios' ? 20 : 12,
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
    backdropFilter: 'blur(20px)',
  },
  tabBarLabel: {
    fontSize: 12,
    fontWeight: '800',
    marginTop: 6,
    marginBottom: 0,
    letterSpacing: 0.5,
    textAlign: 'center',
    // display: Platform.OS === 'ios' ? 'flex' : 'none', // Hide labels on Android
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
