---
name: verify
description: wf phase 6 (Verify) — independent per-AC grounded verification, confirmation reviews, loop records on failure. Invoked via /wf:dev when verify is the active phase.
---

# /wf:verify — Verify (interactive exit)

Author and verifier are different stances: verify against the requirement,
not against what Build intended.

Contract first:
- Per AC: `wf record ac-verdict ac=AC-1 status=pass|fail|deferred`
  - `pass` is REFUSED unless a grounded green test-run tagged with that AC
    exists (run the verification command from the strategy; capture is
    automatic) — passing is earned, not asserted
  - `deferred` requires `wf approve deferral` first
  - `fail` → write the loop: `wf loop --ac AC-1 --cause slip|design|plan
    --evidence "expected X, observed Y"` — cause=slip re-opens Build;
    design/plan re-open those phases. Caps: 10 loops/run, 2 slip-loops/AC
    (the 3rd must name a structural cause)
- Confirmation verdicts (diff/artifact): `@wf:adversary` (red-team) and
  `@wf:design-conformance-reviewer`; ux projects add `@wf:ux-reviewer`
- diff: run the secret scan (e.g. `gitleaks detect`) — auto-captured as
  category=secret-scan; a filtered or exit-less run never grounds
- intent deploy: `wf record smoke-run cmd="…" exit=0 target="…"` +
  `wf record rollback-readiness --json '{"procedure":"…","trigger":"…"}'`
- assessment: the findings report exists and is not a stub
  (`wf record artifact … role=deliverable-report status=present`);
  intent investigate: origin attribution recorded

`wf phase exit` only when every AC has a grounded verdict and no fail is
undispositioned.
