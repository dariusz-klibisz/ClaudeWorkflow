# 01 — Verified Claude-Native Capabilities (Evidence Base)

Every mechanism the redesign relies on, verified against the official Claude Code
documentation (code.claude.com/docs, retrieved 2026-07). Minimum versions noted
where the docs state them. This is the *evidence base*: later documents cite this
one instead of re-arguing feasibility.

## Table of contents
1. [Hook system fundamentals](#1-hook-system-fundamentals)
2. [The blocking matrix](#2-the-blocking-matrix)
3. [Key events in detail](#3-key-events-in-detail)
4. [Prompt and agent hooks](#4-prompt-and-agent-hooks)
5. [Compaction: what survives, what hooks fire](#5-compaction-what-survives-what-hooks-fire)
6. [Subagent lifecycle](#6-subagent-lifecycle)
7. [Native task system](#7-native-task-system)
8. [Skills](#8-skills)
9. [Plugins and marketplaces](#9-plugins-and-marketplaces)
10. [Memory and rules](#10-memory-and-rules)
11. [Runtime environment facts](#11-runtime-environment-facts)
12. [Corrections to earlier assumptions](#12-corrections-to-earlier-assumptions)

---

## 1. Hook system fundamentals

Source: `docs/en/hooks` (reference), `docs/en/hooks-guide`.

- Hook handlers: `command` (shell/exec), `http`, `mcp_tool`, `prompt` (LLM
  yes/no), `agent` (verifier subagent). All matching hooks run **in parallel**;
  identical handlers dedupe.
- **Exec form**: when `args` is present, `command` is spawned directly with no
  shell — use for anything referencing `${CLAUDE_PLUGIN_ROOT}` (no quoting
  hazards). Shell form (`args` absent) uses `sh -c` / Git Bash / PowerShell;
  per-hook `shell: "powershell"` runs PowerShell regardless of platform config.
  ⚠ On Windows, exec form requires `command` to resolve to a **real
  executable** (`.exe`) — shell scripts and `.cmd`/`.bat` shims cannot be
  spawned without a shell. A dispatcher *script* is therefore not a viable
  exec-form hook target on native Windows (see 07 §4).
- **Exit-code protocol**: exit 0 = success (stdout parsed for JSON output);
  **exit 2 = blocking error** (stderr fed to Claude; any JSON on stdout is
  ignored); **any other code (incl. 1) = non-blocking** — "If your hook is meant
  to enforce a policy, use exit 2." JSON output is processed only on exit 0.
  (Sole documented exception: `WorktreeCreate` aborts on any non-zero exit.)
- JSON output universals: `continue:false` (+`stopReason` shown to the user)
  stops processing entirely and takes precedence over event decisions;
  `systemMessage`; `suppressOutput`. Output strings incl. `additionalContext`
  are capped at **10,000 characters** (overflow saved to a file).
- **`additionalContext`** injects a string into Claude's context as a system
  reminder at the point the hook fired (SessionStart → conversation start;
  UserPromptSubmit → alongside the prompt; Pre/PostToolUse → next to the tool
  result; Stop/SubagentStop → end of turn, conversation continues). Caveat:
  "text framed as out-of-band system commands can trigger Claude's
  prompt-injection defenses" — payloads must be factual statements.
- `if` field: one permission rule (`Bash(git *)`, `Edit(*.ts)`,
  `Skill(name *)`) filters tool events before spawning the handler —
  best-effort (fails open on unparseable commands), so hard policy belongs in
  the permission system or the handler itself.
- Hook locations: `~/.claude/settings.json`, `.claude/settings.json`,
  `.claude/settings.local.json`, managed policy, **plugin `hooks/hooks.json`**,
  and **skill/agent frontmatter** (scoped to the component's lifetime; all
  events supported; `Stop` in agent frontmatter auto-converts to
  `SubagentStop`; `once: true` honored only in skill frontmatter).
- Path placeholders, substituted AND exported as env vars:
  `${CLAUDE_PROJECT_DIR}`, `${CLAUDE_PLUGIN_ROOT}` (per-version install dir),
  `${CLAUDE_PLUGIN_DATA}` (persistent across updates). Plugin hooks also get
  `${user_config.*}`.
- `async: true` runs in background (cannot block); **`asyncRewake`**: a
  background hook exiting 2 wakes Claude immediately with its stderr as a
  system reminder.
- `CLAUDE_ENV_FILE` (SessionStart, Setup, CwdChanged, FileChanged): append
  `export` lines → available to all subsequent Bash commands in the session.
- `disableAllHooks` setting exists; `/goal` and hooks require the workspace
  trust dialog to have been accepted.

## 2. The blocking matrix

Which events can actually stop something (verified per-event table):

| Event | Blocks? | Exit-2 / block effect |
|---|---|---|
| `PreToolUse` | ✅ | Tool call denied; reason shown to Claude. JSON: `permissionDecision: allow\|deny\|ask\|defer`, `updatedInput` (replace tool input), `additionalContext`. Precedence when hooks disagree: deny > defer > ask > allow. Silence never approves |
| `UserPromptSubmit` | ✅ | Prompt blocked and erased; or exit 0 + `additionalContext` injection. 30s default timeout — a timed-out hook's context is discarded |
| `Stop` | ✅ | **Prevents Claude from stopping**; `decision:"block"` + `reason` (required) is delivered to Claude as its next instruction. Guards: `stop_hook_active` input flag; hard cap **8 consecutive blocks**; `last_assistant_message`, `background_tasks`, `session_crons` provided. Does not fire on user interrupt; API errors fire `StopFailure` (output ignored) |
| `SubagentStop` | ✅ | Prevents the subagent from finishing; `reason` becomes its next instruction. Input: `agent_id`, `agent_type`, `agent_transcript_path`, `last_assistant_message`. Same 8-block cap |
| `TaskCreated` | ✅ | Rolls back task creation; stderr fed to the model |
| `TaskCompleted` | ✅ | **Task is not marked complete**; stderr fed to the model. Fires on explicit `TaskUpdate` completion (and teammate turn-end with in-progress tasks). Docs example: test-gated completion |
| `PostToolBatch` | ✅ | Stops the agentic loop before the next model call |
| `PreCompact` | ✅ | Blocks compaction (observe-only otherwise; **no documented mechanism to inject summarizer instructions** — decision control is block-only, and PreCompact is absent from the `additionalContext` delivery list) |
| `UserPromptExpansion` | ✅ | Blocks a slash-command expansion |
| `PermissionRequest` | ✅ | Denies the permission |
| `PostToolUse` | ❌ | Tool already ran; exit 2 shows stderr to Claude; JSON `decision:"block"`+`reason`, `additionalContext`, `updatedToolOutput` (replace what Claude sees) |
| `SessionStart`, `SubagentStart`, `PostCompact`, `SessionEnd`, `Notification`, `InstructionsLoaded`, `FileChanged` | ❌ | Observability/injection only |

## 3. Key events in detail

- **`SessionStart`** — matchers `startup | resume | clear | compact`. Output:
  `hookSpecificOutput.additionalContext` (added at conversation start),
  `initialUserMessage` (headless), `watchPaths` (dynamic FileChanged watch
  list), `reloadSkills` (re-scan skills after the hook installs some). Only
  `command`/`mcp_tool` handler types. Re-runs on resume with `source:"resume"`.
- **`SubagentStart`** — cannot block, but **injects `additionalContext` into
  the subagent's context** before its first prompt. Matcher = agent type
  (frontmatter `name`; plugin agents as `^plugin:name$` — colon puts it on the
  regex path, anchor it).
- **`PostToolUse` on the `Agent` tool** — `tool_response` carries the
  subagent's final text + telemetry; **but** as of v2.1.198 subagents run in
  the background by default, returning `status:"async_launched"` with no
  content. Verdict capture therefore anchors on `SubagentStop`, not here.
- **`PreCompact`** — input `trigger` (`manual|auto`) + `custom_instructions`
  (user's `/compact` argument; input-only). Can block; cannot modify the
  summary.
- **`PostCompact`** — input includes `compact_summary` (read-only).
- **`InstructionsLoaded`** — observability of CLAUDE.md/`.claude/rules/*.md`
  loads; matcher on `load_reason` incl. `compact`.
- **`FileChanged`** — matcher is a `|`-separated list of *literal filenames*
  to watch; `watchPaths` extends dynamically.
- **`Setup`** — fires only with `--init-only` / `--init`/`--maintenance` in
  `-p` mode; NOT on normal startup. First-run dependency installs must
  therefore be check-and-install in `SessionStart`.
- Common input on every hook: `session_id`, `transcript_path`, `cwd`,
  `permission_mode`, `hook_event_name`; inside subagents additionally
  `agent_id`, `agent_type`.

## 4. Prompt and agent hooks

- `type:"prompt"`: hook input + prompt sent to a fast model (default Haiku,
  30s); returns `{ok: true|false, reason}`. On Stop/SubagentStop, `ok:false`
  feeds `reason` to Claude and continues the turn (this is exactly what
  `/goal` is — "a wrapper around a session-scoped prompt-based Stop hook").
  On PreToolUse it equals `deny`. `continueOnBlock` controls
  PostToolUse/TeammateIdle behavior.
- `type:"agent"`: same contract but spawns a verifier subagent with
  Read/Grep/Glob, up to 50 turns, 60s default timeout. **Documented as
  experimental** ("prefer command hooks for production").
- Supported on: PreToolUse, PostToolUse(+Failure/Batch), Stop, SubagentStop,
  TaskCreated, TaskCompleted, UserPromptSubmit/Expansion, PermissionRequest.
  NOT on SessionStart/SubagentStart/PreCompact/PostCompact.

## 5. Compaction: what survives, what hooks fire

Source: `docs/en/context-window` ("What survives compaction", verbatim),
`docs/en/memory`.

| Mechanism | After compaction |
|---|---|
| System prompt and output style | Unchanged (not part of message history) |
| **Project-root CLAUDE.md and unscoped `.claude/rules/`** | **Re-injected from disk** |
| **Auto memory (MEMORY.md first 200 lines / 25KB)** | **Re-injected from disk** |
| Rules with `paths:` frontmatter | Lost until a matching file is read again |
| Nested CLAUDE.md in subdirectories | Lost until a file in that subdirectory is read again |
| **Invoked skill bodies** | **Re-injected**, capped at 5,000 tokens per skill and 25,000 total; oldest dropped first; truncation keeps the *start* of the file |
| Skill descriptions index | **Not re-injected** — only skills actually invoked are preserved |
| Hooks | N/A — hooks run as code, not context |
| Conversation | Replaced by a structured summary (keeps: user intent, key concepts, files touched w/ snippets, errors+fixes, **pending tasks**, current work) |

Programmatic re-anchoring channels, in firing order around a compaction:
`PreCompact` (observe/block) → compaction → `PostCompact` (gets summary) →
**`SessionStart` matcher `compact`** (inject `additionalContext` regenerated
from disk state) → `InstructionsLoaded` events with `load_reason:"compact"`.
Additionally `UserPromptSubmit` can re-anchor on *every* prompt. Subagent
transcripts are unaffected by main-conversation compaction; subagents
auto-compact independently.

## 6. Subagent lifecycle

Source: `docs/en/sub-agents`.

- Frontmatter fields: `name`, `description` (required); `tools`,
  `disallowedTools`, `model` (alias/full-ID/`inherit`; default inherit),
  `permissionMode`, `maxTurns`, `skills` (**preload full skill content into
  the subagent at startup**), `mcpServers`, `hooks` (lifecycle-scoped),
  `memory` (`user|project|local` → persistent `agent-memory/<name>/` dir with
  auto-loaded MEMORY.md; `project` is shareable via VCS), `background`,
  `effort`, `isolation: worktree`, `color`, `initialPrompt`.
- **Plugin-shipped agents do NOT support `hooks`, `mcpServers`, or
  `permissionMode`** (ignored for security). Enforcement for plugin agents
  must live in plugin `hooks/hooks.json` (SubagentStart/SubagentStop matchers
  on the scoped name) — verified supported.
- Scope priority: managed > `--agents` CLI > `.claude/agents/` (project) >
  `~/.claude/agents/` (user) > plugin `agents/`. A project agent with the same
  name overrides a plugin agent.
- **Background by default** (v2.1.198); foreground only when the result is
  needed immediately; background permission prompts surface in the main
  session (v2.1.186+). `background: true` forces it. Ctrl+B backgrounds a
  running task.
- Nested subagents (v2.1.172+): allowed by listing `Agent` in `tools`; depth
  limit **5, fixed**. `Agent(worker, researcher)` allowlist syntax restricts
  spawnable types for a main-thread `--agent`; `permissions.deny:
  ["Agent(name)"]` blocks types session-wide.
- **Resume**: completed subagents return an agent ID; `SendMessage` (no agent
  teams needed) resumes them with full prior context. Explore/Plan are
  one-shot.
- Startup context of a custom subagent: its own system prompt + task message +
  **CLAUDE.md hierarchy + project rules + git snapshot** (Explore/Plan skip
  these) + preloaded skills. Not the parent conversation.
- Transcripts: `~/.claude/projects/<project>/<session>/subagents/agent-<id>.jsonl`
  (also handed to SubagentStop as `agent_transcript_path`).
- `claude --agent <name>` / project setting `"agent"` runs a whole session as
  an agent (its prompt replaces the default system prompt).

## 7. Native task system

Source: `docs/en/tools-reference`, hooks reference.

- Tools: `TaskCreate`, `TaskUpdate` (status, **dependencies**, details,
  delete), `TaskGet`, `TaskList` — the default since v2.1.142 (TodoWrite
  disabled unless `CLAUDE_CODE_ENABLE_TASKS=0`). No permission prompts.
- Hook gates: `TaskCreated` (exit 2 rolls back creation) and `TaskCompleted`
  (exit 2 prevents completion, stderr fed to model) fire for **any agent's**
  explicit TaskUpdate completion. Input: `task_id`, `task_subject`,
  `task_description?`, `teammate_name?`.
- ⚠ Verified gap: task-list **persistence/storage is not documented** (the
  compaction summary "keeps pending tasks", but disk location/resume behavior
  is unspecified). Design consequence: the engine's on-disk checklist is
  authoritative; native tasks are an enforced *mirror* (see 04 §3).

## 8. Skills

Source: `docs/en/skills`, plugins guide, context-window doc.

- `SKILL.md` with frontmatter: `name`, `description`,
  `disable-model-invocation: true` (user-only; zero context until invoked),
  `context: fork`, `hooks` (+ `once`), `$ARGUMENTS` placeholder. Supporting
  files alongside enable progressive disclosure.
- Invocation: `/name` (or `/plugin:name` for plugin skills); Claude
  auto-invokes based on description unless disabled; runs through the `Skill`
  tool → **`PreToolUse` with matcher `Skill` and `if: "Skill(name *)"` can
  gate skill invocation**.
- Compaction: invoked bodies re-injected (5k/25k caps, start-of-file kept);
  description index not reloaded (§5).
- Skills preload into subagents via the agent `skills` field (full content).

## 9. Plugins and marketplaces

Source: `docs/en/plugins`, `docs/en/plugins-reference`, `docs/en/plugin-marketplaces`.

- Plugin components: `skills/`, `commands/`, `agents/`, `hooks/hooks.json`,
  `.mcp.json`, `.lsp.json`, `monitors/`, **`bin/` (added to the Bash tool's
  PATH while enabled)**, `settings.json` (only `agent`,
  `subagentStatusLine`), themes. Manifest `.claude-plugin/plugin.json`
  (only `name` required; `version` optional).
  ⚠ `bin/` is documented for the **Bash tool's PATH only** — nothing states
  hook subprocesses inherit it. Hook entries must therefore reference the
  engine by absolute placeholder path (`${CLAUDE_PLUGIN_ROOT}/…` or
  `${CLAUDE_PLUGIN_DATA}/…`), never rely on PATH resolution (07 §4).
- **Repo-as-marketplace**: `.claude-plugin/marketplace.json` at the repo root
  listing plugins with relative sources (`"./"` or `"./plugins/x"`). Users:
  `/plugin marketplace add owner/repo`. Sources: github, git URL, git-subdir
  (sparse clone), npm; `ref`/`sha` pinning.
- **Team auto-install**: project `.claude/settings.json` with
  `extraKnownMarketplaces` + `enabledPlugins` — collaborators are prompted on
  trusting the folder. Private repos: git credential helpers; background
  auto-update needs `GITHUB_TOKEN`/`GITLAB_TOKEN`/`BITBUCKET_TOKEN`.
- **Versioning**: resolution order plugin.json `version` → marketplace entry
  `version` → git commit SHA (omit version ⇒ every commit is a release).
  Release channels via two marketplaces on different refs. `renames` map
  migrates renamed/removed plugins. `/plugin marketplace update` +
  auto-updates refresh.
- **Caching**: marketplace plugins are copied to `~/.claude/plugins/cache`
  per version; **no `../` references** survive install; old versions linger
  ~7 days. `${CLAUDE_PLUGIN_ROOT}` changes per update;
  `${CLAUDE_PLUGIN_DATA}` (`~/.claude/plugins/data/<id>/`) persists across
  updates — documented pattern: diff a bundled manifest against a stored copy
  and (re)install deps on mismatch in a SessionStart hook.
- Dev/fallback loading: `--plugin-dir ./dir` (also `.zip`), `--plugin-url`,
  skills-directory plugins (`.claude/skills/<x>/.claude-plugin/plugin.json`),
  and `CLAUDE_CODE_PLUGIN_SEED_DIR` for containers/CI (read-only,
  pre-populated cache). Local `--plugin-dir` overrides a same-named installed
  plugin for the session.
- Enterprise controls: `strictKnownMarketplaces`, `blockedMarketplaces`,
  `allowManagedHooksOnly`, `disableSideloadFlags`.

## 10. Memory and rules

Source: `docs/en/memory`.

- Project instructions: `./CLAUDE.md` or `./.claude/CLAUDE.md` (team,
  committed); `CLAUDE.local.md` (personal); `~/.claude/CLAUDE.md` (user);
  managed policy file/`claudeMd` key. `@path` imports (4-hop depth,
  approval-gated outside the repo). HTML comments are stripped before
  injection. Target <200 lines.
- `.claude/rules/*.md`: modular rules, recursive discovery, symlinks OK;
  `paths:` frontmatter scopes a rule to file globs (loads on matching reads);
  unscoped rules load at launch with project-CLAUDE.md priority and are
  re-injected after compaction.
- Auto memory (`~/.claude/projects/<project>/memory/MEMORY.md` + topic
  files): Claude-maintained; first 200 lines/25KB loaded each session;
  re-injected post-compaction; `autoMemoryDirectory` relocatable; per-agent
  memory via the subagent `memory` field (§6).
- Instructions are *context, not enforcement*: "To block an action regardless
  of what Claude decides, use a PreToolUse hook instead."

## 11. Runtime environment facts

Source: `docs/en/setup`, tools reference, whats-new.

- Claude Code is a **native binary**; the npm package also just links a
  platform binary — "the installed `claude` binary does not itself invoke
  Node." **No Node, Python, or Git Bash is guaranteed on a user's machine.**
- Native Windows without Git Bash is supported (PowerShell tool auto-enabled;
  hooks can set `shell: "powershell"`, which "works regardless of
  CLAUDE_CODE_USE_POWERSHELL_TOOL"). With Git Bash, hooks default to bash.
- Platforms: macOS 13+, Windows 10 1809+, Ubuntu 20.04+/Debian 10+/Alpine
  3.19+ (musl needs libgcc/libstdc++), x64 + ARM64.
- Bash tool: 2-minute default / 10-minute max timeout; 30k-char output cap;
  per-command processes; `run_in_background`.
- Hook spawn cadence: PreToolUse/PostToolUse fire on **every tool call** —
  handler startup latency is a real cost (see 07).
- `/hooks` menu, `claude --debug-file`, `claude plugin validate` for
  debugging; `/context`, `/memory` for context inspection.

## 12. Corrections to earlier assumptions

Deltas found during verification that changed the design:

1. **PreCompact offers no documented way to inject summarizer instructions**
   (block-only decision control; `custom_instructions` is input-only; implied
   by omission rather than stated verbatim) — re-anchoring must be
   post-compaction (`SessionStart(compact)`), which is the more robust design
   anyway (regenerated from disk, not negotiated with a summarizer).
2. **Plugin agents can't carry frontmatter hooks** — the SubagentStop verdict
   gate lives in the plugin's `hooks/hooks.json`, matched by scoped agent
   names.
3. **Subagents are background-by-default (v2.1.198)** — verdict capture moves
   from `PostToolUse(Agent)` (may be `async_launched`, contentless) to
   `SubagentStop` (always has `last_assistant_message` +
   `agent_transcript_path`). Bonus: review rosters parallelize naturally.
4. **Invoked skill bodies DO survive compaction** (5k/25k token caps,
   start-kept) — phase-procedure skills are safer than previously assumed;
   critical contract goes at the top of each SKILL.md.
5. **Task-list persistence is undocumented** — native tasks are the enforced
   mirror, never the source of truth.
6. **Exit code 1 does not block** — every enforcing hook must use exit 2 (or
   exit 0 + JSON decision); this is a documented footgun.
7. **`/goal` exists** and is implemented exactly as our Stop-gate pattern —
   strong evidence the pattern is supported and load-bearing upstream.
8. **Agent hooks are experimental** — the deterministic command-hook gate is
   the core; prompt/agent hooks are an additive layer only.
