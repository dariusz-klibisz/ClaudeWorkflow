# Unified UI/UX Design Reference — Index

A comprehensive, decision-oriented reference on UI/UX design principles and state-of-the-art best practices. Platform-agnostic foundations plus dedicated coverage of web, desktop, mobile, and CLI/TUI platforms, domain patterns (e-commerce, SaaS, AI products), and an operational layer of agent checklists and acceptance criteria. Written to be consumed by AI agents and developers making concrete design decisions.

**Corpus:** 22 files, ~14,500 lines. Standards baseline: WCAG 2.2, Material Design 3, current Apple HIG, Fluent 2, GNOME HIG, clig.dev, NIST SP 800-63B.

## How to use this reference

- **Every entry follows the same template:** Definition → Reasoning/Evidence → When to use → When NOT to use / exceptions → Pros → Cons → Implementation guidance (with concrete numbers) → optional Decision questions / Failure modes / Verification checklist → Related (cross-links) → Sources. Dark-pattern entries in file 18 substitute *Why it's harmful / Detection heuristic / Ethical alternative / Legal status*.
- **Every file ends with a "Quick reference" table** (topic | one-line guidance | key numbers) — scan these first for a compressed view of the whole file.
- **Cross-links are relative** (`[text](NN-file.md#kebab-anchor)`); anchors are GitHub-style kebab-case of headings.
- **When deciding whether to apply a pattern:** read its *When to use / When NOT to use* lists first — they encode the applicability conditions; the *Cons* list encodes what you trade away; the *Verification checklist* encodes how to prove you applied it correctly.
- **For implementation or review work**, go straight to [21-agent-checklists.md](21-agent-checklists.md) (selection matrices, failure-mode diagnostics, acceptance criteria) and follow its cross-links back into the deep-dive files.

## Agent operating model

Procedure for AI agents using this corpus to design, implement, or review UI:

1. Identify the user goal, business goal, risk level, platform, input methods, and accessibility constraints.
2. Read the relevant platform-agnostic principles first ([01](01-core-principles.md), [02](02-cognitive-foundations.md)).
3. Read the relevant platform-specific guidance second (files [13](13-platform-web.md)–[16](16-platform-cli-tui.md)).
4. Use the decision questions and matrices ([21](21-agent-checklists.md)) to choose a pattern.
5. Use the verification checklists and acceptance criteria before delivering or reviewing UI.
6. If guidance conflicts, prefer user task success, accessibility, platform conventions, and reversibility over novelty.

### Conflict-resolution rules

- **Accessibility vs visual style** — accessibility wins. If a visual treatment reduces contrast, removes visible focus, hides labels, or depends on motion, revise the visual treatment.
- **Platform convention vs brand consistency** — platform convention usually wins for system-level controls, navigation, text input, keyboard shortcuts, window behavior, and permission prompts. Brand can shape color, illustration, tone, and density where it doesn't violate user expectations.
- **Minimalism vs discoverability** — minimalism means removing irrelevant complexity, not hiding necessary controls. If users cannot find or understand required actions, the design is too minimal.
- **Efficiency vs learnability** — support both: visible primary paths for novices, accelerators for experts. Never make expert-only controls the only path for critical work.
- **Personalization vs predictability** — make personalized regions recognizable, reversible, and explainable when consequences matter.

### Risk levels

- **Low:** browsing, reading, preference changes, reversible filtering.
- **Medium:** account settings, publishing, saved edits, low-cost purchases.
- **High:** legal, medical, financial, identity, destructive actions, public posting, privacy exposure, irreversible data changes. High-risk flows require stronger error prevention, confirmation, review, undo, auditability, plain language, and accessibility validation ([07 › High-stakes forms](07-forms-and-input.md), [21 › risk escalation](21-agent-checklists.md)).

### Agent output format for UI decisions

```md
Decision: ...
User goal: ...
Platform: ...
Pattern chosen: ...
Why: ...
Tradeoffs: ...
Accessibility checks: ...
Failure modes: ...
Verification: ...
```

### Minimum review checklist

- Can users tell where they are, what changed, and what to do next?
- Are primary and secondary actions visually and semantically distinct?
- Are destructive or irreversible actions protected by prevention, confirmation, or undo?
- Can the flow be completed with keyboard, screen reader, touch, and pointer where applicable?
- Are labels persistent and explicit (no placeholder-only instructions)?
- Are errors plain-language, specific, near the problem, and actionable?
- Does the design handle loading, empty, error, offline, permission-denied, and partial-success states?
- Does the design survive text resizing, dark mode, reduced motion, localization, and narrow viewports?

Caveat: checklists and AI-generated review verdicts are decision support, not a substitute for testing with users; human expert review is required before shipping verdicts ([17 › Design review rubrics](17-ux-process-research.md)).

## File map

