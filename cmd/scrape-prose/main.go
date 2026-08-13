// Command scrape-prose scrapes Type B (prose) dependency drafts for
// High/Critical severity permissions and writes
// data/generated/scraped-dependencies-draft.yaml with needs_review: true.
//
// Usage:
//
//	go run ./cmd/scrape-prose
//	go run ./cmd/scrape-prose --widen
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wakeward/gh-app-graph/pkg/graph"
	"github.com/wakeward/gh-app-graph/pkg/inventory"
	"github.com/wakeward/gh-app-graph/pkg/scrape"
)

func main() {
	permissionsDir := flag.String("permissions", "data/permissions", "directory of draft permission YAML")
	inventoryPath := flag.String("inventory", "data/generated/endpoints-inventory.json", "normalized endpoint inventory")
	outputPath := flag.String("output", "data/generated/scraped-dependencies-draft.yaml", "Type B draft output")
	widen := flag.Bool("widen", false, "also scrape needs_investigation permissions (quarterly wide sweep)")
	flag.Parse()

	permissions, err := graph.LoadPermissionsDir(*permissionsDir)
	if err != nil {
		fail("load permissions", err)
	}
	inv, err := inventory.LoadJSON(*inventoryPath)
	if err != nil {
		fail("load inventory", err)
	}

	edges, warnings, err := scrape.DetectProseDependencies(permissions, inv, *widen, scrape.DefaultFetcher)
	if err != nil {
		fail("detect prose dependencies", err)
	}
	if err := scrape.WriteDraftYAML(*outputPath, edges); err != nil {
		fail("write output", err)
	}

	fmt.Printf("scrape-prose: wrote %s (%d permission entries", *outputPath, len(edges))
	requireCount := 0
	for _, e := range edges {
		requireCount += len(e.Requires)
	}
	fmt.Printf(", %d requires edges)\n", requireCount)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "scrape-prose: warning: %s\n", w)
	}
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "scrape-prose: %s: %v\n", step, err)
	os.Exit(1)
}
