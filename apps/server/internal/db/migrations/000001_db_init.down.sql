-- -----------------------------
-- Down Migration: Drop all tables, types, and extension
-- -----------------------------

-- Drop indexes for clarity (optional, they are dropped automatically with tables)
DROP INDEX IF EXISTS idx_form_response_files_response_id;
DROP INDEX IF EXISTS idx_form_response_files_form_id;

DROP INDEX IF EXISTS idx_form_responses_form_id;
DROP INDEX IF EXISTS idx_form_responses_respondent_id;
DROP INDEX IF EXISTS idx_form_responses_field_id;
DROP INDEX IF EXISTS idx_form_responses_option_id;

DROP INDEX IF EXISTS idx_invitations_form_id;
DROP INDEX IF EXISTS idx_invitations_email;
DROP INDEX IF EXISTS idx_invitations_invited_by;
DROP INDEX IF EXISTS idx_invitations_status;

DROP INDEX IF EXISTS idx_form_field_options_field_id;
DROP INDEX IF EXISTS idx_form_fields_form_id;

DROP INDEX IF EXISTS idx_forms_created_by;
DROP INDEX IF EXISTS idx_forms_status;
DROP INDEX IF EXISTS idx_forms_access;
DROP INDEX IF EXISTS idx_forms_title_description_gin;

-- Drop triggers
DROP TRIGGER IF EXISTS set_user_updated_at ON users;
DROP FUNCTION IF EXISTS update_user_updated_at();

DROP TRIGGER IF EXISTS set_form_updated_at ON forms;
DROP FUNCTION IF EXISTS update_form_updated_at();

-- Drop tables in dependency order
DROP TABLE IF EXISTS form_response_files CASCADE;
DROP TABLE IF EXISTS form_responses CASCADE;
DROP TABLE IF EXISTS invitations CASCADE;
DROP TABLE IF EXISTS form_field_options CASCADE;
DROP TABLE IF EXISTS form_fields CASCADE;
DROP TABLE IF EXISTS forms CASCADE;
DROP TABLE IF EXISTS user_images CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop types
DROP TYPE IF EXISTS invitation_status;
DROP TYPE IF EXISTS form_status;
DROP TYPE IF EXISTS form_access;

-- Drop extension
DROP EXTENSION IF EXISTS pgcrypto;
