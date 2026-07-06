# Core UI/UX Principles

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers the foundational, platform-agnostic principles of interface design: Nielsen's 10 usability heuristics, Don Norman's design principles, Gestalt perception principles, the "Laws of UX" (empirical rules such as Fitts's and Hick's laws), Shneiderman's Eight Golden Rules, and the ISO 9241-110 dialogue principles. The cognitive psychology underlying many of these is expanded in [02-cognitive-foundations.md](02-cognitive-foundations.md). Applied visual rules are in [03-visual-design.md](03-visual-design.md); interaction mechanics in [04-interaction-design.md](04-interaction-design.md); evaluation methods that use these principles (heuristic evaluation) in [17-ux-process-research.md](17-ux-process-research.md).

## Contents

- [Nielsen's 10 usability heuristics](#nielsens-10-usability-heuristics)
  - [1. Visibility of system status](#1-visibility-of-system-status)
  - [2. Match between system and the real world](#2-match-between-system-and-the-real-world)
  - [3. User control and freedom](#3-user-control-and-freedom)
  - [4. Consistency and standards](#4-consistency-and-standards)
  - [5. Error prevention](#5-error-prevention)
  - [6. Recognition rather than recall](#6-recognition-rather-than-recall)
  - [7. Flexibility and efficiency of use](#7-flexibility-and-efficiency-of-use)
  - [8. Aesthetic and minimalist design](#8-aesthetic-and-minimalist-design)
  - [9. Help users recognize, diagnose, and recover from errors](#9-help-users-recognize-diagnose-and-recover-from-errors)
  - [10. Help and documentation](#10-help-and-documentation)
- [Norman's design principles](#normans-design-principles)
  - [Affordances](#affordances)
  - [Signifiers](#signifiers)
  - [Mapping](#mapping)
  - [Feedback](#feedback)
  - [Constraints](#constraints)
  - [Conceptual models](#conceptual-models)
- [Gestalt principles](#gestalt-principles)
  - [Proximity](#proximity)
  - [Similarity](#similarity)
  - [Closure](#closure)
  - [Continuity](#continuity)
  - [Common region](#common-region)
  - [Figure/ground](#figureground)
  - [Common fate](#common-fate)
- [Laws of UX](#laws-of-ux)
  - [Fitts's Law](#fittss-law)
  - [Hick's Law](#hicks-law)
  - [Jakob's Law](#jakobs-law)
  - [Miller's Law](#millers-law)
  - [Tesler's Law](#teslers-law)
  - [Doherty Threshold](#doherty-threshold)
  - [Peak-End Rule](#peak-end-rule)
  - [Aesthetic-Usability Effect](#aesthetic-usability-effect)
  - [Von Restorff Effect](#von-restorff-effect)
  - [Zeigarnik Effect](#zeigarnik-effect)
  - [Serial Position Effect](#serial-position-effect)
  - [Postel's Law](#postels-law)
  - [Pareto Principle](#pareto-principle)
  - [Goal-Gradient Effect](#goal-gradient-effect)
  - [Parkinson's Law](#parkinsons-law)
- [Rule sets and standards](#rule-sets-and-standards)
  - [Shneiderman's Eight Golden Rules](#shneidermans-eight-golden-rules)
  - [ISO 9241-110 dialogue principles](#iso-9241-110-dialogue-principles)
- [Quick reference](#quick-reference)

---

## Nielsen's 10 usability heuristics

The 10 usability heuristics — the most widely used broad rules of thumb for interaction design — were originally co-authored by Jakob Nielsen and Rolf Molich in 1990; Nielsen alone refined the set in 1994 based on a factor analysis of 249 usability problems (definitions unchanged since 1994, examples updated 2024). They are deliberately general — usable for design guidance and for [heuristic evaluation](17-ux-process-research.md).

### 1. Visibility of system status

**Definition:** The design should always keep users informed about what is going on, through appropriate feedback within a reasonable amount of time.

**Reasoning / Evidence:** Nielsen/Molich heuristic #1. Knowing current state lets users learn the outcome of prior interactions and decide next steps; predictable feedback builds trust. Directly tied to response-time research (0.1s/1s/10s limits — see [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds)).

**When to use:**
- Any operation that is not instantaneous (network calls, file operations, background processing).
- Mode changes (edit mode, offline mode, caps lock, recording) — status must be visible while the mode is active.
- Long-running or multi-step flows (progress indicators, step counts).
- System-initiated changes (sync completed, session about to expire).

**When NOT to use / exceptions:**
- Do not surface internal states users cannot act on and do not care about (e.g., cache invalidation details) — this becomes noise.
- Sub-100ms operations need no explicit feedback beyond the natural UI change.

**Pros:**
- Builds trust; reduces repeated clicks, abandonment, and "is it broken?" support tickets.
- Enables error recovery — users notice when something went wrong.

**Cons:**
- Over-signaling creates noise and habituation (users stop reading toasts).
- Poorly-designed indicators (endless spinners with no ETA) can increase anxiety versus showing nothing for very short waits.

**Implementation guidance:**
- <100ms: no indicator needed. 100ms–1s: subtle cue (button pressed state, cursor). 1–10s: spinner or skeleton. >10s: determinate progress bar with time estimate and cancel option.
- Show system state persistently for modes (banner for offline, dot for unsaved changes), transiently for events (toast for saved).
- Confirm destructive/consequential completions explicitly ("Deleted 3 items — Undo").
- Respond to direct manipulation immediately, even if backend work continues — but distinguish optimistic UI from confirmed persistence when data integrity matters (show a pending state until the write is confirmed).
- Show determinate progress when actual progress is knowable; otherwise indeterminate progress with useful context (what is loading).
- Status must not depend on color, animation, or visual position alone; announce important asynchronous changes programmatically (status regions, `aria-live`/`role="status"`) for screen reader users ([05-accessibility.md](05-accessibility.md)).

**Platform notes:**
- Web: use semantic status messaging; avoid unexpectedly replacing content that holds keyboard focus ([13-platform-web.md](13-platform-web.md)).
- Desktop: use progress bars, status bars, task indicators, and notifications per OS conventions ([14-platform-desktop.md](14-platform-desktop.md)).
- Mobile: use haptics sparingly; avoid covering primary content with transient messages ([15-platform-mobile.md](15-platform-mobile.md)).

**Decision questions:**
- Does the user need to know the action was received, completed, failed, or queued?
- Is the process short, long, interruptible, or irreversible?
- Would missing feedback cause duplicate work or data loss?

**Verification checklist:**
- Trigger each main action and verify feedback arrives within the expected threshold.
- Test slow-network and failure states, not just the happy path.
- Test dynamic status updates with a screen reader.

**Related:** [04-interaction-design.md](04-interaction-design.md#loading-strategies), [#feedback](#feedback), [Doherty Threshold](#doherty-threshold).

**Sources:** [NN/g: 10 Usability Heuristics](https://www.nngroup.com/articles/ten-usability-heuristics/), [NN/g: Visibility of System Status](https://www.nngroup.com/articles/visibility-system-status/), [WCAG 2.2: Status Messages (4.1.3)](https://www.w3.org/WAI/WCAG22/Understanding/status-messages.html).

### 2. Match between system and the real world

**Definition:** The design should speak the users' language — words, phrases, and concepts familiar to the user rather than internal jargon — and follow real-world conventions, making information appear in a natural and logical order.

**Reasoning / Evidence:** Nielsen/Molich heuristic #2. When controls follow real-world conventions and correspond to desired outcomes (natural mapping), interfaces are easier to learn and remember. Terminology mismatches between team vocabulary and user vocabulary are among the most common findability failures.

**When to use:**
- All labeling, naming, and information ordering decisions.
- Domain-specific products: use the domain's established terminology (medical, legal, finance), not your engineering names.
- Icon and metaphor selection (trash can, folder, shopping cart).

**When NOT to use / exceptions:**
- Skeuomorphic metaphors can constrain functionality (the "desktop" metaphor limits at scale); do not force a real-world analogy where the digital operation has no physical equivalent.
- Expert tools may legitimately use dense professional jargon that would fail general audiences — match the actual user, not the general public.

**Pros:**
- Reduces learning cost to near zero when the mapping is good.
- Reduces support and training cost.

**Cons:**
- Real-world metaphors can mislead when the digital behavior diverges (e.g., "delete" that only archives).
- User vocabulary must be researched, not assumed — adds research cost.

**Implementation guidance:**
- Derive labels from user research artifacts (interview transcripts, search logs, card-sort labels), not internal schemas.
- Order fields/steps in the sequence users naturally think about the task (e.g., address forms in the order people say addresses in that locale).
- Test labels with cloze/comprehension checks or first-click tests; see [17-ux-process-research.md](17-ux-process-research.md).
- Provide definitions for unavoidable technical terms; plain language also supports cognitive accessibility and second-language readers.

**Decision questions:**
- Would a new user understand this label without documentation?
- Is this term used by users, or only by the team (does navigation mirror the org chart or backend architecture)?
- Does the workflow order match how users think about the task?

**Related:** [09-content-ux-writing.md](09-content-ux-writing.md), [Conceptual models](#conceptual-models), [02-cognitive-foundations.md](02-cognitive-foundations.md#mental-models).

**Sources:** [NN/g: Match Between System and the Real World](https://www.nngroup.com/articles/match-system-real-world/), [NN/g: 10 Usability Heuristics](https://www.nngroup.com/articles/ten-usability-heuristics/), [GOV.UK: Content design guidance](https://www.gov.uk/guidance/content-design).

### 3. User control and freedom

**Definition:** Users often perform actions by mistake and need a clearly marked "emergency exit" — undo, redo, cancel, back — to leave an unwanted state without an extended process.

**Reasoning / Evidence:** Nielsen/Molich heuristic #3. Exits keep users in control and prevent the fear of exploration; people learn interfaces by trying things, and they only try when mistakes are cheap.

**When to use:**
- Every multi-step flow needs a visible way out (cancel/back) without penalty.
- Destructive or bulk operations: prefer undo over confirmation (see [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- Modals and full-screen takeovers: Escape key, close button, click-outside (where safe).

**When NOT to use / exceptions:**
- Genuinely irreversible external effects (payment sent, email delivered) cannot be undone — use confirmation before, clear status after; consider grace-period holds (e.g., Gmail's 5–30s undo send).
- Do not offer undo when recovery cannot actually be guaranteed — an "undo" that only reverses the UI while backend effects persist is worse than an honest confirmation.
- Legal/compliance flows may forbid skipping steps; still allow saving progress and exiting.

**Pros:**
- Encourages exploration → faster learning.
- Reduces error cost and rage-quits.

**Cons:**
- Undo infrastructure is engineering-expensive (command stacks, soft deletes, tombstones).
- Too many escape routes in critical flows can cause accidental data loss if drafts aren't saved.

**Implementation guidance:**
- Support Ctrl/Cmd+Z globally in editing contexts; multi-level undo (≥20 steps typical for editors).
- Soft-delete with a trash/archive and a restore window (30 days is a common convention).
- Cancel must actually cancel (abort the request), not just hide the dialog.
- Preserve entered data when users navigate back ([07-forms-and-input.md](07-forms-and-input.md)).
- Label exits clearly ("Cancel", "Close", "Discard changes", "Keep editing"); keyboard and assistive-technology users need the same exits as pointer users (predictable Escape behavior for dialogs and transient layers).

**Decision questions:**
- What happens if the user made a mistake here?
- Can the action be reversed reliably, including backend and external side effects?
- Is confirmation, undo, or prevention the best protection for this action?

**Related:** [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs), [#5-error-prevention](#5-error-prevention).

**Sources:** [NN/g: User Control and Freedom](https://www.nngroup.com/articles/user-control-and-freedom/).

### 4. Consistency and standards

**Definition:** Users should not have to wonder whether different words, situations, or actions mean the same thing. Follow platform and industry conventions (external consistency) and keep the product self-consistent (internal consistency).

**Reasoning / Evidence:** Nielsen/Molich heuristic #4, grounded in [Jakob's Law](#jakobs-law): users spend most of their time in other products, which set their expectations. Inconsistency forces new learning and increases cognitive load.

**When to use:**
- Default position for every design decision: follow the platform convention unless you have evidence a deviation helps ([14-platform-desktop.md](14-platform-desktop.md), [15-platform-mobile.md](15-platform-mobile.md)).
- Within a product: same term, icon, placement, and behavior for the same concept everywhere (enforced via a [design system](10-design-systems.md)).

**When NOT to use / exceptions:**
- Deliberate inconsistency is correct when two things that look similar behave differently — making them look different prevents errors.
- Innovation is justified when the convention demonstrably fails the task and testing shows the new pattern wins; expect a transition cost.

**Pros:**
- Near-zero learning cost for conventional patterns; transfers knowledge across products.
- Reduces design/dev effort via reuse.

**Cons:**
- Conventions can entrench mediocre patterns; slavish consistency can block genuinely better designs.
- Cross-platform products face conflicting conventions (Ctrl vs Cmd, dialog button order) — consistency with the platform beats consistency with your other builds.

**Implementation guidance:**
- Maintain a terminology list and component library; audit for duplicate patterns quarterly.
- Follow platform HIGs for chrome, shortcuts, and dialog conventions; diverge only inside your content area.
- When migrating patterns, change everything at once rather than leaving mixed generations.
- Define component *usage rules*, not just visual specs; document intentional deviations with reasons.

**Decision questions:**
- Is there an established platform or industry pattern for this?
- Will users expect a behavior learned from other products?
- Is the deviation worth the learning cost it imposes?

**Related:** [Jakob's Law](#jakobs-law), [10-design-systems.md](10-design-systems.md), [14-platform-desktop.md](14-platform-desktop.md#platform-conventions).

**Sources:** [NN/g: Consistency and Standards](https://www.nngroup.com/articles/consistency-and-standards/), [GOV.UK Design System](https://design-system.service.gov.uk/).

### 5. Error prevention

**Definition:** Better than good error messages is a careful design that prevents problems from occurring: eliminate error-prone conditions or check for them and offer confirmation before the user commits.

**Reasoning / Evidence:** Nielsen/Molich heuristic #5. Distinguishes slips (unconscious, attention failures) from mistakes (conscious, mental-model failures) — see [02-cognitive-foundations.md](02-cognitive-foundations.md#slips-vs-mistakes). Prevention is cheaper than recovery for both user and business.

**When to use:**
- High-cost operations first (payments, deletions, sends, publishes).
- Input constraints: only offer valid choices (date pickers disable invalid dates; disable-with-explanation for unavailable options).
- Slip-prone contexts: adjacent dangerous buttons, similar-looking items, mode-dependent actions.

**When NOT to use / exceptions:**
- Over-constraining blocks legitimate edge cases (name fields rejecting apostrophes, address validators rejecting real addresses). Prefer warnings over hard blocks when the input might be valid.
- Excessive confirmations cause habituation and are themselves an error source.

**Pros:**
- Eliminates entire error classes rather than treating symptoms.
- Reduces support burden and user anxiety.

**Cons:**
- Constraints encode assumptions that may be wrong (formats, locales).
- Guardrails add friction for experts; balance with [flexibility](#7-flexibility-and-efficiency-of-use).

**Implementation guidance:**
- Prioritize: prevent high-cost errors first, then frequent small frustrations.
- Good defaults + helpful constraints prevent slips; undo + warnings prevent mistake damage.
- Separate destructive controls spatially (≥8px gap and visual distinction from primary actions); never make destructive the default focused button.
- Use confirmation only for irreversible, infrequent actions; use undo for reversible ones.
- Validate at the right moment: immediately for format constraints, on blur or submit for complex validation — never flag errors while the user is still typing ([07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing)).
- Use input masks cautiously; avoid trapping users in invalid intermediate states.

**Related:** [#3-user-control-and-freedom](#3-user-control-and-freedom), [07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing), [02-cognitive-foundations.md](02-cognitive-foundations.md#slips-vs-mistakes).

**Sources:** [NN/g: Preventing User Errors (slips)](https://www.nngroup.com/articles/slips/), [NN/g: Confirmation Dialogs](https://www.nngroup.com/articles/confirmation-dialog/).

### 6. Recognition rather than recall

**Definition:** Minimize memory load by making elements, actions, and options visible. Users should not have to remember information from one part of the interface to another.

**Reasoning / Evidence:** Nielsen/Molich heuristic #6. Human recognition memory vastly outperforms recall (it is easier to verify "Is Lisbon the capital of Portugal?" than to answer "What is the capital of Portugal?"). Working memory holds only ~4 chunks ([Miller's Law](#millers-law), Cowan's revision).

**When to use:**
- Menus over memorized commands for mainstream users; visible options over free-text where the option set is known.
- Show previously entered data when it's needed again (order summaries at checkout, comparison views).
- Autocomplete, recents, and history features.
- In-context help at the moment of need rather than upfront tutorials.

**When NOT to use / exceptions:**
- Expert interfaces legitimately trade recognition for recall-based speed (command lines, keyboard shortcuts, vim) — provide both paths ([#7-flexibility-and-efficiency-of-use](#7-flexibility-and-efficiency-of-use)).
- Showing everything at once conflicts with [minimalism](#8-aesthetic-and-minimalist-design); use progressive disclosure ([02-cognitive-foundations.md](02-cognitive-foundations.md#progressive-disclosure)).

**Pros:**
- Dramatically lowers learning curve and error rate.
- Serves intermittent users well (the majority of users of most products).

**Cons:**
- Visible-everything designs become cluttered; recognition requires screen space.
- Recognition-driven UIs can be slower for experts than recall-based accelerators.

**Implementation guidance:**
- Never require users to transcribe information between screens; carry it forward automatically (WCAG 3.3.7 Redundant Entry).
- Field labels stay visible while typing (no placeholder-only labels — [07-forms-and-input.md](07-forms-and-input.md#labels-vs-placeholders)).
- Command palettes (Ctrl/Cmd+K) blend recall entry with recognition via fuzzy-matched suggestions.
- Keep current context visible: selected filters, sort order, scope, and the object being acted on; provide examples and previews at decision points.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#recognition-vs-recall), [Miller's Law](#millers-law).

**Sources:** [NN/g: Recognition and Recall in UX](https://www.nngroup.com/articles/recognition-and-recall/).

### 7. Flexibility and efficiency of use

**Definition:** Accelerators — hidden from novices — speed up interaction for experts, so the design caters to both inexperienced and experienced users. Allow users to tailor frequent actions.

**Reasoning / Evidence:** Nielsen/Molich heuristic #7. Users progress from novice to expert on frequent tasks; a UI that only serves novices caps expert productivity, one that only serves experts excludes newcomers.

**When to use:**
- Tools used daily/professionally: keyboard shortcuts, bulk operations, macros, templates, saved views, command palette.
- Personalization (system-driven) and customization (user-driven) of frequent actions, dashboards, toolbars.
- Touch: gesture shortcuts (swipe actions) duplicating visible buttons.

**When NOT to use / exceptions:**
- Rarely-used products (annual tax filing) gain little from accelerators; invest in guidance instead.
- Do not make accelerators the only path — every gesture/shortcut needs a discoverable visible equivalent.

**Pros:**
- Expert throughput scales; power users become advocates.
- Customization increases fit for diverse workflows.

**Cons:**
- Customization features are costly to build/maintain and can fragment support ("your UI looks different from the docs").
- Hidden accelerators are undiscoverable without deliberate teaching (shortcut hints in menus/tooltips).

**Implementation guidance:**
- Standard shortcut set first (Ctrl/Cmd+S/Z/C/V/F/N/W); show shortcuts in menus and tooltips.
- Support batch selection + bulk actions in any list users manage at scale.
- Persist user choices (view, sort, density) per user; sync across devices where possible.
- Track feature usage: accelerators used by <1% may indicate discoverability failure, not uselessness.

**Related:** [14-platform-desktop.md](14-platform-desktop.md#keyboard-shortcuts-design), [16-platform-cli-tui.md](16-platform-cli-tui.md), [#6-recognition-rather-than-recall](#6-recognition-rather-than-recall).

**Sources:** [NN/g: Flexibility and Efficiency of Use](https://www.nngroup.com/articles/flexibility-efficiency-heuristic/), [NN/g: UI Accelerators](https://www.nngroup.com/articles/ui-accelerators/).

### 8. Aesthetic and minimalist design

**Definition:** Interfaces should not contain information that is irrelevant or rarely needed; every extra unit of information competes with the relevant units and diminishes their relative visibility.

**Reasoning / Evidence:** Nielsen/Molich heuristic #8. Signal-to-noise: attention is finite; each added element taxes the [visual hierarchy](03-visual-design.md#visual-hierarchy). This is about content focus, not flat aesthetics.

**When to use:**
- Every screen: ruthlessly prioritize toward the primary task; move secondary content behind progressive disclosure.
- Marketing/landing pages: one primary action per view.
- Dashboards: show the decision-driving numbers, not everything measurable.

**When NOT to use / exceptions:**
- Data-dense professional tools (trading, ops consoles, EMRs) legitimately maximize information density — minimalism there means removing decoration, not data ([03-visual-design.md](03-visual-design.md#density)).
- Minimalism taken too far removes signifiers users need (icon-only mystery UIs, hidden navigation) — see [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Pros:**
- Faster comprehension and task focus; better conversion on focused pages.
- Less content to maintain and localize.

**Cons:**
- "Minimal" is often misread as hiding essentials (labels, affordances) — harming usability.
- Deciding what to cut requires real prioritization work and stakeholder conflict.

**Implementation guidance:**
- For each element ask: does it support the user's primary goal on this screen? If not, remove, demote, or defer it.
- Word count: aim to cut first drafts by ~50% for web reading (users read ~20–28% of words on a page).
- One primary button per view; secondary actions styled subordinately.

**Related:** [03-visual-design.md](03-visual-design.md#visual-hierarchy), [02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-load-theory), [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Sources:** [NN/g: Aesthetic and Minimalist Design](https://www.nngroup.com/articles/aesthetic-minimalist-design/), [NN/g: How Little Do Users Read?](https://www.nngroup.com/articles/how-little-do-users-read/).

### 9. Help users recognize, diagnose, and recover from errors

**Definition:** Error messages should be expressed in plain language (no error codes alone), precisely indicate the problem, and constructively suggest a solution.

**Reasoning / Evidence:** Nielsen/Molich heuristic #9. An error state is where users are most stressed and most likely to abandon; a good message converts failure into recovery.

**When to use:**
- Every error surface: forms, API failures, empty results, permission blocks, offline.
- Pair visual salience (red, icon, bold) with textual clarity; announce to screen readers (aria-live / role=alert).

**When NOT to use / exceptions:**
- Do not use error styling for warnings or info — reserve red for actual errors to prevent habituation.
- Internal error codes may be appended for support ("Error 5031") but never replace the human explanation.

**Pros:**
- Recovers otherwise-lost conversions and prevents support contacts.
- Builds trust ("the product tells me what's wrong and how to fix it").

**Cons:**
- Specific messages require engineering effort to distinguish failure causes.
- Overly specific auth errors can leak security information ("wrong password" vs "invalid credentials" trade-off).

**Implementation guidance:**
- Structure: what happened → why (if known) → how to fix → escape hatch (retry, support link).
- Place field errors adjacent to the field; add a summary with anchor links atop long forms ([07-forms-and-input.md](07-forms-and-input.md#error-message-design)).
- Preserve all user input through the error.
- Tone: no blame, no humor in high-stakes failures, no ALL CAPS ([09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing)).

**Related:** [07-forms-and-input.md](07-forms-and-input.md#error-message-design), [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing), [#5-error-prevention](#5-error-prevention).

**Sources:** [NN/g: Error Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/), [WCAG 2.2: Input Assistance (Guideline 3.3)](https://www.w3.org/WAI/WCAG22/Understanding/input-assistance.html).

### 10. Help and documentation

**Definition:** Best if the system needs no explanation, but when needed, provide help that is easy to search, focused on the user's task, concise, and lists concrete steps.

**Reasoning / Evidence:** Nielsen/Molich heuristic #10. Users consult help reluctantly, in task context, wanting an answer, not a manual. In-context help outperforms separate documentation for task completion.

**When to use:**
- Contextual help at the point of need: tooltips on unfamiliar controls, "what's this?" popovers, inline examples, empty-state guidance.
- Searchable help center for complex products; onboarding checklists for multi-step setup.
- Complex domain rules (tax, permissions) that no UI simplification can remove ([Tesler's Law](#teslers-law)).

**When NOT to use / exceptions:**
- Help must never excuse bad design ("we'll explain it in the tooltip") — fix the design first.
- Forced tutorial tours on first run are widely skipped and resented; prefer do-first learning with optional guidance ([18-patterns-antipatterns.md](18-patterns-antipatterns.md)).

**Pros:**
- Safety net raising task completion for complex products; deflects support tickets.

**Cons:**
- Documentation rots as UI changes; maintenance is a real ongoing cost.
- Overused tooltips/popovers clutter and annoy.

**Implementation guidance:**
- Write help as tasks ("How to export a report"), not features ("The Export module").
- Keep articles scannable: steps, screenshots, <300 words per task where possible.
- Instrument help search queries — they are a direct map of UX failures.

**Related:** [09-content-ux-writing.md](09-content-ux-writing.md#onboarding-copy-and-tooltips), [#8-aesthetic-and-minimalist-design](#8-aesthetic-and-minimalist-design).

**Sources:** [NN/g: Help and Documentation](https://www.nngroup.com/articles/help-and-documentation/), [GOV.UK Service Manual](https://www.gov.uk/service-manual).

---

## Norman's design principles

From Don Norman's *The Design of Everyday Things* (1988, revised 2013) — the vocabulary of how people understand and operate designed things.

### Affordances

**Definition:** The possible actions an object allows relative to the capabilities of an actor — a button affords pressing, a text field affords typing. In screen design what matters is *perceived* affordance: what the user believes is possible.

**Reasoning / Evidence:** Norman adapted J.J. Gibson's ecological psychology term. In GUIs, affordances are conveyed by convention and rendering (pixels afford nothing physically); Norman later introduced "signifiers" to correct widespread misuse of "affordance" for the visible cue.

**When to use:**
- Design every interactive element so its possible actions match user capabilities and expectations (draggable things must actually be draggable everywhere they appear draggable).
- Audit for false affordances (looks clickable but isn't) and hidden affordances (is clickable but doesn't look it).

**When NOT to use / exceptions:**
- Not everything possible should be signified prominently — expert-only affordances can stay quiet (but never *invisible* as the sole path).

**Pros:**
- Correct perceived affordances make interfaces self-explanatory.

**Cons:**
- Flat design has systematically weakened perceived affordances (borderless buttons, unstyled links) — measurable click-uncertainty cost documented by NN/g.

**Implementation guidance:**
- Interactive elements need ≥1 persistent visual cue (shape, border, color, elevation) plus state feedback (hover/focus/pressed).
- Text that acts must look actionable; text that doesn't act must not (no underlined non-links).
- Test: screenshot a screen, ask users to mark everything clickable — target >90% accuracy.

**Related:** [Signifiers](#signifiers), [03-visual-design.md](03-visual-design.md), [04-interaction-design.md](04-interaction-design.md#affordances-and-signifiers-in-practice).

**Sources:** [Norman, The Design of Everyday Things](https://mitpress.mit.edu/9780262525671/the-design-of-everyday-things/), [NN/g: Affordances](https://www.nngroup.com/articles/affordances-clickability-signifiers/).

### Signifiers

**Definition:** Perceivable cues that communicate where and how to act — labels, icons, highlights, textures, sounds. Signifiers signal affordances.

**Reasoning / Evidence:** Norman (2013 revision) introduced the term because designers can't add affordances to a screen, only *signals* of them. Weak signifiers are the root cause of "mystery meat" UI.

**When to use:**
- Wherever an available action is not self-evident: scroll indicators for below-the-fold content, drag handles, chevrons for expandable rows, shortcut hints.
- Gesture-driven UIs: on-screen hints for available swipes (gestures have no natural signifier).

**When NOT to use / exceptions:**
- Signifier overload (everything shouting) destroys hierarchy; strength of signifier should track importance and frequency of the action.

**Pros:**
- Direct lever on discoverability; cheap to add relative to redesign.

**Cons:**
- Aesthetic pressure to remove signifiers is constant; each removal trades beauty for findability and the cost is invisible until measured.

**Implementation guidance:**
- Minimum clickability signifiers: buttons look raised/filled/bordered; links colored and/or underlined in body text.
- Show partial content at fold edges (cut-off card) to signify scrollability.
- Drag: 6-dot/2-bar handle iconography; resize: corner ribbing; expand: chevron pointing in expansion direction.

**Related:** [Affordances](#affordances), [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures), [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Sources:** [NN/g: Signifiers vs Affordances](https://www.nngroup.com/articles/affordances-clickability-signifiers/).

### Mapping

**Definition:** The relationship between controls and their effects. Natural mapping exploits spatial analogy and cultural convention so the control layout itself explains what each control does (stove knobs matching burner layout).

**Reasoning / Evidence:** Norman's stove-control studies: mismatched knob layouts cause persistent errors even in experienced users. Good mapping eliminates the need for labels and memory.

**When to use:**
- Spatial controls: volume sliders vertical (up=more), brightness, media scrubbers left-to-right (in LTR locales).
- Layout of controls that act on screen regions (alignment buttons showing their result).
- Multi-monitor/multi-device arrangements mirroring physical positions.

**When NOT to use / exceptions:**
- Cultural mappings vary (RTL locales reverse horizontal progressions; light-switch direction conventions differ by country) — verify per market ([09-content-ux-writing.md](09-content-ux-writing.md#internationalization-and-localization)).

**Pros:**
- Zero-learning controls; error resistance.

**Cons:**
- Not all functions have a spatial analog; forcing one confuses more than labels would.

**Implementation guidance:**
- Increase = up/right/clockwise (LTR); keep the mapping consistent across the whole product.
- Place actions next to the objects they act on (contextual toolbars), not in distant global chrome.

**Related:** [#2-match-between-system-and-the-real-world](#2-match-between-system-and-the-real-world), [Continuity](#continuity).

**Sources:** [NN/g: Natural Mappings](https://www.nngroup.com/articles/natural-mappings/).

### Feedback

**Definition:** Communicating the result of an action — immediate, informative confirmation that the system received the input and what changed.

**Reasoning / Evidence:** Norman's action-cycle: users evaluate the world after acting ("gulf of evaluation"); missing feedback leaves them unable to tell success from failure, causing repeated actions and errors.

**When to use:**
- Every user action gets acknowledgment within 100ms (state change, ripple, highlight), even if completion comes later.
- Multi-channel where stakes are high: visual + haptic (mobile) + sound (opt-in).

**When NOT to use / exceptions:**
- Excess feedback (a toast for every trivial action) trains users to ignore all feedback; scale feedback intensity with action importance.

**Pros:**
- Prevents double-submission, builds perceived responsiveness and trust.

**Cons:**
- Feedback choreography takes design effort per action type; sloppy global patterns (toast-everything) backfire.

**Implementation guidance:**
- Tiered: trivial action → inline state change only; significant → transient toast (4–8s) with undo; critical → persistent confirmation requiring acknowledgment.
- Disable-and-spinner submit buttons during in-flight requests; re-enable on failure with error.

**Related:** [#1-visibility-of-system-status](#1-visibility-of-system-status), [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds).

**Sources:** [Norman, The Design of Everyday Things](https://mitpress.mit.edu/9780262525671/the-design-of-everyday-things/).

### Constraints

**Definition:** Restrictions — physical, logical, semantic, cultural — that limit possible actions and thereby guide correct behavior (a plug that only fits one way).

**Reasoning / Evidence:** Norman: constraints reduce the action space users must consider, preventing errors before they occur. Digital analogs: disabled states, input masks, restricted choices, forced sequence.

**When to use:**
- Error-prone inputs: pickers over free text where the valid set is closed (dates, countries).
- Sequence enforcement in flows with real dependencies.
- Preventing invalid combinations (disable conflicting options with explanation).

**When NOT to use / exceptions:**
- Over-constraint blocks legitimate use (rigid name/address/phone validators are notorious). Where reality is messier than your model, warn instead of block.
- Disabled controls without explanation are a usability black hole — always say why and how to enable ([18-patterns-antipatterns.md](18-patterns-antipatterns.md)).

**Pros:**
- Prevents whole error classes; reduces decisions ([Hick's Law](#hicks-law)).

**Cons:**
- Encodes assumptions; every false constraint is a support ticket or lost user.

**Implementation guidance:**
- Prefer soft constraints (warnings, confirmations) when validity is probabilistic; hard constraints only for logically impossible states.
- Pair every disabled state with a tooltip/note explaining the unlock condition.

**Related:** [#5-error-prevention](#5-error-prevention), [07-forms-and-input.md](07-forms-and-input.md#input-type-selection).

**Sources:** [Norman, The Design of Everyday Things](https://mitpress.mit.edu/9780262525671/the-design-of-everyday-things/).

### Conceptual models

**Definition:** The user's internal explanation of how the system works, built from the system image (UI, docs, behavior). Good design projects a system image that produces an accurate, simple conceptual model.

**Reasoning / Evidence:** Norman: designers have a design model; users build their model only from what they see. Gaps between them cause "mistakes" (correct execution of wrong plans). Classic example: thermostat believed to be a valve ("crank it to heat faster") vs its true on/off nature.

**When to use:**
- Novel functionality: explicitly decide what model you want users to form, then design every surface to teach it.
- Sync/offline, permissions, and AI features — areas where wrong models cause severe errors and mistrust.

**When NOT to use / exceptions:**
- The projected model should be *useful*, not technically complete — abstract away internals that don't change user decisions.

**Pros:**
- Accurate models let users predict behavior in situations you never explicitly designed for.

**Cons:**
- Model design is invisible work; teams skip it and ship contradictory metaphors.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#mental-models), [#2-match-between-system-and-the-real-world](#2-match-between-system-and-the-real-world).

**Sources:** [NN/g: Mental Models](https://www.nngroup.com/articles/mental-models/), [Norman, The Design of Everyday Things](https://mitpress.mit.edu/9780262525671/the-design-of-everyday-things/).

---

## Gestalt principles

Perceptual grouping laws from 1920s Gestalt psychology (Wertheimer, Koffka, Köhler). They describe how humans automatically organize visual fields — layout either works with these laws or fights them.

### Proximity

**Definition:** Elements close together are perceived as a group; elements far apart as separate.

**Reasoning / Evidence:** One of the strongest and most reliable grouping principles ([common region](#common-region) is empirically the strongest); operates pre-attentively (before conscious processing). Whitespace is therefore a functional grouping tool, not empty decoration.

**When to use:**
- Group labels with their fields, buttons with the content they act on, related settings into clusters.
- Use spacing *instead of* boxes/dividers where possible — cleaner and equally effective.

**When NOT to use / exceptions:**
- Proximity is overridden by strong dissimilarity; don't rely on closeness alone when items look wildly different.

**Pros:** Free (no added elements); the primary tool for reducing perceived complexity.
**Cons:** Cramped layouts accidentally group unrelated things; the errors are invisible to the team that knows the intent.

**Implementation guidance:**
- Spacing ratio rule: space *within* a group must be clearly smaller than space *between* groups — minimum ~1:1.5, comfortable 1:2 (e.g., 8px intra / 16–24px inter).
- Field label to its input: 4–8px; between form groups: 24–32px ([07-forms-and-input.md](07-forms-and-input.md#form-layout)).

**Related:** [03-visual-design.md](03-visual-design.md#spacing-and-layout), [Common region](#common-region).

**Sources:** [NN/g: Proximity Principle](https://www.nngroup.com/articles/gestalt-proximity/).

### Similarity

**Definition:** Elements sharing visual attributes (color, shape, size, orientation) are perceived as related or same-kind.

**Reasoning / Evidence:** Gestalt grouping; color is the strongest similarity cue, then size, then shape. Basis of every component system: all primary buttons look alike, so users transfer learning.

**When to use:**
- Signal same-behavior via same-appearance across the product (all links one style; all destructive actions one style).
- Data viz: same series = same color.

**When NOT to use / exceptions:**
- Different-behavior elements must NOT look similar — visual similarity between a button and a static badge creates false affordance.

**Pros:** Enables instant category recognition; the mechanism behind design-system consistency.
**Cons:** Accidental similarity (decorative element matching link color) silently miscommunicates.

**Implementation guidance:**
- Define one visual identity per interactive class in tokens ([10-design-systems.md](10-design-systems.md#design-tokens)); lint for off-palette colors on interactive elements.
- Reserve the accent color for interactive elements only.

**Related:** [#4-consistency-and-standards](#4-consistency-and-standards), [03-visual-design.md](03-visual-design.md#color).

**Sources:** [NN/g: Similarity Principle](https://www.nngroup.com/articles/gestalt-similarity/).

### Closure

**Definition:** The mind completes incomplete shapes — perceiving a whole even when parts are missing.

**Reasoning / Evidence:** Gestalt completion; allows minimal iconography (a three-line hamburger reads as a list; a partially cut card implies continuation).

**When to use:**
- Icon design: suggest forms with minimal strokes.
- Scrollability signifiers: cut content at the viewport edge so users infer more below/beside.
- Skeleton screens: rough shapes suffice to imply the coming layout.

**When NOT to use / exceptions:**
- Excessive abstraction breaks closure — if fewer than ~most users complete the intended shape, the icon fails; test recognition.

**Pros:** Enables clean minimal visuals with retained meaning.
**Cons:** Cultural/experience-dependent; young users famously fail to recognize the floppy-disk save icon's referent (though they learn the convention).

**Implementation guidance:**
- Carousel/scroll regions: show 10–20% of the next item peeking at the edge.
- Test icons at target size (16–24px) with recognition tasks before shipping icon-only.

**Related:** [03-visual-design.md](03-visual-design.md#iconography), [Continuity](#continuity).

**Sources:** [NN/g: Closure Principle](https://www.nngroup.com/articles/principle-closure/).

### Continuity

**Definition:** Elements arranged on a line or curve are perceived as related and the eye follows the line.

**Reasoning / Evidence:** Gestalt "good continuation." Alignment creates invisible lines that guide scanning; misalignment breaks the perceived flow and reads as disorder or separation.

**When to use:**
- Align related content to shared edges (left-align form labels/fields, grid columns).
- Horizontal scrolling rows (items aligned on one line read as one set).
- Guide the eye through a flow: sequential steps on a visual line.

**When NOT to use / exceptions:**
- Deliberate misalignment can signal separation or draw attention — use sparingly and intentionally.

**Pros:** Alignment is free structure; strongly improves scannability.
**Cons:** Accidental alignments group unrelated items across container boundaries.

**Implementation guidance:**
- Use a layout grid; snap everything ([03-visual-design.md](03-visual-design.md#spacing-and-layout)).
- Prefer one shared left edge per screen region; every additional alignment axis adds scanning cost.

**Related:** [03-visual-design.md](03-visual-design.md#spacing-and-layout), [02-cognitive-foundations.md](02-cognitive-foundations.md#scanning-patterns).

**Sources:** [NN/g: Continuity](https://www.nngroup.com/articles/gestalt-principles-continuity/).

### Common region

**Definition:** Elements within a shared enclosed boundary (card, panel, outline, background) are perceived as grouped — even overriding proximity.

**Reasoning / Evidence:** Palmer (1992) extension of Gestalt grouping; empirically stronger than proximity. This is why cards work.

**When to use:**
- Cards for heterogeneous repeating content (image+title+meta+actions).
- Panels/wells to bind a control cluster (filter groups, dialog sections).
- Table row banding or hover regions to bind wide rows.

**When NOT to use / exceptions:**
- Boxes-in-boxes-in-boxes creates clutter; if proximity alone can group it, skip the container ([#8-aesthetic-and-minimalist-design](#8-aesthetic-and-minimalist-design)).

**Pros:** Strongest grouping tool; makes complex layouts parse instantly.
**Cons:** Container chrome (borders, shadows) accumulates visual noise; each nesting level costs padding/space.

**Implementation guidance:**
- Max ~2 levels of visible containment per screen region; use whitespace for deeper structure.
- Card minimum padding 16px; consistent corner radius and elevation per level ([03-visual-design.md](03-visual-design.md#elevation-shadows-depth)).

**Related:** [Proximity](#proximity), [03-visual-design.md](03-visual-design.md#spacing-and-layout).

**Sources:** [NN/g: Common Region](https://www.nngroup.com/articles/common-region/).

### Figure/ground

**Definition:** Perception separates a scene into a focal figure and receding background; design controls what reads as figure.

**Reasoning / Evidence:** Foundational Gestalt observation (Rubin's vase). Modals with scrims, elevated cards, and focus overlays all exploit deliberate figure/ground manipulation.

**When to use:**
- Modal/dialog scrims (dim background 40–60% black) to force figure status.
- Elevation/shadow to lift interactive surfaces above the page.
- Spotlight onboarding (dim all but the highlighted element).

**When NOT to use / exceptions:**
- Ambiguous figure/ground (busy backgrounds behind text, low contrast layering) destroys readability — text over images requires scrims/gradients.

**Pros:** Directs attention powerfully; communicates hierarchy of layers.
**Cons:** Overuse of overlays interrupts flow ([18-patterns-antipatterns.md](18-patterns-antipatterns.md) — modal overuse).

**Implementation guidance:**
- Text over imagery: add a ≥40% opacity scrim or check contrast at the worst point (still ≥4.5:1).
- Keep one dominant figure per screen; competing "pop" elements cancel out.

**Related:** [03-visual-design.md](03-visual-design.md#elevation-shadows-depth), [#8-aesthetic-and-minimalist-design](#8-aesthetic-and-minimalist-design).

**Sources:** [NN/g: Figure/Ground](https://www.nngroup.com/articles/figure-ground-gestalt/).

### Common fate

**Definition:** Elements moving together (same direction/timing) are perceived as a group.

**Reasoning / Evidence:** Gestalt motion grouping; in UI, synchronized animation binds elements (a panel and its contents sliding in together read as one object).

**When to use:**
- Transition choreography: animate related elements as one unit; stagger unrelated groups.
- Drag interactions: everything belonging to the dragged item moves with it.
- Expanding/collapsing: children move with the parent.

**When NOT to use / exceptions:**
- Honor `prefers-reduced-motion` — grouping must not depend solely on motion ([05-accessibility.md](05-accessibility.md#reduced-motion)).

**Pros:** Communicates relationships during change, when static grouping cues are in flux.
**Cons:** Motion-dependent meaning excludes reduced-motion users and is easy to overdo.

**Implementation guidance:**
- Group members: identical duration/easing/direction. Separate groups: stagger 20–50ms.
- Keep functional motion ≤300ms ([04-interaction-design.md](04-interaction-design.md#motion-and-animation)).

**Related:** [04-interaction-design.md](04-interaction-design.md#motion-and-animation).

**Sources:** [NN/g: Common Fate](https://www.nngroup.com/articles/common-fate/).

---

## Laws of UX

Named empirical regularities, popularized by Jon Yablonski's lawsofux.com. Treat them as strong defaults with known boundary conditions, not physics.

### Fitts's Law

**Definition:** Time to acquire a target is a function of the distance to and size of the target: `MT = a + b·log2(2D/W)`. Bigger and closer = faster to hit.

**Reasoning / Evidence:** Paul Fitts (1954), one of the most replicated findings in HCI; extends to touch and mouse pointing.

**When to use:**
- Size frequently-used and critical controls generously; place them near where the pointer/thumb already is.
- Screen edges/corners are effectively infinite-size targets for mouse (pointer pins there) — OS-level menus/corners exploit this.
- Context menus and radial menus minimize D by appearing at the cursor.
- Mobile: primary actions in the thumb zone ([15-platform-mobile.md](15-platform-mobile.md#thumb-zone)).

**When NOT to use / exceptions:**
- Deliberately *slow down* destructive actions: smaller, farther from common paths (inverse Fitts).
- Edge-pinning doesn't apply inside browser windows or touchscreens (no pointer pinning).

**Pros:** Quantitative, directly actionable sizing/placement rule.
**Cons:** Says nothing about findability or comprehension — a huge close button still fails if unlabeled.

**Implementation guidance:**
- Touch targets ≥44×44pt (Apple) / 48×48dp (Material) / ≥24×24 CSS px absolute minimum (WCAG 2.5.8) with ≥8px spacing.
- Make the whole card/row clickable, not just the title link.
- Put confirm buttons near the triggering context, not across the screen.

**Related:** [05-accessibility.md](05-accessibility.md#target-size), [15-platform-mobile.md](15-platform-mobile.md#touch-targets), [#5-error-prevention](#5-error-prevention).

**Sources:** [Laws of UX: Fitts's Law](https://lawsofux.com/fitts-law/), [Fitts 1954](https://psycnet.apa.org/record/1955-02059-001).

### Hick's Law

**Definition:** Decision time increases logarithmically with the number and complexity of choices: `RT = a + b·log2(n+1)`.

**Reasoning / Evidence:** Hick (1952), Hyman (1953). Applies to simple reaction choices; complex deliberation (comparing products) doesn't follow the neat log curve but still degrades with option count (choice overload — [02-cognitive-foundations.md](02-cognitive-foundations.md#decision-fatigue-and-choice-overload)).

**When to use:**
- Trim menus, toolbars, plan tables, and settings to the decision-relevant set.
- Break complex decisions into sequenced smaller ones (wizards).
- Highlight a recommended option to short-circuit comparison.

**When NOT to use / exceptions:**
- Experts scanning a known layout (e.g., a dense toolbar they've memorized) aren't making Hick-style decisions — density can beat depth for them.
- Removing options users need just moves the cost elsewhere (search, support).
- Categorized long lists behave better than the raw count suggests (users decide per-category first).

**Pros:** Justifies simplification quantitatively; pairs with progressive disclosure.
**Cons:** Misapplied as "always fewer options," harming expert efficiency and choice fairness.

**Implementation guidance:**
- Top-level nav: ~5–7 items; menus >10 items need grouping headers or search.
- Onboarding: one decision per screen.
- Default + "advanced" split: 80% path visible, tail behind disclosure ([Pareto](#pareto-principle)).

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#decision-fatigue-and-choice-overload), [02-cognitive-foundations.md](02-cognitive-foundations.md#progressive-disclosure), [08-navigation-ia.md](08-navigation-ia.md#hierarchy-depth-vs-breadth).

**Sources:** [Laws of UX: Hick's Law](https://lawsofux.com/hicks-law/).

### Jakob's Law

**Definition:** Users spend most of their time on other sites/apps, so they prefer your product to work the same way as everything they already know.

**Reasoning / Evidence:** Coined by Jakob Nielsen (2000). The empirical basis is transfer of learning: existing mental models dominate first-use behavior.

**When to use:**
- Structural conventions: logo top-left→home, cart top-right, search with magnifier, settings gear, checkout flow order.
- New product categories: anchor to the conventions of the nearest familiar category.

**When NOT to use / exceptions:**
- Differentiation where it's your core value is fine — innovate on the substance, keep the chrome conventional ("be different in what you do, familiar in how you're operated").
- When the dominant convention is a dark pattern or measurably poor, don't copy it.

**Pros:** Free usability from day one; lowers onboarding cost to near zero.
**Cons:** Herd effects entrench mediocrity; conventions differ by region/platform, so "the convention" needs verifying per market.

**Implementation guidance:**
- Before designing any common flow (auth, checkout, search, settings), review 3–5 category leaders and default to their shared skeleton.
- Deviate only with test evidence; measure first-use success against the conventional control.

**Related:** [#4-consistency-and-standards](#4-consistency-and-standards), [02-cognitive-foundations.md](02-cognitive-foundations.md#mental-models).

**Sources:** [Laws of UX: Jakob's Law](https://lawsofux.com/jakobs-law/), [NN/g: Jakob's Law video](https://www.nngroup.com/videos/jakobs-law-internet-ux/).

### Miller's Law

**Definition:** Average working memory holds about 7±2 chunks (Miller 1956); modern estimates say ~4±1 (Cowan 2001). Chunking — grouping items into meaningful units — extends effective capacity.

**Reasoning / Evidence:** Miller's paper is often over-applied ("max 7 menu items" is a myth — visible menus are recognition, not recall). The real lesson is chunking and minimizing *memory* demands, not capping visible options.

**When to use:**
- Anything users must *hold in memory*: confirmation codes, instructions spanning screens, comparison across views.
- Chunk long strings: phone numbers, card numbers (4-4-4-4), license keys.
- Limit steps users must remember without visible state.

**When NOT to use / exceptions:**
- Do NOT use it to cap visible lists/menus — visible items require recognition, not recall ([#6-recognition-rather-than-recall](#6-recognition-rather-than-recall)).

**Pros:** Grounds chunking and reduces memory-based errors.
**Cons:** The 7±2 number is misquoted constantly; design to 3–4 chunks for safety.

**Implementation guidance:**
- Verification codes: ≤6 digits, displayed chunked (123 456), with copy button.
- Never require transcribing between screens; show the needed info in place (WCAG 3.3.7).
- Group form fields into labeled sections of 3–5 fields.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking), [#6-recognition-rather-than-recall](#6-recognition-rather-than-recall).

**Sources:** [Laws of UX: Miller's Law](https://lawsofux.com/millers-law/), [Cowan 2001](https://doi.org/10.1017/S0140525X01003922).

### Tesler's Law

**Definition:** (Law of conservation of complexity) Every application has an irreducible amount of complexity; the only question is who deals with it — the user or the developer/designer.

**Reasoning / Evidence:** Larry Tesler (1980s, Xerox PARC). Arguing that engineers should absorb complexity because millions of users each paying a little attention costs more in aggregate than engineering effort once.

**When to use:**
- Defaults, auto-detection, and inference: detect card type from number, infer date format, parse pasted data flexibly.
- Complex domains: absorb rules into the system (tax software computing instead of explaining).
- API/CLI design: sensible defaults over mandatory configuration ([16-platform-cli-tui.md](16-platform-cli-tui.md)).

**When NOT to use / exceptions:**
- Some complexity is the user's essential decision (medical consent, financial risk) — hiding it is unethical, not simplification.
- Over-automation that guesses wrong and hides the controls is worse than asking.

**Pros:** Differentiates genuinely simple products from superficially minimal ones.
**Cons:** Absorbed complexity = engineering cost and edge-case risk; wrong inference frustrates more than an honest question.

**Implementation guidance:**
- For each user input ask: can we detect, default, defer, or drop it?
- Where inference is used, show what was inferred and make correction one action.

**Related:** [#8-aesthetic-and-minimalist-design](#8-aesthetic-and-minimalist-design), [07-forms-and-input.md](07-forms-and-input.md#defaults-and-prefill).

**Sources:** [Laws of UX: Tesler's Law](https://lawsofux.com/teslers-law/).

### Doherty Threshold

**Definition:** Productivity soars when computer and user interact at a pace (<400ms response) where neither waits on the other.

**Reasoning / Evidence:** Doherty & Thadani (IBM, 1982): sub-400ms system response kept users in flow and improved output quality, not just speed.

**When to use:**
- Performance budgeting for interactive operations: target <400ms perceived response for common actions.
- When >400ms is unavoidable, manage perception: optimistic UI, skeletons, progressive results.

**When NOT to use / exceptions:**
- A deliberate short delay can *increase* trust in specific contexts (security scan, "calculating your plan") — fake-work perception, use sparingly and honestly ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#fake-loading)).

**Pros:** Concrete latency budget tied to productivity evidence.
**Cons:** 400ms is a business-computing figure; direct manipulation needs <100ms, and modern expectations keep tightening (INP <200ms — [13-platform-web.md](13-platform-web.md#core-web-vitals)).

**Implementation guidance:**
- Budget: input feedback <100ms; common operations <400ms; anything >1s gets a progress indicator.
- Precompute/prefetch likely next steps; cache aggressively; render partial results early.

**Related:** [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds), [13-platform-web.md](13-platform-web.md#core-web-vitals).

**Sources:** [Laws of UX: Doherty Threshold](https://lawsofux.com/doherty-threshold/).

### Peak-End Rule

**Definition:** People judge an experience by its most intense moment (peak) and its end, not the average of every moment.

**Reasoning / Evidence:** Kahneman & colleagues (colonoscopy and cold-water studies, 1990s): remembered experience ≠ lived experience; duration is largely neglected.

**When to use:**
- Invest disproportionately in: the worst moment of your journey (kill the deepest pain point) and the final step (order confirmation, completion screens, offboarding).
- Add designed positive peaks: delight at milestone completion.
- Error recovery: a brilliantly handled failure can become a positive peak.

**When NOT to use / exceptions:**
- Not a license to neglect the middle: severe average friction still causes in-session abandonment before any "memory" forms.
- Manipulative peak-engineering to mask systemic problems erodes trust when users compare notes.

**Pros:** Focuses polish investment where memory (and reviews, and retention) is actually formed.
**Cons:** Peaks are subjective; requires journey research to find the real lows.

**Implementation guidance:**
- Map the journey; identify the lowest moment and the last moment. Fix the low, celebrate the end.
- Ends include: checkout confirmation, ticket-resolved message, unsubscribe/offboarding (a graceful exit drives return and referrals).

**Related:** [17-ux-process-research.md](17-ux-process-research.md#journey-mapping), [09-content-ux-writing.md](09-content-ux-writing.md).

**Sources:** [Laws of UX: Peak-End Rule](https://lawsofux.com/peak-end-rule/), [Kahneman et al. 1993](https://doi.org/10.1111/j.1467-9280.1993.tb00589.x).

### Aesthetic-Usability Effect

**Definition:** Users perceive attractive designs as more usable — and are more tolerant of minor usability problems in them.

**Reasoning / Evidence:** Kurosu & Kashimura (1995, ATM layouts); replicated by Tractinsky ("What is beautiful is usable," 2000). Visual appeal creates positive affect that increases tolerance and persistence.

**When to use:**
- Justify visual polish investment: aesthetics measurably buys goodwill, first impressions (~50ms judgment), and error tolerance.
- Competitive parity products: aesthetics can be the deciding factor.

**When NOT to use / exceptions:**
- The effect *masks* problems in usability testing — users report beautiful broken things as "easy." Probe behavior, not just satisfaction ratings ([17-ux-process-research.md](17-ux-process-research.md#usability-testing)).
- It buys tolerance for *minor* issues only; severe task failures burn through goodwill fast.

**Pros:** Real, replicated effect; aligns brand and UX investment.
**Cons:** Can be weaponized to ship pretty-but-broken; inflates self-reported metrics.

**Implementation guidance:**
- In tests, weight task success/time/errors over satisfaction scores for attractive prototypes.
- Never conclude "no usability issues" from a visually polished prototype without behavioral data.

**Related:** [03-visual-design.md](03-visual-design.md), [17-ux-process-research.md](17-ux-process-research.md#usability-testing).

**Sources:** [Laws of UX: Aesthetic-Usability Effect](https://lawsofux.com/aesthetic-usability-effect/), [NN/g: Aesthetic-Usability Effect](https://www.nngroup.com/articles/aesthetic-usability-effect/).

### Von Restorff Effect

**Definition:** (Isolation effect) In a set of similar items, the one that differs is most likely to be noticed and remembered.

**Reasoning / Evidence:** Hedwig von Restorff (1933). Distinctiveness captures attention pre-attentively; basis of CTA highlighting.

**When to use:**
- One primary CTA per view, visually isolated (color/size/weight).
- Highlighting the recommended plan in pricing tables.
- Critical status flags in dense tables (one red dot among grey).

**When NOT to use / exceptions:**
- If everything is highlighted, nothing is — the effect requires a uniform field to stand out from.
- Don't exploit it to make the business-preferred option deceptively prominent over user interest ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#preselection)).

**Pros:** Cheap, reliable attention direction.
**Cons:** Ceiling of one (maybe two) isolates per view; color-only isolation fails color-blind users ([05-accessibility.md](05-accessibility.md#color-independence)).

**Implementation guidance:**
- Isolate with ≥2 combined cues (color + size, color + position).
- Audit each screen: count elements competing for "most distinctive" — target 1.

**Related:** [03-visual-design.md](03-visual-design.md#visual-hierarchy), [#8-aesthetic-and-minimalist-design](#8-aesthetic-and-minimalist-design).

**Sources:** [Laws of UX: Von Restorff Effect](https://lawsofux.com/von-restorff-effect/).

### Zeigarnik Effect

**Definition:** People remember uncompleted or interrupted tasks better than completed ones; open loops create cognitive tension.

**Reasoning / Evidence:** Bluma Zeigarnik (1927, waiters remembering unpaid orders). Replication record is mixed, but the applied pattern (visible incompleteness motivates completion) is widely observed in product metrics.

**When to use:**
- Profile-completion meters, setup checklists, progress-to-reward indicators.
- Resume-where-you-left-off affordances (continue watching/reading).
- Onboarding: leave a visible, easy next step open.

**When NOT to use / exceptions:**
- Manufactured incompleteness (fake progress, endless checklists) is manipulative and produces anxiety, not value ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#gamification-and-streaks)).
- Streaks that punish breaks exploit this effect against user wellbeing.

**Pros:** Ethical motivation for genuinely valuable completion (better profiles → better product experience).
**Cons:** Weak/contested effect size; dark-pattern adjacency.

**Implementation guidance:**
- Start progress bars non-zero when real prior steps exist (signup counts as step 1).
- Show remaining-step count and make each next step one click away.

**Related:** [Goal-Gradient Effect](#goal-gradient-effect), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#gamification-and-streaks).

**Sources:** [Laws of UX: Zeigarnik Effect](https://lawsofux.com/zeigarnik-effect/).

### Serial Position Effect

**Definition:** Items at the beginning (primacy) and end (recency) of a list are recalled best; the middle is forgotten.

**Reasoning / Evidence:** Ebbinghaus; formalized by Murdock (1962). Primacy = more rehearsal into long-term memory; recency = still in working memory.

**When to use:**
- Order nav items, menus, and lists so the most important sit first and last.
- Onboarding sequences: key message first, key action last.
- Bottom-nav apps: home and profile/cart at the two ends.

**When NOT to use / exceptions:**
- Applies to *memory* of sequences, less to visual scanning of persistent lists (where F-pattern and position dominate — [02-cognitive-foundations.md](02-cognitive-foundations.md#scanning-patterns)).
- Alphabetical or logical orders may rightfully override.

**Pros:** Zero-cost ordering heuristic.
**Cons:** Conflicts with frequency-of-use ordering; middle items aren't doomed if visible and labeled well.

**Implementation guidance:**
- Put the 1–2 most critical destinations at positions first and last; bury the least-used in the middle.
- In instructional content, repeat the critical instruction at the end.

**Related:** [08-navigation-ia.md](08-navigation-ia.md), [02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking).

**Sources:** [Laws of UX: Serial Position Effect](https://lawsofux.com/serial-position-effect/).

### Postel's Law

**Definition:** (Robustness principle) "Be conservative in what you send, liberal in what you accept" — accept varied, imperfect input; produce strict, predictable output.

**Reasoning / Evidence:** Jon Postel, TCP specification (RFC 761, 1980). Applied to UX: humans are variable; systems should absorb that variability rather than punish it ([Tesler's Law](#teslers-law) sibling).

**When to use:**
- Input parsing: accept phone numbers/dates/card numbers with or without spaces, dashes, mixed formats; normalize internally.
- Search: tolerate typos, synonyms, singular/plural.
- Paste handling: strip formatting, parse smartly.

**When NOT to use / exceptions:**
- Security contexts require strict validation — liberal acceptance of executable/markup input is an injection vector; be liberal in *format*, strict in *content* sanitization.
- Silent "correction" that changes meaning (auto-"fixing" an ID) is dangerous; confirm ambiguous normalizations.

**Pros:** Eliminates format-fight frustration, a top form-abandonment cause.
**Cons:** Parser complexity; ambiguity risk (is 02/03 Feb 3 or Mar 2? — depends on locale, must not guess silently).

**Implementation guidance:**
- Strip whitespace/dashes from numeric IDs before validating; accept both cases in codes.
- Show the normalized interpretation back to the user when ambiguity existed.
- Never reject input for formatting the machine can fix.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#input-masking-and-formatting), [Tesler's Law](#teslers-law).

**Sources:** [Laws of UX: Postel's Law](https://lawsofux.com/postels-law/).

### Pareto Principle

**Definition:** Roughly 80% of effects come from 20% of causes — 80% of usage concentrates on 20% of features.

**Reasoning / Evidence:** Vilfredo Pareto (1896, wealth distribution); Juran generalized. Feature-usage telemetry across software consistently shows heavy-head distributions.

**When to use:**
- Prioritize design/optimization effort on the top-usage flows.
- Progressive disclosure boundaries: the 20% visible, the tail behind "advanced."
- Performance work: optimize the screens 80% of sessions touch.

**When NOT to use / exceptions:**
- The rarely-used 80% may include contractually or safety-critical functions (emergency stop, GDPR export) — importance ≠ frequency.
- New features always start at 0% usage; don't kill them by Pareto alone.

**Pros:** Rational resource allocation; the argument for defaults-first design.
**Cons:** Misused to justify neglecting minority-but-vital needs, including accessibility.

**Implementation guidance:**
- Instrument feature usage before restructuring menus/toolbars; promote the measured head, demote the tail.
- Review the tail for "rare but critical" before demoting.

**Related:** [Hick's Law](#hicks-law), [02-cognitive-foundations.md](02-cognitive-foundations.md#progressive-disclosure).

**Sources:** [Laws of UX: Pareto Principle](https://lawsofux.com/pareto-principle/).

### Goal-Gradient Effect

**Definition:** Effort and motivation increase as people get closer to a goal; perceived progress accelerates completion.

**Reasoning / Evidence:** Hull (1932, rats); Kivetz, Urminsky & Zheng (2006, café loyalty program): purchases accelerate as customers approach the reward. The endowed-progress effect is the separate Nunes & Drèze study (2006, *Journal of Consumer Research*): a 10-stamp car-wash card with 2 pre-stamps outperformed an 8-stamp empty card despite identical required effort (8 purchases).

**When to use:**
- Progress bars in checkout/onboarding/setup; show completed steps prominently.
- Loyalty and milestone systems: seed initial progress honestly (signup = step 1 done).
- Long forms: "2 of 3 sections complete."

**When NOT to use / exceptions:**
- Fake endowed progress (crediting steps that don't exist) is deceptive.
- Very long journeys with visible distant goals demotivate early; chunk into near-term sub-goals.

**Pros:** Measurably improves completion rates late in flows.
**Cons:** Motivation dips right *after* goal completion ("post-reward reset") — plan re-engagement; manipulative variants shade into dark patterns.

**Implementation guidance:**
- Show progress as steps completed (not just remaining); keep step counts honest.
- Place the hardest/optional steps late, when gradient motivation is highest.

**Related:** [Zeigarnik Effect](#zeigarnik-effect), [07-forms-and-input.md](07-forms-and-input.md#multi-step-vs-single-page-forms).

**Sources:** [Laws of UX: Goal-Gradient Effect](https://lawsofux.com/goal-gradient-effect/), [Kivetz, Urminsky & Zheng 2006](https://doi.org/10.1509/jmkr.43.1.39), [Nunes & Drèze 2006](https://doi.org/10.1086/500480).

### Parkinson's Law

**Definition:** Work expands to fill the time available for its completion; users take as long as the perceived allowance.

**Reasoning / Evidence:** Cyril Parkinson (1955, essay). In UX: unconstrained tasks drift; deadlines and time expectations shape completion behavior.

**When to use:**
- Set honest time expectations ("takes about 2 minutes") — reduces abandonment and expands perceived generosity when you finish early.
- Timeboxing in productivity tools (focus timers, meeting defaults shorter than 60min).
- Autofill/smart defaults compress task time below users' expected allowance → delight.

**When NOT to use / exceptions:**
- Artificial countdown pressure (fake scarcity timers) is a dark pattern ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#fake-scarcity-fake-urgency-and-fake-social-proof)).
- Genuine deliberation tasks (consent, legal review) must not be time-pressured.

**Pros:** Time expectations are cheap to add and measurably reduce drop-off.
**Cons:** Underestimating stated times destroys trust; the "law" is an aphorism, not controlled science.

**Implementation guidance:**
- State expected duration at flow start; err on overestimating slightly (finish "early").
- For recurring tasks, show elapsed-time stats to help users self-timebox.

**Related:** [Goal-Gradient Effect](#goal-gradient-effect), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#fake-scarcity-fake-urgency-and-fake-social-proof).

**Sources:** [Laws of UX: Parkinson's Law](https://lawsofux.com/parkinsons-law/).

---

## Rule sets and standards

### Shneiderman's Eight Golden Rules

**Definition:** Ben Shneiderman's rules of interface design (*Designing the User Interface*, 1987–): (1) strive for consistency; (2) seek universal usability (novice–expert paths, accessibility); (3) offer informative feedback; (4) design dialogs to yield closure (clear start–middle–end); (5) prevent errors; (6) permit easy reversal of actions; (7) keep users in control; (8) reduce short-term memory load.

**Reasoning / Evidence:** Distilled from decades of HCI research and practice; heavily overlaps Nielsen's heuristics (independent convergence increases confidence in both sets). Distinct contributions: *closure* (#4) and *universal usability* (#2) are more explicit here.

**When to use:**
- As an alternative/complementary checklist for heuristic evaluation.
- "Yield closure": design flows with explicit completion signals (confirmation pages, success states) so users know they're done and mentally release the task.

**When NOT to use / exceptions:**
- Like all heuristic sets: guidance for evaluation, not a substitute for user testing.

**Pros:** Adds closure and universal-usability framing missing from other sets.
**Cons:** Redundant with Nielsen for most items; using both sets doubles checklist length for ~30% new coverage.

**Implementation guidance:**
- Every flow ends with an unambiguous completion state (what happened, what's next).
- Design novice path (guided, visible) and expert path (shortcuts, bulk) for each core task.

**Related:** [Nielsen's heuristics](#nielsens-10-usability-heuristics), [17-ux-process-research.md](17-ux-process-research.md#heuristic-evaluation).

**Sources:** [Shneiderman: Eight Golden Rules](https://www.cs.umd.edu/users/ben/goldenrules.html).

### ISO 9241-110 dialogue principles

**Definition:** The international ergonomics standard's seven interaction principles (2020 revision): suitability for the user's tasks; self-descriptiveness; conformity with user expectations; learnability; controllability; robustness against use errors; user engagement.

**Reasoning / Evidence:** ISO 9241 "Ergonomics of human-system interaction" — the standards-track counterpart to heuristics, used in regulated procurement (EU tenders, medical, automotive) and as the basis for usability requirements documents.

**When to use:**
- Regulated/procurement contexts requiring standards conformance claims.
- Writing testable usability requirements ("the system shall be self-descriptive: every screen identifies its purpose without external help").

**When NOT to use / exceptions:**
- For day-to-day design critique, Nielsen/Shneiderman phrasing is more actionable; ISO wording is abstract.

**Pros:** Standards authority; maps to EN 301 549 / regulatory ecosystems ([05-accessibility.md](05-accessibility.md#why-accessibility)).
**Cons:** Paywalled documents; abstraction requires interpretation work.

**Implementation guidance:**
- Map each ISO principle to concrete acceptance criteria per screen during requirements phase; cite the mapping in conformance documentation.

**Related:** [05-accessibility.md](05-accessibility.md), [17-ux-process-research.md](17-ux-process-research.md).

**Sources:** [ISO 9241-110:2020](https://www.iso.org/standard/75258.html).

---

## Quick reference

| Principle | One-line guidance | Key numbers |
|---|---|---|
| Visibility of system status | Always show what's happening, promptly | <100ms none, 1–10s spinner, >10s progress bar |
| Match system/real world | Users' words, users' order | — |
| User control & freedom | Undo, cancel, exit everywhere | ≥20 undo levels; 30-day trash |
| Consistency & standards | Follow platform + self conventions | — |
| Error prevention | Design out errors before messaging them | ≥8px separation for destructive actions |
| Recognition over recall | Show, don't make users remember | ~4 chunks working memory |
| Flexibility & efficiency | Novice path + expert accelerators | Standard shortcut set |
| Aesthetic & minimalist | Every element must earn its place | Cut copy ~50%; 1 primary CTA |
| Error recovery | Plain language: what/why/how to fix | — |
| Help & documentation | Task-based, in context, searchable | <300 words/task |
| Affordances/signifiers | Actionable must look actionable | >90% clickability recognition |
| Mapping | Control layout mirrors effect | Up/right = more (LTR) |
| Feedback | Acknowledge every action | <100ms; toast 4–8s |
| Constraints | Only offer valid actions; explain disabled | — |
| Proximity | Space within < space between groups | ≥1:1.5, ideal 1:2 |
| Similarity | Same look = same behavior | One accent color for interactive |
| Common region | Containers group strongest | ≤2 nesting levels; 16px card padding |
| Figure/ground | One dominant layer | Scrim 40–60% |
| Common fate | Related things move together | Stagger 20–50ms |
| Fitts's Law | Big, close targets for frequent actions | 44pt/48dp/24px targets |
| Hick's Law | Fewer, grouped choices | 5–7 top-level nav items |
| Jakob's Law | Default to the convention | Review 3–5 category leaders |
| Miller's Law | Chunk; don't tax memory | 3–4 safe chunks; 6-digit codes chunked |
| Tesler's Law | System absorbs complexity, not user | detect > default > defer > drop |
| Doherty Threshold | Keep the loop under 400ms | <400ms; input ack <100ms |
| Peak-End Rule | Fix the worst moment; polish the end | ~50ms first impression |
| Aesthetic-Usability | Beauty buys tolerance — and masks bugs | Weight behavior over ratings |
| Von Restorff | One thing pops per view | 1 isolate; ≥2 combined cues |
| Zeigarnik | Open loops drive completion | Non-zero progress start |
| Serial Position | First and last positions win | Ends of nav for critical items |
| Postel's Law | Accept liberally, output strictly | Normalize, then confirm ambiguity |
| Pareto | Optimize the top 20% of usage | Instrument before restructuring |
| Goal-Gradient | Visible progress accelerates finish | Endowed progress (2/10 > 0/8) |
| Parkinson's Law | Set honest time expectations | State duration up front |
| Shneiderman #4 | Every flow yields explicit closure | — |
| ISO 9241-110 | Standards framing for regulated work | 7 principles (2020) |
