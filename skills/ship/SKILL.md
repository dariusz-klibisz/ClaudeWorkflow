---
name: ship
description: wf phase 7 (Ship) — deliver the package, resolve the trace findings, lessons, archive, close. Invoked via /wf:dev when ship is the active phase.
---

# /wf:ship — Ship (interactive)

Contract first:
- Trace: `wf trace` is engine-computed in M2; until then `wf status` lists
  the open items (followups, waivers, forces) — resolve or disposition each:
  open followups become tasks now or are carried:
  `wf record followup updates=<id> status=next-run`
- `@wf:auditor` verdict over the resolved state (HIGH findings block close)
- diff: record the delivery package —
  `wf record artifact path=<PR-or-release-ref> role=delivery status=present`;
  intent deploy additionally `role=delivery-manifest` (target, config diff,
  rollout method)
- Lessons: propose what the run taught
  (`wf record lesson text="…" status=proposed [check="…"]`), the user
  accepts or rejects each (`wf approve lesson` +
  `wf record lesson updates=<id> status=accepted|rejected`) — accepted
  `check:` lessons become enforced contract items next run (M3 wires the
  contracts.d write)

Close-out, in one atomic engine transaction:
1. `wf phase exit` (ship contract met)
2. `wf run close` — archives events to `.workflow/runs/<id>/`, compacts the
   live log (open followups, commit-origins, lessons stay live), clears the
   snapshot. Ordering is engine-owned — nothing to sequence by hand.
