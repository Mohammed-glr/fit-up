import React, { useEffect, useState, useCallback } from 'react';
import { View, StyleSheet, KeyboardAvoidingView, Platform, Alert } from 'react-native';
import { MessageList } from './message-list';
import { MessageComposer } from './message-composer';
import { WorkoutAttachmentPicker } from './workout-attachment-picker';
import { httpClient } from '@/api/client';
import { useToastMethods } from '@/components/ui';
import { useDeleteMessage, useMarkAllAsRead, useSendMessage, useUpdateMessage } from '@/hooks/message/use-conversation';
import { COLORS } from '@/constants/theme';
import type { MessageWithDetails } from '@/types/message';
import type { WorkoutShareSummary } from '@/types/workout-sharing';

interface ChatViewProps {
    conversationId: number;
}

import { useChatWebSocket } from '@/hooks/message/use-chat-websocket';

export const ChatView: React.FC<ChatViewProps> = ({ conversationId }) => {
    useChatWebSocket(conversationId);
    const [messageText, setMessageText] = useState('');
    const [showWorkoutPicker, setShowWorkoutPicker] = useState(false);

    const { mutateAsync: sendMessageAsync, isPending: isSending } = useSendMessage();
    const { mutateAsync: updateMessageAsync, isPending: isUpdating } = useUpdateMessage();
    const { mutateAsync: deleteMessageAsync, isPending: isDeleting } = useDeleteMessage();
    const { mutate: markConversationAsRead } = useMarkAllAsRead();
    const { showError } = useToastMethods();
    const [editingMessage, setEditingMessage] = useState<MessageWithDetails | null>(null);

    // Mock recent workouts - TODO: Replace with actual API call
    const recentWorkouts = [
        {
            session_id: 1,
            workout_title: 'Push Day - Chest & Triceps',
            completed_at: new Date().toISOString(),
            duration_minutes: 65,
            total_exercises: 6,
            total_volume_lbs: 4250,
        },
        {
            session_id: 2,
            workout_title: 'Pull Day - Back & Biceps',
            completed_at: new Date(Date.now() - 86400000).toISOString(),
            duration_minutes: 58,
            total_exercises: 5,
            total_volume_lbs: 3890,
        },
    ];

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

    const handleAttachWorkout = useCallback(() => {
        setShowWorkoutPicker(true);
    }, []);

    const handleSelectWorkout = useCallback(async (sessionId: number) => {
        try {
            setShowWorkoutPicker(false);

            console.log('Fetching workout summary for session:', sessionId);

            // Fetch the workout share summary
            const response = await httpClient.get<WorkoutShareSummary>(
                `/workout-sessions/${sessionId}/share-summary`
            );

            console.log('Workout summary response:', response.data);

            const summary = response.data;

            // Format the workout summary text
            const formattedText = formatWorkoutSummary(summary);

            // Append to message text
            setMessageText(prev => {
                const separator = prev.trim() ? '\n\n' : '';
                return prev + separator + formattedText;
            });
        } catch (error: any) {
            console.error('Error attaching workout:', error);
            console.error('Error response:', error.response?.data);
            console.error('Error status:', error.response?.status);

            const errorMessage = error.response?.data?.error ||
                error.response?.data?.message ||
                error.message ||
                'Failed to attach workout. Please try again.';

            Alert.alert('Error', errorMessage);
        }
    }, []);

    const formatWorkoutSummary = (summary: WorkoutShareSummary) => {
        const lines = [
            `🏋️ ${summary.workout_title}`,
            `📅 ${new Date(summary.completed_at).toLocaleDateString()}`,
            `⏱️ Duration: ${summary.duration_minutes} minutes`,
            `💪 ${summary.total_exercises} exercises • ${summary.total_sets} sets`,
            `📊 Total Volume: ${summary.total_volume_lbs.toFixed(0)} lbs`,
        ];

        if (summary.prs_achieved > 0) {
            lines.push(`🏆 ${summary.prs_achieved} Personal Records!`);
        }

        return lines.join('\n');
    };

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
                    onAttachWorkout={handleAttachWorkout}
                />

                <WorkoutAttachmentPicker
                    visible={showWorkoutPicker}
                    onClose={() => setShowWorkoutPicker(false)}
                    onSelectWorkout={handleSelectWorkout}
                    recentWorkouts={recentWorkouts}
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
