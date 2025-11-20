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
import { useCoachDashboard } from '@/hooks/schema/use-coach';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { MotiView } from 'moti';
import LogoutButton from '@/components/auth/logout-button';
import { Ionicons } from '@expo/vector-icons';
import { EditBioModal } from '@/components/profile/EditBioModal';
import { EditImageModal } from '@/components/profile/EditImageModal';
import { useToastMethods } from '@/components/ui/toast-provider';

export default function CoachProfileScreen() {
  const router = useRouter();
  const { data: currentUser, isLoading } = useCurrentUser();
  const { data: dashboard, isLoading: isLoadingDashboard } = useCoachDashboard();
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
          <ActivityIndicator size="large" color={COLORS.primary} />
          <Text style={styles.loadingText}>Loading your profile...</Text>
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
          style={styles.headerContainer}
        >
          <View style={styles.header}>
            <View style={styles.profileImageContainer}>
              {currentUser?.image ? (
                <Image source={{ uri: currentUser.image }} style={styles.profileImage} />
              ) : (
                <View style={[styles.profileImage, styles.placeholderImage]}>
                  <Ionicons name="person" size={60} color={COLORS.text.tertiary} />
                </View>
              )}
              <View style={styles.coachBadge}>
                <Ionicons name="fitness" size={18} color={COLORS.white} />
              </View>
              <TouchableOpacity
                style={styles.editImageButton}
                onPress={() => setIsImageModalVisible(true)}
              >
                <Ionicons name="camera" size={20} color={COLORS.white} />
              </TouchableOpacity>
            </View>

            <Text style={styles.name}>{currentUser?.name || 'Coach'}</Text>
            <View style={styles.roleContainer}>
              <Ionicons name="star" size={16} color={COLORS.primary} />
              <Text style={styles.roleText}>Professional Coach</Text>
            </View>
            <Text style={styles.username}>@{currentUser?.username || 'username'}</Text>
            <Text style={styles.email}>{currentUser?.email}</Text>
          </View>
        </MotiView>

        <MotiView
          from={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ type: 'timing', duration: 400, delay: 100 }}
          style={styles.statsGrid}
        >
          <View style={styles.statCard}>
            <View style={styles.statIconContainer}>
              <Ionicons name="people" size={24} color={COLORS.primary} />
            </View>
            {isLoadingDashboard ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.statValue}>{dashboard?.total_clients || 0}</Text>
                <Text style={styles.statLabel}>Clients</Text>
              </>
            )}
          </View>

          <View style={styles.statCard}>
            <View style={styles.statIconContainer}>
              <Ionicons name="document-text" size={24} color={COLORS.primary} />
            </View>
            {isLoadingDashboard ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.statValue}>{dashboard?.active_schemas || 0}</Text>
                <Text style={styles.statLabel}>Programs</Text>
              </>
            )}
          </View>

          <View style={styles.statCard}>
            <View style={styles.statIconContainer}>
              <Ionicons name="trophy" size={24} color={COLORS.primary} />
            </View>
            {isLoadingDashboard ? (
              <ActivityIndicator size="small" color={COLORS.primary} />
            ) : (
              <>
                <Text style={styles.statValue}>{Math.round((dashboard?.average_completion || 0) * 100)}%</Text>
                <Text style={styles.statLabel}>Success</Text>
              </>
            )}
          </View>
        </MotiView>

        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 200 }}
          style={styles.section}
        >
          <Text style={styles.sectionTitle}>Professional Bio</Text>
          <Text style={styles.bio}>
            {currentUser?.bio || 'Share your expertise, certifications, and coaching philosophy to attract clients!'}
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
          <Text style={styles.menuSectionTitle}>Quick Actions</Text>

          <MenuItemCard
            icon="people"
            title="My Clients"
            subtitle="Manage your client list"
            badge={dashboard?.total_clients?.toString()}
            onPress={() => router.push('/(coach)/clients')}
          />
          <MenuItemCard
            icon="document-text"
            title="Training Programs"
            subtitle="Create and manage programs"
            badge={dashboard?.active_schemas?.toString()}
            onPress={() => router.push('/(coach)/schema-templates')}
          />
          <MenuItemCard
            icon="chatbubbles"
            title="Messages"
            subtitle="Chat with your clients"
            onPress={() => router.push('/(coach)/conversations')}
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
          <Text style={styles.footerText}>FitUp Coach Edition v1.0.0</Text>
          <Text style={styles.footerText}>Empowering coaches to transform lives 💚</Text>
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
  badge?: string;
  onPress: () => void;
}

const MenuItemCard: React.FC<MenuItemCardProps> = ({ icon, title, subtitle, badge, onPress }) => {
  return (
    <Pressable
      style={({ pressed }) => [styles.menuCard, pressed && styles.menuCardPressed]}
      onPress={onPress}
    >
      <View style={styles.menuIconContainer}>
        <Ionicons name={icon} size={24} color={COLORS.primary} />
      </View>
      <View style={styles.menuContent}>
        <View style={styles.menuTitleRow}>
          <Text style={styles.menuTitle}>{title}</Text>
          {badge && (
            <View style={styles.badge}>
              <Text style={styles.badgeText}>{badge}</Text>
            </View>
          )}
        </View>
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
    marginTop: SPACING.base,
  },
  headerContainer: {
    marginBottom: SPACING.base,
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
    borderColor: COLORS.primaryDark,
  },
  placeholderImage: {
    backgroundColor: COLORS.background.card,
    alignItems: 'center',
    justifyContent: 'center',
  },
  coachBadge: {
    position: 'absolute',
    top: 0,
    left: 0,
    backgroundColor: COLORS.primary,
    width: 36,
    height: 36,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: COLORS.background.auth,
    ...SHADOWS.base,
  },
  editImageButton: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    backgroundColor: COLORS.primary,
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: COLORS.background.auth,
    ...SHADOWS.base,
  },
  name: {
    fontSize: FONT_SIZES['3xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginTop: SPACING.sm,
  },
  roleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.accent,
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    marginTop: SPACING.xs,
  },
  roleText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.primarySoft,
    marginLeft: SPACING.xs,
  },
  username: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.placeholder,
    marginTop: SPACING.sm,
  },
  email: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: SPACING.xs,
  },
  statsGrid: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    gap: SPACING.sm,
  },
  statCard: {
    flex: 1,
    backgroundColor: COLORS.background.card,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    alignItems: 'center',
    ...SHADOWS.base,
  },
  statIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.base,
    backgroundColor: COLORS.primaryDark,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: SPACING.sm,
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
    marginTop: SPACING.xs,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.inverse,
    marginTop: SPACING.xs,
    textAlign: 'center',
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
    fontSize: FONT_SIZES.xl,
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
  menuSectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    marginBottom: SPACING.xs,
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
    width: 52,
    height: 52,
    borderRadius: BORDER_RADIUS.base,
    backgroundColor: COLORS.primaryDark,
    alignItems: 'center',
    justifyContent: 'center',
  },
  menuContent: {
    flex: 1,
    marginLeft: SPACING.base,
  },
  menuTitleRow: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  menuTitle: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.inverse,
  },
  badge: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.sm,
    paddingVertical: 2,
    borderRadius: BORDER_RADIUS.full,
    marginLeft: SPACING.sm,
  },
  badgeText: {
    fontSize: FONT_SIZES.xs,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.white,
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
});
