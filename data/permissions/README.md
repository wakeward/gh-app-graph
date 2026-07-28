# Draft permission seed data

**Status: DRAFT - not canonical.**

These YAML files are starting notes, not a verified permission catalog. They were
seeded from an incomplete working spreadsheet and manual Enterprise UI
observations. Severity ratings, `impact_plane` assignments, and
`needs_investigation` flags reflect work-in-progress judgment, not reviewed
truth.

The real catalog will be built from:

1. **`octokit/app-permissions`** (endpoint inventory via `cmd/fetch-inventory`) - baseline only; currently ~20 permissions vs many more in GitHub's live docs
2. **Type A dependency detection** (`cmd/detect-overlap`) - endpoints listed under multiple permission categories
3. **Type B prose scraping + human review** (`cmd/scrape-prose`) - additional permissions required per endpoint
4. **Manual curation** - severity and toxic combinations layered on top once the technical graph is solid

Do not use these files for scoring or external publication until they have been
reconciled against the generated inventory. Treat every row as provisional.
