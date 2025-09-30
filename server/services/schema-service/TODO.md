# FitUp Workout Plan Generation Service TODO

## 📋 Overview
This document outlines the implementation tasks for FitUp's intelligent workout plan generation service. FitUp is a smart fitness app that creates **personalized, dynamic workout plans** that adapt to user progress, goals, and equipment. This service is the core of FitUp's "personal trainer in your pocket" functionality.

## 🎯 FitUp Core Features
- **Auto-generated workout plans** tailored to fitness level and goals
- **Dynamic weekly updates** that adapt to progress
- **Equipment-aware planning** (bodyweight, dumbbells, full gym, or mix)
- **Progress tracking** with performance history
- **Smart balancing** (missed sessions, exceeded targets)
- **Goal-based programming** (muscle building, weight loss, strength, general fitness)

## ✅ Completed
- [x] Repository layer implementations (all CRUD operations)
- [x] SQL schema alignment fixes
- [x] Type definitions and interfaces
- [x] Database connection patterns
- [x] Basic project structure

## 🚧 Implementation Tasks

### 1. Core Fitness Intelligence Services
**Priority: CRITICAL** | **Estimated Time: 4-5 days**

#### 1.1 Workout Plan Generator Service (`internal/service/plan_generator_service.go`)
- [ ] **Dynamic Plan Generation**
  ```go
  - GenerateWeeklyPlan(ctx, userID, weekStart) (*WeeklyPlanResponse, error)
  - RegenerateFromProgress(ctx, userID, progressData) (*WeeklyPlanResponse, error)
  - AdaptToMissedSessions(ctx, userID, missedSessions) (*WeeklyPlanResponse, error)
  - BalanceNextWeek(ctx, userID, currentWeekData) (*WeeklyPlanResponse, error)
  ```

- [ ] **Goal-Based Programming**
  ```go
  - CreateMuscleGainPlan(ctx, userProfile) (*WorkoutPlan, error)
  - CreateWeightLossPlan(ctx, userProfile) (*WorkoutPlan, error)
  - CreateStrengthPlan(ctx, userProfile) (*WorkoutPlan, error)
  - CreateGeneralFitnessPlan(ctx, userProfile) (*WorkoutPlan, error)
  ```

- [ ] **Equipment-Aware Planning**
  ```go
  - GenerateBodyweightPlan(ctx, userProfile) (*WorkoutPlan, error)
  - GenerateDumbbellPlan(ctx, userProfile, weights) (*WorkoutPlan, error)
  - GenerateGymPlan(ctx, userProfile, equipment) (*WorkoutPlan, error)
  - AdaptPlanToEquipment(ctx, existingPlan, availableEquipment) (*WorkoutPlan, error)
  ```

#### 1.2 Progress Analysis Service (`internal/service/progress_analysis_service.go`)
- [ ] **Performance Tracking**
  ```go
  - AnalyzeWeeklyProgress(ctx, userID, weekData) (*ProgressAnalysis, error)
  - DetectProgressStalls(ctx, userID, exerciseID) (*StallDetection, error)
  - CalculateIntensityAdjustments(ctx, userID, exercisePerformance) (*IntensityAdjustment, error)
  - EvaluateGoalProgress(ctx, userID, goalType) (*GoalProgressReport, error)
  ```

- [ ] **Adaptive Intelligence**
  ```go
  - SuggestProgressions(ctx, userID, exerciseID) ([]ExerciseProgression, error)
  - RecommendDeloadWeek(ctx, userID, fatigueLevel) (*DeloadRecommendation, error)
  - AdjustVolumeBasedOnRecovery(ctx, userID, recoveryMetrics) (*VolumeAdjustment, error)
  ```

#### 1.3 Exercise Selection Service (`internal/service/exercise_selection_service.go`)
- [ ] **Smart Exercise Selection**
  ```go
  - SelectExercisesForGoal(ctx, goal, equipment, level) ([]Exercise, error)
  - RecommendVariations(ctx, exerciseID, reason) ([]ExerciseVariation, error)
  - BalanceMuscleGroups(ctx, currentPlan) (*BalancedPlan, error)
  - SubstituteExercise(ctx, exerciseID, reason, constraints) (*ExerciseSubstitution, error)
  ```

