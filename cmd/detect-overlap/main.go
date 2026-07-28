// Command detect-overlap derives Type A (category-overlap) dependency edges
// from data/generated/endpoints-inventory.json and writes
// data/generated/overlap-dependencies.yaml.
//
// Usage:
//
//	go run ./cmd/detect-overlap
//	go run ./cmd/detect-overlap --inventory data/generated/endpoints-inventory.json --output data/generated/overlap-dependencies.yaml
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wakeward/github-app-permissions-graph/pkg/inventory"
)

func main() {
	inventoryPath := flag.String("inventory", "data/generated/endpoints-inventory.json", "normalized endpoint inventory JSON")
	outputPath := flag.String("output", "data/generated/overlap-dependencies.yaml", "output path for overlap dependency edges")
	flag.Parse()

	inv, err := inventory.LoadJSON(*inventoryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "detect-overlap: load inventory: %v\n", err)
		os.Exit(1)
	}

	edges := inventory.DetectOverlapDependencies(inv)
	if err := inventory.WriteOverlapYAML(*outputPath, edges); err != nil {
		fmt.Fprintf(os.Stderr, "detect-overlap: write output: %v\n", err)
		os.Exit(1)
	}

	requireCount := 0
	for _, e := range edges {
		requireCount += len(e.Requires)
	}
	fmt.Printf("detect-overlap: wrote %s (%d permission entries, %d requires edges from %d multi-permission endpoints)\n",
		*outputPath, len(edges), requireCount, inv.Meta.MultiPermissionCount)
}
