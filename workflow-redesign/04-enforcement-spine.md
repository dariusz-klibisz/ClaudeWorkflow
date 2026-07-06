# 04 — The Enforcement Spine

How "Claude cannot skip steps" becomes a property of the tooling. Every
guarantee maps to a *blocking* mechanism from the verified matrix
([01 §2](01-claude-native-capabilities.md#2-the-blocking-matrix)); prose is
never load-bearing. The deterministic core is command hooks + the `wf` engine;
prompt/agent hooks are an additive semantic layer (01 §12.8).

## Table of contents
1. [The run state machine](#1-the-run-state-machine)
2. [Gate 1 — the Stop gate](#2-gate-1--the-stop-gate)
3. [Gate 2 — task gates](#3-gate-2--task-gates)
4. [Gate 3 — the reviewer verdict gate](#4-gate-3--the-reviewer-verdict-gate)
5. [Gate 4 — tool gates (PreToolUse)](#5-gate-4--tool-gates-pretooluse)
6. [Context injection (the always-visible checklist)](#6-context-injection-the-always-visible-checklist)
7. [Escape hatches and escalation](#7-escape-hatches-and-escalation)
8. [Honest bounds](#8-honest-bounds)
9. [Guarantee → mechanism matrix](#9-guarantee--mechanism-matrix)
10. [Legacy issue cross-reference (A1–G6)](#10-legacy-issue-cross-reference-a1g6)

---

## 1. The run state machine

The engine owns a small, explicit state machine per run
(`.workflow/state/run.json` + append-only event log, 08):

```
(no run) ──wf run start──▶ frame ▶ context ▶ [design] ▶ plan ▶ build ▶ verify ▶ ship ──wf run close──▶ (closed)
                              ▲______________loop (verify→build/design/plan)______________|
   parked ◀──wf park── (any)         branched ──▶ new run (parent linked)
```

- `wf phase exit` evaluates the current phase contract (03 §4) against
  **records only** and either advances or returns the finding list. Exit codes
  0/2/3 keep the pass/gaps/broken distinction. The contract itself is data —
  the declarative spec + `contracts.d/` additions (03 §4.0) — evaluated by
  the engine's closed predicate set; the state machine's phase list is read
  from the same spec.
- Phase *entry* is implicit in exit of the predecessor; re-entry happens only
  via loop/resume events. There is no way to be "in" two phases.
- The engine is the **only writer** of phase transitions. Skills instruct;
  hooks enforce; the engine records.

## 2. Gate 1 — the Stop gate

**The anti-"declared done too early" mechanism.** A `Stop` command hook
(plugin `hooks/hooks.json`) invokes `wf gate stop`, which inspects state and
decides:

| State | Decision |
|---|---|
| No active run | allow |
| Active phase **waiting on the user** (an approval gate reached, a blocking ambiguity posed, AskUserQuestion pending) — engine knows because the *next unmet contract item is an approval/user record* | allow (turn ends so the user can respond); `additionalContext` prints the waiting-on line |
| Active phase with unmet contract items the agent can progress **without** user input (unrecorded verdicts, ungrounded ACs, unclosed tasks, missing records) | **block** — `decision:"block"`, `reason` = the top ≤5 unmet items with the exact `wf` commands/skills to produce them |
| `stop_hook_active: true` and the unmet set is unchanged from the previous block | allow after 3 identical blocks with a system message recommending `/wf:park` (self-imposed cap under the platform's 8; prevents burn-loops on a genuinely stuck item) |
| Run parked/closed | allow |

Design notes:
- Deterministic (command hook). An optional **prompt-hook layer** can be added
  for softer conditions ("is the summary faithful to the diff?") — additive,
  never the core (agent hooks are experimental, 01 §4).
- `background_tasks` from the Stop input are consulted: a turn ending while a
  gating reviewer subagent is still running is allowed (the work is in
  flight), with an anchor note.
- This is the exact `/goal` pattern, generalized and made state-driven —
  verified supported (01 §12.7).

## 3. Gate 2 — task gates

**The anti-"checked it off anyway" mechanism.**

- `TaskCreated` → `wf gate task-create`: enforces task shape (DoD present, AC
  link or explicit `pure-technical`, created under an active Plan/Build
  phase). Exit 2 rolls the creation back with the correction fed to the model.
- `TaskCompleted` → `wf gate task-complete`: looks up the task's DoD in engine
  state and verifies its records exist — e.g. for a diff task: an AC-tagged
  green test-run newer than the red one; for a doc task: the artifact record
  present-not-stub. Exit 2 = "not complete; missing: …" fed straight back.
- Native tasks are a mirror of engine checklist state (03 §7); the engine
  reconciles at SessionStart so a lost native list never loses obligations.

## 4. Gate 3 — the reviewer verdict gate

**The anti-"reviewer finished without a verdict / verdict faked" mechanism.**
Anchored on `SubagentStop` (background-default reality, 01 §12.3):

- Plugin `hooks/hooks.json` registers `SubagentStop` matchers for the gating
  roster's scoped names (`^wf:design-reviewer$`, …) → `wf gate verdict`:
  1. Parse the fenced ```` ```verdict ```` block from
     `last_assistant_message` (fallback: tail of `agent_transcript_path`).
  2. Missing/malformed → **block the subagent finishing** (`decision:"block"`,
     reason = the exact block format) — the reviewer itself is forced to emit
     a parseable verdict; `unparsed` can no longer reach the ledger as a
     terminal state (kills A2/A3 at the source). After 2 blocks, record
     `unparsed` (fails the phase gate) and allow — no wedge.
  3. Valid → **auto-record** the verdict (`auto: true`) with criticals/majors;
     `clean`+criticals>0 auto-downgrades. Vocabulary: single constant
     (07 §4): `clean | changes-required | safe | risky | unsafe | n/a`.
- `SubagentStart` matchers inject the reviewer's inputs (the diff pointer,
  the design/records it must check, its corpus routing) as
  `additionalContext` — a reviewer can no longer run against the wrong or
  missing inputs (kills C5's wrong-phase attribution: the verdict event
  carries the run/phase the engine injected, not "last opened mark").
- **Sticky auto-evidence** (tightens legacy F2): a manual
  `wf record verdict` cannot supersede an `auto:true` failing verdict; only a
  new auto-captured run of the same reviewer, or an explicit
  `wf disposition` record (surfaced in reports), can.

## 5. Gate 4 — tool gates (PreToolUse)

Minimal by design — the heavy lifting moved to gates 1–3:

| Matcher / if | Gate | Behavior |
|---|---|---|
| `Skill` / `Skill(wf:*)` | phase-sequence gate | `wf gate skill`: deny invoking a phase skill that is not the active phase (or a legal loop-back target). Replaces v0.36's `pre_read` ordinal police — procedures are skills now, and skill invocation is a tool call we can deny with a reason |
| `Bash` | catastrophic net | The small always-on blocklist (rm -rf /, force-push default branch, curl\|sh, /etc writes) — `permissionDecision:"deny"`. Duplicated as `permissions.deny` rules where expressible (permission system is the harder boundary; the hook covers pattern gaps). No env-var escape |
| `Edit\|Write` | stray-edit guard | Deny **project-file** edits when no run/phase/task is active ("start or resume a run first; docs edits under Ship are exempt") — replaces the intent ledger (03 §8); `.workflow/` and `docs/` bookkeeping exempt by path anchor (not basename — fixes C7) |
| `PostToolUse(Bash)` | test capture | `wf capture test`: recognized runners → grounded `test-run` records (explicit exit code only; filter-pipe → ungrounded), **skips any command whose head resolves into plugin `bin/`/hooks** and matches runners only in the command head (fixes G1) |
| `PostToolUse(Edit\|Write)` | edit ledger | Append edit→active-task binding records (never blocks) |

## 6. Context injection (the always-visible checklist)

The anti-"didn't know / forgot" mechanism (fixes G3/G6 structurally):

- **`SessionStart` (all matchers: startup/resume/clear/compact)** →
  `wf inject session`: run id, family+intent, active phase, the unmet
  contract items, waiting-on state, open tasks, last 3 events — regenerated
  from disk, ≤60 lines. After compaction this fires with `source:"compact"`,
  so the summary is never the source of truth.
- **`UserPromptSubmit`** → `wf inject turn`: a ≤10-line anchor (phase, top
  unmet items, active task). Cheap, every turn, and immune to
  conversation drift. (30s timeout noted; the engine answers in
  milliseconds, 07.)
- **`SubagentStart`** → reviewer inputs (§4).
- Phase procedures are **skills** whose bodies re-inject after compaction
  (5k/25k caps — contract at the top of each SKILL.md, details in supporting
  files; 05 §3).

## 7. Escape hatches and escalation

| Hatch | Effect | Audit |
|---|---|---|
| `/wf:park` (`wf park --reason`) | Clears all gates for the run; resumable | Recorded, reported |
| `/wf:force` (`wf force-exit --reason`) | Bypass one phase gate | Recorded; **escalates**: 2nd force in a run demands a structural-cause field; 3rd auto-parks with a repair checklist (fixes G4) |
| `WF_ENFORCE=0` | Downgrades Stop/skill/edit gates to warnings (task + verdict gates keep functioning — they protect data integrity, not sequencing). **Provenance-guarded**: honored only in hook-invoked contexts (the env Claude Code spawns hooks with — user-controlled), *ignored* when `wf` runs inside the agent's Bash tool (hook stdin JSON present ⇒ hook context; absent ⇒ agent context). An agent exporting `WF_ENFORCE=0` in its own shell changes nothing — the escape belongs to the human | Loud warning each firing; recorded `escape` event; session counter in report |
| — | Catastrophic Bash net has **no** hatch | n/a |
| Broken engine | Every hook wrapper: if `wf` is missing/crashes, gates that *sequence* fail open with a loud `systemMessage`; gates that *record* fail closed (no fabricated evidence). `wf doctor` is the repair path; park/force never depend on gate evaluation (no-wedge invariant) |

## 8. Honest bounds

Unchanged in kind, sharpened in surfacing:
1. **Approvals are self-attested** — `wf approve` records who/what/payload;
   no hook event proves a human typed it. Reported per run. (Optional
   hardening: an `AskUserQuestion`-based approval skill whose transcript
   position is recorded — still not proof, noted as such.)
2. Manual records are `auto:false` and reported; auto-captured facts
   (verdicts, test-runs) dominate via stickiness (§4).
3. The Stop gate cannot outlast 8 platform blocks; park is always the honest
   terminal state — enforcement makes skipping *loud and recorded*, not
   impossible against a determined operator. That is the design goal: the
   agent cannot skip silently; the human can always decide.

## 9. Guarantee → mechanism matrix

| Guarantee | Mechanism (event) | Failure mode if mechanism down |
|---|---|---|
| Can't end turn with progressable unmet contract | Stop gate (`Stop`, block) | fail-open + loud (§7) |
| Can't mark a task done without its DoD evidence | `TaskCompleted` exit 2 | task stays open (fail-closed) |
| Can't create malformed tasks | `TaskCreated` exit 2 | rolled back |
| Reviewer can't finish without machine verdict | `SubagentStop` block + auto-record | recorded `unparsed` → phase gate blocks |
| Reviewer can't run on wrong inputs | `SubagentStart` injection | verdict recorded without injected scope → flagged |
| Can't exit a phase without its records | `wf phase exit` gate (engine, exit 0/2/3) | fail-closed (3) |
| Can't jump to a later phase's procedure | PreToolUse deny on `Skill(wf:*)` | fail-open + loud |
| Can't edit project files outside a run/task | PreToolUse deny on Edit/Write | fail-open + loud |
| Test results can't be faked/fabricated | Auto-capture only from real runner exits; null/filtered = ungrounded; manual can't override auto-fail | manual records flagged `auto:false` |
| ACs can't pass ungrounded | Verify contract: pass ⇒ linked AC-tagged green run | fail-closed |
| Loops can't run forever / mask design defects | Engine caps (10/run, 2 slip/AC) in the loop command itself | command refuses |
| Nothing is forgotten across compaction | SessionStart(compact)+UserPromptSubmit injection; state on disk | next injection repairs |
| Nothing is forgotten across sessions | SessionStart(startup/resume) injection + native-task resync | idem |
| Run can't close with leaks | Ship contract: engine trace findings resolved; tasks closed/carried | fail-closed |
| Bypasses are visible | park/force/WF_ENFORCE all recorded + reported + escalating | — |

## 10. Legacy issue cross-reference (A1–G6)

How each catalog class from `docs/workflow-reference/07-issues-catalog.md`
becomes structurally impossible (spot checks on the named items):

- **A1 (phantom agent)**: one roster declaration in the workflow spec
  (03 §4.0) generates the agent files, the SubagentStop matchers, and the
  contract requirements at build time (07 §4); a name in a contract without
  an agent file fails spec validation and the plugin's CI
  (`claude plugin validate` + engine selftest).
- **A2/A3 (`n/a` unparseable)**: verdict vocabulary is one constant; the
  SubagentStop gate *forces* emission and parses with the same constant;
  `n/a` is first-class.
- **A5/A6 (archive ordering, task-ref collisions)**: `wf run close` is one
  atomic engine operation; IDs are monotonic engine-issued, never
  line-counts.
- **B1–B4 (route-scoping mismatches)**: families are a single field read by
  both the procedure selection and the gate — there is no second set to
  disagree (03 §2).
- **B5–B7 (stale instruction docs)**: CLAUDE.md ground-rules block and skill
  text ship in the plugin, versioned with the hooks they describe; no
  per-repo copies to go stale.
- **C1/C2/C3/C12 (substring heuristics, header parsing)**: no prose is
  parsed, ever (03 §5); the run header is engine-written.
- **C4 (run-global red→green)**: per-AC/task pairing (03 §4.5).
- **C5/C6 (phase attribution, unnormalized labels)**: phases are engine
  state, not string labels replayed from logs; evidence events get the
  engine's current phase atomically.
- **C10 (empty Context)**: min-content contract item (03 §4.2).
- **C13 (python3 everywhere)**: engine is a static binary (07).
- **D1/D2 (4× rosters, hand-synced versions)**: single-sourced in the
  workflow spec, consumers generated/validated (03 §4.0, 07 §4).
- **D3 (settings drift)**: enforcement wiring lives in plugin hooks.json —
  updates with the plugin.
- **E2 (abandoned runs)**: `wf doctor` offers stale-run close-out; SessionStart
  flags a run idle >N days.
- **F2 (manual overrides auto)**: sticky auto-evidence (§4).
- **G1 (false test capture)**: capture excludes hook self-calls, matches
  command head only (§5).
- **G2 (run-ID mismatch → force cascade)**: run identity lives in committed
  state (08 §5) with an explicit `wf run adopt` re-attach command;
  session_start-style minting divergence can't orphan a run.
- **G3/G6 (checklist not open / lost post-compaction)**: §6 injection.
- **G4 (force cascade)**: escalation (§7).
- **G5 (retro gate empty)**: Ship contract is the strictest, engine-computed
  (03 §4.7).
