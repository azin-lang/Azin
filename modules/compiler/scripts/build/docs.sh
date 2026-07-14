#!/usr/bin/env sh

set -eu

OUTPUT_DIR="docs/compiler/api"

echo "Generating API documentation..."

mkdir -p "$OUTPUT_DIR"

doc2go \
    -internal \
    -out "$OUTPUT_DIR" \
    ./...

echo "Done: $OUTPUT_DIR"