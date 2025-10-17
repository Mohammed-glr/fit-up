import React, { useCallback, useEffect, useMemo, useRef } from 'react';
import { ActivityIndicator, FlatList, RefreshControl, StyleSheet, Text, View } from 'react-native';
import { useConversationMessages } from '@/hooks/message/use-conversation';
import type { MessageWithDetails } from '@/types/message';
import { MessageBubble } from './message-bubble';
import { useAuth } from '@/context/auth-context';

type MessageListProps = {
    conversationId: number;
};

type FlatListType = FlatList<MessageWithDetails>;

export const MessageList: React.FC<MessageListProps> = ({ conversationId }) => {
    const { user } = useAuth();
    const listRef = useRef<FlatListType | null>(null);
    const { data, isLoading, isError, refetch, fetchNextPage, hasNextPage, isFetchingNextPage } = useConversationMessages(conversationId);

    const messages = useMemo(() => {
        if (!data?.pages) {
            return [] as MessageWithDetails[];
        }
        return data.pages.flatMap((page) => page.messages);
    }, [data]);

    const newestMessageId = messages.length > 0 ? messages[0].message_id : null;
    const previousNewestMessageId = useRef<number | null>(null);

    useEffect(() => {
        if (newestMessageId === null) {
            previousNewestMessageId.current = null;
            return;
        }
        if (previousNewestMessageId.current === null) {
            previousNewestMessageId.current = newestMessageId;
            return;
        }
        if (previousNewestMessageId.current !== newestMessageId) {
            listRef.current?.scrollToOffset({ offset: 0, animated: true });
            previousNewestMessageId.current = newestMessageId;
        }
    }, [newestMessageId]);

    const handleLoadMore = useCallback(() => {
        if (hasNextPage && !isFetchingNextPage) {
            fetchNextPage();
        }
    }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

    const renderMessage = useCallback(
        ({ item }: { item: MessageWithDetails }) => (
            <MessageBubble message={item} isOwnMessage={item.sender_id === user?.id} />
        ),
        [user?.id],
    );

    const keyExtractor = useCallback((item: MessageWithDetails) => item.message_id.toString(), []);

    if (isLoading && messages.length === 0) {
        return (
            <View style={styles.centered}>
                <ActivityIndicator size="large" color="#2563EB" />
            </View>
        );
    }

    if (isError) {
        return (
            <View style={styles.centered}>
                <Text style={styles.errorText}>We could not load the conversation. Pull to retry.</Text>
            </View>
        );
    }

    return (
        <FlatList
            ref={listRef}
            data={messages}
            renderItem={renderMessage}
            keyExtractor={keyExtractor}
            inverted
            contentContainerStyle={messages.length === 0 ? styles.emptyContent : undefined}
            style={styles.list}
            onEndReached={handleLoadMore}
            onEndReachedThreshold={0.1}
            refreshControl={(
                <RefreshControl
                    refreshing={isLoading}
                    onRefresh={refetch}
                    tintColor="#2563EB"
                    colors={["#2563EB"]}
                />
            )}
            ListEmptyComponent={!isLoading ? <Text style={styles.emptyText}>No messages yet. Say hello!</Text> : null}
            ListFooterComponent={isFetchingNextPage ? <ActivityIndicator style={styles.footerSpinner} color="#2563EB" /> : null}
        />
    );
};

const styles = StyleSheet.create({
    list: {
        flex: 1,
        backgroundColor: '#030712',
    },
    centered: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#030712',
        padding: 24,
    },
    errorText: {
        color: '#FCA5A5',
        textAlign: 'center',
        fontSize: 15,
    },
    emptyContent: {
        flexGrow: 1,
        justifyContent: 'center',
        alignItems: 'center',
        paddingVertical: 32,
    },
    emptyText: {
        color: '#9CA3AF',
        fontSize: 15,
        textAlign: 'center',
    },
    footerSpinner: {
        paddingVertical: 12,
    },
});
