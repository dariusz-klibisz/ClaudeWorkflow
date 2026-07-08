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
  Build → Verify → Ship. Each phase has an exit contract (recorded facts,
  reviewer verdicts, user approvals) and an entry contract (the previous
  phases' canonical outputs must exist before the transition — force/adopt
  landings get re-checked); document deliverables are engine-created from
  templates (`wf doc new`) and verified AUTHORED ON DISK, not just claimed.
- **Four gates**, wired as Claude Code hooks: the Stop gate (can't end the
  turn with unmet obligations), task gates (can't complete a task without
  captured red→green test evidence), the verdict gate (reviewer subagents
  must end with a parseable verdict, auto-captured), and tool gates
  (phase-skill sequencing, an always-on catastrophic-Bash net, and
  ledger protection: `.workflow/{log,state,runs}` and `config.json` are
  engine-written only — tool writes are denied with no override, and the
  event log is hash-chained so out-of-band edits surface in `wf doctor`).
- **Grounded evidence.** Test runs are captured from the Bash tool by the
  hook itself (`auto:true`) — recognized from a built-in runner list, the
  run's own recorded verification commands (any language), or project
  config `"runners"`. Red→green pairs are runner-matched: a cross-runner
  "pair" (a gitleaks red before a pytest green) never satisfies test-first,
  and same-runner pairs with diverging selectors pass but are surfaced as
  weakly-paired in `wf report`. Manual records stay possible but are marked
  self-attested and surface in `wf report` — and once verdict capture has
  proven itself alive in a run, a hand-recorded gating-reviewer verdict no
  longer satisfies the contract (re-run the reviewer, or disposition it).
- **Content floors, not just checkboxes.** Requirements must carry ≥1 AC
  (write-time refusal — an AC-less requirement used to dodge every per-AC
  gate), context maps and negative-space walks have waivable ≥3-element
  depth floors, option-sets need ≥2 genuine candidates and the engine
  refuses re-selecting a previously rejected option, and assessment runs
  record structured `finding`s that must each appear verbatim in the
  on-disk report.
- **Audited escapes, not hidden ones.** `/wf:park` (honest stop),
  `/wf:force` (bypass one gate; escalates — the 3rd force auto-parks the
  run). Everything is recorded and reported. Loops run Verify→{Build,
  Design, Plan} on failing ACs — and Ship-stage discoveries (a failing
  audit, an open trace finding) loop back to Verify (`--cause audit`)
  instead of being dispositioned around.
- **Anchored approvals — three dials.** Approvals are self-attested by
  design (no hook proves a human typed them), but answers given through
  AskUserQuestion are hook-captured and linked to the approval
  (`answer_ref`) — harder to fabricate, still not proof, and reported as
  such. Scope/design/plan approvals additionally **bind engine-computed
  refs** (the requirement rids, selected options, task ids in scope at
  approval time): records added afterwards without re-approval surface as
  approval-drift findings at Ship. Opt-in strictness in `.workflow/config.json` `flags`:
  `"approvals": "hardened"` refuses un-anchored approvals outright;
  `"approvals": "challenge"` goes further — each approval requires a
  **single-use code the engine shows only in your statusline** (never to
  the model; tool gates deny reads of its storage), which you type into
  the AskUserQuestion answer. The model cannot manufacture its own anchor
  because it never sees the code before you do.
- **Lessons that bite — and are measured.** At Ship, the run's lessons are
  proposed and user-triaged (`wf lessons accept|reject`); accepted lessons
  with a `check:` become ordinary contract items in
  `.workflow/contracts.d/` — **enforced from the next run on** by the same
  evaluator as everything else. Prose lessons regenerate
  `.claude/rules/wf-lessons.md`. `wf lessons status` is the efficacy view:
  per accepted lesson, whether its item was waived (dodged) or its trigger
  recurred in later runs (not working).
