# Coach Features Improvement Plan

## Current Issues

### 1. **Assign Client Implementation**
**Problems:**
- Manual username input is error-prone and slow
- No user search/autocomplete functionality
- No validation if the username exists before attempting assignment
- Limited feedback on why assignment might fail
- Modal UI is basic and doesn't show client information before assignment
- No bulk assignment capability
- No invitation system for clients who aren't yet registered

### 2. **Client Management**
**Problems:**
- Limited client filtering (only search by name/email)
- No sorting options (by activity, progress, last workout, etc.)
- No bulk actions (assign multiple schemas, send messages to multiple clients)
- Client cards show limited actionable information
- No quick actions menu
- Missing client status indicators (active, inactive, needs attention)

### 3. **Schema/Plan Creation for Clients**
**Problems:**
- Creating a schema from the client list just redirects to form
- No template preview before applying to client
- No ability to clone/modify existing successful schemas
- Can't create variations of a template for different clients
- Missing bulk schema assignment

### 4. **Communication & Engagement**
**Problems:**
- No quick message functionality from client list
- No notification system for client milestones
- Missing automated check-ins
- No progress alerts for coaches
- Limited coach activity logging

---

## Recommended Improvements

### Phase 1: Enhanced Client Assignment (High Priority)

#### 1.1 User Search & Discovery
**Frontend:**
```typescript
// New component: ClientSearchModal.tsx
interface ClientSearchModalProps {
  visible: boolean;
  onClose: () => void;
  onAssign: (client: UserSearchResult) => void;
}

// Features:
- Real-time search with debouncing
- Display user avatar, name, email, current fitness level
- Show if user already has a coach
- Display user's goals and preferences
- Preview user's profile before assigning
```

**Backend:**
```go
// New endpoint: GET /api/v1/schema/coach/search-users?query={query}&limit={limit}
// Returns: Users matching search query with relevant profile info
// Filters: Exclude already assigned clients, show registration status
```

**Implementation:**
- Add user search endpoint that returns unassigned users
- Implement debounced search input (wait 300ms after typing stops)
- Show user cards with avatars and key info
- Add "Preview Profile" button before assignment
- Implement assignment confirmation dialog with user details

#### 1.2 Client Invitation System
**Frontend:**
```typescript
// New component: InviteClientModal.tsx
interface InviteClientData {
  email: string;
  firstName: string;
  lastName: string;
  message?: string;
}

// Features:
- Send invitation email to non-registered users
- Pre-assign coach relationship (activated when they register)
- Track invitation status (pending, accepted, expired)
- Resend invitations
```

**Backend:**
```go
// New endpoints:
// POST /api/v1/schema/coach/invitations
// GET /api/v1/schema/coach/invitations
// POST /api/v1/schema/coach/invitations/{id}/resend
// DELETE /api/v1/schema/coach/invitations/{id}

// New table: coach_invitations
CREATE TABLE coach_invitations (
    id SERIAL PRIMARY KEY,
    coach_id TEXT NOT NULL,
    email TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    invitation_token UUID NOT NULL UNIQUE,
    status TEXT NOT NULL, -- 'pending', 'accepted', 'expired', 'cancelled'
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    accepted_at TIMESTAMP,
    UNIQUE(coach_id, email)
);
```

---

### Phase 2: Advanced Client Management (High Priority)

#### 2.1 Smart Filtering & Sorting
**Frontend Updates:**
```typescript
// Enhanced clients.tsx
interface ClientFilters {
  status: 'all' | 'active' | 'inactive' | 'needs_attention';
  sortBy: 'name' | 'last_active' | 'progress' | 'streak' | 'join_date';
  sortOrder: 'asc' | 'desc';
  hasActiveSchema: boolean | null;
  fitnessLevel: FitnessLevel | null;
}

// Add filter chips UI
// Add sort dropdown
// Persist filters in URL params
```

**Backend:**
```go
// Update GET /api/v1/schema/coach/clients with query params:
// ?status=active&sort_by=last_active&sort_order=desc&has_schema=true&fitness_level=intermediate
```

#### 2.2 Client Status Indicators
**Add automatic status calculation:**
- 🟢 **Active**: Worked out in last 3 days
- 🟡 **Needs Attention**: No workout in 3-7 days
- 🔴 **Inactive**: No workout in 7+ days
- ⚠️ **No Schema**: Client has no active workout plan

