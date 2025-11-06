import React, { useCallback, useEffect, useMemo, useRef } from 'react';
import { ActivityIndicator, Alert, FlatList, Pressable, RefreshControl, StyleSheet, Text, View } from 'react-native';
import { useConversationMessages } from '@/hooks/message/use-conversation';
import type { MessageWithDetails } from '@/types/message';
import { MessageBubble } from './message-bubble';
import { useAuth } from '@/context/auth-context';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';


type MessageListProps = {
    conversationId: number;
    onRequestEdit?: (message: MessageWithDetails) => void;
    onRequestDelete?: (message: MessageWithDetails) => void;
};



type FlatListType = FlatList<MessageWithDetails>;

export const MessageList: React.FC<MessageListProps> = ({ conversationId, onRequestEdit, onRequestDelete }) => {
    const { user } = useAuth();
    const listRef = useRef<FlatListType | null>(null);
    const { data, isLoading, isError, refetch, fetchNextPage, hasNextPage, isFetchingNextPage } = useConversationMessages(conversationId);

    const messages = useMemo(() => {
        if (!data?.pages) {
            return [] as MessageWithDetails[];
        }
        const allMessages = data.pages
            .flatMap((page) => page.messages ?? [])
            .filter((message): message is MessageWithDetails => Boolean(message && message.message_id != null));
        
        return allMessages.reverse();
    }, [data]);

    const newestMessageId = messages.length > 0 ? messages[messages.length - 1]?.message_id : null;
    const previousNewestMessageId = useRef<number | null>(null);

    useEffect(() => {
        if (newestMessageId === null) {
            previousNewestMessageId.current = null;
            return;
        }
        if (previousNewestMessageId.current === null) {
            previousNewestMessageId.current = newestMessageId;
            setTimeout(() => {
                listRef.current?.scrollToEnd({ animated: false });
            }, 100);
            return;
        }
        if (previousNewestMessageId.current !== newestMessageId) {
            listRef.current?.scrollToEnd({ animated: true });
            previousNewestMessageId.current = newestMessageId;
        }
    }, [newestMessageId]);

    const handleLoadMore = useCallback(() => {
        if (hasNextPage && !isFetchingNextPage) {
            fetchNextPage();
        }
    }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

    const handleLongPress = useCallback(
        (message: MessageWithDetails) => {
            if (!user?.id || message.sender_id !== user.id) {
                return;
            }

            Alert.alert('Message options', undefined, [
                {
                    text: 'Edit message',
                    onPress: () => onRequestEdit?.(message),
                },
                {
                    text: 'Delete message',
                    style: 'destructive',
                    onPress: () => onRequestDelete?.(message),
                },
                {
                    text: 'Cancel',
                    style: 'cancel',
                },
            ]);
        },
        [onRequestDelete, onRequestEdit, user?.id],
    );

    const renderMessage = useCallback(
        ({ item }: { item: MessageWithDetails }) => {
            const isOwnMessage = item.sender_id === user?.id;
            return (
                <Pressable
                    style={styles.messageWrapper}
                    onLongPress={() => handleLongPress(item)}
                    delayLongPress={250}
                    disabled={!isOwnMessage}
                >
                    <MessageBubble message={item} isOwnMessage={isOwnMessage} />
                </Pressable>
            );
        },
        [handleLongPress, user?.id],
    );

    const keyExtractor = useCallback((item: MessageWithDetails) => item.message_id.toString(), []);

    if (isLoading && messages.length === 0) {
        return (
            <View style={styles.centered}>
                <ActivityIndicator size="large" color={COLORS.primary} />
            </View>
        );
    }

    if (isError) {
        return (
            <View style={styles.centered}>
                <Text style={styles.errorText}>Failed to load messages</Text>
                <Text style={styles.errorSubtext}>Pull down to retry</Text>
            </View>
        );
    }



    return (
        <FlatList
            ref={listRef}
            data={messages}
            renderItem={renderMessage}
            keyExtractor={keyExtractor}
            contentContainerStyle={messages.length === 0 ? styles.emptyContent : styles.contentContainer}
            style={styles.list}

            onEndReached={handleLoadMore}
            onEndReachedThreshold={0.1}
            refreshControl={(
                <RefreshControl
                    refreshing={isLoading}
                    onRefresh={refetch}
                    tintColor={COLORS.primary}
                    colors={[COLORS.primary]}
                />
                
            )}
            ListEmptyComponent={
                !isLoading ? (
                    <View style={styles.emptyState}>
                        <Text style={styles.emptyText}>No messages yet</Text>
                        <Text style={styles.emptySubtext}>Start the conversation!</Text>
                    </View>
                ) : null
            }
            ListFooterComponent={
                isFetchingNextPage ? (
                    <ActivityIndicator style={styles.footerSpinner} color={COLORS.primary} />
                ) : null
            }
        />
    );
};

const styles = StyleSheet.create({
    list: {
        flex: 1,
        backgroundColor: COLORS.background.auth,
    },
    messageWrapper: {
        width: '100%',
    },
    contentContainer: {
        paddingVertical: SPACING.md,
    },
    centered: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: COLORS.background.auth,
        padding: SPACING.xl,
    },
    errorText: {
        color: COLORS.error,
        fontSize: FONT_SIZES.lg,
        fontWeight: FONT_WEIGHTS.semibold,
        textAlign: 'center',
        marginBottom: SPACING.xs,
    },
    errorSubtext: {
        color: COLORS.text.tertiary,
        fontSize: FONT_SIZES.sm,
        textAlign: 'center',
    },
    emptyContent: {
        flexGrow: 1,
        justifyContent: 'center',
        alignItems: 'center',
        paddingVertical: SPACING['3xl'],
    },
    emptyState: {
        alignItems: 'center',
    },
    emptyText: {
        color: COLORS.text.secondary,
        fontSize: FONT_SIZES.lg,
        fontWeight: FONT_WEIGHTS.semibold,
        marginBottom: SPACING.xs,
    },
    emptySubtext: {
        color: COLORS.text.tertiary,
        fontSize: FONT_SIZES.sm,
    },
    footerSpinner: {
        paddingVertical: SPACING.base,
    },
    senderInfo: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: SPACING.xs,
        marginLeft: SPACING.base,
    },
    senderAvatar: {
        width: 24,
        height: 24,
        borderRadius: 12,
        marginRight: SPACING.xs,
    },
    senderName: {
        fontSize: FONT_SIZES.sm,
        fontWeight: FONT_WEIGHTS.semibold,
        color: COLORS.text.placeholder,
    }
});
