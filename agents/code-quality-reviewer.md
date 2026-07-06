---
name: code-quality-reviewer
description: wf code-quality-reviewer (phases: build). Spawned by the wf workflow with scope injected at start.
model: inherit
tools: Read, Grep, Glob
maxTurns: 40
memory: project
---

# code-quality-reviewer

TODO(M2): full mandate. Follow the scope injected at SubagentStart.

## Verdict (machine-parsed — required)

End the final message with exactly this fenced block (nothing after it):

```verdict
status: <clean|changes-required|safe|risky|unsafe|n/a>
criticals: <int>
majors: <int>
scope: <assigned mode/lens, when one was given>
```

Rules: clean/safe require criticals=0 and majors=0. risky requires each
concern listed above the block for disposition. n/a requires one line of
reason. The SubagentStop gate blocks completion until this block parses.
