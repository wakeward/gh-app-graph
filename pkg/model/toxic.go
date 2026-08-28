package model

// BlastRadius is how severe a toxic combination's worst-case outcome is.
type BlastRadius string

const (
	BlastRadiusCritical BlastRadius = "Critical"
	BlastRadiusHigh     BlastRadius = "High"
	BlastRadiusMedium   BlastRadius = "Medium"
)

// PermissionGrant is one (api_key, access) pair required to be present for a
// toxic combination to apply.
type PermissionGrant struct {
	APIKey string      `yaml:"api_key" json:"api_key"`
	Access AccessLevel `yaml:"access" json:"access"`
}

// ToxicCombination is a named attack technique an app installation enables.
// Most entries require two or more independently-grantable permissions
// co-present on the same identity; single-grant entries are also valid when
// the technique/exploit path is worth surfacing explicitly (users often
// accept a scope without understanding what it unlocks).
// This is distinct from a DependencyEdge, which is a technical prerequisite.
type ToxicCombination struct {
	ID                          string            `yaml:"id" json:"id"`
	Technique                   string            `yaml:"technique" json:"technique"`
	Permissions                 []PermissionGrant `yaml:"permissions" json:"permissions"`
	BlastRadius                 BlastRadius       `yaml:"blast_radius" json:"blast_radius"`
	RiskDescription             string            `yaml:"risk_description" json:"risk_description"`
	ExploitPath                 string            `yaml:"exploit_path" json:"exploit_path"`
	OverlapsTechnicalDependency bool              `yaml:"overlaps_technical_dependency" json:"overlaps_technical_dependency"`
	PlatformAvailability        PlatformAvailability `yaml:"platform_availability,omitempty" json:"platform_availability,omitempty"`
	Source                      string            `yaml:"source" json:"source"`
}
