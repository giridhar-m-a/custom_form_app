BEGIN;

-- 1. Rename enum
ALTER TYPE invitation_status RENAME TO invitation_status_old;

-- 2. Create new enum
CREATE TYPE invitation_status AS ENUM (
  'pending',
  'bounced',
  'clicked',
  'opened',
  'delivered',
  'complained',
  'delayed',
  'failed',
  'submitted'
);

-- 3. Normalize data while still on old enum
UPDATE invitations
SET status = 'submitted'
WHERE status IN ('opened', 'invited', 'expired');

-- 4. Convert column via TEXT bridge
ALTER TABLE invitations
  ALTER COLUMN status DROP DEFAULT,
  ALTER COLUMN status TYPE invitation_status
  USING status::text::invitation_status;

-- 5. Drop old enum
DROP TYPE invitation_status_old;

-- 6. Add scheduling columns
ALTER TABLE forms ADD COLUMN IF NOT EXISTS scheduling_id UUID;
ALTER TABLE forms ADD COLUMN IF NOT EXISTS scheduled_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE forms ADD COLUMN IF NOT EXISTS closing_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE forms ADD COLUMN IF NOT EXISTS is_schedule_completed BOOLEAN DEFAULT FALSE;
ALTER TABLE forms ADD COLUMN IF NOT EXISTS is_scheduled BOOLEAN DEFAULT FALSE;
ALTER TABLE forms ADD COLUMN IF NOT EXISTS invitation_schedule_gap INT NULL; -- time gap between the time when the form is scheduled to go live and the time when the invitation is scheduled to be sent
ALTER TABLE forms ADD COLUMN IF NOT EXISTS invitation_schedule_id UUID;

-- 7. Add constraints to ensure that invitation_schedule_gap is not null when is_scheduled is true
ALTER TABLE forms
ADD CONSTRAINT invitation_schedule_gap_constraint
CHECK (
    NOT is_scheduled
    OR (
        scheduled_time IS NOT NULL
        AND invitation_schedule_gap IS NOT NULL
    )
);


CREATE INDEX IF NOT EXISTS idx_forms_scheduling_id ON forms(scheduling_id);
CREATE INDEX IF NOT EXISTS idx_forms_invitation_schedule_id ON forms(invitation_schedule_id);

-- 7. Set new default
ALTER TABLE invitations
  ALTER COLUMN status SET DEFAULT 'pending';

COMMIT;
