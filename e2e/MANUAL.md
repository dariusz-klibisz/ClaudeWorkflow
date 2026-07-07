# Pre-release checklist — manual E2E scenarios

The automated adversarial suite (`e2e/run.sh`, 7 scenarios) is the release
gate's first half. Two scenarios from the validation plan (09 §2) resist
reliable headless scripting and stay MANUAL until they don't. Run both
before tagging a release.

## 1. Mid-Build kill + fresh-clone resume (`wf run adopt`)

1. Start a real `/wf:dev` run in a test repo; let it reach Build with open
   tasks, then kill the session (close the terminal mid-work).
2. Commit `.workflow/`, clone the repo to a NEW directory (fresh machine
   simulation), open a session there.
3. Expect: `wf doctor` names the log/snapshot divergence; `wf run adopt`
   re-attaches; the SessionStart injection restores the exact contract
   checklist (compare against `wf status` pre-kill). **No force needed.**

## 2. Compaction soak (obligations across ≥2 auto-compactions)

1. Start a `/wf:dev` run with a deliberately long Build (several tasks).
2. Capture the obligation baseline: `wf status > /tmp/wf-before.txt`.
3. Work until the session auto-compacts **twice** (or force `/compact` at
   phase boundaries).
4. After each compaction: `wf status > /tmp/wf-after-N.txt` and diff against
   the baseline minus legitimately-completed items.
5. Expect: no obligation lost, the injected [wf] block re-anchors the model
   (it should keep working the checklist, not re-plan from memory), and the
   ledger shows no duplicate/contradictory records around the boundaries.

## Recording the pass

Run both, then note results in the release PR/commit message:
`e2e-manual: adopt-resume PASS, compaction-soak PASS (2 compactions)`.
