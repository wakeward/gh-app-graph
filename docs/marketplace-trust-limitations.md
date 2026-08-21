# Marketplace trust limitations

**Status:** draft reference (Aug 2026). Likelihood and blog material. Not
scored by `gh-app-check`.

Administrators often treat Marketplace listing as a security endorsement.
This document captures **structural limitations** that affect real-world
likelihood of over-trusted installs. It does not change permission severity
scores.

---

## Verification paradox

Marketplace verification emphasizes publisher identity hygiene (domain TXT,
verified email, 2FA on publisher org). It does **not** guarantee:

- Secure backend handling of PEM, webhook secrets, or OAuth callbacks
- Least-privilege manifest design
- Absence of transitive dependency risk in publisher infrastructure
- Immutable release provenance for the App's supporting services

**Likelihood effect:** busy owners may conflate "verified badge" with "safe to
grant `contents: write` org-wide." Attack pattern:
`initial-access-marketplace-verification-paradox`.

**Enterprise mitigations:** app allow lists, mandatory permission review,
`gh-app-check org` baseline (capability only, not trust rating).

---

## Permission granularity ceiling

GitHub App permissions are **repository-level scopes**, not path-level ACLs.

| Requested scope | What GitHub enforces | What admins may wish for |
|---|---|---|
| `contents: read` | Entire repository history and tree | Read only `/src` |
| `contents: write` | Any file including `.github/workflows/` | Annotate PRs without workflow write |
| `actions: write` | Workflows and workflow runs in granted repos | Trigger one named workflow |

A vendor needing PR annotations must often request broad `contents` and
`pull_requests` writes. **Impact** is unchanged in our model (the grant is still
dangerous); **likelihood of over-provisioning** is higher because least privilege
is structurally awkward.

This reinforces why toxic combinations and near-miss warnings matter even for
"legitimate" SaaS categories (SAST, docs bots, release automation).

---

## Publisher compromise cascade

If a widely installed App's publisher account or PEM store is compromised,
every installation's **declared permissions** become immediately exercisable
by the attacker. Phase 1 audit shows the same risk before and after compromise
(`execution-abuse-declared-scopes`, S2b).

Mutable deployment of vendor **Actions** (separate from GitHub Apps) caused
large-scale CI incidents (e.g. `tj-actions/changed-files`). That is adjacent
surface - Actions supply chain, not App manifest - but the **blast radius
narrative** (silent update, thousands of repos) informs enterprise app policy.

---

## Complementary tooling (not replacements)

| Tool | Plane | Relationship to gh-app-check |
|---|---|---|
| **gh-app-check org** | Control plane (installed App permissions) | Canonical for permission/toxic snapshot |
| **Zizmor** | CI workflow static analysis | Finds PRT misuse, unpinned Actions; complements Phase 2 trace |
| **Legitify / Scorecard** | Org/repo policy posture | Branch protection, member privileges |
| **Harden-Runner** | Runtime egress on runners | Blocks exfil after token theft |

Use each in its layer. None substitutes for reviewing what Apps are installed.

---

## References

- Attack patterns: [`data/attack-patterns.yaml`](../data/attack-patterns.yaml)
- Red Team series mapping: [`attack-patterns/red-team-series-mapping.md`](attack-patterns/red-team-series-mapping.md)
- Installation gates: [`installation-gates.md`](installation-gates.md)
