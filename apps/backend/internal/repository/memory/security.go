package memory

import (
	"context"
	"sort"
	"time"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

func (s *Store) UpsertSecurityFindings(_ context.Context, findings []domain.SecurityFinding) ([]domain.SecurityFinding, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	persisted := make([]domain.SecurityFinding, 0, len(findings))
	for _, finding := range findings {
		if id, ok := s.securityFindingByFingerprint[finding.Fingerprint]; ok {
			for i := range s.securityFindings {
				if s.securityFindings[i].ID != id {
					continue
				}
				existing := s.securityFindings[i]
				existing.RuleID = finding.RuleID
				existing.RuleVersion = finding.RuleVersion
				existing.Severity = finding.Severity
				existing.Title = finding.Title
				existing.Description = finding.Description
				existing.FieldPath = finding.FieldPath
				existing.FixHint = finding.FixHint
				existing.Metadata = cloneMap(finding.Metadata)
				existing.LastDetectedAt = finding.LastDetectedAt
				existing.UpdatedAt = now
				if existing.Status == domain.SecurityFindingStatusResolved {
					existing.Status = domain.SecurityFindingStatusOpen
					existing.ResolvedAt = nil
				}
				s.securityFindings[i] = existing
				persisted = append(persisted, cloneSecurityFinding(existing))
				break
			}
			continue
		}

		finding.ID = s.nextSecurityFindingID
		finding.Metadata = cloneMap(finding.Metadata)
		if finding.FirstDetectedAt.IsZero() {
			finding.FirstDetectedAt = now
		}
		if finding.LastDetectedAt.IsZero() {
			finding.LastDetectedAt = finding.FirstDetectedAt
		}
		if finding.UpdatedAt.IsZero() {
			finding.UpdatedAt = now
		}
		s.nextSecurityFindingID++
		s.securityFindings = append(s.securityFindings, finding)
		s.securityFindingByFingerprint[finding.Fingerprint] = finding.ID
		persisted = append(persisted, cloneSecurityFinding(finding))
	}

	return persisted, nil
}

func (s *Store) ListSecurityFindings(_ context.Context, filter domain.SecurityFindingFilter) ([]domain.SecurityFinding, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]domain.SecurityFinding, 0, len(s.securityFindings))
	for _, finding := range s.securityFindings {
		if len(filter.Statuses) > 0 && !containsFindingStatus(filter.Statuses, finding.Status) {
			continue
		}
		if len(filter.Severities) > 0 && !containsFindingSeverity(filter.Severities, finding.Severity) {
			continue
		}
		filtered = append(filtered, cloneSecurityFinding(finding))
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filter.Order == "asc" {
			if filtered[i].LastDetectedAt.Equal(filtered[j].LastDetectedAt) {
				return filtered[i].ID < filtered[j].ID
			}
			return filtered[i].LastDetectedAt.Before(filtered[j].LastDetectedAt)
		}
		if filtered[i].LastDetectedAt.Equal(filtered[j].LastDetectedAt) {
			return filtered[i].ID > filtered[j].ID
		}
		return filtered[i].LastDetectedAt.After(filtered[j].LastDetectedAt)
	})

	return paginate(filtered, filter.Pagination), nil
}

func containsFindingStatus(allowed []domain.SecurityFindingStatus, status domain.SecurityFindingStatus) bool {
	for _, candidate := range allowed {
		if candidate == status {
			return true
		}
	}
	return false
}

func containsFindingSeverity(allowed []domain.SecuritySeverity, severity domain.SecuritySeverity) bool {
	for _, candidate := range allowed {
		if candidate == severity {
			return true
		}
	}
	return false
}

func cloneSecurityFinding(finding domain.SecurityFinding) domain.SecurityFinding {
	finding.Metadata = cloneMap(finding.Metadata)
	if finding.ResolvedAt != nil {
		resolvedAt := finding.ResolvedAt.UTC()
		finding.ResolvedAt = &resolvedAt
	}
	return finding
}
