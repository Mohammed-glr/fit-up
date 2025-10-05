# Coach Management System - Implementation Requirements

## Overzicht
Dit document beschrijft alle benodigde wijzigingen en toevoegingen om het schema service geschikt te maken voor een coach-managed model, waarbij coaches via een aparte beheerdersomgeving sportschema's op maat kunnen maken voor medewerkers.

---

## 📋 Huidige Situatie

### Wat Werkt Goed ✅
- Automatische plan generatie met adaptieve algoritmes
- Uitgebreide personalisatie (fitness level, goals, equipment)
- Complete exercise library en workout templates
- Progress tracking en analytics
- PDF export functionaliteit
- Workout session management

### Wat Ontbreekt ❌
- Coach/Admin role system
- Handmatige schema creatie door coaches
- Coach-to-user assignment systeem
- Aparte coach beheerdersomgeving
- Authorization en permission checks
- Coach dashboard met client overzicht

---

## 🎯 Implementatie Roadmap

### **FASE 1: Role & Permission System (PRIORITEIT HOOG)**

#### 1.1 User Roles Toevoegen

**BELANGRIJK:** User data (email, password, profile) blijft in Auth Service.
Schema Service heeft alleen role cache voor authorization checks.

**Bestand:** `internal/schema/types/types.go`

```go
// User Roles (lokaal voor authorization)
type UserRole string

const (
	RoleUser   UserRole = "user"   // Normale gebruiker/medewerker
	RoleCoach  UserRole = "coach"  // Personal trainer/coach
	RoleAdmin  UserRole = "admin"  // System administrator
)

// GEEN volledige user struct - alleen role cache
type UserRoleCache struct {
	AuthUserID    string    `json:"auth_user_id" db:"auth_user_id"` // Van auth service
	Role          UserRole  `json:"role" db:"role"`
	LastSyncedAt  time.Time `json:"last_synced_at" db:"last_synced_at"`
}

// Voor display in dashboard - data komt van auth service API
type UserDisplayInfo struct {
	AuthUserID string `json:"auth_user_id"` // Reference only
	FirstName  string `json:"first_name"`   // Van auth service API
	LastName   string `json:"last_name"`    // Van auth service API
	Email      string `json:"email"`        // Van auth service API
}
```

#### 1.2 Database Migratie

**BELANGRIJK:** Geen volledige user table - alleen role cache.
User data blijft in Auth Service database.

**Bestand:** `migrations/XXXXXX_add_user_roles_cache.sql`

```sql
-- Role cache voor authorization checks (GEEN volledige user data)
CREATE TABLE IF NOT EXISTS user_roles_cache (
    auth_user_id TEXT PRIMARY KEY,
    role TEXT NOT NULL DEFAULT 'user',
    last_synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_role CHECK (role IN ('user', 'coach', 'admin'))
);

CREATE INDEX idx_user_roles_cache_role ON user_roles_cache(role);
CREATE INDEX idx_user_roles_cache_synced ON user_roles_cache(last_synced_at);

-- OPTIONEEL: Trigger voor auto-sync
CREATE OR REPLACE FUNCTION update_role_sync_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_synced_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_role_sync
BEFORE UPDATE ON user_roles_cache
FOR EACH ROW
EXECUTE FUNCTION update_role_sync_timestamp();

-- NOTE: Sync roles via:
-- 1. Auth service API calls wanneer nodig
-- 2. Event-based updates van auth service
-- 3. Periodieke sync job (bijv. elke 5 minuten voor stale data)
```

#### 1.3 Auth Middleware voor Coach Role

**Bestand:** `internal/schema/middleware/auth_middleware.go` (NIEUW)

```go
package middleware

import (
	"context"
	"net/http"
	"strings"
	
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type AuthMiddleware struct {
	// dependency injection voor user role lookup
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// Middleware om te checken of user een coach is
func (m *AuthMiddleware) RequireCoachRole() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := r.Header.Get("X-User-ID")
			if authUserID == "" {
				http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
				return
			}
			
			// TODO: Lookup user role from database
			// userRole, err := m.getUserRole(r.Context(), authUserID)
			// if err != nil || (userRole != types.RoleCoach && userRole != types.RoleAdmin) {
			//     http.Error(w, "Forbidden: Coach role required", http.StatusForbidden)
			//     return
			// }
			
			ctx := context.WithValue(r.Context(), "userRole", types.RoleCoach)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware om te checken of user admin is
func (m *AuthMiddleware) RequireAdminRole() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := r.Header.Get("X-User-ID")
			if authUserID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			
			// TODO: Check admin role
			ctx := context.WithValue(r.Context(), "userRole", types.RoleAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware om eigenaarschap van resource te valideren
func (m *AuthMiddleware) ValidateResourceOwnership(resourceType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Check if user owns resource or is coach/admin
			next.ServeHTTP(w, r)
		})
	}
}
```

