# Auth Service Integration Guide

## Overzicht

De Schema Service heeft **GEEN** eigen user table. Alle user data (email, password, profile) blijft in de Auth Service. Deze guide beschrijft hoe Schema Service communiceert met Auth Service.

---

## 🏗️ Architectuur

```
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway                            │
└────────────┬───────────────────────────┬────────────────────┘
             │                           │
             ▼                           ▼
┌────────────────────────┐    ┌────────────────────────────┐
│    Auth Service        │    │    Schema Service          │
│  ┌──────────────────┐  │    │  ┌──────────────────────┐  │
│  │ users            │  │    │  │ workout_profiles     │  │
│  │ - user_id (PK)   │  │    │  │ - auth_user_id (FK)  │  │
│  │ - email          │  │    │  │ - level              │  │
│  │ - password_hash  │  │    │  │ - goals              │  │
│  │ - role           │  │    │  └──────────────────────┘  │
│  │ - first_name     │  │    │  ┌──────────────────────┐  │
│  │ - last_name      │  │    │  │ user_roles_cache     │  │
│  └──────────────────┘  │    │  │ - auth_user_id (PK)  │  │
│                        │    │  │ - role (cached)      │  │
│  Database: auth_db     │    │  │ - last_synced_at     │  │
└────────────────────────┘    │  └──────────────────────┘  │
                              │                            │
                              │  Database: schema_db       │
                              └────────────────────────────┘
```

---

## 🔑 Data Ownership

### Auth Service Owns:
- ✅ User authentication (login, logout, tokens)
- ✅ User profile (email, name, avatar)
- ✅ User roles (source of truth)
- ✅ Password management
- ✅ Email verification
- ✅ Account activation/deactivation

### Schema Service Owns:
- ✅ Workout profiles (fitness level, goals)
- ✅ Training schemas
- ✅ Exercise logs
- ✅ Progress tracking
- ✅ Coach assignments
- ✅ **Role cache** (voor performance)

---

## 🔄 Communication Patterns

### Pattern 1: HTTP REST Calls (AANBEVOLEN voor MVP)

**Auth Service exposes endpoints:**
```
GET  /api/v1/users/{userId}              - Get user profile
GET  /api/v1/users/{userId}/role         - Get user role
POST /api/v1/users/batch                 - Get multiple users
GET  /api/v1/coaches                     - List all coaches
```

**Schema Service client:**

```go
// internal/schema/clients/auth_client.go
package clients

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type AuthClient interface {
    GetUserInfo(ctx context.Context, authUserID string) (*UserInfo, error)
    GetUserRole(ctx context.Context, authUserID string) (string, error)
    GetUsersBatch(ctx context.Context, authUserIDs []string) ([]UserInfo, error)
    ListCoaches(ctx context.Context) ([]CoachInfo, error)
}

type UserInfo struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Role      string `json:"role"`
    IsActive  bool   `json:"is_active"`
}

type CoachInfo struct {
    UserID    string `json:"user_id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Specialty string `json:"specialty"`
}

type httpAuthClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewAuthClient(baseURL string) AuthClient {
    return &httpAuthClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *httpAuthClient) GetUserInfo(ctx context.Context, authUserID string) (*UserInfo, error) {
    url := fmt.Sprintf("%s/api/v1/users/%s", c.baseURL, authUserID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Add internal service auth header
    req.Header.Set("X-Service-Token", "internal-secret-token")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to call auth service: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("auth service returned %d", resp.StatusCode)
    }
    
    var userInfo UserInfo
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        return nil, err
    }
    
    return &userInfo, nil
}

func (c *httpAuthClient) GetUsersBatch(ctx context.Context, authUserIDs []string) ([]UserInfo, error) {
    // Batch request voor performance
    url := fmt.Sprintf("%s/api/v1/users/batch", c.baseURL)
    
    payload := map[string][]string{"user_ids": authUserIDs}
    jsonData, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Service-Token", "internal-secret-token")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var users []UserInfo
    if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
        return nil, err
    }
    
    return users, nil
}

