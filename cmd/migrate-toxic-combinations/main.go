// Command migrate-toxic-combinations is a one-time migration: it converts
// the "toxic-combinations" CSV and the threat-oriented writeup markdown into
// data/toxic-combinations.yaml, normalizing human-readable permission names
// to the canonical api_key values defined in data/permissions/*.yaml. It
// does not archive the source documents into the repo - see cmd/migrate-csv
// for why.
//
// TODO(toxic-combinations): implement the one-time conversion.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "migrate-toxic-combinations: not yet implemented")
	os.Exit(1)
}
