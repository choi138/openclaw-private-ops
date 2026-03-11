## Context

This module provisions a WireGuard VM and a private OpenClaw VM. Today, credentials and API tokens are provided as Terraform variables and rendered into startup scripts. Even with `sensitive = true`, secret values can still be exposed through state handling workflows, execution logs, and metadata/script surfaces. The change introduces a secure, runtime retrieval path using Secret Manager while preserving current behavior for existing users.

## Goals / Non-Goals

**Goals:**
- Support Secret Manager references for wg-easy and OpenClaw sensitive values.
- Ensure secrets are resolved at boot time by VM identity, not injected as Terraform-managed plaintext.
- Enforce least-privilege IAM for secret access.
- Preserve backward compatibility for users still using plaintext variables.

**Non-Goals:**
- Automatically creating/managing secret values in Secret Manager.
- Replacing all existing variables with breaking changes in one release.
- Designing a full secret rotation orchestrator.

## Decisions

- **Decision: Add secret-reference variables alongside existing plaintext variables.**
  - Rationale: non-breaking adoption path and incremental migration.
  - Alternative considered: hard switch to Secret Manager-only inputs (rejected as a breaking change).

- **Decision: Use dedicated VM service accounts and per-secret IAM bindings.**
  - Rationale: explicit least privilege and auditable access boundaries.
  - Alternative considered: reusing default Compute Engine service account with broader project-level access (rejected for over-privilege risk).

- **Decision: Retrieve secrets at startup via metadata token + Secret Manager REST API.**
  - Rationale: avoids storing secret values in Terraform state and avoids extra heavyweight dependencies.
  - Alternative considered: reading secrets through Terraform data sources (rejected because secret payloads can end up in state).

- **Decision: Fail fast on required secret retrieval failure.**
  - Rationale: safer than starting services with empty or invalid credentials.
  - Alternative considered: silent fallback behavior (rejected because it hides misconfiguration).

## Risks / Trade-offs

- **[Risk]** Misconfigured secret references or IAM permissions can break first boot.  
  **Mitigation:** Add strict variable validation, clear startup error messages, and README troubleshooting steps.

- **[Risk]** Secret rotation does not take effect until service restart/reboot.  
  **Mitigation:** Document restart procedure and optional operational runbook for rotation windows.

- **[Trade-off]** Additional Terraform resources (service accounts/IAM) increase module complexity.  
  **Mitigation:** Keep defaults simple and expose minimal new variables with clear examples.

## Migration Plan

1. Add new secret-reference inputs and IAM/instance wiring while keeping existing plaintext inputs intact.
2. Update startup scripts to support both paths (secret reference preferred when set).
3. Update docs/examples with migration guidance from plaintext to Secret Manager.
4. Roll out by first creating secrets and IAM, then switching tfvars to secret references.
5. Rollback path: restore previous plaintext vars and re-apply.

## Open Questions

- Should secret references allow explicit version pinning, or default only to `latest` initially?
- Should this module optionally manage `google_project_service` for Secret Manager API enablement, or require pre-enabled projects?
