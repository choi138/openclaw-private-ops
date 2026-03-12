ALTER TABLE infra_snapshots
  ADD COLUMN IF NOT EXISTS event_id TEXT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_infra_snapshots_source_event_id
  ON infra_snapshots(source, event_id)
  WHERE event_id IS NOT NULL;
