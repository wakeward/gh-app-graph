// Package inventory fetches and normalizes the octokit/app-permissions
// generated JSON (the base, single-permission endpoint inventory) into
// data/generated/endpoints-inventory.json, and detects Type A
// (category-overlap) dependency edges from it.
//
// TODO(fetch-inventory, detect-overlap): implement the fetch + normalize +
// overlap-detection logic. Pin an explicit source API version in the output
// so a GitHub API version bump is visible as a distinct diff from a content
// change (see the plan's gap analysis, item 1).
package inventory
