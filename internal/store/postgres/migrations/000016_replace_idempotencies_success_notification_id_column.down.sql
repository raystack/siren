ALTER TABLE
  notifications
DROP COLUMN IF EXISTS receiver_selectors;

ALTER TABLE
  idempotencies
ADD COLUMN IF NOT EXISTS success boolean;