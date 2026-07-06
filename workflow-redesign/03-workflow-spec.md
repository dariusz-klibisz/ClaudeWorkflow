# 03 — Workflow Specification (Lifecycle Re-Derivation)

The lifecycle rebuilt from first principles. Constraint from the user: "phases
were ok, but think deeper whether all routes are needed and whether all steps
were correct — redesign if there is a better way." Method: every phase and step
must (a) name its consumer, (b) name its enforcement mechanism, or be cut;
everything the engine can compute is removed from the agent's job.

## Table of contents
1. [What changed and why](#1-what-changed-and-why)
2. [Route families (replacing 12 routes)](#2-route-families-replacing-12-routes)
3. [The seven phases](#3-the-seven-phases)
4. [Per-phase contracts](#4-per-phase-contracts)
5. [Structured records (abolishing magic prose)](#5-structured-records-abolishing-magic-prose)
6. [Loops, branches, parking](#6-loops-branches-parking)
7. [The no-leak funnel on native tasks](#7-the-no-leak-funnel-on-native-tasks)
8. [What was cut, with rationale](#8-what-was-cut-with-rationale)

---

## 1. What changed and why

Three structural findings from the as-is analysis drive the redesign:

1. **Route explosion caused the bug class.** 12 classifications × 10 phases ×
   route-only/skip markers produced the entire B-class of scoping defects
   (prose demanding gates the gate didn't run, gates blocking phases whose
   procedure was filtered away). The *gate* only ever distinguished 5 route
   sets (`_BUILD_DELIVER/_BUILD/_DIFF/_ASSESSMENT/_NO_DESIGN` in
   `phase_gate.py`), and those collapse to exactly the three family shapes
   plus two orthogonal switches this spec keeps — deliver-packaging (now the
   `deploy` intent tag, §2) and design participation (now the family
   participation table, §2). The other distinctions were prose flavor.
2. **Two phases were mostly machine work.** Trace (route-completeness, no-leak
   reconciliation, RTM cross-check, ledger verification) is arithmetic over
   state the engine owns — making an LLM *perform* it was both wasteful and
   unreliable. Similarly half of Intake (run minting, task triage listing,
   folder creation) is bookkeeping.
3. **Interactive elicitation was split across two phases** (Intake restates +
   confirms classification; Clarify restates again + confirms scope), causing
   duplicated user confirmations and the weakest gate in the system
   (issue C1: "Classification:" substring).

Result: **7 phases, 3 route families, and a hard rule: the agent does judgment
work; the engine does bookkeeping and verification.**

## 2. Route families (replacing 12 routes)

A run is classified by its **deliverable shape** — the only property the
contracts ever needed:

| Family | Deliverable | Legacy labels absorbed |
|---|---|---|
| **`diff`** | A code change (tests, source, config, deploy artifacts) | CODE_NEW, CODE_FIX, CODE_REFACTOR, TEST, DEPLOY |
| **`artifact`** | Authored durable documents (design docs, ADRs, threat models, end-user docs) | ARCH_DESIGN, DOC_CREATE, DOC_UPDATE |
| **`assessment`** | A findings report about existing work (no changes made) | CODE_REVIEW, ARCH_REVIEW, INVESTIGATE, RESEARCH |

The legacy label is kept as a free-form **`intent` tag** on the run record
(`fix`, `refactor`, `deploy`, `investigate`, …). Tags tailor *content within a
step* (e.g. `fix` and `investigate` trigger origin-discovery; `refactor`
suppresses new requirements; `deploy` adds release packaging) but **never
change which phases run or which gates apply** — that is family-only. This is
the structural fix for B1–B4: there is exactly one rendered procedure per
(family, phase), and the gate keys on the same family field of the same record.

Phase participation:

| Phase | diff | artifact | assessment |
|---|---|---|---|
| 1 Frame | ✅ | ✅ | ✅ |
| 2 Context | ✅ | ✅ | ✅ (map-only depth) |
| 3 Design | ✅ | ✅ | ⛔ (assessments evaluate designs; they don't make them) |
| 4 Plan | ✅ | ✅ (outline) | ✅ (report outline) |
| 5 Build | ✅ | ✅ (author artifacts) | ✅ (author report) |
| 6 Verify | ✅ | ✅ (artifact-shaped gates) | ✅ (report-shaped gates) |
| 7 Ship | ✅ | ✅ | ✅ |

Within-family skip (e.g. a trivial diff that doesn't need staged design) is an
explicit, engine-recorded decision (`wf phase waive design --reason "…"`),
allowed only where the family contract marks the phase *waivable*, surfaced in
`wf report`, and re-checked at Ship (§4.7). Never silence.

## 3. The seven phases

| # | Phase | Was (v0.36) | Mode | One-line purpose |
|---|---|---|---|---|
| 1 | **Frame** | Intake + Clarify | interactive | Understand the task with the user: classify (family+intent), risk-screen, restate, elicit requirements/ACs through the lenses |
| 2 | **Context** | Context | interactive exit | Map the code/reality; validate feasibility; baseline + approve scope/requirements |
| 3 | **Design** | Design | interactive exit | Staged option evaluation → selected design, reviewed fixed-to-clean, critic-checked, user-approved |
| 4 | **Plan** | Plan | interactive exit | Decompose into verifiable tasks (mirrored to native task list), verification strategy, dependency check, approval |
| 5 | **Build** | Implement | auto-advance | Execute tasks test-first; per-task verification; the review roster runs here, fixed-to-clean |
| 6 | **Verify** | Verify | interactive exit | Per-AC grounded verification; confirmation gates; loop records on failure |
| 7 | **Ship** | Trace + Deliver + Retro | interactive | Deliver the package; the engine emits the trace report (agent resolves findings); lessons; archive; close |

Merging rationale:
- **Frame = Intake+Clarify**: one continuous elicitation conversation, one
  confirmation flow. What Intake did beyond that was bookkeeping (run minting,
  folders, task triage *listing*) — now engine work in `wf run start`. The
  user confirms (family, intent, scope) once, as a recorded approval, killing
  C1 (substring "confirmation").
- **Ship = Trace+Deliver+Retro**: Trace's checks are computed by the engine at
  Ship entry (`wf trace` → findings list). The agent's remaining Ship work is
  genuine: package/PR authoring, resolving trace findings, lessons
  (agent-authored, user-approved), archive trigger. Retro's PDCA prose is cut
  (§8); its measurable outputs (signals, lessons, close-out) stay.
- **Not merged**: Context vs Frame (elicitation vs codebase-grounding — the
  reclassify-on-map-contact checkpoint needs a boundary between them); Design
  vs Plan (option selection vs decomposition have different reviewers and
  different approval semantics; the combined-for-SIMPLE path remains as one
  presentation + one approval, engine-recorded); Build vs Verify (author vs
  independent confirmation — collapsing them is how verification dies).

## 4. Per-phase contracts

Contract = the records that must exist (engine-verified) for `wf phase exit`
to pass. All record kinds are defined in §5; approvals are discrete events;
verdicts come from the SubagentStop capture path (04 §4). "◇" = waivable with
recorded reason.

### 4.0 The workflow is data (extensibility by design)

The v0.36 postmortem shows *where* extensibility dies: phases, routes, gates,
and rosters as string literals in ≥9 places, so adding one verification step
was a 5–7-file change and adding a phase a ~15-site change. Single-sourcing
into compiled constants (an earlier draft of this design) fixes the drift but
still makes every contract change an engine release. Instead, the workflow
definition is **data**, and the engine is a generic interpreter:

- **One spec** — `workflow/workflow.yaml` in the plugin (JSON-Schema-validated
  at load and in CI) — declares:
  - **phases**: id, order, mode (interactive/auto-advance), waivability,
    owning skill;
  - **families** and the phase-participation table (§2);
  - **record kinds** with their payload schemas (08 §3) — the vocabulary of
    inputs/outputs;
  - **the gating roster** (agents, their phases, verdict policy) — the hooks
    matchers, agent contract requirements, and 06's tables are generated
    views of this one list;
  - **contract items** per (phase, family): each item = `{id, families,
    predicate, params, waivable, remediation}`.
- **Closed predicate vocabulary.** Contract items compose a small, fixed set
  of engine-implemented predicates (final set, fixed by transcribing all of
  §4 before the interpreter was built — which surfaced the last two):
  `record-exists(kind, filter, ≥min)` · `linked-record(kind, link, filter)`
  (inside per-each) · `verdict-in(agent, statuses, scope,
  risky-with-dispositions)` (with the sticky-auto-evidence rule) ·
  `approval(gate)` · `no-open(kind, field, open-values)` (over
  update-folded record state) · `per-each(kind, each, item)` (e.g. "per AC:
  a linked green test-run") · `any-of(items)` · `red-green(link)` (a failing
  then a later passing grounded test-run, both link-tagged) — plus nothing
  else. No expression language, no scripting: a predicate the vocabulary
  can't express is an engine change *on purpose* (that boundary is what
  keeps gates deterministic and testable).
- **What each change costs now**: a new input/output = a record kind + schema
  (spec-only); a new verification step = one contract item (spec-only); a new
  action = a `wf record`/capture producer for an existing or new kind; a new
  phase = one spec entry + one SKILL.md (the state machine, Stop-gate lists,
  injections, and skill gating all derive from the spec). None of these touch
  engine code unless a genuinely new predicate is needed.
- **Project-local extensions**: `.workflow/contracts.d/*.yaml` (committed) may
  *add* record kinds and contract items (never remove or weaken shipped ones;
  the engine enforces add-only). This is also exactly how accepted **lessons
  with `check:`** become enforced next run (§4.7) — a lesson-check *is* a
  contract item in `contracts.d/`, one representation, one evaluator. And it
  is the future home of regulated packs: a companion plugin = a contract pack
  + reviewer agents extending the same spec (09 §4.1).
- **Meta-validation replaces the meta-test**: `wf doctor`/CI verify that every
  contract item references declared record kinds/agents/phases, every gating
  agent in the spec has an agent file and a hooks matcher, and every skill
  named by a phase exists — the v0.36 `gate_contract_map`/parity-test genre,
  now checking one artifact instead of a dozen.

### 4.1 Frame
- `classification` record (family + intent + restated task) **approved by
  user** (approval event `frame`).
- `risk` record — engine `wf risk scan` output (deterministic signal grep) +
  agent-added signals; each signal bound to a lens.
- Per selected lens: ≥1 `ambiguity` record (tagged, dispositioned:
  resolved/logged/deferred) or an explicit `none` with reason.
- diff/artifact: `requirement` records (SWR-style, with ACs; each AC
  `verifiable: true|reason`) + a `completeness` record (negative-space walk:
  error/empty/max/concurrent/unhappy items, each dispositioned).
- Gating agent runs (family diff/artifact): `abuse-case-analyst` verdict;
  security + adversarial lens reviewer verdicts.
- intent `fix`/`investigate`: `origin` record (`wf origin discover` — commit
  attribution, best-effort content, **required-present** at Verify for
  `investigate`).

### 4.2 Context
- `context-map` record: files/modules examined + sufficiency note (min-content
  enforced: ≥N entries or explicit tiny-scope reason — fixes C10).
- `assumption` records; `high-risk` flagged ones must appear in the approval
  payload.
- `reclassify` checkpoint record: `confirmed` or family/intent change (→
  branch, §6).
- diff/artifact: requirements baselined (`active`/`dropped`/`revised` status
  transitions recorded) + **approval event `scope`** (payload = requirement
  diff + high-risk assumptions). assessment: approval event `scope` with the
  map summary (lighter payload, same discrete event — fixes the B1 class:
  same gate, every family).
- `researcher` verdict when external research ran, or `n/a` record.

### 4.3 Design (diff, artifact; waivable ◇ for trivial diffs)
- `option-set` records per stage (system / software / UX-when-applicable):
  2–4 genuine candidates, selection + rejection reasons; loop re-entries must
  reference prior `rejected` options (engine cross-checks IDs — never
  re-proposed).
- Reviewer verdicts fixed-to-clean per stage: `design-reviewer` (merged
  system+software reviewer, 06 §2), `ux-design-reviewer` (when
  `ux: true` in project config and change is UI-bearing; else recorded `n/a`).
- `threat-model` + `attack-tree` records when risk signals/trust boundaries
  present (analyst verdict; high-feasibility paths mitigated or
  ADR-accepted).
- `critic` verdict (pass/`risky`-with-dispositions/fail); dispositions are
  records, criticals unwaivable.
- ADR artifact record when architectural (alternatives = the option-set IDs —
  generated skeleton by engine).
- **Approval event `design`** (payload: selected options + risks + testability
  sketch).

### 4.4 Plan
- `task` records — atomic, each with definition-of-done + AC links, mirrored
  to native `TaskCreate` (§7). diff: first task per AC is its failing test.
- `verification-strategy` record per AC (method + tool/command) — becomes the
  Verify checklist.
- `scope-boundary` record; leftover Frame `deferred` ambiguities dispositioned
  (engine lists them; unresolved ones block).
- diff: `deps` record from `wf deps check` (manifest-driven presence check;
  `missing` blocks; `n/a` honest escape). artifact/assessment: auto-`n/a`
  (fixes B2 — the gate and the procedure share one family key).
- `critic` verdict (waived if the design-phase critic covered a combined
  presentation — the waiver is itself a record, fixing E4's invisibility).
- **Approval event `plan`**.

### 4.5 Build (auto-advance)
- Per task: red→green `test-run` pair **tagged with the task/AC id** (engine
  captures via PostToolUse(Bash); manual grounding via `wf record test`);
  `test-first: n/a` record with reason for genuinely testless tasks (fixes C4
  — red→green is per-AC, not run-global).
- Task completion is gated live by `TaskCompleted` (04 §3): a task closes only
  when its DoD records exist.
- **Roster verdicts fixed-to-clean on the diff/artifact** (family-scoped
  roster, 06 §2), captured at SubagentStop.
- `deviation` records for any departure from the approved design (user-acked).
- Out-of-scope discoveries → `followup` records (§7), never scope expansion.
- Commits tagged `[run:<id>]` + `commit-origin` records (durable attribution).

### 4.6 Verify
- Per AC: a `verdict` record (`pass|fail|deferred`) where `pass` **requires a
  linked green `test-run` for that AC** (diff family; artifact/assessment use
  linked artifact-checks instead). `deferred` requires a user approval event.
- Confirmation verdicts: `adversary` + `design-conformance` (diff/artifact;
  consistency with Build-phase verdicts enforced), `ux-reviewer` (when
  applicable), quality floor metrics if configured **and** measured.
- Security baseline records (diff): secret-scan, SCA — engine-captured
  runner results (`filtered`/`exit:null` = ungrounded, never pass).
- intent `deploy` (closes legacy deferred issue 12): deployment-shaped
  verification items in place of "unit suite green" alone — a `smoke-run`
  record (post-deploy/staging smoke or health-check command, engine-captured
  like any test-run), a `rollback-readiness` record (procedure + trigger
  condition, or explicit `n/a` with reason), and config/target diff captured
  into the delivery manifest at Ship.
- assessment: deliverable-report artifact record present-not-stub;
  `investigate` intent: origin attribution present.
- On any `fail`: a `loop` record (§6) — the engine refuses `exit` while an
  undispositioned `fail` exists.

### 4.7 Ship
- **Engine-generated trace report** (`wf trace`): route-completeness (phases
  entered/exited vs family contract incl. waivers), unresolved records,
  open followups, unconsumed approvals, verdict coverage, forced/parked
  events. Agent resolves each finding or dispositions it; `auditor` agent
  verdict over the resolved report (HIGH findings block).
- Delivery records: family-appropriate package (`pr`, `release`, `report`)
  with a `delivery-manifest` record for diff+deploy (artifact events
  present-not-stub).
- `lesson` proposals (engine `wf lessons suggest` + agent-spotted), user
  accept/reject (approval events); accepted lessons with `check:` become
  engine-enforced contract items next run — written as ordinary contract
  items into `.workflow/contracts.d/lessons.yaml` (§4.0; one evaluator, no
  special case). Accepted **prose-guidance lessons** (no `check:`) get a
  defined delivery channel: the engine regenerates
  `.claude/rules/wf-lessons.md` (committed, marker-delimited) — unscoped
  rules load at launch and re-inject after compaction (01 §10), replacing
  v0.36's `learn --apply` injection into now-read-only plugin skills.
- Signals snapshot (`wf report --run`) written to the run archive.
- `wf run close` — archives, compacts, clears — in one atomic engine command
  (fixes A5/A6/G5: ordering is no longer agent-sequenced prose).

## 5. Structured records (abolishing magic prose)

The single most important change. In v0.36 the gate scanned `decisions.md` for
~20 magic substrings (`**Risk signals**:`, `Suite delta:`, `passish`…), which
produced both false blocks (v0.30's whole theme) and false passes (C1, C2).

In v2:
- Every gate-relevant fact is created by an **engine command** (`wf record
  <kind> …`, `wf approve <gate>`, or automatic capture) that validates its
  vocabulary at write time and appends a typed JSONL event (schema in 08 §3).
- **The gate reads only records.** It never opens any Markdown.
- `decisions.md` is replaced by an engine-*generated* narrative
  (`wf log --md`, regenerated on demand from records) — humans read it; no
  machine ever parses it. There is nothing left for an agent to phrase wrongly
  and nothing for a gate to scan fragile tokens out of.
- Free-text rationale still exists — as a `note` field *on records*, not as
  structure-bearing prose.

## 6. Loops, branches, parking

Carried forward (the proven part of v0.36), with engine mechanization:

- **Loop** (in-run): `wf loop --ac AC-3 --cause slip|design|plan --evidence
  "…"` writes the loop record (discriminating evidence + observed-vs-expected
  + suite-delta captured from the last test-runs automatically) and re-opens
  the target phase. Caps engine-enforced: 10 cycles/run, 2 slip-loops/AC
  (3rd forces cause ≠ slip). Re-entries append iterations; rejected design
  options carry forward by ID.
- **Branch** (new run, inherited context): requirement invalidated or
  reclassification. `wf run branch --reason …` snapshots the parent's failure
  context into the child's Frame inputs; parent lineage is a record
  (`parent_run`), closing the "improvised branch without parent_run" gap.
- **Park**: `wf park --reason …` — the always-available honest stop; clears
  the Stop-gate; resumable. Loop-cap overflow auto-parks with a root-cause
  record prompt.
- **Force**: `wf force-exit --reason …` — audited bypass, counted, surfaced;
  and (new, fixes G4) the engine **escalates**: a second force in one run
  requires naming the structural cause, and a third auto-parks the run with a
  repair checklist instead of another bypass.

## 7. The no-leak funnel on native tasks

- Plan tasks and mid-run `followup` discoveries are mirrored into the native
  task list (`TaskCreate`, with dependencies). The engine's checklist state is
  authoritative (01 §7 persistence caveat); a `SessionStart`/`resume` sync
  re-creates any missing native tasks from state.
- `TaskCompleted` hook = the live no-leak gate (a task cannot close without
  its DoD records). Run-level: `wf phase exit` on Ship requires every task
  resolved or explicitly carried (`followup → next-run` record) — the Phase-0
  triage of v0.36 becomes `wf run start` printing carried followups as Frame
  inputs.

## 8. What was cut, with rationale

| Cut | Rationale |
|---|---|
| 12 rendered route folders, route-only/skip markers, route banners | Replaced by 3 family procedures + intent tags (§2); eliminates B-class entirely |
| Trace as an agent-performed phase | Engine computes it (§4.7); agent only resolves findings |
| The `record_intent` 1:1 intent↔edit ledger + `pre_edit` block | Its goal (rationale per edit) is served better by task binding: every Build edit occurs under an active task with a DoD; PostToolUse ledgers edits→task automatically. The intent gate was the workflow's largest honest-agent friction source (C7, G-report) with basename-bypass holes; task binding gives strictly more context per edit with zero extra agent calls. A PreToolUse guard remains only for edits *outside any active task/phase* (04 §5) |
| Perspective *files* as separate loadables + `pre_read` ordinal police | Lenses become reviewer-agent prompts + a Frame checklist; phase sequencing is enforced by gating the `Skill` tool + the Stop/exit gates (04 §5), not by policing `Read` |
| SYS-NNN Stage A vs SWR Stage B two-tier requirements | One `requirement` record kind with a `level: system|software` field; the two-stage ceremony collapsed into Frame; RTM is engine-generated from records |
| DORA/production-signal manual backfill prose, PDCA retro liturgy | `wf report` computes what's measurable; unmeasurable tiers dropped from agent obligations (kept as optional record kinds) |
| Regulated-profile per-standard reviewers, SBOM emission, 13 domain standards files | Out of scope for the core plugin (Claude-only personal/team tool). Re-introducible later as a companion plugin; noted in 09 open questions |
| `explain`-era outputs, stakeholder-registry RACI ceremony, workflow.md dispatcher file | No consumer / replaced by engine state injection |
| decisions.md as a parsed artifact; the 150-line/TOC token budgets; editable regions | Generated narrative (§5); budgets replaced by skill re-injection caps (05 §3) |
