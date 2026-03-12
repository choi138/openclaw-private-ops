## 1. Input Interface and Validation

- [x] 1.1 Add Secret Manager reference variables for wg-easy and OpenClaw sensitive inputs in `variables.tf`
- [x] 1.2 Add validation/check logic to enforce mutually exclusive plaintext vs secret-reference inputs
- [x] 1.3 Update `examples/basic` variable wiring and sample tfvars for Secret Manager usage

## 2. IAM and Compute Wiring

- [x] 2.1 Add dedicated service accounts for VPN/OpenClaw instances (or documented equivalent identity model)
- [x] 2.2 Add IAM bindings to grant `roles/secretmanager.secretAccessor` only for referenced secrets
- [x] 2.3 Attach service accounts and required scopes to each Compute Engine instance

## 3. Startup Secret Retrieval

- [x] 3.1 Update `templates/startup.sh.tpl` to fetch wg-easy secrets from Secret Manager at boot
- [x] 3.2 Update `templates/startup-openclaw.sh.tpl` to fetch OpenClaw secrets from Secret Manager at boot
- [x] 3.3 Add fail-fast handling and clear log messages for secret retrieval failures

## 4. Documentation and Verification

- [x] 4.1 Update README with Secret Manager setup, IAM prerequisites, and migration guidance
- [x] 4.2 Add verification steps in `tests/README.md` for Secret Manager path and backward-compatible plaintext path
- [x] 4.3 Validate with `terraform fmt`, `terraform validate`, and TFLint for root and example configurations
