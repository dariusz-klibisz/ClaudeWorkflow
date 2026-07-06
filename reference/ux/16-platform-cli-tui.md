# Platform: CLI & TUI

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers user experience for command-line interfaces and terminal UIs: modern human-first CLI philosophy (clig.dev synthesis), command and flag design, help systems, output design, progress indication, errors and exit codes, interactivity rules, configuration precedence, safety, composability, tab completion, full-screen TUI patterns, and performance. General principles (feedback, error prevention, recognition/recall) apply throughout — see [01-core-principles.md](01-core-principles.md) and [02-cognitive-foundations.md](02-cognitive-foundations.md); desktop-terminal integration overlaps [14-platform-desktop.md](14-platform-desktop.md#linux-conventions).

## Contents

- [Philosophy](#philosophy)
  - [Human-first CLI design](#human-first-cli-design)
- [Command surface](#command-surface)
  - [Command structure and naming](#command-structure-and-naming)
  - [Arguments vs flags](#arguments-vs-flags)
  - [Standard flag vocabulary](#standard-flag-vocabulary)
- [Help and guidance](#help-and-guidance)
  - [Help system design](#help-system-design)
  - [Suggestions and did-you-mean](#suggestions-and-did-you-mean)
  - [Tab completion](#tab-completion)
- [Output](#output)
  - [Output design: stdout, stderr, and humans](#output-design-stdout-stderr-and-humans)
  - [Color and formatting](#color-and-formatting)
  - [Machine-readable output](#machine-readable-output)
  - [Progress indication](#progress-indication)
- [Failure](#failure)
  - [Error messages](#error-messages)
  - [Exit codes](#exit-codes)
- [Behavior](#behavior)
  - [Interactivity rules](#interactivity-rules)
  - [Configuration precedence](#configuration-precedence)
  - [Safety: dry-run, idempotency, destructive guards](#safety-dry-run-idempotency-destructive-guards)
  - [Composability](#composability)
  - [Performance and responsiveness](#performance-and-responsiveness)
  - [Versioning and deprecation communication](#versioning-and-deprecation-communication)
- [Full-screen TUIs](#full-screen-tuis)
  - [When a TUI beats a CLI](#when-a-tui-beats-a-cli)
  - [TUI interaction patterns](#tui-interaction-patterns)
- [Quick reference](#quick-reference)

---

## Philosophy

### Human-first CLI design

**Definition:** The modern doctrine (codified by clig.dev, drawing on decades of UNIX practice): CLIs are *human* interfaces first and scripts' interfaces second — conversational, discoverable, forgiving — while preserving UNIX composability (text streams, pipes, small tools) when a machine is on the other end.

**Reasoning / Evidence:** Classic UNIX philosophy (McIlroy: write programs that do one thing well, work together, handle text streams) optimized for machine composition; the modern correction observes that most CLI invocations are humans exploring — and tools like git succeeded *despite* early UX hostility, spawning a generation of gentler tools (gh, cargo, deno) proving discoverability and composability coexist. TTY detection lets one tool serve both audiences dynamically.

**When to use:**
- Every developer/ops tool: assume a stressed human at 2am seeing your tool for the first time — every message should help them succeed.
- Dual-audience discipline: human-friendly by default when stdout is a TTY; stable, parseable when piped ([#machine-readable-output](#machine-readable-output)).

**When NOT to use / exceptions:**
- Pure plumbing tools (build-internal utilities) may reasonably be machine-only — document that explicitly.

**Pros:** Adoption (CLIs live or die by word-of-mouth developer love); reduced support burden; scriptability preserved.
**Cons:** Dual-mode output paths are more code; "friendly" must never compromise output stability contracts.

**Implementation guidance:**
- The empathy checklist per command: does it explain itself (`--help`)? Does failure say what to do next? Does success confirm what happened? Can it be undone/dry-run?
- Design commands as a *conversation*: each output suggests the next step ("Created app 'atlas'. Run `tool deploy` to publish.").

**Related:** [#help-system-design](#help-system-design), [#error-messages](#error-messages), [01-core-principles.md](01-core-principles.md#1-visibility-of-system-status).

**Sources:** [clig.dev](https://clig.dev/), [The Art of Unix Programming (Raymond)](http://www.catb.org/~esr/writings/taoup/).

---

## Command surface

### Command structure and naming

**Definition:** How multi-command tools organize verbs and nouns: `tool <noun> <verb>` (e.g., `docker container ls`, `gh pr create`) vs `tool <verb>` flat structure (e.g., `git commit`) — with noun-verb subcommand grouping as the modern default for tools managing multiple resource types.

**Reasoning / Evidence:** Noun-first grouping scales: related commands cluster under discoverable namespaces (`tool db --help` lists all db operations — recognition over recall [01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall)); flat verb sets work only while the verb count stays small. Consistency beats either choice: mixed schemes (`tool create-db` + `tool user delete`) are the real failure.

**When to use:**
- Noun-verb (`tool pr create`, `tool pr list`): tools with ≥2 resource types; verbs standardized across nouns (list/get/create/update/delete — same verbs everywhere).
- Flat verbs (`tool build`, `tool test`): single-domain tools with <~10 commands.

**When NOT to use / exceptions:**
- Don't nest deeper than 2 levels (`tool a b c d` is a labyrinth); don't create both `tool db create` and `tool create db`.
- Command names: lowercase, no underscores, prefer full words with unambiguous shortest-prefix; the binary name itself short (2–8 chars), memorable, collision-checked against common tools.

**Pros (noun-verb):** Discoverability by namespace, predictable growth, per-noun help.
**Cons:** Slightly longer invocations (mitigate with aliases: `tool ls` → `tool container list`).

**Implementation guidance:**
- Verb vocabulary fixed product-wide: `list` (not ls/show/enumerate mixed), `get`, `create`, `delete` (not remove/rm/destroy mixed — pick one, alias the rest).
- Aliases for frequency: document them in help; muscle-memory compatibility aliases when renaming ([#versioning-and-deprecation-communication](#versioning-and-deprecation-communication)).
- Default command: running bare `tool` shows help or a status dashboard — never nothing, never an error.

**Related:** [#standard-flag-vocabulary](#standard-flag-vocabulary), [#help-system-design](#help-system-design).

**Sources:** [clig.dev: Subcommands](https://clig.dev/#subcommands), [gh CLI design](https://primer.style/cli/getting-started/principles).

### Arguments vs flags

**Definition:** The choice between positional arguments (`cp SRC DST`) and named flags (`--output file`): positionals for the 1–2 obvious primary objects; flags for everything else — clig.dev's "prefer flags to args" lean, because named parameters self-document.

**Reasoning / Evidence:** Positionals are terse but opaque (what's argument 3 of `ln`? — perennial confusion); flags read as documentation in scripts and shell history (`--force` explains itself in a code review; a bare `-f 3 true` doesn't). Order-dependence of multiple positionals is a chronic slip source ([02-cognitive-foundations.md](02-cognitive-foundations.md#slips-vs-mistakes) — `cp` direction errors).

**When to use:**
- Positional: the single obvious object (`tool deploy <app>`, `cat <file>`), or a natural homogeneous list (`rm file1 file2`).
- Flags: everything with a value, anything optional, anything whose meaning isn't obvious from position; required flags acceptable when explicitness matters more than brevity.

**When NOT to use / exceptions:**
- Never >2 heterogeneous positionals; never positionals whose order users must memorize between similar commands.
- Support `--` separator (end of flags) and `-` as stdin/stdout placeholder per convention.

**Pros (flags):** Self-documentation, order-independence, safe extension (new flags don't break old calls).
**Cons:** Verbosity (mitigate: short aliases for the frequent few).

**Implementation guidance:**
- Every flag: long form always (`--output`); short form (`-o`) only for frequently-typed ones.
- Flag values: both `--flag value` and `--flag=value` accepted; repeated flags for lists (`-v -v` or `--tag a --tag b`).
- Validate combinations early with specific errors ("--json and --quiet are mutually exclusive").

**Related:** [#standard-flag-vocabulary](#standard-flag-vocabulary), [01-core-principles.md](01-core-principles.md#postels-law).

**Sources:** [clig.dev: Arguments and flags](https://clig.dev/#arguments-and-flags).

### Standard flag vocabulary

**Definition:** The conventional flag names users expect everywhere: `-h/--help`, `--version`, `-v/--verbose` (repeatable), `-q/--quiet`, `-o/--output <path>`, `-f/--force`, `--json` (or `--format`), `--dry-run`, `--yes/-y` (assume-yes), `--no-color`, `--config <path>`, `--all/-a`, `--recursive/-r`.

**Reasoning / Evidence:** [Jakob's Law](01-core-principles.md#jakobs-law) for terminals: decades of getopt convention mean users *guess* these flags before reading help; honoring the guesses is free usability. Conflicts hurt: `-v` as "version" (some tools) vs "verbose" (most) is a classic trap — prefer `--version` long-only.

**When to use:** Every CLI: implement the standard set with standard meanings before inventing anything; reserve `-h` for help absolutely.

**When NOT to use / exceptions:**
- `-f` ambiguity (force vs file): pick per domain norms, document loudly.
- Never repurpose a standard name for a different meaning (a `--force` that means "refresh cache" will burn someone).

**Pros:** Guessability, script portability, muscle-memory transfer.
**Cons:** None; deviation costs are pure loss.

**Implementation guidance:**
- `--help` on every command and subcommand; bare `-h` never destructive.
- `--version`: version + minimal build info to stdout, exit 0.
- Verbosity ladder: default (result-focused) → `-v` (progress detail) → `-vv` (debug) → `-q` (errors only); `--quiet` still prints errors to stderr.

**Related:** [#help-system-design](#help-system-design), [#interactivity-rules](#interactivity-rules), [01-core-principles.md](01-core-principles.md#jakobs-law).

**Sources:** [clig.dev: standard flags](https://clig.dev/#arguments-and-flags), [GNU coding standards: option table](https://www.gnu.org/prep/standards/html_node/Option-Table.html).

---

## Help and guidance

### Help system design

**Definition:** The layered help stack: bare-invocation guidance (no args → concise help or sensible default), `-h` (compact: usage + common flags + examples), `--help` (full reference), examples-first documentation, and web docs linked — the CLI's entire discoverability surface.

**Reasoning / Evidence:** CLIs lack visible affordances; help *is* the interface map ([01-core-principles.md](01-core-principles.md#10-help-and-documentation) at maximum stakes). Usage studies of man-page-style walls vs examples-first help consistently favor examples: users pattern-match a working invocation faster than they parse grammar ([02-cognitive-foundations.md](02-cognitive-foundations.md#recognition-vs-recall)).

**When to use:**
- No args (when args are required): print concise help, exit 1 — never a bare error, never hang waiting for stdin invisibly ([#interactivity-rules](#interactivity-rules)).
- Help content order: one-line description → usage synopsis → **examples (2–5, real-world, copy-pasteable)** → common flags → less-common flags → footer (docs URL, version).
- Per-subcommand help (`tool db --help`) listing that namespace's commands.

**When NOT to use / exceptions:**
- Tools whose bare invocation has an obvious default job (`ls`, `git status`-like dashboards) run it instead of printing help.

**Pros:** Self-service adoption; support deflection; the help text doubles as your API contract documentation.
**Cons:** Help maintenance drift (generate from the flag definitions — single source of truth).

**Implementation guidance:**
- Width ≤80 columns for readability; sections bold when TTY; examples prefixed `$ ` conventionally.
- `tool help <command>` as alias for `tool <command> --help`; man pages for system-installed tools (generated).
- Suggest next steps in help footers and success output ("Next: run `tool init`").

**Related:** [#suggestions-and-did-you-mean](#suggestions-and-did-you-mean), [09-content-ux-writing.md](09-content-ux-writing.md#plain-language).

**Sources:** [clig.dev: Help](https://clig.dev/#help), [Command Line Interface Guidelines: docs](https://clig.dev/#documentation).

### Suggestions and did-you-mean

**Definition:** Recovering typos and near-misses: unknown command/flag errors that propose the closest match ("unknown command 'stauts' — did you mean 'status'?"), optionally offering interactive acceptance.

**Reasoning / Evidence:** Typos are the CLI's dominant slip class ([02-cognitive-foundations.md](02-cognitive-foundations.md#slips-vs-mistakes)); edit-distance suggestion (git's implementation is canonical) converts a dead-end error into a one-keystroke recovery. Error-with-suggestion embodies [heuristic #9](01-core-principles.md#9-help-users-recognize-diagnose-and-recover-from-errors).

**When to use:** All unknown-command/flag/argument-value errors; also near-miss values ("no environment 'prod' — did you mean 'production'?").

**When NOT to use / exceptions:**
- Never auto-*execute* the guess (destructive potential); interactive confirm at most, and only when TTY ([#interactivity-rules](#interactivity-rules)).

**Pros:** Slip recovery at near-zero cost. **Cons:** Bad suggestions on short strings (threshold the edit distance sensibly: ≤2 for words ≥4 chars).

**Implementation guidance:**
- Levenshtein ≤2 against known commands; show top 1–3 candidates; exit non-zero regardless.
- Framework support (Cobra, clap) makes this nearly free — turn it on.

**Related:** [#error-messages](#error-messages), [02-cognitive-foundations.md](02-cognitive-foundations.md#slips-vs-mistakes).

**Sources:** [clig.dev: errors](https://clig.dev/#errors), [git autocorrect](https://git-scm.com/book/en/v2/Customizing-Git-Git-Configuration).

### Tab completion

**Definition:** Shell-integrated completion of commands, flags, and — the high-value tier — *dynamic values* (branch names, app names, file-type-filtered paths) for bash/zsh/fish/PowerShell.

**Reasoning / Evidence:** Completion converts recall to recognition mid-keystroke ([02-cognitive-foundations.md](02-cognitive-foundations.md#recognition-vs-recall)); dynamic completion (gh completing PR numbers, kubectl completing pod names) eliminates whole lookup-then-retype loops.

**When to use:** Every distributable CLI: ship completion scripts (`tool completion zsh > …` convention), document one-line installation per shell; package managers (brew, apt) install them automatically.

**When NOT to use / exceptions:** Dynamic completions must be fast (<100ms) — cache or degrade to static rather than lag the shell.

**Pros:** Expert speed, typo elimination, discoverability-in-flow (tab-tab lists options).
**Cons:** Per-shell maintenance (framework-generated where possible); slow dynamic completion is worse than none.

**Implementation guidance:**
- Generate from command definitions (Cobra/clap/click do this); complete flags after `-`, values after `=`/space contextually.
- Dynamic sources with timeouts and caching; never hit the network synchronously without a subsecond bound.

**Related:** [#help-system-design](#help-system-design), [01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use).

**Sources:** [clig.dev: tab completion](https://clig.dev/#tab-completions).

---

## Output

### Output design: stdout, stderr, and humans

**Definition:** The stream contract: **stdout = the product** (data, results — what pipes consume); **stderr = the conversation** (progress, warnings, errors, logs — what humans read). Plus brevity discipline: say what happened, concisely, and shut up when all went well (with a confirmation line — silence-on-success is the old-school extreme; modern tools confirm briefly).

**Reasoning / Evidence:** The stream split is what makes `tool | jq` work while progress still displays; polluting stdout with banners/logs breaks every downstream pipe (the cardinal composability sin). Confirmation-on-success (e.g., `Created bucket s3://x`) serves [visibility of status](01-core-principles.md#1-visibility-of-system-status) without breaking pipes when it goes to stderr — or stays on stdout only when it *is* the result.

**When to use:**
- stdout: query results, generated content, IDs of created resources (usable in scripts: `ID=$(tool create …)`).
- stderr: everything else — progress bars, "Fetching…", warnings, deprecation notices, errors.
- Success output: one line stating the outcome + (when apt) the next step; `-q` silences it.

**When NOT to use / exceptions:**
- ASCII-art banners, version headers, and tips on stdout: never (stderr, and only when TTY, and sparingly).
- Timestamps/log-prefixes on human output: don't (that's `-vv` debug mode's job).

**Pros:** Pipe-safety + human clarity simultaneously.
**Cons:** Requires discipline in every print statement (lint it: no logging to stdout).

**Implementation guidance:**
- TTY-detect both streams independently (stdout piped + stderr TTY = show progress on stderr, clean data on stdout).
- Line-buffered stdout when piped (streaming consumers); flush before exit.
- Tables for TTY humans (aligned columns, headers); TSV when piped ([#machine-readable-output](#machine-readable-output)).

**Related:** [#machine-readable-output](#machine-readable-output), [#color-and-formatting](#color-and-formatting), [#composability](#composability).

**Sources:** [clig.dev: Output](https://clig.dev/#output).

### Color and formatting

**Definition:** Terminal color/style used semantically — red errors, yellow warnings, green success, bold emphasis, dim de-emphasis — with mandatory escape hatches: auto-disable when not a TTY, `NO_COLOR` env respect, `--no-color` flag, and never color as the *only* signal.

**Reasoning / Evidence:** Color aids scanning ([03-visual-design.md](03-visual-design.md#semantic-color) semantics transfer); but piped output with ANSI codes corrupts files/parsers, ~8% of men of Northern European descent (~0.5% of women) are color-blind ([05-accessibility.md](05-accessibility.md#color-independence)), and terminal themes vary wildly (yellow-on-white invisibility). The no-color.org convention (`NO_COLOR` env var) is the ecosystem standard.

**When to use:**
- Sparingly, semantically: errors red+**"error:"** prefix, warnings yellow+"warning:" prefix (the text prefix is the color-independent channel), success green ✓, headers bold, metadata dim.
- Symbols with fallbacks: ✓/✗ where the terminal supports Unicode; ASCII fallback (`OK`/`FAIL`) detection or config.

**When NOT to use / exceptions:**
- Disable automatically: stdout not a TTY, `NO_COLOR` set, `TERM=dumb`; honor `FORCE_COLOR`/`--color=always` for the CI-wants-color case.
- Rainbow output (color-per-field) is noise; background colors rarely survive themes.

**Pros:** Scan speed for humans at zero machine cost (when gated correctly).
**Cons:** Theme unpredictability (test dark+light terminals); Windows legacy console quirks (modern Windows Terminal is fine).

**Implementation guidance:**
- Prefix discipline: `error:`/`warning:` textual markers always, color as enhancement.
- 16-color palette for portability (256/truecolor as progressive enhancement).
- Test: pipe to file (no codes), NO_COLOR=1, light theme, dark theme.

**Related:** [05-accessibility.md](05-accessibility.md#color-independence), [#output-design-stdout-stderr-and-humans](#output-design-stdout-stderr-and-humans).

**Sources:** [no-color.org](https://no-color.org/), [clig.dev: color](https://clig.dev/#output).

### Machine-readable output

**Definition:** Stable structured output for scripts: `--json` (or `--format json|yaml|tsv`) emitting well-defined schemas, plain line-oriented output when piped, and the stability contract — machine output is an API with versioning obligations.

**Reasoning / Evidence:** Modern pipelines are JSON-first (`| jq` is the new awk); tools without `--json` force fragile text-scraping that breaks on every cosmetic change. The contract insight: once scripts parse your output, its format is frozen — human output can evolve freely *only if* machine output is a separate, stable channel.

**When to use:**
- Any command producing data: `--json` from day one (retrofitting after scripts scrape human output is too late).
- Piped default: plain, line-oriented, tab-separated, no decoration — grep/cut/awk-friendly.
- Streaming data: JSON Lines (one object per line) over monolithic arrays.

**When NOT to use / exceptions:**
- Purely interactive commands (login flows) may have nothing to structure — still exit-code correctly.

**Pros:** First-class scripting citizenship; human-output freedom preserved; ecosystem integrations (CI, dashboards) unlocked.
**Cons:** Schema maintenance (document it; version it; additive changes only within a major version).

**Implementation guidance:**
- JSON schema: stable field names (snake_case convention), ISO 8601 timestamps, explicit nulls, no locale-dependent formatting; errors in JSON mode also structured (`{"error": {...}}` on stderr or per convention).
- `--format` values: json, yaml, table (explicit human), tsv; `--jq`-style built-in filtering is a power-feature differentiator (gh does this).

**Related:** [#composability](#composability), [#exit-codes](#exit-codes), [10-design-systems.md](10-design-systems.md#versioning-and-deprecation) (same contract thinking).

**Sources:** [clig.dev: machine-readable output](https://clig.dev/#output), [JSON Lines](https://jsonlines.org/).

### Progress indication

**Definition:** Feedback during long operations in a terminal: spinners (unknown duration), progress bars with counts/ETA (known), step logs (multi-stage), and the >1s rule — no operation runs silently past ~1 second ([04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds) applied to text).

**Reasoning / Evidence:** A silent hanging terminal is indistinguishable from a hung process — users Ctrl+C working operations constantly for lack of a spinner. The response-time tiers transfer directly: <1s nothing needed; 1–10s spinner + verb ("Resolving dependencies…"); >10s determinate progress + ETA when computable.

**When to use:**
- Spinner + current-action text for network/IO waits; progress bar (`[####----] 42% 12s`) for enumerable work (downloads, batch operations); stage checklists (✓ Built / ⠋ Deploying / ○ Verifying) for pipelines.
- All progress to **stderr** ([#output-design-stdout-stderr-and-humans](#output-design-stdout-stderr-and-humans)); TTY-only (piped/CI logs get periodic plain lines instead — spinner frames pollute logs).

**When NOT to use / exceptions:**
- CI environments (detect `CI` env): timestamped plain progress lines, no ANSI cursor games.
- Never fake progress; stalled operations must say so ("still waiting on lock…") rather than freeze the bar.

**Pros:** Kills the false-hang Ctrl+C; ETA manages expectations ([01-core-principles.md](01-core-principles.md#1-visibility-of-system-status)).
**Cons:** Terminal-width/redraw handling; concurrency (parallel task rendering needs a real TUI region — or simple per-task lines).

**Implementation guidance:**
- Spinner threshold: appear after ~300ms grace ([04-interaction-design.md](04-interaction-design.md#loading-strategies)); update ≤10fps (terminal redraw cost).
- On completion, replace the progress line with the outcome line (leave a clean scrollback story).
- Long operations: interruptible (Ctrl+C aborts cleanly with state message), resumable where possible ([#safety-dry-run-idempotency-destructive-guards](#safety-dry-run-idempotency-destructive-guards)).

**Related:** [04-interaction-design.md](04-interaction-design.md#loading-strategies), [#interactivity-rules](#interactivity-rules).

**Sources:** [clig.dev: responsiveness](https://clig.dev/#robustness).

---

## Failure

### Error messages

**Definition:** CLI errors that rescue: `error:` prefix, plain-language what-happened, the offending input quoted, the *fix* suggested, and a docs/debug pointer — to stderr, with a correct exit code. ([Heuristic #9](01-core-principles.md#9-help-users-recognize-diagnose-and-recover-from-errors) in its harshest environment: no UI to click, just text and hope.)

**Reasoning / Evidence:** CLI errors are read at moments of maximal frustration and often end up pasted into search engines/issue trackers — greppable, specific error text is both user rescue and your support search-index. Stack traces as default user-facing errors are diagnostic vomit (keep behind `-vv`/`--debug`).

**When to use:**
- Template: `error: couldn't connect to database 'prod-db' (connection refused). Is the server running? Try 'tool db status' or see https://docs…/errors#db-connect`.
- Categorize: user error (fix suggestion) vs environment error (what to check) vs bug (how to report, with `--debug` bundle).
- Machine mode: structured errors with stable `code` fields ([#machine-readable-output](#machine-readable-output)).

**When NOT to use / exceptions:**
- Warnings ≠ errors: warnings to stderr, execution continues, exit 0 (unless `--strict`).
- Don't bury the error under 40 lines of context; error first (or last — the visible-at-prompt position), detail above.

**Pros:** Self-service recovery; searchable support surface.
**Cons:** Per-failure-mode writing effort (worth it: each good message amortizes across every future occurrence).

**Implementation guidance:**
- Error codes (E1234 or slugs) for docs deep-links and grep; stable across versions.
- Signal-to-noise: one error, clearly, not a cascade (fail-fast or aggregate related failures into a summarized report).
- Include remediation commands ready to copy-paste where safe.

**Related:** [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing), [#exit-codes](#exit-codes), [#suggestions-and-did-you-mean](#suggestions-and-did-you-mean).

**Sources:** [clig.dev: errors](https://clig.dev/#errors).

### Exit codes

**Definition:** The numeric success contract: **0 = success, non-zero = failure**, with documented specific codes for programmatically-distinguishable failures (e.g., 1 general, 2 misuse/bad-args — a shell/bash convention ("misuse of shell builtins") widely adopted by tools for usage errors, then domain codes) — the API that `&&`, `||`, `set -e`, and every CI system consume.

**Reasoning / Evidence:** Scripts branch on exit codes; a tool exiting 0 on failure silently corrupts pipelines (the worst CLI bug class); a tool exiting non-zero on success breaks `set -e` automation. Reserved zones: 126–127 (shell: not-executable/not-found), 128+n (killed by signal n), 130 (Ctrl+C) — don't collide.

**When to use:**
- Always 0 iff the requested operation succeeded; partial failure (3 of 5 items processed) = non-zero + per-item report (unless `--continue-on-error` semantics documented otherwise).
- Distinct codes where scripts plausibly branch: not-found (curl-style code families), auth-failure, validation vs runtime — documented in `--help`/man.

**When NOT to use / exceptions:**
- Don't invent 50 codes nobody checks — 3–8 meaningful ones beat an unmaintained taxonomy; JSON error objects carry the rich detail ([#machine-readable-output](#machine-readable-output)).
- Diff-style tools (exit 1 = differences found, 2 = error) follow their genre's convention.

**Pros:** Automation correctness; the simplest API you'll ever ship.
**Cons:** Only discipline — every early-return and exception path must map deliberately.

**Implementation guidance:**
- Document codes in help; test them (CI asserts exit codes per scenario); never exit 0 from a catch-all exception handler.
- `--quiet` mode exists largely so scripts can rely purely on exit codes.

**Related:** [#machine-readable-output](#machine-readable-output), [#composability](#composability).

**Sources:** [clig.dev: exit codes](https://clig.dev/#robustness), [Exit status (Wikipedia)](https://en.wikipedia.org/wiki/Exit_status), [GNU: exit statuses](https://www.gnu.org/software/bash/manual/html_node/Exit-Status.html).

---

## Behavior

### Interactivity rules

**Definition:** When a CLI may prompt: **only when stdin/stdout is a TTY**, always with a flag override (`--yes`, `--force`, explicit value flags) so scripts never hang, and never as the sole path to any capability.

**Reasoning / Evidence:** A prompt in a non-TTY context (CI, cron, pipe) hangs forever or crashes — the automation-killer; conversely, interactive niceties (confirmation prompts, wizards like `npm init`, select-menus) genuinely help humans. TTY detection resolves the tension mechanically.

**When to use:**
- Confirmations for destructive ops when TTY ([#safety-dry-run-idempotency-destructive-guards](#safety-dry-run-idempotency-destructive-guards)); interactive wizards for complex first-time setup (`tool init`) — each answer echoed as the equivalent flag so users learn the scriptable form.
- Secrets: prompt with hidden input when TTY; accept via env var/file/flag(discouraged: leaks into history/ps) for automation.

**When NOT to use / exceptions:**
- No TTY: fail fast with a clear message naming the missing input/flag ("--env required in non-interactive mode"), exit 2.
- Never invisible-hang reading stdin when the user probably didn't mean to pipe ([#help-system-design](#help-system-design): no-args → help).
- `--no-input` flag convention for forcing non-interactive even in TTY (CI images with pseudo-TTYs).

**Pros:** Human convenience + automation safety coexisting.
**Cons:** Dual paths to test; prompt flows must stay in sync with their flag equivalents.

**Implementation guidance:**
- Every prompt: default answer shown (`Proceed? [y/N]` — capital = default), Enter takes default, Ctrl+C aborts cleanly.
- Wizard outputs: print the resulting config/command ("Equivalent: tool create --name x --region eu") — teaches the fast path ([01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use)).

**Related:** [#safety-dry-run-idempotency-destructive-guards](#safety-dry-run-idempotency-destructive-guards), [#standard-flag-vocabulary](#standard-flag-vocabulary).

**Sources:** [clig.dev: interactivity](https://clig.dev/#interactivity).

### Configuration precedence

**Definition:** The layered config resolution order, most-specific wins: **command-line flags > environment variables > project-level config file > user-level config (`$XDG_CONFIG_HOME/tool/`) > system-level > built-in defaults** — each layer inspectable and documented.

**Reasoning / Evidence:** The hierarchy matches intent-specificity (a flag is this-invocation intent; user config is standing preference); violating it (config file overriding an explicit flag) produces baffling "but I passed --env=staging!" bugs. XDG base-dir compliance keeps home directories sane ([14-platform-desktop.md](14-platform-desktop.md#linux-conventions)); dotfile sprawl (`~/.mytool`) is the legacy anti-pattern.

**When to use:**
- Any tool with persistent preferences: implement the full ladder; `tool config` subcommands (get/set/list) beat hand-editing; `tool config list --show-origin` (git-style) reveals which layer set what — the debuggability feature.
- Env vars: `TOOL_*` prefix namespace, documented; respect ecosystem standards (`NO_COLOR`, `EDITOR`, `PAGER`, `HTTP_PROXY`, `CI`).

**When NOT to use / exceptions:**
- Secrets in config files: warn/refuse; point to keychain/env/secret-manager integration.
- Don't auto-write config the user didn't ask to persist (surprise state).

**Pros:** Predictable overrides; per-project reproducibility (committed project config) + personal ergonomics coexisting.
**Cons:** Six layers of "where did this value come from?" — solved only by the show-origin introspection; build it.

**Implementation guidance:**
- Formats: TOML/YAML/JSON — pick one, document schema, validate with helpful errors (line numbers).
- Project config discovery: walk up from cwd (like git); explicit `--config` beats discovery.

**Related:** [14-platform-desktop.md](14-platform-desktop.md#linux-conventions), [#standard-flag-vocabulary](#standard-flag-vocabulary).

**Sources:** [clig.dev: configuration](https://clig.dev/#configuration), [XDG Base Directory spec](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).

### Safety: dry-run, idempotency, destructive guards

**Definition:** The CLI safety toolkit: `--dry-run` previews (print what *would* happen), idempotent operations (safe re-runs — crash-resume without damage), confirmation for destructive scope (TTY prompt; `--force`/typed-name for catastrophes), no-destructive-defaults, and recoverability (trash over unlink, backups before overwrite) where feasible.

**Reasoning / Evidence:** CLIs operate at maximum blast radius (prod databases, file trees, fleets) with minimum feedback; `rm -rf` lore exists because text interfaces skip the [error-prevention](01-core-principles.md#5-error-prevention) layers GUIs grew. Terraform's plan/apply split is the canonical modern pattern: preview as a *first-class artifact*.

**When to use:**
- `--dry-run` on every mutating command — output format identical to real-run reporting ("would delete 3 files:").
- Destructive + irreversible: TTY confirmation showing concrete scope ("This deletes the 'prod' database (2.3 GB, 14 collections). Type the database name to confirm:") — scope-typed confirmation for the catastrophic tier ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs) escalation ladder, text edition).
- Idempotency: `create` that tolerates existing (with `--if-not-exists` semantics or clear states), resumable batch jobs, atomic writes (temp file + rename).

**When NOT to use / exceptions:**
- Read-only commands need none of it — don't confirm-fatigue ([02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)); guard rails scale with irreversibility, not with mutation per se.
- `--force`/`--yes` must genuinely bypass every prompt (automation contract) — a "force" that still prompts breaks CI.

**Pros:** Catastrophe prevention; operator confidence (dry-run before prod is a workflow, not a feature).
**Cons:** Dry-run fidelity is engineering-hard (must simulate accurately or it's a false promise — label estimates as estimates).

**Implementation guidance:**
- Destructive verbs explicit: `delete`, not overloaded flags on innocuous commands; plural scope in output ("3 resources") before confirm.
- Recoverability ladder: soft-delete/trash where the domain allows; `--backup` conventions for in-place file edits; print undo instructions after destructive success where they exist.
- Crash safety: journaling/atomic ops so Ctrl+C mid-run never leaves corruption; state what was and wasn't done on abort.

**Related:** [01-core-principles.md](01-core-principles.md#5-error-prevention), [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs), [#interactivity-rules](#interactivity-rules).

**Sources:** [clig.dev: robustness](https://clig.dev/#robustness), [Terraform plan/apply](https://developer.hashicorp.com/terraform/cli/commands/plan).

### Composability

**Definition:** Playing well in pipelines: clean stdout ([#output-design-stdout-stderr-and-humans](#output-design-stdout-stderr-and-humans)), `-` for stdin/stdout file arguments, streaming over buffering (emit results as produced), line-oriented/JSON-lines output, sane signal handling (SIGPIPE death without stack traces), and doing one thing well enough that other tools want to pipe to you.

**Reasoning / Evidence:** The UNIX pipeline is the original component architecture; tools honoring its contracts inherit an ecosystem (grep/jq/xargs/parallel) as free features. Buffering everything before output breaks `| head` responsiveness and memory on big streams.

**When to use:**
- Accept stdin where a file makes sense (`tool process -` or auto-detect piped stdin); emit progressively (line-buffered) for streamable results.
- Handle SIGPIPE gracefully (downstream `head` closing is normal termination, not an error dump); Ctrl+C cleanup handlers.

**When NOT to use / exceptions:**
- Results requiring global aggregation (sorted output) legitimately buffer — document it.
- Interactive-only commands are outside the pipeline contract (and must fail fast when piped — [#interactivity-rules](#interactivity-rules)).

**Pros:** Ecosystem leverage; power users compose capabilities you never built.
**Cons:** Streaming architectures constrain internal design (worth it for data-processing tools).

**Implementation guidance:**
- Test suite: `tool … | head -1` (SIGPIPE), `echo x | tool` (stdin), `tool | cat` (TTY-detection), 10GB-input memory behavior.
- Filter-style tools: stdin→stdout by default, in-place editing only by explicit flag (with backup).

**Related:** [#machine-readable-output](#machine-readable-output), [#human-first-cli-design](#human-first-cli-design).

**Sources:** [clig.dev: composability](https://clig.dev/#output), [TAOUP: Rule of Composition](http://www.catb.org/~esr/writings/taoup/html/ch01s06.html).

### Performance and responsiveness

**Definition:** CLI speed as UX: startup latency (<100ms to first output feels instant — interpreter/VM startup is the enemy), fast common paths (status/list commands are run hundreds of times daily), and responsiveness signals when slow is unavoidable ([#progress-indication](#progress-indication)).

**Reasoning / Evidence:** CLIs sit in interactive loops (shell prompt cadence) and scripts (invoked thousands of times); 500ms startup makes a tool *feel* broken at the prompt and makes scripts crawl. The 100ms instantaneous threshold ([04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds)) applies to time-to-first-byte of output.

**When to use:**
- Budget: `tool --version` and `--help` <100ms (these gate perceived quality); frequent read commands <300ms; lazy-load heavy dependencies behind the subcommands that need them.
- Cache expensive lookups (with visible staleness controls: `--refresh`, TTLs documented) — cache *transparency* prevents "why is it lying" mistrust.

**When NOT to use / exceptions:**
- Network-bound operations are excused from startup budgets but not from progress feedback.

**Pros:** Prompt-loop fluency; script-scale throughput.
**Cons:** Language/runtime choice consequences (JIT warmup, interpreter startup) are architectural — know them before v1.

**Implementation guidance:**
- Measure cold start in CI; defer imports/plugin scans; consider daemon modes only with full transparency (a background daemon surprising users is a trust cost).
- Show cached-data age when serving stale ("as of 5m ago; --refresh to update").

**Related:** [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds), [#progress-indication](#progress-indication).

**Sources:** [clig.dev: responsiveness](https://clig.dev/#robustness).

### Versioning and deprecation communication

**Definition:** Managing change in a scripted-against interface: semantic versioning of the *CLI surface* (commands, flags, output schemas, exit codes are all API), deprecation warnings (stderr, non-breaking, with migration path and timeline), aliases through renames, and changelogs that call out breaking changes first.

**Reasoning / Evidence:** CLI consumers are scripts that don't read release notes; silent breaking changes detonate in CI weeks later. The deprecation pipeline (warn → alias → remove over major versions) is the only humane path; git's multi-year `checkout`→`switch/restore` migration is the reference case.

**When to use:**
- Renames: old form aliased + stderr deprecation note ("'tool old-cmd' is deprecated, use 'tool new-cmd'; alias removed in v3") for ≥1 major version.
- Output changes: additive only within a major for `--json`; human output free to change (that's why the split exists — [#machine-readable-output](#machine-readable-output)).
- Breaking releases: loud changelog section, migration guide, ideally a `tool migrate`/doctor command.

**When NOT to use / exceptions:**
- Security-forced breaks skip the timeline — communicate maximally instead.
- Deprecation warnings must be suppressible in the interim (`--no-warnings`/env) or CI logs drown.

**Pros:** Trust with the automation ecosystem; upgrade willingness (users who've been burned pin versions forever).
**Cons:** Alias/compat maintenance burden — bounded by the removal schedule.

**Implementation guidance:**
- `--version` semver; document the stability contract explicitly ("human output unstable; --json stable within major").
- Telemetry (if any, opt-in/transparent) on deprecated-path usage to time removals with data.

**Related:** [10-design-systems.md](10-design-systems.md#versioning-and-deprecation), [#machine-readable-output](#machine-readable-output).

**Sources:** [semver.org](https://semver.org/), [clig.dev: future-proofing](https://clig.dev/#future-proofing).

---

## Full-screen TUIs

### When a TUI beats a CLI

**Definition:** Full-screen terminal applications (htop, lazygit, k9s, tig — built on curses/bubbletea/ratatui-class frameworks): persistent visual state, keyboard-driven manipulation, live updates — versus the command-invocation model.

**Reasoning / Evidence:** TUIs win when the task is *session-shaped*: browsing/monitoring (continuous state beats repeated `list` invocations), multi-step manipulation of visible objects (staging hunks in lazygit vs `git add -p` round-trips), dashboards (live refresh). CLIs win for automation (TUIs are unscriptable), one-shot operations, and composition. The best tools ship both: `tool` commands + `tool ui` for the session mode.

**When to use:**
- Monitoring (top-likes), object browsing at scale (k9s over kubectl-list loops), interactive selection/manipulation workflows (interactive rebase, log exploration).
- Terminal-bound audiences (SSH, servers) needing GUI-like tasks without GUI infrastructure.

**When NOT to use / exceptions:**
- Never TUI-only for operations automation needs — every TUI action should have a CLI/flag equivalent ([#interactivity-rules](#interactivity-rules) generalized).
- Simple flows: a well-made prompt sequence beats a full-screen app's complexity.

**Pros:** GUI-class efficiency over SSH; keyboard-driven speed; live state.
**Cons:** Unscriptable, screen-reader-hostile (terminal a11y is poor — keep the CLI parity path as the accessible route), framework/terminal-compat maintenance.

**Implementation guidance:**
- Enter/exit cleanly: alternate screen buffer (restore the user's scrollback on exit), terminal state restoration on crash (raw-mode traps).
- Degrade over SSH/limited terminals: detect capabilities (truecolor→16-color, Unicode→ASCII borders).

**Related:** [#tui-interaction-patterns](#tui-interaction-patterns), [#human-first-cli-design](#human-first-cli-design).

**Sources:** [clig.dev](https://clig.dev/), [Charm (bubbletea) docs](https://github.com/charmbracelet/bubbletea).

### TUI interaction patterns

**Definition:** The converged TUI interaction vocabulary: vim-ish navigation (hjkl + arrows both), `?` for the keybinding overlay, `q`/Esc to back out, `/` to filter/search, Enter to drill in, Tab between panes, `:` or a command palette for command mode, contextual key-hint bar at the bottom, mouse support as enhancement, and resize-responsive pane layouts.

**Reasoning / Evidence:** These conventions emerged from vi/less/mutt lineage and are now cross-tool muscle memory ([Jakob's Law](01-core-principles.md#jakobs-law) for terminal dwellers) — a TUI inventing its own navigation fights every user's spinal cord. The persistent key-hint footer (lazygit/k9s pattern) solves the discoverability problem that killed earlier TUIs ([01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall)).

**When to use:**
- Standard keys: arrows *and* hjkl navigate; `?` overlays full keymap; Esc/q backs out one level (q at root = quit, confirm if state at risk); `/` incremental filter; g/G top/bottom; Ctrl+d/u paging.
- Layout: pane-based (list + detail), focused pane visually marked, Tab/number keys switching; bottom status bar: mode/context + top 5–7 key hints.
- Live data: refresh indicators, pause capability, changed-item highlighting (brief — [02-cognitive-foundations.md](02-cognitive-foundations.md#change-blindness)).

**When NOT to use / exceptions:**
- Don't require modal (vim-insert-style) editing knowledge for basic use — modes need unmistakable indication ([02-cognitive-foundations.md](02-cognitive-foundations.md#mode-errors)) and non-modal alternatives for entry-level operations.
- Mouse: support clicks/scroll where terminals allow, but never require ([05-accessibility.md](05-accessibility.md#keyboard-access-and-no-keyboard-trap) inverted — keyboard is the guaranteed channel here).

**Pros:** Instant familiarity for terminal natives; extreme operation speed.
**Cons:** Steep for non-vim users (the `?` overlay and hint bar are the mitigation); terminal-emulator inconsistencies (key reporting, resize events) demand testing breadth.

**Implementation guidance:**
- Resize: reflow gracefully to 80×24 minimum; collapse panes below thresholds rather than corrupting.
- Config: keybindings remappable (config file), theme-able (and respect terminal palette by default).
- State: remember pane sizes/filters/position across sessions ([14-platform-desktop.md](14-platform-desktop.md#window-management-and-state-restoration) in spirit).

**Related:** [#when-a-tui-beats-a-cli](#when-a-tui-beats-a-cli), [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns).

**Sources:** [lazygit](https://github.com/jesseduffield/lazygit), [k9s](https://k9scli.io/), [Charm design resources](https://charm.sh/).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Philosophy | Human-first, machine-honoring; TTY-detect | — |
| Command structure | noun-verb subcommands; fixed verb vocabulary | ≤2 nesting levels |
| Args vs flags | ≤2 positionals; flags self-document | `--` and `-` conventions |
| Standard flags | -h/--help, --version, -v, -q, --json, --dry-run, -y, --no-color | Never repurpose |
| Help | Examples first; no-args → help; ≤80 cols | 2–5 real examples |
| Did-you-mean | Suggest, never auto-run | Levenshtein ≤2 |
| Completion | Ship for bash/zsh/fish; dynamic values | <100ms dynamic |
| Streams | stdout = data; stderr = conversation | Lint: no logs to stdout |
| Color | Semantic + text prefixes; NO_COLOR; TTY-only | 16-color base |
| Machine output | --json day one; stable schema; JSON Lines | Additive within major |
| Progress | Nothing silent >1s; stderr; CI-plain | Spinner after ~300ms |
| Errors | error: + what + fix + link; traces behind --debug | Stable error codes |
| Exit codes | 0 = success, meaningful non-zero, documented | 2 = usage; 130 = SIGINT |
| Interactivity | Prompt only on TTY; flag override always | [y/N] defaults |
| Config | flags > env > project > user > system > defaults | XDG dirs; show-origin |
| Safety | --dry-run everywhere; typed confirm for catastrophe | Scope shown before confirm |
| Composability | Stream, handle SIGPIPE, accept - | `\| head` test |
| Performance | Instant startup; cache with transparency | --help <100ms |
| Deprecation | Warn → alias → remove across majors | ≥1 major grace |
| TUI when | Sessions: browse/monitor/manipulate; CLI parity kept | 80×24 minimum |
| TUI keys | hjkl+arrows, ?, /, q/Esc, hint bar | 5–7 hints visible |
