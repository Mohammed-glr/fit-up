import { useEffect, useRef, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/auth-context';
import { secureStorage } from '@/api/storage/secure-storage';
import { conversationKeys } from './use-conversation';
import type { WebSocketMessage, MessageWithDetails } from '@/types/message';

const getWebSocketUrl = () => {
    const apiUrl = process.env.EXPO_PUBLIC_API_URL || 'https://api.fitupp.nl';
    const wsProtocol = apiUrl.startsWith('https') ? 'wss' : 'ws';
    const baseUrl = apiUrl.replace(/^https?:\/\//, '');
    return `${wsProtocol}://${baseUrl}/ws`;
};

const WS_URL = getWebSocketUrl();

type WebSocketMessageType = 'new_message' | 'message_edited' | 'message_deleted' | 'message_read' | 'error';

export const useChatWebSocket = (conversationId?: number) => {
    const { user } = useAuth();
    const queryClient = useQueryClient();
    const ws = useRef<WebSocket | null>(null);
    const reconnectTimeout = useRef<ReturnType<typeof setTimeout>>();

    const connect = useCallback(async () => {
        if (!user) {
            return;
        }

        const token = await secureStorage.getToken('access_token');
        if (!token) {
            return;
        }

        if (token.length < 10) {
            return;
        }

        if (ws.current) {
            ws.current.close();
        }

        const encodedToken = encodeURIComponent(token);
        const wsUrl = `${WS_URL}?token=${encodedToken}`;
        
        try {
            ws.current = new WebSocket(wsUrl);
        } catch (error) {
            return;
        }

        ws.current.onopen = () => {}

        ws.current.onmessage = (event) => {
            try {
                if (event.data === 'ping') {
                    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
                        ws.current.send('pong');
                    }
                    return;
                }

                const data = JSON.parse(event.data) as WebSocketMessage;
                handleWebSocketMessage(data);
            } catch (error) {
                // Silent error handling
            }
        };

        ws.current.onerror = (error) => {
            if (ws.current) {
                ws.current.close();
            }
        };

        ws.current.onclose = (event) => {
            if (event.code === 1002) {
                return;
            }
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
        switch (data.type) {
            case 'new_message':
                if (data.message && data.conversation_id) {
                    queryClient.setQueryData(
                        conversationKeys.messages(data.conversation_id),
                        (oldData: any) => {
                            if (!oldData) return oldData;

                            const exists = oldData.pages.some((page: any) =>
                                page.messages.some((m: MessageWithDetails) => m.message_id === data.message?.message_id)
                            );

                            if (exists) return oldData;

                            const newPages = [...oldData.pages];
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
                if (data.conversation_id) {
                    queryClient.invalidateQueries({ queryKey: conversationKeys.unreadCount(data.conversation_id) });
                }
                break;
        }
    };
};