func (c *httpAuthClient) GetUserRole(ctx context.Context, authUserID string) (string, error) {
    url := fmt.Sprintf("%s/api/v1/users/%s/role", c.baseURL, authUserID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    req.Header.Set("X-Service-Token", "internal-secret-token")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var result struct {
        Role string `json:"role"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    
    return result.Role, nil
}
```

**Usage in Coach Service:**

```go
// internal/schema/services/coach_service.go
func (s *coachServiceImpl) GetClientSummary(ctx context.Context, authUserID string) (*types.ClientSummary, error) {
    // 1. Get local fitness data
    profile, err := s.repo.GetWorkoutProfile(ctx, authUserID)
    if err != nil {
        return nil, err
    }
    
    schema, _ := s.repo.GetActiveSchema(ctx, authUserID)
    completion, _ := s.repo.GetCompletionRate(ctx, authUserID)
    
    // 2. Get user info from auth service
    userInfo, err := s.authClient.GetUserInfo(ctx, authUserID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }
    
    // 3. Combine data
    return &types.ClientSummary{
        AuthID:        authUserID,
        FirstName:     userInfo.FirstName,
        LastName:      userInfo.LastName,
        Email:         userInfo.Email,
        CurrentSchema: schema,
        FitnessLevel:  string(profile.Level),
        CompletionRate: completion,
    }, nil
}

func (s *coachServiceImpl) GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error) {
    // 1. Get coach's clients from local DB
    assignments, err := s.repo.GetCoachAssignments(ctx, coachID)
    if err != nil {
        return nil, err
    }
    
    // 2. Get user IDs
    var userIDs []string
    for _, assignment := range assignments {
        userIDs = append(userIDs, assignment.AuthUserID)
    }
    
    // 3. Batch fetch user info from auth service (efficient!)
    usersInfo, err := s.authClient.GetUsersBatch(ctx, userIDs)
    if err != nil {
        return nil, err
    }
    
    // 4. Create lookup map
    userMap := make(map[string]*clients.UserInfo)
    for i := range usersInfo {
        userMap[usersInfo[i].UserID] = &usersInfo[i]
    }
    
    // 5. Build client summaries
    var clients []types.ClientSummary
    for _, assignment := range assignments {
        userInfo := userMap[assignment.AuthUserID]
        if userInfo == nil {
            continue
        }
        
        // Get local data
        completion, _ := s.repo.GetCompletionRate(ctx, assignment.AuthUserID)
        
        clients = append(clients, types.ClientSummary{
            AuthID:         assignment.AuthUserID,
            FirstName:      userInfo.FirstName,
            LastName:       userInfo.LastName,
            Email:          userInfo.Email,
            CompletionRate: completion,
            // ... other local data
        })
    }
    
    return &types.CoachDashboard{
        CoachID: coachID,
        Clients: clients,
    }, nil
}
```

---

### Pattern 2: Role Caching met Sync

**Voor frequente role checks - cache lokaal:**

```go
// internal/schema/services/role_sync_service.go
package service

type RoleSyncService interface {
    SyncUserRole(ctx context.Context, authUserID string) error
    GetCachedRole(ctx context.Context, authUserID string) (string, error)
    SyncAllRoles(ctx context.Context) error
}

type roleSyncServiceImpl struct {
    repo       repository.SchemaRepo
    authClient clients.AuthClient
}

func (s *roleSyncServiceImpl) SyncUserRole(ctx context.Context, authUserID string) error {
    // Get role from auth service
    role, err := s.authClient.GetUserRole(ctx, authUserID)
    if err != nil {
        return err
    }
    
    // Update cache
    return s.repo.UpdateRoleCache(ctx, authUserID, role)
}

func (s *roleSyncServiceImpl) GetCachedRole(ctx context.Context, authUserID string) (string, error) {
    // Check cache first
    cached, err := s.repo.GetCachedRole(ctx, authUserID)
    if err == nil && cached != nil {
        // Check if cache is fresh (< 5 minutes)
        if time.Since(cached.LastSyncedAt) < 5*time.Minute {
            return cached.Role, nil
        }
    }
    
    // Cache miss or stale - fetch from auth service
    role, err := s.authClient.GetUserRole(ctx, authUserID)
    if err != nil {
        return "", err
    }
    
    // Update cache
    _ = s.repo.UpdateRoleCache(ctx, authUserID, role)
    
    return role, nil
}

// Periodieke sync job (run elke 5 minuten)
func (s *roleSyncServiceImpl) SyncAllRoles(ctx context.Context) error {
    // Get all users with stale cache
    staleUsers, err := s.repo.GetStaleRoleCaches(ctx, 5*time.Minute)
    if err != nil {
        return err
    }
    
    for _, user := range staleUsers {
        _ = s.SyncUserRole(ctx, user.AuthUserID)
    }
    
    return nil
}
```

**Repository implementation:**

```go
// internal/schema/repository/role_cache_repo.go
func (s *Store) GetCachedRole(ctx context.Context, authUserID string) (*types.UserRoleCache, error) {
    query := `
        SELECT auth_user_id, role, last_synced_at
        FROM user_roles_cache
        WHERE auth_user_id = $1
    `
    
    var cache types.UserRoleCache
    err := s.db.QueryRow(ctx, query, authUserID).Scan(
        &cache.AuthUserID,
        &cache.Role,
        &cache.LastSyncedAt,
    )
    
    if err == pgx.ErrNoRows {
        return nil, nil
    }
    
    return &cache, err
}

func (s *Store) UpdateRoleCache(ctx context.Context, authUserID, role string) error {
    query := `
        INSERT INTO user_roles_cache (auth_user_id, role, last_synced_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (auth_user_id) 
        DO UPDATE SET role = $2, last_synced_at = NOW()
    `
    
    _, err := s.db.Exec(ctx, query, authUserID, role)
    return err
}

func (s *Store) GetStaleRoleCaches(ctx context.Context, staleDuration time.Duration) ([]types.UserRoleCache, error) {
    query := `
        SELECT auth_user_id, role, last_synced_at
        FROM user_roles_cache
        WHERE last_synced_at < NOW() - $1::interval
        ORDER BY last_synced_at ASC
        LIMIT 100
    `
    
    rows, err := s.db.Query(ctx, query, staleDuration)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var caches []types.UserRoleCache
    for rows.Next() {
        var cache types.UserRoleCache
        if err := rows.Scan(&cache.AuthUserID, &cache.Role, &cache.LastSyncedAt); err != nil {
            continue
        }
        caches = append(caches, cache)
    }
    
    return caches, nil
}
```

---

### Pattern 3: Event-Based Updates (ADVANCED - voor later)

**Auth Service publiceert events bij wijzigingen:**

```go
// Auth Service publiceert events
type UserRoleChangedEvent struct {
    UserID   string    `json:"user_id"`
    OldRole  string    `json:"old_role"`
    NewRole  string    `json:"new_role"`
    ChangedAt time.Time `json:"changed_at"`
}

// Schema Service luistert naar events
func (s *eventHandler) HandleUserRoleChanged(event UserRoleChangedEvent) error {
    // Update local cache immediately
    return s.repo.UpdateRoleCache(context.Background(), event.UserID, event.NewRole)
}
```

---

## 🔐 Authorization Middleware

**Update de auth middleware om role cache te gebruiken:**

```go
// internal/schema/middleware/auth_middleware.go
type AuthMiddleware struct {
    roleService service.RoleSyncService
}

func (m *AuthMiddleware) RequireCoachRole() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authUserID := r.Header.Get("X-User-ID")
            if authUserID == "" {
                http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
                return
            }
            
            // Check cached role
            role, err := m.roleService.GetCachedRole(r.Context(), authUserID)
            if err != nil {
                http.Error(w, "Failed to verify role", http.StatusInternalServerError)
                return
            }
            
            if role != "coach" && role != "admin" {
                http.Error(w, "Forbidden: Coach role required", http.StatusForbidden)
                return
            }
            
            ctx := context.WithValue(r.Context(), "userRole", role)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 📋 Implementatie Checklist

### In Auth Service (moet nog gemaakt worden):
- [ ] GET `/api/v1/users/{userId}` endpoint
- [ ] GET `/api/v1/users/{userId}/role` endpoint
- [ ] POST `/api/v1/users/batch` endpoint voor bulk fetching
- [ ] GET `/api/v1/coaches` endpoint voor alle coaches
- [ ] Internal service authentication (X-Service-Token header)

### In Schema Service:
- [ ] Create `internal/schema/clients/auth_client.go`
- [ ] Implement `GetUserInfo()`
- [ ] Implement `GetUsersBatch()` voor efficiency
- [ ] Implement `GetUserRole()`
- [ ] Implement `ListCoaches()`
- [ ] Create `user_roles_cache` table
- [ ] Create `RoleSyncService`
- [ ] Implement cache sync logic
- [ ] Update middleware om cache te gebruiken
- [ ] Add periodic sync job (cronjob/background task)
- [ ] Add config voor auth service URL

### Configuration:
```yaml
# config/config.yaml
auth_service:
  base_url: "http://auth-service:8080"
  timeout: 10s
  service_token: "secure-internal-token"
  
role_cache:
  ttl: 5m
  sync_interval: 5m
```

---

## 🚀 Startup Sequence

```go
// cmd/server/main.go
func main() {
    // 1. Init database
    db := initDatabase()
    
    // 2. Init auth client
    authClient := clients.NewAuthClient(config.AuthServiceURL)
    
    // 3. Init repositories
    repo := repository.NewStore(db)
    
    // 4. Init services
    roleSyncService := service.NewRoleSyncService(repo, authClient)
    coachService := service.NewCoachService(repo, authClient)
    
    // 5. Init middleware
    authMiddleware := middleware.NewAuthMiddleware(roleSyncService)
    
    // 6. Init handlers
    coachHandler := handlers.NewCoachHandler(coachService)
    
    // 7. Setup routes
    r := chi.NewRouter()
    r.Route("/coach", func(r chi.Router) {
        r.Use(authMiddleware.RequireCoachRole())
        coachHandler.RegisterRoutes(r)
    })
    
    // 8. Start sync job
    go startRoleSyncJob(roleSyncService)
    
    // 9. Start server
    http.ListenAndServe(":8081", r)
}

func startRoleSyncJob(service service.RoleSyncService) {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        _ = service.SyncAllRoles(context.Background())
    }
}
```

---

## 🎯 Best Practices

1. **Cache Intelligently**: Alleen role cachen, niet alle user data
2. **Batch Requests**: Gebruik batch endpoints voor dashboard (niet 100 aparte calls)
3. **Handle Failures**: Auth service down? Use stale cache met waarschuwing
4. **Circuit Breaker**: Implementeer circuit breaker voor auth service calls
5. **Timeouts**: Altijd timeouts op HTTP calls
6. **Retry Logic**: Exponential backoff voor transient failures
7. **Monitoring**: Log alle auth service calls voor debugging
8. **Security**: Gebruik internal service token voor service-to-service calls

---

## 📊 Performance Considerations

### Dashboard met 100 clients:
```go
// ❌ BAD - 100 separate calls
for _, client := range clients {
    userInfo, _ := authClient.GetUserInfo(ctx, client.AuthID)
    // Very slow!
}

// ✅ GOOD - 1 batch call
userIDs := extractUserIDs(clients)
usersInfo, _ := authClient.GetUsersBatch(ctx, userIDs)
// Fast!
```

### Role Checks:
```go
// ❌ BAD - Call auth service voor elke request
role, _ := authClient.GetUserRole(ctx, userID)

// ✅ GOOD - Use cache
role, _ := roleService.GetCachedRole(ctx, userID)
// Cache hit: <1ms, Cache miss: ~10ms
```

---

## 🔍 Troubleshooting

### Problem: Auth service is down
```go
func (s *coachServiceImpl) GetClientSummary(ctx context.Context, authUserID string) (*types.ClientSummary, error) {
    userInfo, err := s.authClient.GetUserInfo(ctx, authUserID)
    if err != nil {
        // Fallback: Use cached data or show "User info unavailable"
        return &types.ClientSummary{
            AuthID:    authUserID,
            FirstName: "User",
            LastName:  authUserID, // Fallback
            Email:     "unavailable@temp.com",
        }, nil
    }
    // ...
}
```

### Problem: Role cache is stale
- Sync job runs every 5 minutes
- Manual sync: `curl -X POST /api/v1/admin/sync-roles`
- Event-based updates (future)

---

## 📝 Summary

1. **GEEN user table in schema service** - alleen role cache
2. **Auth Service is source of truth** voor user data
3. **HTTP REST calls** voor user info (batch voor performance)
4. **Local role cache** voor authorization checks
5. **Periodic sync** om cache fresh te houden
6. **Graceful degradation** als auth service down is

**Laatst bijgewerkt:** 5 oktober 2025
