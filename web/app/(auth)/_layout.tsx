import { useAuth } from "@/context/auth-context";
import { Redirect, Stack } from "expo-router";


export default function AuthLayout() {
    const {isAuthenticated, isLoading} = useAuth();

    if (isAuthenticated && !isLoading) {
        return <Redirect href="/(tabs)" />;
    }


    return (
        <Stack screenOptions={{headerShown: false}}>
            <Stack.Screen name="login" />
            <Stack.Screen name="register" />
            <Stack.Screen name="forgot-password" />
            <Stack.Screen name="reset-password" />
            <Stack.Screen name="oauth-callback" />
        </Stack>
    )
}
