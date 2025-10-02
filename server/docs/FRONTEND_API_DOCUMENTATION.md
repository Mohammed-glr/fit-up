# Fit-Up API Documentation voor Frontend

**Versie:** 1.0  
**Base URL:** `http://localhost:8080/api/v1`  
**Laatste update:** Oktober 2, 2025

---

## 📋 Inhoudsopgave

1. [Authenticatie](#authenticatie)
2. [Fitness Profile](#fitness-profile)
3. [Exercises](#exercises)
4. [Workouts](#workouts)
5. [Workout Sessions](#workout-sessions)
6. [Plan Generation](#plan-generation)
7. [Types & Enums](#types--enums)
8. [Error Handling](#error-handling)

---

## 🔐 Authenticatie

Alle endpoints (behalve `/auth/*`) vereisen JWT authenticatie via de `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

### Auth Endpoints

| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| POST | `/auth/register` | Nieuwe gebruiker aanmaken |
| POST | `/auth/login` | Inloggen en JWT token ontvangen |
| POST | `/auth/logout` | Uitloggen (token invalideren) |
| GET | `/auth/me` | Huidige gebruiker ophalen |
| GET | `/auth/google` | OAuth Google login starten |
| GET | `/auth/google/callback` | OAuth Google callback |

---

## 💪 Fitness Profile

Endpoints voor het beheren van gebruikers fitness profiel, doelen en assessments.

### 1. **Create Fitness Assessment**

**Doel:** Een nieuwe fitness assessment aanmaken voor een gebruiker.  
**Use case:** Bij onboarding of periodieke evaluaties.

```http
POST /fitness-profile/{userID}/assessment
```

**Request Body:**
```json
{
  "assessment_date": "2025-10-02T10:00:00Z",
  "overall_level": "intermediate",
  "strength_level": 7,
  "endurance_level": 6,
  "flexibility_level": 5,
  "mobility_score": 6,
  "body_composition": {
    "weight": 75.5,
    "body_fat_percentage": 18.5,
    "muscle_mass_percentage": 42.0
  },
  "notes": "Goede vorm, focus op flexibiliteit"
}
```

**Response:** `201 Created`
```json
{
  "assessment_id": 1,
  "user_id": 123,
  "assessment_date": "2025-10-02T10:00:00Z",
  "overall_level": "intermediate",
  "strength_level": 7,
  "endurance_level": 6,
  "flexibility_level": 5,
  "mobility_score": 6,
  "body_composition": {...},
  "notes": "Goede vorm, focus op flexibiliteit"
}
```

**Frontend Implementatie:**
```typescript
async function createFitnessAssessment(userId: number, data: FitnessAssessmentRequest) {
  const response = await fetch(`/api/v1/fitness-profile/${userId}/assessment`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify(data)
  });
  
  if (!response.ok) throw new Error('Assessment aanmaken mislukt');
  return await response.json();
}
```

---

### 2. **Get User Fitness Profile**

**Doel:** Huidige fitness profiel van gebruiker ophalen.  
**Use case:** Dashboard, profiel pagina.

```http
GET /fitness-profile/{userID}
```

**Response:** `200 OK`
```json
{
  "user_id": 123,
  "current_level": "intermediate",
  "primary_goal": "muscle_gain",
  "secondary_goals": ["strength", "endurance"],
  "workout_frequency": 4,
  "available_equipment": ["barbell", "dumbbell", "machine"],
  "latest_assessment": {
    "assessment_id": 1,
    "assessment_date": "2025-10-02T10:00:00Z",
    "overall_level": "intermediate",
    "strength_level": 7,
    "endurance_level": 6
  },
  "goals_progress": [...]
}
```

**Frontend Implementatie:**
```typescript
async function getFitnessProfile(userId: number) {
  const response = await fetch(`/api/v1/fitness-profile/${userId}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  
  if (!response.ok) throw new Error('Profiel ophalen mislukt');
  return await response.json();
}
```

---

### 3. **Update Fitness Level**

**Doel:** Fitness level van gebruiker bijwerken.  
**Use case:** Na progressie of nieuwe assessment.

```http
PUT /fitness-profile/{userID}/level
```

**Request Body:**
```json
{
  "level": "advanced"
}
```

**Response:** `200 OK`
```json
{
  "message": "Fitness level updated successfully"
}
```

---

### 4. **Update Fitness Goals**

**Doel:** Fitness doelen van gebruiker aanpassen.  
**Use case:** Settings pagina, doel aanpassingen.

```http
PUT /fitness-profile/{userID}/goals
```

**Request Body:**
```json
[
  {
    "goal": "muscle_gain",
    "priority": 1,
    "target_date": "2026-01-01T00:00:00Z"
  },
  {
    "goal": "strength",
    "priority": 2,
    "target_date": "2026-06-01T00:00:00Z"
  }
]
```

**Response:** `200 OK`
```json
{
  "message": "Fitness goals updated successfully"
}
```

**Frontend Implementatie:**
```typescript
interface FitnessGoalTarget {
  goal: FitnessGoal;
  priority: number;
  target_date?: string;
}

async function updateFitnessGoals(userId: number, goals: FitnessGoalTarget[]) {
  const response = await fetch(`/api/v1/fitness-profile/${userId}/goals`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify(goals)
  });
  
  if (!response.ok) throw new Error('Goals bijwerken mislukt');
  return await response.json();
}
```

---

### 5. **Estimate One Rep Max (1RM)**

**Doel:** Schatting maken van maximale gewicht voor 1 herhaling.  
**Use case:** Kracht tracking, workout planning.

```http
POST /fitness-profile/{userID}/exercises/{exerciseID}/1rm
```

**Request Body:**
```json
{
  "weight": 80.0,
  "reps": 8,
  "rpe": 8.5
}
```

**Response:** `200 OK`
```json
{
  "estimated_1rm": 100.0,
  "confidence": 0.85,
  "formula_used": "epley",
  "recorded_at": "2025-10-02T14:30:00Z"
}
```

**Frontend Implementatie:**
```typescript
async function estimate1RM(userId: number, exerciseId: number, performance: {
  weight: number;
  reps: number;
  rpe?: number;
}) {
  const response = await fetch(
    `/api/v1/fitness-profile/${userId}/exercises/${exerciseId}/1rm`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify(performance)
    }
  );
  
  if (!response.ok) throw new Error('1RM berekenen mislukt');
  return await response.json();
}
```

---

### 6. **Get 1RM History**

**Doel:** Historische 1RM data ophalen voor progressie tracking.  
**Use case:** Progress charts, kracht ontwikkeling.

```http
GET /fitness-profile/{userID}/exercises/{exerciseID}/1rm/history
```

**Response:** `200 OK`
```json
[
  {
    "estimated_1rm": 100.0,
    "recorded_at": "2025-10-02T14:30:00Z",
    "confidence": 0.85
  },
  {
    "estimated_1rm": 95.0,
    "recorded_at": "2025-09-25T14:30:00Z",
    "confidence": 0.82
  }
]
```

---

### 7. **Get Recovery Metrics**

**Doel:** Herstel status van gebruiker ophalen.  
**Use case:** Bepalen of gebruiker klaar is voor training.

```http
GET /fitness-profile/{userID}/recovery
```

**Response:** `200 OK`
```json
{
  "user_id": 123,
  "date": "2025-10-02",
  "recovery_score": 7.5,
  "muscle_soreness": 3,
  "sleep_quality": 8,
  "stress_level": 4,
  "fatigue_level": 3,
  "ready_to_train": true,
  "recommended_intensity": "moderate"
}
```

---

## 🏋️ Exercises

Endpoints voor het ophalen en filteren van oefeningen.

### 1. **Get Exercise by ID**

**Doel:** Details van specifieke oefening ophalen.  
**Use case:** Oefening detail pagina.

```http
GET /exercises/{id}
```

**Response:** `200 OK`
```json
{
  "exercise_id": 1,
  "name": "Barbell Bench Press",
  "muscle_groups": ["chest", "triceps", "shoulders"],
  "difficulty": "intermediate",
  "equipment": "barbell",
  "type": "strength",
  "default_sets": 3,
  "default_reps": "8-12",
  "rest_seconds": 120,
  "description": "Klassieke compound oefening voor borst",
  "instructions": ["Lig op de bank...", "Pak de barbell..."],
  "video_url": "https://example.com/video.mp4"
}
```

**Frontend Implementatie:**
```typescript
async function getExercise(exerciseId: number) {
  const response = await fetch(`/api/v1/exercises/${exerciseId}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  
  if (!response.ok) throw new Error('Oefening ophalen mislukt');
  return await response.json();
}
```

---

### 2. **List Exercises (Paginated)**

**Doel:** Alle oefeningen ophalen met paginatie.  
**Use case:** Oefeningen library, browse pagina.

```http
GET /exercises?page=1&page_size=20
```

**Query Parameters:**
- `page` (optional, default: 1) - Pagina nummer
- `page_size` (optional, default: 20) - Items per pagina

**Response:** `200 OK`
```json
{
  "data": [
    {
      "exercise_id": 1,
      "name": "Barbell Bench Press",
      "muscle_groups": ["chest", "triceps"],
      "difficulty": "intermediate",
      "equipment": "barbell",
      "type": "strength"
    },
    ...
  ],
  "total_count": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

**Frontend Implementatie:**
```typescript
async function listExercises(page = 1, pageSize = 20) {
  const response = await fetch(
    `/api/v1/exercises?page=${page}&page_size=${pageSize}`,
    { headers: { 'Authorization': `Bearer ${getToken()}` } }
  );
  
  if (!response.ok) throw new Error('Oefeningen ophalen mislukt');
  return await response.json();
}
```

---

### 3. **Filter Exercises**

**Doel:** Oefeningen filteren op criteria.  
**Use case:** Zoekfunctie met filters.

```http
GET /exercises/filter?muscle_group=chest&equipment=barbell&difficulty=intermediate
```

**Query Parameters:**
- `muscle_group` (optional) - Spiergroep (chest, back, legs, etc.)
- `equipment` (optional) - Equipment type (barbell, dumbbell, bodyweight, etc.)
- `difficulty` (optional) - Moeilijkheidsgraad (beginner, intermediate, advanced)
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

**Response:** `200 OK` - Zelfde structuur als List Exercises

**Frontend Implementatie:**
```typescript
interface ExerciseFilters {
  muscle_group?: string;
  equipment?: string;
  difficulty?: string;
  page?: number;
  page_size?: number;
}

async function filterExercises(filters: ExerciseFilters) {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, value]) => {
    if (value !== undefined) params.append(key, value.toString());
  });
  
  const response = await fetch(`/api/v1/exercises/filter?${params}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  
  if (!response.ok) throw new Error('Filteren mislukt');
  return await response.json();
}
```

---

### 4. **Search Exercises**

**Doel:** Oefeningen zoeken op naam.  
**Use case:** Zoekbalk functionaliteit.

```http
GET /exercises/search?q=bench press
```

**Query Parameters:**
- `q` (required) - Zoekterm
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

**Response:** `200 OK` - Zelfde structuur als List Exercises

---

### 5. **Get Exercises by Muscle Group**

**Doel:** Alle oefeningen voor specifieke spiergroep.  
**Use case:** Spiergroep specifieke workout planning.

```http
GET /exercises/muscle-groups/{group}
```

**Voorbeeld:**
```http
GET /exercises/muscle-groups/chest
```

**Response:** `200 OK`
```json
[
  {
    "exercise_id": 1,
    "name": "Barbell Bench Press",
    "muscle_groups": ["chest", "triceps"],
    "difficulty": "intermediate",
    "equipment": "barbell"
  },
  ...
]
```

---

### 6. **Get Exercises by Equipment**

**Doel:** Alle oefeningen voor specifiek equipment.  
**Use case:** Workout aanpassen aan beschikbaar equipment.

```http
GET /exercises/equipment/{type}
```

**Voorbeeld:**
```http
GET /exercises/equipment/dumbbell
```

**Response:** `200 OK` - Array van oefeningen

---

### 7. **Get Recommended Exercises**

**Doel:** Gepersonaliseerde oefening aanbevelingen.  
**Use case:** AI-powered workout suggestions.

```http
GET /exercises/users/{userID}/recommended?count=10
```

**Query Parameters:**
- `count` (optional, default: 10) - Aantal aanbevelingen

**Response:** `200 OK`
```json
[
  {
    "exercise_id": 1,
    "name": "Barbell Bench Press",
    "muscle_groups": ["chest", "triceps"],
    "recommendation_score": 0.95,
    "reason": "Matches your fitness level and goals"
  },
  ...
]
```

**Frontend Implementatie:**
```typescript
async function getRecommendedExercises(userId: number, count = 10) {
  const response = await fetch(
    `/api/v1/exercises/users/${userId}/recommended?count=${count}`,
    { headers: { 'Authorization': `Bearer ${getToken()}` } }
  );
  
  if (!response.ok) throw new Error('Aanbevelingen ophalen mislukt');
  return await response.json();
}
```

---

### 8. **Get Most Used Exercises**

**Doel:** Populairste oefeningen ophalen.  
**Use case:** Trending oefeningen tonen.

```http
GET /exercises/popular?limit=10
```

**Query Parameters:**
- `limit` (optional, default: 10) - Aantal resultaten

**Response:** `200 OK`
```json
[
  {
    "exercise_id": 1,
    "name": "Barbell Bench Press",
    "usage_count": 1250,
    "muscle_groups": ["chest", "triceps"]
  },
  ...
]
```

---

## 📅 Workouts

Endpoints voor workout templates en schema's.

### 1. **Get Workout by ID**

**Doel:** Specifieke workout ophalen.  
**Use case:** Workout detail pagina.

```http
GET /workouts/{id}
```

**Response:** `200 OK`
```json
{
  "workout_id": 1,
  "schema_id": 10,
  "day_of_week": 1,
  "focus": "upper_body",
  "created_at": "2025-10-01T10:00:00Z"
}
```

---

### 2. **Get Workout with Exercises**

**Doel:** Workout met alle oefeningen ophalen.  
**Use case:** Workout uitvoeren, detail view.

```http
GET /workouts/{id}/full
```

**Response:** `200 OK`
```json
{
  "workout_id": 1,
  "schema_id": 10,
  "day_of_week": 1,
  "focus": "upper_body",
  "exercises": [
    {
      "we_id": 1,
      "sets": 3,
      "reps": "8-12",
      "rest_seconds": 90,
      "exercise": {
        "exercise_id": 1,
        "name": "Barbell Bench Press",
        "muscle_groups": ["chest", "triceps"],
        "difficulty": "intermediate",
        "equipment": "barbell",
        "type": "strength"
      }
    },
    {
      "we_id": 2,
      "sets": 3,
      "reps": "10-15",
      "rest_seconds": 60,
      "exercise": {
        "exercise_id": 5,
        "name": "Dumbbell Rows",
        "muscle_groups": ["back", "biceps"],
        "difficulty": "intermediate",
        "equipment": "dumbbell",
        "type": "strength"
      }
    }
  ]
}
```

**Frontend Implementatie:**
```typescript
async function getWorkoutWithExercises(workoutId: number) {
  const response = await fetch(`/api/v1/workouts/${workoutId}/full`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  
  if (!response.ok) throw new Error('Workout ophalen mislukt');
  return await response.json();
}
```

---

### 3. **Get Workouts by Schema ID**

**Doel:** Alle workouts van een schema ophalen.  
**Use case:** Weekschema overzicht.

```http
GET /schemas/{schemaID}/workouts
```

**Response:** `200 OK`
```json
[
  {
    "workout_id": 1,
    "schema_id": 10,
    "day_of_week": 1,
    "focus": "upper_body"
  },
  {
    "workout_id": 2,
    "schema_id": 10,
    "day_of_week": 3,
    "focus": "lower_body"
  },
  {
    "workout_id": 3,
    "schema_id": 10,
    "day_of_week": 5,
    "focus": "full_body"
  }
]
```

---

## 🏃 Workout Sessions

Endpoints voor het tracken van workout sessies en prestaties.

### 1. **Start Workout Session**

**Doel:** Nieuwe workout sessie starten.  
**Use case:** Begin van training.

```http
POST /sessions/start
```

**Request Body:**
```json
{
  "user_id": 123,
  "workout_id": 1,
  "started_at": "2025-10-02T10:00:00Z"
}
```

**Response:** `201 Created`
```json
{
  "session_id": 456,
  "user_id": 123,
  "workout_id": 1,
  "started_at": "2025-10-02T10:00:00Z",
  "status": "in_progress"
}
```

**Frontend Implementatie:**
```typescript
async function startWorkoutSession(userId: number, workoutId: number) {
  const response = await fetch('/api/v1/sessions/start', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({
      user_id: userId,
      workout_id: workoutId,
      started_at: new Date().toISOString()
    })
  });
  
  if (!response.ok) throw new Error('Sessie starten mislukt');
  return await response.json();
}
```

---

### 2. **Complete Workout Session**

**Doel:** Workout sessie afronden.  
**Use case:** Einde van training.

```http
POST /sessions/{sessionID}/complete
```

**Request Body:**
```json
{
  "completed_at": "2025-10-02T11:30:00Z",
  "notes": "Goede workout, voelde sterk"
}
```

**Response:** `200 OK`
```json
{
  "session_id": 456,
  "user_id": 123,
  "workout_id": 1,
  "started_at": "2025-10-02T10:00:00Z",
  "completed_at": "2025-10-02T11:30:00Z",
  "duration_minutes": 90,
  "status": "completed",
  "notes": "Goede workout, voelde sterk"
}
```

**Frontend Implementatie:**
```typescript
async function completeWorkoutSession(sessionId: number, notes?: string) {
  const response = await fetch(`/api/v1/sessions/${sessionId}/complete`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({
      completed_at: new Date().toISOString(),
      notes: notes
    })
  });
  
  if (!response.ok) throw new Error('Sessie afronden mislukt');
  return await response.json();
}
```

---

### 3. **Skip Workout**

**Doel:** Workout overslaan en reden registreren.  
**Use case:** Gebruiker kan niet trainen.

```http
POST /sessions/{sessionID}/skip
```

**Request Body:**
```json
{
  "reason": "Sick",
  "notes": "Griep, rust nemen"
}
```

**Response:** `200 OK`
```json
{
  "session_id": 456,
  "status": "skipped",
  "skip_reason": "Sick",
  "notes": "Griep, rust nemen"
}
```

---

### 4. **Log Exercise Performance**

**Doel:** Prestatie van oefening registreren tijdens sessie.  
**Use case:** Real-time tracking tijdens training.

```http
POST /sessions/{sessionID}/exercises/{exerciseID}/log
```

**Request Body:**
```json
{
  "set_number": 1,
  "reps_completed": 10,
  "weight_used": 80.0,
  "rpe": 7.5,
  "notes": "Voelde goed"
}
```

**Response:** `201 Created`
```json
{
  "log_id": 789,
  "session_id": 456,
  "exercise_id": 1,
  "set_number": 1,
  "reps_completed": 10,
  "weight_used": 80.0,
  "rpe": 7.5,
  "notes": "Voelde goed",
  "logged_at": "2025-10-02T10:15:00Z"
}
```

**Frontend Implementatie:**
```typescript
interface ExercisePerformance {
  set_number: number;
  reps_completed: number;
  weight_used?: number;
  rpe?: number;
  notes?: string;
}

async function logExercisePerformance(
  sessionId: number,
  exerciseId: number,
  performance: ExercisePerformance
) {
  const response = await fetch(
    `/api/v1/sessions/${sessionId}/exercises/${exerciseId}/log`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify(performance)
    }
  );
  
  if (!response.ok) throw new Error('Performance loggen mislukt');
  return await response.json();
}
```

---

### 5. **Get Active Session**

**Doel:** Huidige actieve sessie ophalen.  
**Use case:** Resume training, check status.

```http
GET /sessions/users/{userID}/active
```

**Response:** `200 OK`
```json
{
  "session_id": 456,
  "user_id": 123,
  "workout_id": 1,
  "started_at": "2025-10-02T10:00:00Z",
  "status": "in_progress",
  "workout": {
    "workout_id": 1,
    "focus": "upper_body",
    "exercises": [...]
  }
}
```

**Response bij geen actieve sessie:** `404 Not Found`

---

### 6. **Get Session History**

**Doel:** Historische workout sessies ophalen.  
**Use case:** Progress tracking, workout geschiedenis.

```http
GET /sessions/users/{userID}/history?page=1&page_size=20
```

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

**Response:** `200 OK`
```json
{
  "data": [
    {
      "session_id": 456,
      "user_id": 123,
      "workout_id": 1,
      "started_at": "2025-10-02T10:00:00Z",
      "completed_at": "2025-10-02T11:30:00Z",
      "duration_minutes": 90,
      "status": "completed"
    },
    ...
  ],
  "total_count": 45,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

---

### 7. **Get Session Metrics**

**Doel:** Gedetailleerde metrics van sessie ophalen.  
**Use case:** Sessie analyse, statistieken.

```http
GET /sessions/{sessionID}/metrics
```

**Response:** `200 OK`
```json
{
  "session_id": 456,
  "total_volume": 5400.0,
  "total_sets": 18,
  "total_reps": 135,
  "average_rpe": 7.8,
  "duration_minutes": 90,
  "exercises_completed": 6,
  "calories_burned": 450,
  "personal_records": [
    {
      "exercise_id": 1,
      "exercise_name": "Barbell Bench Press",
      "previous_best": 75.0,
      "new_best": 80.0,
      "improvement": 5.0
    }
  ]
}
```

---

### 8. **Get Weekly Session Stats**

**Doel:** Wekelijkse training statistieken ophalen.  
**Use case:** Weekly summary, progress dashboard.

```http
GET /sessions/users/{userID}/weekly?week_start=2025-09-30
```

**Query Parameters:**
- `week_start` (optional) - Start datum van week (YYYY-MM-DD), default: huidige week

**Response:** `200 OK`
```json
{
  "user_id": 123,
  "week_start": "2025-09-30",
  "week_end": "2025-10-06",
  "total_sessions": 4,
  "total_duration_minutes": 360,
  "total_volume": 21600.0,
  "sessions_by_day": {
    "monday": 1,
    "wednesday": 1,
    "thursday": 1,
    "saturday": 1
  },
  "completion_rate": 0.8,
  "average_rpe": 7.5
}
```

**Frontend Implementatie:**
```typescript
async function getWeeklyStats(userId: number, weekStart?: string) {
  const params = weekStart ? `?week_start=${weekStart}` : '';
  const response = await fetch(
    `/api/v1/sessions/users/${userId}/weekly${params}`,
    { headers: { 'Authorization': `Bearer ${getToken()}` } }
  );
  
  if (!response.ok) throw new Error('Statistieken ophalen mislukt');
  return await response.json();
}
```

---

## 🤖 Plan Generation

AI-powered workout plan generatie en aanpassingen.

### 1. **Create Plan Generation**

**Doel:** Nieuwe gepersonaliseerd trainingsplan genereren.  
**Use case:** Onboarding, nieuwe plan aanvragen.

```http
POST /plans/generate
```

**Request Body:**
```json
{
  "user_id": 123,
  "plan_duration_weeks": 8,
  "preferences": {
    "focus_areas": ["strength", "muscle_gain"],
    "avoid_exercises": [15, 23],
    "preferred_days": [1, 3, 5]
  }
}
```

**Response:** `201 Created`
```json
{
  "plan_id": 789,
  "user_id": 123,
  "generated_at": "2025-10-02T10:00:00Z",
  "duration_weeks": 8,
  "status": "active",
  "weekly_workouts": 3,
  "workouts": [
    {
      "workout_id": 1,
      "day_of_week": 1,
      "focus": "upper_body",
      "exercises": [...]
    },
    ...
  ]
}
```

**Frontend Implementatie:**
```typescript
interface PlanGenerationRequest {
  user_id: number;
  plan_duration_weeks: number;
  preferences?: {
    focus_areas?: string[];
    avoid_exercises?: number[];
    preferred_days?: number[];
  };
}

async function generatePlan(request: PlanGenerationRequest) {
  const response = await fetch('/api/v1/plans/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify(request)
  });
  
  if (!response.ok) throw new Error('Plan genereren mislukt');
  return await response.json();
}
```

---

### 2. **Get Active Plan**

**Doel:** Huidige actieve trainingsplan ophalen.  
**Use case:** Dashboard, workout planning.

```http
GET /plans/users/{userID}/active
```

**Response:** `200 OK`
```json
{
  "plan_id": 789,
  "user_id": 123,
  "generated_at": "2025-10-02T10:00:00Z",
  "duration_weeks": 8,
  "status": "active",
  "current_week": 2,
  "completion_rate": 0.75,
  "workouts": [...]
}
```

**Response bij geen actief plan:** `404 Not Found`

---

### 3. **Get Plan History**

**Doel:** Historische plannen ophalen.  
**Use case:** Vorige plannen bekijken.

```http
GET /plans/users/{userID}/history?limit=10
```

**Query Parameters:**
- `limit` (optional, default: 10) - Aantal plannen

**Response:** `200 OK`
```json
[
  {
    "plan_id": 789,
    "generated_at": "2025-10-02T10:00:00Z",
    "duration_weeks": 8,
    "status": "active",
    "completion_rate": 0.75
  },
  {
    "plan_id": 700,
    "generated_at": "2025-08-01T10:00:00Z",
    "duration_weeks": 12,
    "status": "completed",
    "completion_rate": 0.92
  },
  ...
]
```

---

### 4. **Track Plan Performance**

**Doel:** Performance van plan bijhouden voor aanpassingen.  
**Use case:** Auto-adjustment van plan.

```http
POST /plans/{planID}/track
```

**Request Body:**
```json
{
  "week_number": 2,
  "completion_rate": 0.75,
  "average_rpe": 7.8,
  "user_feedback": "Gaat goed maar upper body is zwaar"
}
```

**Response:** `200 OK`
```json
{
  "tracked": true,
  "adjustments_recommended": true,
  "suggestions": [
    "Overweeg upper body intensiteit te verlagen",
    "Extra rustdag tussen upper body workouts"
  ]
}
```

---

### 5. **Get Plan Effectiveness**

**Doel:** Effectiviteit score van plan ophalen.  
**Use case:** Plan evaluatie, AI learning.

```http
GET /plans/{planID}/effectiveness
```

**Response:** `200 OK`
```json
{
  "plan_id": 789,
  "effectiveness_score": 0.85,
  "adherence_rate": 0.9,
  "progress_rate": 0.82,
  "user_satisfaction": 4.5,
  "recommendations": [
    "Plan werkt goed, houd focus",
    "Overweeg intensiteit te verhogen in week 5"
  ]
}
```

---

### 6. **Mark Plan for Regeneration**

**Doel:** Plan markeren voor nieuwe generatie.  
**Use case:** Plan niet passend, nieuwe gewenst.

```http
POST /plans/{planID}/regenerate
```

**Request Body:**
```json
{
  "reason": "too_difficult",
  "feedback": "Te zwaar, kan niet alle oefeningen voltooien"
}
```

**Response:** `200 OK`
```json
{
  "marked_for_regeneration": true,
  "reason": "too_difficult",
  "feedback": "Te zwaar, kan niet alle oefeningen voltooien",
  "new_plan_available_at": "2025-10-03T10:00:00Z"
}
```

---

### 7. **Get Adaptation History**

**Doel:** Historische aanpassingen van plannen ophalen.  
**Use case:** AI learning, pattern detection.

```http
GET /plans/users/{userID}/adaptations
```

**Response:** `200 OK`
```json
[
  {
    "adaptation_id": 1,
    "plan_id": 789,
    "adapted_at": "2025-10-01T10:00:00Z",
    "reason": "performance_improvement",
    "changes": [
      "Increased weight for bench press by 5kg",
      "Added extra set to squats"
    ]
  },
  ...
]
```

---

## 📊 Types & Enums

### Fitness Levels
```typescript
type FitnessLevel = 'beginner' | 'intermediate' | 'advanced';
```

### Fitness Goals
```typescript
type FitnessGoal = 
  | 'strength'
  | 'muscle_gain'
  | 'fat_loss'
  | 'endurance'
  | 'general_fitness';
```

### Exercise Types
```typescript
type ExerciseType = 
  | 'strength'
  | 'cardio'
  | 'mobility'
  | 'hiit'
  | 'stretching';
```

### Equipment Types
```typescript
type EquipmentType = 
  | 'barbell'
  | 'dumbbell'
  | 'bodyweight'
  | 'machine'
  | 'kettlebell'
  | 'resistance_band';
```

### Session Status
```typescript
type SessionStatus = 
  | 'in_progress'
  | 'completed'
  | 'skipped'
  | 'cancelled';
```

### Plan Status
```typescript
type PlanStatus = 
  | 'active'
  | 'completed'
  | 'abandoned'
  | 'regenerating';
```

---

## ⚠️ Error Handling

### Standard Error Response

Alle errors volgen dit format:

```json
{
  "success": false,
  "error": "Error message hier",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional context"
  }
}
```

### HTTP Status Codes

| Code | Betekenis | Use Case |
|------|-----------|----------|
| 200 | OK | Succesvolle request |
| 201 | Created | Resource succesvol aangemaakt |
| 400 | Bad Request | Invalide input data |
| 401 | Unauthorized | Geen of invalide JWT token |
| 403 | Forbidden | Geen toegang tot resource |
| 404 | Not Found | Resource niet gevonden |
| 409 | Conflict | Resource conflict (bijv. duplicate) |
| 500 | Internal Server Error | Server fout |

### Frontend Error Handling

```typescript
async function apiRequest<T>(
  url: string,
  options?: RequestInit
): Promise<T> {
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`,
        ...options?.headers
      }
    });

    if (!response.ok) {
      const error = await response.json();
      throw new APIError(
        error.error || 'Request failed',
        response.status,
        error.code
      );
    }

    return await response.json();
  } catch (error) {
    if (error instanceof APIError) {
      // Handle API errors
      handleAPIError(error);
    } else {
      // Handle network errors
      handleNetworkError(error);
    }
    throw error;
  }
}

class APIError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public errorCode?: string
  ) {
    super(message);
    this.name = 'APIError';
  }
}
```

---

## 🔄 Rate Limiting

- **Default:** 100 requests per 15 minuten per IP
- **Authenticated:** 1000 requests per 15 minuten per gebruiker

Rate limit headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1696249200
```

