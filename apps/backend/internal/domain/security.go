package domain

import (
	"slices"
	"time"
)

const SupportedSecurityAnalysisSchemaVersion = 1

type SecuritySeverity string

const (
	SecuritySeverityCritical SecuritySeverity = "critical"
	SecuritySeverityHigh     SecuritySeverity = "high"
	SecuritySeverityMedium   SecuritySeverity = "medium"
	SecuritySeverityInfo     SecuritySeverity = "info"
)

var allowedSecuritySeverities = []SecuritySeverity{
	SecuritySeverityCritical,
	SecuritySeverityHigh,
	SecuritySeverityMedium,
	SecuritySeverityInfo,
}

func IsAllowedSecuritySeverity(v string) bool {
	return slices.Contains(allowedSecuritySeverities, SecuritySeverity(v))
}

type SecurityFindingStatus string

const (
	SecurityFindingStatusOpen         SecurityFindingStatus = "open"
	SecurityFindingStatusAcknowledged SecurityFindingStatus = "acknowledged"
	SecurityFindingStatusResolved     SecurityFindingStatus = "resolved"
)

var allowedSecurityFindingStatuses = []SecurityFindingStatus{
	SecurityFindingStatusOpen,
	SecurityFindingStatusAcknowledged,
	SecurityFindingStatusResolved,
}

func IsAllowedSecurityFindingStatus(v string) bool {
	return slices.Contains(allowedSecurityFindingStatuses, SecurityFindingStatus(v))
}

type SecurityFinding struct {
	ID              int64                 `json:"id"`
	Fingerprint     string                `json:"fingerprint"`
	RuleID          string                `json:"rule_id"`
	RuleVersion     string                `json:"rule_version"`
	Severity        SecuritySeverity      `json:"severity"`
	Status          SecurityFindingStatus `json:"status"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	FieldPath       string                `json:"field_path,omitempty"`
	FixHint         string                `json:"fix_hint,omitempty"`
	Metadata        map[string]any        `json:"metadata,omitempty"`
	FirstDetectedAt time.Time             `json:"first_detected_at"`
	LastDetectedAt  time.Time             `json:"last_detected_at"`
	ResolvedAt      *time.Time            `json:"resolved_at,omitempty"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

type SecurityAnalysisInput struct {
	SchemaVersion int            `json:"schema_version"`
	Tfvars        map[string]any `json:"tfvars"`
}

type SecurityAnalysisResult struct {
	SchemaVersion    int               `json:"schema_version"`
	RuleBundleID     string            `json:"rule_bundle_id"`
	RuleBundleVersion string           `json:"rule_bundle_version"`
	Findings         []SecurityFinding `json:"findings"`
}

type SecurityFindingFilter struct {
	Statuses   []SecurityFindingStatus
	Severities []SecuritySeverity
	Pagination Pagination
	Order      string
}
