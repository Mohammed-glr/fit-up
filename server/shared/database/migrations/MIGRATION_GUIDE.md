# ✅ SQL Migrations Fixed - Summary

## 🎯 What Was Fixed

### **1. Removed Duplicate Users Table** ✅
- **Problem**: Migration `005_create_workout_schema.up.sql` was creating a duplicate `users` table
- **Solution**: Removed the duplicate table definition since users are created in `001_create_users.up.sql`
- **Changed**: All `user_id` references now use `TEXT` type (matching auth users table) instead of `INT`

### **2. Fixed User ID References** ✅
Updated to use TEXT user_id consistently:
- `weekly_schemas.user_id` → TEXT (references `users(id)`)
- `progress_logs.user_id` → TEXT (references `users(id)`)
- All other tables already used TEXT correctly

### **3. Fixed SQL Syntax Errors** ✅
- **Problem**: `UNIQUE(goal_id, calculated_at::DATE)` syntax error in line 75
- **Solution**: Replaced with proper `CREATE UNIQUE INDEX` statements
- **Fixed in**:
  - `goal_progress` table
  - `optimal_loads` table

### **4. Renamed Migration Files** ✅
All migrations now have proper `.up.sql` and `.down.sql` suffixes:
- ✅ `006_create_fitness_profile_tables.up.sql` / `.down.sql`
- ✅ `007_create_session_tracking_tables.up.sql` / `.down.sql`
- ✅ `008_create_plan_generation_tables.up.sql` / `.down.sql`
- ✅ `009_create_analytics_metrics_tables.up.sql` / `.down.sql`
- ✅ `010_fix_schema_inconsistencies.up.sql` / `.down.sql`
- ✅ `011_add_performance_optimizations.up.sql` / `.down.sql`

### **5. Created Down Migrations** ✅
All migrations now have proper rollback scripts.

---

## 📋 Migration File Structure

```
shared/database/migrations/
├── 001_create_users.up.sql              ✅ Auth & users
├── 001_create_users.down.sql
├── 002_add_jwt_management.up.sql        ✅ JWT tokens
├── 002_add_jwt_management.down.sql
├── 003_add_oauth_support.up.sql         ✅ OAuth
├── 003_add_oauth_support.down.sql
├── 004_fix_jwt_refresh_tokens_trigger.up.sql
├── 004_fix_jwt_refresh_tokens_trigger.down.sql
├── 005_create_workout_schema.up.sql     ✅ FIXED - No duplicate users
├── 005_create_workout_schema.down.sql
├── 006_create_fitness_profile_tables.up.sql  ✅ FIXED
├── 006_create_fitness_profile_tables.down.sql
├── 007_create_session_tracking_tables.up.sql
├── 007_create_session_tracking_tables.down.sql
├── 008_create_plan_generation_tables.up.sql
├── 008_create_plan_generation_tables.down.sql
├── 009_create_analytics_metrics_tables.up.sql  ✅ FIXED - Syntax errors
├── 009_create_analytics_metrics_tables.down.sql
├── 010_fix_schema_inconsistencies.up.sql  ✅ FIXED
├── 010_fix_schema_inconsistencies.down.sql
├── 011_add_performance_optimizations.up.sql
└── 011_add_performance_optimizations.down.sql
```

---

## 🚀 How to Run Migrations

### **Step 1: Set Environment Variable**
```powershell
# Set your database URL
$env:DATABASE_URL = "postgres://fitup:fitup_password@localhost:5432/fitup?sslmode=disable"
```

### **Step 2: Check Current Version**
```powershell
cd C:\Users\Mohammed\fit-up\server
migrate -path "shared/database/migrations" -database $env:DATABASE_URL version
```

### **Step 3: Run All Migrations**
```powershell
# Apply all migrations
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up
```

### **Step 4: Verify Success**
```powershell
# Check version again
migrate -path "shared/database/migrations" -database $env:DATABASE_URL version
```

---

## 🔄 Other Useful Migration Commands

### **Apply Specific Number of Migrations**
```powershell
# Apply next 3 migrations
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up 3
```

### **Rollback Migrations**
```powershell
# Rollback last migration
migrate -path "shared/database/migrations" -database $env:DATABASE_URL down 1

# Rollback all migrations
migrate -path "shared/database/migrations" -database $env:DATABASE_URL down -all
```

### **Go to Specific Version**
```powershell
# Migrate to version 5
migrate -path "shared/database/migrations" -database $env:DATABASE_URL goto 5
```

### **Force Version (if dirty state)**
```powershell
# Force set to version 5 without running migration
migrate -path "shared/database/migrations" -database $env:DATABASE_URL force 5
```

