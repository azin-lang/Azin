#!/usr/bin/env bash
# Usage: ./build.sh

set -euo pipefail

readonly BUILD_DIR="build"
readonly OUTPUT="${BUILD_DIR}/azin_lsp"
readonly SOURCE="./modules/lsp/cmd/lsp"

mkdir -p "$BUILD_DIR"

echo "Building Azin compiler..."
go build \
    -trimpath \
    -o "$OUTPUT" \
    "$SOURCE"

echo "Done: $OUTPUT"
