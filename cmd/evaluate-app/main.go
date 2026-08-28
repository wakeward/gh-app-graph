// Command evaluate-app scores a given app's declared permission set - either
// a hand-written local file under apps/*.yaml for an internal app, or a
// live, unauthenticated fetch of GET /apps/{app_slug} for a public
// third-party app - against data/toxic-combinations.yaml using pkg/eval,
// and renders reports/<app-slug>-risk-report.md.
//
// TODO(evaluate-app): implement using pkg/eval, pkg/ghapps, and pkg/graph
// (for the unsatisfied-dependency data-quality check).
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "evaluate-app: not yet implemented")
	os.Exit(1)
}
