# Fit-Up Server - API Routes Summary

**Last Updated:** October 8, 2025  
**Base URL:** `http://localhost:8080`  
**API Version:** v1

## Overview

All routes are now properly organized under `/api/v1` for consistency, except for the WebSocket endpoint which is at the root level.

---

## 🔐 Authentication Routes
**Base Path:** `/api/v1/auth`

### Public Routes
- `POST /api/v1/auth/login` - User login (rate-limited)
- `POST /api/v1/auth/register` - User registration (rate-limited)
- `POST /api/v1/auth/forgot-password` - Request password reset (rate-limited)
- `POST /api/v1/auth/reset-password` - Reset password with token (rate-limited)
- `GET /api/v1/auth/{username}` - Get user by username
- `POST /api/v1/auth/validate-token` - Validate JWT token
- `POST /api/v1/auth/refresh-token` - Refresh access token (rate-limited)
- `POST /api/v1/auth/logout` - Logout user

### OAuth Routes
- `POST /api/v1/auth/oauth/{provider}` - Initiate OAuth flow
- `GET /api/v1/auth/oauth/callback/{provider}` - OAuth callback

### Protected Routes (Require JWT)
- `POST /api/v1/auth/change-password` - Change password
- `POST /api/v1/auth/link/{provider}` - Link OAuth account
- `DELETE /api/v1/auth/unlink/{provider}` - Unlink OAuth account
- `GET /api/v1/auth/linked-accounts` - Get linked accounts

---

## 💪 Workout/Schema Routes
**Base Path:** `/api/v1`

### Exercise Routes (Public)
- `GET /api/v1/exercises` - List all exercises
- `GET /api/v1/exercises/{id}` - Get exercise by ID
- `POST /api/v1/exercises/filter` - Filter exercises
- `GET /api/v1/exercises/search` - Search exercises
- `GET /api/v1/exercises/muscle-group/{muscleGroup}` - Get by muscle group
- `GET /api/v1/exercises/equipment/{equipment}` - Get by equipment
- `GET /api/v1/exercises/recommended` - Get recommended exercises
- `GET /api/v1/exercises/most-used` - Get most used exercises
- `GET /api/v1/exercises/{id}/usage-stats` - Get exercise usage statistics

### Workout Routes (Protected)
- `GET /api/v1/workouts/{id}` - Get workout by ID
- `GET /api/v1/workouts/{id}/exercises` - Get workout with exercises

### Workout Session Routes (Protected)
- `POST /api/v1/workout-sessions/start` - Start workout session
- `POST /api/v1/workout-sessions/{sessionID}/complete` - Complete session
- `POST /api/v1/workout-sessions/{sessionID}/skip` - Skip workout
- `POST /api/v1/workout-sessions/{sessionID}/log-exercise` - Log exercise performance
- `GET /api/v1/workout-sessions/users/{userID}/active` - Get active session
- `GET /api/v1/workout-sessions/users/{userID}/history` - Get session history
- `GET /api/v1/workout-sessions/users/{userID}/metrics` - Get session metrics
- `GET /api/v1/workout-sessions/users/{userID}/weekly-stats` - Get weekly stats

### Fitness Profile Routes (Protected)
- `POST /api/v1/fitness-profile/users/{userID}/assessment` - Create fitness assessment
- `GET /api/v1/fitness-profile/users/{userID}` - Get fitness profile
- `PUT /api/v1/fitness-profile/users/{userID}/fitness-level` - Update fitness level
- `PUT /api/v1/fitness-profile/users/{userID}/goals` - Update fitness goals
- `POST /api/v1/fitness-profile/users/{userID}/1rm-estimate` - Estimate one-rep max
- `GET /api/v1/fitness-profile/users/{userID}/1rm-history` - Get 1RM history
- `POST /api/v1/fitness-profile/users/{userID}/movement-assessment` - Create movement assessment
- `GET /api/v1/fitness-profile/users/{userID}/movement-limitations` - Get movement limitations
- `POST /api/v1/fitness-profile/users/{userID}/workout-profile` - Create workout profile
- `GET /api/v1/fitness-profile/users/{userID}/workout-profile` - Get workout profile
- `POST /api/v1/fitness-profile/users/{userID}/fitness-goals` - Create fitness goal
- `GET /api/v1/fitness-profile/users/{userID}/active-goals` - Get active goals

### Plan Generation Routes (Protected)
- `POST /api/v1/plans` - Create plan generation
- `GET /api/v1/plans/users/{userID}/active` - Get active plan
- `GET /api/v1/plans/users/{userID}/history` - Get plan history
- `POST /api/v1/plans/{planID}/performance` - Track plan performance
- `GET /api/v1/plans/{planID}/download` - Download plan as PDF
- `POST /api/v1/plans/{planID}/regenerate` - Mark plan for regeneration

### Coach Routes (Protected - Coach Role Required)
- `GET /api/v1/coach/clients` - List coach's clients
- `GET /api/v1/coach/clients/{clientID}` - Get client details
- `POST /api/v1/coach/clients/{clientID}/assign-plan` - Assign plan to client
- Additional coach-specific routes...

---

## 💬 Message Routes
**Base Path:** `/api/v1` (Protected - Require JWT)

