# github-app-permissions-graph

[![Unit Testing](https://github.com/wakeward/github-app-permissions-graph/actions/workflows/unit_tests.yml/badge.svg)](https://github.com/wakeward/github-app-permissions-graph/actions/workflows/unit_tests.yml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/wakeward/github-app-permissions-graph/badge)](https://securityscorecards.dev/viewer/?uri=github.com/wakeward/github-app-permissions-graph)
[![License](https://img.shields.io/github/license/wakeward/github-app-permissions-graph)](LICENSE)

A structured data model for GitHub App permission risk: severity, technical
dependencies between permissions, and named "toxic combination" attack
techniques - plus the tooling to keep it fresh and to score a specific app's
declared permissions against it.

**Current status:** `data/permissions/*.yaml` is **draft seed data** only -
seeded from an incomplete working spreadsheet, not a verified catalog. The
canonical permission graph will come from `octokit/app-permissions` plus
dependency detection (`fetch-inventory`, `detect-overlap`). Severity and
toxic-combination judgments get layered on once that technical base is solid.

**Bootstrap vs canonical:** `cmd/migrate-csv` and
`cmd/migrate-toxic-combinations` are temporary importers for working
spreadsheets and notes while the lists are still in flux. The durable source
of truth is the YAML under `data/` (`permissions/`, `toxic-combinations.yaml`,
`dependencies-manual.yaml`) edited and reviewed in-repo. Once curation stabilizes,
expect those migrate commands (and their parsing logic) to be removed rather
than kept as part of the ongoing refresh pipeline. Ongoing automation is
`fetch-inventory` -> `detect-overlap` -> `build` plus human edits to the YAML.

The model is: **structured YAML (source of truth) -> generated views** (Mermaid
graphs, markdown tables, per-app risk reports).

## Why Go, and why a separate repo from `gh-app-check`

This project and [`gh-app-check`](https://github.com/wakeward/gh-app-check)
solve different problems that happen to compose:

| | `github-app-permissions-graph` (this repo) | `gh-app-check` |
|---|---|---|
| **Question it answers** | "What's risky about this permission, and what combinations of permissions are dangerous together?" | "What permissions does this *installed* app actually have in my org, and how is its token actually used?" |
| **Data source** | Public, unauthenticated: GitHub's docs, `octokit/app-permissions`, and `GET /apps/{app_slug}` for any public app's *declared* permissions | Authenticated: `GET /orgs/{org}/installations` for apps *installed* in your own org, plus Code Search for workflow-token tracing |
| **Owns** | The permission catalog, dependency graph, and toxic-combination definitions (severity, exploit paths, blast radius) | Fetching live installation data and running a generic evaluator against it |

The intent is that `gh-app-check` stays lean - it fetches permissions and
runs a generic evaluator - while the actual judgment calls (what counts as
"toxic," how severe a permission is, what the exploit path looks like) live
here as a single, versioned source of truth that `gh-app-check` (or any
other consumer) imports directly:

```go
import (
    "github.com/wakeward/github-app-permissions-graph/pkg/eval"
    "github.com/wakeward/github-app-permissions-graph/pkg/model"
)

result := eval.Evaluate(installationPermissions, toxicCombinations)
```

This is a Go module (not the originally-considered Python project)
specifically so this import works with no serialization boundary: shared
Go types (`pkg/model`), a shared evaluation engine (`pkg/eval`), and
(once built) compiled-in data via `go:embed` (`pkg/data`) - no JSON
round-trip, no risk of the two repos' definitions of "toxic" drifting apart.

## Repository layout

```
cmd/                              # one small binary per pipeline step
  fetch-inventory/                # pull + normalize octokit/app-permissions JSON
  detect-overlap/                 # Type A (category-overlap) dependency detection
  scrape-prose/                   # Type B (prose) dependency drafts, High/Critical scope by default
  build/                          # merge everything, validate, render docs/
  migrate-csv/                    # optional one-time: import a working spreadsheet as draft seed YAML
  migrate-toxic-combinations/     # one-time: import the toxic-combinations CSV + writeup
  evaluate-app/                   # score one app's permissions against toxic-combinations.yaml
pkg/
  model/                          # shared types - import this to work with the data
  eval/                           # toxic-combination matching engine - import this to evaluate a permission set
  graph/                          # dependency merge/resolution + api_key uniqueness validation
  inventory/                      # octokit/app-permissions fetch + Type A detection support
  scrape/                         # Type B prose scraper (goquery)
  render/                         # markdown/Mermaid doc generation
  ghapps/                         # thin client for public GET /apps/{app_slug}
  data/                           # go:embed of the built, resolved JSON - the compiled-in consumption path
data/
  permissions/                    # *.yaml draft seed (see data/permissions/README.md) - NOT canonical yet
  dependencies-manual.yaml        # human-reviewed "requires" edges (highest precedence)
  toxic-combinations.yaml         # named attack techniques unlocked by permission co-occurrence
  generated/                      # endpoints-inventory.json, overlap-dependencies.yaml, scraped-dependencies-draft.yaml
  permissions.resolved.{yaml,json}  # final merged output
docs/
  permissions-table.md            # generated
  permissions-graph.md            # generated (Mermaid)
  toxic-combinations.md           # generated
  methodology.md                  # how severity/blast_radius/confidence are assigned - see below
  quarterly-review-checklist.md   # the fixed non-scrapable checklist the quarterly Automation works through
apps/                             # hand-written permission sets for internal apps, input to evaluate-app
reports/                          # <app-slug>-risk-report.md output of evaluate-app
.github/workflows/
  refresh.yml                     # weekly, cheap auto-detected drift -> PR
  unit_tests.yml, lint_tests.yml, security_scan.yml, scorecard.yml
```

## The two kinds of edge: dependency vs. toxic combination

A **dependency edge** (`data/dependencies-manual.yaml` and the generated
drafts) is a technical prerequisite - an endpoint literally will not work
without the other permission. A **toxic combination**
(`data/toxic-combinations.yaml`) is a risk-correlation rule: two or more
*independently grantable* permissions that, combined, unlock a named attack
technique that neither enables alone. Some toxic combinations overlap with a
technical dependency; most don't. See `pkg/model/dependency.go` and
`pkg/model/toxic.go` for the exact schemas.

## Methodology

`docs/methodology.md` documents how `Severity`, `BlastRadius`, and
`Confidence` are assigned - required reading before trusting (or disputing)
a rating, and a hard requirement if this data is ever consumed outside this
repo's own tooling.

## Development

```bash
go build ./...
go vet ./...
go test ./...
```

## License

Apache License 2.0 - see [LICENSE](LICENSE). This project consumes the
[`octokit/app-permissions`](https://github.com/octokit/app-permissions)
generated inventory as its base endpoint data; see `pkg/inventory` for
attribution details once the fetch step is implemented.
