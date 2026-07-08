# 08 — State and Storage

Where everything lives: the three storage tiers, the state schema, and the git
policy. Governing rules: the gate reads only engine-written records (03 §5);
nothing persistent lives in `${CLAUDE_PLUGIN_ROOT}` (per-version cache);
per-project state is minimal and mostly committed.

## 1. The three tiers

| Tier | Path | Lifetime | Contents |
|---|---|---|---|
| Plugin (read-only) | `${CLAUDE_PLUGIN_ROOT}` | per plugin version | engine binary, hooks.json, skills, agents, templates, corpora |
| Machine data | `${CLAUDE_PLUGIN_DATA}` (`~/.claude/plugins/data/wf/`) | survives updates | fetched binaries (pattern B), caches, per-machine scratch — nothing authoritative |
| Project | `${CLAUDE_PROJECT_DIR}/.workflow/` + `docs/` | with the repo | all authoritative workflow state + deliverables |

## 2. Project layout

```
.workflow/
├── config.json          # schema, family defaults, ux flag, thresholds, plugin version at init   ✅ committed
├── contracts.d/         # project-local additive contract items + record kinds (03 §4.0),
│   │                    #   incl. engine-written lessons.yaml (accepted lesson checks)           ✅ committed
├── state/
│   ├── run.json         # current run snapshot: id, family, intent, phase, checklist, waiting-on ✅ committed
│   └── lock             # single-writer lockfile                                                 ❌ ignored
├── log/
│   └── events.jsonl     # THE append-only event log (all record kinds, §3), current runs         ✅ committed
├── runs/<run-id>/       # closed-run archives: events slice + generated narrative + signals      ✅ committed
└── local/               # per-machine: session cursors, native-task mirror map, caches           ❌ ignored (.workflow/.gitignore)
docs/                    # deliverables (PROJECT.md, ADRs, design docs, reviews, releases…)       ✅ committed
```

Key departures from v0.36:
- **Run identity is committed** (`state/run.json` + events). A fresh clone or
  second machine sees the in-flight run; `wf run adopt` re-attaches a session
  to it. The G2 failure (gitignored cursors + committed narrative diverging →
  force cascade) cannot occur: there is one source and it travels with the
  repo. Merge conflicts on `events.jsonl` are append-only line unions —
  resolvable by union merge because event identity is ULID-based (§3), so
  parallel branches can't mint colliding IDs; the `run.json` snapshot is
  never hand-merged — `wf doctor` re-derives it from the merged log.
- **One log, not five.** steps/tasks/ledger/evidence/decisions collapse into
  typed events in a single ordered log (+ the derived `run.json` snapshot for
  O(1) gate reads). `decisions.md` no longer exists as a source — `wf log
  --md` renders the narrative on demand; a rendered copy is frozen into each
  run archive for human history.

## 3. Event schema (v2)

Envelope (every event):

```json
{"schema":1,"id":"01JZ8M3H7VQ4T9E2R6W8XKPDNC","seq":184,
 "ts":"2026-07-05T12:34:56Z","run":"20260705-a1b2",
 "phase":"build","kind":"…","auto":false,"actor":"agent|engine|hook","note":"…"}
```

`id` is an engine-issued **ULID** — the identity used for cross-references
and dedup, globally unique so events created on parallel git branches or
machines never collide on merge. `seq` is a per-stream ordering hint
(monotonic per writer; never a line count — A6) — ties across merged streams
are ordered by `id`, and `wf doctor` re-derives snapshots from the merged
log. `phase` is stamped by the engine from live state (never replayed from
labels — C5/C6). `auto:true` only from hook-capture paths; sticky-evidence
rule in 04 §4.

