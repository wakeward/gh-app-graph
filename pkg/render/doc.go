// Package render generates the markdown/Mermaid documentation
// (docs/permissions-table.md, docs/permissions-graph.md,
// docs/toxic-combinations.md) from the resolved data, using stable,
// deterministic ordering (sorted by api_key, then rule id) so weekly CI and
// the quarterly Automation don't produce noisy re-ordering-only diffs when
// they run close together.
//
// TODO(build-render): implement the templates and rendering logic.
package render
