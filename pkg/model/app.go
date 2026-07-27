package model

// AppPermissionSet is a declared set of permissions for one GitHub App:
// either hand-written for an internal app (apps/*.yaml) or fetched live from
// the public, unauthenticated GET /apps/{app_slug} endpoint for a
// third-party app's default requested permissions.
type AppPermissionSet struct {
	Name        string            `yaml:"name" json:"name"`
	Slug        string            `yaml:"slug,omitempty" json:"slug,omitempty"`
	Permissions map[string]string `yaml:"permissions" json:"permissions"` // api_key -> "read" | "write"
}

// HasGrant reports whether the set includes api_key at least at the given
// access level (a "write" grant satisfies a "read" requirement too).
func (s AppPermissionSet) HasGrant(apiKey string, access AccessLevel) bool {
	granted, ok := s.Permissions[apiKey]
	if !ok {
		return false
	}
	if access == AccessRead {
		return granted == string(AccessRead) || granted == string(AccessWrite)
	}
	return granted == string(AccessWrite)
}
