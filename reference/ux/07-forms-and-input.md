# Forms & Data Entry

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers form design end to end: layout, labels, field types, validation, errors, defaults, high-stakes and regulated submissions, payment/identity inputs, authentication, uploads, and length reduction. Error-message *writing* is in [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing); accessibility mechanics in [05-accessibility.md](05-accessibility.md#forms-accessibility); validation interaction timing overview in [04-interaction-design.md](04-interaction-design.md#validation-feedback-timing); wizard navigation in [08-navigation-ia.md](08-navigation-ia.md#sequential--wizard-navigation); mobile input specifics in [15-platform-mobile.md](15-platform-mobile.md#mobile-forms).

## Contents

- [Structure](#structure)
  - [Form layout](#form-layout)
  - [Label placement](#label-placement)
  - [Labels vs placeholders](#labels-vs-placeholders)
  - [Required vs optional marking](#required-vs-optional-marking)
  - [Multi-step vs single-page forms](#multi-step-vs-single-page-forms)
  - [Form length reduction](#form-length-reduction)
- [Input controls](#input-controls)
  - [Input type selection](#input-type-selection)
  - [Date input](#date-input)
  - [Input masking and formatting](#input-masking-and-formatting)
  - [Search inputs](#search-inputs)
  - [File upload](#file-upload)
- [Validation and errors](#validation-and-errors)
  - [Inline validation timing](#inline-validation-timing)
  - [Error message design](#error-message-design)
  - [Destructive action confirmation](#destructive-action-confirmation)
- [Data and state](#data-and-state)
  - [Defaults and prefill](#defaults-and-prefill)
  - [Autosave and draft recovery](#autosave-and-draft-recovery)
- [High-stakes forms](#high-stakes-forms)
  - [High-risk and regulated forms](#high-risk-and-regulated-forms)
  - [Payment and identity inputs](#payment-and-identity-inputs)
- [Authentication](#authentication)
  - [Password and auth UX](#password-and-auth-ux)
- [Quick reference](#quick-reference)

---

## Structure

### Form layout

**Definition:** The spatial organization of a form: single-column vertical flow as the default, logical grouping into labeled sections, and field widths that match expected content length.

**Reasoning / Evidence:** Eye-tracking and completion studies (CXL/Baymard corpus) show single-column forms complete faster and with fewer skipped fields: one visual path, no ambiguity about order (multi-column forms create Z-vs-column ambiguity — do I go right or down?). Field width is a signifier: a 4-character-wide ZIP field communicates the expected input.

**When to use:**
- Single column for virtually all forms (registration, checkout, settings, applications).
- Group related fields under section headings every 3–7 fields ([Miller-safe chunks](02-cognitive-foundations.md#working-memory-and-chunking)); genuinely paired short fields may share a row (first/last name, city/state/ZIP, card expiry/CVC) when the pairing matches the user's mental unit.
- Match widths to content: full-width for open text, sized for constrained inputs (date parts, codes, quantities).

**When NOT to use / exceptions:**
- Dense internal tools/admin panels used daily by trained staff can justify multi-column grids (screen economy for experts — [03-visual-design.md](03-visual-design.md#density-modes)).
- Comparison/matrix input (grading rubrics, per-row settings) is inherently tabular.

**Pros:** Fastest completion, lowest skip rate, trivially responsive, clear tab order.
**Cons:** Longer visual page (mitigate with grouping and multi-step, not columns).

**Implementation guidance:**
- Label→field gap 4–8px; between fields 16–24px; between sections 32–48px ([Proximity](01-core-principles.md#proximity)).
- One primary submit button, left-aligned with the fields or full-width (mobile); never center-stray.
- Tab order strictly top-to-bottom; paired-row fields left-to-right then down.
- Field width tokens: XS (4ch: ZIP, CVC), S (10ch: phone, date), M (25ch: name, email), L (full: address line, textarea).

**Verification checklist:** Before shipping any form, verify it works end-to-end:
- [ ] Keyboard only: every field reachable and operable in logical order; submit on Enter; no focus traps.
- [ ] Screen reader: every label, hint, and error announced; groups exposed (fieldset/legend or equivalent).
- [ ] Browser autofill: correct `autocomplete` tokens; autofill populates the right fields and nothing breaks.
- [ ] Password managers: save and fill both work; no paste- or autofill-blocking.
- [ ] Mobile keyboards: correct `type`/`inputmode` per field (numeric for codes, email keyboard for email, etc.).

**Related:** [#label-placement](#label-placement), [03-visual-design.md](03-visual-design.md#spacing-and-layout), [05-accessibility.md](05-accessibility.md#forms-accessibility).

**Sources:** [NN/g: Form Design White Space](https://www.nngroup.com/articles/form-design-white-space/), [Baymard: Form Field Usability](https://baymard.com/blog/mobile-form-usability).

### Label placement

**Definition:** Where the field label sits: **top-aligned** (above the input), **left-aligned** (beside, in a label column), **inline/floating** (inside the field, animating up on focus).

**Reasoning / Evidence:** Matteo Penzo's eye-tracking study (2006) and subsequent replications provide fixation/saccade evidence: top-aligned labels required the fewest fixations — a single fixation often captured both label and field, with minimal saccade travel down one vertical scan path — while left-aligned layouts demanded heavier saccade movement between the label column and the fields. Left-aligned nonetheless eases *scanning a completed form* (the label column reads vertically) — better for review-heavy/edit-rarely contexts. Top-aligned also works at any label length and any viewport. Floating labels compromise both: small text, animation complexity, no room for hints, poor localization tolerance.

**When to use:**
- **Top-aligned:** default for all user-facing forms, all mobile, all translated products (labels can grow without breaking layout).
- **Left-aligned:** dense enterprise/admin screens where users *review* more than *fill*, and horizontal space is abundant; settings pages with short stable labels.
- **Floating labels:** only when vertical space is severely constrained AND labels are short — and even then prefer top-aligned.

**When NOT to use / exceptions:**
- Never left-aligned on mobile (label column collapses); never floating for complex fields needing helper text.

**Pros (top):** Efficient scanning (fewest fixations), responsiveness, i18n tolerance, room for hint text below label.
**Cons (top):** Tallest layout; long forms get very long (solve with steps/sections).

**Implementation guidance:**
- Top labels: 14–16px, weight 500–600, 4–8px above input, sentence case, no colons needed.
- Left labels (when used): right-align the label text toward the field (proximity binding), fixed column ~30–35% width.
- Helper text below the label or below the input — pick one convention product-wide.

**Related:** [#labels-vs-placeholders](#labels-vs-placeholders), [09-content-ux-writing.md](09-content-ux-writing.md#placeholder-and-helper-text).

**Sources:** [Penzo 2006 (UXmatters)](https://www.uxmatters.com/mt/archives/2006/07/label-placement-in-forms.php), [NN/g: Form Design](https://www.nngroup.com/articles/form-design-white-space/).

### Labels vs placeholders

**Definition:** The rule: placeholders (grey in-field text) must never serve as the only label. Labels persist; placeholders vanish on input.

**Reasoning / Evidence:** NN/g documented the failure set: (1) disappearing context — mid-entry users forget what the field asks, must delete to re-read; (2) memory tax on multi-field forms; (3) the contrast dilemma — typical light-grey placeholder text fails WCAG 1.4.3, but placeholder text darkened to pass 4.5:1 looks pre-filled, so users skip fields that appear completed; (4) no label after entry means reviewing/fixing a completed form is guesswork; (5) many AT combinations don't announce placeholders reliably.

**When to use (placeholders legitimately):**
- Format *examples* supplementing a visible label ("e.g. name@company.com").
- Search boxes (the magnifier + context make the field self-evident — the one conventional exception).

**When NOT to use / exceptions:**
- Never as label replacement; never for critical instructions; never repeat the label verbatim in the placeholder (noise).

**Pros (of proper labels):** Persistent context, review-ability, AT reliability, contrast compliance.
**Cons:** None — this is settled; floating labels are the compromise for teams that "want the clean look," and they're still worse than top labels.

**Implementation guidance:**
- Every input: visible `<label for>`; hints as persistent helper text below label/field, not placeholder.
- If placeholders are used for examples: prefix "e.g.". Note that WCAG 1.4.3 applies to visible placeholder text — it needs 4.5:1 contrast like other text — yet placeholder text that meets 4.5:1 tends to look like a pre-filled value. This tension is one more reason not to rely on placeholders at all: put anything that matters in persistent helper text, and keep examples distinguishable from entered text (e.g. italic styling).

**Related:** [01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall), [05-accessibility.md](05-accessibility.md#forms-accessibility), [09-content-ux-writing.md](09-content-ux-writing.md#placeholder-and-helper-text).

**Sources:** [NN/g: Placeholders in Form Fields Are Harmful](https://www.nngroup.com/articles/form-design-placeholders/).

### Required vs optional marking

**Definition:** Signaling which fields must be completed. Modern rule: when most fields are required, mark the **optional** ones ("(optional)" suffix); the asterisk-on-required convention remains acceptable when required fields are the minority — pick per form, apply consistently.

**Reasoning / Evidence:** NN/g: users don't reliably know asterisk conventions; explicit "(optional)" text is unambiguous and reduces both skipped-required errors and unnecessarily-completed optional fields. Better still: minimize optional fields until the distinction is rare ([#form-length-reduction](#form-length-reduction)).

**When to use:**
- Mostly-required forms (checkout, registration): mark optional fields "(optional)"; no asterisks anywhere.
- Mixed/mostly-optional forms (profiles, settings): mark required with * + a legend, or better, restructure.

**When NOT to use / exceptions:**
- Don't mark anything in single-field or all-required-and-obvious contexts (login).
- Never color-only required indication ([05-accessibility.md](05-accessibility.md#use-of-color)).

**Pros ("optional" marking):** Zero-convention comprehension; nudges teams to justify every optional field.
**Cons:** Slightly longer labels; mixed legacy conventions in one product confuse — migrate wholesale.

**Implementation guidance:**
- "(optional)" in lighter weight after the label text; `required`+`aria-required` on required inputs regardless of visual convention.
- If asterisks: red *, explained once above the form ("* required"), and still programmatically marked.

**Related:** [#form-length-reduction](#form-length-reduction), [05-accessibility.md](05-accessibility.md#forms-accessibility).

**Sources:** [NN/g: Required Fields](https://www.nngroup.com/articles/required-fields/), [GOV.UK: question pages](https://design-system.service.gov.uk/patterns/question-pages/).

### Multi-step vs single-page forms

**Definition:** Splitting a long form into sequential steps (wizard) versus one scrolling page.

**Reasoning / Evidence:** Steps reduce perceived effort (each screen looks achievable), enable per-step validation, and create progress momentum ([Goal-Gradient](01-core-principles.md#goal-gradient-effect)); single pages let users preview total scope and move freely. Checkout research (Baymard) favors clearly-stepped flows with visible progress; GOV.UK's "one thing per page" pattern tests exceptionally well for infrequent, high-stakes citizen tasks.

**When to use (multi-step):**
- >7–10 fields spanning distinct topics (account → shipping → payment).
- High-stakes/infrequent tasks where focus beats overview (government applications: one-thing-per-page).
- When later steps depend on earlier answers (branching).

**When to use (single page):**
- ≤7 fields; expert-repeated tasks (daily data entry — steps are a click tax); forms users need to see whole before starting.

**When NOT to use / exceptions:**
- Never fake steps (splitting 4 fields into 3 screens inflates interaction cost) — split only when it reduces complexity or maps to a real process, not to manufacture artificial progress.
- Never trap: steps must allow backward movement with data preserved (WCAG 3.3.7 — [05-accessibility.md](05-accessibility.md#redundant-entry)).

**Pros (steps):** Lower perceived load, cleaner validation, progress motivation, analytics per-step drop-off visibility.
**Cons (steps):** More clicks, harder whole-form review (mitigate: summary/review step), state management complexity.

**Implementation guidance:**
- Progress indicator: numbered steps with labels ("2 of 4 — Shipping"); completed steps clickable to revisit.
- Validate per step (errors never travel forward); persist all state through back/forward.
- Final review step for consequential submissions (order review before purchase; mandatory for regulated submissions — [#high-risk-and-regulated-forms](#high-risk-and-regulated-forms)).
- Save-and-resume for flows >5 minutes ([#autosave-and-draft-recovery](#autosave-and-draft-recovery)).

**Related:** [08-navigation-ia.md](08-navigation-ia.md#sequential--wizard-navigation), [01-core-principles.md](01-core-principles.md#goal-gradient-effect), [#high-risk-and-regulated-forms](#high-risk-and-regulated-forms).

**Sources:** [GOV.UK: One thing per page](https://design-system.service.gov.uk/patterns/question-pages/), [Baymard: Checkout Usability](https://baymard.com/research/checkout-usability).

### Form length reduction

**Definition:** The discipline of asking only what is needed, when it is needed: every field must justify its existence against measured completion cost.

**Reasoning / Evidence:** Field count correlates inversely with conversion across virtually all published A/B literature (each field is friction + privacy cost + error surface). Expedia's famous single-field removal ("Company" confusing users) reportedly recovered $12M/yr. Progressive profiling (deferring nice-to-have data to later sessions) preserves acquisition while eventually gathering data.

**When to use:**
- Every form review: for each field ask — do we *act* on this data now? Can we detect it (geo, locale), default it, derive it (card type from number), or defer it (profile completion later)?
- Name fields: single "Full name" unless legal/payment processing genuinely requires split fields.
- Address: autocomplete API (type-ahead) over 5 manual fields; country-adaptive formats.
- Phone/secondary email: only if there's a real operational use — and say what it is ("for delivery updates only").

**When NOT to use / exceptions:**
- Legal/compliance minimums (KYC, age verification) are floors — explain *why* you ask.
- Qualification friction is sometimes intentional (enterprise lead forms filter seriousness) — a business choice, made knowingly.

**Pros:** Conversion, privacy posture, less validation surface, faster completion.
**Cons:** Organizational fights (every team wants "just one field"); deferred collection needs later-capture flows built.

**Implementation guidance:**
- Field budget reviews with data: completion analytics per field (hesitation time, error rate, abandonment point).
- Explain sensitive asks inline ("Why do we need this?" disclosure).
- Never ask twice ([05-accessibility.md](05-accessibility.md#redundant-entry)): same-as-shipping defaults checked; derive everything derivable ([Tesler](01-core-principles.md#teslers-law)).

**Related:** [#defaults-and-prefill](#defaults-and-prefill), [01-core-principles.md](01-core-principles.md#teslers-law), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#forced-action-privacy-zuckering).

**Sources:** [Baymard: checkout field reduction](https://baymard.com/blog/checkout-flow-average-form-fields), [NN/g: website forms](https://www.nngroup.com/articles/web-form-design/).

---

## Input controls

### Input type selection

**Definition:** Choosing the control per data type and option count: text inputs (open data), radio groups (≤5–7 visible exclusive options), checkboxes (independent toggles), select dropdowns (long closed lists — a last resort), comboboxes/autocomplete (very long lists), segmented controls (2–5 exclusive, instant effect), steppers (small numeric ranges), toggles/switches (immediate binary state), sliders (continuous exploratory values).

**Reasoning / Evidence:** Luke Wroblewski / NN/g / Baymard convergence: dropdowns are overused — they hide options (extra interaction to see choices), are slow (open→scan→scroll→select), and painful on mobile. Radios expose all options at once (recognition, single tap). Toggle vs checkbox semantics: a switch implies *immediate effect* (like a light switch); a checkbox implies *submitted-with-form* state — mixing them confuses commit expectations.

**When to use:**
- 2–5 exclusive options: radio group (or segmented control for compact instant-effect settings).
- 5–15 options: radio if space allows; select if space-constrained.
- >15 options (countries, states, products): combobox with type-ahead filtering — never a bare 200-item select.
- Numeric small-range (1–10 quantity): stepper; large-range known value: plain number input; exploratory range: slider + numeric input pair.
- Binary immediate setting: switch. Binary form consent/selection: checkbox.

**When NOT to use / exceptions:**
- Native selects still win for platform-familiar OS pickers on mobile when the list is modest and no filtering is needed — don't rebuild what the OS does well.
- `<input type="number">` has UX bugs (scroll-to-increment accidents, letter "e" acceptance) — for codes/IDs use `inputmode="numeric" pattern="[0-9]*"` text inputs instead.
- Sliders for precise required values (exact price, exact date) are wrong — precision needs typing.

**Pros (right control):** Fewer taps, visible options, correct mobile keyboards, lower error rates.
**Cons:** More control types = more design-system surface; segmented/custom controls need full keyboard+ARIA work ([05-accessibility.md](05-accessibility.md#aria)).

**Implementation guidance:**
- Correct HTML types for keyboards & autofill: `email`, `tel`, `url`, plus `inputmode` and `autocomplete` attributes throughout ([15-platform-mobile.md](15-platform-mobile.md#mobile-forms)).
- Radio/checkbox: whole row (label included) clickable; ≥44px row height touch.
- Comboboxes: filter on contains-match not just prefix; show "no matches" with escape hatch.
- Switches: apply instantly, no Save button on the same settings (or use checkboxes + Save — one model per screen).

**Related:** [#date-input](#date-input), [05-accessibility.md](05-accessibility.md#forms-accessibility), [04-interaction-design.md](04-interaction-design.md#component-states).

**Sources:** [NN/g: Drop-Down Usability](https://www.nngroup.com/articles/drop-down-menus/), [NN/g: Toggle Switch Guidelines](https://www.nngroup.com/articles/toggle-switch-guidelines/), [GOV.UK: Select](https://design-system.service.gov.uk/components/select/).

### Date input

**Definition:** Choosing between typed date fields (separate D/M/Y inputs or single masked field) and calendar pickers.

**Reasoning / Evidence:** GOV.UK research: **memorable dates** (birth dates, document dates) are *typed* faster and more accurately — users know them as numbers; forcing calendar navigation to 1987 is torture. **Near-future scheduling dates** (appointments, bookings) benefit from calendars, which show day-of-week and availability context. Format ambiguity (02/03: Feb 3 or Mar 2?) is locale-dependent and must be resolved explicitly.

**When to use:**
- Memorable/past dates: three labeled fields (Day / Month / Year) — GOV.UK pattern — or one field with explicit format hint and liberal parsing ([Postel](01-core-principles.md#postels-law)).
- Scheduling/near dates: calendar picker WITH typed entry alternative (for power users and assistive-technology users).
- Ranges (check-in/out): dual-calendar with visible range highlighting + typed fields.

**When NOT to use / exceptions:**
- Never calendar-only for birth dates; never dropdown-triplets (31+12+100 options = worst of everything).
- Native mobile date inputs (`<input type="date">`) are decent on mobile, inconsistent on desktop — test per audience.

**Pros (typed):** Speed for known dates, AT-friendly, no widget complexity.
**Cons (typed):** Format ambiguity risk (solve with labeled parts or explicit hint + echo-back).

**Implementation guidance:**
- Labeled parts: Day (2ch) Month (2ch) Year (4ch), auto-advance OFF (backspace hostility), numeric inputmode.
- Single field: show format ("DD/MM/YYYY"), parse liberally, echo interpretation ("3 February 2025") for confirmation.
- Calendars: keyboard-navigable (arrows/PageUp/Down per APG), disabled invalid dates with reason, week starts per locale.

**Related:** [#input-type-selection](#input-type-selection), [09-content-ux-writing.md](09-content-ux-writing.md#numbers-dates-and-units-formatting).

**Sources:** [GOV.UK: Dates pattern](https://design-system.service.gov.uk/patterns/dates/), [NN/g: Date Input](https://www.nngroup.com/articles/date-input/).

### Input masking and formatting

**Definition:** Live formatting of input as typed (card number spacing 4-4-4-4, phone grouping, currency separators) and constrained-entry masks.

**Reasoning / Evidence:** Chunked display aids verification ([02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking)); auto-formatting removes the "will it accept spaces?" anxiety by accepting anything and normalizing ([Postel](01-core-principles.md#postels-law)). Danger zone: masks that fight the user (cursor jumps, deleted delimiters resurrecting, rejected paste).

**When to use:**
- Card numbers (auto-space, auto-detect brand icon), phone (locale-formatted), currency (thousands separators), sort codes/IBANs.

**When NOT to use / exceptions:**
- Auto-advance between fields on "completion" — breaks backspace correction and paste; avoid except proven kiosk flows.
- Masks on variable-format international data (phone formats vary wildly — format only after country is known).
- Never block paste ([05-accessibility.md](05-accessibility.md#accessible-authentication)) — parse pasted content liberally instead (strip spaces/dashes from pasted card numbers).

**Pros:** Verification ease, error reduction, perceived polish.
**Cons:** Implementation bugs are common and infuriating (cursor position math); test paste/backspace/mid-edit thoroughly.

**Implementation guidance:**
- Accept-then-format: never reject characters the formatter can strip; format on the display layer, store normalized.
- Preserve caret position through reformatting; format-on-blur is the safe fallback if live masking is buggy.

**Related:** [01-core-principles.md](01-core-principles.md#postels-law), [#inline-validation-timing](#inline-validation-timing).

**Sources:** [Baymard: Credit Card Form Fields](https://baymard.com/blog/credit-card-form-design).

### Search inputs

**Definition:** The search box as a form control: affordance (magnifier icon + placeholder context), submission (Enter + visible button), clearing (× when filled), and assistance (suggestions, recents).

**Reasoning / Evidence:** Search-dominant users exist on every content-rich product ([08-navigation-ia.md](08-navigation-ia.md#search-vs-browse)); the box itself has strong conventions (top placement, magnifier) whose violation costs findability.

**When to use:**
- Content/catalog products: persistent visible box (not icon-only reveal) sized ~27+ characters.
- With suggestions: recent searches (personal) + popular/autocomplete (global), keyboard-navigable.

**When NOT to use / exceptions:**
- Small apps (<~50 destinations) may not need search — a good nav beats a bad search.
- Icon-only collapsed search acceptable only under severe space constraints (mobile headers) — measured discoverability cost.

**Pros:** Direct goal expression; rescues failed navigation.
**Cons:** A bad search (no typo tolerance, bad ranking) is worse than none — users trust it and conclude content doesn't exist.

**Implementation guidance:**
- `type="search"` (built-in clearing on some platforms) + explicit × clear button; Enter submits; suggestions per APG combobox pattern.
- Preserve the query on the results page (editable); typo tolerance and synonyms in the engine ([08-navigation-ia.md](08-navigation-ia.md#no-results-recovery)).

**Related:** [08-navigation-ia.md](08-navigation-ia.md#search), [#input-type-selection](#input-type-selection).

**Sources:** [NN/g: Search — Visible and Simple](https://www.nngroup.com/articles/search-visible-and-simple/).

### File upload

**Definition:** Ingesting user files: click-to-browse + drag-drop zone, constraint communication (types/size *before* attempting), progress, and failure recovery.

**Reasoning / Evidence:** Upload failures after long waits are rage-inducing (sunk time + unclear cause); constraints revealed post-hoc ("PDF only" after selecting a DOCX) are pure preventable error ([01-core-principles.md](01-core-principles.md#5-error-prevention)).

**When to use:**
- Standard pattern everywhere: bordered drop zone with "Drag files here or **browse**" (browse = real button/label for `<input type=file>`), constraints listed inside the zone ("PNG, JPG, PDF · max 10 MB · up to 5 files").

**When NOT to use / exceptions:**
- Drag-drop is enhancement only — click path is the accessible baseline (2.5.7 — [05-accessibility.md](05-accessibility.md#dragging-alternatives)).
- Mobile: "browse" opens the OS picker/camera — design for that path primarily.

**Pros (dual-path zones):** Efficient for desktop power users, accessible baseline for all.
**Cons:** Multi-file state UI (per-file progress/errors/removal) is real work.

**Implementation guidance:**
- Validate type/size client-side at selection (instant rejection with reason), server-side always.
- Per-file: name, size, progress bar, cancel ×, retry-on-failure; overall status for batches.
- Large files: chunked/resumable (tus or equivalent) above ~50–100MB; warn on connection loss, resume on return.
- Show thumbnails for images; scanned/processing states distinct from uploading.

**Related:** [04-interaction-design.md](04-interaction-design.md#loading-strategies), [05-accessibility.md](05-accessibility.md#dragging-alternatives).

**Sources:** [NN/g: Drag-and-Drop](https://www.nngroup.com/articles/drag-drop/), [GOV.UK: File upload](https://design-system.service.gov.uk/components/file-upload/).

---

## Validation and errors

### Inline validation timing

**Definition:** The when of field validation: **on blur** for first judgment, **on keystroke after first error** for fix confirmation ("reward early, punish late"), **debounced async** for server checks, **on submit** as the always-present backstop.

**Reasoning / Evidence:** Keystroke-time flags punish incomplete input (typing "j" → "invalid email!") — measured frustration and *longer* completion (users stop to fight errors mid-thought). On-blur aligns judgment with the user's own "done with this field" signal. Once errored, immediate re-validation on typing gives instant fix feedback — the asymmetry that tests best (Baymard, GOV.UK guidance converge).

**When to use:**
- Per-field on blur (only for fields the user actually entered — tabbing through untouched fields must not flag).
- Errored fields: revalidate per keystroke; clear the error the moment input becomes valid.
- Async uniqueness/availability (username, email-exists): debounce 300–500ms after pause, in-field spinner, then verdict.
- Submit: validate all, summarize, focus the first error.

**When NOT to use / exceptions:**
- Prevent-over-validate where possible: length caps, numeric constraints at entry ([Constraints](01-core-principles.md#constraints)); show character limits (with live count) *before* they are exceeded, not after.
- Password-creation fields: live *requirement checklist* (checking off rules as met) beats blur-error — requirements are rules to satisfy, not errors to fix.
- Never validate-on-load or flag pristine empty fields.

**Pros:** Errors caught adjacent to cause; fixes confirmed instantly; flow preserved.
**Cons:** Multi-state logic (pristine/touched/errored/fixed) must be built correctly; premature blur flags on partial-but-intended-to-return entry (acceptable cost).

**Implementation guidance:**
- State machine per field: pristine → touched(blur) → valid|invalid → (invalid) live-revalidate.
- Success indication: subtle check for high-stakes formats (email, card); skip green-checking every trivial field (noise).
- Error rendering: red border + icon + message below, `aria-invalid` + `aria-describedby` ([05-accessibility.md](05-accessibility.md#forms-accessibility)); message specific to fix, not just "invalid" ([09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing)).

**Related:** [04-interaction-design.md](04-interaction-design.md#validation-feedback-timing), [#error-message-design](#error-message-design).

**Sources:** [NN/g: Errors in Forms](https://www.nngroup.com/articles/errors-forms-design-guidelines/), [GOV.UK: Validation](https://design-system.service.gov.uk/patterns/validation/).

### Error message design

**Definition:** The structural pattern for form errors: field-adjacent specific messages + (for long forms) a top error summary with anchor links, with all user input preserved.

**Reasoning / Evidence:** GOV.UK's tested pattern: on failed submit, focus moves to an error summary box ("There is a problem") listing each error as a link jumping to its field — solving discovery on long/scrolled forms and giving SR users a navigable manifest. Field messages must say *how to fix*, not just *that it's wrong*. Input preservation is absolute: clearing fields on error is among the most conversion-hostile behaviors measurable. A related failure mode: forcing the user to restart the whole form to correct one field — every error must be fixable in place.

**When to use:**
- Field-level message below every errored field: what's wrong + what to do ("Enter a date after 1 January 2025").
- Error summary at top for forms >~5 fields or multi-section; focus it on submit failure.
- Server-side failures: same rendering path as client errors (never a generic "something went wrong" banner alone for field-attributable problems).

**When NOT to use / exceptions:**
- Toast-only for field errors: never (transient + disconnected from the field).
- Don't stack summary + inline for 1–2-field forms; inline suffices.

**Pros:** Discovery, AT navigation, fix guidance, preserved trust.
**Cons:** Summary linking requires per-field anchors/IDs discipline.

**Implementation guidance:**
- Summary: `role="alert"` or focus-target heading, red-bordered box, error count, links (`href="#field-id"`).
- Message copy rules in [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing): specific, blame-free, jargon-free, no ALL CAPS.
- Preserve everything through error round-trips, including file selections where technically possible, scroll position, and step position.
- Analytics: log error type × field — the top of that table is your form redesign backlog ([17-ux-process-research.md](17-ux-process-research.md#analytics-and-behavioral-data)).

**Related:** [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing), [05-accessibility.md](05-accessibility.md#forms-accessibility), [01-core-principles.md](01-core-principles.md#9-help-users-recognize-diagnose-and-recover-from-errors).

**Sources:** [GOV.UK: Error summary](https://design-system.service.gov.uk/components/error-summary/), [NN/g: Error Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/).

### Destructive action confirmation

**Definition:** Protection for delete/discard/overwrite within form contexts: unsaved-changes guards, discard confirmations, and type-to-confirm for catastrophic scope. (General undo-vs-confirm doctrine: [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs).)

**Reasoning / Evidence:** Habituation makes routine confirmations worthless ([02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)); data loss (discarded drafts, overwritten records) is the highest-severity form failure — one incident teaches lasting distrust ([02-cognitive-foundations.md](02-cognitive-foundations.md#learned-helplessness-and-error-tolerant-design)).

**When to use:**
- Unsaved-changes navigation guard: only when real unsaved input exists; offer Save / Discard / Keep editing.
- Discard drafts: confirm with content stakes stated ("Discard 400 words?").
- Overwrite warnings when saving would replace newer/others' data (show whose/when).

**When NOT to use / exceptions:**
- Autosave eliminates most of this category — prefer it ([#autosave-and-draft-recovery](#autosave-and-draft-recovery)).
- Don't guard forms with zero entered data.

**Implementation guidance:**
- Dialog verbs name the action ("Discard draft" / "Keep editing" — never Yes/No); safe action focused by default.
- Browser `beforeunload` guard only while dirty; clear it on save.

**Related:** [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs), [09-content-ux-writing.md](09-content-ux-writing.md#confirmation-and-destructive-action-copy).

**Sources:** [NN/g: Confirmation Dialogs](https://www.nngroup.com/articles/confirmation-dialog/).

---

## Data and state

### Defaults and prefill

**Definition:** Pre-populating fields with detected, remembered, most-common, or safest values — the single most powerful form accelerator and the most ethically loaded ([default bias](02-cognitive-foundations.md#cognitive-biases-in-ux)).

**Reasoning / Evidence:** Defaults stick: effective organ-donation consent rates ran ~4.25%–27.5% in opt-in countries (Denmark as low as 4.25%) versus 85–99% in opt-out countries — from one checkbox's default (Johnson & Goldstein 2003). Applied honestly: detected country/locale/currency, remembered addresses, sensible quantity=1 save real effort ([Tesler](01-core-principles.md#teslers-law)). Applied dishonestly: prechecked newsletters, preselected add-on insurance = dark pattern, increasingly illegal (GDPR consent must be unchecked; EU consumer law bans pre-ticked extras).

**When to use:**
- Detectable context: country (geo/locale), language, currency, date format, phone prefix — always with easy correction.
- Remembered data: addresses, payment methods (tokenized), previous choices — the return-user experience.
- Most-common values: quantity 1, current date where sensible, standard shipping.
- Browser autofill: full `autocomplete` attribute coverage (this is also WCAG 1.3.5) — never disable autofill on address/payment/identity fields.

**When NOT to use / exceptions:**
- Never precheck: marketing consent, paid add-ons, data sharing, recurring upgrades ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#preselection)).
- No default where the wrong guess is costly and undetectable (title/salutation, gender) — empty beats wrong for identity fields.
- High-stakes irreversible choices (legal elections): deliberate unselected state forces conscious choice.

**Pros:** Massive effort reduction, error reduction (detected > typed), conversion.
**Cons:** Wrong defaults propagate silently (users don't check prefilled fields — verify detection quality); ethics/legal exposure when misused.

**Implementation guidance:**
- Visually distinguish nothing — prefilled fields look normal (users correct what they notice; excessive "we guessed this" chrome causes doubt); DO make correction one click/tap.
- `autocomplete` tokens per HTML spec (`given-name`, `email`, `street-address`, `cc-number`…); test real browser autofill flows (they break silently).
- Log correction rates on defaults: >20% corrections = bad default, remove or improve detection.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux), [01-core-principles.md](01-core-principles.md#teslers-law), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#preselection).

**Sources:** [Johnson & Goldstein 2003](https://doi.org/10.1126/science.1091721), [NN/g: The Power of Defaults](https://www.nngroup.com/articles/the-power-of-defaults/), [HTML autocomplete spec](https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#autofill).

### Autosave and draft recovery

**Definition:** Continuously persisting in-progress input (locally and/or server-side) so no crash, timeout, navigation, or accident loses work — with visible save-state indication and recovery on return.

**Reasoning / Evidence:** Data loss is the most severe form-UX failure ([learned helplessness](02-cognitive-foundations.md#learned-helplessness-and-error-tolerant-design)); post-completion errors make manual saving unreliable ([02-cognitive-foundations.md](02-cognitive-foundations.md#post-completion-errors)). Editor-class products (Docs, Notion, Figma) made continuous save the expectation; forms lag behind at users' expense.

**When to use:**
- Any input session >2–3 minutes: long forms, applications, content creation, settings with many changes.
- Session-expiry-prone contexts (banking, government): draft survival across re-auth is mandatory kindness (and 2.2.5/2.2.6-adjacent).

**When NOT to use / exceptions:**
- Sensitive data (card numbers, health details) may be inappropriate to persist locally — scope autosave per-field sensitivity.
- Single-field or sub-minute forms don't need it.
- Autosave of *destructive* edits without version history is dangerous — pair with undo/history in editors.

**Pros:** Eliminates the worst failure class; enables device-switching resume; removes save-anxiety.
**Cons:** Conflict handling (two tabs/devices editing), storage/privacy scope, "saved" state honesty (indicator must reflect real persistence).

**Implementation guidance:**
- Save triggers: debounced ~1–2s after typing pause + on blur + on step change; local (localStorage/IndexedDB) instantly, server as available.
- Indicator: quiet "Saved"/"Saving…" text near the form title (Docs pattern) — no toast per save.
- Recovery: on return, restore silently into the form OR offer explicitly ("Resume your application from March 3?") when staleness is possible.
- Conflicts: last-write-wins with visible warning for forms; real merge/history for documents.
- Session expiry: warn ≥20s before, one-click extend, and draft survives regardless ([05-accessibility.md](05-accessibility.md#cognitive-accessibility)).

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#post-completion-errors), [#destructive-action-confirmation](#destructive-action-confirmation), [14-platform-desktop.md](14-platform-desktop.md#offline-first-and-autosave).

**Sources:** [NN/g: autosave in complex applications](https://www.nngroup.com/articles/errors-forms-design-guidelines/).

---

## High-stakes forms

### High-risk and regulated forms

**Definition:** Forms whose submission carries real-world consequences — legal filings, financial transactions, medical declarations, government applications, identity verification. These require safeguards beyond ordinary form hygiene: explicit consequence disclosure before commit, a review step with edit paths, outcome-named commit buttons, confirmation/receipt, and timeout protection.

**Reasoning / Evidence:** Ordinary form errors cost a retry; high-stakes form errors cost money, legal standing, or eligibility — and are often irreversible after submission. WCAG 3.3.4 (Error Prevention: Legal, Financial, Data) makes reversibility, checking, or confirmation a normative requirement for exactly this class. GOV.UK's tested "Check your answers" pattern operationalizes it: a full review page with a change-link per answer before the commit action. The signature failure mode: the final button says "Continue" but actually commits legally — the user believes another step follows and discovers they have signed/filed/paid. Related failure modes: session expiry destroying an hour of application work; being unable to correct one field without restarting the entire form.

**When to use:**
- Any submission that creates a legal obligation, moves money, files with an authority, declares medical/identity facts, or cannot be freely undone.
- Multi-step regulated flows: always end with a review step ([#multi-step-vs-single-page-forms](#multi-step-vs-single-page-forms)).

**When NOT to use / exceptions:**
- Don't apply the full ceremony to low-stakes forms — a newsletter signup with a review step is friction theater; habituation to unnecessary confirmation erodes attention where it matters ([02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)).

**Pros:** Prevents the worst-severity error class; legal/compliance alignment (WCAG 3.3.4, consumer law); user trust in consequential moments.
**Cons:** Longer flows; review-step upkeep as fields change; teams tempted to water down consequence language (resist — see [09-content-ux-writing.md](09-content-ux-writing.md)).

**Implementation guidance:**
- **Consequence disclosure before commit:** state plainly what submitting does ("Submitting files your return with HMRC. You can't edit it afterwards.") adjacent to the commit button, not buried in terms.
- **Review step:** every answer displayed with a "Change" link returning to that field and then *back to review* (not restarting the flow); avoid redundant re-entry unless security/validity genuinely requires it ([05-accessibility.md](05-accessibility.md#redundant-entry)).
- **Outcome-named commit button:** the final button names the legal/financial outcome — "Sign and submit application", "Confirm and pay £142", "File return" — never "Continue", "Next", or "OK" on an action that commits.
- **Confirmation/receipt:** post-submit page + durable record (reference number, email/PDF receipt) stating what was submitted and when, and what happens next.
- **Timeout/data-loss protection:** autosave drafts through session expiry and re-auth ([#autosave-and-draft-recovery](#autosave-and-draft-recovery)).
- **Explain sensitive asks:** why each piece of legal/medical/identity data is needed, inline ([#form-length-reduction](#form-length-reduction)).
- Minimize collection: regulated context is not a license to over-collect.

**Decision questions:**
- Can the user state, before pressing the final button, exactly what will happen when they press it? (If not, the disclosure or button label is failing.)
- Can they fix any single answer from the review page without losing other work?
- What survives a session timeout at minute 40 of the application?

**Related:** [#multi-step-vs-single-page-forms](#multi-step-vs-single-page-forms), [#autosave-and-draft-recovery](#autosave-and-draft-recovery), [#destructive-action-confirmation](#destructive-action-confirmation), [09-content-ux-writing.md](09-content-ux-writing.md#confirmation-and-destructive-action-copy), [05-accessibility.md](05-accessibility.md#redundant-entry).

**Sources:** [Understanding WCAG 3.3.4: Error Prevention (Legal, Financial, Data)](https://www.w3.org/WAI/WCAG21/Understanding/error-prevention-legal-financial-data.html), [GOV.UK: Check answers pattern](https://design-system.service.gov.uk/patterns/check-answers/).

### Payment and identity inputs

**Definition:** The field cluster with the highest error cost and the strongest diversity of valid inputs: payment amounts and card data, names, and addresses. Core rules: show total cost before the payment commit; never encode restrictive assumptions about what a "valid" name looks like; never force one country's address format on the world.

**Reasoning / Evidence:** Surprise costs at payment are a leading documented cause of checkout abandonment (Baymard's checkout research; also the drip-pricing dark pattern — [18-patterns-antipatterns.md](18-patterns-antipatterns.md)). Name validation built on Western assumptions (required first+last, Latin-only characters, minimum lengths) rejects real people: single-name (mononym) users, names in non-Latin scripts, names with apostrophes/hyphens/spaces, very short or very long names — the W3C i18n personal-names guidance documents the breadth. Address formats vary structurally by country (postal-code formats and placement, administrative divisions, line order); a mandatory US-style State dropdown + 5-digit ZIP field simply blocks non-US customers.

**When to use:**
- Every checkout, donation, subscription, KYC, or account-identity flow.

**When NOT to use / exceptions:**
- Genuine legal/payment-processing requirements for split or exact-document names exist (KYC, card networks) — when they apply, say so ("as it appears on your card/ID") rather than silently rejecting; the constraint is the processor's, so name it.
- Domestic-only services may reasonably fix the country — but still validate against that country's real formats, and say the service is domestic-only up front.

**Pros:** Fewer abandoned payments, global reach, fewer support tickets from rejected legitimate identities.
**Cons:** Country-adaptive address forms are real engineering; legal-name edge cases need policy decisions, not just UI.

**Implementation guidance:**
- **Payment:**
  - Show the full total — item, shipping, taxes, fees — *before* the commit action; the commit button states the amount ("Pay $48.20").
  - Preserve non-sensitive entered data through payment errors (never make a declined card cost the whole form); card numbers/CVC may legitimately require re-entry.
  - Explain payment errors helpfully without leaking sensitive detail ("Your card was declined — try another card or contact your bank," not raw processor codes).
  - Card fields: auto-format 4-4-4-4, brand auto-detect, paste allowed ([#input-masking-and-formatting](#input-masking-and-formatting)); `autocomplete="cc-number"` etc.
- **Names:**
  - Default to a single "Full name" field ([#form-length-reduction](#form-length-reduction)); split only for genuine processing needs.
  - Accept single-word names, non-Latin scripts, diacritics, apostrophes, hyphens, and spaces; no arbitrary min/max lengths; never force capitalization.
  - Don't require "legal name" unless a legal process consumes it — and explain the verification requirement when it does.
- **Addresses:**
  - Country first; adapt subsequent fields (labels, order, which fields exist, postal-code validation) to the selected country — formats detail in [09-content-ux-writing.md](09-content-ux-writing.md).
  - Address-lookup/autocomplete with a manual-entry fallback always available (new buildings, rural addresses, lookup failures).
  - No US-only State/ZIP assumptions globally; make postal-code fields free-text-tolerant per country (UK/Canada/Netherlands formats include letters and spaces).

**Related:** [#high-risk-and-regulated-forms](#high-risk-and-regulated-forms), [#form-length-reduction](#form-length-reduction), [#input-masking-and-formatting](#input-masking-and-formatting), [09-content-ux-writing.md](09-content-ux-writing.md), [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md).

**Sources:** [Baymard: Checkout Usability](https://baymard.com/research/checkout-usability), [W3C i18n: Personal names around the world](https://www.w3.org/International/questions/qa-personal-names), [Baymard: Credit Card Form Fields](https://baymard.com/blog/credit-card-form-design).

---

## Authentication

### Password and auth UX

**Definition:** The authentication stack's UX: password fields done right (show toggle, paste allowed, length-first policy), progressive alternatives (magic links, OTP), platform biometrics, and the passkey (WebAuthn) transition.

**Reasoning / Evidence:** NIST SP 800-63B reversed decades of harmful folklore: **length over composition** (no mandatory symbol-soup, no periodic forced rotation, check against breached-password lists instead); composition rules produce predictable patterns (Password1!) while blocking passphrases. WCAG 3.3.8 mandates manager-compatibility (no paste-blocking, no autofill-blocking). Passkeys eliminate phishing and recall entirely — the strategic direction Apple/Google/Microsoft ship today.

**When to use:**
- Passwords (while they remain): minimum 8 chars (12+ encouraged) with no composition demands, live requirement checklist during creation, show-password toggle, `autocomplete="new-password"/"current-password"`, breached-list check with helpful rejection.
- Passkeys: offer as primary for new registrations where audience devices support; always with fallback.
- Magic links: low-frequency products where password memory is unreasonable; know the trade (email dependency, delay, device-switching friction).
- OTP/2FA: TOTP apps > SMS (SIM-swap risk); 6-digit codes chunked (123 456), `autocomplete="one-time-code"` for OS auto-fill, generous expiry (≥5 min), resend with cooldown timer shown.

**When NOT to use / exceptions:**
- Never: paste-blocking, autofill-blocking, aggressive short expiry on codes, security questions (guessable/researchable — NIST deprecated), forced periodic rotation without breach evidence.
- Don't demand re-auth mid-flow without preserving all state ([#autosave-and-draft-recovery](#autosave-and-draft-recovery)).
- CAPTCHA: last resort, with non-visual alternatives; invisible/risk-based first ([05-accessibility.md](05-accessibility.md#accessible-authentication)).

**Pros (modern stack):** Security *and* usability rise together (the old trade-off was mostly folklore); passkeys remove the phishing class entirely.
**Cons:** Passkey UX still maturing (cross-device/recovery flows confuse); magic links add latency; every added factor adds support burden.

**Implementation guidance:**
- Login form: single visible email→password (or email-first flow for SSO routing), labels not placeholders, submit on Enter, error preserving the email, "Forgot password" adjacent (not hidden).
- Show-password: eye toggle, keyboard-accessible, doesn't reset on submit failure.
- Errors: "Email or password didn't match" (don't confirm account existence on login; DO say clearly on registration if email exists — with sign-in link).
- Session UX: "remember me" honest scope; re-auth for sensitive ops only (step-up), not blanket short sessions.
- Password reset: expire tokens sanely (≥1h), single-use, land in a set-new-password form with the same creation aids.

**Related:** [05-accessibility.md](05-accessibility.md#accessible-authentication), [#inline-validation-timing](#inline-validation-timing), [15-platform-mobile.md](15-platform-mobile.md#mobile-forms).

**Sources:** [NIST SP 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html), [WebAuthn / passkeys](https://www.w3.org/TR/webauthn-2/), [Understanding WCAG 3.3.8](https://www.w3.org/WAI/WCAG22/Understanding/accessible-authentication-minimum.html).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Layout | Single column, grouped sections, width = content | Groups of 3–7 fields; gaps 16–24px |
| Label placement | Top-aligned default; left for dense review UIs | Label gap 4–8px |
| Placeholders | Never as labels; examples only | — |
| Required/optional | Mark the minority; prefer "(optional)" | — |
| Multi-step | Split >7–10 fields by topic; show progress | Review step for purchases |
| Length | Every field justified; detect/default/defer/drop | >20% default-correction = bad default |
| Control choice | Radios ≤7 visible; combobox >15; select last resort | Rows ≥44px touch |
| Toggle vs checkbox | Switch = instant; checkbox = submitted | One model per screen |
| Dates | Type memorable; pick near-future | D/M/Y labeled parts |
| Masking | Accept anything, format display, never block paste | Cards 4-4-4-4 |
| Validation timing | Blur first; keystroke after error; async debounced | 300–500ms debounce |
| Errors | Adjacent message + summary with links; preserve input | Summary for >5 fields |
| Defaults | Honest best-interest only; full autocomplete attrs | Never precheck consent/add-ons |
| Autosave | Any session >2–3 min; quiet indicator | Debounce 1–2s; warn 20s pre-expiry |
| High-risk forms | Review step + edit paths; commit button names the outcome, never "Continue" | Receipt + reference after submit |
| Payment & identity | Total cost before commit; international names/addresses; country-adaptive formats | Commit button states amount |
| Passwords | Length over composition; paste/autofill sacred | ≥8 min (12+ better); no rotation |
| 2FA/OTP | TOTP > SMS; chunked codes, OS autofill | 6 digits; ≥5 min expiry |
| Passkeys | Offer as primary where supported | — |