#### 2.3 Quick Actions Menu
**Frontend:**
```typescript
// Add context menu to each client card
interface QuickAction {
  icon: string;
  label: string;
  action: () => void;
  variant?: 'default' | 'danger';
}

const quickActions: QuickAction[] = [
  { icon: 'calendar', label: 'Assign Schema', action: () => {...} },
  { icon: 'chatbubble', label: 'Send Message', action: () => {...} },
  { icon: 'stats-chart', label: 'View Progress', action: () => {...} },
  { icon: 'clipboard', label: 'View Workouts', action: () => {...} },
  { icon: 'notifications', label: 'Set Reminder', action: () => {...} },
  { icon: 'remove-circle', label: 'Remove Client', action: () => {...}, variant: 'danger' },
];
```

#### 2.4 Bulk Actions
**Frontend:**
```typescript
// Add selection mode to client list
interface BulkActionState {
  selectionMode: boolean;
  selectedClients: Set<string>;
}

// Bulk actions toolbar
const bulkActions = [
  'Assign Schema from Template',
  'Send Group Message',
  'Export Client Data',
  'Add to Group',
];
```

---

### Phase 3: Enhanced Schema Management (Medium Priority)

#### 3.1 Template Preview & Customization
**Frontend:**
```typescript
// New component: SchemaTemplatePreview.tsx
interface SchemaTemplatePreviewProps {
  template: WorkoutTemplate;
  clientId: string;
  onApply: (customizations: SchemaCustomization) => void;
}

interface SchemaCustomization {
  startDate: Date;
  adjustDifficulty: 'easier' | 'same' | 'harder';
  modifyDays: { [key: string]: boolean };
  notes: string;
}

// Features:
- Preview full template structure
- Adjust difficulty multiplier
- Enable/disable specific workout days
- Set custom start date
- Add personalized notes
```

#### 3.2 Schema Quick Apply
**Frontend:**
```typescript
// Add to client card dropdown
const handleQuickApplyTemplate = (clientId: string, templateId: number) => {
  // Show quick confirmation modal with template name
  // Apply with default settings
  // Show success toast
};
```

#### 3.3 Schema Cloning & Variations
**Frontend:**
```typescript
// New feature: Clone existing client's schema to another client
interface SchemaCloneRequest {
  sourceClientId: string;
  targetClientId: string;
  adjustments: {
    scaleDifficulty: number; // 0.8 - 1.2
    swapEquipment?: { from: EquipmentType; to: EquipmentType }[];
  };
}
```

**Backend:**
```go
// New endpoint: POST /api/v1/schema/coach/schemas/clone
// Copies schema structure with adjustments for client needs
```

---

### Phase 4: Communication Hub (Medium Priority)

#### 4.1 Quick Message System
**Frontend:**
```typescript
// New component: QuickMessageModal.tsx
interface QuickMessageProps {
  clientIds: string[];
  templates?: MessageTemplate[];
}

interface MessageTemplate {
  id: string;
  name: string;
  content: string;
  category: 'encouragement' | 'check_in' | 'reminder' | 'custom';
}

// Pre-built templates:
const messageTemplates = [
  { name: 'Weekly Check-in', content: 'Hey {firstName}, how did your week go? 💪' },
  { name: 'Missed Workout', content: 'Haven't seen you at the gym lately. Everything okay?' },
  { name: 'Great Progress', content: 'Amazing progress this week! Keep it up! 🔥' },
  { name: 'Schema Update', content: 'I've updated your workout plan. Check it out!' },
];
```

#### 4.2 Automated Notifications
**Backend:**
```go
// New service: coach_notifications_service.go

// Notification triggers:
- Client completes milestone (e.g., 10 workouts, 30-day streak)
- Client misses 2+ consecutive scheduled workouts
- Client achieves new PR
- Client completes full schema week
- Client requests help/feedback

// Notification delivery:
- In-app notifications
- Email digest (configurable)
- Push notifications (optional)
```

