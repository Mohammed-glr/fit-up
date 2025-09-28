import { useAuth } from "@/context/auth-context";
import { Redirect } from "expo-router";
import { ReactNode } from "react";
import { View, ActivityIndicator, StyleSheet } from "react-native";

interface AuthGuardProps {
    children: ReactNode;
    requireAuth?: boolean;
}

export default function AuthGuard({ children, requireAuth = true }: AuthGuardProps) {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return (
            <View style={styles.loadingContainer}>
                <ActivityIndicator size="large" color="#007AFF" />
            </View>
        );
    }

    if (requireAuth && !isAuthenticated) {
        return <Redirect href="/(auth)/login" />;
    }

    if (!requireAuth && isAuthenticated) {
        return <Redirect href="/(tabs)" />;
    }

    return <>{children}</>;
}

const styles = StyleSheet.create({
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#FFFFFF',
    },
});