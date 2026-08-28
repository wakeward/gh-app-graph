# AI coding agents are GitHub Apps: the external bot trust mistake

**Status:** skeleton  
**Sources:** patterns `execution-ai-agent-external-bot-trust`, `execution-indirect-prompt-injection-agent`

---

## Hook

[TODO: Workflow checks `endsWith('[bot]')`. Attacker's own App opens issue. Claude/Codex/Gemini workflow runs with secrets.]

## GitHub Apps as bot identities

[TODO: Any user can register an App. Install on attacker repo. Interact with public issues on victim repo.]

## Prompt injection second hop

[TODO: Issue body → agent reads `/proc/self/environ` or OIDC tokens → exfil to issue comment.]

## Relationship to permission audit

[TODO: Installed agent App may have reasonable manifest; failure is workflow trust model.]

## Mitigations

- [TODO: Allow-list App slug / installation ID]
- [TODO: Org-owned Apps only for automation triggers]
- [TODO: Restrict `id-token: write` on agent jobs]

## References

- `docs/references/ecosystem-research-synthesis.md` (verify claims before publish)
