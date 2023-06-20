ALTER TABLE
  subscriptions
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by;