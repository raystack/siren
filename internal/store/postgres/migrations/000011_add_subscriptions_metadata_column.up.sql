ALTER TABLE
  subscriptions
ADD COLUMN IF NOT EXISTS metadata jsonb;
CREATE INDEX IF NOT EXISTS subscriptions_idx_metadata ON subscriptions USING GIN(metadata jsonb_path_ops);