package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuditEvent struct {
	ID         string  `json:"id"`
	Service    string  `json:"service"`
	ActorType  string  `json:"actor_type"`
	ActorID    *string `json:"actor_id,omitempty"`
	TargetType string  `json:"target_type"`
	TargetID   string  `json:"target_id"`
	SiteID     *string `json:"site_id,omitempty"`
	// TenantDomain rides alongside SiteID during the v1->v2 envelope
	// migration. See specs/2026-05-12-site-id-to-tenant-id-migration.md §9.
	// Column add lands in D3 3.5; no consumer reads this field yet.
	TenantDomain *string         `json:"tenant_domain,omitempty"`
	Action       string          `json:"action"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Note         *string         `json:"note,omitempty"`
	ReferenceID  *string         `json:"reference_id,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
}

type EventOption func(*AuditEvent)

func NewEvent(service, action, targetType, targetID string, opts ...EventOption) AuditEvent {
	e := AuditEvent{
		ID:         uuid.Must(uuid.NewV7()).String(),
		Service:    service,
		ActorType:  "system",
		TargetType: targetType,
		TargetID:   targetID,
		Action:     action,
		Timestamp:  time.Now().UTC(),
	}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

func WithActor(actorType, actorID string) EventOption {
	return func(e *AuditEvent) {
		e.ActorType = actorType
		if actorID != "" {
			e.ActorID = &actorID
		}
	}
}

func WithSiteID(siteID string) EventOption {
	return func(e *AuditEvent) {
		if siteID != "" {
			e.SiteID = &siteID
		}
	}
}

func WithTenantDomain(tenantDomain string) EventOption {
	return func(e *AuditEvent) {
		if tenantDomain != "" {
			e.TenantDomain = &tenantDomain
		}
	}
}

func WithMetadata(metadata json.RawMessage) EventOption {
	return func(e *AuditEvent) {
		e.Metadata = metadata
	}
}

func WithNote(note string) EventOption {
	return func(e *AuditEvent) {
		if note != "" {
			e.Note = &note
		}
	}
}

func WithReferenceID(refID string) EventOption {
	return func(e *AuditEvent) {
		if refID != "" {
			e.ReferenceID = &refID
		}
	}
}

func (e AuditEvent) PartitionKey() string {
	return fmt.Sprintf("%s:%s", e.TargetType, e.TargetID)
}
