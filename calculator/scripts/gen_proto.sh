#!/usr/bin/env bash
set -euo pipefail

PROTO_DIR="proto"
OUT_DIR="internal/gen"

mkdir -p "$OUT_DIR"

# Ensure required plugins are installed
if ! command -v protoc-gen-go >/dev/null || ! command -v protoc-gen-go-grpc >/dev/null; then
  echo "Error: protoc-gen-go and/or protoc-gen-go-grpc not found in PATH."
  echo "Install with:"
  echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
  exit 1
fi

# Find all .proto files and generate code
find "$PROTO_DIR" -name '*.proto' | while read -r proto_file; do
  protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$OUT_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
    "$proto_file"
done