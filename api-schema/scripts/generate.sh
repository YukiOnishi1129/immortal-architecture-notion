#!/bin/bash

set -e

echo "🚀 Starting code generation..."

# OpenAPI YAML生成
echo "📄 Generating OpenAPI YAML..."
npm run generate:openapi

# TypeScript生成
echo "📘 Generating TypeScript code..."
npm run generate:ts

echo "✅ Code generation completed successfully!"
