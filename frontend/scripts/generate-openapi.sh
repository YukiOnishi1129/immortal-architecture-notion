#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)
FRONTEND_DIR=$(cd -- "${SCRIPT_DIR}/.." && pwd)
REPO_ROOT=$(cd -- "${FRONTEND_DIR}/.." && pwd)

OPENAPI_SPEC="${REPO_ROOT}/api-schema/generated/openapi.yaml"
TARGET_DIR="${FRONTEND_DIR}/src/external/client/api/generated"
API_SCHEMA_BIN="${REPO_ROOT}/api-schema/node_modules/.bin/openapi-generator-cli"

if [ ! -f "${OPENAPI_SPEC}" ]; then
  echo "❌ OpenAPI spec not found at ${OPENAPI_SPEC}"
  echo "   Run the TypeSpec/OpenAPI generation pipeline first (e.g. pnpm build in api-schema)."
  exit 1
fi

if [ -x "${API_SCHEMA_BIN}" ]; then
  GENERATOR_CMD="${API_SCHEMA_BIN}"
elif command -v openapi-generator-cli >/dev/null 2>&1; then
  GENERATOR_CMD="$(command -v openapi-generator-cli)"
else
  echo "❌ openapi-generator-cli not found."
  echo "   Install it in api-schema (pnpm install) or globally, then retry."
  exit 1
fi

echo "🌀 Generating TypeScript API client from ${OPENAPI_SPEC}"
rm -rf "${TARGET_DIR}"

"${GENERATOR_CMD}" generate \
  -i "${OPENAPI_SPEC}" \
  -g typescript-fetch \
  -o "${TARGET_DIR}" \
  --additional-properties=withSeparateModelsAndApi=false,supportsES6=true,useSingleRequestParameter=true

echo "✅ TypeScript client generated in ${TARGET_DIR}"
