## ADDED Requirements

### Requirement: Internal ingest endpoints are provided for operational producers
The system MUST provide internal ingestion endpoints for conversation events, infra snapshots, and request attempts.

#### Scenario: Conversation event ingestion
- **WHEN** an authorized internal producer calls `POST /v1/ingest/conversation-events` with a valid payload
- **THEN** the API persists the event and returns a success response

#### Scenario: Infra snapshot ingestion
- **WHEN** an authorized internal producer calls `POST /v1/ingest/infra-snapshot` with a valid payload
- **THEN** the API stores the snapshot fields with the reported capture timestamp

### Requirement: Ingest payloads are validated against contract rules
The system MUST validate required fields, field types, and schema version compatibility before processing ingest requests.

#### Scenario: Invalid ingest payload
- **WHEN** a producer sends an ingest request missing required fields or with unsupported schema version
- **THEN** the API rejects the request with a validation error and does not persist data

### Requirement: Event ingestion is idempotent
The system MUST prevent duplicate persistence when the same producer event is delivered multiple times.

#### Scenario: Duplicate event delivery
- **WHEN** the API receives two ingest requests with the same producer identity and event identifier
- **THEN** only one logical event is committed and the duplicate request is treated as already processed

### Requirement: Transient persistence failures are retried safely
The system MUST retry transient ingest persistence failures using bounded retry attempts and backoff.

#### Scenario: Temporary database outage
- **WHEN** an ingest write fails due to a transient database error
- **THEN** the event is scheduled for retry without violating idempotency guarantees

#### Scenario: Retry budget exhausted
- **WHEN** retries exceed the configured maximum attempts
- **THEN** the event is marked failed in a dead-letter path with operator-visible metadata

### Requirement: Ingest throughput target is supported
The system MUST support at least 100 ingest events per second for MVP traffic profiles.

#### Scenario: Sustained ingest load
- **WHEN** producers submit events at 100 events per second under valid contract inputs
- **THEN** the API processes events without sustained backlog growth beyond configured lag thresholds

### Requirement: Ingest access is restricted to authorized internal callers
The system MUST reject ingest requests from unauthorized or non-internal callers.

#### Scenario: Unauthorized ingest request
- **WHEN** a caller without valid ingest authorization calls an ingest endpoint
- **THEN** the API returns an authorization error and does not process payload data
