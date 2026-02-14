#!/usr/bin/env bash
# next-version.sh — Compute the next semver tag for this repo.
# Usage: ./scripts/next-version.sh [major|minor|patch]
set -euo pipefail

BUMP="${1:-patch}"
PREFIX="v"

LATEST=$(git tag --sort=-v:refname | grep -E "^${PREFIX}[0-9]+\.[0-9]+\.[0-9]+$" | head -1 || true)

if [ -z "$LATEST" ]; then
  echo "${PREFIX}0.1.0"
  exit 0
fi

VERSION="${LATEST#$PREFIX}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

case "$BUMP" in
  major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
  patch) PATCH=$((PATCH + 1)) ;;
  *) echo "Usage: $0 [major|minor|patch]" >&2; exit 1 ;;
esac

echo "${PREFIX}${MAJOR}.${MINOR}.${PATCH}"
