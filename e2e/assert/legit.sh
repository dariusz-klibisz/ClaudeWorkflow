#!/bin/sh
# Invariant: a compliant session in a clean fixture suffers no spurious
# blocks — a run gets started and Frame work happens; whatever end state it
# reached (open, parked, closed) is coherent; report signals compute; a
# closed run has its signals snapshot; nothing was forced.
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

starts=$(count '.kind=="run" and ((.data.action=="start") or (.data.action=="branch"))')
frame=$(count '.phase=="frame" and .kind!="phase" and .kind!="run"')
forces=$(count '.kind=="phase" and .data.action=="force"')
closes=$(count '.kind=="run" and .data.action=="close"')

note "run starts=$starts frame records=$frame forces=$forces closes=$closes"
[ "$starts" -ge 1 ] || fail "no run was started — the benign path is blocked somewhere"
[ "$frame" -ge 1 ] || fail "no Frame work was recorded"
[ "$forces" -eq 0 ] || fail "a compliant session needed $forces force(s) — spurious blocking"

# report must compute over whatever state resulted
( cd "$FIXTURE" && CLAUDE_PLUGIN_ROOT="$REPO" "$WF" report --json >/dev/null ) || fail "wf report --json failed"

# a closed run must carry its signals snapshot (03 §4.7)
if [ "$closes" -ge 1 ]; then
  ls "$FIXTURE"/.workflow/runs/*/signals.json >/dev/null 2>&1 || fail "closed run without signals.json"
fi
echo "OK: compliant path ran without spurious blocks"
