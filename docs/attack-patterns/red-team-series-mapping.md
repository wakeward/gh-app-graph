# Red Team GitHub series → GitHub App provenance and initial access

**Status:** draft mapping (Aug 2026).
**Scope:** GitHub Apps only. This program models App manifests, installations,
declared permissions, and installation-token capability - not OAuth Apps, PATs,
gh CLI tokens, VSCode extensions, or IDE tasks.

The [Red Teaming GitHub](https://wakeward.uk/tags/github/) series is **attack
surface inspiration**: where Apps enter the environment, who approves them, and
what remote outcomes matter once installed. Parts 2 and 3 are largely out of
scope except where they clarify trust mistakes that also apply to App installs.

---

## App provenance (where an installation comes from)

Every pattern in `data/attack-patterns.yaml` should eventually tag one or more
provenance channels:

| Provenance | Description | Example patterns |
|---|---|---|
| `marketplace` | App discovered and installed via GitHub Marketplace | `initial-access-marketplace-typosquat`, `initial-access-marketplace-install-clickthrough` |
| `third_party_direct` | External publisher; install via OAuth link or vendor docs, not Marketplace | `initial-access-phishing-security-app`, `initial-access-direct-install-link` |
| `internal_custom` | Org-owned or team-built App registered on the org account | `initial-access-internal-app-registration` |
| `enterprise` | Enterprise-scoped App or enterprise-managed install policy | (TBD - theorizing session) |

Part 1's remote attack surface ([Attack Surface post](https://wakeward.uk/security/20240513_red_team_github_1/)) lists Marketplace, Actions, and webhooks. For **this program**, Marketplace and org-owned App registration are the primary **install sources**. Webhooks and Actions describe **post-install capability**, mapped to toxic combinations.

---

## Part 1 → App install paths (in scope)

| Blog theme | App-native pattern |
|---|---|
| Repository / vendor confusion | `initial-access-marketplace-typosquat` |
| Social engineering to run untrusted code | `initial-access-phishing-security-app` (Gitloker-class) |
| Marketplace as distribution | `initial-access-marketplace-install-clickthrough` |
| Org owner vs repo admin gate | `initial-access-owner-approval`, `initial-access-selected-repo-secrets` |
| Internal developer registers tooling | `initial-access-internal-app-registration` |
| Remote outcomes (ransomware, CI abuse, exfil) | Toxic combinations - capability after install |

**Field-observed chain:** fake security App → owner install → `administration` +
`contents` write → `ransomware-code-destruction` (Gitloker, cited in toxic YAML).

---

## Part 2 and Part 3 (out of scope for this catalog)

| Series content | Relationship to GitHub Apps program |
|---|---|
| Part 2 - `gh auth token`, OAuth in keyring | **Out of scope.** User-delegated OAuth, not App installation tokens. |
| Part 2 - gh CLI extension typosquat | **Analogy only** for Marketplace App typosquat trust mistakes. Not catalogued as a pattern. |
| Part 3 - VSCode `tasks.json` / extensions | **Out of scope.** Developer IDE vectors, not App OAuth installs. |

**In scope from Part 2 (App-only):** GitHub's policy that exposed **GitHub App**
credentials in public repos/gists are revoked → `credential-access-leaked-app-private-key`
(how an existing App identity is abused, not how a new App lands).

---

## Initial access priority list (Apps only)

1. Marketplace typosquat or trusted-vendor mimic
2. Phishing / notification-driven install of rogue App (direct or Marketplace)
3. Marketplace or vendor install without permission review
4. Internal custom App registered by developer, approved by owner
5. Selected-repo install targeting CI/secrets repos
6. Publisher manifest permission upgrade (S2a)
7. Leaked App private key → mint installation tokens (credential access, not new install)

Install approver requirements: [`installation-gates.md`](../installation-gates.md).

---

## Gaps still open (Apps only)

- Enterprise App approval workflows and `enterprise_*` permission scope
- App ownership transfer between publishers
- Suspended App reactivation
- GHES-only App permissions
- Verified audit log events for `integration_installation.*`

---

## References

- [Part 1 - Attack Surface](https://wakeward.uk/security/20240513_red_team_github_1/)
- [Part 2 - GitHub CLI](https://wakeward.uk/security/20240821_red_team_github_2/) (trust analogies only)
- [Part 3 - VSCode](https://wakeward.uk/security/20241211_red_team_github_3/) (out of scope)
