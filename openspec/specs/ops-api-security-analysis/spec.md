## Purpose

Define the internal Ops API requirements for tfvars-based security analysis, including deterministic rules, persisted finding lifecycle state, safe redaction, and auditable operator workflows.

## Requirements

### Requirement: Tfvars security analysis endpoint is available
The system MUST provide an endpoint that accepts tfvars JSON input and returns security findings, and it MUST require an admin-authenticated bearer token for access.

#### Scenario: Analyze valid tfvars payload
- **WHEN** an admin-authenticated client calls `POST /v1/security/analyze-tfvars` with valid tfvars JSON
- **THEN** the API returns a structured list of findings including rule identifier, severity, and description fields

#### Scenario: Analyze request is unauthenticated or uses the wrong token
- **WHEN** a client calls `POST /v1/security/analyze-tfvars` without the admin bearer token
- **THEN** the API responds with `401 Unauthorized`

### Requirement: MVP baseline includes at least six security rules
The system MUST evaluate input against a baseline rule set containing at least six distinct security checks.

#### Scenario: Input violates multiple baseline rules
- **WHEN** analyzed tfvars contains conditions that violate several baseline checks
- **THEN** the response includes findings for each triggered rule from the baseline set

### Requirement: Finding severities and lifecycle states are normalized
The system MUST represent finding severity as one of `critical`, `high`, `medium`, or `info`, and state as `open`, `acknowledged`, or `resolved`.

#### Scenario: New finding persistence
- **WHEN** analysis generates a previously unseen finding fingerprint
- **THEN** the persisted finding is created with a normalized severity and default state `open`

### Requirement: Findings are persisted and queryable
The system MUST persist findings and provide query access with status and severity filtering, and it MUST require an admin-authenticated bearer token for reads.

#### Scenario: Query open findings
- **WHEN** an admin-authenticated client calls `GET /v1/security/findings?status=open`
- **THEN** the API returns only findings whose current lifecycle state is `open`

#### Scenario: Query by severity
- **WHEN** a client filters findings by severity
- **THEN** the API returns findings limited to the requested severity values

#### Scenario: Findings query is unauthenticated or uses the wrong token
- **WHEN** a client calls `GET /v1/security/findings` without the admin bearer token
- **THEN** the API responds with `401 Unauthorized`

### Requirement: Sensitive input values are not exposed in findings storage
The system MUST avoid storing raw sensitive tfvars values in finding titles, descriptions, or metadata fields.

#### Scenario: Rule triggers on sensitive field
- **WHEN** a rule evaluates a sensitive configuration field
- **THEN** persisted and returned finding data contains redacted or masked representations instead of raw secret values

### Requirement: Security analysis actions are auditable
The system MUST record audit events for security analysis execution and findings retrieval.

#### Scenario: Analysis request succeeds
- **WHEN** a client successfully invokes tfvars analysis
- **THEN** the API writes an audit event containing actor, action, and timestamp

#### Scenario: Findings query succeeds
- **WHEN** a client reads findings from the query endpoint
- **THEN** the API writes an audit event for the read action
