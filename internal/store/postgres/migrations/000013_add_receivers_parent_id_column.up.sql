ALTER TABLE
  receivers
ADD COLUMN IF NOT EXISTS parent_id bigint;

CREATE INDEX IF NOT EXISTS receivers_idx_parent_id ON receivers(parent_id);