# 07 — The Engine and the Implementation Language

The `wf` engine is the single executable behind every hook, gate, and record
command. This doc fixes its responsibilities and command surface, evaluates
implementation languages against the verified runtime facts
([01 §11](01-claude-native-capabilities.md#11-runtime-environment-facts)), and
recommends one.

## 1. Engine responsibilities

One binary, five jobs — everything else in the system is markdown or the
declarative workflow spec (03 §4.0), which the engine loads, validates, and
interprets (spec + `contracts.d/` merge) on every invocation:

1. **State store** (08): append-only event log + derived state under
   `.workflow/`; atomic writes; schema-versioned; monotonic IDs.
2. **Gates**: `wf gate stop|task-create|task-complete|verdict|skill|edit` —
   the hook entry points; `wf phase exit` — the phase contract evaluator
   (exit 0/2/3).
3. **Recording**: `wf record <kind>`, `wf approve <gate>`, `wf loop`,
   `wf park`, `wf force-exit`, `wf capture test` (PostToolUse payload →
   grounded test-run), `wf run start|branch|adopt|close`, `wf phase waive`,
   `wf deps check`, `wf origin discover`, `wf doc new`, `wf risk scan`.
4. **Injection**: `wf inject session|turn|agent <name>` — the context
   payloads (05 §5), generated from state in milliseconds.
5. **Reporting/maintenance**: `wf status`, `wf trace`, `wf report [--json]`,
   `wf lessons suggest|accept|reject|apply`, `wf doctor`, `wf selftest`.

Every hook in `hooks/hooks.json` is exec-form referencing the engine by
**absolute placeholder path**:
`{"type":"command","command":"${CLAUDE_PLUGIN_DATA}/bin/wf","args":["gate","stop"]}`.
Not bare `wf` on PATH: plugin `bin/` is documented for the *Bash tool's* PATH
only, not hook subprocesses (01 §9) — the Bash-tool PATH remains a
convenience for the agent's own `wf record …` calls. Hook stdin JSON goes
straight to the engine — no shell parsing layer at all (removes the entire
bash+heredoc hazard class, C13/C14, and the G1 quoting/matching bugs'
habitat). How the binary lands in `${CLAUDE_PLUGIN_DATA}/bin/` is §4.

## 2. Language requirements

| Requirement | Why |
|---|---|
| R1: Start fast | PreToolUse/PostToolUse/UserPromptSubmit fire on **every** tool call/prompt; UserPromptSubmit has a 30s budget but perceived latency matters at ~100s of invocations per session |
| R2: Zero runtime dependency | Claude Code guarantees **no Node, no Python, no Git Bash** (native binary; Windows-without-Git-Bash supported). The engine must run where Claude Code runs: mac/linux/win, x64+arm64, incl. musl |
| R3: One codebase, heavily testable | The v0.36 postmortem is unambiguous: hand-mirrored logic drifts; enforcement code needs unit + integration tests (selftest) |
| R4: Single-file distribution | Ships inside a plugin that is cache-copied per version |
| R5: JSON-native ergonomics | Everything is JSON in/out |

## 3. Options analysis

| Option | R1 startup | R2 deps | R3 testability | R4 distribution | Notes |
|---|---|---|---|---|---|
| **Bash + embedded python3** (status quo) | ~50–120ms | ✗ python3 not guaranteed; ✗ native Windows | poor (proven by history) | text files (easy) | Eliminated by R2/R3 |
| **Python (zipapp .pyz)** | ~60–150ms (interpreter start) | ✗ python3 required — absent on native Windows and many mac defaults | good | single .pyz | Familiar (current codebase), but fails R2 |
| **TypeScript on Node** | ~50–100ms | ✗ node required; plugin can `npm install` into `CLAUDE_PLUGIN_DATA` but that presumes npm/node exist | good | needs node_modules or bundling | Ecosystem-adjacent, fails R2 |
| **Bun-compiled single binary** | ~10–30ms | ✅ self-contained | good | ~50–90MB **per platform** | Works, but binary size × 8 platforms is unwieldy in a cache-copied plugin |
| **Go static binary** | **~5–15ms** | ✅ none; trivial cross-compile to all 8 supported targets incl. musl/arm64/windows | very good (table-driven tests, race detector) | **~6–10MB per platform** | Recommended |
| **Rust static binary** | ~5ms | ✅ | very good | ~3–8MB | Equal fitness; higher authoring cost; no advantage over Go for this workload |

## 4. Recommendation: Go

**Go static binaries, cross-compiled per platform.** It is the only option
that meets R1+R2 outright with small artifacts, and its plain, test-friendly
style suits a tool whose whole value is *correct, boring enforcement logic*.
Rust is an acceptable substitute if maintainer preference dictates;
TypeScript/Bun is the fallback if a compiled toolchain is unacceptable
(accepting the size cost); Python is viable **only** if native-Windows support
is explicitly waived — record that trade-off if taken.

**The workflow spec is the single source; the engine is its interpreter**
(03 §4.0): `workflow/workflow.yaml` declares phases, families, record kinds,
contract items, the gating roster, and the verdict vocabulary. Go types +
JSON Schemas validate the spec (load-time and CI); `go generate` emits the
spec-derived views — `hooks/hooks.json` (SubagentStart/Stop matchers from the
roster), the agents' verdict-block section, the record JSON schemas shipped
in `workflow/schemas/` — and CI fails if generated files drift (structurally
kills A1/D1/D2/D5). Contract *changes* are spec edits, not engine releases;
only a new predicate type touches Go.

### Binary distribution inside the plugin

Both patterns converge on one invariant: hooks reference a **stable absolute
path** — `${CLAUDE_PLUGIN_DATA}/bin/wf` — because exec-form hooks cannot rely
on PATH and, on Windows, cannot spawn dispatcher *scripts* (01 §1/§9). A
`SessionStart` **bootstrap hook** (the only shell scripts in the system:
paired `.sh` + `.ps1` with `shell:"powershell"` on the Windows entry, since
it must run before the binary exists) installs the engine there,
checksum-verified, then re-runs `wf doctor --bootstrap` to confirm:

- **A — bundle all platforms (default)**: per-platform binaries under plugin
  `bin/` (`wf-darwin-arm64`, `wf-windows-x64.exe`, …); bootstrap copies the
  matching one to `${CLAUDE_PLUGIN_DATA}/bin/wf` (plus `wf.exe` on Windows —
  the copy at the hook-referenced path is a real executable image on every
  platform, so exec form spawns it directly). ~60–80MB plugin total:
  acceptable for a per-version cache; fully offline; binaries are CI-built
  release artifacts on the release tag (or attached via the marketplace
  `sha`-pinned source). A thin `wf` shim in plugin `bin/` keeps the bare
  command available inside the Bash tool.
- **B — fetch-on-first-use**: bootstrap downloads the platform binary into
  `${CLAUDE_PLUGIN_DATA}/bin/` (the documented dependency pattern), verified
  against SHA256 checksums committed in the plugin. Small plugin, but
  first-run needs network.

Update handling: bootstrap compares the installed binary's version against
the plugin's bundled manifest and reinstalls on mismatch (the documented
`CLAUDE_PLUGIN_DATA` diff-and-reinstall pattern) — hooks keep one stable path
across plugin updates.

## 5. Engine design constraints

- **Latency budget**: <20ms for `gate`/`inject`/`capture` paths (state reads
  from a compact current-state file, not log replay; log replay only in
  `doctor`/`report`).
- **Atomicity**: single-writer lockfile per `.workflow/`; every mutation =
  append event + rewrite derived snapshot atomically (temp+rename). `wf run
  close` is one transaction (fixes the A5/A6 ordering class).
- **Fail-safe split** (04 §7): sequencing gates fail *open+loud* when the
  engine is broken; recording paths fail *closed* (never fabricate).
  `parked`/`force` handling is independent of contract evaluation (no-wedge).
- **Schema versioning**: state events carry `schema`; engine reads N and N-1,
  upgrades lazily (02 §5). No other migration machinery exists.

## 6. Testing strategy

The as-is system's core lesson: *enforcement claims must be adversarially
tested*. Three layers:

1. **Unit**: the predicate evaluators (each of the closed vocabulary,
   table-driven against record fixtures), verdict parsing, capture filters
   (incl. the G1 cases: hook self-calls, runner names inside quoted args,
   filter pipes), record schema validation. Plus **spec validation**: every
   contract item references declared record kinds/agents/phases; every
   gating-roster entry has an agent file and a hooks matcher; `contracts.d/`
   merging is add-only (03 §4.0).
2. **Golden/integration** (`wf selftest`): a scripted fake session drives the
   hook entry points with recorded Claude-Code JSON payloads (Stop with
   `stop_hook_active`, SubagentStop with verdict/malformed/missing blocks,
   TaskCompleted with unmet DoD…) and asserts block/allow decisions and state
   effects. Runs in CI on all platforms; also invocable in a live scaffold.
3. **Adversarial E2E** (release gate, 09): drive real Claude Code with a
   deliberately non-compliant system prompt ("skip whatever you can") against
   a sample repo and assert from the resulting state that every gate held —
   the mechanized version of the manual audits that produced
   `claude_issues.md`.
