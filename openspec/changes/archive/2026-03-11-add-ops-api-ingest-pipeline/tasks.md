## 1. Ingest Contract and Storage Foundation

- [x] 1.1 Define request schemas for conversation events, infra snapshots, and request attempts
- [x] 1.2 Add migration tables/indexes for idempotency ledger and ingest metadata
- [x] 1.3 Implement payload validation and schema-version checks in ingest DTO layer

## 2. Endpoint and Persistence Implementation

- [x] 2.1 Implement `POST /v1/ingest/conversation-events` with atomic dedup + write path
- [x] 2.2 Implement `POST /v1/ingest/infra-snapshot` with latest snapshot upsert/append logic
- [x] 2.3 Implement `POST /v1/ingest/request-attempt` with provider/model metrics persistence

## 3. Retry and Reliability Controls

- [x] 3.1 Add transient-failure classification and retry scheduling with bounded exponential backoff
- [x] 3.2 Implement dead-letter handling and operator-visible failure counters
- [x] 3.3 Add ingest lag and queue depth metrics for alert integration

## 4. Verification and Producer Rollout

- [x] 4.1 Add tests for idempotency, duplicate deliveries, and retry behavior
- [x] 4.2 Document producer contract requirements (event identity, schema version, auth expectations)
- [x] 4.3 Validate ingestion throughput target (>=100 events/sec) under representative load
