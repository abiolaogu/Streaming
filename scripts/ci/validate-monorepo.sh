#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

TARGET="${1:-all}"

run_go_checks() {
  echo "==> Running Go module checks"
  go_modules="$(rg --files | rg 'go\.mod$' | sed 's#/go\.mod##' | sort)"

  while IFS= read -r module; do
    [[ -z "$module" ]] && continue
    echo "--> go test ${module}"
    (
      cd "$module"
      go test ./...
    )
  done <<< "$go_modules"
}

run_web_checks() {
  echo "==> Running web checks"
  npm ci
  npm run check
}

run_node_service_checks() {
  echo "==> Running Node service checks"
  local node_services=(
    "services/websocket-service"
    "services/notification-service"
  )

  for service in "${node_services[@]}"; do
    if [[ -f "${service}/package.json" ]]; then
      echo "--> npm checks ${service}"
      (
        cd "$service"
        npm install --no-audit --no-fund
        npm run lint --if-present
        npm run test --if-present
      )
    fi
  done
}

run_python_checks() {
  echo "==> Running Python syntax checks"
  python_targets="$(rg --files | rg 'requirements\.txt$' | sed 's#/requirements\.txt##' | sort)"

  while IFS= read -r target; do
    [[ -z "$target" ]] && continue
    echo "--> python compileall ${target}"
    python3 -m compileall -q "$target"
  done <<< "$python_targets"
}

run_security_baseline() {
  echo "==> Running security baseline"
  ./scripts/ci/check-no-insecure-defaults.sh
}

case "$TARGET" in
  go)
    run_go_checks
    ;;
  web)
    run_web_checks
    ;;
  node)
    run_node_service_checks
    ;;
  python)
    run_python_checks
    ;;
  security)
    run_security_baseline
    ;;
  all)
    run_go_checks
    run_web_checks
    run_node_service_checks
    run_python_checks
    run_security_baseline
    ;;
  *)
    echo "Unknown target '${TARGET}'. Use one of: go|web|node|python|security|all"
    exit 1
    ;;
esac
