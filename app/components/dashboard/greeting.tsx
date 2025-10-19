import { useCurrentUser } from "@/hooks/user/use-current-user";
import React from "react";
import { View, Text, StyleSheet, Platform } from "react-native";
import { MotiView, MotiText } from "moti";
import { COLORS } from "@/constants/theme";

export const DashboardGreeting: React.FC = () => {
    const { data: user, isLoading } = useCurrentUser(); 
    const { timeOfDay, greeting } = (() => {
        const hour = new Date().getHours();
        if (hour < 12) return { timeOfDay: "morning", greeting: "Good morning" };
        if (hour < 18) return { timeOfDay: "afternoon", greeting: "Good afternoon" };
        return { timeOfDay: "evening", greeting: "Good evening" };
    })();
    const { clock } = (() => {
        const now = new Date();
        return { clock: now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) };
    })();

    if (isLoading) return <div>Loading...</div>;
    if (!user) return <div>User not found</div>;


    return (
        <MotiView
            style={styles.container}
            from={{
                opacity: 0,
                translateY: -20,
                scale: 0.95,
            }}
            animate={{
                opacity: 1,
                translateY: 0,
                scale: 1,
            }}
            transition={{
                type: "timing",
                duration: 600,
                delay: 100,
            }}
        >
            <MotiText
                style={styles.greeting}
                from={{
                    opacity: 0,
                    translateX: -30,
                }}
                animate={{
                    opacity: 1,
                    translateX: 0,
                }}
                transition={{
                    type: "spring",
                    damping: 15,
                    stiffness: 150,
                    delay: 300,
                }}
            >
                {greeting} {user.name}!
            </MotiText>
            <MotiText
                style={styles.time}
                from={{
                    opacity: 0,
                    translateY: 20,
                }}
                animate={{
                    opacity: 1,
                    translateY: 0,
                }}
                transition={{
                    type: "spring",
                    damping: 15,
                    stiffness: 150,
                    delay: 300,
                }}
            >
                {clock}
            </MotiText>
        </MotiView>
    );
};

const styles = StyleSheet.create({
    container: {
        borderRadius: 8,
        position: "absolute",
        top: 20,
        left: 20,
        right: 20,
        ...Platform.select({
            ios: {
                shadowColor: COLORS.black,
                shadowOffset: { width: 0, height: 2 },
                shadowOpacity: 0.1,
                shadowRadius: 4,
            },
            android: {
                elevation: 2,
            },
        }),
    },
    greeting: {
        fontSize: 30,
        fontWeight: "600",
        color: COLORS.text.inverse,
    },
    time: {
        fontSize: 24,
        fontWeight: "600",
        color: COLORS.text.inverse,
        position: 'absolute',
        top: 10,
        right: 16,
        opacity: 0.8,
        backgroundColor: COLORS.primary,
        borderRadius: 8,
        paddingHorizontal: 8,
        paddingVertical: 4,

    },
});