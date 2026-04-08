-- ============================================================================
-- HabitFlow RLS Policy Fix
-- ============================================================================
-- Created: 2026-04-08
-- Purpose: Fix database access after RLS was enabled on all tables
--
-- BACKGROUND:
-- This app uses direct PostgreSQL connection (GORM), not Supabase client SDK.
-- Connection string format: postgresql://USER:PASSWORD@HOST:PORT/DB
--
-- SOLUTION APPROACH:
-- 1. First, check if the connecting user has BYPASSRLS privilege
-- 2. If using 'postgres' superuser or service_role: grant BYPASSRLS
-- 3. If using a regular user: create permissive policies for backend operations
-- ============================================================================

-- ----------------------------------------------------------------------------
-- OPTION 1: Grant BYPASSRLS to the backend connection user
-- ----------------------------------------------------------------------------
-- This is the recommended approach for backend services that need full access.
-- Replace 'postgres' with your actual DATABASE_URL username if different.
--
-- To find your username, check your DATABASE_URL:
-- postgresql://[THIS_IS_USERNAME]:password@host:port/db
--
-- Common Supabase usernames:
-- - postgres (pooler connection)
-- - service_role (if using service_role key via API)
-- - authenticated (if using authenticated key via API)
-- - anon (if using anon key via API)

-- Uncomment ONE of the following based on your DATABASE_URL username:
-- ALTER ROLE postgres BYPASSRLS;
-- ALTER ROLE service_role BYPASSRLS;
-- ALTER ROLE authenticated BYPASSRLS;

-- If you're not sure which user, run this query first to see current user:
-- SELECT current_user, current_role;

-- ----------------------------------------------------------------------------
-- OPTION 2: Create permissive RLS policies (if BYPASSRLS is not an option)
-- ----------------------------------------------------------------------------
-- Use this if:
-- - You cannot grant BYPASSRLS to your user
-- - You want RLS enabled for security but need backend to work
-- - You're using 'authenticated' or 'anon' role

