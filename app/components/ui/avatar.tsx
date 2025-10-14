import {
    View,
    Text,
    Image,
    StyleSheet,
    Platform
} from "react-native";
import { MotiView, MotiText } from "moti";
import { COLORS } from "@/constants/theme";
import { useCurrentUser } from "@/hooks/user/use-current-user"; 
import React from "react";


export const Avatar = () => {
    const { data: user, isLoading } = useCurrentUser();

    const initial = user?.name ? user.name.charAt(0).toUpperCase() : "?";

    if (isLoading) return (
        <View style={styles.avatarContainer}>
            <View style={styles.avatarPlaceholder}>
                <Text style={styles.avatarInitial}>...</Text>
            </View>
        </View>
    );
    if (!user) return null;
    return (
        <MotiView
            style={styles.avatarContainer}
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
            {user.image ? (
                <Image
                    source={{ uri: user.image }}
                    style={styles.avatarImage}
                />
            ) : (
                <View style={styles.avatarPlaceholder}>
                    <Text style={styles.avatarInitial}>{initial}</Text>
                </View>
            )}
        </MotiView>
    );
};

const styles = StyleSheet.create({
    avatarContainer: {
        width: 48,
        height: 48,
        borderRadius: 24,
        overflow: "hidden",
    },
    avatarImage: {
        width: "100%",
        height: "100%",
        borderRadius: 24,
    },
    avatarPlaceholder: {
        width: "100%",
        height: "100%",
        justifyContent: "center",
        alignItems: "center",
        backgroundColor: 'rgba(28, 28, 30, 0.95)',
    },
    avatarInitial: {
        fontSize: 24,
        color: COLORS.white,
    },
});
