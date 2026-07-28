// Command build merges data/permissions/*.yaml, dependency sources, and
// optional toxic-combinations into data/permissions.resolved.{yaml,json},
// validates api_key uniqueness, and renders docs/.
//
// Usage:
//
//	go run ./cmd/build
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wakeward/github-app-permissions-graph/pkg/graph"
	"github.com/wakeward/github-app-permissions-graph/pkg/inventory"
	"github.com/wakeward/github-app-permissions-graph/pkg/render"
)

func main() {
	permissionsDir := flag.String("permissions", "data/permissions", "directory of draft permission YAML")
	manualPath := flag.String("manual-deps", "data/dependencies-manual.yaml", "human-reviewed dependency edges")
	scrapedPath := flag.String("scraped-deps", "data/generated/scraped-dependencies-draft.yaml", "Type B draft dependencies")
	overlapPath := flag.String("overlap-deps", "data/generated/overlap-dependencies.yaml", "Type A overlap dependencies")
	toxicPath := flag.String("toxic", "data/toxic-combinations.yaml", "toxic combination rules")
	resolvedYAML := flag.String("resolved-yaml", "data/permissions.resolved.yaml", "merged output YAML")
	resolvedJSON := flag.String("resolved-json", "data/permissions.resolved.json", "merged output JSON")
	toxicJSON := flag.String("toxic-json", "data/generated/toxic-combinations.json", "toxic combinations JSON")
	docsDir := flag.String("docs", "docs", "generated markdown output directory")
	flag.Parse()

	permissions, err := graph.LoadPermissionsDir(*permissionsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build: load permissions: %v\n", err)
		os.Exit(1)
	}

	overlap, err := graph.LoadDependencyEdges(*overlapPath)
	if err != nil {
		fail("overlap deps", err)
	}
	scraped, err := graph.LoadDependencyEdges(*scrapedPath)
	if err != nil {
		fail("scraped deps", err)
	}
	manual, err := graph.LoadDependencyEdges(*manualPath)
	if err != nil {
		fail("manual deps", err)
	}
	toxic, err := graph.LoadToxicCombinations(*toxicPath)
	if err != nil {
		fail("toxic combinations", err)
	}

	resolved := graph.Build(permissions, overlap, scraped, manual, toxic, inventory.GitHubDocsAPIVersion)

	if err := graph.WriteResolved(*resolvedYAML, *resolvedJSON, resolved); err != nil {
		fail("write resolved", err)
	}
	if err := graph.WriteToxicJSON(*toxicJSON, resolved.ToxicCombinations); err != nil {
		fail("write toxic json", err)
	}
	if err := render.WriteAll(*docsDir, resolved); err != nil {
		fail("render docs", err)
	}

	fmt.Printf("build: wrote %s + %s (%d permissions, %d dependency edges, %d toxic combos)\n",
		*resolvedYAML, *resolvedJSON,
		resolved.BuildMeta.PermissionCount,
		resolved.BuildMeta.DependencyEdgeCount,
		resolved.BuildMeta.ToxicCombinationCount,
	)
	fmt.Printf("build: rendered docs in %s\n", *docsDir)
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "build: %s: %v\n", step, err)
	os.Exit(1)
}
