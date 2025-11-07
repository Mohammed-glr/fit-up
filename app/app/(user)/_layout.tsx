import React from 'react';
import { UserTabNavigator } from '@/components/navigation/user-tab-navigator';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { Redirect } from 'expo-router';
import { RecipeProvider } from '@/context/recipe-context';

export default function UserLayout() {
  const { data: currentUser, isLoading } = useCurrentUser();

  if (isLoading) {
    return null;
  }

  if (!currentUser) {
    return <Redirect href="/(auth)/login" />;
  }

  if (currentUser.role === 'coach') {
    return <Redirect href="/(coach)" />;
  }

  return (
    <RecipeProvider>
      <UserTabNavigator />
    </RecipeProvider>
  );
}
