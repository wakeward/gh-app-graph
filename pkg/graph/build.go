package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wakeward/github-app-permissions-graph/pkg/inventory"
	"github.com/wakeward/github-app-permissions-graph/pkg/model"
	"gopkg.in/yaml.v3"
)

// BuildMeta records provenance for a build run.
type BuildMeta struct {
	BuiltAt               time.Time `json:"built_at" yaml:"built_at"`
	GitHubDocsAPIVersion  string    `json:"github_docs_api_version" yaml:"github_docs_api_version"`
	PermissionCount       int       `json:"permission_count" yaml:"permission_count"`
	DependencyEdgeCount   int       `json:"dependency_edge_count" yaml:"dependency_edge_count"`
	ToxicCombinationCount int       `json:"toxic_combination_count" yaml:"toxic_combination_count"`
	DraftDataNotice       string    `json:"draft_data_notice" yaml:"draft_data_notice"`
}

// Resolved is the merged output written to data/permissions.resolved.{yaml,json}.
type Resolved struct {
	BuildMeta         BuildMeta                `json:"build_meta" yaml:"build_meta"`
	Permissions       []model.Permission       `json:"permissions" yaml:"permissions"`
	Dependencies      []model.DependencyEdge   `json:"dependencies" yaml:"dependencies"`
	ToxicCombinations []model.ToxicCombination `json:"toxic_combinations" yaml:"toxic_combinations"`
}

const draftDataNotice = "permissions metadata is DRAFT seed data; dependencies include auto-detected Type A overlaps from octokit/app-permissions (subset of live GitHub permissions)."

// LoadPermissionsDir loads every *.yaml file in dir as permissions and
// validates global api_key uniqueness.
func LoadPermissionsDir(dir string) ([]model.Permission, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)

	seen := make(map[string]string)
	var out []model.Permission
	for _, path := range matches {
		var batch []model.Permission
		if err := readYAML(path, &batch); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		for _, p := range batch {
			if prev, ok := seen[p.APIKey]; ok {
				return nil, fmt.Errorf("duplicate api_key %q in %s and %s", p.APIKey, prev, path)
			}
			seen[p.APIKey] = path
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].APIKey < out[j].APIKey })
	return out, nil
}

// LoadDependencyEdges reads dependency edges. Missing files yield empty.
func LoadDependencyEdges(path string) ([]model.DependencyEdge, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}
	var edges []model.DependencyEdge
	if err := readYAML(path, &edges); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return edges, nil
}

// LoadToxicCombinations reads toxic combination rules. Missing files yield empty.
func LoadToxicCombinations(path string) ([]model.ToxicCombination, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}
	var combos []model.ToxicCombination
	if err := readYAML(path, &combos); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	sort.Slice(combos, func(i, j int) bool { return combos[i].ID < combos[j].ID })
	return combos, nil
}

func readYAML(path string, dest any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, dest)
}

// MergeDependencies merges layers in ascending precedence. Each later layer
// replaces the entire requires list for a permission+access key when present
// (manual > scraped-draft > overlap).
func MergeDependencies(layers ...[]model.DependencyEdge) []model.DependencyEdge {
	byKey := make(map[string]model.DependencyEdge)
	for _, layer := range layers {
		for _, edge := range layer {
			byKey[edgeKey(edge.Permission, edge.Access)] = edge
		}
	}
	out := make([]model.DependencyEdge, 0, len(byKey))
	for _, edge := range byKey {
		sort.Slice(edge.Requires, func(i, j int) bool {
			if edge.Requires[i].Permission != edge.Requires[j].Permission {
				return edge.Requires[i].Permission < edge.Requires[j].Permission
			}
			return edge.Requires[i].Access < edge.Requires[j].Access
		})
		out = append(out, edge)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Permission != out[j].Permission {
			return out[i].Permission < out[j].Permission
		}
		return out[i].Access < out[j].Access
	})
	return out
}

func edgeKey(permission string, access model.AccessLevel) string {
	return permission + "\x00" + string(access)
}

// Build assembles the resolved graph from loaded inputs.
func Build(permissions []model.Permission, overlap, scraped, manual []model.DependencyEdge, toxic []model.ToxicCombination, apiVersion string) Resolved {
	deps := MergeDependencies(overlap, scraped, manual)
	if toxic == nil {
		toxic = []model.ToxicCombination{}
	}
	if apiVersion == "" {
		apiVersion = inventory.GitHubDocsAPIVersion
	}
	return Resolved{
		BuildMeta: BuildMeta{
			BuiltAt:               time.Now().UTC(),
			GitHubDocsAPIVersion:  apiVersion,
			PermissionCount:       len(permissions),
			DependencyEdgeCount:   len(deps),
			ToxicCombinationCount: len(toxic),
			DraftDataNotice:       draftDataNotice,
		},
		Permissions:       permissions,
		Dependencies:      deps,
		ToxicCombinations: toxic,
	}
}

// WriteResolved writes resolved data as YAML and JSON.
func WriteResolved(yamlPath, jsonPath string, resolved Resolved) error {
	if err := writeYAML(yamlPath, resolved); err != nil {
		return err
	}
	data, err := json.MarshalIndent(resolved, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", jsonPath, err)
	}
	return nil
}

// WriteToxicJSON writes toxic combinations alone for lightweight consumers.
func WriteToxicJSON(path string, combos []model.ToxicCombination) error {
	if combos == nil {
		combos = []model.ToxicCombination{}
	}
	data, err := json.MarshalIndent(combos, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func writeYAML(path string, v any) error {
	var buf strings.Builder
	buf.WriteString("# Generated by cmd/build - do not edit by hand\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode yaml: %w", err)
	}
	enc.Close()
	return os.WriteFile(path, []byte(buf.String()), 0o644)
}

// DependencyMap indexes merged dependencies by permission+access.
func DependencyMap(deps []model.DependencyEdge) map[string]model.DependencyEdge {
	m := make(map[string]model.DependencyEdge, len(deps))
	for _, d := range deps {
		m[edgeKey(d.Permission, d.Access)] = d
	}
	return m
}
