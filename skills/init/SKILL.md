---
name: init
description: Adopt the wf workflow in this project (one-time setup). User-invoked only.
disable-model-invocation: true
---

# /wf:init — project adoption

Idempotent; safe to re-run after plugin updates. The `wf` command is on the
Bash tool's PATH while the plugin is enabled — call it bare, no env setup.

0. **Refuse on a legacy scaffold**: if `.workflow/manifest.json` exists, this
   directory was scaffolded by the old ClaudeInit generator, which used the
   same `.workflow/` layout. Stop and tell the user to remove or rename the
   old `.workflow/` tree first — wf does not migrate or share state with it.
1. `wf init` — creates `.workflow/` (config.json, log/, contracts.d/,
   .gitignore) and records the plugin version.
2. **Merge** (never overwrite) `.claude/settings.json` so collaborators get
   the plugin on folder trust, and `wf` calls need no permission prompts:
   ```json
   {
     "extraKnownMarketplaces": {
       "claude-workflow": {
         "source": { "source": "github", "repo": "dariusz-klibisz/ClaudeWorkflow" }
       }
     },
     "enabledPlugins": { "wf@claude-workflow": true },
     "permissions": { "allow": ["Bash(wf *)"] }
   }
   ```
3. Append the wf block to `CLAUDE.md` (create the file if absent), delimited
   by `<!-- wf:begin -->` / `<!-- wf:end -->` so the engine can refresh it:

   ```markdown
   <!-- wf:begin -->
   ## Workflow (wf)
   - Work happens inside runs: start every task with /wf:dev.
   - .workflow/ on disk is the source of truth; after compaction or resume,
     trust the injected [wf] status block over memory.
   - Record facts only via wf commands — never edit .workflow/state|log by hand.
   - Deliverable documents go under docs/.
   - Audited escapes: /wf:park (honest stop), /wf:force (bypass, escalates).
   <!-- wf:end -->
   ```
4. Verify: `wf doctor --bootstrap` prints the engine version and contract
   count, AND checks that the hook engine is installed at the plugin data
   path — if it reports the engine was missing, it installs it on the spot
   (exit 2 means hooks are still dead; do not proceed until fixed). Confirm
   the CLAUDE.md block markers were written verbatim. Commit
   `.workflow/config.json`, `.workflow/.gitignore`, `.claude/settings.json`,
   and `CLAUDE.md`.
5. Tell the user, explicitly, both restart caveats:
   - **Hooks**: if the plugin was installed or reloaded in THIS session
     (`/plugin`, `/reload-plugins`), the SessionStart bootstrap never fired —
     gates and verdict capture are dead until `wf doctor --bootstrap` heals
     them or the session is restarted. Any hook error mentioning a missing
     `.../data/wf-*/bin/wf` is this exact condition: re-run
     `wf doctor --bootstrap`.
   - **Permissions**: rules written to settings.json take effect on the
     **next session** — restart (or check `/permissions`) so `wf` calls run
     prompt-free. If a `wf` command still prompts mid-session, accepting
     "don't ask again for `wf record *`" is safe: the engine only writes
     under `.workflow/`.

Native Windows (no Git Bash) note: the SessionStart bootstrap ships as a sh
script, so the FIRST install needs one manual step — either
`powershell -File <plugin-root>/scripts/bootstrap.ps1` or
`wf doctor --bootstrap` from any working wf. After that, updates are
automatic on every platform: the installed engine detects version skew at
SessionStart (`wf inject session`) and re-runs the bootstrap itself
(sh or PowerShell, whichever exists).
