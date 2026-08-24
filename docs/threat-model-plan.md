# GitHub App threat model - planning document

**Status:** planning only - methodology draft started; tooling unchanged until reviewed.
**Owner:** Wakeward (human judgment required before any ratings are treated as canonical).
**Tracking:** [gh-app-graph#3](https://github.com/wakeward/gh-app-graph/issues/3)

This document captures agreed direction from design discussion (Aug 2026). It
deliberately separates **capability** (what installed permissions enable) from
**likelihood** (who can get those permissions in front of real repos, and through
which gates). Implementation in YAML, `methodology.md`, lab traces, and
`gh-app-check` scoring comes *after* this planning work.

---

## Goals

1. Document Wakeward's perspective on **standalone permission risk** (severity at
   read vs write, impact plane, security notes).
2. Define what **additional factors** an adversary needs beyond a single
   permission - technical co-requisites, toxic combinations, actor gates, and
   execution-plane preconditions.
3. Build a **scenario catalog** covering ~90% of common GitHub App attack
   vectors, with emphasis on **initial access** and **permission change**
   paths (not only "what the API allows once installed").
4. **Validate** high-impact chains beyond theory - step-by-step desk/lab traces,
   including investigation of which GitHub license tiers are required.
5. Only then publish **`docs/methodology.md`** and promote draft YAML from
   "spreadsheet seed" to reviewed truth.

---

## What exists today (do not confuse with finished threat model)

| Layer | Question | Location | Maturity |
|---|---|---|---|
| Standalone severity | How bad is this permission alone? | `data/permissions/*.yaml` | Draft seed |
| Technical dependency | What else must be granted for an endpoint to work? | `dependencies-manual.yaml`, generated overlap/scrape | Partial |
| Toxic combination | What named technique unlocks when permissions co-occur? | `data/toxic-combinations.yaml` (21 entries) | Curated; review in progress |
| Control-plane snapshot | What is installed in an org right now? | `gh-app-check` Phase 1 | Implemented |
| Actor / likelihood / validation | Who enables this, how likely, proven how? | **Not modeled yet** | This plan |

---

## Core design principle

```
Impact     ≈ f(blast_radius, repository_selection, reachable secrets/data)
Likelihood ≈ f(initial_access_path, enterprise_controls, detectability)
```

**`gh-app-check org` ranks impact (capabilities installed).** It does not rank
likelihood and does not assess user install authorization. Install gates live in
[`installation-gates.md`](installation-gates.md) for blog and manual review only.

Three different meanings of "additional permissions required to exploit":

1. **Technical co-requisite** - API will not succeed without permission X.
   Model as dependency edges.
2. **Technique co-requisite** - permissions are independently grantable; together
   they unlock an attack. Model as toxic combinations.
3. **Actor co-requisite** - not a GitHub App permission at all (org owner approval,
   ability to modify installation, compromised publisher, leaked app key, workflow
   that mints installation tokens). Model in scenario / pattern docs - **do not**
   encode as fake `api_key` entries in toxic YAML.

---

## Draft scenarios (non-exhaustive - from initial brainstorm)

These four scenarios came from an informal brainstorm. They are **starting
examples**, not a complete catalog. A dedicated theorizing session is required
before treating any list as ~90% coverage.

### S1 - App author without install rights (approval-gated)

| Field | Notes |
|---|---|
| **Actor** | Developer (or compromised dev) who can create/edit an app definition |
| **Path** | Register app with elevated requested permissions → persuade or routine-approve an org owner → installation lands |
| **Enabling event** | Installation approval (human/process gate) |
| **Install gates** | Reference only (blog, likelihood reasoning). See [`installation-gates.md`](installation-gates.md). Not in `gh-app-check` output. |
| **Capability vs likelihood** | Manifest shows intent; impact is latent until install |
| **Related toxic combos** | Any combo - depends on what was approved |
| **Detection (investigate)** | `integration_installation.*` audit events; who approved |
| **Validation** | Not yet run |

### S2 - Third-party supply chain

Two sub-variants to keep separate:

**S2a - Permission escalation by publisher**

| Field | Notes |
|---|---|
| **Actor** | External app publisher |
| **Path** | Benign install → publisher updates manifest → org owner re-approves expanded permissions |
| **Enabling event** | Permission upgrade approval |
| **Related tooling** | Phase 4 drift guard (future) |
| **Validation** | Not yet run |

**S2b - Malicious behavior within existing grants**

| Field | Notes |
|---|---|
| **Actor** | External app publisher (or compromised publisher infrastructure) |
| **Path** | Install within declared permissions; backend abuses granted scopes |
| **Enabling event** | None (permissions unchanged) |
| **Gap** | Phase 1 static audit will not catch; needs behavior trust, execution-plane trace (Phase 2), or external reputation |
| **Validation** | Not yet run |

### S3 - Temporary permission change / hide activity (insider or compromised admin)

