package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type PublisherConfig struct {
	PollInterval time.Duration
	BatchSize    int
	MaxRetries   int
	Topic        string
	CleanupDays  int
}

func DefaultPublisherConfig() PublisherConfig {
	return PublisherConfig{
		PollInterval: 2 * time.Second,
		BatchSize:    100,
		MaxRetries:   10,
		Topic:        "audit.events",
		CleanupDays:  7,
	}
}

func PendingQuerySQL() string {
	return `SELECT id, payload FROM audit_outbox WHERE status = 'pending' ORDER BY created_at LIMIT $1`
}

func MarkPublishedSQL() string {
	return `UPDATE audit_outbox SET status = 'published', published_at = NOW() WHERE id = $1`
}

func IncrementRetrySQL() string {
	return `UPDATE audit_outbox SET retry_count = retry_count + 1, status = CASE WHEN retry_count + 1 >= $1 THEN 'failed' ELSE 'pending' END WHERE id = $2`
}

func MarkFailedSQL() string {
	return `UPDATE audit_outbox SET status = 'failed' WHERE id = $1`
}

func CleanupSQL(days int) string {
	return fmt.Sprintf("DELETE FROM audit_outbox WHERE status = 'published' AND published_at < NOW() - INTERVAL '%d days'", days)
}

func StartPublisher(ctx context.Context, db *sql.DB, client *kgo.Client, logger *slog.Logger, cfg PublisherConfig) {
	go func() {
		ticker := time.NewTicker(cfg.PollInterval)
		defer ticker.Stop()

		cleanupTicker := time.NewTicker(1 * time.Hour)
		defer cleanupTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				publishPending(ctx, db, client, logger, cfg)
			case <-cleanupTicker.C:
				cleanup(ctx, db, logger, cfg)
			}
		}
	}()
}

func publishPending(ctx context.Context, db *sql.DB, client *kgo.Client, logger *slog.Logger, cfg PublisherConfig) {
	rows, err := db.QueryContext(ctx, PendingQuerySQL(), cfg.BatchSize)
	if err != nil {
		logger.Error("failed to query outbox", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var payload []byte
		if err := rows.Scan(&id, &payload); err != nil {
			logger.Error("failed to scan outbox row", "error", err)
			continue
		}

		var event AuditEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			logger.Error("failed to unmarshal outbox payload, marking failed", "id", id, "error", err)
			db.ExecContext(ctx, MarkFailedSQL(), id)
			continue
		}

		record := &kgo.Record{
			Topic: cfg.Topic,
			Key:   []byte(event.PartitionKey()),
			Value: payload,
		}

		results := client.ProduceSync(ctx, record)
		if err := results.FirstErr(); err != nil {
			logger.Error("failed to publish audit event", "id", id, "error", err)
			db.ExecContext(ctx, IncrementRetrySQL(), cfg.MaxRetries, id)
			continue
		}

		if _, err := db.ExecContext(ctx, MarkPublishedSQL(), id); err != nil {
			logger.Error("failed to mark published", "id", id, "error", err)
		}
	}
}

func cleanup(ctx context.Context, db *sql.DB, logger *slog.Logger, cfg PublisherConfig) {
	result, err := db.ExecContext(ctx, CleanupSQL(cfg.CleanupDays))
	if err != nil {
		logger.Error("failed to cleanup outbox", "error", err)
		return
	}
	if rows, _ := result.RowsAffected(); rows > 0 {
		logger.Info("cleaned up published outbox records", "count", rows)
	}
}
