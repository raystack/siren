CREATE TABLE IF NOT EXISTS idempotencies (
    id bigserial PRIMARY KEY,
    scope text not null,
    key text not null,
    success boolean,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

CREATE UNIQUE INDEX idempotencies_keys_scope_key ON idempotencies (scope, key);
