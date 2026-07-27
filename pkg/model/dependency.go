package model

// DependencyType records how a "requires" edge was discovered.
type DependencyType string

const (
	// DependencyTypeCategoryOverlap (Type A): the same endpoint is listed
	// under two permission categories at once - free, automatable, high
	// confidence.
	DependencyTypeCategoryOverlap DependencyType = "category-overlap"
	// DependencyTypeProseScraped (Type B): the requirement is only stated in
	// an endpoint's own doc-page prose - needs scraping plus human review.
	DependencyTypeProseScraped DependencyType = "prose-scraped"
	// DependencyTypeManual: hand-confirmed by a human, highest precedence.
	DependencyTypeManual DependencyType = "manual"
)

// Confidence records how much a dependency edge can be trusted without
// further human review.
type Confidence string

const (
	ConfidenceHigh Confidence = "high"
	ConfidenceLow  Confidence = "low"
)

// RequiredPermission is one permission+access pair required by another
// permission in order to function.
type RequiredPermission struct {
	Permission  string         `yaml:"permission" json:"permission"`
	Access      AccessLevel    `yaml:"access" json:"access"`
	Condition   string         `yaml:"condition,omitempty" json:"condition,omitempty"`
	Confidence  Confidence     `yaml:"confidence" json:"confidence"`
	Type        DependencyType `yaml:"type" json:"type"`
	SourceURL   string         `yaml:"source_url,omitempty" json:"source_url,omitempty"`
	NeedsReview bool           `yaml:"needs_review,omitempty" json:"needs_review,omitempty"`
}

// DependencyEdge lists everything a given permission+access requires to
// function. This is a technical prerequisite, distinct from a
// ToxicCombination (see toxic.go) which is a risk-correlation rule rather
// than a functional requirement.
type DependencyEdge struct {
	Permission string                `yaml:"permission" json:"permission"`
	Access     AccessLevel           `yaml:"access" json:"access"`
	Requires   []RequiredPermission  `yaml:"requires" json:"requires"`
}
