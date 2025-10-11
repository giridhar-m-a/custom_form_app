-- =========================
-- DOWN MIGRATION
-- =========================

-- 1. Drop triggers and functions for `forms`
DROP TRIGGER IF EXISTS set_form_updated_at ON forms;
DROP FUNCTION IF EXISTS update_form_updated_at;

-- 2. Drop triggers and functions for `users`
DROP TRIGGER IF EXISTS set_user_updated_at ON users;
DROP FUNCTION IF EXISTS update_user_updated_at;

-- 3. Drop foreign key + unique constraints related to user_images
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_user_profile_pic;
ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_user_profile_pic;

-- 4. Drop tables in reverse dependency order (child → parent)

DROP TABLE IF EXISTS form_responses CASCADE;
DROP TABLE IF EXISTS invitations CASCADE;
DROP TABLE IF EXISTS form_field_options CASCADE;
DROP TABLE IF EXISTS form_fields CASCADE;
DROP TABLE IF EXISTS forms CASCADE;
DROP TABLE IF EXISTS user_images CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- 5. Drop ENUM types
DROP TYPE IF EXISTS invitation_status;
DROP TYPE IF EXISTS form_status;
DROP TYPE IF EXISTS form_access;

-- 6. Drop extension (optional, if not used elsewhere)
-- DROP EXTENSION IF EXISTS pgcrypto;
