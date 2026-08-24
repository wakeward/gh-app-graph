# Publish readiness checklist

Living checklist for publishing `gh-app-graph` and consumer `gh-app-check`.
Consumer-side steps: [`gh-app-check` publish doc](https://github.com/wakeward/gh-app-check/blob/main/docs/PUBLISH-READINESS.md).

## Phase A - Policy, scrub, OSS hygiene (in progress)

- [ ] `organization-takeover` single-grant policy decided (`docs/calibration-notes.md`)
- [x] README typo + `evaluate-app` marked not implemented
- [x] `CONTRIBUTING.md`, issue templates, `.github/CODEOWNERS`
- [x] `docs/blog/` banner already marks internal WIP skeletons

## Phase B - Release prep (private)

- [ ] CI triggers on `data/**` and `pkg/data/bundled/**`
- [ ] `go test` in `refresh.yml` before auto-PR
- [ ] Tag v0.1.0 after policy call + consumer pin

## Data disclaimer for v0.1.0

Permission YAML is **draft seed data**. Toxic combinations and methodology
are more mature but still evolving. README and `data/permissions/README.md`
state this explicitly.
