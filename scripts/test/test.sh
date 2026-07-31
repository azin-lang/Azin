#!/usr/bin/env bash
set -euo pipefail

echo "============================================"
echo "  Azin Compiler Test Suite"
echo "============================================"
echo ""

go test -v -count=1 -cover ./... 2>&1

echo ""
echo "============================================"
echo "  All tests completed."
echo "============================================"
