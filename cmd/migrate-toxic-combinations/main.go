// Command migrate-toxic-combinations is a one-time migration: it converts
// the "toxic-combinations" CSV and the threat-oriented writeup markdown into
// data/toxic-combinations.yaml, normalizing human-readable permission names
// to the canonical api_key values defined in data/permissions/*.yaml, and
// archives both source documents into data/sources/raw/.
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
