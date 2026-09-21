#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(cd -- "${SCRIPT_DIR}/.." && pwd)

OPENAPI_SPEC="${PROJECT_ROOT}/../api-schema/generated/openapi.yaml"
CONFIG_FILE="${PROJECT_ROOT}/oapi-codegen.yaml"

if ! command -v oapi-codegen >/dev/null 2>&1; then
  ALT_BIN="$(go env GOPATH 2>/dev/null)/bin/oapi-codegen"
  if [ ! -x "${ALT_BIN}" ]; then
    echo "oapi-codegen not found. Install via 'go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest'"
    exit 1
  fi
  OAPI_BIN="${ALT_BIN}"
else
  OAPI_BIN="$(command -v oapi-codegen)"
fi

echo "🚧 Generating Go server/types from ${OPENAPI_SPEC}"
"${OAPI_BIN}" -config "${CONFIG_FILE}" "${OPENAPI_SPEC}"
echo "✅ Generated files written under ${PROJECT_ROOT}/internal/generated/openapi"
