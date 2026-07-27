// Package model defines the shared data types for the GitHub App permissions
// graph: permissions, their dependencies, and toxic combinations. Other Go
// tools (e.g. gh-app-check) are meant to import this package directly to
// share these types, rather than re-deriving them from a serialized format -
// see the root README's "Relationship to gh-app-check" section.
package model

// Severity is the risk severity assigned to a permission on its own, before
// considering any combination with other permissions.
type Severity string

const (
	SeverityCritical Severity = "Critical"
	SeverityHigh     Severity = "High"
	SeverityMedium   Severity = "Medium"
	SeverityLow      Severity = "Low"
)

// ImpactPlane categorizes a permission by what it affects if abused: the
// Control Plane (settings/configuration such as branch protection), the
// Data/Execution Plane (source code, CI runs), or the Governance/Identity
// Plane (who can access what).
type ImpactPlane string

const (
	ImpactPlaneControl            ImpactPlane = "control"
	ImpactPlaneDataExecution      ImpactPlane = "data_execution"
	ImpactPlaneGovernanceIdentity ImpactPlane = "governance_identity"
)

// DocStatus records how well-documented a permission is in GitHub's public
// docs. Some Enterprise permissions are visible in the App creation UI but
// have no published REST endpoints yet.
type DocStatus string

const (
	DocStatusDocumented          DocStatus = "documented"
	DocStatusUndocumentedPreview DocStatus = "undocumented_preview"
	DocStatusUnconfirmedKey      DocStatus = "unconfirmed_key"
)

// AccessLevel is the grant level for a permission. GitHub Apps can only ever
// be granted "read" or "write" - endpoints whose docs show an "admin"
// requirement map to AccessWrite (see docs/methodology.md).
type AccessLevel string

const (
	AccessRead  AccessLevel = "read"
	AccessWrite AccessLevel = "write"
)

// Permission is one entry from data/permissions/*.yaml: a single named
// GitHub App permission and everything known about its standalone risk.
type Permission struct {
	Name               string      `yaml:"name" json:"name"`
	APIKey             string      `yaml:"api_key" json:"api_key"`
	Category           string      `yaml:"category" json:"category"`
	AccessLevels       []string    `yaml:"access_levels" json:"access_levels"`
	Overview           string      `yaml:"overview" json:"overview"`
	Severity           Severity    `yaml:"severity" json:"severity"`
	SecurityNotes      string      `yaml:"security_notes" json:"security_notes"`
	NeedsInvestigation bool        `yaml:"needs_investigation" json:"needs_investigation"`
	DocStatus          DocStatus   `yaml:"doc_status" json:"doc_status"`
	ImpactPlane        ImpactPlane `yaml:"impact_plane" json:"impact_plane"`
}
