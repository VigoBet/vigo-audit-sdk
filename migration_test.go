package audit

import (
	"strings"
	"testing"
)

func TestAuditLogExtensionsSQL(t *testing.T) {
	sql := AuditLogExtensionsSQL()
	wanted := []string{
		"ALTER TABLE admin_audit_log",
		"ADD COLUMN IF NOT EXISTS target_type",
		"ADD COLUMN IF NOT EXISTS payload_diff JSONB",
		"ADD COLUMN IF NOT EXISTS permission_key",
		"ADD COLUMN IF NOT EXISTS outcome",
		"ADD COLUMN IF NOT EXISTS failure_reason",
		"ADD COLUMN IF NOT EXISTS trace_id",
		"ADD COLUMN IF NOT EXISTS ip INET",
		"ADD COLUMN IF NOT EXISTS parent_actor_account_id UUID",
	}
	for _, w := range wanted {
		if !strings.Contains(sql, w) {
			t.Errorf("expected SQL to contain %q, got:\n%s", w, sql)
		}
	}
}

func TestReadAuditMigrationSQL(t *testing.T) {
	sql := ReadAuditMigrationSQL()
	wanted := []string{
		"CREATE TABLE IF NOT EXISTS admin_read_audit",
		"id UUID PRIMARY KEY",
		"actor_account_id",
		"tenant_id",
		"endpoint",
		"target_type",
		"target_id",
		"accessed_at",
		"trace_id",
		"idx_read_audit_actor",
		"idx_read_audit_target",
	}
	for _, w := range wanted {
		if !strings.Contains(sql, w) {
			t.Errorf("expected SQL to contain %q, got:\n%s", w, sql)
		}
	}
}