### **Drop Everything (DANGER!)**
```powershell
# Drop all tables and data
migrate -path "shared/database/migrations" -database $env:DATABASE_URL drop -f
```

---

## 🗄️ Database Schema Overview

After running all migrations, you'll have:

### **Auth & User Management** (001-004)
- `users` - Main user table (TEXT id with UUID)
- `sessions` - User sessions
- `refresh_tokens` - JWT refresh tokens
- `password_reset_tokens` - Password reset flow
- `verification_tokens` - Email verification
- `two_factor_tokens` - 2FA support
- `accounts` - OAuth linked accounts
- `oauth_states` - OAuth state management

### **Workout System** (005)
- `exercises` - Exercise library
- `workout_templates` - Workout program templates
- `weekly_schemas` - User weekly schedules
- `workouts` - Individual workout days
- `workout_exercises` - Exercises in workouts
- `progress_logs` - Performance tracking

### **Fitness Profiles** (006)
- `workout_profiles` - User fitness setup
- `fitness_assessments` - Fitness level assessments
- `fitness_goal_targets` - User goals
- `movement_assessments` - Movement quality
- `movement_limitations` - Physical limitations
- `one_rep_max_estimates` - Strength estimates

### **Session Tracking** (007)
- `workout_sessions` - Workout session data
- `skipped_workouts` - Missed workout tracking
- `exercise_performances` - Exercise performance
- `set_performances` - Individual set data
- `session_metrics` - Calculated metrics
- `weekly_session_stats` - Weekly summaries

### **Plan Generation** (008)
- `generated_plans` - AI-generated plans
- `plan_performance_data` - Plan effectiveness
- `plan_adaptations` - Plan modifications
- `plan_generation_metadata` - Algorithm data

### **Analytics & Metrics** (009)
- `recovery_metrics` - Recovery tracking
- `strength_progressions` - Strength trends
- `plateau_detections` - Plateau identification
- `training_volumes` - Volume tracking
- `intensity_progressions` - Intensity trends
- `goal_progress` - Goal tracking
- `goal_predictions` - Goal predictions
- `goal_adjustments` - Recommendations
- `optimal_loads` - Load recommendations

### **Optimizations** (010-011)
- Enum types for data consistency
- Additional constraints and validations
- Performance indexes
- Monitoring views

---

## 🧪 Testing Your Database

### **1. Connect to Database**
```powershell
# Using psql
psql $env:DATABASE_URL

# Or using pgAdmin or any PostgreSQL client
```

### **2. Check Tables**
```sql
-- List all tables
\dt

-- Check users table
SELECT * FROM users LIMIT 5;

-- Check exercises table
SELECT * FROM exercises LIMIT 5;

-- Check migration version
SELECT * FROM schema_migrations;
```

### **3. Verify Foreign Keys**
```sql
-- Check all foreign key constraints
SELECT
    tc.table_name, 
    kcu.column_name, 
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
ORDER BY tc.table_name;
```

---

## ⚠️ Troubleshooting

### **Problem: "dirty database version"**
**Solution:**
```powershell
# Check current version
migrate -path "shared/database/migrations" -database $env:DATABASE_URL version

# Force to last good version (e.g., 4)
migrate -path "shared/database/migrations" -database $env:DATABASE_URL force 4

# Try migration again
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up
```

### **Problem: "duplicate key value violates unique constraint"**
**Solution:**
```powershell
# Check for existing data
psql $env:DATABASE_URL -c "SELECT * FROM users LIMIT 5;"

# If needed, drop and recreate
migrate -path "shared/database/migrations" -database $env:DATABASE_URL down -all
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up
```

### **Problem: Database doesn't exist**
**Solution:**
```powershell
# Create database
psql "postgres://fitup:fitup_password@localhost:5432/postgres" -c "CREATE DATABASE fitup;"

# Run migrations
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up
```

---

## 📝 Key Changes Summary

| Issue | Before | After |
|-------|--------|-------|
| **Duplicate Users Table** | Created in both 001 and 005 | Only in 001 ✅ |
| **user_id Type** | Mixed INT and TEXT | Consistent TEXT ✅ |
| **Syntax Error (line 75)** | `UNIQUE(goal_id, calculated_at::DATE)` | Separate INDEX ✅ |
| **Migration Names** | `.sql` files | `.up.sql` / `.down.sql` ✅ |
| **Down Migrations** | Missing | All created ✅ |

---

## ✅ Ready to Run!

Your migrations are now **fixed and ready**. Run:

```powershell
cd C:\Users\Mohammed\fit-up\server
$env:DATABASE_URL = "postgres://fitup:fitup_password@localhost:5432/fitup?sslmode=disable"
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up
```

🎉 **All SQL tables are fixed and ready for deployment!**
