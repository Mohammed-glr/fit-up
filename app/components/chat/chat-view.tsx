import React, { useEffect, useState, useCallback } from 'react';
import { View, StyleSheet, KeyboardAvoidingView, Platform } from 'react-native';
import { MessageList } from './message-list';
import { MessageComposer } from './message-composer';
import { useMarkAllAsRead, useSendMessage } from '@/hooks/message/use-conversation';

interface ChatViewProps {
    conversationId: number;
}

export const ChatView: React.FC<ChatViewProps> = ({ conversationId }) => {
    const [messageText, setMessageText] = useState('');
    const { mutateAsync: sendMessageAsync, isPending: isSending } = useSendMessage();
    const { mutate: markConversationAsRead } = useMarkAllAsRead();

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
            await sendMessageAsync({
                conversation_id: conversationId,
                message_text: trimmed,
            });
            setMessageText('');
        } catch (error) {
        }
    }, [conversationId, messageText, sendMessageAsync]);

    return (
        <KeyboardAvoidingView
            style={styles.keyboardAvoiding}
            behavior={Platform.OS === 'ios' ? 'padding' : undefined}
            keyboardVerticalOffset={Platform.OS === 'ios' ? 64 : 0}
        >
            <View style={styles.container}>
                <MessageList conversationId={conversationId} />
                <MessageComposer
                    value={messageText}
                    onChangeText={setMessageText}
                    onSend={handleSend}
                    isSending={isSending}
                />
            </View>
        </KeyboardAvoidingView>
    );
};

const styles = StyleSheet.create({
    keyboardAvoiding: {
        flex: 1,
        backgroundColor: '#030712',
    },
    container: {
        flex: 1,
    },
});
