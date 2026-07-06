# UI/UX Design Principles & Best Practices Reference

A comprehensive, **decision-oriented, citation-backed** reference on UI/UX design
principles and state-of-the-art best practices. It covers platform-agnostic
foundations — core principles, cognition, visual design, interaction design,
accessibility, forms, navigation, content — plus dedicated coverage of web, desktop,
mobile, and CLI/TUI platforms, domain patterns (e-commerce, SaaS/enterprise, AI
products), and an operational layer of agent checklists, pattern-selection matrices,
and acceptance criteria.

This repository is **documentation, not code.** There is no build or test suite —
only Markdown. It is designed to be **referenced from other repositories** while you
design, implement, or review user interfaces, and to be equally useful to human
engineers and AI coding agents. Every entry states **what** to do and **why**, with
evidence, concrete numbers, **When to use / When NOT to use** conditions, and
**Pros/Cons** where a pattern has genuine trade-offs.

> **Deep entry point:** [`00-index.md`](00-index.md) contains the full navigation —
> the agent operating model (conflict-resolution rules, risk levels, output format),
> task-based decision trees, an alphabetical topic lookup, the key-numbers cheat
> sheet, a glossary, and the annotated bibliography. Start there for anything beyond
> a quick orientation.

---

## File Map

| File | Scope |
|---|---|
| [`00-index.md`](00-index.md) | Navigation, agent operating model, decision trees, topic lookup, key numbers, glossary, bibliography |
| [`01-core-principles.md`](01-core-principles.md) | Nielsen's heuristics, Norman's principles, Gestalt, Laws of UX (Fitts, Hick, Doherty…), Shneiderman, ISO 9241-110 |
| [`02-cognitive-foundations.md`](02-cognitive-foundations.md) | Cognitive load, memory, recognition vs recall, progressive disclosure, mental models, biases, slips vs mistakes, flow |
| [`03-visual-design.md`](03-visual-design.md) | Hierarchy, typography, color (contrast, OKLCH), 8pt grid, density, elevation, iconography, dark mode, trends |
| [`04-interaction-design.md`](04-interaction-design.md) | Component states, response-time thresholds, modes, motion, loading, optimistic UI, gestures, keyboard, undo vs confirmation, error taxonomy |
| [`05-accessibility.md`](05-accessibility.md) | Legal case, WCAG 2.2 criteria with numbers, semantic HTML, ARIA principles, screen readers, COGA, testing methodology |
| [`06-aria-widget-reference.md`](06-aria-widget-reference.md) | ~24 APG widget patterns: roles, keyboard maps, failure modes; native-first matrix; acceptance checklist |
| [`07-forms-and-input.md`](07-forms-and-input.md) | Layout, labels, input types, validation timing, error design, defaults, autosave, auth UX (NIST, passkeys), high-risk forms |
| [`08-navigation-ia.md`](08-navigation-ia.md) | IA fundamentals, nav patterns, breadcrumbs, faceted filtering, search, sorting, command palette, scope safety, pagination, URLs |
| [`09-content-ux-writing.md`](09-content-ux-writing.md) | Plain language, voice & tone, microcopy, error/confirmation copy, capitalization, formatting, i18n/l10n |
| [`10-design-systems.md`](10-design-systems.md) | When (not) to build, design tokens, component APIs, documentation, test hooks, theming, versioning, governance, adoption |
| [`11-components-and-overlays.md`](11-components-and-overlays.md) | Per-component behavioral contracts; overlay layer-selection matrix; modal/dialog/drawer/popover/tooltip specs |
| [`12-data-tables-dashboards.md`](12-data-tables-dashboards.md) | Table vs grid selection, table a11y, sorting/filtering, grid editing, bulk actions, chart selection, dashboards |
| [`13-platform-web.md`](13-platform-web.md) | Responsive design, browser conventions, Core Web Vitals, SPA vs MPA, progressive enhancement, PWA, consent UX |
| [`14-platform-desktop.md`](14-platform-desktop.md) | Windows/macOS/Linux conventions, shortcuts, menus, command palettes, windowing, native vs Electron/Tauri, update UX |
| [`15-platform-mobile.md`](15-platform-mobile.md) | Thumb zone, touch targets, Material 3 / Apple HIG, back behavior, gestures, haptics, offline, tablets/foldables |
| [`16-platform-cli-tui.md`](16-platform-cli-tui.md) | Command structure, flags, help, stdout/stderr, NO_COLOR, --json, exit codes, TTY rules, XDG config, TUI patterns |
| [`17-ux-process-research.md`](17-ux-process-research.md) | Research methods, usability testing, 5-user rule, heuristic evaluation, A/B testing, UX metrics, ethics, review rubrics |
| [`18-patterns-antipatterns.md`](18-patterns-antipatterns.md) | Dark patterns (legal status, detection, alternatives), anti-patterns, trust/security notes, controversial patterns |
| [`19-domain-ecommerce-saas.md`](19-domain-ecommerce-saas.md) | Checkout, conversion ethics, enterprise/admin UX, settings, login patterns, roles, account deletion/export |
| [`20-ai-product-ux.md`](20-ai-product-ux.md) | AI interaction paradigms, risk tiers, prompt UX, output states, agentic-action rules, AI accessibility, chatbots |
| [`21-agent-checklists.md`](21-agent-checklists.md) | Decision checklists, pattern-selection matrices, failure-mode diagnostics, acceptance criteria, PR-review template |

