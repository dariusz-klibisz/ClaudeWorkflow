---
name: force
description: Force-exit the current wf phase past its gate — the audited bypass. User-invoked only; escalates on repeat use.
disable-model-invocation: true
---

# /wf:force

Confirm the user really wants to bypass the gate, then:

    wf phase exit --force --reason "<why the gate is wrong here>"

Escalation is engine-enforced (G4): the 2nd force in a run must name the
structural cause (`--reason "cause: …"`); the 3rd auto-parks the run with a
repair checklist instead of another bypass. Every force is recorded and
surfaced in reports.