- [ ] **Progressive Overload Management**
  ```go
  - CalculateProgression(ctx, userID, exerciseID, currentWeek) (*ProgressionScheme, error)
  - DetermineSetRepScheme(ctx, goal, exerciseType, userLevel) (*SetRepScheme, error)
  - AdjustIntensity(ctx, userID, exerciseID, performance) (*IntensityAdjustment, error)
  ```

#### 1.4 User Fitness Profile Service (`internal/service/fitness_profile_service.go`)
- [ ] **Fitness Assessment**
  ```go
  - CreateFitnessProfile(ctx, userID, assessmentData) (*FitnessProfile, error)
  - UpdateFitnessLevel(ctx, userID, progressData) (*FitnessProfile, error)
  - EstimateOneRepMax(ctx, userID, exerciseID, performanceData) (*OneRepMaxEstimate, error)
  - AssessMovementPatterns(ctx, userID, movementData) (*MovementAssessment, error)
  ```

- [ ] **Goal Management**
  ```go
  - SetFitnessGoals(ctx, userID, goals) error
  - TrackGoalProgress(ctx, userID) (*GoalProgress, error)
  - UpdateGoalTimeline(ctx, userID, goalID, progress) error
  - SuggestRealisticGoals(ctx, userID, currentMetrics) ([]GoalSuggestion, error)
  ```

### 2. Enhanced Service Layer
**Priority: HIGH** | **Estimated Time: 3-4 days**

#### 2.1 User Service (`internal/service/user_service.go`)
- [ ] **Enhanced User Management**
  ```go
  - CreateUserWithFitnessProfile(ctx, req) (*UserResponse, error)
  - UpdateUserFitnessGoals(ctx, userID, goals) error
  - GetUserDashboard(ctx, userID) (*DashboardResponse, error)
  - UpdateUserEquipment(ctx, userID, equipment) error
  ```

#### 2.2 Workout Session Service (`internal/service/workout_session_service.go`)
- [ ] **Session Management**
  ```go
  - StartWorkoutSession(ctx, userID, workoutID) (*SessionResponse, error)
  - LogExercisePerformance(ctx, sessionID, exerciseID, performance) error
  - CompleteWorkoutSession(ctx, sessionID, sessionData) (*SessionSummary, error)
  - SkipWorkout(ctx, userID, workoutID, reason) error
  ```

#### 2.3 Recovery & Adaptation Service (`internal/service/recovery_service.go`)
- [ ] **Recovery Tracking**
  ```go
  - LogRecoveryMetrics(ctx, userID, metrics) error
  - AssessRecoveryStatus(ctx, userID) (*RecoveryStatus, error)
  - RecommendRestDay(ctx, userID, fatigueLevel) (*RestDayRecommendation, error)
  - AdjustNextWorkout(ctx, userID, recoveryData) (*WorkoutAdjustment, error)
  ```
### 3. API Handler Layer (FitUp-Focused)
**Priority: HIGH** | **Estimated Time: 2-3 days**

#### 3.1 Workout Plan Handler (`internal/handlers/workout_plan_handler.go`)
- [ ] **Plan Generation Endpoints**
  ```
  POST   /api/v1/plans/generate              - Generate new weekly plan
  POST   /api/v1/plans/regenerate            - Regenerate plan based on progress
  GET    /api/v1/plans/current/{userID}      - Get current week's plan
  POST   /api/v1/plans/adapt                 - Adapt plan for missed sessions
  PUT    /api/v1/plans/{planID}/feedback     - Provide feedback for plan adaptation
  ```

#### 3.2 Fitness Dashboard Handler (`internal/handlers/dashboard_handler.go`)
- [ ] **Dashboard Endpoints**
  ```
  GET    /api/v1/dashboard/{userID}          - Get complete fitness dashboard
  GET    /api/v1/dashboard/{userID}/stats    - Get fitness statistics
  GET    /api/v1/dashboard/{userID}/progress - Get progress overview
  GET    /api/v1/dashboard/{userID}/goals    - Get goal progress
  ```

#### 3.3 Workout Session Handler (`internal/handlers/session_handler.go`)
- [ ] **Session Management Endpoints**
  ```
  POST   /api/v1/sessions/start              - Start workout session
  PUT    /api/v1/sessions/{sessionID}/log    - Log exercise performance
  POST   /api/v1/sessions/{sessionID}/complete - Complete session
  POST   /api/v1/sessions/{sessionID}/skip   - Skip workout with reason
  GET    /api/v1/sessions/{sessionID}/summary - Get session summary
  ```

