---
name: frame
description: wf phase 1 (Frame) — understand the task with the user, classify, risk-screen, elicit requirements and acceptance criteria through lenses. Invoked via /wf:dev when frame is the active phase.
---

# /wf:frame — Frame (interactive)

Contract first — the exit gate verifies these records exist (`wf status`
lists what is still missing, with commands):

- `wf record classification family=<f> intent=<i> restated="…"` — your
  one-paragraph restatement of the task, confirmed with the user
- `wf approve frame --payload "<family/intent: restatement>"` — record it
  ONLY after the user explicitly confirms (never infer approval). Pose the
  confirmation via AskUserQuestion, naming the classification in the
  question — the hook captures the answer, infers the topic, and anchors
  the approval (`answer_ref`); with config `approvals: hardened` an
  un-anchored or topic-mismatched approval is refused; with
  `approvals: challenge` additionally ask the user to type the code shown
  in THEIR statusline (the engine never shows it to you) — the captured
  answer must contain it
- `wf risk scan --text "<the task + restatement>" [--add signal]…` — the
  deterministic screen; add signals your judgment finds (auth, network,
  data, boundary, destructive, concurrency, ui)
- Per lens the scan selected: ≥1 `wf record ambiguity lens=<l> text="…"
  disposition=resolved|logged|deferred` — or an explicit
  `lens=<l> none=true disposition=none` with the reason in `text`
- diff/artifact only:
  - `wf record requirement rid=SWR-1 level=software text="…" status=active
    --json '{"acs":[{"id":"AC-1","text":"…","verifiable":true}]}'` — every AC
    must be verifiable in principle (name how)
  - `wf record completeness --json '{"items":[{"case":"empty input","disposition":"covered:AC-2"}]}'`
    — the negative-space walk: error, empty, max, concurrent, unhappy paths.
    Dispositions carry a traceable vocabulary (write-time enforced):
    `covered:<AC id>` (that AC's tests carry the case — add the AC if none
    does), `out-of-scope: <reason>`, or `accepted-risk: <reason>`. The gate
    expects ≥3 dispositioned cases (`frame.completeness-depth`); genuinely
    smaller spaces: waive it with the reason
  - spawn `@wf:adversary` (abuse-case mode) and `@wf:lens-reviewer`
    (security lens) on the framed requirements — their verdicts are captured
    automatically at completion; ux-enabled projects with a UI-bearing task
    also spawn `@wf:ux-reviewer` (usability-lens mode) and record its
    ambiguities like any other lens
- intent fix/investigate: `wf origin discover --path <file> --text "<code
  fragment>"` (git-grounded attribution; falls back to
  `wf record origin attribution="…"` with reduced confidence when git is
  inconclusive)
- intent fix additionally: a REGRESSION requirement traced to the origin —
  `wf record requirement rid=SWR-n … origin=<origin-record-id>` whose AC
  states the defect no longer reproduces (`frame.fix-regression`). Its red
  test at Build IS the reproduction; a fix without one never proves the bug
  existed. Origin genuinely inconclusive: waive with the reason

Procedure:
1. Restate the task in your own words. Ask the user targeted questions per
   selected lens until each lens yields a real ambiguity or a reasoned none.
2. Blocking ambiguities stop here (the Stop gate lets the turn end while an
   approval/question is pending); deferrable ones are recorded `deferred` —
   Plan forces their disposition later.
3. When the contract is met: `wf phase exit` (the engine advances to
   Context). Do not summarize state from memory — `wf status` is the truth.
