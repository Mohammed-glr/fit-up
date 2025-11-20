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
import { LinearGradient } from 'expo-linear-gradient';
import { BlurView } from 'expo-blur';

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
          headerTransparent: true,
          headerTitle: '',
          headerTintColor: COLORS.white,
          headerRight: () => (
            <TouchableOpacity onPress={handleShare} style={styles.headerButton}>
              <Ionicons name="share-outline" size={24} color={COLORS.white} />
            </TouchableOpacity>
          ),
        }}
      />
      <StatusBar barStyle="light-content" />
      <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false} bounces={false}>
        <View style={styles.headerContainer}>
          <LinearGradient colors={[COLORS.primary, COLORS.primaryDark]} start={{ x: 0, y: 0 }} end={{ x: 1, y: 1 }} style={styles.headerGradient} />
          <View style={styles.headerPattern} />
        </View>
        <View style={styles.contentContainer}>
          <MotiView from={{ opacity: 0, translateY: 50 }} animate={{ opacity: 1, translateY: 0 }} transition={{ type: 'spring', damping: 15 }} style={styles.profileCardContainer}>
            <BlurView intensity={20} tint="light" style={styles.blurContainer}>
              <View style={styles.avatarContainer}>
                <View style={styles.avatarWrapper}>
                  {profile.image ? (
                    <Image source={{ uri: profile.image }} style={styles.avatarImage} />
                  ) : (
                    <View style={styles.avatarFallback}>
                      <Text style={styles.avatarInitial}>{getInitial(profile.name || profile.username)}</Text>
                    </View>
                  )}
                </View>
                <View style={styles.roleBadgeContainer}>
                  <LinearGradient colors={[COLORS.primarySoft, COLORS.primary]} start={{ x: 0, y: 0 }} end={{ x: 1, y: 0 }} style={styles.roleBadge}>
                    <Text style={styles.roleText}>{formatRole(profile.role)}</Text>
                  </LinearGradient>
                </View>
              </View>
              <View style={styles.infoContainer}>
                <Text style={styles.nameText}>{profile.name || profile.username}</Text>
                <Text style={styles.usernameText}>@{profile.username}</Text>
                <View style={styles.metaContainer}>
                  <Ionicons name="calendar-outline" size={14} color={COLORS.text.tertiary} />
                  <Text style={styles.metaText}>Member since {joinDate}</Text>
                </View>
                <View style={styles.divider} />
                <Text style={styles.sectionTitle}>About</Text>
                {profile.bio ? (
                  <Text style={styles.bioText}>{profile.bio}</Text>
                ) : (
                  <Text style={styles.bioEmpty}>No bio provided yet.</Text>
                )}
                <View style={styles.actionContainer}>
                  {!isOwnProfile && (
                    <Button
                      title={createConversation.isPending ? 'Starting…' : 'Message'}
                      variant="primary"
                      onPress={handleMessage}
                      disabled={isMessageDisabled}
                      loading={createConversation.isPending}
                      style={styles.messageButton}
                      icon={<Ionicons name="chatbubble-ellipses-outline" size={20} color={COLORS.white} />}
                    />
                  )}
                  {!isOwnProfile && (
                    <Button
                      title={isFollowing ? 'Following' : 'Follow'}
                      variant={isFollowing ? 'secondary' : 'primary'}
                      onPress={handleFollowToggle}
                      style={styles.followButton}
                    />
                  )}
                  {messageHelper && (
                    <MotiView from={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 300 }} style={styles.helperContainer}>
                      <Ionicons name="information-circle-outline" size={16} color={COLORS.text.tertiary} />
                      <Text style={styles.helperText}>{messageHelper}</Text>
                    </MotiView>
                  )}
                </View>
              </View>
            </BlurView>
          </MotiView>
        </View>
      </ScrollView>
    </View>
  );
};

export default PublicProfileScreen;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.primary,
  },
  centeredContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: SPACING['2xl'],
    backgroundColor: COLORS.background.primary,
  },
  scrollContent: {
    flexGrow: 1,
  },
  headerContainer: {
    height: 200,
    width: '100%',
    position: 'relative',
  },
  headerGradient: {
    ...StyleSheet.absoluteFillObject,
  },
  headerPattern: {
    ...StyleSheet.absoluteFillObject,
    opacity: 0.1,
    backgroundColor: '#000',
  },
  headerButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(0,0,0,0.2)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  contentContainer: {
    flex: 1,
    marginTop: -60,
    paddingHorizontal: SPACING.lg,
    paddingBottom: SPACING['4xl'],
  },
  profileCardContainer: {
    borderRadius: BORDER_RADIUS['3xl'],
    ...SHADOWS.lg,
    overflow: 'hidden',
    backgroundColor: 'transparent',
  },
  blurContainer: {
    padding: SPACING.xl,
    alignItems: 'center',
    backgroundColor: 'rgba(255,255,255,0.15)',
  },
  avatarContainer: {
    marginTop: -SPACING['4xl'],
    alignItems: 'center',
    marginBottom: SPACING.lg,
  },
  avatarWrapper: {
    width: 120,
    height: 120,
    borderRadius: 60,
    padding: 4,
    backgroundColor: COLORS.background.card,
    ...SHADOWS.modern,
  },
  avatarImage: {
    width: '100%',
    height: '100%',
    borderRadius: 60,
    backgroundColor: COLORS.background.card,
  },
  avatarFallback: {
    width: '100%',
    height: '100%',
    borderRadius: 60,
    backgroundColor: COLORS.primarySoft,
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarInitial: {
    fontSize: 48,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.primary,
  },
  roleBadgeContainer: {
    marginTop: -16,
  },
  roleBadge: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
  },
  roleText: {
    color: COLORS.white,
    fontWeight: FONT_WEIGHTS.bold,
    fontSize: FONT_SIZES.xs,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  infoContainer: {
    width: '100%',
    alignItems: 'center',
  },
  nameText: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    textAlign: 'center',
    marginBottom: 2,
  },
  usernameText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.primary,
    marginBottom: SPACING.sm,
    fontWeight: FONT_WEIGHTS.medium,
  },
  metaContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.lg,
    backgroundColor: COLORS.background.card,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
  },
  metaText: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginLeft: SPACING.xs,
  },
  divider: {
    width: '100%',
    height: 1,
    backgroundColor: COLORS.border.subtle,
    marginBottom: SPACING.lg,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.inverse,
    alignSelf: 'flex-start',
    marginBottom: SPACING.sm,
  },
  bioText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    lineHeight: 24,
    textAlign: 'left',
    alignSelf: 'stretch',
    marginBottom: SPACING['2xl'],
  },
  bioEmpty: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    fontStyle: 'italic',
    marginBottom: SPACING['2xl'],
    alignSelf: 'flex-start',
  },
  actionContainer: {
    width: '100%',
    alignItems: 'center',
  },
  messageButton: {
    width: '100%',
    marginBottom: SPACING.sm,
  },
  followButton: {
    width: '100%',
  },
  helperContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: SPACING.md,
    paddingHorizontal: SPACING.md,
  },
  helperText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginLeft: SPACING.xs,
    textAlign: 'center',
  },
  feedbackText: {
    fontSize: FONT_SIZES.lg,
    color: COLORS.text.secondary,
    textAlign: 'center',
    marginBottom: SPACING.lg,
  },
  retryButton: {
    marginTop: SPACING.sm,
    minWidth: 120,
  },
});