---

### **FASE 2: Coach-to-User Assignment System (PRIORITEIT HOOG)**

#### 2.1 Coach Assignment Types

**Bestand:** `internal/schema/types/types.go` (toevoegen)

```go
// =============================================================================
// COACH MANAGEMENT TYPES
// =============================================================================

// Coach-to-User assignment
type CoachAssignment struct {
	AssignmentID  int       `json:"assignment_id" db:"assignment_id"`
	CoachID       string    `json:"coach_id" db:"coach_id"`           // Auth user ID van coach
	UserID        int       `json:"user_id" db:"user_id"`             // Internal user ID
	AssignedAt    time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy    string    `json:"assigned_by" db:"assigned_by"`     // Admin die assignment deed
	IsActive      bool      `json:"is_active" db:"is_active"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty" db:"deactivated_at"`
	Notes         string    `json:"notes" db:"notes"`
}

type CoachAssignmentRequest struct {
	CoachID string `json:"coach_id" validate:"required"`
	UserID  int    `json:"user_id" validate:"required"`
	Notes   string `json:"notes"`
}

// Client info voor coach dashboard
// NOTE: FirstName, LastName, Email komen van Auth Service API
type ClientSummary struct {
	UserID            int        `json:"user_id"`
	AuthID            string     `json:"auth_id"`
	FirstName         string     `json:"first_name"`      // Van Auth Service API
	LastName          string     `json:"last_name"`       // Van Auth Service API
	Email             string     `json:"email"`           // Van Auth Service API
	AssignedAt        time.Time  `json:"assigned_at"`     // Lokaal
	CurrentSchemaID   *int       `json:"current_schema_id,omitempty"` // Lokaal
	ActiveGoals       int        `json:"active_goals"`    // Lokaal
	CompletionRate    float64    `json:"completion_rate"` // Lokaal
	LastWorkoutDate   *time.Time `json:"last_workout_date,omitempty"` // Lokaal
	TotalWorkouts     int        `json:"total_workouts"`  // Lokaal
	CurrentStreak     int        `json:"current_streak"`  // Lokaal
	FitnessLevel      string     `json:"fitness_level"`   // Lokaal
}

// Coach dashboard overzicht
type CoachDashboard struct {
	CoachID           string          `json:"coach_id"`
	TotalClients      int             `json:"total_clients"`
	ActiveClients     int             `json:"active_clients"`
	ActiveSchemas     int             `json:"active_schemas"`
	TotalWorkouts     int             `json:"total_workouts_this_month"`
	AverageCompletion float64         `json:"average_completion_rate"`
	Clients           []ClientSummary `json:"clients"`
	RecentActivity    []CoachActivity `json:"recent_activity"`
}

type CoachActivity struct {
	ActivityID   int       `json:"activity_id"`
	ActivityType string    `json:"activity_type"` // "schema_created", "schema_updated", "client_assigned"
	UserID       int       `json:"user_id"`
	UserName     string    `json:"user_name"`
	Description  string    `json:"description"`
	Timestamp    time.Time `json:"timestamp"`
}
```

#### 2.2 Database Migratie

**Bestand:** `migrations/XXXXXX_add_coach_assignments.sql`

```sql
-- Coach assignments table
CREATE TABLE IF NOT EXISTS coach_assignments (
    assignment_id SERIAL PRIMARY KEY,
    coach_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    deactivated_at TIMESTAMP,
    notes TEXT,
    
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES fitness_profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT unique_active_assignment UNIQUE (coach_id, user_id, is_active)
);

CREATE INDEX idx_coach_assignments_coach_id ON coach_assignments(coach_id);
CREATE INDEX idx_coach_assignments_user_id ON coach_assignments(user_id);
CREATE INDEX idx_coach_assignments_active ON coach_assignments(is_active);

-- Coach activity log
CREATE TABLE IF NOT EXISTS coach_activity_log (
    activity_id SERIAL PRIMARY KEY,
    coach_id TEXT NOT NULL,
    user_id INTEGER,
    activity_type TEXT NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_activity_type CHECK (activity_type IN (
        'schema_created', 'schema_updated', 'schema_deleted',
        'client_assigned', 'client_removed', 'goal_created',
        'assessment_created', 'note_added'
    ))
);

