package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/config"
	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/http/handler"
	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/http/middleware"
)

type AuditWriter interface {
	middleware.ReadAuditWriter
	handler.AuditWriter
}

type Dependencies struct {
	Readiness    handler.ReadinessChecker
	Dashboard    handler.DashboardReader
	Conversation handler.ConversationReader
	Infra        handler.InfraReader
	Security     handler.SecurityAnalyzer
	Ingest       handler.IngestWriter
	Audit        AuditWriter
}

func NewRouter(cfg config.Config, deps Dependencies, logger *slog.Logger) http.Handler {
	api := handler.New(handler.Dependencies{
		Readiness:            deps.Readiness,
		Dashboard:            deps.Dashboard,
		Conversation:         deps.Conversation,
		Infra:                deps.Infra,
		Security:             deps.Security,
		Ingest:               deps.Ingest,
		Audit:                deps.Audit,
		IngestMaxBodyBytes:   cfg.IngestMaxBodyBytes,
		SecurityMaxBodyBytes: cfg.SecurityMaxBodyBytes,
	}, logger)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", api.Healthz)
	mux.HandleFunc("GET /v1/readyz", api.Readyz)

	registerProtected := func(pattern string, h http.HandlerFunc) {
		var wrapped http.Handler = h
		wrapped = middleware.WithReadAudit(wrapped, deps.Audit, logger)
		wrapped = middleware.WithBearerAuth(wrapped, cfg.AdminToken, "admin")
		mux.Handle(pattern, wrapped)
	}

	registerIngest := func(pattern string, h http.HandlerFunc) {
		var wrapped http.Handler = h
		wrapped = middleware.WithBearerAuth(wrapped, cfg.IngestToken, "ingest")
		mux.Handle(pattern, wrapped)
	}

	registerProtected("GET /v1/dashboard/summary", api.DashboardSummary)
	registerProtected("GET /v1/dashboard/timeseries", api.DashboardTimeseries)
	registerProtected("GET /v1/conversations", api.ConversationsList)
	registerProtected("GET /v1/conversations/{id}", api.ConversationDetail)
	registerProtected("GET /v1/conversations/{id}/messages", api.ConversationMessages)
	registerProtected("GET /v1/conversations/{id}/attempts", api.ConversationAttempts)
	registerProtected("GET /v1/infra/status", api.InfraStatus)
	registerProtected("GET /v1/infra/snapshots", api.InfraSnapshots)
	registerProtected("GET /v1/ingest/status", api.IngestStatus)
	registerProtected("GET /v1/security/findings", api.SecurityFindings)
	registerProtected("POST /v1/security/analyze-tfvars", api.AnalyzeSecurityTfvars)
	registerIngest("POST /v1/ingest/conversation-events", api.IngestConversationEvents)
	registerIngest("POST /v1/ingest/infra-snapshot", api.IngestInfraSnapshot)
	registerIngest("POST /v1/ingest/request-attempt", api.IngestRequestAttempt)

	return mux
}
