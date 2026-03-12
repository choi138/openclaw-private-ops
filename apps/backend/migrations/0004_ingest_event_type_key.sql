ALTER TABLE ingest_events
  DROP CONSTRAINT IF EXISTS ingest_events_pkey;

ALTER TABLE ingest_events
  ADD CONSTRAINT ingest_events_pkey PRIMARY KEY (event_type, source, event_id);
