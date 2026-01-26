-- -----------------------------
-- DOWN Migration: Restore user_profile_pic_id and remove cascade delete
-- -----------------------------

BEGIN;

-- 1. Remove cascade delete from user_images
ALTER TABLE user_images
DROP CONSTRAINT IF EXISTS user_images_user_id_fkey;

ALTER TABLE user_images
ADD CONSTRAINT user_images_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(user_id);

-- 2. Re-add user_profile_pic_id column
ALTER TABLE users
ADD COLUMN user_profile_pic_id UUID;

-- 3. Restore UNIQUE constraint
ALTER TABLE users
ADD CONSTRAINT uq_user_profile_pic UNIQUE (user_profile_pic_id);

-- 4. Restore FOREIGN KEY constraint
ALTER TABLE users
ADD CONSTRAINT fk_user_profile_pic
FOREIGN KEY (user_profile_pic_id)
REFERENCES user_images(file_id);

COMMIT;
