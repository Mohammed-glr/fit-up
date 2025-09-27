import { useAuth } from "@/context/auth-context";
import { Redirect, Stack } from "expo-router";
import { StyleSheet } from 'react-native';
import { ThemedView } from '@/components/themed-view';

export default function AuthLayout() {
    const {isAuthenticated, isLoading} = useAuth();

    if (isAuthenticated && !isLoading) {
        return <Redirect href="/(tabs)" />;
    }

    return (
        <ThemedView style={styles.container} fullScreen>
            <Stack screenOptions={{
                headerShown: false,
                contentStyle: { backgroundColor: 'transparent' }
            }}>
                <Stack.Screen name="login" />
                <Stack.Screen name="register" />
                <Stack.Screen name="forgot-password" />
                <Stack.Screen name="reset-password" />
                <Stack.Screen name="oauth-callback" />
            </Stack>
        </ThemedView>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
});
