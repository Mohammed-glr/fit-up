import React, { useEffect, useState, useCallback } from 'react';
import { View, StyleSheet, KeyboardAvoidingView, Platform, Alert } from 'react-native';
import { MessageList } from './message-list';
import { MessageComposer } from './message-composer';
import { useToastMethods } from '@/components/ui';
import { useDeleteMessage, useMarkAllAsRead, useSendMessage, useUpdateMessage } from '@/hooks/message/use-conversation';
import { COLORS } from '@/constants/theme';
import type { MessageWithDetails } from '@/types/message';

interface ChatViewProps {
    conversationId: number;
}

export const ChatView: React.FC<ChatViewProps> = ({ conversationId }) => {
    const [messageText, setMessageText] = useState('');
    const { mutateAsync: sendMessageAsync, isPending: isSending } = useSendMessage();
    const { mutateAsync: updateMessageAsync, isPending: isUpdating } = useUpdateMessage();
    const { mutateAsync: deleteMessageAsync, isPending: isDeleting } = useDeleteMessage();
    const { mutate: markConversationAsRead } = useMarkAllAsRead();
    const { showError } = useToastMethods();
    const [editingMessage, setEditingMessage] = useState<MessageWithDetails | null>(null);

    const resetComposerState = useCallback(() => {
        setMessageText('');
        setEditingMessage(null);
    }, []);

    useEffect(() => {
        if (!conversationId) {
            return;
        }
        markConversationAsRead(conversationId);
    }, [conversationId, markConversationAsRead]);

    const handleSend = useCallback(async () => {
        const trimmed = messageText.trim();
        if (!trimmed) {
            return;
        }

        try {
            if (editingMessage) {
                await updateMessageAsync({
                    message_id: editingMessage.message_id,
                    conversation_id: conversationId,
                    data: { message_text: trimmed },
                });
            } else {
                await sendMessageAsync({
                    conversation_id: conversationId,
                    message_text: trimmed,
                });
            }
            resetComposerState();
        } catch (error) {
            console.error('Failed to submit message:', error);
            showError('Failed to submit message. Please try again.');
        }
    }, [conversationId, editingMessage, messageText, resetComposerState, sendMessageAsync, showError, updateMessageAsync]);

    const handleEditRequest = useCallback((message: MessageWithDetails) => {
        setEditingMessage(message);
        setMessageText(message.message_text ?? '');
    }, []);

    const handleDeleteRequest = useCallback((message: MessageWithDetails) => {
        Alert.alert('Delete message?', 'This action will remove the message for everyone in the chat.', [
            { text: 'Cancel', style: 'cancel' },
            {
                text: 'Delete',
                style: 'destructive',
                onPress: async () => {
                    try {
                        await deleteMessageAsync({
                            message_id: message.message_id,
                            conversation_id: conversationId,
                        });
                        if (editingMessage?.message_id === message.message_id) {
                            resetComposerState();
                        }
                    } catch (error) {
                        console.error('Failed to delete message:', error);
                        showError('Failed to delete message. Please try again.');
                    }
                },
            },
        ]);
    }, [conversationId, deleteMessageAsync, editingMessage?.message_id, resetComposerState, showError]);

    const handleCancelEdit = useCallback(() => {
        resetComposerState();
    }, [resetComposerState]);

    const isComposerBusy = isSending || isUpdating || isDeleting;

    return (
        <KeyboardAvoidingView
            style={styles.keyboardAvoiding}
            behavior={Platform.OS === 'ios' ? 'padding' : undefined}
            keyboardVerticalOffset={Platform.OS === 'ios' ? 64 : 0}
        >
            <View style={styles.container}>
                <MessageList
                    conversationId={conversationId}
                    onRequestEdit={handleEditRequest}
                    onRequestDelete={handleDeleteRequest}
                />
                <MessageComposer
                    value={messageText}
                    onChangeText={setMessageText}
                    onSend={handleSend}
                    isSending={isComposerBusy}
                    isEditing={Boolean(editingMessage)}
                    onCancelEdit={handleCancelEdit}
                />
            </View>
        </KeyboardAvoidingView>
    );
};

const styles = StyleSheet.create({
    keyboardAvoiding: {
        flex: 1,
        backgroundColor: COLORS.background.auth,
        paddingTop: 10,

    },
    container: {
        flex: 1,
    },
});
