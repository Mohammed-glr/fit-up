# Fix for 404 Error on Workout Sharing Endpoint

## Problem
The `/workout-sessions/{sessionId}/share-summary` endpoint is returning a 404 error because the handler methods in the Go backend were not exported (they started with lowercase letters).

## Changes Made

### Backend Changes (Go)
1. **File**: `server/internal/schema/handlers/workout_sharing.go`
   - Changed `handleGetWorkoutShareSummary` → `HandleGetWorkoutShareSummary`
   - Changed `handleShareWorkout` → `HandleShareWorkout`
   - Updated route registrations to use the exported method names

2. **File**: `server/internal/schema/handlers/routes.go`
   - Updated route handlers to call `HandleGetWorkoutShareSummary` and `HandleShareWorkout`

### Frontend Changes
3. **File**: `app/components/chat/chat-view.tsx`
   - Added better error logging to show detailed error messages
   - Console logs now show the exact error response from the API

## Deployment Steps

### Option 1: Local Development (Docker Compose)
If you're running locally with Docker Compose:

```bash
cd server
docker-compose down
docker-compose up --build -d
```

### Option 2: Production Deployment (if using external server)
If the server is running on a remote host (`api.fitupp.nl`):

#### Step 1: Build the Docker image
```bash
cd server
docker build -t tdmdh/fitup-server:latest .
```

#### Step 2: Push to Docker Hub (if applicable)
```bash
docker push tdmdh/fitup-server:latest
```

#### Step 3: Pull and restart on production server
SSH into your production server and run:
```bash
cd /path/to/fit-up/server
docker-compose -f docker-compose.production.yml pull
docker-compose -f docker-compose.production.yml up -d --force-recreate api-server
```

### Option 3: Alternative - Restart without rebuild
If Docker Desktop is installed and running:
```powershell
cd server
docker-compose restart api-server
```

## Verification

After deploying, test the endpoint:

### 1. Check if the endpoint is accessible
```bash
curl -X GET "https://api.fitupp.nl/api/v1/workout-sessions/2/share-summary" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2. Check the logs
```bash
docker logs fitup-api-server -f
```

### 3. Test in the app
1. Open the chat view in the app
2. Click the attachment button (dumbbell icon)
3. Select a workout from the picker
4. Check the console logs for detailed error information

## Potential Additional Issues

### Issue 1: Workout session doesn't exist
If the session with ID 2 doesn't exist in the database, the endpoint will return 404.

**Solution**: Create a test workout session or use an existing session ID.

### Issue 2: Workout session belongs to different user
The endpoint validates that the workout session belongs to the authenticated user.

**Solution**: Ensure you're using a session ID that belongs to the logged-in user.

### Issue 3: Workout session not completed
The query filters for `status = 'completed'` workouts.

**Solution**: Ensure the workout session has been marked as completed.

## Database Check

To verify what workout sessions exist:

```sql
-- Connect to the database
docker exec -it fitup-postgres psql -U fitup -d fitup

-- Check workout sessions
SELECT 
    session_id, 
    user_id, 
    workout_id, 
    status, 
    start_time, 
    end_time 
FROM workout_sessions 
WHERE user_id = 'YOUR_USER_ID'
ORDER BY start_time DESC
LIMIT 10;

-- Exit psql
\q
```

## Mock Data Alternative

If no workout sessions exist, you can create mock data in the picker component temporarily:

**File**: `app/components/chat/chat-view.tsx`

The component already has mock workout data:
```typescript
const recentWorkouts = [
    {
        session_id: 1,
        workout_title: 'Push Day - Chest & Triceps',
        completed_at: new Date().toISOString(),
        duration_minutes: 65,
        total_exercises: 6,
        total_volume_lbs: 4250,
    },
    // ... more mock data
];
```

However, these mock sessions won't work with the API because they don't exist in the database.

## Next Steps

1. **Deploy the updated backend** using one of the options above
2. **Verify the endpoint** responds correctly
3. **Test the workout attachment** feature in the app
4. **Create real workout sessions** if none exist in the database

## Contact

If issues persist after deployment, check:
- Server logs: `docker logs fitup-api-server`
- Database connection: Verify the database has workout sessions
- Authentication: Ensure JWT token is valid and not expired
- Network: Verify the app can reach `api.fitupp.nl`