#### 3.4 Fitness Profile Handler (`internal/handlers/fitness_profile_handler.go`)
- [ ] **Profile Management Endpoints**
  ```
  POST   /api/v1/profile/{userID}/assessment - Initial fitness assessment
  PUT    /api/v1/profile/{userID}/goals      - Update fitness goals
  PUT    /api/v1/profile/{userID}/equipment  - Update available equipment
  GET    /api/v1/profile/{userID}/level      - Get current fitness level
  POST   /api/v1/profile/{userID}/1rm        - Estimate one-rep max
  ```

#### 3.5 Exercise Library Handler (`internal/handlers/exercise_library_handler.go`)
- [ ] **Enhanced Exercise Endpoints**
  ```
  GET    /api/v1/exercises/recommendations/{userID} - Get personalized exercise recommendations
  GET    /api/v1/exercises/substitutions/{exerciseID} - Get exercise substitutions
  GET    /api/v1/exercises/progressions/{exerciseID} - Get exercise progressions
  GET    /api/v1/exercises/by-equipment     - Filter by available equipment
  GET    /api/v1/exercises/by-goal         - Filter by fitness goal
  ```

### 4. Smart Logic & Analytics Components
**Priority: MEDIUM** | **Estimated Time: 3-4 days**

#### 4.1 Plan Generation Algorithm (`internal/algorithm/plan_generator.go`)
- [ ] **Core Rule-Based Algorithm Implementation**
  ```go
  - GeneratePlanBasedOnGoal(goal, profile, equipment) (*WorkoutPlan, error)
  - BalanceWeeklyVolume(exercises, targetVolume) (*BalancedPlan, error)
  - ApplyProgressiveOverload(currentPlan, progressData) (*UpdatedPlan, error)
  - OptimizeRecoveryTime(exercises, userRecovery) (*OptimizedSchedule, error)
  ```

#### 4.2 Adaptation Engine (`internal/algorithm/adaptation_engine.go`)
- [ ] **Dynamic Adaptation Logic**
  ```go
  - AdaptForMissedWorkouts(plan, missedSessions) (*AdaptedPlan, error)
  - AdjustForPerformance(plan, performanceData) (*AdjustedPlan, error)
  - RebalanceForWeakPoints(plan, assessmentData) (*RebalancedPlan, error)
  ```

#### 4.3 Progress Calculator (`internal/algorithm/progress_calculator.go`)
- [ ] **Rule-Based Progress Calculations**
  ```go
  - CalculateGoalAchievement(userID, goalID, currentProgress) (*Calculation, error)
  - EstimateTimeToGoal(userID, goalType, targetMetrics) (*TimeEstimate, error)
  - SuggestGoalAdjustments(userID, currentTrajectory) ([]GoalAdjustment, error)
  ```

### 5. FitUp-Specific Middleware & Utilities
**Priority: MEDIUM** | **Estimated Time: 1-2 days**

#### 5.1 Fitness Context Middleware (`internal/middleware/fitness_context.go`)
- [ ] **User fitness context injection**
- [ ] **Equipment availability checking**
- [ ] **Goal-based access control**

#### 5.2 Progress Tracking Middleware (`internal/middleware/progress_tracking.go`)
- [ ] **Automatic progress logging**
- [ ] **Session state management**
- [ ] **Performance metrics collection**

#### 5.3 Plan Validation Service (`internal/service/plan_validation_service.go`)
- [ ] **Plan safety validation**
- [ ] **Volume appropriateness checking**
- [ ] **Recovery time validation**
- [ ] **Equipment requirement validation**
### 6. Testing Strategy (FitUp-Focused)
**Priority: MEDIUM** | **Estimated Time: 2-3 days**

#### 6.1 Algorithm Testing
- [ ] **Plan Generation Tests** (`tests/algorithm/plan_generator_test.go`)
  - Test plan generation for different goals
  - Test equipment-based plan adaptation
  - Test progressive overload calculations
  - Test volume balancing

- [ ] **Adaptation Logic Tests** (`tests/algorithm/adaptation_test.go`)
  - Test missed workout handling
  - Test performance-based adaptations
  - Test recovery-based adjustments

