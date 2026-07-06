# Decision Guides, Checklists & Templates

Use these templates and checklists to turn the principles in [`01`](01-architecture-principles.md)–[`07`](07-security-reliability-operations.md) into concrete, recorded decisions. Adapt each to your project's risk and team size — a small internal tool and a regulated platform need different levels of rigor.

---

## 1. Architecture Decision Record Template

```md
# ADR-0001: Title

## Status
Proposed | Accepted | Superseded | Deprecated

## Date
YYYY-MM-DD

## Context
What problem are we solving?
What constraints matter?
What quality attributes are affected?
What assumptions are we making?

## Decision
What did we decide?

## Alternatives Considered
- Option A: pros, cons, why rejected or accepted.
- Option B: pros, cons, why rejected or accepted.

## Consequences
Positive consequences.
Negative consequences.
New risks.
Operational impact.
Security impact.
Cost impact.

## Validation
How will we know if this works?
What metrics, tests, or fitness functions apply?

## Review Date
When should we revisit this?

## References
Links to tickets, diagrams, prototypes, benchmarks, standards, or research.
```

---

## 2. Architecture Review Checklist

- [ ] System context is documented (C4 context view).
- [ ] Users, external systems, and trust boundaries are shown.
- [ ] Architecturally significant requirements are explicit and measurable.
- [ ] Data ownership is clear.
- [ ] Consistency model is defined.
- [ ] Security model is documented.
- [ ] Failure modes are identified.
- [ ] Observability exists for critical paths.
- [ ] Deployment and rollback are defined.
- [ ] Cost drivers are understood.
- [ ] Accessibility is considered for any human-facing UI.
- [ ] Privacy-sensitive data is inventoried.
- [ ] Key decisions have ADRs.
- [ ] Operational ownership is clear.
- [ ] Compliance requirements are mapped.

---

## 3. Quality Attribute Scenario Template

```md
## Quality Attribute Scenario

Attribute: Reliability | Performance | Security | Accessibility | etc.

Source:           Who or what triggers the scenario?
Stimulus:         What happens?
Environment:      Under what conditions?
Artifact:         What part of the system is affected?
Response:         What should the system do?
Response Measure: How will we measure success?

Example:
Under normal production load, when a user submits checkout, the system
returns a confirmed order or actionable failure within 2 seconds at p95,
without double-charging the user.
```

---

## 4. Technology Selection Rubric

Score each category from 1 to 5 and write notes.

| Category | Score | Notes |
| --- | --- | --- |
| Solves the actual problem |  |  |
| Team expertise |  |  |
| Operational maturity required |  |  |
| Security posture |  |  |
| Ecosystem and maintenance |  |  |
| Performance fit |  |  |
| Cost |  |  |
| Vendor lock-in |  |  |
| Migration path |  |  |
| Observability / supportability |  |  |
| Compliance / licensing |  |  |
| Reversibility |  |  |

Decision notes:
- What is the smallest experiment that reduces uncertainty?
- What would make us reverse this decision?
- Who owns the technology long-term?

---

## 5. Microservices Readiness Checklist

