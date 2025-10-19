import { Redirect } from 'expo-router';
import { useAuth } from '@/context/auth-context';
import { View, ActivityIndicator } from 'react-native';
import { COLORS } from '@/constants/theme';

export default function Index() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: COLORS.background.auth }}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (user) {
    if (user.role === 'coach') {
      return <Redirect href="/(coach)" />;
    }
    return <Redirect href="/(user)" />;
  }

  return <Redirect href="/(auth)/login" />;
}
