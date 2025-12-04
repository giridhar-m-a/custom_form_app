-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-----------------------------
-- Table: users
-----------------------------
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_full_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) UNIQUE NOT NULL,
    user_google_id VARCHAR(255) UNIQUE,
    user_profile_pic_id UUID,
    user_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_password VARCHAR(255)
);

-- Trigger to update user_updated_at
CREATE OR REPLACE FUNCTION update_user_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.user_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_user_updated_at();

-----------------------------
-- Table: user_images
-----------------------------
CREATE TABLE IF NOT EXISTS user_images (
    file_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(255) NOT NULL,
    file_uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_id UUID UNIQUE REFERENCES users(user_id)
);

ALTER TABLE users
ADD CONSTRAINT fk_user_profile_pic FOREIGN KEY (user_profile_pic_id) REFERENCES user_images(file_id);

ALTER TABLE users
ADD CONSTRAINT uq_user_profile_pic UNIQUE (user_profile_pic_id);

-----------------------------
-- Enums
-----------------------------
CREATE TYPE form_status AS ENUM ('draft', 'published', 'archived', 'closed');
CREATE TYPE form_access AS ENUM ('public', 'restricted');
CREATE TYPE invitation_status AS ENUM ('opened', 'submitted', 'invited', 'expired');

-----------------------------
-- Table: forms
-----------------------------
CREATE TABLE IF NOT EXISTS forms (
    form_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_title VARCHAR(255) NOT NULL,
    form_description TEXT,
    form_status form_status DEFAULT 'draft',
    form_access form_access DEFAULT 'restricted',
    form_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    form_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(user_id) ON DELETE CASCADE
);

-- Trigger to update form_updated_at
CREATE OR REPLACE FUNCTION update_form_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.form_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_form_updated_at
BEFORE UPDATE ON forms
FOR EACH ROW
EXECUTE FUNCTION update_form_updated_at();

-- Indexes
CREATE INDEX IF NOT EXISTS idx_forms_created_by ON forms(created_by);
CREATE INDEX IF NOT EXISTS idx_forms_status ON forms(form_status);
CREATE INDEX IF NOT EXISTS idx_forms_access ON forms(form_access);
CREATE INDEX IF NOT EXISTS idx_forms_title_description_gin
ON forms USING gin (to_tsvector('english', form_title || ' ' || form_description));

-----------------------------
-- Table: form_fields
-----------------------------
CREATE TABLE IF NOT EXISTS form_fields (
    field_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    field_label TEXT NOT NULL,
    field_type VARCHAR(255) NOT NULL,
    form_id UUID NOT NULL REFERENCES forms(form_id) ON DELETE CASCADE,
    is_required BOOLEAN DEFAULT FALSE,
    ordering INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_form_fields_form_id ON form_fields(form_id);

-----------------------------
-- Table: form_field_options
-----------------------------
CREATE TABLE IF NOT EXISTS form_field_options (
    option_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    option_label TEXT NOT NULL,
    field_id UUID REFERENCES form_fields(field_id) ON DELETE CASCADE,
    ordering INT NOT NULL,
    is_answer BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_form_field_options_field_id ON form_field_options(field_id);

-----------------------------
-- Table: invitations
-----------------------------
CREATE TABLE IF NOT EXISTS invitations (
    invitation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES forms(form_id) ON DELETE CASCADE,
    invited_email VARCHAR(255) NOT NULL,
    invitation_token UUID DEFAULT gen_random_uuid(),
    invited_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    invited_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    status invitation_status DEFAULT 'invited',
    opened_at TIMESTAMP WITH TIME ZONE,
    submitted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_invitations_form_id ON invitations(form_id);
CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(invited_email);
CREATE INDEX IF NOT EXISTS idx_invitations_invited_by ON invitations(invited_by);
CREATE INDEX IF NOT EXISTS idx_invitations_status ON invitations(status);

-----------------------------
-- Table: form_responses
-----------------------------
CREATE TABLE IF NOT EXISTS form_responses (
    response_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES forms(form_id) ON DELETE CASCADE,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    respondent_id UUID REFERENCES invitations(invitation_id) ON DELETE SET NULL,
    form_option_id UUID REFERENCES form_field_options(option_id) ON DELETE SET NULL,
    form_field_id UUID NOT NULL REFERENCES form_fields(field_id) ON DELETE CASCADE,
    response_text TEXT
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_form_responses_form_id ON form_responses(form_id);
CREATE INDEX IF NOT EXISTS idx_form_responses_respondent_id ON form_responses(respondent_id);
CREATE INDEX IF NOT EXISTS idx_form_responses_field_id ON form_responses(form_field_id);
CREATE INDEX IF NOT EXISTS idx_form_responses_option_id ON form_responses(form_option_id);

-----------------------------
-- Table: form_response_files
-----------------------------
CREATE TABLE IF NOT EXISTS form_response_files (
    response_file_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(1024) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(255) NOT NULL,
    file_uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    response_id UUID NOT NULL REFERENCES form_responses(response_id) ON DELETE CASCADE,
    form_id UUID NOT NULL REFERENCES forms(form_id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_form_response_files_response_id ON form_response_files(response_id);
CREATE INDEX IF NOT EXISTS idx_form_response_files_form_id ON form_response_files(form_id);
