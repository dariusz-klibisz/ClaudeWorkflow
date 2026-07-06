# ARIA & Semantic Widget Reference

Purpose: AI-agent reference for choosing native HTML vs ARIA/custom widgets and enforcing WAI-ARIA Authoring Practices Guide (APG)-derived behavior.

**Scope note.** This file covers widget semantics: which role structure, states/properties, and keyboard maps a widget must implement, and when to prefer native HTML. For WCAG success criteria, testing procedures, and broader accessibility requirements (contrast, focus visibility, reflow, screen-reader testing), see [05-accessibility.md](05-accessibility.md). For behavioral and content contracts of composed components and overlays (modals, toasts, popovers, dropdowns — what they say, when they open, how they stack), see [11-components-and-overlays.md](11-components-and-overlays.md). Related: [04-interaction-design.md](04-interaction-design.md) for interaction states, [07-forms-and-input.md](07-forms-and-input.md) for form controls in context, [12-data-tables-dashboards.md](12-data-tables-dashboards.md) for tables and grids at scale.

## Prime Directive: The First Rule of ARIA

Use native HTML controls first. Use ARIA only when native semantics cannot express the required widget and when the full keyboard, focus, name, role, state, and value contract can be implemented.

This is the "first rule of ARIA" from the W3C's [Using ARIA](https://www.w3.org/TR/using-aria/): if you can use a native HTML element or attribute with the semantics and behavior you require already built in, do so instead of re-purposing an element and adding an ARIA role, state, or property. The second rule follows directly: do not change native semantics unless you truly must (e.g., do not put `role="tab"` on an `<h2>` when a button inside the heading works).

## ARIA Risk Rules

- ARIA changes semantics, not behavior. Agents must ensure behavior is implemented separately.
- Bad ARIA is worse than no ARIA when it misleads assistive technology.
- Do not use `role="button"` on a non-button unless keyboard activation and disabled semantics are implemented.
- Do not use `role="application"` unless building a complex app that intentionally takes over screen-reader navigation.
- Do not hide focusable content with `aria-hidden="true"`.
- Do not use ARIA to paper over invalid interaction design.

## Native First Matrix

| Need | Prefer Native | Use ARIA/Custom Only If | Required Custom Behavior |
| --- | --- | --- | --- |
| Action button | `<button>` | Native button impossible due to framework constraints | Enter/Space activation, disabled, focus, name |
| Navigation | `<a href>` | Non-URL virtual command | Keyboard/focus semantics appropriate to role |
| Text input | `<input>`/`<textarea>` | Rich editor/custom input | Name, value, selection, keyboard, IME, AT support |
| Checkbox | `<input type="checkbox">` | Tri-state/custom visuals require it | checked/mixed state, Space toggles |
| Radio | `<input type="radio">` | Composite custom radio group | Arrow navigation, checked state, group label |
| Select | `<select>` | Rich searchable/multi-select needed | Combobox/listbox pattern fully implemented |
| Dialog | native `<dialog>` or robust dialog component | Native unsupported/insufficient | Modal semantics, focus trap, restore focus |
| Table | `<table>` | Virtualized/interactive grid | Grid semantics, row/col info, keyboard model |
| Disclosure | `<button>` + region or `<details>` | Custom animation/layout | `aria-expanded`, focus, state |

## Pattern Reference

