import React from 'react';
import { useLocalSearchParams } from 'expo-router';

import { ChatView } from '@/components/chat/chat-view';

type ChatScreenProps = {
    route?: { params?: { conversationId?: number | string } };
};

const ChatScreen: React.FC<ChatScreenProps> = (props) => {
    const route = props?.route;
    const searchParams = useLocalSearchParams<{ conversationId?: string }>();

    const paramValue =
        searchParams.conversationId ?? route?.params?.conversationId?.toString();

    const conversationId = paramValue ? Number(paramValue) : undefined;

    if (!conversationId || Number.isNaN(conversationId)) {
        return null;
    }

    return <ChatView conversationId={conversationId} />;
};

export default ChatScreen;