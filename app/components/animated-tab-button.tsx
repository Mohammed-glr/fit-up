import { Pressable, View, StyleSheet } from 'react-native';
import { MotiView } from 'moti';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { useEffect, useState } from 'react';
import { COLORS } from '@/constants/theme';

export function AnimatedTabButton({ children, onPress, focused }: any) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';
  const [pressed, setPressed] = useState(false);

  return (
    <Pressable
      onPress={onPress}
      onPressIn={() => setPressed(true)}
      onPressOut={() => setPressed(false)}
      style={styles.container}
    >
      <MotiView
        style={styles.content}
        animate={{
          scale: pressed ? 0.88 : (focused ? 1.05 : 0.92),
          translateY: focused ? -6 : 0,
        }}
        transition={{
          type: 'spring',
          damping: 15,
          stiffness: 150,
        }}
      > 
        <MotiView 
          style={[
            styles.background,
            {
              backgroundColor: isDark 
                ? 'rgba(143, 229, 7, 0.18)' 
                : 'rgba(143, 229, 7, 0.12)',
            }
          ]}
          animate={{
            opacity: focused ? 1 : 0,
            scale: focused ? 1 : 0.8,
          }}
          transition={{
            type: 'spring',
            damping: 15,
            stiffness: 150,
          }}
        />
        <View style={styles.iconContainer}>
          {children}
        </View>
      </MotiView>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 8,
    
  },
  content: {
    alignItems: 'center',
    justifyContent: 'center',
    position: 'relative',
    backgroundColor: 'transparent',
  },
  background: {
    position: 'absolute',
    width: 56,
    height: 56,
    borderRadius: 28,
    top: -8,
  },
  iconContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1,
  },
});