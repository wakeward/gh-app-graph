# What GitHub App permission auditing can (and cannot) see

**Status:** skeleton  
**Sources:** `docs/methodology.md`, `docs/credential-lifecycle.md`, `gh-app-check` Phase 1

---

## Hook

[TODO: Org runs `gh app-check org`, finds Critical toxic combo on a "legitimate" CI App. Reader asks: "Are we breached?" Frame: you learned **capability**, not **runtime abuse**.]

## The two planes

| Plane | Question | Tool today |
|---|---|---|
| Control | What permissions are installed? | `gh-app-check org` |
| Execution | How are tokens minted and used in workflows? | Phase 2 `trace` (planned) |

[TODO: Simple diagram - manifest/installation API vs Actions workflow.]

## What Phase 1 answers well

- Toxic permission combinations (from `gh-app-graph`)
- All-repositories blast radius
- `administration: write` and god-mode write counts
- GHES-only scope highlighting

[TODO: One sanitized table row example.]

## What Phase 1 does not answer

- Whether a publisher backend is compromised (S2b)
- `pull_request_target` + App token exfiltration (S5)
- Vendor OAuth callback bugs (S6)
- Platform CVE scope expansion

[TODO: Link each to attack pattern ID, not deep exploit prose.]

## How to use findings without a witch hunt

[TODO: Remediation = narrow permissions, repo selection, remove App. Do not hunt approvers. Link `installation-gates.md` as optional reading.]

## Call to action

[TODO: Run baseline audit; pair with Zizmor for workflows; track publisher upgrades.]

## References

- https://github.com/wakeward/gh-app-check
- https://github.com/wakeward/gh-app-graph/blob/main/docs/methodology.md
