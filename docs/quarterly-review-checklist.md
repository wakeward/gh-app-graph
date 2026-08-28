# Quarterly review checklist

This is the fixed list of things that cannot be fully automated and that the
quarterly Cursor Automation (or a human, if the Automation isn't set up yet)
must work through by hand each run. Weekly CI (`.github/workflows/refresh.yml`)
only re-runs the two cheap, fully-automatable steps (`fetch-inventory`,
`detect-overlap`) - everything below is deliberately out of scope for it.

## 1. Wide Type B (prose-scraped dependency) sweep

Re-run `scrape-prose` widened to *every* flagged permission, not just
High/Critical severity (the default weekly/on-demand scope). Review the
resulting `data/generated/scraped-dependencies-draft.yaml` and promote
confirmed edges into `data/dependencies-manual.yaml`.

## 2. Re-check undocumented/unconfirmed permissions

For every permission in `data/permissions/*.yaml` with
`doc_status: undocumented_preview` or `doc_status: unconfirmed_key`, check
the current [Permissions required for GitHub Apps](https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps)
page to see if it now has documented endpoints. Update `doc_status` and
`needs_investigation` accordingly.

## 3. GitHub changelog review

Check the [GitHub Changelog](https://github.blog/changelog/) and the Apps
changelog for permission-related announcements since the last run (new
permissions, renamed permissions, deprecated endpoints).

## 4. Open `needs_investigation` review

Review every row across `data/permissions/*.yaml` with
`needs_investigation: true` and either resolve it (update severity, notes,
doc_status) or leave a comment explaining what's still blocking it.

## 5. New toxic-combination candidates

Given anything promoted in step 1 or newly announced in step 3, consider
whether `data/toxic-combinations.yaml` needs a new entry. Not every new
dependency edge implies a new toxic combination - see
[docs/methodology.md](methodology.md) for how a combination earns its own
entry (independently-grantable permissions, a named attack technique, a
concrete exploit path).

## 6. Live App-creation UI diff (cannot be automated)

Open a GitHub App's "Permissions" settings tab in a browser
(**Settings > Developer settings > GitHub Apps > (an app) > Permissions**)
and diff the permission picker against `data/permissions/*.yaml`. This is
how the undocumented Enterprise permissions (Custom enterprise roles,
Enterprise AI controls, Enterprise credentials, and others) were originally
found - there's no public API that enumerates the picker independent of
documented endpoints, so this step needs a human (or a future
browser-capable agent run) every quarter.

## Output

Open a single PR with everything found/changed, clearly separating
"high-confidence, mechanically derived" changes (steps 1-2) from "judgment
calls - please review" (steps 4-5), and leave a checklist comment noting
whether step 6 (the live UI diff) was completed this run.
