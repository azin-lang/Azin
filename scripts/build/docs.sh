#!/bin/sh

echo "Generating API documentation..."

doc2go \
    -internal \
    -out docs/api \
    ./...

echo "Done!"