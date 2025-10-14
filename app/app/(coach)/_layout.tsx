import React from 'react';
import { CoachTabNavigator } from '@/components/navigation/coach-tab-navigator';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { Redirect } from 'expo-router';

export default function CoachLayout() {
  const { data: currentUser, isLoading } = useCurrentUser();

  if (isLoading) {
    return null; // Or a loading screen
  }

  if (!currentUser) {
    return <Redirect href="/(auth)/login" />;
  }

  // If user is not a coach, redirect to user routes
  if (currentUser.role !== 'coach') {
    return <Redirect href="/(user)" />;
  }

  return <CoachTabNavigator />
}