| kind | payload (abridged) | producers |
|---|---|---|
| `run` | start/branch/adopt/close, family, intent, parent | `wf run *` |
| `phase` | enter/exit/waive/loop/park/force, target, reason | engine transitions |
| `classification` | family, intent, restated task | Frame |
| `risk` | signals[] + lens bindings | `wf risk scan` + agent |
| `ambiguity` | lens, text, disposition | Frame |
| `requirement` | id, level, text, acs[] (each `verifiable`), status transitions | Frame/Context |
| `completeness` | items[] + dispositions | Frame |
| `assumption` | text, high_risk | Context |
| `context-map` | entries[], sufficiency | Context |
| `option-set` | stage, candidates[], selected, rejected[]+reasons | Design |
| `threat` / `attack-path` | STRIDE/tree entries, mitigation status | Design |
| `task` | id, dod[], ac_links[], status | Plan/Build (mirrored to TaskCreate) |
| `verification-strategy` | ac, method, command | Plan |
| `deps` | verdict present/missing/n-a, detail | `wf deps check` |
| `test-run` | cmd, exit (int\|null), summary, ac/task tag, tests, coverage, filtered | `wf capture test` / manual |
| `verdict` | agent, status, criticals, majors, scope | SubagentStop gate / manual |
| `approval` | gate (frame/scope/design/plan/deferral/lesson), approver, payload_hash | `wf approve` — always `auto:false` |
| `artifact` | path, status present/stub/missing, template | `wf doc new` / `wf record artifact` |
| `deviation` / `disposition` / `followup` / `origin` / `commit-origin` / `metric` / `lesson` | as in 03 | various |
| `escape` | park/force/enforce-off, reason, escalation level | `wf park` / `wf force-exit` / hooks |

JSON Schemas for all kinds are generated from the engine's type definitions
(07 §4) and shipped in the plugin for external validation.

## 4. Derived views (never sources)

- `wf status` — the injection payloads (05 §5).
- `wf trace` — Ship findings: contract coverage, open items, waivers, escapes.
- `wf log --md [--run id]` — human narrative (frozen per archive).
- `docs/requirements/RTM.md` — generated requirements-traceability view
  (`wf rtm render`), committed for reviewers but headed "generated — edit
  records, not this file".
- `wf report [--json]` — the health signals (kept from v0.36's honest tiers:
  loops per run, escapes, self-attested counts, ungrounded ACs, lesson
  efficacy, deliver-reached).

## 5. Git policy

- Committed: `config.json`, `contracts.d/**`, `state/run.json`,
  `log/events.jsonl`, `runs/**`, `docs/**` — the full audit trail and the
  project's contract extensions travel with the repo.
- Ignored: `state/lock`, `local/**` (session cursors, mirrors, caches),
  everything under `${CLAUDE_PLUGIN_DATA}`.
- The engine never writes outside `.workflow/`, `docs/`, and
  `${CLAUDE_PLUGIN_DATA}`.

## 6. Growth and archival

- `wf run close` moves the run's events from `log/events.jsonl` into
  `runs/<id>/events.jsonl` (atomic transaction: archive → verify counts →
  compact → clear snapshot — one command, fixing the A5 ordering gap
  including the terminal event, which is written *before* the move as part of
  the same transaction).
- Preserved in the live log across compaction: open `followup`s ONLY — the
  bounded-live-log rule (revised: the original design kept `commit-origin`
  and `lesson` events live forever, so the log grew without bound). Lesson
  and commit-origin events archive with their run; their readers (lesson
  regeneration, `wf origin discover`) fold archived events back in via the
  committed `runs/<id>/` slices. Run close also prunes the per-machine
  `local/` counters (tasks-mirror, verdict-attempts, stop-gate).
- The live log is hash-chained (`prev` per line, sha256/16) and re-anchored
  by engine rewrites; `wf doctor` verifies the chain and reports torn or
  foreign lines — combined with the tool-gate denies on
  `.workflow/{log,state,runs}` and `config.json`, direct ledger forgery is
  blocked at the tool surface and tamper-evident past it.
- Abandoned runs: `wf doctor` flags runs idle >30 days and offers
  park-and-archive (E2). `runs/**` grows with history by design; a
  `wf archive prune --before` exists for repos that outgrow it (squashes
  event slices, keeps narratives + signals).

## 7. Concurrency and multi-session

- Lockfile single-writer; gates are read-only and lock-free.
- Two live sessions on one repo: both see the same run; the Stop gate and
  injections are per-session but state-consistent. Worktrees: `.workflow/`
  lives per-worktree (each worktree = its own run stream); `wf report` can
  aggregate across worktrees later (09 open question).
