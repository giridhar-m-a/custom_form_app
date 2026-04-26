ALTER TABLE users
DROP COLUMN is_temp;

ALTER TABLE users
DROP COLUMN is_deleted;

ALTER TABLE forms
DROP COLUMN is_deleted;

ALTER TABLE users
ALTER COLUMN user_email SET NOT NULL;
