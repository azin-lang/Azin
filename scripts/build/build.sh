#!/bin/bash

AZC_SOURCE=./cmd/compiler
BUILD_DIR=./build

mkdir -p "$BUILD_DIR"

go build -o "$BUILD_DIR/azc" "$AZC_SOURCE"