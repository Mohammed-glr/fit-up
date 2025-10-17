import React, { useMemo } from 'react';
import {
  ActivityIndicator,
  Image,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';

import { Button } from '@/components/forms';
import { useToastMethods } from '@/components/ui';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS } from '@/constants/theme';
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
  const rawUsername = params.username;
  const username = Array.isArray(rawUsername) ? rawUsername[0] : rawUsername;

  const { user } = useAuth();
  const { showInfo, showError, showSuccess } = useToastMethods();

  const { data, isLoading, isError, refetch } = usePublicProfile(username);
  const createConversation = useCreateConversation();

  const profile = data ?? null;

  const isOwnProfile = useMemo(() => {
    if (!profile || !user) {
      return false;
    }
    return profile.id === user.id;
  }, [profile, user]);

  const roleMismatch = useMemo(() => {
    if (!profile || !user) {
      return false;
    }
    return !canMessage(user.role, profile.role);
  }, [profile, user]);

  const isMessageDisabled =
    !profile || createConversation.isPending || isOwnProfile || roleMismatch;

  const messageHelper = useMemo(() => {
    if (!user) {
      return 'Sign in to send direct messages.';
    }

    if (isOwnProfile) {
      return 'You cannot start a conversation with yourself.';
    }

    if (roleMismatch) {
      return roleRestrictionMessage(user.role);
    }

    return undefined;
  }, [isOwnProfile, roleMismatch, user]);

  const handleMessage = async () => {
    if (!profile) {
      return;
    }

    if (!user) {
      showInfo('Log in to start a conversation.');
      router.push('/(auth)/login');
      return;
    }

    if (isOwnProfile) {
      showInfo('You cannot message yourself.');
      return;
    }

    if (roleMismatch) {
      showInfo(roleRestrictionMessage(user.role));
      return;
    }

    const participants =
      user.role === 'coach'
        ? { coach_id: user.id, client_id: profile.id }
        : { coach_id: profile.id, client_id: user.id };

    try {
      const result = await createConversation.mutateAsync(participants);
      const conversationId = result.conversation.conversation_id;
      const targetRoute = user.role === 'coach' ? '/(coach)/chat' : '/(user)/chat';

      if (result.message) {
        showInfo(result.message);
      } else {
        showSuccess('Conversation created successfully.');
      }

      router.push({
        pathname: targetRoute,
        params: { conversationId: String(conversationId) },
      });
    } catch (error) {
      console.error('Failed to open conversation from profile:', error);
      showError('Failed to start conversation. Please try again.');
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
        <Button
          title="Try Again"
          variant="secondary"
          onPress={() => refetch()}
          style={styles.retryButton}
        />
      </View>
    );
  }

  return (
    <ScrollView contentContainerStyle={styles.scrollContent} style={styles.screen}>
      <View style={styles.card}>
        <View style={styles.avatarWrapper}>
          {profile.image ? (
            <Image source={{ uri: profile.image }} style={styles.avatarImage} />
          ) : (
            <View style={styles.avatarFallback}>
              <Text style={styles.avatarInitial}>
                {getInitial(profile.name || profile.username)}
              </Text>
            </View>
          )}
        </View>

        <Text style={styles.nameText}>{profile.name || profile.username}</Text>
        <Text style={styles.usernameText}>@{profile.username}</Text>
        <Text style={styles.roleBadge}>{formatRole(profile.role)}</Text>

        {profile.bio ? (
          <Text style={styles.bioText}>{profile.bio}</Text>
        ) : (
          <Text style={styles.bioEmpty}>No bio provided yet.</Text>
        )}

        <Button
          title={createConversation.isPending ? 'Starting…' : 'Message'}
          variant="primary"
          onPress={handleMessage}
          disabled={isMessageDisabled}
          loading={createConversation.isPending}
          style={styles.messageButton}
        />

        {messageHelper ? (
          <Text style={styles.helperText}>{messageHelper}</Text>
        ) : null}
      </View>
    </ScrollView>
  );
};

export default PublicProfileScreen;

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: COLORS.background.primary,
  },
  scrollContent: {
    padding: SPACING['2xl'],
    paddingBottom: SPACING['4xl'],
    alignItems: 'center',
    justifyContent: 'center',
  },
  card: {
    width: '100%',
    maxWidth: 420,
    backgroundColor: COLORS.background.surface,
    borderRadius: BORDER_RADIUS['3xl'],
    padding: SPACING['2xl'],
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 12 },
    shadowOpacity: 0.16,
    shadowRadius: 18,
    elevation: 14,
  },
  avatarWrapper: {
    width: 96,
    height: 96,
    borderRadius: 48,
    marginBottom: SPACING.lg,
    overflow: 'hidden',
    backgroundColor: 'rgba(0,0,0,0.08)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarImage: {
    width: '100%',
    height: '100%',
  },
  avatarFallback: {
    width: '100%',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'rgba(0,0,0,0.1)',
  },
  avatarInitial: {
    fontSize: 36,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
  },
  nameText: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.primary,
    textAlign: 'center',
  },
  usernameText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.secondary,
    marginTop: SPACING.xs,
  },
  roleBadge: {
    marginTop: SPACING.sm,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.primary,
    color: COLORS.text.inverse,
    fontWeight: FONT_WEIGHTS.semibold,
    fontSize: FONT_SIZES.sm,
  },
  bioText: {
    marginTop: SPACING['2xl'],
    fontSize: FONT_SIZES.base,
    color: COLORS.text.primary,
    textAlign: 'center',
    lineHeight: 22,
  },
  bioEmpty: {
    marginTop: SPACING['2xl'],
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
  messageButton: {
    width: '100%',
    marginTop: SPACING['2xl'],
  },
  helperText: {
    marginTop: SPACING.sm,
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
  centeredContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: SPACING['2xl'],
    backgroundColor: COLORS.background.primary,
  },
  feedbackText: {
    fontSize: FONT_SIZES.lg,
    color: COLORS.text.primary,
    textAlign: 'center',
    marginBottom: SPACING.lg,
  },
  retryButton: {
    marginTop: SPACING.sm,
  },
});
