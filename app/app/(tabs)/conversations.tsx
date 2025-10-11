import React from "react";

import { View, Text, FlatList, TouchableOpacity, ActivityIndicator } from 'react-native';
import { useConversations } from "@/hooks/message/use-conversation";
import type { ConversationOverview } from "@/types";


export const ConversationsScreen = ({ navigation }: any) => {
    const {
        data,
        isLoading,
        isError,
        refetch,
    } = useConversations();

    if (isLoading) {
        return <View style={styles.container}><Text>Loading...</Text></View>;
    }

    if (isError) {
        return <View style={styles.container}><Text>Error loading conversations.</Text></View>;
    }

    const renderItem = ({ item }: { item: ConversationOverview}) => {
        return (
            <TouchableOpacity
                onPress={() => navigation.navigate('Chat', { conversationId: item.conversation_id })}
                style={{ padding: 16, borderBottomWidth: 1, borderBottomColor: '#eee' }}
            >
                <View style={{ flexDirection: 'row', justifyContent: 'space-between' }}>
        <Text style={{ fontWeight: 'bold' }}>
          {item.coach_name} - {item.client_name}
        </Text>
        {item.total_messages > 0 && (
          <View style={{ backgroundColor: 'red', borderRadius: 10, padding: 4 }}>
            <Text style={{ color: 'white', fontSize: 12 }}>{item.total_messages}</Text>
          </View>
        )}
      </View>
      {item.last_message_text && (
        <Text style={{ color: '#666', marginTop: 4 }} numberOfLines={1}>
          {item.last_message_text}
        </Text>
      )}
    </TouchableOpacity>
    );  
  };

     return (
    <View style={{ flex: 1 }}>
      <FlatList
        data={data?.conversations || []}
        renderItem={renderItem}
        keyExtractor={(item) => item.conversation_id.toString()}
        refreshing={isLoading}
        onRefresh={refetch}
      />
    </View>
  );
}

const styles = {
    container: {
        flex: 1,
    },
};