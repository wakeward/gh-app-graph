# Catalog calibration notes (sanitized)

**Status:** implemented (2026-08-24, org-admin decision 2026-08-28). Informs gh-app-graph before tagging v0.1.0.

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

## Decision: `organization_administration: write` (single-grant toxic retained)

**Standalone permission:** `organization_administration: write` → **Critical** (unchanged).

**Single-grant toxic:** `organization-takeover` → **retained** (2026-08-28).

Unlike `contents: write`, standalone Critical severity alone does not convey the
full org-wide blast radius. The toxic surfaces explicit exploit paths and
attack-pattern linkage.

**Rationale (Tier 0 / organization-level compromise):**

- **Global ruleset manipulation** - modify, disable, or delete org-wide rulesets;
  strip branch protection across all repositories at once
- **CI/CD hijacking via org runners** - manage organization-level self-hosted
  runners; register a malicious runner to intercept jobs, harvest secrets, or
  inject code into builds
- **Base permission alteration** - change default repository permissions
  org-wide (e.g. Read → Write or Admin), instantly broadening member access
- **Member/owner manipulation** - invite malicious users, promote to Owner,
  demote legitimate admins (original writeup path)

Scoring below Critical would understate the risk: a compromised token or app
with this grant effectively hands over the keys to the org security posture
and CI/CD infrastructure.

---

## Gate before v0.1.0 tag

- [x] `contents: write` standalone → High
- [x] Remove single-grant supply-chain Critical toxics
- [x] Re-run org scan; confirm internal PR bots no longer CRITICAL from contents alone
- [x] Decide org-administration single-grant policy (`organization-takeover` retained)
