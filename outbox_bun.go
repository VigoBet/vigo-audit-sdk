package audit

import (
	"context"
	"encoding/json"

	"github.com/uptrace/bun"
)

func RecordWithBun(ctx context.Context, db bun.IDB, event AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = db.NewRaw(
		"INSERT INTO audit_outbox (id, payload) VALUES (?, ?::jsonb)",
		event.ID, payload,
	).Exec(ctx)
	return err
}
