// Command build merges data/permissions/*.yaml, data/dependencies-manual.yaml,
// data/generated/overlap-dependencies.yaml, and data/toxic-combinations.yaml
// into data/permissions.resolved.{yaml,json} and
// data/generated/toxic-combinations.json, validates api_key uniqueness, and
// renders docs/permissions-table.md, docs/permissions-graph.md, and
// docs/toxic-combinations.md.
//
// TODO(build-render): implement using pkg/graph and pkg/render.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "build: not yet implemented")
	os.Exit(1)
}
