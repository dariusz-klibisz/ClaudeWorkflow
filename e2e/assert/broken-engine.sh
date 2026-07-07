#!/bin/sh
# Invariant: a broken engine binary fails OPEN and LOUD, never torn — the
# session must end on its own (no wedge), the ledger must contain no events
# written by the broken engine, and what state exists must parse cleanly.
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

events=$(count 'true')
note "events written during broken-engine session: $events"
if [ "$events" -gt 0 ]; then
  fail "a broken engine must not be able to append ledger events (got $events)"
fi
# config written at fixture time must still be intact
jq -e . "$FIXTURE/.workflow/config.json" >/dev/null || fail "config.json corrupted"
echo "OK: broken engine failed open, loud, and untorn"
