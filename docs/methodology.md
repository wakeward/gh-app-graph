# GitHub App permission methodology (draft)

**Status:** draft - started Aug 2026. Not yet canonical. Requires theorizing
session sign-off ([#3](https://github.com/wakeward/gh-app-graph/issues/3)).

This document explains how Wakeward scores GitHub App **capabilities** in
`gh-app-graph` and what `gh-app-check org` reports today. It separates
capability (impact) from likelihood (actor path, enterprise controls).

For planning context see [`threat-model-plan.md`](threat-model-plan.md). For
who can install an App at enterprise, org, or repository level see
[`installation-gates.md`](installation-gates.md). For named attack paths with
actor context see [`data/attack-patterns.yaml`](../data/attack-patterns.yaml).

---

## What we model

| Layer | Question | Where | Used by |
|---|---|---|---|
| Standalone severity | How bad is one permission at read vs write? | `data/permissions/*.yaml` | Catalog, future marketplace eval |
| Technical dependency | What grant is required for an API to work? | `dependencies-manual.yaml`, generated overlap | Overlap flags on toxic combos |
| Toxic combination | What named technique do co-granted permissions enable? | `data/toxic-combinations.yaml` | `gh-app-check org` |
| Attack pattern | Who enables it, how does it land, how detectable? | `data/attack-patterns.yaml` | Human review, blog, future likelihood work |
| Control-plane rules | Structural red flags (all repos, admin write, god-mode)? | `gh-app-check/pkg/rules` | `gh-app-check org` |

**Out of scope for this program:** OAuth Apps, user PATs, gh CLI OAuth tokens,
VSCode extensions/tasks, and any vector that does not involve a GitHub App
manifest, installation, or installation token.

**Out of scope for Phase 1 scoring:** likelihood scoring, publisher trust ratings, runtime behavior analysis, user install authorization (see below).

### Installation gates (reference only - not assessment)

Who can approve an install depends on user role and org policy (enterprise
owner, organization owner, or repository admin under strict conditions). See
[`installation-gates.md`](installation-gates.md).

This material supports **blog posts and manual likelihood reasoning**. It is
**not** parsed, scored, or emitted by `gh-app-check`. The tool reports what
permissions Apps have, not which user role approved the install.

### Credential lifecycle and ecosystem context (reference only)

How PEM, JWT, IAT, and U2S tokens are minted and abused is documented in
[`credential-lifecycle.md`](credential-lifecycle.md). Marketplace trust
limitations and third-party research synthesis live in
[`marketplace-trust-limitations.md`](marketplace-trust-limitations.md) and
[`references/ecosystem-research-synthesis.md`](references/ecosystem-research-synthesis.md).
None of this changes Phase 1 scoring.

---

## Standalone permission severity

Each permission has per-access-level severity:

| Severity | Meaning |
|---|---|
| Critical | Direct path to org/repo takeover, secret overwrite, or irreversible data loss |
| High | Significant abuse potential; often needs pairing or context to exploit fully |
| Medium | Material misuse or disruption; reconnaissance with follow-on potential |
| Low | Limited abuse; nuisance or narrow reconnaissance |
| Informational | No meaningful standalone risk (e.g. public metadata) |
| Unknown | Not yet assessed - always paired with `needs_investigation: true` |

Severity answers: *if this scope alone were granted and exercised, how bad is the worst plausible outcome?*

It does **not** answer: *how likely is an org to grant it?*

### Escalation example: `contents: write`

Standalone `contents: write` is **High** (branch protection and the
workflows/contents boundary limit direct pipeline takeover). It becomes
**Critical** only via toxic combinations that add bypass permissions
(`administration: write`, `workflows: write`, or PR/check forgery triples).
See [`calibration-notes.md`](calibration-notes.md).

### Impact plane

Each permission is classified into one impact plane:

- **control** - settings, branch protection, RBAC, network config
- **data_execution** - code, CI, artifacts, secrets at runtime
- **governance_identity** - membership, credentials-as-identity, directory data

### Platform availability

| Value | Meaning |
|---|---|
| `all` (default) | GitHub.com, Enterprise Cloud, Enterprise Server |
| `ghes_only` | Self-hosted Enterprise Server only (e.g. pre-receive hooks) |

`gh-app-check org` excludes `ghes_only` toxic rules on cloud scans unless
`--platform ghes` is set.

### Doc status

| Value | Meaning |
|---|---|
| `documented` | Confirmed in public docs and/or live UI |
| `undocumented_preview` | Visible in UI as Preview |
| `unconfirmed_key` | Inferred api_key - verify before relying on automated rules |
| `disputed` | Conflicting signals (schema vs UI vs docs) |

Permissions with `needs_investigation: true` must not drive hard CI fail gates
until reviewed.

---

## Toxic combinations

A toxic combination is a **named attack technique** an installation enables when
**every** listed permission grant is present on that installation.

### Matching rule (AND, not OR)

Each toxic entry lists one or more `(api_key, access)` pairs. The engine checks
**all** of them against the installed permission map:

```
Match  ⇔  grant₁ present  AND  grant₂ present  AND  …
```

- If the combo lists `contents: write` and `workflows: write`, **both** must be
  granted (as write) for a match. `contents: write` alone does **not** fire that combo.
- A granted **write** satisfies a combo requirement for **read** on the same key
  (write is strictly stronger).
- There is no partial credit: missing one grant means no match (but may produce a
  near miss - see below).

Implementation: `pkg/eval/evaluate.go` calls `HasGrant` for each required grant.

### Single-grant vs multi-grant entries

| Type | When it fires | Near miss? |
|---|---|---|
| **Multi-grant** (most entries) | All listed grants present | Yes - if exactly one grant missing |
| **Single-grant** (rare) | That one grant present | No - either matches or does not |

Most toxics require **two or more** permissions because the technique needs
co-granted capabilities (e.g. forge checks **and** push code **and** open PRs).

Single-grant entries are reserved for cases where reviewers might accept a scope
without understanding what it unlocks, and standalone permission severity does
**not** already capture it. We removed single-grant Critical entries where
standalone severity suffices (e.g. `contents: write` - see
[`calibration-notes.md`](calibration-notes.md), [`design-decisions.md`](design-decisions.md)).

### Near misses

For **multi-grant** combos only: if the installation is missing **exactly one**
required grant, the combo appears as a near miss ("one permission away from X").
This is not a toxic match - no blast radius from that combo is applied.

Single-grant combos do not generate near misses.

### What toxics do not consider

Toxic matching is **permissions only**. It does not use:

- `repository_selection` (all vs selected)
- Which repos are selected
- Who approved the install
- Publisher trust or runtime behavior

Control-plane rules (all-repos, administration write, god-mode write count) apply
**in addition** to toxic matches.

Fields:

| Field | Purpose |
|---|---|
| `technique` | Human-readable name shown in auditor output |
| `permissions` | Required `(api_key, access)` grants |
| `blast_radius` | Critical / High / Medium - worst-case outcome if exploited |
| `exploit_path` | Short narrative of the abuse chain |
| `platform_availability` | Optional; inherits from permission grants when omitted |
| `overlaps_technical_dependency` | True when combo includes a technical prerequisite pair |

Toxic combinations are **capability** statements. They do not encode actor
preconditions (org owner approval, compromised publisher, etc.).

---

## Control-plane rules (`gh-app-check`)

Phase 1 applies structural predicates in addition to toxic combos:

| Rule | Default risk | Rationale |
|---|---|---|
| All repositories + any write | HIGH | Maximum blast radius |
| All repositories, read-only only | WARN | Broad reconnaissance surface |
| `administration: write` | CRITICAL | Repo lifecycle and branch protection control |
| More than 5 write scopes | HIGH | God-mode portfolio |

These are impact shortcuts, not replacements for toxic technique naming.

---

## Capability vs likelihood

```
Impact     ≈ f(granted permissions, repository_selection, reachable secrets)
Likelihood ≈ f(initial_access_path, enterprise_controls, publisher trust, detectability)
```

| Tool | Measures today |
|---|---|
| `gh-app-check org` | Impact (what is installed) |
| Attack patterns catalog | Likelihood context (draft) |
| Phase 2 `trace` (planned) | Execution-plane preconditions |
| Phase 4 drift guard (planned) | Permission change over time |

Do not treat PASS or absence of toxic matches as "safe." Phase 1 is a
control-plane snapshot only.

---

## Validation ladder

Evidence quality for claims in this repo:

| Level | Proves | Typical needs |
|---|---|---|
| L0 | Docs/API mapping | Public documentation |
| L1 | Single permission to endpoint | Test org |
| L2 | Multi-step toxic combo end-to-end | Org owner test org |
| L3 | Ephemeral change / hiding behavior | Audit log (often Enterprise) |
| L4 | Enterprise-only permissions / approval workflows | GHEC/GHES trial org |
| L5 | Execution-plane token abuse | Code Search, Actions (Phase 2) |

Toxic combinations and permissions should cite validation level in desk traces
when promoted from draft seed to reviewed.

---

## Promotion criteria (draft seed → reviewed)

Before treating a permission row or toxic combo as canonical:

1. api_key confirmed against live App permissions UI or manifest
2. Severity reviewed for read and write independently
3. At least one realistic initial-access path documented in attack patterns
4. For Critical blast radius: desk trace started or incident reference cited
5. `needs_investigation` cleared or explicitly deferred with reason

---

## Explicit non-goals (current phase)

- OAuth Apps, PATs, gh CLI, IDE extensions (adjacent attack surface, not this catalog)
- Encoding actor preconditions as fake permissions in toxic YAML
- Automated likelihood scoring in `gh-app-check`
- Malicious publisher runtime behavior detection (static manifest/install only)

---

## References

- Threat model plan: [`threat-model-plan.md`](threat-model-plan.md)
- Attack patterns (draft): [`data/attack-patterns.yaml`](../data/attack-patterns.yaml)
- Toxic combinations (generated): [`toxic-combinations.md`](toxic-combinations.md)
- Credential lifecycle: [`credential-lifecycle.md`](credential-lifecycle.md)
- Marketplace trust: [`marketplace-trust-limitations.md`](marketplace-trust-limitations.md)
- Research synthesis: [`references/ecosystem-research-synthesis.md`](references/ecosystem-research-synthesis.md)
- Design decisions: [`design-decisions.md`](design-decisions.md)
- Control-plane auditor: [`gh-app-check`](https://github.com/wakeward/gh-app-check)