| File | Title | Contents |
|---|---|---|
| [01-core-principles.md](01-core-principles.md) | Core UI/UX Principles | Nielsen's 10 heuristics, Norman's principles (affordances, signifiers, mapping, feedback, constraints, conceptual models), 7 Gestalt principles, 15 Laws of UX (Fitts, Hick, Jakob, Miller, Tesler, Doherty, Peak-End…), Shneiderman's 8 Golden Rules, ISO 9241-110 |
| [02-cognitive-foundations.md](02-cognitive-foundations.md) | Cognitive Foundations | Cognitive load, working memory & chunking, recognition vs recall, progressive disclosure, mental models, banner blindness, scanning patterns, change blindness, habituation, choice overload, satisficing, System 1/2, biases, slips vs mistakes, mode errors, post-completion errors, flow, interruption cost |
| [03-visual-design.md](03-visual-design.md) | Visual Design | Visual hierarchy, typography, color (palettes, semantic color, contrast, color blindness, OKLCH), spacing (8pt grid), density modes, elevation/shadows, border radius, iconography, imagery, dark mode, trend analysis |
| [04-interaction-design.md](04-interaction-design.md) | Interaction Design | Affordances in practice, direct manipulation, component states, response-time thresholds (0.1/1/10s), modes & mode errors, microinteractions, motion, loading strategies, optimistic UI, gestures, drag & drop, keyboard patterns, undo vs confirmation, error taxonomy & feedback, validation timing, scroll, empty states, notifications, performance measurement (field vs lab) |
| [05-accessibility.md](05-accessibility.md) | Accessibility & Inclusive Design | Legal/business case (ADA, EAA, EN 301 549), WCAG 2.2 criteria with numbers (contrast, target size, focus, keyboard, reflow, dragging, redundant entry, accessible auth incl. AAA 2.4.12/3.3.9), semantic HTML, ARIA principles, screen readers, forms a11y, reduced motion, COGA, testing methodology + manual test scripts |
| [06-aria-widget-reference.md](06-aria-widget-reference.md) | ARIA & Semantic Widget Reference | ~24 APG widget patterns with roles, required states, keyboard maps, failure modes; native-first matrix; ARIA risk rules; custom-widget acceptance checklist |
| [07-forms-and-input.md](07-forms-and-input.md) | Forms & Data Entry | Layout, labels, required/optional, multi-step vs single, input types, date input, masking, file upload, validation timing, error design, destructive confirmation, defaults, autosave, password/auth UX (NIST, passkeys), high-risk & regulated forms, payment/identity inputs |
| [08-navigation-ia.md](08-navigation-ia.md) | Information Architecture & Navigation | IA fundamentals (schemes, depth vs breadth, card sorting, tree testing, information scent), nav patterns, breadcrumbs, faceted filtering, search (autocomplete, scoped, no-results), sorting, command palette, scope/environment safety, pagination vs infinite scroll, URL design, menu ergonomics |
| [09-content-ux-writing.md](09-content-ux-writing.md) | Content & UX Writing | Plain language, reading levels, voice & tone, terminology, inclusive language, button labels, link text, error/confirmation/notification copy, empty states, capitalization, numbers/dates/units, truncation, scannability, i18n/l10n (expansion, ICU plurals, RTL, data-format matrix, culturalization) |
| [10-design-systems.md](10-design-systems.md) | Design Systems | When (not) to build, design tokens (tiers, naming, theming, DTCG), atomic design critique, component API design, composition vs configuration, documentation, stable test hooks, theming, versioning & deprecation, governance, adoption metrics, reference systems compared |
| [11-components-and-overlays.md](11-components-and-overlays.md) | Components & Overlays | Per-component behavioral contracts (buttons, links, inputs, checkbox, radio, select, combobox, switch, menus, tabs, accordion, cards, alerts w/ severity mapping, progress indicators, empty states); overlay layer-selection matrix; modal/alert-dialog/drawer/popover/tooltip specs; z-index & focus rules |
| [12-data-tables-dashboards.md](12-data-tables-dashboards.md) | Data Tables, Grids & Dashboards | Table vs cards vs grid selection, table semantics/a11y, sorting/filtering/column management, grid editing models, bulk actions, density rules, metric cards, chart-selection matrix, viz accessibility, dashboard layout |
| [13-platform-web.md](13-platform-web.md) | Platform: Web | Responsive design, viewport quirks, browser conventions, multi-tab, session expiry, Core Web Vitals (field vs lab), perceived performance, web typography, URLs, SPA vs MPA, progressive enhancement, PWA, scroll UX, consent/cookie UX, permission prompts, SEO-UX, agent-friendly content |
| [14-platform-desktop.md](14-platform-desktop.md) | Platform: Desktop | Windows/macOS/Linux conventions + shortcut mapping table, menus, context menus, command palettes, window management, modality, density, desktop drag & drop, native vs Electron/Tauri, system integration, hi-DPI/multi-monitor, offline/autosave, update UX, first-run |
| [15-platform-mobile.md](15-platform-mobile.md) | Platform: Mobile | Thumb zone, touch targets (44pt/48dp/24px), Material 3, Apple HIG, navigation patterns, back behavior, deep links, gestures, mobile forms, haptics, orientation & safe areas, notifications, onboarding, offline, performance, tablets/foldables, responsive vs adaptive, cross-platform comparison matrix |
| [16-platform-cli-tui.md](16-platform-cli-tui.md) | Platform: CLI & TUI | Human-first philosophy, command structure, flags, help design, stdout/stderr contract, color (NO_COLOR), --json output, progress, errors & exit codes, TTY interactivity rules, config precedence (XDG), safety (dry-run), composability, TUI patterns |
| [17-ux-process-research.md](17-ux-process-research.md) | UX Process & Research | Method landscape, interviews, surveys, personas, JTBD, journey mapping, usability testing, 5-user rule, heuristic evaluation, cognitive walkthrough, A/B testing, analytics, UX metrics (SUS/SEQ/NPS/HEART), research ethics, GDS phases, prototyping, double diamond, design sprints, critique, design debt, prioritization, AI-assisted research, design review rubrics |
| [18-patterns-antipatterns.md](18-patterns-antipatterns.md) | Patterns to Avoid & Controversial Patterns | Deceptive–neutral–persuasive spectrum + verification questions, 15 dark patterns (legal status, detection heuristics, ethical alternatives), 20 anti-patterns, security/privacy/public-service trust notes, 8 controversial patterns with balanced trade-offs |
| [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md) | Domain: E-commerce, SaaS & Accounts | Listing/PDP/cart/checkout, conversion ethics, enterprise/admin UX (scope visibility, permission-aware UI, bulk safeguards, audit), settings IA + save-behavior models, login patterns, permission prompts, roles, account deletion/export |
| [20-ai-product-ux.md](20-ai-product-ux.md) | AI Product UX | Interaction paradigms (GUI vs prompting vs agent), HAX-style lifecycle, AI risk tiers, prompt UX, output states (streaming, citations, uncertainty), agentic-action rules (plan preview, activity log, stop), AI accessibility, support chatbots (liability) |
| [21-agent-checklists.md](21-agent-checklists.md) | Agent Checklists & Acceptance Criteria | Decision checklists, ~40-row pattern-selection matrices, symptom→fix failure-mode diagnostics, testable acceptance criteria per platform/domain, PR-review template, principle tradeoff tables |

