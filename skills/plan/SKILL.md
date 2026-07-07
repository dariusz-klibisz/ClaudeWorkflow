---
name: plan
description: wf phase 4 (Plan) — decompose into verifiable tasks with a per-AC verification strategy, dependency check, and user approval. Invoked via /wf:dev when plan is the active phase.
---

# /wf:plan — Plan (interactive exit)

Contract first:

- `wf record task tid=T-1 subject="…" status=open --json '{"dod":["…"],"ac_links":["AC-1"]}'`
  — atomic tasks, each with a definition-of-done and AC links; for the diff
  family the FIRST task per AC is its failing test. Mirror each into the
  native task list (TaskCreate) **with the tid as subject prefix** —
  `"T-1: <same subject>"` — that prefix is how the gate links the native
  task to its wf record (never duplicate-record)
- Per AC: `wf record verification-strategy --json '{"ac":"AC-1","method":"…","command":"…"}'`
  — this becomes the Verify checklist (single-quoted --json avoids
  permission-prompt friction on embedded commands). The `command` also
  teaches Bash test capture this run's runner: record the REAL invocation
  you will run (e.g. `python3 -m unittest test_app.TestApp.test_x -v`), and
  Build/Verify runs of that runner get auto-captured in any language. Teams
  with custom wrappers (`./scripts/test.sh`) can add them to
  `.workflow/config.json` under `"runners": […]`.
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
  — pose the confirmation via AskUserQuestion (the hook anchors the answer)

`wf phase exit` when met — Build is auto-advance: from here the engine
expects execution, not conversation.
