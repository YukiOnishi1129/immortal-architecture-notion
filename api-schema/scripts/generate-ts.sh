#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(cd -- "${SCRIPT_DIR}/.." && pwd)
REPO_ROOT=$(cd -- "${PROJECT_ROOT}/.." && pwd)

TS_TARGET_DIR="${REPO_ROOT}/frontend/src/external/client/api/generated"

echo "📘 Generating TypeScript code from OpenAPI..."

rm -rf "${TS_TARGET_DIR}"

pnpm exec openapi-generator-cli generate \
  -i "${PROJECT_ROOT}/generated/openapi.yaml" \
  -g typescript-fetch \
  -o "${TS_TARGET_DIR}" \
  --additional-properties=withSeparateModelsAndApi=false,supportsES6=true,useSingleRequestParameter=true

echo "✅ TypeScript code generated at: ${TS_TARGET_DIR}"
