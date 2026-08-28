# GitHub App Security program - long-term vision

**Status:** planning and direction - not committed implementation.
**Owner:** [Wakeward](https://wakeward.uk) (Kevin Ward).
**Prerequisite work:** threat model ([#3](https://github.com/wakeward/gh-app-graph/issues/3)),
`gh-app-graph` curation, and `gh-app-check` org validation must land first.

This document captures the long-term program beyond the current repos. It
complements [`docs/threat-model-plan.md`](threat-model-plan.md) (how we judge
risk) with **what we build and publish** once that judgment is defensible.

---

## Program in one diagram

```text
gh-app-graph (rules + catalog knowledge)
      │
      ├── gh-app-check        → installed apps in *your* org (today)
      ├── gh-app-catalog      → *all* Marketplace apps, public manifests (future)
      │
      ├── Lane A → GitHub platform security research (responsible disclosure)
      ├── Lane B → suspicious Marketplace app signals → triage → partners
      └── Lane C → wakeward.uk blog (methodology, tooling, research writeups)

Product site (domain TBD, e.g. gh-app-sec.*) → daily static pages, public diffs
Personal blog (wakeward.uk)                     → attribution, depth, trust
```

---

## Near-term focus (do this first)

1. Complete threat model work ([#3](https://github.com/wakeward/gh-app-graph/issues/3)) -
   theorizing session, attack-pattern catalog, `docs/methodology.md`, lab traces.
2. Harden **`gh-app-graph`** - promote YAML from draft seed to reviewed catalog;
   tag a release so consumers can drop local `replace` directives.
3. Validate **`gh-app-check org`** against a real company org - calibrate rules
   and toxic combos against lived experience (sanitized notes only in public docs).
4. **`gh-app-check` Phase 2-3** as needed - execution-plane trace, SARIF/strict CI.

Nothing in the marketplace catalog or product site blocks (1)-(3). Those are
foundations for credible public output.

---

## Pillar 1 - Marketplace catalog and public drift

**User question answered:** For every app on the
[GitHub Marketplace (Apps)](https://github.com/marketplace?type=apps), what
permissions does it *request*, how does that score against our rules, and **what
changed since yesterday?**

This is **not** the same as `gh-app-check org`:

| | Marketplace catalog | `gh-app-check org` |
|---|---|---|
| Data | Public `GET /apps/{slug}` manifest | Authenticated org installations |
| Audience | Anyone considering an install | Your org's security team |
| Drift | Publisher changes requested permissions | Your org's approved grants change |

### Full crawl

- Discover the complete Marketplace app slug set daily (paginated listing crawl;
  persist discovery metadata).
- Fetch each app's public manifest; normalize permissions to a canonical form.
- Evaluate against embedded `gh-app-graph` rules (standalone severity, toxic
  combinations).
- Publish static pages (see **Product site** below).

Tracking issue: [#4](https://github.com/wakeward/gh-app-graph/issues/4).

### Hash-first diff (fast daily scans)

Avoid re-running heavy evaluation and page generation when nothing changed.

**Canonical fingerprint per app:**

```text
permissions_hash = SHA-256(normalized JSON of:
  sorted map of api_key → read|write
  + repository_selection (if present)
  + optional future fields: events[], etc.)
```

**Daily pipeline tiers:**

| Tier | Action |
|---|---|
| Discover | Refresh full slug list from Marketplace |
| Fingerprint | Fetch manifests in parallel; compute `permissions_hash` |
| Deep | Toxic eval + page regen **only** when hash changed or app is new |

**Eval cache:** results keyed by `permissions_hash` so apps sharing identical
permission sets reuse the same rule output.

**Snapshot store:** append-only history (JSONL per day or SQLite) enabling
time-series diffs and static site generation from latest + delta.

### Public drift alerts (policy)

**Drift is public by default.** Permission changes are surfaced on the site as
awareness, **not** as accusations of malice.

Framing for users:

- "This app's **requested permissions changed** on {date}."
- "If you already installed it, your org may be prompted to re-approve."
- "Review whether you still need this app and whether the new scopes match its purpose."

**Not in scope for drift copy:** labeling an app malicious, vendor shaming, or
CVE-style severity for permission additions alone.

A machine-readable **feed** (RSS/Atom/webhook) is deferred. The product goal is
to bring people to the **site** first; a feed can follow once diff quality and
storm handling are proven.

### Storm conditions (platform-wide permission floods)

GitHub occasionally introduces permissions that many or all apps must adopt
(e.g. a new baseline scope). A naive diff engine would report "every app
changed" and overwhelm users.

**Design requirements:**

1. **Correlated-change detection** - if more than a configurable threshold
   (e.g. 30-50%) of apps gain the same permission on the same day, classify as a
   **platform-wide landscape event**, not N individual emergencies.
2. **Prominent site banner** - e.g. "GitHub-wide permission change detected on
   {date}: {permission} added across {N} apps - likely a platform requirement,
   not individual app behavior."
3. **Per-app truth preserved** - still record each app's diff in the data store;
   de-emphasize per-app alert styling on the homepage during a storm.
4. **Changelog correlation** - manual or semi-automated link to GitHub
   changelog / docs when a storm is recognized; optional `storm_event` records
   in snapshot metadata.
5. **Rate-limited homepage** - cap how many individual "changed yesterday" items
   appear above the fold; full list remains on a dedicated daily diff page.

Legitimate broad change and suspicious single-app escalation should look
different in the UX.

### Product site and attribution

**Domain:** undecided; candidate dedicated domain (e.g. `gh-app-sec.*`) for the
data product. Register when ready; implementation can start on static paths
under wakeward.uk if needed.

**Attribution:** every catalog page links clearly to
[Wakeward / wakeward.uk](https://wakeward.uk) as author of the research and
tooling. The personal blog remains the voice for methodology and platform
research; the product site is the daily reference.

**Static generation:** daily job produces HTML (Hugo/Blowfish or similar) -
per-app pages, daily "what changed" summary, methodology links, repo links.

---

## Pillar 2 - Blog and community (Lane C)

A blog series on [wakeward.uk](https://wakeward.uk) covering:

- The GitHub App threat model (capability vs likelihood, initial access)
- Entry points and what to watch in audit logs (once validated)
- How `gh-app-graph`, `gh-app-check`, and the future catalog fit together
- Platform research findings (Lane A) when disclosure allows

**Publish gate:** `docs/methodology.md`, attack-pattern catalog, at least one lab
trace, and working Phase 1 tooling - so the article is evidence-backed, not
spreadsheet opinion.

Tracking issue: TBD ([#7](https://github.com/wakeward/gh-app-graph/issues/7) when filed).

---

## Publication strategy (private → public)

Both **`gh-app-graph`** and **`gh-app-check`** are **private** today. Only
collaborators see GitHub history, issues, and file contents. Flipping either
repo to **Public** exposes **everything ever pushed** on that repo, including
old commits - deleting a file later does not remove it from git history.

This section defines **what to open when**, so graph/check can ship without
giving away the full product roadmap or half-baked judgments.

### What is defensible vs easy to copy

| Harder to replicate quickly | Easier to fork once public |
|---|---|
| Curated threat model + lab traces | Generic "score app permissions" idea |
| Methodology and attack-pattern catalog | YAML toxic combos + Go eval library |
| Daily Marketplace snapshot **history** | A point-in-time permission scrape |
| Storm-aware diff UX and trust | Basic `gh-app-check org` clone |
| wakeward.uk narrative and research credibility | — |

The moat is **judgment, execution, and data over time** - not the headline concept.

### Visibility tiers

| Tier | Artifacts | When |
|---|---|---|
| **Private** (now) | `docs/vision.md`, roadmap issues #4-#7, draft TM, unvalidated YAML | Until gates below are met |
| **Limited** | Company org validation notes (internal); shared with trusted reviewers if needed | During TM + calibration |
| **Public - tooling** | `gh-app-graph` library + data, `gh-app-check` CLI, `docs/methodology.md`, `threat-model-plan.md` (once reviewed) | After TM v1 + org validation |
| **Public - narrative** | wakeward.uk blog series ([#7](https://github.com/wakeward/gh-app-graph/issues/7)) | After methodology + at least one lab trace |
| **Public - product** | Marketplace catalog site, daily drift pages ([#4](https://github.com/wakeward/gh-app-graph/issues/4)) | After catalog spike + storm handling v1 |
| **Coordinated** | Lane A platform findings ([#5](https://github.com/wakeward/gh-app-graph/issues/5)) | GitHub Security disclosure process |
| **Partner-only** | Lane B triage ([#6](https://github.com/wakeward/gh-app-graph/issues/6)) | Never public malice labels without triage |

### Gates before making repos public

**`gh-app-graph` public:**

- [ ] Threat model theorizing session complete ([#3](https://github.com/wakeward/gh-app-graph/issues/3))
- [ ] `docs/methodology.md` written and self-reviewed
- [ ] At least one toxic combo or high-severity permission validated (desk or lab trace)
- [ ] Comfortable defending severity ratings in public

**`gh-app-check` public:**

- [ ] Phase 1 run against company org; false-positive rate understood
- [ ] Depends on public `gh-app-graph` release tag (or documented version pin)

**`docs/vision.md` in public git history:**

- [ ] Optional until catalog/blog launch approaches
- [ ] May stay private-repo-only longer than graph/check if product roadmap is sensitive
- [ ] If ever published, treat as intentional - history is permanent

**`gh-app-catalog` (future repo):**

- [ ] Prefer **private until site MVP**; open-source generator later if desired
- [ ] Snapshot **data store** can remain private or published as artifacts - separate decision

### Operational rules

1. **`docs/vision.md` may stay in a private repo** until catalog/blog launch; if the repo is still private, pushing this file is backup-only, not public exposure.
2. **Before changing repo visibility to Public**, decide whether `docs/vision.md` should be in history; removing it later requires history rewrite.
3. **Private repo pushes are fine** for backup and solo work - they are not public exposure.
4. **Blog before or with tooling launch** - establishes attribution on [wakeward.uk](https://wakeward.uk) even if code is forked.
5. **No rush to public** - private repos still support Issues, Projects, and Actions.
6. Before **Settings → Change visibility → Public**, run `git log --oneline` and confirm no files you regret are in history.

### If someone forks early

Apache 2.0 (current license) allows forks. Response is quality and speed of
curation, not secrecy of the eval idea. Historical Marketplace diffs and
published research are the durable differentiators.

---

## Pillar 3 - Research lanes

### Lane A - GitHub platform disclosure (primary)

Research **weaknesses in how GitHub Apps work as a platform**, not vendor
malice. Examples:

- `client_id` exposure and app identity handling in URLs, logs, or public pages
- Private key generation, storage, rotation - can key material ever be viewed
  or inferred inappropriately?
- Install / OAuth flow integrity (confused deputy, redirect handling)
- Manifest vs effective token scope
- Enterprise approval and restriction bypass

**Output:** desk traces, lab reproduction, responsible disclosure to GitHub
Security. Blog post after coordinated disclosure where applicable.

Tracking issue: [#5](https://github.com/wakeward/gh-app-graph/issues/5).

### Lane B - Suspicious Marketplace apps (secondary, careful)

**Separate from Lane A.** Use catalog diffs and rules to flag apps worth human
triage - sudden permission expansion, typosquat patterns, toxic combos plus
broad scope - without public malice labels until reviewed.

**Partners:** potential collaboration with
[OpenSourceMalware.com](https://opensourcemalware.com/) to supply structured,
reproducible detections (permission hash, diff timeline, archived manifest
evidence). Outreach **after** the catalog has diff history and a measured
false-positive rate.

Public drift pages remain neutral ("permissions changed"); Lane B triage stays
internal until a partner or GitHub Trust & Safety path accepts a report.

Tracking issue: [#6](https://github.com/wakeward/gh-app-graph/issues/6).

---

## Repository map (current and future)

| Repo | Role | Status |
|---|---|---|
| **`gh-app-graph`** | Permission catalog, dependencies, toxic combos, eval library, embedded data | Active |
| **`gh-app-check`** | Org installation auditor (`gh` CLI extension) | Phase 1 shipped |
| **`gh-app-catalog`** (name TBD) | Marketplace discover, snapshot, hash diff, static site input | Not started |
| **wakeward.uk** (blog) | Methodology, research, attribution | External |

---

## Phased roadmap

| Phase | Deliverable |
|---|---|
| **Now** | Threat model ([#3](https://github.com/wakeward/gh-app-graph/issues/3)), graph curation, org validation |
| **Next** | `gh-app-graph` release tag; `gh-app-check` Phase 2-3 as prioritized |
| **Then** | Catalog spike - discover + hash + 100-app trial ([#4](https://github.com/wakeward/gh-app-graph/issues/4)) |
| **Then** | Daily full crawl + static site MVP + public drift pages (storm handling v1) |
| **Parallel** | Lane A hypothesis list from TM; blog outline ([#7](https://github.com/wakeward/gh-app-graph/issues/7)) |
| **Later** | Product domain; OSMalware partnership ([#6](https://github.com/wakeward/gh-app-graph/issues/6)); feed if site proves useful |

---

## Explicit non-goals (for now)

- RSS/Atom feed (site-first; revisit after storm handling works)
- Public accusations of malicious apps on drift pages
- Likelihood scoring in `gh-app-check` before methodology exists
- Full marketplace catalog before threat model and org validation complete

---

## Open questions

1. Final product domain name (`gh-app-sec.*` vs subdirectory on wakeward.uk for MVP)?
2. Storm threshold percentage and auto vs manual storm classification?
3. Marketplace discover implementation (official API vs crawl - legal/ToS spike)?
4. When to approach OpenSourceMalware (after N days of diff data?)?

---

## References

- Threat model plan: [`docs/threat-model-plan.md`](threat-model-plan.md)
- GitHub Marketplace Apps: https://github.com/marketplace?type=apps
- Capability diff inspiration: [google/capslock](https://github.com/google/capslock)
- Community threat intel (future partner): [OpenSourceMalware.com](https://opensourcemalware.com/)
- Org drift guard backlog: [gh-app-check/docs/BACKLOG.md](https://github.com/wakeward/gh-app-check/blob/main/docs/BACKLOG.md)
