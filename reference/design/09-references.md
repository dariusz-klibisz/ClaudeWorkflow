# References & Authoritative Sources

This reference synthesizes established, widely recognized sources in software architecture and engineering. The guide summarizes and integrates these sources; it does not reproduce them wholesale. Use them as **inputs to decisions, not absolute rules** — every source assumes a context.

*Last verified: 2026-06-25.*

---

## Architecture & Design

- **Martin Fowler** — *Software Architecture Guide*; *Patterns of Enterprise Application Architecture (PoEAA)*; *Microservices* (with James Lewis); articles on *PresentationDomainDataLayering*, *MonolithFirst*, *Strangler Fig / Patterns of Legacy Displacement*, *Micro Frontends*, *Serverless Architectures*, *Feature Toggles*, *GUI Architectures*. https://martinfowler.com/architecture/
- **Robert C. Martin (Uncle Bob)** — *Clean Architecture*; *Clean Code*; *Agile Software Development: Principles, Patterns, and Practices*; SOLID. https://en.wikipedia.org/wiki/SOLID
- **Eric Evans** — *Domain-Driven Design*; **Vaughn Vernon** — *Implementing Domain-Driven Design*.
- **Gang of Four** (Gamma, Helm, Johnson, Vlissides) — *Design Patterns*.
- **Neal Ford & Mark Richards** — *Fundamentals of Software Architecture*; *Software Architecture: The Hard Parts* (First/Second Laws of Software Architecture).
- **Neal Ford, Rebecca Parsons, Patrick Kua** — *Building Evolutionary Architectures* (fitness functions).
- **Alistair Cockburn** — *Hexagonal Architecture (Ports & Adapters)*. **Jeffrey Palermo** — *Onion Architecture*.
- **Simon Brown** — *The C4 Model*. https://c4model.com
- **Gregor Hohpe** — *The Software Architect Elevator*; *Enterprise Integration Patterns*.
- **Sam Newman** — *Building Microservices*; *Monolith to Microservices*.
- **Unmesh Joshi** — *Patterns of Distributed Systems*.
- **Melvin Conway** — *Conway's Law*; **Skelton & Pais** — *Team Topologies*.
- **Craig Larman** — *Applying UML and Patterns* (GRASP).
- **Meilir Page-Jones** — *connascence*; **Sandi Metz** — *Practical Object-Oriented Design* (AHA / "the wrong abstraction").
- **John Ousterhout** — *A Philosophy of Software Design* (deep vs shallow modules).
- **David Parnas** — *On the Criteria To Be Used in Decomposing Systems into Modules* (information hiding).

*Used for:* monolith vs microservices, layering, structural patterns, distributed-systems trade-offs, legacy modernization, UI architecture, evolutionary architecture, DDD, and code-level design principles.

---

## Architecture Patterns & Documentation

- **Azure Cloud Design Patterns:** https://learn.microsoft.com/en-us/azure/architecture/patterns/
- **microservices.io Pattern Language:** https://microservices.io/patterns/
- **Architecture Decision Records (ADR):** https://adr.github.io/
- **C4 Model:** https://c4model.com/

*Used for:* cache-aside, queues, CQRS, saga, transactional outbox, circuit breaker, rate limiting, deployment stamps, database-per-service, ADR rationale, and C4 context/container/component/code views.

---

## Cloud & Well-Architected Frameworks

- **Microsoft Azure Well-Architected Framework:** https://learn.microsoft.com/en-us/azure/well-architected/
- **AWS Well-Architected Framework:** https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html
- **Google Cloud Architecture / Well-Architected Framework:** https://cloud.google.com/architecture/framework

*Used for:* quality-attribute framing; the six pillars (operational excellence, security, reliability, performance efficiency, cost optimization, sustainability); cloud workload review structure; cross-pillar trade-offs.

---

## Process, Delivery & Reliability

- **Adam Wiggins / Heroku** — *The Twelve-Factor App*. https://12factor.net/
- **DORA Capabilities Catalog:** https://dora.dev/devops-capabilities/
- **Nicole Forsgren, Jez Humble, Gene Kim** — *Accelerate* (DORA metrics).
- **Jez Humble, David Farley** — *Continuous Delivery*.
- **Google** — *Site Reliability Engineering (SRE)* (SLI/SLO/SLA, error budgets, incident response, simplicity). https://sre.google/sre-book/table-of-contents/

*Used for:* continuous delivery, trunk-based development, deployment automation, observability, SLOs, toil reduction, incident response, configuration, logs, disposability, and cloud-native app practices.

---

## Distributed Systems Theory

