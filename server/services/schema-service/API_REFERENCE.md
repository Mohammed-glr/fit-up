# FitUp Schema Service - API Quick Reference

## Base URL
```
http://localhost:8083/api/v1
```

## Health Check
```bash
GET /health
```

## 🏋️ Exercises

### List Exercises
```bash
GET /exercises?page=1&limit=20
```

### Search Exercises
```bash
GET /exercises/search?q=squat&page=1&limit=10
```

### Filter Exercises
```bash
GET /exercises/filter?muscle_group=chest&equipment=barbell&difficulty=intermediate
```

### Get Exercise by ID
```bash
GET /exercises/123
```

### Get Exercises by Muscle Group
```bash
GET /exercises/muscle-group/chest
```

### Get Exercises by Equipment
```bash
GET /exercises/equipment/dumbbell
```

### Get Recommended Exercises
```bash
GET /exercises/recommended/1?count=10
```

### Get Popular Exercises
```bash
GET /exercises/popular?limit=20
```

### Get Exercise Stats
```bash
GET /exercises/123/stats
```

## 💪 Workouts

### Get Workout
```bash
GET /workouts/5
```

### Get Workout with Exercises
```bash
GET /workouts/5/full
```

### Get Schema Workouts
```bash
GET /schemas/10/workouts
```

## 👤 Fitness Profiles

### Get User Fitness Profile
```bash
GET /fitness-profiles/1
```

### Create Fitness Assessment
```bash
POST /fitness-profiles/1/assessment
Content-Type: application/json

{
  "overall_level": "intermediate",
  "strength_level": "intermediate",
  "cardio_level": "beginner",
  "flexibility_level": "beginner",
  "assessment_data": {
    "push_ups": 25,
    "pull_ups": 5,
    "squat_max": 100
  }
}
```

### Update Fitness Level
```bash
PUT /fitness-profiles/1/level
Content-Type: application/json

{
  "level": "advanced"
}
```

### Create Fitness Goal
```bash
POST /fitness-profiles/1/goals
Content-Type: application/json

{
  "goal_type": "strength",
  "target_value": 150,
  "current_value": 100,
  "target_date": "2025-12-31T00:00:00Z",
  "unit": "kg"
}
```

### Get Active Goals
```bash
GET /fitness-profiles/1/goals/active
```

### Update Fitness Goals
```bash
PUT /fitness-profiles/1/goals
Content-Type: application/json

[
  {
    "goal_type": "strength",
    "target_value": 150,
    "current_value": 100,
    "target_date": "2025-12-31T00:00:00Z",
    "unit": "kg"
  }
]
```

### Estimate One Rep Max
```bash
POST /fitness-profiles/1/1rm/exercise/10
Content-Type: application/json

{
  "weight": 80,
  "reps": 5,
  "sets": 3,
  "rpe": 8
}
```

### Get 1RM History
```bash
GET /fitness-profiles/1/1rm/exercise/10/history
```

### Create Movement Assessment
```bash
POST /fitness-profiles/1/movement-assessment
Content-Type: application/json

{
  "movement_data": {
    "overhead_squat": "pass",
    "single_leg_balance": "fail"
  },
  "limitations": ["limited_ankle_mobility", "tight_hip_flexors"]
}
```

### Get Movement Limitations
```bash
GET /fitness-profiles/1/limitations
```

### Create Workout Profile
```bash
POST /workout-profiles
X-User-ID: auth-user-123
Content-Type: application/json

{
  "current_level": "intermediate",
  "primary_goal": "muscle_gain",
  "available_equipment": ["barbell", "dumbbell"],
  "workout_frequency": 4,
  "session_duration": 60
}
```

### Get Workout Profile
```bash
GET /workout-profiles
X-User-ID: auth-user-123
```

## 📅 Plan Generation

### Generate Workout Plan
```bash
POST /plans/generate
Content-Type: application/json

{
  "user_id": 1,
  "metadata": {
    "user_goals": ["muscle_gain"],
    "available_equipment": ["barbell", "dumbbell", "cable"],
    "fitness_level": "intermediate",
    "weekly_frequency": 4,
    "time_per_workout": 60
  }
}
```

### Get Active Plan
```bash
GET /plans/active/1
```

### Get Plan History
```bash
GET /plans/history/1?limit=10
```

### Track Plan Performance
```bash
POST /plans/5/performance
Content-Type: application/json

{
  "completion_rate": 0.85,
  "average_rpe": 7.5,
  "total_volume": 15000,
  "week_number": 2
}
```

