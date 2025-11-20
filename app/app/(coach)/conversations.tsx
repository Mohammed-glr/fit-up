import React from "react";
import { View, Text, FlatList, ActivityIndicator, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useConversations } from "@/hooks/message/use-conversation";
import type { ConversationOverview } from "@/types";
import { ConversationItem } from "../../components/chat/conversation-item";
import { CreateConversationFAB } from "../../components/chat/createConversationFAB";
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';

export default function ConversationsScreen() {
    const router = useRouter();
    const {
        data,
        isLoading,
        isError,
        refetch,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
    } = useConversations();

    const conversations = (
        data?.pages.flatMap((page) => page.conversations ?? []) ?? []
    ).filter((conversation): conversation is ConversationOverview => conversation != null);

    const renderItem = ({ item }: { item: ConversationOverview }) => {
        return (
            <ConversationItem
                conversation={item}
                onPress={() =>
                    router.push({
                        pathname: '/(coach)/chat',
                        params: { conversationId: String(item.conversation_id) },
                    })
                }
            />
        );
    };

    const renderEmptyState = () => {
        if (isLoading) {
            return (
                <View style={styles.centerContainer}>
                    <ActivityIndicator size="large" color={COLORS.primary} />
                </View>
            );
        }

        return (
            <View style={styles.emptyState}>
                <Text style={styles.emptyTitle}>No conversations yet</Text>
                <Text style={styles.emptySubtitle}>
                    {isError
                        ? 'Failed to load conversations. Pull to retry.'
                        : 'Tap the + button to start a new conversation'}
                </Text>
            </View>
        );
    };

    return (
        <View style={styles.container}>
            <FlatList
                data={conversations}
                renderItem={renderItem}
                keyExtractor={(item, index) =>
                    item?.conversation_id != null
                        ? String(item.conversation_id)
                        : `conversation-${index}`
                }
                refreshing={isLoading}
                onRefresh={refetch}
                onEndReached={() => {
                    if (hasNextPage && !isFetchingNextPage) {
                        fetchNextPage();
                    }
                }}
                onEndReachedThreshold={0.6}
                ListEmptyComponent={renderEmptyState}
                ListFooterComponent={
                    isFetchingNextPage ? (
                        <View style={styles.footerLoader}>
                            <ActivityIndicator size="small" color={COLORS.primary} />
                        </View>
                    ) : null
                }
                contentContainerStyle={conversations.length === 0 ? styles.emptyContent : undefined}
            />

            <CreateConversationFAB />
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background.auth,
    },
    centerContainer: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        paddingVertical: SPACING['4xl'],
    },
    emptyContent: {
        flexGrow: 1,
    },
    emptyState: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        paddingHorizontal: SPACING.xl,
        paddingVertical: SPACING['4xl'],
    },
    emptyTitle: {
        fontSize: FONT_SIZES['2xl'],
        fontWeight: FONT_WEIGHTS.bold,
        color: COLORS.text.auth.primary,
        marginBottom: SPACING.sm,
        textAlign: 'center',
    },
    emptySubtitle: {
        fontSize: FONT_SIZES.base,
        color: COLORS.text.tertiary,
        textAlign: 'center',
        lineHeight: 24,
    },
    footerLoader: {
        paddingVertical: SPACING.base,
        alignItems: 'center',
    },
});