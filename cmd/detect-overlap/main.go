// Command detect-overlap derives Type A (category-overlap) dependency edges
// from data/generated/endpoints-inventory.json and writes
// data/generated/overlap-dependencies.yaml.
//
// TODO(detect-overlap): implement using pkg/inventory and pkg/graph.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "detect-overlap: not yet implemented")
	os.Exit(1)
}
