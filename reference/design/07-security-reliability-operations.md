# Security, Reliability, Operations & Delivery

Production software is not complete when it works on a developer's machine. It is complete when it can be **built, deployed, operated, secured, observed, recovered, and changed safely.** This file covers the secure software development lifecycle, threat modeling, identity, supply-chain security, reliability engineering (SLOs, observability, alerting, incident response), continuous delivery, and the operational qualities — cost and sustainability — that keep a system viable.

It complements the application-specific security/observability content in [`04`](04-web-application-design.md) and [`05`](05-desktop-application-design.md), and the quality-attribute framing in [`06`](06-quality-attributes-tradeoffs.md).

---

## 1. Secure Software Development Lifecycle (SSDF)

#### Summary
Security must be integrated across the lifecycle — planning, design, implementation, verification, release, operations, and vulnerability response — not bolted on before launch.

#### NIST SSDF practice groups
- **Prepare the Organization (PO):** people, processes, and technology are ready for secure development.
- **Protect the Software (PS):** protect code, build artifacts, credentials, and release integrity.
- **Produce Well-Secured Software (PW):** design, implement, review, and test to reduce vulnerabilities.
- **Respond to Vulnerabilities (RV):** identify, remediate, disclose, and prevent recurrence.

#### Guidance
Define security requirements early; threat-model significant features; use secure coding standards; automate dependency and secret scanning; review high-risk code manually; protect CI/CD credentials; generate SBOM/provenance where required; define vulnerability intake and response SLAs.

#### Common mistakes
Security review only right before launch; treating scanner output *as* the security program; no owner for vulnerability response; production secrets available to all developers.

