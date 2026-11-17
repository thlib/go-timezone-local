#!/usr/bin/env bash
set -euo pipefail

script="$(cd "$(dirname "$0")" && pwd)/release-tzmapping.sh"
fixtures=$(mktemp -d)
trap 'rm -rf "$fixtures"' EXIT
mkdir "$fixtures/bin"
cat > "$fixtures/bin/gh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
  'release view') test -f "$RELEASE_STATE/$3" ;;
  'release create')
    if [ "${FAIL_RELEASE:-0}" = 1 ]; then exit 1; fi
    touch "$RELEASE_STATE/$3"
    echo "$3" >> "$RELEASE_STATE/created"
    ;;
  *) exit 99 ;;
esac
EOF
chmod +x "$fixtures/bin/gh"
export PATH="$fixtures/bin:$PATH"
# Keep fixtures independent of the developer's global signing/identity settings.
export GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null
export GIT_AUTHOR_NAME=Test GIT_AUTHOR_EMAIL=test@example.test
unset GIT_COMMITTER_NAME GIT_COMMITTER_EMAIL

fixture() {
  local name=$1
  mkdir "$fixtures/$name"
  cd "$fixtures/$name"
  git init -q --bare remote.git
  git init -q -b main work
  cd work
  git config user.name Test
  git config user.email test@example.test
  mkdir tzlocal
  echo original > tzlocal/tzmapping.go
  git add .
  git commit -qm initial
  git tag v0.0.9
  git tag v9.0.0-rc.1
  git remote add origin ../remote.git
  git push -q -u origin main --tags
  export RELEASE_STATE="$fixtures/$name/releases"
  mkdir "$RELEASE_STATE"
  base=$(git rev-parse HEAD)
}

fixture unchanged
bash "$script"
test "$(git rev-parse HEAD)" = "$base"
test ! -e "$RELEASE_STATE/created"

fixture release
echo updated > tzlocal/tzmapping.go
# The workflow sets identity through Git config, so remove the fixture's author override.
unset GIT_AUTHOR_NAME GIT_AUTHOR_EMAIL
bash "$script"
test "$(git --git-dir=../remote.git rev-parse main)" = "$(git rev-parse v0.1.0)"
test "$(cat "$RELEASE_STATE/created")" = v0.1.0
bash "$script"
test "$(wc -l < "$RELEASE_STATE/created")" -eq 1

fixture retry
echo updated > tzlocal/tzmapping.go
if FAIL_RELEASE=1 bash "$script"; then echo 'Expected release failure' >&2; exit 1; fi
test "$(git --git-dir=../remote.git rev-parse main)" = "$(git rev-parse v0.1.0)"
test ! -e "$RELEASE_STATE/created"
bash "$script"
test "$(cat "$RELEASE_STATE/created")" = v0.1.0
test "$(git tag -l v0.2.0)" = ''

fixture stale
echo unrelated > other
git add other
git commit -qm advance
git push -q origin main
git checkout -q --detach "$base"
echo updated > tzlocal/tzmapping.go
if bash "$script"; then echo 'Expected stale main rejection' >&2; exit 1; fi
test ! -e "$RELEASE_STATE/created"
test "$(git tag -l v0.1.0)" = ''

fixture atomic
cat > ../remote.git/hooks/update <<'EOF'
#!/usr/bin/env bash
case "$1" in refs/tags/*) exit 1 ;; esac
EOF
chmod +x ../remote.git/hooks/update
echo updated > tzlocal/tzmapping.go
if bash "$script"; then echo 'Expected protected tag rejection' >&2; exit 1; fi
test "$(git --git-dir=../remote.git rev-parse main)" = "$base"
test ! -e "$RELEASE_STATE/created"

echo 'Release tests passed: no-op, minor bump, idempotence, retry, stale main, atomic rejection.'
