CREATE TABLE IF NOT EXISTS notifications (
    id text PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_id bigint,
    type text,
    data jsonb,
    labels jsonb,
    valid_duration text,
    template text,
    created_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS notifications_idx_labels ON notifications USING GIN(labels jsonb_path_ops);
CREATE INDEX IF NOT EXISTS notifications_idx_template ON notifications (template);