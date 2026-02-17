#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

echo "Running insecure-default baseline checks..."

# Scan runtime source code for hardcoded JWT fallbacks and placeholder secrets.
# Documentation/examples are intentionally excluded from this check.
if rg -n \
  "JWT_SECRET_KEY\\s*\\|\\|\\s*['\"]secret['\"]|your-secret-key|mvp_test_secret_change_in_production" \
  services packages src \
  --glob '*.go' \
  --glob '*.js' \
  --glob '*.ts' \
  --glob '*.tsx'; then
  echo "Insecure defaults detected in runtime source code."
  exit 1
fi

echo "No insecure runtime defaults detected."
