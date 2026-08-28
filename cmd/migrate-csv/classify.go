package main

import "github.com/wakeward/gh-app-graph/pkg/model"

// impactPlaneByAPIKey is a first-pass, hand-reviewed classification of every
// permission's impact plane (see the threat-oriented breakdown in
// docs/methodology.md once written). This is deliberately a literal,
// reviewable map rather than a keyword heuristic: each entry is a specific
// judgment call, not a pattern match, and is expected to be corrected in
// place (in the generated YAML, not here) as understanding improves -
// this map only runs once, at migration time.
//
// Rule of thumb used throughout: Control = configures policy/settings
// (branch protection, RBAC, network config); Data/Execution = touches code,
// CI execution, artifacts, or secrets used at runtime; Governance/Identity =
// who has access, membership, credentials-as-identity, and auditability.
var impactPlaneByAPIKey = map[string]model.ImpactPlane{
	// Account
	"blocking":                    model.ImpactPlaneGovernanceIdentity,
	"codespaces":                  model.ImpactPlaneDataExecution,
	"codespaces_lifecycle_admin":  model.ImpactPlaneDataExecution,
	"codespaces_secrets":          model.ImpactPlaneDataExecution,
	"codespaces_user_secrets":     model.ImpactPlaneDataExecution,
	"copilot_messages":            model.ImpactPlaneDataExecution,
	"emails":                      model.ImpactPlaneGovernanceIdentity,
	"followers":                   model.ImpactPlaneGovernanceIdentity,
	"gists":                       model.ImpactPlaneDataExecution,
	"git_signing_ssh_public_keys": model.ImpactPlaneDataExecution,
	"gpg_keys":                    model.ImpactPlaneDataExecution,
	"interaction_limits":          model.ImpactPlaneControl,
	"keys":                        model.ImpactPlaneDataExecution,
	"notifications":               model.ImpactPlaneGovernanceIdentity,
	"plan":                        model.ImpactPlaneGovernanceIdentity,
	"profile":                     model.ImpactPlaneGovernanceIdentity,
	"starring":                    model.ImpactPlaneGovernanceIdentity,
	"user_events":                 model.ImpactPlaneGovernanceIdentity,
	"watching":                    model.ImpactPlaneGovernanceIdentity,

	// Enterprise (existing, documented to some degree)
	"enterprise_custom_properties":                      model.ImpactPlaneControl,
	"enterprise_organization_installation_repositories": model.ImpactPlaneControl,
	"enterprise_organization_installations":             model.ImpactPlaneControl,
	"enterprise_scim":                                   model.ImpactPlaneGovernanceIdentity,

	// Organization
	"members":                                     model.ImpactPlaneGovernanceIdentity,
	"organization_actions_variables":              model.ImpactPlaneDataExecution,
	"organization_administration":                 model.ImpactPlaneControl,
	"organization_agent_secrets":                  model.ImpactPlaneDataExecution,
	"organization_agent_variables":                model.ImpactPlaneDataExecution,
	"organization_announcement_banners":           model.ImpactPlaneControl,
	"organization_api_insights":                   model.ImpactPlaneGovernanceIdentity,
	"organization_campaigns":                      model.ImpactPlaneGovernanceIdentity,
	"organization_codespaces":                     model.ImpactPlaneDataExecution,
	"organization_codespaces_lifecycle_admin":     model.ImpactPlaneDataExecution,
	"organization_codespaces_secrets":             model.ImpactPlaneDataExecution,
	"organization_codespaces_settings":            model.ImpactPlaneControl,
	"organization_copilot_agent_settings":         model.ImpactPlaneDataExecution,
	"organization_copilot_content_exclusion":      model.ImpactPlaneDataExecution,
	"organization_copilot_metrics":                model.ImpactPlaneGovernanceIdentity,
	"organization_copilot_seat_management":        model.ImpactPlaneGovernanceIdentity,
	"organization_copilot_spaces":                 model.ImpactPlaneDataExecution,
	"organization_custom_org_roles":               model.ImpactPlaneControl,
	"organization_custom_properties":              model.ImpactPlaneControl,
	"organization_custom_roles":                   model.ImpactPlaneControl,
	"organization_dependabot_secrets":             model.ImpactPlaneDataExecution,
	"organization_dependabot_variables":           model.ImpactPlaneDataExecution,
	"organization_events":                         model.ImpactPlaneGovernanceIdentity,
	"organization_hooks":                          model.ImpactPlaneControl,
	"organization_knowledge_bases":                model.ImpactPlaneDataExecution,
	"organization_network_configurations":         model.ImpactPlaneControl,
	"organization_packages":                       model.ImpactPlaneDataExecution,
	"organization_personal_access_token_requests": model.ImpactPlaneGovernanceIdentity,
	"organization_personal_access_tokens":         model.ImpactPlaneGovernanceIdentity,
	"organization_plan":                           model.ImpactPlaneGovernanceIdentity,
	"organization_pre_receive_hooks":              model.ImpactPlaneDataExecution,
	"organization_private_registries":             model.ImpactPlaneDataExecution,
	"organization_projects":                       model.ImpactPlaneDataExecution,
	"organization_secrets":                        model.ImpactPlaneDataExecution,
	"organization_self_hosted_runners":            model.ImpactPlaneDataExecution,
	"organization_user_blocking":                  model.ImpactPlaneGovernanceIdentity,
	"team_discussions":                            model.ImpactPlaneDataExecution,

	// Repository
	"actions":                      model.ImpactPlaneDataExecution,
	"actions_variables":            model.ImpactPlaneDataExecution,
	"administration":               model.ImpactPlaneControl,
	"agent_secrets":                model.ImpactPlaneDataExecution,
	"agent_variables":              model.ImpactPlaneDataExecution,
	"attestations":                 model.ImpactPlaneDataExecution,
	"checks":                       model.ImpactPlaneDataExecution,
	"code_quality":                 model.ImpactPlaneDataExecution,
	"codespaces_metadata":          model.ImpactPlaneDataExecution,
	"content_references":           model.ImpactPlaneDataExecution,
	"contents":                     model.ImpactPlaneDataExecution,
	"copilot_agent_settings":       model.ImpactPlaneDataExecution,
	"dependabot_secrets":           model.ImpactPlaneDataExecution,
	"deployments":                  model.ImpactPlaneDataExecution,
	"discussions":                  model.ImpactPlaneDataExecution,
	"environments":                 model.ImpactPlaneControl,
	"issues":                       model.ImpactPlaneDataExecution,
	"merge_queues":                 model.ImpactPlaneDataExecution,
	"metadata":                     model.ImpactPlaneDataExecution,
	"packages":                     model.ImpactPlaneDataExecution,
	"pages":                        model.ImpactPlaneDataExecution,
	"pull_requests":                model.ImpactPlaneDataExecution,
	"repository_advisories":        model.ImpactPlaneDataExecution,
	"repository_custom_properties": model.ImpactPlaneControl,
	"repository_hooks":             model.ImpactPlaneControl,
	"repository_pre_receive_hooks": model.ImpactPlaneDataExecution,
	"repository_projects":          model.ImpactPlaneDataExecution,
	"secret_scanning_alerts":       model.ImpactPlaneDataExecution,
	"secrets":                      model.ImpactPlaneDataExecution,
	"security_events":              model.ImpactPlaneDataExecution,
	"single_file":                  model.ImpactPlaneDataExecution,
	"statuses":                     model.ImpactPlaneDataExecution,
	"vulnerability_alerts":         model.ImpactPlaneDataExecution,
	"workflows":                    model.ImpactPlaneDataExecution,

	// Enterprise (new, seeded from the App creation UI screenshot - see
	// enterprise_seed.go)
	"enterprise_copilot_usage_records":          model.ImpactPlaneGovernanceIdentity,
	"enterprise_custom_roles":                   model.ImpactPlaneControl,
	"enterprise_ai_controls":                    model.ImpactPlaneControl,
	"enterprise_copilot_metrics":                model.ImpactPlaneGovernanceIdentity,
	"enterprise_credentials":                    model.ImpactPlaneGovernanceIdentity,
	"enterprise_custom_organization_roles":      model.ImpactPlaneControl,
	"enterprise_organization_custom_properties": model.ImpactPlaneControl,
	"enterprise_innersource_vulnerabilities":    model.ImpactPlaneDataExecution,
	"enterprise_organizations":                  model.ImpactPlaneGovernanceIdentity,
	"enterprise_people":                         model.ImpactPlaneGovernanceIdentity,
	"enterprise_single_sign_on":                 model.ImpactPlaneGovernanceIdentity,
	"enterprise_teams":                          model.ImpactPlaneGovernanceIdentity,
}

// classifyImpactPlane looks up apiKey's impact plane. ok is false if
// apiKey isn't in the map - every permission in the source CSV plus the
// enterprise seed list is expected to be classified, so callers should treat
// a miss as a data gap (flag needs_investigation) rather than silently
// defaulting.
func classifyImpactPlane(apiKey string) (plane model.ImpactPlane, ok bool) {
	plane, ok = impactPlaneByAPIKey[apiKey]
	return plane, ok
}
