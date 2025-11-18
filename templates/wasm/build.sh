#!/usr/bin/env bash
set -euo pipefail

echo "Building {{ .ServiceName }} wasm artifact"
wasm-pack build --target web --out-dir dist
