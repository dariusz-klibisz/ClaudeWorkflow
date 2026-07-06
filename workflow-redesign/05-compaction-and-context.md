# 05 — Compaction Survival and Context Discipline

Requirement: the workflow must survive compaction with no lost phases, steps,
documents, or obligations. Evidence base:
[01 §5](01-claude-native-capabilities.md#5-compaction-what-survives-what-hooks-fire).

## 1. The principle: context is a cache

Nothing the workflow depends on may live *only* in the conversation. The
conversation is a cache of `.workflow/` state; compaction invalidates the
cache; hooks repopulate it. An agent that "remembers wrong" is corrected by
the next injection; an agent that "forgot" is re-told before it can act.

## 2. The survival plan, layer by layer

| Layer | Mechanism | Survives compaction because… |
|---|---|---|
| Obligations (phases, contract items, tasks) | Engine state on disk; **`SessionStart` matcher `compact`** injects the regenerated status block | fires *after* every compaction; content comes from disk, not the summary |
| Per-turn anchoring | **`UserPromptSubmit`** injects a ≤10-line anchor (phase, top unmet items, active task) | fires on every prompt, unconditionally |
| Ground rules (recording contract, escapes, path roots) | **Project-root CLAUDE.md block** (≤40 lines) + unscoped `.claude/rules/` — incl. the engine-generated `wf-lessons.md` (accepted prose lessons, 03 §4.7) | both re-inject from disk after compaction (verified table) |
| Phase procedures | **Skills** — invoked bodies re-inject post-compaction (5k tokens/skill, 25k total, oldest-dropped, start-of-file kept) | verified; see §3 for the budget discipline |
| Long-lived knowledge | Auto memory MEMORY.md (first 200 lines/25KB) + per-agent `memory: project` dirs | auto memory re-injects; agent memory reloads per spawn |
| Reviewer context | `SubagentStart` injection per spawn; subagent transcripts unaffected by main-session compaction; subagents auto-compact independently | verified |
| Documents | On disk under `docs/` + artifact records; never "remembered" | n/a |
| Pending work items | Engine checklist + native tasks (the compaction summary explicitly preserves pending tasks — belt and suspenders) | verified |

What we deliberately do **not** rely on: the compaction summary's fidelity,
`PreCompact` injection (impossible — 01 §12.1), skill *descriptions* index
(not re-injected), path-scoped rules or nested CLAUDE.md (lost until
re-triggered — the workflow uses neither for anything load-bearing).

## 3. Skill design for the 5k/25k re-injection budget

- Each phase SKILL.md: **contract first** (the exit checklist and record
  commands in the first ~80 lines — truncation keeps the start), procedure
  detail after, deep guidance in supporting files loaded on demand
  (progressive disclosure). Target ≤4k tokens per skill body.
- 7 phase skills ≈ within the 25k total cap even if all were invoked in one
  session; in practice ≤3 are live at once. If the oldest is dropped, the
  Stop gate + injections still carry the *contract*; only the *how-to* needs
  re-invoking — and the injection names the skill to re-invoke
  ("resume with `/wf:build`").
- Entry skill `/wf:dev` stays tiny (read state → decide → invoke phase
  skill); escape skills (`/wf:park`, `/wf:force`) are
  `disable-model-invocation: true` (zero context until the user calls them,
  and Claude cannot invoke them on its own — an escape the model can't take
  unilaterally).

## 4. The CLAUDE.md block (per-project, ≤40 lines)

Written once by `/wf:init` (marker-delimited, engine can refresh it):
project facts (name, language, build/test commands) + the invariant ground
rules — state lives in `.workflow/`; work happens inside runs
(`/wf:dev`); record via `wf` commands, never by editing state files; docs go
under `docs/`; escapes are `/wf:park` and `/wf:force` (audited); after any
compaction or resume, trust the injected status block over memory. Everything
else (procedures, rosters, gates) deliberately lives in the plugin — nothing
version-sensitive is in the repo (fixes the B5–B7 stale-docs class). The one
other engine-maintained repo file is `.claude/rules/wf-lessons.md` (accepted
prose lessons, 03 §4.7): genuinely per-project knowledge, regenerated from
lesson records — never hand-edited, so it cannot go stale the B5 way.

## 5. Injection payloads (specs)

**`wf inject session`** (SessionStart, all matchers; ≤60 lines, hard-capped
well under the 10k-char limit):

```
[wf] run 20260705-a1b2 · diff/fix · phase: build (3/7) · started 2h ago
waiting-on: nothing — 4 contract items open
next actions:
  1. task T-3 "empty-file handling" in_progress — red test recorded, green missing
     → make it pass, then TaskUpdate completed
  2. roster: code-security-reviewer verdict missing → spawn @wf:code-security-reviewer
  3. record: wf record test "<cmd>" --ac AC-2 (AC-2 pass ungrounded)
  4. deviation D-1 awaiting user ack
open tasks: 2/6 · loops: 1 (AC-2, slip) · forces: 0
resume procedure: /wf:build   escapes: /wf:park /wf:force
```

**`wf inject turn`** (UserPromptSubmit; ≤10 lines): first 4 lines of the
above. **Formatting rule** (verified caveat): factual, declarative statements
— no imperative "SYSTEM:" framing, which can trip prompt-injection defenses.

**`wf inject agent <name>`** (SubagentStart): the reviewer's scope — run id,
what to review (diff ref/artifact paths/records), which corpus files to load
(06 §5), the verdict-block format, and its criticals policy.

## 6. Context budget

Steady-state overhead per session: CLAUDE.md block (~0.4k tokens) + session
injection (~0.6k) + per-turn anchor (~0.15k × turns) + one live phase skill
(≤4k) + auto-memory (≤1k). ≈6k tokens versus v0.36's ENTRY.md + workflow.md +
step file + perspectives + agent prose (~12–20k), while being *more*
persistent. Heavy content (corpora, templates, reviews) stays in subagent
contexts by design — the roster is subagent-only, and Explore handles
codebase mapping.

## 7. Resume and cross-session continuity

- `SessionStart` matchers `startup`/`resume`/`clear` get the same injection
  as `compact` — a fresh terminal, a `--resume`, and a `/clear` all
  re-anchor identically.
- Run identity is committed state (08 §5) with `wf run adopt` for
  re-attachment — a second machine or fresh clone resumes the same run
  instead of forcing a bypass cascade (fixes G2).
- `session_end`/idle: the Stop gate already guarantees no phase is silently
  abandoned mid-turn; `wf doctor` flags runs idle >N days for close-out
  (E2).
