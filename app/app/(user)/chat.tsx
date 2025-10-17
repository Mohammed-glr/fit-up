import React from 'react';
import { ChatView } from '@/components/chat/chat-view';

const ChatScreen = ({ route }: any) => {
    const { conversationId } = route.params;
    return <ChatView conversationId={conversationId} />;
};

export default ChatScreen;