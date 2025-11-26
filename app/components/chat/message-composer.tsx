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
    isEditing?: boolean;
    onCancelEdit?: () => void;
    onAttachWorkout?: () => void; // New prop for workout attachment
}

export const MessageComposer: React.FC<MessageComposerProps> = ({ 
    value, 
    onChangeText, 
    onSend, 
    isSending, 
    isEditing, 
    onCancelEdit,
    onAttachWorkout 
}) => {
    const isDisabled = isSending || value.trim().length === 0;

    return (
        <View style={[styles.container, isEditing && styles.containerEditing]}>
            {isEditing ? (
                <View style={styles.editBanner}>
                    <Text style={styles.editBannerText}>Editing message</Text>
                    <TouchableOpacity
                        onPress={onCancelEdit}
                        disabled={!onCancelEdit}
                        accessibilityLabel="Cancel editing"
                        hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                    >
                        <Text style={styles.cancelEditText}>Cancel</Text>
                    </TouchableOpacity>
                </View>
            ) : null}
            
            {/* {onAttachWorkout && !isEditing && (
                <TouchableOpacity
                    style={styles.attachButton}
                    onPress={onAttachWorkout}
                    disabled={isSending}
                    accessibilityLabel="Share workout"
                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                    <Ionicons
                        name="barbell-outline"
                        size={24}
                        color={COLORS.white}
                    />
                </TouchableOpacity>
            )} */}
            
            <TextInput
                style={styles.input}
                value={value}
                onChangeText={onChangeText}
                placeholder={isEditing ? 'Edit your message...' : 'Type a message...'}
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
        margin: 5,
        marginBottom: 50,
        borderRadius: BORDER_RADIUS['3xl'],
        borderColor:'rgba(28, 28, 30, 0.95)',
        borderWidth: 1,
        position: 'relative',
    },
    containerEditing: {
        paddingTop: SPACING.lg + SPACING.xs,
    },
    editBanner: {
        position: 'absolute',
        top: -SPACING.lg,
        left: SPACING.md,
        right: SPACING.md,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    editBannerText: {
        color: COLORS.warning,
        fontWeight: FONT_WEIGHTS.semibold,
        fontSize: FONT_SIZES.sm,
    },
    cancelEditText: {
        color: COLORS.text.auth.secondary,
        fontWeight: FONT_WEIGHTS.semibold,
        fontSize: FONT_SIZES.sm,
    },
    attachButton: {
        paddingHorizontal: SPACING.md,
        paddingVertical: SPACING.md,
        backgroundColor: COLORS.background.accent,
        borderRadius: BORDER_RADIUS.full,
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: SPACING.xs,
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
