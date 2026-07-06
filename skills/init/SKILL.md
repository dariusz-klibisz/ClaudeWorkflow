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
   count; confirm the CLAUDE.md block markers were written verbatim. Commit
   `.workflow/config.json`, `.workflow/.gitignore`, `.claude/settings.json`,
   and `CLAUDE.md`.
5. Tell the user: permission rules written to settings.json take effect on
   the **next session** — restart (or check `/permissions`) so `wf` calls run
   prompt-free. If a `wf` command still prompts mid-session, accepting
   "don't ask again for `wf record *`" is safe: the engine only writes under
   `.workflow/`.

Native Windows (no Git Bash) note: until M5, run the engine installer once by
hand — `powershell -File <plugin-root>/scripts/bootstrap.ps1` — since the
SessionStart bootstrap ships as a sh script only.
