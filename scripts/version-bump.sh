#!/usr/bin/env bash
# Bump semver based on Conventional Commits since last tag.
# Usage: ./scripts/version-bump.sh [--dry-run]
#
# Rules (from conventionalcommits.org):
#   BREAKING CHANGE footer or type! → major
#   feat                            → minor
#   fix / perf / refactor / ...     → patch

set -euo pipefail

DRY_RUN=false
[[ "${1:-}" == "--dry-run" ]] && DRY_RUN=true

current_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Current tag: $current_tag"

# Strip leading 'v'
version="${current_tag#v}"
IFS='.' read -r major minor patch <<< "$version"

# Collect commits since last tag
commits=$(git log "${current_tag}..HEAD" --pretty=format:"%s" 2>/dev/null || git log --pretty=format:"%s")

if [[ -z "$commits" ]]; then
  echo "No commits since $current_tag — nothing to bump."
  exit 0
fi

bump="none"

while IFS= read -r msg; do
  # Breaking: type! or BREAKING CHANGE in footer
  if echo "$msg" | grep -qE '^[a-z]+(\([^)]+\))?!:' || echo "$msg" | grep -q "BREAKING CHANGE"; then
    bump="major"
    break
  fi
  if echo "$msg" | grep -qE '^feat(\([^)]+\))?:'; then
    [[ "$bump" != "major" ]] && bump="minor"
  fi
  if echo "$msg" | grep -qE '^(fix|perf|refactor|docs|test|chore|ci|build)(\([^)]+\))?:'; then
    [[ "$bump" == "none" ]] && bump="patch"
  fi
done <<< "$commits"

if [[ "$bump" == "none" ]]; then
  echo "No releasable commits since $current_tag."
  exit 0
fi

case "$bump" in
  major) major=$((major + 1)); minor=0; patch=0 ;;
  minor) minor=$((minor + 1)); patch=0 ;;
  patch) patch=$((patch + 1)) ;;
esac

new_tag="v${major}.${minor}.${patch}"
echo "Bump: $bump → $new_tag"

if [[ "$DRY_RUN" == "true" ]]; then
  echo "(dry-run) Would tag: $new_tag"
  exit 0
fi

git tag -a "$new_tag" -m "chore(release): $new_tag"
echo "Tagged $new_tag — push with: git push origin $new_tag"
