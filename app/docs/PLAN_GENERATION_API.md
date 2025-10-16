# Plan Generation API Documentation

## Overview
The Plan Generation Service creates personalized workout plans based on user preferences, fitness level, available equipment, and goals. The system uses an adaptive algorithm that selects appropriate exercises, sets, reps, and rest periods tailored to each user.

## Complete Workflow

### 1. Generate a New Plan

**Endpoint:** `POST /api/v1/plans/generate`

**Request Body:**
```json
{
  "user_id": 123,
  "metadata": {
    "user_goals": ["muscle_gain"],
    "available_equipment": ["bodyweight", "dumbbell", "barbell"],
    "fitness_level": "intermediate",
    "weekly_frequency": 4,
    "time_per_workout": 60
  }
}
```

**User Goals Options:**
- `muscle_gain` - Build muscle mass
- `strength` - Increase maximum strength
- `fat_loss` - Reduce body fat
- `endurance` - Improve cardiovascular fitness
- `general_fitness` - Overall health and fitness

**Fitness Level Options:**
- `beginner` - New to fitness training
- `intermediate` - 6+ months of training experience
- `advanced` - 2+ years of consistent training

**Equipment Types:**
- `bodyweight` - No equipment needed
- `dumbbell` - Dumbbells available
- `barbell` - Barbell available
- `resistance_band` - Resistance bands
- `kettlebell` - Kettlebells
- `cable_machine` - Cable machines
- `gym_full` - Full gym access

**Response:**
```json
{
  "plan_id": 1,
  "user_id": 123,
  "week_start": "2025-10-01T00:00:00Z",
  "generated_at": "2025-10-01T10:30:00Z",
  "algorithm": "fitup_adaptive_v1",
  "is_active": true,
  "metadata": {
    "parameters": {
      "template_used": "intermediate_muscle_gain_4day",
      "total_exercises": 24,
      "muscle_groups_targeted": ["chest", "back", "shoulders", "legs", "arms"],
      "equipment_utilized": ["bodyweight", "dumbbell", "barbell"],
      "estimated_volume": 96,
      "progression_method": "progressive_overload",
      "generated_plan": [
        {
          "day_title": "Upper Body Push",
          "focus": "chest_shoulders_triceps",
          "exercises": [
            {
              "exercise_id": 1,
              "name": "Barbell Bench Press",
              "sets": 4,
              "reps": "8-12",
              "rest": 120,
              "notes": "Focus on controlled eccentric phase"
            },
            {
              "exercise_id": 2,
              "name": "Dumbbell Shoulder Press",
              "sets": 3,
              "reps": "10-12",
              "rest": 90
            }
          ]
        },
        {
          "day_title": "Lower Body",
          "focus": "legs_glutes",
          "exercises": [...]
        },
        {
          "day_title": "Rest Day",
          "focus": "recovery",
          "exercises": []
        }
      ]
    }
  }
}
```

### 2. Get Active Plan

**Endpoint:** `GET /api/v1/plans/active/{userID}`

**Response:** Returns the user's currently active workout plan with the same structure as above.

### 3. Download Plan as PDF

**Endpoint:** `GET /api/v1/plans/{planID}/download`

**Response:** Binary PDF file

**Headers:**
- `Content-Type: application/pdf`
- `Content-Disposition: attachment; filename=workout_plan_{planID}.pdf`

**PDF Contents:**
1. **Cover Page**
   - FIT-UP branding
   - Plan title
   - Generation date

2. **Plan Information**
   - Plan ID
   - Generated date/time
   - Week start date
   - Algorithm used
   - Status (Active/Inactive)

3. **Training Parameters**
   - Template used
   - Total exercises
   - Muscle groups targeted
   - Equipment utilized
   - Estimated weekly volume

4. **Weekly Training Schedule**
   - Each day with:
     - Day title and focus area
     - Exercise table (name, sets, reps, rest)
     - Exercise notes and tips
     - Estimated duration per workout

5. **Training Guidelines & Tips**
   - Warm-up protocol
   - Proper form guidance
   - Rest period recommendations
   - Progressive overload principles
   - Recovery strategies
   - Nutrition tips
   - Injury prevention

6. **Exercise Modifications**
   - Alternative exercises for different equipment
   - Easier/harder variations

7. **Weekly Summary**
   - Total workout days
   - Total exercises
   - Estimated weekly training time

### 4. Track Plan Performance

**Endpoint:** `POST /api/v1/plans/{planID}/performance`

**Request Body:**
```json
{
  "completion_rate": 0.85,
  "average_rpe": 7.5,
  "workouts_completed": 3,
  "workouts_planned": 4,
  "adherence_notes": "Missed one workout due to work"
}
```

### 5. Get Plan Effectiveness Score

**Endpoint:** `GET /api/v1/plans/{planID}/effectiveness`

**Response:**
```json
{
  "plan_id": 1,
  "effectiveness_score": 85.5
}
```

