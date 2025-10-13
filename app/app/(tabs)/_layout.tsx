import { Redirect } from 'expo-router';
import React from 'react';

import { useAuth } from '@/context/auth-context';
import { TabNavigator } from '@/components/navigation/tab-navigator';

export default function TabLayout() {
  const { isAuthenticated, isLoading } = useAuth();

  if (!isAuthenticated && !isLoading) {
    return <Redirect href="../(auth)/login" />;
  }

  return <TabNavigator />;
}
