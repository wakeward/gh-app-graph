# Catalog calibration notes (sanitized)

**Status:** implemented (2026-08-24). Informs gh-app-graph before tagging v0.1.0.

No org names or customer identifiers.

---

## Decision: `contents: write` standalone severity

**Standalone permission:** `contents: write` → **High** (was Critical).

**Rationale:**

- Branch protection on default branches blocks direct production pushes; apps must use branches and PRs.
- `contents: write` alone **cannot** modify `.github/workflows/` (requires `workflows: write`).
- Residual High risk: build-script poisoning on feature branches (package.json, Makefile) when CI or developers check out the branch.

**Removed single-grant toxics:**

- `supply-chain-poisoning-via-releases-contents` (contents:write alone)
- `supply-chain-poisoning-via-releases-packages` (packages:write alone; standalone already High)

Critical escalation remains in **multi-grant** toxic combinations:

| Combo | Binding permissions |
|---|---|
| `ransomware-code-destruction` | `contents: write` + `administration: write` |
| `arbitrary-code-execution` | `contents: write` + `workflows: write` |
| `ci-cd-pipeline-takeover-exfiltration-contents-write-workflows-write` | same |
| `autonomous-gatekeeper-bypass-*` | `contents: write` + `pull_requests: write` + `checks:write` / `statuses:write` |

---

## What validated well (2026-08-24 org scan)

| Signal | Observation |
|---|---|
| REST permission mapping | Spot-checked installs match scan rows |
| `administration: write` | CRITICAL control-plane rule matches intuition |
| All-repos read vs write | WARN vs HIGH split works |
| God-mode write count | HIGH when writes > 5 |
| GHES rule exclusion | Cloud scans exclude GHES-only toxics |

---

## Open: `organization-takeover`

Single grant `organization_administration: write` still maps to Critical (GitHub bundles rulesets). Options unchanged; not adjusted in this pass.

---

## Gate before v0.1.0 tag

- [x] `contents: write` standalone → High
- [x] Remove single-grant supply-chain Critical toxics
- [ ] Re-run org scan; confirm internal PR bots no longer CRITICAL from contents alone
- [ ] Decide org-administration single-grant policy (optional)
