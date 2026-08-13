// Command fetch-inventory pulls and normalizes the octokit/app-permissions
// generated JSON into data/generated/endpoints-inventory.json.
//
// Usage:
//
//	go run ./cmd/fetch-inventory
//	go run ./cmd/fetch-inventory --input /path/to/api.github.com.json --output data/generated/endpoints-inventory.json
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wakeward/gh-app-graph/pkg/inventory"
)

func main() {
	sourceURL := flag.String("url", inventory.DefaultSourceURL, "URL to fetch octokit/app-permissions JSON from")
	inputPath := flag.String("input", "", "local JSON file instead of fetching from --url")
	outputPath := flag.String("output", "data/generated/endpoints-inventory.json", "output path for normalized inventory")
	flag.Parse()

	var (
		raw        []byte
		provenance string
		err        error
	)
	switch {
	case *inputPath != "":
		raw, err = inventory.ReadFile(*inputPath)
		provenance = "file://" + *inputPath
	case *sourceURL != "":
		raw, err = inventory.FetchHTTP(*sourceURL)
		provenance = *sourceURL
	default:
		fmt.Fprintln(os.Stderr, "fetch-inventory: specify --input or --url")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch-inventory: load source: %v\n", err)
		os.Exit(1)
	}

	inv, err := inventory.BuildFromRaw(provenance, raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch-inventory: normalize: %v\n", err)
		os.Exit(1)
	}

	if err := inventory.WriteJSON(*outputPath, inv); err != nil {
		fmt.Fprintf(os.Stderr, "fetch-inventory: write output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("fetch-inventory: wrote %s (%d permissions, %d endpoints, %d multi-permission endpoints, api_version=%s)\n",
		*outputPath,
		inv.Meta.PermissionCount,
		inv.Meta.EndpointCount,
		inv.Meta.MultiPermissionCount,
		inv.Meta.GitHubDocsAPIVersion,
	)
}
