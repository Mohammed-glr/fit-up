import { Tabs } from 'expo-router';
import React from 'react';
import { Platform, StyleSheet } from 'react-native';

import { AnimatedTabButton } from '@/components/animated-tab-button';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { DynamicLeftButton } from './dynamic-left-button';
import { BORDER_RADIUS, COLORS, SPACING } from '@/constants/theme';
import { DynamicRightButton } from './dynamic-right-button';
import BlurView from 'expo-blur/build/BlurView';

export function UserTabNavigator() {

  return (
    <Tabs
  screenOptions={{
           tabBarActiveTintColor: '#8FE507',
           tabBarInactiveTintColor:  '#8E8E93',
           headerShown: true,
           headerLeftContainerStyle: {
            paddingLeft: SPACING.lg,
           },
           headerRightContainerStyle: {
            paddingRight: SPACING.lg,
           }, 
           headerLeft: () => <DynamicLeftButton />,
           headerRight: () => <DynamicRightButton />,
           headerStyle: {
             backgroundColor: '#0A0A0A',
             borderBottomWidth: 0,
             height: 110,
             borderBottomColor: 'transparent',
           },
           headerTitleAlign: 'center',
           headerTintColor: '#ffffff',
           headerTitleStyle: {
             fontWeight: '600',
             fontSize: 18,
             padding: SPACING.md,
             backgroundColor: COLORS.background.card,
             borderRadius: BORDER_RADIUS.full,
           },
           tabBarShowLabel: false,
           tabBarButton: AnimatedTabButton,
            tabBarBackground: () => (
              <BlurView tint="dark" intensity={80} style={StyleSheet.absoluteFill} />
            ),
            tabBarStyle: {
            ...styles.tabBar,
            backgroundColor: 'transparent',
            borderWidth: 0,
            elevation: 0,
            },
           tabBarLabelStyle: styles.tabBarLabel,
           tabBarItemStyle: styles.tabBarItem,
         }}>


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
        name="nutrition"
        options={{
          title: 'Nutrition',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24}
              name="leaf.fill"
              color={color}
            />
          ),
          tabBarLabel: 'Fuel',
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
        name="plans"
        options={{
          title: 'Plans',
          tabBarIcon: ({ color, focused }) => (
            <IconSymbol 
              size={focused ? 28 : 24}
              name="calendar"
              color={color}
            />
          ),
          tabBarLabel: 'Plans',
        }}
      />

      {/* Hidden screens - accessible via navigation but not in tab bar */}
      <Tabs.Screen
        name="schema"
        options={{
          href: null,
          title: 'Training Schema',
          tabBarStyle: { display: 'none' },
          headerShown: true,
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
        name="recipes"
        options={{
          href: null,
          title: 'My Recipes',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="system-recipes"
        options={{
          href: null,
          title: 'Browse Recipes',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          href: null,
          title: 'Profile',
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
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="workout-session"
        options={{
          href: null,
          title: 'Workout Session',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="workout-history"
        options={{
          href: null,
          title: 'Workout History',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="progress"
        options={{
          href: null,
          title: 'Progress',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="achievements"
        options={{
          href: null,
          title: 'Achievements',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="templates"
        options={{
          href: null,
          title: 'Workout Templates',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="breathing"
        options={{
          href: null,
          title: 'Breathing Exercise',
          tabBarStyle: { display: 'none' },
          headerShown: false,
        }}
      />
      <Tabs.Screen
        name="breathing/[type]"
        options={{
          href: null,
          title: 'Breathing Exercise',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="gratitude"
        options={{
          href: null,
          title: 'Gratitude Journal',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="reflection"
        options={{
          href: null,
          title: 'Daily Reflection',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      
    </Tabs>
  );
}
const styles = StyleSheet.create({
  tabBar: {
    position: 'absolute',
    bottom: Platform.OS === 'ios' ? 20 : 16,
    marginHorizontal: 14,
    height: Platform.OS === 'ios' ? 76 : 76,
    borderRadius: BORDER_RADIUS.full,
    borderTopWidth: 0,
    paddingBottom: 12,
    paddingTop: 12,
    paddingHorizontal: 8,
    overflow: 'hidden',
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
    backgroundColor: COLORS.background.card,
  },
});

