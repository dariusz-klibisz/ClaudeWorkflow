#!/bin/sh
# Invariant: no task reaches done/completed without red→green evidence, a
# recorded waiver, or a recorded force — fake completions leave no legal
# trace, honest escapes leave an audited one.
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

completed=$(count '.kind=="task" and ((.data.status // "")=="done" or (.data.status // "")=="completed")')
greens=$(count '.kind=="test-run" and (.data.grounded==true) and ((.data.exit // 1)==0)')
reds=$(count '.kind=="test-run" and (.data.grounded==true) and ((.data.exit // 0)!=0)')
waivers=$(count '.kind=="waiver"')
forces=$(count '.kind=="phase" and .data.action=="force"')

note "completed-task records=$completed grounded greens=$greens reds=$reds waivers=$waivers forces=$forces"
if [ "$completed" -gt 0 ] && [ "$greens" -eq 0 ] && [ "$waivers" -eq 0 ] && [ "$forces" -eq 0 ]; then
  fail "$completed task(s) marked done with zero grounded test evidence and zero audited escapes"
fi
echo "OK: tasks could not be fake-completed silently"
