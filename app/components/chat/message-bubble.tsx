import React, { memo, useMemo } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import type { MessageWithDetails } from '@/types/message';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';

type MessageBubbleProps = {
    message: MessageWithDetails;
    isOwnMessage: boolean;
};

const formatTime = (isoDate: string) => {
    const date = new Date(isoDate);
    if (Number.isNaN(date.getTime())) {
        return '';
    }
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

export const MessageBubble = memo<MessageBubbleProps>(({ message, isOwnMessage }) => {
    const sentAt = useMemo(() => formatTime(message.sent_at), [message.sent_at]);
    const isDeleted = message.is_deleted || message.message_text === '[Message deleted]';

    return (
        <View style={[styles.container, isOwnMessage ? styles.containerOwn : styles.containerOther]}>
            <View style={[styles.bubble, isOwnMessage ? styles.bubbleOwn : styles.bubbleOther]}>
                {!isOwnMessage && message.sender_name ? (
                    <Text style={styles.sender}>{message.sender_name}</Text>
                ) : null}
                <Text
                    style={[
                        styles.messageText,
                        isDeleted && styles.deletedText,
                    ]}
                >
                    {message.message_text}
                </Text>
                {message.attachments && message.attachments.length > 0 ? (
                    <View style={styles.attachments}>
                        {message.attachments.map((attachment, index) => (
                            <Text
                                key={attachment.attachment_id}
                                style={[styles.attachmentLabel, index > 0 && styles.attachmentSpacing]}
                            >
                                {attachment.file_name}
                            </Text>
                        ))}
                    </View>
                ) : null}
            </View>
                <View style={styles.meta}>
                    <Text style={styles.timestamp}>{sentAt}</Text>
                    {isOwnMessage ? (
                        <Text style={[styles.readReceipt, message.is_read && styles.readReceiptRead]}>
                            {message.is_read ? 'read' : 'delivered'}
                        </Text>
                    ) : null}
                </View>
        </View>
    );
});

const styles = StyleSheet.create({
    container: {
        width: '100%',
        paddingHorizontal: SPACING.base,
        marginBottom: SPACING.md,
    },
    containerOwn: {
        alignItems: 'flex-end',
    },
    containerOther: {
        alignItems: 'flex-start',
    },
    bubble: {
        maxWidth: '80%',
        borderRadius: BORDER_RADIUS.lg,
        paddingVertical: SPACING.sm + 2,
        paddingHorizontal: SPACING.base - 2,
    },
    bubbleOwn: {
        backgroundColor: COLORS.primary,
    },
    bubbleOther: {
        backgroundColor: COLORS.background.card,
    },
    sender: {
        fontSize: FONT_SIZES.xs,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.auth.secondary,
        marginBottom: SPACING.xs,
    },
    messageText: {
        fontSize: FONT_SIZES.base,
        lineHeight: 22,
        color: COLORS.text.auth.primary,
    },
    deletedText: {
        fontStyle: 'italic',
        color: COLORS.text.tertiary,
    },
    attachments: {
        marginTop: SPACING.sm,
    },
    attachmentLabel: {
        fontSize: FONT_SIZES.xs,
        color: COLORS.primary,
        textDecorationLine: 'underline',
    },
    attachmentSpacing: {
        marginTop: SPACING.xs,
    },
    meta: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginTop: SPACING.sm,
    },
    timestamp: {
        fontSize: FONT_SIZES.xs - 1,
        color: COLORS.text.tertiary,
    },
    readReceipt: {
        fontSize: FONT_SIZES.xs,
        color: COLORS.text.tertiary,
        marginLeft: SPACING.xs,
    },
    readReceiptRead: {
        color: COLORS.primary,
        fontWeight: FONT_WEIGHTS.semibold,
    },
});

MessageBubble.displayName = 'MessageBubble';
