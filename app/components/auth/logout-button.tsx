import { useAuth } from "@/context/auth-context";
import { useState } from "react";
import { router } from "expo-router";
import { useQueryClient } from "@tanstack/react-query";

import { 
    Button,
} from '@/components/forms';
export default function LogoutButton() {
    const { logout } = useAuth();
    const queryClient = useQueryClient();
    const [isLoading, setIsLoading] = useState(false);

    const handleLogout = async () => {
        if (isLoading) return;
        setIsLoading(true);
        try {
            await logout();
            queryClient.clear();
            router.replace('/(auth)/login');
        } catch (error) {
            console.error("Logout failed:", error);
            queryClient.clear();
            router.replace('/(auth)/login');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <>
            <Button 
                title="Logout" 
                onPress={handleLogout} 
                loading={isLoading} 
                disabled={isLoading} 
                variant="outline"
            />
        </>
    );
}