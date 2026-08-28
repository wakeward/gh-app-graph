# Marketplace verification is not a security review

**Status:** skeleton  
**Sources:** `docs/marketplace-trust-limitations.md`, pattern `initial-access-marketplace-verification-paradox`

---

## Hook

[TODO: Verified badge on Marketplace listing. Admin installs with org-wide `contents: write`. What verification actually checked.]

## Verification paradox

[TODO: TXT record, 2FA, email - vs PEM storage, callback validation, manifest design.]

## Permission granularity ceiling

[TODO: Why SAST bots need broad scopes; no `/src`-only ACL. Tie to over-provisioning **likelihood**, not new severity scores.]

## Publisher compromise cascade

[TODO: S2b - same permissions before/after compromise. Optional adjacent note on Actions mutable tags (out of App scope).]

## What to do instead

- [TODO: Allow lists, approval policies]
- [TODO: `gh-app-check org` before/after install - capability baseline]
- [TODO: Vendor architecture review for OAuth callbacks]

## References

- Red Teaming GitHub Part 1 (Marketplace click-through)
- `docs/marketplace-trust-limitations.md`
