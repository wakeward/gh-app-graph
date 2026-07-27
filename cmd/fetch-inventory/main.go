// Command fetch-inventory pulls and normalizes the octokit/app-permissions
// generated JSON into data/generated/endpoints-inventory.json.
//
// TODO(fetch-inventory): implement using pkg/inventory.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "fetch-inventory: not yet implemented")
	os.Exit(1)
}
