# AI-Agent Decision Checklists & Acceptance Criteria

This file is the operational layer of the reference: deterministic checklists, pattern-selection matrices, failure-mode diagnostics, testable acceptance criteria, and principle tradeoff tables for AI agents (and humans) proposing, implementing, or reviewing UI. Every row links back to the deep-dive file that explains *why*; use this file to decide and verify fast, and the topic files to understand and to handle cases the matrices don't cover. Start at [00-index.md](00-index.md) for the map of the whole reference.

**Caveat before anything else:** these checklists are decision support, not a substitute for testing with users. A matrix can tell you the conventional pattern; only observation tells you whether *your* users complete *their* task ([17-ux-process-research.md](17-ux-process-research.md#usability-testing)). And AI-generated review verdicts require human expert review before they gate or ship anything — an agent's confident pass is a claim to verify, not a certification ([17-ux-process-research.md](17-ux-process-research.md#design-review-rubrics-and-ai-assisted-review) and its caveats on single-evaluator coverage).

## Contents

- [1. Decision checklists](#1-decision-checklists)
- [2. Pattern selection matrices](#2-pattern-selection-matrices)
- [3. Failure-mode diagnostics](#3-failure-mode-diagnostics)
- [4. Implementation acceptance criteria](#4-implementation-acceptance-criteria)
- [5. Principle tradeoff tables](#5-principle-tradeoff-tables)
- [Quick reference](#quick-reference)

---

## 1. Decision checklists

Use these before proposing, implementing, or reviewing UI. They are intentionally direct and operational.

### Universal UI checklist

- **State:** The user can tell current location, progress, selection, saved/unsaved status, and system status ([01-core-principles.md](01-core-principles.md)).
- **Actions:** The user can identify primary, secondary, dangerous, disabled, and unavailable actions ([03-visual-design.md](03-visual-design.md)).
- **Consequences:** The user understands what will happen before committing to meaningful actions ([18-patterns-antipatterns.md](18-patterns-antipatterns.md)).
- **Recovery:** The user can undo, cancel, go back, retry, edit, or contact support when appropriate ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- **Consistency:** Labels, icons, layouts, and behaviors match product and platform conventions ([10-design-systems.md](10-design-systems.md)).
- **Accessibility:** Keyboard, focus order, screen reader names, contrast, target size, reduced motion, and error identification are handled ([05-accessibility.md](05-accessibility.md)).
- **Responsiveness:** Layout and interaction work across viewport sizes and input modes ([13-platform-web.md](13-platform-web.md)).
- **Performance:** The UI gives immediate feedback and avoids blocking interaction unnecessarily ([04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds)).
- **Content:** Text is plain-language, task-focused, and localized-ready ([09-content-ux-writing.md](09-content-ux-writing.md)).
- **Ethics:** No dark patterns, hidden costs, coercive defaults, or deceptive urgency ([18-patterns-antipatterns.md](18-patterns-antipatterns.md)).

### Pattern selection questions

Ask these before opening the matrices in [section 2](#2-pattern-selection-matrices):

- Is the task frequent or rare?
- Is the task reversible or irreversible?
- Is the user novice, expert, or mixed?
- Is the content simple, hierarchical, tabular, spatial, time-based, or conversational?
- Does the user need comparison, exploration, editing, monitoring, or completion?
- Is the context mobile, desktop, web, embedded, kiosk, CLI, or assistive technology-heavy?
- What is the cost of delay, error, disclosure, or abandonment?

### Component choice rules

Detail on every rule: [11-components-and-overlays.md](11-components-and-overlays.md) and [07-forms-and-input.md](07-forms-and-input.md).

- Use a **button** for actions that change state or submit work.
- Use a **link** for navigation to another page, route, or document.
- Use a **checkbox** for independent yes/no choices or multi-select.
- Use **radio buttons** when one choice among a small visible set is required (~5–7 options is the practical ceiling before a list gets unwieldy; see [07-forms-and-input.md](07-forms-and-input.md#input-type-selection)).
- Use a **select/dropdown** only as a last resort — when options are many, familiar, and not worth comparing side by side, and radio/checkbox/autocomplete/typed input are all worse fits ([07-forms-and-input.md](07-forms-and-input.md#input-type-selection)).
- Use **tabs** for peer sections within the same context, not for step-by-step flows ([08-navigation-ia.md](08-navigation-ia.md)).
- Use an **accordion** to reduce scanning cost, not to hide mandatory information.
- Use a **modal** only when the user must resolve a focused task before continuing ([11-components-and-overlays.md](11-components-and-overlays.md)).
- Use **toast notifications** only for non-blocking status that does not require immediate action and does not need to be retained ([04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)).
- Use **inline validation** for field-specific problems; use **summary validation** for form-level failure ([07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing)).

### High-risk flow checklist

For payments, deletion, legal/identity submissions, and anything irreversible ([19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md) for checkout specifics):

- The user can review all critical details before commit.
- The action label names the outcome, not just `OK` or `Submit` ([09-content-ux-writing.md](09-content-ux-writing.md)).
- Destructive actions are separated from safe actions.
- Confirmation is used only when needed and includes the object, consequence, and escape route ([07-forms-and-input.md](07-forms-and-input.md#destructive-action-confirmation)).
- Undo is preferred over confirmation when safe and reliable ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- Errors cannot silently fail.
- Audit trails, receipts, or confirmation records are available where relevant.
- Accessibility is manually checked, not only linted ([05-accessibility.md](05-accessibility.md#testing-methodology)).

### Accessibility acceptance checklist

Deep dives: [05-accessibility.md](05-accessibility.md); widget-level keyboard/ARIA contracts: [06-aria-widget-reference.md](06-aria-widget-reference.md).

- Every interactive element has an accessible name, role, value, and visible/focus state (WCAG 4.1.2, 2.4.7).
- Focus order follows visual and task order (WCAG 2.4.3).
- No keyboard trap exists (WCAG 2.1.2).
- Focus is moved intentionally after route changes, modal opens/closes, dynamic validation, and async updates.
- Text contrast meets WCAG 2.2 AA minimums — 4.5:1 for normal text, 3:1 for large text (SC 1.4.3); non-text UI indicators meet 3:1 (SC 1.4.11) ([05-accessibility.md](05-accessibility.md#contrast-minimum)).
- Information is not conveyed by color alone (SC 1.4.1) ([05-accessibility.md](05-accessibility.md#use-of-color)).
- Touch/pointer targets are at least 24×24 CSS px (SC 2.5.8) or meet a valid spacing exception; platform guidance is larger (~44–48px) ([05-accessibility.md](05-accessibility.md#target-size)).
- Motion respects reduced-motion preferences and is not required to understand content ([05-accessibility.md](05-accessibility.md#reduced-motion)).
- Error messages are associated with fields and announced when appropriate (SC 3.3.1) ([05-accessibility.md](05-accessibility.md#forms-accessibility)).
- Custom widgets follow WAI-ARIA Authoring Practices or use native controls instead ([06-aria-widget-reference.md](06-aria-widget-reference.md)).

### Web checklist

Full treatment: [13-platform-web.md](13-platform-web.md).

- Use semantic HTML first; ARIA only fills semantic gaps.
- Preserve browser behaviors: links, buttons, forms, history, focus, zoom, text selection, context menus.
- Provide responsive layouts with content-driven breakpoints where possible.
- Avoid hover-only interactions ([04-interaction-design.md](04-interaction-design.md#hover-dependence)).
- Optimize Core Web Vitals — LCP ≤ 2.5s, INP ≤ 200ms, CLS ≤ 0.1 at the 75th percentile — and perceived responsiveness ([13-platform-web.md](13-platform-web.md)).
- Make form inputs compatible with autocomplete, labels, validation, and password managers ([07-forms-and-input.md](07-forms-and-input.md)).

### Desktop checklist

Full treatment: [14-platform-desktop.md](14-platform-desktop.md).

- Respect OS window, menu, shortcut, file, drag/drop, and notification conventions.
- Provide keyboard shortcuts for frequent expert tasks and visible menu paths for discoverability.
- Handle multi-window, unsaved changes, large screens, dense data, and precision pointer use ([12-data-tables-dashboards.md](12-data-tables-dashboards.md) for dense data).
- Do not import mobile-only patterns when native desktop controls would be clearer.

### Mobile checklist

Full treatment: [15-platform-mobile.md](15-platform-mobile.md).

- Primary actions are reachable and have adequate target size.
- Critical controls do not depend on hover, precision pointer, or multi-key shortcuts.
- Gestures have visible alternatives ([05-accessibility.md](05-accessibility.md#dragging-alternatives)).
- Text input minimizes typing through appropriate keyboards, suggestions, scanning, pickers, and defaults.
- The UI handles interruption, rotation, poor connectivity, and one-handed use.

### CLI/TUI checklist

The command line is a platform too; full treatment: [16-platform-cli-tui.md](16-platform-cli-tui.md).

- Human-first output on stdout; diagnostics and progress on stderr.
- Provide `--json` (or equivalent) machine-readable output for scripting ([16-platform-cli-tui.md](16-platform-cli-tui.md#machine-readable-output)).
- Exit codes are meaningful: 0 for success, nonzero for distinct failure classes ([16-platform-cli-tui.md](16-platform-cli-tui.md#exit-codes)).
- Detect TTY: disable color/interactivity/prompts when output is piped or non-interactive.
- Destructive commands support `--dry-run`, confirmation, or explicit force flags.
- Errors state what failed and what to try next; suggest corrections for near-miss commands.

---

## 2. Pattern selection matrices

Purpose: choose UI patterns quickly and defensibly. These matrices favor deterministic decision rules over explanatory prose. Use them when choosing a component before implementation, reviewing a proposed UI for pattern mismatch, translating a user story into UI structure, or resolving conflicts between candidate components. When a row's stance surprises you, read the linked deep-dive file before overriding it.

### Global decision rules

| Situation | Choose | Avoid | Required checks | Deep dive |
| --- | --- | --- | --- | --- |
| Navigates to a new page, route, document, or anchor | Link | Button with navigation handler | Link has meaningful text, visited/focus states, correct URL/history behavior | [13-platform-web.md](13-platform-web.md) |
| Performs an action, mutates state, submits data, opens a dialog | Button | Link styled as action | Button label names outcome, disabled/loading states handled | [11-components-and-overlays.md](11-components-and-overlays.md) |
| One choice from a small visible set (2–7 options; ~5–7 is the practical ceiling) | Radio group | Select/dropdown | Group label, one selected if required, keyboard arrow behavior | [07-forms-and-input.md](07-forms-and-input.md#input-type-selection) |
| Multiple independent choices | Checkbox group | Multi-select dropdown unless space constrained | Fieldset/group label, independent checked states | [07-forms-and-input.md](07-forms-and-input.md) |
| Binary setting that takes effect immediately | Switch/toggle | Checkbox if setting is not immediate or label is ambiguous | Current state clear, label describes setting, state persists | [11-components-and-overlays.md](11-components-and-overlays.md) |
| Binary form value submitted later | Checkbox | Switch | Label says what checked means, form validation handles required agreement | [07-forms-and-input.md](07-forms-and-input.md) |
| Many familiar options with low comparison need | Select/combobox — as a last resort, after radio, checkbox, autocomplete, and typed input are ruled out | Radio list that becomes too long; select as the default reach | Label, keyboard search, clear selected state | [07-forms-and-input.md](07-forms-and-input.md#input-type-selection) |
| Many options where search matters | Combobox/autocomplete | Plain select | Accessible name, popup behavior, keyboard support, no forced exact match unless required | [06-aria-widget-reference.md](06-aria-widget-reference.md) |
| Approximate numeric value | Slider/range | Text field only | Text/numeric fallback if precision matters, keyboard support | [07-forms-and-input.md](07-forms-and-input.md) |
| Precise numeric value | Number input/text input with validation | Slider only | Units, range, step, locale, invalid state | [07-forms-and-input.md](07-forms-and-input.md) |
| Date the user already knows (birth date, document date) | Typed day/month/year input as primary (separate fields or parsed text), optional picker | Calendar-only picker | Locale, format hint, keyboard entry, validation | [07-forms-and-input.md](07-forms-and-input.md#date-input) |
| Date chosen by browsing availability (scheduling) | Calendar/date picker | Raw text only | Keyboard-accessible calendar, manual entry fallback | [07-forms-and-input.md](07-forms-and-input.md#date-input) |
| Short free text | Text field | Textarea | Label, autocomplete/inputmode where relevant | [07-forms-and-input.md](07-forms-and-input.md) |
| Long free text | Textarea/rich editor | Single-line input | Character count if constrained, autosave for long effort | [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery) |
| High-risk final submission | Review page + explicit commit button | Single blind submit | Summary, edit links, consequence text, error recovery | [07-forms-and-input.md](07-forms-and-input.md#high-risk-and-regulated-forms), [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md) |
| Reversible destructive action | Undo snackbar/toast or recycle state | Repeated confirmation | Undo duration/accessibility, persistence model | [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs) |
| Irreversible destructive action | Confirmation dialog/page | Toast-only undo | Object named, consequence stated, safe default focus/action | [07-forms-and-input.md](07-forms-and-input.md#destructive-action-confirmation) |
| Action prerequisites unmet | Enabled control + validation explaining what's missing on activation (GOV.UK stance) | Disabled control with no visible reason | If disabled is unavoidable, explanation is adjacent and discoverable | [07-forms-and-input.md](07-forms-and-input.md), [18-patterns-antipatterns.md](18-patterns-antipatterns.md) |
| Local field problem | Inline field error | Global-only error | Error associated with field, message actionable | [07-forms-and-input.md](07-forms-and-input.md#error-message-design) |
| Form submission has multiple errors | Error summary + inline errors | Alert with no field links | Focus summary, links to fields, preserve input | [07-forms-and-input.md](07-forms-and-input.md#error-message-design) |
| Brief non-blocking status | Toast/flag/status region | Modal dialog | Does not hide critical info, accessible announcement if meaningful | [04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design) |
| Important persistent page/section status | Inline alert/banner/infobar | Toast that disappears | Severity, action, dismiss rules, accessible semantics | [11-components-and-overlays.md](11-components-and-overlays.md) |
| Blocking decision required before continuing | Modal dialog | Tooltip/popover | Focus trap, labelled dialog, close/cancel, restore focus | [11-components-and-overlays.md](11-components-and-overlays.md), [06-aria-widget-reference.md](06-aria-widget-reference.md) |
| Supplemental contextual info | Popover/inline disclosure | Modal | Trigger relation, dismissal, keyboard, no critical hover-only content | [11-components-and-overlays.md](11-components-and-overlays.md) |
| Short non-interactive label/help | Tooltip | Popover with controls | Trigger on focus/hover, not sole source of required info | [11-components-and-overlays.md](11-components-and-overlays.md) |
| Long optional detail | Disclosure/accordion | Tooltip | State announced, heading/button semantics | [11-components-and-overlays.md](11-components-and-overlays.md) |
| Peer content sections in one context | Tabs | Accordion for unrelated steps | Tablist semantics, selected tab, keyboard behavior | [06-aria-widget-reference.md](06-aria-widget-reference.md) |
| Sequential process | Stepper/progress indicator + pages | Tabs | Current step, completed steps, back/edit behavior | [08-navigation-ia.md](08-navigation-ia.md) |
| Hierarchical navigation location | Breadcrumb | Back button only | Current item, parent links, not primary nav alone | [08-navigation-ia.md](08-navigation-ia.md) |
| Large known-item content | Search | Deep-only navigation | Scope, filters, no-results recovery | [08-navigation-ia.md](08-navigation-ia.md) |
| Narrow content set with categories | Browse/navigation | Search-only | Clear categories, selected state | [08-navigation-ia.md](08-navigation-ia.md) |
| Compare records by attributes | Table/data grid | Cards | Headers, sorting, responsive behavior, keyboard if interactive | [12-data-tables-dashboards.md](12-data-tables-dashboards.md) |
| Browse visually rich heterogeneous items | Cards/grid | Dense table | Accessible names, actions, selection state | [12-data-tables-dashboards.md](12-data-tables-dashboards.md) |
| Edit many cells rapidly | Data grid | Form-per-row | Keyboard model, focus, validation, undo/bulk operations | [12-data-tables-dashboards.md](12-data-tables-dashboards.md) |
| Show KPIs or trend monitoring | Dashboard | Static report only | Freshness, definitions, empty/error states, drilldown | [12-data-tables-dashboards.md](12-data-tables-dashboards.md) |
| User needs configuration | Settings page | Hidden modal chain | Grouping, search, defaults, reset/undo, scope | [08-navigation-ia.md](08-navigation-ia.md) |
| Need expert command speed | Command palette | Only hidden shortcuts | Discoverability, permissions, undo, keyboard | [14-platform-desktop.md](14-platform-desktop.md) |
| Automation/scripting consumer | CLI with human-first output, `--json` for machines, meaningful exit codes, TTY detection | GUI-only path for automatable workflows | Piped output stays parseable, non-interactive mode works, errors on stderr | [16-platform-cli-tui.md](16-platform-cli-tui.md) |
| AI generates uncertain output | Draft/recommendation with review | Silent autonomous commit | Confidence/provenance, edit/undo, user confirmation for risk | [20-ai-product-ux.md](20-ai-product-ux.md) |

### Navigation pattern matrix

Full treatment: [08-navigation-ia.md](08-navigation-ia.md).

| Information shape | Primary pattern | Use when | Avoid when | Required states |
| --- | --- | --- | --- | --- |
| 3–7 top-level destinations, mobile | Bottom navigation | Destinations are peers and frequent | More than 7 items, deep hierarchy | Current item, labels, badges with accessible names ([15-platform-mobile.md](15-platform-mobile.md)) |
| Many app sections, desktop/web | Side navigation | Persistent app workspace | Marketing/simple sites | Collapsed/expanded, selected, nested, keyboard |
| Few top-level website categories | Top navigation | Broad website navigation | Complex app workflows | Current item, responsive menu, skip links |
| Hierarchical location | Breadcrumb | Users navigate nested structures | Flat apps | Parent links, current page, truncation handling |
| Peer views of same object | Tabs | Content sections are same-level | Sequential process | Selected, focus, panel labelled by tab ([06-aria-widget-reference.md](06-aria-widget-reference.md)) |
| Process flow | Step indicator | User must complete steps | Random-access peers | Current, completed, blocked, error step |
| Expert action access | Command palette | Many commands, keyboard-heavy users | Novice-only critical path | Search, shortcuts, permissions, no destructive surprise ([14-platform-desktop.md](14-platform-desktop.md)) |

### Overlay pattern matrix

Full treatment: [11-components-and-overlays.md](11-components-and-overlays.md); ARIA contracts: [06-aria-widget-reference.md](06-aria-widget-reference.md).

| Need | Pattern | Use when | Avoid when | Required checks |
| --- | --- | --- | --- | --- |
| Interrupt for required decision | Modal dialog | User must decide before continuing | Optional info or background status | Focus management, escape/cancel, labelled, restore focus |
| Interrupt for urgent destructive/legal/security issue | Alert dialog | Consequence is important and immediate | Routine validation | Initial focus safe, message concise, action explicit |
| Side task while preserving page context | Drawer/side panel | User references underlying context | Task is simple or page navigation is clearer | Focus, close, width/responsiveness |
| Small contextual interactive content | Popover/inline dialog | Content relates to trigger and may include controls | Critical workflow, large content | Trigger relation, dismissal, keyboard |
| Brief text explanation | Tooltip | Nonessential, noninteractive, short | Required info, mobile-only, interactive content | Focus/hover, delay, escape, no keyboard trap ([05-accessibility.md](05-accessibility.md#content-on-hover-or-focus)) |
| Inline optional detail | Disclosure/accordion | User chooses to reveal details | Required warnings/pricing/consent | Button semantics, expanded state, persistent content |

### Feedback pattern matrix

Full treatment: [04-interaction-design.md](04-interaction-design.md).

| Event | Pattern | Use when | Avoid when | Accessibility requirement |
| --- | --- | --- | --- | --- |
| Immediate control activation | Pressed/selected/loading state | Any direct interaction | None | State exposed programmatically when meaningful |
| Short async work | Inline spinner/skeleton | Region loading | Whole page not blocked | Loading text or status if not visually obvious ([04-interaction-design.md](04-interaction-design.md#loading-strategies)) |
| Long determinate work | Progress bar | Real progress measurable | Fake progress | Name, value, ability to cancel if possible |
| Long indeterminate work | Progress indicator + explanatory text | Duration unknown | No user value in waiting | Describe task, timeout/retry path |
| Instant-feel local update | Optimistic UI (apply immediately, reconcile in background) | Fast, low-risk, high-success operations (likes, toggles, reordering) | **Never for payments or irreversible actions** — show real state and real confirmation, no "reliable rollback" exception | Failure reconciliation visible and announced; user's view never silently diverges from server ([04-interaction-design.md](04-interaction-design.md#optimistic-ui)) |
| Autosave | Inline saved/saving/error status | Editing drafts | Toast-only save state | Status visible near editor and announced if important ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)) |
| Successful low-risk action | Toast/flag | Non-blocking confirmation | User must keep record | Accessible status, duration enough or persistent log ([04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)) |
| High-risk success | Confirmation page/receipt | Payment, legal, booking, identity | Disappearing toast | Reference ID, summary, next steps ([19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md)) |
| Recoverable error | Inline alert/error message | User can fix now | Vague generic message | Problem + action + field association ([07-forms-and-input.md](07-forms-and-input.md#error-message-design)) |
| System outage/permission issue | Persistent banner/empty state | User cannot complete task | Field-level error | Reason, recovery, support/contact |

### Form control matrix

Full treatment: [07-forms-and-input.md](07-forms-and-input.md).

| Data need | Control | Use when | Avoid when | Required agent checks |
| --- | --- | --- | --- | --- |
| Legal agreement | Checkbox | User must explicitly agree | Pre-checked consent | Label includes agreement, link text meaningful |
| Contact email | Email input | Email syntax expected | Username can be non-email | autocomplete, inputmode, validation message |
| Password | Password input | Secret text | One-time code pasted from app? consider text/code field | Password manager, paste allowed, show/hide control accessible ([07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)) |
| One-time code | Text input with autocomplete | OTP required | Manual transcription only | Paste allowed, `autocomplete=one-time-code` where web |
| Address | Structured fields or lookup + manual edit | Delivery/billing needed | One format globally | Locale support, autocomplete, manual fallback |
| Name | Flexible text fields | Personal name | First/last required globally | International naming support ([09-content-ux-writing.md](09-content-ux-writing.md#internationalization-and-localization)) |
| Phone | Tel input | Phone contact required | Numeric-only validation | Country/region, spaces/symbols accepted |
| File upload | File input/drop zone | User supplies file | Drag-only upload | Browse button, type/size limits, progress, remove ([07-forms-and-input.md](07-forms-and-input.md#file-upload)) |
| Date of birth | Typed separate fields or flexible text — memorable dates are typed, not browsed | Calendar-only picker | Calendar-only picker | Format hints, accessible errors ([07-forms-and-input.md](07-forms-and-input.md#date-input)) |
| Appointment date | Calendar/date picker + text entry | Choosing availability | Text only with many constraints | Disabled dates explained, keyboard support ([07-forms-and-input.md](07-forms-and-input.md#date-input)) |

### Risk-based pattern escalation

Cross-reference: [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs), [07-forms-and-input.md](07-forms-and-input.md#high-risk-and-regulated-forms).

| Risk | Examples | Minimum safeguards | Escalate to |
| --- | --- | --- | --- |
| Low | Filtering, theme toggle, sort | Immediate feedback, undo if useful; optimistic UI acceptable | No confirmation |
| Medium | Save settings, publish internal note | Clear label, validation, undo/edit | Confirmation if externally visible |
| High | Payment, delete, public post, legal submission | Review, explicit action, recovery, receipt; never optimistic UI | Confirmation page/dialog, audit trail |
| Critical | Medical, financial transfer, identity, irreversible data loss | Strong verification, human support, logging, redundant checks | Multi-step review, out-of-band confirmation where appropriate |

---

## 3. Failure-mode diagnostics

Purpose: diagnose likely UX failures from observed symptoms and recommend targeted fixes with acceptance criteria.

### Diagnostic workflow

1. Identify the user-visible symptom.
2. Map it to a likely principle or pattern failure ([01-core-principles.md](01-core-principles.md), [18-patterns-antipatterns.md](18-patterns-antipatterns.md)).
3. Check platform/accessibility obligations ([05-accessibility.md](05-accessibility.md), files 13–16).
4. Recommend the smallest fix that removes the failure.
5. Add acceptance criteria ([section 4](#4-implementation-acceptance-criteria)).

### Symptom-to-fix matrix

| Symptom | Likely failure | Diagnose | Recommended fix | Acceptance criteria | Deep dive |
| --- | --- | --- | --- | --- | --- |
| Users click repeatedly (rage clicks) | No feedback, slow response, disabled state unclear | Test slow network and async path | Add pressed/loading/success/error states; prevent duplicate submit | One activation gives immediate feedback and duplicate action is blocked or harmless | [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds) |
| Users abandon form after errors | Error messages vague, fields cleared, errors hard to find | Submit invalid form | Add error summary, inline errors, preserve input, focus summary | User can identify and fix every error without retyping valid fields | [07-forms-and-input.md](07-forms-and-input.md#error-message-design) |
| Keyboard user loses location | Focus invisible/obscured or unexpected focus movement | Tab through page at zoom | Add visible focus, prevent sticky overlap, logical focus order | Focus is visible and not hidden at all interactive elements (SC 2.4.7, 2.4.11) | [05-accessibility.md](05-accessibility.md#focus-visible) |
| Screen reader user hears unlabeled controls | Missing accessible names | Inspect accessibility tree | Add visible labels or accessible names; prefer visible labels | Every control has name/role/value (SC 4.1.2) | [05-accessibility.md](05-accessibility.md#screen-readers) |
| User cannot finish without mouse | Custom widget lacks keyboard behavior | Keyboard-only test | Use native control or implement APG keyboard model | Full task complete via keyboard (SC 2.1.1) | [06-aria-widget-reference.md](06-aria-widget-reference.md) |
| User cannot use touch accurately | Targets too small/crowded | Test mobile and tremor-prone scenarios | Increase target size/spacing, separate destructive actions | Targets meet 24×24 CSS px minimum (SC 2.5.8) or spacing exceptions; ~44–48px on touch platforms | [05-accessibility.md](05-accessibility.md#target-size), [15-platform-mobile.md](15-platform-mobile.md) |
| User cannot discover feature | Hidden gesture/hover/shortcut only | Observe novice task | Add visible path, menu item, button, or onboarding hint | Critical feature has a visible non-gesture path | [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures) |
| User edits wrong object | Scope/location unclear | Check header, breadcrumbs, object labels | Add object title, environment/account/workspace indicator | User can state affected object before acting | [08-navigation-ia.md](08-navigation-ia.md) |
| User loses work on navigation | No draft/unsaved warning | Navigate away mid-edit | Autosave, draft, or unsaved-changes guard | User can recover or intentionally discard changes | [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery) |
| User cannot interpret disabled button | No prerequisite explanation | Inspect disabled controls | Prefer enabling the action with explanatory validation on activation (default stance); add inline prerequisite text if it must stay disabled | User knows what to do to enable action | [07-forms-and-input.md](07-forms-and-input.md), [18-patterns-antipatterns.md](18-patterns-antipatterns.md) |
| Users ignore confirmations | Too many generic dialogs | Count confirmation frequency | Replace reversible confirmations with undo; keep confirmations for high risk | Confirmations appear only for meaningful risk and name consequence | [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs) |
| Users misunderstand primary action | Vague labels like Submit/Continue | Review action labels | Rename with outcome-specific verbs | Label predicts result without surrounding context | [09-content-ux-writing.md](09-content-ux-writing.md) |
| Users cannot recover from error | Error lacks next step | Trigger failure states | Add actionable recovery: retry, edit, contact, save copy | Every error offers feasible next step | [04-interaction-design.md](04-interaction-design.md#error-feedback-and-error-taxonomy) |
| Users miss important status | Toast disappears or status not near task | Trigger success/failure | Use persistent inline status or confirmation page | Required information remains available | [04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design) |
| Users get lost in flow | No progress/location/back | Walk multi-step path | Add step indicator, page titles, back/edit, save progress | User can identify current step and return safely | [08-navigation-ia.md](08-navigation-ia.md) |
| Users re-enter same data | Process asks redundant info | Trace process | Auto-populate or provide selection from prior entry | Previously entered info reused unless essential/security/invalid (SC 3.3.7) | [05-accessibility.md](05-accessibility.md#redundant-entry) |
| Authentication excludes users | Memory/transcription/puzzle required | Test password manager, paste, alternate login | Allow password managers, paste, magic link/passkey/SSO alternative | Login does not require an unsupported cognitive function test (SC 3.3.8) | [05-accessibility.md](05-accessibility.md#accessible-authentication), [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux) |
| Drag/drop blocks users | No non-drag alternative | Attempt with keyboard/single pointer | Add up/down buttons, move menu, select destination | Same function possible without dragging (SC 2.5.7) | [05-accessibility.md](05-accessibility.md#dragging-alternatives) |
| AI output trusted blindly | No uncertainty/provenance/review | Review AI feature flow | Add confidence, sources, editable draft, confirmation for high risk | User can inspect, edit, reject, or verify AI output | [20-ai-product-ux.md](20-ai-product-ux.md) |
| AI failure is dead end | No fallback when AI wrong | Force bad/empty AI result | Add retry, refine, manual path, escalation | User can complete task without AI success | [20-ai-product-ux.md](20-ai-product-ux.md) |
| Layout shifts while reading/clicking | CLS/layout instability | Throttle load and interact | Reserve space, avoid insertions above focus, stabilize ads/media | CLS ≤ 0.1; content does not jump during critical interaction | [13-platform-web.md](13-platform-web.md) |
| Page feels slow after load | Poor INP/main-thread blocking | Field performance + interaction testing | Reduce JS, split work, defer noncritical, optimize handlers | Common interactions respond within INP ≤ 200ms at p75 | [13-platform-web.md](13-platform-web.md) |

### Accessibility-specific diagnostics

Detail: [05-accessibility.md](05-accessibility.md#manual-test-scripts) and [06-aria-widget-reference.md](06-aria-widget-reference.md).

| Test | Failure signal | Likely cause | Fix |
| --- | --- | --- | --- |
| Tab through all controls | Focus disappears | CSS removes outline, overlay covers focus, tabindex misuse | Restore visible focus, remove positive tabindex, manage overlays |
| Screen reader reads button as `button` only | Icon-only control lacks name | Missing label | Add visible text or `aria-label` with outcome |
| Screen reader repeats confusing text | Redundant labels/descriptions | Misused aria-labelledby/describedby | Keep accessible name concise, description supplemental |
| Form error not announced | Error not associated | Missing `aria-describedby`/invalid state or focus summary | Associate errors and focus error summary after submit |
| Modal background still navigable | Focus not trapped/inert | Dialog implementation incomplete | Use native dialog or robust modal pattern ([06-aria-widget-reference.md](06-aria-widget-reference.md)) |
| Tooltip unavailable on keyboard/touch | Hover-only trigger | Mouse-only behavior | Trigger on focus and provide inline alternative for required info |
| Custom select unusable | ARIA role without full behavior | Incomplete combobox/listbox pattern | Prefer native select or implement APG pattern fully ([06-aria-widget-reference.md](06-aria-widget-reference.md)) |

### Platform mismatch diagnostics

| Platform | Symptom | Likely mismatch | Fix |
| --- | --- | --- | --- |
| Web | Browser back loses app state unexpectedly | SPA route/state not integrated | Use URL/history state and restore focus/scroll intentionally ([13-platform-web.md](13-platform-web.md)) |
| Web | Password manager fails | Custom inputs or blocked paste | Use semantic inputs, autocomplete, allow paste ([07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)) |
| Desktop | Power users complain workflow is slow | No shortcuts/menu paths/bulk actions | Add standard shortcuts, menus, multi-select, command palette ([14-platform-desktop.md](14-platform-desktop.md)) |
| Desktop | App feels like mobile port | Oversized sparse UI, hidden menus | Increase density options, native menus, resizable panes ([14-platform-desktop.md](14-platform-desktop.md)) |
| Mobile | Users miss actions in overflow | Desktop nav compressed poorly | Use mobile-first primary action placement and bottom/nav patterns ([15-platform-mobile.md](15-platform-mobile.md)) |
| Mobile | Keyboard covers field/action | View not adjusted for virtual keyboard | Scroll focused input into view, keep submit reachable ([15-platform-mobile.md](15-platform-mobile.md)) |
| CLI | Scripts break on output changes or piped color codes | Human output treated as the machine interface | Provide stable `--json` output, detect TTY, keep diagnostics on stderr ([16-platform-cli-tui.md](16-platform-cli-tui.md#machine-readable-output)) |
| CLI | Automation hangs waiting for input | Interactive prompt in non-interactive context | Detect non-TTY, fail with actionable error or accept flags/env for all inputs ([16-platform-cli-tui.md](16-platform-cli-tui.md)) |

### Severity scoring for agents

| Severity | Definition | Examples | Response |
| --- | --- | --- | --- |
| Blocker | Prevents task completion or excludes assistive technology users | Keyboard trap, data loss, inaccessible login | Must fix before shipping |
| High | Causes high error risk, privacy/security issue, or major abandonment | Hidden destructive consequence, invalid payment recovery | Fix before release unless explicitly accepted |
| Medium | Slows task or causes confusion with workaround | Weak labels, poor empty state, inconsistent nav | Fix in current iteration if feasible |
| Low | Cosmetic or minor efficiency issue | Slightly verbose copy, minor alignment | Backlog if no larger impact |

Severity here aligns with the prioritization guidance in [17-ux-process-research.md](17-ux-process-research.md); blockers and accessibility failures skip the queue.

### Agent recommendation format

```md
Finding: ...
Symptom: ...
Likely failure mode: ...
Affected users: ...
Severity: Blocker | High | Medium | Low
Recommended fix: ...
Acceptance criteria: ...
Related guidance: ...
```

Remember: a diagnostic verdict produced by an AI agent is an input to human judgment, not a shipped conclusion ([17-ux-process-research.md](17-ux-process-research.md#design-review-rubrics-and-ai-assisted-review)).

---

## 4. Implementation acceptance criteria

Testable criteria for implementing or reviewing product UI. Each criterion should be verifiable by inspection or a concrete test.

### Universal criteria

- Every interactive element has a visible affordance and a programmatic role ([04-interaction-design.md](04-interaction-design.md#affordances-and-signifiers-in-practice)).
- Every meaningful action produces immediate feedback.
- Every async action has loading, success, and failure states appropriate to its duration and consequence.
- Every destructive or high-risk action names the object and consequence before commit.
- Every recoverable action offers undo, edit, retry, cancel, or back where feasible.
- Every screen has a clear title or primary purpose.
- Primary action, secondary action, and destructive action are visually and semantically distinguishable ([03-visual-design.md](03-visual-design.md)).
- Empty, loading, error, partial, permission-denied, offline, and success states are designed where applicable ([04-interaction-design.md](04-interaction-design.md#empty-states)).
- Text is plain-language and does not require internal product knowledge ([09-content-ux-writing.md](09-content-ux-writing.md)).
- The UI remains usable with localization, text expansion, and user scaling. Layouts must tolerate roughly 30–50% expansion for paragraph-length text **and 100–300% for short labels and buttons** (short strings expand disproportionately in translation); design with at least +40% width headroom on labels and controls ([09-content-ux-writing.md](09-content-ux-writing.md#internationalization-and-localization)).

### Accessibility criteria

Detail and test scripts: [05-accessibility.md](05-accessibility.md).

- Keyboard users can complete all tasks without pointer input where keyboard input exists (SC 2.1.1).
- Focus order follows visual and task order (SC 2.4.3).
- Focus indicator is visible (SC 2.4.7) and not obscured by sticky or overlay content (SC 2.4.11).
- Modal dialogs move focus inside on open, trap focus while open, and restore focus on close ([06-aria-widget-reference.md](06-aria-widget-reference.md)).
- No keyboard trap exists unless explicit instructions and escape are provided (SC 2.1.2).
- Every input has a persistent visible label or equivalent visible naming pattern ([07-forms-and-input.md](07-forms-and-input.md#labels-vs-placeholders)).
- Every form error is programmatically associated with the relevant field (SC 3.3.1).
- Failed form submission moves focus to an error summary or first error in a predictable way.
- Color is not the only indicator of state, status, error, requiredness, or selection (SC 1.4.1).
- Text contrast meets WCAG 2.2 AA — 4.5:1 normal text, 3:1 large text (SC 1.4.3); non-text UI contrast meets 3:1 (SC 1.4.11) — unless a documented exception applies ([05-accessibility.md](05-accessibility.md#contrast-minimum)).
- Pointer/touch targets meet the 24×24 CSS px minimum (SC 2.5.8) or valid spacing exceptions ([05-accessibility.md](05-accessibility.md#target-size)).
- Dragging, swiping, and complex gestures have non-gesture alternatives (SC 2.5.7).
- Authentication does not require unsupported memory, transcription, or puzzle-solving without an accessible alternative or assistive mechanism (SC 3.3.8) ([05-accessibility.md](05-accessibility.md#accessible-authentication)).
- Dynamic status changes are announced when they affect task progress or outcome (SC 4.1.3).
- Motion respects reduced-motion preferences and is not required for comprehension ([05-accessibility.md](05-accessibility.md#reduced-motion)).

### Web criteria

Detail: [13-platform-web.md](13-platform-web.md).

- Navigation uses links; state-changing actions use buttons.
- Semantic HTML is used before ARIA ([05-accessibility.md](05-accessibility.md#semantic-html-first)).
- Page/route changes update title, focus, and meaningful history state.
- Browser back/forward does not destroy unsaved work without warning or autosave.
- Forms use native form semantics where practical.
- Inputs use appropriate `type`, `autocomplete`, `inputmode`, and validation constraints ([07-forms-and-input.md](07-forms-and-input.md)).
- Password and one-time-code fields allow paste and password manager/autofill support.
- The page remains usable at 200% zoom (SC 1.4.4) and at 320 CSS px width equivalent / common mobile widths (SC 1.4.10 reflow; [05-accessibility.md](05-accessibility.md#reflow)).
- Core Web Vitals meet targets at the 75th percentile: LCP ≤ 2.5s, INP ≤ 200ms, CLS ≤ 0.1.
- Layout reserves space for images, ads, embeds, and async content to prevent unexpected shifts.
- Hover-only content is not required for task completion (SC 1.4.13).

### Desktop criteria

Detail: [14-platform-desktop.md](14-platform-desktop.md).

- Standard commands have platform-appropriate shortcuts where expected.
- Shortcut commands also have visible menu or toolbar paths.
- Menus, context menus, title bars, file dialogs, notifications, and window behavior follow OS conventions.
- Unsaved changes are protected by autosave, drafts, or explicit save/discard prompts.
- Multi-window and multi-document scope is visible.
- Dense data views support keyboard navigation, selection, sorting, resizing, and horizontal/vertical overflow intentionally ([12-data-tables-dashboards.md](12-data-tables-dashboards.md)).
- Preferences/settings apply at a visible scope: app, account, workspace, project, or document.

### Mobile criteria

Detail: [15-platform-mobile.md](15-platform-mobile.md).

- Primary actions are reachable and not hidden behind keyboard or system UI.
- Touch targets have adequate size and spacing (platform guidance ~44pt iOS / 48dp Android; WCAG floor 24 CSS px, SC 2.5.8).
- Gestures have visible alternatives.
- Virtual keyboard type matches expected input.
- Layout respects safe areas, orientation changes, dynamic type, and narrow widths.
- The flow tolerates interruption, app backgrounding, poor connectivity, and retry.
- Permission prompts are contextual and explain value before the system prompt where appropriate.

### CLI/TUI criteria

Detail: [16-platform-cli-tui.md](16-platform-cli-tui.md).

- Human-readable output by default; `--json` or equivalent stable format for machine consumers.
- stdout carries results; stderr carries diagnostics, progress, and errors.
- Exit codes distinguish success from failure classes; scripts can branch on them.
- TTY detection disables color, spinners, and prompts in non-interactive contexts.
- Destructive operations support `--dry-run` and require explicit confirmation or force flags.
- Help (`-h`/`--help`) shows usage, common examples, and flag documentation.

### Forms criteria

Detail: [07-forms-and-input.md](07-forms-and-input.md).

- Each field has a reason to be collected at that point in the flow.
- Required/optional status is clear.
- Help text explains non-obvious constraints before error.
- Validation occurs at a helpful time and does not punish incomplete typing ([07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing)).
- Errors include the problem and how to fix it.
- Valid fields retain values after failed submission.
- Submit controls default to enabled, with validation explaining problems on activation; disabled-until-valid is the fallback and requires adjacent explanation.
- Long forms can be saved, resumed, or safely navigated where user effort is high.
- Review pages summarize high-risk submissions with edit paths ([07-forms-and-input.md](07-forms-and-input.md#high-risk-and-regulated-forms)).
- Redundant entry is avoided within a process unless security, validity, or essential task logic requires it (SC 3.3.7).

### Feedback criteria

Detail: [04-interaction-design.md](04-interaction-design.md).

- Loading states are local when only local content is affected.
- Determinate progress is used only when actual progress is knowable.
- Indeterminate loading includes useful context when delay may exceed a short wait.
- Optimistic UI is used only for low-risk, reversible operations — never for payments or irreversible actions ([04-interaction-design.md](04-interaction-design.md#optimistic-ui)).
- Success messages persist when users need to retain a reference, receipt, or next step.
- Toasts do not contain the only copy of critical information.
- Error states preserve user work and provide recovery.
- Offline and server-failure states distinguish user-fixable vs system-fixable problems.

### AI feature criteria

Detail: [20-ai-product-ux.md](20-ai-product-ux.md).

- AI capability and limitations are communicated before reliance.
- User can review, edit, reject, retry, or regenerate AI output.
- High-risk AI output is never committed silently.
- Sources, provenance, or confidence are shown when trust depends on them.
- Failure modes include fallback paths not dependent on AI success.
- Feedback/reporting mechanisms are available for incorrect or harmful output.
- User data use is explained and minimized.
- Latency is handled with progress, cancel, and partial results where useful.
- AI-generated content is distinguishable from human-authored content when authorship matters.

### Design system criteria

Detail: [10-design-systems.md](10-design-systems.md).

- Component usage matches documented intent.
- Component states include default, hover, focus, active, disabled, loading, selected, error, success where relevant.
- Component APIs require accessible names or labels when needed.
- Design tokens are used for color, spacing, typography, elevation, motion, and z-index where a system exists.
- Deprecated components are not introduced unless migration is impossible and documented.
- Custom variants document why existing variants are insufficient.

### Research and measurement criteria

Detail: [17-ux-process-research.md](17-ux-process-research.md).

- High-risk flows have usability testing or documented equivalent evidence.
- Accessibility testing includes automated, manual keyboard, and assistive technology checks where feasible ([05-accessibility.md](05-accessibility.md#testing-methodology)).
- Metrics are tied to user outcomes, not only business conversion.
- Analytics events do not collect unnecessary personal data.
- Findings are prioritized by severity, frequency, and user/business impact.

### PR review checklist for agents

```md
UI states complete: yes/no
Keyboard path complete: yes/no
Accessible names/roles complete: yes/no
Focus management complete: yes/no
Responsive behavior complete: yes/no
Error/recovery states complete: yes/no
High-risk safeguards complete: yes/no/not applicable
Platform conventions followed: yes/no
Known tradeoffs documented: yes/no
Required follow-up tests: ...
```

A completed checklist from an AI agent is a review artifact for a human to verify, not sign-off ([17-ux-process-research.md](17-ux-process-research.md#design-review-rubrics-and-ai-assisted-review)).

---

## 5. Principle tradeoff tables

Compact decision support for selecting between principles or resolving conflicts. Column meanings: **Principle** (the rule), **Description** (what it means in practice), **Use when** (usually appropriate), **Avoid/limit when** (blind use can harm), **Pros**/**Cons** (likely benefits and costs), **Reasoning** (why it exists), **Verify** (how to check application). The evidence and nuance behind each row live in the linked files.

### Universal principles

Deep dive: [01-core-principles.md](01-core-principles.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Visibility of system status | Communicate current state, progress, result, and next action. | Any action has consequence, delay, async work, save state, or failure risk. | Low-risk changes that are already obvious. | Builds trust, prevents duplicate actions, improves recovery. | Too many messages create noise; inaccurate progress damages trust. | Users need feedback to connect action to outcome. | Slow network, failed request, screen reader status, and repeated-click behavior are tested. |
| Match user's world | Use user language, mental models, and natural task order. | Labels, navigation, workflows, domain objects, onboarding, errors. | Real-world metaphor is outdated, misleading, or culturally narrow. | Improves comprehension and learnability. | Multiple user groups may use different terms. | Users should not translate internal system concepts. | New users can explain labels and workflow without internal knowledge. |
| User control and freedom | Provide undo, cancel, back, exit, edit, and recovery. | Multi-step flows, destructive actions, editing, dialogs, modes. | Recovery is impossible and would be dishonest. | Reduces anxiety, supports exploration, lowers error cost. | Undo can be technically complex; too many exits can confuse regulated flows. | Users make mistakes and need safe recovery. | Accidental navigation, dialog closure, destructive action, and unsaved changes are tested. |
| Consistency and standards | Reuse familiar patterns internally and externally. | Common controls, navigation, shortcuts, forms, platform-specific behavior. | Existing pattern demonstrably harms users or conflicts with accessibility. | Improves learnability, speed, and maintainability. | Can inhibit innovation; platform conventions may differ. | Users transfer knowledge from other products. | Same labels/actions behave the same; platform conventions are followed or deviations documented. |
| Error prevention | Remove, constrain, warn, or guide before errors happen. | High-risk actions, forms, invalid combinations, permissions, destructive work. | Constraints block legitimate workflows or fire while user is still typing. | Reduces failures, support, and user frustration. | Can feel restrictive; confirmations can become ignored. | Prevention is less costly than recovery. | Invalid states are hard to reach; high-risk actions have review, undo, or confirmation. |
| Recognition over recall | Keep options, state, and required information visible or retrievable. | Menus, forms, comparisons, multi-step flows, configuration. | Showing everything creates clutter or hides priority. | Reduces cognitive load and supports interruption. | Takes space; can overwhelm if unstructured. | Working memory is limited ([02-cognitive-foundations.md](02-cognitive-foundations.md)). | Users do not need information from previous screens to continue. |
| Aesthetic and minimalist design | Remove irrelevant elements and prioritize essentials. | Any screen with competing content or actions. | Minimalism hides labels, affordances, instructions, or required controls. | Improves focus, scanning, and perceived quality. | Can reduce discoverability and density. | Every extra element competes for attention. | Primary task is visually obvious; removed elements are not needed for completion. |
| Plain-language error recovery | Explain problem, location, and fix in human terms. | Validation, permissions, network failure, conflicts, unavailable actions. | Security requirements require limited detail; still provide safe next step. | Faster recovery, less blame, more trust. | Detailed messages can be long; security details may leak information. | Users cannot recover from vague failures. | Errors preserve data, identify issue, and give action. |
| Contextual help | Provide concise help near difficult decisions. | Complex concepts, risky choices, unfamiliar fields, enterprise tools. | Help text repeats obvious labels or clutters simple flows. | Reduces support and errors. | Excess help increases reading burden. | Help is most useful at the moment of need. | Users can complete task without leaving flow for basic questions. |

### Cognitive principles

Deep dive: [02-cognitive-foundations.md](02-cognitive-foundations.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Reduce cognitive load | Remove unnecessary memory, comparison, translation, and ambiguity. | Complex decisions, forms, dashboards, new-user flows. | Simplification hides essential expert controls. | Faster comprehension, fewer errors, broader inclusion. | Can add steps or reduce density. | Users have limited attention and working memory. | Users can state purpose and next action quickly. |
| Progressive disclosure | Reveal details and advanced controls as needed. | Advanced settings, optional filters, setup, expert features. | Information is mandatory, risky, or affects consent/pricing. | Reduces overwhelm while preserving depth. | Hidden features are less discoverable. | Users need different levels of detail at different times. | Novices can start; experts can still reach advanced controls. |
| Limit choice complexity | Group, prioritize, or filter choices. | Large menus, plan selection, filters, command lists. | Removing options prevents valid tasks. | Faster decisions and easier scanning. | Poor grouping creates ambiguity. | Choice time rises with number and similarity of options (Hick's law; [02-cognitive-foundations.md](02-cognitive-foundations.md)). | Users can distinguish options and find likely choice. |
| Make targets easy to acquire | Size and place controls for the input method. | Touch, frequent actions, destructive controls, dense lists. | Dense expert tools need compactness but must preserve operability. | Fewer accidental taps/clicks, faster interaction. | Larger targets consume space. | Motor effort affects task success (Fitts's law). | Target size and spacing meet platform/WCAG expectations (24px floor, SC 2.5.8). |
| Respect mental models | Present objects, relationships, and actions as users expect. | Domain tools, onboarding, settings, workflow design. | User expectations are unsafe or impossible; teach clearly instead. | Improves prediction and error prevention. | Different roles may have different models. | Users act based on their model of the system. | Users can predict outcome before acting. |

### Information architecture and navigation

Deep dive: [08-navigation-ia.md](08-navigation-ia.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Task-based IA | Structure around user goals, not org charts. | Navigation, settings, docs, service flows. | Compliance or domain constraints require a legal structure; still translate labels. | Better findability and completion. | May not match backend ownership. | Users come to complete tasks, not understand organizations. | Tree testing or usability testing confirms users find destinations ([17-ux-process-research.md](17-ux-process-research.md)). |
| Visible location and scope | Show current page, object, account, project, environment, filters. | Multi-tenant, admin, dashboards, editors, deep navigation. | Single-screen simple utilities. | Prevents wrong-object/wrong-environment errors. | Extra headers consume space. | Users need orientation to act safely. | User can identify where they are and what object actions affect. |
| Appropriate navigation pattern | Match top nav, side nav, tabs, breadcrumbs, search, or bottom nav to structure. | Designing app or section navigation. | Copying a trendy pattern despite content mismatch. | Improves wayfinding and platform fit. | Wrong pattern creates hidden depth or clutter. | Navigation is a map of user tasks. | Navigation supports top tasks within expected clicks/taps. |
| Search with recovery | Search should reveal scope, tolerate realistic queries, and recover from no results. | Large content, known-item finding, catalogs, docs. | Small information spaces where browse is faster. | Speeds known-item retrieval. | Search quality is hard; poor results damage trust. | Users often know what they want but not where it lives. | Query, filter, empty result, and return-to-results flows work. |

### Interaction and input

Deep dive: [04-interaction-design.md](04-interaction-design.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Direct manipulation feedback | Show response to click, tap, drag, type, select, resize. | Any interactive control or object. | Purely passive content. | Feels responsive and understandable. | Many states require design/engineering discipline. | Users need to know input was received. | Hover/press/focus/selected/loading/disabled states exist. |
| Keyboard operability | Ensure keyboard access and logical focus. | Web, desktop, enterprise, accessibility, power-user workflows. | Touch-only hardware with no keyboard path, but avoid excluding assistive tech. | Accessibility and expert efficiency. | Custom keyboard behavior can conflict with platform shortcuts. | Keyboard is both assistive and productive. | Full task can be completed keyboard-only ([05-accessibility.md](05-accessibility.md#keyboard-navigation)). |
| Gesture alternatives | Do not make hidden gestures the only path. | Mobile, touch, drag/drop, canvas tools. | Gesture is supplemental and non-critical. | Improves accessibility and discoverability. | Visible alternatives take space. | Gestures are hard to discover and inaccessible to some users. | Swipe/drag/long-press tasks have buttons/menu alternatives (SC 2.5.7). |
| Clear modes | Make modes visible and escapable. | Edit mode, bulk select, drawing, admin impersonation. | Mode can be avoided with modeless interaction. | Prevents mode errors. | Persistent indicators consume space. | Same input doing different things causes mistakes ([04-interaction-design.md](04-interaction-design.md#modes-and-mode-errors)). | User can identify mode and exit it. |
| Prefer enabled controls over unexplained disabled ones | Default to enabled actions that validate and explain on activation; if a control must be disabled, show the prerequisite nearby. | Temporarily unavailable actions, gated submissions. | Activation would cause harm even with validation; then disable with visible reason. | Users always get an explanation path; prevents dead ends. | Validation-on-activation requires good error messaging. | Unexplained disabled controls strand users; users need to know how to proceed. | No disabled control lacks a nearby or discoverable reason; enabled-submit paths explain failures actionably. |

### Visual design

Deep dive: [03-visual-design.md](03-visual-design.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Visual hierarchy | Use visual attributes to show importance and relationships. | Every screen. | Over-styling creates false importance. | Faster scanning and action selection. | Requires restraint and consistency. | Users scan before reading. | Primary content/action stands out at a glance. |
| Proximity and grouping | Place related items together and separate unrelated items. | Forms, cards, toolbars, dashboards. | Groups become too large; then subdivide. | Improves comprehension and reduces errors. | Extra spacing may reduce density. | Users infer relationships from spatial layout (Gestalt; [03-visual-design.md](03-visual-design.md)). | Labels, controls, and actions are near affected content. |
| Readable typography | Use type for legibility, hierarchy, and localization. | Text-heavy and form-heavy UI. | Brand display type harms readability. | Better comprehension and accessibility. | Larger type consumes space. | Text is the primary UI material. | Text scales, wraps, and remains readable — including with 100–300% expansion on short labels ([09-content-ux-writing.md](09-content-ux-writing.md#internationalization-and-localization)). |
| Color with contrast | Use color to support, not replace, meaning. | Status, hierarchy, branding, charts. | Color is the sole indicator or fails contrast. | Fast recognition and brand expression. | Poor palettes exclude users and fail dark/high-contrast modes. | Visual perception varies widely. | Contrast passes 4.5:1 text / 3:1 non-text (SC 1.4.3, 1.4.11; [05-accessibility.md](05-accessibility.md#contrast-minimum)); non-color cues exist. |
| Purposeful motion | Motion explains change, state, and spatial continuity. | Transitions, feedback, loading, object movement. | Motion is decorative, slow, distracting, or vestibular-triggering. | Improves orientation and perceived polish. | Can harm comfort and performance. | Motion should serve cognition. | Reduced-motion mode works and meaning remains ([04-interaction-design.md](04-interaction-design.md#motion-and-animation)). |

### Forms and validation

Deep dive: [07-forms-and-input.md](07-forms-and-input.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Ask only what is needed | Keep forms scoped to current task. | Sign-up, checkout, settings, admin data entry. | Regulatory or operational needs require more data; explain why. | Higher completion and trust. | May require later progressive collection. | Every field is a cost. | Each field has a user/business reason and timing rationale. |
| Persistent labels | Labels remain visible and programmatically associated. | All form fields. | None for required fields. | Accessibility and fewer input mistakes. | Takes more vertical space than placeholder-only design. | Users need context before, during, and after input. | Labels remain visible with values and errors ([07-forms-and-input.md](07-forms-and-input.md#labels-vs-placeholders)). |
| Right control for input | Match control to choice structure. | Any data entry. | Native control is insufficient; custom control must replicate semantics. | Faster input, fewer errors. | Wrong controls cause friction. | Controls encode constraints and expectations. | Users can complete input without extra explanation ([section 2 matrices](#2-pattern-selection-matrices)). |
| Helpful validation timing | Validate at a moment that helps correction. | Format fields, required fields, server checks. | Premature validation during incomplete typing. | Faster correction. | Debounced/server validation can be complex. | Feedback timing affects flow. | Errors appear neither too early nor too late ([07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing)). |
| Preserve data on error | Never punish failed submit by clearing work. | Any form or editor. | Sensitive data may need secure handling; preserve safely where possible. | Trust and recovery. | Requires state management. | Data loss is a severe usability failure. | Failed submit keeps user input. |

### Platform-specific principles

Deep dives: [13-platform-web.md](13-platform-web.md), [14-platform-desktop.md](14-platform-desktop.md), [15-platform-mobile.md](15-platform-mobile.md), [16-platform-cli-tui.md](16-platform-cli-tui.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Preserve web fundamentals | Use semantic HTML, URLs, browser history, links/buttons correctly. | Web apps and sites. | Highly specialized canvas/app experiences; still expose accessibility paths. | Accessibility, interoperability, user trust. | SPA complexity can break fundamentals. | The browser is part of the UX. | Back, zoom, keyboard, forms, and screen readers work. |
| Respect desktop conventions | Follow OS windows, menus, shortcuts, files, notifications. | Desktop apps. | Cross-platform framework imposes differences; document and mitigate. | Familiarity and productivity. | More platform-specific implementation. | Desktop users expect OS integration. | Native conventions and accessibility settings are respected. |
| Design for mobile context | Optimize touch, reach, interruptions, connectivity, and text entry. | Mobile apps and mobile web. | Tablet/desktop contexts need richer density. | Better completion in real contexts. | Can oversimplify complex tasks. | Mobile use is constrained and interrupted. | One-handed, offline/poor network, rotation, and text scaling are tested. |
| Honor CLI conventions | Human-first output, machine-readable modes, exit codes, TTY detection, composability. | CLI tools, TUIs, developer tooling, automation surfaces. | The tool is purely interactive/visual with no scripting use. | Scriptability, trust, ecosystem fit. | Two output modes to maintain. | The shell is a platform with its own conventions. | Piped, scripted, and interactive use all work; exit codes are meaningful ([16-platform-cli-tui.md](16-platform-cli-tui.md)). |
| Cross-platform adaptation | Keep concepts coherent while adapting controls and behavior. | Products across web, desktop, mobile. | Forcing identical UI harms platform fit. | Brand and conceptual continuity. | Requires multiple variants. | Users expect both product familiarity and platform familiarity. | Same task feels coherent but platform-native. |

### Trust, privacy, and ethics

Deep dives: [18-patterns-antipatterns.md](18-patterns-antipatterns.md), [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Consequence clarity | Explain costs, visibility, commitments, data use, and irreversible effects. | Purchases, sharing, deletion, permissions, account changes. | Low-risk obvious actions. | Trust and fewer severe mistakes. | More text can slow flow. | Consent requires understanding. | Users can state what will happen before acting. |
| Privacy by design | Minimize data collection and explain sensitive requests. | Any personal or sensitive data. | Legal requirements require collection; still explain and protect. | Trust, compliance, reduced harm. | May limit personalization. | Data exposure has real consequences. | Each data request has clear purpose and control. |
| Secure by default | Make safe choices easiest. | Authentication, sharing, permissions, admin tools. | Expert users need advanced options; keep them available but safe. | Reduces security failures. | Strong security can add friction. | Users should not need security expertise. | Password managers, accessible auth, warnings, and recovery work ([05-accessibility.md](05-accessibility.md#accessible-authentication)). |
| Avoid dark patterns | Do not deceive, coerce, obstruct, or hide costs. | All product decisions. | None. | Long-term trust and ethical compliance. | May reduce short-term conversion tricks. | User autonomy is core to good UX. | Decline/cancel/opt-out is clear and fair ([18-patterns-antipatterns.md](18-patterns-antipatterns.md)). |

### Research and evaluation

Deep dive: [17-ux-process-research.md](17-ux-process-research.md).

| Principle | Description | Use when | Avoid or limit when | Pros | Cons and tradeoffs | Reasoning | Verify |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Start with user needs | Ground decisions in observed goals and context. | New features, redesigns, high-risk flows. | Emergency fixes may use best judgment; validate later. | Better product-market and task fit. | Research takes time. | Teams are biased by internal knowledge. | Findings map to design decisions. |
| Test tasks, not screens | Evaluate whether users can complete real work. | Usability testing and QA. | Visual-only review is sufficient only for cosmetic changes. | Finds real friction. | Requires scenario design. | UX success is task success. | Representative users complete representative tasks. |
| Combine methods | Use qualitative and quantitative evidence together. | Product iteration and measurement. | Very small changes may need lightweight validation. | More reliable decisions. | More coordination. | Metrics show what; research helps explain why. | Decisions cite both behavioral signal and user evidence where possible. |

---

## Quick reference

How to use this file:

- **Designing something new?** Answer the [pattern selection questions](#pattern-selection-questions), then pick from the matrices in [section 2](#2-pattern-selection-matrices); check the row's "Required checks" and follow its deep-dive link before implementing.
- **Diagnosing a problem?** Find the symptom in [section 3](#3-failure-mode-diagnostics), apply the smallest fix, attach the row's acceptance criteria, and score severity before recommending.
- **Reviewing an implementation or PR?** Walk [section 4](#4-implementation-acceptance-criteria) group by group (universal → accessibility → platform → domain), then fill the [PR review checklist](#pr-review-checklist-for-agents).
- **Resolving a principle conflict?** Use the tradeoff tables in [section 5](#5-principle-tradeoff-tables) — the "Avoid or limit when" and "Verify" columns are the decision points.
- **In all cases:** these checklists inform judgment; they do not replace testing with users, and AI-generated verdicts require human expert sign-off ([17-ux-process-research.md](17-ux-process-research.md#design-review-rubrics-and-ai-assisted-review)).

Key thresholds worth memorizing: contrast 4.5:1 text / 3:1 large text and non-text (SC 1.4.3, 1.4.11) · targets 24×24 CSS px floor (SC 2.5.8), ~44–48px on touch · Core Web Vitals LCP ≤ 2.5s, INP ≤ 200ms, CLS ≤ 0.1 at p75 · text expansion ~30–50% for paragraphs, 100–300% for short labels.
