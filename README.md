# Claw Engineering Template

![CI](https://github.com/1mf-claw/claw-engineering-template/actions/workflows/ci-release.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/1mf-claw/claw-engineering-template)
![License](https://img.shields.io/badge/license-MIT-green)

Production-ready full‑stack template (Go + Vue) with CI/CD, SRE baseline, and governance baked in.

## What You Get
- Go backend (clean HTTP service)
- Vue 3 + Vite frontend
- GitHub Actions CI (lint / test / cross‑platform build)
- semantic‑release (automated versioning)
- Multi‑environment config (dev / staging / prod)
- SRE baseline (metrics, alerting, systemd ready)
- Governance (PR template, CODEOWNERS, SECURITY policy)

## Use This Template
1. Click **Use this template**
2. Create a new repository
3. Clone the new repo
4. Update:
   - `backend/go.mod` (module path)
   - Environment files (`.env.*`)
5. Push to `main`

CI will run automatically.

## Local Development
### Backend
```bash
cd backend
go run ./cmd/app
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## CI / Release
- CI runs on every push / PR
- Releases are generated from commit messages (Conventional Commits)

## Standards
See `docs/ENGINEERING_STANDARD.md` for the full production and governance baseline.

---
This repository is intended to be the **single source of truth** for new services.
