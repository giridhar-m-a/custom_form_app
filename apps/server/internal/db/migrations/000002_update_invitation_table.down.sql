--1. Add invitation_token column to invitations table
ALTER TABLE invitations
ADD COLUMN invitation_token UUID DEFAULT gen_random_uuid();

--2. Drop invited_name column from invitations table
ALTER TABLE invitations
DROP COLUMN invited_name;

ALTER TABLE invitations
DROP CONSTRAINT IF EXISTS invitations_invited_email_form_id_key;

--3. Add form_option_id column to form_responses table
ALTER TABLE form_responses
ADD COLUMN form_option_id UUID REFERENCES form_field_options(option_id) ON DELETE SET NULL;

--4. Update form_option_id column in form_responses table
UPDATE form_responses
SET form_option_id = response_options.form_option_id
FROM response_options
WHERE response_options.response_id = form_responses.response_id;

--5. Create index on form_option_id column in form_responses table
CREATE INDEX IF NOT EXISTS idx_form_responses_option_id ON form_responses(form_option_id);

--6. Drop response_options table
DROP TABLE response_options;
