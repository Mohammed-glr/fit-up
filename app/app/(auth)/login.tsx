import LoginForm from "@/components/auth/login-form";
import { useAuth } from "@/context/auth-context";
import { Redirect } from "expo-router";
import React from "react";

const Login = () => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isAuthenticated && !isLoading) {
        return <Redirect href="/(tabs)" />;
    }
    
    return (
        <LoginForm />
    );
};

export default Login;