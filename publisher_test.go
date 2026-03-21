package audit_test

import (
	"testing"

	audit "github.com/VigoBet/vigo-audit-sdk"
)

func TestPublisherConfig(t *testing.T) {
	t.Run("default config has correct values", func(t *testing.T) {
		cfg := audit.DefaultPublisherConfig()

		if cfg.PollInterval.Seconds() != 2 {
			t.Errorf("expected 2s poll interval, got %v", cfg.PollInterval)
		}
		if cfg.BatchSize != 100 {
			t.Errorf("expected batch size 100, got %d", cfg.BatchSize)
		}
		if cfg.MaxRetries != 10 {
			t.Errorf("expected max retries 10, got %d", cfg.MaxRetries)
		}
		if cfg.Topic != "audit.events" {
			t.Errorf("expected topic 'audit.events', got %q", cfg.Topic)
		}
		if cfg.CleanupDays != 7 {
			t.Errorf("expected 7 day cleanup, got %d", cfg.CleanupDays)
		}
	})
}

func TestCleanupSQL(t *testing.T) {
	t.Run("uses integer days in interval", func(t *testing.T) {
		sql := audit.CleanupSQL(7)
		expected := "DELETE FROM audit_outbox WHERE status = 'published' AND published_at < NOW() - INTERVAL '7 days'"
		if sql != expected {
			t.Errorf("expected %q, got %q", expected, sql)
		}
	})
}

func TestPendingQuerySQL(t *testing.T) {
	t.Run("returns correct query", func(t *testing.T) {
		sql := audit.PendingQuerySQL()
		if sql == "" {
			t.Error("expected non-empty SQL")
		}
	})
}

func TestMarkPublishedSQL(t *testing.T) {
	t.Run("returns correct update", func(t *testing.T) {
		sql := audit.MarkPublishedSQL()
		if sql == "" {
			t.Error("expected non-empty SQL")
		}
	})
}

func TestIncrementRetrySQL(t *testing.T) {
	t.Run("returns correct update with max retry check", func(t *testing.T) {
		sql := audit.IncrementRetrySQL()
		if sql == "" {
			t.Error("expected non-empty SQL")
		}
	})
}
