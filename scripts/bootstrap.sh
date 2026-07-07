#!/bin/sh
# wf bootstrap (SessionStart): install the platform engine binary into
# ${CLAUDE_PLUGIN_DATA}/bin/wf — the stable path every hook references
# (07 §4). Runtime-free by design (POSIX sh; the PowerShell twin covers
# native Windows). Re-runs are cheap no-ops on version match.
set -eu

root="${CLAUDE_PLUGIN_ROOT:-}"
data="${CLAUDE_PLUGIN_DATA:-}"
[ -n "$root" ] || { echo "[wf bootstrap] CLAUDE_PLUGIN_ROOT unset" >&2; exit 0; }
[ -n "$data" ] || { echo "[wf bootstrap] CLAUDE_PLUGIN_DATA unset" >&2; exit 0; }

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux*)  os=linux ;;
  darwin*) os=darwin ;;
  mingw*|msys*|cygwin*) os=windows ;;
esac
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
esac

src="$root/bin/wf-$os-$arch"
[ "$os" = windows ] && src="$src.exe"

if [ ! -f "$src" ]; then
  # dev checkout fallback: build from source when Go is available. Stamp the
  # version from plugin.json — an unstamped build reports "0.1.0-dev" and
  # masquerades as ancient code in doctor/status output.
  if [ -d "$root/engine" ] && command -v go >/dev/null 2>&1; then
    ver="0.0.0"
    if [ -f "$root/.claude-plugin/plugin.json" ]; then
      ver=$(sed -n 's/.*"version"[^"]*"\([^"]*\)".*/\1/p' "$root/.claude-plugin/plugin.json" | head -1)
    fi
    echo "[wf bootstrap] no bundled binary for $os/$arch — building $ver+src from source" >&2
    ( cd "$root/engine" && go build -ldflags "-X main.Version=$ver+src" -o "$src" ./cmd/wf )
  else
    echo "[wf bootstrap] no engine binary for $os/$arch under $root/bin — wf gates will fail open" >&2
    exit 0
  fi
fi

# Version stamp: bin/VERSION when it ships (release/dev builds); otherwise
# fall back to the binary's checksum so git installs (which .gitignore the
# VERSION file) still get no-op re-runs and a written $data/bin/VERSION.
want=""
[ -f "$root/bin/VERSION" ] && want=$(cat "$root/bin/VERSION")
if [ -z "$want" ] && command -v sha256sum >/dev/null 2>&1; then
  want="sha256:$(sha256sum "$src" | cut -c1-16)"
fi
have=""
[ -f "$data/bin/VERSION" ] && have=$(cat "$data/bin/VERSION")
if [ -n "$want" ] && [ "$want" = "$have" ] && [ -x "$data/bin/wf" ]; then
  exit 0
fi

# checksum verification when the sums file ships (release builds); a sums
# file with no entry for this platform is skipped (not a refusal) — only an
# entry that exists and MISmatches blocks the install.
if [ -f "$root/bin/SHA256SUMS" ] && command -v sha256sum >/dev/null 2>&1; then
  if grep -q " $(basename "$src")\$" "$root/bin/SHA256SUMS"; then
    ( cd "$root/bin" && grep " $(basename "$src")\$" SHA256SUMS | sha256sum -c - >/dev/null ) || {
      echo "[wf bootstrap] checksum mismatch for $(basename "$src") — refusing to install" >&2
      exit 0
    }
  else
    echo "[wf bootstrap] SHA256SUMS has no entry for $(basename "$src") — skipping verification" >&2
  fi
fi

mkdir -p "$data/bin"
cp "$src" "$data/bin/wf.tmp"
chmod +x "$data/bin/wf.tmp"
mv "$data/bin/wf.tmp" "$data/bin/wf"
if [ "$os" = windows ]; then cp "$data/bin/wf" "$data/bin/wf.exe"; fi
if [ -n "$want" ]; then printf '%s' "$want" > "$data/bin/VERSION"; fi
echo "[wf bootstrap] installed wf ($os/$arch) -> $data/bin/wf" >&2
exit 0
