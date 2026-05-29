#!/bin/bash
#
# install.sh — wire the in-repo `githooks/` directory as the active
# hooks path for this clone.
#
# Sets `core.hooksPath` so git uses githooks/ directly — pulls pick
# up hook updates automatically, no re-run needed.
#
# Idempotent: re-runs are safe.

set -e

if [ ! -d ".git" ]; then
    echo "❌ Not in a Git repository. Run from the repo root."
    exit 1
fi

if [ ! -d "githooks" ]; then
    echo "❌ githooks/ directory not found."
    exit 1
fi

echo "🔧 Wiring core.hooksPath = githooks for this clone..."
git config core.hooksPath githooks
chmod +x githooks/pre-commit githooks/pre-push githooks/commit-msg 2>/dev/null || true

# Clean up any stale copies in .git/hooks/ from a previous installer.
for old in .git/hooks/pre-commit .git/hooks/pre-push .git/hooks/commit-msg; do
    if [ -f "$old" ] && [ ! -L "$old" ]; then
        rm -f "$old"
    fi
done

echo "✅ Hooks active. The active versions are now in githooks/ — pulls"
echo "   pick up updates automatically; no re-run of this script needed."
echo
echo "📋 Active hooks (SDK minimal variant):"
echo "   pre-commit  — fast gofmt + goimports check on staged Go files"
echo "   pre-push    — go vet + go test + go build"
echo "   commit-msg  — conventional commit format validation"
echo
echo "Skip pre-commit ad-hoc with: git commit --no-verify"
