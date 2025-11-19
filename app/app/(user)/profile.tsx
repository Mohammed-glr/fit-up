import React, { useState } from 'react';
import {
  SafeAreaView,
  StyleSheet,
  Text,
  View,
  ScrollView,
  Image,
  TouchableOpacity,
  Pressable,
  ActivityIndicator,
} from 'react-native';
import { useRouter } from 'expo-router';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { useUpdateProfile } from '@/hooks/user/use-update-profile';
import { useUserStats } from '@/hooks/user/use-user-stats';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { MotiView } from 'moti';
import LogoutButton from '@/components/auth/logout-button';
import { Ionicons } from '@expo/vector-icons';
import { EditBioModal } from '@/components/profile/EditBioModal';
import { EditImageModal } from '@/components/profile/EditImageModal';
import { useToastMethods } from '@/components/ui/toast-provider';

export default function ProfileScreen() {
  const router = useRouter();
  const { data: currentUser, isLoading } = useCurrentUser();
  const { data: userStats, isLoading: isLoadingStats } = useUserStats();
  const { mutate: updateProfile, isPending: isUpdating } = useUpdateProfile();
  const { showSuccess, showError } = useToastMethods();
  
  const [isBioModalVisible, setIsBioModalVisible] = useState(false);
  const [isImageModalVisible, setIsImageModalVisible] = useState(false);

  const handleBioSave = (bio: string) => {
    updateProfile(
      { bio },
      {
        onSuccess: () => {
          setIsBioModalVisible(false);
          showSuccess('Bio updated successfully!');
        },
        onError: (error: any) => {
          showError(error?.message || 'Failed to update bio');
        },
      }
    );
  };

  const handleImageSelect = (imageUri: string) => {
    updateProfile(
      { image: imageUri },
      {
        onSuccess: () => {
          setIsImageModalVisible(false);
          showSuccess('Profile image updated successfully!');
        },
        onError: (error: any) => {
          showError(error?.message || 'Failed to update profile image');
        },
      }
    );
  };

  if (isLoading) {
    return (
      <View style={styles.container}>
        <SafeAreaView style={styles.safeArea}>
          <Text style={styles.loadingText}>Loading...</Text>
        </SafeAreaView>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false}>
        <MotiView
          from={{ opacity: 0, translateY: -20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400 }}
          style={styles.header}
        >
          <View style={styles.profileImageContainer}>
            {currentUser?.image ? (
              <Image source={{ uri: currentUser.image }} style={styles.profileImage} />
            ) : (
              <View style={[styles.profileImage, styles.placeholderImage]}>
                <Ionicons name="person" size={60} color={COLORS.text.tertiary} />
              </View>
            )}
            <TouchableOpacity 
              style={styles.editImageButton}
              onPress={() => setIsImageModalVisible(true)}
            >
              <Ionicons name="camera" size={20} color={COLORS.white} />
            </TouchableOpacity>
          </View>

          <Text style={styles.name}>{currentUser?.name || 'User'}</Text>
          <Text style={styles.username}>@{currentUser?.username || 'username'}</Text>
          <Text style={styles.email}>{currentUser?.email}</Text>
        </MotiView>

        <MotiView
          from={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ type: 'timing', duration: 400, delay: 100 }}
          style={styles.statsContainer}
        >
          <View style={styles.statItem}>
            {isLoadingStats ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.statValue}>{userStats?.total_workouts || 0}</Text>
                <Text style={styles.statLabel}>Workouts</Text>
              </>
            )}
          </View>
          <View style={styles.statDivider} />
          <View style={styles.statItem}>
            {isLoadingStats ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.statValue}>{userStats?.active_programs || 0}</Text>
                <Text style={styles.statLabel}>Programs</Text>
              </>
            )}
          </View>
          <View style={styles.statDivider} />
          <View style={styles.statItem}>
            {isLoadingStats ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.statValue}>{userStats?.days_active || 0}</Text>
                <Text style={styles.statLabel}>Days Active</Text>
              </>
            )}
          </View>
        </MotiView>

        {userStats?.assigned_coach && (
          <MotiView
            from={{ opacity: 0, translateY: 20 }}
            animate={{ opacity: 1, translateY: 0 }}
            transition={{ type: 'timing', duration: 400, delay: 200 }}
            style={styles.section}
          >
            <Text style={styles.sectionTitle}>Your Coach</Text>
            <View style={styles.coachCard}>
              <View style={styles.coachHeader}>
                {userStats.assigned_coach.image ? (
                  <Image source={{ uri: userStats.assigned_coach.image }} style={styles.coachImage} />
                ) : (
                  <View style={[styles.coachImage, styles.placeholderCoachImage]}>
                    <Ionicons name="person" size={30} color={COLORS.text.tertiary} />
                  </View>
                )}
                <View style={styles.coachInfo}>
                  <Text style={styles.coachName}>{userStats.assigned_coach.name}</Text>
                  {userStats.assigned_coach.specialty && (
                    <Text style={styles.coachSpecialty}>{userStats.assigned_coach.specialty}</Text>
                  )}
                  <Text style={styles.coachAssignedDate}>
                    Assigned {new Date(userStats.assigned_coach.assigned_at).toLocaleDateString()}
                  </Text>
                </View>
              </View>
              <TouchableOpacity 
                style={styles.messageCoachButton}
                onPress={() => router.push('/(user)/conversations')}
              >
                <Ionicons name="chatbubble-outline" size={18} color={COLORS.primary} />
                <Text style={styles.messageCoachText}>Message Coach</Text>
              </TouchableOpacity>
            </View>
          </MotiView>
        )}

        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 300 }}
          style={styles.section}
        >
          <Text style={styles.sectionTitle}>About</Text>
          <Text style={styles.bio}>
            {currentUser?.bio || 'No bio added yet. Share something about your fitness journey!'}
          </Text>
          <TouchableOpacity 
            style={styles.editButton}
            onPress={() => setIsBioModalVisible(true)}
          >
            <Ionicons name="create-outline" size={18} color={COLORS.primary} />
            <Text style={styles.editButtonText}>Edit Bio</Text>
          </TouchableOpacity>
        </MotiView>

        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 300 }}
        >
          <MenuItemCard
            icon="fitness"
            title="My Goals"
            subtitle="Set and track your fitness goals"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="bar-chart"
            title="Progress"
            subtitle="View your progress over time"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="calendar"
            title="Workout History"
            subtitle="See your past workouts"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="nutrition"
            title="Nutrition"
            subtitle="Track your meals and calories"
            onPress={() => router.push('/(user)/nutrition')}
          />
          <MenuItemCard
            icon="heart"
            title="Mindfulness"
            subtitle="Meditation and breathing exercises"
            onPress={() => router.push('/(user)/mindfullness')}
          />
          <MenuItemCard
            icon="settings"
            title="Settings"
            subtitle="Manage your account settings"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="ribbon"
            title="Achievements"
            subtitle="Check out your achievements"
            onPress={() => router.push('/(user)/achievements')}
          />


        </MotiView>

        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 400 }}
          style={styles.logoutSection}
        >
          <LogoutButton />
        </MotiView>

        <View style={styles.footer}>
          <Text style={styles.footerText}>FitUp v1.0.0</Text>
          <Text style={styles.footerText}>Made with 💚 for fitness enthusiasts</Text>
        </View>
      </ScrollView>

      <EditBioModal
        visible={isBioModalVisible}
        onClose={() => setIsBioModalVisible(false)}
        onSave={handleBioSave}
        currentBio={currentUser?.bio || ''}
        isLoading={isUpdating}
      />

      <EditImageModal
        visible={isImageModalVisible}
        onClose={() => setIsImageModalVisible(false)}
        onImageSelected={handleImageSelect}
      />
    </View>
  );
}

