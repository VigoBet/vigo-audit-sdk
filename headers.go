package audit

import "net/http"

const (
	HeaderActorType = "X-Audit-Actor-Type"
	HeaderActorID   = "X-Audit-Actor-ID"
	HeaderSiteID    = "X-Audit-Site-ID"
)

func InjectHeaders(r *http.Request, actorType, actorID, siteID string) {
	r.Header.Set(HeaderActorType, actorType)
	if actorID != "" {
		r.Header.Set(HeaderActorID, actorID)
	}
	if siteID != "" {
		r.Header.Set(HeaderSiteID, siteID)
	}
}

func FromHeaders(r *http.Request) (actorType, actorID, siteID string) {
	actorType = r.Header.Get(HeaderActorType)
	if actorType == "" {
		actorType = "system"
	}
	actorID = r.Header.Get(HeaderActorID)
	siteID = r.Header.Get(HeaderSiteID)
	return
}
