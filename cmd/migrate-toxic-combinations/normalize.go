package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/wakeward/github-app-permissions-graph/pkg/model"
)

// nameToAPIKey maps human-readable permission labels from the CSV/writeup to
// canonical api_key values in data/permissions/*.yaml.
var nameToAPIKey = map[string]string{
	"administration":              "administration",
	"contents":                    "contents",
	"workflows":                   "workflows",
	"organization administration": "organization_administration",
	"organization admin":          "organization_administration",
	"members":                     "members",
	"deployments":                 "deployments",
	"packages":                    "packages",
	"checks":                      "checks",
	"pull requests":               "pull_requests",
	"statuses":                    "statuses",
}

var slugRE = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeName(name string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(name))), " ")
}

func resolveAPIKey(name string, known map[string]struct{}) (string, error) {
	norm := normalizeName(name)
	if key, ok := nameToAPIKey[norm]; ok {
		return key, nil
	}
	// Fall back to snake_case of the normalized label.
	key := strings.ReplaceAll(norm, " ", "_")
	if _, ok := known[key]; ok {
		return key, nil
	}
	return "", fmt.Errorf("unknown permission name %q (normalized %q)", name, norm)
}

func parseAccessLevel(text string) (model.AccessLevel, error) {
	switch strings.ToLower(strings.TrimSpace(text)) {
	case "read", "read-only", "ro":
		return model.AccessRead, nil
	case "write", "read / write", "read/write":
		return model.AccessWrite, nil
	default:
		return "", fmt.Errorf("unrecognized access level %q", text)
	}
}

func slugify(s string) string {
	return strings.Trim(slugRE.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func grantSetKey(grants []model.PermissionGrant) string {
	parts := make([]string, len(grants))
	for i, g := range grants {
		parts[i] = g.APIKey + ":" + string(g.Access)
	}
	sort.Strings(parts)
	return strings.Join(parts, "+")
}