---

## 🌐 CORS

De API accepteert requests van:
- `http://localhost:3000` (development)
- `http://localhost:19006` (Expo development)
- Productie domein (nog te bepalen)

---

## 📱 WebSocket Support (Toekomstig)

Voor real-time updates tijdens workouts:

```typescript
const ws = new WebSocket('ws://localhost:8080/ws/sessions/{sessionID}');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  // Handle real-time updates
};
```

---

## 🔧 Development Tips

### 1. TypeScript Types Genereren

Gebruik deze types in je frontend:

```typescript
// src/types/api.ts
export interface Exercise {
  exercise_id: number;
  name: string;
  muscle_groups: string[];
  difficulty: FitnessLevel;
  equipment: EquipmentType;
  type: ExerciseType;
  default_sets: number;
  default_reps: string;
  rest_seconds: number;
}

export interface WorkoutSession {
  session_id: number;
  user_id: number;
  workout_id: number;
  started_at: string;
  completed_at?: string;
  duration_minutes?: number;
  status: SessionStatus;
  notes?: string;
}

// etc...
```

### 2. API Client Setup

```typescript
// src/lib/api-client.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options?.headers
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return await response.json();
  }

  // Exercise endpoints
  async getExercise(id: number) {
    return this.request<Exercise>(`/exercises/${id}`);
  }

  async listExercises(page = 1, pageSize = 20) {
    return this.request<PaginatedResponse<Exercise>>(
      `/exercises?page=${page}&page_size=${pageSize}`
    );
  }

  // Workout session endpoints
  async startSession(userId: number, workoutId: number) {
    return this.request<WorkoutSession>('/sessions/start', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        workout_id: workoutId,
        started_at: new Date().toISOString()
      })
    });
  }

  // etc...
}

export const apiClient = new ApiClient();
```

