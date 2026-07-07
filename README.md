# wf — enforced development workflow for Claude Code

wf is a Claude Code plugin that turns "please follow the process" into a
process that is actually enforced. Work happens inside **runs** that move
through gated phases; an engine (`wf`, a single Go binary behind every hook)
blocks the model from skipping reviews, faking test results, or silently
expanding scope — and records every escape it grants.

## How it works

- **Runs and phases.** Every task is a run in one of three families —
  `diff` (code changes), `artifact` (authored documents), `assessment`
  (findings reports) — moving through Frame → Context → Design → Plan →
  Build → Verify → Ship. Each phase has a contract: recorded facts,
  reviewer verdicts, and user approvals that must exist before the phase
  can exit.
- **Four gates**, wired as Claude Code hooks: the Stop gate (can't end the
  turn with unmet obligations), task gates (can't complete a task without
  captured red→green test evidence), the verdict gate (reviewer subagents
  must end with a parseable verdict, auto-captured), and tool gates
  (phase-skill sequencing, an always-on catastrophic-Bash net).
- **Grounded evidence.** Test runs are captured from the Bash tool by the
  hook itself (`auto:true`) — recognized from a built-in runner list, the
  run's own recorded verification commands (any language), or project
  config `"runners"`. Manual records stay possible but are marked
  self-attested and surface in `wf report`.
- **Audited escapes, not hidden ones.** `/wf:park` (honest stop),
  `/wf:force` (bypass one gate; escalates — the 3rd force auto-parks the
  run). Everything is recorded and reported.
- **Lessons that bite.** At Ship, the run's lessons are proposed and
  user-triaged (`wf lessons accept|reject`); accepted lessons with a
  `check:` become ordinary contract items in `.workflow/contracts.d/` —
  **enforced from the next run on** by the same evaluator as everything
  else. Prose lessons regenerate `.claude/rules/wf-lessons.md`.

## Quickstart

1. Trust the project folder in Claude Code and accept the marketplace
   prompt (or add it once):

   ```json
   // .claude/settings.json
   {
     "extraKnownMarketplaces": {
       "claude-workflow": {
         "source": { "source": "github", "repo": "dariusz-klibisz/ClaudeWorkflow" }
       }
     },
     "enabledPlugins": { "wf@claude-workflow": true }
   }
   ```

2. Run `/wf:init` — idempotent one-time adoption. It creates `.workflow/`,
   merges the settings above, writes the CLAUDE.md ground-rules block, and
   verifies the install (`wf doctor --bootstrap`).
3. **Restart the session** (two reasons: the `Bash(wf *)` permission rule
   and, if the plugin was installed mid-session, the SessionStart bootstrap
   that arms the hooks — `wf doctor --bootstrap` can also heal that on the
   spot).
4. Start every task with `/wf:dev <task>`.

## What lands in your repo

| Path | Content | Committed |
|---|---|---|
| `.claude/settings.json` | marketplace + plugin enable + `Bash(wf *)` allow — merged, never overwritten | yes |
| `CLAUDE.md` | a small marker-delimited ground-rules block | yes |
| `.workflow/` | engine state: event ledger, run snapshots, config, `contracts.d/` (incl. engine-written `lessons.yaml`) | mixed (`.workflow/.gitignore` handles it) |
| `.claude/rules/wf-lessons.md` | accepted prose lessons — engine-generated | yes |
| `docs/` | deliverables created by runs (ADRs, reviews, incidents, release notes…) | yes |

No hooks, agents, skills, or templates land in the repo — those are
plugin-resident and update with the plugin.

## Day-to-day commands

The model drives these; you mostly approve gates. Useful directly:

```
wf status                      where the run stands (authoritative, from disk)
wf report [--run <id|current>] health signals: loops, escapes, self-attested
                               counts, ungrounded ACs, lesson efficacy
wf trace                       ship close-out findings
wf lessons suggest|accept|reject|apply
wf doctor [--bootstrap]        state health · verifies AND heals the hook engine
wf selftest                    22 in-scaffold enforcement scenarios
```

## Updating

- `/plugin marketplace update` + `/plugin update wf`, then `/reload-plugins`
  (mid-session updates need it). No migrations: the engine reads state by
  `schema` with a one-release grace and upgrades lazily on write.
- **Mid-session installs leave hooks dead** (SessionStart never fired):
  any hook error naming a missing `.../data/wf-*/bin/wf` means exactly
  that — run `wf doctor --bootstrap` (installs the engine on the spot) or
  restart the session.

## Fallback installs

- **Air-gapped / pinned**: `git submodule add <this-repo> tools/workflow`,
  launch with `claude --plugin-dir ./tools/workflow`. Same plugin, no
  marketplace.
- **CI**: bake a seed directory into the image
  (`CLAUDE_CODE_PLUGIN_SEED_DIR`), read-only.
- **No marketplace reachable**: the plugin also loads from a `.zip` via
  `--plugin-url`.

## Troubleshooting

| Symptom | Meaning | Fix |
|---|---|---|
| Hook errors: `ENOENT … data/wf-*/bin/wf` | hooks dead — bootstrap never ran (mid-session install) | `wf doctor --bootstrap` or restart |
| `wf report` shows all verdicts/test-runs self-attested | capture hooks dead, or your test runner isn't recognized | doctor; make verification-strategy commands the real invocations; config `"runners"` |
| Adoption refuses: legacy scaffold | `.workflow/manifest.json` is from the old generator | remove/rename the old `.workflow/` first (no migration ships) |
| `wf` prompts for permission mid-session | settings rules apply next session | restart, or accept "don't ask again for `wf record *`" — the engine only writes under `.workflow/` |
| Native Windows (no sh) | SessionStart bootstrap is sh-only until M5 | run `powershell -File <plugin-root>/scripts/bootstrap.ps1` once |

## Notes for reviewers of this repo

- Engine source: `engine/` (Go, zero runtime deps; `make test` runs vet +
  race tests; `make check` verifies generated files; `wf selftest` drives
  the gates with recorded hook payloads).
- Design docs: `workflow-redesign/` (01–09) — the spec this plugin
  implements; `workflow/workflow.yaml` is the machine-readable contract.
- Agent memory (`memory: project` on design-reviewer,
  code-quality-reviewer, adversary) is self-curated recall across runs;
  **lessons** are user-approved contract changes. They complement, never
  replace, each other.
