# Who can install a GitHub App (and why our auditor does not score it)

**Status:** skeleton  
**Sources:** `docs/installation-gates.md`

---

## Hook

[TODO: Reader sees Critical finding, asks "which admin failed?" Reframe: assess **what the App can do**, not **who clicked install**.]

## Three install levels

| Level | Who | Constraints |
|---|---|---|
| Enterprise | Enterprise Owner | App Manager insufficient |
| Organization | Organization Owner | Org-level perms, admin write |
| Repository | Repo admin | Zero org perms, no admin write, policy allows |

[TODO: Expand each with one sentence. No naming individuals.]

## Reader's guide (optional inference)

[TODO: Org-level grant → not repo-admin-only path. Heuristic for **blog readers**, not tool output.]

## Enterprise-owned App permissions

[TODO: `enterprise_organization_installations` third-party restriction.]

## Why gh-app-check excludes this

[TODO: Witch hunt risk; inference breaks on upgrades; audit attribution is different product.]

## References

- `docs/installation-gates.md`
