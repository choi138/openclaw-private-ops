## ADDED Requirements

### Requirement: Secret Manager references for sensitive bootstrap inputs
The system MUST support Secret Manager reference inputs for sensitive values used by wg-easy and OpenClaw bootstrap (admin credentials, gateway password, and optional API/bot tokens).

#### Scenario: Secret references are provided
- **WHEN** a user supplies supported Secret Manager reference variables instead of plaintext variables
- **THEN** Terraform plans/applies using only references and not raw secret payload values

### Requirement: Runtime secret resolution through instance identity
The system MUST resolve Secret Manager-backed values at VM startup using the VM service account identity before starting wg-easy or OpenClaw services.

#### Scenario: Required secret can be resolved
- **WHEN** startup script fetches a required secret successfully from Secret Manager
- **THEN** the target service starts using the retrieved value

#### Scenario: Required secret cannot be resolved
- **WHEN** the required secret reference is invalid or access is denied
- **THEN** startup exits with a non-zero status and logs an actionable error

### Requirement: Least-privilege secret access
The system MUST grant secret access permissions only to the specific secrets required by each VM.

#### Scenario: VM reads an authorized secret
- **WHEN** a VM service account requests a configured secret that it is explicitly bound to
- **THEN** Secret Manager access is granted

#### Scenario: VM reads an unauthorized secret
- **WHEN** a VM service account requests a secret outside its configured bindings
- **THEN** Secret Manager access is denied

### Requirement: Backward compatibility for plaintext inputs
The system MUST continue to support existing plaintext sensitive variables when Secret Manager reference variables are not provided.

#### Scenario: Existing plaintext configuration is used
- **WHEN** users keep current plaintext-sensitive inputs and do not set secret references
- **THEN** module behavior remains compatible with current deployments

### Requirement: Mutual exclusivity of secret sources
The system MUST enforce mutually exclusive source rules for each credential that supports both plaintext and secret-reference forms.

#### Scenario: Both sources are set for one credential
- **WHEN** a plaintext value and a secret-reference value are both configured for the same credential
- **THEN** Terraform validation fails with a clear error message
