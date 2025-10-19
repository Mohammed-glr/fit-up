import React from 'react';
import { View, TextInput, TouchableOpacity, Text, StyleSheet } from 'react-native';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import {
    Button
} from '@/components/forms/button';

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
                placeholderTextColor={COLORS.text.placeholder}
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
        // position: 'absolute',
        // bottom: 0,
        // left: 0,
        // right: 0,
        flexDirection: 'row',
        alignItems: 'flex-end',
        paddingHorizontal: SPACING.md,
        paddingVertical: SPACING.sm,    
        marginBottom: 30,
        borderTopWidth: 1,
        borderTopColor: COLORS.border.dark,
        backgroundColor: COLORS.background.card,
        margin: 10,
        borderRadius: BORDER_RADIUS['3xl']
        
    },
    input: {
        flex: 1,
        minHeight: 40,
        maxHeight: 120,
        paddingHorizontal: SPACING.base,
        paddingVertical: SPACING.md,
        borderRadius: BORDER_RADIUS.lg,
        borderWidth: 1,
        borderColor: COLORS.border.dark,
        backgroundColor: COLORS.background.dark,
        color: COLORS.text.auth.primary,
        fontSize: FONT_SIZES.base,
    },
    sendButton: {
        paddingHorizontal: SPACING.lg,
        paddingVertical: SPACING.md,
        borderRadius: BORDER_RADIUS.lg,
        backgroundColor: COLORS.primary,
        marginLeft: SPACING.md,
    },
    sendButtonDisabled: {
        opacity: 0.5,
    },
    sendLabel: {
        color: COLORS.text.primary,
        fontWeight: FONT_WEIGHTS.semibold,
        fontSize: FONT_SIZES.base,
    },
});
