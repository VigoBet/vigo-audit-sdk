package audit

import (
	"context"
	"database/sql"
	"encoding/json"
)

func OutboxInsertSQL() string {
	return `INSERT INTO audit_outbox (id, payload) VALUES ($1, $2)`
}

func RecordWithSQL(ctx context.Context, tx *sql.Tx, event AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, OutboxInsertSQL(), event.ID, payload)
	return err
}

func RecordWithSQLDB(ctx context.Context, db *sql.DB, event AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, OutboxInsertSQL(), event.ID, payload)
	return err
}
