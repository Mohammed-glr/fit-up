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
  useGratitudeEntries,
  useCreateGratitudeEntry,
  useDeleteGratitudeEntry,
} from '@/hooks/mindfulness/use-mindfulness';
import { GRATITUDE_TAGS, MOOD_SCALE } from '@/types/mindfulness';
import { useMindfulnessContext } from '@/context/mindfulness-context';
import { BORDER_RADIUS, SPACING, COLORS } from '@/constants/theme';

export default function GratitudeJournalScreen() {
  const { gratitudeMode, setGratitudeMode, isGratitudeWriting, setIsGratitudeWriting, setOnSaveGratitude, setIsSavingGratitude } = useMindfulnessContext();
  const [isWriting, setIsWriting] = useState(false);
  const [entryText, setEntryText] = useState('');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [selectedMood, setSelectedMood] = useState<number | undefined>();

  const { data: entries = [], isLoading } = useGratitudeEntries(50);
  const createEntry = useCreateGratitudeEntry();
  const deleteEntry = useDeleteGratitudeEntry();

  useEffect(() => {
    if (gratitudeMode === 'create') {
      setIsWriting(true);
      setGratitudeMode('list');
    }
  }, [gratitudeMode]);

  useEffect(() => {
    setIsGratitudeWriting(isWriting);
  }, [isWriting, setIsGratitudeWriting]);

  useEffect(() => {
    setIsSavingGratitude(createEntry.isPending);
  }, [createEntry.isPending, setIsSavingGratitude]);

  const handleSave = async () => {
    if (!entryText.trim()) {
      Alert.alert('Error', 'Please write something');
      return;
    }

    try {
      await createEntry.mutateAsync({
        entry_text: entryText,
        tags: selectedTags.length > 0 ? selectedTags : undefined,
        mood: selectedMood,
      });

      setEntryText('');
      setSelectedTags([]);
      setSelectedMood(undefined);
      setIsWriting(false);
    } catch (error) {
      Alert.alert('Error', 'Failed to save entry');
    }
  };

  useEffect(() => {
    if (isWriting) {
      setOnSaveGratitude(() => handleSave);
    } else {
      setOnSaveGratitude(() => undefined);
    }
  }, [isWriting, entryText, selectedTags, selectedMood, setOnSaveGratitude]);

  const toggleTag = (tag: string) => {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]
    );
  };

  const handleDelete = (entryId: number) => {
    Alert.alert('Delete Entry', 'Are you sure you want to delete this entry?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: () => deleteEntry.mutate(entryId),
      },
    ]);
  };

  const renderEntry = ({ item }: { item: any }) => {
    const date = new Date(item.created_at);
    const formattedDate = date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });

    const mood = item.mood ? MOOD_SCALE.find((m: any) => m.value === item.mood) : null;

    return (
      <View style={styles.entryCard}>
        <View style={styles.entryHeader}>
          <Text style={styles.entryDate}>{formattedDate}</Text>
          <TouchableOpacity
            onPress={() => handleDelete(item.entry_id)}
            hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
          >
            <Text style={styles.deleteButton}>✕</Text>
          </TouchableOpacity>
        </View>

        <Text style={styles.entryText}>{item.entry_text}</Text>

        {mood && (
          <View style={styles.moodContainer}>
            <Text style={styles.moodEmoji}>{mood.emoji}</Text>
            <Text style={styles.moodLabel}>{mood.label}</Text>
          </View>
        )}

        {item.tags && item.tags.length > 0 && (
          <View style={styles.tagsContainer}>
            {item.tags.map((tag: string, index: number) => (
              <View key={index} style={styles.tag}>
                <Text style={styles.tagText}>{tag}</Text>
              </View>
            ))}
          </View>
        )}
      </View>
    );
  };

  if (isWriting) {
    return (
      <View style={styles.container}>
        <SafeAreaView style={styles.safeArea} edges={['bottom']}>
          <ScrollView
            style={styles.scrollView}
            contentContainerStyle={styles.scrollContent}
            showsVerticalScrollIndicator={false}
          >
            <Text style={styles.sectionTitle}>What are you grateful for?</Text>
            <TextInput
              style={styles.textInput}
              value={entryText}
              onChangeText={setEntryText}
              placeholder="Today, I'm grateful for..."
              placeholderTextColor="#555555"
              multiline
              autoFocus
              textAlignVertical="top"
            />

            <Text style={styles.sectionTitle}>How are you feeling?</Text>
            <View style={styles.moodSelector}>
              {MOOD_SCALE.map((mood: any) => (
                <TouchableOpacity
                  key={mood.value}
                  style={[
                    styles.moodOption,
                    selectedMood === mood.value && styles.moodOptionSelected,
                  ]}
                  onPress={() => setSelectedMood(mood.value)}
                  activeOpacity={0.7}
                >
                  <Text style={styles.moodOptionEmoji}>{mood.emoji}</Text>
                </TouchableOpacity>
              ))}
            </View>

            <Text style={styles.sectionTitle}>Tags (optional)</Text>
            <View style={styles.tagsSelector}>
              {GRATITUDE_TAGS.map((tag: string) => (
                <TouchableOpacity
                  key={tag}
                  style={[
                    styles.tagOption,
                    selectedTags.includes(tag) && styles.tagOptionSelected,
                  ]}
                  onPress={() => toggleTag(tag)}
                  activeOpacity={0.7}
                >
                  <Text
                    style={[
                      styles.tagOptionText,
                      selectedTags.includes(tag) && styles.tagOptionTextSelected,
                    ]}
                  >
                    {tag}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </ScrollView>
        </SafeAreaView>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['bottom']}>
        {isLoading ? (
          <View style={styles.loadingContainer}>
            <Text style={styles.loadingText}>Loading...</Text>
          </View>
        ) : !entries || entries.length === 0 ? (
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyEmoji}>✍️</Text>
            <Text style={styles.emptyTitle}>Start Your Gratitude Journey</Text>
            <Text style={styles.emptyText}>
              Reflect on the positive moments in your life
            </Text>
            <TouchableOpacity
              style={styles.emptyButton}
              onPress={() => setIsWriting(true)}
              activeOpacity={0.8}
            >
              <Text style={styles.emptyButtonText}>Write First Entry</Text>
            </TouchableOpacity>
          </View>
        ) : (
          <FlatList
            data={entries}
            renderItem={renderEntry}
            keyExtractor={(item) => item.entry_id.toString()}
            contentContainerStyle={styles.listContent}
            showsVerticalScrollIndicator={false}
          />
        )}
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
  newButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#6C63FF',
    borderRadius: BORDER_RADIUS.full,
  },
  newButtonText: {
    fontSize: 28,
    color: '#FFFFFF',
    marginTop: -2,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
    marginBottom: 12,
    marginTop: 16,
  },
  textInput: {
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    color: '#FFFFFF',
    fontSize: 16,
    minHeight: 200,
  },
  moodSelector: {
    flexDirection: 'row',
    gap: 12,
  },
  moodOption: {
    flex: 1,
    aspectRatio: 1,
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS['2xl'],
    alignItems: 'center',
    justifyContent: 'center',
  },
  moodOptionSelected: {
    borderColor: '#6C63FF',
    backgroundColor: '#1A1A2A',
  },
  moodOptionEmoji: {
    fontSize: 32,
  },
  tagsSelector: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  tagOption: {
    backgroundColor: '#1A1A1A',
    borderRadius: BORDER_RADIUS.full,
    paddingVertical: 10,
    paddingHorizontal: 16,
  },
  tagOptionSelected: {
    backgroundColor: '#1A1A2A',
    borderColor: '#6C63FF',
  },
  tagOptionText: {
    fontSize: 14,
    color: '#888888',
    fontWeight: '500',
  },
  tagOptionTextSelected: {
    color: '#6C63FF',
  },
  listContent: {
    padding: 20,
  },
  entryCard: {
    backgroundColor: COLORS.background.card,
    borderRadius: BORDER_RADIUS['2xl'],
    padding: SPACING.lg,
    marginBottom: 16,
  },
  entryHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  entryDate: {
    fontSize: 14,
    color: '#888888',
  },
  deleteButton: {
    fontSize: 20,
    color: '#666666',
  },
  entryText: {
    fontSize: 16,
    color: '#FFFFFF',
    lineHeight: 24,
    marginBottom: 12,
  },
  moodContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 12,
  },
  moodEmoji: {
    fontSize: 20,
  },
  moodLabel: {
    fontSize: 14,
    color: '#888888',
  },
  tagsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  tag: {
    backgroundColor: '#2A2A2A',
    borderRadius: 8,
    paddingVertical: 4,
    paddingHorizontal: 12,
  },
  tagText: {
    fontSize: 12,
    color: '#6C63FF',
    fontWeight: '500',
  },
  loadingContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
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
    marginBottom: 24,
  },
  emptyButton: {
    backgroundColor: '#6C63FF',
    borderRadius: 16,
    paddingVertical: 16,
    paddingHorizontal: 32,
  },
  emptyButtonText: {
    fontSize: 16,
    fontWeight: '700',
    color: '#FFFFFF',
  },
});
