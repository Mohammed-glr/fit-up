import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Animated,
  Dimensions,
} from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BREATHING_PATTERNS, BreathingPattern } from '../../../types/mindfulness';
import { useCreateBreathingExercise } from '../../../hooks/mindfulness/use-mindfulness';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';

const { width } = Dimensions.get('window');
const CIRCLE_SIZE = width * 0.6;

type BreathingPhase = 'inhale' | 'hold1' | 'exhale' | 'hold2';

export default function BreathingExerciseScreen() {
  const params = useLocalSearchParams();
  const breathingType = params.type as string;
  const pattern = BREATHING_PATTERNS[breathingType];

  const [isActive, setIsActive] = useState(false);
  const [currentPhase, setCurrentPhase] = useState<BreathingPhase>('inhale');
  const [cycleCount, setCycleCount] = useState(0);
  const [timeInPhase, setTimeInPhase] = useState(0);
  const [totalTime, setTotalTime] = useState(0);

  const scaleAnim = useRef(new Animated.Value(0.5)).current;
  const opacityAnim = useRef(new Animated.Value(0.3)).current;

  const createExercise = useCreateBreathingExercise();

  const getPhaseConfig = (phase: BreathingPhase): number => {
    switch (phase) {
      case 'inhale':
        return pattern.pattern[0];
      case 'hold1':
        return pattern.pattern[1];
      case 'exhale':
        return pattern.pattern[2];
      case 'hold2':
        return pattern.pattern[3];
    }
  };

  const nextPhase = (current: BreathingPhase): BreathingPhase => {
    if (current === 'inhale') return pattern.pattern[1] > 0 ? 'hold1' : 'exhale';
    if (current === 'hold1') return 'exhale';
    if (current === 'exhale') return pattern.pattern[3] > 0 ? 'hold2' : 'inhale';
    return 'inhale';
  };

  const getPhaseText = (phase: BreathingPhase): string => {
    switch (phase) {
      case 'inhale':
        return 'Breathe In';
      case 'hold1':
      case 'hold2':
        return 'Hold';
      case 'exhale':
        return 'Breathe Out';
    }
  };

  const animatePhase = (phase: BreathingPhase) => {
    const duration = getPhaseConfig(phase) * 1000;

    if (phase === 'inhale') {
      Animated.parallel([
        Animated.timing(scaleAnim, {
          toValue: 1,
          duration,
          useNativeDriver: true,
        }),
        Animated.timing(opacityAnim, {
          toValue: 0.8,
          duration,
          useNativeDriver: true,
        }),
      ]).start();
    } else if (phase === 'exhale') {
      Animated.parallel([
        Animated.timing(scaleAnim, {
          toValue: 0.5,
          duration,
          useNativeDriver: true,
        }),
        Animated.timing(opacityAnim, {
          toValue: 0.3,
          duration,
          useNativeDriver: true,
        }),
      ]).start();
    }
  };

  useEffect(() => {
    if (!isActive) return;

    const interval = setInterval(() => {
      setTimeInPhase((prev) => {
        const phaseDuration = getPhaseConfig(currentPhase);
        const newTime = prev + 1;

        if (newTime >= phaseDuration) {
          const next = nextPhase(currentPhase);
          setCurrentPhase(next);
          animatePhase(next);

          if (next === 'inhale') {
            setCycleCount((c) => c + 1);
          }

          return 0;
        }

        return newTime;
      });

      setTotalTime((t) => t + 1);
    }, 1000);

    return () => clearInterval(interval);
  }, [isActive, currentPhase]);

  useEffect(() => {
    if (isActive) {
      animatePhase(currentPhase);
    }
  }, [isActive]);

  const handleStart = () => {
    setIsActive(true);
    setCurrentPhase('inhale');
    animatePhase('inhale');
  };

  const handleStop = () => {
    setIsActive(false);
  };

  const handleComplete = async () => {
    setIsActive(false);

    try {
      await createExercise.mutateAsync({
        breathing_type: pattern.type,
        duration_seconds: totalTime,
        cycles_completed: cycleCount,
      });

      router.back();
    } catch (error) {
      console.error('Failed to save breathing exercise:', error);
    }
  };

  if (!pattern) {
    return (
      <View style={styles.container}>
        <SafeAreaView style={styles.safeArea}>
          <Text style={styles.errorText}>Pattern not found</Text>
          <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
            <Text style={styles.backButtonText}>Go Back</Text>
          </TouchableOpacity>
        </SafeAreaView>
      </View>
    );
  }

  const phaseDuration = getPhaseConfig(currentPhase);
  const progress = timeInPhase / phaseDuration;

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        {/* <View style={styles.header}>
          <TouchableOpacity onPress={() => router.back()} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>✕</Text>
          </TouchableOpacity>
          <Text style={styles.title}>{pattern.name}</Text>
          <View style={styles.placeholder} />
        </View> */}

        <View style={styles.stats}>
          <View style={styles.statItem}>
            <Text style={styles.statValue}>{cycleCount}</Text>
            <Text style={styles.statLabel}>Cycles</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statValue}>{Math.floor(totalTime / 60)}:{(totalTime % 60).toString().padStart(2, '0')}</Text>
            <Text style={styles.statLabel}>Time</Text>
          </View>
        </View>

        <View style={styles.circleContainer}>
          <Animated.View
            style={[
              styles.circle,
              {
                transform: [{ scale: scaleAnim }],
                opacity: opacityAnim,
              },
            ]}
          />

          <View style={styles.phaseTextContainer}>
            <Text style={styles.phaseText}>{getPhaseText(currentPhase)}</Text>
            <Text style={styles.countdownText}>
              {phaseDuration - timeInPhase}
            </Text>
          </View>
        </View>

        <View style={styles.progressBarContainer}>
          <View style={[styles.progressBar, { width: `${progress * 100}%` }]} />
        </View>

        <View style={styles.patternInfo}>
          <Text style={styles.patternText}>
            {pattern.pattern[0]}s in • {pattern.pattern[1] > 0 ? `${pattern.pattern[1]}s hold • ` : ''}
            {pattern.pattern[2]}s out{pattern.pattern[3] > 0 ? ` • ${pattern.pattern[3]}s hold` : ''}
          </Text>
        </View>

        <View style={styles.controls}>
          {!isActive ? (
            <TouchableOpacity
              style={styles.startButton}
              onPress={handleStart}
              activeOpacity={0.8}
            >
              <Text style={styles.startButtonText}>Start</Text>
            </TouchableOpacity>
          ) : (
            <View style={styles.activeControls}>
              <TouchableOpacity
                style={styles.pauseButton}
                onPress={handleStop}
                activeOpacity={0.8}
              >
                <Text style={styles.pauseButtonText}>Pause</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={styles.completeButton}
                onPress={handleComplete}
                activeOpacity={0.8}
              >
                <Text style={styles.completeButtonText}>Complete</Text>
              </TouchableOpacity>
            </View>
          )}
        </View>

        <View style={styles.benefits}>
          <Text style={styles.benefitsTitle}>Benefits:</Text>
          {pattern.benefits.map((benefit, index) => (
            <Text key={index} style={styles.benefitItem}>
              • {benefit}
            </Text>
          ))}
        </View>
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
    padding: 20,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 32,
  },
  closeButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
  },
  closeButtonText: {
    fontSize: 24,
    color: '#FFFFFF',
  },
  title: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  placeholder: {
    width: 40,
  },
  stats: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 48,
  },
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 28,
    fontWeight: '700',
    color: '#6C63FF',
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 14,
    color: '#888888',
  },
  circleContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    height: CIRCLE_SIZE,
    marginBottom: 32,
  },
  circle: {
    width: CIRCLE_SIZE,
    height: CIRCLE_SIZE,
    borderRadius: CIRCLE_SIZE / 2,
    backgroundColor: '#6C63FF',
    position: 'absolute',
  },
  phaseTextContainer: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  phaseText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 8,
  },
  countdownText: {
    fontSize: 48,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  progressBarContainer: {
    height: 4,
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS.sm,
    marginBottom: 16,
    overflow: 'hidden',
    width: '96%',
    alignSelf: 'center',
  },
  progressBar: {
    height: '100%',
    backgroundColor: '#6C63FF',
    borderRadius: 2,
  },
  patternInfo: {
    alignItems: 'center',
    marginBottom: 32,
  },
  patternText: {
    fontSize: 14,
    color: '#888888',
  },
  controls: {
    marginBottom: 32,
  },
  startButton: {
    backgroundColor: '#6C63FF',
    borderRadius: BORDER_RADIUS['2xl'],
    paddingVertical: SPACING.lg,
    alignItems: 'center',
  },
  startButtonText: {
    fontSize: 18,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  activeControls: {
    flexDirection: 'row',
    gap: 12,
  },
  pauseButton: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
    borderRadius: BORDER_RADIUS.full,
    paddingVertical: 18,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#6C63FF',
  },
  pauseButtonText: {
    fontSize: 18,
    fontWeight: '700',
    color: '#6C63FF',
  },
  completeButton: {
    flex: 1,
    backgroundColor: '#6C63FF',
    borderRadius: BORDER_RADIUS.full,
    paddingVertical: 18,
    alignItems: 'center',
  },
  completeButtonText: {
    fontSize: 18,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  benefits: {
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
  },
  benefitsTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 12,
  },
  benefitItem: {
    fontSize: 14,
    color: '#888888',
    marginBottom: 8,
  },
  errorText: {
    fontSize: 18,
    color: '#FF4444',
    textAlign: 'center',
    marginBottom: 20,
  },
  backButton: {
    backgroundColor: '#6C63FF',
    borderRadius: 16,
    paddingVertical: 16,
    paddingHorizontal: 32,
    alignItems: 'center',
  },
  backButtonText: {
    fontSize: 16,
    fontWeight: '700',
    color: '#FFFFFF',
  },
});
