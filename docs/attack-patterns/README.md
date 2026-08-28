# GitHub App attack patterns (draft catalog)

**Status:** draft seed for theorizing session ([#3](https://github.com/wakeward/gh-app-graph/issues/3)).
**Scope:** GitHub Apps only - manifests, installations, declared permissions,
installation tokens. Not OAuth Apps, PATs, gh CLI, or IDE tooling.

**Format:** structured YAML in [`attack-patterns.yaml`](../../data/attack-patterns.yaml).

Each pattern documents **how** an App installation lands (actor, gates, detection).
Toxic combinations document **what** granted permissions enable. Cross-link by
`toxic_combination_ids` only - do not duplicate exploit prose.

## App provenance (install source)

Where an App enters the org - the primary lens for initial access:

| Provenance | Meaning |
|---|---|
| `marketplace` | Installed from GitHub Marketplace |
| `third_party_direct` | External publisher; direct OAuth install link or vendor channel |
| `internal_custom` | Org-owned or team-built App on the org account |
| `enterprise` | Enterprise-owned App; enterprise-level install |

Optional YAML field: `app_provenance[]` on each pattern.

## Pattern scope

Not every pattern is an **install** path. Optional field: `pattern_scope`.

| Scope | Meaning | In gh-app-check? |
|---|---|---|
| `app_install` | How permissions become an installation | Indirectly (permissions audited after install) |
| `execution_plane` | Workflow/token abuse with unchanged manifest | Phase 2 `trace` only |
| `vendor_integration` | Third-party SaaS mishandles App OAuth/webhooks | No |
| `platform_defect` | GitHub platform CVE expands token scope | No |
| `infrastructure` | DNS, callback hosting | No |

## Install level (who can approve)

Documented in [`installation-gates.md`](../installation-gates.md). Optional YAML
field: `install_level[]`. **Reference and blog material only** - not part of
`gh-app-check` assessment.

| Level | Approver | Constraint summary |
|---|---|---|
| `enterprise` | Enterprise Owner only | App Manager cannot install at enterprise level; some enterprise permissions require enterprise-owned App |
| `organization` | Organization Owner | Org-level permissions, `administration` write, org-wide repo selection |
| `repository` | Repository admin | Zero org-level permissions; no `administration` write; org policy allows |

For **likelihood reasoning** after an audit: org-level grants or `administration: write`
on an installation imply a user with org (or enterprise) owner authority approved
it, not a repository-admin-only path. See the reader's guide in
[`installation-gates.md`](../installation-gates.md).

## Prior art: Red Teaming GitHub series

The [Red Teaming GitHub](https://wakeward.uk/tags/github/) posts inform **attack
surface and install trust mistakes**, not the full pattern catalog. Mapping:
[`red-team-series-mapping.md`](red-team-series-mapping.md).

| Provenance | Example patterns |
|---|---|
| Marketplace | `initial-access-marketplace-typosquat`, `initial-access-marketplace-install-clickthrough`, `initial-access-marketplace-verification-paradox` |
| Third-party direct | `initial-access-phishing-security-app`, `initial-access-direct-install-link` |
| Internal custom | `initial-access-internal-app-registration` |

## Tactics (working taxonomy)

| Tactic ID | Question |
|---|---|
| `app_registration` | How is the App manifest defined? |
| `initial_access` | How do requested permissions become an installation? |
| `permission_acquisition` | How do grants expand after install? |
| `credential_access` | How is an existing App identity abused (e.g. leaked PEM)? |
| `execution` | How is installed capability exercised? |
| `persistence` | How is App access maintained? |
| `defense_evasion` | How is abuse hidden? |
| `impact` | What is the outcome? |

## Scenario mapping (S1-S6)

| Ref | Summary | Patterns |
|---|---|---|
| S1 | Author lacks install rights; approver gate | `initial-access-owner-approval`, `initial-access-repository-admin-install`, `initial-access-enterprise-owner-install`, provenance-specific paths |
| S2a | Publisher expands manifest | `permission-acquisition-publisher-upgrade` |
| S2b | Malicious use of unchanged grants | `execution-abuse-declared-scopes`, `platform-defect-u2s-scope-creep` |
| S3 | Ephemeral settings change | `defense-evasion-ephemeral-admin-change` |
| S4 | Portfolio of narrow Apps | `impact-portfolio-spray` |
| S5 | Workflow / AI agent execution abuse | `credential-access-prt-fork-iatt-exfiltration`, `execution-ai-agent-external-bot-trust`, `execution-indirect-prompt-injection-agent` |
| S6 | Vendor callback / DNS infrastructure | `credential-access-vendor-install-callback-binding`, `credential-access-webhook-handler-ssrf-jwt`, `persistence-dangling-dns-oauth-webhook` |

## External research

Third-party ecosystem synthesis (Aug 2026):
[`references/ecosystem-research-synthesis.md`](../references/ecosystem-research-synthesis.md).
Verify CVE IDs before blog citation.

## Next steps

1. Theorizing session - gap-check App-only paths; enterprise provenance
2. Verify audit log event names per license tier
3. Desk traces for top initial-access patterns
4. Map every Critical toxic combo to at least one install or upgrade pattern
