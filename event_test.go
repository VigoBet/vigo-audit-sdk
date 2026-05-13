package audit_test

import (
	"encoding/json"
	"testing"

	audit "github.com/VigoBet/vigo-audit-sdk"
)

func TestNewEvent(t *testing.T) {
	t.Run("creates event with all fields", func(t *testing.T) {
		result := audit.NewEvent(
			"core-wallets",
			"balance.credit",
			"user",
			"456",
			audit.WithActor("admin", "123"),
			audit.WithSiteID("somchaifanclub"),
			audit.WithTenantDomain("alamak.world"),
			audit.WithMetadata(json.RawMessage(`{"amount":100}`)),
			audit.WithNote("VIP compensation"),
		)
		if result.TenantDomain == nil || *result.TenantDomain != "alamak.world" {
			t.Errorf("expected tenant_domain 'alamak.world', got %v", result.TenantDomain)
		}

		if result.ID == "" {
			t.Error("expected non-empty ID")
		}
		if result.Timestamp.IsZero() {
			t.Error("expected non-zero timestamp")
		}
		if result.Service != "core-wallets" {
			t.Errorf("expected service 'core-wallets', got %q", result.Service)
		}
		if result.ActorType != "admin" {
			t.Errorf("expected actor_type 'admin', got %q", result.ActorType)
		}
		if result.ActorID == nil || *result.ActorID != "123" {
			t.Errorf("expected actor_id '123', got %v", result.ActorID)
		}
		if result.TargetType != "user" {
			t.Errorf("expected target_type 'user', got %q", result.TargetType)
		}
		if result.TargetID != "456" {
			t.Errorf("expected target_id '456', got %q", result.TargetID)
		}
		if result.Action != "balance.credit" {
			t.Errorf("expected action 'balance.credit', got %q", result.Action)
		}
		if result.SiteID == nil || *result.SiteID != "somchaifanclub" {
			t.Errorf("expected site_id 'somchaifanclub', got %v", result.SiteID)
		}
		if result.Note == nil || *result.Note != "VIP compensation" {
			t.Errorf("expected note 'VIP compensation', got %v", result.Note)
		}
	})

	t.Run("creates system event without actor_id", func(t *testing.T) {
		result := audit.NewEvent(
			"core-wallets",
			"bonus.expire",
			"user",
			"456",
			audit.WithActor("system", ""),
		)

		if result.ActorType != "system" {
			t.Errorf("expected actor_type 'system', got %q", result.ActorType)
		}
		if result.ActorID != nil {
			t.Errorf("expected nil actor_id, got %v", result.ActorID)
		}
	})

	t.Run("defaults to system actor when no actor option", func(t *testing.T) {
		result := audit.NewEvent("svc", "act", "user", "1")
		if result.ActorType != "system" {
			t.Errorf("expected default actor_type 'system', got %q", result.ActorType)
		}
	})

	t.Run("serializes to JSON with correct field names", func(t *testing.T) {
		event := audit.NewEvent(
			"backoffice-api",
			"user.update",
			"user",
			"789",
			audit.WithActor("admin", "1"),
		)

		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var parsed map[string]any
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requiredFields := []string{"id", "service", "actor_type", "target_type", "target_id", "action", "timestamp"}
		for _, field := range requiredFields {
			if _, ok := parsed[field]; !ok {
				t.Errorf("missing required field %q in JSON output", field)
			}
		}
	})

	t.Run("partition key returns target_type:target_id", func(t *testing.T) {
		event := audit.NewEvent("svc", "act", "user", "123")
		expected := "user:123"
		result := event.PartitionKey()
		if result != expected {
			t.Errorf("expected partition key %q, got %q", expected, result)
		}
	})

	t.Run("WithTenantDomain populates the field with snake_case JSON key", func(t *testing.T) {
		event := audit.NewEvent("core-auth", "tenant.create", "tenant", "t-1",
			audit.WithTenantDomain("alamak.world"))
		if event.TenantDomain == nil || *event.TenantDomain != "alamak.world" {
			t.Fatalf("expected tenant_domain populated, got %v", event.TenantDomain)
		}
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}
		var parsed map[string]any
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if got, ok := parsed["tenant_domain"]; !ok || got != "alamak.world" {
			t.Errorf("expected tenant_domain in JSON, got %v (ok=%v)", got, ok)
		}
		if _, ok := parsed["tenantDomain"]; ok {
			t.Error("camelCase tenantDomain leaked onto the wire")
		}
	})

	t.Run("WithTenantDomain skips empty string", func(t *testing.T) {
		event := audit.NewEvent("svc", "act", "user", "1", audit.WithTenantDomain(""))
		if event.TenantDomain != nil {
			t.Errorf("expected nil tenant_domain for empty input, got %v", event.TenantDomain)
		}
	})

	t.Run("omitempty: tenant_domain absent from JSON when nil", func(t *testing.T) {
		event := audit.NewEvent("svc", "act", "user", "1")
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}
		var parsed map[string]any
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if _, ok := parsed["tenant_domain"]; ok {
			t.Error("expected tenant_domain absent (omitempty), but present")
		}
	})

	t.Run("site_id and tenant_domain coexist (transition envelope)", func(t *testing.T) {
		event := audit.NewEvent("core-auth", "user.update", "user", "u-1",
			audit.WithSiteID("somchaifanclub"),
			audit.WithTenantDomain("alamak.world"))
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}
		var parsed map[string]any
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if parsed["site_id"] != "somchaifanclub" {
			t.Errorf("expected site_id, got %v", parsed["site_id"])
		}
		if parsed["tenant_domain"] != "alamak.world" {
			t.Errorf("expected tenant_domain, got %v", parsed["tenant_domain"])
		}
	})

	t.Run("backwards compat: unmarshal payload without tenant_domain", func(t *testing.T) {
		raw := []byte(`{"id":"x","service":"s","actor_type":"system","target_type":"user","target_id":"1","action":"a","timestamp":"2026-05-12T00:00:00Z"}`)
		var event audit.AuditEvent
		if err := json.Unmarshal(raw, &event); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if event.TenantDomain != nil {
			t.Errorf("expected nil tenant_domain for legacy payload, got %v", event.TenantDomain)
		}
	})

	t.Run("wire-compat smoke: TS-style payload unmarshals cleanly", func(t *testing.T) {
		raw := []byte(`{"id":"01900000-0000-7000-8000-000000000000","service":"core-auth","actor_type":"staff","actor_id":"a-1","target_type":"user","target_id":"u-1","site_id":"somchaifanclub","tenant_domain":"alamak.world","action":"user.update","timestamp":"2026-05-12T00:00:00Z"}`)
		var event audit.AuditEvent
		if err := json.Unmarshal(raw, &event); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if event.SiteID == nil || *event.SiteID != "somchaifanclub" {
			t.Errorf("site_id roundtrip failed: %v", event.SiteID)
		}
		if event.TenantDomain == nil || *event.TenantDomain != "alamak.world" {
			t.Errorf("tenant_domain roundtrip failed: %v", event.TenantDomain)
		}
	})
}
