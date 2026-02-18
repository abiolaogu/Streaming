#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

MODE="${1:-changed}" # changed|all
CONTRACT_FILE="${CONTRACT_GATES_FILE:-scripts/ci/contract-gates.txt}"
BASE_SHA="${BASE_SHA:-}"
HEAD_SHA="${HEAD_SHA:-HEAD}"

if [[ ! -f "$CONTRACT_FILE" ]]; then
  echo "Contract gate config not found: $CONTRACT_FILE"
  exit 1
fi

configured_services() {
  awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    $0 !~ /^[[:space:]]*#/ && NF >= 3 {
      svc = trim($1)
      if (svc != "") print svc
    }
  ' "$CONTRACT_FILE" | sort -u
}

service_metadata() {
  local service="$1"
  python3 - "$CONTRACT_FILE" "$service" <<'PY'
import sys

path = sys.argv[1]
target = sys.argv[2]

with open(path, "r", encoding="utf-8") as handle:
    for raw in handle:
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        parts = [part.strip() for part in line.split("|", 2)]
        if len(parts) != 3:
            continue
        service, required, command = parts
        if service == target:
            print(f"{required}\t{command}")
            raise SystemExit(0)
raise SystemExit(1)
PY
}

all_configured="$(configured_services)"
if [[ -z "$all_configured" ]]; then
  echo "No contract gates configured in ${CONTRACT_FILE}"
  exit 1
fi

resolve_changed_services() {
  local base="$BASE_SHA"
  local head="$HEAD_SHA"

  if [[ -z "$base" || "$base" == "0000000000000000000000000000000000000000" ]]; then
    if git rev-parse --verify HEAD~1 >/dev/null 2>&1; then
      base="HEAD~1"
    else
      return 0
    fi
  fi

  git diff --name-only "$base" "$head" \
    | awk -F/ 'NF >= 3 && $1 == "services" {print $1 "/" $2}' \
    | sort -u
}

candidate_services="$(resolve_changed_services)"
if [[ "$MODE" == "all" || -z "$candidate_services" ]]; then
  candidate_services="$all_configured"
fi

echo "Running contract-test gates..."
failures=0
executed=0

while IFS= read -r service; do
  [[ -z "$service" ]] && continue

  metadata="$(service_metadata "$service" 2>/dev/null || true)"
  [[ -z "$metadata" ]] && continue
  required="${metadata%%$'\t'*}"
  command="${metadata#*$'\t'}"
  [[ -z "$required" ]] && continue
  [[ "$required" != "required" ]] && continue

  if [[ -z "$command" ]]; then
    echo "FAIL: Missing contract command for required service ${service}"
    failures=1
    continue
  fi

  if ! rg -n '^func TestContract' "$service" --glob '*_test.go' >/dev/null; then
    echo "FAIL: ${service} is required but has no TestContract* tests."
    failures=1
    continue
  fi

  echo "--> ${service}"
  if ! (
    cd "$service"
    eval "$command"
  ); then
    echo "    FAIL: contract tests failed for ${service}"
    failures=1
    continue
  fi

  executed=$((executed + 1))
  echo "    PASS"
done <<< "$candidate_services"

if [[ "$failures" -ne 0 ]]; then
  echo "Contract-test gates failed."
  exit 1
fi

if [[ "$executed" -eq 0 ]]; then
  echo "No required contract gates were applicable for this run."
else
  echo "All required contract gates passed."
fi