## Decision tree: "I'm working on X — what do I read?"

| Task | Read (in order) |
|---|---|
| Starting any new product/feature design | [01](01-core-principles.md) → [02](02-cognitive-foundations.md) → relevant platform file → [17](17-ux-process-research.md) for validation plan |
| Building a web app | [13-platform-web.md](13-platform-web.md) → [04](04-interaction-design.md) → [05](05-accessibility.md) → [08](08-navigation-ia.md) |
| Building a desktop app | [14-platform-desktop.md](14-platform-desktop.md) → [04](04-interaction-design.md) → [03](03-visual-design.md#density-modes) |
| Building a mobile app | [15-platform-mobile.md](15-platform-mobile.md) → [04](04-interaction-design.md) → [08](08-navigation-ia.md) |
| Building a CLI tool | [16-platform-cli-tui.md](16-platform-cli-tui.md) (self-contained) → [01](01-core-principles.md#postels-law) |
| Choosing/implementing a UI component | [21](21-agent-checklists.md) matrices → [11-components-and-overlays.md](11-components-and-overlays.md) → [06-aria-widget-reference.md](06-aria-widget-reference.md) |
| Designing a form / checkout / signup | [07-forms-and-input.md](07-forms-and-input.md) → [19](19-domain-ecommerce-saas.md) → [09](09-content-ux-writing.md) → [05](05-accessibility.md#forms-accessibility) |
| Choosing navigation structure | [08-navigation-ia.md](08-navigation-ia.md) → platform file's navigation section |
| Picking colors / type / spacing | [03-visual-design.md](03-visual-design.md) → [05](05-accessibility.md#contrast-minimum) → [10](10-design-systems.md#design-tokens) |
| Adding animation / loading states | [04-interaction-design.md](04-interaction-design.md#motion-and-animation) → [05](05-accessibility.md#reduced-motion) |
| Writing UI copy / error messages | [09-content-ux-writing.md](09-content-ux-writing.md) → [07](07-forms-and-input.md#error-message-design) |
| Accessibility compliance work | [05-accessibility.md](05-accessibility.md) → [06](06-aria-widget-reference.md) → [17](17-ux-process-research.md#accessibility-auditing-in-process) |
| Starting/evaluating a design system | [10-design-systems.md](10-design-systems.md) → [11](11-components-and-overlays.md) → [03](03-visual-design.md) |
| Dialogs / modals / tooltips | [11](11-components-and-overlays.md) → [06](06-aria-widget-reference.md) → [18](18-patterns-antipatterns.md#modal-overuse) |
| Data tables / dashboards | [12-data-tables-dashboards.md](12-data-tables-dashboards.md) → [03](03-visual-design.md#density-modes) → [14](14-platform-desktop.md#density-and-information-rich-layouts) |
| Search / filtering / sorting | [08](08-navigation-ia.md#search) → [12](12-data-tables-dashboards.md) → [07](07-forms-and-input.md#search-inputs) |
| Dark mode implementation | [03](03-visual-design.md#dark-mode-design) → [10](10-design-systems.md) → [15](15-platform-mobile.md#dark-mode-on-mobile) |
| Enterprise/admin UX, multi-tenant scope | [19](19-domain-ecommerce-saas.md) → [08](08-navigation-ia.md) scope safety → [12](12-data-tables-dashboards.md) bulk actions |
| Settings / preferences | [19](19-domain-ecommerce-saas.md) settings section → [08](08-navigation-ia.md) → [07](07-forms-and-input.md#autosave-and-draft-recovery) |
| Authentication / passwords / accounts | [07](07-forms-and-input.md#password-and-auth-ux) → [19](19-domain-ecommerce-saas.md) accounts section → [05](05-accessibility.md#accessible-authentication) |
| Designing an AI feature / agent UX | [20-ai-product-ux.md](20-ai-product-ux.md) → [18](18-patterns-antipatterns.md) → [21](21-agent-checklists.md) |
| Planning research / testing a design | [17-ux-process-research.md](17-ux-process-research.md#the-research-method-landscape) → method-specific entries |
| Reviewing a design for problems | [21](21-agent-checklists.md) diagnostics → [01](01-core-principles.md) heuristics → [18](18-patterns-antipatterns.md) detection heuristics → [17](17-ux-process-research.md#heuristic-evaluation) |
| Reviewing a PR that touches UI | [21](21-agent-checklists.md) acceptance criteria + PR template → relevant deep-dive files |
| Ethics / legal review of persuasion mechanics | [18-patterns-antipatterns.md](18-patterns-antipatterns.md) → [02](02-cognitive-foundations.md#cognitive-biases-in-ux) → [19](19-domain-ecommerce-saas.md) conversion ethics |
| Performance as UX | [13](13-platform-web.md#core-web-vitals) → [04](04-interaction-design.md#feedback-and-response-time-thresholds) → [15](15-platform-mobile.md) |
| Notifications / interruptions | [04](04-interaction-design.md#notifications-and-interruption-design) → [15](15-platform-mobile.md#notifications-ux) → [02](02-cognitive-foundations.md#interruption-cost) |
| Onboarding / first-run / empty states | [15](15-platform-mobile.md#onboarding) or [14](14-platform-desktop.md#first-run-experience) → [04](04-interaction-design.md#empty-states) → [09](09-content-ux-writing.md) |
| Destructive actions / undo | [04](04-interaction-design.md#undoredo-vs-confirmation-dialogs) → [07](07-forms-and-input.md#destructive-action-confirmation) → [11](11-components-and-overlays.md) alert dialogs |
| Internationalization / localization | [09](09-content-ux-writing.md) i18n section → [15](15-platform-mobile.md) localization across layouts → [21](21-agent-checklists.md) i18n criteria |

## Topic lookup

Alphabetical map of common topics to their primary entry (secondary locations in parentheses).

| Topic | Location |
|---|---|
| 8pt grid | [03 › The 8pt spacing grid](03-visual-design.md#the-8pt-spacing-grid) |
| A/B testing | [17 › A/B testing](17-ux-process-research.md#ab-testing) |
| Affordances / signifiers | [01 › Affordances](01-core-principles.md#affordances), [04](04-interaction-design.md#affordances-and-signifiers-in-practice) |
| AI risk tiers / agentic actions | [20](20-ai-product-ux.md) |
| ARIA (principles) | [05 › ARIA](05-accessibility.md#aria); per-widget: [06](06-aria-widget-reference.md) |
| Autofill / autocomplete attributes | [07 › Defaults and prefill](07-forms-and-input.md#defaults-and-prefill) |
| Autosave / drafts | [07 › Autosave](07-forms-and-input.md#autosave-and-draft-recovery) ([14](14-platform-desktop.md#offline-first-and-autosave)) |
| Back button behavior | [13 › Browser conventions](13-platform-web.md#browser-conventions), [15 › Back behavior](15-platform-mobile.md#back-behavior) |
| Breadcrumbs | [08 › Breadcrumbs](08-navigation-ia.md#breadcrumbs) |
| Bulk actions | [12 › Bulk actions](12-data-tables-dashboards.md), [19](19-domain-ecommerce-saas.md) |
| Buttons: labels | [09 › Button labels](09-content-ux-writing.md#button-labels) |
| Buttons: links vs buttons | [11](11-components-and-overlays.md), [13 › Links vs buttons](13-platform-web.md#links-vs-buttons) |
| Card sorting / tree testing | [08 › IA fundamentals](08-navigation-ia.md#ia-fundamentals) |
| Carousels | [18 › Auto-rotating](18-patterns-antipatterns.md#auto-rotating-carousels) |
| Charts / data visualization | [12 › Dashboards](12-data-tables-dashboards.md) |
| Checkout / cart | [19](19-domain-ecommerce-saas.md), [07](07-forms-and-input.md) |
| Cognitive load | [02 › Cognitive load theory](02-cognitive-foundations.md#cognitive-load-theory) |
| Color palettes / color blindness | [03 › Building a UI palette](03-visual-design.md#building-a-ui-palette), [03](03-visual-design.md#color-blindness-and-color-independence), [05](05-accessibility.md#color-independence) |
| Command palettes | [14 › Command palettes](14-platform-desktop.md#command-palettes), [08](08-navigation-ia.md) |
| Component states | [04 › Component states](04-interaction-design.md#component-states), [11](11-components-and-overlays.md) |
| Confirmation dialogs | [04 › Undo/redo vs confirmation](04-interaction-design.md#undoredo-vs-confirmation-dialogs), [11](11-components-and-overlays.md) |
| Consent / cookie banners | [13 › Consent and cookie UX](13-platform-web.md#consent-and-cookie-ux) |
| Contrast ratios | [05 › Contrast minimum](05-accessibility.md#contrast-minimum), [03](03-visual-design.md#contrast) |
| Core Web Vitals | [13 › Core Web Vitals](13-platform-web.md#core-web-vitals) |
| Dark mode | [03 › Dark mode design](03-visual-design.md#dark-mode-design) |
| Dark patterns | [18 › Dark patterns catalog](18-patterns-antipatterns.md#dark-patterns-catalog) |
| Data tables / data grids | [12](12-data-tables-dashboards.md) |
| Date pickers | [07 › Date input](07-forms-and-input.md#date-input) |
| Deep links | [15 › Deep links](15-platform-mobile.md#deep-links), [13 › URLs as UX](13-platform-web.md#urls-as-ux) |
| Defaults | [07 › Defaults and prefill](07-forms-and-input.md#defaults-and-prefill), [02](02-cognitive-foundations.md#cognitive-biases-in-ux) |
| Density (compact/comfortable) | [03 › Density modes](03-visual-design.md#density-modes), [12](12-data-tables-dashboards.md) |
| Design tokens | [10 › Design tokens](10-design-systems.md#design-tokens) |
| Disabled buttons | [11](11-components-and-overlays.md), [18](18-patterns-antipatterns.md) |
| Drag and drop | [04 › Drag and drop](04-interaction-design.md#drag-and-drop), [14](14-platform-desktop.md#desktop-drag-and-drop) |
| Empty states | [04 › Empty states](04-interaction-design.md#empty-states), [09](09-content-ux-writing.md), [11](11-components-and-overlays.md) |
| Environment/tenant scope safety | [08 › Scope and context safety](08-navigation-ia.md), [19](19-domain-ecommerce-saas.md) |
| Error messages (writing) | [09 › Error message writing](09-content-ux-writing.md#message-design) |
| Error messages (forms) / error taxonomy | [07 › Error message design](07-forms-and-input.md#error-message-design), [04](04-interaction-design.md) |
| Errors: slips vs mistakes | [02 › Slips vs mistakes](02-cognitive-foundations.md#slips-vs-mistakes) |
| Exit codes | [16 › Exit codes](16-platform-cli-tui.md#exit-codes) |
| F-pattern / scanning | [02 › Scanning patterns](02-cognitive-foundations.md#scanning-patterns) |
| Fitts's Law | [01 › Fitts's Law](01-core-principles.md#fittss-law) |
| Focus indicators | [05 › Focus visible](05-accessibility.md#focus-visible) |
| Gestalt principles | [01 › Gestalt principles](01-core-principles.md#gestalt-principles) |
| Gestures | [04 › Touch and pointer gestures](04-interaction-design.md#touch-and-pointer-gestures), [15](15-platform-mobile.md#gestures) |
| Hamburger menu | [08 › Hamburger menu](08-navigation-ia.md#hamburger-menu), [18](18-patterns-antipatterns.md#hamburger-menu) |
| Haptics | [15 › Haptics](15-platform-mobile.md#haptics) |
| Hick's Law / choice overload | [01 › Hick's Law](01-core-principles.md#hicks-law), [02](02-cognitive-foundations.md#decision-fatigue-and-choice-overload) |
| Icons | [03 › Iconography](03-visual-design.md#iconography) |
| Infinite scroll vs pagination | [08 › Content traversal](08-navigation-ia.md#content-traversal), [18](18-patterns-antipatterns.md#infinite-scroll) |
| Internationalization (i18n) | [09 › i18n/l10n](09-content-ux-writing.md#global-content) |
| Keyboard shortcuts | [14 › Keyboard shortcuts design](14-platform-desktop.md#keyboard-shortcuts-design) |
| Keyboard navigation / focus management | [04](04-interaction-design.md#keyboard-interaction-patterns), [05](05-accessibility.md#keyboard-navigation), [06](06-aria-widget-reference.md) |
| Label placement (forms) | [07 › Label placement](07-forms-and-input.md#label-placement) |
| Loading: spinners/skeletons/progress | [04 › Loading strategies](04-interaction-design.md#loading-strategies), [11](11-components-and-overlays.md) |
| Login / SSO / step-up auth | [19](19-domain-ecommerce-saas.md), [07](07-forms-and-input.md#password-and-auth-ux) |
| Mental models | [02 › Mental models](02-cognitive-foundations.md#mental-models-and-conceptual-models) |
| Menus (desktop) | [14 › Menus](14-platform-desktop.md#menus), [08](08-navigation-ia.md#menu-ergonomics) |
| Metric cards / dashboards | [12](12-data-tables-dashboards.md) |
| Microinteractions | [04 › Microinteractions](04-interaction-design.md#microinteractions) |
| Modals / dialogs | [11](11-components-and-overlays.md), [14](14-platform-desktop.md#modality), [18](18-patterns-antipatterns.md#modal-overuse) |
| Modes / mode errors | [04](04-interaction-design.md), [02](02-cognitive-foundations.md#mode-errors) |
| Motion / animation durations | [04 › Motion and animation](04-interaction-design.md#motion-and-animation) |
| Nielsen's heuristics | [01 › Nielsen's 10 usability heuristics](01-core-principles.md#nielsens-10-usability-heuristics) |
| Notifications / push | [04](04-interaction-design.md#notifications-and-interruption-design), [15](15-platform-mobile.md#notifications-ux) |
| Offline design | [15 › Offline](15-platform-mobile.md#offline-and-poor-connectivity), [14](14-platform-desktop.md#offline-first-and-autosave) |
| Onboarding | [15 › Onboarding](15-platform-mobile.md#onboarding), [14](14-platform-desktop.md#first-run-experience) |
| Optimistic UI | [04 › Optimistic UI](04-interaction-design.md#optimistic-ui) |
| Orientation (mobile) | [15](15-platform-mobile.md) |
| Passwords / passkeys | [07 › Password and auth UX](07-forms-and-input.md#password-and-auth-ux) |
| Permission prompts | [13](13-platform-web.md), [15](15-platform-mobile.md), [19](19-domain-ecommerce-saas.md) |
| Personas / JTBD | [17 › Personas](17-ux-process-research.md#personas), [17 › JTBD](17-ux-process-research.md#jobs-to-be-done) |
| Placeholders | [07 › Labels vs placeholders](07-forms-and-input.md#labels-vs-placeholders) |
| Plain language | [09 › Plain language](09-content-ux-writing.md#language-foundations) |
| Progressive disclosure | [02 › Progressive disclosure](02-cognitive-foundations.md#progressive-disclosure) |
| Progressive enhancement / PWA | [13](13-platform-web.md#progressive-enhancement), [13](13-platform-web.md#pwa-capabilities) |
| Reduced motion | [05 › Reduced motion](05-accessibility.md#reduced-motion) |
| Response time limits (0.1/1/10s) | [04 › Feedback and response time thresholds](04-interaction-design.md#feedback-and-response-time-thresholds) |
| Responsive design / breakpoints | [13 › Responsive design](13-platform-web.md#responsive-design), [15](15-platform-mobile.md) |
| Screen readers / alt text | [05 › Screen readers](05-accessibility.md#screen-readers) |
| Search | [08 › Search](08-navigation-ia.md#search) |
| Settings / preferences | [19](19-domain-ecommerce-saas.md) |
| Sorting | [08 › Sorting](08-navigation-ia.md), [12](12-data-tables-dashboards.md) |
| SPA vs MPA | [13 › SPA vs MPA](13-platform-web.md#spa-vs-mpa) |
| SUS / UX metrics | [17 › UX metrics frameworks](17-ux-process-research.md#ux-metrics-frameworks) |
| Tabs (navigation) | [08 › Tabs](08-navigation-ia.md#tabs), [11](11-components-and-overlays.md) |
| Target sizes (touch) | [05 › Target size](05-accessibility.md#target-size), [15](15-platform-mobile.md#touch-targets) |
| Test hooks (`data-testid`) | [10 › Stable test hooks](10-design-systems.md#stable-test-hooks) |
| Toasts / snackbars | [04 › Notifications](04-interaction-design.md#notifications-and-interruption-design), [09](09-content-ux-writing.md) |
| Tooltips | [11](11-components-and-overlays.md), [04](04-interaction-design.md#hover-dependence), [06](06-aria-widget-reference.md) |
| Typography | [03 › Typography](03-visual-design.md#typography) |
| Undo / redo | [04 › Undo/redo vs confirmation](04-interaction-design.md#undoredo-vs-confirmation-dialogs) |
| URLs | [13 › URLs as UX](13-platform-web.md#urls-as-ux) |
| Usability testing / 5-user rule | [17 › Usability testing](17-ux-process-research.md#usability-testing), [17](17-ux-process-research.md#the-5-user-rule) |
| Validation timing | [07 › Inline validation timing](07-forms-and-input.md#inline-validation-timing) |
| WCAG 2.2 | [05 › WCAG 2.2 structure](05-accessibility.md#wcag-22-structure-and-conformance-levels) |
| Wizards / multi-step | [07 › Multi-step vs single-page](07-forms-and-input.md#multi-step-vs-single-page-forms), [08](08-navigation-ia.md) |

## Key numbers cheat sheet

The most-cited hard numbers across the corpus:

| Number | Meaning | Source entry |
|---|---|---|
| 0.1s / 1s / 10s | Response-time limits: instant / flow / attention | [04](04-interaction-design.md#feedback-and-response-time-thresholds) |
| 400ms | Doherty threshold for productivity loops | [01](01-core-principles.md#doherty-threshold) |
| 4.5:1 / 3:1 / 7:1 | WCAG contrast: text AA / large-text & UI AA / AAA | [05](05-accessibility.md#contrast-minimum) |
| 44pt / 48dp / 24px | Touch targets: Apple / Material / WCAG floor | [15](15-platform-mobile.md#touch-targets) |
| 45–75ch, 1.5 | Line length and body line-height | [03](03-visual-design.md#line-height-and-line-length) |
| 16px | Minimum web body font (also prevents iOS zoom) | [03](03-visual-design.md#type-scale) |
| 320px / 400% | WCAG reflow width / zoom | [05](05-accessibility.md#reflow) |
| LCP 2.5s / INP 200ms / CLS 0.1 | Core Web Vitals "good" thresholds (75th percentile, field data) | [13](13-platform-web.md#core-web-vitals) |
| 100–300ms (≤500ms) | Functional animation durations | [04](04-interaction-design.md#motion-and-animation) |
| 4±1 | Working-memory chunks (modern estimate) | [02](02-cognitive-foundations.md#working-memory-and-chunking) |
| ~5 users | Formative usability-test round size | [17](17-ux-process-research.md#the-5-user-rule) |
| SUS 68 | Average usability score (50th percentile) | [17](17-ux-process-research.md#ux-metrics-frameworks) |
| 3–5 / 5–7 | Bottom-nav tabs / top-level nav items | [15](15-platform-mobile.md), [01](01-core-principles.md#hicks-law) |
| ~8% | Men (N. European descent) with color vision deficiency | [03](03-visual-design.md#color-blindness-and-color-independence) |
| ~57% / ~30–40% | A11y issues automated tools catch (Deque, weighted / issue-type studies) | [05](05-accessibility.md#testing-methodology) |
| #121212 | Dark-mode base surface (not pure black) | [03](03-visual-design.md#dark-mode-design) |
| 30 days | Soft-delete/trash retention convention | [04](04-interaction-design.md#undoredo-vs-confirmation-dialogs) |
| +30–50% / +100–300% | i18n text expansion: paragraphs / short labels | [09](09-content-ux-writing.md#global-content) |

## Glossary

- **A11y** — accessibility. **AT** — assistive technology (screen readers, switch devices).
- **Affordance / Signifier** — what an object allows / the perceivable cue signaling it ([01](01-core-principles.md#normans-design-principles)).
- **APG** — W3C ARIA Authoring Practices Guide (widget keyboard/ARIA patterns) ([06](06-aria-widget-reference.md)).
- **CLS / LCP / INP** — Core Web Vitals: layout shift / largest paint / interaction latency.
- **Cognitive load** — total working-memory demand a task imposes ([02](02-cognitive-foundations.md#cognitive-load-theory)).
- **Dark pattern (deceptive pattern)** — design manipulating users against their interest ([18](18-patterns-antipatterns.md)).
- **Design tokens** — named, tiered design decisions (color.primary, space.4) ([10](10-design-systems.md#design-tokens)).
- **DTCG** — W3C Design Tokens Community Group (token format standard).
- **EAA / EN 301 549** — European Accessibility Act (in force 2019, compliance deadline June 28, 2025) and its harmonized technical standard.
- **Focus order / focus trap** — sequence of keyboard focus; a trap confines focus (correct in modals, a bug elsewhere) ([05](05-accessibility.md), [11](11-components-and-overlays.md)).
- **HAX** — Microsoft's Human-AI eXperience guidelines ([20](20-ai-product-ux.md)).
- **HIG** — Human Interface Guidelines (Apple's, GNOME's, generically any platform's).
- **IA** — information architecture ([08](08-navigation-ia.md)).
- **JTBD** — Jobs-to-be-Done framework ([17](17-ux-process-research.md#jobs-to-be-done)).
- **M3** — Material Design 3 (Google/Android design system).
- **Mental model** — a user's internal explanation of how a system works ([02](02-cognitive-foundations.md#mental-models-and-conceptual-models)).
- **MPA / SPA** — multi-page vs single-page web application architecture ([13](13-platform-web.md#spa-vs-mpa)).
- **POUR** — WCAG principles: Perceivable, Operable, Understandable, Robust.
- **Progressive disclosure** — show essentials, reveal the rest on demand ([02](02-cognitive-foundations.md#progressive-disclosure)).
- **RUM** — real user monitoring (field performance data, vs lab) ([13](13-platform-web.md#core-web-vitals)).
- **SC** — WCAG Success Criterion (e.g., SC 1.4.3).
- **SR** — screen reader. **SUS** — System Usability Scale ([17](17-ux-process-research.md#ux-metrics-frameworks)).
- **TTY** — terminal session (interactive); TTY-detection gates CLI interactivity ([16](16-platform-cli-tui.md#interactivity-rules)).
- **TUI** — text/terminal user interface (full-screen, e.g., htop) ([16](16-platform-cli-tui.md#full-screen-tuis)).
- **VPAT/ACR** — Voluntary Product Accessibility Template / Accessibility Conformance Report.
- **WCAG** — Web Content Accessibility Guidelines (2.2 current) ([05](05-accessibility.md)).
- **XDG** — freedesktop.org base-directory standard for config/cache paths ([16](16-platform-cli-tui.md#configuration-precedence)).

## Bibliography (primary sources, annotated)

Consolidated list of the most-cited sources across all files, with authority notes. Per-entry citations appear inline in each file.

### Standards and guidelines

- **WCAG 2.2** — https://www.w3.org/TR/WCAG22/ (quick ref: https://www.w3.org/WAI/WCAG22/quickref/; per-criterion: https://www.w3.org/WAI/WCAG22/Understanding/). *Normative accessibility standard; AA is the practical legal baseline.*
- **ARIA Authoring Practices Guide** — https://www.w3.org/WAI/ARIA/apg/patterns/. *De facto widget keyboard/ARIA behavior spec; guidance, not normative.* **Using ARIA** — https://www.w3.org/TR/using-aria/.
- **W3C COGA: Making Content Usable** — https://www.w3.org/TR/coga-usable/. *Cognitive accessibility beyond WCAG.*
- **EN 301 549** — https://www.etsi.org/deliver/etsi_en/301500_301599/301549/ · **ISO 9241-110:2020** — https://www.iso.org/standard/75258.html
- **NIST SP 800-63B** (authentication) — https://pages.nist.gov/800-63-3/sp800-63b.html. *Overrides folk password wisdom (composition rules, rotation).*
- HTML autocomplete spec — https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#autofill · XDG Base Directory — https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html · Semantic Versioning — https://semver.org/ · NO_COLOR — https://no-color.org/ · JSON Lines — https://jsonlines.org/

### Platform design systems and guidelines

- **Material Design 3** — https://m3.material.io/ · **Apple HIG** — https://developer.apple.com/design/human-interface-guidelines/ · **Fluent 2** — https://fluent2.microsoft.design/. *Authoritative for their platforms; check currency — they change yearly.*
- GNOME HIG — https://developer.gnome.org/hig/ · KDE HIG — https://develop.kde.org/hig/
- **GOV.UK Design System** — https://design-system.service.gov.uk/ · **USWDS** — https://designsystem.digital.gov/. *Research-backed, accessibility-first; strongest evidence discipline of any public system.*
- Carbon (IBM) — https://carbondesignsystem.com/ · Polaris (Shopify) — https://polaris.shopify.com/ · Atlassian — https://atlassian.design/
- **Command Line Interface Guidelines** — https://clig.dev/. *The CLI/TUI baseline.*
- Android design/vitals — https://developer.android.com/design/ui · https://developer.android.com/topic/performance/vitals

### Research organizations and evidence bases

- **Nielsen Norman Group** — https://www.nngroup.com/articles/. *Largest body of applied usability research; heuristics, eyetracking, pattern studies. Some older findings need currency checks.*
- **Baymard Institute** — https://baymard.com/. *E-commerce/forms/checkout; large-N usability evidence.*
- **WebAIM** — https://webaim.org/ (the Million report: https://webaim.org/projects/million/). *Accessibility prevalence data and SR surveys.*
- **MeasuringU** — https://measuringu.com/. *Quantitative methods, SUS benchmarks.*
- **web.dev** — https://web.dev/. *Core Web Vitals definitions and thresholds (normative for CWV).*
- MDN — https://developer.mozilla.org/ · Laws of UX — https://lawsofux.com/
- **Deceptive Design** — https://www.deceptive.design/types. *Brignull's dark-pattern taxonomy (revised 2023).* · FTC dark-patterns report — https://www.ftc.gov/reports/bringing-dark-patterns-light · EDPB guidelines — https://www.edpb.europa.eu/our-work-tools/our-documents/guidelines/guidelines-032022-deceptive-design-patterns-social-media_en
- **Microsoft HAX Toolkit** — https://www.microsoft.com/en-us/haxtoolkit/ · **Google PAIR Guidebook** — https://pair.withgoogle.com/guidebook/. *Human-AI interaction frameworks.*

### Foundational books and papers

- Norman, *The Design of Everyday Things* (2013 rev.) — affordances, signifiers, mapping, conceptual models
- Nielsen & Molich 1990 (original heuristics); Nielsen 1994 (refined set, factor analysis)
- Nielsen & Landauer 1993 — the 5-user mathematics
- Shneiderman, *Designing the User Interface* — Eight Golden Rules; direct manipulation (1983)
- Kahneman, *Thinking, Fast and Slow*; Tversky & Kahneman 1974; Kahneman & Tversky 1979 (prospect theory)
- Fitts 1954; Hick 1952 — pointing and choice-reaction laws
- Miller 1956; Cowan 2001 — working-memory capacity
- Sweller 1988 — cognitive load theory
- Reason, *Human Error* (1990) — slips/mistakes taxonomy
- Csikszentmihalyi, *Flow* (1990); Mark et al. 2005/2008 — interruption cost (resumption lag / faster-but-stressed)
- Iyengar & Lepper 2000; Scheibehenne et al. 2010 — choice overload and its meta-analysis
- Johnson & Goldstein 2003 — default effects (organ donation, opt-in ~4–28% vs opt-out 85–99%)
- Nunes & Drèze 2006 — endowed progress; Kivetz et al. 2006 — goal-gradient acceleration
- Saffer, *Microinteractions* (2013) · Krug, *Don't Make Me Think*
- Kohavi et al., *Trustworthy Online Controlled Experiments* — A/B testing practice
- Knapp, *Sprint* — design sprints · Christensen, "Know Your Customers' Jobs to Be Done" (HBR 2016)
- Buell & Norton — the labor illusion
- Ottosson — OKLab color space: https://bottosson.github.io/posts/oklab/
- Ink & Switch, "Local-first software" — https://www.inkandswitch.com/local-first/
