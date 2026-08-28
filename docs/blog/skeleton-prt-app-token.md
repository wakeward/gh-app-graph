# pull_request_target and GitHub App installation tokens

**Status:** skeleton  
**Sources:** GHSA-9g93-rxr5-xhqw, pattern `credential-access-prt-fork-iatt-exfiltration`

---

## Hook

[TODO: Fork PR triggers docs workflow. App token ends up in git-credentials temp file. Fork code reads and exfiltrates.]

## The antipattern (step trace)

1. [TODO: `pull_request_target` - base repo context + secrets]
2. [TODO: `actions/create-github-app-token` with org PEM]
3. [TODO: `actions/checkout` with head ref + token]
4. [TODO: `cargo run` / build untrusted code]

[TODO: Mermaid or numbered diagram.]

## Why Phase 1 audit still looks "fine"

[TODO: Installed permissions unchanged. Abuse is execution-plane.]

## Detection and prevention

- [TODO: Zizmor / workflow review]
- [TODO: Never pass IAT to checkout of untrusted refs]
- [TODO: Phase 2 trace grep targets for gh-app-check]

## References

- https://github.com/advisories/GHSA-9g93-rxr5-xhqw
- `docs/credential-lifecycle.md`
