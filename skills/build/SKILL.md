---
name: build
description: wf phase 5 (Build) — execute tasks test-first under the task gates; the review roster runs here, fixed to clean. Invoked via /wf:dev when build is the active phase.
---

# /wf:build — Build (auto-advance)

Per task, in order:
1. `TaskUpdate` the native task to `in_progress` and update the wf record:
   `wf record task updates=<task-event-id> status=in_progress` — the test
   capture binds runs to the single in-progress task automatically.
2. **Red first**: write the failing test, run it (the failing run is
   auto-captured and tagged with the task/AC). Then implement until green
   (also auto-captured). Genuinely testless task?
   `wf contract waive <tid> --reason "…"` — otherwise TaskCompleted will
   refuse the checkbox.
3. Complete the native task (TaskUpdate → completed). The gate verifies the
   red→green pair for this task and marks the wf record done.

Rules while building:
- Departures from the approved design are recorded, then user-acked:
  `wf record deviation text="…" status=pending`, ack via
  `wf approve deviation` + `wf record deviation updates=<id> status=acked`.
- Out-of-scope discoveries: `wf record followup text="…" status=open` —
  never silent scope expansion.
- Commit messages carry `[run:<id>]` — commits are auto-captured as
  `commit-origin` records; an untagged commit gets a visible reminder.

Roster (diff family) — spawn on the accumulated diff, fix findings to clean;
verdicts are captured automatically, failing auto verdicts are sticky:
`@wf:code-quality-reviewer` · `@wf:code-security-reviewer` ·
`@wf:code-testing-reviewer` · `@wf:design-conformance-reviewer` ·
`@wf:adversary` (red-team mode) · `@wf:ux-reviewer` (ux projects).
Rosters parallelize — spawn them together, they run in the background.

artifact/assessment families: author the deliverable instead;
`wf record artifact path=… status=present role=deliverable-report` (reports)
and keep tasks gated the same way.

`wf phase exit` when the contract is met (the Stop gate will keep pointing
at what is missing).
