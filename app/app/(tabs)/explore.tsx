import {
  View,
  StyleSheet,
  useColorScheme,
} from 'react-native';

export default function TabTwoScreen() {
  // placeholder for now
  const colorScheme = useColorScheme();
  return <View style={[styles.container, { backgroundColor: colorScheme === 'dark' ? '#000' : '#fff' }]} />;
}


const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});
