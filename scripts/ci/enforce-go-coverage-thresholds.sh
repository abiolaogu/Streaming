#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

THRESHOLD_FILE="${GO_COVERAGE_THRESHOLD_FILE:-scripts/ci/go-coverage-thresholds.txt}"
DEFAULT_THRESHOLD="${DEFAULT_GO_COVERAGE_THRESHOLD:-0}"

if [[ ! -f "$THRESHOLD_FILE" ]]; then
  echo "Coverage threshold file not found: $THRESHOLD_FILE"
  exit 1
fi

threshold_for_module() {
  local module="$1"
  local threshold
  threshold="$(awk -v module="$module" '
    $0 !~ /^[[:space:]]*#/ && NF >= 2 && $1 == module { print $2; exit }
  ' "$THRESHOLD_FILE")"
  if [[ -z "$threshold" ]]; then
    threshold="$DEFAULT_THRESHOLD"
  fi
  echo "$threshold"
}

go_modules="$(rg --files | rg 'go\.mod$' | sed 's#/go\.mod##' | sort)"

if [[ -z "$go_modules" ]]; then
  echo "No Go modules found."
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

failures=0
echo "Running Go coverage threshold checks..."

while IFS= read -r module; do
  [[ -z "$module" ]] && continue
  threshold="$(threshold_for_module "$module")"
  profile_file="${tmp_dir}/$(echo "$module" | tr '/.' '__').cover.out"
  test_log_file="${tmp_dir}/$(echo "$module" | tr '/.' '__').test.log"

  echo "--> ${module} (threshold: ${threshold}%)"
  if ! (
    cd "$module"
    go test ./... -covermode=atomic -coverprofile="$profile_file" >"$test_log_file"
  ); then
    echo "    go test failed for ${module}"
    cat "$test_log_file"
    failures=1
    continue
  fi

  coverage="$(
    cd "$module"
    go tool cover -func="$profile_file" | awk '/^total:/ {gsub("%", "", $3); print $3}'
  )"
  coverage="${coverage:-0}"

  if ! python3 - "$coverage" "$threshold" <<'PY'
import sys

coverage = float(sys.argv[1])
threshold = float(sys.argv[2])
if coverage + 1e-9 < threshold:
    raise SystemExit(1)
PY
  then
    echo "    FAIL: coverage ${coverage}% is below threshold ${threshold}%"
    failures=1
    continue
  fi

  echo "    PASS: coverage ${coverage}%"
done <<< "$go_modules"

if [[ "$failures" -ne 0 ]]; then
  echo "Go coverage threshold checks failed."
  exit 1
fi

echo "All Go coverage thresholds satisfied."
