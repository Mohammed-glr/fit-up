import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { SplashScreen, Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import 'react-native-reanimated';

import { useColorScheme } from '@/hooks/use-color-scheme';
import { AuthProvider } from '@/context/auth-context';
import { ToastProvider } from '@/components/ui/toast-provider';

SplashScreen.preventAutoHideAsync();

export const unstable_settings = {
  anchor: '(tabs)',
};

export default function RootLayout() {
  const colorScheme = useColorScheme();

  return (
    <AuthProvider>
      <ToastProvider maxToasts={3}>
        <ThemeProvider value={colorScheme === 'dark' ? DarkTheme : DefaultTheme}>
          <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="(tabs)" />
            <Stack.Screen name="(auth)" />
            <Stack.Screen name="modal" options={{ presentation: 'modal', title: 'Modal' }} />
            <Stack.Screen name="not-found" options={{ title: 'Oops!' }} />
          </Stack>
          <StatusBar style="auto" />
        </ThemeProvider>
      </ToastProvider>
    </AuthProvider>
  );
}
