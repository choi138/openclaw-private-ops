package dto

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

func TestDecodeSecurityAnalysisRejectsEmptyTfvars(t *testing.T) {
	body := `{
		"schema_version": 1,
		"tfvars": {}
	}`
	req := httptest.NewRequest("POST", "/v1/security/analyze-tfvars", strings.NewReader(body))

	_, err := DecodeSecurityAnalysis(req, 1024)
	if err == nil {
		t.Fatal("expected validation error")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validationErr.Messages[0] != "tfvars must contain at least one field" {
		t.Fatalf("unexpected validation message: %+v", validationErr.Messages)
	}
}

func TestDecodeSecurityAnalysisHonorsMaxBodyBytes(t *testing.T) {
	body := `{
		"schema_version": 1,
		"tfvars": {
			"ui_source_ranges": ["0.0.0.0/0"],
			"ssh_source_ranges": ["0.0.0.0/0"]
		}
	}`
	req := httptest.NewRequest("POST", "/v1/security/analyze-tfvars", strings.NewReader(body))

	_, err := DecodeSecurityAnalysis(req, 16)
	if err == nil {
		t.Fatal("expected body size validation error")
	}
}

func TestParseSecurityFindingFilterSupportsCSVAndPagination(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/security/findings?status=open,resolved&severity=critical,high&page=2&page_size=25&order=asc", nil)

	filter, err := ParseSecurityFindingFilter(req)
	if err != nil {
		t.Fatalf("expected filter to parse, got %v", err)
	}

	if filter.Order != "asc" {
		t.Fatalf("expected asc order, got %q", filter.Order)
	}
	if filter.Pagination.Page != 2 || filter.Pagination.PageSize != 25 {
		t.Fatalf("unexpected pagination: %+v", filter.Pagination)
	}
	if len(filter.Statuses) != 2 || filter.Statuses[0] != domain.SecurityFindingStatusOpen || filter.Statuses[1] != domain.SecurityFindingStatusResolved {
		t.Fatalf("unexpected statuses: %+v", filter.Statuses)
	}
	if len(filter.Severities) != 2 || filter.Severities[0] != domain.SecuritySeverityCritical || filter.Severities[1] != domain.SecuritySeverityHigh {
		t.Fatalf("unexpected severities: %+v", filter.Severities)
	}
}
