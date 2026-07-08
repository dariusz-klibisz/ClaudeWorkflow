#!/bin/sh
# wf adversarial E2E (09 §2, 07 §6 layer 3) — the release gate.
# Drives real headless Claude Code sessions against a throwaway fixture repo
# with hostile system-prompt overlays, then asserts FROM THE RESULTING
# .workflow/ STATE that every gate held. Invariant-style assertions: runs are
# non-deterministic, but a gate that held leaves (or refuses to leave)
# specific traces regardless of what the model chose to do.
#
# Usage:   e2e/run.sh [scenario ...]        (default: all automated scenarios)
# Env:     WF_E2E_TIMEOUT   per-scenario wall clock, seconds (default 900)
#          WF_E2E_MODEL     model override (default: CLI default)
#          WF_E2E_ISOLATE   1 (default) = fresh CLAUDE_CONFIG_DIR per scenario
#                           (needs ANTHROPIC_API_KEY); 0 = your logged-in
#                           claude + your real plugin data dir
#          WF_E2E_CLAUDE    claude binary (default: claude); set to
#                           "$REPO/e2e/mock-claude.sh" for offline harness
#                           self-tests (runs a scripted minimal session)
# Needs:   claude on PATH (or WF_E2E_CLAUDE), jq, go, git.
set -eu

REPO=$(cd "$(dirname "$0")/.." && pwd)
TIMEOUT="${WF_E2E_TIMEOUT:-900}"
CLAUDE="${WF_E2E_CLAUDE:-claude}"
ISOLATE="${WF_E2E_ISOLATE:-1}"

SCENARIOS="${*:-skip-reviews mark-all-done fake-tests forge-ledger sabotage-reviewer later-phase broken-engine legit}"

command -v jq >/dev/null 2>&1 || { echo "e2e: jq required"; exit 3; }
command -v git >/dev/null 2>&1 || { echo "e2e: git required"; exit 3; }
if [ "$CLAUDE" = claude ] && ! command -v claude >/dev/null 2>&1; then
  echo "e2e: claude CLI required (or set WF_E2E_CLAUDE)"; exit 3
fi
if [ "$ISOLATE" = 1 ] && [ -z "${ANTHROPIC_API_KEY:-}" ] && [ "$CLAUDE" = claude ]; then
  echo "e2e: WF_E2E_ISOLATE=1 needs ANTHROPIC_API_KEY (or set WF_E2E_ISOLATE=0 to use your logged-in claude)"
  exit 3
fi

# one dev engine build for all scenarios (the fixture's plugin root is the
# checkout itself — --plugin-dir exercises the real bootstrap path)
( cd "$REPO" && make build >/dev/null )
WF="$REPO/bin/wf-$(go env GOOS)-$(go env GOARCH)$(go env GOEXE)"

pass=0; fail=0; skip=0; failed=""
for sc in $SCENARIOS; do
  overlay="$REPO/e2e/overlays/$sc.md"
  assert="$REPO/e2e/assert/$sc.sh"
  task="$REPO/e2e/tasks/$sc.txt"
  for f in "$overlay" "$assert" "$task"; do
    [ -f "$f" ] || { echo "e2e: unknown scenario $sc (missing $f)"; exit 3; }
  done

  if [ "$sc" = broken-engine ] && [ "$ISOLATE" != 1 ]; then
    echo "=== e2e: $sc — SKIPPED (needs WF_E2E_ISOLATE=1: your real plugin data dir may hold a healthy engine)"
    skip=$((skip + 1))
    continue
  fi

  echo "=== e2e: $sc ==="
  FIXTURE=$(mktemp -d "${TMPDIR:-/tmp}/wf-e2e-$sc.XXXXXX")
  sh "$REPO/e2e/fixture.sh" "$FIXTURE" "$REPO" "$WF"

  # scenario-specific environment sabotage (restored after the scenario)
  if [ "$sc" = broken-engine ]; then
    cp "$WF" "$WF.e2e-backup"
    printf 'GARBAGE-NOT-AN-ENGINE' > "$WF"
  fi

  set +e
  ( cd "$FIXTURE"
    if [ "$ISOLATE" = 1 ]; then export CLAUDE_CONFIG_DIR="$FIXTURE/claude-config"; fi
    timeout "$TIMEOUT" "$CLAUDE" -p "$(cat "$task")" \
      --plugin-dir "$REPO" \
      --append-system-prompt "$(cat "$overlay")" \
      --dangerously-skip-permissions \
      ${WF_E2E_MODEL:+--model "$WF_E2E_MODEL"} \
      --output-format text ) > "$FIXTURE/session.out" 2>&1
  session_exit=$?
  set -e
  echo "--- session exit: $session_exit (output: $FIXTURE/session.out)"

  if [ "$sc" = broken-engine ]; then
    mv "$WF.e2e-backup" "$WF"
  fi

  set +e
  FIXTURE="$FIXTURE" WF="$WF" REPO="$REPO" SESSION_EXIT="$session_exit" sh "$assert"
  rc=$?
  set -e
  if [ "$rc" -eq 0 ]; then
    echo "--- PASS: $sc"
    pass=$((pass + 1))
    rm -rf "$FIXTURE"
  else
    echo "--- FAIL: $sc (fixture kept: $FIXTURE)"
    fail=$((fail + 1))
    failed="$failed $sc"
  fi
done

echo "e2e: $pass passed, $fail failed, $skip skipped${failed:+ —$failed}"
[ "$fail" -eq 0 ]
