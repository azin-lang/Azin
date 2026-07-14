#!/usr/bin/env bash
# Usage: ./build.sh

set -euo pipefail

readonly BUILD_DIR="build"
readonly OUTPUT="${BUILD_DIR}/azc"
readonly SOURCE="./modules/compiler/cmd/azc"

mkdir -p "$BUILD_DIR"

echo "Building Azin compiler..."
go build \
    -trimpath \
    -o "$OUTPUT" \
    "$SOURCE"

echo "Done: $OUTPUT"
