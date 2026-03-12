package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

func (s *Store) UpsertSecurityFindings(ctx context.Context, findings []domain.SecurityFinding) ([]domain.SecurityFinding, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	const upsertQuery = `
INSERT INTO security_findings (
  fingerprint, rule_id, rule_version, severity, status, title, description, field_path, fix_hint, metadata_json,
  first_detected_at, last_detected_at, resolved_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10::jsonb, $11, $12, $13, $14)
ON CONFLICT (fingerprint) DO UPDATE
SET rule_id = EXCLUDED.rule_id,
    rule_version = EXCLUDED.rule_version,
    severity = EXCLUDED.severity,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    field_path = EXCLUDED.field_path,
    fix_hint = EXCLUDED.fix_hint,
    metadata_json = EXCLUDED.metadata_json,
    last_detected_at = EXCLUDED.last_detected_at,
    updated_at = EXCLUDED.updated_at,
    status = CASE
      WHEN security_findings.status = 'resolved' THEN 'open'
      ELSE security_findings.status
    END,
    resolved_at = CASE
      WHEN security_findings.status = 'resolved' THEN NULL
      ELSE security_findings.resolved_at
    END
RETURNING id, fingerprint, rule_id, rule_version, severity, status, title, description, COALESCE(field_path, ''), COALESCE(fix_hint, ''), metadata_json, first_detected_at, last_detected_at, resolved_at, updated_at
`

	persisted := make([]domain.SecurityFinding, 0, len(findings))
	for _, finding := range findings {
		metadata := finding.Metadata
		if metadata == nil {
			metadata = map[string]any{}
		}
		payload, err := json.Marshal(metadata)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		var (
			item       domain.SecurityFinding
			fieldPath  string
			fixHint    string
			metadataJSON []byte
			resolvedAt sql.NullTime
		)
		if err := tx.QueryRowContext(
			ctx,
			upsertQuery,
			finding.Fingerprint,
			finding.RuleID,
			finding.RuleVersion,
			string(finding.Severity),
			string(finding.Status),
			finding.Title,
			finding.Description,
			finding.FieldPath,
			finding.FixHint,
			string(payload),
			finding.FirstDetectedAt,
			finding.LastDetectedAt,
			finding.ResolvedAt,
			finding.UpdatedAt,
		).Scan(
			&item.ID,
			&item.Fingerprint,
			&item.RuleID,
			&item.RuleVersion,
			&item.Severity,
			&item.Status,
			&item.Title,
			&item.Description,
			&fieldPath,
			&fixHint,
			&metadataJSON,
			&item.FirstDetectedAt,
			&item.LastDetectedAt,
			&resolvedAt,
			&item.UpdatedAt,
		); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		item.FieldPath = fieldPath
		item.FixHint = fixHint
		if resolvedAt.Valid {
			ts := resolvedAt.Time.UTC()
			item.ResolvedAt = &ts
		}
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		persisted = append(persisted, item)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return persisted, nil
}

func (s *Store) ListSecurityFindings(ctx context.Context, filter domain.SecurityFindingFilter) ([]domain.SecurityFinding, error) {
	if filter.Pagination.Page < 1 {
		filter.Pagination.Page = 1
	}
	if filter.Pagination.PageSize <= 0 {
		filter.Pagination.PageSize = 50
	}
	offset := (filter.Pagination.Page - 1) * filter.Pagination.PageSize

	args := make([]any, 0)
	where := make([]string, 0)
	argPos := 1

	if len(filter.Statuses) > 0 {
		statuses := make([]string, 0, len(filter.Statuses))
		for _, status := range filter.Statuses {
			statuses = append(statuses, string(status))
		}
		where = append(where, fmt.Sprintf("status = ANY($%d)", argPos))
		args = append(args, pq.Array(statuses))
		argPos++
	}

	if len(filter.Severities) > 0 {
		severities := make([]string, 0, len(filter.Severities))
		for _, severity := range filter.Severities {
			severities = append(severities, string(severity))
		}
		where = append(where, fmt.Sprintf("severity = ANY($%d)", argPos))
		args = append(args, pq.Array(severities))
		argPos++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	order := "DESC"
	if filter.Order == "asc" {
		order = "ASC"
	}

	query := fmt.Sprintf(`
SELECT id, fingerprint, rule_id, rule_version, severity, status, title, description, COALESCE(field_path, ''), COALESCE(fix_hint, ''), metadata_json, first_detected_at, last_detected_at, resolved_at, updated_at
FROM security_findings
%s
ORDER BY last_detected_at %s, id %s
LIMIT $%d OFFSET $%d
`, whereClause, order, order, argPos, argPos+1)
	args = append(args, filter.Pagination.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	findings := make([]domain.SecurityFinding, 0)
	for rows.Next() {
		var (
			item       domain.SecurityFinding
			fieldPath  string
			fixHint    string
			metadataJSON []byte
			resolvedAt sql.NullTime
		)
		if err := rows.Scan(
			&item.ID,
			&item.Fingerprint,
			&item.RuleID,
			&item.RuleVersion,
			&item.Severity,
			&item.Status,
			&item.Title,
			&item.Description,
			&fieldPath,
			&fixHint,
			&metadataJSON,
			&item.FirstDetectedAt,
			&item.LastDetectedAt,
			&resolvedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.FieldPath = fieldPath
		item.FixHint = fixHint
		if resolvedAt.Valid {
			ts := resolvedAt.Time.UTC()
			item.ResolvedAt = &ts
		}
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			return nil, err
		}
		findings = append(findings, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return findings, nil
}