-- Enable RLS on all tables (already done, but ensuring it's set)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE refresh_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE habits ENABLE ROW LEVEL SECURITY;
ALTER TABLE habit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE streaks ENABLE ROW LEVEL SECURITY;
ALTER TABLE push_subscriptions ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- USERS TABLE POLICIES
-- ============================================================================

-- Drop existing policies if any (idempotent)
DROP POLICY IF EXISTS "users_backend_all" ON users;
DROP POLICY IF EXISTS "users_select_all" ON users;
DROP POLICY IF EXISTS "users_insert_all" ON users;
DROP POLICY IF EXISTS "users_update_all" ON users;
DROP POLICY IF EXISTS "users_delete_all" ON users;

-- Create permissive policy for backend operations
-- This allows the backend service to manage all users
CREATE POLICY "users_backend_all" ON users
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- REFRESH_TOKENS TABLE POLICIES
-- ============================================================================

DROP POLICY IF EXISTS "refresh_tokens_backend_all" ON refresh_tokens;
DROP POLICY IF EXISTS "refresh_tokens_select_all" ON refresh_tokens;
DROP POLICY IF EXISTS "refresh_tokens_insert_all" ON refresh_tokens;
DROP POLICY IF EXISTS "refresh_tokens_update_all" ON refresh_tokens;
DROP POLICY IF EXISTS "refresh_tokens_delete_all" ON refresh_tokens;

CREATE POLICY "refresh_tokens_backend_all" ON refresh_tokens
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- HABITS TABLE POLICIES
-- ============================================================================

DROP POLICY IF EXISTS "habits_backend_all" ON habits;
DROP POLICY IF EXISTS "habits_select_all" ON habits;
DROP POLICY IF EXISTS "habits_insert_all" ON habits;
DROP POLICY IF EXISTS "habits_update_all" ON habits;
DROP POLICY IF EXISTS "habits_delete_all" ON habits;

CREATE POLICY "habits_backend_all" ON habits
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- HABIT_LOGS TABLE POLICIES
-- ============================================================================

DROP POLICY IF EXISTS "habit_logs_backend_all" ON habit_logs;
DROP POLICY IF EXISTS "habit_logs_select_all" ON habit_logs;
DROP POLICY IF EXISTS "habit_logs_insert_all" ON habit_logs;
DROP POLICY IF EXISTS "habit_logs_update_all" ON habit_logs;
DROP POLICY IF EXISTS "habit_logs_delete_all" ON habit_logs;

CREATE POLICY "habit_logs_backend_all" ON habit_logs
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- STREAKS TABLE POLICIES
-- ============================================================================

DROP POLICY IF EXISTS "streaks_backend_all" ON streaks;
DROP POLICY IF EXISTS "streaks_select_all" ON streaks;
DROP POLICY IF EXISTS "streaks_insert_all" ON streaks;
DROP POLICY IF EXISTS "streaks_update_all" ON streaks;
DROP POLICY IF EXISTS "streaks_delete_all" ON streaks;

CREATE POLICY "streaks_backend_all" ON streaks
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- PUSH_SUBSCRIPTIONS TABLE POLICIES
-- ============================================================================

DROP POLICY IF EXISTS "push_subscriptions_backend_all" ON push_subscriptions;
DROP POLICY IF EXISTS "push_subscriptions_select_all" ON push_subscriptions;
DROP POLICY IF EXISTS "push_subscriptions_insert_all" ON push_subscriptions;
DROP POLICY IF EXISTS "push_subscriptions_update_all" ON push_subscriptions;
DROP POLICY IF EXISTS "push_subscriptions_delete_all" ON push_subscriptions;

CREATE POLICY "push_subscriptions_backend_all" ON push_subscriptions
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================
-- After running this migration, execute these queries to verify:

-- 1. Check current user and RLS bypass status:
-- SELECT current_user,
--        rolbypassrls
-- FROM pg_roles
-- WHERE rolname = current_user;

-- 2. Verify RLS is enabled on all tables:
-- SELECT schemaname, tablename, rowsecurity
-- FROM pg_tables
-- WHERE schemaname = 'public'
-- AND tablename IN ('users', 'refresh_tokens', 'habits', 'habit_logs', 'streaks', 'push_subscriptions');

-- 3. List all policies:
-- SELECT schemaname, tablename, policyname, permissive, roles, cmd, qual, with_check
-- FROM pg_policies
-- WHERE schemaname = 'public'
-- ORDER BY tablename, policyname;

-- 4. Test a simple query to verify access:
-- SELECT COUNT(*) FROM users;
-- SELECT COUNT(*) FROM habits;

-- ============================================================================
-- USAGE INSTRUCTIONS
-- ============================================================================
--
-- STEP 1: Determine your connection user
-- Run this in Supabase SQL Editor or via psql:
--   SELECT current_user;
--
-- STEP 2: Choose the appropriate fix:
--
-- SCENARIO A: Using 'postgres' user (pooler connection)
--   → Uncomment line 30: ALTER ROLE postgres BYPASSRLS;
--   → Run only that line, skip the policies section
--   → This is the cleanest solution for backend services
--
-- SCENARIO B: Using 'service_role'
--   → Uncomment line 31: ALTER ROLE service_role BYPASSRLS;
--   → Run only that line, skip the policies section
--
-- SCENARIO C: Using 'authenticated' or 'anon' role (uncommon for backend)
--   → Run the entire policies section (lines 48-161)
--   → This creates permissive policies allowing all operations
--
-- SCENARIO D: Not sure which user
--   → Run the entire file (it's safe, idempotent)
--   → BYPASSRLS grants will fail if user doesn't exist (safe to ignore)
--   → Policies will be created as fallback
--
-- STEP 3: Run this file
-- Option A - Via Supabase Dashboard:
--   1. Go to SQL Editor in Supabase Dashboard
--   2. Paste and run the appropriate section
--
-- Option B - Via psql:
--   psql "postgresql://user:pass@host:port/db" -f fix_rls_policies.sql
--
-- STEP 4: Verify
--   Run the verification queries (lines 165-184)
--   Test your application to ensure it can read/write to the database
--
-- ============================================================================
-- SECURITY NOTES
-- ============================================================================
--
-- The policies created here are PERMISSIVE (allow all operations).
-- This is appropriate because:
-- 1. HabitFlow backend handles ALL authorization logic in application code
-- 2. The backend uses JWT middleware to authenticate users
-- 3. All API endpoints are protected with user context from JWT
-- 4. There is no direct database access from frontend
--
-- If you want more granular RLS policies (e.g., users can only see their own data),
-- you would need to:
-- 1. Set up Supabase Auth instead of custom JWT
-- 2. Use auth.uid() in policy conditions
-- 3. Modify the backend to pass user context to database queries
--
-- Current architecture: Backend is trusted, RLS serves as safety net only.
-- ============================================================================
