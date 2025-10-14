import React from 'react';
import { UserTabNavigator } from '@/components/navigation/user-tab-navigator';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { Redirect } from 'expo-router';

export default function UserLayout() {
  const { data: currentUser, isLoading } = useCurrentUser();

  if (isLoading) {
    return null; // Or a loading screen
  }

  if (!currentUser) {
    return <Redirect href="/(auth)/login" />;
  }

  // If user is a coach, redirect to coach routes
  if (currentUser.role === 'coach') {
    return <Redirect href="/(coach)" />;
  }

  return <UserTabNavigator />
}
