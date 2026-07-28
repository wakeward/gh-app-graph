package inventory

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Normalize parses octokit/app-permissions JSON bytes into a sorted Inventory.
func Normalize(sourceURL string, raw []byte, fetchedAt time.Time) (*Inventory, error) {
	var src octokitRaw
	if err := json.Unmarshal(raw, &src); err != nil {
		return nil, fmt.Errorf("decode octokit JSON: %w", err)
	}
	if len(src.Paths) == 0 || len(src.Permissions) == 0 {
		return nil, fmt.Errorf("octokit JSON missing paths or permissions sections")
	}

	endpoints := buildEndpoints(src.Paths)
	permissions := buildPermissions(src.Permissions)
	endpointIndex := buildEndpointIndex(src.Permissions)

	multi := 0
	for _, grants := range endpointIndex {
		if len(grants) > 1 {
			multi++
		}
	}

	return &Inventory{
		Meta: Meta{
			SourceURL:            sourceURL,
			SourceRepo:           "octokit/app-permissions",
			FetchedAt:            fetchedAt.UTC(),
			GitHubDocsAPIVersion: GitHubDocsAPIVersion,
			PermissionCount:      len(permissions),
			EndpointCount:        len(endpoints),
			MultiPermissionCount: multi,
			SourceLimitations:    sourceLimitationsNote,
		},
		Endpoints:     endpoints,
		Permissions:   permissions,
		EndpointIndex: endpointIndex,
	}, nil
}

const sourceLimitationsNote = "octokit/app-permissions maps endpoints to a single primary permission and currently enumerates a subset of GitHub App permissions vs the live docs/UI. Use endpoint_index multi-permission entries for Type A overlap detection; permissions absent here may still exist in GitHub docs or the App creation UI."

func buildEndpoints(paths map[string]map[string]octokitPathGrant) []Endpoint {
	var out []Endpoint
	for path, methods := range paths {
		for method, grant := range methods {
			out = append(out, Endpoint{
				Method:     strings.ToUpper(method),
				Path:       path,
				Access:     normalizeAccess(grant.Access),
				Permission: grant.Permission,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		if out[i].Method != out[j].Method {
			return out[i].Method < out[j].Method
		}
		return out[i].Permission < out[j].Permission
	})
	return out
}

func buildPermissions(perms map[string]octokitPermission) map[string]PermissionEndpoints {
	out := make(map[string]PermissionEndpoints, len(perms))
	for name, p := range perms {
		read := append([]string(nil), p.Read...)
		write := append([]string(nil), p.Write...)
		sort.Strings(read)
		sort.Strings(write)
		out[name] = PermissionEndpoints{
			Read:   read,
			Write:  write,
			DocURL: p.URL,
		}
	}
	return out
}

func buildEndpointIndex(perms map[string]octokitPermission) map[string][]PermissionGrant {
	index := make(map[string][]PermissionGrant)
	for permName, p := range perms {
		addIndexEntries(index, permName, "read", p.Read)
		addIndexEntries(index, permName, "write", p.Write)
	}

	for key, grants := range index {
		sort.Slice(grants, func(i, j int) bool {
			if grants[i].Permission != grants[j].Permission {
				return grants[i].Permission < grants[j].Permission
			}
			return grants[i].Access < grants[j].Access
		})
		index[key] = grants
	}
	return index
}

func addIndexEntries(index map[string][]PermissionGrant, permission, access string, endpoints []string) {
	for _, ep := range endpoints {
		ep = strings.TrimSpace(ep)
		if ep == "" {
			continue
		}
		grants := index[ep]
		if !grantContains(grants, permission, access) {
			index[ep] = append(grants, PermissionGrant{
				Permission: permission,
				Access:     access,
			})
		}
	}
}

func grantContains(grants []PermissionGrant, permission, access string) bool {
	for _, g := range grants {
		if g.Permission == permission && g.Access == access {
			return true
		}
	}
	return false
}

// normalizeAccess maps octokit access values to read/write. GitHub endpoint
// docs occasionally show "admin"; App grants only expose read/write, so admin
// folds to write (see docs/methodology.md).
func normalizeAccess(access string) string {
	switch strings.ToLower(access) {
	case "read":
		return "read"
	case "admin":
		return "write"
	default:
		return "write"
	}
}
