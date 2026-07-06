# Cognitive Foundations

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers the psychology underlying interface design decisions: cognitive load, memory, attention, scanning, decision-making, error psychology, biases, and flow. Named "laws" derived from this research (Fitts, Hick, Miller, etc.) are in [01-core-principles.md](01-core-principles.md); their visual application in [03-visual-design.md](03-visual-design.md) and interaction application in [04-interaction-design.md](04-interaction-design.md). Cognitive *accessibility* requirements are in [05-accessibility.md](05-accessibility.md#cognitive-accessibility). Research methods to observe these phenomena are in [17-ux-process-research.md](17-ux-process-research.md).

## Contents

- [Load and memory](#load-and-memory)
  - [Cognitive load theory](#cognitive-load-theory)
  - [Working memory and chunking](#working-memory-and-chunking)
  - [Recognition vs recall](#recognition-vs-recall)
  - [Progressive disclosure](#progressive-disclosure)
- [Mental models](#mental-models)
  - [Mental models and conceptual models](#mental-models-and-conceptual-models)
  - [Curse of knowledge](#curse-of-knowledge)
- [Attention and perception](#attention-and-perception)
  - [Selective attention and banner blindness](#selective-attention-and-banner-blindness)
  - [Scanning patterns](#scanning-patterns)
  - [Change blindness](#change-blindness)
  - [Habituation](#habituation)
- [Decision-making](#decision-making)
  - [Decision fatigue and choice overload](#decision-fatigue-and-choice-overload)
  - [Satisficing](#satisficing)
  - [System 1 and System 2](#system-1-and-system-2)
  - [Cognitive biases in UX](#cognitive-biases-in-ux)
- [Error psychology](#error-psychology)
  - [Slips vs mistakes](#slips-vs-mistakes)
  - [Mode errors](#mode-errors)
  - [Post-completion errors](#post-completion-errors)
  - [Learned helplessness and error-tolerant design](#learned-helplessness-and-error-tolerant-design)
- [Engagement and interruption](#engagement-and-interruption)
  - [Flow state](#flow-state)
  - [Interruption cost](#interruption-cost)
- [Quick reference](#quick-reference)

---

## Load and memory

### Cognitive load theory

**Definition:** The total working-memory effort a task imposes, split into intrinsic load (the task's essential difficulty), extraneous load (effort wasted on how information is presented), and germane load (effort building understanding). UX's job: minimize extraneous, manage intrinsic, support germane.

**Reasoning / Evidence:** Sweller (1988), instructional psychology; applied to HCI by NN/g and others. Working memory is the bottleneck of all interaction; when load exceeds capacity, users slow down, err, forget goals, and abandon.

**When to use:**
- Every design review: identify what forces users to think about the *interface* instead of their *task* (unfamiliar patterns, ambiguous labels, visual noise, memory demands between screens).
- Complex domains (finance, healthcare, config): spend the user's limited budget on the domain decision, not on decoding the UI.
- Onboarding: new users carry maximal load; defer nonessential concepts.

**When NOT to use / exceptions:**
- Desirable difficulty: learning tools may deliberately add germane load (retrieval practice); don't optimize it away in education products.
- Intrinsic load can't be removed, only sequenced ([Tesler's Law](01-core-principles.md#teslers-law)) — oversimplifying essential complexity misleads.

**Pros:**
- Unifying framework: most heuristics (consistency, recognition, minimalism) are load-reduction strategies.
- Predicts failure modes testing will find.

**Cons:**
- Load is hard to measure directly (proxies: task time, errors, pupilometry in labs, self-report NASA-TLX).
- "Reduce load" can be misread as "dumb down" — the target is extraneous load only.
- Reducing visible complexity can add navigation depth — hidden complexity is still complexity to reach.

**Implementation guidance:**
- Kill the top extraneous-load sources: unfamiliar patterns (use conventions — [Jakob's Law](01-core-principles.md#jakobs-law)), visual clutter, ambiguous labels, cross-screen memory demands, mid-task decisions that could be defaults.
- Sequence intrinsic load: wizards, staged onboarding, sensible defaults with later refinement.
- Keep context visible; avoid simultaneous competing calls to action.
- Instrument: rising error rates and time-on-step localize overload points.

**Verification checklist:**
- A new user can explain what the screen is for.
- The required next action is visible.
- Important state does not require memory from a previous screen.

**Related:** [01-core-principles.md](01-core-principles.md#8-aesthetic-and-minimalist-design), [#working-memory-and-chunking](#working-memory-and-chunking), [#progressive-disclosure](#progressive-disclosure).

**Sources:** [NN/g: Minimize Cognitive Load](https://www.nngroup.com/articles/minimize-cognitive-load/), [Sweller 1988](https://doi.org/10.1207/s15516709cog1202_4).

### Working memory and chunking

**Definition:** Working memory holds the information currently in conscious use — classically 7±2 items (Miller 1956), by modern estimates ~4±1 chunks (Cowan 2001), decaying within ~10–30 seconds without rehearsal. Chunking groups elements into meaningful units, multiplying effective capacity.

**Reasoning / Evidence:** Miller 1956; Cowan 2001; Baddeley's working-memory model. A "chunk" is defined by meaning: `FBI-CIA-NASA` is 3 chunks for those who know the acronyms, 10 letters for those who don't — chunking capacity depends on user knowledge.

**When to use:**
- Any information users must carry across time or screens: codes, instructions, comparisons.
- Formatting long strings: group digits/characters into 3–4-unit chunks.
- Complex processes: name and group steps so each stage is one chunk.

**When NOT to use / exceptions:**
- Irrelevant to visible, persistent information (recognition, not memory) — do not cap menu items at 7 on memory grounds ([Miller's Law](01-core-principles.md#millers-law) misuse).

**Pros:** Concrete, testable formatting rules; large error-rate reductions on transcription tasks.
**Cons:** Chunk boundaries must match user-meaningful units; arbitrary grouping doesn't help. Chunking can also make comparison harder if related information is split across chunks — keep attributes users compare within the same chunk or view.

**Implementation guidance:**
- Verification codes ≤6 digits, displayed as `123 456`, with copy-to-clipboard; never make users retype what the system already knows (WCAG 3.3.7).
- Card numbers 4-4-4-4; phone numbers in locale format; IBANs in 4-char groups.
- Instructions: max 3–4 steps visible per stage; keep the data users need on the same screen as where they need it.
- Comparison tasks: side-by-side views, never alternate-and-remember.

**Related:** [01-core-principles.md](01-core-principles.md#millers-law), [01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall), [07-forms-and-input.md](07-forms-and-input.md#input-masking-and-formatting).

**Sources:** [Cowan 2001](https://doi.org/10.1017/S0140525X01003922), [NN/g: Working Memory and External Memory](https://www.nngroup.com/articles/working-memory-external-memory/).

### Recognition vs recall

**Definition:** Recognition (identifying something when you see it) is dramatically easier and more reliable than recall (retrieving it from memory unaided). Interfaces should convert recall demands into recognition opportunities.

**Reasoning / Evidence:** Classic memory research: recognition performance exceeds recall across ages and conditions because visible cues do the retrieval work. Basis of GUI superiority over command lines for casual users, of autocomplete, of history/recents.

**When to use:**
- Menus/pickers over free-text entry when the option set is known and small.
- Autocomplete, search suggestions, recently-used lists, favorites.
- Re-entry contexts: show previous values, saved addresses, order history.
- Icon+label defaults; visible shortcuts hints (menus teach the shortcuts).

**When NOT to use / exceptions:**
- Recall-based interfaces (CLI, shortcuts, search-first UIs) are faster for experts with practiced knowledge — provide as accelerator, not replacement ([01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use)).
- Security deliberately uses recall (passwords) — though modern auth is moving to recognition/possession (passkeys) partly for this reason ([07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)).

**Pros:** Slashes learning curve and error rates; serves intermittent users.
**Cons:** Recognition surfaces consume space; long recognition lists trade memory load for scanning load — thus autocomplete (type-to-filter hybrid) often optimal.

**Implementation guidance:**
- Country/state/long lists: combobox with type-ahead, not plain 200-item dropdown.
- Command palettes: fuzzy match + recents + descriptions = recognition over pure recall.
- Show format examples next to fields rather than expecting format recall.

**Related:** [01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall), [08-navigation-ia.md](08-navigation-ia.md#autocomplete-and-suggestions).

**Sources:** [NN/g: Recognition and Recall in UX](https://www.nngroup.com/articles/recognition-and-recall/).

### Progressive disclosure

**Definition:** Show only the information and options needed for the current step; reveal advanced or secondary content on explicit request (accordion, "advanced" panel, drill-down, staged flow).

**Reasoning / Evidence:** Nielsen: progressive disclosure lets one interface serve novices (simple surface) and experts (full depth) simultaneously; it manages Hick-type decision cost and cognitive load by deferring the tail of options ([Pareto](01-core-principles.md#pareto-principle)).

**When to use:**
- Settings: common options visible, advanced behind a disclosure.
- Complex creation flows: required fields first, optional metadata expandable.
- Feature-rich toolbars: primary tools visible, overflow menu for the rest.
- Help: summary first, details expandable.

**When NOT to use / exceptions:**
- Don't hide information needed for the *current* decision (hidden costs at checkout = dark pattern — [18-patterns-antipatterns.md](18-patterns-antipatterns.md#hidden-costs-and-drip-pricing)). The same applies to required fields, critical warnings, privacy consequences, and destructive outcomes.
- Expert-dense tools where all-visible beats click-to-reveal (trading desks, DAWs).
- Comparison content: hiding attributes behind per-item disclosure prevents comparison.
- Frequently needed items behind disclosure = repeated click tax; disclosure boundaries must follow usage data.

**Pros:**
- Serves mixed skill levels with one UI; reduces initial overwhelm; keeps error-prone advanced options away from novices.

**Cons:**
- Discoverability cost: features behind disclosure are found less (hidden-navigation research shows large usage drops).
- Wrong split = experts click constantly; requires usage instrumentation to place the cut.
- Each disclosure level adds interaction cost and state to manage.

**Implementation guidance:**
- Maximum 2 disclosure levels for settings-type content; label the disclosure with what's inside ("Advanced export options", not "More").
- Persist disclosure state per user across sessions for repeated workflows.
- Staged disclosure (wizard steps) for linear tasks; on-demand disclosure (accordions) for reference content.
- Ensure disclosed content is reachable by keyboard and announced to AT ([05-accessibility.md](05-accessibility.md#aria)).

**Related:** [01-core-principles.md](01-core-principles.md#hicks-law), [#cognitive-load-theory](#cognitive-load-theory), [08-navigation-ia.md](08-navigation-ia.md#hierarchy-depth-vs-breadth).

**Sources:** [NN/g: Progressive Disclosure](https://www.nngroup.com/articles/progressive-disclosure/).

---

## Mental models

### Mental models and conceptual models

**Definition:** A mental model is what the user believes about how a system works — built from prior experience with similar systems, the UI itself, and word of mouth. It drives predictions and behavior, and it is usually incomplete and partially wrong.

**Reasoning / Evidence:** Craik (1943); Norman and Nielsen applied to HCI. Users act on their model, not on your implementation; most "user errors" are correct executions of a wrong model. Jakob's Law is mental-model transfer from other products.

**When to use:**
- Feature design start: research the existing model (interviews, card sorts) and either conform to it or explicitly teach a better one.
- Diagnosing recurring "user errors": ask what model would make that behavior correct.
- Sync/sharing/permissions/AI features — chronic model-mismatch zones (e.g., users modeling "the cloud" as a place files move *to*, versus replication).

**When NOT to use / exceptions:**
- Conforming to a wrong-but-common model can be right (folders as metaphor) or limiting; when teaching a new model, invest heavily in the system image and expect a transition cost.

**Pros:** Explains and predicts behavior that pure UI inspection misses; guides which conventions are safe to break.
**Cons:** Models are invisible; eliciting them requires research skill; teams project their own (see [curse of knowledge](#curse-of-knowledge)).

**Implementation guidance:**
- Elicit: ask users to explain/draw how they think it works; note their verbs and metaphors, reuse them in UI copy.
- Design the system image deliberately: every label, animation, and structure teaches the model — check they teach the *same* one.
- Mismatch fixes in order of preference: change design to match user model → teach the model via onboarding/microcopy → accept and mitigate errors.

**Related:** [01-core-principles.md](01-core-principles.md#conceptual-models), [01-core-principles.md](01-core-principles.md#jakobs-law), [17-ux-process-research.md](17-ux-process-research.md#interviews-and-contextual-inquiry).

**Sources:** [NN/g: Mental Models](https://www.nngroup.com/articles/mental-models/).

### Curse of knowledge

**Definition:** Once you know something, you cannot accurately simulate not knowing it. Design teams systematically overestimate users' knowledge of their product, domain terms, and conventions.

**Reasoning / Evidence:** Newton's 1990 "tappers and listeners" study (tappers predicted 50% melody recognition; actual 2.5%). In product teams: internal jargon leaks into UI, "obvious" flows fail testing, edge cases go unimagined.

**When to use:**
- Always as a standing correction: every internal review is biased optimistic.
- Copy review: flag every term the team uses daily (org-chart names, feature codenames, domain abbreviations).
- Justifying user testing to stakeholders who "already know it's fine."

**When NOT to use / exceptions:**
- Domain-expert products: matching professional jargon is correct — the curse is misjudging *which* audience you have.

**Pros:** Explains why testing with 5 outsiders finds what 50 insiders missed.
**Cons:** Cannot be self-corrected by trying harder; only external observation works.

**Implementation guidance:**
- Test with people who have zero product exposure; hallway tests catch the worst cases in minutes.
- Reading-level check UI copy ([09-content-ux-writing.md](09-content-ux-writing.md#reading-level-targets)); ban internal codenames from UI.
- New-hire test: have every new employee use the product on day 1 and log confusion.

**Related:** [#mental-models-and-conceptual-models](#mental-models-and-conceptual-models), [17-ux-process-research.md](17-ux-process-research.md#usability-testing).

**Sources:** [Heath & Heath, Made to Stick — Curse of Knowledge](https://heathbrothers.com/books/made-to-stick/), [Newton 1990 (Stanford dissertation)](https://searchworks.stanford.edu/view/2144755).

---

## Attention and perception

### Selective attention and banner blindness

**Definition:** Users attend to a narrow task-relevant channel and filter the rest. Banner blindness: learned ignoring of anything shaped/placed like advertising or perceived as promotional — including your own important content if it looks promotional.

**Reasoning / Evidence:** NN/g eyetracking across two decades: users skip right-rail content, image-heavy banner-shaped regions, and anything near ad positions; the filter transfers across sites. Selective attention demos (Simons & Chabris' invisible gorilla) show unattended content simply isn't processed.

**When to use:**
- Placing critical content: never style announcements, security notices, or key features like banners/promos.
- Auditing layouts: anything in classic ad positions (top banner, right rail) with graphic styling will be ignored.
- Understanding why "but it's right there on the page" support tickets happen.

**When NOT to use / exceptions:**
- Actual ads: banner blindness is why display advertising underperforms — a business problem outside UX scope here.

**Pros:** Predicts a large, systematic class of "users don't see it" failures.
**Cons:** The main mitigation (make it look like content) is exactly what deceptive advertising does — keep the ethical line: style honestly, don't disguise promos as content ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#disguised-ads)).

**Implementation guidance:**
- Critical messages: place inline in the content/task flow, styled as system UI, not as colorful cards with stock imagery.
- Avoid right-rail placement for anything users must see; integrate into the main column.
- Test findability with first-click/5-second tests, not team review.

**Related:** [#scanning-patterns](#scanning-patterns), [03-visual-design.md](03-visual-design.md#visual-hierarchy).

**Sources:** [NN/g: Banner Blindness Revisited](https://www.nngroup.com/articles/banner-blindness-old-and-new-findings/).

### Scanning patterns

**Definition:** Users don't read screens; they scan in predictable patterns: F-pattern (text-heavy pages: two horizontal sweeps then a left vertical scan), Z-pattern (sparse layouts), layer-cake (scanning headings/subheadings only), spotted (hunting for specific words/links).

**Reasoning / Evidence:** NN/g eyetracking corpus (2006 F-pattern study; 2019 update confirming across devices). Users read ~20–28% of words on an average page. Crucially: F-pattern is *default behavior on badly structured text* — good formatting (headings, bullets, front-loading) produces the more efficient layer-cake scan.

**When to use:**
- All text-bearing screens: front-load first 2 words of headings/links; put key info in the first 2 paragraphs; left-align scannable anchors.
- Layout: critical content top-left region (LTR locales); don't hide essentials bottom-right.
- Lists/tables: first column carries the identifying information.

**When NOT to use / exceptions:**
- RTL locales mirror the patterns — design must mirror too ([09-content-ux-writing.md](09-content-ux-writing.md#internationalization-and-localization)).
- Task-driven hunting (spotted pattern) overrides geometry: users find a distinctively worded link anywhere — distinct labels beat position.
- F-pattern describes failure to be scannable; the goal is to *break* it with structure, not to design for it.

**Pros:** Directly actionable content-structure rules with strong eyetracking evidence.
**Cons:** Patterns vary by task, layout, and culture; misread as immutable geometry ("put everything top-left") rather than formatting guidance.

**Implementation guidance:**
- Headings every 2–4 paragraphs; front-load keywords ("Export reports" not "How you can export your reports").
- Bullets for parallel items; bold key phrases (sparingly — [03-visual-design.md](03-visual-design.md#typography)).
- One idea per paragraph; inverted pyramid ordering ([09-content-ux-writing.md](09-content-ux-writing.md#content-hierarchy-and-front-loading)).

**Related:** [09-content-ux-writing.md](09-content-ux-writing.md#scannability-headings-bullets-bold), [03-visual-design.md](03-visual-design.md#visual-hierarchy).

**Sources:** [NN/g: F-Shaped Pattern](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/), [NN/g: How Little Do Users Read?](https://www.nngroup.com/articles/how-little-do-users-read/).

### Change blindness

**Definition:** Failure to notice changes in a visual scene when the change coincides with a disruption (page load, scroll, saccade) or occurs outside the attended area — even large changes.

**Reasoning / Evidence:** Rensink et al. (1997) flicker paradigm; robust lab phenomenon. In UI: silently updated totals, inline errors appearing far from the button clicked, content swapping during scroll — all routinely missed.

**When to use:**
- Any state change happening away from the user's current focus point: cart totals, background validation, auto-save status, async results arriving.
- Post-submit feedback placement: users look at the button they clicked; feedback appearing only at top of a long page is missed.

**When NOT to use / exceptions:**
- Deliberately quiet updates (ambient presence indicators) may be fine to miss; not every change needs attention.

**Pros:** Explains "the error message was right there" failures; motivates motion/announcement design.
**Cons:** The fix (motion, highlighting) conflicts with reduced-motion needs — pair with ARIA live regions ([05-accessibility.md](05-accessibility.md#aria)).

**Implementation guidance:**
- Announce changes near the user's locus of attention (inline by the button/field), not only in distant regions.
- Use brief highlight animation (e.g., 1–2s background fade, the "yellow-fade technique") on remotely-updated values.
- Async errors after navigation: never show on the previous page only; carry to current view.
- Pair visual change with `aria-live` announcements for AT users.

**Related:** [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds), [05-accessibility.md](05-accessibility.md#aria).

**Sources:** [Rensink et al. 1997](https://doi.org/10.1111/j.1467-9280.1997.tb00427.x), [NN/g: Change Blindness](https://www.nngroup.com/articles/change-blindness-definition/).

### Habituation

**Definition:** Repeated exposure to a stimulus reduces response to it. In UI: users stop reading recurring dialogs, warnings, and toasts, clicking through on autopilot — including the one time it matters.

**Reasoning / Evidence:** Basic learning psychology; documented in security-warning research (users dismiss browser certificate warnings within seconds after few exposures). Confirmation-dialog habituation is why "Are you sure?" fails as protection.

**When to use:**
- Auditing recurring interruptions: any dialog seen weekly is already invisible.
- Designing destructive-action protection: assume habituated click-through; use friction that can't be automated by muscle memory (type-to-confirm) for the rare catastrophic cases.
- Warning design: reserve strong styling for rare events.

**When NOT to use / exceptions:**
- Habituation is also *good*: it's how UIs become fluent and low-effort — the goal is protecting the exceptional moments, not preventing habit.

**Pros:** Explains failure of confirmation-based safety; motivates undo-over-confirm.
**Cons:** Anti-habituation friction (typing the project name to delete) is real friction — reserve for irreversible high-value operations.

**Implementation guidance:**
- Prefer undo to confirmation for reversible actions ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- Vary/strengthen only the rare critical dialog: show consequences concretely ("This deletes 1,204 records permanently"), require typed confirmation for catastrophes.
- Never use the same dialog pattern for routine and critical confirmations.
- Rate-limit toasts/notifications; batch low-priority ones ([04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)).

**Related:** [#slips-vs-mistakes](#slips-vs-mistakes), [01-core-principles.md](01-core-principles.md#5-error-prevention).

**Sources:** [NN/g: Confirmation Dialogs](https://www.nngroup.com/articles/confirmation-dialog/), [Sunshine et al. 2009 (SSL warnings)](https://www.usenix.org/legacy/event/sec09/tech/full_papers/sunshine.pdf).

---

## Decision-making

### Decision fatigue and choice overload

**Definition:** Decision quality and willingness degrade over a sequence of decisions (fatigue); too many options reduce satisfaction and can reduce action (overload / "paradox of choice").

**Reasoning / Evidence:** Iyengar & Lepper's jam study (2000): 24-jam display drew more browsers but 6-jam display converted 10× better. Caveat: meta-analyses (Scheibehenne et al. 2010) show the overload effect is context-dependent — strongest with unfamiliar, hard-to-compare options, no clear preference, and time pressure; near zero when options are easily comparable or users have expertise. Decision fatigue similarly contested in lab replication but consistently observed in long-flow analytics.

**When to use:**
- Long configuration/checkout flows: count the decisions; cut, default, or defer.
- Unfamiliar-choice contexts (plan selection, first-time settings): curate to 3–5 options, mark a recommendation.
- Onboarding: minimize day-1 decisions; ask later, in context.

**When NOT to use / exceptions:**
- Expert users with clear preferences handle large assortments fine (filtering matters more than count).
- Reducing choice below user needs just relocates the cost to search/support.
- Don't cite the jam study as settled law; treat as risk factor, test in your context.

**Pros:** Motivates defaults, recommendations, and flow-shortening with plausible mechanism.
**Cons:** Effect-size uncertainty; "curate for the user" can slide into limiting fair choice ([18-patterns-antipatterns.md](18-patterns-antipatterns.md)).

**Implementation guidance:**
- Decision budget per flow: count every choice a user must make to finish; challenge each (default it? defer it? drop it?).
- Plans/pricing: 3–4 tiers, one visually recommended, comparison table for the detail-hungry.
- Long option lists: add filters/sort/search rather than truncating; support comparability (consistent attributes side-by-side).
- Put important decisions early in flows (before fatigue), trivial ones late or defaulted.

**Related:** [01-core-principles.md](01-core-principles.md#hicks-law), [07-forms-and-input.md](07-forms-and-input.md#defaults-and-prefill), [#satisficing](#satisficing).

**Sources:** [Iyengar & Lepper 2000](https://doi.org/10.1037/0022-3514.79.6.995), [Scheibehenne et al. 2010 meta-analysis](https://doi.org/10.1086/651235).

### Satisficing

**Definition:** Users choose the first acceptable option rather than searching for the optimal one — "good enough" beats "best" when search is costly.

**Reasoning / Evidence:** Herbert Simon (1956, bounded rationality). Web behavior: users click the first plausible link, not the best one (information foraging); they don't read alternatives once one option clears the threshold.

**When to use:**
- Ordering: put the best answer first — users won't scan past a plausible earlier one (search results, suggested actions, nav items).
- Link/label writing: each option must be self-evidently right or wrong; vague labels get satisficed into wrong clicks (pogo-sticking).
- Accepting that users skim instructions: design for partial reading (safe defaults, recoverable errors).

**When NOT to use / exceptions:**
- High-stakes comparison shopping: users do maximize when consequences justify it — support thorough comparison there.

**Pros:** Realistic behavioral model; explains first-position dominance and why "they didn't read it."
**Cons:** Designing for skimmers can under-serve deliberate readers — layer depth beneath scannable surface.

**Implementation guidance:**
- Best/most-likely option first in every list you control; measure first-click accuracy.
- Distinct, front-loaded link labels ([09-content-ux-writing.md](09-content-ux-writing.md#link-text)); avoid near-synonym siblings ("Settings" vs "Preferences").
- Assume zero instruction reading: validate inline, make errors recoverable, put must-know info in the control's label not a paragraph above.

**Related:** [#scanning-patterns](#scanning-patterns), [08-navigation-ia.md](08-navigation-ia.md#information-scent).

**Sources:** [Simon 1956](https://doi.org/10.1037/h0042769), [NN/g: Satisficing](https://www.nngroup.com/articles/satisficing/).

### System 1 and System 2

**Definition:** Kahneman's dual-process model: System 1 — fast, automatic, effortless, pattern-matching; System 2 — slow, deliberate, effortful, logical. Users run interfaces on System 1 whenever possible.

**Reasoning / Evidence:** Kahneman, *Thinking, Fast and Slow* (2011), synthesizing decades of judgment research. UI implication: familiar patterns are processed automatically; anything requiring System 2 (novel patterns, reading, math) feels effortful and is error-prone or skipped.

**When to use:**
- Routine flows: design for System 1 — conventions, clear affordances, recognition, defaults; don't demand reading or calculation.
- Critical irreversible decisions: deliberately *engage* System 2 — break the automatic flow, show concrete consequences, add productive friction (typed confirmation).
- Reviewing designs: identify where you accidentally demand System 2 (mental math on prices, parsing dense text mid-task) and where you dangerously permit System 1 (one-click irreversible).

**When NOT to use / exceptions:**
- The dichotomy is a simplification (contested neuroscience) — use as design heuristic, not brain fact.

**Pros:** Powerful framing for matching interaction cost to decision stakes.
**Cons:** Overused as pop-psych justification for anything; friction placement still needs judgment and testing.

**Implementation guidance:**
- Do the math for users: show totals, differences, "you save X", converted units.
- Match friction to stakes: zero friction for reversible routine; speed bumps (confirmation with consequences, review step) only for irreversible/costly.
- Keep dangerous and routine actions visually and spatially distinct so System 1 pattern-matching doesn't misfire ([#slips-vs-mistakes](#slips-vs-mistakes)).

**Related:** [#habituation](#habituation), [#cognitive-biases-in-ux](#cognitive-biases-in-ux), [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs).

**Sources:** [Kahneman, Thinking Fast and Slow](https://www.penguinrandomhouse.com/books/89308/thinking-fast-and-slow-by-daniel-kahneman/).

### Cognitive biases in UX

**Definition:** Systematic deviations from rational judgment that shape how users perceive and decide in interfaces. Most design-relevant: anchoring (first number sets the scale), framing (equivalent facts, different presentations, different choices), loss aversion (losses weigh ~2× gains), default bias/status quo (pre-set options stick), social proof (others' behavior guides ours), scarcity (rarity implies value), sunk cost (past investment traps continued investment), primacy/recency (position shapes memory — see [Serial Position](01-core-principles.md#serial-position-effect)).

**Reasoning / Evidence:** Tversky & Kahneman (1974) and successors; default bias documented at massive scale (organ-donation consent in opt-in countries ranged from ~4.25% (Denmark) to 27.5%, vs 85–99% in opt-out countries, Johnson & Goldstein 2003); loss-aversion coefficient ~2:1 (Kahneman & Tversky 1979).

**When to use:**
- Setting defaults: the default is a decision you make for most users — make it the honest best-interest option.
- Pricing/plan pages: anchoring (show the reference price), framing ("90% fat-free" vs "10% fat"), decoy effects — powerful, use transparently.
- Adoption/engagement: genuine social proof (real reviews, real usage counts); loss-framing for genuinely at-risk value ("your draft will be discarded").
- Bias-auditing your own design decisions (confirmation bias in research reads, sunk cost in feature investment).

**When NOT to use / exceptions:**
- Every one of these has a dark-pattern twin: fake scarcity, manufactured social proof, preselected paid add-ons, loss-framed manipulation ("your family will be unprotected"). The line: does the nudge serve the user's informed interest, and does it survive disclosure? ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#the-deceptive-neutral-persuasive-spectrum)).
- Regulated contexts (EU DSA, GDPR consent) legally restrict exploiting defaults/framing in consent UIs.

**Pros:** Predicts real behavior better than rational-actor assumptions; defaults alone can drive outsized beneficial outcomes.
**Cons:** Ethically radioactive when user and business interests diverge; effect sizes vary; users increasingly recognize and resent manipulation.

**Implementation guidance:**
- Defaults: safest/most-common-interest option preselected; never prechecked marketing or paid extras.
- Anchors must be real (actual former price, actual competitor price) — fake anchors are illegal in many jurisdictions.
- Social proof must be verifiable; scarcity claims must be true and specific.
- Frame consistently: don't loss-frame to scare, gain-frame the same fact elsewhere.

**Related:** [18-patterns-antipatterns.md](18-patterns-antipatterns.md#dark-patterns-catalog), [07-forms-and-input.md](07-forms-and-input.md#defaults-and-prefill), [#system-1-and-system-2](#system-1-and-system-2).

**Sources:** [Tversky & Kahneman 1974](https://doi.org/10.1126/science.185.4157.1124), [Johnson & Goldstein 2003](https://doi.org/10.1126/science.1091721), [Kahneman & Tversky 1979](https://doi.org/10.2307/1914185).

---

## Error psychology

### Slips vs mistakes

**Definition:** Slips: right intention, wrong execution — automatic behavior misfires (clicking the adjacent button, typo, wrong-window paste). Mistakes: wrong intention formed from a wrong mental model — the plan itself was flawed. They need different countermeasures.

**Reasoning / Evidence:** Norman/Reason human-error taxonomy (Reason 1990, *Human Error*). Slips happen to *experts* on autopilot (System 1); mistakes happen more to novices with immature models. Most "user error" analysis fails by not distinguishing them.

**When to use:**
- Every error analysis: classify observed errors before designing fixes.
- Slip countermeasures: spacing/size/distinctiveness of adjacent controls, undo, input constraints, confirmation-resistant design for destructive neighbors.
- Mistake countermeasures: better system image, clearer labels, teaching the model, previews of consequences.

**When NOT to use / exceptions:**
- Some errors are system errors misattributed to users (misleading UI); fix the design, don't taxonomize the victim.

**Pros:** Correct diagnosis prevents wasted fixes (adding a tutorial won't stop slips; enlarging buttons won't fix model errors).
**Cons:** Classification requires observing intent (think-aloud), not just error logs.

**Implementation guidance:**
- Slips: separate destructive from routine controls (never adjacent Save/Delete); generous targets ([Fitts](01-core-principles.md#fittss-law)); undo everywhere; constrain inputs.
- Mistakes: previews ("this will affect 34 users"), plain-language labels, consequence disclosure before commit, onboarding for genuinely novel models.
- Log-plus-observe: analytics find error hotspots; usability sessions classify them.

**Related:** [01-core-principles.md](01-core-principles.md#5-error-prevention), [#mental-models-and-conceptual-models](#mental-models-and-conceptual-models), [#habituation](#habituation).

**Sources:** [NN/g: Slips vs Mistakes](https://www.nngroup.com/articles/slips/), [Reason, Human Error (1990)](https://doi.org/10.1017/CBO9781139062367).

### Mode errors

**Definition:** Acting correctly for the mode you *think* you're in, while the system is in another mode (typing into the wrong window, caps-lock passwords, editing the wrong environment/production, vim insert-vs-normal).

**Reasoning / Evidence:** Classic HCI failure class (Tesler campaigned against modes at PARC: "Don't Mode Me In"). Aviation/medical incident literature is full of mode errors; severity scales with mode consequences.

**When to use:**
- Any modal design: recording, edit vs view, admin vs user, environment switches (prod/staging), input modes.
- Prefer modeless alternatives: quasimodes (spring-loaded: mode held only while key held, e.g., Shift), direct per-object tools, non-modal editing.

**When NOT to use / exceptions:**
- Modes are sometimes right (focused writing modes, vim's efficiency for experts) — the requirement is unmistakable mode visibility, not mode elimination.

**Pros:** Naming the class drives systematic protection in high-stakes tools.
**Cons:** Persistent mode indicators add chrome; subtle ones fail their purpose.

**Implementation guidance:**
- Mode state must be visible at the locus of attention, not just a distant status bar (cursor change, colored border/frame around the whole workspace for dangerous modes).
- High-stakes modes (production env, live recording): unmistakable persistent indicator (colored top banner) + entry confirmation + auto-exit timeouts where sensible.
- Keyboard focus follows visible cues; avoid focus-stealing that redirects typed input silently.

**Related:** [#slips-vs-mistakes](#slips-vs-mistakes), [01-core-principles.md](01-core-principles.md#1-visibility-of-system-status).

**Sources:** [NN/g: Modes in UI](https://www.nngroup.com/articles/modes/).

### Post-completion errors

**Definition:** Omitting a final step after the main goal is achieved — forgetting the card in the ATM after taking cash, the original in the copier, clicking away before pressing Save.

**Reasoning / Evidence:** Byrne & Bovair (1997): goal completion releases working-memory maintenance of remaining subgoals. Fix proven in ATMs worldwide: force card retrieval *before* dispensing cash (process redesign, not reminders).

**When to use:**
- Flow design: any step after the user's perceived goal-completion moment is at risk — reorder so system-required steps precede user-goal payoff.
- Save/submit design: users who "finished writing" forget submission — autosave, or prominent unsaved-state indicators.

**When NOT to use / exceptions:**
- If reordering is impossible, escalate reminders at exit (unsaved-changes prompt on close) — weaker but necessary fallback.

**Pros:** Process reordering eliminates the error class entirely, beating any reminder.
**Cons:** Reordering may fight business logic ordering; requires flow-map thinking.

**Implementation guidance:**
- Put payoff last: download starts after survey; cash after card.
- Autosave continuously where possible ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)); else block navigation with unsaved-changes dialog (and make it reliable — flaky prompts train dismissal).
- Cleanup steps (log out of shared terminals): auto-timeout as backstop.

**Related:** [#working-memory-and-chunking](#working-memory-and-chunking), [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery).

**Sources:** [Byrne & Bovair 1997](https://doi.org/10.1207/s15516709cog2101_2).

### Learned helplessness and error-tolerant design

**Definition:** After repeated uncontrollable failures, people stop trying — transferring to avoidance of whole categories ("I'm bad at computers"). Error-tolerant design keeps failures cheap, explained, and recoverable so users keep exploring.

**Reasoning / Evidence:** Seligman (1967) learned-helplessness paradigm; in HCI, users with histories of punishing software show reduced exploration and feature adoption. Exploration is how users learn ([01-core-principles.md](01-core-principles.md#3-user-control-and-freedom)); punishing it caps product mastery.

**When to use:**
- Products for low-confidence/broad audiences: government, banking, healthcare portals.
- Error experience design: every failure is a trust event; blame-free, specific, recoverable.
- Data-loss prevention: losing typed work once teaches permanent distrust.

**When NOT to use / exceptions:**
- None practically; even expert tools benefit — experts just have higher failure tolerance thresholds.

**Pros:** Long-term engagement and feature adoption; reduced support avoidance behaviors.
**Cons:** Requires investment in undo, drafts, and error-state design that's invisible when working.

**Implementation guidance:**
- Never lose user input — through errors, timeouts, navigation, or crashes ([07-forms-and-input.md](07-forms-and-input.md#error-message-design)).
- Blame-free copy: "That code didn't match" not "You entered an invalid code" ([09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing)).
- Make exploration visibly safe: undo, preview, sandbox/trial modes, "you can change this later" notes at decision points.

**Related:** [01-core-principles.md](01-core-principles.md#3-user-control-and-freedom), [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing).

**Sources:** [Seligman & Maier 1967](https://doi.org/10.1037/h0024514).

---

## Engagement and interruption

### Flow state

**Definition:** Full absorption in an activity where challenge matches skill; time distorts, self-consciousness drops, performance peaks. Flow-supporting software gets out of the way: clear goals, immediate feedback, no interruptions.

**Reasoning / Evidence:** Csikszentmihalyi (1975–). Preconditions map directly to UI qualities: clear goals (visible state/next steps), immediate feedback (<100ms–400ms loop — [Doherty](01-core-principles.md#doherty-threshold)), challenge-skill balance (progressive complexity), no interruptions.

**When to use:**
- Creation and deep-work tools (editors, IDEs, design tools): latency budgets, distraction-free modes, suppressed non-critical notifications.
- Games/learning: difficulty curves matching growing skill.
- Any product where session productivity is the value proposition.

**When NOT to use / exceptions:**
- Engagement-maximizing "flow" in feeds (infinite scroll + variable reward) is the exploitative twin — flow in service of user goals vs flow as retention trap ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#infinite-scroll)).
- Transactional tools should optimize completion speed, not immersion.

**Pros:** North-star framing for tool UX quality; ties latency, feedback, and interruption policies together.
**Cons:** Hard to measure directly; "engagement" metrics can masquerade as flow.

**Implementation guidance:**
- Response budget: keystroke echo <50ms, common ops <400ms.
- Focus modes: full-screen, notification suppression honoring OS focus/do-not-disturb states.
- Batch non-urgent product messages to natural boundaries (session start/end, task completion), never mid-task.

**Related:** [#interruption-cost](#interruption-cost), [01-core-principles.md](01-core-principles.md#doherty-threshold).

**Sources:** [Csikszentmihalyi, Flow (1990)](https://www.harpercollins.com/products/flow-mihaly-csikszentmihalyi).

### Interruption cost

**Definition:** Interruptions impose a resumption lag — time and errors reacquiring task context — far exceeding the interruption's own duration. Office fieldwork (Mark et al. 2005, "No Task Left Behind?") found workers took on the order of ~25 minutes to return to an interrupted task; the widely quoted "23 minutes 15 seconds" figure comes from Gloria Mark's subsequent interviews and summaries of this research line.

**Reasoning / Evidence:** Mark, González & Harris (2005) documented heavily fragmented work and long resumption delays. Mark, Gudith & Klocke (2008) found a more nuanced result: interrupted work was actually completed *faster* (people compensate by working more intensely), but at the price of significantly higher stress, frustration, time pressure, and effort. Altmann & Trafton's memory-for-goals theory explains the mechanism: task context decays during interruption and must be rebuilt from environmental cues. Interruptions at subtask boundaries cost far less than mid-subtask.

**When to use:**
- Notification policy: default to batching; interrupt only for user-urgent items (security, direct human contact, time-critical).
- Modal/dialog policy: never steal focus while a user is typing; no unsolicited modals mid-task ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#modal-overuse)).
- Resumption support: preserve exact state (scroll, focus, selection, drafts) across interruptions the OS/browser imposes.

**When NOT to use / exceptions:**
- Genuinely critical interrupts (data-loss imminent, security) justify the cost — the bar is user-critical, not business-desired.

**Pros:** Quantifies the harm of interruptive design; strong ammunition against notification-driven engagement tactics.
**Cons:** Timing interrupts to task boundaries requires knowing task structure — approximate via idle detection. The 2008 speed finding is easily misquoted as "interruptions are fine" — the stress/frustration cost is the point.

**Implementation guidance:**
- Defer non-critical notifications to idle moments/boundaries; digest low-priority into summaries.
- Never focus-steal: queue attention requests as badges/subtle banners; let the user choose when.
- Resumption cues: on return, show what changed and where you were ("3 new comments — jump back to your draft").
- Preserve state obsessively: drafts, scroll, form fields survive app switches, session refreshes, crashes.

**Related:** [#flow-state](#flow-state), [04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design), [15-platform-mobile.md](15-platform-mobile.md#notifications-ux).

**Sources:** [Mark, González & Harris 2005](https://doi.org/10.1145/1054972.1055017), [Mark, Gudith & Klocke 2008](https://doi.org/10.1145/1357054.1357072), [Altmann & Trafton 2002](https://doi.org/10.1207/s15516709cog2601_2).

---

## Quick reference

| Concept | One-line guidance | Key numbers |
|---|---|---|
| Cognitive load | Cut extraneous load; sequence intrinsic | NASA-TLX for measurement |
| Working memory | Chunk everything; never demand transcription | ~4±1 chunks; 10–30s decay |
| Recognition vs recall | Show options; type-to-filter hybrids for long lists | — |
| Progressive disclosure | 80% visible, tail on request; label what's hidden | ≤2 levels; usage-data boundaries |
| Mental models | Match or deliberately teach; reuse user vocabulary | — |
| Curse of knowledge | Insiders can't judge clarity — test outsiders | 2.5% vs 50% (tappers study) |
| Banner blindness | Don't style important content like promos | Avoid right rail for essentials |
| Scanning patterns | Front-load; structure beats F-pattern | ~20–28% of words read; first 2 words |
| Change blindness | Change at the locus of attention + announce | 1–2s highlight fade |
| Habituation | Routine dialogs are invisible; undo > confirm | Type-to-confirm for catastrophes |
| Choice overload | Curate unfamiliar choices; recommend one | 3–5 options; jam study 6 vs 24 (contested) |
| Satisficing | Best option first; distinct labels | First plausible link wins |
| System 1/2 | Match friction to stakes; do the math for users | — |
| Biases | Honest defaults, real anchors, true scarcity | Loss aversion ~2:1; defaults: ~4–28% opt-in vs 85–99% opt-out |
| Slips vs mistakes | Spacing/undo for slips; model-teaching for mistakes | Never adjacent Save/Delete |
| Mode errors | Mode state visible at locus of attention | Colored frame for dangerous modes |
| Post-completion errors | Payoff last; autosave as backstop | ATM: card before cash |
| Learned helplessness | Never lose input; blame-free errors | — |
| Flow | Clear goals, fast feedback, no interrupts | Keystroke <50ms; ops <400ms |
| Interruption cost | Batch notifications; preserve state | ~25min resumption (Mark 2005); 2008: faster but more stress |
