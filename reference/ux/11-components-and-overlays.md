# Components & Overlays

## Scope

This file is the behavioral contract reference for common UI components (Part 1) and for layered/overlay UI — modals, alert dialogs, drawers, popovers, tooltips (Part 2). Each entry answers: when to use it, when not to, what states and behaviors it must implement, what content rules apply, and how it typically fails. Per-widget ARIA roles, attributes, and full keyboard specs live in [06-aria-widget-reference.md](06-aria-widget-reference.md); the generic component-state matrix (default/hover/focus/…), loading strategies, and notifications/toasts live in [04-interaction-design.md](04-interaction-design.md); form-field behavior and validation live in [07-forms-and-input.md](07-forms-and-input.md); tables and data-heavy views live in [12-data-tables-dashboards.md](12-data-tables-dashboards.md); modal overuse and dark patterns are cataloged in [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

## Contents

- [Part 1 — Component pattern library](#part-1--component-pattern-library)
  - [How to read this library](#how-to-read-this-library)
  - [Buttons](#buttons) · [Links](#links) · [Text inputs](#text-inputs) · [Checkbox](#checkbox) · [Radio group](#radio-group) · [Select / dropdown](#select--dropdown) · [Combobox / autocomplete](#combobox--autocomplete) · [Toggle / switch](#toggle--switch) · [Menus](#menus) · [Tabs](#tabs) · [Accordion / disclosure](#accordion--disclosure) · [Cards](#cards) · [Alerts, banners, inline messages](#alerts-banners-inline-messages) · [Progress indicators](#progress-indicators) · [Empty states](#empty-states)
  - [Component quality checklist](#component-quality-checklist)
- [Part 2 — Overlays and layered UI](#part-2--overlays-and-layered-ui)
  - [Layer selection matrix](#layer-selection-matrix)
  - [Modal dialog](#modal-dialog) · [Alert dialog / confirmation](#alert-dialog--confirmation) · [Drawer / side panel](#drawer--side-panel) · [Popover / inline dialog](#popover--inline-dialog) · [Tooltip](#tooltip)
  - [Backdrop and z-index rules](#backdrop-and-z-index-rules)
  - [Overlay acceptance criteria](#overlay-acceptance-criteria)
- [Quick reference](#quick-reference)

---

## Part 1 — Component pattern library

### How to read this library

Each entry follows one template: **Use when / Avoid when / Behavior contract / Content rules / Failure modes / Verification.** The behavior contract states what the component must do, not how to build it — implementation-level ARIA and keyboard specs are in [06-aria-widget-reference.md](06-aria-widget-reference.md). The failure modes and decision questions are the core value: they encode the mistakes that recur across real products.

Two corpus-wide stances apply throughout:

1. **Disabled controls.** The default for form submission is to **keep the button enabled** and validate on activation, showing an error summary and field-level messages (GOV.UK design system stance; see [18-patterns-antipatterns.md#disabled-submit-without-explanation](18-patterns-antipatterns.md#disabled-submit-without-explanation)). Disabling is the fallback for cases where activation is genuinely impossible or harmful (in-flight requests, missing prerequisites outside the form) — and then it must come **with a visible, adjacent explanation** of what unlocks it, never a bare grey button. Momentary in-flight disabling with a spinner is fine ([04-interaction-design.md#component-states](04-interaction-design.md#component-states)).
2. **State completeness.** A component without designed states is unfinished. The canonical interactive set is **default, hover (pointer only), focus-visible, active/pressed, disabled, loading**, plus error/selected/read-only where applicable — full matrix and styling rules in [04-interaction-design.md#component-states](04-interaction-design.md#component-states) and [10-design-systems.md](10-design-systems.md).

### Buttons

**Use when:** the control performs an action — submits, executes a command, opens a dialog, mutates state.

**Avoid when:** the control navigates to a URL or route → use a [link](#links).

**Behavior contract:**

- Native `<button>` (or platform button) wherever possible; Enter and Space both activate.
- States: default, hover (where a pointer exists), focus-visible, active/pressed, disabled, loading. **Success/error feedback usually belongs to the surrounding flow** (inline message, redirect, toast — [04-interaction-design.md#notifications-and-interruption-design](04-interaction-design.md#notifications-and-interruption-design)), not to the button itself; a brief button-level success state ("Saved ✓") is an optional nicety for async in-place actions, never the sole feedback channel.
- Loading state prevents duplicate submission or makes duplicates harmless (idempotency); preserve button width and change the label ("Saving…") rather than collapsing to a spinner alone.
- Submit buttons: enabled by default; validate on activation (stance 1 above).
- Icon-only buttons have an accessible name and, ideally, a visible label or tooltip.

**Content rules:** the label names the outcome — "Delete project", "Save changes", not "OK", "Submit", or "Click here". Destructive actions get explicit labels and visual/spatial separation from safe actions ([07-forms-and-input.md#destructive-action-confirmation](07-forms-and-input.md#destructive-action-confirmation)).

**Decision questions:**

- Does it mutate state? → button. Does it navigate? → link. (A button styled as a link or vice versa breaks keyboard behavior, open-in-new-tab, and user expectations both ways.)
- Is it destructive? → explicit label, separation, and undo-vs-confirm decision ([04-interaction-design.md#undoredo-vs-confirmation-dialogs](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- Is it async? → define the loading state and the duplicate-click behavior before shipping.

**Failure modes:** generic labels hiding the outcome; disabled button with no explanation and no path to enabling; icon-only button with no accessible name; loading state that lets a second click fire the action twice.

**Verification:** activate by Enter, Space, and pointer; check every state renders; double-click a slow async action; read the label out of context — does it say what happens?

### Links

**Use when:** navigation — pages, routes, anchors, documents, external destinations.

**Avoid when:** the control mutates state or executes a command → button. `<a href="#" onclick=…>` for a data mutation breaks middle-click, copy-link, and history semantics.

**Behavior contract:**

- Native `<a href>` so open-in-new-tab, copy address, and history work.
- States: default, hover, focus-visible, active; visited where the content benefits (documentation, search results).
- Links remain recognizable without relying on color alone (underline or equivalent — WCAG 1.4.1).
- New-window/external behavior disclosed where it would surprise (icon + accessible-name suffix).

**Content rules:** link text describes the destination or purpose and works out of context — screen-reader users navigate by link lists. Never repeat bare "Learn more" / "Read more" / "Here" without differentiating context ([09-content-ux-writing.md](09-content-ux-writing.md)).

**Failure modes:** link styled as button performing a mutation; ambiguous repeated link text; color-only link affordance; `target="_blank"` without disclosure.

**Verification:** middle-click and Ctrl/Cmd+click behave sanely; tab to each link and confirm a visible focus indicator; scan the page's link texts in isolation.

### Text inputs

**Use when:** free-form user-entered text. **Avoid when:** the value comes from a small fixed choice set — use radio/checkbox/select ([07-forms-and-input.md#input-type-selection](07-forms-and-input.md#input-type-selection)).

**Behavior contract:**

- States: empty, focused, filled, disabled, read-only, required, invalid/error, plus help text; success indication where confirmation genuinely helps.
- Persistent visible label, programmatically associated; placeholder is never the label ([07-forms-and-input.md#labels-vs-placeholders](07-forms-and-input.md#labels-vs-placeholders)).
- Errors programmatically associated with the field ([07-forms-and-input.md#error-message-design](07-forms-and-input.md#error-message-design)); validation timing per [04-interaction-design.md#validation-feedback-timing](04-interaction-design.md#validation-feedback-timing).
- On web: correct `type`, `inputmode`, and `autocomplete` attributes.
- Accept forgiving formats (spaces in card numbers, mixed phone formats) — normalize on the backend, don't reject.

**Failure modes:** placeholder-only labeling; overly strict format rejection; input value cleared on validation error (never destroy user input); error message not announced to AT.

**Verification:** submit an invalid value — is the input preserved and the error associated? Zoom to 200%; check label survives autofill styling.

### Checkbox

**Use when:** independent boolean choices, "select all that apply" multi-select, explicit agreements. **Avoid when:** options are mutually exclusive → [radio group](#radio-group); the setting takes effect immediately → [toggle](#toggle--switch).

**Behavior contract:**

- States: checked, unchecked, focus-visible, disabled, error; **mixed/indeterminate** when a parent checkbox aggregates children (parent shows mixed when children disagree; activating a mixed parent checks all, then toggles all).
- Label clickable; adequate touch target ([05-accessibility.md](05-accessibility.md)).

**Content rules:**

- Phrase labels so *checked = affirmative* and avoid negations ("Don't not send emails" bugs).
- **Never pre-check consent** for privacy, marketing, or legal agreement — pre-selection of consent is a dark pattern and a GDPR problem ([18-patterns-antipatterns.md#preselection](18-patterns-antipatterns.md#preselection)).

**Failure modes:** pre-checked consent; checkbox used for either/or choices; mixed state that isn't conveyed programmatically (`aria-checked="mixed"` — [06-aria-widget-reference.md](06-aria-widget-reference.md)).

**Verification:** toggle by Space and by clicking the label; if a parent aggregate exists, verify the mixed state visually and via AT.

### Radio group

**Use when:** exactly one choice from a small, visible set the user should see and compare at a glance. **Avoid when:** there are many options or search is needed.

**Behavior contract:**

- Group has a label/legend; each option has a visible label; one tab stop for the group with arrow-key movement between options (native behavior; required if custom — [06-aria-widget-reference.md](06-aria-widget-reference.md)).
- A default selection is optional; if the question is consequential (billing frequency, consent-adjacent), consider no default and validate on submit.

**Sizing rule:** radios shine at roughly **2–7 options**; common guidance is to reconsider around **~5–7** — beyond that, scanning cost grows and a different control may serve better. Even then, **select/dropdown remains a last resort**, not the automatic next step — consider a combobox, filtering, or restructuring the question first ([07-forms-and-input.md#input-type-selection](07-forms-and-input.md#input-type-selection)).

**Failure modes:** using a dropdown for 3 important options users should compare side-by-side (hides the choice space behind a click); missing group label so AT reads options without context; custom radios without arrow-key behavior.

**Verification:** tab into the group (one stop), arrow between options, confirm the legend is announced.

### Select / dropdown

**Use when:** many familiar options where side-by-side comparison isn't important (country, state, year). **Avoid when:** the set is small and comparable (→ radio), the set is large and searchable (→ combobox), or the items are *actions* (→ [menu](#menus)).

**Behavior contract:**

- Prefer the native element — free keyboard support, mobile pickers, AT compatibility.
- Visible persistent label; a meaningful default or an explicit placeholder option that cannot be submitted as a value.
- Selected value visible when closed.

**Decision rule (choice controls):** radio for 2–7 important visible choices → combobox for searchable/suggested large sets → select for familiar long lists where typing isn't natural. Full decision table: [07-forms-and-input.md#input-type-selection](07-forms-and-input.md#input-type-selection).

**Failure modes:** select as navigation (jump-menus that navigate on change break keyboard scanning); select for binary choices; 200-item unsearchable select; custom select that drops keyboard type-ahead.

**Verification:** open and choose with keyboard only; check behavior on a touch device (native picker vs. custom listbox).

### Combobox / autocomplete

**Use when:** large option sets, typeahead search, suggestion-assisted entry. **Avoid when:** the choice set is small and fixed — simpler controls win.

**Behavior contract (W3C APG combobox pattern):**

- Typing filters/suggests; arrow keys move an active option; Enter selects; Escape closes the list without selecting.
- Visually distinct highlighted state that is also conveyed programmatically (`aria-activedescendant` or focus — [06-aria-widget-reference.md](06-aria-widget-reference.md)).
- Explicit **no-results state** with guidance, never a silent empty list.
- If free-text entry is legitimate, allow values not in the list; if not, say so on rejection.
- Announce result counts to AT on filter.

**Failure modes:** custom combobox unusable with a screen reader (the single most-failed custom widget); user cannot enter a valid value that isn't suggested; suggestions obscured behind the on-screen keyboard on mobile; slow remote filtering with no pending indication.

**Verification:** complete a selection with screen reader on at least one platform; type a value not in the list; test with a throttled network.

### Toggle / switch

**Use when:** an on/off setting that takes effect **immediately** (no separate save step). **Avoid when:** the boolean is part of a form submitted later (→ checkbox) or represents agreement/consent (→ checkbox, unchecked).

**Behavior contract:**

- Immediate effect; if the change is async, show pending state and revert with an error message on failure.
- On/off state conveyed by more than color (position + label or text).

**Content rules:** **the label describes the state/thing controlled, not the action** — "Email notifications" (with the switch showing on/off), not "Turn on email notifications" (now what does "on" mean when it's already on?). Never use a toggle where the user must guess whether the label reflects the current state or the state after toggling (NN/g toggle-switch guidance).

**Failure modes:** toggle inside a form that still requires "Save" (checkbox territory); action-phrased ambiguous labels; state indicated by color alone; silent failure of the async change.

**Verification:** read the label with the switch in both positions — is the current state unambiguous? Kill the network and flip it — does it revert with a message?

### Menus

**Use when:** command lists — overflow ("…") actions, context menus, action groups on an item. **Avoid when:** it's site/app navigation (use nav links; `role="menu"` semantics are for command menus) — unless the platform convention is menus (desktop menu bars, [14-platform-desktop.md](14-platform-desktop.md)).

**Behavior contract (W3C APG menu/menu-button pattern):**

- Trigger exposes expanded/collapsed state (`aria-expanded`); arrows move through items; Escape closes and returns focus to the trigger ([06-aria-widget-reference.md](06-aria-widget-reference.md)).
- Destructive items visually separated (separator + danger styling) from safe items.
- Prefer disable-with-reason over hiding items whose preconditions aren't met — vanished items confuse; disabled items teach "exists, not now" ([14-platform-desktop.md](14-platform-desktop.md)).

**Content rules:** items are verbs naming their action; items opening dialogs end with "…" per desktop convention.

**Failure modes:** menu semantics slapped on navigation links (breaks AT expectations); Escape not closing; focus lost on close; destructive item adjacent to a common safe item with no separation.

**Verification:** open with keyboard, arrow through, Escape out — focus back on trigger; screen reader announces "menu" context and item count.

### Tabs

**Use when:** peer panels of alternative content in the same context — one visible at a time, same page. **Avoid when:** the steps are sequential (→ wizard/stepper), the sections are unrelated navigation (→ nav links/routes), or users need to compare panel contents side by side.

**Behavior contract (W3C APG tabs pattern):**

- Selected tab visually distinct and programmatically exposed (`aria-selected`); panel labelled by its tab.
- Arrow keys move between tabs; Tab moves into the panel ([06-aria-widget-reference.md](06-aria-widget-reference.md)).
- Deep-linkable state where content is share-worthy (URL reflects selected tab on web).

**Failure modes:** tabs used for wizard steps (users expect free order in tabs); form fields split across tabs and validated together (errors invisible in unselected tabs); selected state conveyed by color only; tab list wrapping/scrolling with no overflow affordance.

**Verification:** operate entirely by keyboard; trigger a validation error in a background tab — does anything indicate where the problem is?

### Accordion / disclosure

**Use when:** optional details, FAQ-style content, advanced settings — progressive disclosure of *secondary* content. **Avoid when:** the content is required to proceed or to decide.

**Behavior contract (W3C APG disclosure pattern):**

- Header is a real button carrying `aria-expanded`; Enter/Space toggles; the expanded/collapsed state has a visible indicator (chevron etc.).
- Collapsed content must be genuinely optional: **never hide pricing, fees, consent terms, required fields, or critical warnings in collapsed panels** — hiding costs behind interaction is the "drip pricing" dark pattern ([18-patterns-antipatterns.md#hidden-costs-and-drip-pricing](18-patterns-antipatterns.md#hidden-costs-and-drip-pricing)).
- Decide deliberately whether multiple panels can be open at once (usually yes); auto-collapsing siblings punishes comparison.
- Expanded state must survive printing and in-page find where feasible.

**Failure modes:** required fields collapsed at validation time (error points to invisible content); expanded state not conveyed to AT; whole-header click target but only the icon is interactive.

**Verification:** toggle by keyboard; trigger an error inside a collapsed panel — does it open/scroll/announce? Check nothing consequential lives only in a collapsed panel.

### Cards

**Use when:** heterogeneous browsable items — mixed media, summary + action group, visual scanning. **Avoid when:** users must compare many attributes across many items — that's a table ([12-data-tables-dashboards.md](12-data-tables-dashboards.md)).

**Behavior contract:**

- Clear title (a real heading), summary content, and a **deliberate click model**: either the whole card is one link with a single well-defined destination, or the card contains discrete actions — **not both ambiguously**. The whole-card-link ambiguity is the classic failure: users can't tell what's clickable, and nesting interactive elements inside a link breaks keyboard and AT behavior.
- Recommended pattern: the title is the real link; the card's hit area is extended via CSS (pseudo-element), and secondary actions sit above it in z-order. One predictable focus target per card, not five nested ones.
- Text must be selectable — a whole-card click handler that fires on text selection is hostile.

**Failure modes:** card is "sort of clickable" (some regions navigate, others don't, no affordance); nested links/buttons inside an anchor; card grid unusable by keyboard because every card is a div with a click handler.

**Verification:** tab through a card — how many stops, and does each make sense? Try to select text; middle-click the card.

### Alerts, banners, inline messages

**Use when:** status, warnings, errors, and contextual information tied to a page or region. For transient event notifications (toasts) and interruption budgets, see [04-interaction-design.md#notifications-and-interruption-design](04-interaction-design.md#notifications-and-interruption-design). **Avoid when:** it's decoration or marketing dressed as system status.

**Severity mapping (use all five levels deliberately):**

| Level | Meaning | Typical treatment |
| --- | --- | --- |
| Info | Neutral context; nothing wrong | Quiet inline/banner, dismissible |
| Success | Completed result confirmation | Inline or toast; auto-dismiss acceptable |
| Warning | Risk *before* an action; degraded but working | Persistent until resolved or acknowledged |
| Error | Failure requiring user attention | Persistent, actionable, associated with cause ([04-interaction-design.md#error-feedback-and-error-taxonomy](04-interaction-design.md#error-feedback-and-error-taxonomy)) |
| Critical | Service/security/legal interruption | Blocking or page-level; never auto-dismisses |

**Behavior contract:** severity conveyed by icon + text, never color alone; errors that appear dynamically are announced (live region — [06-aria-widget-reference.md](06-aria-widget-reference.md)); dismissible messages don't take dismissal as resolution of the underlying problem.

**Failure modes:** everything styled as a red error (severity inflation trains users to ignore all of it); warning used where error is meant; success toast as the only confirmation of a high-stakes action; message with no next action.

**Verification:** does each message name the cause and the next step? Turn on a screen reader and trigger a dynamic error.

### Progress indicators

**Use when:** the system is loading, processing, or making the user wait. Full loading strategy (thresholds, optimistic UI) in [04-interaction-design.md#loading-strategies](04-interaction-design.md#loading-strategies).

**Selection rule:**

- **Determinate bar** when progress is measurable (upload bytes, steps completed) — always prefer it when you can measure.
- **Indeterminate bar** when duration is unknown but the operation is ongoing and page/section-scoped.
- **Skeleton screen** when the incoming content's structure is predictable — it sets layout expectations and reduces perceived wait.
- **Spinner** only for **short, localized** waits (a button, a small region); a full-page spinner is a last resort.

**Behavior contract:** appear only past the ~1 s perceived-delay threshold (flashing indicators for fast operations add noise — [04-interaction-design.md#feedback-and-response-time-thresholds](04-interaction-design.md#feedback-and-response-time-thresholds)); long operations (>10 s) get progress + time estimate + cancel where possible; state exposed programmatically.

**Failure modes:** fake progress bars ([18-patterns-antipatterns.md#fake-loading](18-patterns-antipatterns.md#fake-loading)); indeterminate spinner for a 2-minute operation with no reassurance; skeleton that doesn't match the loaded layout (causes shift); spinner with no timeout/error path.

**Verification:** throttle the network — does every wait have an indicator, and does every indicator have an error/timeout path?

### Empty states

**Use when:** a view has no content — first use, cleared filters, no permissions, no search results, post-deletion, or an error. Design one per *cause*; a generic "No data" serves none of them. Fuller treatment: [04-interaction-design.md#empty-states](04-interaction-design.md#empty-states).

**Content rules (all four required):**

1. **Cause** — why is this empty ("No results for 'foo'", "You haven't created any projects yet").
2. **Next action** — the primary thing to do now (create, clear filters, request access).
3. **Recovery path** — how to get back/out if this is unexpected.
4. **Tone matched to context** — playful is fine for first-use, wrong for "your data failed to load".

**Failure modes:** blank region indistinguishable from a loading failure; empty state for filtered-to-zero that doesn't offer "clear filters"; illustration with no action.

**Verification:** enumerate the causes that can produce this empty view — does each render a distinguishable, actionable state?

### Component quality checklist

Every shipped component passes all of these (definition-of-done; see also [10-design-systems.md](10-design-systems.md)):

- [ ] Correct semantic element or ARIA pattern ([06-aria-widget-reference.md](06-aria-widget-reference.md))
- [ ] Full visual **and programmatic** state matrix ([04-interaction-design.md#component-states](04-interaction-design.md#component-states))
- [ ] Complete keyboard behavior; visible focus indicator (WCAG 2.4.7)
- [ ] Screen-reader behavior verified, not assumed
- [ ] Touch target size ≥ platform minimum ([05-accessibility.md](05-accessibility.md))
- [ ] Error, empty, and loading behavior defined
- [ ] Responsive behavior at breakpoints; long-content overflow handled
- [ ] Content rules documented (label phrasing, truncation)
- [ ] Platform-convention deviations are deliberate and documented ([13-platform-web.md](13-platform-web.md), [14-platform-desktop.md](14-platform-desktop.md), [15-platform-mobile.md](15-platform-mobile.md))

---

## Part 2 — Overlays and layered UI

Layered UI trades context preservation against interruption. The recurring sins: overlays that steal or lose focus, hide required information, trap keyboard users, or interrupt when nothing needed interrupting ([18-patterns-antipatterns.md#modal-overuse](18-patterns-antipatterns.md#modal-overuse)). Choose the *lightest* layer that satisfies the need.

### Layer selection matrix

| Need | Use | Do **not** use | Reason |
| --- | --- | --- | --- |
| User must decide before continuing | Modal dialog | Tooltip, popover, toast | Requires focus and an explicit response |
| Destructive / legal / security confirmation | Alert dialog or confirmation page | Toast, inline hint | Consequence must be understood *before* the action ([07-forms-and-input.md#destructive-action-confirmation](07-forms-and-input.md#destructive-action-confirmation)) |
| Secondary task while preserving page context | Drawer / side panel | Modal (if the task isn't blocking) | User can reference the underlying content |
| Short supplemental non-interactive explanation | Tooltip | Modal, popover | Lightweight help; nothing to operate |
| Interactive contextual controls near a trigger | Popover / inline dialog | Tooltip | Tooltips must not contain controls (see [Tooltip](#tooltip)) |
| Long or complex workflow | Page / route | Modal | Pages handle navigation, saving, deep links, and errors better |
| Optional details within the flow | Disclosure ([accordion](#accordion--disclosure)) | Modal | Keeps the user in context |
| Passive event notification | Toast / status region ([04-interaction-design.md#notifications-and-interruption-design](04-interaction-design.md#notifications-and-interruption-design)) | Modal, alert dialog | No decision required — don't block |

Decision heuristic: **must the user act before anything else?** No → don't use a modal. Yes, and the action is destructive/irreversible → alert dialog or review page. Yes, and it's a short contained task → modal. Anything longer → page or drawer.

### Modal dialog

**Use when:** focused blocking tasks — required decisions, short forms (1–5 fields), confirmations that genuinely must interrupt. **Avoid when:** long workflows, optional information, repeated routine actions, or ambient system status.

**Behavior contract (W3C APG dialog-modal pattern):**

1. **Focus moves in on open** — to the first sensible element (first field, or the dialog container/heading for reading-first dialogs). Never leave focus behind the dialog.
2. **Background is inert** — not clickable, not focusable, not readable by AT (`inert` / `aria-modal="true"` + correct structure).
3. **Focus is trapped** while open — Tab cycles within the dialog only.
4. **Escape closes** (and a visible close control exists) unless closing mid-operation is genuinely unsafe; decide backdrop-click behavior explicitly (dismiss for low-stakes, ignore for dialogs containing user input — silent data loss on a stray click is worse than a second Escape press).
5. **Focus returns** to the trigger (or a logical successor if the trigger is gone) on close — WCAG 2.4.3 focus order applies across the open/close cycle.

Lifecycle states to design: opening, open, validation error (resolvable *inside* the modal), submitting, failure, closing.

**Content rules:** accessible name (visible title); one primary action; errors surfaced within the dialog, never behind it.

**Decision questions:** Must the user finish this before continuing? Is the content short enough? Can every error be resolved inside? What exactly happens on Escape and backdrop click?

**Failure modes:** modal opens but focus stays behind it; user can Tab to the background; native `<dialog>` opened via the `open` attribute instead of `showModal()` → non-modal, background still interactive, no focus trap, Escape dead, `aria-modal` false advertising ([06-aria-widget-reference.md](06-aria-widget-reference.md)); modal taller than the viewport with no internal scroll (buttons unreachable — test at small viewport and 200–400% zoom); destructive action receives default focus (a stray Enter deletes something); modal used because it was the easy component, not because blocking was required ([18-patterns-antipatterns.md#modal-overuse](18-patterns-antipatterns.md#modal-overuse); NN/g: modals are for when the user's full attention is genuinely required).

**Verification:** open → Tab around (trapped?) → Escape (closed? focus restored?) → reopen with screen reader (title announced? background silent?) → shrink viewport to 320 px and zoom to 400%.

### Alert dialog / confirmation

**Use when:** an important interruption requiring acknowledgement or a decision — destructive, irreversible, or high-consequence actions. **Avoid when:** the action is low-risk and reversible — **prefer undo**; routine confirmations train reflexive clicking and stop protecting anything ([04-interaction-design.md#undoredo-vs-confirmation-dialogs](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).

**Content contract — all four elements, every time:**

1. **Object** — what exactly is affected ("the project *Acme Website* and its 14 deployments").
2. **Consequence** — what will happen, concretely.
3. **Reversibility** — whether undo/recovery exists, and for how long.
4. **Explicit action label** — the confirm button names the action: **"Delete project", never a bare "OK"** or "Yes". The cancel option is a plain "Cancel", not a guilt-trip ([18-patterns-antipatterns.md#confirmshaming](18-patterns-antipatterns.md#confirmshaming)).

**Behavior contract:** modal mechanics above, plus: default focus on the *safe* action (or the dialog itself), never on the destructive button; for the highest-stakes cases add typed confirmation or a review page ([07-forms-and-input.md#destructive-action-confirmation](07-forms-and-input.md#destructive-action-confirmation)); use the alert-dialog role so AT announces it assertively ([06-aria-widget-reference.md](06-aria-widget-reference.md)).

**Decision rule:** reversible + low-risk → do it and offer undo. Irreversible or high-risk → confirmation dialog or review page. Frequent + destructive → redesign so it's reversible instead of confirming constantly.

**Verification:** read the dialog aloud — can a user who ignored the body still make the right choice from the button labels alone? Press Enter immediately on open — what fires?

### Drawer / side panel

**Use when:** a supplementary task or detail view related to the current page — edit-in-context, item details, filters — where seeing the underlying content helps. **Avoid when:** critical blocking confirmation (→ alert dialog) or unrelated navigation (→ route).

**Behavior contract:**

- Clear title and a visible close control; Escape closes.
- Decide modality explicitly: a modal drawer follows the full [modal mechanics](#modal-dialog); a non-modal drawer must still manage focus in/out sanely and not strand keyboard users.
- Responsive: on narrow screens a drawer typically becomes full-screen — plan that transformation, don't let it happen by CSS accident.
- Unsaved-changes handling defined: warn, autosave, or discard-with-undo — never silent loss on close/backdrop-click.
- Deep-linkable where the panel represents shareable state (URL reflects the open item on web).

**Failure modes:** drawer contains a form, backdrop click discards it silently; keyboard focus stays on the page behind a modal drawer; drawer content unreachable at 400% zoom; nested drawer opening over a drawer (same escalation smell as nested modals).

**Verification:** open, edit a field, click the backdrop — what happened to the data? Keyboard-only pass; 320 px viewport pass.

### Popover / inline dialog

**Use when:** interactive contextual content anchored to a trigger — filter panels, date pickers, share controls, rich previews. **Avoid when:** the info is critical (users may never open it), the form is large, or the workflow is complex (→ drawer or page).

**Behavior contract:**

- Trigger exposes expanded state; the popover's spatial and programmatic relationship to the trigger is clear (focus moves into it or it's next in reading order — [06-aria-widget-reference.md](06-aria-widget-reference.md)).
- Dismissible by Escape and (usually) outside-click; on dismiss, focus returns to the trigger.
- All content keyboard-reachable; the popover never obscures the trigger or the current focus unexpectedly.
- Positioning survives edge collisions: flips/shifts at viewport edges, does not clip inside scroll containers, does not render behind sticky headers.

**Failure modes:** popover contains the only path to required information; opens on hover with interactive content inside (hover users can't reach it — hover-triggered content follows WCAG 1.4.13: dismissible, hoverable, persistent); focus lost to `<body>` on close; popover behind a sticky header at certain scroll positions.

**Verification:** open near every viewport edge; open, Escape, confirm focus restore; try to operate it with keyboard only and with a screen reader.

### Tooltip

**Use when:** short, **non-interactive**, supplementary explanation — icon-button names, abbreviated-value expansions, minor clarifications. **Avoid when:** required instructions, error messages, form labels, interactive content, or anything mobile users must see (touch has no hover).

**Behavior contract (W3C APG tooltip pattern; WCAG 1.4.13 Content on Hover or Focus):**

- Appears on **both** hover **and** keyboard focus of the trigger.
- **Dismissible with Escape** without moving focus (1.4.13).
- **Persistent and hoverable**: stays visible while the pointer is over the tooltip itself, and until hover/focus leaves or the user dismisses (1.4.13).
- **Never receives focus itself** and contains **no interactive content** — no links, no buttons. If you need controls, it's a [popover](#popover--inline-dialog), full stop.
- Short text only (roughly a sentence); delay before showing on hover (~300–600 ms) to avoid flicker, no delay on focus.
- Programmatically associated with the trigger (`aria-describedby` — [06-aria-widget-reference.md](06-aria-widget-reference.md)); don't duplicate the accessible name into the description.

**Failure modes:** tooltip contains a button or link; the tooltip is the only way to learn a required field's format (put it in visible help text — [07-forms-and-input.md](07-forms-and-input.md)); tooltip cannot be triggered by keyboard or touch; tooltip covers the very input it describes; tooltip vanishes when the pointer moves toward it.

**Verification:** trigger by Tab alone; press Escape; hover the trigger then move the pointer *onto* the tooltip; check the content is either non-essential or duplicated somewhere touch users can reach.

### Backdrop and z-index rules

The layer stack must be **deterministic** — an ad-hoc pile of z-indexes eventually puts a tooltip behind a backdrop or a dropdown under a sticky header.

- Maintain a single ordered layer scale (e.g., page < sticky chrome < dropdown/popover < drawer < modal < toast) as design tokens ([10-design-systems.md](10-design-systems.md)); render overlays through a portal/top-layer mechanism, not wherever the DOM happens to be.
- **Only the topmost modal traps focus.** Everything beneath it is inert.
- **Escape closes the topmost dismissible layer first**, one layer per press: tooltip → popover → menu → modal, innermost outward.
- **Nested modals are strongly discouraged.** A modal spawning a modal is usually a workflow crammed into the wrong container — restructure into a page, drawer, or in-dialog step. If truly unavoidable (a destructive confirmation atop an editing modal), the inner one alone traps focus and closes first.
- Popovers/dropdowns must not clip inside `overflow: hidden` scroll containers or render behind sticky headers — position against the viewport/top layer.
- One backdrop per modal layer; stacking multiple dimming backdrops signals the nesting problem above.

### Overlay acceptance criteria

Every layered surface, before ship:

- [ ] User can **open, understand, act, dismiss, and recover** from the layer with keyboard alone *and* with pointer alone (parity — no layer is hover-only or pointer-only)
- [ ] Focus is never lost (to `<body>`) or hidden (behind a layer) at any point in the open/close cycle — WCAG 2.4.3
- [ ] No required information is tooltip-only or hover-only
- [ ] Works at **320 px viewport width** and **400% zoom** (WCAG 1.4.10 reflow): content scrolls within the layer, actions remain reachable
- [ ] Escape and close-control behavior defined and consistent; unsaved input is never silently destroyed
- [ ] Screen reader announces the layer's name and role on open; background is silent while a modal is up
- [ ] The lightest sufficient layer was chosen (justify every modal against the [matrix](#layer-selection-matrix))

---

## Quick reference

| Component / layer | One-line rule | Key check |
| --- | --- | --- |
| Button | Mutates state; label names the outcome; keep submit enabled + validate on activation | Double-click a slow async action |
| Link | Navigates; text works out of context; never color-only | Middle-click behaves sanely |
| Text input | Persistent label; forgiving formats; never clear on error | Invalid submit preserves input |
| Checkbox | Independent booleans; mixed for aggregates; never pre-check consent | Consent boxes start unchecked |
| Radio group | One of a small visible set (~2–7, reconsider past ~5–7); dropdown is last resort | Group label announced; arrows move |
| Select | Long familiar lists, no comparison needed; prefer native | Keyboard-only selection works |
| Combobox | Large/searchable sets; needs no-results state | Screen-reader pass completes |
| Toggle | Immediate effect; label describes the **state**, not the action | Label unambiguous in both positions |
| Menu | Commands, not navigation; Escape closes, focus returns | Destructive items separated |
| Tabs | Peer panels, same context; never wizard steps | Errors in hidden tabs surfaced |
| Accordion | Optional details only — never pricing/consent/required fields | Error inside collapsed panel opens it |
| Card | One deliberate click model — whole-card link *or* discrete actions | Tab stops per card make sense |
| Alert/message | Five severities, used honestly; icon + text, not color alone | Each message has cause + next step |
| Progress | Determinate > skeleton > indeterminate > spinner (short/local only) | Every wait has an error/timeout path |
| Empty state | Cause + next action + recovery + right tone, per cause | Distinguishable from load failure |
| Modal | Only when the user must act first; native `<dialog>` via `showModal()`; focus in, inert behind, trap, Escape, restore | 320 px + 400% zoom pass |
| Alert dialog | Object + consequence + reversibility + explicit label ("Delete project", not "OK") | Enter-on-open doesn't destroy |
| Drawer | Secondary task with context; unsaved changes never silently lost | Backdrop click vs. dirty form |
| Popover | Interactive contextual content; Escape + focus restore | Edge collision, no clipping |
| Tooltip | Hover **and** focus; Escape-dismissible; never focused; zero interactive content (WCAG 1.4.13) | Keyboard + touch can reach the info |
| Layer stack | Deterministic order; topmost modal traps; Escape peels topmost; no nested modals | One Escape press = one layer |
