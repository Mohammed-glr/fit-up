import React from 'react';
import { View, TextInput, TouchableOpacity, Text, StyleSheet } from 'react-native';

interface MessageComposerProps {
    value: string;
    onChangeText: (text: string) => void;
    onSend: () => void;
    isSending?: boolean;
}

export const MessageComposer: React.FC<MessageComposerProps> = ({ value, onChangeText, onSend, isSending }) => {
    const isDisabled = isSending || value.trim().length === 0;

    return (
        <View style={styles.container}>
            <TextInput
                style={styles.input}
                value={value}
                onChangeText={onChangeText}
                placeholder="Type a message..."
                placeholderTextColor="#9CA3AF"
                multiline
                editable={!isSending}
            />
            <TouchableOpacity
                style={[styles.sendButton, isDisabled && styles.sendButtonDisabled]}
                onPress={onSend}
                disabled={isDisabled}
                accessibilityLabel="Send message"
            >
                <Text style={styles.sendLabel}>{isSending ? 'Sending…' : 'Send'}</Text>
            </TouchableOpacity>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'flex-end',
        paddingHorizontal: 12,
        paddingVertical: 16,
        borderTopWidth: 1,
        borderTopColor: '#1F2937',
        backgroundColor: '#0B1120',
    },
    input: {
        flex: 1,
        minHeight: 40,
        maxHeight: 120,
        paddingHorizontal: 16,
        paddingVertical: 12,
        borderRadius: 20,
        borderWidth: 1,
        borderColor: '#374151',
        backgroundColor: '#111827',
        color: '#F9FAFB',
        fontSize: 16,
    },
    sendButton: {
        paddingHorizontal: 18,
        paddingVertical: 12,
        borderRadius: 20,
        backgroundColor: '#2563EB',
        marginLeft: 12,
    },
    sendButtonDisabled: {
        backgroundColor: '#1E3A8A',
    },
    sendLabel: {
        color: '#F9FAFB',
        fontWeight: '600',
    },
});
