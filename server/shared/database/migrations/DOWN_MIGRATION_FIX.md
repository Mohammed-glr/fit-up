# 🔧 Quick Fix Applied - Down Migration Order

## ✅ Fixed: `001_create_users.down.sql`

### **Problem:**
The down migration was trying to drop tables in the wrong order:
```sql
DROP TABLE IF EXISTS subroles;  -- ❌ Can't drop this first!
DROP TABLE IF EXISTS users;     -- Has FK to subroles
```

### **Solution:**
Fixed the order to drop dependent tables first:
```sql
DROP TABLE IF EXISTS users CASCADE;     -- Drop first (has FK)
DROP TABLE IF EXISTS subroles CASCADE;  -- Drop second (referenced)
```

---

## 🚀 Commands to Run

### **Step 1: Force Clean State**
```powershell
cd C:\Users\Mohammed\fit-up\server

# Set database URL
$env:DATABASE_URL = "postgres://fitup:fitup_password@localhost:5432/fitup?sslmode=disable"

# Force to a clean version (since migration 1 failed)
migrate -path "shared/database/migrations" -database $env:DATABASE_URL force 1
```

### **Step 2: Try Down Migration Again**
```powershell
# Now try down migration again
migrate -path "shared/database/migrations" -database $env:DATABASE_URL down
```

### **Step 3: Or Start Fresh**
If you want to start completely fresh:

```powershell
# Drop everything and start over
migrate -path "shared/database/migrations" -database $env:DATABASE_URL drop -f

# Run all migrations from scratch
migrate -path "shared/database/migrations" -database $env:DATABASE_URL up
```

---

## 📋 What Was Changed

**File:** `001_create_users.down.sql`

**Changes:**
1. ✅ Added `CASCADE` to all DROP TABLE statements
2. ✅ Reordered table drops - dependent tables first
3. ✅ Added DROP INDEX statements
4. ✅ Fixed function name `trigger_set_timestamp()` → `trigger_set_timestamp()`

**New Order:**
1. Drop triggers and functions
2. Drop indexes
3. Drop tables with foreign keys (`login_activities`, `sessions`, etc.)
4. Drop `users` table (which has FK to subroles)
5. Drop `subroles` table
6. Drop enum types

---

## ✅ Ready to Test!

Run the commands above to test the fix. The migration should now work correctly! 🎉
