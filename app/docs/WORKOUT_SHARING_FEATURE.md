# Workout Sharing Feature

## Overview
The Workout Sharing feature allows users to share their completed workout sessions with their coach directly through the chat interface. Users can attach workout summaries to messages, which display comprehensive information about the workout including exercises, sets, volume, and personal records.

## Architecture

### Backend Components

#### 1. Types (`server/internal/schema/types/workout_sharing.go`)
- `WorkoutShareSummary`: Complete workout summary with exercises, stats, and PRs
- `ShareWorkoutRequest`: Request to share a workout
- `ShareWorkoutResponse`: Response after sharing
- `ShareExerciseDetail`: Detailed exercise information

#### 2. Repository (`server/internal/schema/repository/workout_sharing_repo.go`)
- `WorkoutSharingRepo` interface with `GetWorkoutShareSummary` method
- Complex SQL query that:
  - JOINs workout_sessions, workouts, exercise_performances, set_performances, exercises
  - Aggregates exercise data (volume, sets, reps)
  - Calculates total duration and workout stats
  - Identifies personal records by comparing to historical data
  - Returns formatted `WorkoutShareSummary`

#### 3. Handlers (`server/internal/schema/handlers/workout_sharing.go`)
- `handleGetWorkoutShareSummary`: GET /workout-sessions/{sessionId}/share-summary
  - Validates user authentication and session ownership
  - Returns workout summary for sharing
- `handleShareWorkout`: POST /workout-sessions/share
  - Processes share requests
  - Can send to coach, export as image, or copy as text

#### 4. Routes (`server/internal/schema/handlers/routes.go`)
- Routes registered under JWT authentication middleware:
  ```
  GET  /workout-sessions/{sessionId}/share-summary
  POST /workout-sessions/share
  ```

### Frontend Components

#### 1. Types (`app/types/workout-sharing.ts`)
- `WorkoutShareSummary`: Matches backend type
  - session_id, workout_title, completed_at
  - duration_minutes, total_exercises, total_sets, total_reps
  - total_volume_lbs, prs_achieved
  - exercises array with detailed information
- `ShareWorkoutRequest`, `ShareWorkoutResponse`
- `WorkoutShareExercise` with set details and PR indicators

#### 2. React Query Hooks (`app/hooks/workout/use-workout-sharing.ts`)
- `useWorkoutShareSummary(sessionId, enabled)`
  - Fetches workout summary from backend
  - Query key: `['workout-share-summary', sessionId]`
  - 5-minute stale time
  - Enabled only when sessionId provided
  
- `useShareWorkout()`
  - Mutation to share workout
  - Shows success/error alerts
  - POST to `/workout-sessions/share`

#### 3. Components

##### ShareWorkoutModal (`app/components/dashboard/share-workout-modal.tsx`)
- Modal for sharing completed workouts
- Features:
  - Workout summary preview (title, date, duration)
  - Stats cards (duration, exercises, sets, volume)
  - PR achievements section with badges
  - Exercise list with sets/reps/weight details
  - Share options:
    - Copy as Text (clipboard)
    - Share to Social Media
    - Send to Coach
    - Export as Image
- Uses `useWorkoutShareSummary` to fetch summary
- Uses `useShareWorkout` mutation to send to backend

##### WorkoutAttachmentPicker (`app/components/chat/workout-attachment-picker.tsx`)
- Modal to select recent workout sessions
- Features:
  - Lists recent completed workouts (last 30 days)
  - Shows: workout title, date, duration, exercises count, volume
  - Search/filter functionality (future enhancement)
  - Tap to select and close
  - Empty state for no workouts
  - Loading state with spinner
  - Error state with retry option
- Props:
  - `visible`: boolean
  - `onClose`: callback
  - `onSelectWorkout`: callback with sessionId
  - `recentWorkouts`: optional pre-loaded workouts

##### MessageComposer Enhancement (`app/components/chat/message-composer.tsx`)
- Added attachment button with dumbbell icon
- Button positioned before input field
- Optional `onAttachWorkout` prop
- Only shows when callback provided

##### ChatView Enhancement (`app/components/chat/chat-view.tsx`)
- Added `showWorkoutPicker` state
- Added `handleAttachWorkout` callback (opens picker)
- Added `handleSelectWorkout` callback:
  - Fetches workout summary via httpClient
  - Formats summary as text
  - Inserts into message input
  - Closes picker
- Added `formatWorkoutSummary` helper function:
  - Creates formatted text with emojis
  - Includes: title, date, duration, exercises, sets, volume
  - Shows PR count if any achieved
- Integrated WorkoutAttachmentPicker modal

## User Flow

### Sharing Workout in Chat
1. User opens chat with coach
2. User clicks attachment button (dumbbell icon) in message composer
3. WorkoutAttachmentPicker modal opens
4. User sees list of recent completed workouts
5. User selects a workout
6. Formatted workout summary is inserted into message input
7. User can edit the message or add additional text
8. User sends message to coach
9. Coach receives message with workout summary

