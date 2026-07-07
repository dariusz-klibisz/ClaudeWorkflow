---
name: design
description: wf phase 3 (Design) — staged option evaluation to a reviewed, critic-checked, user-approved design. Invoked via /wf:dev when design is the active phase (diff and artifact families).
---

# /wf:design — Design (interactive exit; waivable for trivial diffs)

Trivial diff with no design decisions? `wf phase waive design --reason "…"`
(recorded, surfaced at Ship). Otherwise, contract first:

- `wf record option-set stage=system --json '{"candidates":[…],"selected":"…","rejected":[{"id":"…","reason":"…"}]}'`
  — 2–4 GENUINE candidates with selection and rejection reasons
- `wf record option-set stage=software --json …` — same at software level
- ux-enabled projects with UI-bearing change: `stage=ux` set too
- spawn `@wf:design-reviewer` on the SELECTED option — fix findings until
  its verdict is clean (verdicts auto-captured; a failing auto verdict is
  sticky — fix the design, re-run the reviewer; manual records cannot
  override it)
- security signals present: `wf record threat --json '{"entries":[…]}'`
  (STRIDE over the trust boundaries) and spawn `@wf:adversary`
  (attack-tree mode); mitigate or ADR-accept high-feasibility paths
- spawn `@wf:critic` — independent go/no-go; a `risky` verdict passes only
  after `wf record disposition ref=<verdict-id> text="…"` per concern
- architectural decision: `wf doc new adr --slug <decision-name>` (engine
  copies the template into `docs/architecture/adr/` and records the stub),
  author it, then flip the record:
  `wf record artifact updates=<id> status=present`
- `wf approve design --payload "<selected options + risks + testability>"`
  — pose the confirmation via AskUserQuestion (the hook anchors the answer)
  — after presenting to the user

Loop re-entries (from Verify) must reference previously rejected option IDs
— never re-propose a rejected option as new.

`wf phase exit` when met.
