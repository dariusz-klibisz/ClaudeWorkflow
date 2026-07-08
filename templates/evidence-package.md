# Evidence package: <run-id> — <standards in force>

> Engineering-discipline evidence assembled from this run's ledger. **This
> workflow is NOT a compliance tool and this package is NOT a certification
> artifact** — it is input for the accredited assessor your standard
> requires. Every item below must reference a record, artifact, or file
> that exists; the compliance-reviewer treats an item contradicting the
> ledger as a critical.

## Scope

- **Run**: <run-id> · **Family/intent**: <family>/<intent>
- **Standards in force**: <from the installed packs — e.g. ISO 26262 (ASIL <x>), IEC 62304 (Class <x>)>
- **Safety/criticality classification**: <recorded where — record id / ADR>

## Requirements & traceability

- Requirements + ACs: <record ids / docs/requirements/>
- RTM: <docs/requirements/RTM-<run-id>.md — generate with `wf trace --rtm --write`>

## Design evidence

- Selected design + ADRs: <docs/design/…, docs/architecture/adr/…>
- Required analyses for the classification (per the standard's checklist):
  <threat model / attack tree / safety analyses — path or ADR-accepted waiver each>

## Verification evidence

- Test plan: <docs/test/…>
- Per-AC grounded results: <from `wf report --run <id>`: N ACs, N grounded greens>
- Coverage / required metrics: <metric records; thresholds from config>
- Reviews: <verdict records — quality, security, testing, conformance, adversary, compliance>

## Configuration & release

- Delivery artifact: <PR/release ref — the role=delivery artifact record>
- Deviations & waivers: <each with its approval/disposition record id>
- Escapes (forces/parks/loops): <from `wf report` — none is a claim the ledger must back>

## Standard-specific items

<one subsection per standard in force: walk the "Evidence" section of
.workflow/packs/<pack>/<standard>.md and reference each listed item>
