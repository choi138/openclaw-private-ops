package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

const (
	RuleBundleID      = "ops-api-security-analysis"
	RuleBundleVersion = "2026.03.11"
)

type Repository interface {
	UpsertSecurityFindings(ctx context.Context, findings []domain.SecurityFinding) ([]domain.SecurityFinding, error)
	ListSecurityFindings(ctx context.Context, filter domain.SecurityFindingFilter) ([]domain.SecurityFinding, error)
}

type Rule interface {
	ID() string
	Version() string
	Evaluate(tfvars map[string]any) []Match
}

type Match struct {
	RuleID      string
	RuleVersion string
	Severity    domain.SecuritySeverity
	Title       string
	Description string
	FieldPath   string
	MatchKey    string
	FixHint     string
	Metadata    map[string]any
}

type Service struct {
	repo  Repository
	rules []Rule
	now   func() time.Time
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("security.NewService requires Repository")
	}

	return &Service{
		repo:  repo,
		rules: defaultRules(),
		now:   func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Analyze(ctx context.Context, input domain.SecurityAnalysisInput) (domain.SecurityAnalysisResult, error) {
	findings := make([]domain.SecurityFinding, 0)
	now := s.now()

	for _, rule := range s.rules {
		for _, match := range rule.Evaluate(input.Tfvars) {
			finding := domain.SecurityFinding{
				Fingerprint:     computeFingerprint(match),
				RuleID:          match.RuleID,
				RuleVersion:     match.RuleVersion,
				Severity:        match.Severity,
				Status:          domain.SecurityFindingStatusOpen,
				Title:           match.Title,
				Description:     match.Description,
				FieldPath:       match.FieldPath,
				FixHint:         match.FixHint,
				Metadata:        cloneMetadata(match.Metadata),
				FirstDetectedAt: now,
				LastDetectedAt:  now,
				UpdatedAt:       now,
			}
			findings = append(findings, finding)
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			return severityRank(findings[i].Severity) < severityRank(findings[j].Severity)
		}
		if findings[i].RuleID != findings[j].RuleID {
			return findings[i].RuleID < findings[j].RuleID
		}
		if findings[i].FieldPath != findings[j].FieldPath {
			return findings[i].FieldPath < findings[j].FieldPath
		}
		return findings[i].Fingerprint < findings[j].Fingerprint
	})

	persisted, err := s.repo.UpsertSecurityFindings(ctx, findings)
	if err != nil {
		return domain.SecurityAnalysisResult{}, err
	}

	return domain.SecurityAnalysisResult{
		SchemaVersion:     domain.SupportedSecurityAnalysisSchemaVersion,
		RuleBundleID:      RuleBundleID,
		RuleBundleVersion: RuleBundleVersion,
		Findings:          persisted,
	}, nil
}

func (s *Service) ListFindings(ctx context.Context, filter domain.SecurityFindingFilter) ([]domain.SecurityFinding, error) {
	return s.repo.ListSecurityFindings(ctx, filter)
}

func severityRank(severity domain.SecuritySeverity) int {
	switch severity {
	case domain.SecuritySeverityCritical:
		return 0
	case domain.SecuritySeverityHigh:
		return 1
	case domain.SecuritySeverityMedium:
		return 2
	default:
		return 3
	}
}

func computeFingerprint(match Match) string {
	payload := map[string]any{
		"rule_id":    match.RuleID,
		"field_path": match.FieldPath,
		"match_key":  stableMatchKey(match),
	}
	encoded, _ := json.Marshal(payload)
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

func stableMatchKey(match Match) string {
	if trimmed := strings.TrimSpace(match.MatchKey); trimmed != "" {
		return trimmed
	}
	if trimmed := strings.TrimSpace(match.FieldPath); trimmed != "" {
		return trimmed
	}
	return strings.TrimSpace(match.RuleID)
}

func cloneMetadata(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		switch value := v.(type) {
		case map[string]any:
			out[k] = cloneMetadata(value)
		case []string:
			out[k] = append([]string(nil), value...)
		case []any:
			copied := make([]any, len(value))
			copy(copied, value)
			out[k] = copied
		default:
			out[k] = value
		}
	}
	return out
}

func trimStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
