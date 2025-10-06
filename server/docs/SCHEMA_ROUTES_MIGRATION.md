# Schema Routes Migration

## Overview
All schema-related routes have been centralized in `/internal/schema/handlers/routes.go` and are now using the JWT authentication middleware from `/shared/middleware/auth_middleware.go`.

## Changes Made

### 1. Updated Auth Middleware (`/shared/middleware/auth_middleware.go`)
- Added JWT validation using `ValidateJWT` from auth service
- Added new middleware methods:
  - `RequireJWTAuth()` - Validates JWT token and adds claims to context
  - `OptionalJWTAuth()` - Validates JWT if present, allows request without it
- Added helper functions:
  - `extractTokenFromHeader()` - Extracts Bearer token from Authorization header
  - `validateJWTToken()` - Validates JWT and returns claims
  - `GetUserClaimsFromContext()` - Retrieves JWT claims from context
  - `GetAuthUserIDFromContext()` - Retrieves authenticated user ID from context

### 2. Created Centralized Schema Routes (`/internal/schema/handlers/routes.go`)
Created a new `SchemaRoutes` struct that holds all schema-related handlers and registers all routes:

**Public Routes (No Auth):**
- `GET /exercises` - List all exercises
- `GET /exercises/{id}` - Get exercise by ID
- `POST /exercises/filter` - Filter exercises
- `GET /exercises/search` - Search exercises
- `GET /exercises/muscle-group/{muscleGroup}` - Get exercises by muscle group
- `GET /exercises/equipment/{equipment}` - Get exercises by equipment
- `GET /exercises/recommended` - Get recommended exercises
- `GET /exercises/most-used` - Get most used exercises
- `GET /exercises/{id}/usage-stats` - Get exercise usage stats

**Authenticated Routes (Require JWT):**

**Workouts:**
- `GET /workouts/{id}` - Get workout by ID
- `GET /workouts/{id}/exercises` - Get workout with exercises

**Workout Sessions:**
- `POST /workout-sessions/start` - Start a new workout session
- `POST /workout-sessions/{sessionID}/complete` - Complete a session
- `POST /workout-sessions/{sessionID}/skip` - Skip a workout
- `POST /workout-sessions/{sessionID}/log-exercise` - Log exercise performance
- `GET /workout-sessions/users/{userID}/active` - Get active session
- `GET /workout-sessions/users/{userID}/history` - Get session history
- `GET /workout-sessions/users/{userID}/metrics` - Get session metrics
- `GET /workout-sessions/users/{userID}/weekly-stats` - Get weekly stats

**Fitness Profile:**
- `POST /fitness-profile/users/{userID}/assessment` - Create fitness assessment
- `GET /fitness-profile/users/{userID}` - Get user fitness profile
- `PUT /fitness-profile/users/{userID}/fitness-level` - Update fitness level
- `PUT /fitness-profile/users/{userID}/goals` - Update fitness goals
- `POST /fitness-profile/users/{userID}/1rm-estimate` - Estimate one-rep max
- `GET /fitness-profile/users/{userID}/1rm-history` - Get 1RM history
- `POST /fitness-profile/users/{userID}/movement-assessment` - Create movement assessment
- `GET /fitness-profile/users/{userID}/movement-limitations` - Get movement limitations
- `POST /fitness-profile/users/{userID}/workout-profile` - Create workout profile
- `GET /fitness-profile/users/{userID}/workout-profile` - Get workout profile
- `POST /fitness-profile/users/{userID}/fitness-goals` - Create fitness goal
- `GET /fitness-profile/users/{userID}/active-goals` - Get active goals

**Plans:**
- `POST /plans` - Create plan generation
- `GET /plans/users/{userID}/active` - Get active plan
- `GET /plans/users/{userID}/history` - Get plan history
- `POST /plans/{planID}/performance` - Track plan performance