| Field | Notes |
|---|---|
| **Actor** | Org owner, app manager, or compromised high-privilege account |
| **Path** | Temporarily widen permissions or use already-wide grants → act → revert settings to reduce visibility |
| **Example technique** | `ransomware-code-destruction` (`administration:write` + `contents:write`) - branch protection off, push, protection restored |
| **Enabling event** | App permission change and/or repository settings change |
| **Detection (investigate)** | Correlate branch protection events with app permission updates |
| **Validation** | Not yet run |

### S4 - Spray many apps / mask in noise

| Field | Notes |
|---|---|
| **Actor** | Party with repeated install rights, or decentralized teams |
| **Path** | Many low-noise apps, each narrowly scoped; collective exposure is large |
| **Model note** | Governance / portfolio effect - not a new per-permission severity score |
| **Gap** | Org-level aggregation; inventory + drift over time |
| **Validation** | Not yet run |

### S5 - Execution plane: workflow and agent abuse

| Field | Notes |
|---|---|
| **Actor** | Fork PR author, external GitHub App publisher, or prompt injection via issue body |
| **Path** | Installed permissions unchanged; IAT minted in CI and exfiltrated (PRT + fork checkout), or agent workflow trusts `[bot]` actors |
| **Enabling event** | Workflow run, not new install |
| **Gap** | Phase 1 static audit unchanged; Phase 2 `trace` and workflow static analysis (Zizmor) |
| **Patterns** | `credential-access-prt-fork-iatt-exfiltration`, `execution-ai-agent-external-bot-trust`, `execution-indirect-prompt-injection-agent` |
| **Validation** | GHSA-9g93-rxr5-xhqw field observed; agent paths desk trace pending |

### S6 - Vendor integration and infrastructure

| Field | Notes |
|---|---|
| **Actor** | Unauthenticated callback attacker, SSRF to webhook handler, DNS hijacker |
| **Path** | Victim org installs App correctly; vendor binding or callback endpoint is poisoned |
| **Enabling event** | OAuth callback, webhook delivery, or dangling DNS on integration subdomain |
| **Gap** | Outside org install audit; reference for blog and vendor review checklists |
| **Patterns** | `credential-access-vendor-install-callback-binding`, `credential-access-webhook-handler-ssrf-jwt`, `persistence-dangling-dns-oauth-webhook` |
| **Validation** | CVE references cited; verify before external publication |

---

## Planned deliverable: GH Apps ATT&CK-style catalog

Working title: **GitHub App Attack Patterns** (or similar - name TBD in
theorizing session).

Purpose: a structured list analogous to MITRE ATT&CK, but scoped to GitHub Apps
and adjacent execution paths (Actions tokens, installation tokens, webhooks).
Emphasis on **initial access** and **permission acquisition** tactics, not only
post-install impact.

### Proposed tactic columns (draft - revise in session)

| Tactic | Question it answers | Example techniques (placeholder) |
|---|---|---|
| **App registration** | How does the adversary define what permissions are *requested*? | Create app, fork legitimate app manifest, typosquat app name |
| **Initial access to org** | How do requested permissions become an *installation*? | Owner approval, selected-repo admin install, pre-approved app allow-list abuse |
| **Permission acquisition** | How do grants expand after install? | Publisher manifest update, owner click-through re-approval, misconfigured broad `repository_selection: all` |
| **Credential access** | How does the adversary obtain usable tokens? | Leaked PEM/private key, over-scoped workflow `GITHUB_TOKEN`, stolen installation token from logs |
| **Execution** | How is capability exercised? | REST API abuse, malicious workflow dispatch, webhook-driven pipeline |
| **Persistence** | How is access maintained? | Long-lived installation, hidden workflow, secondary app |
| **Defense evasion** | How is abuse hidden? | Ephemeral permission/settings change (S3), high-volume low-severity apps (S4), inactive app suddenly activated |
| **Impact** | What is the outcome? | Org takeover, supply chain artifact swap, secret exfil, CI compute abuse |

Each technique entry should eventually include:

- `technique_id` (stable slug)
- `tactic`
- `description`
- `app_permissions[]` (links to toxic combo IDs where applicable)
- `actor_preconditions[]`
- `execution_preconditions[]`
- `enterprise_controls_that_reduce_likelihood[]`
- `observable_artifacts[]` (audit log event names - **to be verified per license**)
- `related_scenarios[]` (S1-S4 mapping)
- `validation_status`: `theoretical` | `desk-trace-draft` | `lab-reproduced` | `field-observed`
- `license_notes` (what GH tier was required to reproduce or observe logs)

**Relationship to existing data:**

- Toxic combinations remain the **capability** layer (what co-granted permissions enable).
- Attack patterns add **context** (who, when, how it lands, how likely, how to detect).
- Do not duplicate exploit prose in three places - cross-link by ID.

---

## Theorizing session (required before claiming ~90% coverage)

### Session outcome

A reviewed table of attack patterns with:

- No major obvious gaps in **initial access** and **permission change** paths
- Explicit "out of scope / rare" bucket for edge cases
- Mapping from each high-impact toxic combo to at least one realistic initial-access path
- Priority list for lab validation (top N by impact × plausible likelihood)

