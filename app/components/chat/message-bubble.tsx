import React, { memo, useMemo } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import type { MessageWithDetails } from '@/types/message';

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
                <View style={styles.meta}>
                    <Text style={styles.timestamp}>{sentAt}</Text>
                    {isOwnMessage ? (
                        <Text style={[styles.readReceipt, message.is_read && styles.readReceiptRead]}>
                            {message.is_read ? 'Read' : 'Sent'}
                        </Text>
                    ) : null}
                </View>
            </View>
        </View>
    );
});

const styles = StyleSheet.create({
    container: {
        width: '100%',
        paddingHorizontal: 16,
        marginBottom: 12,
    },
    containerOwn: {
        alignItems: 'flex-end',
    },
    containerOther: {
        alignItems: 'flex-start',
    },
    bubble: {
        maxWidth: '80%',
        borderRadius: 16,
        paddingVertical: 10,
        paddingHorizontal: 14,
    },
    bubbleOwn: {
        backgroundColor: '#2563EB',
    },
    bubbleOther: {
        backgroundColor: '#1F2937',
    },
    sender: {
        fontSize: 12,
        fontWeight: '600',
        color: '#F3F4F6',
        marginBottom: 4,
    },
    messageText: {
        fontSize: 15,
        lineHeight: 20,
        color: '#F9FAFB',
    },
    deletedText: {
        fontStyle: 'italic',
        color: '#9CA3AF',
    },
    attachments: {
        marginTop: 8,
    },
    attachmentLabel: {
        fontSize: 12,
        color: '#93C5FD',
        textDecorationLine: 'underline',
    },
    attachmentSpacing: {
        marginTop: 4,
    },
    meta: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginTop: 8,
    },
    timestamp: {
        fontSize: 11,
        color: '#CBD5F5',
    },
    readReceipt: {
        fontSize: 11,
        color: '#BFDBFE',
    },
    readReceiptRead: {
        fontWeight: '600',
    },
});

MessageBubble.displayName = 'MessageBubble';
