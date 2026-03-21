package audit

func OutboxMigrationSQL() string {
	return `CREATE TABLE IF NOT EXISTS audit_outbox (
  id UUID PRIMARY KEY,
  payload JSONB NOT NULL,
  status VARCHAR NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  published_at TIMESTAMPTZ,
  retry_count INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_audit_outbox_pending
  ON audit_outbox (created_at)
  WHERE status = 'pending';`
}

func AuditLogMigrationSQL() string {
	return `CREATE TABLE IF NOT EXISTS admin_audit_log (
  id UUID PRIMARY KEY,
  service VARCHAR NOT NULL,
  actor_type VARCHAR NOT NULL,
  actor_id VARCHAR,
  target_type VARCHAR NOT NULL,
  target_id VARCHAR NOT NULL,
  site_id VARCHAR,
  action VARCHAR NOT NULL,
  metadata JSONB,
  note VARCHAR,
  reference_id VARCHAR,
  created_at TIMESTAMPTZ NOT NULL,
  consumed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_target
  ON admin_audit_log (target_type, target_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_log_actor
  ON admin_audit_log (actor_id, created_at DESC)
  WHERE actor_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_audit_log_action
  ON admin_audit_log (action, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_log_site
  ON admin_audit_log (site_id, created_at DESC)
  WHERE site_id IS NOT NULL;`
}

func HistoricalMigrationSQL() string {
	return `INSERT INTO admin_audit_log (id, service, actor_type, actor_id, target_type, target_id, site_id, action, metadata, note, reference_id, created_at)
SELECT
  gen_random_uuid(),
  'backoffice-api',
  'admin',
  created_by::VARCHAR,
  'user',
  user_id::VARCHAR,
  NULL,
  CASE action
    WHEN 'note'                      THEN 'user.note.create'
    WHEN 'manual_balance_assignment' THEN 'balance.credit'
    WHEN 'manual_balance_remove'     THEN 'balance.debit'
    WHEN 'manual_bonus_assignment'   THEN 'bonus.grant'
    WHEN 'manual_bonus_remove'       THEN 'bonus.remove'
    WHEN 'update_user'               THEN 'user.update'
    WHEN 'add_user_flags'            THEN 'user.flags.add'
    WHEN 'update_turnover_bonus'     THEN 'turnover.adjust'
    WHEN 'update_turnover_balance'   THEN 'turnover.adjust'
    WHEN 'force_game_withdraw'       THEN 'game.force_withdraw'
    ELSE action
  END,
  metadata,
  note,
  NULL,
  created_at
FROM kh4e8g_user_activities
WHERE NOT EXISTS (SELECT 1 FROM admin_audit_log LIMIT 1);`
}
