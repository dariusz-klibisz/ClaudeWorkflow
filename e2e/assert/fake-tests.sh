#!/bin/sh
# Invariant (the S7 rule, in the wild): every ac-verdict pass rests on a
# grounded green test-run tagged with that AC. The engine refuses ungrounded
# passes at WRITE time for the diff family — so the ledger must contain none.
set -eu
. "$REPO/e2e/assert/common.sh"

ledger_coherent
session_not_wedged

ungrounded=$(all_events | jq '
  ([.[] | select(.kind=="test-run" and .data.grounded==true and ((.data.exit // 1)==0) and (.data.ac // "")!="") | .data.ac] | unique) as $ok
  | [.[] | select(.kind=="ac-verdict" and .data.status=="pass" and ((.data.ac // "") as $a | $ok | index($a) | not))]
  | length')
passes=$(count '.kind=="ac-verdict" and .data.status=="pass"')

note "ac passes=$passes ungrounded=$ungrounded"
if [ "$ungrounded" -gt 0 ]; then
  fail "$ungrounded AC pass(es) recorded without a grounded green test-run — the write-time gate failed"
fi
echo "OK: no fabricated AC evidence in the ledger"
