---
name: park
description: Park the active wf run — the audited honest stop. User-invoked only; the model cannot take this escape unilaterally.
disable-model-invocation: true
---

# /wf:park

Ask the user for the reason if not given, then:

    wf park --reason "<why the run stops here>"

Parking clears every sequencing gate for the run, is recorded and surfaced
in reports, and is always resumable: `wf run resume`. Parking is the honest
terminal state when genuinely blocked — prefer it over forcing gates.