#### 6.2 Service Integration Tests
- [ ] **End-to-End Plan Generation** (`tests/integration/plan_e2e_test.go`)
  - Complete user journey from profile creation to plan generation
  - Test plan regeneration based on progress
  - Test multi-week plan evolution

#### 6.3 Performance Tests
- [ ] **Plan Generation Performance** (`tests/performance/`)
  - Test plan generation speed under load
  - Test database query optimization
  - Test concurrent user plan generation

### 7. FitUp Configuration & Deployment
**Priority: MEDIUM** | **Estimated Time: 1-2 days**

#### 7.1 Fitness-Specific Configuration
- [ ] **Algorithm Parameters** (`internal/config/fitness_config.go`)
  ```go
  type FitnessConfig struct {
      PlanGeneration PlanGenerationConfig
      ProgressionRules ProgressionConfig
      RecoverySettings RecoveryConfig
      GoalParameters GoalConfig
  }
  ```

#### 7.2 Exercise Database Seeding
- [ ] **Exercise Library Setup** (`scripts/seed_exercises.sql`)
  - Comprehensive exercise database
  - Equipment mappings
  - Muscle group categorizations
  - Difficulty progressions

#### 7.3 Plan Templates
- [ ] **Default Plan Templates** (`scripts/seed_templates.sql`)
  - Beginner templates for each goal
  - Intermediate and advanced progressions
  - Equipment-specific templates

## 🎯 FitUp Implementation Roadmap

### Phase 1: Core Logic Engine (Week 1-2)
**Focus: The "Personal Trainer" Brain**
1. **Plan Generator Service** - Auto-generate personalized plans
2. **Progress Analysis Service** - Track and analyze user performance
3. **Exercise Selection Service** - Smart exercise selection and progression
4. **Fitness Profile Service** - User assessment and goal management

### Phase 2: Dynamic Adaptation (Week 2-3)
**Focus: The "Evolves With You" Feature**
1. **Adaptation Engine** - Handle missed workouts and performance changes
2. **Recovery Service** - Monitor and adjust for recovery
3. **Workout Session Service** - Real-time session management
4. **Progress Predictor** - Predictive analytics for goal achievement

### Phase 3: User Experience (Week 3-4)
**Focus: "Simple, Clean Interface"**
1. **Dashboard Handler** - Comprehensive fitness dashboard
2. **Session Handler** - Seamless workout tracking
3. **Plan Handler** - Easy plan management
4. **Profile Handler** - Intuitive fitness profile management

### Phase 4: Production Ready (Week 4-5)
**Focus: "Available 24/7"**
1. **Comprehensive Testing** - Ensure reliability
2. **Performance Optimization** - Fast plan generation
3. **Documentation** - Clear API documentation
4. **Monitoring & Analytics** - Production observability

## 🎪 FitUp Unique Features Implementation

### 1. Dynamic Weekly Re-generation
- [ ] **Smart Plan Evolution**: Plans that adapt based on your performance data
- [ ] **Missed Session Balancing**: Automatically adjust next week for missed workouts
- [ ] **Performance-Based Progression**: Increase intensity when you excel

### 2. Equipment-Aware Logic
- [ ] **Adaptive Equipment Planning**: Switch between home and gym seamlessly
- [ ] **Substitution Engine**: Automatic exercise substitutions based on available equipment
- [ ] **Progressive Equipment Introduction**: Gradually introduce new equipment as user advances

### 3. Goal-Driven Programming
- [ ] **Muscle Building Focus**: Hypertrophy-optimized plans with proper volume distribution
- [ ] **Fat Loss Optimization**: High-intensity circuits combined with strength training
- [ ] **Strength Development**: Progressive overload with optimal rest periods
- [ ] **General Fitness**: Balanced approach for overall health and wellness

## 📋 Implementation Checklist

### Daily Priorities
- [ ] **Start with Plan Generator**: This is FitUp's core differentiator
- [ ] **Focus on User Journey**: Profile → Assessment → Plan → Adaptation
- [ ] **Test with Real Scenarios**: Different goals, equipment, experience levels
- [ ] **Validate Algorithm Logic**: Ensure plans make fitness sense

### Key Success Metrics
- [ ] **Plan Personalization**: No two users should get identical plans
- [ ] **Adaptation Responsiveness**: Plans should evolve weekly based on performance
- [ ] **Goal Alignment**: Generated plans should clearly support user goals
- [ ] **Equipment Optimization**: Plans should maximize available equipment usage