### Conversation Routes
- `POST /api/v1/conversations` - Create conversation
- `GET /api/v1/conversations` - List conversations
- `GET /api/v1/conversations/{conversation_id}` - Get conversation details
- `GET /api/v1/conversations/{conversation_id}/unread-count` - Get unread message count
- `GET /api/v1/conversations/{conversation_id}/messages` - Get messages in conversation
- `POST /api/v1/conversations/{conversation_id}/messages/read-all` - Mark all as read

### Message Routes
- `POST /api/v1/messages` - Send message
- `PUT /api/v1/messages/{message_id}` - Update message
- `DELETE /api/v1/messages/{message_id}` - Delete message
- `POST /api/v1/messages/{message_id}/read` - Mark message as read

### WebSocket Route
- `WS /ws` - WebSocket connection for real-time messaging (requires JWT token in query param)

---

## 🍽️ Food Tracker Routes
**Base Path:** `/api/v1/food-tracker`

### Recipe Routes (Public)
- `GET /api/v1/food-tracker/recipes/system` - List system recipes
- `GET /api/v1/food-tracker/recipes/system/{id}` - Get system recipe
- `GET /api/v1/food-tracker/recipes/search` - Search recipes

### Recipe Management Routes (Admin Only)
- `POST /api/v1/food-tracker/recipes/system` - Create system recipe
- `PUT /api/v1/food-tracker/recipes/system/{id}` - Update system recipe
- `DELETE /api/v1/food-tracker/recipes/system/{id}` - Delete system recipe

### User Recipe Routes (Protected)
- `GET /api/v1/food-tracker/recipes/user` - List user's recipes
- `POST /api/v1/food-tracker/recipes/user` - Create user recipe
- `GET /api/v1/food-tracker/recipes/user/{id}` - Get user recipe
- `PUT /api/v1/food-tracker/recipes/user/{id}` - Update user recipe
- `DELETE /api/v1/food-tracker/recipes/user/{id}` - Delete user recipe

### Favorite Routes (Protected)
- `GET /api/v1/food-tracker/recipes/favorites` - Get favorite recipes
- `PATCH /api/v1/food-tracker/recipes/favorites/{recipeID}` - Toggle favorite

### Food Log Routes (Protected)
- `POST /api/v1/food-tracker/food-logs` - Log food entry
- `POST /api/v1/food-tracker/food-logs/recipe` - Log recipe as meal
- `GET /api/v1/food-tracker/food-logs/date/{date}` - Get logs by date
- `GET /api/v1/food-tracker/food-logs/range` - Get logs by date range
- `GET /api/v1/food-tracker/food-logs/{id}` - Get food log entry
- `PUT /api/v1/food-tracker/food-logs/{id}` - Update food log
- `DELETE /api/v1/food-tracker/food-logs/{id}` - Delete food log

### Nutrition Routes (Protected)
- `GET /api/v1/food-tracker/nutrition/daily/{date}` - Get daily nutrition
- `GET /api/v1/food-tracker/nutrition/weekly` - Get weekly nutrition
- `GET /api/v1/food-tracker/nutrition/monthly` - Get monthly nutrition
- `GET /api/v1/food-tracker/nutrition/goals` - Get nutrition goals
- `POST /api/v1/food-tracker/nutrition/goals` - Create/update nutrition goals
- `GET /api/v1/food-tracker/nutrition/comparison/{date}` - Get nutrition comparison
- `GET /api/v1/food-tracker/nutrition/insights/{date}` - Get nutrition insights

---

## 🏥 Health Check
- `GET /health` - Server health check

---

## Route Analysis Summary

### ✅ Properly Structured Services
All services now follow the `/api/v1/{service}/*` pattern:

1. **Auth Service** → `/api/v1/auth/*`
2. **Schema/Workout Service** → `/api/v1/{exercises,workouts,workout-sessions,fitness-profile,plans,coach}/*`
3. **Message Service** → `/api/v1/{conversations,messages}/*`
4. **Food Tracker Service** → `/api/v1/food-tracker/*`

### 🔧 Integration Status
- ✅ Auth Service - Integrated and working
- ✅ Schema/Workout Service - Integrated and working
- ✅ Message Service - **NEWLY INTEGRATED** with WebSocket support
- ✅ Food Tracker Service - **NEWLY INTEGRATED** with consistent API versioning

### 🎯 Key Improvements Made
1. **Message Service Integration** - Added complete message and conversation handling with real-time WebSocket support
2. **Food Tracker Integration** - Moved from inconsistent `/food-tracker` to proper `/api/v1/food-tracker` structure
3. **Consistent API Versioning** - All services now under `/api/v1` except WebSocket
4. **Simple Ingredient Database** - Created in-memory nutrition database with 20+ common ingredients

### 🔒 Authentication Middleware
- Most routes properly protected with JWT authentication
- Admin routes properly protected with admin role check
- Coach routes properly protected with coach role check
- Rate limiting applied to sensitive auth endpoints

### 📝 Notes
- WebSocket endpoint (`/ws`) is intentionally at root level for simplicity
- All services use Chi router with proper middleware stacking
- Services are modular and follow clean architecture patterns
- Food tracker includes both system recipes (managed by admins) and user recipes
