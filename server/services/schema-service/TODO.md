# Schema Service Implementation TODO

## 📋 Overview
This document outlines the remaining implementation tasks for the fit-up schema service. The repository layer has been implemented and aligned with the SQL schema. The following layers and components need to be completed.

## ✅ Completed
- [x] Repository layer implementations (all CRUD operations)
- [x] SQL schema alignment fixes
- [x] Type definitions and interfaces
- [x] Database connection patterns
- [x] Basic project structure

## 🚧 Implementation Tasks

### 1. Service Layer Implementation
**Priority: HIGH** | **Estimated Time: 3-4 days**

#### 1.1 User Service (`internal/service/user_service.go`)
- [ ] **UserService struct and constructor**
  - Dependency injection for repository
  - Configuration management
  - Logger integration

- [ ] **Core user operations**
  ```go
  - CreateUser(ctx, req) (*UserResponse, error)
  - GetUserByID(ctx, id) (*UserResponse, error)
  - GetUserByEmail(ctx, email) (*UserResponse, error)
  - UpdateUserProfile(ctx, id, req) (*UserResponse, error)
  - DeleteUser(ctx, id) error
  ```

- [ ] **User preferences and settings**
  ```go
  - UpdateUserPreferences(ctx, id, preferences) error
  - GetUserStats(ctx, id) (*UserStatsResponse, error)
  - ValidateUserEquipment(ctx, equipment) error
  ```

#### 1.2 Exercise Service (`internal/service/exercise_service.go`)
- [ ] **Exercise management**
  ```go
  - CreateExercise(ctx, req) (*ExerciseResponse, error)
  - GetExercisesByFilter(ctx, filter) (*PaginatedExerciseResponse, error)
  - SearchExercises(ctx, query, pagination) (*PaginatedExerciseResponse, error)
  - GetRecommendedExercises(ctx, userID, muscleGroup) ([]ExerciseResponse, error)
  ```

- [ ] **Exercise categorization**
  ```go
  - GetExercisesByMuscleGroup(ctx, muscleGroup) ([]ExerciseResponse, error)
  - GetExercisesByEquipment(ctx, equipment) ([]ExerciseResponse, error)
  - GetExercisesByDifficulty(ctx, level) ([]ExerciseResponse, error)
  ```

- [ ] **Bulk operations**
  ```go
  - BulkImportExercises(ctx, exercises) (*BulkImportResult, error)
  - ValidateExerciseData(ctx, exercise) (*ValidationResult, error)
  ```

#### 1.3 Workout Service (`internal/service/workout_service.go`)
- [ ] **Workout management**
  ```go
  - CreateWorkout(ctx, req) (*WorkoutResponse, error)
  - GetWorkoutWithExercises(ctx, workoutID) (*WorkoutWithExercisesResponse, error)
  - UpdateWorkout(ctx, id, req) (*WorkoutResponse, error)
  - DeleteWorkout(ctx, id) error
  ```

- [ ] **Workout-Exercise relationships**
  ```go
  - AddExerciseToWorkout(ctx, workoutID, exerciseReq) error
  - RemoveExerciseFromWorkout(ctx, workoutID, exerciseID) error
  - UpdateWorkoutExercise(ctx, weID, req) error
  - ReorderWorkoutExercises(ctx, workoutID, order) error
  ```

#### 1.4 Schema Service (`internal/service/schema_service.go`)
- [ ] **Weekly schema management**
  ```go
  - CreateWeeklySchema(ctx, req) (*WeeklySchemaResponse, error)
  - GetUserActiveSchema(ctx, userID) (*WeeklySchemaResponse, error)
  - GetSchemaWithWorkouts(ctx, schemaID) (*WeeklySchemaWithWorkoutsResponse, error)
  - CloneSchema(ctx, schemaID, newWeekStart) (*WeeklySchemaResponse, error)
  ```

