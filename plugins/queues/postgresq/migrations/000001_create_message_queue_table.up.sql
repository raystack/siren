CREATE TABLE IF NOT EXISTS message_queue (
   id text NOT NULL PRIMARY KEY,
   status text NOT NULL, -- ENQUEUED/RUNNING/FAILED/DONE
   receiver_type text NOT NULL,
   configs jsonb,
   details jsonb,
   last_error text,
   max_tries integer,
   try_count integer,
   retryable boolean NOT NULL DEFAULT false,
   expired_at timestamptz,
   created_at timestamptz NOT NULL,
   updated_at timestamptz NOT NULL
);