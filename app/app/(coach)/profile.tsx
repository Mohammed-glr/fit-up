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
} from 'react-native';
import { useCurrentUser } from '@/hooks/user/use-current-user';
import { useUpdateProfile } from '@/hooks/user/use-update-profile';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { MotiView } from 'moti';
import LogoutButton from '@/components/auth/logout-button';
import { Ionicons } from '@expo/vector-icons';
import { EditBioModal } from '@/components/profile/EditBioModal';
import { EditImageModal } from '@/components/profile/EditImageModal';
import { useToastMethods } from '@/components/ui/toast-provider';

export default function CoachProfileScreen() {
  const { data: currentUser, isLoading } = useCurrentUser();
  const { mutate: updateProfile, isPending: isUpdating } = useUpdateProfile();
  const { showSuccess, showError } = useToastMethods();
  
  const [activeTab, setActiveTab] = useState<'about' | 'stats' | 'settings'>('about');
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
        {/* Header Section with Coach Badge */}
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
            <View style={styles.coachBadge}>
              <Ionicons name="fitness" size={16} color={COLORS.white} />
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
        </MotiView>

        {/* Coach Stats Section */}
        <MotiView
          from={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ type: 'timing', duration: 400, delay: 100 }}
          style={styles.statsContainer}
        >
          <View style={styles.statItem}>
            <Text style={styles.statValue}>0</Text>
            <Text style={styles.statLabel}>Clients</Text>
          </View>
          <View style={styles.statDivider} />
          <View style={styles.statItem}>
            <Text style={styles.statValue}>0</Text>
            <Text style={styles.statLabel}>Programs</Text>
          </View>
          <View style={styles.statDivider} />
          <View style={styles.statItem}>
            <Text style={styles.statValue}>0</Text>
            <Text style={styles.statLabel}>Reviews</Text>
          </View>
        </MotiView>

        {/* Rating Section */}
        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 150 }}
          style={styles.ratingCard}
        >
          <View style={styles.ratingHeader}>
            <Text style={styles.ratingTitle}>Coach Rating</Text>
            <View style={styles.starsContainer}>
              {[1, 2, 3, 4, 5].map((star) => (
                <Ionicons
                  key={star}
                  name="star"
                  size={20}
                  color={star <= 0 ? COLORS.border.medium : COLORS.primary}
                />
              ))}
            </View>
          </View>
          <Text style={styles.ratingSubtitle}>No ratings yet</Text>
        </MotiView>

        {/* Bio Section */}
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

        {/* Specializations Section */}
        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 250 }}
          style={styles.section}
        >
          <Text style={styles.sectionTitle}>Specializations</Text>
          <View style={styles.tagsContainer}>
            <View style={styles.tag}>
              <Text style={styles.tagText}>Strength Training</Text>
            </View>
            <View style={styles.tag}>
              <Text style={styles.tagText}>Weight Loss</Text>
            </View>
            <View style={styles.tag}>
              <Text style={styles.tagText}>Nutrition</Text>
            </View>
          </View>
          <TouchableOpacity style={styles.editButton}>
            <Ionicons name="add-circle-outline" size={18} color={COLORS.primary} />
            <Text style={styles.editButtonText}>Add Specializations</Text>
          </TouchableOpacity>
        </MotiView>

        {/* Coach Menu Items */}
        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 300 }}
        >
          <MenuItemCard
            icon="people"
            title="My Clients"
            subtitle="Manage your client list"
            badge="0"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="document-text"
            title="Training Programs"
            subtitle="Create and manage programs"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="calendar"
            title="Schedule"
            subtitle="View and manage your schedule"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="bar-chart"
            title="Analytics"
            subtitle="Track your business performance"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="wallet"
            title="Earnings"
            subtitle="View your revenue and payments"
            onPress={() => {}}
          />
          <MenuItemCard
            icon="settings"
            title="Settings"
            subtitle="Manage your account settings"
            onPress={() => {}}
          />
        </MotiView>

        {/* Logout Section */}
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

      {/* Edit Modals */}
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
    backgroundColor: COLORS.background.primary,
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
    color: COLORS.text.secondary,
  },
  header: {
    alignItems: 'center',
    paddingTop: SPACING['3xl'],
    paddingHorizontal: SPACING.xl,
    paddingBottom: SPACING.xl,
    backgroundColor: COLORS.background.secondary,
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
    borderColor: COLORS.primary,
  },
  placeholderImage: {
    backgroundColor: COLORS.background.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  coachBadge: {
    position: 'absolute',
    top: 0,
    left: 0,
    backgroundColor: COLORS.primary,
    width: 32,
    height: 32,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: COLORS.white,
  },
  editImageButton: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    backgroundColor: COLORS.primary,
    width: 36,
    height: 36,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: COLORS.white,
  },
  name: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
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
    color: COLORS.primary,
    marginLeft: SPACING.xs,
  },
  username: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.text.secondary,
    marginTop: SPACING.sm,
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
    backgroundColor: COLORS.white,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.lg,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS.base,
    ...SHADOWS.base,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
  },
  statValue: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  statLabel: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.secondary,
    marginTop: SPACING.xs,
  },
  statDivider: {
    width: 1,
    height: 40,
    backgroundColor: COLORS.border.subtle,
  },
  ratingCard: {
    backgroundColor: COLORS.white,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS.base,
    ...SHADOWS.base,
  },
  ratingHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.xs,
  },
  ratingTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  starsContainer: {
    flexDirection: 'row',
    gap: 4,
  },
  ratingSubtitle: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  section: {
    backgroundColor: COLORS.white,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS.base,
    ...SHADOWS.base,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
    marginBottom: SPACING.sm,
  },
  bio: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.secondary,
    lineHeight: 24,
  },
  tagsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  tag: {
    backgroundColor: COLORS.background.accent,
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    borderWidth: 1,
    borderColor: COLORS.primary,
  },
  tagText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium,
    color: COLORS.primary,
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
    backgroundColor: COLORS.white,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS.base,
    ...SHADOWS.base,
  },
  menuCardPressed: {
    backgroundColor: COLORS.background.secondary,
    transform: [{ scale: 0.98 }],
  },
  menuIconContainer: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.base,
    backgroundColor: COLORS.background.accent,
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
    color: COLORS.text.primary,
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
