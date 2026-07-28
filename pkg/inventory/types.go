package inventory

import "time"

// GitHubDocsAPIVersion is the GitHub REST API docs version this inventory is
// reconciled against (recorded in generated output so a version bump is a
// distinct diff from endpoint churn). Update deliberately when GitHub ships a
// new apiVersion and we re-verify scrapers against it.
const GitHubDocsAPIVersion = "2026-03-10"

// DefaultSourceURL is the octokit/app-permissions generated inventory we
// consume as the base endpoint-permission mapping.
const DefaultSourceURL = "https://raw.githubusercontent.com/octokit/app-permissions/main/generated/api.github.com.json"

// Meta describes where the inventory came from and when it was built.
type Meta struct {
	SourceURL            string    `json:"source_url"`
	SourceRepo           string    `json:"source_repo"`
	FetchedAt            time.Time `json:"fetched_at"`
	GitHubDocsAPIVersion string    `json:"github_docs_api_version"`
	PermissionCount      int       `json:"permission_count"`
	EndpointCount        int       `json:"endpoint_count"`
	MultiPermissionCount int       `json:"multi_permission_endpoint_count"`
	// SourceLimitations notes known gaps in the upstream octokit inventory
	// (useful baseline, not exhaustive vs GitHub's live docs/UI).
	SourceLimitations string `json:"source_limitations"`
}

// PermissionGrant records one permission+access pair required for an endpoint.
type PermissionGrant struct {
	Permission string `json:"permission"`
	Access     string `json:"access"`
}

// Endpoint is one REST operation with its path-level primary mapping from
// octokit's paths section (single permission per method+path).
type Endpoint struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Access     string `json:"access"`
	Permission string `json:"permission"`
}

// PermissionEndpoints lists documented endpoints for one GitHub App permission.
type PermissionEndpoints struct {
	Read   []string `json:"read,omitempty"`
	Write  []string `json:"write,omitempty"`
	DocURL string   `json:"doc_url,omitempty"`
}

// Inventory is the normalized endpoint catalog written to
// data/generated/endpoints-inventory.json.
type Inventory struct {
	Meta        Meta                           `json:"meta"`
	Endpoints   []Endpoint                     `json:"endpoints"`
	Permissions map[string]PermissionEndpoints `json:"permissions"`
	// EndpointIndex maps "METHOD /path" (octokit endpoint string format) to
	// every permission category that lists that operation. Entries with len>1
	// are Type A (category-overlap) dependency candidates for detect-overlap.
	EndpointIndex map[string][]PermissionGrant `json:"endpoint_index"`
}