Each pattern below links to its W3C APG page. Full pattern index: [ARIA Authoring Practices Guide — Patterns](https://www.w3.org/WAI/ARIA/apg/patterns/).

### Accordion

ID: ARIA-PATTERN-ACCORDION
APG: [Accordion Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/accordion/)
Use for: vertically stacked sections with show/hide content.
Avoid for: unrelated page navigation, required content hidden from users, very short content where headings suffice.

#### Required Semantics
- Header trigger is a button.
- Button has `aria-expanded` matching panel visibility.
- Button has accessible name from header text.
- Panel relationship exposed with `aria-controls` or labelled region when useful.

#### Keyboard
- Tab moves between focusable controls.
- Enter/Space toggles focused accordion header.
- Optional arrow/Home/End navigation between headers.

#### Failure Modes
- Header is a clickable `<div>`.
- Expanded state not announced.
- Hidden panel focusable while collapsed.

### Alert

ID: ARIA-PATTERN-ALERT
APG: [Alert Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/alert/)
Use for: brief important message that should be announced without moving focus.
Avoid for: routine status spam, messages requiring immediate response.

#### Required Semantics
- Use `role="alert"` or assertive live region sparingly.
- Live-region container exists in the DOM before content changes (see failure modes).
- Message text is concise and actionable.

#### Failure Modes
- Live-region container injected at the same time as the message. A live region must already exist in the DOM before its content changes; screen readers announce *changes* inside an existing live region, so inserting the `role="alert"`/live-region container together with the message often fails to announce anything. Render the empty container first, then insert or change the message text.
- Alert steals focus unnecessarily.

### Alert Dialog

ID: ARIA-PATTERN-ALERTDIALOG
APG: [Alert Dialog Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/alertdialog/)
Use for: urgent modal message requiring user response.
Avoid for: ordinary validation or low-risk confirmation.

#### Required Semantics
- Dialog role/alertdialog role.
- Accessible name and description.
- Modal focus management.
- Safe default focus unless destructive action requires explicit movement.

### Breadcrumb

ID: ARIA-PATTERN-BREADCRUMB
APG: [Breadcrumb Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/breadcrumb/)
Use for: hierarchical location trail.
Avoid for: primary navigation alone.

#### Required Semantics
- Navigation landmark labelled as breadcrumb.
- Ordered list of links.
- Current page marked with `aria-current="page"`.

### Button

ID: ARIA-PATTERN-BUTTON
APG: [Button Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/button/)
Use for: actions.
Avoid for: navigation to new URL; use link.

#### Native Recommendation
- Use `<button>`.

#### Keyboard
- Enter and Space activate.

#### State
- Toggle buttons expose pressed state.
- Disabled state is clear and explained where needed.

### Carousel

ID: ARIA-PATTERN-CAROUSEL
APG: [Carousel Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/carousel/)
Use for: small set of featured items where rotation is not essential.
Avoid for: critical information, large content sets, auto-rotating ads that distract.

#### Required Behavior
- Pause/stop controls for auto-rotation.
- Do not auto-advance while keyboard focus is inside.
- Slide controls have meaningful names.
- Current slide/state exposed.

### Checkbox

ID: ARIA-PATTERN-CHECKBOX
APG: [Checkbox Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/checkbox/)
Use for: independent choices or multi-select.
Avoid for: mutually exclusive choices; use radio.

#### Native Recommendation
- Use native checkbox.

#### Keyboard
- Space toggles.

#### State
- Supports checked/unchecked; tri-state only when parent aggregate state exists.

### Combobox

ID: ARIA-PATTERN-COMBOBOX
APG: [Combobox Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/combobox/)
Use for: input with associated popup suggestions/options.
Avoid for: small option sets visible as radio buttons or simple native select.

#### Required Semantics
- Input/combobox has accessible name.
- Popup relationship and expanded state exposed.
- Active option exposed with `aria-activedescendant` or managed focus.
- Autocomplete behavior is clear.

#### Keyboard
- Typing filters/suggests.
- Down/Up navigates options.
- Enter accepts option when appropriate.
- Escape closes popup.

#### Failure Modes
- Screen reader cannot tell popup exists.
- Arrow keys move cursor instead of list unexpectedly.
- Selected value and highlighted suggestion are confused.

### Dialog Modal

ID: ARIA-PATTERN-DIALOG-MODAL
APG: [Dialog (Modal) Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/)
Use for: focused blocking task or required decision.
Avoid for: optional help, long complex workflows better served by page, non-blocking messages.
See [11-components-and-overlays.md](11-components-and-overlays.md) for modal content and stacking contracts.

#### Required Semantics
- Dialog has accessible name.
- Background is inert/unavailable to keyboard and assistive tech.
- Focus moves inside on open.
- Focus is trapped while open.
- Focus restores to trigger or logical next point on close.
- Native `<dialog>` is opened with `showModal()`, which provides top-layer rendering, `::backdrop`, an inert background, and Escape handling for free (WHATWG HTML / MDN `<dialog>`).

#### Keyboard
- Escape closes when safe.
- Tab/Shift+Tab cycle within modal.

#### Failure Modes
- Focus remains behind modal.
- Screen reader can navigate background.
- Close button unlabeled.
- Destructive button receives default focus.
- Native `<dialog>` opened by setting the `open` attribute or calling `show()`: the result is non-modal — no focus trap, no inert background, Escape does not close — while `aria-modal="true"` then actively misreports modality to assistive technology.

### Disclosure

ID: ARIA-PATTERN-DISCLOSURE
APG: [Disclosure Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/disclosure/)
Use for: simple show/hide content.
Avoid for: complex multi-panel accordion with extra semantics unless needed.

#### Required Behavior
- Button toggles visibility.
- `aria-expanded` reflects state.
- Hidden content not focusable.

### Grid

ID: ARIA-PATTERN-GRID
APG: [Grid Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/grid/)
Use for: interactive tabular data requiring cell navigation or embedded controls.
Avoid for: static data; use native table.
See [12-data-tables-dashboards.md](12-data-tables-dashboards.md) for data-table design at scale.

#### Required Behavior
- Arrow-key cell navigation if using grid model.
- Row/column headers exposed.
- Selection/editing mode clear.
- Virtualization preserves row/column context.

#### Failure Modes
- Data table incorrectly given grid role without keyboard model.
- Screen reader loses row/column position.

### Landmarks

ID: ARIA-PATTERN-LANDMARKS
APG: [Landmark Regions](https://www.w3.org/WAI/ARIA/apg/practices/landmark-regions/)
Use for: major page regions.

#### Native Recommendation
- Use `<header>`, `<nav>`, `<main>`, `<aside>`, `<footer>`, `<form>`, `<section>` with labels where needed.

#### Acceptance Criteria
- Page has one main landmark.
- Repeated nav landmarks are labelled distinctly.

### Listbox

ID: ARIA-PATTERN-LISTBOX
APG: [Listbox Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/listbox/)
Use for: selecting one or more options from a list.
Avoid for: actions menu; use menu.

#### Keyboard
- Arrow keys navigate.
- Typeahead when useful.
- Space/Enter selection depending on model.
- Multi-select behavior documented and conventional.

### Menu And Menubar

ID: ARIA-PATTERN-MENU
APG: [Menu and Menubar Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/menubar/)
Use for: command lists, application menus, context menus.
Avoid for: site navigation lists unless menu behavior is truly needed.

#### Keyboard
- Arrow keys move between items.
- Enter/Space activates.
- Escape closes.
- Submenus have predictable open/close behavior.

#### Failure Modes
- Site nav uses menu roles and breaks normal link navigation.

### Radio Group

ID: ARIA-PATTERN-RADIO
APG: [Radio Group Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/radio/)
Use for: mutually exclusive option set.

#### Native Recommendation
- Use native radio buttons and fieldset/legend.

#### Keyboard
- Arrow keys move/select within group.

### Slider

ID: ARIA-PATTERN-SLIDER
APG: [Slider Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/slider/)
Use for: approximate value from range.
Avoid for: precise values without text alternative.

#### Required Behavior
- Keyboard increments/decrements.
- Min/max/current value exposed.
- Text value exposed when numeric value is not meaningful.

### Switch

ID: ARIA-PATTERN-SWITCH
APG: [Switch Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/switch/)
Use for: immediate on/off setting.
Avoid for: submitted boolean form value where checkbox is clearer.

#### State
- Expose on/off checked state.
- Label describes setting, not action.

### Table

ID: ARIA-PATTERN-TABLE
APG: [Table Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/table/)
Use for: static tabular data.
Avoid for: layout.
See [12-data-tables-dashboards.md](12-data-tables-dashboards.md) for table design guidance.

#### Native Recommendation
- Use `<table>`, `<th>`, captions, scope/headers as needed.

### Tabs

ID: ARIA-PATTERN-TABS
APG: [Tabs Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/tabs/)
Use for: peer panels within same context.
Avoid for: sequential steps or navigation to unrelated pages.

#### Required Semantics
- `tablist`, `tab`, `tabpanel` or native/component equivalent.
- Selected tab exposed.
- Tab controls labelled.

#### Keyboard
- Tab enters/leaves tablist.
- Arrow keys move between tabs.
- Home/End optional but recommended.
- Activation model automatic only if panel loads instantly; manual if expensive.

### Toolbar

ID: ARIA-PATTERN-TOOLBAR
APG: [Toolbar Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/toolbar/)
Use for: grouped controls for editor/app commands.

#### Keyboard
- Roving focus or standard tab sequence, documented.
- Arrow navigation where APG toolbar model used.

### Tooltip

ID: ARIA-PATTERN-TOOLTIP
APG: [Tooltip Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/tooltip/)
Use for: short noninteractive supplementary text.
Avoid for: required information, controls, long help, mobile-critical content.
See [11-components-and-overlays.md](11-components-and-overlays.md) for tooltip vs popover decisions.

#### Required Behavior
- Appears on hover and focus.
- Dismissible with Escape.
- Does not receive focus.
- Trigger has accessible relationship if needed.

### Tree View

ID: ARIA-PATTERN-TREEVIEW
APG: [Tree View Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/treeview/)
Use for: hierarchical list navigation/selection.
Avoid for: simple flat navigation.

#### Keyboard
- Arrow keys expand/collapse and move.
- Home/End support recommended.
- Typeahead recommended.

### Treegrid

ID: ARIA-PATTERN-TREEGRID
APG: [Treegrid Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/treegrid/)
Use for: hierarchical interactive tabular data.
Avoid for: static hierarchy; use tree or table.

#### Required Behavior
- Combines grid cell navigation with hierarchy expand/collapse.
- Row/column/hierarchy state exposed.

### Window Splitter

ID: ARIA-PATTERN-WINDOWSPLITTER
APG: [Window Splitter Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/windowsplitter/)
Use for: resizable panes.
Avoid for: decorative separators.

#### Required Behavior
- Keyboard resize support.
- Current value/position exposed.
- Min/max constraints respected.

## Custom Widget Acceptance Checklist

- Native alternative was considered and documented.
- Widget has accessible name.
- Role matches behavior.
- State/value exposed and updated.
- Keyboard behavior matches user expectation/APG.
- Focus order and focus visibility are correct.
- Disabled/read-only states are distinguishable.
- Pointer/touch behavior has keyboard equivalent.
- Screen reader smoke test passes.
- Zoom/reflow does not break interaction.

For WCAG-level verification and screen-reader test procedures, see [05-accessibility.md](05-accessibility.md). For agent-facing review checklists, see [21-agent-checklists.md](21-agent-checklists.md).

## Quick reference

| Widget | Native element first | Key keyboard commands | Top failure mode |
| --- | --- | --- | --- |
| Accordion | `<button>` headers (+ `<details>` for simple cases) | Enter/Space toggle; optional arrows/Home/End | Header is a clickable `<div>`; state not announced |
| Alert | — (live region) | None (no focus move) | Container injected with message; nothing announced |
| Alert dialog | `<dialog>` | Escape (when safe), Tab cycle | Ordinary confirmations misusing alertdialog |
| Breadcrumb | `<nav>` + `<ol>` of `<a>` | Standard link navigation | Missing `aria-current="page"` |
| Button | `<button>` | Enter, Space | Non-button element without keyboard activation |
| Carousel | — | Standard controls; no auto-advance on focus | Auto-rotation with no pause/stop |
| Checkbox | `<input type="checkbox">` | Space toggles | Custom visual with no checked-state exposure |
| Combobox | `<select>`/`<input>` where possible | Type to filter; Up/Down; Enter; Escape | Popup existence not exposed to screen reader |
| Dialog (modal) | `<dialog>` opened with `showModal()` (never the `open` attribute / `show()`) | Escape; Tab/Shift+Tab trapped | Focus escapes behind modal; background not inert |
| Disclosure | `<details>`/`<summary>` or `<button>` | Enter/Space toggle | Hidden content still focusable |
| Grid | `<table>` unless interactive | Arrow keys per cell | Grid role without keyboard model |
| Landmarks | `<header>`/`<nav>`/`<main>`/`<footer>` etc. | Screen-reader landmark navigation | Multiple unlabelled nav landmarks; no main |
| Listbox | `<select>` | Arrow keys; typeahead; Space/Enter select | Unconventional multi-select behavior |
| Menu/menubar | — (commands only) | Arrows; Enter/Space; Escape | Menu roles on site navigation links |
| Radio group | `<input type="radio">` + fieldset/legend | Arrow keys move/select | Missing group label |
| Slider | `<input type="range">` | Arrows increment/decrement | Value not exposed or not precise enough |
| Switch | `<input type="checkbox">` (styled) | Space toggles | Label describes action instead of setting |
| Table | `<table>` with `<th>`/caption | Standard reading navigation | Table used for layout; missing headers |
| Tabs | — | Tab in/out; arrows between tabs; Home/End | Auto-activation on slow-loading panels |
| Toolbar | — | Roving focus; arrow navigation | Every control in tab order with no roving model |
| Tooltip | — | Shows on focus; Escape dismisses | Hover-only; required info hidden in tooltip |
| Tree view | — | Arrows expand/collapse/move; Home/End; typeahead | Hierarchy state not exposed |
| Treegrid | — | Grid arrows + expand/collapse | Row/column/hierarchy context lost |
| Window splitter | — | Keyboard resize keys | No keyboard resize; position not exposed |

## Sources

- [WAI-ARIA Authoring Practices Guide (APG) — Patterns](https://www.w3.org/WAI/ARIA/apg/patterns/), W3C. Per-pattern keyboard interaction and role/state/property requirements.
- [WAI-ARIA 1.2 Specification](https://www.w3.org/TR/wai-aria-1.2/), W3C Recommendation. Normative role, state, and property definitions.
- [Using ARIA](https://www.w3.org/TR/using-aria/), W3C. The rules of ARIA use, including the first rule (prefer native HTML) and second rule (don't change native semantics).
- [HTML Standard — the `dialog` element](https://html.spec.whatwg.org/multipage/interactive-elements.html#the-dialog-element), WHATWG, and [MDN: `<dialog>`](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog). Native dialog semantics: `showModal()` vs `show()`/the `open` attribute, top-layer rendering, `::backdrop`, and Escape behavior.
- [Landmark Regions](https://www.w3.org/WAI/ARIA/apg/practices/landmark-regions/), W3C APG practices.
