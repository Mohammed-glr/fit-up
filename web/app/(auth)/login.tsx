import LoginForm from "@/components/auth/login-form";
import { useAuth } from "@/context/auth-context";
import React from "react";


const Login = () => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isAuthenticated && !isLoading) {
        return null;
    }
    return (
        <LoginForm />
    )
}