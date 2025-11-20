# Backend Fixes Summary

## Issues Fixed

### 1. Missing Assigned Coach Data in /auth/stats Endpoint

**Problem:** The `/auth/stats` endpoint was not returning `assigned_coach` data even though a coach was assigned to the client.

**Root Causes:** 
1. The SQL query in `getAssignedCoach()` function was using incorrect column names that don't exist in the `users` table:
   - `display_name` → should be `name`
   - `image_url` → should be `image`  
   - `specialty` → doesn't exist in users table
2. The query tried to join messages using non-existent `receiver_id` column
3. Messages are linked through `conversations` table, not directly between users

**Fixes Applied:**
- Updated `server/internal/auth/repository/store.go` line 576-614
- Changed SQL query to use correct column names: `u.name` and `u.image`
- Removed `specialty` field from the SELECT and Scan statements
- Fixed column name `created_at` to `assigned_at` for coach assignment date
- Fixed message counting to use `LATERAL` join through `conversations` table:
  ```sql
  LEFT JOIN LATERAL (
    SELECT COUNT(*) as total
    FROM conversations c
    INNER JOIN messages m ON m.conversation_id = c.conversation_id
    WHERE c.coach_id = u.id AND c.client_id = $1
  ) msg_count ON true
  ```
- Updated `server/internal/auth/types/stats.go` to make `Specialty` a pointer type (`*string`) for true optionality

**Result:** The `/auth/stats` endpoint will now properly return assigned coach information when a client has a coach assigned, allowing the coach card to display on the client profile page.

---

### 2. Coach-Created Schemas Not Appearing in Workout Source Selector

**Problem:** When a coach creates a training schema/plan for a client, it doesn't appear in the client's workout source selector.

**Root Causes:** 
1. The `GetWeeklySchemasByUserID()` function was querying `weekly_schemas` table directly using `user_id`, but:
   - The `user_id` column in `weekly_schemas` stores the `workout_profile_id` (integer)
   - The frontend passes `auth_user_id` from the current user (UUID string)
   - There was no join to map `auth_user_id` to `workout_profile_id`
2. Type mismatch: comparing integer `workout_profile_id` with text `auth_user_id` caused SQL error

**Fixes Applied:**
- Updated `server/internal/schema/repository/schema_repo.go` lines 87-134
- Added `INNER JOIN` with `workout_profiles` table to properly match:
  ```sql
  INNER JOIN workout_profiles wp ON wp.workout_profile_id = ws.user_id
  WHERE wp.auth_user_id = $1
  ```
- Updated both the main query and count query to use the join

**Result:** Clients will now see schemas/plans created for them by their coach in the workout source selector, allowing them to select and start coach-assigned workouts.

---

## Files Modified

1. `server/internal/auth/repository/store.go` - Fixed `getAssignedCoach()` SQL query (column names + message counting)
2. `server/internal/auth/types/stats.go` - Made `Specialty` field optional with pointer type
3. `server/internal/schema/repository/schema_repo.go` - Fixed `GetWeeklySchemasByUserID()` to join with workout_profiles and cast types

---

## Testing Recommendations

### Test 1: Assigned Coach Display
1. Assign a coach to a client user
2. Login as the client
3. Navigate to the profile page
4. Verify the "Your Coach" card appears with coach name, image, and assigned date
5. Verify the message count is accurate
6. Verify the "Message Coach" button works

### Test 2: Coach-Created Schemas
1. Login as a coach
2. Create a weekly schema/plan for a client
3. Logout and login as that client
4. Navigate to workout source selector
5. Verify the coach-created schema appears as an option
6. Verify you can select it and view the workouts

---

## Database Schema Notes

The fixes assume the following database structure:
- `users` table has columns: `id` (UUID), `name`, `image`
- `workout_profiles` table links `auth_user_id` (UUID) to `workout_profile_id` (integer)
- `coach_assignments` table links coaches to clients via `workout_profile_id`
- `weekly_schemas` table stores schemas with `user_id` = `workout_profile_id` (integer)
- `conversations` table has `coach_id` and `client_id` (both UUIDs)
- `messages` table has `conversation_id` linking to conversations (no direct receiver_id)

---

### 4. Fix Workout Completion for Coach Schemas
**Issue:**
- Users could not complete workouts from coach-created schemas.
- The backend returned `400 Bad Request: plan_id is required`.
- Coach schemas do not have a `plan_id` (which is specific to AI-generated plans).

**Fix:**
- Modified `server/internal/auth/handlers/stats.go` to remove the strict validation `if completion.PlanID <= 0`.
- The `SaveWorkoutCompletion` function in the repository does not use `PlanID` (it logs progress at the exercise level), so this validation was unnecessary and blocking valid use cases.

**Files Modified:**
- `server/internal/auth/handlers/stats.go`
