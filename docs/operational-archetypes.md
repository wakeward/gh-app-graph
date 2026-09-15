# Operational archetypes (u2s / s2s / r2w)

**Status:** backlog capture (2026-09-15). Not implemented in the YAML catalog yet.
**Consumers:** `gh-app-catalog` (public Marketplace classification + IOC),
`gh-app-check` (installed-app profile, later), `gh-app-sec-research` (platform
questions only).

GitHub documents three ways an App can act. They are **not mutually exclusive**.
A Marketplace bot is usually **s2s + r2w**. A SaaS dashboard is often all three.

Source: [About creating GitHub Apps](https://docs.github.com/en/apps/creating-github-apps/about-creating-github-apps/about-creating-github-apps).

| Code | Docs name | Token / auth | Typical trigger |
|---|---|---|---|
| **u2s** | Act on behalf of a user | User access token (OAuth / device flow) | Browser, CLI, IDE |
| **s2s** | Act on their own behalf | App JWT → installation access token | Cron, CI, webhook handler |
| **r2w** | Respond to webhooks | HMAC `X-Hub-Signature-256` (no GitHub API token required to *receive*) | Inbound POST from GitHub |

## Why this belongs in gh-app-graph

Permission severity and toxic combos answer **what the install can do**.
Archetypes answer **how the app is supposed to wake up**. That is the join key
when a new anti-pattern appears: filter Marketplace apps by `events` +
`permissions` + inferred mode, instead of re-scoring the whole catalog by hand.

This is **not** a substitute for installation context. Public
`GET /apps/{slug}` does **not** prove the publisher implemented OAuth, stored a
PEM, or validates webhook signatures. It only exposes **declared** surface.

## What the public API actually gives

From unauthenticated (or PAT-rate-limited) `GET /apps/{slug}`:

| Field | Use for |
|---|---|
| `permissions` | Already scored in catalog / check |
| `events` | **r2w** indicator (non-empty subscribed events) |
| `external_url` | Publisher site; later dangling-DNS IOC |
| `installations_count` | Scale, not privilege |
| `owner` | Publisher identity drift |
| Numeric `id` | Stable identity if slug changes (see T08) |

**Not public** (do not invent them from Marketplace HTML):

- `callback_url` / redirect URIs
- `hook_attributes.url` / webhook URL
- Whether a PEM exists
- Whether `webhook_secret` is set
- User-permission block vs installation-permission block as separate maps
  (the public payload is the installation permission map)

**Inference, not proof:**

- Non-empty `events` → treat as **r2w-capable** (reactive bot if also write perms).
- Any installation-scoped permission → **s2s-capable** (every GitHub App that
  can be installed is s2s-capable; the interesting split is events + writes).
- **u2s** cannot be confirmed from `/apps/{slug}` alone. User-level scopes and
  callback URLs are owner-only. Do not scrape the install HTML consent page as
  a product feature (ToS / brittle / easy to confuse with OAuth Apps).

Legacy **OAuth Apps** have no public scopes endpoint. Out of scope for
`gh-app-graph` GitHub-App catalog. If Marketplace rows are Actions-only or
OAuth, `gh-app-catalog` already classifies them as non-scoreable.

## Anti-pattern join (graph side)

Keep exploit prose in `data/attack-patterns.yaml`. Add optional selector fields
later (do not duplicate RCE write-ups here):

| Selector | Maps toward |
|---|---|
| `events` intersects `issue_comment`, `pull_request`, `issues` + write on contents/issues/PRs | Untrusted input → confused deputy (`execution-abuse-declared-scopes`) |
| `events` intersects `workflow_run`, `check_run` + `actions`/`checks` write | CI privilege path |
| Non-empty `events` + any high write | Webhook handler is in the blast radius (`credential-access-webhook-handler-ssrf-jwt`) |
| High write, empty `events` | Pure s2s / batch; PEM theft is the story, not inbound GitHub POSTs |
| Permission drift on those same keys | Catalog IOC Phase 1 (escalation), not a new pattern ID |

Modes stack. **Reactive bot** = r2w + s2s. **Interactive platform** = u2s + s2s
+ r2w (u2s inferred only inside `gh-app-check` or owner-token research).

## Related ideas (split by repo)

### Enumerate IDs / abandoned apps / "take ownership"

- **Catalog:** flag apps whose numeric `id` is stable, slug/owner changed, or
  inferred upstream repo has no push in 12+ months **and** still requests write
  scopes. That is IOC Phase 2+ (`docs/IOC-TRACKING.md` in gh-app-catalog).
- **Graph:** pattern already covers typosquat and verification paradox; add
  "abandoned high-privilege Marketplace app" as a **detection note**, not a
  new exploit path.
- **Do not:** register lookalike apps, squat expired publisher domains for
  catch, or "safely take ownership" of third-party Marketplace listings.
  That is B3, outside GitHub safe harbor, and not a catalog feature.

### App registration credentials vs a normal user

Own-account App JWT / IAT can read **installation** APIs a user PAT cannot.
That does **not** mean Marketplace enumeration gets extra fields: catalog
already uses the public `/apps/{slug}` path. Extra-vs-user questions live in
`gh-app-sec-research` (T01 never-public read, T18 GraphQL/REST split), on
**self-owned** apps only.

### MitM / sniff traffic to an App server

Out of scope for graph implementation and for catalog. Existing patterns:

- `credential-access-webhook-handler-ssrf-jwt`
- `persistence-dangling-dns-oauth-webhook`
- `credential-access-vendor-install-callback-binding`

Those are **vendor/DNS** failure modes. No traffic interception, no attacking
third-party webhook hosts.

### Repo maintenance (no commit in 12 months)

Catalog metadata job once `events` + owner/repo heuristic exist. Graph only
needs the rule: abandoned + still `contents`/`actions`/`secrets` write →
elevated **publisher-compromise likelihood**, not a higher toxic-combo score.

## Implementation order (when we build it)

1. **Catalog (done in tree):** persist `events` and numeric `id` from `GET /apps/{slug}`
   on each scan. UI treats non-empty events as **r2w-capable**, not proof of a live webhook.
2. **Catalog IOC Phase 1 (done in tree):** permission escalation + storm flag.
   Overlay later: `workflows:write` + `workflow_run` subscription louder than write alone.
3. **Graph:** optional `events_selectors` / `archetype` on attack-pattern YAML
   after snapshots contain events so queries are real.
4. **Abandoned / ID lifecycle:** after owner + events are in snapshots.

## Open questions

1. Is `events: []` on a Marketplace app a useful "pure s2s" signal, or do many
   publishers omit events in the public payload even when they listen?
2. Should archetype affect the 0-10 exploitability score, or stay a **tag**
   next to IOC? Recommendation: **tag** (same as IOC vs risk).
3. Numeric App `id` in snapshots for T08 slug-reuse detection - confirm field
   stability on `GET /apps/{slug}`.
