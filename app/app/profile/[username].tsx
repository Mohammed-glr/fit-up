import React, { useMemo, useState } from 'react';
import {
  ActivityIndicator,
  Image,
  ScrollView,
  StyleSheet,
  Text,
  View,
  TouchableOpacity,
  Share,
  StatusBar,
} from 'react-native';
import { useLocalSearchParams, useRouter, Stack } from 'expo-router';
import { MotiView } from 'moti';
import { Ionicons } from '@expo/vector-icons';

import { Button } from '@/components/forms';
import { useToastMethods } from '@/components/ui';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';
import { useAuth } from '@/context/auth-context';
import { useCreateConversation } from '@/hooks/message/use-conversation';
import { usePublicProfile } from '@/hooks/user/use-public-profile';
import {
  canMessage,
  formatRole,
  getInitial,
  roleRestrictionMessage,
} from '@/utils/conversation';

const PublicProfileScreen: React.FC = () => {
  const router = useRouter();
  const params = useLocalSearchParams<{ username?: string | string[] }>();
  const [isFollowing, setIsFollowing] = useState(false);
  const rawUsername = params.username;
  const username = Array.isArray(rawUsername) ? rawUsername[0] : rawUsername;

  const { user } = useAuth();
  const { showInfo, showError, showSuccess } = useToastMethods();

  const { data, isLoading, isError, refetch } = usePublicProfile(username);
  const createConversation = useCreateConversation();

  const profile = data ?? null;

  const isOwnProfile = useMemo(() => {
    if (!profile || !user) return false;
    return profile.id === user.id;
  }, [profile, user]);

  const roleMismatch = useMemo(() => {
    if (!profile || !user) return false;
    return !canMessage(user.role, profile.role);
  }, [profile, user]);

  const isMessageDisabled = !profile || createConversation.isPending || isOwnProfile || roleMismatch;

  const messageHelper = useMemo(() => {
    if (!user) return 'Sign in to send direct messages.';
    if (isOwnProfile) return 'This is your public profile preview.';
    if (roleMismatch) return roleRestrictionMessage(user.role);
    return undefined;
  }, [isOwnProfile, roleMismatch, user]);

  const handleMessage = async () => {
    if (!profile) return;
    if (!user) {
      showInfo('Log in to start a conversation.');
      router.push('/(auth)/login');
      return;
    }
    if (isOwnProfile) return;
    if (roleMismatch) {
      showInfo(roleRestrictionMessage(user.role));
      return;
    }
    const participants = user.role === 'coach'
      ? { coach_id: user.id, client_id: profile.id }
      : { coach_id: profile.id, client_id: user.id };
    try {
      const result = await createConversation.mutateAsync(participants);
      const conversationId = result.conversation.conversation_id;
      const targetRoute = user.role === 'coach' ? '/(coach)/chat' : '/(user)/chat';
      if (result.message) showInfo(result.message);
      else showSuccess('Conversation created successfully.');
      router.push({ pathname: targetRoute, params: { conversationId: String(conversationId) } });
    } catch (e) {
      console.error(e);
      showError('Failed to start conversation. Please try again.');
    }
  };

  const handleFollowToggle = () => {
    setIsFollowing(!isFollowing);
    showSuccess(isFollowing ? 'Unfollowed' : 'Followed');
  };

  const handleShare = async () => {
    if (!profile) return;
    try {
      await Share.share({
        message: `Check out ${profile.name || profile.username}'s profile on FitUp!`,
      });
    } catch (e) {
      console.error(e);
    }
  };

  if (!username) {
    return (
      <View style={styles.centeredContainer}>
        <Text style={styles.feedbackText}>No username provided.</Text>
      </View>
    );
  }

  if (isLoading) {
    return (
      <View style={styles.centeredContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (isError || !profile) {
    return (
      <View style={styles.centeredContainer}>
        <Text style={styles.feedbackText}>Unable to load this profile.</Text>
        <Button title="Try Again" variant="secondary" onPress={() => refetch()} style={styles.retryButton} />
      </View>
    );
  }

  const joinDate = new Date(profile.created_at).toLocaleDateString(undefined, { month: 'long', year: 'numeric' });

  return (
    <View style={styles.container}>
      <Stack.Screen
        options={{
          headerShown: true,
          title: `@${profile.username}`,
          headerStyle: {
            backgroundColor: COLORS.background.auth,
          },
          headerTintColor: COLORS.text.inverse,
          headerShadowVisible: false,
        }}
      />
      <StatusBar barStyle="light-content" backgroundColor={COLORS.background.auth} />
      
      <ScrollView 
        contentContainerStyle={styles.scrollContent} 
        showsVerticalScrollIndicator={false}
      >
        <MotiView
          from={{ opacity: 0, translateY: -20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400 }}
          style={styles.profileHeader}
        >
          <View style={styles.avatarSection}>
            <View style={styles.avatarWrapper}>
              {profile.image ? (
                <Image source={{ uri: profile.image }} style={styles.avatarImage} />
              ) : (
                <View style={styles.avatarFallback}>
                  <Text style={styles.avatarInitial}>{getInitial(profile.name || profile.username)}</Text>
                </View>
              )}
            </View>
            
            <View style={styles.roleBadge}>
              <Ionicons 
                name={profile.role === 'coach' ? 'trophy' : 'person'} 
                size={12} 
                color={COLORS.primary} 
              />
              <Text style={styles.roleText}>{formatRole(profile.role)}</Text>
            </View>
          </View>

          <Text style={styles.nameText}>{profile.name || profile.username}</Text>
          <Text style={styles.usernameText}>@{profile.username}</Text>
          
          <View style={styles.metaRow}>
            <Ionicons name="calendar-outline" size={14} color={COLORS.text.tertiary} />
            <Text style={styles.metaText}>Joined {joinDate}</Text>
          </View>
        </MotiView>

        {!isOwnProfile && (
          <MotiView
            from={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ type: 'timing', duration: 400, delay: 100 }}
            style={styles.actionSection}
          >
            <TouchableOpacity
              style={[styles.actionButton, styles.primaryButton]}
              onPress={handleMessage}
              disabled={isMessageDisabled}
            >
              <Ionicons 
                name="chatbubble-ellipses-outline" 
                size={20} 
                color={createConversation.isPending ? COLORS.text.tertiary : COLORS.white} 
              />
              <Text style={[styles.actionButtonText, styles.primaryButtonText]}>
                {createConversation.isPending ? 'Starting...' : 'Message'}
              </Text>
            </TouchableOpacity>

            {/* <TouchableOpacity
              style={[styles.actionButton, styles.secondaryButton]}
              onPress={handleFollowToggle}
            >
              <Ionicons 
                name={isFollowing ? 'checkmark-circle' : 'person-add-outline'} 
                size={20} 
                color={COLORS.primary} 
              />
              <Text style={[styles.actionButtonText, styles.secondaryButtonText]}>
                {isFollowing ? 'Following' : 'Follow'}
              </Text>
            </TouchableOpacity> */}

            <TouchableOpacity
              style={[styles.actionButton, styles.iconButton]}
              onPress={handleShare}
            >
              <Ionicons name="share-outline" size={20} color={COLORS.primary} />
            </TouchableOpacity>
          </MotiView>
        )}

        {/* Helper Message */}
        {messageHelper && (
          <MotiView
            from={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 200 }}
            style={styles.helperCard}
          >
            <Ionicons name="information-circle" size={16} color={COLORS.info} />
            <Text style={styles.helperText}>{messageHelper}</Text>
          </MotiView>
        )}

        {/* Bio Section */}
        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 200 }}
          style={styles.section}
        >
          <View style={styles.sectionHeader}>
            <Ionicons name="document-text-outline" size={20} color={COLORS.primary} />
            <Text style={styles.sectionTitle}>About</Text>
          </View>
          {profile.bio ? (
            <Text style={styles.bioText}>{profile.bio}</Text>
          ) : (
            <View style={styles.emptyState}>
              <Ionicons name="create-outline" size={32} color={COLORS.text.tertiary} />
              <Text style={styles.emptyText}>No bio added yet</Text>
            </View>
          )}
        </MotiView>

        {/* Stats Section (if you want to add stats later) */}
        <MotiView
          from={{ opacity: 0, translateY: 20 }}
          animate={{ opacity: 1, translateY: 0 }}
          transition={{ type: 'timing', duration: 400, delay: 300 }}
          style={styles.section}
        >
          <View style={styles.sectionHeader}>
            <Ionicons name="stats-chart-outline" size={20} color={COLORS.primary} />
            <Text style={styles.sectionTitle}>Activity</Text>
          </View>
          <View style={styles.statsGrid}>
            <View style={styles.statCard}>
              <Ionicons name="barbell" size={24} color={COLORS.primary} />
              <Text style={styles.statValue}>0</Text>
              <Text style={styles.statLabel}>Workouts</Text>
            </View>
            <View style={styles.statCard}>
              <Ionicons name="flame" size={24} color={COLORS.error} />
              <Text style={styles.statValue}>0</Text>
              <Text style={styles.statLabel}>Streak</Text>
            </View>
            <View style={styles.statCard}>
              <Ionicons name="trophy" size={24} color={COLORS.warning} />
              <Text style={styles.statValue}>0</Text>
              <Text style={styles.statLabel}>Achievements</Text>
            </View>
          </View>
        </MotiView>

        {isOwnProfile && (
          <MotiView
            from={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 400 }}
            style={styles.ownProfileNote}
          >
            <Ionicons name="eye-outline" size={16} color={COLORS.info} />
            <Text style={styles.ownProfileText}>This is how your profile appears to others</Text>
          </MotiView>
        )}
      </ScrollView>
    </View>
  );
};

