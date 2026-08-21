# Platform CVEs with GitHub App impact

**Status:** draft synthesis (Aug 2026). **Verify CVE IDs and descriptions**
against NVD/GitHub advisories before blog or external citation. Sourced from
third-party research digest; some entries may be imprecise.

**Scope:** platform authorization and infrastructure defects that change
**effective** App token behavior. Not install permission misconfigurations.
**Not assessed by `gh-app-check`.**

---

## Authorization and scope enforcement

| CVE | CVSS (reported) | Type | App-relevant impact |
|---|---|---|---|
| CVE-2024-6337 | 6.8 | Incorrect authorization | App with `contents: read` + `pull_requests: write` could read private issue content via user access token path |
| CVE-2026-14340 | 5.3 | Scope creep (CWE-863) | U2S token scoped to one installation could write to public repos outside installation repo set |
| CVE-2025-6600 | Moderate | Information exposure | Zero-scope U2S + Search API could disclose private repo names given malicious App install |
| CVE-2022-23741 | 7.2 | Privilege escalation | Scoped U2S token escalation via admin installing malicious App (historical GHES) |
| CVE-2026-1999 | 7.1 | Authorization bypass | `enable_auto_merge` mutation allowed fork PR merge without push access (branch protection dependent) |

**Pattern:** `platform-defect-u2s-scope-creep` in attack patterns catalog.

**Blog angle:** least-privilege install review does not help if the platform
mis-enforces token bounds. Patch cadence and GHES upgrade path matter.

---

## Infrastructure and supply chain (adjacent)

These affect the hosting platform or CI substrate, not App manifest design:

| CVE | Type | App-relevant note |
|---|---|---|
| CVE-2024-6800 | SAML XML signature wrapping (GHES) | Site admin → arbitrary App installs |
| CVE-2026-3854 | Git backend RCE (X-Stat header) | Full platform compromise; all installations |
| CVE-2026-15343 | Dependabot updater path traversal | Write to `.github/workflows/` from container context |
| CVE-2026-54167 | Pipelines-as-Code webhook SSRF | Exfil signed App JWT before webhook verification |

Vendor integration flaws (OneUptime CVE-2026-30920, Dokploy CVE-2026-72871) are
**OAuth callback / state validation** on the vendor side - see
`credential-access-vendor-install-callback-binding`.

---

## Explicitly out of scope for gh-app-graph catalog

| Topic | Why excluded |
|---|---|
| Heroku / Travis CI OAuth token breaches (2022) | Stored **OAuth user tokens**, not GitHub App installations |
| Clone2Leak (GitHub Desktop / git credential protocol) | Local client, not App permission model |
| Generic GHES DoS / path traversal without App token angle | Platform availability, not App assessment |

Listed in [`ecosystem-research-synthesis.md`](ecosystem-research-synthesis.md)
for completeness.

---

## References

- Credential lifecycle: [`../credential-lifecycle.md`](../credential-lifecycle.md)
- GitHub Security Advisories: https://github.com/advisories
