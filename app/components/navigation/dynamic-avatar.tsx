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

interface DynamicAvatarProps {
    onNavigate?: (conversationId: number) => void;
}


export const DynamicAvatar: React.FC<DynamicAvatarProps> = ({ onNavigate }) => {
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
            ? `/(${routeContext})/chat`
            : `/(${routeContext})/chat`;
        router.push({
            pathname: chatPath,
            params: { conversationId: String(conversationId) },
        });
    }, [onNavigate, routeContext, router]);

    return ( 
        <CreateConversationFAB
            onConversationCreated={handleConversationCreated}
            iconProps={ {