### Share Modal (Alternative Flow)
1. User completes a workout
2. User opens workout summary screen
3. User clicks "Share" button
4. ShareWorkoutModal opens
5. User sees workout summary with stats and exercises
6. User chooses share option:
   - Copy Text: Copies to clipboard
   - Share to Social: Opens system share sheet
   - Send to Coach: Opens chat with pre-filled message
   - Export Image: Generates and saves workout image

## Data Format

### Workout Summary Text Format
```
🏋️ [Workout Title]
📅 [Date]
⏱️ Duration: [X] minutes
💪 [Y] exercises • [Z] sets
📊 Total Volume: [N] lbs
🏆 [M] Personal Records! (if any)
```

### Example
```
🏋️ Push Day - Chest & Triceps
📅 1/15/2024
⏱️ Duration: 65 minutes
💪 6 exercises • 24 sets
📊 Total Volume: 4250 lbs
🏆 2 Personal Records!
```

## Database Schema

### Required Tables
- `workout_sessions`: Stores completed workout sessions
  - session_id, user_id, workout_id
  - start_time, end_time, status
  - notes, mood_before, mood_after
  
- `exercise_performances`: Stores exercise completion data
  - performance_id, session_id, exercise_id
  - order, notes, completed
  
- `set_performances`: Stores individual set data
  - set_performance_id, performance_id
  - set_number, weight, reps, completed
  - rpe, rest_seconds
  
- `session_metrics`: Stores aggregate session metrics
  - session_id, total_duration_seconds
  - total_exercises, total_sets, total_reps
  - total_volume_lbs

## API Endpoints

### GET /workout-sessions/{sessionId}/share-summary
**Authentication**: Required (JWT)

**Response**: 200 OK
```json
{
  "session_id": 1,
  "workout_title": "Push Day - Chest & Triceps",
  "completed_at": "2024-01-15T10:30:00Z",
  "duration_minutes": 65,
  "total_exercises": 6,
  "total_sets": 24,
  "total_reps": 96,
  "total_volume_lbs": 4250,
  "prs_achieved": 2,
  "exercises": [
    {
      "exercise_name": "Bench Press",
      "sets_completed": 4,
      "total_reps": 32,
      "total_volume_lbs": 1600,
      "pr_achieved": true,
      "best_set": {
        "weight": 225,
        "reps": 8
      }
    }
  ]
}
```

### POST /workout-sessions/share
**Authentication**: Required (JWT)

**Request Body**:
```json
{
  "session_id": 1,
  "share_type": "coach",
  "message": "Great workout today!"
}
```

**Response**: 200 OK
```json
{
  "success": true,
  "message": "Workout shared successfully",
  "share_url": "https://..."
}
```

## Testing

### Manual Testing Steps
1. Complete a workout session
2. Open chat with coach
3. Click attachment button
4. Verify picker shows recent workouts
5. Select a workout
6. Verify formatted summary appears in message input
7. Send message
8. Verify coach receives message with summary

### Edge Cases
- No recent workouts (empty state)
- Network error fetching summary
- Workout with no PRs achieved
- Very long workout title (truncation)
- Workout with incomplete exercises

## Future Enhancements
1. **Image Export**: Generate visual workout summary as image
2. **Social Sharing**: Share to Instagram, Twitter, etc.
3. **Workout Comparison**: Compare current workout to previous sessions
4. **Advanced Filtering**: Filter workouts by date, type, muscle group
5. **Workout Analytics**: Show trends and progress over time
6. **Custom Templates**: Allow users to create custom share formats
7. **Coach Feedback**: Allow coach to comment on shared workouts

## Security Considerations
- All endpoints require JWT authentication
- Users can only access their own workout sessions
- Session ownership verified before returning summary
- Sensitive user data (notes, mood) excluded from shares (optional)

## Performance Optimizations
- Query results cached with React Query (5-minute stale time)
- Workout list limited to 30 recent sessions
- Lazy loading for exercise details
- Optimized SQL query with proper indexes on:
  - workout_sessions(user_id, session_id)
  - exercise_performances(session_id)
  - set_performances(performance_id)

## Implementation Status
✅ Backend types and structures
✅ Repository layer with SQL queries
✅ HTTP handlers with authentication
✅ Routes registration
✅ Frontend types matching backend
✅ React Query hooks
✅ ShareWorkoutModal component
✅ WorkoutAttachmentPicker component
✅ MessageComposer enhancement
✅ ChatView integration
✅ Type error fixes
✅ API client integration

## Related Files
- Backend:
  - `server/internal/schema/types/workout_sharing.go`
  - `server/internal/schema/repository/workout_sharing_repo.go`
  - `server/internal/schema/handlers/workout_sharing.go`
  - `server/internal/schema/handlers/routes.go`

- Frontend:
  - `app/types/workout-sharing.ts`
  - `app/hooks/workout/use-workout-sharing.ts`
  - `app/components/dashboard/share-workout-modal.tsx`
  - `app/components/chat/workout-attachment-picker.tsx`
  - `app/components/chat/message-composer.tsx`
  - `app/components/chat/chat-view.tsx`
