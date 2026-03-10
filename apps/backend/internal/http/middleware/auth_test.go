package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/http/middleware"
)

func TestWithBearerAuthRejectsInvalidToken(t *testing.T) {
	h := middleware.WithBearerAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), "expected-token")

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/summary", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestWithBearerAuthPassesValidToken(t *testing.T) {
	h := middleware.WithBearerAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor := middleware.ActorFromContext(r.Context())
		if actor == "" {
			t.Fatalf("expected actor in context")
		}
		w.WriteHeader(http.StatusNoContent)
	}), "expected-token")

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/summary", nil)
	req.Header.Set("Authorization", "Bearer expected-token")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}
