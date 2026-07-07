---
name: status
description: Show the wf workflow status — run, phase, unmet contract items, open tasks. Use when the user asks where the workflow stands.
---

# /wf:status

Run `wf status` and relay it verbatim (it is regenerated from disk and
authoritative). For state problems, run `wf doctor` and follow its
remediations.

For health signals across runs (loops, escapes, self-attested counts,
ungrounded ACs, lesson efficacy, deliver-reached), run `wf report`
(`--run <id|current>` for one run, `--json` for machines) and relay the
⚠ concern lines verbatim — they are the dishonesty signatures.
