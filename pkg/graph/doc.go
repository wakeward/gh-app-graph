// Package graph resolves the three dependency data sources into the final
// merged view (data/permissions.resolved.{yaml,json}): manual edges take
// precedence over scraped-draft edges, which take precedence over
// auto-detected category-overlap edges. It also validates that every
// permission's api_key is globally unique across data/permissions/*.yaml,
// failing loudly if that invariant is ever violated.
//
// TODO(build-render, gap-fixes): implement the merge-precedence resolver and
// the api_key uniqueness check, with fixture-based tests per the plan's gap
// analysis (merge-precedence bugs would silently produce a wrong risk
// conclusion).
package graph