CREATE INDEX idx_coach_activity_coach_id ON coach_activity_log(coach_id);
CREATE INDEX idx_coach_activity_timestamp ON coach_activity_log(created_at);
```

#### 2.3 Repository Interface & Implementation

**Bestand:** `internal/schema/repository/interfaces.go` (toevoegen)

```go
// Coach Assignment Repository
type CoachAssignmentRepo interface {
	// Assignments
	CreateCoachAssignment(ctx context.Context, assignment *types.CoachAssignmentRequest) (*types.CoachAssignment, error)
	GetCoachAssignment(ctx context.Context, assignmentID int) (*types.CoachAssignment, error)
	GetClientsByCoachID(ctx context.Context, coachID string) ([]types.ClientSummary, error)
	GetCoachByUserID(ctx context.Context, userID int) (*types.CoachAssignment, error)
	DeactivateAssignment(ctx context.Context, assignmentID int) error
	IsCoachForUser(ctx context.Context, coachID string, userID int) (bool, error)
	
	// Coach Dashboard
	GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error)
	
	// Activity Logging
	LogCoachActivity(ctx context.Context, activity *types.CoachActivity) error
	GetCoachActivityLog(ctx context.Context, coachID string, limit int) ([]types.CoachActivity, error)
}
```

**Bestand:** `internal/schema/repository/coach_assignment_repo.go` (NIEUW)

```go
package repository