#### Sources
NIST SSDF; OWASP ASVS; OWASP Top 10 ([04 §7](04-web-application-design.md#7-web-application-security)).

---

## 2. Threat Modeling

#### Summary
Identify what can go wrong, how likely and damaging it is, and what controls are needed — before implementation.

#### When to use
New architecture; authentication/authorization changes; payment, PII, secrets, file upload, admin features; desktop IPC and native integration; third-party integrations.

#### Guidance
Identify assets, actors, trust boundaries, data flows, and abuse cases. For each dependency, ask: what if it is malicious, compromised, slow, or unavailable? Record mitigations and residual risk; revisit when architecture changes. (Frameworks like STRIDE help enumerate threat categories.)

#### Common mistakes
No diagram of trust boundaries; only considering external attackers; ignoring insider, supply-chain, and client-side threats.

---

## 3. Authentication, Authorization, and Identity

#### Summary
**Authentication** proves *who* a subject is; **authorization** decides *what* they may do. Design them separately.

#### Guidance
Centralize identity integration where possible; enforce authorization **server-side** (or privileged-side for desktop); prefer least privilege; model permissions in business terms; use short-lived tokens; protect refresh tokens and sessions; log security-relevant decisions; test authorization paths explicitly.

**Multi-tenant isolation.** Logical tenant isolation must not rest solely on hand-written per-request filters executed under an omnipotent credential (admin/service token): one missed or null-valued filter is a cross-tenant breach with zero defense in depth. Prefer a centralized tenancy guard (one audited helper that resolves tenant context for every data access), least-privilege credentials per access pattern, and row-level security where the datastore supports it. A null or absent tenant identifier must fail closed — never degrade into an empty filter value that matches everything.

#### Common mistakes
UI hides actions but the API allows them; role checks scattered through code; tenant ID trusted from the client; long-lived bearer tokens in insecure storage; tenant filtering re-implemented by hand in each query under a shared omnipotent credential.

(Web-specific mechanics — sessions vs JWT, OAuth/OIDC, RBAC vs ABAC — are in [04 §8](04-web-application-design.md#8-authentication--authorization).)

---

## 4. Supply Chain Security

#### Summary
Your software includes dependencies, build tools, CI/CD actions, containers, runtimes, installers, and update infrastructure — all of which can be attacked. (This rose to **A03** in OWASP Top 10:2025.)

#### Guidance
Maintain a dependency inventory; use lockfiles and reproducible builds where feasible; scan dependencies and containers; pin CI actions and base images by digest; protect signing keys; require code review for dependency changes; generate an **SBOM** for serious products; monitor vulnerability disclosures.

#### Common mistakes
CI tokens with excessive permissions; pulling `latest` images without digest pinning; no owner for transitive dependency upgrades; build artifacts not signed or verified.

---

## 5. SLOs, SLIs, and Error Budgets

#### Summary
Define reliability targets using **user-relevant indicators**, and spend the resulting **error budget** deliberately.

#### Key terms
- **SLI (Indicator):** a measurement, e.g., the ratio of successful requests, or the share served under 300 ms.
- **SLO (Objective):** a target for an SLI, e.g., 99.9% success over 30 days.
- **SLA (Agreement):** a contractual promise (often with penalties) — usually looser than the internal SLO.
- **Error budget:** `1 − SLO` — the allowed unreliability, spent on releases and risk. When it's exhausted, slow down risky changes.

#### Guidance
Define SLOs for **critical user journeys**, not just infrastructure uptime; use **percentiles**, not averages, for latency; separate availability from correctness where needed; review error-budget burn to balance feature velocity vs reliability work.

#### Common mistakes
100% availability targets without business justification; infrastructure uptime used as a proxy for user success; alerts not tied to user impact. (*100% is the wrong target* — see availability nines in [06 §2.3](06-quality-attributes-tradeoffs.md#2-the-quality-attributes-catalog).)

#### Sources
Google SRE book; Well-Architected reliability pillars.

---

## 6. Observability

#### Summary
Observability lets teams understand system behavior from outputs: logs, metrics, traces, events, and profiles. (Application-level detail and the three pillars are in [04 §13](04-web-application-design.md#13-observability).)

#### Guidance
Use structured logs with correlation/request IDs; emit metrics for latency, traffic, errors, and saturation (the "four golden signals"); trace distributed requests and propagate context across queues; capture deployment events; redact secrets and personal data; define dashboards by **user journey** and service dependency. Use **RED** for services and **USE** for resources.

#### Common mistakes
Logs without context; metrics too high-level to diagnose; traces not propagated across asynchronous boundaries; alerting on symptoms no user cares about.

---

## 7. Alerting

#### Summary
Alerts should indicate **actionable, user-impacting** problems.

#### Guidance
Alert on **SLO burn**, not every low-level anomaly; every page should require human action; include runbook links; tune thresholds and silence noisy alerts; separate tickets, notifications, and pages.

#### Common mistakes
Paging on CPU alone; no owner for an alert; alerts that auto-resolve before anyone can act; no post-incident alert review.

---

## 8. Incident Response & Postmortems

#### Summary
Incidents are inevitable; the goal is **fast mitigation** and **learning without blame**.

#### Guidance
Define severity levels; assign roles (incident commander, communications, operations, subject-matter experts); keep a timeline; communicate externally when required; hold **blameless postmortems** for significant incidents; track corrective actions to completion.

#### Common mistakes
Searching for a person to blame; no timeline; vague or never-completed postmortem actions; incidents hidden from product/business stakeholders.

#### Sources
Google SRE book.

---

## 9. Continuous Delivery

#### Summary
Continuous delivery makes releasing software a reliable, low-risk, repeatable process.

#### Guidance
**Build once, promote the same artifact** across environments; automate tests and security checks; use deployment strategies (rolling, blue/green, canary, feature flags — [04 §14](04-web-application-design.md#14-deployment-strategies)); separate **deploy** from **release** where useful; make rollback or roll-forward reliable; track the **DORA metrics** (deployment frequency, lead time, change failure rate, time to restore).

#### Common mistakes
Manual production steps; different artifacts per environment; no database migration strategy; feature flags never removed.

#### Sources
DORA capabilities; Continuous Delivery (Humble & Farley); Twelve-Factor App.

---

## 10. Trunk-Based Development

#### Summary
Developers integrate small changes into a shared mainline frequently — fewer merge conflicts, faster feedback, and a foundation for continuous delivery.

#### Requirements / trade-offs
Requires strong automated tests and **feature flags or branch-by-abstraction** for incomplete work. Teams without CI/automated tests should first improve feedback loops.

#### Common mistakes
Long-lived branches mislabeled "trunk-based"; tolerating a broken main; feature flags without lifecycle management.

---

## 11. Feature Flags

#### Summary
Feature flags change behavior without redeploying code, decoupling deploy from release and enabling canaries, experiments, and kill switches.

#### Guidance
Classify flags — **release, experiment, ops, permission**; set owners and expiry dates; test both important states; avoid deeply nested flags; remove release flags after rollout.

#### Common mistakes
Permanent flags with no owner; security decisions implemented only as client-side flags; no record of flag changes during incidents.

#### Sources
Martin Fowler, *Feature Toggles*.

---

## 12. Database Change Management

#### Summary
Database changes must deploy safely alongside application changes.

#### Guidance
Version migrations in source control; use **expand-and-contract** for breaking schema changes ([04 §9.3](04-web-application-design.md#93-database-migrations)); make migrations idempotent where possible; test on realistic data volume; back up before destructive migrations; monitor migration duration and locks.

#### Common mistakes
Dropping columns still used by old app versions; long locks during peak traffic; no rollback plan for data transformations; manual schema changes not captured in code.

---

## 13. Twelve-Factor-Informed App Practices

The Twelve-Factor App ([04 §1](04-web-application-design.md#1-the-twelve-factor-app-and-beyond)) still influences modern service design: one codebase/many deploys; explicit dependencies; config in the environment/external stores; backing services as attached resources; separate build/release/run; stateless processes; port binding; concurrency via the process model; fast startup and graceful shutdown; dev/prod parity; logs as event streams; one-off admin processes.

**Trade-offs:** some practices need adaptation for desktop apps, serverless, or regulated environments; environment variables alone are insufficient for complex or secret configuration.

---

## 14. Cost Optimization

#### Summary
Cost is an architectural quality. Waste reduces product optionality.

#### Guidance
Track **unit economics** (cost per user, request, tenant, build, report, or job); right-size resources; use autoscaling where appropriate; remove idle resources; set budgets and alerts; design data retention intentionally; weigh managed services vs operational labor.

#### Common mistakes
Optimizing cloud spend while wasting engineering months; no owner for cost; no tags/labels for attribution; over-provisioning for rare peaks without business need.

---

## 15. Sustainability

#### Summary
Efficient software reduces infrastructure, energy, and device-resource waste — overlapping heavily with performance and cost.

#### Guidance
Reduce unnecessary computation and data movement; optimize storage retention; use efficient regions/services; prefer event-driven updates over polling; optimize frontend bundle size and CPU usage; shut down idle environments.

#### Trade-offs
Sustainability work should align with performance and cost priorities; some redundancy is necessary for reliability.

#### Sources
AWS/Azure/Google Cloud Well-Architected sustainability and cost guidance.

---

## Key Cross-References
- **Quality attributes & trade-offs:** [`06`](06-quality-attributes-tradeoffs.md).
- **Web & desktop security/observability:** [`04`](04-web-application-design.md), [`05`](05-desktop-application-design.md).
- **Delivery process & flags in code terms:** [`03` §11](03-software-design-principles.md#11-development-process--collaboration).
- **Operational checklists:** [`08`](08-checklists-and-templates.md) (security, reliability, operational readiness).