The `00`–`21` numbering is a stable ordering; treat the filenames as durable anchors
for cross-references.

---

## How to Navigate

Start from the **task in front of you**, then load only what you need.

- **Orienting / finding a topic** → [`00-index.md`](00-index.md), especially its
  decision tree ("I'm working on X — what do I read?") and alphabetical topic lookup.
- **Building for a platform** → the platform file
  ([`13` web](13-platform-web.md), [`14` desktop](14-platform-desktop.md),
  [`15` mobile](15-platform-mobile.md), [`16` CLI/TUI](16-platform-cli-tui.md)),
  then [`04`](04-interaction-design.md) and [`05`](05-accessibility.md).
- **Designing a form / checkout / auth flow** → [`07`](07-forms-and-input.md) →
  [`19`](19-domain-ecommerce-saas.md) → [`09`](09-content-ux-writing.md).
- **Choosing or implementing a component** → [`21`](21-agent-checklists.md) matrices →
  [`11`](11-components-and-overlays.md) → [`06`](06-aria-widget-reference.md).
- **Accessibility work** → [`05`](05-accessibility.md) → [`06`](06-aria-widget-reference.md).
- **Reviewing a design or PR** → [`21`](21-agent-checklists.md) (diagnostics +
  acceptance criteria) → [`01`](01-core-principles.md) heuristics →
  [`18`](18-patterns-antipatterns.md) detection heuristics.
- **Planning research / validation** → [`17`](17-ux-process-research.md).

Every file is self-contained: it states its scope at the top and ends with a
machine-scannable **Quick reference** table — read that first for an actionable
summary of the whole file.

---

## Conventions

- **Entry template:** Definition → Reasoning/Evidence → When to use → When NOT to
  use → Pros → Cons → Implementation guidance (with concrete numbers) → optional
  Decision questions / Failure modes / Verification checklist → Related → Sources.
- **Conflict resolution** (documented in [`00-index.md`](00-index.md)): user task
  success and accessibility win over visual style; platform conventions usually win
  over brand; minimalism must not hide necessary controls; reversibility over novelty.
- **Risk levels:** low / medium / high — high-risk flows (legal, financial, medical,
  destructive, irreversible) require stronger prevention, confirmation, review, undo,
  and auditability.
- **Cross-links** are relative with GitHub-style kebab-case anchors
  (`[text](NN-file.md#heading-anchor)`).

---

## For AI Coding Agents

This repo is structured for retrieval:

1. **Begin at [`00-index.md`](00-index.md).** Its decision trees and topic lookup map
   a task or topic to the exact file and section; the key-numbers cheat sheet gives
   the most-cited hard values (contrast ratios, touch targets, response-time limits).
2. **Load only what you need** — the relevant deep-dive file(s), plus
   [`21-agent-checklists.md`](21-agent-checklists.md) for selection matrices and
   acceptance criteria on the final pass.
3. **Read *When to use / When NOT to use* before applying a pattern** — they encode
   the applicability conditions; *Cons* encodes what you trade away; *Verification
   checklist* encodes how to prove you applied it correctly.
4. **Cite the entries and their sources** when justifying a decision; do not invent
   citations or statistics.
5. **Checklists are decision support, not a substitute for testing with users.**
   AI-generated review verdicts require human expert review
   ([`17`](17-ux-process-research.md)).

Before editing this repository, read [`AGENTS.md`](AGENTS.md) for content conventions
and maintenance rules.
