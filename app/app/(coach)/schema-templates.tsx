import React, { useState } from 'react';
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  Alert,
  TextInput,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';

import {
  useCoachTemplates,
  useDeleteTemplate,
} from '@/hooks/schema/use-coach';
import type { WorkoutTemplate } from '@/types/schema';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS } from '@/constants/theme';

export default function SchemaTemplatesScreen() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');

  const { data: templatesData, isLoading, refetch, isRefetching } = useCoachTemplates();
  const deleteTemplateMutation = useDeleteTemplate();

  const templates = templatesData?.templates || [];

  const filteredTemplates = templates.filter((template) => {
    const query = searchQuery.toLowerCase();
    return (
      template.name.toLowerCase().includes(query) ||
      (template.description && template.description.toLowerCase().includes(query))
    );
  });

  const handleDeleteTemplate = (templateId: number) => {
    Alert.alert(
      'Delete Template',
      'Permanently remove this workout template?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteTemplateMutation.mutateAsync(templateId);
              Alert.alert('Deleted', 'Template removed from your library.');
            } catch (error: any) {
              Alert.alert('Delete Failed', error?.message || 'Try again later.');
            }
          },
        },
      ]
    );
  };

  const handleCreateFromTemplate = (templateId: number) => {
    Alert.alert(
      'Assign Template',
      'Navigate to client list to assign this template to a user.',
      [{ text: 'OK' }]
    );
  };

  const renderTemplateCard = ({ item }: { item: WorkoutTemplate }) => {
    return (
      <View style={styles.templateCard}>
        <View style={styles.templateHeader}>
          <View style={styles.templateIcon}>
            <Ionicons name="document-text" size={24} color={COLORS.primary} />
          </View>
          <View style={styles.templateInfo}>
            <Text style={styles.templateName}>{item.name}</Text>
            {item.description && (
              <Text style={styles.templateDescription} numberOfLines={2}>
                {item.description}
              </Text>
            )}
          </View>
        </View>

        <View style={styles.templateMeta}>
          <View style={styles.metaItem}>
            <Text style={styles.metaLabel}>Days</Text>
            <Text style={styles.metaValue}>{item.days_per_week}/week</Text>
          </View>
          <View style={styles.metaItem}>
            <Text style={styles.metaLabel}>Level</Text>
            <Text style={styles.metaValue}>
              {item.min_level === item.max_level
                ? item.min_level
                : `${item.min_level}–${item.max_level}`}
            </Text>
          </View>
          <View style={styles.metaItem}>
            <Text style={styles.metaLabel}>Goals</Text>
            <Text style={styles.metaValue} numberOfLines={1}>
              {item.suitable_goals || 'All'}
            </Text>
          </View>
        </View>

        <View style={styles.templateActions}>
          <TouchableOpacity
            style={[styles.actionButton, styles.primaryAction]}
            onPress={() => handleCreateFromTemplate(item.template_id)}
          >
            <Ionicons name="person-add" size={18} color={COLORS.text.primary} />
            <Text style={styles.primaryActionText}>Use Template</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.actionButton, styles.secondaryAction]}
            onPress={() => handleDeleteTemplate(item.template_id)}
            disabled={deleteTemplateMutation.isPending}
          >
            <Ionicons name="trash-outline" size={18} color={COLORS.error} />
            <Text style={styles.secondaryActionText}>Remove</Text>
          </TouchableOpacity>
        </View>
      </View>
    );
  };

  

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
          <Ionicons name="arrow-back" size={22} color={COLORS.text.primary} />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Schema Templates</Text>
        <View style={styles.headerSpacer} />
      </View>

      <View style={styles.searchContainer}>
        <Ionicons
          name="search"
          size={20}
          color={COLORS.text.tertiary}
          style={styles.searchIcon}
        />
        <TextInput
          style={styles.searchInput}
          placeholder="Search templates..."
          placeholderTextColor={COLORS.text.placeholder}
          value={searchQuery}
          onChangeText={setSearchQuery}
        />
        {searchQuery.length > 0 && (
          <TouchableOpacity onPress={() => setSearchQuery('')}>
            <Ionicons name="close-circle" size={20} color={COLORS.text.tertiary} />
          </TouchableOpacity>
        )}
      </View>

      <View style={styles.statsBar}>
        <Text style={styles.statsText}>
          {filteredTemplates.length} {filteredTemplates.length === 1 ? 'template' : 'templates'}
        </Text>
      </View>

      <FlatList
        data={filteredTemplates}
        keyExtractor={(item) => item.template_id.toString()}
        renderItem={renderTemplateCard}
        contentContainerStyle={styles.listContent}
        refreshControl={
          <RefreshControl
            refreshing={isRefetching}
            onRefresh={refetch}
            tintColor={COLORS.primary}
            colors={[COLORS.primary]}
          />
        }
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Ionicons name="folder-open-outline" size={64} color={COLORS.text.tertiary} />
            <Text style={styles.emptyTitle}>No templates</Text>
            <Text style={styles.emptySubtitle}>
              {searchQuery
                ? 'Try a different search term'
                : 'Save your first schema as a template'}
            </Text>
          </View>
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingTop: SPACING['4xl'],
    paddingHorizontal: SPACING.lg,
    paddingBottom: SPACING.lg,
  },
  backButton: {
    width: 44,
    height: 44,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    alignItems: 'center',
    justifyContent: 'center',
    ...SHADOWS.sm,
  },
  headerTitle: {
    fontSize: FONT_SIZES['2xl'],
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
  },
  headerSpacer: {
    width: 44,
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    marginHorizontal: SPACING.lg,
    marginTop: SPACING.md,
    paddingHorizontal: SPACING.md,
    borderRadius: BORDER_RADIUS.md,
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  searchIcon: {
    marginRight: SPACING.sm,
  },
  searchInput: {
    flex: 1,
    paddingVertical: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.auth.primary,
  },
  statsBar: {
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
  },
  statsText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
  },
  listContent: {
    padding: SPACING.base,
  },
  templateCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.lg,
    padding: SPACING.base,
    marginBottom: SPACING.md,
    ...SHADOWS.sm,
  },
  templateHeader: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: SPACING.md,
  },
  templateIcon: {
    width: 48,
    height: 48,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  templateInfo: {
    flex: 1,
  },
  templateName: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    marginBottom: 4,
  },
  templateDescription: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    lineHeight: 18,
  },
  templateMeta: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: SPACING.md,
    paddingTop: SPACING.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.dark,
  },
  metaItem: {
    flex: 1,
    alignItems: 'center',
  },
  metaLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    marginBottom: 4,
  },
  metaValue: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.auth.primary,
    textTransform: 'capitalize',
  },
  templateActions: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  actionButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
    gap: SPACING.xs,
  },
  primaryAction: {
    backgroundColor: COLORS.primary,
  },
  secondaryAction: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  primaryActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.text.primary,
  },
  secondaryActionText: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold,
    color: COLORS.error,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING['5xl'],
  },
  emptyTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold,
    color: COLORS.text.auth.primary,
    marginTop: SPACING.base,
    marginBottom: SPACING.xs,
  },
  emptySubtitle: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
  },
});
