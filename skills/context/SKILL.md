---
name: context
description: wf phase 2 (Context) — map the code/reality the task touches, validate feasibility, baseline and approve scope. Invoked via /wf:dev when context is the active phase.
---

# /wf:context — Context (interactive exit)

Contract first:

- `wf record context-map --json '{"entries":["path — why it matters", …],"sufficiency":"…"}'`
  — the files/modules/systems actually examined (use the Explore subagent
  for the mapping; record what you verified, not what you assume)
- `wf record assumption text="…" high_risk=true|false` — for every
  assumption the map could not verify
- `wf record reclassify result=confirmed` — the checkpoint: does the
  classification survive contact with the map? If not:
  `wf run branch --family <f> --intent <i> --reason "reclassify: …"`
- `wf approve scope --payload "<requirement baseline + high-risk assumptions>"`
  — pose the confirmation via AskUserQuestion (the hook anchors the answer)
  — present the user: surviving requirements (active/dropped/revised),
  high-risk assumptions, and the feasibility read; record after explicit OK
- `@wf:researcher` verdict when external research ran; otherwise
  `wf contract waive context.research-grounded --reason "…"`

Procedure:
1. Map before deciding: read the code paths the task touches; verify the
   claimed behavior actually exists (grep + targeted reads beat memory).
2. assessment family: the map itself is the core deliverable input — go one
   level deeper than feels necessary.
3. Baseline requirements: mark each Frame requirement active, dropped, or
   revised (`wf record requirement --json …` updates by new record; keep
   rid stable).
4. `wf phase exit` when the contract is met.
