# 06 — Agents, Skills, and Reference Corpora

Roster v2 (re-justified from the v0.36 roster of 28+6), the verdict protocol,
the skills catalog, and the bundled corpora. Constraints: plugin agents can't
carry frontmatter hooks (01 §6) — all gating lives in plugin `hooks/hooks.json`
via SubagentStart/SubagentStop matchers; subagents are background-by-default
(rosters parallelize).

## 1. Roster derivation principles

- **Merge agents that differ only by injected scope.** Static per-flavor files
  were the v0.36 duplication engine (7 perspective reviewers, 3 adversarial
  analysts, 2 design reviewers…). `SubagentStart` injection + the task prompt
  now parameterize one agent per *competence*, not per *occasion*.
- **Every gating agent must exist** and is declared once in the workflow
  spec's roster (03 §4.0); the hooks matchers, contract requirements, and the
  tables below are generated/validated views of that one list (kills A1/D1/D5
  by construction).
- Read-only enforced by `tools:` frontmatter; only `implementer` writes.
- Explore/Plan are **native built-ins** — no shipped `explore` agent (removes
  the D6 shadowing).

## 2. Roster v2 (15 agents; 11 gating)

### Author-side (not gated)

| Agent | Model | Tools | Phase | Mandate / corpus routing |
|---|---|---|---|---|
| `researcher` | inherit | WebSearch, WebFetch, Read, Grep | Context, Design | Sourced external research; Design mode: 2–4 genuine candidates per stage with sources; honors carried `rejected` option IDs |
| `designer` | opus | Read, Grep, Glob | Design (stage 1 system, stage 2 software — one agent, staged via injected scope) | Enumerates/validates/selects; corpus: `design/06` (decide what matters) → `01`/`02` (principles/patterns) → `03` (software design); cites file+section |
| `ux-designer` ◇ | opus | Read, Grep, Glob | Design stage 3 (`ux: true` projects) | UI/interaction candidates → `docs/design/ux-<slug>.md`; corpus `ux/00-index` → deep dives |
| `implementer` | inherit | Read, Grep, Glob, **Edit, Write, Bash** | Build | Executes tasks test-first under the task gates; corpus: `coding/` root files + `languages/<ext>` per the extension routing table; cites rule IDs (`GEN-SEC-03`, `PY-NAME-01`) |

### Gating (verdict block enforced at SubagentStop)

| Agent | Phase(s) | Mandate | Absorbs (v0.36) |
|---|---|---|---|
| `critic` | Design, Plan | Independent go/no-go: Safe/Risky/Unsafe + concerns; risky ⇒ dispositions | critic |
| `design-reviewer` | Design (per stage, fixed-to-clean) | Reviews the *selected* option against `design/` corpus (arch principles, patterns, quality attributes); never re-approves rejected options | system-design-reviewer + software-design-reviewer |
| `ux-design-reviewer` ◇ | Design stage 3 | `docs/design/ux-*.md` vs `ux/21-agent-checklists`, `05-accessibility`; a11y criticals unwaivable | ux-design-reviewer |
| `code-quality-reviewer` | Build roster | Quality + error handling + **concurrency + performance** (per-area findings; `n/a`-with-note for absent surfaces); corpus `coding/01,03,06,07,09` + `checklists/<lang>` | code-quality + code-concurrency + code-performance |
| `code-security-reviewer` | Build roster | `coding/04-security` (`GEN-SEC`), OWASP; leaked secret always critical | code-security-reviewer |
| `code-testing-reviewer` | Build roster | `coding/05-testing` (`GEN-TEST`); asserts the red→green records match the diff's claims | code-testing-reviewer |
| `design-conformance-reviewer` | Build roster + Verify confirmation | Diff implements the approved design/ADR (or standing architecture for `refactor` intent, inferred-with-reduced-confidence) | design-conformance-reviewer |
| `adversary` | Frame (abuse-case mode), Design (attack-tree mode), Build roster + Verify confirmation (red-team mode) | Break-it specialist; mode injected at SubagentStart; findings carry break recipes; exploitable authz/injection/secret always critical | adversary + abuse-case-analyst + attack-tree-analyst + the phantom `adversarial-reviewer` |
| `lens-reviewer` | Frame (ambiguities), Verify (diff/report pass) | One agent, lens injected (user/maintainer/security/reliability/compliance/stakeholder/operator); output contract: 1–3 tagged findings or explicit reasoned none | the 7 `<perspective>-reviewer` agents |
| `ux-reviewer` ◇ | Frame (usability), Build roster, Verify | WCAG 2.2 AA + interaction vs approved ux design; `n/a` for no-UI diffs; corpus `ux/05,06,18,21` | ux-reviewer |
| `auditor` | Ship | Reviews the engine-generated trace report + resolutions; HIGH findings block close | auditor |

