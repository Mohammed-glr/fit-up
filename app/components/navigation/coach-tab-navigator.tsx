import { Tabs } from 'expo-router';
import React from 'react';
import { Platform, StyleSheet } from 'react-native';

import { AnimatedTabButton } from '@/components/animated-tab-button';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { BORDER_RADIUS, COLORS, SPACING } from '@/constants/theme';
import { DynamicLeftButton } from './dynamic-left-button';
import { DynamicRightButton } from './dynamic-right-button';


export function CoachTabNavigator() {


  return (
       <Tabs
         screenOptions={{
           tabBarActiveTintColor: '#8FE507',
           tabBarInactiveTintColor:  '#8E8E93',
           headerShown: true,
           headerLeft: () => <DynamicLeftButton />,
           headerRight: () => <DynamicRightButton />,
           headerStyle: {
             backgroundColor: '#0A0A0A',
             borderBottomColor: 'transparent',
             shadowColor: 'transparent',
             height: 110,
           },
           headerTintColor: '#ffffff',
           headerTitleStyle: {
             fontWeight: '600',
             fontSize: 18,
             padding: SPACING.md,
             backgroundColor: COLORS.background.card,
             borderRadius: BORDER_RADIUS.full,
           },
           tabBarButton: AnimatedTabButton,
           tabBarStyle: [
             styles.tabBar,
             {
               backgroundColor: 'rgba(28, 28, 30, 0.95)',
               borderWidth: 1,
               borderColor: 'rgba(255, 255, 255, 0.1)',
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
                name="house.fill"
                color={color} 
              />
            ),
            tabBarLabel: 'Home',
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
        name="clients"
        options={{
          href: null,
          title: 'clients',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
        <Tabs.Screen
        name="client-details"
        options={{
          href: null,
          title: 'Client Details',
          tabBarStyle: { display: 'none' },
        headerShown: true,
        }}
      />
        <Tabs.Screen
        name="schema-create"
        options={{
          href: null,
          title: 'Create Schema',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
        <Tabs.Screen
        name="schema-templates"
        options={{
          href: null,
          title: 'Templates',
          tabBarStyle: { display: 'none' },
          headerShown: true,
        }}
      />
      <Tabs.Screen
        name="workout-editor"
        options={{
          href: null,
          title: 'Edit Workout',
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
    </Tabs>
  );
}

const styles = StyleSheet.create({
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
