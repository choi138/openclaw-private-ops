## Why

The Ops Console needs clear, actionable security risk visibility rather than manual `tfvars` inspection. We need a backend security analysis API that evaluates configuration input with deterministic rules and persists findings for follow-up.

## What Changes

- Implement `POST /v1/security/analyze-tfvars` to analyze provided tfvars JSON and return structured findings.
- Introduce a rule engine with severity levels and at least six baseline security checks.
- Persist findings with lifecycle state (`open`, `acknowledged`, `resolved`) and timestamps.
- Provide findings query endpoint with status/severity filtering.
- Add audit records for security analysis and findings-read operations.

## Capabilities

### New Capabilities
- `ops-api-security-analysis`: Rule-driven tfvars security analysis, finding lifecycle persistence, and query APIs for Ops workflows.

### Modified Capabilities
- (none)

## Impact

- New backend security domain/service/repository packages and API handlers.
- Additional database schema for findings state/history and rule metadata.
- Frontend security panel switches from local validation to backend findings API.
- Operational process updates for finding triage and remediation tracking.
