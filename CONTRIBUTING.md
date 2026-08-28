# Contributing to gh-app-graph

Thank you for helping improve the GitHub App permission risk catalog.

## What belongs here

- Permission severity, dependencies, and toxic combination definitions
  (`data/permissions/`, `data/toxic-combinations.yaml`)
- Methodology docs (`docs/methodology.md`, `docs/design-decisions.md`)
- Build/eval tooling (`cmd/build`, `pkg/eval`, `pkg/data`)

CLI org scanning lives in
[`gh-app-check`](https://github.com/wakeward/gh-app-check).

## Data status

`data/permissions/*.yaml` is **draft seed data** until reconciled with
`octokit/app-permissions`. Read `data/permissions/README.md` before changing
severity or adding permissions.

After editing YAML, run:

```bash
go run ./cmd/build
go test ./...
```

Commit both source YAML and updated `pkg/data/bundled/` artifacts.

## Pull requests

- Open a PR against `main` from a branch.
- Explain the security rationale for severity or toxic-combo changes; link to
  methodology or calibration notes when relevant.
- Do not commit customer org data, private app manifests, or raw spreadsheets
  from `data/sources/raw/` (gitignored).

## Reporting security issues

See [SECURITY.md](SECURITY.md). Do **not** open public issues for
vulnerabilities in the catalog logic or tooling.

## Maintainer note

Solo-maintainer project; issues and PRs welcome.
