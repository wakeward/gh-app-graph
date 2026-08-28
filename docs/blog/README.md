# Blog authoring skeletons (GitHub App Security program)

**Status:** internal drafts. These files are **skeletons** for Wakeward blog
posts under `/security/`. They are not published content.

## Workflow

1. Pick a skeleton below (or copy `_template.md`).
2. Fill sections marked `[TODO]`; link to canonical program docs rather than
   duplicating permission tables.
3. Run `gh app-check org` on a **sanitized** example if you need a screenshot
   (never commit org-specific output).
4. Verify CVE/advisory IDs before publish (`docs/references/platform-cves-app-adjacent.md`).
5. Publish via the Wakeward blog repo (Hugo / Blowfish); link back to
   `gh-app-graph` and `gh-app-check` when tools are public.

## Scope reminder (every post)

| Say | Do not say |
|---|---|
| `gh-app-check` reports **installed App permissions** | The tool identifies who installed an App |
| Install gates are **likelihood context** for readers | Name or blame specific approvers |
| Phase 1 is a **control-plane snapshot** | PASS means safe |

## Skeleton index

| File | Working title | Primary docs |
|---|---|---|
| [`skeleton-control-plane-vs-execution-plane.md`](skeleton-control-plane-vs-execution-plane.md) | What GitHub App permission auditing can and cannot see | `methodology.md`, `credential-lifecycle.md` |
| [`skeleton-marketplace-trust.md`](skeleton-marketplace-trust.md) | Marketplace verification is not a security review | `marketplace-trust-limitations.md` |
| [`skeleton-prt-app-token.md`](skeleton-prt-app-token.md) | pull_request_target + App tokens | `credential-access-prt-fork-iatt-exfiltration` pattern, GHSA |
| [`skeleton-install-gates-likelihood.md`](skeleton-install-gates-likelihood.md) | Who can install a GitHub App (and why we do not score it) | `installation-gates.md` |
| [`skeleton-ai-bots-as-apps.md`](skeleton-ai-bots-as-apps.md) | AI agents are GitHub Apps: bot trust mistakes | S5 patterns in `attack-patterns.yaml` |
| [`skeleton-platform-scope-creep.md`](skeleton-platform-scope-creep.md) | When least-privilege install review is not enough | `platform-cves-app-adjacent.md` |

## Template

Copy [`_template.md`](_template.md) for new posts.

## Series placement

Suggested tag: `github`, `github-apps`, `supply-chain`.

Optional series link: [Red Teaming GitHub](https://wakeward.uk/tags/github/) (Apps-only follow-on, not gh CLI / IDE).
