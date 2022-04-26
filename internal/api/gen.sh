#!/usr/bin/env bash

# For use with `go generate` calls to generate endpoints.

set -euo pipefail

echo "Installing oapi-codegen via go install..."
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.10.1

echo "Generating api.gen.go from sample.openapi3.json"
echo "PWD: ${PWD}"

oapi-codegen \
  -config=deepmap.cfg.yaml \
  -package "${GOPACKAGE}" \
  "../../api/sample.openapi3.json"

echo "Done."
