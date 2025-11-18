import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ActivityIndicator } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { MotiView } from 'moti';
import { useWorkoutHistory } from '@/hooks/user/use-workout-history';
import { COLORS, SPACING, FONT_SIZES, BORDER_RADIUS, FONT_WEIGHTS } from '@/constants/theme';

// Simple date formatter
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  return `${months[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()}`;
};

export const WorkoutHistorySummary = () => {
  const router = useRouter();
  const { data: historyData, isLoading } = useWorkoutHistory({
    page: 1,
    pageSize: 3, // Just show last 3 workouts
  });

  if (isLoading) {
    return (
      <View style={styles.container}>
        <View style={styles.header}>
          <Text style={styles.title}>Recent Workouts</Text>
        </View>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="small" color={COLORS.primary} />
        </View>
      </View>
    );
  }

  if (!historyData || historyData.workouts.length === 0) {
    return null; 
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Recent Workouts</Text>
        <TouchableOpacity
          onPress={() => router.push('/(user)/workout-history')}
          style={styles.viewAllButton}
        >
          <Text style={styles.viewAllText}>View All</Text>
          <Ionicons name="chevron-forward" size={16} color={COLORS.primary} />
        </TouchableOpacity>
      </View>

      <View>
        {historyData.workouts.map((workout, index) => (
         <MotiView
                     key={workout.date}
                     from={{ opacity: 0, translateY: 20 }}
                     animate={{ opacity: 1, translateY: 0 }}
                     transition={{ type: 'timing', duration: 300, delay: index * 50 }}
                     style={styles.workoutCard}
                   >
                     {/* Date Header */}
                     <View style={styles.dateHeader}>
                       <View style={styles.dateIconContainer}>
                         <Ionicons name="calendar" size={20} color={COLORS.primary} />
                       </View>
                       <View style={styles.dateInfo}>
                         <Text style={styles.dateText}>{formatDate(workout.date)}</Text>
                         {workout.day_title && (
                           <Text style={styles.dayTitle}>{workout.day_title}</Text>
                         )}
                       </View>
                     </View>
         
                     {/* Workout Stats */}
                     <View style={styles.workoutStats}>
                       <View style={styles.statItem}>
                         <Ionicons name="fitness" size={16} color={COLORS.text.secondary} />
                         <Text style={styles.statItemText}>
                           {workout.total_exercises} exercises
                         </Text>
                       </View>
                       <View style={styles.statItem}>
                         <Ionicons name="checkmark-circle" size={16} color={COLORS.success} />
                         <Text style={styles.statItemText}>
                           {workout.completed_sets} sets
                         </Text>
                       </View>
                       <View style={styles.statItem}>
                         <Ionicons name="barbell" size={16} color={COLORS.warning} />
                         <Text style={styles.statItemText}>
                           {Math.round(workout.total_volume)} lbs
                         </Text>
                       </View>
                       <View style={styles.statItem}>
                         <Ionicons name="time" size={16} color={COLORS.info} />
                         <Text style={styles.statItemText}>
                           {workout.duration_minutes} min
                         </Text>
                       </View>
                     </View>
         
                     {/* Exercises List */}
                     {workout.exercises && workout.exercises.length > 0 && (
                       <View style={styles.exercisesList}>
                         <Text style={styles.exercisesLabel}>Exercises:</Text>
                         <View style={styles.exerciseTags}>
                           {workout.exercises.slice(0, 3).map((exercise, idx) => (
                             <View key={idx} style={styles.exerciseTag}>
                               <Text style={styles.exerciseTagText} numberOfLines={1}>
                                 {exercise}
                               </Text>
                             </View>
                           ))}
                           {workout.exercises.length > 3 && (
                             <View style={styles.exerciseTag}>
                               <Text style={styles.exerciseTagText}>
                                 +{workout.exercises.length - 3} more
                               </Text>
                             </View>
                           )}
                         </View>
                       </View>
                     )}
                   </MotiView>
        ))}
      </View>

      {/* <TouchableOpacity
        style={styles.viewAllCard}
        onPress={() => router.push('/(user)/workout-history')}
        activeOpacity={0.7}
      >
        <Text style={styles.viewAllCardText}>View Full History</Text>
        <Ionicons name="arrow-forward" size={20} color={COLORS.primary} />
      </TouchableOpacity> */}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: SPACING.xl,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: '700',
    color: COLORS.text.inverse,
  },
  viewAllButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  viewAllText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: '600',
    color: COLORS.primary,
  },
  loadingContainer: {
    padding: SPACING.xl,
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
  },
 workoutCard: {
     backgroundColor: COLORS.background.card,
     borderRadius: BORDER_RADIUS['2xl'],
     padding: SPACING.base,
     marginBottom: SPACING.md,
     shadowColor: '#000',
     shadowOffset: { 
       width: 0, 
       height: 12 
     },
     shadowOpacity: 0.35,
   shadowRadius: 24,
   },
   dateHeader: {
     flexDirection: 'row',
     alignItems: 'center',
     marginBottom: SPACING.md,
   },
   dateIconContainer: {
     width: 40,
     height: 40,
     borderRadius: BORDER_RADIUS.full,
     backgroundColor: `${COLORS.primary}15`,
     justifyContent: 'center',
     alignItems: 'center',
     marginRight: SPACING.md,
   },
   dateInfo: {
     flex: 1,
   },
   dateText: {
     fontSize: FONT_SIZES.base,
     fontWeight: FONT_WEIGHTS.semibold as any,
     color: COLORS.text.inverse,
   },
   dayTitle: {
     fontSize: FONT_SIZES.sm,
     color: COLORS.text.tertiary,
     marginTop: 2,
   },
   workoutStats: {
     flexDirection: 'row',
     flexWrap: 'wrap',
     gap: SPACING.md,
     marginBottom: SPACING.md,
   },
   statItem: {
     flexDirection: 'row',
     alignItems: 'center',
     gap: SPACING.xs,
   },
   statItemText: {
     fontSize: FONT_SIZES.sm,
     color: COLORS.text.tertiary,
   },
   exercisesList: {
   },
   exercisesLabel: {
     fontSize: FONT_SIZES.sm,
     fontWeight: FONT_WEIGHTS.medium as any,
     color: COLORS.text.secondary,
     marginBottom: SPACING.xs,
   },
   exerciseTags: {
     flexDirection: 'row',
     flexWrap: 'wrap',
     gap: SPACING.xs,
   },
   exerciseTag: {
     backgroundColor: `${COLORS.primary}10`,
     paddingHorizontal: SPACING.sm,
     paddingVertical: SPACING.xs / 2,
     borderRadius: BORDER_RADIUS.sm,
     maxWidth: 150,
   },
   exerciseTagText: {
     fontSize: FONT_SIZES.xs,
     color: COLORS.primary,
   },
  viewAllCard: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: SPACING.md,
    gap: SPACING.sm,
  },
  viewAllCardText: {
    fontSize: FONT_SIZES.base,
    fontWeight: '600',
    color: COLORS.primary,
  },
});
