# Git Hooks (SDK Minimal Variant)

Lightweight git hooks for this SDK. This is the minimal SDK variant
of the vigo-core-games hooks pattern — SDKs ship as libraries, so we
skip the heavier checks (race detection, security scan, gitleaks)
that live in the service-repo variant.

## 🔧 Installation

```bash
./githooks/install.sh
```

Sets `core.hooksPath=githooks` so git points directly at the in-repo
directory. Updates propagate via `git pull` — no re-run needed.

## 📋 Hooks

### pre-commit
Fast (<1s) formatting check on staged Go files:
- `gofmt -s -l` on every staged `.go` file
- `goimports -l` if available

### pre-push
Minimal SDK suite:
- `go vet ./...`
- `go test ./...`
- `go build ./...`

No race detection, no gitleaks, no gosec — those live in the
service-repo variant (see `vigo-core-games/githooks/`).

### commit-msg
Validates conventional-commits format:
- `feat: add new feature`
- `fix: resolve bug`
- `docs: update README`
- etc.

## 🚫 Bypassing

```bash
git commit --no-verify   # skip pre-commit + commit-msg
git push --no-verify     # skip pre-push
```

Use sparingly.
