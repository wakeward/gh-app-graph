// Command migrate-toxic-combinations is a one-time migration: it converts
// the "toxic-combinations" CSV and the threat-oriented writeup markdown into
// data/toxic-combinations.yaml, normalizing human-readable permission names
// to the canonical api_key values defined in data/permissions/*.yaml. It
// does not archive the source documents into the repo - see cmd/migrate-csv
// for why.
//
// Usage:
//
//	go run ./cmd/migrate-toxic-combinations \
//	  --csv "/path/to/GitHub Apps Permissions - toxic-combinations.csv" \
//	  --writeup "/path/to/GitHub permissions (...).md"
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/wakeward/gh-app-graph/pkg/fileio"
	"github.com/wakeward/gh-app-graph/pkg/graph"
	"github.com/wakeward/gh-app-graph/pkg/model"
	"github.com/wakeward/gh-app-graph/pkg/platform"
)

func main() {
	csvPath := flag.String("csv", "", "path to the toxic-combinations CSV (required)")
	writeupPath := flag.String("writeup", "", "optional path to the threat-oriented writeup markdown")
	outPath := flag.String("out", "data/toxic-combinations.yaml", "output YAML path")
	permissionsDir := flag.String("permissions", "data/permissions", "directory of permission YAML for api_key validation")
	overlapPath := flag.String("overlap-deps", "data/generated/overlap-dependencies.yaml", "Type A overlap deps for overlaps_technical_dependency flag")
	flag.Parse()

	if *csvPath == "" {
		fmt.Fprintln(os.Stderr, "migrate-toxic-combinations: --csv is required")
		os.Exit(1)
	}

	known, err := loadKnownAPIKeys(*permissionsDir)
	if err != nil {
		fail("load permissions", err)
	}

	ghesKeys, err := loadGHESOnlyKeys(*permissionsDir)
	if err != nil {
		fail("load GHES-only permissions", err)
	}

	rows, err := readCSV(*csvPath)
	if err != nil {
		fail("read CSV", err)
	}

	combos, warnings, err := buildFromCSV(rows, known)
	if err != nil {
		fail("parse CSV", err)
	}

	if *writeupPath != "" {
		supplement, err := supplementFromWriteup(*writeupPath, combos, known)
		if err != nil {
			fail("parse writeup", err)
		}
		combos = append(combos, supplement...)
	}

	overlap, err := graph.LoadDependencyEdges(*overlapPath)
	if err != nil {
		fail("load overlap deps", err)
	}
	for i := range combos {
		combos[i].OverlapsTechnicalDependency = comboOverlapsDependency(combos[i], overlap)
		combos[i].PlatformAvailability = inferComboPlatform(combos[i], ghesKeys)
		validateGrants(&combos[i], known, &warnings)
	}

	sort.Slice(combos, func(i, j int) bool { return combos[i].ID < combos[j].ID })

	if err := writeToxicYAML(*outPath, combos); err != nil {
		fail("write output", err)
	}

	fmt.Printf("migrate-toxic-combinations: wrote %d toxic combinations to %s\n", len(combos), *outPath)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "migrate-toxic-combinations: warning: %s\n", w)
	}
}

func loadKnownAPIKeys(dir string) (map[string]struct{}, error) {
	perms, err := graph.LoadPermissionsDir(dir)
	if err != nil {
		return nil, err
	}
	known := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		known[p.APIKey] = struct{}{}
	}
	return known, nil
}

func loadGHESOnlyKeys(dir string) (map[string]struct{}, error) {
	perms, err := graph.LoadPermissionsDir(dir)
	if err != nil {
		return nil, err
	}
	return platform.GHESOnlyKeys(perms), nil
}

func validateGrants(combo *model.ToxicCombination, known map[string]struct{}, warnings *[]string) {
	for _, g := range combo.Permissions {
		if _, ok := known[g.APIKey]; !ok {
			*warnings = append(*warnings, fmt.Sprintf("%s: api_key %q not found in %s", combo.ID, g.APIKey, "data/permissions"))
		}
	}
}

func comboOverlapsDependency(combo model.ToxicCombination, deps []model.DependencyEdge) bool {
	depMap := graph.DependencyMap(deps)
	grants := make(map[string]struct{}, len(combo.Permissions))
	for _, g := range combo.Permissions {
		grants[grantKey(g)] = struct{}{}
	}
	for _, g := range combo.Permissions {
		edge, ok := depMap[depKey(g.APIKey, g.Access)]
		if !ok {
			continue
		}
		for _, req := range edge.Requires {
			if _, ok := grants[grantKey(model.PermissionGrant{APIKey: req.Permission, Access: req.Access})]; ok {
				return true
			}
		}
	}
	return false
}

func grantKey(g model.PermissionGrant) string {
	return g.APIKey + "\x00" + string(g.Access)
}

func depKey(apiKey string, access model.AccessLevel) string {
	return apiKey + "\x00" + string(access)
}

func writeToxicYAML(path string, combos []model.ToxicCombination) error {
	if err := fileio.MkdirAll(filepath.Dir(path)); err != nil {
		return err
	}
	header := `# Toxic permission combinations
#
# Status: curated in-repo. Seeded initially via cmd/migrate-toxic-combinations
# from a working CSV + writeup; edit this file directly once the list stabilizes.
# That migrate command is bootstrap-only and planned for removal.
`
	return fileio.WriteYAML(path, header, combos)
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "migrate-toxic-combinations: %s: %v\n", step, err)
	os.Exit(1)
}
