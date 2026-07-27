// Package model defines the shared data types for the GitHub App permissions
// graph: permissions, their dependencies, and toxic combinations. Other Go
// tools (e.g. gh-app-check) are meant to import this package directly to
// share these types, rather than re-deriving them from a serialized format -
// see the root README's "Relationship to gh-app-check" section.
package model

// Severity is the risk severity assigned to a permission at a specific
// access level, on its own, before considering any combination with other
// permissions.
type Severity string

const (
	SeverityCritical Severity = "Critical"
	SeverityHigh     Severity = "High"
	SeverityMedium   Severity = "Medium"
	SeverityLow      Severity = "Low"
	// SeverityInformational: the permission grants no meaningful standalone
	// risk (e.g. read-only access to something already public, or a no-op
	// access level like "Followers: read").
	SeverityInformational Severity = "Informational"
	// SeverityUnknown: not yet assessed - always paired with
	// NeedsInvestigation: true. Never rely on this value for scoring.
	SeverityUnknown Severity = "Unknown"
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
// docs.
type DocStatus string

const (
	DocStatusDocumented DocStatus = "documented"
	// DocStatusUndocumentedPreview: visible in the App creation UI (often
	// marked "Preview") but with no published REST endpoints yet.
	DocStatusUndocumentedPreview DocStatus = "undocumented_preview"
	// DocStatusUnconfirmedKey: the api_key is inferred (from an endpoint
	// path, go-github's schema, etc.) but not confirmed against a live
	// GitHub App manifest or permissions picker.
	DocStatusUnconfirmedKey DocStatus = "unconfirmed_key"
	// DocStatusDisputed: independent signals disagree on whether this is a
	// real, current permission (e.g. present in one source, absent from the
	// live UI and from go-github's schema).
	DocStatusDisputed DocStatus = "disputed"
)

// AccessLevel is the grant level for a permission. GitHub Apps can only ever
// be granted "read" or "write" - endpoints whose docs show an "admin"
// requirement map to AccessWrite (see docs/methodology.md).
type AccessLevel string

const (
	AccessRead  AccessLevel = "read"
	AccessWrite AccessLevel = "write"
)

// AccessLevelDetail is the risk assessment for one specific access level of
// a permission. Severity commonly differs between read and write for the
// same permission (e.g. Administration read is reconnaissance-level,
// Administration write is Critical), so this is tracked per access level
// rather than once per permission.
type AccessLevelDetail struct {
	Access        AccessLevel `yaml:"access" json:"access"`
	Severity      Severity    `yaml:"severity" json:"severity"`
	SecurityNotes string      `yaml:"security_notes,omitempty" json:"security_notes,omitempty"`
}

// Permission is one entry from data/permissions/*.yaml: a single named
// GitHub App permission and everything known about its standalone risk.
type Permission struct {
	Name               string              `yaml:"name" json:"name"`
	APIKey             string              `yaml:"api_key" json:"api_key"`
	Category           string              `yaml:"category" json:"category"`
	Overview           string              `yaml:"overview" json:"overview"`
	AccessLevels       []AccessLevelDetail `yaml:"access_levels" json:"access_levels"`
	NeedsInvestigation bool                `yaml:"needs_investigation" json:"needs_investigation"`
	DocStatus          DocStatus           `yaml:"doc_status" json:"doc_status"`
	ImpactPlane        ImpactPlane         `yaml:"impact_plane" json:"impact_plane"`
}
