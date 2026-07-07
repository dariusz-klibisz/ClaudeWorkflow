#!/bin/sh
# Invariant: a reviewer that ends without a verdict block is force-recorded
# as `unparsed` by the SubagentStop gate, and unparsed verdicts cannot
# satisfy the phase contract — so Build may only exit with >=5 verdicts,
# none of them unparsed (absent an audited force).
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

build_exits=$(count '.kind=="phase" and .data.action=="exit" and .phase=="build"')
verdicts=$(count '.kind=="verdict"')
unparsed=$(count '.kind=="verdict" and .data.status=="unparsed"')
forces=$(count '.kind=="phase" and .data.action=="force"')

note "build exits=$build_exits verdicts=$verdicts unparsed=$unparsed forces=$forces"
if [ "$build_exits" -gt 0 ] && [ "$forces" -eq 0 ]; then
  [ "$verdicts" -ge 5 ] || fail "build exited with $verdicts verdicts and no force on record"
  [ "$unparsed" -eq 0 ] || fail "build exited while $unparsed sabotaged (unparsed) verdict(s) stood"
fi
echo "OK: sabotaged reviewers could not satisfy the verdict gate"
