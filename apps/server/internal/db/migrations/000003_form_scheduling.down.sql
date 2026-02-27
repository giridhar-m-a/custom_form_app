BEGIN;

-- 1. Temporarily convert enum column to TEXT
ALTER TABLE invitations
  ALTER COLUMN status DROP DEFAULT;

ALTER TABLE invitations
  ALTER COLUMN status TYPE TEXT
  USING status::text;

-- 2. Normalize values safely
UPDATE invitations
SET status = 'invited'
WHERE status IN ('pending', 'accepted', 'failed');

-- 3. Drop new enum type
DROP TYPE invitation_status;

-- 4. Recreate old enum
CREATE TYPE invitation_status AS ENUM ('opened', 'submitted', 'invited', 'expired');

-- 5. Convert column back to old enum
ALTER TABLE invitations
  ALTER COLUMN status TYPE invitation_status
  USING status::invitation_status;

-- 6. Restore default
ALTER TABLE invitations
  ALTER COLUMN status SET DEFAULT 'invited';

-- ALTER TABLE invitations
--   DROP COLUMN resend_id;

-- 7. Drop index
DROP INDEX IF EXISTS idx_forms_scheduling_id;
DROP INDEX IF EXISTS idx_forms_invitation_schedule_id;

-- 8. Drop constraints
ALTER TABLE forms DROP CONSTRAINT invitation_schedule_gap_constraint;

-- 9. Drop scheduling columns
ALTER TABLE forms DROP COLUMN IF EXISTS scheduling_id;
ALTER TABLE forms DROP COLUMN IF EXISTS scheduled_time;
ALTER TABLE forms DROP COLUMN IF EXISTS closing_time;
ALTER TABLE forms DROP COLUMN IF EXISTS is_schedule_completed;
ALTER TABLE forms DROP COLUMN IF EXISTS is_scheduled;
ALTER TABLE forms DROP COLUMN IF EXISTS invitation_schedule_id;
ALTER TABLE forms DROP COLUMN IF EXISTS invitation_schedule_gap;

-- 10. Drop resend_id
ALTER TABLE invitations DROP COLUMN IF EXISTS resend_id;

COMMIT;
