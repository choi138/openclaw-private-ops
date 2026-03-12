## 1. Authentication and Authorization Hardening

- [ ] 1.1 Implement OIDC token validation middleware and provider configuration checks
- [ ] 1.2 Define route-level RBAC policies for `admin`, `viewer`, and `auditor` roles
- [ ] 1.3 Add guarded break-glass token flow with explicit enable/disable controls and audit logging

## 2. Observability Baseline

- [ ] 2.1 Add Prometheus-compatible `/metrics` endpoint for request, latency, error, and ingest lag metrics
- [ ] 2.2 Add OpenTelemetry tracing middleware with configurable exporter and sampling defaults
- [ ] 2.3 Standardize structured log fields (`request_id`, `principal_id`, `path`, `latency_ms`, `status`, `error_code`) and keep end-user identifiers optional to avoid over-logging PII

## 3. Retention Automation

- [ ] 3.1 Implement scheduled retention jobs for raw message payloads and aged aggregates
- [ ] 3.2 Add dry-run mode and deletion summary reporting before enforced cleanup
- [ ] 3.3 Add retention configuration and safety guards (minimum window, max rows per run)

## 4. Operational Validation and Documentation

- [ ] 4.1 Add tests for RBAC policy enforcement, OIDC validation failures, and break-glass behavior
- [ ] 4.2 Validate alert signals for error-rate, ingest-lag, and service-down conditions
- [ ] 4.3 Update runbooks and deployment docs for observability and identity configuration
