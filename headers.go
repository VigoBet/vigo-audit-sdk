package audit

import "net/http"

const (
	HeaderActorType    = "X-Audit-Actor-Type"
	HeaderActorID      = "X-Audit-Actor-ID"
	HeaderSiteID       = "X-Audit-Site-ID"
	HeaderTenantDomain = "X-Audit-Tenant-Domain"
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

// InjectHeadersWithTenant is the v2 (additive) variant of InjectHeaders that
// also propagates the tenant domain. Kept separate so existing call sites
// stay source-compatible during the v1->v2 envelope migration.
func InjectHeadersWithTenant(r *http.Request, actorType, actorID, siteID, tenantDomain string) {
	InjectHeaders(r, actorType, actorID, siteID)
	if tenantDomain != "" {
		r.Header.Set(HeaderTenantDomain, tenantDomain)
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

// FromHeadersWithTenant returns the same context as FromHeaders plus the
// optional X-Audit-Tenant-Domain value. Empty string when absent.
func FromHeadersWithTenant(r *http.Request) (actorType, actorID, siteID, tenantDomain string) {
	actorType, actorID, siteID = FromHeaders(r)
	tenantDomain = r.Header.Get(HeaderTenantDomain)
	return
}
