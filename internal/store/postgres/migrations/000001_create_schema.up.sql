CREATE TABLE IF NOT EXISTS providers (
    id bigserial PRIMARY KEY,
    host text,
    urn text UNIQUE,
    name text,
    type text,
    credentials jsonb,
    labels jsonb,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS alerts (
    id bigserial PRIMARY KEY,
    provider_id bigint REFERENCES providers(id),
    resource_name text,
    metric_name text,
    metric_value text,
    severity text,
    rule text,
    triggered_at timestamp,
    created_at timestamptz,
    updated_at timestamptz
 );

 CREATE TABLE IF NOT EXISTS namespaces (
    id bigserial PRIMARY KEY,
    provider_id bigint REFERENCES providers(id),
    urn text,
    name text,
    credentials text,
    labels jsonb,
    created_at timestamptz,
    updated_at timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS urn_provider_id_unique ON namespaces (provider_id,urn);

CREATE TABLE IF NOT EXISTS receivers (
    id bigserial PRIMARY KEY,
    name text,
    type text,
    labels jsonb,
    configurations jsonb,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS templates (
    id bigserial PRIMARY KEY,
    name text UNIQUE,
    body text,
    tags text[],
    variables jsonb,
    created_at timestamptz,
    updated_at timestamptz
);
CREATE INDEX IF NOT EXISTS idx_tags  ON templates USING gin(tags);

CREATE TABLE IF NOT EXISTS rules (
    id bigserial PRIMARY KEY,
    name text UNIQUE,
    namespace text,
    group_name text,
    template text,
    enabled boolean,
    variables jsonb,
    provider_namespace bigint REFERENCES namespaces(id),
    created_at timestamptz,
    updated_at timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_name ON rules (namespace,group_name,template,provider_namespace);

CREATE TABLE IF NOT EXISTS subscriptions (
    id bigserial PRIMARY KEY,
    namespace_id bigint REFERENCES namespaces(id),
    urn text UNIQUE,
    receiver jsonb,
    match jsonb,
    created_at timestamptz,
    updated_at timestamptz
);
