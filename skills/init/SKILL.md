---
name: init
description: Adopt the wf workflow in this project (one-time setup). User-invoked only.
disable-model-invocation: true
---

# /wf:init — project adoption

Idempotent; safe to re-run after plugin updates. Steps:

1. `wf init` — creates `.workflow/` (config.json, log/, contracts.d/,
   .gitignore) and records the plugin version.
2. **Merge** (never overwrite) `.claude/settings.json` so collaborators get
   the plugin on folder trust:
   ```json
   {
     "extraKnownMarketplaces": {
       "claude-workflow": {
         "source": { "source": "github", "repo": "dariusz-klibisz/ClaudeWorkflow" }
       }
     },
     "enabledPlugins": { "wf@claude-workflow": true }
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
4. If a v0.36 scaffold exists (`.workflow/manifest.json` present): offer to
   archive the old `hooks/ steps/ agents/` trees to `.workflow/legacy/` and
   strip old hook wiring from `.claude/settings.json`. Old audit artifacts
   stay readable; the new engine never parses them.
5. Verify: `wf doctor --bootstrap` prints the engine version and contract
   count. Commit `.workflow/config.json`, `.workflow/.gitignore`,
   `.claude/settings.json`, and `CLAUDE.md`.