import (
	"context"
	"fmt"
	
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateCoachAssignment(ctx context.Context, req *types.CoachAssignmentRequest) (*types.CoachAssignment, error) {
	query := `
		INSERT INTO coach_assignments (coach_id, user_id, assigned_by, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING assignment_id, coach_id, user_id, assigned_at, assigned_by, is_active, notes
	`
	
	var assignment types.CoachAssignment
	err := s.db.QueryRow(ctx, query,
		req.CoachID,
		req.UserID,
		req.CoachID, // For now, coach assigns themselves
		req.Notes,
	).Scan(
		&assignment.AssignmentID,
		&assignment.CoachID,
		&assignment.UserID,
		&assignment.AssignedAt,
		&assignment.AssignedBy,
		&assignment.IsActive,
		&assignment.Notes,
	)
	
	return &assignment, err
}

func (s *Store) GetClientsByCoachID(ctx context.Context, coachID string) ([]types.ClientSummary, error) {
	query := `
		SELECT 
			ca.user_id,
			wp.auth_user_id,
			'FirstName' as first_name,  -- TODO: Add user profile table
			'LastName' as last_name,
			'email@example.com' as email,
			ca.assigned_at,
			ws.schema_id,
			COUNT(DISTINCT fg.goal_id) as active_goals,
			AVG(CASE WHEN sess.status = 'completed' THEN 1.0 ELSE 0.0 END) as completion_rate,
			MAX(sess.start_time) as last_workout_date,
			COUNT(DISTINCT sess.session_id) as total_workouts,
			0 as current_streak,  -- TODO: Calculate streak
			wp.level as fitness_level
		FROM coach_assignments ca
		LEFT JOIN workout_profiles wp ON ca.user_id = wp.workout_profile_id
		LEFT JOIN weekly_schemas ws ON ca.user_id = ws.user_id AND ws.active = true
		LEFT JOIN fitness_goals fg ON ca.user_id = fg.user_id AND fg.is_active = true
		LEFT JOIN workout_sessions sess ON ca.user_id = sess.user_id
		WHERE ca.coach_id = $1 AND ca.is_active = true
		GROUP BY ca.user_id, wp.auth_user_id, ca.assigned_at, ws.schema_id, wp.level
		ORDER BY ca.assigned_at DESC
	`
	
	rows, err := s.db.Query(ctx, query, coachID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var clients []types.ClientSummary
	for rows.Next() {
		var client types.ClientSummary
		err := rows.Scan(
			&client.UserID,
			&client.AuthID,
			&client.FirstName,
			&client.LastName,
			&client.Email,
			&client.AssignedAt,
			&client.CurrentSchemaID,
			&client.ActiveGoals,
			&client.CompletionRate,
			&client.LastWorkoutDate,
			&client.TotalWorkouts,
			&client.CurrentStreak,
			&client.FitnessLevel,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	
	return clients, nil
}

func (s *Store) IsCoachForUser(ctx context.Context, coachID string, userID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM coach_assignments 
			WHERE coach_id = $1 AND user_id = $2 AND is_active = true
		)
	`
	
	var exists bool
	err := s.db.QueryRow(ctx, query, coachID, userID).Scan(&exists)
	return exists, err
}

func (s *Store) GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error) {
	// Implementation for complete dashboard data
	// TODO: Implement full dashboard aggregation
	
	clients, err := s.GetClientsByCoachID(ctx, coachID)
	if err != nil {
		return nil, err
	}
	
	dashboard := &types.CoachDashboard{
		CoachID:      coachID,
		TotalClients: len(clients),
		Clients:      clients,
	}
	
	return dashboard, nil
}

func (s *Store) LogCoachActivity(ctx context.Context, activity *types.CoachActivity) error {
	query := `
		INSERT INTO coach_activity_log (coach_id, user_id, activity_type, description)
		VALUES ($1, $2, $3, $4)
	`
	
	_, err := s.db.Exec(ctx, query,
		// TODO: Add coach_id to activity
		activity.UserID,
		activity.ActivityType,
		activity.Description,
	)
	
	return err
}
```

---

### **FASE 3: Handmatige Schema Creatie (PRIORITEIT HOOG)**

#### 3.1 Manual Schema Types

**Bestand:** `internal/schema/types/types.go` (toevoegen)

```go
// =============================================================================
// MANUAL SCHEMA CREATION TYPES
// =============================================================================

// Manual schema request van coach
type ManualSchemaRequest struct {
	UserID      int                     `json:"user_id" validate:"required"`
	CoachID     string                  `json:"coach_id" validate:"required"`
	Name        string                  `json:"name" validate:"required,min=3,max=100"`
	Description string                  `json:"description" validate:"max=500"`
	StartDate   time.Time               `json:"start_date" validate:"required"`
	EndDate     *time.Time              `json:"end_date"`
	IsTemplate  bool                    `json:"is_template"`
	Workouts    []ManualWorkoutRequest  `json:"workouts" validate:"required,min=1"`
}

type ManualWorkoutRequest struct {
	DayOfWeek    int                      `json:"day_of_week" validate:"required,min=1,max=7"`
	WorkoutName  string                   `json:"workout_name" validate:"required"`
	Focus        string                   `json:"focus" validate:"required"`
	Notes        string                   `json:"notes"`
	EstimatedMin int                      `json:"estimated_minutes"`
	Exercises    []ManualExerciseRequest  `json:"exercises" validate:"required,min=1"`
}

type ManualExerciseRequest struct {
	ExerciseID    int    `json:"exercise_id" validate:"required"`
	Sets          int    `json:"sets" validate:"required,min=1,max=10"`
	Reps          string `json:"reps" validate:"required"`
	RestSeconds   int    `json:"rest_seconds" validate:"required,min=0,max=600"`
	Weight        string `json:"weight"`           // "60kg", "BW", "RPE 7"
	Tempo         string `json:"tempo"`            // "3-0-1-0"
	Notes         string `json:"notes"`
	OrderIndex    int    `json:"order_index"`
	IsSuperSet    bool   `json:"is_superset"`
	SuperSetGroup int    `json:"superset_group"`
}

// Schema metadata voor tracking
type SchemaMetadata struct {
	CreatedBy      string                 `json:"created_by"`       // "coach", "system", "user"
	CreatorID      string                 `json:"creator_id"`       // Auth ID van creator
	IsCustom       bool                   `json:"is_custom"`        // Handmatig vs automatisch
	BaseTemplateID *int                   `json:"base_template_id"` // Als van template
	LastModifiedBy string                 `json:"last_modified_by"`
	ModifiedAt     *time.Time             `json:"modified_at"`
	Version        int                    `json:"version"`
	Tags           []string               `json:"tags"`
	CustomData     map[string]interface{} `json:"custom_data"`
}

// Extended weekly schema met coach info
type WeeklySchemaExtended struct {
	WeeklySchema
	CoachID      *string         `json:"coach_id,omitempty"`
	CoachName    string          `json:"coach_name,omitempty"`
	Metadata     SchemaMetadata  `json:"metadata"`
	Workouts     []WorkoutDetail `json:"workouts"`
}

type WorkoutDetail struct {
	Workout
	Exercises     []WorkoutExerciseDetail `json:"exercises"`
	EstimatedMin  int                     `json:"estimated_minutes"`
	Notes         string                  `json:"notes"`
}
```

#### 3.2 Database Migratie

**Bestand:** `migrations/XXXXXX_add_manual_schema_support.sql`

```sql
-- Add metadata columns to weekly_schemas
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS coach_id TEXT;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS created_by TEXT DEFAULT 'system';
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS is_custom BOOLEAN DEFAULT FALSE;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS base_template_id INTEGER;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS metadata JSONB;
ALTER TABLE weekly_schemas ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;

-- Add enhanced workout fields
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS workout_name TEXT;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS estimated_minutes INTEGER;

-- Add enhanced workout exercise fields
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS weight TEXT;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS tempo TEXT;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS order_index INTEGER DEFAULT 0;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS is_superset BOOLEAN DEFAULT FALSE;
ALTER TABLE workout_exercises ADD COLUMN IF NOT EXISTS superset_group INTEGER;

CREATE INDEX idx_weekly_schemas_coach_id ON weekly_schemas(coach_id);
CREATE INDEX idx_weekly_schemas_created_by ON weekly_schemas(created_by);
```

#### 3.3 Service Layer

**Bestand:** `internal/schema/services/coach_service.go` (NIEUW)

```go
package service

import (
	"context"
	"fmt"
	"time"
	
	"github.com/go-playground/validator/v10"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type CoachService interface {
	// Coach Management
	AssignClientToCoach(ctx context.Context, req *types.CoachAssignmentRequest) (*types.CoachAssignment, error)
	GetCoachClients(ctx context.Context, coachID string) ([]types.ClientSummary, error)
	GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error)
	RemoveClientFromCoach(ctx context.Context, assignmentID int) error
	
	// Manual Schema Creation
	CreateManualSchemaForClient(ctx context.Context, coachID string, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error)
	UpdateManualSchema(ctx context.Context, coachID string, schemaID int, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error)
	DeleteSchema(ctx context.Context, coachID string, schemaID int) error
	CloneSchemaToClient(ctx context.Context, coachID string, sourceSchemaID int, targetUserID int) (*types.WeeklySchemaExtended, error)
	
	// Schema Templates voor Coaches
	SaveSchemaAsTemplate(ctx context.Context, coachID string, schemaID int, templateName string) error
	GetCoachTemplates(ctx context.Context, coachID string) ([]types.WorkoutTemplate, error)
	CreateSchemaFromCoachTemplate(ctx context.Context, coachID string, templateID int, userID int) (*types.WeeklySchemaExtended, error)
	
	// Client Progress Review
	GetClientProgress(ctx context.Context, coachID string, userID int) (*types.UserProgressSummary, error)
	AddCoachNote(ctx context.Context, coachID string, userID int, note string) error
	
	// Validation
	ValidateCoachPermission(ctx context.Context, coachID string, userID int) error
}

type coachServiceImpl struct {
	repo      repository.SchemaRepo
	validator *validator.Validate
}

func NewCoachService(repo repository.SchemaRepo) CoachService {
	return &coachServiceImpl{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *coachServiceImpl) CreateManualSchemaForClient(ctx context.Context, coachID string, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	// Check if coach is authorized for this user
	if err := s.ValidateCoachPermission(ctx, coachID, req.UserID); err != nil {
		return nil, err
	}
	
	// Create weekly schema
	schemaReq := &types.WeeklySchemaRequest{
		UserID:    req.UserID,
		WeekStart: req.StartDate,
	}
	
	schema, err := s.repo.Schemas().CreateWeeklySchema(ctx, schemaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}
	
	// Create workouts
	var workouts []types.Workout
	for _, workoutReq := range req.Workouts {
		workout := &types.WorkoutRequest{
			SchemaID:  schema.SchemaID,
			DayOfWeek: workoutReq.DayOfWeek,
			Focus:     workoutReq.Focus,
		}
		
		createdWorkout, err := s.repo.Workouts().CreateWorkout(ctx, workout)
		if err != nil {
			return nil, fmt.Errorf("failed to create workout: %w", err)
		}
		
		// Add exercises to workout
		for _, exReq := range workoutReq.Exercises {
			exerciseReq := &types.WorkoutExerciseRequest{
				WorkoutID:   createdWorkout.WorkoutID,
				ExerciseID:  exReq.ExerciseID,
				Sets:        exReq.Sets,
				Reps:        exReq.Reps,
				RestSeconds: exReq.RestSeconds,
			}
			
			_, err := s.repo.WorkoutExercises().CreateWorkoutExercise(ctx, exerciseReq)
			if err != nil {
				return nil, fmt.Errorf("failed to add exercise: %w", err)
			}
		}
		
		workouts = append(workouts, *createdWorkout)
	}
	
	// Log activity
	activity := &types.CoachActivity{
		ActivityType: "schema_created",
		UserID:       req.UserID,
		Description:  fmt.Sprintf("Created manual schema: %s", req.Name),
		Timestamp:    time.Now(),
	}
	_ = s.repo.CoachAssignments().LogCoachActivity(ctx, activity)
	
	// TODO: Return complete extended schema
	return &types.WeeklySchemaExtended{
		WeeklySchema: *schema,
		CoachID:      &coachID,
	}, nil
}

func (s *coachServiceImpl) ValidateCoachPermission(ctx context.Context, coachID string, userID int) error {
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, userID)
	if err != nil {
		return fmt.Errorf("failed to check coach permission: %w", err)
	}
	
	if !isCoach {
		return fmt.Errorf("coach %s is not authorized for user %d", coachID, userID)
	}
	
	return nil
}

func (s *coachServiceImpl) GetCoachClients(ctx context.Context, coachID string) ([]types.ClientSummary, error) {
	return s.repo.CoachAssignments().GetClientsByCoachID(ctx, coachID)
}

func (s *coachServiceImpl) GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error) {
	return s.repo.CoachAssignments().GetCoachDashboard(ctx, coachID)
}

// TODO: Implement remaining methods
```

---

### **FASE 4: Coach Handler/Routes (PRIORITEIT HOOG)**

#### 4.1 Coach Handler

**Bestand:** `internal/schema/handlers/coach_handler.go` (NIEUW)

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type CoachHandler struct {
	service service.CoachService
}

func NewCoachHandler(service service.CoachService) *CoachHandler {
	return &CoachHandler{
		service: service,
	}
}

// GET /api/v1/coach/dashboard
func (h *CoachHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	dashboard, err := h.service.GetCoachDashboard(r.Context(), coachID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, dashboard)
}

// GET /api/v1/coach/clients
func (h *CoachHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	clients, err := h.service.GetCoachClients(r.Context(), coachID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"clients": clients,
		"total":   len(clients),
	})
}

// POST /api/v1/coach/clients/{userID}/schemas
func (h *CoachHandler) CreateSchemaForClient(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	
	var req types.ManualSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	req.UserID = userID
	req.CoachID = coachID
	
	schema, err := h.service.CreateManualSchemaForClient(r.Context(), coachID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, schema)
}

// GET /api/v1/coach/clients/{userID}
func (h *CoachHandler) GetClientDetails(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	
	// Validate permission
	if err := h.service.ValidateCoachPermission(r.Context(), coachID, userID); err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized for this client")
		return
	}
	
	progress, err := h.service.GetClientProgress(r.Context(), coachID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, progress)
}

// PUT /api/v1/coach/schemas/{schemaID}
func (h *CoachHandler) UpdateSchema(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}
	
	var req types.ManualSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	schema, err := h.service.UpdateManualSchema(r.Context(), coachID, schemaID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, schema)
}

// DELETE /api/v1/coach/schemas/{schemaID}
func (h *CoachHandler) DeleteSchema(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}
	
	if err := h.service.DeleteSchema(r.Context(), coachID, schemaID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Schema deleted successfully",
	})
}

