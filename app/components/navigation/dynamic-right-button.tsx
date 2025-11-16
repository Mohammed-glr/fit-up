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
    Platform,
    View,
    Text,
    Image
} from "react-native";
import { MotiView } from 'moti';
import { UserMenu } from './user-menu';
import { AssignClientButton } from '../coach/assign-client-button';
import { useRecipeContext } from '@/context/recipe-context';
import { useTemplateContext } from '@/context/template-context';
import { useMindfulnessContext } from '@/context/mindfulness-context';
import { useWorkoutEditorContext } from '@/context/workout-editor-context';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { useConversation } from '@/hooks/message/use-conversation';
import { useClientDetails, useCoachClients } from '@/hooks/schema/use-coach';


type RouteContext = 'coach' | 'user';

interface DynamicButtonProps {
    onNavigate?: (conversationId: number) => void;
}


export const DynamicRightButton: React.FC<DynamicButtonProps> = ({ onNavigate }) => {
    const navigation = useNavigation();
    const router = useRouter();
    const { onCreateRecipe } = useRecipeContext();
    const { onCreateTemplate } = useTemplateContext();
    const { 
        triggerGratitudeCreate, 
        triggerReflectionHistory, 
        isGratitudeWriting,
        isReflectionResponding,
        isReflectionHistory,
        onSaveGratitude,
        onSaveReflection,
        isSavingGratitude,
        isSavingReflection
    } = useMindfulnessContext();
    const { onSaveWorkout, isSavingWorkout } = useWorkoutEditorContext();
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
    
    const clientId = routeParams?.userId ? parseInt(routeParams.userId, 10) : 0;
    const { data: clientData } = useClientDetails(clientId);
    
    const { data: clientsData } = useCoachClients();
    const clientFromList = clientsData?.clients.find(c => c.user_id === clientId);

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

    if (currentRouteName === 'gratitude') {
        if (isGratitudeWriting && onSaveGratitude) {
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
                    <TouchableOpacity
                        onPress={onSaveGratitude}
                        disabled={isSavingGratitude}
                        style={styles.headerButton}
                        accessibilityLabel="Save entry"
                        accessibilityRole="button"
                        hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                    >
                        <IconSymbol 
                            name={isSavingGratitude ? "hourglass" : "checkmark"} 
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
                    type: 'timing',
                    duration: 200,
                }}
            >
                <TouchableOpacity
                    onPress={triggerGratitudeCreate}
                    style={styles.headerButton}
                    accessibilityLabel="Create new entry"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="plus" 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
        );
    }

    if (currentRouteName === 'reflection') {
        if (isReflectionResponding && onSaveReflection) {
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
                    <TouchableOpacity
                        onPress={onSaveReflection}
                        disabled={isSavingReflection}
                        style={styles.headerButton}
                        accessibilityLabel="Save response"
                        accessibilityRole="button"
                        hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                    >
                        <IconSymbol 
                            name={isSavingReflection ? "hourglass" : "checkmark"} 
                            size={24} 
                            style={styles.icon}
                            color={COLORS.text.inverse}
                        />
                    </TouchableOpacity>
                </MotiView>
            );
        }

        if (isReflectionHistory) {
            return null;
        }
        
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
                <TouchableOpacity
                    onPress={triggerReflectionHistory}
                    style={styles.headerButton}
                    accessibilityLabel="View history"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="book.fill" 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
        );
    }

    if (currentRouteName === 'chat') {
        const conversation = conversationData?.conversation;
        const isCoach = currentUser?.role === 'coach';
        const receiverName = isCoach ? conversation?.client_name : conversation?.coach_name;
        const receiverImage = isCoach ? conversation?.client_image : conversation?.coach_image;
        
        const displayName = receiverName || (isCoach ? 'Client' : 'Coach');
        const receiverInitial = displayName.charAt(0).toUpperCase();
        
        console.log('Chat Avatar Debug:', {
            conversationId,
            hasConversationData: !!conversationData,
            conversation,
            currentUserRole: currentUser?.role,
            isCoach,
            receiverName,
            displayName,
            receiverImage
        });
        
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                transition={{
                    type: 'timing',
                    duration: 250,
                }}
            >
                <View style={styles.avatarContainer}>
                    {receiverImage ? (
                        <Image
                            source={{ uri: receiverImage }}
                            style={styles.avatarImage}
                        />
                    ) : (
                        <View style={styles.avatarPlaceholder}>
                            <Text style={styles.avatarInitial}>{receiverInitial}</Text>
                        </View>
                    )}
                </View>
            </MotiView>
        );
    }
    
    if (currentRouteName === 'client-details') {

        const client = clientFromList || clientData;
        
        console.log('Client Details Avatar Debug:', {
            clientId,
            hasClientData: !!clientData,
            hasClientFromList: !!clientFromList,
            clientData,
            clientFromList,
            routeParams
        });
        
        const clientName = client && 'first_name' in client && 'last_name' in client
            ? `${client.first_name} ${client.last_name}` 
            : '';
        const clientInitial = clientName ? clientName.charAt(0).toUpperCase() : '?';
        
        return (
            <MotiView
                from={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                transition={{
                    type: 'timing',
                    duration: 250,
                }}
            >
                <View style={styles.avatarContainer}>
                    <View style={styles.avatarPlaceholder}>
                        <Text style={styles.avatarInitial}>{clientInitial}</Text>
                    </View>
                </View>
            </MotiView>
        );
    }

    if (currentRouteName === 'workout-editor') {
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
                <TouchableOpacity
                    onPress={() => {
                        if (onSaveWorkout) {
                            onSaveWorkout();
                        }
                    }}
                    style={styles.headerButton}
                    accessibilityLabel="Save workout"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                    disabled={isSavingWorkout}
                >
                    <IconSymbol 
                        name={isSavingWorkout ? "hourglass" : "checkmark"} 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
        );
    }

    if (currentRouteName === 'clients' ) {
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
                <AssignClientButton onAssigned={navigation.goBack} />
            </MotiView>
        );
    }

    if (currentRouteName === 'recipes' && onCreateRecipe) {
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
                <TouchableOpacity
                    onPress={onCreateRecipe}
                    style={styles.createRecipeButton}
                    accessibilityLabel="Create new recipe"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="plus" 
                        size={24} 
                        style={styles.icon}
                        color={COLORS.text.inverse}
                    />
                </TouchableOpacity>
            </MotiView>
        );
    }

    if (currentRouteName === 'templates' && onCreateTemplate) {
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
                <TouchableOpacity
                    onPress={onCreateTemplate}
                    style={styles.createTemplateButton}
                    accessibilityLabel="Create new template"
                    accessibilityRole="button"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <IconSymbol 
                        name="plus" 
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
           <UserMenu />
        </MotiView>
    );
};

const styles = StyleSheet.create({
    headerButton: {
        backgroundColor: COLORS.background.accent,
        padding: SPACING.md,
        marginRight: SPACING.md,
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
    createRecipeButton: {
       backgroundColor: COLORS.background.accent,
        padding: SPACING.md,
        marginRight: SPACING.md,
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
    createTemplateButton: {
        backgroundColor: COLORS.background.accent,
        padding: SPACING.md,
        marginRight: SPACING.md,
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
        width: 45,
        height: 45,
        borderRadius: BORDER_RADIUS.full,
        overflow: 'hidden',
        marginRight: SPACING.md,
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
        fontSize: 24,
        fontWeight: '600',
        color: COLORS.text.inverse,
    },
});