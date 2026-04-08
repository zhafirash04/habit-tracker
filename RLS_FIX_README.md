# RLS Fix Summary - HabitFlow

## Problem Diagnosis

✅ **Supabase Key Type**: **Direct PostgreSQL Connection** (not using Supabase client SDK)
- Connection method: GORM with `DATABASE_URL` environment variable
- Format: `postgresql://USER:PASSWORD@HOST:PORT/DATABASE`
- Location: `internal/config/config.go:34` and `.env` file

✅ **Tables Affected**: 6 tables (found in `internal/database/database.go:33-39`)
1. `users`
2. `refresh_tokens`
3. `habits`
4. `habit_logs`
5. `streaks`
6. `push_subscriptions`

## Root Cause

When RLS (Row-Level Security) was enabled on all tables without policies:
- **All database queries are blocked** (default deny)
- Even though the app uses a direct PostgreSQL connection, the connecting user lacks `BYPASSRLS` privilege
- GORM queries fail because RLS is blocking read/write operations

## Solution

I've created **`fix_rls_policies.sql`** with two approaches:

### Approach 1: Grant BYPASSRLS (Recommended ⭐)
Best for backend services that need full database access.

```sql
-- Determine your user first:
SELECT current_user;

-- Then grant BYPASSRLS (replace 'postgres' with your user):
ALTER ROLE postgres BYPASSRLS;
```

**When to use**:
- ✅ Your DATABASE_URL uses `postgres` user (pooler connection)
- ✅ You have permission to grant BYPASSRLS
- ✅ Backend is trusted and handles all authorization

### Approach 2: Create Permissive Policies (Fallback)
Creates `FOR ALL` policies allowing all operations.

**When to use**:
- ❌ Cannot grant BYPASSRLS to your user
- ❌ Using `authenticated` or `anon` role (uncommon for backend)
- ✅ Want RLS enabled as safety net with permissive policies

## Quick Start

### Step 1: Check your current database user
```sql
SELECT current_user, rolbypassrls FROM pg_roles WHERE rolname = current_user;
```

### Step 2: Apply the fix

**Option A - Via Supabase Dashboard** (Easiest):
1. Go to **Supabase Dashboard** → **SQL Editor**
2. Open `fix_rls_policies.sql`
3. If your user is `postgres`: Uncomment line 30 and run just that line
4. If unsure: Run the entire file (safe, idempotent)

**Option B - Via psql**:
```bash
psql "$DATABASE_URL" -f fix_rls_policies.sql
```

### Step 3: Verify it worked
```sql
-- Should return your username and 't' for bypassrls:
SELECT current_user, rolbypassrls
FROM pg_roles
WHERE rolname = current_user;

-- Should work without errors:
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM habits;
```

### Step 4: Test your application
```bash
# Restart your app and verify it works:
go run cmd/server/main.go

# Test API endpoints:
curl http://localhost:8080/api/v1/health
```

## Files Created

1. **`fix_rls_policies.sql`** - Ready-to-run SQL migration
   - Comprehensive comments explaining each approach
   - Idempotent (safe to run multiple times)
   - Includes verification queries

## Architecture Notes

HabitFlow's security model:
- ✅ **Backend handles ALL authorization** via JWT middleware
- ✅ **No direct database access from frontend**
- ✅ **User context validated in every protected endpoint**
- ✅ **RLS serves as defense-in-depth, not primary security**

This is why permissive RLS policies (or BYPASSRLS) are appropriate for this architecture.

## Why This Happened

Supabase RLS is designed for **Supabase client SDK** usage where:
- Frontend uses `anon` key to connect
- Supabase Auth provides `auth.uid()` for RLS policies
- Policies enforce user-specific access

But HabitFlow uses a different architecture:
- Backend uses **direct PostgreSQL connection** (GORM)
- Custom JWT auth (not Supabase Auth)
- Application-layer authorization

When RLS was enabled without considering this architecture, it blocked the backend's database access.

## Next Steps

1. ✅ Run the SQL migration
2. ✅ Verify database access is restored
3. ✅ Test the application end-to-end
4. ✅ Consider documenting this in your deployment runbook

## Need Different Policies?

If you want truly restrictive RLS policies (users can only access their own data):

**Requirements**:
- Switch from custom JWT to Supabase Auth
- Modify all queries to use Supabase client
- Update policies to use `auth.uid()` in conditions
- Major architecture change

**Current approach is recommended** for your architecture.