export default PublicProfileScreen;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  centeredContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: SPACING['2xl'],
    backgroundColor: COLORS.background.auth,
  },
  scrollContent: {
    paddingBottom: SPACING['4xl'],
  },
  profileHeader: {
    alignItems: 'center',
    paddingHorizontal: SPACING.lg,
    paddingTop: SPACING.xl,
    paddingBottom: SPACING.lg,
  },
  avatarSection: {
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  avatarWrapper: {
    width: 100,
    height: 100,
    borderRadius: 50,
    borderWidth: 4,
    borderColor: COLORS.primaryDark,
    padding: 2,
    backgroundColor: COLORS.background.card,
    ...SHADOWS.base,
  },
  avatarImage: {
    width: '100%',
    height: '100%',
    borderRadius: 50,
  },
  avatarFallback: {
    width: '100%',
    height: '100%',
    borderRadius: 50,
    backgroundColor: COLORS.primaryDark,
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarInitial: {
    fontSize: 40,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  roleBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    backgroundColor: COLORS.primaryDark,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    marginTop: SPACING.sm,
  },
  roleText: {
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.bold,
    fontSize: FONT_SIZES.xs,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  nameText: {
    fontSize: 28,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    textAlign: 'center',
    marginTop: SPACING.sm,
  },
  usernameText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    marginTop: 4,
    fontWeight: FONT_WEIGHTS.medium,
  },
  metaRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    marginTop: SPACING.sm,
  },
  metaText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  actionSection: {
    flexDirection: 'row',
    paddingHorizontal: SPACING.lg,
    marginTop: SPACING.lg,
    gap: SPACING.sm,
  },
  actionButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.xs,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.xl,
    ...SHADOWS.sm,
  },
  primaryButton: {
    backgroundColor: COLORS.primaryDark,
  },
  secondaryButton: {
    backgroundColor: COLORS.darkGray,
    borderWidth: 1,
    borderColor: COLORS.primary + '30',
  },
  iconButton: {
    flex: 0,
    width: 48,
    backgroundColor: COLORS.darkGray,
    borderWidth: 1,
    borderColor: COLORS.primary + '30',
  },
  actionButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold,
  },
  primaryButtonText: {
    color: COLORS.white,
  },
  secondaryButtonText: {
    color: COLORS.primary,
  },
  helperCard: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    backgroundColor: COLORS.darkGray,
    marginHorizontal: SPACING.lg,
    marginTop: SPACING.md,
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS.xl,
    borderWidth: 1,
    borderColor: COLORS.info + '30',
  },
  helperText: {
    flex: 1,
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    lineHeight: 20,
  },
  section: {
    backgroundColor: COLORS.background.card,
    marginHorizontal: SPACING.lg,
    marginTop: SPACING.lg,
    padding: SPACING.lg,
    borderRadius: BORDER_RADIUS['2xl'],
    ...SHADOWS.base,
  },
  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    marginBottom: SPACING.md,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
  },
  bioText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.placeholder,
    lineHeight: 24,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: SPACING.xl,
  },
  emptyText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginTop: SPACING.sm,
    fontStyle: 'italic',
  },
  statsGrid: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  statCard: {
    flex: 1,
    alignItems: 'center',
    backgroundColor: COLORS.darkGray,
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS.xl,
  },
  statValue: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    marginTop: SPACING.xs,
  },
  statLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginTop: 2,
  },
  ownProfileNote: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    backgroundColor: COLORS.darkGray,
    marginHorizontal: SPACING.lg,
    marginTop: SPACING.lg,
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS.xl,
    borderWidth: 1,
    borderColor: COLORS.info + '30',
  },
  ownProfileText: {
    flex: 1,
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
  },
  feedbackText: {
    fontSize: FONT_SIZES.lg,
    color: COLORS.text.placeholder,
    textAlign: 'center',
    marginBottom: SPACING.lg,
  },
  retryButton: {
    marginTop: SPACING.sm,
    minWidth: 120,
  },
});
