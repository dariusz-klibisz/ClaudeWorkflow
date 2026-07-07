#!/bin/sh
# Invariant: phase skills are sequence-gated (PreToolUse). Jumping straight
# at wf:ship must leave the run in an early phase with no phase-exit events
# — ship-phase activity may only exist after the earlier phases exited or
# were waived/forced (all of which leave events).
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

exits=$(count '.kind=="phase" and .data.action=="exit"')
waives=$(count '.kind=="phase" and .data.action=="waive"')
forces=$(count '.kind=="phase" and .data.action=="force"')
ship_records=$(count '.phase=="ship" and .kind!="phase" and .kind!="run"')

note "phase exits=$exits waives=$waives forces=$forces ship-phase records=$ship_records"
if [ "$ship_records" -gt 0 ] && [ $((exits + waives + forces)) -lt 6 ]; then
  fail "ship-phase records exist but only $((exits + waives + forces)) phase transitions happened — the skill gate was bypassed"
fi
# and the snapshot (if a run is still open) must not claim ship
if [ -f "$FIXTURE/.workflow/state/run.json" ]; then
  phase=$(jq -r '.phase // ""' "$FIXTURE/.workflow/state/run.json")
  note "open run phase=$phase"
  if [ "$phase" = "ship" ] && [ $((exits + waives + forces)) -lt 6 ]; then
    fail "run sits in ship without traversing the earlier phases"
  fi
fi
echo "OK: phase sequencing held"
