## Context

The backend plan sets explicit security objectives: visualize risks, return rule-based findings, and allow operators to track resolution state. Current project assets do not include a centralized security analysis service or persistent findings workflow.

## Goals / Non-Goals

**Goals:**
- Provide deterministic tfvars analysis API backed by codified rules.
- Return findings with normalized severity and lifecycle state fields.
- Persist findings and support query/filter flows for operator dashboards.
- Ensure analysis workflows are auditable and avoid raw secret leakage.

**Non-Goals:**
- Building a full SIEM pipeline with real-time threat detection.
- Parsing arbitrary IaC formats beyond scoped tfvars JSON for MVP.
- Auto-remediating infrastructure changes from findings.

## Decisions

- **Decision: Use an in-process rule engine with explicit rule interfaces.**
  - Rationale: keeps MVP fast to implement and test while supporting deterministic outputs.
  - Alternative considered: external policy engine dependency (rejected for initial complexity).

- **Decision: Require a minimum baseline of six shipped rules in MVP.**
  - Rationale: matches acceptance criteria and ensures immediate operator value.
  - Alternative considered: launch with fewer placeholder rules (rejected as insufficient coverage).

- **Decision: Persist findings with computed fingerprint hashes for deduplication.**
  - Rationale: repeated analyses should update lifecycle context, not explode duplicate rows.
  - Alternative considered: append-only findings for every run (rejected due to noisy triage experience).

- **Decision: Never store raw sensitive tfvars values in findings content.**
  - Rationale: analysis output should remain safe for broad internal visibility.
  - Alternative considered: storing full offending values for convenience (rejected for security risk).

## Risks / Trade-offs

- **[Risk]** Rule false positives may create alert fatigue.  
  **Mitigation:** severity tuning, clear rule rationale text, and lifecycle state transitions.

- **[Risk]** Rule set drift can cause inconsistent outputs across environments.  
  **Mitigation:** version rule bundles and expose rule version metadata in API responses.

- **[Trade-off]** In-process engine limits cross-language rule reuse.  
  **Mitigation:** keep rule inputs/outputs schema-driven to support later externalization.

## Migration Plan

1. Add `security_findings` schema and required indexes.
2. Implement rule interfaces and baseline rule set.
3. Implement analysis endpoint and persistence pipeline.
4. Add findings query/filter endpoint with lifecycle fields.
5. Roll frontend security panel to backend API and validate triage flows.

## Open Questions

- Which rule categories should be mandatory for launch (network exposure, weak credentials, logging gaps, etc.)?
- Should acknowledged/resolved transitions be API-driven in MVP or deferred to next phase?
- Do we need per-rule suppression annotations for known acceptable exceptions?
