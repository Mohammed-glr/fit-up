# Plan Generation Algorithm

## Overview
The **FitUp Adaptive Algorithm** (`fitup_adaptive_v1`) automatically creates personalized workout plans based on user goals, fitness level, and available equipment. It follows a data-driven approach using pre-defined templates and exercise databases.

---

## How It Works

### 1. **Input Collection**
The algorithm receives user preferences:
- **Fitness Goals** (muscle gain, strength, fat loss, endurance, general fitness)
- **Fitness Level** (beginner, intermediate, advanced)
- **Weekly Frequency** (number of workout days per week)
- **Time Per Workout** (minutes available per session)
- **Available Equipment** (barbell, dumbbell, bodyweight, machines, etc.)

### 2. **Template Selection**
The system searches the workout template database:
- Filters templates matching the user's **goal** and **level**
- Finds the best template matching the **weekly frequency**
- Falls back to the closest frequency match if exact match isn't found

**Example:** For an intermediate user wanting muscle gain with 4 workouts/week → selects `intermediate_muscle_gain_4day` template

### 3. **Exercise Selection**
For each workout day in the template:
- Loads exercises from the template structure
- **Filters by equipment availability** (ensures user has the required equipment)
- **Filters by fitness level** (prevents advanced exercises for beginners)
- **Adds bodyweight exercises** (always available as fallback options)
- **Finds substitutes** when original exercise isn't available

### 4. **Progressive Overload Application**
The algorithm adjusts workout parameters based on user level:
- Sets appropriate **sets and reps** based on the goal (e.g., 8-12 reps for muscle gain)
- Defines **rest periods** (60-90s for hypertrophy, 120-180s for strength)
- Calculates **intensity guidelines** (RPE ranges, progression rates)

### 5. **Muscle Group Balance**
Ensures balanced training:
- Tracks which muscle groups are targeted across the week
- Distributes exercises evenly to prevent overtraining
- Maintains proper **push/pull/legs** split ratios

### 6. **Plan Structure Generation**
Creates the final workout structure:
- Organizes workouts by day (Day 1, Day 2, etc.)
- Assigns each day a **focus** (e.g., "Upper Push", "Lower Body")
- Lists exercises with specific parameters:
  - Exercise name and ID
  - Sets and reps
  - Rest time
  - Form cues/notes

### 7. **Metadata Enhancement**
The plan is enriched with analytics:
- Total exercises count
- Muscle groups targeted
- Equipment utilized
- Estimated weekly volume (total sets)
- Progression method (e.g., progressive overload)
- Week start date

### 8. **Database Storage**
The generated plan is saved to the database:
- Creates a `GeneratedPlan` record
- Stores all workout days and exercises
- Marks the plan as **active**
- Links it to the user's profile

---

## Key Features

### **Adaptive Intelligence**
- Automatically substitutes exercises when equipment is unavailable
- Ensures safety by filtering exercises above user's skill level
- Balances muscle groups to prevent imbalances

### **Progressive Structure**
- Builds from the user's current fitness level
- Provides clear progression paths for each exercise
- Includes regression options for difficult movements

### **Time-Aware**
- Respects user's time constraints per workout
- Adjusts volume based on available session duration
- Calculates total weekly training time

### **Equipment Flexibility**
- Works with any equipment combination
- Includes bodyweight alternatives
- Optimizes based on what's available

---

## Algorithm Flow

```
User Input
    ↓
Load Fitness Data (exercises, templates, goals, levels)
    ↓
Select Optimal Template (match goal + level + frequency)
    ↓
Generate Exercise Selection (filter by equipment & level)
    ↓
Apply Progressive Overload (adjust sets/reps/rest)
    ↓
Optimize Muscle Balance (distribute work evenly)
    ↓
Build Plan Structure (organize days & exercises)
    ↓
Save to Database
    ↓
Return Generated Plan
```

---

## Example Output

For an **intermediate** user wanting **muscle gain** with **4 days/week**, **dumbbell** equipment:

**Generated Plan:**
- **Day 1:** Upper Push (Chest, Shoulders, Triceps)
  - Dumbbell Bench Press: 4 sets × 8-12 reps
  - Dumbbell Shoulder Press: 3 sets × 10-12 reps
  - Tricep Extensions: 3 sets × 12-15 reps

- **Day 2:** Lower Body (Quads, Glutes, Hamstrings)
  - Goblet Squats: 4 sets × 10-12 reps
  - Dumbbell Lunges: 3 sets × 12-15 reps
  - Romanian Deadlifts: 3 sets × 8-10 reps

- **Day 3:** Upper Pull (Back, Biceps)
- **Day 4:** Full Body (Mixed muscle groups)

---

## Success Factors

1. **Validation:** Checks user input before generation
2. **Limit Control:** Maximum 3 active plans per user
3. **Error Handling:** Falls back gracefully if template/exercises unavailable
4. **Persistence:** All data stored for tracking and adaptation
5. **PDF Export:** Plans can be downloaded as formatted documents

---

## Future Enhancements

The algorithm monitors plan effectiveness through:
- **Completion rate** (% of workouts finished)
- **Average RPE** (perceived exertion)
- **Progress rate** (strength/volume gains)
- **User satisfaction** ratings

Based on performance data, the system can trigger **plan adaptations** to optimize results over time.
