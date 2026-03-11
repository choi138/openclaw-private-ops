ALTER TABLE accounts
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'legacy';

ALTER TABLE conversations
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'legacy',
  ADD COLUMN IF NOT EXISTS external_id TEXT NULL;

ALTER TABLE messages
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'legacy',
  ADD COLUMN IF NOT EXISTS external_id TEXT NULL;

ALTER TABLE request_attempts
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'legacy',
  ADD COLUMN IF NOT EXISTS external_id TEXT NULL;

ALTER TABLE infra_snapshots
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'legacy';

ALTER TABLE accounts
  DROP CONSTRAINT IF EXISTS accounts_external_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_source_external_id
  ON accounts(source, external_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_conversations_source_external_id
  ON conversations(source, external_id)
  WHERE external_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_messages_source_external_id
  ON messages(source, external_id)
  WHERE external_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_request_attempts_source_external_id
  ON request_attempts(source, external_id)
  WHERE external_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_conversations_source_started
  ON conversations(source, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_infra_snapshots_source_captured
  ON infra_snapshots(source, captured_at DESC);

CREATE TABLE IF NOT EXISTS infra_status_latest (
  source TEXT PRIMARY KEY,
  snapshot_id BIGINT NOT NULL REFERENCES infra_snapshots(id) ON DELETE CASCADE,
  captured_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ingest_events (
  event_type TEXT NOT NULL,
  source TEXT NOT NULL,
  event_id TEXT NOT NULL,
  schema_version INTEGER NOT NULL,
  status TEXT NOT NULL,
  payload_json JSONB NOT NULL,
  last_error TEXT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 1,
  first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  next_retry_at TIMESTAMPTZ NULL,
  processed_at TIMESTAMPTZ NULL,
  dead_lettered_at TIMESTAMPTZ NULL,
  PRIMARY KEY (source, event_id)
);

CREATE INDEX IF NOT EXISTS idx_ingest_events_status_retry
  ON ingest_events(status, next_retry_at);

CREATE INDEX IF NOT EXISTS idx_ingest_events_first_seen
  ON ingest_events(first_seen_at);

CREATE INDEX IF NOT EXISTS idx_ingest_events_dead_letter
  ON ingest_events(dead_lettered_at)
  WHERE status = 'dead_letter';
