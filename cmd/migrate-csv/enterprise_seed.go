package main

import "github.com/wakeward/github-app-permissions-graph/pkg/model"

// seedUndocumentedEnterprisePermissions returns the Enterprise permissions
// visible in the live GitHub App creation UI ("Enterprise permissions"
// section, marked "Preview") as of a 2026-07-27 observation, that are NOT
// already present in the source CSV. These have no published REST
// endpoints and no confirmed manifest api_key - the api_key values below are
// placeholders derived from the UI label (snake_case), not confirmed
// against any GitHub App manifest or endpoint. All are marked
// doc_status: undocumented_preview and needs_investigation: true.
//
// docs/quarterly-review-checklist.md item 6 (live App-creation UI diff) is
// how this list gets kept current - re-open the same screen and diff
// against this file each quarter.
//
// Four Enterprise permissions already existed in the source CSV
// (enterprise_custom_properties, enterprise_organization_installation_repositories,
// enterprise_organization_installations, enterprise_scim) and are handled by
// the normal CSV row parsing instead - they're deliberately not duplicated
// here.
func seedUndocumentedEnterprisePermissions() []model.Permission {
	const provenance = " (No public API key confirmed. Seeded from the GitHub App creation UI's " +
		"\"Enterprise permissions\" screen, observed 2026-07-27. api_key is a placeholder derived " +
		"from the UI label pending real endpoint/manifest confirmation.)"

	entries := []struct {
		name, apiKey, overview string
		severity               model.Severity
		notes                  string
	}{
		{
			name:     "Copilot usage records",
			apiKey:   "enterprise_copilot_usage_records",
			overview: "View enterprise Copilot API usage records." + provenance,
			severity: model.SeverityLow,
			notes:    "Information disclosure only (usage stats, not code content) - but at enterprise-wide scope.",
		},
		{
			name:     "Custom enterprise roles",
			apiKey:   "enterprise_custom_roles",
			overview: "Manage enterprise custom roles and assignments." + provenance,
			severity: model.SeverityHigh,
			notes:    "Privilege escalation vector - can grant elevated roles across the whole enterprise, not just one org.",
		},
		{
			name:     "Enterprise AI controls",
			apiKey:   "enterprise_ai_controls",
			overview: "Manage enterprise-wide AI controls configuration." + provenance,
			severity: model.SeverityMedium,
			notes:    "Could weaken AI safety/data-handling controls (e.g. Copilot data retention, model access) enterprise-wide.",
		},
		{
			name:     "Enterprise Copilot metrics",
			apiKey:   "enterprise_copilot_metrics",
			overview: "View enterprise Copilot metrics." + provenance,
			severity: model.SeverityLow,
			notes:    "Information disclosure only, enterprise-wide scope - distinct from the org-scoped organization_copilot_metrics.",
		},
		{
			name:     "Enterprise credentials",
			apiKey:   "enterprise_credentials",
			overview: "View and manage enterprise credentials." + provenance,
			severity: model.SeverityCritical,
			notes:    "UI label is ambiguous about scope, but \"manage\" + \"credentials\" + enterprise-wide reach implies direct identity/credential exposure or manipulation - treat as Critical until the real endpoint clarifies scope.",
		},
		{
			name:     "Enterprise custom organization roles",
			apiKey:   "enterprise_custom_organization_roles",
			overview: "Create, edit, delete and list custom organization roles at the enterprise level. View system organization roles." + provenance,
			severity: model.SeverityHigh,
			notes:    "Privilege escalation across every organization in the enterprise, not just one.",
		},
		{
			name:     "Enterprise custom properties for organizations",
			apiKey:   "enterprise_organization_custom_properties",
			overview: "View organization custom properties and administer definitions at the enterprise level." + provenance,
			severity: model.SeverityMedium,
			notes:    "Configuration manipulation, enterprise-wide - distinct from the repo-scoped enterprise_custom_properties already in the CSV.",
		},
		{
			name:     "Enterprise innersource vulnerabilities",
			apiKey:   "enterprise_innersource_vulnerabilities",
			overview: "View and manage innersource vulnerabilities across the enterprise." + provenance,
			severity: model.SeverityMedium,
			notes:    "Same class of risk as repository-level vulnerability_alerts (suppress/hide findings) but at enterprise-wide blast radius.",
		},
		{
			name:     "Enterprise organizations",
			apiKey:   "enterprise_organizations",
			overview: "Create and remove enterprise organizations." + provenance,
			severity: model.SeverityCritical,
			notes:    "Can create a rogue organization as an exfiltration/staging area, or delete a legitimate one - highly destructive and hard to fully undo.",
		},
		{
			name:     "Enterprise people",
			apiKey:   "enterprise_people",
			overview: "Manage user access to the enterprise." + provenance,
			severity: model.SeverityCritical,
			notes:    "Identity/access control at the highest level - add or remove any user's enterprise-wide access.",
		},
		{
			name:     "Enterprise single sign-on",
			apiKey:   "enterprise_single_sign_on",
			overview: "View and manage enterprise single sign-on configuration." + provenance,
			severity: model.SeverityCritical,
			notes:    "Could disable or reconfigure SSO enforcement, enabling account takeover or bypass of the enterprise's primary identity control.",
		},
		{
			name:     "Enterprise teams",
			apiKey:   "enterprise_teams",
			overview: "Create, edit, remove and view enterprise teams." + provenance,
			severity: model.SeverityHigh,
			notes:    "Could restructure team-based access grants across the enterprise to establish persistence or remove legitimate access.",
		},
	}

	permissions := make([]model.Permission, 0, len(entries))
	for _, e := range entries {
		plane, ok := classifyImpactPlane(e.apiKey)
		if !ok {
			// Programmer error: every seeded key must have a classification
			// alongside it in classify.go.
			panic("migrate-csv: enterprise seed entry " + e.apiKey + " has no impact_plane classification")
		}
		permissions = append(permissions, model.Permission{
			Name:     e.name,
			APIKey:   e.apiKey,
			Category: "enterprise",
			Overview: e.overview,
			AccessLevels: []model.AccessLevelDetail{
				// Access levels aren't confirmed either (the screenshot only
				// shows the current "No access" selection, not the dropdown
				// options) - read+write is the typical default for
				// admin-style permissions, flagged for quarterly re-check.
				{Access: model.AccessRead, Severity: severityForRead(e.severity), SecurityNotes: e.notes + " (Read-only access level and its severity are inferred, not confirmed.)"},
				{Access: model.AccessWrite, Severity: e.severity, SecurityNotes: e.notes},
			},
			NeedsInvestigation: true,
			DocStatus:          model.DocStatusUndocumentedPreview,
			ImpactPlane:        plane,
		})
	}
	return permissions
}

// severityForRead derives a plausible read-only severity from a write
// severity for the enterprise seed entries, since only the write-level risk
// was reasoned about above. This is a one-step-down heuristic, not a
// researched judgment - it exists purely so the read entry isn't left blank,
// and is exactly the kind of row the quarterly review's needs_investigation
// pass should revisit.
func severityForRead(write model.Severity) model.Severity {
	switch write {
	case model.SeverityCritical:
		return model.SeverityHigh
	case model.SeverityHigh:
		return model.SeverityMedium
	default:
		return model.SeverityLow
	}
}