**Coach Routes (Require Coach Role):**
- `GET /coach/dashboard` - Get coach dashboard
- `GET /coach/stats` - Get coach statistics
- `GET /coach/activity` - Get recent activity
- `GET /coach/clients` - Get all clients
- `POST /coach/clients/assign` - Assign a client
- `GET /coach/clients/{userID}` - Get client details
- `DELETE /coach/clients/{assignmentID}` - Remove client
- `GET /coach/clients/{userID}/progress` - Get client progress
- `GET /coach/clients/{userID}/workouts` - Get client workouts
- `GET /coach/clients/{userID}/schemas` - Get client schemas
- `POST /coach/clients/{userID}/notes` - Add client note
- `POST /coach/clients/{userID}/schemas` - Create schema for client
- `PUT /coach/schemas/{schemaID}` - Update schema
- `DELETE /coach/schemas/{schemaID}` - Delete schema
- `POST /coach/schemas/{schemaID}/clone` - Clone schema
- `GET /coach/templates` - Get templates
- `POST /coach/templates` - Save template
- `POST /coach/templates/{templateID}/create-schema` - Create from template
- `DELETE /coach/templates/{templateID}` - Delete template

**Admin Routes (Require Admin Role):**
- Placeholder for future admin-specific routes

### 3. Updated Main Server (`/cmd/main.go`)
- Added imports for schema handlers and services
- Initialized all schema services (exercise, workout, session, profile, plan, coach)
- Created `SchemaRoutes` instance with all dependencies
- Registered schema routes in the main router at `/api/v1`
- Updated server startup logs to show all available endpoints

## How to Use JWT Authentication in Your Handlers

### 1. Extract User ID from Context
```go
import "github.com/tdmdh/fit-up-server/shared/middleware"

func (h *Handler) YourHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetAuthUserIDFromContext(r.Context())
    if !ok {
        // User ID not found - shouldn't happen if middleware is applied
        respondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    // Use userID...
}
```

### 2. Extract Full JWT Claims from Context
```go
import (
    "github.com/tdmdh/fit-up-server/shared/middleware"
    authTypes "github.com/tdmdh/fit-up-server/internal/auth/types"
)

func (h *Handler) YourHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    // Access claims
    userID := claims.UserID
    email := claims.Email
    role := claims.Role // authTypes.UserRole (admin, coach, user)
    
    // Use claims...
}
```

### 3. Check User Role
```go
func (h *Handler) YourHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    if claims.Role == authTypes.RoleAdmin {
        // Admin-specific logic
    } else if claims.Role == authTypes.RoleCoach {
        // Coach-specific logic
    }
}
```

## Testing

### Test Public Routes (No Auth Required)
```bash
curl http://localhost:8080/api/v1/exercises
```

### Test Authenticated Routes (With JWT)
```bash
# First, get a JWT token by logging in
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.access_token')

# Then use the token
curl http://localhost:8080/api/v1/workouts/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Test Coach Routes (Requires Coach Role)
```bash
# Login as coach
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"coach@example.com","password":"password"}' \
  | jq -r '.access_token')

# Access coach endpoints
curl http://localhost:8080/api/v1/coach/dashboard \
  -H "Authorization: Bearer $TOKEN"
```

## Benefits

1. **Centralized Route Management**: All schema routes are now in one place
2. **Proper JWT Authentication**: All routes use proper JWT validation
3. **Role-Based Access Control**: Coach and admin routes are protected by role middleware
4. **Consistent Error Handling**: JWT errors are handled consistently across all routes
5. **Type-Safe Claims**: Full access to JWT claims in handlers via context
6. **Easy to Extend**: Adding new routes is straightforward with the new structure
7. **Better Security**: Proper token validation on every request

## Migration Notes

- The old `X-User-ID` header approach is still supported by existing middleware methods for backward compatibility
- New routes should use JWT authentication (`RequireJWTAuth()`)
- Coach handlers are properly protected by `RequireCoachRole()` middleware
- Admin routes placeholder is ready for future admin functionality