◇ = emitted/required only when project config `ux: true`.

Cut without replacement: the 6 regulated `<standard>-reviewer`s (out of core
scope, 03 §8 — a future companion plugin can add a gating reviewer by
extending the spec roster + a contract pack, 03 §4.0), `explore` (native),
standalone `adversarial-reviewer` (was a phantom; competence lives in
`adversary`).

### Frontmatter conventions

`model:` aliases only (`opus`/`inherit` — never dated IDs; the v0.35 lesson);
`maxTurns` on every gating agent; `memory: project` on `design-reviewer`,
`code-quality-reviewer`, `adversary` (accumulate project conventions/attack
surface across runs — complements, not replaces, the lessons system: lessons
are *user-approved contract changes*; agent memory is *self-curated recall*);
`skills:` preloads the agent's corpus-routing skill where one exists.

## 3. Verdict protocol v2

- The fenced block (single-sourced in the workflow spec, 03 §4.0 / 07 §4):

  ```verdict
  status: <clean|changes-required|safe|risky|unsafe|n/a>
  criticals: <int>
  majors: <int>
  ```

  `n/a` is first-class (fixes A2/A3). Pass set: `clean|safe|n/a`. `risky`
  passes only with recorded dispositions; `clean|safe` requires
  criticals=0 ∧ majors=0 (auto-downgrade otherwise).
- **Emission is enforced, not hoped for**: the SubagentStop gate blocks the
  reviewer until it emits a parseable block (2 attempts, then recorded
  `unparsed` = fail) — see 04 §4. Auto-recorded with `auto: true`; sticky
  against manual override.
- Severity vocabulary unified to `minor|major|critical` everywhere (drops the
  4-level regulated scale — D4).

## 4. Skills catalog

| Skill | Invocation | Content |
|---|---|---|
| `wf:dev` | `/wf:dev` (+ auto) | Session entry: read injected status → resume/start/branch decision → invoke the active phase skill. ≤40 lines |
| `wf:init` | user-only (`disable-model-invocation`) | Project adoption/upgrade (02 §4/§6) |
| `wf:frame` … `wf:ship` | phase skills (7) | The phase procedure: contract-first layout (05 §3), record commands inline, family-specific sections (one file per phase; family branching is small enough inline once routes collapsed to 3) |
| `wf:park`, `wf:force` | user-only | The audited escapes — model cannot self-invoke |
| `wf:status` | `/wf:status` | Pretty-print `wf status` + report signals |

Phase skills are gated by PreToolUse on `Skill(wf:*)` (04 §5): only the
active phase's skill (or a legal loop target) can be invoked — sequencing is
enforced at the invocation, not by policing file reads.

## 5. Bundled corpora

Per the distribution decision, snapshots live in the plugin:

```
reference/
├── design/    ← ../Design  (00-index … 09-references; ~10 files)
├── coding/    ← ../Coding  (00-index, 01–12, languages/{python,typescript,vue,csharp,embedded-c}.md, checklists/, references.md)
└── ux/        ← UI_UX      (optional; 00-index … 21-agent-checklists)
```

- **Sync**: `scripts/sync-corpora.sh` (maintainer-side) copies from the source
  repos, stamps `reference/<name>/VERSION` (source SHA + date). Corpora
  update with plugin releases — always present, never a submodule to forget
  (and plugin cache-copying forbids external references anyway).
- **Routing**: both corpora are built for this ("load only what you need"):
  - `coding/00-index.md` defines the **file-extension routing table**
    (`.py → languages/python.md + checklists/python.md`, etc.), stable rule
    IDs, and the priority rule "language-specific wins." `implementer` and
    the code reviewers follow it; the final pre-delivery pass uses
    `checklists/`.
  - `design/00-index.md` defines **reading paths** ("designing a new system →
    06 then 01/02 then 03…") and the decision quick-reference matrix;
    `designer`/`design-reviewer` follow those paths; ADRs use the `08`
    template.
  - Routing lives in each agent's prompt as a small table, and
    `wf inject agent` names the concrete files for the current scope, so a
    reviewer never greps the whole corpus.
- Agents state the fallback explicitly: corpus absent/unreadable ⇒ use own
  knowledge and note it in the verdict.

## 6. Document templates

The v0.36 clean-template discipline survives with fewer moving parts:
templates live in the plugin (`templates/`), the engine copies them
(`wf doc new adr|design|threat-model|ux|review|incident|release-notes|
delivery-manifest …`) into the correct `docs/…` destination and records the
artifact — the copy step is engine-mediated, so "authored the template in
place" is no longer possible. RTM and the run narrative are *generated* views
over records, not maintained documents.
