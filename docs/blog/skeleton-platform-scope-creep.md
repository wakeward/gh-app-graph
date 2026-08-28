# When least-privilege GitHub App installs are not enough

**Status:** skeleton  
**Sources:** `docs/references/platform-cves-app-adjacent.md`, pattern `platform-defect-u2s-scope-creep`

---

## Hook

[TODO: Team did everything right on install review. Platform CVE expands U2S token scope anyway.]

## Platform vs manifest risk

| Type | Mitigation |
|---|---|
| Over-broad manifest | `gh-app-check`, install review |
| Platform auth bug | Patch GHES / monitor advisories |

[TODO: Do not conflate the two in audit reports.]

## Example CVE table (verify before publish)

[TODO: Pull 2-3 rows from platform-cves-app-adjacent.md after NVD verification. CVE-2024-6337, CVE-2026-14340 style.]

## GHES upgrade cadence

[TODO: Practical enterprise advice.]

## References

- GitHub Security Advisories
- `docs/references/platform-cves-app-adjacent.md`
