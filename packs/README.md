# wf contract packs

Official add-only contract packs. Install into a project with:

```
wf pack install <plugin-root>/packs/<name>        # e.g. packs/sbom
wf pack install <plugin-root>/packs/regulated/<std>
```

Install is validated (the merged spec must load strictly) and add-only — a
pack can never weaken or replace shipped contract items. YAML lands in
`.workflow/contracts.d/`, docs land in `.workflow/packs/<name>/`; commit
both so the pack travels with the repo. Uninstall = remove those files.

| Pack | Domain | Adds |
|---|---|---|
| `sbom` | supply chain | `x-sbom` record + a Ship gate requiring a generated SBOM per delivered diff |
| `regulated/iso-26262` | automotive (ASIL QM–D) | compliance-reviewer gates at Design/Verify + evidence package at Ship |
| `regulated/iec-62304` | medical (Class A–C) | same shape |
| `regulated/do-178c` | aviation (DAL E–A) | same shape |
| `regulated/iec-61508` | E/E/PE safety (SIL 1–4) | same shape |
| `regulated/en-50128` | rail (SIL 0–4) | same shape |
| `regulated/nist-800-53` | US federal (FIPS 199) | same shape |

The regulated packs activate the shipped `compliance-reviewer` agent (one
competence, many standards — the injected scope names the standard; its
checklist is the pack's `.md`, injected at spawn from `.workflow/packs/`).
Several packs may be installed together: each standard gets its own scoped
verdict gates; one evidence package covers the standards in force.

**NOT A COMPLIANCE TOOL.** The regulated packs provide structured,
agent-facing reminders and evidence discipline only. They do NOT certify
or demonstrate conformance with any standard — engage an accredited
assessor before claiming compliance.
