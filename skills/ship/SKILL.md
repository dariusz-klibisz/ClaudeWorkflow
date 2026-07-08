---
name: ship
description: wf phase 7 (Ship) — deliver the package, resolve the trace findings, lessons, archive, close. Invoked via /wf:dev when ship is the active phase.
---

# /wf:ship — Ship (interactive)

Contract first:
- `wf trace` — the engine computes the close-out findings (forced exits,
  open followups, unacked deviations, unresolved ambiguities) and prints
  phase coverage + the escape/waiver inventory. Resolve or disposition each
  finding: `wf record trace-finding updates=<id> status=resolved|dispositioned
  note="…"`; open followups become tasks now or are carried
  (`wf record followup updates=<id> status=next-run`). Re-run `wf trace`
  until no finding is open
- `@wf:auditor` verdict over the resolved trace (HIGH findings block close)
- diff/artifact: `wf trace --rtm --write` — generates the
  requirements-traceability matrix (`docs/requirements/RTM-<run>.md`):
  requirement → AC → verification → grounded evidence → verdict → tasks →
  loops. A derived view over the ledger (never hand-edit; re-run to
  regenerate) — this is the auditor/stakeholder-facing rendering of the
  chain the gates enforced
- diff: record the delivery package —
  `wf record artifact path=<PR-or-release-ref> role=delivery status=present`;
  intent deploy: `wf doc new delivery-manifest --slug <release>` (target,
  config diff, rollout, smoke, rollback), author + flip to present
- Lessons: `wf lessons status` first — the efficacy view (per accepted
  lesson: was its contract item waived? did its trigger recur in later
  runs?). A dodged or recurring lesson is a process gap to raise with the
  user, not a formality. Then `wf lessons suggest` (the engine proposes
  from the run's health signals) and add agent-spotted ones
  (`wf record lesson text="…" status=proposed [check="…"]`). The `check`
  field is a contract-item fragment (YAML or JSON: phase, predicate, params,
  remediation — the shipped predicate vocabulary). Ask the user, then
  disposition each: `wf lessons accept|reject <id>` — one command records
  the approval, flips the status, and regenerates the delivery channels:
  accepted `check:` lessons become enforced `lesson.*` contract items in
  `.workflow/contracts.d/lessons.yaml` (blocking from the NEXT run), prose
  lessons regenerate `.claude/rules/wf-lessons.md`. Commit both files.

Close-out, in one atomic engine transaction:
1. `wf phase exit` (ship contract met)
2. `wf run close` — archives events to `.workflow/runs/<id>/`, freezes the
   signals snapshot (`signals.json`), compacts the live log (ONLY open
   followups stay live; lessons and commit-origins archive with the run and
   their readers fold them back in), prunes the per-machine `local/`
   counters, clears the snapshot. Ordering is engine-owned — nothing to
   sequence by hand.

Optional but valued: `wf doc new retro --slug <run-id>` for runs worth a
written retrospective — author it before close (a stub blocks the
ship contract's artifact sweep).
