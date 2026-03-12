package memory

import (
	"context"
	"testing"
	"time"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/domain"
)

func TestInsertReadAuditCopiesMetadata(t *testing.T) {
	store := NewStore()
	event := domain.AuditEvent{
		Actor:        "admin",
		Action:       "read",
		ResourceType: "dashboard",
		Metadata: map[string]any{
			"path": "/v1/dashboard/summary",
			"nested": map[string]any{
				"status": 200,
			},
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := store.InsertReadAudit(context.Background(), event); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	event.Metadata["path"] = "/mutated"
	event.Metadata["nested"].(map[string]any)["status"] = 500

	events := store.AuditEvents()
	if got := events[len(events)-1].Metadata["path"]; got != "/v1/dashboard/summary" {
		t.Fatalf("expected stored path to remain unchanged, got %v", got)
	}
	nested := events[len(events)-1].Metadata["nested"].(map[string]any)
	if got := nested["status"]; got != 200 {
		t.Fatalf("expected stored nested status to remain unchanged, got %v", got)
	}
}

func TestAuditEventsReturnsCopies(t *testing.T) {
	store := NewStore()
	if err := store.InsertReadAudit(context.Background(), domain.AuditEvent{
		Actor:        "admin",
		Action:       "read",
		ResourceType: "dashboard",
		Metadata: map[string]any{
			"path": "/v1/dashboard/summary",
		},
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events := store.AuditEvents()
	events[0].Metadata["path"] = "/mutated"

	eventsAgain := store.AuditEvents()
	if got := eventsAgain[0].Metadata["path"]; got != "/v1/dashboard/summary" {
		t.Fatalf("expected stored path to remain unchanged, got %v", got)
	}
}

func TestEncodedKeysAvoidSeparatorCollisions(t *testing.T) {
	if first, second := ingestKey("source:a", "id"), ingestKey("source", "a:id"); first == second {
		t.Fatalf("expected external identity keys to remain distinct, got %q", first)
	}
	if first, second := ingestEventKey("conversation:event", "source", "id"), ingestEventKey("conversation", "event:source", "id"); first == second {
		t.Fatalf("expected ingest event keys to remain distinct, got %q", first)
	}
}

func TestPersistInfraSnapshotReusesEventIdentity(t *testing.T) {
	store := NewStore()
	snapshot := domain.InfraSnapshotInput{
		SchemaVersion: 1,
		Source:        "wireguard",
		EventID:       "evt-1",
		CapturedAt:    time.Date(2026, 3, 11, 8, 0, 0, 0, time.UTC),
		VPNPeerCount:  3,
		OpenClawUp:    true,
		CPUPct:        12.5,
		MemPct:        31.2,
	}

	if err := store.PersistInfraSnapshot(context.Background(), snapshot); err != nil {
		t.Fatalf("expected first snapshot persist to succeed, got %v", err)
	}

	snapshot.CPUPct = 18.0
	if err := store.PersistInfraSnapshot(context.Background(), snapshot); err != nil {
		t.Fatalf("expected repeated snapshot persist to succeed, got %v", err)
	}

	snapshots, err := store.ListSnapshots(context.Background(), time.Time{}, time.Now().UTC().Add(time.Hour), domain.Pagination{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected snapshot list to load, got %v", err)
	}
	if len(snapshots) != 2 {
		t.Fatalf("expected one seed snapshot plus one event-backed snapshot, got %d", len(snapshots))
	}
	found := false
	for _, item := range snapshots {
		if item.CapturedAt.Equal(snapshot.CapturedAt) && item.CPUPct == 18.0 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected repeated persist to update existing snapshot, got %+v", snapshots)
	}
}

func TestPersistRequestAttemptDoesNotClearConversationEndTime(t *testing.T) {
	store := NewStore()
	startedAt := time.Date(2026, 3, 11, 8, 0, 0, 0, time.UTC)
	endedAt := time.Date(2026, 3, 11, 8, 5, 0, 0, time.UTC)

	if err := store.PersistConversationEvent(context.Background(), domain.ConversationEventInput{
		SchemaVersion: 1,
		Source:        "openclaw",
		EventID:       "evt-conv",
		OccurredAt:    endedAt,
		Account: domain.AccountInput{
			ExternalID: "acct-1",
			Email:      "ops@example.com",
			Status:     "active",
		},
		Conversation: domain.ConversationInput{
			ExternalID: "conv-1",
			Channel:    "telegram",
			Status:     "completed",
			StartedAt:  startedAt,
			EndedAt:    &endedAt,
		},
	}); err != nil {
		t.Fatalf("expected conversation event to persist, got %v", err)
	}

	if err := store.PersistRequestAttempt(context.Background(), domain.RequestAttemptEventInput{
		SchemaVersion: 1,
		Source:        "openclaw",
		EventID:       "evt-attempt",
		OccurredAt:    endedAt.Add(-time.Minute),
		Account: domain.AccountInput{
			ExternalID: "acct-1",
			Email:      "ops@example.com",
			Status:     "active",
		},
		Conversation: domain.ConversationInput{
			ExternalID: "conv-1",
			Channel:    "telegram",
			Status:     "in_progress",
			StartedAt:  startedAt.Add(time.Minute),
		},
		Attempt: domain.RequestAttemptInput{
			ExternalID: "attempt-1",
			Provider:   "anthropic",
			Model:      "claude",
			TokensIn:   1,
			TokensOut:  1,
			CostUSD:    0.01,
			LatencyMS:  50,
			Success:    true,
			CreatedAt:  endedAt.Add(-time.Minute),
		},
	}); err != nil {
		t.Fatalf("expected request attempt to persist, got %v", err)
	}

	conversations, err := store.ListConversations(context.Background(), domain.ConversationFilter{}, domain.Pagination{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected conversations to load, got %v", err)
	}

	var foundConversation *domain.Conversation
	for i := range conversations {
		if conversations[i].AccountID != 1001 {
			foundConversation = &conversations[i]
			break
		}
	}
	if foundConversation == nil {
		t.Fatal("expected persisted conversation to exist")
	}
	if foundConversation.Status != "completed" {
		t.Fatalf("expected terminal status to be preserved, got %q", foundConversation.Status)
	}
	if !foundConversation.StartedAt.Equal(startedAt) {
		t.Fatalf("expected earliest started_at to be preserved, got %v", foundConversation.StartedAt)
	}
	if foundConversation.EndedAt == nil || !foundConversation.EndedAt.Equal(endedAt) {
		t.Fatalf("expected ended_at to be preserved, got %+v", foundConversation.EndedAt)
	}
}