interface MenuItemCardProps {
  icon: keyof typeof Ionicons.glyphMap;
  title: string;
  subtitle: string;
  onPress: () => void;
}

const MenuItemCard: React.FC<MenuItemCardProps> = ({ icon, title, subtitle, onPress }) => {
  return (
    <Pressable
      style={({ pressed }) => [styles.menuCard, pressed && styles.menuCardPressed]}
      onPress={onPress}
    >
      <View style={styles.menuIconContainer}>
        <Ionicons name={icon} size={24} color={COLORS.primary} />
      </View>
      <View style={styles.menuContent}>
        <Text style={styles.menuTitle}>{title}</Text>
        <Text style={styles.menuSubtitle}>{subtitle}</Text>
      </View>
      <Ionicons name="chevron-forward" size={20} color={COLORS.text.tertiary} />
    </Pressable>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  safeArea: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scrollContent: {
    paddingBottom: SPACING['4xl'],
  },
  loadingText: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.inverse,
  },
  header: {
    alignItems: 'center',
    paddingTop: SPACING['3xl'],
    paddingHorizontal: SPACING.xl,
    paddingBottom: SPACING.xl,
  },
  profileImageContainer: {
    position: 'relative',
    marginBottom: SPACING.base,
  },
  profileImage: {
    width: 120,
    height: 120,
    borderRadius: BORDER_RADIUS.full,
    borderWidth: 4,
    borderColor: COLORS.lightGray,
  },
  placeholderImage: {
    backgroundColor: COLORS.background.card,
    alignItems: 'center',
    justifyContent: 'center',
  },
  editImageButton: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    backgroundColor: COLORS.primaryDark,
    width: 36,
    height: 36,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: COLORS.lightGray,
  },
  name: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginTop: SPACING.sm,
  },
  username: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.placeholder,
    marginTop: SPACING.xs,
  },
  email: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: SPACING.xs,
  },
  statsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.lg,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    ...SHADOWS.base,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primaryDark,
  },
  statLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.inverse,
    marginTop: SPACING.xs,
  },
  statDivider: {
    width: 1,
    height: 40,
    backgroundColor: COLORS.border.subtle,
  },
  section: {
    backgroundColor: COLORS.background.card,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    ...SHADOWS.base,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginBottom: SPACING.sm,
  },
  bio: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    lineHeight: 24,
  },
  editButton: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: SPACING.base,
    paddingVertical: SPACING.sm,
  },
  editButtonText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.medium,
    marginLeft: SPACING.xs,
  },
  menuCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    ...SHADOWS.base,
  },
  menuCardPressed: {
    backgroundColor: COLORS.background.accent,
    transform: [{ scale: 0.98 }],
  },
  menuIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS['2xl'],
    backgroundColor: COLORS.primaryDark,
    alignItems: 'center',
    justifyContent: 'center',
  },
  menuContent: {
    flex: 1,
    marginLeft: SPACING.base,
  },
  menuTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
  },
  menuSubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  logoutSection: {
    marginHorizontal: SPACING.base,
    marginTop: SPACING['2xl'],
  },
  footer: {
    alignItems: 'center',
    marginTop: SPACING['2xl'],
    paddingTop: SPACING.lg,
  },
  footerText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: SPACING.xs,
  },
  coachCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    ...SHADOWS.base,
  },
  coachHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.base,
  },
  coachImage: {
    width: 64,
    height: 64,
    borderRadius: 32,
    marginRight: SPACING.base,
  },
  placeholderCoachImage: {
    backgroundColor: COLORS.background.accent,
    alignItems: 'center',
    justifyContent: 'center',
  },
  coachInfo: {
    flex: 1,
  },
  coachName: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginBottom: 4,
  },
  coachSpecialty: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.primary,
    marginBottom: 2,
  },
  coachAssignedDate: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
  },
  messageCoachButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.primaryDark,
    borderRadius: BORDER_RADIUS.xl,
    paddingVertical: SPACING.sm,
    paddingHorizontal: SPACING.base,
  },
  messageCoachText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primary,
    marginLeft: SPACING.xs,
  },
});