Use this **before** splitting services ([02 §4.2](02-architecture-patterns.md#42-microservices-architecture)).

- [ ] Business boundaries are understood.
- [ ] Each proposed service has clear data ownership.
- [ ] Teams can deploy independently.
- [ ] CI/CD is automated.
- [ ] Observability includes logs, metrics, traces, and correlation IDs.
- [ ] On-call ownership exists.
- [ ] Service contracts are versioned.
- [ ] Contract tests exist or are planned.
- [ ] Idempotency is designed for messaging/retries.
- [ ] Security between services is designed.
- [ ] Local development and test strategy exist.
- [ ] Database migration strategy exists.
- [ ] Event/schema evolution strategy exists.
- [ ] Operational cost is accepted.

> If many answers are "no," prefer a **modular monolith** first.

---

## 6. Modular Monolith Checklist

- [ ] Modules align with business capabilities or subdomains.
- [ ] Dependency direction is enforced.
- [ ] Each module has an explicit public API.
- [ ] Internal module implementation is hidden.
- [ ] Database access rules are defined.
- [ ] Shared kernel / common code is small and owned.
- [ ] Tests protect module boundaries.
- [ ] Build structure reflects boundaries.
- [ ] Architecture tests prevent cycles.
- [ ] An extraction path exists for modules likely to become services.

---

## 7. Web Application Checklist

### Performance
- [ ] Core Web Vitals measured in **field** data (LCP ≤ 2.5 s, INP ≤ 200 ms, CLS ≤ 0.1 at p75).
- [ ] Performance budgets exist.
- [ ] Images optimized and dimensioned.
- [ ] JavaScript bundle size monitored.
- [ ] Third-party scripts reviewed.

### Accessibility
- [ ] Semantic HTML used.
- [ ] Keyboard navigation works.
- [ ] Focus states visible.
- [ ] Forms have labels and errors.
- [ ] Color contrast sufficient (≥ 4.5:1 for text).
- [ ] Screen-reader smoke tests cover critical flows.
- [ ] Motion preferences respected.

### Security (OWASP Top 10:2025)
- [ ] Authentication and authorization server-enforced.
- [ ] CSRF protection where needed.
- [ ] CSP defined; CORS narrow.
- [ ] Inputs validated and outputs encoded; injection prevented.
- [ ] Secrets not exposed to client bundles.
- [ ] Dependencies/supply chain scanned (SBOM); build pipeline protected.
- [ ] Rate limits protect abuse-prone endpoints.
- [ ] Errors handled without leaking internals.

### Privacy
- [ ] Personal data inventory exists.
- [ ] Analytics and third-party scripts reviewed.
- [ ] Sensitive data not in URLs or logs.
- [ ] Retention/deletion rules defined.

---

## 8. Desktop Application Checklist

### Platform UX
- [ ] Menus, shortcuts, dialogs, notifications follow platform expectations.
- [ ] High DPI and scaling tested.
- [ ] Light/dark/system theme behavior works.
- [ ] Keyboard navigation works.
- [ ] Accessibility APIs supported.

### Security
- [ ] Trust boundaries documented.
- [ ] IPC command surface minimal and typed.
- [ ] IPC sender/origin and parameters validated.
- [ ] Privileged APIs not exposed directly to untrusted UI.
- [ ] Auto-updates signed and verified.
- [ ] Code signing/notarization/store requirements handled.
- [ ] Secrets use OS credential storage.
- [ ] Logs/crash dumps exclude sensitive data.

### Distribution
- [ ] Installer and update paths tested.
- [ ] Migration from older versions tested.
- [ ] Rollback or recovery path exists.
- [ ] Release channels defined.

### Performance
- [ ] Cold and warm startup measured.
- [ ] Memory usage measured over long sessions.
- [ ] UI thread blocking avoided.
- [ ] Background work respects CPU/battery.

---

## 9. Security Review Checklist

- [ ] Assets and trust boundaries documented.
- [ ] Threat model exists for high-risk features.
- [ ] Authentication strong enough for the risk.
- [ ] Authorization enforced server-side or privileged-side.
- [ ] Secrets stored and rotated safely.
- [ ] Dependencies inventoried and scanned.
- [ ] Build pipeline protected.
- [ ] Logs avoid sensitive data.
- [ ] Security events monitored.
- [ ] Vulnerability response process exists.
- [ ] Backups protected and tested.
- [ ] Admin functions have additional controls.
- [ ] File upload and parsing sandboxed or constrained.
- [ ] SSRF and injection risks mitigated.
- [ ] Desktop/webview bridges constrained if applicable.

---

## 10. Reliability Review Checklist

- [ ] Critical user journeys identified.
- [ ] SLOs and SLIs exist.
- [ ] Timeouts set for remote calls.
- [ ] Retries use backoff and jitter.
- [ ] Circuit breakers or bulkheads protect dependencies where needed.
- [ ] Queues have dead-letter and retry policies.
- [ ] Backpressure designed.
- [ ] Graceful degradation defined.
- [ ] Backups and restores tested.
- [ ] Disaster recovery objectives (RTO/RPO) defined.
- [ ] Runbooks exist for common failures.
- [ ] Alerts actionable and tied to user impact.

---

## 11. API Design Checklist

- [ ] API uses domain language.
- [ ] Authentication and authorization specified.
- [ ] Error model consistent and machine-readable.
- [ ] Pagination exists for collections.
- [ ] Rate limits documented.
- [ ] Idempotency supported for retryable state changes.
- [ ] Versioning strategy exists.
- [ ] Deprecation policy exists.
- [ ] Sensitive fields classified.
- [ ] Request/response examples documented.
- [ ] Contract tests exist for important consumers.

---

## 12. Data Architecture Checklist

- [ ] Source of truth is clear.
- [ ] Data owner is clear.
- [ ] Schema evolution process exists.
- [ ] Consistency model defined.
- [ ] Backup/restore tested.
- [ ] Retention and deletion defined.
- [ ] Privacy classification done.
- [ ] Access controls enforced.
- [ ] Reporting/analytics do not couple operational services unnecessarily.
- [ ] Migration strategy exists for large changes.

---

## 13. Build vs Buy vs Rent Decision Template

```md
## Capability
What capability is needed?

## Strategic Importance
Is this differentiating or commodity?

## Requirements
Functional requirements.
Quality requirements.
Compliance requirements.

## Options
Build internally.
Buy product.
Use managed service.
Outsource.

## Evaluation
Cost.
Time to market.
Control.
Security.
Compliance.
Integration.
Exit path.
Vendor lock-in.
Operational burden.

## Decision
Chosen option and rationale.

## Review Trigger
What would cause reconsideration?
```

---

## 14. Pattern Selection Matrix

| Problem | Candidate Pattern | Avoid If |
| --- | --- | --- |
| Codebase hard to change but one deployment is acceptable | Modular monolith | Teams require independent deployment now |
| Teams need independent delivery and service ownership | Microservices | Domain boundaries or operations are immature |
| Legacy replacement is risky | Strangler fig | Coexistence is impossible |
| External model would pollute domain | Anti-corruption layer | External model should be canonical |
| Burst traffic overloads downstream | Queue-based load leveling | User needs immediate completion |
| Remote dependency fails intermittently | Retry with backoff | Operation is not idempotent |
| Remote dependency failure cascades | Circuit breaker | Local deterministic call |
| Read model differs from write model | CQRS | Simple CRUD is enough |
| Need full audit/history | Event sourcing | Schema evolution/privacy deletion conflicts unresolved |
| Multi-service transaction | Saga | Atomic consistency is mandatory |
| Database update plus event publish | Transactional outbox | Event loss acceptable or event store is primary |
| Different clients need different API shapes | Backend for frontend | All clients need identical API |
| Large shared deployment blast radius | Deployment stamps / cells | Small system, low isolation needs |

---

## 15. Operational Readiness Review

- [ ] Service owner identified.
- [ ] Runbook exists.
- [ ] Dashboards exist.
- [ ] Alerts exist and are actionable.
- [ ] SLOs defined.
- [ ] Deployment and rollback tested.
- [ ] Backup/restore tested if stateful.
- [ ] Load test completed for expected demand.
- [ ] Security review completed.
- [ ] Dependency failure behavior tested.
- [ ] On-call handoff completed.
- [ ] Cost estimate and budget alerts configured.
- [ ] Incident communication path defined.

---

## 16. Legacy Modernization Decision Log (Lightweight)

```md
## Capability being modernized
## Current behavior & hidden rules discovered
## Slice strategy (route / user group / workflow)
## Coexistence & data-sync approach
## Anti-corruption boundary
## Rollback / route-control plan
## End-state ownership
## Risks & mitigations
```

Pairs with the **Strangler Fig** and **Anti-Corruption Layer** patterns ([02 §4.8](02-architecture-patterns.md#48-strangler-fig-legacy-modernization), [§8.5](02-architecture-patterns.md#85-context-mapping--anti-corruption-layer-acl)).

---

## Key Cross-References
- **Principles & patterns** the checklists enforce: [`01`](01-architecture-principles.md), [`02`](02-architecture-patterns.md).
- **Quality targets:** [`06`](06-quality-attributes-tradeoffs.md). **Security/reliability/delivery:** [`07`](07-security-reliability-operations.md).
