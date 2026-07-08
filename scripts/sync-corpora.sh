#!/bin/sh
# Maintainer-side: refresh the bundled reference corpora from their source
# repos (06 §5). Corpora ship inside the plugin because the cache-copy forbids
# ../ references. Each snapshot gets a VERSION stamp (remote + SHA + date) and
# a SHA256SUMS manifest (gen -check verifies the snapshot against it — a
# hand-edited bundled corpus is CI-visible drift, not a silent divergence).
#
# Usage:
#   scripts/sync-corpora.sh [source-parent-dir]           refresh snapshots
#   scripts/sync-corpora.sh --check [source-parent-dir]   drift check only:
#     exit 2 when a source repo's HEAD differs from the VERSION sha
#     (skips sources that aren't present locally)
set -eu

mode=sync
if [ "${1:-}" = "--check" ]; then
  mode=check
  shift
fi

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(dirname "$here")
parent="${1:-$(dirname "$root")}"
drift=0

# sha256 tool portability (linux: sha256sum, macOS: shasum -a 256)
sum_cmd() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$@"; else shasum -a 256 "$@"; fi
}

stamp_sums() {
  dst=$1
  (cd "$dst" && find . -name '*.md' | LC_ALL=C sort | while IFS= read -r f; do
    sum_cmd "$f"
  done) > "$dst/SHA256SUMS.tmp"
  mv "$dst/SHA256SUMS.tmp" "$dst/SHA256SUMS"
}

remote_of() {
  git -C "$1" remote get-url origin 2>/dev/null || echo "$1"
}

sync_one() {
  name=$1
  src=$2
  dst="$root/reference/$name"
  if [ ! -d "$src" ]; then
    echo "skip $name: $src not found" >&2
    return 0
  fi
  if [ "$mode" = check ]; then
    want=$(sed -n 's/^sha: //p' "$dst/VERSION" 2>/dev/null || echo "")
    have=$(git -C "$src" rev-parse --short HEAD 2>/dev/null || echo "unversioned")
    if [ "$want" != "$have" ]; then
      echo "DRIFT $name: bundled sha $want != source HEAD $have ($src) — run scripts/sync-corpora.sh" >&2
      drift=2
    else
      echo "ok $name: $have"
    fi
    return 0
  fi
  rm -rf "$dst"
  mkdir -p "$dst"
  # copy markdown trees, excluding VCS and tool-config files
  (cd "$src" && find . -name '*.md' ! -path './.git/*' ! -name 'AGENTS.md' ! -name 'CLAUDE.md' | while IFS= read -r f; do
    mkdir -p "$dst/$(dirname "$f")"
    cp "$f" "$dst/$f"
  done)
  sha=$(git -C "$src" rev-parse --short HEAD 2>/dev/null || echo "unversioned")
  # remote URL (not the local path): the scheduled corpora workflow clones it
  printf 'source: %s\nsha: %s\ndate: %s\n' "$(remote_of "$src")" "$sha" "$(date -u +%Y-%m-%d)" > "$dst/VERSION"
  stamp_sums "$dst"
  echo "synced $name <- $src ($sha): $(find "$dst" -name '*.md' | wc -l) files"
}

sync_one design "$parent/Design"
sync_one coding "$parent/Coding"
sync_one ux "$parent/UI_UX"

exit $drift
