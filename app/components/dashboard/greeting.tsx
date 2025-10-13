import { useCurrentUser } from "@/hooks/user/use-current-user";
import React from "react";

export const DashboardGreeting: React.FC = () => {
    const { data: user, isLoading } = useCurrentUser();

    if (isLoading) return <div>Loading...</div>;
    if (!user) return <div>User not found</div>;

    return <div>Welcome back, {user.name}!</div>;
};