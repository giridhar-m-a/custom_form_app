--1. Add invitation_token column to invitations table
ALTER TABLE invitations
ADD COLUMN invitation_token UUID DEFAULT gen_random_uuid();

--2. Drop invited_name column from invitations table
ALTER TABLE invitations
DROP COLUMN invited_name;

ALTER TABLE invitations
DROP CONSTRAINT IF EXISTS invitations_invited_email_form_id_key;

--3. Drop response_options table
DROP TABLE response_options;
