import { Pressable, View, StyleSheet } from 'react-native';
import { MotiView } from 'moti';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { useEffect, useState } from 'react';
import { COLORS } from '@/constants/theme';

export function AnimatedTabButton({ children, onPress, focused }: any) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';
  const [pressed, setPressed] = useState(false);
  const [rippleKey, setRippleKey] = useState(0);

  useEffect(() => {
    if (focused) {
      setRippleKey(prev => prev + 1);
    }
  }, [focused]);

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
            styles.glow,
            {
              backgroundColor: isDark 
                ? 'rgba(143, 229, 7, 0.3)' 
                : 'rgba(143, 229, 7, 0.25)',
            }
          ]}
          animate={{
            opacity: focused ? 0.6 : 0,
            scale: focused ? 1 : 0.8,
          }}
          transition={{
            type: 'timing',
            duration: 300,
          }}
        />
        
        <MotiView 
          key={rippleKey}
          style={[
            styles.ripple,
            {
              backgroundColor: isDark 
                ? 'rgba(143, 229, 7, 0.2)' 
                : 'rgba(143, 229, 7, 0.15)',
            }
          ]}
          from={{
            opacity: 1,
            scale: 0.5,
          }}
          animate={{
            opacity: 0,
            scale: 2,
          }}
          transition={{
            type: 'timing',
            duration: 600,
          }}
        />
        
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
  },
  glow: {
    position: 'absolute',
    width: 64,
    height: 64,
    borderRadius: 32,
    top: -8,
  },
  ripple: {
    position: 'absolute',
    width: 56,
    height: 56,
    borderRadius: 28,
    top: -8,
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