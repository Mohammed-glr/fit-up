export type MessageTemplateCategory = 'check_in' | 'encouragement' | 'reminder' | 'workout' | 'milestone' | 'custom';

export interface MessageTemplate {
  id: string;
  category: MessageTemplateCategory;
  title: string;
  content: string;
  variables: string[]; // e.g., ['clientName', 'workoutName']
  icon: string;
}

export const MESSAGE_TEMPLATES: MessageTemplate[] = [
  // Check-in Templates
  {
    id: 'check_in_weekly',
    category: 'check_in',
    title: 'Weekly Check-in',
    content: 'Hey {{clientName}}! 👋 How was your week? I noticed you completed {{completedWorkouts}} workouts. How are you feeling about your progress?',
    variables: ['clientName', 'completedWorkouts'],
    icon: 'chatbubble-ellipses',
  },
  {
    id: 'check_in_missed',
    category: 'check_in',
    title: 'Missed Workouts Check-in',
    content: 'Hi {{clientName}}, I noticed you missed a few workouts this week. Is everything okay? Let me know if you need any adjustments to your plan. 💪',
    variables: ['clientName'],
    icon: 'help-circle',
  },
  {
    id: 'check_in_monthly',
    category: 'check_in',
    title: 'Monthly Progress Review',
    content: 'Hey {{clientName}}! 🎯 It\'s been a month since we started. Let\'s review your progress and set new goals. When\'s a good time to chat?',
    variables: ['clientName'],
    icon: 'calendar',
  },

  // Encouragement Templates
  {
    id: 'encouragement_streak',
    category: 'encouragement',
    title: 'Workout Streak Celebration',
    content: '🔥 Amazing job, {{clientName}}! You\'re on a {{streakDays}}-day workout streak! Keep up the fantastic work! 💪',
    variables: ['clientName', 'streakDays'],
    icon: 'flame',
  },
  {
    id: 'encouragement_progress',
    category: 'encouragement',
    title: 'Progress Recognition',
    content: '🌟 Great progress this week, {{clientName}}! I can see real improvements in your {{focusArea}}. You\'re crushing it!',
    variables: ['clientName', 'focusArea'],
    icon: 'trophy',
  },
  {
    id: 'encouragement_comeback',
    category: 'encouragement',
    title: 'Welcome Back',
    content: 'Welcome back, {{clientName}}! 🎉 Great to see you getting back on track. Remember, consistency is key. You\'ve got this!',
    variables: ['clientName'],
    icon: 'heart',
  },
  {
    id: 'encouragement_tough_day',
    category: 'encouragement',
    title: 'Tough Workout Motivation',
    content: 'Hey {{clientName}}, I know today\'s workout was challenging, but you powered through! That\'s what separates good from great. 💪🔥',
    variables: ['clientName'],
    icon: 'flash',
  },

  // Reminder Templates
  {
    id: 'reminder_workout',
    category: 'reminder',
    title: 'Workout Reminder',
    content: '⏰ Friendly reminder, {{clientName}}! You have {{workoutName}} scheduled for today. Let me know if you need any modifications!',
    variables: ['clientName', 'workoutName'],
    icon: 'alarm',
  },
  {
    id: 'reminder_rest_day',
    category: 'reminder',
    title: 'Rest Day Reminder',
    content: '🛌 Don\'t forget, {{clientName}} - today is your rest day! Recovery is just as important as training. Take it easy!',
    variables: ['clientName'],
    icon: 'bed',
  },
  {
    id: 'reminder_hydration',
    category: 'reminder',
    title: 'Hydration Reminder',
    content: '💧 Quick reminder, {{clientName}}: Stay hydrated today, especially before and after your workout!',
    variables: ['clientName'],
    icon: 'water',
  },

  // Workout-specific Templates
  {
    id: 'workout_new_plan',
    category: 'workout',
    title: 'New Workout Plan',
    content: '📋 Hey {{clientName}}! I\'ve created a new workout plan tailored to your goals. Check it out and let me know if you have any questions!',
    variables: ['clientName'],
    icon: 'document-text',
  },
  {
    id: 'workout_modification',
    category: 'workout',
    title: 'Workout Modification',
    content: 'Hi {{clientName}}, I\'ve made some adjustments to your workout based on your progress. The new plan focuses more on {{focusArea}}.',
    variables: ['clientName', 'focusArea'],
    icon: 'create',
  },
  {
    id: 'workout_form_check',
    category: 'workout',
    title: 'Form Check Request',
    content: 'Hey {{clientName}}! When you do {{exerciseName}} today, could you send me a quick form check video? Want to make sure you\'re getting the most out of it! 🎥',
    variables: ['clientName', 'exerciseName'],
    icon: 'videocam',
  },

  // Milestone Templates
  {
    id: 'milestone_pr',
    category: 'milestone',
    title: 'Personal Record',
    content: '🎉 Congratulations, {{clientName}}! You just hit a new PR on {{exerciseName}}! Your hard work is paying off! 💪🔥',
    variables: ['clientName', 'exerciseName'],
    icon: 'medal',
  },
  {
    id: 'milestone_goal',
    category: 'milestone',
    title: 'Goal Achievement',
    content: '🏆 Incredible, {{clientName}}! You\'ve reached your goal of {{goalDescription}}! Time to set a new challenge! 🚀',
    variables: ['clientName', 'goalDescription'],
    icon: 'ribbon',
  },
  {
    id: 'milestone_consistency',
    category: 'milestone',
    title: 'Consistency Milestone',
    content: '⭐ Amazing, {{clientName}}! You\'ve completed {{totalWorkouts}} workouts with me. Your dedication is inspiring! Keep it up!',
    variables: ['clientName', 'totalWorkouts'],
    icon: 'star',
  },
];

