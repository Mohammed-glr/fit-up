
import { IconSymbol } from "@/components/ui/icon-symbol";
import { BORDER_RADIUS, COLORS, SPACING } from "@/constants/theme";
import { useNavigation } from "@react-navigation/native";
import { useRouter } from "expo-router";
import React, { useMemo, useCallback } from "react";
import {
    StyleSheet,
    TouchableOpacity,
    Platform,
    View,
    Text,
    Image
} from "react-native";
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { useConversation } from '@/hooks/message/use-conversation';

type RouteContext = 'coach' | 'user';

interface DynamicButtonProps {
    onNavigate?: (conversationId: number) => void;
}


export const DynamicLeftButton: React.FC<DynamicButtonProps> = ({ onNavigate }) => {
    const navigation = useNavigation();
    const router = useRouter();
    const { data: currentUser } = useCurrentUser();
    
    const { currentRouteName, routeContext, routeParams } = useMemo(() => {
        const navState = navigation.getState();
        const route = navState?.routes[navState?.index || 0];
        const routeName = route?.name || '';
        
        const routePath = route?.params as any;
        const pathSegment = routePath?.screen || routeName;
        const context: RouteContext = pathSegment.toString().includes('coach') ? 'coach' : 'user';
        
        return {
            currentRouteName: routeName,
            routeContext: context,
            routeParams: route?.params as any
        };
    }, [navigation]);
    
    const conversationId = routeParams?.conversationId ? Number(routeParams.conversationId) : 0;
    const { data: conversationData } = useConversation(conversationId);

    const handleBack = useCallback(() => {
        if (navigation.canGoBack()) {
            router.back();
        }
    }, [navigation, router]);
    
    if (currentRouteName === 'workout-session') {
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                transition={{
                    type: 'spring',
                    damping: 15,
                    stiffness: 150,
                }}
            > <TouchableOpacity onPress={() => router.back()}>
                  <Ionicons 
                  name="close" 
                  size={24} 
                  color={COLORS.text.inverse} 
                  style={{ marginLeft: SPACING.md, padding: SPACING.md, backgroundColor: COLORS.primaryDark, borderRadius: BORDER_RADIUS.full}} />
                </TouchableOpacity>
            </MotiView>
        )
    }

    if (currentRouteName === 'gratitude') {
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                transition={{
                    type: 'spring',
                    damping: 15,
                    stiffness: 150,
                }}
            >
                <TouchableOpacity
                    onPress={handleBack}
                    style={styles.headerButton}
                    accessibilityLabel="Go back"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="chevron.left" 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
        );
    }

    if (currentRouteName === 'reflection' || currentRouteName === 'breathing' || currentRouteName === 'workout-editor') {
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                transition={{
                    type: 'spring',
                    damping: 15,
                    stiffness: 150,
                }}
            >
                <TouchableOpacity
                    onPress={handleBack}
                    style={styles.headerButton}
                    accessibilityLabel="Go back"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="chevron.left" 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
        );
    }

   
    
    if ( navigation.canGoBack() === false ) {
        return null;
    }

    return (

        <MotiView
            from={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.8 }}
            transition={{
                type: 'spring',
                damping: 15,
                stiffness: 150,
            }}
        >
            <TouchableOpacity
                onPress={handleBack}
                style={styles.headerButton}
                accessibilityLabel="Go back"
                accessibilityRole="button"
                hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
            >
                <IconSymbol 
                    name="chevron.left" 
                    size={24} 
                    style={styles.icon}
                    color={COLORS.text.inverse}
                />
            </TouchableOpacity>
        </MotiView>
    );
};

const styles = StyleSheet.create({
    headerButton: {
        backgroundColor: COLORS.background.accent,
        padding: SPACING.md,
        marginLeft: SPACING.md,
        borderRadius: BORDER_RADIUS.full,
        justifyContent: 'center',
        alignItems: 'center',
        minWidth: 40,
        minHeight: 40,
        ...Platform.select({
            ios: {
            },
            android: {
                elevation: 0,
            },
        }),
    },
    icon: {
        fontWeight: '600',
    },
    avatarContainer: {
        width: 40,
        height: 40,
        borderRadius: 20,
        overflow: 'hidden',
        marginLeft: SPACING.md,
    },
    avatarImage: {
        width: '100%',
        height: '100%',
    },
    avatarPlaceholder: {
        width: '100%',
        height: '100%',
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: COLORS.background.accent,
    },
    avatarInitial: {
        fontSize: 18,
        fontWeight: '600',
        color: COLORS.text.inverse,
    },
});