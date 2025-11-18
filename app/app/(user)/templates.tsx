import React, { useState } from 'react';
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  TextInput,
  Platform,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useUserTemplates, usePublicTemplates } from '@/hooks/user/use-templates';
import { useDeleteTemplate, useCreateTemplate } from '@/hooks/user/use-template-mutations';
import { COLORS, FONT_SIZES, FONT_WEIGHTS, SPACING, BORDER_RADIUS } from '@/constants/theme';
import type { UserWorkoutTemplate, CreateUserTemplateRequest } from '@/types/workout-template';
import { TemplateDetailModal } from '@/components/dashboard/template-detail-modal';
import { TemplateFormModal } from '@/components/dashboard/template-form-modal';
import { useTemplateContext } from '@/context/template-context';

type TabType = 'user' | 'public';

export default function TemplatesScreen() {
  const [activeTab, setActiveTab] = useState<TabType>('user');
  const [searchTerm, setSearchTerm] = useState('');
  const [page] = useState(1);
  const [selectedTemplate, setSelectedTemplate] = useState<UserWorkoutTemplate | null>(null);
  const [isDetailModalVisible, setIsDetailModalVisible] = useState(false);
  const [isFormModalVisible, setIsFormModalVisible] = useState(false);

  const { setOnCreateTemplate } = useTemplateContext();

  const {
    data: userTemplatesData,
    isLoading: isLoadingUser,
    refetch: refetchUser,
    isRefetching: isRefetchingUser,
  } = useUserTemplates({ page, page_size: 20 });

  const {
    data: publicTemplatesData,
    isLoading: isLoadingPublic,
    refetch: refetchPublic,
    isRefetching: isRefetchingPublic,
  } = usePublicTemplates({ page, page_size: 20 });

  const deleteTemplateMutation = useDeleteTemplate();
  const createTemplateMutation = useCreateTemplate();

  React.useEffect(() => {
    if (activeTab === 'user') {
      setOnCreateTemplate(() => handleCreateTemplate);
    } else {
      setOnCreateTemplate(null);
    }
    
    return () => {
      setOnCreateTemplate(null);
    };
  }, [activeTab, setOnCreateTemplate]);

  const currentData = activeTab === 'user' ? userTemplatesData : publicTemplatesData;
  const isLoading = activeTab === 'user' ? isLoadingUser : isLoadingPublic;
  const isRefreshing = activeTab === 'user' ? isRefetchingUser : isRefetchingPublic;
  const refetch = activeTab === 'user' ? refetchUser : refetchPublic;

  const filteredTemplates = React.useMemo(() => {
    if (!currentData?.templates) return [];
    if (!searchTerm.trim()) return currentData.templates;

    const search = searchTerm.toLowerCase();
    return currentData.templates.filter(
      (template) =>
        template.name.toLowerCase().includes(search) ||
        template.description.toLowerCase().includes(search)
    );
  }, [currentData?.templates, searchTerm]);

  const handleTemplatePress = (template: UserWorkoutTemplate) => {
    setSelectedTemplate(template);
    setIsDetailModalVisible(true);
  };

  const handleCreateTemplate = () => {
    setIsFormModalVisible(true);
  };

  const handleEditTemplate = (template: UserWorkoutTemplate) => {
    setSelectedTemplate(template);
    setIsFormModalVisible(true);
  };

  const handleSaveTemplate = (data: CreateUserTemplateRequest) => {
    createTemplateMutation.mutate(data, {
      onSuccess: () => {
        setIsFormModalVisible(false);
        setSelectedTemplate(null);
        Alert.alert('Success', 'Template created successfully!');
      },
      onError: (error) => {
        Alert.alert('Error', `Failed to create template: ${error.message}`);
      },
    });
  };

  const handleDeleteTemplate = (templateId: string, templateName: string) => {
    const confirmDelete = () => {
      deleteTemplateMutation.mutate(templateId, {
        onSuccess: () => {
          Alert.alert('Success', `"${templateName}" has been deleted`);
        },
        onError: (error) => {
          Alert.alert('Error', `Failed to delete template: ${error.message}`);
        },
      });
    };

    if (Platform.OS === 'web') {
      if (confirm(`Are you sure you want to delete "${templateName}"?`)) {
        confirmDelete();
      }
    } else {
      Alert.alert(
        'Delete Template',
        `Are you sure you want to delete "${templateName}"?`,
        [
          { text: 'Cancel', style: 'cancel' },
          {
            text: 'Delete',
            style: 'destructive',
            onPress: confirmDelete,
          },
        ]
      );
    }
  };

  const renderTemplateItem = ({ item }: { item: UserWorkoutTemplate }) => (
    <TouchableOpacity
      style={styles.templateCard}
      onPress={() => handleTemplatePress(item)}
      activeOpacity={0.7}
    >
      <View style={styles.templateHeader}>
        <View style={styles.templateInfo}>
          <Text style={styles.templateName}>{item.name}</Text>
          {item.description && (
            <Text style={styles.templateDescription} numberOfLines={2}>
              {item.description}
            </Text>
          )}
          <View style={styles.templateMeta}>
            <Ionicons name="barbell-outline" size={14} color={COLORS.text.tertiary} />
            <Text style={styles.templateMetaText}>
              {item.exercises.length} exercise{item.exercises.length !== 1 ? 's' : ''}
            </Text>
            {item.is_public && (
              <>
                <Ionicons name="globe-outline" size={14} color={COLORS.primary} />
                <Text style={[styles.templateMetaText, { color: COLORS.primary }]}>Public</Text>
              </>
            )}
          </View>
        </View>

        {activeTab === 'user' && (
          <View style={styles.templateActions}>
            <TouchableOpacity
              onPress={(e) => {
                e.stopPropagation();
                handleEditTemplate(item);
              }}
              style={styles.actionButton}
            >
              <Ionicons name="create-outline" size={20} color={COLORS.primary} />
            </TouchableOpacity>
            <TouchableOpacity
              onPress={(e) => {
                e.stopPropagation();
                handleDeleteTemplate(item.template_id, item.name);
              }}
              style={styles.actionButton}
            >
              <Ionicons name="trash-outline" size={20} color={COLORS.error} />
            </TouchableOpacity>
          </View>
        )}
      </View>

      {/* Exercise preview */}
      {item.exercises.length > 0 && (
        <View style={styles.exercisePreview}>
          <Text style={styles.exercisePreviewTitle}>Exercises:</Text>
          {item.exercises.slice(0, 3).map((exercise, index) => (
            <Text key={index} style={styles.exercisePreviewText}>
              • {exercise.exercise_name} - {exercise.sets}x{exercise.target_reps}
            </Text>
          ))}
          {item.exercises.length > 3 && (
            <Text style={styles.exercisePreviewMore}>
              +{item.exercises.length - 3} more
            </Text>
          )}
        </View>
      )}
    </TouchableOpacity>
  );

  const renderEmptyState = () => (
    <View style={styles.emptyState}>
      <Ionicons
        name={activeTab === 'user' ? 'document-text-outline' : 'globe-outline'}
        size={64}
        color={COLORS.text.tertiary}
      />
      <Text style={styles.emptyStateTitle}>
        {activeTab === 'user' ? 'No Templates Yet' : 'No Public Templates'}
      </Text>
      <Text style={styles.emptyStateMessage}>
        {activeTab === 'user'
          ? 'Create your first workout template to get started'
          : 'Check back later for community templates'}
      </Text>
      {activeTab === 'user' && (
        <TouchableOpacity style={styles.createButton} onPress={handleCreateTemplate}>
          <Ionicons name="add-circle" size={20} color={COLORS.white} />
          <Text style={styles.createButtonText}>Create Template</Text>
        </TouchableOpacity>
      )}
    </View>
  );

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <View style={styles.titleContainer}>
          <Ionicons
            name="document-text"
            size={28}
            color={COLORS.primary}
          />
          <Text style={styles.headerTitle}>My Templates</Text>
        </View>
        <Text style={styles.headerSubtitle}>Create and manage your workout templates.</Text>
      </View>

      {/* Tabs */}
      <View style={styles.tabsContainer}>
        <TouchableOpacity
          style={[styles.tab, activeTab === 'user' && styles.activeTab]}
          onPress={() => setActiveTab('user')}
        >
          <Text style={[styles.tabText, activeTab === 'user' && styles.activeTabText]}>
            My Templates
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.tab, activeTab === 'public' && styles.activeTab]}
          onPress={() => setActiveTab('public')}
        >
          <Text style={[styles.tabText, activeTab === 'public' && styles.activeTabText]}>
            Public Templates
          </Text>
        </TouchableOpacity>
      </View>

      {/* Search Bar */}
      <View style={styles.searchContainer}>
        <Ionicons name="search" size={20} color={COLORS.text.tertiary} />
        <TextInput
          style={styles.searchInput}
          placeholder="Search templates..."
          placeholderTextColor={COLORS.text.tertiary}
          value={searchTerm}
          onChangeText={setSearchTerm}
        />
        {searchTerm.length > 0 && (
          <TouchableOpacity onPress={() => setSearchTerm('')}>
            <Ionicons name="close-circle" size={20} color={COLORS.text.tertiary} />
          </TouchableOpacity>
        )}
      </View>

      {/* Templates List */}
      {isLoading ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={COLORS.primary} />
          <Text style={styles.loadingText}>Loading templates...</Text>
        </View>
      ) : (
        <FlatList
          data={filteredTemplates}
          renderItem={renderTemplateItem}
          keyExtractor={(item) => item.template_id}
          contentContainerStyle={styles.listContent}
          ListEmptyComponent={renderEmptyState}
          refreshControl={
            <RefreshControl
              refreshing={isRefreshing}
              onRefresh={refetch}
              tintColor={COLORS.primary}
            />
          }
        />
      )}

      {/* Template Detail Modal */}
      <TemplateDetailModal
        visible={isDetailModalVisible}
        template={selectedTemplate}
        onClose={() => {
          setIsDetailModalVisible(false);
          setSelectedTemplate(null);
        }}
        onEdit={activeTab === 'user' ? handleEditTemplate : undefined}
        canEdit={activeTab === 'user'}
      />

      {/* Template Form Modal */}
      <TemplateFormModal
        visible={isFormModalVisible}
        onClose={() => {
          setIsFormModalVisible(false);
          if (!createTemplateMutation.isPending) {
            setSelectedTemplate(null);
          }
        }}
        onSave={handleSaveTemplate}
        initialData={selectedTemplate}
        isLoading={createTemplateMutation.isPending}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.auth,
    padding: SPACING.base,
  },
    header: {
    marginBottom: SPACING.lg,
  },
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 4,
  },
  title: {
    fontSize: 32,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  subtitle: {
    fontSize: 16,
    color: '#888888',
  },
  headerInfo: {
    flex: 1,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  headerSubtitle: {
    fontSize: 14,
    color: '#888888',
    marginTop: 2,
  },
  tabsContainer: {
    flexDirection: 'row',
    padding: SPACING.sm,
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS.full,
    alignItems: 'center',
    gap: SPACING.sm,
  },
  tab: {
    flex: 1,
    paddingVertical: SPACING.sm,
    alignItems: 'center',
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: COLORS.border.dark,
  },
  activeTab: {
    backgroundColor: COLORS.primarySoft,
  },
  tabText: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    fontWeight: FONT_WEIGHTS.medium as any,
  },
  activeTabText: {
    color: COLORS.primaryDark,
    fontWeight: FONT_WEIGHTS.bold as any,
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background.card,
    // marginHorizontal: SPACING.lg,
    marginVertical: SPACING.md,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.md,
    gap: SPACING.sm,
  },
  searchInput: {
    flex: 1,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.inverse,
  },
  listContent: {
    flexGrow: 1,
  },
  templateCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: SPACING.md,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  templateHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: SPACING.sm,
  },
  templateInfo: {
    flex: 1,
  },
  templateName: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.primarySoft,
    marginBottom: SPACING.xs,
  },
  templateDescription: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.placeholder,
    marginBottom: SPACING.sm,
  },
  templateMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
  },
  templateMetaText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginRight: SPACING.sm,
  },
  templateActions: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  actionButton: {
    padding: SPACING.xs,
  },
  exercisePreview: {
    marginTop: SPACING.sm,
    paddingTop: SPACING.sm,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.light,
  },
  exercisePreviewTitle: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.secondary,
    marginBottom: SPACING.xs,
  },
  exercisePreviewText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    marginBottom: 2,
  },
  exercisePreviewMore: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.medium as any,
    marginTop: SPACING.xs,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: SPACING.md,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
  },
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: SPACING.xl,
  },
  emptyStateTitle: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.inverse,
    marginTop: SPACING.md,
    marginBottom: SPACING.sm,
  },
  emptyStateMessage: {
    fontSize: FONT_SIZES.base,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginBottom: SPACING.lg,
  },
  createButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.full,
    gap: SPACING.sm,
  },
  createButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.white,
  },
});
