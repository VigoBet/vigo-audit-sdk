package audit_test

import (
	"strings"
	"testing"

	audit "github.com/VigoBet/vigo-audit-sdk"
)

func TestOutboxMigrationSQL(t *testing.T) {
	t.Run("returns SQL with table and index", func(t *testing.T) {
		sql := audit.OutboxMigrationSQL()

		if !strings.Contains(sql, "audit_outbox") {
			t.Error("expected SQL to contain 'audit_outbox'")
		}
		if !strings.Contains(sql, "idx_audit_outbox_pending") {
			t.Error("expected SQL to contain index name")
		}
		if !strings.Contains(sql, "IF NOT EXISTS") {
			t.Error("expected SQL to be idempotent")
		}
	})
}

func TestAuditLogMigrationSQL(t *testing.T) {
	t.Run("returns SQL with table and indexes", func(t *testing.T) {
		sql := audit.AuditLogMigrationSQL()

		if !strings.Contains(sql, "admin_audit_log") {
			t.Error("expected SQL to contain 'admin_audit_log'")
		}
		if !strings.Contains(sql, "idx_audit_log_target") {
			t.Error("expected target index")
		}
		if !strings.Contains(sql, "idx_audit_log_actor") {
			t.Error("expected actor index")
		}
	})
}

func TestHistoricalMigrationSQL(t *testing.T) {
	t.Run("returns migration SQL with action mapping", func(t *testing.T) {
		sql := audit.HistoricalMigrationSQL()

		if !strings.Contains(sql, "kh4e8g_user_activities") {
			t.Error("expected reference to legacy table")
		}
		if !strings.Contains(sql, "user.note.create") {
			t.Error("expected mapped action name")
		}
		if !strings.Contains(sql, "NOT EXISTS") {
			t.Error("expected idempotency guard")
		}
	})
}

func TestOutboxInsertSQL(t *testing.T) {
	t.Run("returns parameterized insert", func(t *testing.T) {
		sql := audit.OutboxInsertSQL()

		if !strings.Contains(sql, "INSERT INTO audit_outbox") {
			t.Error("expected INSERT statement")
		}
		if !strings.Contains(sql, "$1") {
			t.Error("expected parameterized query")
		}
	})
}
