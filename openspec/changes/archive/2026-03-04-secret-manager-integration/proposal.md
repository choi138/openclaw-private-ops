## Why

The module currently accepts sensitive values as Terraform inputs and passes them into startup metadata, which increases exposure in state files and operational logs. We need a production-safe secret path that keeps credentials in GCP Secret Manager and retrieves them only at VM boot time.

## What Changes

- Add optional Secret Manager reference inputs for sensitive values used by wg-easy and OpenClaw.
- Keep existing plaintext inputs for backward compatibility, but enforce mutually exclusive source rules where both plain and secret-ref forms exist.
- Add VM identity and secret-scoped IAM wiring so each VM can read only the specific secrets it needs, without project-wide Secret Manager access.
- Update startup scripts to fetch secrets from Secret Manager at runtime and fail fast with sanitized logs on missing access or invalid references; never log secret values, and redact secret identifiers or full resource paths where practical.
- Update docs/examples to show secure migration from plaintext variables to secret references.

## Capabilities

### New Capabilities
- (none)

### Modified Capabilities
- `openclaw-vpn-only`: Extend provisioning and bootstrap flows to support Secret Manager-backed secrets without breaking existing VPN-only behavior.

## Impact

- Terraform files: `variables.tf`, `main.tf`, `templates/startup.sh.tpl`, `templates/startup-openclaw.sh.tpl`, `README.md`, `examples/basic/*`.
- GCP dependencies: Secret Manager access via instance identity and secret-level IAM bindings (`roles/secretmanager.secretAccessor` only on explicitly allowed secrets, or equivalent conditional IAM).
- Operations: reduced sensitive data exposure in Terraform state/metadata; stricter bootstrap validation for secret access.
