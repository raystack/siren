ALTER TABLE
  idempotencies
ADD COLUMN IF NOT EXISTS notification_id text;

ALTER TABLE
  idempotencies
DROP COLUMN IF EXISTS success;