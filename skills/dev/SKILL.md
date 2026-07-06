---
name: dev
description: Session entry for the wf enforced workflow. Use when the user asks to start, resume, or continue development work in a wf-adopted project (any coding, design, review, or investigation task).
---

# /wf:dev — session entry

State lives on disk in `.workflow/`; the injected `[wf]` status block is
authoritative over conversation memory. This skill only routes — the engine
decides, phase skills instruct, hooks enforce.

1. Run `wf status`. Three cases:
   - **No active run** → classify the user's task and start one:
     `wf run start --family diff|artifact|assessment --intent <tag>`
     - `diff` = the deliverable is a code change (intents: new, fix,
       refactor, test, deploy)
     - `artifact` = authored documents (arch-design, doc-create, doc-update)
     - `assessment` = a findings report, nothing modified (code-review,
       arch-review, investigate, research)
     The classification is provisional — Frame confirms it with the user.
   - **Active run** → invoke the phase skill named on the
     `resume procedure:` line (e.g. `/wf:frame`). Do not re-plan from memory.
   - **Parked run** → tell the user why it parked (the status shows the
     reason); on their go-ahead `wf run resume`, or branch for new work:
     `wf run branch --reason …`.
2. If `wf status` reports state problems, run `wf doctor` and follow its
   remediation (a fresh clone re-attaches with `wf run adopt`).
3. Never edit files under `.workflow/state|log|runs` by hand; every fact is
   recorded through `wf` commands. Escapes exist and are audited:
   `/wf:park` (honest stop), `/wf:force` (bypass one gate, escalates).
