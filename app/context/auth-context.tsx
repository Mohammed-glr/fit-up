import { authService } from "@/services/api/auth-service";
import { secureStorage } from "@/services/storage/secure-storage";
import { LoginRequest, RegisterRequest, User } from "@/types/auth";
import { createContext, useContext, useEffect, useState } from "react";

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (credentials: LoginRequest) => Promise<void>;
    logout: () => Promise<void>;
    register: (userData: RegisterRequest) => Promise<void>;
    refreshToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    useEffect(() => {
        initializeAuth();
    }, []);


    const initializeAuth = async () => {
        try {
        const token = await secureStorage.getToken('access_token');
        if (token) {
            const response = await authService.validateToken();
            setUser(response.user);
        }
        } catch (error) {
        await secureStorage.clearTokens();
        } finally {
        setIsLoading(false);
        }
    };

  const login = async (credentials: LoginRequest) => {
    try {
        setIsLoading(true);
        const response = await authService.login(credentials);
        setUser(response.user);
        } catch (error) {
            console.error("Login failed:", error);
            throw error;
        } finally {
            setIsLoading(false);
        }
    }

  const logout = async () => {
    try {
        setIsLoading(true);
        await authService.logout();
        setUser(null);
        await secureStorage.clearTokens();
    } catch (error) {
        console.error("Logout failed:", error);
        setUser(null);
        await secureStorage.clearTokens();
    } finally {
        setIsLoading(false);
    }
    }

    const register = async (userData: RegisterRequest) => {
        try {
            setIsLoading(true);
            const response = await authService.register(userData);
            setUser(response.user);
        } catch (error) {
            console.error("Registration failed:", error);
            throw error;
        } finally {
            setIsLoading(false);
        }
    }

    const refreshToken = async () => {
        try {
            setIsLoading(true);
            await authService.refreshToken();
        } catch (error) {
            console.error("Token refresh failed:", error);
            setUser(null);
        } finally {
            setIsLoading(false);
        }
    }


const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    register,
    refreshToken,
};

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