// POST /api/v1/coach/schemas/{schemaID}/clone
func (h *CoachHandler) CloneSchema(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	schemaID, err := strconv.Atoi(chi.URLParam(r, "schemaID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid schema ID")
		return
	}
	
	var req struct {
		TargetUserID int `json:"target_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	schema, err := h.service.CloneSchemaToClient(r.Context(), coachID, schemaID, req.TargetUserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, schema)
}

// POST /api/v1/coach/clients/assign
func (h *CoachHandler) AssignClient(w http.ResponseWriter, r *http.Request) {
	coachID := r.Header.Get("X-User-ID")
	if coachID == "" {
		respondWithError(w, http.StatusUnauthorized, "Coach ID not found")
		return
	}
	
	var req types.CoachAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	req.CoachID = coachID
	
	assignment, err := h.service.AssignClientToCoach(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, assignment)
}

func (h *CoachHandler) RegisterRoutes(r chi.Router) {
	r.Route("/coach", func(r chi.Router) {
		// TODO: Add coach auth middleware
		// r.Use(authMiddleware.RequireCoachRole())
		
		// Dashboard
		r.Get("/dashboard", h.GetDashboard)
		
		// Client Management
		r.Get("/clients", h.GetClients)
		r.Get("/clients/{userID}", h.GetClientDetails)
		r.Post("/clients/assign", h.AssignClient)
		
		// Schema Management
		r.Post("/clients/{userID}/schemas", h.CreateSchemaForClient)
		r.Put("/schemas/{schemaID}", h.UpdateSchema)
		r.Delete("/schemas/{schemaID}", h.DeleteSchema)
		r.Post("/schemas/{schemaID}/clone", h.CloneSchema)
		
		// Templates
		r.Get("/templates", h.GetTemplates)
		r.Post("/templates", h.SaveTemplate)
		r.Post("/templates/{templateID}/create-schema", h.CreateFromTemplate)
	})
}

// TODO: Implement remaining handlers
func (h *CoachHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotImplemented, "Not implemented yet")
}

func (h *CoachHandler) SaveTemplate(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotImplemented, "Not implemented yet")
}

func (h *CoachHandler) CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotImplemented, "Not implemented yet")
}
```

---

### **FASE 5: Frontend Dashboard (PRIORITEIT MIDDEL)**

#### 5.1 Coach Dashboard UI Vereisten

**Bestand:** `docs/COACH_DASHBOARD_UI_REQUIREMENTS.md` (NIEUW)

```markdown
# Coach Dashboard UI Requirements

## Dashboard Overview Page

### Key Metrics (Top Cards)
- Total Clients
- Active Schemas
- Completion Rate This Week
- Workouts Completed This Month

### Client List
- Searchable/Filterable table
- Columns:
  - Name
  - Current Schema
  - Completion Rate
  - Last Workout
  - Current Streak
  - Actions (View, Edit, Create Schema)

### Recent Activity Feed
- Schema created/updated
- Workouts completed by clients
- Goals achieved
- Notes added

### Quick Actions
- Create New Schema
- Assign New Client
- View Templates
- Generate Reports

## Client Detail Page

### Client Profile
- Personal info
- Fitness level
- Active goals
- Current limitations/injuries

### Current Schema
- Week view with workouts
- Exercise details
- Progress tracking
- Ability to edit inline

### Progress Charts
- Completion rate over time
- Workout frequency
- Exercise progress (weight/reps)
- Goal progress

### Actions
- Create New Schema
- Clone Existing Schema
- Add Note
- Schedule Assessment

## Schema Builder Page

### Drag & Drop Interface
- Exercise library sidebar
- Week view (7 days)
- Drag exercises to days
- Set sets/reps/rest inline
- Reorder exercises
- Superset grouping

### Exercise Search
- Filter by:
  - Muscle group
  - Equipment
  - Difficulty
  - Type (strength/cardio)
- Preview exercise details

### Save Options
- Save as draft
- Publish to client
- Save as template
- Clone from existing

## Templates Library

### Template Cards
- Template name/description
- Days per week
- Target fitness level
- Suitable goals
- Preview workouts

### Actions
- Use Template
- Edit Template
- Delete Template
- Share Template (future)
```

---

### **FASE 6: Testing & Documentation (PRIORITEIT LAAG)**

#### 6.1 Unit Tests

**Bestand:** `internal/schema/services/coach_service_test.go` (NIEUW)

```go
package service_test

import (
	"context"
	"testing"
	"time"
	
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func TestCreateManualSchemaForClient(t *testing.T) {
	// TODO: Implement tests
	t.Run("should create schema successfully", func(t *testing.T) {
		// Setup
		// Execute
		// Assert
	})
	
	t.Run("should fail when coach not authorized", func(t *testing.T) {
		// Setup
		// Execute
		// Assert
	})
	
	t.Run("should validate request properly", func(t *testing.T) {
		// Setup
		// Execute
		// Assert
	})
}

func TestValidateCoachPermission(t *testing.T) {
	// TODO: Implement tests
}
```

#### 6.2 Integration Tests

**Bestand:** `internal/schema/integration/coach_integration_test.go` (NIEUW)

```go
package integration_test

import (
	"testing"
)

func TestCoachWorkflow_EndToEnd(t *testing.T) {
	// 1. Assign client to coach
	// 2. Create manual schema
	// 3. Client completes workouts
	// 4. Coach reviews progress
	// 5. Coach updates schema
}
```

---

## 📋 Implementation Checklist

### Database
- [ ] Create `user_roles` table
- [ ] Create `coach_assignments` table
- [ ] Create `coach_activity_log` table
- [ ] Add metadata columns to `weekly_schemas`
- [ ] Add enhanced fields to `workouts` table
- [ ] Add enhanced fields to `workout_exercises` table
- [ ] Create indexes for performance

### Backend Types & Models
- [ ] Add `UserRole` types
- [ ] Add `CoachAssignment` types
- [ ] Add `ManualSchemaRequest` types
- [ ] Add `SchemaMetadata` types
- [ ] Add `CoachDashboard` types
- [ ] Add `ClientSummary` types

### Repository Layer
- [ ] Implement `CoachAssignmentRepo` interface
- [ ] Implement `CreateCoachAssignment`
- [ ] Implement `GetClientsByCoachID`
- [ ] Implement `IsCoachForUser`
- [ ] Implement `GetCoachDashboard`
- [ ] Implement `LogCoachActivity`
- [ ] Update existing repos with coach support

### Service Layer
- [ ] Create `CoachService` interface
- [ ] Implement `CreateManualSchemaForClient`
- [ ] Implement `UpdateManualSchema`
- [ ] Implement `ValidateCoachPermission`
- [ ] Implement `GetCoachClients`
- [ ] Implement `GetCoachDashboard`
- [ ] Implement `CloneSchemaToClient`
- [ ] Implement `SaveSchemaAsTemplate`

### Middleware
- [ ] Create `AuthMiddleware`
- [ ] Implement `RequireCoachRole`
- [ ] Implement `RequireAdminRole`
- [ ] Implement `ValidateResourceOwnership`

### Handlers & Routes
- [ ] Create `CoachHandler`
- [ ] Implement GET `/coach/dashboard`
- [ ] Implement GET `/coach/clients`
- [ ] Implement GET `/coach/clients/{userID}`
- [ ] Implement POST `/coach/clients/{userID}/schemas`
- [ ] Implement PUT `/coach/schemas/{schemaID}`
- [ ] Implement DELETE `/coach/schemas/{schemaID}`
- [ ] Implement POST `/coach/clients/assign`
- [ ] Implement POST `/coach/schemas/{schemaID}/clone`
- [ ] Wire up routes in main router

### Testing
- [ ] Unit tests voor CoachService
- [ ] Unit tests voor coach repository methods
- [ ] Integration tests voor coach workflow
- [ ] E2E tests voor complete user journey
- [ ] Load tests voor dashboard queries

### Documentation
- [ ] API documentation voor coach endpoints
- [ ] Coach dashboard UI requirements
- [ ] User guide voor coaches
- [ ] Permission model documentation
- [ ] Database schema documentation

### Frontend (Separate Project)
- [ ] Coach login/authentication
- [ ] Coach dashboard page
- [ ] Client list component
- [ ] Client detail page
- [ ] Schema builder UI
- [ ] Exercise library component
- [ ] Drag & drop workout builder
- [ ] Progress charts & analytics
- [ ] Template library UI

---

## 🎯 Prioriteit & Volgorde

### Week 1: Foundation
1. Database migrations (alle tabellen)
2. Types & models (alle coach types)
3. Auth middleware (role checking)

### Week 2: Core Functionality
4. Repository layer (coach assignments)
5. Basic coach service (permission checking)
6. Coach assignment endpoints

### Week 3: Schema Creation
7. Manual schema creation types
8. Manual schema service methods
9. Schema creation endpoints
10. Schema update/delete endpoints

### Week 4: Dashboard & Analytics
11. Dashboard queries & aggregations
12. Client summary generation
13. Activity logging
14. Dashboard endpoint

### Week 5: Advanced Features
15. Schema cloning
16. Template system
17. Bulk operations
18. Performance optimization

### Week 6: Testing & Polish
19. Unit tests
20. Integration tests
21. Documentation
22. Code review & refactoring

---

## 🚀 Quick Start Commands

```bash
# Create migrations
migrate create -ext sql -dir migrations -seq add_user_roles
migrate create -ext sql -dir migrations -seq add_coach_assignments
migrate create -ext sql -dir migrations -seq add_manual_schema_support

# Run migrations
migrate -path migrations -database "postgres://..." up

# Generate mocks for testing
mockgen -source=internal/schema/repository/interfaces.go -destination=internal/schema/mocks/mock_repo.go

# Run tests
go test ./internal/schema/... -v -cover

# Run specific test
go test ./internal/schema/services -run TestCreateManualSchema -v
```

---

## 📞 Support & Questions

Voor vragen over implementatie:
1. Check deze documentatie
2. Bekijk de API documentation in `/docs`
3. Review de existing code voor patterns
4. Raadpleeg het team

---

## 📝 Notes

- Alle coach actions moeten gelogd worden in `coach_activity_log`
- Permission checks zijn cruciaal - altijd valideren
- Use transactions voor schema creation (multiple tables)
- Consider caching voor dashboard queries
- Rate limiting voor coach endpoints
- Backup before any schema deletion
- Version control voor schema updates
- Audit trail voor compliance

---

**Laatst bijgewerkt:** 5 oktober 2025
**Status:** Implementation Pending
**Eigenaar:** Development Team