- **Martin Kleppmann** — *Designing Data-Intensive Applications* (consistency, distributed data).
- **Eric Brewer** — *CAP Theorem*; **Daniel Abadi** — *PACELC*.
- **Deutsch & Gosling** — *Fallacies of Distributed Computing*.
- **Gene Amdahl** — *Amdahl's Law*; **Neil Gunther** — *Universal Scalability Law*.
- **Zhamak Dehghani** — *Data Mesh*.
- **Fred Brooks** — *No Silver Bullet* (essential vs accidental complexity); *The Mythical Man-Month*.
- **Ward Cunningham** — *Technical Debt* metaphor.

---

## Security & Secure Development

- **OWASP Top 10:2025** (current release): https://owasp.org/Top10/2025/ — primary web application risk categories. Previous: **OWASP Top 10:2021** https://owasp.org/Top10/2021/
- **OWASP Application Security Verification Standard (ASVS):** https://owasp.org/www-project-application-security-verification-standard/
- **OWASP Cheat Sheet Series:** https://cheatsheetseries.owasp.org/
- **NIST Secure Software Development Framework (SSDF), SP 800-218:** https://csrc.nist.gov/Projects/ssdf

*Used for:* secure SDLC, verification levels, application security requirements, vulnerability prevention, supply-chain security, secure environments, vulnerability response, and web app risk categories.

> **OWASP Top 10 note (2025 vs 2021):** Broken Access Control remains #1. **Security Misconfiguration** rose to #2; **Software Supply Chain Failures** (A03) and **Mishandling of Exceptional Conditions** (A10) are new/reframed categories; the 2021 standalone **SSRF** is folded into Broken Access Control. Rankings evolve roughly every few years — verify against the OWASP source for production decisions.

---

## Web Application Design

- **Google web.dev — Web Vitals:** https://web.dev/articles/vitals — Core Web Vitals and thresholds.
- **web.dev Learn / Performance / Accessibility:** https://web.dev/learn/
- **W3C WCAG (Web Content Accessibility Guidelines):** https://www.w3.org/WAI/standards-guidelines/wcag/
- **Kent C. Dodds** — *Testing Trophy*; **Mike Cohn** — *Test Pyramid*.

*Used for:* rendering strategies, performance, Core Web Vitals, accessibility (POUR), responsive design, privacy, progressive web apps, and testing strategy.

> **Core Web Vitals note (verified 2026-06-25):** the current stable Core Web Vitals are **LCP ≤ 2.5 s**, **INP ≤ 200 ms**, and **CLS ≤ 0.1**, assessed at the **75th percentile** across mobile and desktop. INP replaced First Input Delay (FID) as a Core Web Vital in 2024. Metrics evolve on an annual cadence; verify against web.dev.

---

## Desktop Application Design

- **Microsoft — Windows App Design / Windows App SDK / WinUI / MSIX / Authenticode:** https://learn.microsoft.com/en-us/windows/apps/
- **Apple — Human Interface Guidelines / App Sandbox / notarization / Keychain:** https://developer.apple.com/design/human-interface-guidelines
- **GNOME Human Interface Guidelines:** https://developer.gnome.org/hig/
- **KDE Human Interface Guidelines:** https://develop.kde.org/hig/
- **freedesktop.org** — XDG Base Directory spec; Flatpak / Snap / AppImage.
- **Flatpak Sandbox Permissions:** https://docs.flatpak.org/en/latest/sandbox-permissions.html
- **Electron Security:** https://www.electronjs.org/docs/latest/tutorial/security
- **Tauri Security:** https://v2.tauri.app/security/
- **Qt:** https://doc.qt.io — Widgets vs Qt Quick/QML, deployment.
- **Flutter desktop, .NET MAUI, Avalonia, Uno** — respective official docs.
- **Martin Fowler** — *GUI Architectures* (MVC/MVP/MVVM): https://martinfowler.com/eaaDev/uiArchs.html

*Used for:* platform design conventions, desktop UX, webview-based app security, trust boundaries, IPC, sandboxing, permissions, portals, code signing, update risk, and platform integration.

> *Note:* the Apple HIG page serves a JavaScript-only shell to automated fetchers. It remains an authoritative reference; this guide reflects general platform-design knowledge rather than verbatim Apple text.

---

## Source Evaluation Notes

The strongest sources are **vendor/operator frameworks and standards bodies** (Well-Architected frameworks, NIST SSDF, OWASP ASVS, WCAG, Google SRE) because they are maintained and reflect production experience across many systems — but they are intentionally broad. **Practitioner books and pattern catalogs** (Fowler, Newman, Kleppmann, Evans, GoF) are best for trade-off language and failure modes.

Use sources as **inputs to decisions, not absolute rules**. Every source assumes a context: a small internal tool, a regulated financial platform, a global consumer web app, and a local-first desktop app each need a different level of rigor.

Specific quantitative figures in this reference (Core Web Vitals thresholds, availability "nines," framework binary sizes, OWASP rankings) reflect commonly published industry values at the verification date and may evolve — **always confirm against the current primary source for production decisions.**
