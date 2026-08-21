# GitHub App installation gates

**Status:** draft reference (Aug 2026). Actor requirements for who can install
a GitHub App and at what scope.

**Purpose:** reference material for threat modeling, blog posts, and human
reasoning about **real-world likelihood** ("who had to click approve?"). This
is **not** App permission data and **not** part of `gh-app-check` assessment
output.

These are **human/process gates**, not GitHub App permissions. Document in
attack patterns and blog content. Do **not** encode as toxic YAML entries or
`gh-app-check` rules.

Source: GitHub Apps documentation and platform behavior (verify against current
docs before treating as canonical).

---

## Relationship to `gh-app-check`

| Tool / artifact | Scope |
|---|---|
| **`gh-app-check org`** | Installed **App permissions** only (capability / impact) |
| **This document** | **User** roles and policies required to complete an install (likelihood context) |

`gh-app-check` does not infer, score, or report install approver paths. A
security reviewer (or blog reader) may use the tables below **after** an audit
to reason about how an installation likely landed - that is outside the tool.

---

## Install levels

| Level | Who can install | Scope granted |
|---|---|---|
| **Enterprise** | Enterprise Owner only | Entire enterprise account |
| **Organization** | Organization Owner | One organization (all or selected repos per app config) |
| **Repository** | Repository admin (conditional) | Only repos where the admin has admin rights |

**App Manager role:** can manage apps in some contexts but **cannot** install
apps at the **enterprise** level.

**Enterprise Owner + Org Owner:** on credit-card-billed enterprises, an Enterprise
Owner who is also an Organization Owner may install apps on orgs within the
enterprise.

---

## Enterprise level

**Requires:** Enterprise Owner.

**Cannot install at enterprise level:** users with only the App Manager role.

**Third-party restriction:** if a third-party App requests either of these
permissions, you **cannot** install it on your enterprise - the App is
restricted to installation on the enterprise that **owns** the App:

- `enterprise_organization_installations`
- `enterprise_organization_installation_repositories`

Implication for threat modeling: those enterprise permissions imply an
**enterprise-owned** (internal/custom) App, not a arbitrary Marketplace vendor App.

---

## Organization level

**Requires:** Organization Owner.

Org-wide installs (`repository_selection: all`) and org-level permission grants
require this path (or enterprise-level install above).

---

## Repository level (repository admin)

A **repository admin** may install an App on repos they administer only when
**all** of the following hold:

1. The App requests **zero organization-level permissions**.
2. The App does **not** request the repository **`administration`** permission.
3. Organization Owners have **not** enabled the setting that prevents
   repository admins from installing apps.

**Scope:** access is limited to repositories where the installing user is a
repository admin.

Implication for **likelihood reasoning** (blog, threat model, manual review):

- Installed Apps with **org-level permissions** or **`administration: write`**
  could not have been installed via the repository-admin-only path (owner or
  enterprise install required).
- The repository-admin path is a **lower-friction initial access** route -
  targeted installs on sensitive repos without org owner involvement, but
  constrained in what the App manifest can request.

---

## Mapping to attack patterns

| Install level | Example pattern IDs |
|---|---|
| Organization | `initial-access-owner-approval`, `initial-access-marketplace-*`, `initial-access-phishing-security-app`, `initial-access-internal-app-registration` |
| Repository (conditional) | `initial-access-repository-admin-install` |
| Enterprise | `initial-access-enterprise-owner-install` |

---

## Reader's guide: inferring install path from permission grants

When interpreting an org audit or writing about risk, a human reviewer can use
installed **App** permissions to narrow which **user** install gate likely applied.
This is educational context only - not automated by `gh-app-check`.

| Observation (App permissions installed) | Likely user install gate |
|---|---|
| Any org-level permission grant | Not repository-admin-only install |
| `administration: write` | Not repository-admin-only install |
| `organization_administration: write` | Organization or enterprise owner path |
| `enterprise_*` permissions | Enterprise-owned App context (see restriction above) |
| Narrow repo-only grants, no admin write | May include repository-admin install path |

---

## Open verification items

- [ ] Exact org setting name for "prevent repository admins from installing apps"
- [ ] Audit log events per install level (enterprise vs org vs repo)
- [ ] App Manager capabilities at org level vs enterprise level
- [ ] GHEC vs GHES differences for enterprise install policy

---

## References

- Attack patterns: [`data/attack-patterns.yaml`](../data/attack-patterns.yaml)
- Methodology: [`methodology.md`](methodology.md)
- GitHub Docs: [Installing GitHub Apps](https://docs.github.com/en/apps/using-github-apps/installing-github-apps)
