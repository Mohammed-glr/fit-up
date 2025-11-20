import React from "react";
import { View, Text, FlatList, ActivityIndicator, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useConversations } from "@/hooks/message/use-conversation";
import type { ConversationOverview } from "@/types";
import { CreateConversationFAB } from "../../components/chat/createConversationFAB";
import { ConversationItem } from "../../components/chat/conversation-item";
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
import Ionicons from "@expo/vector-icons/build/Ionicons";

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

    const conversations = data?.pages.flatMap((page) => page.conversations) ?? [];

    const renderItem = ({ item }: { item: ConversationOverview }) => {
        return (
            <ConversationItem
                conversation={item}
                onPress={() =>
                    router.push({
                        pathname: '/(user)/chat',
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

            <View style={styles.header}>
                <View style={styles.titleContainer}>
                    <Ionicons name="chatbubbles" size={32} color={COLORS.primary} />
                    <Text style={styles.headerTitle}>Conversation</Text>
                </View>
                <Text style={styles.headerSubtitle}>
                    Chat with your coach and stay connected!
                </Text>
            </View>

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
    header: {
        marginBottom: 24,
        margin: SPACING.base,
    },
    titleContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 12,
        marginBottom: 4,
    },
    title: {
        fontSize: 32,
        fontWeight: '700',
        color: '#FFFFFF',
    },
    subtitle: {
        fontSize: 16,
        color: '#888888',
    },
    headerInfo: {
        flex: 1,
    },
    headerTitle: {
        fontSize: 28,
        fontWeight: '700',
        color: '#FFFFFF',
    },
    headerSubtitle: {
        fontSize: 14,
        color: '#888888',
        marginTop: 2,
    },
});