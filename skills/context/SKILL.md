---
name: context
description: wf phase 2 (Context) — map the code/reality the task touches, validate feasibility, baseline and approve scope. Invoked via /wf:dev when context is the active phase.
---

# /wf:context — Context (interactive exit)

Contract first:

- `wf record context-map --json '{"entries":["path — why it matters", …],"sufficiency":"…"}'`
  — the files/modules/systems actually examined (use the Explore subagent
  for the mapping; record what you verified, not what you assume). The gate
  expects ≥3 entries (`context.map-depth`) — the code, its tests, its
  callers/config all count; a genuinely smaller surface: waive with reason
- `wf record assumption text="…" status=open high_risk=true|false` — for
  every assumption the map could not verify. Assumptions have a lifecycle:
  recorded `open` here, discharged `validated|invalidated` at Verify
  (high-risk ones gate the exit there — an assumption is a debt, not a note)
- `wf record reclassify result=confirmed` — the checkpoint: does the
  classification survive contact with the map? If not:
  `wf run branch --family <f> --intent <i> --reason "reclassify: …"`
- `wf approve scope --payload "<requirement baseline + high-risk assumptions>"`
  — pose the confirmation via AskUserQuestion, naming the SCOPE in the
  question (the hook infers the topic and anchors the answer to this gate)
  — present the user: surviving requirements (active/dropped/revised),
  high-risk assumptions, and the feasibility read; record after explicit OK.
  The engine binds the approval to the current requirement/assumption
  baseline — records added later without re-approval surface as approval
  drift at Ship
- `@wf:researcher` verdict when external research ran; otherwise
  `wf contract waive context.research-grounded --reason "…"`. Fold the
  researcher's load-bearing answers into context-map entries (with the
  source named in the entry text) — transcripts die at compaction; the map
  is what Design and Plan actually read

Procedure:
1. Map before deciding: read the code paths the task touches; verify the
   claimed behavior actually exists (grep + targeted reads beat memory).
2. assessment family: the map itself is the core deliverable input — go one
   level deeper than feels necessary.
3. Baseline requirements: mark each Frame requirement active, dropped, or
   revised (`wf record requirement --json …` updates by new record; keep
   rid stable).
4. `wf phase exit` when the contract is met.
