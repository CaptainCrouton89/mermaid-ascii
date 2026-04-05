#!/bin/bash
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PYTHON_DIR="$(dirname "$SCRIPT_DIR")"
GO_ROOT="$(dirname "$PYTHON_DIR")"
TARGET="$PYTHON_DIR/src/mermaid_ascii/mermaid-ascii"

cd "$GO_ROOT"
go build -o "$TARGET" .
chmod +x "$TARGET"
