# 02 — Distribution and Updates

How the workflow reaches a project, how it updates, and what remains in the
target repository. Requirement: "easier to update — maybe a submodule instead of
migrating all the time," with the constraint that some files must exist in the
repo root. Evidence: [01 §9](01-claude-native-capabilities.md#9-plugins-and-marketplaces).

## 1. Options analysis

### Option A — Status quo: generator + per-project emission + migrations

How v0.36 works: a Python CLI renders ~150 files into each repo; every tooling
change requires either `update` (template re-render) or a versioned migration.

| Pros | Cons |
|---|---|
| Fully self-contained repos (air-gap friendly) | **The problem being solved**: every fix = re-emission/migration in every repo; 36 versions of migration code accumulated |
| Per-project customization trivial (edit the emitted file) | Emitted copies drift; editable-region protocol is fragile; doctor needed to detect drift |
| No runtime dependency on a plugin system | Huge generator surface; prose and logic version-skew across projects |

### Option B — Git submodule mount

The tool lives in one repo, mounted at e.g. `tools/workflow/`; the project's
`.claude/settings.json` points hooks at `${CLAUDE_PROJECT_DIR}/tools/workflow/…`.

| Pros | Cons |
|---|---|
| One source of truth; `git submodule update` = upgrade | **Agents and skills are not discovered from arbitrary paths** — they must live in `.claude/agents/`, `~/.claude/agents/`, or a plugin. A submodule needs symlinks or a copy step into `.claude/` (an install script — the generator again, smaller) |
| Works fully offline / air-gapped | Submodule UX friction (init/update, detached heads, contributors forgetting `--recurse-submodules`) |
| Version pinned per project by SHA (good for reproducibility) | Root files (CLAUDE.md, settings.json) still must be materialized in the repo — partial generation remains |
| | No update notification/channel machinery — you build it |

### Option C — Claude Code plugin via marketplace (repo-as-marketplace)

The tool is a plugin; its own repo doubles as the marketplace
(`.claude-plugin/marketplace.json` at root). Projects opt in via two settings
keys; Claude Code handles fetch, cache, versioning, updates.

| Pros | Cons |
|---|---|
| **Native discovery** of skills, agents, hooks, MCP, `bin/` — zero copy step | Plugins are cache-copied; no `../` references (corpora must be bundled — already decided) |
| **Update = version bump** (or just a commit, SHA-versioned); users get it via auto-update or `/plugin marketplace update`. `renames` map handles restructures. **No migrations ever** | Requires network on first install (mitigations: `--plugin-dir`, seed dir) |
| Team adoption = 2 keys in committed `.claude/settings.json` (`extraKnownMarketplaces` + `enabledPlugins`); collaborators prompted on folder trust | `${CLAUDE_PLUGIN_ROOT}` changes per update — anything persistent must use `${CLAUDE_PLUGIN_DATA}` or the repo (a discipline, enforced by design in 08) |
| Release channels (stable/latest refs), SHA pinning, private-repo support, enterprise allowlisting — all built in | Plugin skills are namespaced (`/wf:dev` not `/dev`) — cosmetic |
| `${CLAUDE_PLUGIN_DATA}` gives a blessed persistent, update-surviving data dir | Hook/agent trust rides Claude's plugin trust model (acceptable: same trust as any settings hook) |

### Option D — Hybrid: plugin layout, multiple load paths

Design the artifact *as a plugin directory* and load it three ways:
1. **Marketplace install** (primary — full update machinery);
2. **`--plugin-dir ./tools/workflow`** on a submodule checkout (air-gapped /
   pinned-by-SHA projects; same layout, no marketplace);
3. **`CLAUDE_CODE_PLUGIN_SEED_DIR`** (CI/containers; read-only pre-populated
   cache).

Cost: nothing — a plugin directory is loadable by all three mechanisms
unchanged. The only per-mode difference is who triggers updates.

## 2. Recommendation

**Option D: build one plugin; distribute via marketplace as the primary
channel; document `--plugin-dir`-on-submodule and seed-dir as first-class
fallbacks.**

Rationale:
- C beats A and B on the stated requirement (update pain) by a wide margin —
  it deletes the migration system, the emitter, drift detection, and the
  editable-region protocol in one move, because *there are no per-project
  copies of logic or prose anymore*.
- B's only real advantage over C (offline/pinning) is fully recovered by D's
  fallback modes at zero design cost, because `--plugin-dir` accepts exactly
  the plugin layout. A submodule checkout of the plugin repo + one line in a
  wrapper script = the air-gapped story.
- The root-files concern that motivated the submodule question dissolves: the
  repo keeps only *state and instructions that are genuinely per-project*
  (§4) — everything with update churn lives in the plugin.

## 3. The plugin (`claude-workflow`, working name `wf`)

```
claude-workflow/                       ← the tool's repo = the marketplace
├── .claude-plugin/
│   ├── plugin.json                    ← name: "wf", version: X.Y.Z (explicit; see §5)
│   └── marketplace.json               ← lists "./" as the single plugin
├── skills/                            ← entry + per-phase procedures (06 §4)
│   ├── dev/SKILL.md                   ← /wf:dev — session entry
│   ├── init/SKILL.md                  ← /wf:init — project adoption (disable-model-invocation)
│   ├── frame/SKILL.md … ship/SKILL.md ← phase procedures
│   └── park/SKILL.md, force/SKILL.md  ← audited escapes (disable-model-invocation)
├── agents/                            ← roster v2 (06 §2) — no hooks frontmatter (plugin restriction)
├── hooks/hooks.json                   ← the enforcement wiring (04)
├── workflow/                          ← the declarative workflow spec (03 §4.0)
│   ├── workflow.yaml                  ← phases, families, contract items, roster
│   └── schemas/*.json                 ← record-kind JSON schemas (08 §3)
├── bin/                               ← per-platform `wf` engine binaries (07 §4)
├── reference/
│   ├── design/                        ← bundled snapshot of ../Design
│   ├── coding/                        ← bundled snapshot of ../Coding (+ languages/, checklists/)
│   └── ux/                            ← optional snapshot of UI_UX
├── templates/                         ← document templates (read-only sources)
└── scripts/sync-corpora.sh            ← maintainer-side: refresh reference/ from source repos
```

All hook entries use **exec form** referencing the engine by **absolute
placeholder path** (`${CLAUDE_PLUGIN_DATA}/bin/wf` after the SessionStart
bootstrap, 07 §4) — never bare `wf` on PATH: plugin `bin/` is documented for
the Bash tool's PATH only, and Windows exec form can only spawn real
executables (01 §9/§1). A `SessionStart` check-and-install step selects the
bundled platform binary (or fetches it, checksum-verified) into
`${CLAUDE_PLUGIN_DATA}/bin/` — the documented dependency pattern.

## 4. Repo footprint after adoption

`/wf:init` (a user-invoked, `disable-model-invocation` skill) materializes the
*minimum* per-project surface:

| Path | Content | Committed? |
|---|---|---|
| `.claude/settings.json` | `extraKnownMarketplaces` + `enabledPlugins` (+ optionally `permissions` hardening) — merged, not overwritten | ✅ |
| `CLAUDE.md` (or `.claude/CLAUDE.md`) | ≤40 lines: project facts + the workflow ground rules block (05 §4) | ✅ |
| `.workflow/` | Engine state (08): run records, checklist state, evidence ledger, config | mixed (08 §5) |
| `docs/` | Deliverables (PROJECT.md, RTM, ADRs, reviews…) created on demand by runs | ✅ |

Nothing else. No hooks in the repo, no agents, no step files, no templates —
those are all plugin-resident and update with it. `/wf:init` is idempotent and
version-aware (it records the plugin version in `.workflow/config.json`; the
engine warns on major-version skew — see §6).

## 5. Versioning and update flow

- **Explicit semver in `plugin.json`** (not SHA-versioning): the engine, hooks,
  skills, and agents ship as one atomically-versioned unit; users update on
  releases, not on every commit. A `latest` marketplace ref can be added later
  for early adopters (two-marketplace channel pattern, verified).
- Update paths: background auto-update (with `GITHUB_TOKEN` for private
  hosting) or `/plugin marketplace update` + `/plugin update wf`. Mid-session
  updates require `/reload-plugins` (hooks/servers switch to the new
  `${CLAUDE_PLUGIN_ROOT}`).
- **State compatibility instead of migrations**: the engine reads
  `.workflow/` state tagged with a `schema` field; it must read schema N and
  N-1 (one-release grace) and upgrade state lazily on first write. Because no
  prose/templates live in the repo, "migration" reduces to this single,
  testable state-schema concern.
- `renames` map in marketplace.json reserved for any future plugin split
  (e.g. extracting corpora into companion plugins).

## 6. Adoption and coexistence

- **New project**: trust folder → accept marketplace prompt → `/wf:init` →
  first run.
- **Existing v0.36 scaffold**: `/wf:init` detects `.workflow/manifest.json`,
  offers to (a) import durable state (decisions history, lessons,
  `trace/runs/`, RTM/docs stay in place — they're plain docs), (b) archive the
  old `.workflow/hooks|steps|agents` trees, and (c) strip the old hook wiring
  from `.claude/settings.json`. Old committed audit artifacts remain readable;
  the new engine never parses them (fresh ledger).
- **Air-gapped**: `git submodule add <tool-repo> tools/workflow` + launch via
  `claude --plugin-dir ./tools/workflow` (alias in a project script). Same
  plugin, no marketplace.
- **CI**: seed directory baked into the image (`CLAUDE_CODE_PLUGIN_SEED_DIR`),
  read-only, verified pattern.

## 7. Risks and mitigations

| Risk | Mitigation |
|---|---|
| Plugin cache path changes each update while a session is open | Engine resolves its own state via `${CLAUDE_PROJECT_DIR}/.workflow` + `${CLAUDE_PLUGIN_DATA}`; never stores `${CLAUDE_PLUGIN_ROOT}` paths in state; `/reload-plugins` documented in the update note |
| User hand-edits `.claude/settings.json` and breaks wiring | All enforcement wiring lives in plugin `hooks/hooks.json` (not project settings) — the project file only enables the plugin; `wf doctor` verifies |
| Marketplace unreachable on first install | `--plugin-dir` fallback; seed dir; plugin also loadable from a `.zip` via `--plugin-url` |
| Enterprise policy blocks plugins (`allowManagedHooksOnly`, `strictKnownMarketplaces`) | Documented requirement; org can force-enable via managed `enabledPlugins` (exempt from allowManagedHooksOnly per docs) |
| Trust prompt fatigue → user declines marketplace | `/wf:init` re-checks and prints exact remediation |
