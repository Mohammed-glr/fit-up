import React from 'react';
import { View, TextInput, TouchableOpacity, Text, StyleSheet } from 'react-native';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
import {
    Ionicons
} from '@expo/vector-icons'

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
                <Ionicons
                    name={isSending ? 'arrow-up-circle' : 'arrow-up-circle'}
                    size={24}
                    color={COLORS.text.inverse}
                />
            </TouchableOpacity>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'flex-end',
        paddingHorizontal: SPACING.md,
        paddingVertical: SPACING.sm,            
        backgroundColor: COLORS.background.card,
        margin: 20,
        borderRadius: BORDER_RADIUS['3xl'],
        borderColor:'rgba(28, 28, 30, 0.95)',
        borderWidth: 1,
        
    },
    input: {
        flex: 1,
        minHeight: 40,
        maxHeight: 120,
        paddingHorizontal: SPACING.md,
        paddingVertical: SPACING.sm,
        borderRadius: BORDER_RADIUS.lg,
        backgroundColor: COLORS.background.dark,
        color: COLORS.text.auth.primary,
        fontSize: FONT_SIZES.base,
    },
    sendButton: {
        paddingHorizontal: SPACING.md,
        paddingVertical: SPACING.md,
        borderRadius: BORDER_RADIUS.full,
        backgroundColor: COLORS.background.accent,
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
