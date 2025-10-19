import React from 'react';
import { SafeAreaView, StyleSheet, Text, View } from 'react-native';
import {
  DashboardGreeting
} from '@/components/dashboard/greeting';
export default function DashboardScreen() {
  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        <DashboardGreeting />
        <Text style={styles.text}>Dashboard placeholder</Text>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0A0A0A',
  },
  safeArea: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  text: {
    fontSize: 20,
    fontWeight: '600',
    color: '#FFFFFF',
  },
});
