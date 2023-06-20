ALTER TABLE
  subscriptions
ADD COLUMN IF NOT EXISTS created_by text,
ADD COLUMN IF NOT EXISTS updated_by text;