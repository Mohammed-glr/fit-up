
import { authService } from "@/api/services/auth-service";
import { secureStorage } from "@/api/storage/secure-storage";
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
    getCurrentUser: () => Promise<User | null>;
    updateRole: (role: 'user' | 'coach') => Promise<void>;
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
            const response = await authService.ValidateToken();
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
        const response = await authService.Login(credentials);
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
        await authService.Logout();
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
            const response = await authService.Register(userData);
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
            await authService.RefreshToken();
        } catch (error) {
            console.error("Token refresh failed:", error);
            setUser(null);
        } finally {
            setIsLoading(false);
        }
    }

    const getCurrentUser = async (): Promise<User | null> => {
        if (user) return user;
        return null;
    };

    const updateRole = async (role: 'user' | 'coach') => {
        try {
            setIsLoading(true);
            const response = await authService.UpdateRole(role);
            setUser(response.user);
        } catch (error) {
            console.error("Update role failed:", error);
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    register,
    refreshToken,
    getCurrentUser,
    updateRole,
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
