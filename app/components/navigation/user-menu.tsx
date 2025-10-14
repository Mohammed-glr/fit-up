import {
    View,
    Text,
    StyleSheet,
    Platform,
    TouchableOpacity,
    Modal,
    Pressable
} from "react-native";
import { MotiView } from "moti";
import { COLORS, SPACING, FONT_SIZES } from "@/constants/theme";
import { useCurrentUser } from "@/hooks/user/use-current-user"; 
import React from "react";
import { Avatar } from "@/components/ui/avatar";
import { useRouter } from "expo-router";
import LogoutButton from "../auth/logout-button";

export const UserMenu: React.FC = () => {
    const { data: user, isLoading } = useCurrentUser();
    const router = useRouter();
    const [isOpen, setIsOpen] = React.useState(false);

    if (isLoading) {
        return (
            <View style={styles.container}>
                <View style={styles.avatarPlaceholder}>
                    <Text style={styles.loadingText}>...</Text>
                </View>
            </View>
        );
    }
    
    if (!user) return null

    const handleProfile = () => {
        setIsOpen(false);
        router.push('/(tabs)/profile');
    };

    return (
        <View style={styles.container}>
            <TouchableOpacity 
                onPress={() => setIsOpen(!isOpen)}
                activeOpacity={0.7}
            >
                <Avatar />
            </TouchableOpacity>

            <Modal
                visible={isOpen}
                transparent={true}
                animationType="fade"
                onRequestClose={() => setIsOpen(false)}
            >
                <Pressable 
                    style={styles.modalOverlay}
                    onPress={() => setIsOpen(false)}
                >
                    <MotiView
                        from={{ opacity: 0, translateY: -20, scale: 0.9 }}
                        animate={{ opacity: 1, translateY: 0, scale: 1 }}
                        exit={{ opacity: 0, translateY: -20, scale: 0.9 }}
                        transition={{ type: 'timing', duration: 200 }}
                        style={styles.dropdown}
                    >
                        <View style={styles.userInfo}>
                            <Text style={styles.userName}>{user.username}</Text>
                            <Text style={styles.userEmail}>{user.email}</Text>
                            <Text style={styles.userRole}>{user.role}</Text>
                        </View>

                        <View style={styles.divider} />

                        <View style={styles.divider} />

                            <LogoutButton />
                    </MotiView>
                </Pressable>
            </Modal>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        paddingRight: 16,
        justifyContent: 'center',
        alignItems: 'center',
    },
    avatarPlaceholder: {
        width: 40,
        height: 40,
        borderRadius: 20,
        backgroundColor: COLORS.border.light,
        justifyContent: 'center',
        alignItems: 'center',
    },
    loadingText: {
        color: COLORS.text.primary,
        fontSize: 14,
    },
    modalOverlay: {
        flex: 1,
        backgroundColor: 'rgba(0, 0, 0, 0.5)',
        justifyContent: 'flex-start',
        alignItems: 'flex-end',
        paddingTop: Platform.OS === 'ios' ? 90 : 70,
        paddingRight: 16,
    },
    dropdown: {
        backgroundColor: COLORS.background.surface,
        borderRadius: 24,
        minWidth: 250,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.3,
        shadowRadius: 8,
        elevation: 8,
        overflow: 'hidden',
    },
    userInfo: {
        padding: SPACING.base,
        backgroundColor: COLORS.background.secondary,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.border.subtle,
    },
    userName: {
        fontSize: FONT_SIZES.lg,
        fontWeight: '600',
        color: COLORS.text.primary,
        marginBottom: 4,
    },
    userEmail: {
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.tertiary,
        marginBottom: 4,
    },
    userRole: {
        fontSize: FONT_SIZES.xs,
        color: COLORS.primary,
        fontWeight: '600',
        textTransform: 'uppercase',
    },
    divider: {
        height: 1,
        backgroundColor: COLORS.border.subtle,
    },
    menuItem: {
        paddingVertical: SPACING.base,
        paddingHorizontal: SPACING.base,
    },
    menuItemText: {
        fontSize: FONT_SIZES.base,
        color: COLORS.text.primary,
    },
    logoutItem: {
        backgroundColor: COLORS.background.secondary,
    },
    logoutText: {
        color: COLORS.error,
    },
});