- [ ] **Schema operations**
  ```go
  - ActivateSchema(ctx, schemaID) error
  - DeactivateSchema(ctx, schemaID) error
  - GenerateSchemaFromTemplate(ctx, userID, templateID, weekStart) (*WeeklySchemaResponse, error)
  ```

#### 1.5 Template Service (`internal/service/template_service.go`)
- [ ] **Template management**
  ```go
  - GetRecommendedTemplates(ctx, userID) ([]WorkoutTemplateResponse, error)
  - GetTemplatesByFilter(ctx, filter) (*PaginatedTemplateResponse, error)
  - CreateCustomTemplate(ctx, userID, req) (*WorkoutTemplateResponse, error)
  ```

#### 1.6 Progress Service (`internal/service/progress_service.go`)
- [ ] **Progress tracking**
  ```go
  - LogProgress(ctx, req) (*ProgressLogResponse, error)
  - GetUserProgressSummary(ctx, userID) (*ProgressSummaryResponse, error)
  - GetProgressTrend(ctx, userID, exerciseID, days) (*ProgressTrendResponse, error)
  - CalculatePersonalBests(ctx, userID) ([]PersonalBestResponse, error)
  ```

#### 1.7 Recommendation Service (`internal/service/recommendation_service.go`)
- [ ] **AI-powered recommendations**
  ```go
  - RecommendWorkoutPlan(ctx, userID) (*WorkoutPlanResponse, error)
  - RecommendExerciseProgression(ctx, userID, exerciseID) (*ProgressionResponse, error)
  - RecommendRestDays(ctx, userID) (*RestDayResponse, error)
  ```

#### 1.8 Validation Service (`internal/service/validation_service.go`)
- [ ] **Input validation**
  ```go
  - ValidateUserRequest(req) error
  - ValidateExerciseRequest(req) error
  - ValidateWorkoutRequest(req) error
  - ValidateProgressLogRequest(req) error
  ```

### 2. Handler Layer Implementation
**Priority: HIGH** | **Estimated Time: 2-3 days**

#### 2.1 Base Handler Structure (`internal/handlers/handlers.go`)
- [ ] **Handler struct and dependencies**
  ```go
  type Handlers struct {
      services *service.Services
      logger   *logger.Logger
      config   *config.Config
  }
  ```

- [ ] **Common middleware**
  - Request ID middleware
  - Logging middleware
  - CORS middleware
  - Rate limiting
  - Authentication middleware

#### 2.2 Individual Handlers

##### User Handler (`internal/handlers/user_handler.go`)
- [ ] **HTTP endpoints**
  ```
  POST   /api/v1/users              - Create user
  GET    /api/v1/users/{id}         - Get user by ID
  PUT    /api/v1/users/{id}         - Update user
  DELETE /api/v1/users/{id}         - Delete user
  GET    /api/v1/users/{id}/stats   - Get user stats
  ```

##### Exercise Handler (`internal/handlers/exercise_handler.go`)
- [ ] **HTTP endpoints**
  ```
  GET    /api/v1/exercises                    - List exercises (with filters)
  POST   /api/v1/exercises                    - Create exercise
  GET    /api/v1/exercises/{id}               - Get exercise by ID
  PUT    /api/v1/exercises/{id}               - Update exercise
  DELETE /api/v1/exercises/{id}               - Delete exercise
  GET    /api/v1/exercises/search             - Search exercises
  GET    /api/v1/exercises/recommended/{userID} - Get recommended exercises
  ```

##### Workout Handler (`internal/handlers/workout_handler.go`)
- [ ] **HTTP endpoints**
  ```
  POST   /api/v1/workouts                     - Create workout
  GET    /api/v1/workouts/{id}                - Get workout
  PUT    /api/v1/workouts/{id}                - Update workout
  DELETE /api/v1/workouts/{id}                - Delete workout
  GET    /api/v1/workouts/{id}/exercises      - Get workout with exercises
  POST   /api/v1/workouts/{id}/exercises      - Add exercise to workout
  DELETE /api/v1/workouts/{id}/exercises/{exerciseId} - Remove exercise
  ```

