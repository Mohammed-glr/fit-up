
import { authService } from "@/api/services/auth-service";
import { secureStorage } from "@/api/storage/secure-storage";
import { LoginRequest, RegisterRequest, User } from "@/types/auth";
import { createContext, useContext, useEffect, useState } from "react";

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    isEmailVerified: boolean;
    verificationMessage: string | null;
    verificationError: string | null;
    login: (credentials: LoginRequest) => Promise<User>;
    logout: () => Promise<void>;
    register: (userData: RegisterRequest) => Promise<void>;
    refreshToken: () => Promise<void>;
    getCurrentUser: () => Promise<User | null>;
    updateRole: (role: 'user' | 'coach') => Promise<void>;
    verifyEmail: (token: string) => Promise<string>;
    resendVerification: (email: string) => Promise<string>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [isEmailVerified, setIsEmailVerified] = useState<boolean>(false);
    const [verificationMessage, setVerificationMessage] = useState<string | null>(null);
    const [verificationError, setVerificationError] = useState<string | null>(null);

    useEffect(() => {
        initializeAuth();
    }, []);

    
    const initializeAuth = async () => {
        try {
        const token = await secureStorage.getToken('access_token');
        if (token) {
            const response = await authService.ValidateToken();
            setUser(response.user);
            setIsEmailVerified(!!response.user?.email_verified);
        }
        } catch (error) {
        await secureStorage.clearTokens();
        } finally {
        setIsLoading(false);
        }
    };

  const login = async (credentials: LoginRequest): Promise<User> => {
    try {
        setIsLoading(true);
        setVerificationMessage(null);
        setVerificationError(null);
        const response = await authService.Login(credentials);
        setUser(response.user);
        setIsEmailVerified(!!response.user?.email_verified);
        return response.user;
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
        setIsEmailVerified(false);
        setVerificationMessage(null);
        setVerificationError(null);
        await secureStorage.clearTokens();
    } catch (error) {
        console.error("Logout failed:", error);
        setUser(null);
        setIsEmailVerified(false);
        setVerificationMessage(null);
        setVerificationError(null);
        await secureStorage.clearTokens();
    } finally {
        setIsLoading(false);
    }
    }

    const register = async (userData: RegisterRequest) => {
        try {
            setIsLoading(true);
            setVerificationMessage(null);
            setVerificationError(null);
            await authService.Register(userData);
            setIsEmailVerified(false);
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
            setIsEmailVerified(false);
        } finally {
            setIsLoading(false);
        }
    }
    const verifyEmail = async (token: string): Promise<string> => {
        setVerificationMessage(null);
        setVerificationError(null);
        try {
            setIsLoading(true);
            const res = await authService.VerifyEmail(token);
            const message = res.message || 'Email verified successfully.';
            setVerificationMessage(message);
            if (user) {
                const response = await authService.ValidateToken();
                setUser(response.user);
                setIsEmailVerified(!!response.user?.email_verified);
            } else {
                setIsEmailVerified(true);
            }
            return message;
        } catch (error: any) {
            setVerificationError(error?.response?.data?.message || 'Verification failed');
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

    const resendVerification = async (email: string): Promise<string> => {
        setVerificationMessage(null);
        setVerificationError(null);
        try {
            setIsLoading(true);
            const res = await authService.ResendVerification(email);
            const message = res.message || 'Verification email sent.';
            setVerificationMessage(message);
            return message;
        } catch (error: any) {
            setVerificationError(error?.response?.data?.message || 'Resend failed');
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

    const getCurrentUser = async (): Promise<User | null> => {
        if (user) return user;
        return null;
    };

    const updateRole = async (role: 'user' | 'coach') => {
        try {
            setIsLoading(true);
            const response = await authService.UpdateRole(role);
            setUser(response.user);
            setIsEmailVerified(!!response.user?.email_verified);
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
        isEmailVerified,
        verificationMessage,
        verificationError,
        login,
        logout,
        register,
        refreshToken,
        getCurrentUser,
        updateRole,
        verifyEmail,
        resendVerification,
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
