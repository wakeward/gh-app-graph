# Ecosystem research synthesis (Aug 2026)

**Status:** internal digest. Summarizes external research on GitHub App
vulnerabilities with **Wakeward scope labels**. Use to enrich attack patterns,
blog drafts, and Phase 2 backlog. Does not change `gh-app-check` scoring rules.

**Provenance:** third-party LLM research report (Gemini), cross-checked against
existing catalog and program boundaries.

---

## What we incorporated

| Theme | Where it lives | Assessment impact |
|---|---|---|
| Token typology (PEM, JWT, IAT, U2S) | [`credential-lifecycle.md`](../credential-lifecycle.md) | None (reference) |
| Marketplace verification paradox | [`marketplace-trust-limitations.md`](../marketplace-trust-limitations.md), attack pattern | None |
| Permission granularity ceiling | `marketplace-trust-limitations.md` | Explains over-provisioning likelihood in blog |
| `pull_request_target` + App token exfil | Attack pattern `credential-access-prt-fork-iatt-exfiltration` | Phase 2 trace candidate |
| AI agent `[bot]` trust bypass | Attack patterns (Claude/Codex/Gemini CLI) | Phase 2 trace candidate |
| Vendor OAuth callback / state flaws | Attack pattern `credential-access-vendor-install-callback-binding` | Reference only |
| Webhook SSRF → JWT exfil | Attack pattern `credential-access-webhook-handler-ssrf-jwt` | Reference only |
| Dangling DNS on App callbacks/webhooks | Attack pattern `persistence-dangling-dns-oauth-webhook` | Reference only |
| Platform U2S / IAT scope CVEs | [`platform-cves-app-adjacent.md`](platform-cves-app-adjacent.md) | Reference only |
| Complementary tools (Zizmor, Legitify, Harden-Runner) | `marketplace-trust-limitations.md`, gh-app-check backlog | None |

---

## What we deliberately excluded from assessment

| Theme | Reason |
|---|---|
| User install authorization / "who could install" scoring | Witch-hunt risk; stays in [`installation-gates.md`](../installation-gates.md) as blog reference |
| OAuth App / PAT storage breaches (Heroku, Travis) | Different credential model; out of program scope |
| Git client / Desktop credential leaks | Local client surface |
| Pure platform RCE/DoS without App permission lesson | Patch management topic; minimal install-review value |
| Marketplace **Actions** mutable tag incidents | Actions supply chain; note as adjacent in marketplace doc only |
| Automated likelihood or publisher trust scores | Explicit non-goal |

---

## New attack patterns added (summary)

| ID | Tactic | Scope |
|---|---|---|
| `initial-access-marketplace-verification-paradox` | initial_access | app_install |
| `credential-access-prt-fork-iatt-exfiltration` | credential_access | execution_plane |
| `execution-ai-agent-external-bot-trust` | execution | execution_plane |
| `execution-indirect-prompt-injection-agent` | execution | execution_plane |
| `credential-access-vendor-install-callback-binding` | credential_access | vendor_integration |
| `credential-access-webhook-handler-ssrf-jwt` | credential_access | vendor_integration |
| `persistence-dangling-dns-oauth-webhook` | persistence | infrastructure |
| `platform-defect-u2s-scope-creep` | execution | platform_defect |

---

## Suggested blog series outline (from synthesis)

1. **Control plane vs execution plane** - what `gh-app-check` measures and what it cannot
2. **Marketplace is not a security review** - verification paradox + permission granularity
3. **The PRT + App token antipattern** - skim-rs / GHSA walkthrough
4. **When least-privilege install is not enough** - platform CVE table (U2S scope)
5. **AI bots are GitHub Apps** - external `[bot]` trust mistakes

---

## Open verification

- [ ] Confirm each CVE ID in `platform-cves-app-adjacent.md` against primary sources
- [ ] Desk trace for `credential-access-prt-fork-iatt-exfiltration`
- [ ] Audit log event names for AI agent workflow triggers
