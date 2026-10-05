
-- 1. Add a column to flag an user as temp user
ALTER TABLE users
ADD COLUMN is_temp BOOLEAN DEFAULT FALSE;

--2. Add a column to flag an user as an deleted user
ALTER TABLE users
ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

--3. Add a column to flag an form as deleted
ALTER TABLE forms
ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

--4. Make Email field in users optional
ALTER TABLE users
ALTER COLUMN user_email DROP NOT NULL;

