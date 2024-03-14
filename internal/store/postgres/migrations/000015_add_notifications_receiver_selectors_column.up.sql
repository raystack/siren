ALTER TABLE
  notifications
ADD COLUMN IF NOT EXISTS receiver_selectors jsonb;