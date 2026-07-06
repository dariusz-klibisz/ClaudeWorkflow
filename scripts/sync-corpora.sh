#!/bin/sh
# Maintainer-side: refresh the bundled reference corpora from their source
# repos (06 §5). Corpora ship inside the plugin because the cache-copy forbids
# ../ references. Each snapshot gets a VERSION stamp (source SHA + date).
#
# Usage: scripts/sync-corpora.sh [source-parent-dir]   (default: ../..)
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(dirname "$here")
parent="${1:-$(dirname "$root")}"

sync_one() {
  name=$1
  src=$2
  dst="$root/reference/$name"
  if [ ! -d "$src" ]; then
    echo "skip $name: $src not found" >&2
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
  printf 'source: %s\nsha: %s\ndate: %s\n' "$src" "$sha" "$(date -u +%Y-%m-%d)" > "$dst/VERSION"
  echo "synced $name <- $src ($sha): $(find "$dst" -name '*.md' | wc -l) files"
}

sync_one design "$parent/Design"
sync_one coding "$parent/Coding"
sync_one ux "$parent/UI_UX"
