import {
    CreateConversationFAB
} from '@/components/chat/createConversationFAB';
import { IconSymbol } from "@/components/ui/icon-symbol";
import { BORDER_RADIUS, COLORS, SPACING } from "@/constants/theme";
import { useNavigation } from "@react-navigation/native";
import { useRouter } from "expo-router";
import React, { useMemo, useCallback } from "react";
import {
    StyleSheet,
    TouchableOpacity,
    Platform
} from "react-native";
import { MotiView } from 'moti';

type RouteContext = 'coach' | 'user';

interface DynamicButtonProps {
    onNavigate?: (conversationId: number) => void;
}


export const DynamicButton: React.FC<DynamicButtonProps> = ({ onNavigate }) => {
    const navigation = useNavigation();
    const router = useRouter();
    
    const { currentRouteName, routeContext } = useMemo(() => {
        const navState = navigation.getState();
        const route = navState?.routes[navState?.index || 0];
        const routeName = route?.name || '';
        
        const routePath = route?.params as any;
        const pathSegment = routePath?.screen || routeName;
        const context: RouteContext = pathSegment.toString().includes('coach') ? 'coach' : 'user';
        
        return {
            currentRouteName: routeName,
            routeContext: context
        };
    }, [navigation]);

    const handleBack = useCallback(() => {
        if (navigation.canGoBack()) {
            navigation.goBack();
        } else {
            router.back();
        }
    }, [navigation, router]);

    const handleConversationCreated = useCallback((conversationId: number) => {
        if (onNavigate) {
            onNavigate(conversationId);
            return;
        }

        const chatPath = routeContext === 'coach' 
            ? '/(coach)/chat' 
            : '/(user)/chat';
        
        router.push({ 
            pathname: chatPath as any, 
            params: { conversationId } 
        });
    }, [routeContext, router, onNavigate]);

    if (currentRouteName === 'conversations') {
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                transition={{
                    type: 'timing',
                    duration: 200,
                }}
            >
                <CreateConversationFAB 
                    onConversationCreated={handleConversationCreated}
                />
            </MotiView>
        );
    }

    if (currentRouteName === 'chat') {
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8, }}
                animate={{ opacity: 1, scale: 1, }}
                exit={{ opacity: 0, scale: 0.8, }}
                transition={{
                    type: 'timing',
                    duration: 250,
                }}
            >
                <TouchableOpacity
                    onPress={handleBack}
                    style={styles.headerButton}
                    accessibilityLabel="Close chat"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="xmark" 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
            
        );
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
    }

});