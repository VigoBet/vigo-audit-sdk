package audit_test

import (
	"net/http/httptest"
	"testing"

	audit "github.com/VigoBet/vigo-audit-sdk"
)

func TestInjectHeaders(t *testing.T) {
	t.Run("injects all audit headers", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		audit.InjectHeaders(req, "admin", "123", "somchaifanclub")

		if req.Header.Get("X-Audit-Actor-Type") != "admin" {
			t.Error("missing X-Audit-Actor-Type")
		}
		if req.Header.Get("X-Audit-Actor-ID") != "123" {
			t.Error("missing X-Audit-Actor-ID")
		}
		if req.Header.Get("X-Audit-Site-ID") != "somchaifanclub" {
			t.Error("missing X-Audit-Site-ID")
		}
	})

	t.Run("skips empty values", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		audit.InjectHeaders(req, "system", "", "")

		if req.Header.Get("X-Audit-Actor-Type") != "system" {
			t.Error("missing X-Audit-Actor-Type")
		}
		if req.Header.Get("X-Audit-Actor-ID") != "" {
			t.Error("expected empty X-Audit-Actor-ID")
		}
	})
}

func TestFromHeaders(t *testing.T) {
	t.Run("extracts audit context from headers", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("X-Audit-Actor-Type", "admin")
		req.Header.Set("X-Audit-Actor-ID", "123")
		req.Header.Set("X-Audit-Site-ID", "somchaifanclub")

		actorType, actorID, siteID := audit.FromHeaders(req)

		if actorType != "admin" {
			t.Errorf("expected 'admin', got %q", actorType)
		}
		if actorID != "123" {
			t.Errorf("expected '123', got %q", actorID)
		}
		if siteID != "somchaifanclub" {
			t.Errorf("expected 'somchaifanclub', got %q", siteID)
		}
	})

	t.Run("returns system default when headers missing", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		actorType, actorID, siteID := audit.FromHeaders(req)

		if actorType != "system" {
			t.Errorf("expected default 'system', got %q", actorType)
		}
		if actorID != "" {
			t.Errorf("expected empty, got %q", actorID)
		}
		if siteID != "" {
			t.Errorf("expected empty, got %q", siteID)
		}
	})
}