- **Generated traceability.** `wf trace --rtm [--write]` renders the
  requirements-traceability matrix (requirement → AC → verification →
  grounded evidence → verdict → tasks → loops) straight from the ledger —
  the auditor-facing view of the chain the gates enforce, never
  hand-authored.
- **Compliance packs (opt-in).** `packs/` ships add-only contract packs:
  `sbom` and six regulated standards (ISO 26262, IEC 62304, DO-178C,
  IEC 61508, EN 50128, NIST 800-53) that arm a per-standard
  compliance-reviewer at Design/Verify plus an evidence package at Ship.
  **Not a compliance tool** — see `packs/README.md`.

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
wf statusline                  one-line statusLine payload (run · phase · unmet);
                               /wf:init wires it into .claude/settings.json
                               unless you already have a statusLine
wf report [--run <id|current>] health signals: loops (by cause and per AC),
          [--worktrees]        escapes, self-attested counts, ungrounded ACs,
                               lesson counters; --worktrees groups across trees
wf trace [--rtm [--write]]     ship close-out findings; --rtm renders the
                               requirements-traceability matrix (--write emits
                               docs/requirements/RTM-<run>.md)
wf lessons suggest|accept|reject|apply|status
                               status = efficacy: dodged items, recurring triggers
wf doc new <type> --slug …     17 engine-mediated document templates (ADR,
                               design, threat-model, abuse-cases, attack-tree,
                               test-plan, runbook, retro, evidence-package,
                               findings reports…); status=present is refused
                               until the file is authored on disk (never a stub)
wf pack install <dir-or-yaml>  add-only contract packs (validated before merge);
                               official packs ship under packs/ (sbom + 6
                               regulated standards)
wf doctor [--bootstrap]        state health, ledger hash-chain verification,
                               corpus snapshot age · verifies AND heals the
                               hook engine
wf selftest                    39 in-scaffold enforcement scenarios
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
| Native Windows (no sh) | SessionStart bootstrap is sh-only | one-time: `powershell -File <plugin-root>/scripts/bootstrap.ps1` (or `wf doctor --bootstrap`); afterwards the engine self-updates on version skew |
| Session block shows "engine self-updated" | plugin was updated; the engine re-bootstrapped itself | nothing — hooks run the new engine from the next event on |

## Releasing (maintainers)

Binaries never live in git — they are attached to GitHub Releases, and the
bootstrap **fetches the platform binary on first use, verified against the
committed `bin/MANIFEST`** (the trust anchor; 07 §4-B). Release builds are
reproducible (`-trimpath`, plain semver stamp), so CI can verify that a
rebuild at the tag matches the committed checksums exactly:

```sh
# 1. bump "version" in .claude-plugin/plugin.json
make dist RELEASE=1 manifest   # 2. build all 6 platforms + write bin/MANIFEST
git commit -am "release: vX.Y.Z"
git tag vX.Y.Z && git push --follow-tags
# 3. release.yml: tests + selftest + rebuild + verify-manifest + gh release
```

If `verify-manifest` fails in CI, the build is not reproducing (toolchain
skew or a stale manifest) — nothing is published.

## Notes for reviewers of this repo

- Engine source: `engine/` (Go, zero runtime deps; `make test` runs vet +
  race tests; `make check` verifies generated files; `wf selftest` drives
  the gates with recorded hook payloads).
- Design docs: `workflow-redesign/` (01–09) — the spec this plugin
  implements; `workflow/workflow.yaml` is the machine-readable contract.
- Bundled corpora (`reference/{design,coding,ux}`) are versioned snapshots:
  `VERSION` names the source remote + sha, `SHA256SUMS` is verified by
  `make check` (hand-edits fail CI), and the scheduled `corpora` workflow
  fails when a source repo drifts from its snapshot (release-blocking in
  `release.yml`). Refresh with `scripts/sync-corpora.sh`.
- Agent memory (`memory: project` on design-reviewer,
  code-quality-reviewer, adversary) is self-curated recall across runs;
  **lessons** are user-approved contract changes. They complement, never
  replace, each other.
