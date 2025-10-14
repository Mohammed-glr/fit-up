import LoginForm from "@/components/auth/login-form";
import { useCurrentUser } from "@/hooks/user/use-current-user";
import { Redirect } from "expo-router";
import React from "react";

const Login = () => {
    const { data: currentUser, isLoading } = useCurrentUser();

    if (!isLoading && currentUser) {
        if (currentUser.role === 'coach') {
            return <Redirect href="/(coach)" />;
        }
        return <Redirect href="/(user)" />;
    }
    
    return (
        <LoginForm />
    );
};

export default Login;