### 6. Mark Plan for Regeneration

**Endpoint:** `POST /api/v1/plans/{planID}/regenerate`

**Request Body:**
```json
{
  "reason": "User plateau - need more variety"
}
```

### 7. Get Plan History

**Endpoint:** `GET /api/v1/plans/history/{userID}?limit=10`

**Response:** Array of historical plans with analysis of plan evolution.

### 8. Get Adaptation History

**Endpoint:** `GET /api/v1/plans/adaptations/{userID}`

**Response:** List of all plan adaptations made for the user.

## How the System Works

### Plan Generation Process

1. **User Input Validation**
   - Validates user preferences and requirements
   - Ensures all required fields are present

2. **Template Selection**
   - Loads fitness data (exercises, templates, goals, levels)
   - Finds templates matching user's fitness level and goals
   - Selects optimal template based on weekly frequency

3. **Exercise Selection**
   - Filters exercises by available equipment
   - Ensures exercises match fitness level
   - Substitutes exercises if needed

4. **Progressive Overload Application**
   - Adjusts sets/reps based on fitness level
   - Applies goal-specific rep ranges
   - Sets appropriate rest periods

5. **Muscle Group Balance**
   - Ensures balanced muscle development
   - Optimizes for recovery patterns

6. **Plan Storage**
   - Saves plan with complete metadata
   - Marks as active plan for user

### Adaptive Features

The system automatically adapts plans based on performance:

- **Low Completion Rate (<60%)**: Reduces volume and intensity
- **High RPE + Low Completion (>8.5 RPE, <80%)**: Adds rest days, reduces intensity
- **High Completion + Low RPE (>90%, <6.0 RPE)**: Increases volume and intensity

## Example Usage Flow

```bash
# 1. Generate a new plan
curl -X POST http://localhost:8080/api/v1/plans/generate \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "metadata": {
      "user_goals": ["muscle_gain"],
      "available_equipment": ["bodyweight", "dumbbell"],
      "fitness_level": "intermediate",
      "weekly_frequency": 4,
      "time_per_workout": 60
    }
  }'

# Response: { "plan_id": 1, ... }

# 2. Download the PDF
curl -X GET http://localhost:8080/api/v1/plans/1/download \
  -o my_workout_plan.pdf

# 3. Track performance after a week
curl -X POST http://localhost:8080/api/v1/plans/1/performance \
  -H "Content-Type: application/json" \
  -d '{
    "completion_rate": 0.75,
    "average_rpe": 7.0,
    "workouts_completed": 3,
    "workouts_planned": 4
  }'

# 4. Get effectiveness score
curl -X GET http://localhost:8080/api/v1/plans/1/effectiveness
```

## PowerShell Testing Scripts

### Generate and Download Plan
```powershell
# Generate plan
$body = @{
    user_id = 123
    metadata = @{
        user_goals = @("muscle_gain")
        available_equipment = @("bodyweight", "dumbbell")
        fitness_level = "intermediate"
        weekly_frequency = 4
        time_per_workout = 60
    }
} | ConvertTo-Json -Depth 10

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/plans/generate" `
    -Method POST `
    -Body $body `
    -ContentType "application/json"

$planId = $response.plan_id
Write-Host "Plan ID: $planId"

# Download PDF
Invoke-WebRequest -Uri "http://localhost:8080/api/v1/plans/$planId/download" `
    -OutFile "workout_plan_$planId.pdf"

Write-Host "PDF downloaded: workout_plan_$planId.pdf"
```

## Data Flow Diagram

```
User Preferences
     ↓
[Validate Input]
     ↓
[Load Fitness Data]
     ↓
[Select Template] → Based on Level, Goals, Frequency
     ↓
[Filter Exercises] → Based on Equipment, Level
     ↓
[Apply Progressive Overload] → Adjust Sets/Reps/Rest
     ↓
[Balance Muscle Groups]
     ↓
[Save Plan with Metadata]
     ↓
[Return Generated Plan]
     ↓
User can → [View Plan] → [Download PDF] → [Track Performance]
```

## Important Notes

1. **One Active Plan Per User**: Users can only have one active plan at a time. To generate a new plan, the old one must be deactivated.

2. **Plan Metadata**: All plan details are stored in the metadata JSON field, including the complete workout structure.

3. **PDF Generation**: PDFs are generated on-demand and include comprehensive workout details, guidelines, and tips.

4. **Adaptive Algorithm**: The system monitors performance and can automatically suggest plan adaptations.

5. **Equipment Flexibility**: The system will substitute exercises if required equipment isn't available.

6. **Level Appropriate**: Exercises are filtered to match or be below the user's fitness level for safety.

## Error Handling

Common error responses:

- `400 Bad Request`: Invalid input data
- `404 Not Found`: Plan not found
- `409 Conflict`: User already has an active plan
- `500 Internal Server Error`: Server-side error

All errors return JSON:
```json
{
  "error": "Error description"
}
```
