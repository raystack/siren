DROP INDEX IF EXISTS receivers_idx_parent_id;

ALTER TABLE
  receivers
DROP COLUMN IF EXISTS parent_id;