#!/usr/bin/env bash
set -euo pipefail

# Runs only after the candidate has passed tests on all supported operating systems.
git fetch origin main --tags
if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/main)" ]; then
  echo 'main advanced after validation; retry the workflow against the new main.' >&2
  exit 1
fi

latest=$(git tag --merged HEAD --sort=-version:refname | { grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' || true; } | head -n 1)
if [ -z "$latest" ]; then
  echo 'No stable version tag found; choose the initial version manually.' >&2
  exit 1
fi

release() {
  local version=$1
  if ! gh release view "$version" >/dev/null 2>&1; then
    gh release create "$version" --verify-tag --title "$version" --generate-notes
  fi
}

if git diff --quiet -- tzlocal/tzmapping.go; then
  # Recover from a previous run that pushed the commit/tag but failed to publish.
  if [ "$(git log -1 --format=%s "$latest")" = "chore: update time zone mappings ($latest)" ] &&
     [ "$(git log -1 --format=%ae "$latest")" = '41898282+github-actions[bot]@users.noreply.github.com' ] &&
     [ "$(git diff-tree --no-commit-id --name-only -r "$latest")" = 'tzlocal/tzmapping.go' ]; then
    release "$latest"
  else
    echo 'Mappings unchanged; no commit or release needed.'
  fi
  exit 0
fi

# Finish any pending automated release before creating another version.
if [ "$(git log -1 --format=%s "$latest")" = "chore: update time zone mappings ($latest)" ]; then
  release "$latest"
fi

IFS=. read -r major minor patch <<< "${latest#v}"
version="v${major}.$((10#$minor + 1)).0"
git config user.name 'github-actions[bot]'
git config user.email '41898282+github-actions[bot]@users.noreply.github.com'
git add -- tzlocal/tzmapping.go
git commit -m "chore: update time zone mappings ($version)"
git tag "$version"
# Publish the branch and version together. Protection/race failures change neither.
git push --atomic origin HEAD:refs/heads/main "refs/tags/$version"
release "$version"
