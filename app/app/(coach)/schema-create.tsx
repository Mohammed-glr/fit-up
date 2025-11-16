import React from 'react';
import { View, StyleSheet, Alert } from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import { SchemaForm } from '@/components/schema/schema-form';
import { useCreateSchema } from '@/hooks/schema/use-coach';
import { useAuth } from '@/context/auth-context';
import type { ManualSchemaRequest } from '@/types/schema';
import { COLORS } from '@/constants/theme';

export default function SchemaCreateScreen() {
  const { user } = useAuth();
  const params = useLocalSearchParams<{ userId?: string }>();
  const createSchemaMutation = useCreateSchema();

  const handleSubmit = async (data: ManualSchemaRequest) => {
    if (!user?.id) {
      Alert.alert('Error', 'User not authenticated');
      return;
    }

    const userId = params.userId ? parseInt(params.userId) : 0;
    console.log('Schema Create - Raw userId param:', params.userId);
    console.log('Schema Create - Parsed userId:', userId);
    console.log('Schema Create - Coach ID:', user.id);
    
    if (!userId) {
      Alert.alert('Error', 'No user selected');
      return;
    }

    const requestData: ManualSchemaRequest = {
      ...data,
      user_id: userId,
      coach_id: user.id,
      start_date: new Date().toISOString(),
    };

    console.log('Schema Create - Final request data:', JSON.stringify(requestData, null, 2));

    try {
      await createSchemaMutation.mutateAsync({
        userID: userId,
        schema: requestData,
      });
      Alert.alert('Success', 'Schema created successfully', [
        {
          text: 'OK',
          onPress: () => router.back(),
        },
      ]);
    } catch (error: any) {
      console.error('Schema Create - Error:', error);
      Alert.alert('Error', error.message || 'Failed to create schema');
    }
  };

  const handleCancel = () => {
    Alert.alert(
      'Discard Changes',
      'Are you sure you want to discard this schema?',
      [
        { text: 'Continue Editing', style: 'cancel' },
        {
          text: 'Discard',
          style: 'destructive',
          onPress: () => router.back(),
        },
      ]
    );
  };

  return (
    <View style={styles.container}>
      <SchemaForm
        onSubmit={handleSubmit}
        onCancel={handleCancel}
        isLoading={createSchemaMutation.isPending}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
});
