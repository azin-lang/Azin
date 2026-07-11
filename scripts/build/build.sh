#!/bin/bash
# Usage: ./build.sh

AZC_SOURCE=./cmd/compiler

go build -o azc $AZC_SOURCE
