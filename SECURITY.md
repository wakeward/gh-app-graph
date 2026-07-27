# Security Policy

## Supported Versions

I support the latest major release of `github-app-permissions-graph`.

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < Latest| :x:                |

## Reporting a Vulnerability

I take the security of `github-app-permissions-graph` seriously. If you believe you have found a security vulnerability, please report it to me as described below.

### How to Report

Please do **not** report security vulnerabilities through public GitHub issues.

Instead, please report them via:
- **GitHub Security Advisories**: Use the "Report a vulnerability" button on the [Security tab](https://github.com/wakeward/github-app-permissions-graph/security/advisories).

### Response Timeline

I will acknowledge your report within **48 hours** and aim to provide a resolution or mitigation plan within **5 business days**.

### Scope

**In Scope**:
- The `cmd/*` binaries and importable `pkg/*` library packages.
- The build/refresh pipeline (`.github/workflows/refresh.yml`) and release integrity.

**Out of Scope**:
- Vulnerabilities in third-party libraries (unless actionable by me updating them).
- GitHub API behavior itself.
- The accuracy of a specific severity, blast_radius, or toxic-combination judgment call - see [docs/methodology.md](docs/methodology.md) and open an issue instead of a security advisory for those.
