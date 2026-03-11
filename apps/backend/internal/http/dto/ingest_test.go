package dto

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

func TestDecodeConversationEventValidatesSchemaVersion(t *testing.T) {
	body := `{
		"schema_version": 99,
		"source": "openclaw",
		"event_id": "evt-1",
		"occurred_at": "2026-03-11T08:00:05Z",
		"account": {"external_id":"acct-1","email":"ops@example.com","status":"active"},
		"conversation": {"external_id":"conv-1","channel":"telegram","status":"completed","started_at":"2026-03-11T08:00:00Z"},
		"message": {"external_id":"msg-1","role":"user","content_masked":"hello","created_at":"2026-03-11T08:00:05Z"}
	}`
	req := httptest.NewRequest("POST", "/v1/ingest/conversation-events", strings.NewReader(body))

	_, err := DecodeConversationEvent(req, 1024)
	if err == nil {
		t.Fatal("expected schema version validation error")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validationErr.Messages[0] != "schema_version must be 1" {
		t.Fatalf("unexpected validation message: %+v", validationErr.Messages)
	}
}

func TestDecodeRequestAttemptEventRejectsNegativeMetrics(t *testing.T) {
	body := `{
		"schema_version": 1,
		"source": "openclaw",
		"event_id": "evt-2",
		"occurred_at": "2026-03-11T08:00:05Z",
		"account": {"external_id":"acct-1","email":"ops@example.com","status":"active"},
		"conversation": {"external_id":"conv-1","channel":"telegram","status":"completed","started_at":"2026-03-11T08:00:00Z"},
		"attempt": {
			"external_id":"attempt-1",
			"provider":"anthropic",
			"model":"claude",
			"tokens_in":-1,
			"tokens_out":0,
			"cost_usd":0,
			"latency_ms":5,
			"success":true,
			"created_at":"2026-03-11T08:00:05Z"
		}
	}`
	req := httptest.NewRequest("POST", "/v1/ingest/request-attempt", strings.NewReader(body))

	_, err := DecodeRequestAttemptEvent(req, 2048)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestDecodeInfraSnapshotReturnsDomainPayload(t *testing.T) {
	body := `{
		"schema_version": 1,
		"source": "wireguard",
		"event_id": "evt-3",
		"captured_at": "2026-03-11T08:00:05Z",
		"vpn_peer_count": 3,
		"openclaw_up": true,
		"cpu_pct": 12.5,
		"mem_pct": 31.2
	}`
	req := httptest.NewRequest("POST", "/v1/ingest/infra-snapshot", strings.NewReader(body))

	payload, err := DecodeInfraSnapshot(req, 2048)
	if err != nil {
		t.Fatalf("expected payload to decode, got %v", err)
	}
	if payload.SchemaVersion != domain.SupportedIngestSchemaVersion {
		t.Fatalf("unexpected schema version %d", payload.SchemaVersion)
	}
	if payload.Source != "wireguard" || payload.EventID != "evt-3" {
		t.Fatalf("unexpected payload identity: %+v", payload)
	}
}
