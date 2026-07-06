---
name: plan
description: wf phase 4 (Plan) — decompose into verifiable tasks with a per-AC verification strategy, dependency check, and user approval. Invoked via /wf:dev when plan is the active phase.
---

# /wf:plan — Plan (interactive exit)

Contract first:

- `wf record task tid=T-1 subject="…" status=open --json '{"dod":["…"],"ac_links":["AC-1"]}'`
  — atomic tasks, each with a definition-of-done and AC links; for the diff
  family the FIRST task per AC is its failing test. Mirror each into the
  native task list (TaskCreate) — the TaskCompleted gate holds them to
  their DoD live
- Per AC: `wf record verification-strategy ac=AC-1 method="…" command="…"`
  — this becomes the Verify checklist
- `wf record scope-boundary --json '{"in_scope":[…],"out_of_scope":[…]}'`
- Deferred Frame ambiguities: disposition each (the gate lists them) — a
  new ambiguity record with `updates=<id>` and the final disposition, or a
  user-approved deferral (`wf approve deferral` first)
- diff family: `wf deps check` — verifies manifests + that every
  verification-strategy command's tool resolves, and records the deps
  verdict (missing blocks the gate; n/a is the honest escape). Run it AFTER
  recording the verification strategies
- spawn `@wf:critic` on the plan — or, when Design and Plan were presented
  and approved together, `wf contract waive plan.critic --reason "combined
  presentation covered at design"` (the waiver is recorded — E4)
- `wf approve plan --payload "<task list + verification strategy>"`

`wf phase exit` when met — Build is auto-advance: from here the engine
expects execution, not conversation.
