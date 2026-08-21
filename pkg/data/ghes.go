package data

import "github.com/wakeward/gh-app-graph/pkg/platform"

// LoadGHESOnlyAPIKeys returns api_key values for permissions that apply only
// on GitHub Enterprise Server, from the embedded permission catalog.
func LoadGHESOnlyAPIKeys() (map[string]struct{}, error) {
	perms, err := LoadResolvedPermissions()
	if err != nil {
		return nil, err
	}
	return platform.GHESOnlyKeys(perms), nil
}
