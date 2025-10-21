import React, { memo, useMemo } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Image } from 'react-native';
import type { ConversationOverview } from '@/types';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { useAuth } from '@/context/auth-context';

interface ConversationItemProps {
    conversation: ConversationOverview;
    onPress: () => void;
}

const formatTime = (isoDate: string | null | undefined) => {
    if (!isoDate) return '';
    
    const date = new Date(isoDate);
    if (Number.isNaN(date.getTime())) return '';
    
    const now = new Date();
    const diffInMs = now.getTime() - date.getTime();
    const diffInHours = diffInMs / (1000 * 60 * 60);
    
    if (diffInHours < 24) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    
    if (diffInHours < 168) {
        return date.toLocaleDateString([], { weekday: 'short' });
    }
    
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
};

const getInitials = (name: string) => {
    if (!name) return '?';
    return name
        .split(' ')
        .map(word => word[0])
        .join('')
        .toUpperCase()
        .slice(0, 2);
};

export const ConversationItem = memo<ConversationItemProps>(({ conversation, onPress }) => {
    const { user } = useAuth();

    const { displayName, displayImage } = useMemo(() => {
        const fallbackName = conversation.coach_name || conversation.client_name || 'Unknown';
        const fallbackImage = conversation.coach_image ?? conversation.client_image ?? null;

        if (!user) {
            return { displayName: fallbackName, displayImage: fallbackImage };
        }

        if (user.id === conversation.coach_id) {
            return {
                displayName: conversation.client_name || fallbackName,
                displayImage: conversation.client_image ?? fallbackImage,
            };
        }

        if (user.id === conversation.client_id) {
            return {
                displayName: conversation.coach_name || fallbackName,
                displayImage: conversation.coach_image ?? fallbackImage,
            };
        }

        if (user.role === 'coach') {
            return {
                displayName: conversation.client_name || fallbackName,
                displayImage: conversation.client_image ?? fallbackImage,
            };
        }

        return {
            displayName: conversation.coach_name || fallbackName,
            displayImage: conversation.coach_image ?? fallbackImage,
        };
    }, [conversation, user]);

    const formattedTime = formatTime(conversation.last_message_at);
    const isLastMessageFromCurrentUser = Boolean(
        user?.id && conversation.last_message_sender_id === user.id
    );
    const messagePreview = conversation.last_message_text
        ? `${isLastMessageFromCurrentUser ? 'You: ' : ''}${conversation.last_message_text}`
        : null;
    const hasUnread = !isLastMessageFromCurrentUser && conversation.total_messages > 0;
    
    return (
        <TouchableOpacity 
            style={styles.container} 
            onPress={onPress}
            activeOpacity={0.7}
        >
            <View style={styles.content}>
                {displayImage ? (
                    <Image 
                        source={{ uri: displayImage }} 
                        style={styles.avatar}
                    />
                ) : (
                    <View style={styles.avatarPlaceholder}>
                        <Text style={styles.avatarText}>{getInitials(displayName)}</Text>
                    </View>
                )}
                
                <View style={styles.main}>
                    <View style={styles.header}>
                        <Text style={styles.name} numberOfLines={1}>
                            {displayName}
                        </Text>
                        {formattedTime ? (
                            <Text style={[styles.time, hasUnread && styles.timeUnread]}>
                                {formattedTime}
                            </Text>
                        ) : null}
                    </View>
                    
                    <View style={styles.footer}>
                        {messagePreview ? (
                            <Text 
                                style={[styles.message, hasUnread && styles.messageUnread]} 
                                numberOfLines={1}
                            >
                                {messagePreview}
                            </Text>
                        ) : (
                            <Text style={styles.messagePlaceholder}>No messages yet</Text>
                        )}
                        
                        {hasUnread ? (
                            <View style={styles.badge}>
                                <Text style={styles.badgeText}>
                                    {conversation.total_messages > 99 ? '99+' : Math.max(conversation.total_messages, 1)}
                                </Text>
                            </View>
                        ) : null}
                    </View>
                </View>
            </View>
            
            <View style={styles.separator} />
        </TouchableOpacity>
    );
});

const styles = StyleSheet.create({
    container: {
        backgroundColor: COLORS.background.auth,
    },
    content: {
        flexDirection: 'row',
        padding: SPACING.base,
        alignItems: 'center',
    },
    avatar: {
        width: 56,
        height: 56,
        borderRadius: 28,
        backgroundColor: COLORS.background.card,
    },
    avatarPlaceholder: {
        width: 56,
        height: 56,
        borderRadius: 28,
        backgroundColor: COLORS.background.card,
        justifyContent: 'center',
        alignItems: 'center',
    },
    avatarText: {
        fontSize: FONT_SIZES.lg,
        fontWeight: FONT_WEIGHTS.bold,
        color: COLORS.primary,
    },
    main: {
        flex: 1,
        marginLeft: SPACING.md,
    },
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: SPACING.xs,
    },
    name: {
        flex: 1,
        fontSize: FONT_SIZES.base,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.auth.primary,
        marginRight: SPACING.sm,
    },
    time: {
        fontSize: FONT_SIZES.xs,
        color: COLORS.text.tertiary,
    },
    timeUnread: {
        color: COLORS.primary,
        fontWeight: FONT_WEIGHTS.semibold,
    },
    footer: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    message: {
        flex: 1,
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.tertiary,
        marginRight: SPACING.sm,
    },
    messageUnread: {
        color: COLORS.text.auth.secondary,
        fontWeight: FONT_WEIGHTS.medium,
    },
    messagePlaceholder: {
        flex: 1,
        fontSize: FONT_SIZES.sm,
        color: COLORS.text.placeholder,
        fontStyle: 'italic',
    },
    badge: {
        minWidth: 20,
        height: 20,
        borderRadius: 10,
        backgroundColor: COLORS.primary,
        justifyContent: 'center',
        alignItems: 'center',
        paddingHorizontal: SPACING.xs,
    },
    badgeText: {
        fontSize: FONT_SIZES.xs - 1,
        fontWeight: FONT_WEIGHTS.bold,
        color: COLORS.text.primary,
    },
    separator: {
        height: 1,
        backgroundColor: COLORS.border.dark,
        marginLeft: 56 + SPACING.base + SPACING.md,
    },
});

ConversationItem.displayName = 'ConversationItem';
