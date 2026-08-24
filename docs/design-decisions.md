# Design decisions (ADR-style)

Short records of deliberate **non-goals** and rejected approaches. Prevents
re-litigating the same ideas without context.

---

## DD-001: Do not modulate severity by install reach (`repository_selection`)

**Date:** 2026-08-24  
**Status:** rejected (for now)

### Proposal

Use where an App is installed (all org repos vs selected repos, and eventually
which repos) to raise or lower reported severity / toxic blast radius.

### Why it sounds attractive

- Blast radius is real: org-wide reach vs one sandbox repo is different.
- Validation showed selected-repo PR bots were over-scored before `contents: write`
  calibration.
- `repository_selection` is objective installation metadata (unlike user install gates).

### Why we rejected it (for now)

1. **Partial overlap with existing rules.** Phase 1 already escalates `all` repos
   (WARN for read-only, HIGH when any write). Further reach-based de-escalation
   adds complexity for marginal gain after catalog calibration.

2. **Selected ≠ low risk.** An App on two production repos is `selected` but may
   be worse than a read-only org-wide scanner. Without repo identity we cannot
   score sensitivity.

3. **Cannot list selected repos with a user token.** Validation confirmed 403 on
   installation repo lists. Modifiers would rely on `selected` alone, which
   hides whether reach is one repo or forty.

4. **Wrong layer for the catalog.** Permission YAML should stay context-free.
   Reach modifiers belong in `gh-app-check` eval if ever added - not in
   `gh-app-graph` severity rows.

5. **Conservative reporting.** Automatic downgrade on `selected` could under-report
   Apps on high-value targets. Better to fix coarse **capability** rules (standalone
   severity, single-grant toxics) than guess reach.

### What we do instead

| Concern | Mechanism |
|---|---|
| Org-wide blast radius | Control-plane rules: all-repos + write → HIGH; read-only → WARN |
| Coarse permission scoring | Standalone severity in `data/permissions/*.yaml` |
| Single-grant false Critical | Remove or avoid toxics when standalone severity suffices (see DD-002) |
| Future reach + sensitivity | Revisit only with installation repo list API or org-owned scanner App |

### Revisit when

- Phase 1 can reliably list installation repositories, **or**
- Org provides a sensitive-repo inventory to cross-walk, **and**
- Retest shows capability calibration alone is insufficient.

---

## DD-002: Single-grant toxic combinations (narrow use)

**Date:** 2026-08-24  
**Status:** accepted constraint

Single-grant toxics (e.g. `organization_administration: write` alone) are allowed
only when standalone permission severity does **not** already tell the story.
Do **not** duplicate standalone High permissions as Critical single-grant toxics
(e.g. removed `supply-chain-poisoning-via-releases-contents`).

See [`calibration-notes.md`](calibration-notes.md).

---

## Related docs

- [`methodology.md`](methodology.md) - toxic combination matching rules
- [`installation-gates.md`](installation-gates.md) - user install authorization (reference only, not scored)
- [`calibration-notes.md`](calibration-notes.md) - org validation learnings
