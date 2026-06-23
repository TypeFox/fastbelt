#!/usr/bin/env bash

set -euo pipefail

# Require tokens for both marketplaces before doing any work.
if [[ -z "${VSCE_PAT:-}" ]]; then
  echo "VSCE_PAT is required."
  exit 1
fi

if [[ -z "${OVSX_PAT:-}" ]]; then
  echo "OVSX_PAT is required."
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Derive output file names from extension metadata.
PACKAGE_NAME="$(node -p "require('./package.json').name")"
PACKAGE_VERSION="$(node -p "require('./package.json').version")"

# Format: "<vsce-target> <GOOS> <GOARCH>"
targets=(
  "win32-x64 windows amd64"
  "win32-arm64 windows arm64"
  "linux-x64 linux amd64"
  "linux-arm64 linux arm64"
  "alpine-x64 linux amd64"
  "alpine-arm64 linux arm64"
  "darwin-x64 darwin amd64"
  "darwin-arm64 darwin arm64"
)

for target_info in "${targets[@]}"; do
  read -r target goos goarch <<<"$target_info"
  # Keep a stable filename so the exact same artifact can be
  # published to both VS Marketplace and Open VSX.
  vsix_file="${PACKAGE_NAME}-${PACKAGE_VERSION}-${target}.vsix"

  echo "==> Building and packaging for ${target} (GOOS=${goos}, GOARCH=${goarch})"
  # Rebuild from a clean dist directory for this platform.
  rm -rf dist

  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" npx vsce package --target "$target" --out "$vsix_file"

  echo "==> Publishing ${vsix_file} to VS Marketplace"
  npx --no-install vsce publish --packagePath "$vsix_file"

  echo "==> Publishing ${vsix_file} to Open VSX"
  npx --no-install ovsx publish "$vsix_file"
done

echo "Done publishing all platform-specific extensions."