##### Schema Handler (`internal/handlers/schema_handler.go`)
- [ ] **HTTP endpoints**
  ```
  POST   /api/v1/schemas                      - Create weekly schema
  GET    /api/v1/schemas/{id}                 - Get schema
  PUT    /api/v1/schemas/{id}                 - Update schema
  DELETE /api/v1/schemas/{id}                 - Delete schema
  GET    /api/v1/users/{userID}/active-schema - Get active schema
  POST   /api/v1/schemas/{id}/activate        - Activate schema
  POST   /api/v1/schemas/{id}/clone           - Clone schema
  ```

##### Template Handler (`internal/handlers/template_handler.go`)
- [ ] **HTTP endpoints**
  ```
  GET    /api/v1/templates                    - List templates
  POST   /api/v1/templates                    - Create template
  GET    /api/v1/templates/{id}               - Get template
  PUT    /api/v1/templates/{id}               - Update template
  DELETE /api/v1/templates/{id}               - Delete template
  GET    /api/v1/templates/recommended/{userID} - Get recommended templates
  ```

##### Progress Handler (`internal/handlers/progress_handler.go`)
- [ ] **HTTP endpoints**
  ```
  POST   /api/v1/progress                     - Log progress
  GET    /api/v1/progress/user/{userID}       - Get user progress
  GET    /api/v1/progress/summary/{userID}    - Get progress summary
  GET    /api/v1/progress/trends/{userID}     - Get progress trends
  GET    /api/v1/progress/personal-bests/{userID} - Get personal bests
  ```

##### Health Handler (`internal/handlers/health_handler.go`)
- [ ] **Health check endpoints**
  ```
  GET    /health                              - Basic health check
  GET    /health/ready                        - Readiness probe
  GET    /health/live                         - Liveness probe
  ```

### 3. Configuration and Setup
**Priority: MEDIUM** | **Estimated Time: 1-2 days**

#### 3.1 Configuration Management
- [ ] **Config structure** (`internal/config/config.go`)
  ```go
  type Config struct {
      Server   ServerConfig
      Database DatabaseConfig
      Auth     AuthConfig
      Logging  LoggingConfig
      Redis    RedisConfig
  }
  ```

- [ ] **Environment-based configuration**
  - Development config
  - Testing config
  - Production config
  - Docker config

#### 3.2 Dependency Injection
- [ ] **Service container** (`internal/container/container.go`)
  ```go
  type Container struct {
      Config     *config.Config
      DB         *pgxpool.Pool
      Repository repository.SchemaRepo
      Services   *service.Services
      Handlers   *handlers.Handlers
  }
  ```

#### 3.3 Main Application (`cmd/main.go`)
- [ ] **Application bootstrap**
  - Configuration loading
  - Database connection
  - Service initialization
  - Route setup
  - Graceful shutdown

### 4. Middleware Implementation
**Priority: MEDIUM** | **Estimated Time: 1 day**

#### 4.1 Authentication Middleware (`internal/middleware/auth.go`)
- [ ] **JWT token validation**
- [ ] **User context injection**
- [ ] **Role-based access control**

#### 4.2 Validation Middleware (`internal/middleware/validation.go`)
- [ ] **Request body validation**
- [ ] **Parameter validation**
- [ ] **Custom validation rules**

#### 4.3 Logging Middleware (`internal/middleware/logging.go`)
- [ ] **Request/response logging**
- [ ] **Performance metrics**
- [ ] **Error tracking**

### 5. Testing Implementation
**Priority: MEDIUM** | **Estimated Time: 2-3 days**

#### 5.1 Unit Tests
- [ ] **Repository tests** (`internal/repository/*_test.go`)
  - Mock database tests
  - Integration tests with test database
  - Error scenario testing

- [ ] **Service tests** (`internal/service/*_test.go`)
  - Business logic testing
  - Mock repository tests
  - Edge case handling

- [ ] **Handler tests** (`internal/handlers/*_test.go`)
  - HTTP endpoint testing
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