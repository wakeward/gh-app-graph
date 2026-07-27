// Command migrate-csv is a one-time migration: it converts the original
// "GitHub Apps Permissions" CSV into data/permissions/*.yaml, seeds the
// undocumented Enterprise permissions visible only in the App creation UI
// screenshot, and archives the original CSV into data/sources/raw/ as an
// immutable source document.
//
// TODO(migrate-csv): implement the one-time conversion.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "migrate-csv: not yet implemented")
	os.Exit(1)
}
