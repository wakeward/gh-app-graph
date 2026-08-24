# Catalog calibration notes (sanitized)

**Status:** post org-validation learnings (2026-08-24). Informs gh-app-graph
toxic rule tuning before tagging v0.1.0 or dropping `gh-app-check`'s `replace`
directive.

No org names or customer identifiers. Source: read-only Phase 1 validation on
Enterprise Cloud (16 installations).

---

## What validated well

| Signal | Observation |
|---|---|
| REST permission mapping | Spot-checked installs: API keys and read/write levels match scan rows |
| `administration: write` | CRITICAL control-plane rule matches intuition |
| All-repos read vs write | WARN vs HIGH split works (read-only metrics App vs StepSecurity-style writes) |
| God-mode write count | HIGH when writes > 5; apps with exactly 5 writes rely on other rules |
| GHES rule exclusion | `server-side-remote-code-execution` excluded on cloud scans |
| Runner poisoning toxic | Useful on IT/ops-style Apps with self-hosted runner write |

---

## Severity too hot (catalog, not tool bug)

### `supply-chain-poisoning-via-releases-contents`

**Current:** single grant `contents: write` → Critical blast radius.

**Problem:** Every PR bot, design-review bot, or committer App on **selected**
repos becomes CRITICAL. Capability is real (releases/tags are contents scope)
but unconditional Critical is wrong for triage on narrow internal bots.

**Options (pick one in theorizing session):**

1. Require `contents: write` **and** `repository_selection: all` for Critical
2. Downrank single-grant to **High** when `selected` repos only
3. Split combo: `contents: write` alone → High; add `packages: write` or
   `actions: write` companion for Critical
4. Keep Critical but add `gh-app-check` context flag (not severity change)

### `organization-takeover`

**Current:** single grant `organization_administration: write` → Critical.

**Problem:** GitHub bundles org rulesets into org administration. A ruleset-only
automation App looks like full org takeover.

**Options:**

1. Downrank to HIGH when org-admin write is the **only** write grant
2. Add methodology note: reviewers must read GitHub's bundled permission meaning
3. Near-miss only (weaker - loses education value)

---

## Near misses

Almost every App with `contents: read` surfaces the same hook/variable near-miss
cluster. Useful for education; too noisy as a default table column.

**Mitigation (gh-app-check):** `--no-near-misses` flag (implemented). Consider
default off for table format only.

---

## Out of scope for catalog

| Topic | Why |
|---|---|
| Selected repo list in output | User token cannot list installation repos without App token |
| Parent org 404 | Org membership role, not OAuth scope |
| Private App name 404 on GET /apps/{slug} | Expected for unpublished Apps |

---

## Gate before v0.1.0 tag

- [ ] Decide `contents: write` single-grant severity policy
- [ ] Decide org-administration single-grant policy
- [ ] Re-run org scan; confirm CRITICAL count aligns with triage expectations
- [ ] Update `docs/methodology.md` promotion criteria if rules change
