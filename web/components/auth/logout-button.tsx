import { useAuth } from "@/context/auth-context";
import { useState } from "react";

import { 
    Button,
} from '@/components/forms';
export default function LogoutButton() {
    const { logout } = useAuth();
    const [isLoading, setIsLoading] = useState(false);

    const handleLogout = async () => {
        if (isLoading) return;
        setIsLoading(true);
        try {
            await logout();
        } catch (error) {
            console.error("Logout failed:", error);
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
                variant="secondary"
            />
        </>
    );
}