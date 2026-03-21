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
			audit.WithMetadata(json.RawMessage(`{"amount":100}`)),
			audit.WithNote("VIP compensation"),
		)

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

		var parsed map[string]interface{}
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
}
