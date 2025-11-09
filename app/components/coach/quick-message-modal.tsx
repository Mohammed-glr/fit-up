import React, { useState, useMemo } from 'react';
import {
  Modal,
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING, BORDER_RADIUS, FONT_SIZES, FONT_WEIGHTS } from '@/constants/theme';
import {
  MESSAGE_TEMPLATES,
  TEMPLATE_CATEGORIES,
  getTemplatesByCategory,
  replaceTemplateVariables,
  getSuggestedVariableValues,
  type MessageTemplate,
  type MessageTemplateCategory,
} from '@/utils/message-templates';

interface QuickMessageModalProps {
  visible: boolean;
  onClose: () => void;
  onSend: (message: string) => Promise<void>;
  clientName?: string;
  clientData?: any;
}

export const QuickMessageModal: React.FC<QuickMessageModalProps> = ({
  visible,
  onClose,
  onSend,
  clientName,
  clientData,
}) => {
  const [selectedCategory, setSelectedCategory] = useState<MessageTemplateCategory | 'all'>('all');
  const [selectedTemplate, setSelectedTemplate] = useState<MessageTemplate | null>(null);
  const [message, setMessage] = useState('');
  const [variables, setVariables] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  // Get filtered templates
  const filteredTemplates = useMemo(() => {
    return getTemplatesByCategory(selectedCategory);
  }, [selectedCategory]);

  // Handle template selection
  const handleTemplateSelect = (template: MessageTemplate) => {
    setSelectedTemplate(template);
    
    // Auto-fill variables with suggested values
    const autoFilledVariables: Record<string, string> = {};
    template.variables.forEach(variable => {
      autoFilledVariables[variable] = getSuggestedVariableValues(variable, {
        ...clientData,
        name: clientName,
        first_name: clientName,
      });
    });
    
    setVariables(autoFilledVariables);
    
    // Generate preview message
    const previewMessage = replaceTemplateVariables(template.content, autoFilledVariables);
    setMessage(previewMessage);
  };

  // Handle custom message
  const handleCustomMessage = () => {
    setSelectedTemplate(null);
    setVariables({});
    setMessage('');
  };

  // Update variable
  const handleVariableChange = (key: string, value: string) => {
    const newVariables = { ...variables, [key]: value };
    setVariables(newVariables);
    
    // Update message preview
    if (selectedTemplate) {
      const previewMessage = replaceTemplateVariables(selectedTemplate.content, newVariables);
      setMessage(previewMessage);
    }
  };

  // Send message
  const handleSend = async () => {
    if (!message.trim()) {
      Alert.alert('Error', 'Please enter a message');
      return;
    }

    setLoading(true);
    try {
      await onSend(message.trim());
      
      // Reset state
      setSelectedCategory('all');
      setSelectedTemplate(null);
      setMessage('');
      setVariables({});
      
      Alert.alert('Success', 'Message sent successfully!');
      onClose();
    } catch (error: any) {
      Alert.alert('Error', error.message || 'Failed to send message');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      setSelectedCategory('all');
      setSelectedTemplate(null);
      setMessage('');
      setVariables({});
      onClose();
    }
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="slide"
      onRequestClose={handleClose}
    >
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.container}
      >
        <TouchableOpacity
          style={styles.overlay}
          activeOpacity={1}
          onPress={handleClose}
        />
        
        <View style={styles.modalContent}>
          {/* Header */}
          <View style={styles.header}>
            <View style={styles.headerLeft}>
              <View style={styles.iconContainer}>
                <Ionicons name="paper-plane" size={24} color={COLORS.primary} />
              </View>
              <Text style={styles.title}>Quick Message</Text>
            </View>
            <TouchableOpacity
              onPress={handleClose}
              style={styles.closeButton}
              disabled={loading}
            >
              <Ionicons name="close" size={24} color={COLORS.text.primary} />
            </TouchableOpacity>
          </View>

          <ScrollView
            style={styles.scrollView}
            showsVerticalScrollIndicator={false}
            keyboardShouldPersistTaps="handled"
          >
            {/* Category Filters */}
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              style={styles.categoryScroll}
              contentContainerStyle={styles.categoryContainer}
            >
              {TEMPLATE_CATEGORIES.map((category) => {
                const isActive = selectedCategory === category.key;
                return (
                  <TouchableOpacity
                    key={category.key}
                    style={[styles.categoryChip, isActive && styles.categoryChipActive]}
                    onPress={() => setSelectedCategory(category.key as any)}
                    activeOpacity={0.7}
                  >
                    <Ionicons
                      name={category.icon as any}
                      size={16}
                      color={isActive ? '#fff' : COLORS.text.tertiary}
                    />
                    <Text style={[styles.categoryLabel, isActive && styles.categoryLabelActive]}>
                      {category.label}
                    </Text>
                  </TouchableOpacity>
                );
              })}
            </ScrollView>

            {/* Template Selection */}
            {selectedCategory !== 'custom' && (
              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Choose a Template</Text>
                <View style={styles.templateGrid}>
                  {filteredTemplates.map((template) => {
                    const isSelected = selectedTemplate?.id === template.id;
                    return (
                      <TouchableOpacity
                        key={template.id}
                        style={[styles.templateCard, isSelected && styles.templateCardActive]}
                        onPress={() => handleTemplateSelect(template)}
                        activeOpacity={0.7}
                      >
                        <View style={styles.templateIcon}>
                          <Ionicons
                            name={template.icon as any}
                            size={20}
                            color={isSelected ? COLORS.primary : COLORS.text.secondary}
                          />
                        </View>
                        <Text
                          style={[styles.templateTitle, isSelected && styles.templateTitleActive]}
                          numberOfLines={2}
                        >
                          {template.title}
                        </Text>
                      </TouchableOpacity>
                    );
                  })}
                  
                  {/* Custom Message Option */}
                  <TouchableOpacity
                    style={[
                      styles.templateCard,
                      !selectedTemplate && message === '' && styles.templateCardActive,
                    ]}
                    onPress={handleCustomMessage}
                    activeOpacity={0.7}
                  >
                    <View style={styles.templateIcon}>
                      <Ionicons
                        name="create"
                        size={20}
                        color={!selectedTemplate && message === '' ? COLORS.primary : COLORS.text.secondary}
                      />
                    </View>
                    <Text
                      style={[
                        styles.templateTitle,
                        !selectedTemplate && message === '' && styles.templateTitleActive,
                      ]}
                      numberOfLines={2}
                    >
                      Custom Message
                    </Text>
                  </TouchableOpacity>
                </View>
              </View>
            )}

            {/* Variable Inputs (if template has variables) */}
            {selectedTemplate && selectedTemplate.variables.length > 0 && (
              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Customize</Text>
                {selectedTemplate.variables.map((variable) => (
                  <View key={variable} style={styles.variableInput}>
                    <Text style={styles.variableLabel}>
                      {variable.replace(/([A-Z])/g, ' $1').replace(/^./, (str) => str.toUpperCase())}
                    </Text>
                    <TextInput
                      style={styles.input}
                      value={variables[variable] || ''}
                      onChangeText={(value) => handleVariableChange(variable, value)}
                      placeholder={`Enter ${variable}`}
                      editable={!loading}
                    />
                  </View>
                ))}
              </View>
            )}

            {/* Message Preview/Edit */}
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Message</Text>
              <TextInput
                style={[styles.input, styles.messageInput]}
                value={message}
                onChangeText={setMessage}
                placeholder="Type your message here..."
                multiline
                numberOfLines={6}
                textAlignVertical="top"
                editable={!loading}
              />
              <Text style={styles.characterCount}>{message.length} characters</Text>
            </View>
          </ScrollView>

          {/* Actions */}
          <View style={styles.actions}>
            <TouchableOpacity
              style={[styles.button, styles.cancelButton]}
              onPress={handleClose}
              disabled={loading}
            >
              <Text style={styles.cancelButtonText}>Cancel</Text>
            </TouchableOpacity>
            
            <TouchableOpacity
              style={[styles.button, styles.sendButton, loading && styles.buttonDisabled]}
              onPress={handleSend}
              disabled={loading || !message.trim()}
            >
              {loading ? (
                <ActivityIndicator color="#fff" />
              ) : (
                <>
                  <Ionicons name="send" size={20} color="#fff" />
                  <Text style={styles.sendButtonText}>Send Message</Text>
                </>
              )}
            </TouchableOpacity>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'flex-end',
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  modalContent: {
    backgroundColor: COLORS.background.primary,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    maxHeight: '90%',
    paddingBottom: Platform.OS === 'ios' ? 20 : 0,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: SPACING.lg,
    paddingTop: SPACING.lg,
    paddingBottom: SPACING.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border.light,
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    flex: 1,
  },
  iconContainer: {
    width: 44,
    height: 44,
    borderRadius: BORDER_RADIUS.lg,
    backgroundColor: COLORS.primary + '15',
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: FONT_SIZES.xl,
    fontWeight: FONT_WEIGHTS.bold as any,
    color: COLORS.text.primary,
    flex: 1,
  },
  closeButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scrollView: {
    paddingHorizontal: SPACING.lg,
    paddingTop: SPACING.md,
  },
  categoryScroll: {
    marginBottom: SPACING.lg,
  },
  categoryContainer: {
    gap: SPACING.sm,
    paddingRight: SPACING.lg,
  },
  categoryChip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.xs,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: BORDER_RADIUS.full,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.light,
  },
  categoryChipActive: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  categoryLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.tertiary,
  },
  categoryLabelActive: {
    color: '#fff',
  },
  section: {
    marginBottom: SPACING.xl,
  },
  sectionTitle: {
    fontSize: FONT_SIZES.lg,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.primary,
    marginBottom: SPACING.md,
  },
  templateGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  templateCard: {
    width: '48%',
    padding: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.light,
    alignItems: 'center',
    gap: SPACING.xs,
  },
  templateCardActive: {
    borderColor: COLORS.primary,
    backgroundColor: COLORS.primary + '10',
  },
  templateIcon: {
    width: 40,
    height: 40,
    borderRadius: BORDER_RADIUS.md,
    backgroundColor: COLORS.background.secondary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  templateTitle: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
    textAlign: 'center',
  },
  templateTitleActive: {
    color: COLORS.primary,
    fontWeight: FONT_WEIGHTS.semibold as any,
  },
  variableInput: {
    marginBottom: SPACING.md,
  },
  variableLabel: {
    fontSize: FONT_SIZES.sm,
    fontWeight: FONT_WEIGHTS.medium as any,
    color: COLORS.text.secondary,
    marginBottom: SPACING.xs,
  },
  input: {
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.light,
    borderRadius: BORDER_RADIUS.lg,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    fontSize: FONT_SIZES.base,
    color: COLORS.text.primary,
  },
  messageInput: {
    minHeight: 120,
    paddingTop: SPACING.sm,
  },
  characterCount: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.tertiary,
    textAlign: 'right',
    marginTop: SPACING.xs,
  },
  actions: {
    flexDirection: 'row',
    gap: SPACING.sm,
    paddingHorizontal: SPACING.lg,
    paddingTop: SPACING.md,
    paddingBottom: SPACING.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.border.light,
  },
  button: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.md,
    borderRadius: BORDER_RADIUS.lg,
    gap: SPACING.xs,
  },
  cancelButton: {
    backgroundColor: COLORS.background.card,
    borderWidth: 1,
    borderColor: COLORS.border.light,
  },
  cancelButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: COLORS.text.primary,
  },
  sendButton: {
    backgroundColor: COLORS.primary,
  },
  sendButtonText: {
    fontSize: FONT_SIZES.base,
    fontWeight: FONT_WEIGHTS.semibold as any,
    color: '#fff',
  },
  buttonDisabled: {
    opacity: 0.6,
  },
});
