import React, { useEffect, useState } from 'react';
import { View, FlatList, TextInput, TouchableOpacity, Text } from 'react-native';
import { useConversationMessages, useSendMessage, useMarkAllAsRead } from '@/hooks/message/use-conversation'
import type { MessageWithDetails } from '@/types/message';

export const ChatScreen = ({ route }: any) => {
    const { conversationId } = route.params;
    const [messageText, setMessageText] = useState('');
    const { data, isLoading, fetchNextPage, hasNextPage } = useConversationMessages(conversationId);
    const sendMessage = useSendMessage();
    const markAllAsRead = useMarkAllAsRead();

    useEffect(() => {
        markAllAsRead.mutate(conversationId);
    }, [conversationId]);

    const allMessages = data?.pages.flatMap(page => page.messages) || [];

    const handleSend = async () => {
        if (!messageText.trim()) return;

        await sendMessage.mutateAsync({ conversation_id: conversationId, message_text: messageText.trim() });
        setMessageText('');
    };

    const renderMessage = ({ item }: { item: MessageWithDetails }) => (
    <View style={{ padding: 12, marginVertical: 4 }}>
      <Text style={{ fontWeight: 'bold' }}>{item.sender_name}</Text>
      <Text>{item.message_text}</Text>
      <Text style={{ fontSize: 10, color: '#999' }}>
        {new Date(item.sent_at).toLocaleTimeString()}
      </Text>
    </View>
  );

  return (
    <View style={{ flex: 1 }}>
      <FlatList
        data={allMessages}
        renderItem={renderMessage}
        keyExtractor={(item) => item.message_id.toString()}
        inverted
        onEndReached={() => hasNextPage && fetchNextPage()}
        onEndReachedThreshold={0.5}
      />
      <View style={{ flexDirection: 'row', padding: 12, borderTopWidth: 1 }}>
        <TextInput
          style={{ flex: 1, borderWidth: 1, borderRadius: 20, paddingHorizontal: 16 }}
          value={messageText}
          onChangeText={setMessageText}
          placeholder="Type a message..."
          multiline
        />
        <TouchableOpacity
          onPress={handleSend}
          disabled={sendMessage.isPending}
          style={{ marginLeft: 8, justifyContent: 'center' }}
        >
          <Text style={{ color: 'blue', fontWeight: 'bold' }}>Send</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};