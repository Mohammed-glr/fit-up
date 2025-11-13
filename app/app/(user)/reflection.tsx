import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  TextInput,
  ScrollView,
  FlatList,
  Alert,
} from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import {
  useReflectionPrompts,
  useReflectionResponses,
  useCreateReflectionResponse,
} from '@/hooks/mindfulness/use-mindfulness';
import { useMindfulnessContext } from '@/context/mindfulness-context';
import { BORDER_RADIUS, SPACING, COLORS } from '@/constants/theme';

export default function ReflectionScreen() {
  const { reflectionMode, setReflectionMode, isReflectionResponding, setIsReflectionResponding, isReflectionHistory, setIsReflectionHistory, setOnSaveReflection, setIsSavingReflection } = useMindfulnessContext();
  const [isResponding, setIsResponding] = useState(false);
  const [selectedPrompt, setSelectedPrompt] = useState<any>(null);
  const [responseText, setResponseText] = useState('');
  const [showHistory, setShowHistory] = useState(false);

  const { data: prompts = [], isLoading: promptsLoading } = useReflectionPrompts();
  const { data: responses = [], isLoading: responsesLoading } =
    useReflectionResponses(50);
  const createResponse = useCreateReflectionResponse();

  const todayPrompt = prompts.length > 0 ? prompts[0] : null;

  useEffect(() => {
    if (reflectionMode === 'history') {
      setShowHistory(true);
      setReflectionMode('main');
    }
  }, [reflectionMode]);

  useEffect(() => {
    setIsReflectionResponding(isResponding);
  }, [isResponding, setIsReflectionResponding]);

  useEffect(() => {
    setIsReflectionHistory(showHistory);
  }, [showHistory, setIsReflectionHistory]);

  useEffect(() => {
    setIsSavingReflection(createResponse.isPending);
  }, [createResponse.isPending, setIsSavingReflection]);

  const handleRespond = (prompt: any) => {
    setSelectedPrompt(prompt);
    setIsResponding(true);
  };

  const handleSave = async () => {
    if (!responseText.trim()) {
      Alert.alert('Error', 'Please write a response');
      return;
    }

    try {
      await createResponse.mutateAsync({
        prompt_id: selectedPrompt?.prompt_id,
        response_text: responseText,
      });

      setResponseText('');
      setIsResponding(false);
      setSelectedPrompt(null);
    } catch (error) {
      Alert.alert('Error', 'Failed to save response');
    }
  };

  useEffect(() => {
    if (isResponding) {
      setOnSaveReflection(() => handleSave);
    } else {
      setOnSaveReflection(() => undefined);
    }
  }, [isResponding, responseText, selectedPrompt, setOnSaveReflection]);

  const renderResponse = ({ item }: { item: any }) => {
    const date = new Date(item.created_at);
    const formattedDate = date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });

    const prompt = prompts.find((p) => p.prompt_id === item.prompt_id);

    return (
      <View style={styles.responseCard}>
        <Text style={styles.responseDate}>{formattedDate}</Text>
        {prompt && (
          <Text style={styles.responsePrompt}>{prompt.prompt_text}</Text>
        )}
        <Text style={styles.responseText}>{item.response_text}</Text>
      </View>
    );
  };

  if (isResponding) {
    return (
      <View style={styles.container}>
        <SafeAreaView style={styles.safeArea} edges={['bottom']}>
          <ScrollView
            style={styles.scrollView}
            contentContainerStyle={styles.scrollContent}
            showsVerticalScrollIndicator={false}
          >
            <View style={styles.promptCard}>
              <Text style={styles.promptIcon}>💭</Text>
              <Text style={styles.promptText}>
                {selectedPrompt?.prompt_text || 'Free-form reflection'}
              </Text>
            </View>

            <TextInput
              style={styles.textInput}
              value={responseText}
              onChangeText={setResponseText}
              placeholder="Share your thoughts..."
              placeholderTextColor="#555555"
              multiline
              autoFocus
              textAlignVertical="top"
            />

            <View style={styles.tipCard}>
              <Text style={styles.tipTitle}>💡 Reflection Tips</Text>
              <Text style={styles.tipText}>
                • Be honest and authentic{'\n'}
                • Take your time{'\n'}
                • There are no wrong answers{'\n'}
                • Focus on your feelings and experiences
              </Text>
            </View>
          </ScrollView>
        </SafeAreaView>
      </View>
    );
  }

  if (showHistory) {
    return (
      <View style={styles.container}>
        <SafeAreaView style={styles.safeArea} edges={['bottom']}>
          {responsesLoading ? (
            <View style={styles.loadingContainer}>
              <Text style={styles.loadingText}>Loading...</Text>
            </View>
          ) : responses.length === 0 ? (
            <View style={styles.emptyContainer}>
              <Text style={styles.emptyEmoji}>📝</Text>
              <Text style={styles.emptyTitle}>No Reflections Yet</Text>
              <Text style={styles.emptyText}>
                Start reflecting to see your history
              </Text>
            </View>
          ) : (
            <FlatList
              data={responses}
              renderItem={renderResponse}
              keyExtractor={(item) => item.response_id.toString()}
              contentContainerStyle={styles.listContent}
              showsVerticalScrollIndicator={false}
            />
          )}
        </SafeAreaView>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['bottom']}>
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          <View style={styles.mainContent}>
            {promptsLoading ? (
              <View style={styles.loadingContainer}>
                <Text style={styles.loadingText}>Loading...</Text>
              </View>
            ) : (
              <>
                {todayPrompt && (
                  <View style={styles.todaySection}>
                    <Text style={styles.todayLabel}>Today's Prompt</Text>
                    <View style={styles.todayPromptCard}>
                      <Text style={styles.todayPromptIcon}>💭</Text>
                      <Text style={styles.todayPromptText}>
                        {todayPrompt.prompt_text}
                      </Text>
                      <TouchableOpacity
                        style={styles.respondButton}
                        onPress={() => handleRespond(todayPrompt)}
                        activeOpacity={0.8}
                      >
                        <Text style={styles.respondButtonText}>
                          Reflect on This
                        </Text>
                      </TouchableOpacity>
                    </View>
                  </View>
                )}

                <View style={styles.freeFormSection}>
                  <Text style={styles.sectionTitle}>Free-form Reflection</Text>
                  <TouchableOpacity
                    style={styles.freeFormCard}
                    onPress={() => handleRespond(null)}
                    activeOpacity={0.7}
                  >
                    <Text style={styles.freeFormIcon}>✍️</Text>
                    <Text style={styles.freeFormTitle}>
                      Write Your Own Thoughts
                    </Text>
                    <Text style={styles.freeFormText}>
                      Express whatever is on your mind
                    </Text>
                  </TouchableOpacity>
                </View>

                {prompts.length > 1 && (
                  <View style={styles.morePromptsSection}>
                    <Text style={styles.sectionTitle}>More Prompts</Text>
                    {prompts.slice(1, 5).map((prompt) => (
                      <TouchableOpacity
                        key={prompt.prompt_id}
                        style={styles.promptListCard}
                        onPress={() => handleRespond(prompt)}
                        activeOpacity={0.7}
                      >
                        <Text style={styles.promptListText}>
                          {prompt.prompt_text}
                        </Text>
                        <Text style={styles.promptListArrow}>→</Text>
                      </TouchableOpacity>
                    ))}
                  </View>
                )}

                {responses.length > 0 && (
                  <View style={styles.recentSection}>
                    <View style={styles.recentHeader}>
                      <Text style={styles.sectionTitle}>Recent Reflections</Text>
                      <TouchableOpacity onPress={() => setShowHistory(true)}>
                        <Text style={styles.viewAllText}>View All</Text>
                      </TouchableOpacity>
                    </View>
                    {responses.slice(0, 3).map((item) => {
                      const date = new Date(item.created_at);
                      const formattedDate = date.toLocaleDateString('en-US', {
                        month: 'short',
                        day: 'numeric',
                      });

                      return (
                        <View key={item.response_id} style={styles.recentCard}>
                          <Text style={styles.recentDate}>{formattedDate}</Text>
                          <Text style={styles.recentText} numberOfLines={2}>
                            {item.response_text}
                          </Text>
                        </View>
                      );
                    })}
                  </View>
                )}
              </>
            )}
          </View>
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0A0A0A',
  },
  safeArea: {
    flex: 1,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 20,
    paddingVertical: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#1A1A1A',
  },
  closeButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
  },
  closeButtonText: {
    fontSize: 24,
    color: '#FFFFFF',
  },
  title: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  saveButton: {
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  saveButtonText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#6C63FF',
  },
  historyButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
  },
  historyButtonText: {
    fontSize: 24,
  },
  placeholder: {
    width: 40,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: 20,
  },
  mainContent: {
    flex: 1,
  },
  todaySection: {
    marginBottom: 32,
  },
  todayLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#6C63FF',
    marginBottom: 12,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  todayPromptCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: 20,
    padding: 24,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#2A2A2A',
  },
  todayPromptIcon: {
    fontSize: 48,
    marginBottom: 16,
  },
  todayPromptText: {
    fontSize: 20,
    fontWeight: '600',
    color: '#FFFFFF',
    textAlign: 'center',
    marginBottom: 24,
    lineHeight: 28,
  },
  respondButton: {
    backgroundColor: '#6C63FF',
    borderRadius: 16,
    paddingVertical: 16,
    paddingHorizontal: 32,
  },
  respondButtonText: {
    fontSize: 16,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  freeFormSection: {
    marginBottom: 32,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#FFFFFF',
    marginBottom: 12,
  },
  freeFormCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: 16,
    padding: 24,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#2A2A2A',
  },
  freeFormIcon: {
    fontSize: 40,
    marginBottom: 12,
  },
  freeFormTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 8,
  },
  freeFormText: {
    fontSize: 14,
    color: '#888888',
    textAlign: 'center',
  },
  morePromptsSection: {
    marginBottom: 32,
  },
  promptListCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: 16,
    padding: 20,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#2A2A2A',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  promptListText: {
    flex: 1,
    fontSize: 16,
    color: '#FFFFFF',
    marginRight: 16,
  },
  promptListArrow: {
    fontSize: 20,
    color: '#6C63FF',
  },
  recentSection: {
    marginBottom: 32,
  },
  recentHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  viewAllText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#6C63FF',
  },
  recentCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#2A2A2A',
  },
  recentDate: {
    fontSize: 12,
    color: '#888888',
    marginBottom: 8,
  },
  recentText: {
    fontSize: 14,
    color: '#FFFFFF',
    lineHeight: 20,
  },
  promptCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS['2xl'],
    padding: 24,
    alignItems: 'center',
    marginBottom: 24,
  },
  promptIcon: {
    fontSize: 48,
    marginBottom: 16,
  },
  promptText: {
    fontSize: 20,
    fontWeight: '600',
    color: '#FFFFFF',
    textAlign: 'center',
    lineHeight: 28,
  },
  textInput: {
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    color: '#FFFFFF',
    fontSize: 16,
    minHeight: 300,
    marginBottom: 24,
  },
  tipCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
  },
  tipTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 12,
  },
  tipText: {
    fontSize: 14,
    color: '#888888',
    lineHeight: 22,
  },
  listContent: {
    padding: 20,
  },
  responseCard: {
    backgroundColor: '#1A1A1A',
    borderRadius: 16,
    padding: 20,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#2A2A2A',
  },
  responseDate: {
    fontSize: 14,
    color: '#888888',
    marginBottom: 12,
  },
  responsePrompt: {
    fontSize: 16,
    fontWeight: '600',
    color: '#6C63FF',
    marginBottom: 12,
  },
  responseText: {
    fontSize: 16,
    color: '#FFFFFF',
    lineHeight: 24,
  },
  loadingContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 40,
  },
  loadingText: {
    fontSize: 16,
    color: '#888888',
  },
  emptyContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 40,
  },
  emptyEmoji: {
    fontSize: 64,
    marginBottom: 16,
  },
  emptyTitle: {
    fontSize: 24,
    fontWeight: '700',
    color: '#FFFFFF',
    marginBottom: 8,
  },
  emptyText: {
    fontSize: 16,
    color: '#888888',
    textAlign: 'center',
  },
});
