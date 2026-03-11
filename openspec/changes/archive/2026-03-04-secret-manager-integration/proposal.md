## Why

The module currently accepts sensitive values as Terraform inputs and passes them into startup metadata, which increases exposure in state files and operational logs. We need a production-safe secret path that keeps credentials in GCP Secret Manager and retrieves them only at VM boot time.

## What Changes

- Add optional Secret Manager reference inputs for sensitive values used by wg-easy and OpenClaw.
- Keep existing plaintext inputs for backward compatibility, but enforce mutually exclusive source rules where both plain and secret-ref forms exist.
- Add VM identity and IAM wiring so each VM can read only the secrets it needs.
- Update startup scripts to fetch secrets from Secret Manager at runtime and fail fast with clear logs on missing access or invalid references.
- Update docs/examples to show secure migration from plaintext variables to secret references.

## Capabilities

### New Capabilities
- (none)

### Modified Capabilities
- `openclaw-vpn-only`: Extend provisioning and bootstrap flows to support Secret Manager-backed secrets without breaking existing VPN-only behavior.

## Impact

- Terraform files: `variables.tf`, `main.tf`, `templates/startup.sh.tpl`, `templates/startup-openclaw.sh.tpl`, `README.md`, `examples/basic/*`.
- GCP dependencies: Secret Manager access via instance identity and IAM (`roles/secretmanager.secretAccessor`).
- Operations: reduced sensitive data exposure in Terraform state/metadata; stricter bootstrap validation for secret access.
