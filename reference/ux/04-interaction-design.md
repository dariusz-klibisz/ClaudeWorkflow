# Interaction Design

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers the behavior layer: affordances in practice, component states, modes and mode errors, feedback timing and measurement, microinteractions, motion, loading strategies, gestures (touch, pointer, and pen), drag and drop, undo vs confirmation, direct manipulation, validation feedback, error taxonomy, hover problems, scrolling, keyboard interaction, empty states, and notifications. Perceptual/cognitive grounding is in [01-core-principles.md](01-core-principles.md) and [02-cognitive-foundations.md](02-cognitive-foundations.md); visual state styling in [03-visual-design.md](03-visual-design.md); form-specific validation in [07-forms-and-input.md](07-forms-and-input.md); keyboard accessibility requirements in [05-accessibility.md](05-accessibility.md); Core Web Vitals thresholds in [13-platform-web.md](13-platform-web.md#core-web-vitals); platform-specific input in [14-platform-desktop.md](14-platform-desktop.md) and [15-platform-mobile.md](15-platform-mobile.md).

## Contents

- [Foundations](#foundations)
  - [Affordances and signifiers in practice](#affordances-and-signifiers-in-practice)
  - [Direct manipulation](#direct-manipulation)
  - [Component states](#component-states)
  - [Modes and mode errors](#modes-and-mode-errors)
- [Feedback and time](#feedback-and-time)
  - [Feedback and response time thresholds](#feedback-and-response-time-thresholds)
  - [Microinteractions](#microinteractions)
  - [Motion and animation](#motion-and-animation)
  - [View transitions](#view-transitions)
- [Loading and waiting](#loading-and-waiting)
  - [Loading strategies](#loading-strategies)
  - [Optimistic UI](#optimistic-ui)
- [Input modalities](#input-modalities)
  - [Touch and pointer gestures](#touch-and-pointer-gestures)
  - [Drag and drop](#drag-and-drop)
  - [Hover dependence](#hover-dependence)
  - [Keyboard interaction patterns](#keyboard-interaction-patterns)
- [Safety and reversibility](#safety-and-reversibility)
  - [Undo/redo vs confirmation dialogs](#undoredo-vs-confirmation-dialogs)
  - [Validation feedback timing](#validation-feedback-timing)
  - [Error feedback and error taxonomy](#error-feedback-and-error-taxonomy)
- [Structure of movement](#structure-of-movement)
  - [Scroll behaviors](#scroll-behaviors)
  - [Empty states](#empty-states)
  - [Notifications and interruption design](#notifications-and-interruption-design)
- [Quick reference](#quick-reference)

---

## Foundations

### Affordances and signifiers in practice

**Definition:** The operational discipline of making every interactive element visibly interactive, every non-interactive element visibly inert, and every possible manipulation discoverable. (Theory: [01-core-principles.md](01-core-principles.md#affordances).)

**Reasoning / Evidence:** NN/g flat-design research measured click uncertainty from weak signifiers: users hover-hunt, hesitate, and miss features entirely. Interaction cost of ambiguity compounds across every session.

**When to use:**
- Audit each screen: (1) everything clickable looks clickable; (2) nothing unclickable looks clickable; (3) draggable/editable/expandable states are signposted.
- Cursor discipline (web/desktop): `pointer` for links/buttons, `text` for editable, `grab/grabbing` for draggable, `default` otherwise — cursor is a first-class signifier.

**When NOT to use / exceptions:**
- Expert-tool hidden accelerators (right-click menus, keyboard-only commands) may lack visual signifiers if a discoverable equivalent path exists.

**Pros:** Directly reduces time-on-task and missed-feature rates.
**Cons:** Signifier weight fights minimal aesthetics — resolve toward function.

**Implementation guidance:**
- Buttons: filled or bordered; links in text: color + underline; inputs: visible border ≥3:1 contrast.
- Hover feedback within 100ms on all interactive elements (see [#component-states](#component-states)).
- Editable-in-place text: pencil-on-hover is weak on touch; prefer persistent affordance or explicit edit mode.
- Test: static-screenshot clickability marking, >90% accuracy target.

**Related:** [01-core-principles.md](01-core-principles.md#signifiers), [03-visual-design.md](03-visual-design.md#flat-design-and-minimalism).

**Sources:** [NN/g: Clickability Signifiers](https://www.nngroup.com/articles/clickable-elements/), [NN/g: Flat Design](https://www.nngroup.com/articles/flat-design/).

### Direct manipulation

**Definition:** Interaction style where users act on visible objects with continuous representation, physical actions instead of command syntax, and rapid, incremental, reversible operations whose effects are immediately visible (Shneiderman, 1983).

**Reasoning / Evidence:** Shneiderman's formulation explains GUI/touch success: visibility + immediate feedback + reversibility minimize the gulfs of execution and evaluation (Norman). Drag-to-reorder, canvas editing, sliders, pinch-zoom are all direct manipulation.

**When to use:**
- Spatial/visual tasks: arranging, resizing, cropping, drawing, reordering.
- Continuous-value adjustment (sliders with live preview beat numeric entry for exploration; pair with numeric input for precision).
- Touch interfaces (fingers on objects are the paradigm's purest form).

**When NOT to use / exceptions:**
- Bulk/repetitive operations: manipulating 500 items directly is torture — provide command-style bulk actions, search-and-replace, multi-select.
- Precision requirements: dragging to exactly 42.5% fails; always pair manipulation with numeric/keyboard entry.
- Abstract operations with no natural object representation.

**Pros:** Low learning curve, error visibility, engagement; changes feel safe (continuous + reversible).
**Cons:** Doesn't scale to bulk; can be slower than commands for experts; accessibility requires parallel non-pointer paths (WCAG 2.5.7 — [#drag-and-drop](#drag-and-drop)).

**Implementation guidance:**
- Live preview during manipulation (ghost/target indicators while dragging, live value while sliding).
- Immediate commit + undo, rather than manipulate-then-apply, where operations are safe.
- Pair every manipulation with keyboard/numeric equivalent.
- The system must respond immediately and visibly to every manipulation: pressed, selected, and changed states must be visually distinct, and toggles must show current vs changed state unambiguously.

**Related:** [#drag-and-drop](#drag-and-drop), [01-core-principles.md](01-core-principles.md#3-user-control-and-freedom).

**Sources:** [Shneiderman 1983](https://doi.org/10.1109/MC.1983.1654471), [NN/g: Direct Manipulation](https://www.nngroup.com/articles/direct-manipulation/).

### Component states

**Definition:** The complete set of visual/behavioral conditions every interactive component must define: default, hover, focus, active/pressed, selected, disabled, loading, error, success, read-only, empty, skeleton. A component without designed states is an unfinished component.

**Reasoning / Evidence:** States are the feedback vocabulary of an interface ([Norman's feedback](01-core-principles.md#feedback)); missing states are a top source of "is it broken?" moments. Focus state is a WCAG requirement (2.4.7), not a nicety.

**When to use:**
- Every interactive component spec/design review: run the state checklist.
- Design systems: state matrix per component is the definition of done ([10-design-systems.md](10-design-systems.md#variants-and-states)).

**When NOT to use / exceptions:**
- Hover states don't exist on touch — never encode information exclusively in hover ([#hover-dependence](#hover-dependence)).
- Disabled states are frequently the wrong answer (see guidance).

**Pros:** Complete feedback loop; predictable, testable components.
**Cons:** Multiplies design/QA surface (states × variants × themes × densities).

**Implementation guidance:**
- Hover: perceptible within 100ms; subtle (background shift, elevation +1); never layout-shifting.
- Focus (keyboard): visible ring ≥3:1 contrast, ≥2px, offset from the element; never `outline: none` without replacement ([05-accessibility.md](05-accessibility.md#focus-visible)).
- Active/pressed: distinct from hover (darker/inset), <100ms response.
- Disabled: explain *why* and how to enable (tooltip/help text); prefer enabled-with-validation-error over silently disabled submit buttons ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#disabled-submit-without-explanation)); keep ≥3:1 contrast where possible; consider read-only (focusable, copyable) instead. Never use disabled styling for read-only information — they are different states with different semantics.
- Loading (per-component): spinner-in-button with disabled interaction + preserved width; label change ("Saving…").
- Error: color + icon + message, associated programmatically (`aria-describedby`).
- Selected: must survive theme changes and remain distinguishable from hover and focus simultaneously (a focused, hovered, selected row shows all three).

**Related:** [03-visual-design.md](03-visual-design.md), [05-accessibility.md](05-accessibility.md#focus-visible), [10-design-systems.md](10-design-systems.md#variants-and-states).

**Sources:** [WCAG 2.4.7](https://www.w3.org/WAI/WCAG22/Understanding/focus-visible.html), [Material: Interaction states](https://m3.material.io/foundations/interaction/states/overview).

### Modes and mode errors

**Definition:** A mode is a system state that changes what the same user action does: selection mode, edit mode, drawing mode, bulk-action mode, admin impersonation mode, caps lock, insert vs overwrite. A **mode error** occurs when the user acts believing they are in one mode while the system is in another. Design rule: make modes obvious, constrained, and easy to exit.

**Reasoning / Evidence:** Mode errors are a classic HCI failure category (Norman): the user's action is correct for their mental model but wrong for the system's actual state, so the error is invisible until damage is done. The cognitive mechanics are covered in [02-cognitive-foundations.md](02-cognitive-foundations.md#mode-errors). Modes are most dangerous when the mode indicator is subtle, peripheral, or absent — the user has no cue that the action mapping changed.

**When to use (modes legitimately):**
- Genuinely distinct activities on the same surface: view vs edit on a document, select-items vs act-on-items in a list, draw vs pan on a canvas.
- Bulk-action mode: entering a checkbox-selection state with a contextual action bar is clearer than hover-revealed per-row bulk controls.
- Admin impersonation ("viewing as user X"): a mode by definition — the same clicks now act on someone else's account.

**When NOT to use / exceptions:**
- Avoid modes when a modeless alternative exists at similar cost (e.g., direct inline editing instead of a global edit mode).
- Never place dangerous or irreversible actions inside ambiguous modes — if the user might not know which mode they're in, destructive actions must not be one click away.
- Quasimodes (spring-loaded modes held by continuous action, like holding Shift) are safer than latched modes because the user's own body maintains the state — prefer them for short-lived mode needs.

**Pros:** Modes let one surface support many activities without duplicating chrome; bulk modes scale actions efficiently.
**Cons:** Every mode is a standing invitation to mode errors; each mode multiplies the state space to design, test, and document.

**Implementation guidance:**
- Persistent, unmissable mode indicator: change the cursor, the toolbar, the header color, or add a banner — the signal must live where the user is looking, not in a distant status bar.
- Admin impersonation must show a persistent, high-contrast banner ("You are viewing as jane@example.com — Exit") for the entire session; it is both a mode signal and a safety rail.
- Provide a clear, always-visible exit control (Esc where conventional, plus a visible "Done"/"Exit" button); never require users to remember an exit gesture.
- Change available actions with the mode: hide or disable actions that don't apply, so the visible toolbar itself communicates the mode.
- On mode entry, make the transition perceptible (brief animation, chrome change) so accidental entry is noticed immediately.

**Verification checklist:**
- [ ] Every mode has a persistent visual indicator visible in the user's working area.
- [ ] Every mode has an obvious exit control (not gesture-only or memory-dependent).
- [ ] No destructive/irreversible action is reachable from an ambiguously indicated mode.
- [ ] Admin/impersonation states display a persistent banner for the full session.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#mode-errors), [01-core-principles.md](01-core-principles.md#1-visibility-of-system-status), [#component-states](#component-states).

**Sources:** [Norman, The Design of Everyday Things](https://mitpress.mit.edu/9780262525671/the-design-of-everyday-things/).

---

## Feedback and time

### Feedback and response time thresholds

**Definition:** The three canonical human time limits (Miller 1968; Nielsen 1993): **0.1s** — feels instantaneous, no feedback needed beyond the result; **1s** — flow limit, delay noticed but thought uninterrupted, no spinner needed but no other feedback either; **10s** — attention limit, beyond it users task-switch; determinate progress + cancel required. Plus the **Doherty threshold (~0.4s)** for productivity-critical loops.

**Reasoning / Evidence:** Robert Miller (1968) response-time taxonomy; Nielsen's synthesis has held for 30+ years because the limits are human constants, not hardware facts. Modern correlates: Google/CrUX INP threshold 200ms "good" ([13-platform-web.md](13-platform-web.md#core-web-vitals)).

**When to use:**
- Performance budgeting: classify every operation into a tier; design the feedback per tier *before* optimizing.
- Choosing indicator type: none (<1s) → indeterminate spinner (1–10s) → determinate progress + time estimate + cancel (>10s).

**When NOT to use / exceptions:**
- Sub-second spinners cause flash-of-spinner jank — delay indicator appearance by 300–500ms so fast responses never show one.
- Artificial delay for trust (security theater) — rare deliberate exception ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#fake-loading)).

**Pros:** Simple, evidence-backed decision rule covering all waiting UX.
**Cons:** Tiers assume attention on the task; background operations follow notification logic instead.

**Implementation guidance:**
- <100ms: input acknowledgment only (state change). Keystroke echo <50ms.
- 100ms–1s: show result when ready; cursor/subtle cue acceptable; no spinner.
- 1–10s: indeterminate indicator (spinner/skeleton), appearing after a 300ms grace delay.
- >10s: determinate progress bar, step description ("Uploading 3 of 12"), time estimate, cancel; allow backgrounding + notify on completion.
- Never freeze the UI thread; interaction must stay responsive during waits.

**Success feedback (what to show when the operation completes):**
- Match assurance to stakes: inline "Saved" status for autosave (quiet, ambient); toast for non-critical background success; **confirmation page or receipt for high-risk transactions** (payments, bookings, legal submissions) — a persistent, referenceable record, never an auto-dismissing message.
- Never put information the user must retain (order numbers, confirmation codes) in a disappearing success message.
- Avoid celebration for routine actions — reserve emphasis for milestones ([#microinteractions](#microinteractions)).

**Measurement methodology (how to know you're meeting the thresholds):**
- Distinguish **field data (RUM — real user monitoring)** from **lab data**: field data reflects what real users experience on real devices/networks and should drive production priorities; lab data (Lighthouse, synthetic tests) catches regressions before release but is not the experience itself.
- Evaluate field metrics at the **75th percentile**, segmented mobile vs desktop — averages hide the slow tail that drives abandonment.
- Lighthouse cannot directly measure real INP (it requires actual user interaction); use TBT (Total Blocking Time) as a lab proxy only, and confirm with field INP.
- Treat **performance budgets as acceptance criteria**: critical pages/flows get explicit field or lab budgets (e.g., "checkout LCP ≤2.5s at p75"), enforced in review/CI, not aspirational footnotes.
- Layout stability is a UX rule, not just a metric: **never insert content above the user's current focus** (late-loading banners, ads, images without reserved space) — it moves the thing they were about to tap. Reserve space for media and dynamic content; keep CLS ≤0.1. Threshold table in [13-platform-web.md](13-platform-web.md#core-web-vitals).

**Verification checklist:**
- [ ] Every operation classified into a response-time tier with designed feedback for that tier.
- [ ] Critical pages have explicit performance budgets (field or lab) treated as acceptance criteria.
- [ ] Metrics evaluated at the 75th percentile, mobile and desktop segmented.
- [ ] Layout shifts never move controls during interaction; space reserved for late-arriving content.
- [ ] Common interactions remain responsive under realistic device/network throttling.

**Related:** [01-core-principles.md](01-core-principles.md#doherty-threshold), [#loading-strategies](#loading-strategies), [13-platform-web.md](13-platform-web.md#core-web-vitals).

**Sources:** [NN/g: Response Times — 3 Important Limits](https://www.nngroup.com/articles/response-times-3-important-limits/), [NN/g: Website Response Times](https://www.nngroup.com/articles/website-response-times/), [Miller 1968](https://doi.org/10.1145/1476589.1476628), [web.dev: Core Web Vitals](https://web.dev/articles/vitals), [web.dev: Lab and field data differences](https://web.dev/articles/lab-and-field-data-differences).

### Microinteractions

**Definition:** Contained product moments accomplishing one small task, structured (Dan Saffer) as: **trigger** (user action or system event) → **rules** (what happens) → **feedback** (how you know) → **loops & modes** (what changes over time/repetition). Examples: like button, toggle, pull-to-refresh, password-strength meter, mute switch.

**Reasoning / Evidence:** Saffer, *Microinteractions* (2013); NN/g. The feel of a product is disproportionately carried by these small loops — they run thousands of times, so tiny costs/delights multiply.

**When to use:**
- Every state-changing control: design the full T-R-F-L loop, not just the icon.
- High-frequency actions (like, save, complete task): worth bespoke feedback polish.
- Communicating rules implicitly: strength meter teaches password rules better than a rule list.

**When NOT to use / exceptions:**
- Novelty animation on high-frequency actions ages fast — charm on 1st use, friction on 1000th; keep frequent-loop feedback fast and subtle, save flourish for rare milestones.
- Loops must respect reduced motion ([#motion-and-animation](#motion-and-animation)).

**Pros:** Product feel and perceived quality; implicit teaching; emotional connection at low cost.
**Cons:** Polish time; inconsistent microinteractions (each feature hand-rolling toggles) fragment the feel — systematize in the design system.

**Implementation guidance:**
- Feedback within 100ms of trigger; total loop ≤300ms for frequent actions.
- Long-term loops: consider repetition (first-time tooltip disappears after N uses; celebration animation only on milestones).
- Document T-R-F-L per component in the design system ([10-design-systems.md](10-design-systems.md)).

**Related:** [#motion-and-animation](#motion-and-animation), [#component-states](#component-states).

**Sources:** [NN/g: Microinteractions](https://www.nngroup.com/articles/microinteractions/), [Saffer, Microinteractions](https://www.oreilly.com/library/view/microinteractions/9781449342760/).

### Motion and animation

**Definition:** Purposeful interface motion serving four functions: **orientation** (where did that come from/go), **causality** (this action caused that result), **feedback** (state confirmation), and **attention** (look here). Everything else is decoration — budget it separately and sparingly.

**Reasoning / Evidence:** Change blindness ([02-cognitive-foundations.md](02-cognitive-foundations.md#change-blindness)) motivates transitional motion; spatial model maintenance (where panels live) depends on consistent movement. Costs are real: vestibular disorders make large motion physically sickening (WCAG 2.3.3), and time spent animating is time the user waits.

**When to use:**
- State transitions that would otherwise teleport (panel open/close, item added/removed from list, reorder).
- Micro-feedback (button press, toggle slide).
- Attention direction for remote changes (highlight fade on updated value).

**When NOT to use / exceptions:**
- `prefers-reduced-motion`: replace movement with opacity cross-fades; never ignore it ([05-accessibility.md](05-accessibility.md#reduced-motion)).
- Frequent expert loops: motion that delays interaction (400ms panel slide on an action done 200×/day) is a tax — keep functional motion ≤200–300ms or instant for power surfaces.
- Autoplaying/looping decorative motion near content (distraction; WCAG 2.2.2 pause requirement).

**Pros:** Spatial comprehension, causality clarity, perceived quality.
**Cons:** Waiting cost at scale; accessibility risk; performance (jank destroys the benefit — 60fps or don't).

**Implementation guidance:**
- Durations: micro-feedback 100–150ms; small transitions (dropdowns, toggles) 150–200ms; medium (panels, dialogs) 200–300ms; large/full-screen 300–400ms; nothing functional >500ms.
- Easing: ease-out (fast start, gentle stop) for entering elements; ease-in for exiting; standard curve (ease-in-out) for on-screen movement; never linear for spatial motion (feels mechanical).
- Choreography: related elements move together ([common fate](01-core-principles.md#common-fate)); stagger groups by 20–50ms; total sequence still ≤500ms.
- Performance: animate only `transform` and `opacity` (compositor-friendly); avoid layout-triggering properties (width/height/top/left).
- Enter vs exit asymmetry: exits faster (150ms) than entrances (250ms) — users wait for what arrives, not what leaves.

**Related:** [01-core-principles.md](01-core-principles.md#common-fate), [05-accessibility.md](05-accessibility.md#reduced-motion), [#view-transitions](#view-transitions).

**Sources:** [NN/g: Animation Purpose](https://www.nngroup.com/articles/animation-purpose-ux/), [Material: Motion](https://m3.material.io/styles/motion/overview), [web.dev: prefers-reduced-motion](https://web.dev/articles/prefers-reduced-motion).

### View transitions

**Definition:** Motion connecting whole-view changes: shared-element transitions (the tapped thumbnail grows into the detail hero), directional slides encoding hierarchy (drill-in slides left, back slides right), cross-fades for lateral moves.

**Reasoning / Evidence:** Maintains the user's spatial model across navigation — the answer to "where am I and how did I get here" encoded physically. Platform APIs now support it natively (View Transitions API on web, shared element transitions on Android/iOS).

**When to use:**
- Master→detail navigation (shared element); hierarchical drill-in/out (directional consistency); tab switches (subtle cross-fade or slide matching tab order).

**When NOT to use / exceptions:**
- Direction must be consistent: if drill-in slides left, back MUST slide right — inconsistent direction is worse than none.
- Reduced motion → cross-fade fallback.
- Slow transitions on frequent navigation are a tax; ≤300ms.

**Pros:** Spatial continuity, perceived speed (motion masks load), polish.
**Cons:** Implementation complexity; breaks if content isn't ready (janky half-transitions) — coordinate with data loading.

**Implementation guidance:**
- Web: View Transitions API (`document.startViewTransition`) with CSS `view-transition-name` for shared elements; feature-detect and degrade to instant swap.
- Encode hierarchy: forward = new content from right (LTR); backward = reverse; modal = from bottom/scale; sibling = fade.
- Keep 200–300ms; run skeletons beneath if data lags.

**Related:** [#motion-and-animation](#motion-and-animation), [08-navigation-ia.md](08-navigation-ia.md), [13-platform-web.md](13-platform-web.md#spa-vs-mpa).

**Sources:** [MDN: View Transitions API](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API).

---

## Loading and waiting

### Loading strategies

**Definition:** The choice among: nothing (fast ops), spinner (short indeterminate), skeleton screen (layout-shaped placeholder), determinate progress bar (long/known), and progressive/streaming render (show what's ready).

**Reasoning / Evidence:** Perceived duration ≠ actual duration: skeletons shift perception by promising structure and appearing busy; progress bars with accelerating fill are perceived faster (Harrison et al. 2007); early visible content anchors the "it's working" judgment (why LCP matters). Spinner-staring makes waits feel *longer* than blank for very short waits.

**When to use:**
- <1s expected: nothing (with 300–500ms grace before any indicator appears).
- 1–10s, unknown duration: skeleton for content-shaped loads (feeds, cards, tables); spinner for small in-component loads (button, dropdown).
- >10s or known progress: determinate bar + step text + estimate + cancel.
- Any duration: progressive render beats any placeholder — stream header first, data as it arrives.

**When NOT to use / exceptions:**
- Skeletons that mismatch the real layout cause worse jank than spinners (CLS at resolve); only skeleton what you can shape-match, and never let a skeleton misrepresent the final content density.
- Multiple concurrent skeletons/spinners on one screen read as broken — coordinate to one loading boundary per region.
- Never skeleton-then-spinner-then-content chains; one strategy per load.
- Never full-page blocking spinners when only one region is loading — indicators should be local to the affected region, with the rest of the UI still usable.

**Pros (skeletons):** Perceived speed, layout stability, structure preview.
**Cons (skeletons):** Maintenance (drift from real layout), false promise if content differs, engineering cost vs a spinner.

**Implementation guidance:**
- Grace delay: show no indicator for the first 300ms; if shown, keep it ≥500ms (no flash).
- Skeletons: match final layout dimensions exactly (aspect-ratio reserved), subtle shimmer 1–1.5s cycle, ≤3 distinct shapes.
- Progress bars: never stall at 99%; move continuously (front-load perceived progress); show step labels for multi-stage.
- Long operations: allow backgrounding; notify completion ([#notifications-and-interruption-design](#notifications-and-interruption-design)).
- Keep controls responsive during background work where safe; explain long-running operations (what's happening, what's affected, whether the user can continue).

**Related:** [#feedback-and-response-time-thresholds](#feedback-and-response-time-thresholds), [13-platform-web.md](13-platform-web.md#core-web-vitals).

**Sources:** [NN/g: Skeleton Screens](https://www.nngroup.com/articles/skeleton-screens/), [NN/g: Progress Indicators](https://www.nngroup.com/articles/progress-indicators/), [Harrison et al. 2007](https://doi.org/10.1145/1240624.1240743).

### Optimistic UI

**Definition:** Reflecting an action's expected result immediately, before server confirmation — the like fills, the message appears, the item moves — reconciling (or rolling back with notice) if the request fails.

**Reasoning / Evidence:** Perceived latency drops to ~0 for the 95–99% of requests that succeed; pioneered at scale by social apps (Instagram likes). The trade: rollback complexity and rare user-visible corrections. (For the underlying human time limits that make instant-feel valuable, see the response-time evidence base in [#feedback-and-response-time-thresholds](#feedback-and-response-time-thresholds).)

**When to use:**
- High-success-rate (>95%), low-stakes, reversible mutations: likes, toggles, reorderings, message sends (with visible pending → sent → failed states).
- Local-first architectures where sync is background by design.

**When NOT to use / exceptions:**
- **Never for payments or irreversible operations** — no exceptions — nor for anything the user might act upon downstream before confirmation (booking a seat that wasn't actually free). The same prohibition covers legal, medical, permission-sensitive, and externally visible actions where a rollback would itself be a visible event.
- Low-reliability networks with high failure rates: constant rollbacks destroy trust faster than honest waiting.
- Server-computed outcomes (search results, validations) can't be optimistic — nothing to predict.

**Pros:** Instant feel; flow preservation ([02-cognitive-foundations.md](02-cognitive-foundations.md#flow-state)).
**Cons:** Rollback UX is hard (user has scrolled away — where does the error go?); state reconciliation complexity; can mask real system health issues.

**Implementation guidance:**
- Show subtle pending state where feasible (message clock icon → check) so "optimistic" isn't "dishonest."
- Failure: revert visibly + toast with retry ("Couldn't save — Retry"); never silently drop.
- Queue offline mutations; badge unsynced state ([15-platform-mobile.md](15-platform-mobile.md#offline-and-poor-connectivity)).

**Decision questions:**
- Is the action low-risk, rare-to-fail, and clearly rollback-able? If any answer is no, wait for the server.
- Is the action financial, legal, destructive, public, medical, or permission-sensitive? If yes, never optimistic.

**Related:** [#loading-strategies](#loading-strategies), [#feedback-and-response-time-thresholds](#feedback-and-response-time-thresholds).

**Sources:** [NN/g: Website Response Times](https://www.nngroup.com/articles/website-response-times/) — response-time evidence base (why perceived latency matters), not an optimistic-UI-specific study.

---

## Input modalities

### Touch and pointer gestures

**Definition:** The standard gesture vocabulary — tap (activate), double-tap (zoom), long-press (context/secondary), swipe (scroll/dismiss/reveal actions), pinch (zoom), two-finger rotate — and the discipline of using it without inventing undiscoverable custom gestures. Also covers pen/stylus as a distinct pointer class.

**Reasoning / Evidence:** Gestures have zero visual signifier: discoverability is entirely conventional or taught. NN/g: custom gestures fail silently — users never find them. The standard set works because platform-wide repetition taught it.

**When to use:**
- Standard gestures for their standard meanings — never repurpose (swipe-left on a list item = actions/delete convention).
- Gestures as accelerators duplicating visible controls (swipe-to-archive AND a visible archive button) — the redundancy rule.
- Pull-to-refresh in feed contexts (established convention).

**When NOT to use / exceptions:**
- Never gesture-only for any important function (discoverability + motor accessibility, WCAG 2.5.1: multipoint/path-based gestures need single-pointer alternatives).
- Edge swipes conflict with system gestures (iOS back-swipe, Android gesture nav) — avoid app gestures at screen edges ([15-platform-mobile.md](15-platform-mobile.md#gestures)).
- Hover-like long-press hides functionality on desktop-migrated designs.

**Pros:** Speed, screen-space economy, fluent expert feel.
**Cons:** Invisible, undiscoverable, motor-demanding, conflict-prone with system gestures.

**Implementation guidance:**
- Touch targets ≥44×44pt / 48×48dp; gesture zones generous ([15-platform-mobile.md](15-platform-mobile.md#touch-targets)).
- Immediate visual tracking during gesture (content follows finger 1:1; rubber-band at limits).
- Teach once: first-use hint (subtle animation of the swipe) then never again.
- Every gesture action has a visible equivalent; WCAG 2.5.2 pointer cancellation: trigger on release, not press, so users can slide off to cancel.
- **Design for the input methods users actually have, not the ones you assume.** Use the `pointer: coarse` / `pointer: fine` (and `any-pointer`) media queries to adapt: coarse pointers (finger) get larger targets and persistent affordances; fine pointers (mouse, trackpad) may get denser layouts and hover accelerators. A device can have both (touchscreen laptop) — `any-pointer` catches that.
- Pen/stylus input: treat as a fine pointer with extra channels (pressure, tilt, hover-before-contact on some hardware); don't gate functionality on those channels — they're enhancements. Palm rejection and generous drag thresholds matter for pen-first surfaces (drawing, annotation, signatures).

**Related:** [15-platform-mobile.md](15-platform-mobile.md#gestures), [05-accessibility.md](05-accessibility.md#touchpointer-accessibility), [#direct-manipulation](#direct-manipulation).

**Sources:** [NN/g: Touch Gestures](https://www.nngroup.com/articles/gestures/), [WCAG 2.5.1](https://www.w3.org/WAI/WCAG22/Understanding/pointer-gestures.html), [MDN: pointer media feature](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/pointer).

### Drag and drop

**Definition:** Direct-manipulation transfer/reordering: press, move with continuous visual tracking, release on a valid target. Requires signifying draggability, showing drag state, marking valid targets, and providing non-drag alternatives.

**Reasoning / Evidence:** Powerful for spatial tasks (reordering, categorizing, file movement) but among the least accessible interactions: motor-impaired users, touch+scroll conflicts, and screen readers all struggle — hence WCAG 2.5.7 (Dragging Movements, AA): every drag operation needs a single-pointer non-drag alternative.

**When to use:**
- Reordering lists/cards (kanban), moving items between containers, file upload (drop zones), spatial arrangement (canvas).

**When NOT to use / exceptions:**
- Never the *only* path: pair with cut/paste, move-to menus, arrow-key reordering, or order-number editing.
- Long-distance drags (across scroll regions/screens) are error-prone — offer select-then-move-to for those.
- Touch lists: drag conflicts with scroll — use explicit drag handles, not whole-row drag.

**Pros:** Natural for spatial mental models; fast for short-distance rearrangement.
**Cons:** Accessibility burden; discoverability (draggability is invisible without handles); fiddly precision at drop boundaries.

**Implementation guidance:**
- Signify: drag handle (⠿ six-dot) on hover/always on touch; cursor `grab`→`grabbing`.
- During drag: lifted appearance (shadow/scale 1.02–1.05), original position ghost, drop target highlight + insertion line; auto-scroll near edges (30–60px activation zone).
- On drop: settle animation ≤200ms; on invalid: spring back with brief reason.
- Keyboard path: focus item → Space/Enter picks up → arrows move → Space drops, Esc cancels — with `aria-live` announcements throughout.

**Related:** [#direct-manipulation](#direct-manipulation), [05-accessibility.md](05-accessibility.md#dragging-alternatives), [#touch-and-pointer-gestures](#touch-and-pointer-gestures).

**Sources:** [WCAG 2.5.7](https://www.w3.org/WAI/WCAG22/Understanding/dragging-movements.html), [NN/g: Drag and Drop](https://www.nngroup.com/articles/drag-drop/).

### Hover dependence

**Definition:** The anti-pattern family of information or actions available only on mouse hover: hover-revealed buttons, hover tooltips carrying essential info, hover menus.

**Reasoning / Evidence:** Touch devices have no hover; keyboard users don't hover; even mouse users can't discover what's invisible until pointed at. Hover-revealed row actions measurably reduce feature discovery. WCAG 1.4.13 constrains hover-triggered content (dismissable, hoverable, persistent).

**When to use (hover legitimately):**
- Enhancement-only: hover states confirming interactivity; shortcut hints in tooltips duplicating visible info; preview-on-hover as accelerator with click-through as the real path.
- Progressive density: hover-revealed *secondary* actions acceptable IF also reachable via visible overflow menu, keyboard focus reveal, and touch tap-reveal.

**When NOT to use / exceptions:**
- Essential info in hover tooltips only (truncation revealing critical text — restructure instead).
- Hover-only navigation menus without click/tap equivalents.
- Anything hover-revealed must appear on keyboard focus too.

**Pros (restrained hover reveal):** Visual density for scanning; less chrome per row.
**Cons:** Discovery loss; platform inequality; interaction asymmetry bugs (focus path forgotten).

**Implementation guidance:**
- Rule: hover may *accelerate*, never *gate*. Everything hover-reachable must be focus-reachable and touch-reachable.
- Row actions: visible overflow (⋯) always; individual icons may fade in on hover/focus as shortcut.
- Tooltips: 300–500ms open delay, no delay when moving between adjacent tooltips, dismiss on Esc (WCAG 1.4.13), never put actions inside plain tooltips (use popovers/toggletips for interactive content).

**Related:** [#component-states](#component-states), [05-accessibility.md](05-accessibility.md), [15-platform-mobile.md](15-platform-mobile.md).

**Sources:** [WCAG 1.4.13](https://www.w3.org/WAI/WCAG22/Understanding/content-on-hover-or-focus.html), [NN/g: Tooltip Guidelines](https://www.nngroup.com/articles/tooltip-guidelines/).

### Keyboard interaction patterns

**Definition:** Complete keyboard operability: logical Tab order across the page, arrow-key navigation *within* composite widgets (menus, tabs, grids, listboxes — the roving-tabindex pattern), Enter/Space activation, Escape dismissal, and focus management across dynamic changes.

**Reasoning / Evidence:** WCAG 2.1.1 (all functionality keyboard-operable) plus power-user efficiency ([01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use)). The ARIA Authoring Practices Guide defines expected key maps per widget — deviation breaks both AT users and muscle memory.

**When to use:**
- Every interactive surface. Composite widgets: one Tab stop for the whole widget, arrows move inside (a 50-item menu must not be 50 Tab stops).
- Dynamic UI: focus moves to opened dialogs (trap inside, return on close), to revealed content when contextually right, and gets announced route changes in SPAs.

**When NOT to use / exceptions:**
- Don't trap focus outside modal contexts.
- Custom shortcuts: avoid single-character global shortcuts (WCAG 2.1.4 — must be remappable/disableable); never override browser/OS standards ([14-platform-desktop.md](14-platform-desktop.md#keyboard-shortcuts-design)).

**Pros:** Accessibility compliance + expert speed with one investment.
**Cons:** Custom widgets require implementing full APG key maps — significant work (argument for using proven component libraries).

**Implementation guidance:**
- Tab order = visual order (LTR/RTL aware); DOM order drives it — avoid positive `tabindex`.
- Standard keys: Enter activates buttons/links; Space toggles/activates buttons and checkboxes; Esc closes topmost layer; arrows navigate within widgets; Home/End to extremes; type-ahead in long lists.
- Modals: focus first meaningful element on open, trap Tab, return focus to trigger on close.
- Skip link as first tab stop on content-heavy pages ([05-accessibility.md](05-accessibility.md#keyboard-navigation)).
- Prefer native controls where possible — they carry the key maps for free; shortcuts accelerate frequent actions but never replace visible paths; document custom shortcuts and expose commands through menus as well (desktop).

**Related:** [05-accessibility.md](05-accessibility.md#keyboard-navigation), [14-platform-desktop.md](14-platform-desktop.md#keyboard-shortcuts-design).

**Sources:** [ARIA Authoring Practices Guide](https://www.w3.org/WAI/ARIA/apg/patterns/), [WCAG 2.1.1](https://www.w3.org/WAI/WCAG22/Understanding/keyboard.html).

---

## Safety and reversibility

### Undo/redo vs confirmation dialogs

**Definition:** Two protection models for consequential actions: **confirm-before** (dialog interrupting to ask "are you sure?") vs **act-then-undo** (do it immediately, offer reversal). Modern default: undo for anything reversible; confirmation only for the irreversible.

**Reasoning / Evidence:** Habituation kills confirmations: routine dialogs are clicked through on autopilot, offering zero protection exactly when needed ([02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)). Undo protects *after* the user realizes the error — which is when errors are actually noticed. Gmail's undo-send is canonical.

**When to use (undo):**
- Deletions (soft-delete + trash), moves, bulk edits, sends with grace window, settings changes.
- Toast with undo (4–10s) + persistent recovery (trash, version history) for the slower realization.

**When to use (confirmation):**
- Truly irreversible + high cost: permanent purge, payment execution, cascade deletions affecting others — plus costly, public, or privacy-affecting actions where a mistake is visible or expensive even if technically reversible.
- Escalate friction to match stakes: consequence-stating dialog → type-to-confirm (resource name) for catastrophic scope.

**When NOT to use / exceptions:**
- Never confirm routine actions (habituation); never *both* confirm and undo for the same trivial action.
- Undo needs engineering honesty: if "undo" can't fully restore (external side effects), don't offer it — confirm instead.

**Pros (undo):** Non-interruptive, protects real errors, encourages exploration.
**Cons (undo):** Infrastructure cost (soft deletes, command stacks, sync semantics); grace-window sends delay actual delivery.

**Implementation guidance:**
- Undo toast: 5–10s visible, keyboard accessible, plus menu-level Undo (Ctrl/Cmd+Z) where a session model exists.
- Soft-delete retention: 30 days convention; show trash location in the toast ("Moved to Trash — Undo").
- Confirmation dialogs (when justified): title states the action ("Delete project 'Atlas'?"), body states consequence concretely ("This permanently deletes 1,204 records"), confirm button restates the verb ("Delete project" — never "OK"/"Yes"), destructive styling, safe button focused by default. Avoid generic "Are you sure?" dialogs — they state no consequence and train click-through.
- Warnings should appear *before* consequence, not after.
- Type-to-confirm reserved for irreversible + large blast radius.

**Related:** [01-core-principles.md](01-core-principles.md#3-user-control-and-freedom), [02-cognitive-foundations.md](02-cognitive-foundations.md#habituation), [09-content-ux-writing.md](09-content-ux-writing.md#confirmation-and-destructive-action-copy).

**Sources:** [NN/g: Confirmation Dialogs](https://www.nngroup.com/articles/confirmation-dialog/), [NN/g: Reset/Cancel Buttons](https://www.nngroup.com/articles/reset-and-cancel-buttons/).

### Validation feedback timing

**Definition:** When to tell users their input is invalid: on keystroke (immediate), on blur (leaving the field), on submit (batch). Best practice: **validate on blur, reward early / punish late** — show success promptly, but don't flag errors while the user is still typing.

**Reasoning / Evidence:** Keystroke-time error flags punish half-typed input ("invalid email" after one character) — measurably raising frustration and correction time. On-blur catches errors at the natural checkpoint. Research (Baymard, GOV.UK) supports hybrid: instant *positive* confirmation, deferred *negative*.

**When to use:**
- Field-level: validate on blur; once a field has erred, re-validate on every keystroke so the fix is confirmed instantly (punish-late-reward-immediately-on-fix).
- Async checks (username availability): debounced 300–500ms after typing pauses, spinner in-field.
- Submit-level: full validation with error summary regardless (never rely on client checks alone).

**When NOT to use / exceptions:**
- Hard-constraint inputs (max length, numeric-only) prevent rather than flag — constrain at entry ([01-core-principles.md](01-core-principles.md#constraints)).
- Cross-field validations (end date > start date) run when the *pair* completes, not on first field's blur.
- Never validate untouched empty fields (page-load error walls).

**Pros:** Errors caught near cause; flow preserved; fixes confirmed instantly.
**Cons:** More states to build than submit-only; premature blur (tabbing through) can still flag — only validate blurred fields the user actually entered.

**Implementation guidance:**
- Full pattern details in [07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing); message writing in [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing).
- Error rendering: border + icon + message below field, `aria-describedby` linkage, focus moves to first error on submit.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#inline-validation-timing), [01-core-principles.md](01-core-principles.md#5-error-prevention), [#error-feedback-and-error-taxonomy](#error-feedback-and-error-taxonomy).

**Sources:** [NN/g: Errors in Forms](https://www.nngroup.com/articles/errors-forms-design-guidelines/), [GOV.UK: Validation pattern](https://design-system.service.gov.uk/patterns/validation/).

### Error feedback and error taxonomy

**Definition:** The discipline of classifying errors by cause and handling each class differently. The three core classes: **field validation errors** (the user's input doesn't meet a rule — fixable by the user, shown at the field), **server/system failures** (the request failed for reasons unrelated to the input — often transient, needs retry), and **permission errors** (the user isn't allowed — needs explanation of who can, or how to request access). Errors must be noticeable, specific, recoverable, and never destructive to user work.

**Reasoning / Evidence:** A single generic error treatment ("Something went wrong") forces users to guess which of three very different situations they're in — and each has a different correct next step (fix input / retry / request access). Misclassified handling is common and corrosive: showing a red field border for a server timeout sends the user hunting for an input mistake that doesn't exist. Losing user input on error is among the most-hated failures in usability studies — retyping is pure punishment for a system fault.

**When to use:**
- Every error path in every flow: classify before designing the message. If engineering can't distinguish classes at the API boundary, that's the first fix.
- Field validation → inline at the field ([#validation-feedback-timing](#validation-feedback-timing)) plus an **error summary at the top of the form** on submit, with links that move focus to each offending field (GOV.UK pattern; details in [07-forms-and-input.md](07-forms-and-input.md)).
- Transient server failures → non-blaming message + **Retry** action; automatic retry with backoff for idempotent reads.
- Permission errors → state what's not allowed, why (if known), and the path forward (who to ask, how to request access); never render a permission wall as a generic failure.

**When NOT to use / exceptions:**
- Don't expose raw technical detail (stack traces, error codes alone) as the primary message; codes may append for support reference.
- Don't retry non-idempotent operations automatically (double-charge risk).
- Don't use error styling for warnings or info — severity levels are part of the taxonomy too.

**Pros:** Users always know their next step; support burden drops; input preservation removes the worst failure mode.
**Cons:** Requires API/error-contract discipline across the stack; more message variants to write and localize.

**Implementation guidance:**
- **Preserve input, always.** On any error — validation, server, session expiry — the user's entered data survives: keep form state client-side, restore after re-auth, autosave drafts for long input.
- Explain next steps in every error message: what happened, why (if known), what to do now ([09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing)).
- Field errors near fields; form-level summary at top on submit, focus moved to it, each item a link to its field.
- Transient failures: include a Retry control; distinguish "your connection" from "our server" when detectable (offline banner vs error state).
- `role="alert"`/`aria-live="assertive"` for errors that must interrupt; ensure error states are perceivable without color alone (icon + text).

**Verification checklist:**
- [ ] Validation, server-failure, and permission errors render distinct treatments with distinct next steps.
- [ ] No error path discards user input (test: fill a long form, kill the network, submit).
- [ ] Form submits with multiple errors show a summary at top with focus management and per-field links.
- [ ] Transient system failures offer Retry; permission errors offer an access path.

**Related:** [07-forms-and-input.md](07-forms-and-input.md), [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing), [#validation-feedback-timing](#validation-feedback-timing), [#empty-states](#empty-states) (error-type empty state).

**Sources:** [NN/g: Error Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/), [GOV.UK: Error summary component](https://design-system.service.gov.uk/components/error-summary/).

---

## Structure of movement

### Scroll behaviors

**Definition:** How content extends beyond the viewport and how the interface behaves during scroll: natural scrolling, sticky elements, scroll restoration, infinite vs paginated content, and the prohibition of scrolljacking.

**Reasoning / Evidence:** Scrolling is the most practiced interaction in existence; interference with its physics (scrolljacking: hijacked speed/direction/snapping) reliably tests poorly and triggers strong dislike (NN/g "Scrolljacking 101"). Users scroll comfortably when content signals continuation ([closure](01-core-principles.md#closure)).

**When to use:**
- Sticky headers/toolbars: keep orientation and actions available on long content — minimized height while scrolled (≤64px), optionally hide-on-scroll-down/show-on-scroll-up for reading surfaces.
- Scroll restoration: back navigation returns to exact prior position (browser default — don't break it in SPAs; [13-platform-web.md](13-platform-web.md#browser-conventions)).
- Scroll-snap for carousels/media (CSS scroll-snap: assistive, interruptible), not for page sections.
- Content choice: pagination vs infinite vs load-more — decision table in [08-navigation-ia.md](08-navigation-ia.md#pagination-vs-infinite-scroll-vs-load-more).

**When NOT to use / exceptions:**
- Scrolljacking (overriding scroll speed/direction, forced full-page section snapping, scroll-triggered takeover animations): avoid; the rare marketing exceptions must keep scroll interruptible and honor reduced-motion.
- Nested scroll regions confuse (scroll-within-scroll traps); minimize; when required, visually bound the inner region clearly.
- Fixed elements consuming >20–25% of small viewports.

**Pros (sticky/restoration):** Orientation, action availability, back-navigation integrity.
**Cons:** Sticky chrome costs viewport (esp. mobile); scroll-linked effects cost performance (jank).

**Implementation guidance:**
- Never `overflow: hidden` on body for layout hacks that break scroll; keep native momentum.
- Signal continuation: cut content at fold edges; avoid full-height sections that look like page ends ("false floor").
- Scroll-linked animations: passive listeners / CSS-driven; test at 6× CPU throttle.
- Preserve scroll position across: back nav, tab switches, list-detail-return, refresh where feasible.

**Related:** [08-navigation-ia.md](08-navigation-ia.md#pagination-vs-infinite-scroll-vs-load-more), [13-platform-web.md](13-platform-web.md#scroll-ux), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#custom-scrollbars-and-scrolljacking).

**Sources:** [NN/g: Scrolljacking 101](https://www.nngroup.com/articles/scrolljacking-101/), [NN/g: Scrolling and Attention](https://www.nngroup.com/articles/scrolling-and-attention/).

### Empty states

**Definition:** What a container shows when it has nothing: first-use (never had content), cleared (user emptied it), no-results (filter/search matched nothing), error (couldn't load). Each is a distinct design with distinct jobs.

**Reasoning / Evidence:** Empty states are disproportionately common in real usage (new users see *only* empty states) yet designed last. First-use empties are onboarding real estate; no-results empties are recovery junctions ([08-navigation-ia.md](08-navigation-ia.md#no-results-recovery)).

**When to use:**
- First-use: explain what will live here + primary CTA to create/import/connect (+ optional example/demo content).
- Cleared: confirm the state is intentional ("All caught up") — celebration beats blankness for task-completion empties.
- No-results: restate the query, offer corrections (spelling, filter loosening, clear-all), suggest alternatives.
- Error: what failed + retry + fallback path ([09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing)).

**When NOT to use / exceptions:**
- Don't confuse types: showing first-use marketing copy for a no-results filter state misleads.
- Illustration without action is decoration — every empty state needs the *next step*.

**Pros:** Converts dead ends into onboarding/recovery; prevents "is it broken?" confusion.
**Cons:** Four+ variants per container to design/localize/test.

**Implementation guidance:**
- Structure: (optional small illustration) + one-line what-this-is + one-line how-to-fill + primary action button.
- No-results: keep filters visible and editable; "Clear all filters" one click away; log queries hitting empty ([17-ux-process-research.md](17-ux-process-research.md#analytics-and-behavioral-data)).
- Loading ≠ empty: never flash empty state while data is in flight (skeleton first — [#loading-strategies](#loading-strategies)).

**Related:** [09-content-ux-writing.md](09-content-ux-writing.md#empty-states), [08-navigation-ia.md](08-navigation-ia.md#no-results-recovery).

**Sources:** [NN/g: Empty States](https://www.nngroup.com/articles/empty-state-interface-design/).

### Notifications and interruption design

**Definition:** The system's attention-requesting vocabulary, by escalating interruption cost: passive badge/dot → inline banner → toast/snackbar (transient) → modal (blocking) → OS push notification (out-of-app). Matching channel to urgency is the core skill.

**Reasoning / Evidence:** Interruption cost research ([02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost)): unnecessary interrupts destroy task performance; notification fatigue causes global dismissal/disabling — over-notifying is self-defeating (push opt-out is permanent-ish). NN/g: modals for genuinely blocking decisions only.

**When to use:**
- Badge: awaiting-attention counts, no urgency (unread counts).
- Inline banner: persistent context-level state (offline, unsaved, degraded service) — stays until resolved.
- Toast: transient outcome confirmations (saved, sent, undone) — self-dismissing, non-blocking.
- Modal: decisions that must block progress (irreversible confirmation, required auth).
- Push/OS: user-relevant events while away — only with earned permission and per-category control ([15-platform-mobile.md](15-platform-mobile.md#notifications-ux)); respect OS notification settings; don't notify for events users will discover naturally in-app.

**When NOT to use / exceptions:**
- Toasts for critical/must-act info (they vanish; AT users may miss them) — critical = persistent banner/modal. **Never put critical or irreversible information in an auto-dismissing toast.**
- Modals for anything non-blocking (newsletter prompts, ratings) — [18-patterns-antipatterns.md](18-patterns-antipatterns.md#modal-overuse).
- Never stack toasts into a wall; queue, collapse to "+3 more", rate-limit.

**Pros (matched escalation):** Attention aligned to urgency; trust in the channel preserved.
**Cons:** Requires product-wide severity taxonomy and governance — every team wants their event to be a modal.

**Implementation guidance:**
- Toast duration: minimum ~5s even for short messages (average reading speed is ~200–250 wpm, and users must first notice the toast); as a floor, allow roughly 300ms per word plus a ~2s buffer. Toasts carrying actions (Undo) or important information should persist longer or until explicitly dismissed. Bottom or top-corner consistent placement, max 1–2 visible, pause-on-hover, `role="status"`/`aria-live="polite"` (assertive only for errors), actionable toasts keyboard-reachable.
- Banner: dismissible only if the condition is (offline banner isn't); one banner slot per app region.
- Modal: focus-trapped, Esc-closable (unless truly blocking), consequence-stating buttons ([#undoredo-vs-confirmation-dialogs](#undoredo-vs-confirmation-dialogs)).
- Global budget: define per-session interrupt caps; batch low-priority to digests; classify every notification by urgency and required action before choosing its channel; give users frequency and channel controls.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost), [15-platform-mobile.md](15-platform-mobile.md#notifications-ux), [09-content-ux-writing.md](09-content-ux-writing.md#notification-and-toast-copy).

**Sources:** [NN/g: Indicators, Validations, Notifications](https://www.nngroup.com/articles/indicators-validations-notifications/), [NN/g: Toasts](https://www.nngroup.com/articles/toast/), [NN/g: Modal & Nonmodal Dialogs](https://www.nngroup.com/articles/modal-nonmodal-dialog/).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Affordances | Clickable looks clickable; cursor discipline | >90% recognition; hover <100ms |
| Direct manipulation | Continuous, visible, reversible; pair with precision input | — |
| Component states | Full matrix or unfinished; explain disabled | Focus ring ≥3:1, ≥2px |
| Modes | Persistent indicator, easy exit, no danger in ambiguity | Impersonation banner: full session |
| Response thresholds | none <1s → spinner 1–10s → progress >10s | 0.1s / 1s / 10s; 300ms grace |
| Performance measurement | Field (RUM) over lab; budgets as acceptance criteria | p75; INP ≤200ms; CLS ≤0.1 |
| Microinteractions | Design trigger→rules→feedback→loops | Feedback <100ms; loop ≤300ms |
| Motion | Purpose only; ease-out in, ease-in out | 100–300ms; ≤500ms cap; stagger 20–50ms |
| View transitions | Direction encodes hierarchy, consistently | ≤300ms; reduced-motion → fade |
| Loading | Skeleton for content, spinner for widgets, bar for long | Grace 300ms; min shown 500ms |
| Optimistic UI | Instant when reversible + >95% success; never payments/irreversible | Visible rollback + retry |
| Gestures | Standard set only; always visible equivalent; adapt to coarse/fine pointer | Targets 44pt/48dp |
| Drag & drop | Handles, targets, keyboard path | Auto-scroll zone 30–60px; WCAG 2.5.7 |
| Hover | Accelerates, never gates | Tooltip delay 300–500ms |
| Keyboard | Tab between, arrows within; APG key maps | One tab stop per composite widget |
| Undo vs confirm | Undo for reversible; confirm for irreversible | Toast 5–10s; 30-day trash |
| Validation timing | On blur; reward early, punish late | Async debounce 300–500ms |
| Error taxonomy | Validation vs server vs permission — distinct handling; preserve input | Summary at top + focus |
| Scrolling | Never hijack; restore position | Sticky ≤64px; ≤25% viewport fixed |
| Empty states | Four types, each with a next step | — |
| Notifications | Badge→banner→toast→modal by urgency | Toast ≥5s (~300ms/word + 2s floor); 1–2 visible max |
