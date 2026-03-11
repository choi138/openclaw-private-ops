## Why

As the Ops API matures, operating it safely requires stronger observability, role-based authentication, and retention controls. We need to move beyond basic token-only access and add production-grade operations safeguards.

## What Changes

- Add observability foundations: Prometheus metrics endpoint, OpenTelemetry tracing hooks, and standardized structured logs.
- Introduce OIDC-based authentication and role-based authorization (`admin`, `viewer`, `auditor`).
- Define controlled break-glass token behavior for emergency access paths.
- Implement automated retention jobs for raw message payloads and aged aggregate records.
- Establish alert signal contracts for error rate, ingest lag, and service availability.

## Capabilities

### New Capabilities
- `ops-api-operations`: Operational hardening for authn/authz, observability, and retention lifecycle controls.

### Modified Capabilities
- (none)

## Impact

- Backend middleware and auth/session components for OIDC token validation and role checks.
- Telemetry configuration (`/metrics`, trace exporter wiring, log field normalization).
- Background worker/job components for retention cleanup policies.
- Runbook and deployment configuration updates for identity provider and observability stack integration.
