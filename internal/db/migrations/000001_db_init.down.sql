-- 1. Drop the trigger from users table
DROP TRIGGER IF EXISTS set_user_updated_at ON users;

-- 2. Drop the trigger function
DROP FUNCTION IF EXISTS update_user_updated_at;

-- 3. Drop tables (child first to avoid FK constraints)
DROP TABLE IF EXISTS user_images;
DROP TABLE IF EXISTS users;

-- 4. (Optional) Remove pgcrypto extension
DROP EXTENSION IF EXISTS pgcrypto;
