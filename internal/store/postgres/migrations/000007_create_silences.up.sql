CREATE TABLE IF NOT EXISTS silences (
    id text PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_id bigint REFERENCES namespaces(id),
    type text,
    target_id text,
    target_expression jsonb,
    creator text,
    comment text,
    created_at timestamptz NOT NULL,
    deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS silences_idx_namespace_id_type_target_id ON silences (namespace_id, type, target_id) WHERE type = 'subscription';
CREATE INDEX IF NOT EXISTS silences_idx_namespace_id_type ON silences (namespace_id, type);
CREATE INDEX IF NOT EXISTS silences_idx_target_expression ON silences USING GIN(target_expression jsonb_path_ops) WHERE type = 'labels';