export const TEMPLATE_CATEGORIES = [
  { key: 'all', label: 'All Templates', icon: 'apps' },
  { key: 'check_in', label: 'Check-ins', icon: 'chatbubble-ellipses' },
  { key: 'encouragement', label: 'Encouragement', icon: 'heart' },
  { key: 'reminder', label: 'Reminders', icon: 'alarm' },
  { key: 'workout', label: 'Workout Updates', icon: 'barbell' },
  { key: 'milestone', label: 'Milestones', icon: 'trophy' },
  { key: 'custom', label: 'Custom Message', icon: 'create' },
] as const;

/**
 * Get templates by category
 */
export function getTemplatesByCategory(category: MessageTemplateCategory | 'all'): MessageTemplate[] {
  if (category === 'all') {
    return MESSAGE_TEMPLATES;
  }
  return MESSAGE_TEMPLATES.filter(template => template.category === category);
}

/**
 * Get a template by ID
 */
export function getTemplateById(id: string): MessageTemplate | undefined {
  return MESSAGE_TEMPLATES.find(template => template.id === id);
}

/**
 * Replace template variables with actual values
 */
export function replaceTemplateVariables(
  content: string,
  variables: Record<string, string | number>
): string {
  let result = content;
  
  Object.entries(variables).forEach(([key, value]) => {
    const regex = new RegExp(`{{${key}}}`, 'g');
    result = result.replace(regex, String(value));
  });
  
  return result;
}

/**
 * Extract variables from a template content
 */
export function extractTemplateVariables(content: string): string[] {
  const regex = /{{(\w+)}}/g;
  const variables: string[] = [];
  let match;
  
  while ((match = regex.exec(content)) !== null) {
    if (!variables.includes(match[1])) {
      variables.push(match[1]);
    }
  }
  
  return variables;
}

/**
 * Get suggested values for template variables based on client data
 */
export function getSuggestedVariableValues(
  variable: string,
  clientData?: any
): string {
  const suggestions: Record<string, string> = {
    clientName: clientData?.name || clientData?.first_name || 'there',
    completedWorkouts: clientData?.total_workouts?.toString() || '0',
    streakDays: clientData?.current_streak?.toString() || '0',
    focusArea: 'strength training',
    workoutName: clientData?.current_schema_name || 'your workout',
    exerciseName: 'squats',
    goalDescription: 'your fitness goal',
    totalWorkouts: clientData?.total_workouts?.toString() || '0',
  };
  
  return suggestions[variable] || `[${variable}]`;
}

/**
 * Validate if all required variables are provided
 */
export function validateTemplateVariables(
  template: MessageTemplate,
  providedVariables: Record<string, string | number>
): { valid: boolean; missing: string[] } {
  const missing = template.variables.filter(
    variable => !providedVariables[variable] || providedVariables[variable] === ''
  );
  
  return {
    valid: missing.length === 0,
    missing,
  };
}