### Get Plan Effectiveness
```bash
GET /plans/5/effectiveness
```

### Mark Plan for Regeneration
```bash
POST /plans/5/regenerate
Content-Type: application/json

{
  "reason": "user_feedback_too_difficult"
}
```

### Get Adaptation History
```bash
GET /plans/adaptations/1
```

### Get Current Week Schema
```bash
GET /schemas/current/1
```

### Create Weekly Schema from Template
```bash
POST /schemas/from-template
Content-Type: application/json

{
  "user_id": 1,
  "template_id": 3,
  "week_start": "2025-10-06T00:00:00Z"
}
```

## 🏃 Workout Sessions

### Start Workout Session
```bash
POST /sessions/start
Content-Type: application/json

{
  "user_id": 1,
  "workout_id": 5
}
```

### Complete Workout Session
```bash
POST /sessions/10/complete
Content-Type: application/json

{
  "exercises_completed": 6,
  "total_duration": 3600,
  "average_rpe": 7.5,
  "total_volume": 5000,
  "exercises": [
    {
      "exercise_id": 1,
      "sets_completed": 3,
      "best_set": {
        "weight": 100,
        "reps": 8,
        "rpe": 8
      }
    }
  ]
}
```

### Skip Workout
```bash
POST /sessions/skip
Content-Type: application/json

{
  "user_id": 1,
  "workout_id": 5,
  "reason": "injury"
}
```

### Log Exercise Performance
```bash
POST /sessions/10/exercises/15/log
Content-Type: application/json

{
  "sets_completed": 3,
  "rpe": 8,
  "best_set": {
    "weight": 100,
    "reps": 8,
    "rpe": 8
  }
}
```

### Get Active Session
```bash
GET /sessions/active/1
```

### Get Session History
```bash
GET /sessions/history/1?page=1&limit=20
```

### Get Session Metrics
```bash
GET /sessions/10/metrics
```

### Get Weekly Session Stats
```bash
GET /sessions/stats/1?week_start=2025-10-01
```

## 📊 Performance Analytics

### Calculate Strength Progression
```bash
GET /analytics/strength/1/exercise/10?timeframe=30
```

### Detect Performance Plateau
```bash
GET /analytics/plateau/1/exercise/10
```

### Predict Goal Achievement
```bash
GET /analytics/goals/1/5/prediction
```

### Calculate Training Volume
```bash
GET /analytics/volume/1?week_start=2025-10-01
```

### Track Intensity Progression
```bash
GET /analytics/intensity/1/exercise/10
```

### Get Optimal Training Load
```bash
GET /analytics/optimal-load/1
```

## 📝 Common Query Parameters

### Pagination
```
?page=1&limit=20
```

### Date Filtering
```
?week_start=2025-10-01
```

### Timeframe
```
?timeframe=30  (days)
```

### Count/Limit
```
?count=10
?limit=20
```

## 🔐 Authentication Headers

```
X-User-ID: auth-user-123
```

## 📦 Response Format

### Success Response
```json
{
  "data": { ... },
  "success": true
}
```

### Error Response
```json
{
  "error": "Error message here"
}
```

### Paginated Response
```json
{
  "data": [...],
  "total_count": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

## 🚀 Example Workflows

### Complete Workout Flow
```bash
# 1. Start session
POST /sessions/start
{"user_id": 1, "workout_id": 5}

# 2. Log each exercise
POST /sessions/10/exercises/15/log
{"sets_completed": 3, "best_set": {"weight": 100, "reps": 8, "rpe": 8}}

# 3. Complete session
POST /sessions/10/complete
{"exercises_completed": 6, "total_duration": 3600, "average_rpe": 7.5}

# 4. View metrics
GET /sessions/10/metrics
```

### Setup User Profile Flow
```bash
# 1. Create fitness assessment
POST /fitness-profiles/1/assessment
{"overall_level": "intermediate", ...}

# 2. Set goals
POST /fitness-profiles/1/goals
{"goal_type": "strength", "target_value": 150, ...}

# 3. Create workout profile
POST /workout-profiles
{"current_level": "intermediate", "primary_goal": "muscle_gain", ...}

# 4. Generate plan
POST /plans/generate
{"user_id": 1, "metadata": {...}}
```

### Track Progress Flow
```bash
# 1. Get strength progression
GET /analytics/strength/1/exercise/10?timeframe=30

# 2. Check for plateau
GET /analytics/plateau/1/exercise/10

# 3. Get goal prediction
GET /analytics/goals/1/5/prediction

# 4. Get training volume
GET /analytics/volume/1
```
