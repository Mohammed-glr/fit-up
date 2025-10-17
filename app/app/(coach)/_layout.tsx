import React from 'react';
import { CoachTabNavigator } from '@/components/navigation/coach-tab-navigator';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { Redirect } from 'expo-router';

export default function CoachLayout() {
  const { data: currentUser, isLoading } = useCurrentUser();

  if (isLoading) {
    return null;
  }

  if (!currentUser) {
    return <Redirect href="/(auth)/login" />;
  }

  if (currentUser.role !== 'coach') {
    return <Redirect href="/(user)" />;
  }

  return <CoachTabNavigator />
}