## 🚀 Getting Started with FitUp Intelligence

1. **Start with Fitness Profile Service**: Build the foundation of user understanding
2. **Implement Plan Generator**: The core "personal trainer" algorithm
3. **Add Progress Analysis**: Enable the system to learn from user performance
4. **Build Adaptation Engine**: Make plans truly dynamic and responsive
5. **Create Seamless APIs**: Ensure mobile app integration is smooth

## 📞 FitUp-Specific Considerations

- **Fitness Safety**: Always validate plans for safety and appropriateness
- **Progressive Overload**: Implement scientifically-backed progression principles
- **Recovery Management**: Balance training stress with adequate recovery
- **User Experience**: Keep the complexity hidden behind a simple interface
- **Scalability**: Ensure the system can handle thousands of users generating plans simultaneously

---

**Total Estimated Time: 15-20 days for complete FitUp intelligence implementation**

**Current Status: Repository layer complete ✅ | FitUp Smart Logic pending 🚧**

**Next Priority: Plan Generator Service - The Heart of FitUp's "Personal Trainer" 🎯**
  - Mock service tests
  - Request/response validation

#### 5.2 Integration Tests
- [ ] **End-to-end API tests** (`tests/integration/`)
- [ ] **Database migration tests**
- [ ] **Performance tests**

#### 5.3 Test Utilities
- [ ] **Test fixtures** (`tests/fixtures/`)
- [ ] **Mock generators** (`tests/mocks/`)
- [ ] **Test database setup** (`tests/setup/`)

### 6. Documentation
**Priority: LOW** | **Estimated Time: 1-2 days**

#### 6.1 API Documentation
- [ ] **OpenAPI/Swagger specification**
- [ ] **Endpoint documentation**
- [ ] **Request/response examples**

#### 6.2 Developer Documentation
- [ ] **Architecture overview**
- [ ] **Database schema documentation**
- [ ] **Deployment guide**
- [ ] **Contributing guidelines**

### 7. DevOps and Deployment
**Priority: LOW** | **Estimated Time: 1-2 days**

#### 7.1 Containerization
- [ ] **Dockerfile optimization**
- [ ] **Docker Compose for development**
- [ ] **Multi-stage builds**

#### 7.2 CI/CD Pipeline
- [ ] **GitHub Actions workflow**
- [ ] **Automated testing**
- [ ] **Security scanning**
- [ ] **Deployment automation**

#### 7.3 Monitoring and Observability
- [ ] **Health check endpoints**
- [ ] **Metrics collection**
- [ ] **Logging configuration**
- [ ] **Error tracking**

## 🎯 Implementation Priority Order

### Phase 1: Core Functionality (Week 1)
1. Service layer implementation (Users, Exercises, Workouts)
2. Basic handler implementation
3. Configuration setup
4. Main application bootstrap

### Phase 2: Advanced Features (Week 2)
1. Schema and Template services
2. Progress tracking service
3. Recommendation service
4. Complete handler implementation

### Phase 3: Quality and Production (Week 3)
1. Comprehensive testing
2. Documentation
3. DevOps setup
4. Performance optimization

## 📋 Implementation Checklist

### Daily Tasks
- [ ] Review and implement one service at a time
- [ ] Write tests alongside implementation
- [ ] Update documentation as you go
- [ ] Test endpoints manually with tools like Postman

### Weekly Reviews
- [ ] Code review and refactoring
- [ ] Performance benchmarking
- [ ] Security review
- [ ] Documentation updates

## 🚀 Getting Started

1. **Start with User Service**: It's the foundation for other services
2. **Implement one complete vertical slice**: User creation flow from handler → service → repository
3. **Add tests immediately**: Don't accumulate technical debt
4. **Use TDD approach**: Write tests first, then implement
5. **Regular commits**: Small, focused commits with clear messages

## 📞 Need Help?

- Check existing implementation patterns in the repository layer
- Refer to the interfaces for method signatures
- Use the SQL schema as the source of truth for data structures
- Follow Go best practices for error handling and naming conventions

---

**Total Estimated Time: 10-15 days for complete implementation**

**Current Status: Repository layer complete ✅ | Service layer pending 🚧**