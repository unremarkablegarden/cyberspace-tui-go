#!/usr/bin/env bash
set -euo pipefail

APP="cyberspace-cli"
OUT="dist"

mkdir -p "$OUT"

echo "Building $APP..."
go build -o "$OUT/$APP" .

echo "Done: $OUT/$APP"

read -rsn1 -p "Run? [Y/n] " answer
echo
if [[ -z "$answer" || "$answer" =~ ^[Yy]$ ]]; then
  exec "$OUT/$APP"
fi
