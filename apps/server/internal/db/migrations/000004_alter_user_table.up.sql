-- -----------------------------
-- Migration: Remove user_profile_pic_id from users
-- -----------------------------

BEGIN;

-- 1. Drop UNIQUE constraint (if exists)
ALTER TABLE users
DROP CONSTRAINT IF EXISTS uq_user_profile_pic;

-- 2. Drop FOREIGN KEY constraint (if exists)
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_user_profile_pic;

-- 3. Drop column
ALTER TABLE users
DROP COLUMN IF EXISTS user_profile_pic_id;

-- 4. Make user image cascade delete
ALTER TABLE user_images
DROP CONSTRAINT IF EXISTS user_images_user_id_fkey;

ALTER TABLE user_images
ADD CONSTRAINT user_images_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(user_id)
ON DELETE CASCADE;

ALTER TABLE user_images
ALTER COLUMN user_id SET NOT NULL;

COMMIT;
