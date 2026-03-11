## Context

The backend plan expects dashboard and timeline features to consume near-real-time operational data. Without a durable ingest layer, UI views will drift from runtime state and duplicate event delivery can corrupt aggregates.

## Goals / Non-Goals

**Goals:**
- Provide stable internal ingest endpoints for all MVP producer types.
- Enforce schema and payload validation before persistence.
- Guarantee idempotent writes for retried or duplicated events.
- Add retry logic that recovers transient storage failures without unbounded loops.

**Non-Goals:**
- Building a full stream-processing platform.
- Introducing external public ingestion access.
- Solving long-term archival analytics beyond operational windows.

## Decisions

- **Decision: Define three explicit ingest endpoints matching producer intent.**
  - Rationale: keeps contracts simple (`conversation-events`, `infra-snapshot`, `request-attempt`) and aligns with plan.
  - Alternative considered: generic single ingest endpoint with dynamic payload type (rejected due to validation complexity).

- **Decision: Require producer event identity (`source + event_id`) for idempotency.**
  - Rationale: deterministic dedup prevents double-counted metrics under retries.
  - Alternative considered: timestamp-window dedup heuristics (rejected due to false positives/negatives).

- **Decision: Persist idempotency ledger in Postgres with unique constraints.**
  - Rationale: reuse primary datastore consistency guarantees and avoid extra infra for MVP.
  - Alternative considered: Redis-only dedup cache (rejected because restarts can lose dedup history).

- **Decision: Apply bounded retry with exponential backoff for transient DB failures.**
  - Rationale: improves delivery reliability while protecting DB from thundering herd retries.
  - Alternative considered: immediate fail-fast only (rejected because transient outages are expected).

## Risks / Trade-offs

- **[Risk]** Strict contract validation may drop producer events during rollout mismatches.  
  **Mitigation:** version event schemas and return actionable validation errors.

- **[Risk]** Retry queues can grow under prolonged DB incidents.  
  **Mitigation:** enforce queue depth limits, dead-lettering, and lag alerts.

- **[Trade-off]** Postgres-backed idempotency adds write amplification.  
  **Mitigation:** index event identity columns and prune old idempotency rows by retention policy.

## Migration Plan

1. Add ingest ledger/metadata tables and indexes.
2. Implement endpoint handlers with validation and identity extraction.
3. Add repository methods for atomic dedup + persistence.
4. Introduce retry worker path for transient failure cases.
5. Roll producers to send event identity headers and validate end-to-end ingest lag.

## Open Questions

- Should Redis be introduced in MVP for retry buffering, or deferred until observed load requires it?
- What is the acceptable dead-letter threshold before paging operators?
- Do we need per-producer authentication tokens in addition to VPN/internal routing constraints?
