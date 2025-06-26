#!/bin/bash
set -e

PROTO_DIR="internal/transport/grpc/protos"
OUT_DIR="$PROTO_DIR"

function build() {
  echo "[PROTO BUILD]"
  for file in $PROTO_DIR/*.proto; do
    echo "Generating for $file..."
    protoc --go_out=$OUT_DIR --go_opt=paths=source_relative \
           --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
           "$file"
  done
  echo "[DONE]"
}

function clean() {
  echo "[PROTO CLEAN]"
  find "$OUT_DIR" -type f \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -delete
  echo "[DONE]"
}

case "$1" in
  build)
    build
    ;;
  clean)
    clean
    ;;
  *)
    echo "Usage: $0 {build|clean}"
    exit 1
    ;;
esac 