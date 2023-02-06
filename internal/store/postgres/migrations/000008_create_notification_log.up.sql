CREATE TABLE IF NOT EXISTS notification_log (
    id text PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_id bigint,
    notification_id text REFERENCES notifications(id),
    subscription_id bigint,
    receiver_id bigint,
    alert_ids bigint[],
    silence_ids text[],
    created_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS notification_log_idx_silence_ids ON notification_log USING GIN(silence_ids);