### 3. React Query Integration

```typescript
// src/hooks/useExercises.ts
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';

export function useExercises(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['exercises', page, pageSize],
    queryFn: () => apiClient.listExercises(page, pageSize)
  });
}

export function useExercise(id: number) {
  return useQuery({
    queryKey: ['exercise', id],
    queryFn: () => apiClient.getExercise(id),
    enabled: !!id
  });
}
```

---

## 📚 Veelgestelde Vragen

### Q: Hoe start ik een workout sessie?
A: 1. Haal actieve plan op (`GET /plans/users/{userID}/active`)  
   2. Kies workout voor vandaag  
   3. Start sessie (`POST /sessions/start`)  
   4. Log exercises tijdens training  
   5. Complete sessie (`POST /sessions/{sessionID}/complete`)

### Q: Hoe filter ik oefeningen op spiergroep?
A: Gebruik `GET /exercises/filter?muscle_group=chest` of `GET /exercises/muscle-groups/chest`

### Q: Hoe track ik progressie?
A: Log elke set met `POST /sessions/{sessionID}/exercises/{exerciseID}/log` en bekijk metrics met `GET /sessions/{sessionID}/metrics`

### Q: Hoe genereer ik een nieuw plan?
A: `POST /plans/generate` met user preferences

---

## 📞 Support

Voor vragen of issues:
- **GitHub:** [github.com/tdmdh/fit-up](https://github.com/tdmdh/fit-up)
- **Email:** support@fit-up.app

---

**Document Versie:** 1.0  
**Laatste Update:** Oktober 2, 2025  
**Auteur:** Fit-Up Backend Team
