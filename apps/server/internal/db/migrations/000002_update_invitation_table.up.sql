-- 1. Add invited_name column to invitations table
ALTER TABLE invitations
ADD COLUMN invited_name VARCHAR(255);

-- 2. Update invited_name column in invitations table
UPDATE invitations
SET invited_name = invited_email;

-- 3. Make invited_name column NOT NULL
ALTER TABLE invitations
ALTER COLUMN invited_name SET NOT NULL;

-- 4. Drop invitation_token column from invitations table
ALTER TABLE invitations
DROP COLUMN invitation_token;

-- Add unique constraint on invited_email + form_id
ALTER TABLE invitations
ADD CONSTRAINT invitations_invited_email_form_id_key
UNIQUE (invited_email, form_id);

-- 5. Create response_options table
CREATE TABLE response_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    response_id UUID NOT NULL
        REFERENCES form_responses(response_id) ON DELETE CASCADE,
    form_option_id UUID NOT NULL
        REFERENCES form_field_options(option_id) ON DELETE CASCADE,
    UNIQUE (response_id, form_option_id)
);

-- 6. Create indexes on response_options table
CREATE INDEX IF NOT EXISTS idx_response_options_response_id ON response_options(response_id);
CREATE INDEX IF NOT EXISTS idx_response_options_form_option_id ON response_options(form_option_id);

-- 7. Insert data into response_options table
INSERT INTO response_options (response_id, form_option_id)
SELECT response_id, form_option_id
FROM form_responses
WHERE form_option_id IS NOT NULL;

-- 8. Drop idx_form_responses_option_id index
DROP INDEX IF EXISTS idx_form_responses_option_id;

-- 9. Drop form_option_id column from form_responses table
ALTER TABLE form_responses
DROP COLUMN form_option_id;