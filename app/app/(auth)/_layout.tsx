import { Redirect, Stack } from "expo-router";
import { StyleSheet } from 'react-native';
import { ThemedView } from '@/components/themed-view';
import { useCurrentUser } from '@/hooks/user/use-current-user';

export default function AuthLayout() {
    const { data: currentUser, isLoading } = useCurrentUser();

    if (!isLoading && currentUser) {
        if (currentUser.role === 'coach') {
            return <Redirect href="/(coach)" />;
        }
        return <Redirect href="/(user)" />;
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
