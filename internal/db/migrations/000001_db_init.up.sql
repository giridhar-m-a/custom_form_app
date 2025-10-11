-- Enable the pgcrypto extension to generate UUIDs
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_full_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) UNIQUE NOT NULL,
    user_google_id VARCHAR(255) UNIQUE NOT NULL,
    user_profile_pic_id UUID,
    user_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the function that updates user_updated_at
CREATE OR REPLACE FUNCTION update_user_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.user_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to the users table
CREATE TRIGGER set_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_user_updated_at();


CREATE TABLE IF NOT EXISTS user_images (
    file_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(255) NOT NULL,
    file_uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_id UUID UNIQUE REFERENCES  users(user_id)
);

ALTER TABLE users
ADD CONSTRAINT fk_user_profile_pic FOREIGN KEY (user_profile_pic_id) REFERENCES user_images(file_id);

ALTER TABLE users
ADD CONSTRAINT uq_user_profile_pic UNIQUE (user_profile_pic_id);

CREATE TYPE form_status AS ENUM ('draft', 'published', 'archived', 'closed');
CREATE TYPE form_access AS ENUM ('public', 'restricted');

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

CREATE OR REPLACE FUNCTION update_form_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.form_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to the forms table
CREATE TRIGGER set_form_updated_at
BEFORE UPDATE ON forms
FOR EACH ROW
EXECUTE FUNCTION update_form_updated_at();
    
CREATE TABLE IF NOT EXISTS form_fields (
    field_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    field_label TEXT NOT NULL,
    field_type VARCHAR(255) NOT NULL,
    form_id UUID REFERENCES forms(form_id) ON DELETE CASCADE NOT NULL,
    is_required BOOLEAN DEFAULT FALSE,
    ordering INT NOT NULL
);

CREATE TABLE IF NOT EXISTS form_field_options (
    option_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    option_label TEXT NOT NULL,
    field_id UUID REFERENCES form_fields(field_id) ON DELETE CASCADE,
    ordering INT NOT NULL
);

CREATE TYPE invitation_status AS ENUM ('opened', 'submitted', 'invited', 'expired');

CREATE TABLE IF NOT EXISTS invitations(
    invitation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID REFERENCES forms(form_id) ON DELETE CASCADE NOT NULL,
    invited_email VARCHAR(255) NOT NULL,
    invitation_token UUID DEFAULT gen_random_uuid(),
    invited_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    invited_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    status invitation_status DEFAULT 'invited',
    opened_at TIMESTAMP WITH TIME ZONE,
    submitted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS form_responses (
    response_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID REFERENCES forms(form_id) ON DELETE CASCADE NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    respondent_id UUID REFERENCES invitations(invitation_id) ON DELETE SET NULL,
    form_option_id UUID REFERENCES form_field_options(option_id) ON DELETE SET NULL,
    form_field_id UUID REFERENCES form_fields(field_id) ON DELETE CASCADE NOT NULL,
    response_text TEXT
);