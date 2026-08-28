# Publish readiness checklist

Living checklist for publishing `gh-app-graph` and consumer `gh-app-check`.
Consumer-side steps: [`gh-app-check` publish doc](https://github.com/wakeward/gh-app-check/blob/main/docs/PUBLISH-READINESS.md).

## Phase A - Policy, scrub, OSS hygiene (done)

- [x] `organization-takeover` single-grant policy decided (`docs/calibration-notes.md`)
- [x] README typo + `evaluate-app` marked not implemented
- [x] `CONTRIBUTING.md`, issue templates, `.github/CODEOWNERS`
- [x] `docs/blog/` banner already marks internal WIP skeletons

## Phase B - Release prep (done)

- [x] CI triggers on `data/**` and `pkg/data/bundled/**`
- [x] `go test` in `refresh.yml` before auto-PR
- [x] Tag v0.1.1 (GPG-signed; module pin for gh-app-check)
- [x] Consumer pin in `gh-app-check` validated in CI
- [x] Token-permissions: job-scoped writes in `refresh.yml`

Signed tags: use `git tag -s` before push. See
[`gh-app-check` setup runbook](https://github.com/wakeward/gh-app-check/blob/main/docs/GITHUB_SETUP_RUNBOOK.md#signed-git-tags-required).

## Phase C - History reset (done)

- [x] Orphan squash to single public commit
- [x] GPG-signed `v0.1.0` / `v0.1.1` tags

## Phase D - Publish (waiting on gh-app-check release)

- [x] Flip repo visibility public
- [ ] Merge finish-publish PR (token-permissions on refresh.yml)
- [ ] Consumer `gh-app-check` v0.1.0 GitHub release cut

## Data disclaimer for v0.1.0

Permission YAML is **draft seed data**. Toxic combinations and methodology
are more mature but still evolving. README and `data/permissions/README.md`
state this explicitly.
