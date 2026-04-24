package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuditEvent struct {
	ID          string          `json:"id"`
	Service     string          `json:"service"`
	ActorType   string          `json:"actor_type"`
	ActorID     *string         `json:"actor_id,omitempty"`
	TargetType  string          `json:"target_type"`
	TargetID    string          `json:"target_id"`
	SiteID      *string         `json:"site_id,omitempty"`
	Action      string          `json:"action"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	Note        *string         `json:"note,omitempty"`
	ReferenceID *string         `json:"reference_id,omitempty"`
	Timestamp   time.Time       `json:"timestamp"`

	// Plan 0.3 extensions — optional; pre-existing producers won't set these.
	TenantID             *string         `json:"tenant_id,omitempty"`
	TargetTypeV2         *string         `json:"target_type_v2,omitempty"`
	PayloadDiff          json.RawMessage `json:"payload_diff,omitempty"`
	PermissionKey        *string         `json:"permission_key,omitempty"`
	Outcome              *string         `json:"outcome,omitempty"`
	FailureReason        *string         `json:"failure_reason,omitempty"`
	TraceID              *string         `json:"trace_id,omitempty"`
	IP                   *string         `json:"ip,omitempty"`
	ParentActorAccountID *string         `json:"parent_actor_account_id,omitempty"`
	SourceService        *string         `json:"source_service,omitempty"`
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

func WithPayloadDiff(raw json.RawMessage) EventOption {
	return func(e *AuditEvent) { e.PayloadDiff = raw }
}

func WithPermissionKey(k string) EventOption {
	return func(e *AuditEvent) {
		if k != "" {
			e.PermissionKey = &k
		}
	}
}

func WithOutcome(o string) EventOption {
	return func(e *AuditEvent) {
		if o != "" {
			e.Outcome = &o
		}
	}
}

func WithFailureReason(r string) EventOption {
	return func(e *AuditEvent) {
		if r != "" {
			e.FailureReason = &r
		}
	}
}

func WithTraceID(t string) EventOption {
	return func(e *AuditEvent) {
		if t != "" {
			e.TraceID = &t
		}
	}
}

func WithIP(ip string) EventOption {
	return func(e *AuditEvent) {
		if ip != "" {
			e.IP = &ip
		}
	}
}

func WithParentActorAccountID(id string) EventOption {
	return func(e *AuditEvent) {
		if id != "" {
			e.ParentActorAccountID = &id
		}
	}
}

func WithTenantID(id string) EventOption {
	return func(e *AuditEvent) {
		if id != "" {
			e.TenantID = &id
		}
	}
}

func WithTargetTypeV2(t string) EventOption {
	return func(e *AuditEvent) {
		if t != "" {
			e.TargetTypeV2 = &t
		}
	}
}

func WithSourceService(s string) EventOption {
	return func(e *AuditEvent) {
		if s != "" {
			e.SourceService = &s
		}
	}
}

func (e AuditEvent) PartitionKey() string {
	return fmt.Sprintf("%s:%s", e.TargetType, e.TargetID)
}