#### 4.3 Coach Activity Dashboard
**Frontend:**
```typescript
// Enhanced coach dashboard with:
interface CoachActivity {
  todaysScheduledWorkouts: WorkoutEvent[];
  clientsNeedingAttention: ClientAlert[];
  recentMessages: Message[];
  upcomingMilestones: Milestone[];
  weeklyStats: {
    totalWorkoutsCompleted: number;
    averageCompletionRate: number;
    newClients: number;
    messagesExchanged: number;
  };
}
```

---

### Phase 5: Advanced Features (Low Priority)

#### 5.1 Client Groups/Tags
```typescript
// Organize clients into groups
interface ClientGroup {
  id: string;
  name: string;
  color: string;
  clientIds: string[];
  description?: string;
}

// Examples: "Competition Prep", "Beginners", "Advanced Athletes"
```

#### 5.2 Progress Tracking & Reports
```typescript
// Generate progress reports
interface ProgressReport {
  clientId: string;
  period: { start: Date; end: Date };
  metrics: {
    workoutsCompleted: number;
    comparisonRate: number;
    strengthGains: ExerciseProgress[];
    consistency: number;
    goalsAchieved: Goal[];
  };
  exportFormats: ['pdf', 'csv', 'json'];
}
```

#### 5.3 AI-Powered Insights
```typescript
// Client recommendations
interface CoachInsight {
  clientId: string;
  type: 'progress' | 'regression' | 'plateau' | 'injury_risk' | 'opportunity';
  severity: 'low' | 'medium' | 'high';
  title: string;
  description: string;
  recommendations: string[];
  confidence: number;
}
```

---

## Implementation Priority

### Sprint 1 (Week 1-2) ⚡ HIGHEST PRIORITY - COMPLETED ✅
1. ✅ User search endpoint for client assignment - COMPLETED
2. ✅ Enhanced AssignClientModal with search functionality - COMPLETED
3. ✅ Client status indicators (Active/Inactive/Needs Attention/No Schema) - COMPLETED
4. ✅ Status filter chips with real-time counts - COMPLETED
5. ⏳ Client invitation system (database + API) - PENDING

### Sprint 2 (Week 3-4) 🔥 HIGH PRIORITY - IN PROGRESS
1. ✅ Client filtering and sorting - COMPLETED
   - Smart filtering by status (Active, Needs Attention, Inactive, No Schema)
   - Multi-criteria sorting (Name, Last Active, Join Date, Workouts, Streak)
   - Combined search + filter + sort
2. ✅ Quick actions menu on client cards - COMPLETED
   - Assign Schema, Send Message, View Progress, View Workouts, Remove Client
   - Confirmation dialogs for destructive actions
3. ⏳ Template preview before applying to client - PENDING
4. ⏳ Quick message functionality - PENDING

### Sprint 3 (Week 5-6) 📊 MEDIUM PRIORITY
1. ✅ Bulk client actions
2. ✅ Schema cloning feature
3. ✅ Message templates
4. ✅ Enhanced coach dashboard with activity feed

### Sprint 4 (Week 7-8) 🎯 NICE TO HAVE
1. ⏳ Automated coach notifications
2. ⏳ Client groups/tags
3. ⏳ Progress report generation
4. ⏳ AI-powered insights

---

## Technical Debt to Address

1. **Error Handling**: Add proper error types and user-friendly messages
2. **Loading States**: Improve loading indicators and skeleton screens
3. **Optimistic Updates**: Use React Query's optimistic updates for better UX
4. **Caching**: Implement proper cache invalidation strategies
5. **Real-time**: Consider WebSocket for real-time coach-client updates
6. **Mobile UX**: Optimize for mobile screens (bottom sheets, swipe actions)

---

## Key UX Principles

1. **Reduce Clicks**: Common actions should be 1-2 clicks away
2. **Show Context**: Display relevant info before requiring action
3. **Provide Feedback**: Always confirm success/failure clearly
4. **Enable Undo**: Allow coaches to reverse accidental actions
5. **Smart Defaults**: Pre-fill forms with sensible defaults
6. **Progressive Disclosure**: Show advanced options only when needed

---

## Success Metrics

- **Time to Assign Client**: Reduce from ~30s to ~5s
- **Client Engagement**: Increase coach-client messages by 50%
- **Schema Assignment**: Increase weekly schema assignments by 40%
- **Coach Satisfaction**: Target 4.5+ stars from coach feedback
- **Feature Adoption**: 80%+ coaches using quick actions within 2 weeks
