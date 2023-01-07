ALTER TABLE
  alerts
ADD COLUMN IF NOT EXISTS namespace_id bigint REFERENCES namespaces(id);