### Suggested agenda (half day)

1. **Frame** - capability vs likelihood split; review S1-S4; agree tactic taxonomy.
2. **Initial access brainstorm** - whiteboard paths to first installation (internal dev, vendor, compromised owner, stolen credentials, marketplace/social engineering).
3. **Permission change brainstorm** - upgrades, repo selection changes, suspension/reactivation, transfer of app ownership.
4. **Execution plane** - overlap with Phase 2 (`gh-app-check trace`): workflow misuse, PEM in repo, OIDC/workload identity (if in scope).
5. **Enterprise controls** - what reduces likelihood per path (app approval policies, SAML, audit log streaming, allow-listing).
6. **Gap check** - compare against known public incidents, GitHub security docs, and existing toxic combos; mark missing patterns.
7. **Prioritize validation** - pick 3-5 patterns for first lab traces; note license questions.

### Seed prompts (not exhaustive)

- Compromised org owner approves malicious install
- Compromised repo admin installs on selected repos that hold production secrets
- Legitimate CI app with `actions: write` used to pivot
- Webhook secret leakage enabling forged deliveries (if in scope)
- GitHub App converted from OAuth App / confusion between credential types
- Suspended app reactivated after long dormancy
- App transferred to new owner (publisher change)
- Marketplace listing vs direct install URL
- Internal app with overly broad `repository_selection: all` drift
- Two apps whose **combined** portfolio exceeds any single review threshold (S4)

---

## Validation program (after theorizing)

### Desk trace template (per pattern)

Create `docs/desk-traces/<technique-id>.md` when ready:

1. Preconditions (org config, app permissions, actor role)
2. Numbered attack steps (API calls or UI actions)
3. Expected observable artifacts
4. Cleanup / rollback
5. License tier used
6. Status and date

### Validation ladder

| Level | Proves | Typical GH needs |
|---|---|---|
| L0 | Docs/API mapping | Public documentation |
| L1 | Single permission → endpoint | Test org, any plan |
| L2 | Multi-step toxic combo end-to-end | Org owner test org |
| L3 | Ephemeral change / hiding behavior | Audit log access (often Enterprise) |
| L4 | Enterprise-only permissions / approval workflows | GHEC/GHES (trial org TBD) |
| L5 | Execution-plane token abuse | Code Search, Actions (Phase 2 scope) |

### License investigation spike (open)

Before lab work, confirm for company and trial orgs:

- [ ] Audit log API availability and event names for app install / permission change
- [ ] GitHub App approval policies - plan requirements
- [ ] Third-party app restrictions
- [ ] Enterprise-only permission picker entries (quarterly checklist step 6)
- [ ] Whether a dedicated **test Enterprise org** is approved for destructive tests

---

## Deliverables checklist (order matters)

- [x] **This plan** reviewed and agreed
- [ ] **Theorizing session** held; draft ATT&CK-style catalog populated (seed + Red Team series mapping done; session still required for gap-check)
- [x] **`docs/methodology.md`** - severity, blast radius, confidence, actor/likelihood rules (draft started)
- [x] **`data/attack-patterns.yaml`** + `docs/attack-patterns/` - structured catalog seed cross-linked to toxic combos
- [ ] **Permission YAML review** - Wakeward perspective merged from spreadsheet; `needs_investigation` cleared where possible
- [ ] **Desk traces** for top-priority patterns (validation_status updated)
- [ ] **Company org calibration** - run `gh-app-check org` and compare to intuition (separate milestone; does not replace lab proof)
- [ ] **GitHub Project items** - break out session, methodology, validation spikes as trackable issues

---

## Explicit non-goals (for now)

- Extending `pkg/model` or `gh-app-check` scoring for likelihood
- Adding actor preconditions into toxic-combination YAML as permissions
- Claiming the four draft scenarios are complete
- Publishing external/blog content before methodology + at least one lab trace exist

---

## Open questions

1. Final name for the ATT&CK-style catalog?
2. Single YAML file vs markdown per technique for attack patterns?
3. Include OAuth Apps / PAT abuse adjacent paths, or strict GitHub Apps scope only?
4. Include GitHub Apps on **GitHub Enterprise Server** differences?
5. Who attends the theorizing session (solo vs peer review)?

---

## References

- Existing toxic combinations: `data/toxic-combinations.yaml`, generated `docs/toxic-combinations.md`
- Permission draft seed: `data/permissions/README.md`
- Quarterly human review: `docs/quarterly-review-checklist.md`
- Control-plane auditor: [`gh-app-check`](https://github.com/wakeward/gh-app-check) Phase 1
- Drift guard backlog: [`gh-app-check/docs/BACKLOG.md`](https://github.com/wakeward/gh-app-check/blob/main/docs/BACKLOG.md)
- Prior art (initial access): [`docs/attack-patterns/red-team-series-mapping.md`](attack-patterns/red-team-series-mapping.md) and the [Red Teaming GitHub](https://wakeward.uk/tags/github/) series
