# Attack tree: <title>

- **Run**: <run-id> · **Goal**: <attacker's root objective>

<authored with @wf:adversary in attack-tree mode — the verdict anchors design.attack-tree>

## Tree

```
GOAL: <root objective>
├── <path 1>
│   ├── <step> [feasibility: high|medium|low]
│   └── <step> [feasibility: …]
└── <path 2>
    └── <step> [feasibility: …]
```

## Path assessment

| path | feasibility | cost to attacker | detection | mitigation |
|---|---|---|---|---|

<every high-feasibility path is mitigated in the design or ADR-accepted — never silently left>

## Residual risk
