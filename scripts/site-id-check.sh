#!/usr/bin/env bash
# site-id-check.sh — block NEW site_id / siteId / X-Site-Id references during the
# site_id -> tenant_id migration (spec docs/superpowers/specs/2026-05-12-site-id-to-tenant-id-migration.md).
#
# Behaviour:
#   - Greps the tree (or staged files) for site_id / siteId / SiteID / SiteId / X-Site-Id.
#   - Compares against .site-id-migration-baseline at the repo root.
#   - Exits 1 if any match appears in a file NOT in the baseline.
#   - --update-baseline regenerates the baseline from the current tree.
#
# Skip in emergencies:  SITE_ID_LINT_SKIP=1 git commit ...
#
# Canonical copy maintained alongside D1 1.11 plan; vendored per repo so the
# hook works without network access.
set -euo pipefail

MODE=ci
BASELINE_FILE=".site-id-migration-baseline"
ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

usage() {
  cat <<EOF
Usage: site-id-check.sh [--ci|--staged|--update-baseline] [--baseline FILE]

  --ci               Scan the full working tree (default; used by CI).
  --staged           Scan only files staged for commit (used by pre-commit hook).
  --update-baseline  Regenerate the baseline file from the current tree.
  --baseline FILE    Override baseline path (default: \$REPO/.site-id-migration-baseline).
  -h, --help         Show this help.

Exit codes: 0 clean, 1 new references introduced, 2 misuse.

Set SITE_ID_LINT_SKIP=1 to bypass (will still fail in CI).
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --ci) MODE=ci; shift ;;
    --staged) MODE=staged; shift ;;
    --update-baseline) MODE=update; shift ;;
    --baseline) BASELINE_FILE="$2"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown arg: $1" >&2; usage; exit 2 ;;
  esac
done

if [[ "${SITE_ID_LINT_SKIP:-}" = "1" && "$MODE" != "update" ]]; then
  echo "site-id-check: SKIPPED (SITE_ID_LINT_SKIP=1)" >&2
  exit 0
fi

cd "$ROOT"

# Pattern matches: site_id, site-id, siteId, SiteId, SiteID, X-Site-Id, x-site-id.
# Case-insensitive on the literal "site" plus underscore/dash/camel boundary.
PATTERN='([Ss]ite[_-][Ii]d|[Ss]ite[Ii]d|[Ss]ite[Ii][Dd]|X-[Ss]ite-[Ii]d)'

# Paths always excluded (vendored / generated / fixtures).
# Using pathspec exclusions (works in BSD and GNU git-grep alike).
EXCLUDES=(
  ':!vendor/'
  ':!node_modules/'
  ':!.git/'
  ':!dist/'
  ':!build/'
  ':!.next/'
  ':!out/'
  ':!coverage/'
  ':!*.stories.*'
  ':!*.snap'
  ':!*.lock'
  ':!*.lockb'
  ':!go.sum'
  ':!pnpm-lock.yaml'
  ':!yarn.lock'
  ':!package-lock.json'
  ':!bun.lockb'
  ':!.site-id-migration-baseline'
  ':!scripts/site-id-check.sh'
)

scan_full() {
  # Emit "path" entries (one per file) — baseline is path-level, not line-level,
  # so reformatting an existing file doesn't trip the lint.
  git grep -InE "$PATTERN" -- "${EXCLUDES[@]}" 2>/dev/null \
    | awk -F: '{print $1}' \
    | LC_ALL=C sort -u
}

scan_staged() {
  # Only files staged for commit; intersect with full scan output.
  staged="$(git diff --cached --name-only --diff-filter=ACMR 2>/dev/null || true)"
  [[ -z "$staged" ]] && return 0
  # shellcheck disable=SC2086
  git grep -InE "$PATTERN" --cached -- $staged 2>/dev/null \
    | awk -F: '{print $1}' \
    | LC_ALL=C sort -u
}

case "$MODE" in
  update)
    scan_full > "$BASELINE_FILE"
    count=$(wc -l < "$BASELINE_FILE" | tr -d ' ')
    echo "site-id-check: wrote $count grandfathered paths to $BASELINE_FILE" >&2
    exit 0
    ;;
  ci|staged)
    if [[ ! -f "$BASELINE_FILE" ]]; then
      echo "site-id-check: baseline $BASELINE_FILE not found — generate with --update-baseline" >&2
      exit 2
    fi
    if [[ "$MODE" == "ci" ]]; then
      current="$(scan_full)"
    else
      current="$(scan_staged || true)"
    fi
    [[ -z "$current" ]] && exit 0
    baseline="$(LC_ALL=C sort -u "$BASELINE_FILE")"
    # New = lines in current but not in baseline.
    new="$(comm -23 <(printf '%s\n' "$current") <(printf '%s\n' "$baseline") || true)"
    if [[ -n "${new// /}" ]]; then
      echo "" >&2
      echo "site-id-check: NEW site_id / siteId references detected outside the baseline:" >&2
      echo "" >&2
      while IFS= read -r path; do
        [[ -z "$path" ]] && continue
        echo "  $path:" >&2
        git grep -nE "$PATTERN" -- "$path" 2>/dev/null | sed 's/^/    /' >&2 || true
      done <<< "$new"
      echo "" >&2
      echo "These are blocked by the site_id->tenant_id migration (spec 2026-05-12)." >&2
      echo "If the new reference is part of the migration itself, regenerate the baseline:" >&2
      echo "  ./scripts/site-id-check.sh --update-baseline" >&2
      echo "Emergency bypass:  SITE_ID_LINT_SKIP=1 git commit ..." >&2
      exit 1
    fi
    exit 0
    ;;
esac
