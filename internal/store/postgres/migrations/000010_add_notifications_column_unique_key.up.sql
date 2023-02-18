ALTER TABLE
  notifications
ADD COLUMN IF NOT EXISTS unique_key text;

CREATE INDEX IF NOT EXISTS notifications_idx_unique_key ON notifications(unique_key);