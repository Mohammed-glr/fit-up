import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Modal,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as Clipboard from 'expo-clipboard';
import * as Sharing from 'expo-sharing';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import type { WorkoutShareSummary } from '@/types/workout-sharing';
import { Button } from '../forms/button';

interface ShareWorkoutModalProps {
  visible: boolean;
  onClose: () => void;
  sessionId: number;
  shareSummary?: WorkoutShareSummary;
  isLoading?: boolean;
}

export const ShareWorkoutModal: React.FC<ShareWorkoutModalProps> = ({
  visible,
  onClose,
  sessionId,
  shareSummary,
  isLoading = false,
}) => {
  const [sharing, setSharing] = useState(false);

  const generateShareText = () => {
    if (!shareSummary) return '';

    let text = `💪 ${shareSummary.workout_title}\n\n`;
    text += `⏱️ Duration: ${shareSummary.duration_minutes} minutes\n`;
    text += `🏋️ Exercises: ${shareSummary.total_exercises}\n`;
    text += `📊 Total Sets: ${shareSummary.total_sets}\n`;
    text += `🔢 Total Reps: ${shareSummary.total_reps}\n`;
    text += `💯 Total Volume: ${shareSummary.total_volume_lbs.toFixed(1)} lbs\n\n`;

    if (shareSummary.prs_achieved > 0) {
      text += `🎉 ${shareSummary.prs_achieved} PRs Achieved!\n\n`;
    }

    text += 'Exercises:\n';
    shareSummary.exercises.forEach((ex) => {
      text += `• ${ex.exercise_name}: ${ex.sets_completed} sets`;
      if (ex.best_set) {
        text += ` (best: ${ex.best_set.weight} lbs × ${ex.best_set.reps} reps)`;
      }
      if (ex.pr_achieved) {
        text += ' 🔥 PR!';
      }
      text += '\n';
    });

    text += '\n#FitUp #WorkoutComplete';
    return text;
  };

  const handleCopyToClipboard = async () => {
    try {
      setSharing(true);
      const text = generateShareText();
      await Clipboard.setStringAsync(text);
      Alert.alert('Copied!', 'Workout summary copied to clipboard');
      onClose();
    } catch (error) {
      Alert.alert('Error', 'Failed to copy to clipboard');
    } finally {
      setSharing(false);
    }
  };

  const handleShareToSocial = async () => {
    try {
      setSharing(true);
      const text = generateShareText();
      
      if (await Sharing.isAvailableAsync()) {
        await Sharing.shareAsync('data:text/plain;base64,' + btoa(text), {
          mimeType: 'text/plain',
          dialogTitle: 'Share Workout',
        });
      } else {
        // Fallback to clipboard
        await Clipboard.setStringAsync(text);
        Alert.alert('Copied!', 'Workout summary copied to clipboard');
      }
      onClose();
    } catch (error) {
      Alert.alert('Error', 'Failed to share workout');
    } finally {
      setSharing(false);
    }
  };

  const handleShareWithCoach = () => {
    // TODO: Implement coach sharing
    Alert.alert('Coming Soon', 'Share with coach feature will be available soon!');
    onClose();
  };

  const handleExportAsImage = () => {
    // TODO: Implement image export
    Alert.alert('Coming Soon', 'Export as image feature will be available soon!');
    onClose();
  };

  if (!shareSummary) {
    return null;
  }

  return (
    <Modal
      visible={visible}
      transparent
      animationType="slide"
      onRequestClose={onClose}
    >
      <View style={styles.overlay}>
        <View style={styles.modalContainer}>
          {/* Header */}
          <View style={styles.header}>
            <Text style={styles.title}>Share Workout</Text>
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color={COLORS.text.primary} />
            </TouchableOpacity>
          </View>

          {isLoading ? (
            <View style={styles.loadingContainer}>
              <ActivityIndicator size="large" color={COLORS.primary} />
              <Text style={styles.loadingText}>Preparing workout summary...</Text>
            </View>
          ) : (
            <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
              {/* Workout Summary Preview */}
              <View style={styles.summaryCard}>
                <Text style={styles.workoutTitle}>{shareSummary.workout_title}</Text>
                
                <View style={styles.statsGrid}>
                  <View style={styles.statItem}>
                    <Ionicons name="time-outline" size={24} color={COLORS.primary} />
                    <Text style={styles.statValue}>{shareSummary.duration_minutes}</Text>
                    <Text style={styles.statLabel}>Minutes</Text>
                  </View>
                  
                  <View style={styles.statItem}>
                    <Ionicons name="barbell-outline" size={24} color={COLORS.primary} />
                    <Text style={styles.statValue}>{shareSummary.total_exercises}</Text>
                    <Text style={styles.statLabel}>Exercises</Text>
                  </View>
                  
                  <View style={styles.statItem}>
                    <Ionicons name="fitness-outline" size={24} color={COLORS.primary} />
                    <Text style={styles.statValue}>{shareSummary.total_sets}</Text>
                    <Text style={styles.statLabel}>Sets</Text>
                  </View>
                  
                  <View style={styles.statItem}>
                    <Ionicons name="trending-up-outline" size={24} color={COLORS.primary} />
                    <Text style={styles.statValue}>{shareSummary.total_volume_lbs.toFixed(0)}</Text>
                    <Text style={styles.statLabel}>lbs</Text>
                  </View>
                </View>

                {shareSummary.prs_achieved > 0 && (
                  <View style={styles.prBadge}>
                    <Ionicons name="trophy" size={20} color={COLORS.warning} />
                    <Text style={styles.prText}>{shareSummary.prs_achieved} PRs Achieved!</Text>
                  </View>
                )}
              </View>

              {/* Share Options */}
              <View style={styles.optionsContainer}>
                <TouchableOpacity
                  style={styles.shareOption}
                  onPress={handleCopyToClipboard}
                  disabled={sharing}
                >
                  <View style={styles.shareIconContainer}>
                    <Ionicons name="copy-outline" size={28} color={COLORS.primary} />
                  </View>
                  <Text style={styles.shareOptionTitle}>Copy Text</Text>
                  <Text style={styles.shareOptionSubtitle}>Copy to clipboard</Text>
                </TouchableOpacity>

                <TouchableOpacity
                  style={styles.shareOption}
                  onPress={handleShareToSocial}
                  disabled={sharing}
                >
                  <View style={styles.shareIconContainer}>
                    <Ionicons name="share-social-outline" size={28} color={COLORS.primary} />
                  </View>
                  <Text style={styles.shareOptionTitle}>Share</Text>
                  <Text style={styles.shareOptionSubtitle}>Share to social</Text>
                </TouchableOpacity>

                <TouchableOpacity
                  style={styles.shareOption}
                  onPress={handleShareWithCoach}
                  disabled={sharing}
                >
                  <View style={styles.shareIconContainer}>
                    <Ionicons name="person-outline" size={28} color={COLORS.primary} />
                  </View>
                  <Text style={styles.shareOptionTitle}>Send to Coach</Text>
                  <Text style={styles.shareOptionSubtitle}>Share with your coach</Text>
                </TouchableOpacity>

                <TouchableOpacity
                  style={styles.shareOption}
                  onPress={handleExportAsImage}
                  disabled={sharing}
                >
                  <View style={styles.shareIconContainer}>
                    <Ionicons name="image-outline" size={28} color={COLORS.primary} />
                  </View>
                  <Text style={styles.shareOptionTitle}>Export Image</Text>
                  <Text style={styles.shareOptionSubtitle}>Save as picture</Text>
                </TouchableOpacity>
              </View>
            </ScrollView>
          )}

          {/* Footer */}
          {!isLoading && (
            <View style={styles.footer}>
              <Button
                onPress={onClose}
                title="Close"
                variant="outline"
              />
            </View>
          )}
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  modalContainer: {
    backgroundColor: COLORS.background.secondary,
    borderTopLeftRadius: BORDER_RADIUS['2xl'],
    borderTopRightRadius: BORDER_RADIUS['2xl'],
    maxHeight: '90%',
    ...SHADOWS.lg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: SPACING.lg,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
  },
  closeButton: {
    padding: SPACING.xs,
  },
  loadingContainer: {
    padding: SPACING['3xl'],
    alignItems: 'center',
  },
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.secondary,
  },
  content: {
    flex: 1,
  },
  summaryCard: {
    margin: SPACING.lg,
    padding: SPACING.lg,
    backgroundColor: COLORS.background.auth,
    borderRadius: BORDER_RADIUS.lg,
    ...SHADOWS.sm,
  },
  workoutTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING.lg,
    textAlign: 'center',
  },
  statsGrid: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: SPACING.md,
  },
  statItem: {
    alignItems: 'center',
    flex: 1,
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    marginTop: SPACING.xs,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
    marginTop: SPACING['2xs'],
  },
  prBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: `${COLORS.warning}20`,
    padding: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
    marginTop: SPACING.md,
  },
  prText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.warning,
    marginLeft: SPACING.xs,
  },
  optionsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    padding: SPACING.md,
    gap: SPACING.md,
  },
  shareOption: {
    width: '48%',
    backgroundColor: COLORS.background.auth,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.lg,
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  shareIconContainer: {
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: `${COLORS.primary}20`,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: SPACING.md,
  },
  shareOptionTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING['2xs'],
  },
  shareOptionSubtitle: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
    textAlign: 'center',
  },
  footer: {
    padding: SPACING.lg,
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
  },
});
