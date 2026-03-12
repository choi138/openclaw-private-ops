package dto

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

type securityAnalysisRequest struct {
	SchemaVersion int            `json:"schema_version"`
	Tfvars        map[string]any `json:"tfvars"`
}

func DecodeSecurityAnalysis(r *http.Request, maxBytes int64) (domain.SecurityAnalysisInput, error) {
	var payload securityAnalysisRequest
	if err := decodeJSON(r, maxBytes, &payload); err != nil {
		return domain.SecurityAnalysisInput{}, err
	}

	messages := make([]string, 0)
	if payload.SchemaVersion != domain.SupportedSecurityAnalysisSchemaVersion {
		messages = append(messages, fmt.Sprintf("schema_version must be %d", domain.SupportedSecurityAnalysisSchemaVersion))
	}
	if payload.Tfvars == nil || len(payload.Tfvars) == 0 {
		messages = append(messages, "tfvars must contain at least one field")
	}
	if len(messages) > 0 {
		return domain.SecurityAnalysisInput{}, ValidationError{Messages: messages}
	}

	return domain.SecurityAnalysisInput{
		SchemaVersion: payload.SchemaVersion,
		Tfvars:        payload.Tfvars,
	}, nil
}

func ParseSecurityFindingFilter(r *http.Request) (domain.SecurityFindingFilter, error) {
	filter := domain.SecurityFindingFilter{
		Pagination: parsePaginationQuery(r),
		Order:      "desc",
	}

	if raw := strings.TrimSpace(r.URL.Query().Get("order")); raw != "" {
		if raw != "asc" && raw != "desc" {
			return domain.SecurityFindingFilter{}, ValidationError{Messages: []string{"order must be one of: asc,desc"}}
		}
		filter.Order = raw
	}

	statuses := splitCSV(r.URL.Query().Get("status"))
	for _, status := range statuses {
		if !domain.IsAllowedSecurityFindingStatus(status) {
			return domain.SecurityFindingFilter{}, ValidationError{Messages: []string{"status must be one of: open,acknowledged,resolved"}}
		}
		filter.Statuses = append(filter.Statuses, domain.SecurityFindingStatus(status))
	}

	severities := splitCSV(r.URL.Query().Get("severity"))
	for _, severity := range severities {
		if !domain.IsAllowedSecuritySeverity(severity) {
			return domain.SecurityFindingFilter{}, ValidationError{Messages: []string{"severity must be one of: critical,high,medium,info"}}
		}
		filter.Severities = append(filter.Severities, domain.SecuritySeverity(severity))
	}

	return filter, nil
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func parsePaginationQuery(r *http.Request) domain.Pagination {
	page := parsePaginationIntOrDefault(r.URL.Query().Get("page"), 1)
	pageSize := parsePaginationIntOrDefault(r.URL.Query().Get("page_size"), 50)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	return domain.Pagination{Page: page, PageSize: pageSize}
}

func parsePaginationIntOrDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return fallback
	}
	return value
}
