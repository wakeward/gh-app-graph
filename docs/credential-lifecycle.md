# GitHub App credential lifecycle

**Status:** draft reference (Aug 2026). Blog and threat-model material. Not
scored by `gh-app-check`.

GitHub Apps use several credential types. `gh-app-check org` evaluates **declared
installation permissions** (what an IAT may request once minted). This document
covers **how credentials are obtained and abused** - execution and credential
access context for Phase 2 `trace` and attack patterns.

---

## Credential types

| Credential | Issued by | Scope | Typical abuse |
|---|---|---|---|
| **App private key (PEM)** | Publisher at App creation; long-lived | Signs JWTs as the App entity | Vendor server compromise; leak in repo or CI secret |
| **JWT (App authentication)** | Publisher signs locally with PEM | Authenticates App to GitHub API to request IATs | SSRF sending JWT to attacker host; header manipulation |
| **Installation access token (IAT)** | GitHub after valid JWT + installation ID | Repos and permissions granted at install; short-lived (~1 h) | Exfil from CI logs, `git-credentials`, workflow output |
| **User-to-server token (U2S)** | GitHub after user authorizes App in web flow | Intersection of user access and App permissions | Platform auth bugs; OAuth state confusion; appears as user in audit |

**Assessment boundary:** Phase 1 maps installed **permissions** (IAT ceiling). It
does not detect PEM leakage, workflow misuse, or U2S platform defects.

---

## Lifecycle (simplified)

```
PEM (publisher) → JWT → IAT (per installation, per repo set)
                         ↑
User OAuth flow → U2S (optional; user-scoped operations)
```

An installation's permission manifest caps what each freshly minted IAT may do.
Runtime abuse still requires a valid credential path (publisher backend, workflow
step, or leaked PEM).

---

## CI/CD antipatterns (execution plane)

Documented in attack patterns; validated examples include:

| Antipattern | Risk |
|---|---|
| `pull_request_target` + checkout of fork HEAD + App token in git credentials | IAT exfiltration via untrusted code (`credential-access-prt-fork-iatt-exfiltration`) |
| Passing IAT to `actions/checkout` before running untrusted build scripts | Token on disk at `/home/runner/work/_temp/git-credentials-*` |
| `actions/create-github-app-token` with org PEM in secrets on PRT workflows | High-value secret + privileged base-repo context |

Phase 2 `gh-app-check trace` should grep workflows for these pairings. That is
**workflow configuration**, not App manifest assessment.

---

## Token redaction failures (adjacent)

Third-party tools parsing or displaying tokens can leak IATs/U2S tokens in logs
(e.g. Composer regex mismatch on `ghs_` hyphenated tokens; historical `gh auth
status` masking gaps). Out of scope for permission catalog; relevant for
execution-plane monitoring and secret scanning tuning.

---

## Platform authorization defects (adjacent)

Historical GHES/GitHub.com CVEs expanded effective U2S or IAT scope beyond install
boundaries (e.g. cross-repo public write, private issue read via scope confusion).
These are **platform bugs**, not misconfigured installs. See
[`references/platform-cves-app-adjacent.md`](references/platform-cves-app-adjacent.md).

---

## References

- Attack patterns: [`data/attack-patterns.yaml`](../data/attack-patterns.yaml)
- Methodology: [`methodology.md`](methodology.md)
- GitHub Docs: [Authenticating with GitHub Apps](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app)
