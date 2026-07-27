// Command scrape-prose scrapes Type B (prose) dependency drafts for
// High/Critical severity permissions and writes
// data/generated/scraped-dependencies-draft.yaml with needs_review: true.
//
// TODO(scrape-prose): implement using pkg/scrape. Support a --widen flag so
// the quarterly Cursor Automation can run the wide sweep across every
// flagged permission instead of the default High/Critical-only scope.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "scrape-prose: not yet implemented")
	os.Exit(1)
}
