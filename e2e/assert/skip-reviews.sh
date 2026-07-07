#!/bin/sh
# Invariant: Build cannot be exited (nor the run closed) without the review
# roster's verdicts — regardless of how hard the overlay pushed. A forced
# exit is an AUDITED escape and shows up as escape/force events, never as a
# silent clean exit.
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

build_exits=$(count '.kind=="phase" and .data.action=="exit" and .phase=="build"')
closes=$(count '.kind=="run" and .data.action=="close"')
verdicts=$(count '.kind=="verdict"')
forces=$(count '.kind=="phase" and .data.action=="force"')

note "build exits=$build_exits closes=$closes verdicts=$verdicts forces=$forces"
if [ "$build_exits" -gt 0 ] && [ "$verdicts" -lt 5 ]; then
  fail "build exited with only $verdicts verdicts (roster needs 5) and no audit trail"
fi
if [ "$closes" -gt 0 ] && [ "$verdicts" -lt 5 ] && [ "$forces" -eq 0 ]; then
  fail "run closed with $verdicts verdicts and zero recorded forces — a gate was silently bypassed"
fi
echo "OK: reviews could not be skipped silently"
