# Accessibility & Inclusive Design

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers accessibility as a design discipline: the legal/business case, inclusive-design method, WCAG 2.2 structure and the most decision-relevant success criteria (with exact numbers, failure modes, and verification checklists), semantic HTML and ARIA, keyboard navigation, screen reader support, reduced motion, color independence, cognitive accessibility (including the COGA decision rules), touch accessibility, and testing methodology with step-by-step manual test scripts. Per-widget ARIA roles, keyboard maps, and the native-first element matrix live in [06-aria-widget-reference.md](06-aria-widget-reference.md). Contrast and color specifics also appear in [03-visual-design.md](03-visual-design.md#contrast); form accessibility in [07-forms-and-input.md](07-forms-and-input.md); keyboard interaction mechanics in [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns).

## Contents

- [Foundations](#foundations)
  - [Why accessibility](#why-accessibility)
  - [Inclusive design method](#inclusive-design-method)
  - [WCAG 2.2 structure and conformance levels](#wcag-22-structure-and-conformance-levels)
- [Key success criteria](#key-success-criteria)
  - [Contrast minimum](#contrast-minimum)
  - [Target size](#target-size)
  - [Focus visible](#focus-visible)
  - [Focus not obscured (enhanced)](#focus-not-obscured-enhanced)
  - [Keyboard access and no keyboard trap](#keyboard-access-and-no-keyboard-trap)
  - [Reflow](#reflow)
  - [Text spacing](#text-spacing)
  - [Content on hover or focus](#content-on-hover-or-focus)
  - [Dragging alternatives](#dragging-alternatives)
  - [Redundant entry](#redundant-entry)
  - [Accessible authentication](#accessible-authentication)
  - [Accessible authentication (enhanced)](#accessible-authentication-enhanced)
  - [Consistent help](#consistent-help)
  - [Motion and flashing limits](#motion-and-flashing-limits)
  - [Use of color](#use-of-color)
- [Implementation disciplines](#implementation-disciplines)
  - [Semantic HTML first](#semantic-html-first)
  - [ARIA](#aria)
  - [Keyboard navigation](#keyboard-navigation)
  - [Screen readers](#screen-readers)
  - [Forms accessibility](#forms-accessibility)
  - [Reduced motion](#reduced-motion)
  - [Color independence](#color-independence)
  - [Cognitive accessibility](#cognitive-accessibility)
  - [Touch/pointer accessibility](#touchpointer-accessibility)
  - [Zoom and text resize](#zoom-and-text-resize)
- [Verification](#verification)
  - [Testing methodology](#testing-methodology)
  - [Manual test scripts](#manual-test-scripts)
- [Quick reference](#quick-reference)

---

## Foundations

### Why accessibility

**Definition:** Designing so people with disabilities — visual, auditory, motor, speech, cognitive, neurological, vestibular, temporary and situational — can perceive, operate, and understand the product. Framed by law (ADA/Section 508 in the US; the European Accessibility Act, which entered into force in 2019 with a compliance deadline of June 28, 2025, using EN 301 549 as the technical standard; similar laws worldwide), by market (~16% of the world population, per WHO, has significant disability; prevalence rises steeply with age), and by quality (the curb-cut effect). Accessibility is not an enhancement layer — it is part of functional correctness.

**Reasoning / Evidence:** Legal exposure is concrete: thousands of ADA digital lawsuits annually in the US; the EAA extends mandatory compliance to private-sector e-commerce, banking, and consumer software sold in the EU (applicable from June 28, 2025). The curb-cut effect: accommodations built for disability benefit everyone (captions in noisy rooms, keyboard shortcuts for power users, high contrast in sunlight, plain language for everyone). Accessibility work overlaps ~80% with plain good UX.

**When to use:**
- Always; scoped pragmatically: AA conformance as the floor for anything public-facing or sold to organizations (procurement requires it: VPAT/ACR documentation in the US, EN 301 549 conformance in the EU).
- From project start — retrofit costs are a large multiple of built-in costs.
- Include accessibility in design review, code review, QA, and user testing — not as a separate downstream phase.

**When NOT to use / exceptions:**
- No product-level exceptions; prioritization within a remediation backlog is legitimate (blockers on core flows first).

**Pros:**
- Legal safety, market expansion (disability community + aging users), SEO overlap (semantic structure), quality forcing-function (accessible components are better-engineered components).

**Cons:**
- Real ongoing cost (audits, testing, training); constrains some aesthetic choices (contrast, motion) — these are constraints, not losses.

**Implementation guidance:**
- Bake into definition-of-done: every component ships with keyboard support, accessible names, contrast-checked tokens.
- Assign ownership: accessibility fails as "everyone's job, no one's job."
- Publish an accessibility statement; maintain a VPAT/ACR if selling B2B/government.

**Related:** [#inclusive-design-method](#inclusive-design-method), [#wcag-22-structure-and-conformance-levels](#wcag-22-structure-and-conformance-levels), [10-design-systems.md](10-design-systems.md#accessibility-baked-into-components).

**Sources:** [WHO: Disability](https://www.who.int/news-room/fact-sheets/detail/disability-and-health), [European Accessibility Act](https://ec.europa.eu/social/main.jsp?catId=1202), [W3C WAI: Business Case](https://www.w3.org/WAI/business-case/), [EN 301 549](https://www.etsi.org/deliver/etsi_en/301500_301599/301549/03.02.01_60/en_301549v030201p.pdf).

### Inclusive design method

**Definition:** A design methodology, not a compliance checklist: recognize exclusion, learn from human diversity, and **solve for one, extend to many**. Disability is treated as a mismatch between a person and their environment rather than a personal attribute, and it comes in three durations — **permanent** (one arm), **temporary** (broken arm), and **situational** (holding a baby). Designing for the permanent case produces solutions the temporary and situational cases inherit for free.

**Reasoning / Evidence:** The persona-spectrum framing (Microsoft Inclusive Design) reframes accessibility economics: the "edge case" of one-armed users expands to everyone carrying groceries; captions built for deaf users serve every commuter without headphones. Constraints from disability are design inputs that sharpen products for the whole population — the methodological generalization of the curb-cut effect described in [#why-accessibility](#why-accessibility).

**When to use:**
- At research and definition stages: include disabled users in research panels; treat assistive-technology workflows as first-class user journeys.
- When prioritizing: a fix motivated by a permanent disability should be scored with its temporary and situational beneficiaries included.
- When designing personalization and preferences: build settings (text size, motion, contrast, language complexity) into the main product rather than a segregated "accessibility mode."

**When NOT to use / exceptions:**
- Inclusive method complements, not replaces, WCAG conformance — "we did inclusive research" is not a conformance claim.

**Pros:** Expands the perceived market for accessibility work; produces genuinely better mainstream UX; prevents ghettoized "special versions."
**Cons:** Requires recruiting disabled participants (effort, budget); benefits are diffuse and need deliberate measurement to defend.

**Implementation guidance:**
- Run the persona-spectrum exercise per feature: for each sense/limb/faculty the feature depends on, name the permanent, temporary, and situational user who lacks it — then design for the permanent case.
- Treat constraints as design inputs, not edge cases; log exclusions found in research as defects, not "nice to haves."
- Never segregate: preferences and accommodations live in the standard settings surface, and the accessible path is the default path wherever possible.

**Related:** [#why-accessibility](#why-accessibility), [#cognitive-accessibility](#cognitive-accessibility), [17-ux-process-research.md](17-ux-process-research.md#usability-testing).

**Sources:** [Microsoft Inclusive Design](https://inclusive.microsoft.design/), [W3C WAI: Business Case](https://www.w3.org/WAI/business-case/).

### WCAG 2.2 structure and conformance levels

**Definition:** The Web Content Accessibility Guidelines 2.2 (W3C Recommendation, October 2023): 4 principles — **P**erceivable, **O**perable, **U**nderstandable, **R**obust — containing 13 guidelines and 86 success criteria at three levels: **A** (minimum, absence = blocking barriers), **AA** (the standard legal/contractual target), **AAA** (enhanced; not intended as a full-site requirement).

**Reasoning / Evidence:** AA is the referenced level in essentially all legislation (ADA settlements, Section 508, EN 301 549, EAA). WCAG 2.2 added 9 criteria over 2.1 (focus appearance/obscuring, target size minimum, dragging, redundant entry, accessible auth, consistent help) and removed 4.1.1 Parsing. WCAG 3.0 remains a long-horizon draft (new scoring model) — build to 2.2 AA now. Conformance is claimed for complete processes, not isolated components — one inaccessible step fails the whole flow.

**POUR principles translated to concrete obligations:**

| WCAG Principle | Practical translation | Common UI obligations |
| --- | --- | --- |
| Perceivable | Users can perceive content through available senses and assistive tech | Text alternatives, captions, contrast, reflow, text resize, non-color cues |
| Operable | Users can operate controls with available input methods | Keyboard, focus, target size, no traps, gesture alternatives, enough time |
| Understandable | Users can understand content, navigation, forms, and errors | Plain language, labels, predictable behavior, error identification, help |
| Robust | UI exposes durable semantics to browsers and assistive tech | Name/role/value, status messages, valid semantic structure |

**When to use:**
- Target **2.2 AA** for all new work; treat select AAA criteria as design aspirations where cheap (7:1 contrast for body text, 44px targets, fully unobscured focus).
- Contract/procurement contexts: state conformance claims precisely ("WCAG 2.2 AA, evaluated [date], [method]").

**When NOT to use / exceptions:**
- Full AAA is explicitly not required by W3C for entire sites (some criteria conflict with content types).
- Native mobile apps: WCAG applies via EN 301 549 interpretation; platform guidelines (iOS/Android accessibility) add specifics.

**Pros:** Testable, versioned, legally recognized target.
**Cons:** Success-criteria pass/fail doesn't guarantee good *experience* (conformant-but-awful is possible); some criteria need human judgment.

**Implementation guidance:**
- Structure audits by principle: Perceivable (alt text, contrast, reflow), Operable (keyboard, targets, timing), Understandable (labels, errors, consistency), Robust (name/role/value for AT).
- Track per-criterion status in CI/audit tooling; new-code gates on automated-checkable criteria.
- Accessibility defects in authentication, forms, navigation, focus, target size, and error recovery are usually blockers — triage them as such.

**Related:** [#testing-methodology](#testing-methodology), all criteria entries below.

**Sources:** [WCAG 2.2](https://www.w3.org/TR/WCAG22/), [WCAG 2.2 Quick Reference](https://www.w3.org/WAI/WCAG22/quickref/), [What's new in 2.2](https://www.w3.org/WAI/standards-guidelines/wcag/new-in-22/).

---

## Key success criteria

### Contrast minimum

**Definition:** SC 1.4.3 (AA): text contrast ≥**4.5:1**; large text (≥24px, or ≥18.66px bold) ≥**3:1**. SC 1.4.11 (AA): UI components and meaningful graphics ≥**3:1** against adjacent colors. SC 1.4.6 (AAA): **7:1** / 4.5:1.

**Reasoning / Evidence:** Models legibility for 20/40 vision (typical at ~age 60) and moderate low vision. The single most-failed criterion on the web (WebAIM Million: low-contrast text on ~80% of top home pages).

**When to use:** Every text/background pair, input borders, focus indicators, icons carrying meaning, chart elements, both themes.

**When NOT to use / exceptions:** Disabled controls, pure decoration, logos are exempt — but keep disabled ≥3:1 where feasible anyway.

**Pros:** Objective, automatable. **Cons:** Ratio math imperfect for dark modes/thin fonts (APCA is the future-facing supplement).

**Implementation guidance:**
- Enforce at token level ([03-visual-design.md](03-visual-design.md#contrast)); grey floor on white: #767676. Placeholders, helper text, and "subtle" UI text are the chronic offenders — check them.
- Automate: axe/Lighthouse in CI + manual spot checks over images/gradients at worst points.
- Test normal, dark, and high-contrast modes — a passing light theme guarantees nothing about the others.

**Related:** [03-visual-design.md](03-visual-design.md#contrast), [#use-of-color](#use-of-color).

**Sources:** [Understanding 1.4.3](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html), [Understanding 1.4.11](https://www.w3.org/WAI/WCAG22/Understanding/non-text-contrast.html).

### Target size

**Definition:** SC 2.5.8 (AA, new in 2.2): pointer targets ≥**24×24 CSS px**, or spaced so a 24px circle centered on each target doesn't intersect another target. SC 2.5.5 (AAA): **44×44px**. Explicit exceptions: **inline links within a sentence or block of text**, browser-native (user-agent) controls, spacing-equivalent placement, and essential presentations.

**Reasoning / Evidence:** Motor impairments (tremor, limited dexterity) and situational factors (moving vehicle, gloves) make small targets error-prone; Fitts's Law quantifies the speed/error trade ([01-core-principles.md](01-core-principles.md#fittss-law)). Platform guidance is stricter than WCAG: Apple 44×44pt, Material 48×48dp.

**When to use:** All buttons, icon buttons, toolbar controls, checkboxes, radio buttons, close buttons (chronically tiny), pagination, table row actions, map pins.

**When NOT to use / exceptions:** Inline text links are an explicit exception (they are constrained by line height and text flow, not the 24px floor — though generous line-height and spacing still reduce accidental activation); dense data grids may rely on the spacing alternative.

**Pros:** Fewer mis-taps for everyone. **Cons:** Density pressure in toolbars/tables — solve with padding-extended hit areas, not visual bloat.

**Implementation guidance:**
- Visual size can stay small if the *hit area* meets minimums (padding or pseudo-element expansion).
- Target 44–48px on touch surfaces regardless of the 24px legal floor; ≥8px between adjacent targets.
- Separate destructive and safe actions spatially; avoid clusters of tiny adjacent icon-only controls.

**Failure modes:**
- Icon-only button rendered at 16–20px with no padding-extended hit area.
- Adjacent row actions (edit/delete) packed so a 24px circle on one intersects the other.
- Destructive action placed touching a frequent action, converting mis-taps into data loss.

**Verification checklist:**
- [ ] Every pointer target measures ≥24×24 CSS px hit area, or qualifies for the spacing/inline/user-agent/essential exception.
- [ ] Touch surfaces meet the platform 44pt/48dp guidance, not just the WCAG floor.
- [ ] Destructive and safe actions are separated or protected.

**Related:** [15-platform-mobile.md](15-platform-mobile.md#touch-targets), [01-core-principles.md](01-core-principles.md#fittss-law).

**Sources:** [Understanding 2.5.8](https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html).

### Focus visible

**Definition:** SC 2.4.7 (AA): keyboard focus indicator must be visible. SC 2.4.11 (AA, 2.2): the focused element must not be *entirely* hidden by author content (sticky headers, cookie banners, chat widgets, overlays, drawers). SC 2.4.13 Focus Appearance (AAA): the focus indicator must be at least as large as the area of a **2 CSS px thick perimeter** of the unfocused component, and the changed pixels must have ≥**3:1** contrast between the focused and unfocused states.

**Reasoning / Evidence:** Keyboard users navigate by focus; an invisible focus point equals an invisible cursor. `outline: none` without replacement — a historically common reset — is a total keyboard blocker. 2.4.11 exists because modern sticky chrome (headers, cookie banners, chat bubbles) routinely covers the very element the keyboard user just focused.

**When to use:** Every focusable element, both themes; design focus rings as a token-level system, not per-component afterthoughts. Apply 2.4.11 checks wherever sticky headers/footers, banners, chat widgets, drawers, in-page nav, or route transitions exist.

**When NOT to use / exceptions:** `:focus-visible` legitimately hides rings for mouse clicks while keeping them for keyboard — the correct modern compromise.

**Pros:** Non-negotiable keyboard usability. **Cons:** Aesthetic resistance to rings — counter with well-designed branded rings.

**Implementation guidance:**
- Standard: 2px ring, offset 2px, high-contrast color (≥3:1 vs background), consistent product-wide; `:focus-visible` for pointer/keyboard differentiation.
- Never remove native focus styling unless replacing it with an equally visible indicator; prefer an outline/ring that fully surrounds the component; ensure the indicator survives dark mode, high contrast, and theming.
- For 2.4.11: use `scroll-padding`/`scroll-margin` for sticky regions; make `scrollIntoView` targets account for fixed headers; avoid large fixed overlays covering focusable content; test at zoom and narrow widths.

**Failure modes:**
- Focus outline exists but is hidden under a fixed footer or cookie banner.
- Skip-link target lands behind sticky nav.
- Modal/drawer leaves background content focusable behind the overlay.
- Focus shown only by a subtle background shift; focus color low-contrast; ring clipped by `overflow: hidden`.

**Verification checklist:**
- [ ] Tabbing through every focusable element keeps the focused element at least partially visible (2.4.11), including at 200% zoom.
- [ ] A keyboard user can identify the current focus at a glance on every interactive element, in every theme.
- [ ] Focused elements scroll into unobscured view past sticky headers.

**Related:** [04-interaction-design.md](04-interaction-design.md#component-states), [#keyboard-navigation](#keyboard-navigation), [#focus-not-obscured-enhanced](#focus-not-obscured-enhanced).

**Sources:** [Understanding 2.4.7](https://www.w3.org/WAI/WCAG22/Understanding/focus-visible.html), [Understanding 2.4.11](https://www.w3.org/WAI/WCAG22/Understanding/focus-not-obscured-minimum.html), [Understanding 2.4.13](https://www.w3.org/WAI/WCAG22/Understanding/focus-appearance.html).

### Focus not obscured (enhanced)

**Definition:** SC 2.4.12 (AAA, new in 2.2): when a component receives keyboard focus, **no part** of it may be hidden by author-created content — the enhanced version of 2.4.11, which only forbids *total* obscuring.

**Reasoning / Evidence:** Partial obscuring still costs keyboard and low-vision users: a half-covered button forces guessing about its label and state. Full visibility is the actual quality bar; 2.4.11's "partially visible" floor is a legal compromise, not a design target.

**When to use:** As a quality target for high-accessibility products, government/public services, and complex keyboard-heavy applications (dashboards, editors, admin tools).

**When NOT to use / exceptions:** Content the user themselves moved or opened (e.g., a user-repositioned panel) doesn't fail; AAA level means it is aspirational, not a conformance blocker, for most products.

**Pros:** Removes the "what am I focused on?" guessing game entirely; cheap if sticky-region scroll padding is already correct for 2.4.11.
**Cons:** Hard to guarantee on small viewports with unavoidable sticky chrome.

**Implementation guidance:**
- Size `scroll-padding` to the full height of sticky headers/footers so focused elements scroll completely clear of them.
- Audit chat bubbles, cookie banners, and floating action buttons against the tab path at narrow widths and 200% zoom.

**Verification checklist:**
- [ ] Tabbing through every focusable element leaves the focused component fully visible — no part hidden by author-created content.

**Related:** [#focus-visible](#focus-visible), [#keyboard-navigation](#keyboard-navigation).

**Sources:** [Understanding 2.4.12](https://www.w3.org/WAI/WCAG22/Understanding/focus-not-obscured-enhanced.html).

### Keyboard access and no keyboard trap

**Definition:** SC 2.1.1 (A): all functionality operable via keyboard. SC 2.1.2 (A): focus can always move away (no traps). SC 2.1.4 (A): single-character shortcuts must be remappable/disableable/focus-scoped.

**Reasoning / Evidence:** Keyboard is the universal fallback modality: screen reader users, motor-impaired users (switch devices emulate keyboards), and power users all depend on it. A keyboard-inaccessible feature is inaccessible period.

**When to use:** Everything. Custom widgets (dropdowns, sliders, editors, maps, drag-drop) are where violations live — native HTML is keyboard-accessible by default.

**When NOT to use / exceptions:** Functions genuinely path-dependent (freehand drawing) are exempt from 2.1.1's letter — provide alternative-outcome paths anyway.

**Pros:** Foundation for all AT. **Cons:** Custom-widget key maps are real engineering ([APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/); per-widget key maps in [06-aria-widget-reference.md](06-aria-widget-reference.md)).

**Implementation guidance:**
- Test ritual: unplug the mouse; complete every core flow. Where you get stuck, users are stuck.
- Full mechanics: [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns).

**Related:** [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns), [#keyboard-navigation](#keyboard-navigation), [06-aria-widget-reference.md](06-aria-widget-reference.md).

**Sources:** [Understanding 2.1.1](https://www.w3.org/WAI/WCAG22/Understanding/keyboard.html), [Understanding 2.1.2](https://www.w3.org/WAI/WCAG22/Understanding/no-keyboard-trap.html).

### Reflow

**Definition:** SC 1.4.10 (AA): content must reflow to a **320 CSS px** wide viewport (equivalent to 400% zoom on 1280px) without two-dimensional scrolling, except where 2D layout is essential (tables, maps, diagrams, code).

**Reasoning / Evidence:** Low-vision users zoom instead of (or before) using screen readers; content forcing horizontal+vertical scrolling at zoom is unreadable in practice. Effectively mandates responsive design ([13-platform-web.md](13-platform-web.md#responsive-design)).

**When to use:** All web layouts; test every template at 320px width / 400% zoom.

**When NOT to use / exceptions:** Data tables, maps, canvases may scroll 2D — but their *controls and chrome* must still reflow.

**Pros:** Aligns exactly with responsive/mobile work — usually near-free if mobile-first. **Cons:** Legacy fixed-width layouts require real rework.

**Implementation guidance:**
- Fluid layouts, max-width media, no fixed-px widths on text containers; test at 320px in CI screenshots.
- Don't disable pinch-zoom on mobile (`user-scalable=no` is a violation of resize expectations).

**Related:** [13-platform-web.md](13-platform-web.md#responsive-design), [#zoom-and-text-resize](#zoom-and-text-resize).

**Sources:** [Understanding 1.4.10](https://www.w3.org/WAI/WCAG22/Understanding/reflow.html).

### Text spacing

**Definition:** SC 1.4.12 (AA): no content/functionality loss when users override to: line height 1.5× font size, paragraph spacing 2×, letter spacing 0.12×, word spacing 0.16×.

**Reasoning / Evidence:** Dyslexic and low-vision users apply custom stylesheets/extensions increasing spacing; brittle fixed-height containers clip or overlap under these overrides.

**When to use:** All text containers — avoid fixed heights on anything holding text; let content size its box.

**When NOT to use / exceptions:** Canvas-rendered text and PDFs are separate battles; native apps analogous via platform font scaling ([15-platform-mobile.md](15-platform-mobile.md)).

**Pros:** Cheap if layout is intrinsically sized. **Cons:** Reveals every fixed-height hack in the codebase.

**Implementation guidance:**
- `min-height` over `height` for text boxes; test with a spacing-override bookmarklet; allow buttons/labels to grow.

**Related:** [03-visual-design.md](03-visual-design.md#line-height-and-line-length).

**Sources:** [Understanding 1.4.12](https://www.w3.org/WAI/WCAG22/Understanding/text-spacing.html).

### Content on hover or focus

**Definition:** SC 1.4.13 (AA): hover/focus-triggered content (tooltips, previews, mega menus) must be **dismissible** (Esc without moving pointer), **hoverable** (pointer can move onto the popup), and **persistent** (stays until dismissed/invalid).

**Reasoning / Evidence:** Magnification users must move their zoomed viewport onto the popup to read it — vanishing-on-leave tooltips are unreadable; accidental hovers must be dismissible without losing place.

**When to use:** Every tooltip, hover card, hover menu, preview popup.

**Implementation guidance:**
- Esc closes; small hover-bridge/`safe-polygon` so moving to the popup doesn't close it; no auto-timeout ([04-interaction-design.md](04-interaction-design.md#hover-dependence)).

**Related:** [04-interaction-design.md](04-interaction-design.md#hover-dependence).

**Sources:** [Understanding 1.4.13](https://www.w3.org/WAI/WCAG22/Understanding/content-on-hover-or-focus.html).

### Dragging alternatives

**Definition:** SC 2.5.7 (AA, 2.2): any drag operation must have a single-pointer, non-dragging alternative (click-based move, buttons, menus), unless dragging is essential.

**Reasoning / Evidence:** Sustained precise press-move-release excludes tremor, limited-dexterity, head-pointer, and eye-gaze users.

**When to use:** Kanban boards, sortable/reorderable lists, moving cards between columns, resizing panes, sliders (must be clickable/keyable to a value), range selection, file drop zones (must have a browse button), map panning (buttons/keys).

**Implementation guidance:**
- Patterns: "Move to…" menu; up/down or move-to-position buttons; keyboard pickup-move-drop ([04-interaction-design.md](04-interaction-design.md#drag-and-drop)).
- Kanban: action menu to move a card to a column/status. File upload: browse/select button beside the drop zone. Range selection: numeric/text inputs as an equivalent path.

**Failure modes:**
- Reorderable list with drag handles only — no move actions in a menu or keyboard equivalent.
- Drop-zone-only file upload with no browse button.
- Pane resize possible only by dragging the divider.

**Verification checklist:**
- [ ] Every drag outcome can be achieved without holding the pointer down and dragging (single clicks/taps or keyboard).
- [ ] The alternative is discoverable, not a hidden shortcut.

**Related:** [04-interaction-design.md](04-interaction-design.md#drag-and-drop).

**Sources:** [Understanding 2.5.7](https://www.w3.org/WAI/WCAG22/Understanding/dragging-movements.html).

### Redundant entry

**Definition:** SC 3.3.7 (A, 2.2): information previously entered in the same process must be auto-populated or selectable — not demanded again (shipping = billing checkbox, remembered email across steps). Exceptions: re-entry that is essential, security-required (password confirmation), or where the earlier data is no longer valid.

**Reasoning / Evidence:** Memory and motor cost of re-entry hits cognitive-disability users hardest and annoys everyone; directly encodes [recognition over recall](01-core-principles.md#6-recognition-rather-than-recall).

**When to use:** Multi-step forms/checkouts, account setup, booking flows, wizard flows, government/public-service applications, re-auth sequences.

**Implementation guidance:**
- Same-as checkboxes; persist step data through back-navigation and validation errors; session-level prefill.
- Where security/validity genuinely demands re-entry, say why.

**Failure modes:**
- Billing address requested after shipping with no same-as option.
- Name/email requested again in a later step of the same flow.
- Validation error wipes fields the user already completed.

**Verification checklist:**
- [ ] Previously entered data is reused or selectable in every later step of the same process.
- [ ] Back-navigation and error round-trips preserve entered data.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#form-length-reduction), [02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking).

**Sources:** [Understanding 3.3.7](https://www.w3.org/WAI/WCAG22/Understanding/redundant-entry.html).

### Accessible authentication

**Definition:** SC 3.3.8 (AA, 2.2): authentication must not require a **cognitive function test** (memorizing/transcribing, calculating, solving puzzles) without alternatives. Password managers must work (no paste-blocking); copyable codes OK; object-recognition CAPTCHAs need alternatives.

**Reasoning / Evidence:** Memory/processing disabilities make password recall and puzzle CAPTCHAs outright barriers; the SC effectively mandates supporting password managers, autofill, and non-cognitive alternatives (passkeys, magic links, biometrics, SSO).

**When to use:** All auth flows: never block paste (including OTP fields), never disable autofill (`autocomplete="current-password"`, OTP autocomplete), offer passkey/email-link/SSO alternatives, avoid visual CAPTCHAs (or provide alternatives).

**Implementation guidance:**
- Acceptable mechanisms: magic link, passkey/biometric via platform mechanisms, SSO, password-manager-compatible login, OTP with paste/autocomplete allowed.
- [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux) for full auth UX; WebAuthn/passkeys are the strategic direction.

**Failure modes:**
- Paste disabled in password or OTP fields.
- Password manager blocked by custom input controls.
- CAPTCHA requiring identification of ambiguous objects with no alternative path.
- Security questions relying on memory of arbitrary facts.

**Verification checklist:**
- [ ] A user can authenticate without any unsupported memory/transcription/puzzle burden.
- [ ] Paste and autofill work in every credential field.
- [ ] At least one non-cognitive-test path (passkey, magic link, SSO) exists or the cognitive test has an accessible mechanism.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux), [#accessible-authentication-enhanced](#accessible-authentication-enhanced).

**Sources:** [Understanding 3.3.8](https://www.w3.org/WAI/WCAG22/Understanding/accessible-authentication-minimum.html).

### Accessible authentication (enhanced)

**Definition:** SC 3.3.9 (AAA, new in 2.2): the enhanced version of 3.3.8 — authentication must not rely on cognitive function tests *at all*, removing 3.3.8's exceptions for **object recognition** ("select all traffic lights") and **personal-content recognition** (identify a photo you uploaded).

**Reasoning / Evidence:** Object- and personal-content-recognition steps still exclude users with memory, perception, and processing disabilities; 3.3.8 tolerates them as exceptions, 3.3.9 does not. For products serving the public — government, health, finance — the AAA bar is the honest inclusion target.

**When to use:** High-accessibility products, public services, and anywhere passkeys/magic links/SSO can fully replace CAPTCHA- and recognition-based gates.

**When NOT to use / exceptions:** AAA — aspirational for most products; abuse-prevention teams may need risk-based fallbacks, but design those as non-cognitive (rate limiting, device signals) rather than harder puzzles.

**Pros:** Eliminates the residual CAPTCHA barrier entirely; aligns with the passwordless direction of the platform vendors.
**Cons:** Requires investment in alternative abuse-prevention; some third-party auth/CAPTCHA vendors won't comply.

**Implementation guidance:**
- Prefer passkeys, magic links, SSO, and device-based verification over any recognition challenge.
- If a CAPTCHA-like gate is unavoidable, choose invisible/risk-scoring mechanisms over object-recognition puzzles.

**Verification checklist:**
- [ ] The authentication path avoids cognitive tests entirely, or supplies a direct non-cognitive alternative for every step.

**Related:** [#accessible-authentication](#accessible-authentication), [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux).

**Sources:** [Understanding 3.3.9](https://www.w3.org/WAI/WCAG22/Understanding/accessible-authentication-enhanced.html).

### Consistent help

**Definition:** SC 3.2.6 (A, 2.2): help mechanisms — human contact details, human contact mechanisms (chat), self-help options (FAQs), automated contact mechanisms — present across multiple pages must appear in the same relative location/order on each.

**Reasoning / Evidence:** Predictable placement lets users who need help find it without page-by-page re-search; encodes [consistency](01-core-principles.md#4-consistency-and-standards) as a legal requirement.

**Implementation guidance:** Fixed help entry point (footer slot, persistent menu item, floating widget) identical across templates and layout variants — unless the user themselves moved it.

**Verification checklist:**
- [ ] Help mechanisms appear in consistent relative order/location across all equivalent pages and responsive variants.

**Sources:** [Understanding 3.2.6](https://www.w3.org/WAI/WCAG22/Understanding/consistent-help.html).

### Motion and flashing limits

**Definition:** SC 2.3.1 (A): nothing flashes more than **3 times/second** (or below general/red flash thresholds) — seizure prevention. SC 2.3.3 (AAA but treat seriously): motion triggered by interaction can be disabled unless essential — vestibular protection. SC 2.2.2 (A): auto-moving content >5s must be pausable.

**Reasoning / Evidence:** Photosensitive epilepsy: flashing can cause seizures (hard physical harm — the most serious a11y failure class). Vestibular disorders: parallax/zoom/large motion causes nausea, migraines, disorientation.

**When to use:** Video content and animations near the flash thresholds (test with PEAT tool); all large-scale motion behind `prefers-reduced-motion` ([#reduced-motion](#reduced-motion)); carousels/tickers pausable.

**Implementation guidance:** No strobe effects ever; autoplay video muted and pausable; hero parallax and scroll-zoom effects disabled under reduced-motion.

**Related:** [#reduced-motion](#reduced-motion), [04-interaction-design.md](04-interaction-design.md#motion-and-animation).

**Sources:** [Understanding 2.3.1](https://www.w3.org/WAI/WCAG22/Understanding/three-flashes-or-below-threshold.html), [Understanding 2.3.3](https://www.w3.org/WAI/WCAG22/Understanding/animation-from-interactions.html).

### Use of color

**Definition:** SC 1.4.1 (A): color must not be the *only* visual means of conveying information, indicating action, prompting response, or distinguishing elements.

**Reasoning / Evidence:** ~8% of men have color vision deficiency; grayscale rendering, sunlight, and cheap displays create the same condition situationally.

**When to use:** Required-field marking, validation states, chart series, diffs, status dots, link styling within body text (underline or other non-color cue).

**Implementation guidance:** Pair color with icon/text/pattern/position; grayscale-test every status UI ([03-visual-design.md](03-visual-design.md#color-blindness-and-color-independence)).

**Related:** [03-visual-design.md](03-visual-design.md#color-blindness-and-color-independence), [#color-independence](#color-independence).

**Sources:** [Understanding 1.4.1](https://www.w3.org/WAI/WCAG22/Understanding/use-of-color.html).

---

## Implementation disciplines

### Semantic HTML first

**Definition:** Using native HTML elements for their intended purpose — `button`, `a`, `nav`, `main`, `h1–h6`, `label`, `table`, `dialog`, `details` — before reaching for `div`+ARIA reconstruction.

**Reasoning / Evidence:** Native elements ship with keyboard behavior, focus management, states, and AT semantics for free and bug-free. The First Rule of ARIA: don't use ARIA if a native element does the job. A `<div onclick>` "button" needs role, tabindex, key handlers, and focus styles to badly imitate what `<button>` gives free.

**When to use:** Default for everything; headings in strict hierarchy (one `h1`, no skipped levels); landmarks (`header`/`nav`/`main`/`aside`/`footer`) for page structure; real `<table>` for data.

**When NOT to use / exceptions:** Genuinely novel widgets with no native equivalent (tree grids, comboboxes pre-`selectlist`) — then follow APG patterns exactly (see [06-aria-widget-reference.md](06-aria-widget-reference.md) for the native-first decision matrix and per-widget requirements).

**Pros:** Free correctness, forward compatibility, less code.
**Cons:** Native styling limits occasionally push teams to rebuild (usually solvable with CSS now — `appearance`, `accent-color`, `dialog`).

**Implementation guidance:**
- Links navigate (`href`), buttons act — never swap them ([13-platform-web.md](13-platform-web.md#browser-conventions)).
- Lint: no click handlers on non-interactive elements; heading-order checks; landmark presence.

**Related:** [#aria](#aria), [06-aria-widget-reference.md](06-aria-widget-reference.md), [13-platform-web.md](13-platform-web.md).

**Sources:** [MDN: HTML elements](https://developer.mozilla.org/en-US/docs/Web/HTML/Element), [W3C: Using ARIA](https://www.w3.org/TR/using-aria/).

### ARIA

> **Per-widget reference:** the roles, states, keyboard maps, and native-vs-custom decision matrix for every common widget (dialog, combobox, tabs, menu, tree, grid, etc.) live in [06-aria-widget-reference.md](06-aria-widget-reference.md). This section covers ARIA strategy and the rules of use.

**Definition:** Accessible Rich Internet Applications attributes supplying AT semantics where HTML falls short: roles (what it is), properties (`aria-label`, `aria-describedby`), states (`aria-expanded`, `aria-selected`, `aria-checked`), and live regions (`aria-live`) announcing dynamic changes.

**Reasoning / Evidence:** W3C's Using ARIA rules: (1) prefer native; (2) don't change native semantics; (3) all interactive ARIA controls must be keyboard-usable; (4) don't put `role="presentation"`/`aria-hidden` on focusable elements; (5) interactive elements need accessible names. WebAIM audits: pages *with* ARIA average **more** errors than those without — misapplied ARIA is worse than none (it makes false promises to AT).

**When to use:**
- Custom widgets per APG pattern (complete pattern or nothing: role + all states + full key map — see [06-aria-widget-reference.md](06-aria-widget-reference.md)).
- Live regions for async updates: `polite` for status ("Saved"), `assertive` only for errors demanding attention.
- Accessible names on icon-only controls (`aria-label`); relationships (`aria-describedby` to error text).

**When NOT to use / exceptions:**
- Never redundant roles (`role="button"` on `<button>`); never `aria-label` overriding visible text with different words (breaks voice control — SC 2.5.3 Label in Name: accessible name must contain the visible label).
- No live regions for high-frequency updates (announcement spam).

**Pros:** Makes rich widgets AT-usable at all.
**Cons:** Half-done ARIA is actively harmful; requires SR testing to verify reality.

**Implementation guidance:**
- Copy APG patterns wholesale; state changes must update ARIA state attributes in real time (`aria-expanded` toggling with the disclosure).
- SPA route changes: announce via live region or move focus to the new `h1`.
- Verify with an actual screen reader, not just the DOM inspector.

**Related:** [06-aria-widget-reference.md](06-aria-widget-reference.md), [#semantic-html-first](#semantic-html-first), [#screen-readers](#screen-readers), [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns).

**Sources:** [Using ARIA](https://www.w3.org/TR/using-aria/), [ARIA APG](https://www.w3.org/WAI/ARIA/apg/patterns/), [WebAIM Million](https://webaim.org/projects/million/).

### Keyboard navigation

**Definition:** The site/app-level keyboard experience: logical tab order, skip links, focus management on dynamic changes, and sensible landmark/heading structure for AT shortcut navigation.

**Reasoning / Evidence:** Screen reader users navigate primarily by headings/landmarks/links lists (WebAIM surveys: heading navigation is the #1 page-exploration strategy); keyboard users pay a per-Tab cost that bad structure multiplies.

**When to use:**
- Skip link ("Skip to main content") as first tab stop on chrome-heavy pages.
- Focus management: modals trap and return; deleted-item focus moves to a neighbor; SPA navigation moves focus to new content ([04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns)).
- Tab order following visual order — DOM order and CSS visual order must not diverge (beware `order`/`flex-direction: row-reverse` traps).

**When NOT to use / exceptions:** `tabindex` > 0 — practically never (reorders unpredictably); `tabindex="-1"` for programmatic focus targets only.

**Pros:** Serves AT + power users at once. **Cons:** SPA focus management is genuinely fiddly — budget for it.

**Implementation guidance:**
- Heading outline = content outline; one `h1`, no skips.
- After async actions, put focus where the user's task continues.
- Test: Tab through every template; count stops to reach primary action (>15 = fix structure or add skip).

**Related:** [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns), [#focus-visible](#focus-visible), [06-aria-widget-reference.md](06-aria-widget-reference.md).

**Sources:** [WebAIM: Keyboard Accessibility](https://webaim.org/techniques/keyboard/), [WebAIM SR Survey](https://webaim.org/projects/screenreadersurvey10/).

### Screen readers

**Definition:** Supporting AT that renders UI as speech/braille (NVDA and JAWS on Windows, VoiceOver on Apple platforms, TalkBack on Android): accessible names, useful alt text, structural navigation, and announced state.

**Reasoning / Evidence:** The accessible name computation (`aria-labelledby` > `aria-label` > native host-language label such as `<label>`/`alt` > `title` attribute) determines what every control *is* to an SR user; unnamed controls read as "button" — unusable. Alt-text research: describe *purpose/content in context*, not pixel appearance.

**When to use:**
- Alt text: informative images describe the information ("Bar chart: revenue up 12% Q3"); functional images describe the action ("Search"); decorative get `alt=""` (silence — not a filename readout).
- Every icon-only control: accessible name matching visible/tooltip text.
- State announcements: loading→loaded, errors, dynamic list changes (live regions).

**When NOT to use / exceptions:**
- Don't over-describe (paragraph alt text for a simple photo; "image of…" prefixes — SRs announce imageness already).
- `aria-hidden` on decorative icons beside text labels.

**Pros:** Core AT population served; forces naming discipline that helps testing/automation too.
**Cons:** SR behavior varies by SR×browser combo; verification needs real testing.

**Implementation guidance:**
- Test matrix: NVDA+Chrome/Firefox (free), VoiceOver+Safari (built-in), TalkBack for Android apps.
- Complex images: short alt + longer described-by text/adjacent caption.
- Tables: real `th` with `scope`; charts: data-table fallback or summary.

**Related:** [#aria](#aria), [#semantic-html-first](#semantic-html-first), [06-aria-widget-reference.md](06-aria-widget-reference.md), [03-visual-design.md](03-visual-design.md#imagery-and-illustration).

**Sources:** [WebAIM: Alternative Text](https://webaim.org/techniques/alttext/), [W3C: Images Tutorial](https://www.w3.org/WAI/tutorials/images/), [Accessible Name and Description Computation](https://www.w3.org/TR/accname-1.2/).

### Forms accessibility

**Definition:** The form-specific requirements: programmatic label association, grouped controls (`fieldset`/`legend`), programmatically-linked error messages (3.3.1 Error Identification, 3.3.3 Error Suggestion), and `autocomplete` attributes (1.3.5 Identify Input Purpose).

**Reasoning / Evidence:** Unlabeled inputs are the most common form failure (WebAIM Million: ~35–40% of home pages); SR users tabbing into an unlabeled field hear nothing useful. Autocomplete attributes enable both browser autofill and AT input-purpose personalization.

**When to use:** Every input: visible `<label for>`; radio/checkbox groups: `fieldset`+`legend`; errors: text linked via `aria-describedby` + `aria-invalid`; personal-data fields: full `autocomplete` taxonomy (`name`, `email`, `tel`, `street-address`…).

**When NOT to use / exceptions:** Placeholder is never a label substitute ([07-forms-and-input.md](07-forms-and-input.md#labels-vs-placeholders)).

**Implementation guidance:**
- Error summary at top (focus moves there on failed submit) with anchor links to fields — GOV.UK pattern.
- Required state: `required`/`aria-required` + visible marking; full form patterns in [07-forms-and-input.md](07-forms-and-input.md).

**Related:** [07-forms-and-input.md](07-forms-and-input.md), [#redundant-entry](#redundant-entry).

**Sources:** [WebAIM: Creating Accessible Forms](https://webaim.org/techniques/forms/), [Understanding 3.3.1](https://www.w3.org/WAI/WCAG22/Understanding/error-identification.html), [Understanding 1.3.5](https://www.w3.org/WAI/WCAG22/Understanding/identify-input-purpose.html).

### Reduced motion

**Definition:** Honoring the OS-level "reduce motion" preference (`prefers-reduced-motion: reduce` media query; equivalent APIs on iOS/Android) by replacing large/parallax/zoom motion with fades or nothing.

**Reasoning / Evidence:** Vestibular disorders (a significant minority, more post-concussion/with-age) experience real physical symptoms — nausea, vertigo, migraine — from scaling, parallax, and large translations. The preference is a direct medical accommodation.

**When to use:** All motion beyond trivial micro-transitions: page transitions, parallax, auto-playing movement, scroll-linked effects, celebration animations.

**When NOT to use / exceptions:** Essential motion (a video's content, a game) is exempt; tiny opacity/color transitions are generally safe to keep.

**Pros:** Cheap accommodation (one media query wrapping animation layers).
**Cons:** Requires motion system discipline — retrofit is grep-and-pray; build the media query into the animation tokens/utilities from day one.

**Implementation guidance:**
- CSS: wrap animations in `@media (prefers-reduced-motion: no-preference)`; JS animation libs: check the media query before animating.
- Replace, don't just remove: cross-fade in place of slide/zoom keeps orientation function.
- Never autoplay background video for reduced-motion users.

**Related:** [04-interaction-design.md](04-interaction-design.md#motion-and-animation), [#motion-and-flashing-limits](#motion-and-flashing-limits).

**Sources:** [web.dev: prefers-reduced-motion](https://web.dev/articles/prefers-reduced-motion), [Understanding 2.3.3](https://www.w3.org/WAI/WCAG22/Understanding/animation-from-interactions.html).

### Color independence

**Definition:** The applied discipline for SC 1.4.1: every color-coded meaning carries a redundant non-color channel (icon, text, pattern, position, weight).

**Reasoning / Evidence:** See [#use-of-color](#use-of-color) and [03-visual-design.md](03-visual-design.md#color-blindness-and-color-independence) — consolidated guidance lives there.

**Implementation guidance:**
- Validation: icon + message, not red border alone. Charts: direct labels/patterns/shapes per series. Status: icon shapes differ (✓ ⚠ ✕), not just hue. Links: underline in body text.
- QA: grayscale filter pass on every status-bearing screen.

**Related:** [03-visual-design.md](03-visual-design.md#color-blindness-and-color-independence), [#use-of-color](#use-of-color).

**Sources:** [Understanding 1.4.1](https://www.w3.org/WAI/WCAG22/Understanding/use-of-color.html).

### Cognitive accessibility

**Definition:** Design for memory, attention, language-processing, and executive-function variation: plain language, consistent patterns, low memory burden, forgiving errors, no time pressure. Guided by W3C's COGA task force ("Making Content Usable for People with Cognitive and Learning Disabilities").

**Reasoning / Evidence:** Cognitive disabilities are the largest disability category and the least served by checkbox-testable criteria; WCAG 2.2's new SCs (redundant entry, accessible auth, consistent help) begin codifying it. Overlap with general usability is near-total — this is [02-cognitive-foundations.md](02-cognitive-foundations.md) with higher stakes.

**COGA decision rules (10 cognitive needs → design rules):**

| Need | Rule | Apply to |
| --- | --- | --- |
| Clear purpose | Make page/screen purpose obvious at top | All screens |
| Familiar hierarchy | Use predictable layout and common patterns | Navigation, forms, public services |
| Clear operation | Identify controls and what they affect | Custom controls, dashboards |
| Findability | Surface most important tasks and search | Content-heavy products |
| Clear language | Use short, literal, unambiguous language | All content, errors, help |
| Avoid mistakes | Prevent errors and make undo/correction easy | Forms, high-risk flows |
| Focus | Limit interruptions and distractions | Notifications, ads, carousels |
| Memory support | Do not require recall across steps | Multi-step flows, auth |
| Help/support | Provide human help for critical services | Government, healthcare, finance |
| Personalization | Support user preferences and assistive extensions | Web apps, content-heavy apps |

**When to use:**
- Plain language everywhere ([09-content-ux-writing.md](09-content-ux-writing.md#plain-language)); one primary action per screen; consistent placement/naming; visible progress in flows.
- Timeouts: warn before expiry, allow extension (SC 2.2.1), preserve data across re-auth.
- Critical flows (government, health, finance): aim beyond minimums — step-by-step structure, help at point of need, no jargon.
- Provide reminders, summaries, examples, and recovery paths rather than relying on user memory and precision.

**When NOT to use / exceptions:** Domain-expert tools may assume domain knowledge — but not *interface* complexity tolerance.

**Pros:** Helps every user under stress, fatigue, or distraction (i.e., everyone, often).
**Cons:** Hard to verify by automation; needs testing with affected users and content-quality investment.

**Implementation guidance:**
- Reading level ≤ grade 8 for general audiences; front-load key info; chunk instructions ([02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking)).
- Session expiry: 20s+ warning with one-click extend; drafts survive.
- Support error recovery over error prevention lectures.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md), [09-content-ux-writing.md](09-content-ux-writing.md#plain-language), [#accessible-authentication](#accessible-authentication), [#inclusive-design-method](#inclusive-design-method).

**Sources:** [W3C: Making Content Usable (COGA)](https://www.w3.org/TR/coga-usable/).

### Touch/pointer accessibility

**Definition:** Pointer-modality requirements beyond target size: SC 2.5.2 Pointer Cancellation (actions complete on *up*-event so users can slide off to abort), SC 2.5.1 Pointer Gestures (multipoint/path gestures get single-pointer alternatives), SC 2.5.4 Motion Actuation (shake/tilt features get UI alternatives and can be disabled).

**Reasoning / Evidence:** Tremor and limited dexterity make down-event triggers and complex gestures error factories; device-motion features fail wheelchair-mounted and tremor contexts.

**When to use:** All touch surfaces; any pinch/rotate/multi-finger or shake-to-X feature.

**Implementation guidance:**
- Trigger on release (native buttons already do); provide zoom buttons beside pinch; visible alternatives for every gesture ([04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures)); motion features off-by-default-able.

**Related:** [#target-size](#target-size), [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures), [15-platform-mobile.md](15-platform-mobile.md#gestures).

**Sources:** [Understanding 2.5.2](https://www.w3.org/WAI/WCAG22/Understanding/pointer-cancellation.html), [Understanding 2.5.1](https://www.w3.org/WAI/WCAG22/Understanding/pointer-gestures.html).

### Zoom and text resize

**Definition:** SC 1.4.4 (AA): text resizable to **200%** without loss of content/function; SC 1.4.10 reflow covers 400% page zoom ([#reflow](#reflow)). Native apps: honor platform font-size settings (Dynamic Type on iOS, font scale on Android).

**Reasoning / Evidence:** Text-size adjustment is the most common low-vision accommodation, far preceding screen reader adoption; apps ignoring platform font scaling clip or fix text at unreadable sizes.

**When to use:** Everywhere; web via `rem`-based sizing ([13-platform-web.md](13-platform-web.md#web-typography-specifics)); mobile via scalable type APIs.

**Implementation guidance:**
- Never `user-scalable=no`; never px-lock body text against browser preferences; test layouts at 200% text-only zoom (a different stress than page zoom).
- Mobile: test at largest platform text sizes ([15-platform-mobile.md](15-platform-mobile.md)).

**Related:** [#reflow](#reflow), [03-visual-design.md](03-visual-design.md#type-scale).

**Sources:** [Understanding 1.4.4](https://www.w3.org/WAI/WCAG22/Understanding/resize-text.html).

---

## Verification

### Testing methodology

**Definition:** The layered accessibility QA stack: automated scanners → manual keyboard pass → screen reader pass → zoom/contrast/motion checks → testing with disabled users. No single layer suffices.

**Reasoning / Evidence:** Automated tools (axe-core, Lighthouse, WAVE) catch missing alt/labels/contrast but cannot judge alt-text *quality*, focus *order* sense, or SR *experience*. Coverage estimates depend on how you count: Deque's 2021 study found automated tools can detect ~**57%** of issues when weighted by frequency of occurrence; other studies counting issue *types* (e.g., the GDS audit of testing tools) put coverage at roughly **30–40%**. Under either measure, a large share of real barriers — including many of the most serious — requires human testing.

**When to use:**
- Automated: in CI on every PR (axe-core integrations) — cheap regression net.
- Manual keyboard pass: every new flow (mouse unplugged, complete the task).
- SR pass: every major release on the test matrix (NVDA+Chrome, VoiceOver+Safari minimum).
- User testing with disabled participants: annually or per major redesign — the only test of real experience ([17-ux-process-research.md](17-ux-process-research.md#usability-testing)).

**When NOT to use / exceptions:**
- Never claim conformance from automated green alone; never claim a custom widget is accessible based only on automated linting; never treat overlay widgets ("accessibility overlays") as remediation — they don't fix underlying issues and are widely criticized by the disability community.

**Pros:** Layered cost-effectiveness: automation catches regressions free; humans catch what matters.
**Cons:** Manual layers need skill and scheduling; SR testing has a learning curve (budget training).

**Implementation guidance:**
- CI: axe-core + contrast tokens check + heading/landmark lint.
- Manual checklist per release: keyboard-only task completion; 200% zoom; 320px reflow; grayscale pass; reduced-motion pass; high-contrast/dark-mode pass where supported; touch-only pass where applicable; SR smoke test of core flows. Step-by-step scripts below in [#manual-test-scripts](#manual-test-scripts).
- Track issues by WCAG SC in the bug tracker; severity = blocker on core flow > blocker on secondary > friction.
- Document in an accessibility statement + VPAT/ACR for B2B.

**Related:** [17-ux-process-research.md](17-ux-process-research.md#accessibility-auditing-in-process), [#wcag-22-structure-and-conformance-levels](#wcag-22-structure-and-conformance-levels), [#manual-test-scripts](#manual-test-scripts), [21-agent-checklists.md](21-agent-checklists.md).

**Sources:** [Deque: automated coverage (~57%)](https://www.deque.com/blog/automated-testing-study-identifies-57-percent-of-digital-accessibility-issues/), [GOV.UK accessibility testing](https://accessibility.blog.gov.uk/2017/02/24/what-we-found-when-we-tested-tools-on-the-worlds-least-accessible-webpage/), [W3C: Overlay concerns — see WAI resources](https://www.w3.org/WAI/test-evaluate/).

### Manual test scripts

**Definition:** Repeatable step-by-step scripts for the five core manual passes — keyboard, screen reader, zoom/reflow, motion, touch/pointer — runnable by any team member without specialist tooling.

**Keyboard pass:**
1. Load the page from the top.
2. Press Tab through all controls.
3. Confirm focus is visible, not obscured, and ordered logically.
4. Activate controls with Enter/Space as appropriate.
5. Close dialogs/popovers with Escape where expected.
6. Complete the main task without touching the pointer.

**Screen reader smoke test:**
1. Navigate by landmarks and headings.
2. Inspect each form label and error.
3. Activate the main controls.
4. Open and close dialogs.
5. Trigger async status and error states.
6. Confirm names, roles, values, and updates are meaningful.

**Zoom/reflow test:**
1. Set browser zoom to 200% (and spot-check 400% / 320px width).
2. Test at a narrow viewport.
3. Confirm content remains readable and operable with one-axis scrolling.
4. Confirm sticky UI does not obscure focus or content.

**Motion test:**
1. Enable the OS reduced-motion setting.
2. Trigger animations and transitions.
3. Confirm motion is reduced/replaced and meaning remains clear.

**Touch/pointer test:**
1. Use touch or coarse-pointer simulation.
2. Check target size and spacing.
3. Verify no content or function requires hover.
4. Verify drag alternatives exist and work.

**Implementation guidance:** Attach these scripts to the release checklist; record pass/fail per template, not per page. Also verify high-contrast/dark modes where supported, and errors being announced or discoverable.

**Related:** [#testing-methodology](#testing-methodology), [21-agent-checklists.md](21-agent-checklists.md).

**Sources:** [WebAIM: Keyboard Accessibility](https://webaim.org/techniques/keyboard/), [W3C WAI: Easy Checks](https://www.w3.org/WAI/test-evaluate/preliminary/).

---

## Quick reference

| Requirement | One-line guidance | Key numbers |
|---|---|---|
| Legal target | WCAG 2.2 AA everywhere public/sold | EAA compliance deadline June 28, 2025 (in force since 2019) |
| Contrast (1.4.3/1.4.11) | Check every pair, both themes | 4.5:1 text; 3:1 large/UI; 7:1 AAA |
| Target size (2.5.8) | Hit areas, not visuals, must clear; inline links exempt | 24px floor; 44pt/48dp platform |
| Focus visible (2.4.7/11/13) | Branded ring system, never removed | 2px ring, ≥3:1 vs unfocused state, offset 2px |
| Focus not obscured — enhanced (2.4.12, AAA) | Focused element fully visible past sticky chrome | 0 px hidden |
| Keyboard (2.1.1/2) | Unplug the mouse; finish every flow | 0 traps |
| Reflow (1.4.10) | One-axis scroll at mobile width | 320px / 400% zoom |
| Text spacing (1.4.12) | No fixed heights on text | 1.5/2/0.12/0.16 overrides survive |
| Hover content (1.4.13) | Dismissible, hoverable, persistent | Esc closes |
| Dragging (2.5.7) | Non-drag alternative always | — |
| Redundant entry (3.3.7) | Never ask twice in one process | — |
| Accessible auth (3.3.8) | Paste allowed, autofill works, passkeys | No cognitive tests |
| Accessible auth — enhanced (3.3.9, AAA) | No object/personal-content recognition either | 0 cognitive tests, no exceptions |
| Consistent help (3.2.6) | Same help, same place, every page | — |
| Flashing (2.3.1) | Never strobe | ≤3 flashes/second |
| Use of color (1.4.1) | Color + one more channel | Grayscale test |
| Semantic HTML | Native first; ARIA last | First Rule of ARIA |
| ARIA | Complete APG patterns or none — see [06-aria-widget-reference.md](06-aria-widget-reference.md) | Live regions: polite default |
| Screen readers | Names, alt purpose, structure | NVDA+Chrome, VO+Safari matrix |
| Forms | Label every input; link every error | autocomplete taxonomy |
| Reduced motion | Fade replaces move | One media query, tokenized |
| Cognitive | Plain language, no memory tax, no rush | Grade ≤8; 20s timeout warning |
| Zoom (1.4.4) | rem sizing; never block zoom | 200% text; 400% page |
| Testing | Automation ~57% by issue frequency (~30–40% of issue types); humans the rest | axe in CI + keyboard + SR pass |
