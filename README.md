# FitUp - Intelligent Fitness Training Platform

**FitUp** is a comprehensive fitness training platform that combines a React Native mobile application with an intelligent backend system powered by advanced adaptive algorithms. The platform provides personalized workout plans, real-time performance tracking, and smart adaptations based on user progress.

---

## 🧠 FitUp Smart Logic Systemm

FitUp's core innovation lies in its **Smart Logic System** - a collection of intelligent algorithms that continuously analyze user performance, detect patterns, and automatically adapt training programs to optimize results while preventing overtraining and injuries.

### 🎯 Core Components of FitUp Smart Logic

The Smart Logic system is built on four main pillars:

#### 1. **Adaptive Plan Generation** (`fitup_adaptive_v1`)

The plan generation algorithm creates personalized workout programs by:

- **Template Selection**: Automatically selects optimal workout templates based on:
  - User's fitness level (beginner, intermediate, advanced)
  - Primary training goals (muscle gain, strength, fat loss, endurance, general fitness)
  - Weekly training frequency (1-7 days per week)
  - Available equipment (bodyweight, dumbbells, barbells, cables, full gym)

- **Exercise Selection & Substitution**:
  - Filters exercises based on available equipment
  - Ensures exercises match or are below user's fitness level for safety
  - Intelligently substitutes unavailable exercises with equivalents based on:
    - Movement pattern matching (push/pull/squat/hinge)
    - Muscle group overlap (minimum 50% overlap)
    - Primary muscle targeting

- **Progressive Overload Application**:
  - Adjusts training volume based on fitness level:
    - Beginner: 0.8x multiplier (reduced volume)
    - Intermediate: 1.0x multiplier (standard volume)
    - Advanced: 1.2x multiplier (increased volume)
  - Customizes rep ranges for specific goals:
    - Strength: 1-5 reps, 180-300s rest
    - Muscle Gain: 6-12 reps, 60-120s rest
    - Endurance: 12-20+ reps, 30-60s rest
  - Sets appropriate rest periods based on exercise type

- **Muscle Group Balance Optimization**:
  - Tracks muscle groups across the week
  - Ensures push/pull balance
  - Optimizes for recovery patterns
  - Prevents muscle group overtraining

**Example Output**:
```json
{
  "algorithm": "fitup_adaptive_v1",
  "parameters": {
    "template_used": "intermediate_muscle_gain_4day",
    "total_exercises": 24,
    "muscle_groups_targeted": ["chest", "back", "shoulders", "legs", "arms"],
    "equipment_utilized": ["bodyweight", "dumbbell", "barbell"],
    "estimated_volume": 96,
    "progression_method": "progressive_overload"
  }
}
```

---

#### 2. **Real-Time Performance Analysis**

During workout sessions, the Smart Logic system continuously monitors and analyzes:

- **Exercise Performance Tracking**:
  - Logs sets, reps, weight, and RPE (Rate of Perceived Exertion) in real-time
  - Detects form warnings based on RPE inconsistency
  - Identifies plateau indicators (no improvement over multiple sessions)
  - Flags progression readiness (consistent performance above target)

- **1RM (One-Rep Max) Estimation**:
  - Uses multiple validated formulas for accuracy:
    - **Epley Formula**: Best for 1-5 reps (95% confidence)
    - **Brzycki Formula**: Best for 6-10 reps (85% confidence)
    - **McGlothin & Lombardi**: Additional validation
  - Adjusts confidence based on rep range and training history
  - Validates estimates against historical data
  - Detects impossible jumps (>10% increase in 7 days)

- **Session Quality Metrics**:
  - **Completion Rate**: Percentage of planned exercises completed
  - **Quality Score**: Composite metric based on:
    - Form consistency (RPE variance)
    - Set completion (planned vs actual)
    - Time efficiency (duration vs target)
  - **Performance Level Classification**:
    - Excellent: ≥90% completion, RPE 7-8
    - Good: 70-89% completion, RPE 6-8
    - Poor: <70% completion, RPE >8.5 or <5

**Real-Time Feedback Example**:
```
✓ Form consistent (RPE: 7.5 avg, variance: 0.3)
⚠ Plateau detected - no weight increase in 3 weeks
✓ Ready for progression - 3 consecutive sessions above target
```

---

#### 3. **Intelligent Plan Adaptation**

The system automatically adapts workout plans based on performance patterns:

- **Low Completion Detection** (<60% completion rate):
  - **Action**: Reduce workout intensity and volume
  - **Changes**: 
    - Fewer sets per exercise (-20%)
    - Lower intensity zones
    - Additional rest between sets
  - **Trigger**: "volume_reduction"

- **Overtraining Prevention** (High RPE + Low Completion):
  - **Condition**: Average RPE >8.5 AND completion <80%
  - **Action**: Add recovery focus
  - **Changes**:
    - Insert additional rest days
    - Reduce training intensity by 15-20%
    - Increase rest periods
    - Add recovery modalities (stretching, mobility)
  - **Trigger**: "potential_overtraining"

- **Progression Trigger** (High Performance):
  - **Condition**: Completion >90% AND RPE <6.0
  - **Action**: Increase training stimulus
  - **Changes**:
    - Increase sets/reps by 10-15%
    - Add weight to exercises
    - Reduce rest periods for density
  - **Trigger**: "ready_for_progression"

- **Skip Pattern Analysis**:
  - Tracks workout skip frequency and reasons
  - Identifies common obstacles (time, fatigue, soreness)
  - Suggests schedule adjustments if skip rate >20%
  - Recommends alternative workout durations

**Adaptation History Tracking**:
```json
{
  "total_adaptations": 12,
  "recent_adaptations": 3,
  "most_common_reason": "ready_for_progression",
  "adaptation_rate_30_days": 0.1
}
```

---


### 🔄 How the Smart Logic Works Together

**Workflow Example**:

1. **User Onboarding**:
   - Fitness level assessment (beginner/intermediate/advanced)
   - Goal selection (muscle gain, strength, etc.)
   - Equipment availability input
   - Schedule preferences (4 days/week, 60 min/session)

2. **Initial Plan Generation**:
   - Smart Logic analyzes inputs
   - Selects "Intermediate Muscle Gain 4-Day Split" template
   - Filters 24 appropriate exercises
   - Generates week 1 with 96 total sets
   - Sets target RPE ranges (7-8 for main lifts)

3. **Week 1-2 Execution**:
   - User completes workouts
   - Real-time performance tracking logs all sets
   - System calculates 1RM estimates
   - Completion rate: 92%, Avg RPE: 7.2
   - **Result**: "Excellent performance, continue as planned"

4. **Week 3-4 Analysis**:
   - Plateau detected on bench press (3 weeks no progress)
   - Skip rate increasing (2 missed workouts in 2 weeks)
   - Recovery score declining (avg 68)
   - **Smart Adaptation Triggered**:
     - Add bench press variation (incline press)
     - Reduce weekly volume by 10%
     - Add extra rest day mid-week

5. **Week 5-8 Progression**:
   - Plateau broken, bench press increases 5kg
   - Completion rate improves to 95%
   - Recovery score stabilizes at 75
   - **Progression Triggered**:
     - Increase volume back to baseline + 5%
     - Add progressive overload (2.5% weight increase)
     - Introduce advanced technique (drop sets on final set)

6. **Week 12 Goal Achievement**:
   - Bench press 1RM increased from 80kg to 95kg
   - Goal prediction: 100kg achievable in 3 more weeks
   - New plan generated with advanced template
   - Success probability: 88%

---

### 📊 Smart Logic Benefits

✅ **Personalized**: Every aspect tailored to individual user  
✅ **Adaptive**: Continuously adjusts based on real performance  
✅ **Predictive**: Forecasts progress and prevents plateaus  
✅ **Safe**: Detects overtraining and enforces recovery  
✅ **Scientific**: Based on validated formulas and research  
✅ **Automated**: No manual plan adjustments needed  
✅ **Comprehensive**: Covers all aspects of training  

---

### 🎓 Scientific Foundation

The FitUp Smart Logic is built on proven principles from exercise science:

- **Progressive Overload**: Gradual increase in training stimulus
- **Specificity**: Training specific to goals
- **Periodization**: Planned variation in training
- **Recovery**: Adequate rest for adaptation
- **Individual Differences**: Personalized programming

**Key Algorithms**:
- 1RM estimation (Epley, Brzycki formulas)
- Linear regression for progression analysis
- Statistical plateau detection
- Weighted scoring for recovery
- Predictive modeling for goal achievement

---


## 🚀 Getting Started

See individual README files for setup:
- [Server Setup](server/README.md)
- [App Setup](app/README.md)

---

## 📄 License

MIT License - see LICENSE file for details.

---

**Built with intelligence. Trained with purpose. FitUp - Your smart training companion.** 🏋️‍♂️💪