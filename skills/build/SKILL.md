---
name: build
description: wf phase 5 (Build) — execute tasks test-first under the task gates; the review roster runs here, fixed to clean. Invoked via /wf:dev when build is the active phase.
---

# /wf:build — Build (auto-advance)

Per task, in order:
1. `TaskUpdate` the native task (subject `"T-<n>: …"` — the tid prefix links
   it to the wf record) to `in_progress`, and update the wf record:
   `wf record task updates=<task-event-id> status=in_progress` — the test
   capture binds runs to the single in-progress task automatically. Exactly
   ONE task in_progress at a time: the implementer briefing targets it.
2. Spawn `@wf:implementer` for the task (diff family). The SubagentStart
   briefing injects the task's tid, DoD, AC texts, the verification
   commands, the approved design selections, and the out-of-scope boundary
   — you route work, it executes. **Red first**: it writes the failing test
   that encodes the AC, runs it (auto-captured, tagged with the task/AC),
   then implements until green (also auto-captured). A genuinely testless
   DIFF task needs `wf contract waive <tid> --reason "…"` — otherwise
   TaskCompleted will refuse the checkbox. **artifact/assessment tasks need
   NO test evidence and NO waiver** — the test-first gate is diff-only;
   waiving doc tasks is pure ceremony (the arch-design run waived 7 for
   nothing); author those directly instead of spawning the implementer.
   **After the FIRST test run, confirm it was auto-captured** (`wf trace`
   or `wf status` shows the test-run with auto/hook provenance). If it
   wasn't: the runner isn't recognized — run it exactly as recorded in the
   verification strategy, or add the wrapper to config `"runners"`, and
   record the missed run manually (`wf record test-run …`) so evidence
   stays honest. Don't discover this at Verify.
3. When the implementer returns, act on what it surfaced (deviations to
   ack, followups recorded), then complete the native task (TaskUpdate →
   completed). The gate verifies the red→green pair for this task and marks
   the wf record done.

Rules while building (they bind the implementer too — it records, you
present to the user):
- Departures from the approved design are recorded, then user-acked:
  `wf record deviation text="…" status=pending`, ack via
  `wf approve deviation` + `wf record deviation updates=<id> status=acked`.
- Out-of-scope discoveries: `wf record followup text="…" status=open` —
  never silent scope expansion. Paths the scope-boundary declares
  out-of-scope are hard ground: edits there surface as high trace findings
  at Ship.
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
and keep tasks gated the same way. Assessments: record findings as you make
them (`wf record finding fid=F-1 severity=… text=…`) — Verify requires the
report to name every recorded fid verbatim.

`wf phase exit` when the contract is met (the Stop gate will keep pointing
at what is missing).
