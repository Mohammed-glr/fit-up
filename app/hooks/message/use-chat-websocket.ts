import { useEffect, useRef, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/auth-context';
import { secureStorage } from '@/api/storage/secure-storage';
import { conversationKeys } from './use-conversation';
import type { WebSocketMessage, MessageWithDetails } from '@/types/message';

const WS_URL = process.env.EXPO_PUBLIC_WS_URL || 'ws://localhost:8080/ws';

type WebSocketMessageType = 'new_message' | 'message_edited' | 'message_deleted' | 'message_read' | 'error';

export const useChatWebSocket = (conversationId?: number) => {
    const { user } = useAuth();
    const queryClient = useQueryClient();
    const ws = useRef<WebSocket | null>(null);
    const reconnectTimeout = useRef<ReturnType<typeof setTimeout>>();

    const connect = useCallback(async () => {
        if (!user) return;

        const token = await secureStorage.getToken('access_token');
        if (!token) return;

        // Close existing connection if any
        if (ws.current) {
            ws.current.close();
        }

        const wsUrl = `${WS_URL}?token=${token}`;
        ws.current = new WebSocket(wsUrl);

        ws.current.onopen = () => {
            console.log('WebSocket Connected');
        };

        ws.current.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data) as WebSocketMessage;
                handleWebSocketMessage(data);
            } catch (error) {
                console.error('Failed to parse WebSocket message:', error);
            }
        };

        ws.current.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        ws.current.onclose = () => {
            console.log('WebSocket Disconnected');
            // Attempt to reconnect after a delay
            reconnectTimeout.current = setTimeout(() => {
                connect();
            }, 3000);
        };
    }, [user]);

    useEffect(() => {
        connect();

        return () => {
            if (ws.current) {
                ws.current.close();
            }
            if (reconnectTimeout.current) {
                clearTimeout(reconnectTimeout.current);
            }
        };
    }, [connect]);

    const handleWebSocketMessage = (data: WebSocketMessage) => {
        console.log('Received WebSocket message:', data.type, data);

        switch (data.type) {
            case 'new_message':
                if (data.message && data.conversation_id) {
                    // Update messages list
                    queryClient.setQueryData(
                        conversationKeys.messages(data.conversation_id),
                        (oldData: any) => {
                            if (!oldData) return oldData;

                            // Check if message already exists (optimistic update)
                            const exists = oldData.pages.some((page: any) =>
                                page.messages.some((m: MessageWithDetails) => m.message_id === data.message?.message_id)
                            );

                            if (exists) return oldData;

                            const newPages = [...oldData.pages];
                            // Add to the first page (newest messages)
                            if (newPages.length > 0) {
                                newPages[0] = {
                                    ...newPages[0],
                                    messages: [data.message, ...newPages[0].messages],
                                };
                            }

                            return {
                                ...oldData,
                                pages: newPages,
                            };
                        }
                    );

                    // Update conversation details (last message)
                    queryClient.setQueryData(
                        conversationKeys.detail(data.conversation_id),
                        (oldData: any) => {
                            if (!oldData) return oldData;
                            return {
                                ...oldData,
                                conversation: {
                                    ...oldData.conversation,
                                    last_message: data.message,
                                    last_message_at: data.timestamp,
                                },
                            };
                        }
                    );

                    // Update conversations list
                    queryClient.invalidateQueries({ queryKey: conversationKeys.lists() });
                }
                break;

            case 'message_edited':
                if (data.message && data.conversation_id) {
                    queryClient.setQueryData(
                        conversationKeys.messages(data.conversation_id),
                        (oldData: any) => {
                            if (!oldData) return oldData;

                            const newPages = oldData.pages.map((page: any) => ({
                                ...page,
                                messages: page.messages.map((m: MessageWithDetails) =>
                                    m.message_id === data.message?.message_id ? data.message : m
                                ),
                            }));

                            return { ...oldData, pages: newPages };
                        }
                    );
                }
                break;

            case 'message_deleted':
                if (data.message_id && data.conversation_id) {
                    queryClient.setQueryData(
                        conversationKeys.messages(data.conversation_id),
                        (oldData: any) => {
                            if (!oldData) return oldData;

                            const newPages = oldData.pages.map((page: any) => ({
                                ...page,
                                messages: page.messages.filter((m: MessageWithDetails) => m.message_id !== data.message_id),
                            }));

                            return { ...oldData, pages: newPages };
                        }
                    );
                }
                break;

            case 'message_read':
                // Invalidate unread count
                if (data.conversation_id) {
                    queryClient.invalidateQueries({ queryKey: conversationKeys.unreadCount(data.conversation_id) });
                }
                break;
        }
    };
};
