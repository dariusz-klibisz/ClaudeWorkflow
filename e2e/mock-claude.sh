#!/bin/sh
# Offline stand-in for the claude CLI — harness self-tests only. Plays a
# minimal scripted "session" so e2e/run.sh's fixture → session → assert
# plumbing can be verified without network, auth, or token spend:
#   WF_E2E_CLAUDE="$REPO/e2e/mock-claude.sh" e2e/run.sh legit later-phase broken-engine
# It is NOT a substitute for the live adversarial pass (e2e/MANUAL.md).
set -eu

REPO=$(cd "$(dirname "$0")/.." && pwd)
WF="$REPO/bin/wf-$(go env GOOS)-$(go env GOARCH)$(go env GOEXE)"
export CLAUDE_PLUGIN_ROOT="$REPO"

# first arg after -p is the task prompt
task=""
prev=""
for a in "$@"; do
  if [ "$prev" = "-p" ] || [ "$prev" = "--print" ]; then task="$a"; fi
  prev="$a"
done

case "$task" in
  *Investigate*)
    # the compliant minimal path: start a run, do some Frame work
    "$WF" run start --family assessment --intent investigate >/dev/null
    "$WF" record classification family=assessment intent=investigate restated="mock investigation" >/dev/null
    echo "mock: started assessment run and recorded classification"
    ;;
  *"Ship the current work"*)
    # a hostile mock would jump at ship — the real gate is PreToolUse, which
    # a mock cannot exercise; it simply does nothing (state stays empty)
    echo "mock: refused to fake ship-phase work"
    ;;
  *)
    # broken-engine (and anything else): try one wf call; a corrupt engine
    # must fail without touching the ledger
    "$WF" run start --family diff --intent new >/dev/null 2>&1 || echo "mock: wf failed (expected under broken-engine)"
    ;;
esac
exit 0
