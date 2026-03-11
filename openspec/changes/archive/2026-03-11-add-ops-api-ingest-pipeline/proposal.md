## Why

Core read APIs alone cannot provide reliable real-time operational insight unless backend data ingestion is standardized. We need an internal ingestion pipeline that persists OpenClaw and infra events with idempotent handling.

## What Changes

- Add internal ingest endpoints for conversation events, infra snapshots, and request attempts.
- Introduce payload validation and schema-version checks for ingest contracts.
- Implement idempotent write behavior using event identity keys to prevent duplicate persistence.
- Add retry handling for transient persistence failures with bounded backoff strategy.
- Expose ingest health indicators so operators can detect lag and drop conditions.

## Capabilities

### New Capabilities
- `ops-api-ingestion`: Internal event ingestion path with validation, idempotency, retry behavior, and ingest observability.

### Modified Capabilities
- (none)

## Impact

- Backend API write-path expansion under `backend/internal/ingest`, `backend/internal/worker`, and repository layers.
- New database structures for event idempotency state and ingest processing metadata.
- Upstream OpenClaw/WireGuard collectors must send contract-compliant payloads and event identifiers.
- Additional operational tuning for queue depth, retry policies, and ingestion alerting.
