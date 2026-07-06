# UX Process & Research

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers how design decisions get made and validated: the research-method landscape, research ethics and inclusive research operations, interviews, surveys, personas and Jobs-to-be-Done, journey mapping, usability testing (including the 5-user rule and its limits), heuristic evaluation, cognitive walkthroughs, A/B testing, analytics, UX metrics (SUS, SEQ, NPS, HEART, and operational metrics), prototyping fidelity, iterative process models (double diamond and the GDS phase model), design sprints, design critique, design review rubrics, debt, prioritization, and AI-assisted research and review. The principles these methods verify live in [01-core-principles.md](01-core-principles.md); accessibility-specific testing in [05-accessibility.md](05-accessibility.md#testing-methodology); executable agent-facing review checklists in [21-agent-checklists.md](21-agent-checklists.md).

## Contents

- [Choosing methods](#choosing-methods)
  - [The research-method landscape](#the-research-method-landscape)
  - [Research ethics and inclusive research operations](#research-ethics-and-inclusive-research-operations)
- [Discovery methods](#discovery-methods)
  - [Interviews and contextual inquiry](#interviews-and-contextual-inquiry)
  - [Surveys](#surveys)
  - [Personas](#personas)
  - [Jobs-to-be-Done](#jobs-to-be-done)
  - [Journey mapping](#journey-mapping)
- [Evaluation methods](#evaluation-methods)
  - [Usability testing](#usability-testing)
  - [The 5-user rule](#the-5-user-rule)
  - [Heuristic evaluation](#heuristic-evaluation)
  - [Cognitive walkthrough](#cognitive-walkthrough)
  - [A/B testing](#ab-testing)
  - [Analytics and behavioral data](#analytics-and-behavioral-data)
- [Measurement](#measurement)
  - [UX metrics frameworks](#ux-metrics-frameworks)
  - [Accessibility auditing in process](#accessibility-auditing-in-process)
- [Ways of working](#ways-of-working)
  - [Prototyping fidelity ladder](#prototyping-fidelity-ladder)
  - [Iterative design and the double diamond](#iterative-design-and-the-double-diamond)
  - [The GDS phase model: Discovery, Alpha, Beta, Live](#the-gds-phase-model-discovery-alpha-beta-live)
  - [Design sprints](#design-sprints)
  - [Design critique](#design-critique)
  - [Design debt](#design-debt)
  - [Prioritizing UX work](#prioritizing-ux-work)
  - [AI-assisted design and research](#ai-assisted-design-and-research)
  - [Design review rubrics and AI-assisted review](#design-review-rubrics-and-ai-assisted-review)
- [Quick reference](#quick-reference)

---

## Choosing methods

### The research-method landscape

**Definition:** The two-axis map for method selection (Rohrer/NN/g): **attitudinal** (what people say: interviews, surveys) vs **behavioral** (what people do: usability tests, analytics, A/B), and **qualitative** (why/how, small-N, direct observation) vs **quantitative** (how many/how much, large-N, statistics). Plus product-phase fit: discovery (generative) → design (formative) → post-launch (summative).

**Reasoning / Evidence:** Method mismatch is the root research failure: asking users what they'd do (attitudinal) to predict behavior — stated preferences diverge from actions notoriously; or A/B testing (what wins) when the question is *why* everything loses. The axes force the right question-to-method mapping.

**When to use:**
- What to build → qualitative attitudinal+contextual (interviews, field studies).
- Does the design work → qualitative behavioral (usability testing).
- How much/for whom → quantitative behavioral (analytics, A/B).
- How do they feel at scale → quantitative attitudinal (surveys, SUS).

**When NOT to use / exceptions:**
- No research beats bad research only when "bad" means biased-and-trusted; imperfect-but-humble beats opinion wars.
- Time-boxed pragmatism: a 3-user hallway test this week beats a 20-user lab study never.

**Pros:** Prevents category errors; gives stakeholders a shared method vocabulary.
**Cons:** Real questions often need method *pairs* (A/B tells you what won; session replays tell you why) — budget accordingly.

**Implementation guidance:**
- Write the decision the research must inform *first*; pick the cheapest method that can change that decision.
- Triangulate big bets: one attitudinal + one behavioral source minimum.
- Map methods to product phase deliberately — the [GDS phase model](#the-gds-phase-model-discovery-alpha-beta-live) gives a per-phase default set.

**Related:** [#usability-testing](#usability-testing), [#ab-testing](#ab-testing), [#interviews-and-contextual-inquiry](#interviews-and-contextual-inquiry).

**Sources:** [NN/g: When to Use Which UX Research Methods](https://www.nngroup.com/articles/which-ux-research-methods/).

### Research ethics and inclusive research operations

**Definition:** The operational floor beneath every research method: **informed consent** (participants understand what's collected, why, by whom, and can withdraw), **accommodations for disabled participants** (accessible recruitment, materials, sessions, and compensation processes), **sensitive-topic care** (extra safeguards when research touches health, finances, trauma, identity, or other emotionally charged areas), and **participant data protection** (secure storage, minimization, retention limits, anonymization in reports).

**Reasoning / Evidence:** Research borrows people's time, candor, and personal data; without consent and protection discipline, it becomes extraction. Practically, ethics failures are also validity failures: participants who don't trust the setup self-censor; inaccessible research operations systematically exclude disabled users from the evidence base, so "validated" designs ship untested against the very users most likely to hit blockers ([05-accessibility.md](05-accessibility.md#testing-methodology)). Regulatory regimes (GDPR and peers) make participant-data mishandling a legal exposure, not just a courtesy lapse.

**When to use:**
- Every study, without exception — consent and data protection are not conditional on study size or formality (a "quick hallway test" that records someone still needs consent).
- Inclusive recruitment as a standing default: include disabled participants for critical flows, not only in dedicated "accessibility studies."
- Heightened protocols for sensitive topics: pre-briefing, explicit opt-outs per question, support-resource signposting, debrief.

**When NOT to use / exceptions:**
- There is no exception to consent; but *proportionality* applies — an anonymous unmoderated survey needs lighter apparatus than recorded in-home contextual inquiry.
- Don't let process weight become an excuse to skip research; a lightweight standing consent-and-data protocol makes ethical research cheap to run repeatedly.

**Pros:** Trustworthy data (participants speak freely), legal and reputational safety, evidence base that actually covers the user population, easier recruitment over time (participants and communities remember how they were treated).
**Cons:** Operational overhead (consent tooling, secure storage, accommodation logistics); sensitive-topic studies need trained researchers, not just enthusiasm.

**Implementation guidance:**
- Consent: plain-language forms ([09-content-ux-writing.md](09-content-ux-writing.md)) covering purpose, recording, data use, retention, withdrawal rights; re-confirm at session start; separate consent for recording vs participation.
- Accommodations: ask about access needs at recruitment (not at the session); budget for interpreters, extra session time, accessible venues/platforms; test your research tooling with assistive technology before the participant does.
- Data protection: store recordings/transcripts in access-controlled locations; anonymize quotes in reports; set and honor deletion schedules; feeding participant data into third-party AI tools requires consent coverage and privacy review ([#ai-assisted-design-and-research](#ai-assisted-design-and-research)).
- Remote vs in-person: choose per participant and task needs, not researcher convenience — remote widens geographic/mobility access; in-person suits contextual observation and some access needs.
- Compensation: pay all participants fairly, via methods disabled participants can actually receive.

**Related:** [#interviews-and-contextual-inquiry](#interviews-and-contextual-inquiry), [#usability-testing](#usability-testing), [#analytics-and-behavioral-data](#analytics-and-behavioral-data), [05-accessibility.md](05-accessibility.md#testing-methodology).

**Sources:** [W3C WAI: Involving Users in Evaluating Web Accessibility](https://www.w3.org/WAI/test-evaluate/involving-users/), [GOV.UK Service Manual: user research](https://www.gov.uk/service-manual/user-research), [NN/g: Obtaining Consent for User Research](https://www.nngroup.com/articles/informed-consent/).

---

## Discovery methods

### Interviews and contextual inquiry

**Definition:** Semi-structured one-on-one conversations about experiences, goals, and workflows; contextual inquiry adds *place* — observing users doing real work in their real environment, interrupting to ask why.

**Reasoning / Evidence:** Interviews surface goals, vocabulary, and mental models ([02-cognitive-foundations.md](02-cognitive-foundations.md#mental-models)) that no analytics can; contextual observation catches what interviews miss — workarounds users don't mention because they've stopped noticing them (habituation), environmental constraints (interruptions, second monitors, sticky notes of passwords). The gap between self-report and observed behavior is the method's justification.

**When to use:**
- Discovery/problem-definition; before building anything; when entering an unfamiliar domain; diagnosing *why* metrics look wrong.
- Contextual inquiry for workflow tools (the environment is half the design input).

**When NOT to use / exceptions:**
- Don't ask users to design ("would you use feature X?" — speculative answers are noise; ask about past behavior instead).
- Not for frequency/magnitude questions (that's surveys/analytics).

**Pros:** Depth, vocabulary capture, unknown-unknowns discovery.
**Cons:** Small-N (don't quote percentages from 8 interviews); interviewer-skill-dependent; recruiting/time cost.

**Implementation guidance:**
- 5–8 per user segment for discovery rounds; 45–60 min; record + transcribe (with consent — [#research-ethics-and-inclusive-research-operations](#research-ethics-and-inclusive-research-operations)).
- Ask about specifics, not generalities ("walk me through the last time you…" not "how do you usually…"); five-whys on surprises; silence is a tool.
- Synthesize within 48h (affinity mapping); outputs: goals, pain points, quotes, vocabulary list ([09-content-ux-writing.md](09-content-ux-writing.md) input).

**Related:** [#personas](#personas), [#jobs-to-be-done](#jobs-to-be-done), [02-cognitive-foundations.md](02-cognitive-foundations.md#curse-of-knowledge).

**Sources:** [NN/g: User Interviews 101](https://www.nngroup.com/articles/user-interviews/), [NN/g: Contextual Inquiry](https://www.nngroup.com/articles/contextual-inquiry/).

### Surveys

**Definition:** Structured questionnaires at scale — powerful for frequency, satisfaction, and segmentation; treacherous because question wording, order, and sampling silently manufacture results.

**Reasoning / Evidence:** Survey pathologies are systematic: leading questions ("How much do you love…"), double-barreled items ("fast and reliable?"), acquiescence bias (agree-tendency), recall error (people misreport past behavior), and survivorship sampling (you survey who stayed). Well-built surveys (validated scales like SUS, neutral wording, piloted) produce trustworthy trends.

**When to use:**
- Quantifying known things (satisfaction tracking, feature-priority ranking among *known* options, segment sizing).
- Standardized instruments post-task/post-experience (SUS, SEQ — [#ux-metrics-frameworks](#ux-metrics-frameworks)).

**When NOT to use / exceptions:**
- Discovering unknowns (open-ends get thin data; interview instead); predicting behavior (stated willingness to pay substantially overstates actual purchase behavior — practitioner experience is consistent on the direction, though inflation magnitude varies and lacks a single citable figure); anything answerable from analytics (asking what you can measure).

**Pros:** Scale, trending over time, statistical segmentation.
**Cons:** Garbage-in sensitivity; response bias (extremes respond); survey fatigue lowering future response rates.

**Implementation guidance:**
- ≤5 min completion; one concept per question; balanced scales (5/7-point Likert, labeled endpoints); randomize option order where order biases; pilot with 5 people thinking aloud.
- Report with N, response rate, and confidence framing; segment before averaging (means hide bimodal truths).

**Related:** [#ux-metrics-frameworks](#ux-metrics-frameworks), [#the-research-method-landscape](#the-research-method-landscape).

**Sources:** [NN/g: Survey Best Practices](https://www.nngroup.com/articles/survey-best-practices/), [MeasuringU: survey design](https://measuringu.com/survey-questions/).

### Personas

**Definition:** Archetypal user representations synthesized from research — goals, behaviors, contexts, pain points — used to keep "the user" concrete in decisions. Evidence-based personas (from real research) vs proto-personas (assumption-documenting placeholders, labeled as such).

**Reasoning / Evidence:** Cooper's original case: engineers design for themselves absent a concrete alternative ([curse of knowledge](02-cognitive-foundations.md#curse-of-knowledge)); a memorable shared referent ("would Marta understand this?") counters it. Critiques are substantive: personas ossify, get decorated with fictional demographics that invite stereotyping, and — the JTBD critique — attributes don't cause behavior, situations do.

**When to use:**
- Multi-segment products where teams chronically design for the wrong (or their own) profile; onboarding new team members to who-users-are.
- Built from real interview/analytics synthesis; refreshed on a schedule; focused on behavior/goals, not demographics-and-stock-photos.

**When NOT to use / exceptions:**
- Proto-personas presented as research = fiction with authority — label or don't use.
- When segments behave identically for your decisions, personas add ceremony without cuts.
- Consider [JTBD](#jobs-to-be-done) instead/alongside when the situation drives behavior more than the person.

**Pros:** Shared referent, empathy artifact, decision tiebreaker.
**Cons:** Maintenance decay, stereotyping risk, weak causal power vs situational frames.

**Implementation guidance:**
- 3–5 max (more = none memorable); one page each: goals, behaviors, context, pains, quote, design implications; behavioral dimensions over demographics.
- Tie to data: each claim traceable to research; version and date them.

**Related:** [#jobs-to-be-done](#jobs-to-be-done), [#interviews-and-contextual-inquiry](#interviews-and-contextual-inquiry).

**Sources:** [NN/g: Personas 101](https://www.nngroup.com/articles/persona/), [Cooper: The Inmates Are Running the Asylum](https://www.pearson.com/en-us/subject-catalog/p/inmates-are-running-the-asylum-the-why-high-tech-products-drive-us-crazy-and-how-to-restore-the-sanity/P200000009348).

### Jobs-to-be-Done

**Definition:** Framing demand around the *progress users are trying to make in a situation* — "when [situation], I want to [motivation], so I can [outcome]" — rather than user attributes: people "hire" products for jobs (Christensen's milkshake: hired for the boring commute, competing with bananas and bagels, not other milkshakes).

**Reasoning / Evidence:** Situations predict behavior better than demographics (the same person hires different solutions in different contexts); JTBD reframes competition realistically (spreadsheets and doing-nothing are competitors) and generates requirements ("make progress on X in situation Y") that personas' attribute-lists don't.

**When to use:**
- Positioning and roadmap decisions (what job, against which real alternatives); feature definition (job stories over user stories when situation matters); innovation hunting (underserved jobs).
- Alongside personas: personas say who and vocabulary; JTBD says when-why-what-for.

**When NOT to use / exceptions:**
- Micro-interaction design decisions need task-level detail, not job-level abstraction.
- JTBD orthodoxy wars (forces-diagram factions) are noise; take the framing, skip the theology.

**Pros:** Causal framing, honest competitive lens, stable over time (jobs outlive solutions).
**Cons:** Abstraction skill required (too-broad jobs: "communicate better"; too-narrow: feature restatements).

**Implementation guidance:**
- Switch interviews (why did you fire the old solution and hire this?) are the method's research engine; map forces: push of problem, pull of new, anxiety of change, habit of present.
- Job statements: situation + motivation + outcome, solution-free.

**Related:** [#personas](#personas), [#interviews-and-contextual-inquiry](#interviews-and-contextual-inquiry).

**Sources:** [Christensen: Know Your Customers' Jobs to Be Done (HBR)](https://hbr.org/2016/09/know-your-customers-jobs-to-be-done), [NN/g: Personas vs Jobs-to-Be-Done](https://www.nngroup.com/articles/personas-jobs-be-done/).

### Journey mapping

**Definition:** Visualizing a user's end-to-end experience toward a goal across stages, touchpoints, actions, thoughts, and emotions — exposing cross-silo failures; service blueprints extend it with the backstage (systems, staff, processes) that produces each touchpoint.

**Reasoning / Evidence:** Organizations are structured by function; experiences run across functions — the seams (handoffs between app/email/support/store) are where journeys break, and no single team sees them. The emotion curve locates the [peak and end](01-core-principles.md#peak-end-rule) worth fixing/polishing.

**When to use:**
- Multi-touchpoint experiences (onboarding→activation, purchase→delivery→support); prioritizing which broken moment matters most; aligning cross-functional owners on one picture.
- Blueprint when the fix is operational (the backstage causes the frontstage pain).

**When NOT to use / exceptions:**
- Single-screen problems don't need journey ceremony.
- Maps built from assumptions in a workshop = shared fiction; build from research (interviews, analytics, support logs) or label as hypothesis.

**Pros:** Cross-silo visibility, prioritization by emotional low-point, stakeholder alignment artifact.
**Cons:** Wall-art risk (beautiful, ignored) — tie each pain to an owner and a backlog item or don't bother.

**Implementation guidance:**
- Scope one persona/job + one goal per map; stages → actions → touchpoints → quotes/emotions → pains → opportunities.
- Overlay data: drop-off rates per stage, support-ticket volumes per touchpoint ([#analytics-and-behavioral-data](#analytics-and-behavioral-data)).

**Related:** [01-core-principles.md](01-core-principles.md#peak-end-rule), [#interviews-and-contextual-inquiry](#interviews-and-contextual-inquiry).

**Sources:** [NN/g: Journey Mapping 101](https://www.nngroup.com/articles/journey-mapping-101/), [NN/g: Service Blueprints](https://www.nngroup.com/articles/service-blueprints-definition/).

---

## Evaluation methods

### Usability testing

**Definition:** Watching representative users attempt realistic tasks with the product/prototype, thinking aloud, while you observe struggles — the core empirical UX method. Variants: moderated (facilitator probes) vs unmoderated (recorded solo via platform), remote vs in-person, formative (find problems, small-N, iterative) vs summative (measure against benchmarks, larger-N).

**Reasoning / Evidence:** Fifty years of practice: watching 5 users attempt tasks finds problems that shock teams every single time ([curse of knowledge](02-cognitive-foundations.md#curse-of-knowledge) is invisible from inside). Behavioral observation beats all opinion sources — including the users' own post-task opinions ([aesthetic-usability effect](01-core-principles.md#aesthetic-usability-effect) inflates ratings of pretty-but-broken).

**When to use:**
- Every significant flow before launch; prototypes before code ([#prototyping-fidelity-ladder](#prototyping-fidelity-ladder)); recurring cadence (monthly-ish) rather than launch-gated events.
- Moderated for complex/exploratory questions; unmoderated for speed/scale on well-defined tasks.
- Include disabled participants for critical flows ([#research-ethics-and-inclusive-research-operations](#research-ethics-and-inclusive-research-operations)).

**When NOT to use / exceptions:**
- Frequency/preference questions (wrong tool); visual-appeal comparisons (preference tests instead).
- Don't test with colleagues/friends for real studies (contaminated by context knowledge).

**Pros:** Finds real problems with certainty analytics can't (you see the *why*); cheap relative to shipping failures.
**Cons:** Artificial-setting effects (mitigate with realistic tasks/context); small-N limits quantification; moderation skill matters (leading users invalidates).

**Implementation guidance:**
- Tasks: goal-phrased, not click-instructions ("You need to expense this receipt" not "Click Expenses"); no vocabulary from the UI in the task (echo effect).
- Think-aloud: instruct once, prompt neutrally ("what are you thinking?"), never rescue early — the struggle is the data.
- Capture: task success (3-level: success/struggled/fail), time, errors, quotes; severity-rate findings (blocker/major/minor/cosmetic × frequency).
- Report within days: video clips of failures persuade stakeholders better than any deck.

**Related:** [#the-5-user-rule](#the-5-user-rule), [#heuristic-evaluation](#heuristic-evaluation), [05-accessibility.md](05-accessibility.md#testing-methodology).

**Sources:** [NN/g: Usability Testing 101](https://www.nngroup.com/articles/usability-testing-101/), [NN/g: Thinking Aloud](https://www.nngroup.com/articles/thinking-aloud-the-1-usability-tool/).

### The 5-user rule

**Definition:** Nielsen's heuristic: ~5 users per formative test round find ~85% of the problems *that affect users at that problem-frequency profile*, with diminishing returns beyond — so run small rounds, fix, and re-test, rather than one big study.

**Reasoning / Evidence:** Nielsen & Landauer (1993): problem-discovery follows 1−(1−L)^n with average per-user detection L≈31%, giving ~85% at n=5. The honest caveats: L varies by problem subtlety and product complexity (complex domains: L lower, need more users); the math is per-*segment* (2 distinct user groups ≈ 3–5 each); it's about *finding problems*, not measuring metrics (summative measurement needs 20+ for stable numbers); rare-context problems (specific data states) escape any small-N test.

**When to use:**
- Formative iteration: 5-ish users → fix → 5 more on the revision (3 rounds of 5 beats 1 round of 15, always).

**When NOT to use / exceptions:**
- Quantitative benchmarking (n≥20/segment); card sorting (~15 participants reaches ~0.90 correlation with larger-sample category structures — Tullis & Wood, NN/g); tree testing (benefits from larger samples, ~50); A/B (thousands — [#ab-testing](#ab-testing)); highly heterogeneous audiences (per-segment math).
- Don't wield "5 is enough" to end research after one round — the rule's whole point is *iteration*.

**Pros:** Makes testing affordable → makes it *happen*; iteration catches fix-induced regressions.
**Cons:** Misquoted constantly (as "5 users suffice for anything"); false confidence when segments/complexity ignored.

**Implementation guidance:**
- Budget frame: same spend = more rounds × fewer users; recruit 6 to seat 5 (no-shows).
- Track cumulative: if round 3 still surfaces new blockers, the design (or task coverage) needs more than testing cadence.

**Related:** [#usability-testing](#usability-testing), [#the-research-method-landscape](#the-research-method-landscape).

**Sources:** [NN/g: Why You Only Need to Test with 5 Users](https://www.nngroup.com/articles/why-you-only-need-to-test-with-5-users/), [Nielsen & Landauer 1993](https://doi.org/10.1145/169059.169166), [Tullis & Wood 2004: How Many Users Are Enough for a Card-Sorting Study?](https://www.nngroup.com/articles/card-sorting-how-many-users-to-test/).

### Heuristic evaluation

**Definition:** Expert inspection: 3–5 evaluators independently audit the interface against a heuristic set ([Nielsen's 10](01-core-principles.md#nielsens-10-usability-heuristics)), then aggregate findings with severity ratings — a discount method finding many problems without recruiting users.

**Reasoning / Evidence:** Nielsen & Molich (1990): single evaluators find ~35% of problems; 3–5 independent evaluators aggregate to ~60–75% (the "evaluator effect": different experts find different problems — independence before aggregation is essential). Complements testing: catches consistency/standards violations experts see instantly and users can't articulate; misses domain-specific and first-use problems only real users reveal.

**When to use:**
- Early designs (before user-testable fidelity); pre-test cleanup (don't spend test participants on findable-by-inspection issues); audits of legacy surfaces; when recruitment is impossible.

**When NOT to use / exceptions:**
- Never as a *replacement* for user testing on major flows (experts systematically miss novel-user confusion — they're not novices).
- One-evaluator "heuristic evaluations" are opinion memos; the method requires multiple independent passes.

**Pros:** Fast (days), cheap, no recruitment, teaches the heuristics to the team.
**Cons:** False positives (expert pet peeves users never hit); evaluator quality variance; misses real-user surprises.

**Implementation guidance:**
- Process: each evaluator solo-passes twice (flow overview, then element detail), logs issue + violated heuristic + location; aggregate, dedupe, severity-rate (0–4 scale: cosmetic→catastrophe) by frequency × impact × persistence.
- Severity dimensions worth scoring explicitly (beyond frequency × impact × persistence): task blockage, error consequence, affected user groups, and **accessibility exclusion** (an issue that fully blocks assistive-technology users is a blocker regardless of how few users it "affects" in raw counts — [05-accessibility.md](05-accessibility.md)). Track **fix complexity** as a *separate* field: it belongs in scheduling ([#prioritizing-ux-work](#prioritizing-ux-work)), never discounted into the severity score itself.
- Mixed evaluator panel (UX + domain expert) widens coverage; validate top findings in the next user test.

**Related:** [01-core-principles.md](01-core-principles.md#nielsens-10-usability-heuristics), [#usability-testing](#usability-testing), [#cognitive-walkthrough](#cognitive-walkthrough), [#design-review-rubrics-and-ai-assisted-review](#design-review-rubrics-and-ai-assisted-review).

**Sources:** [NN/g: How to Conduct a Heuristic Evaluation](https://www.nngroup.com/articles/how-to-conduct-a-heuristic-evaluation/), [Nielsen & Molich 1990](https://doi.org/10.1145/97243.97281).

### Cognitive walkthrough

**Definition:** Task-focused inspection simulating a *first-time user's* learning path: for each step of a target task, evaluators ask — will users know what to try? Will they notice the right control? Will they connect it to their goal? Will they understand the feedback? — a learnability microscope where heuristic evaluation is a general-usability floodlight.

**Reasoning / Evidence:** Polson/Lewis/Wharton (1990s), grounded in exploratory-learning theory: users don't read manuals, they act on the interface's visible promises ([signifiers](01-core-principles.md#signifiers), [information scent](08-navigation-ia.md#information-scent)). The four questions operationalize where exploration breaks.

**When to use:**
- Walk-up-and-use products (kiosks, government forms, onboarding flows); critical first-run paths; when novice learnability is the risk (vs expert efficiency).

**When NOT to use / exceptions:**
- Expert-tool efficiency questions (wrong lens); whole-product audits (too slow per step — pick the 2–3 critical tasks).

**Pros:** Cheap, pre-prototype-capable (works on wireframes), pinpoints the exact failing step and *why*.
**Cons:** Tedious at scale; evaluator empathy limits (experts simulating novices imperfectly — [02-cognitive-foundations.md](02-cognitive-foundations.md#curse-of-knowledge) again).

**Implementation guidance:**
- Define: user's starting knowledge, the task, the correct action sequence; walk each action against the four questions; log failure stories ("user won't know the gear icon holds export").
- Fix categories map to causes: wrong label → [09-content-ux-writing.md](09-content-ux-writing.md); invisible control → [signifiers](01-core-principles.md#signifiers); missing feedback → [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds).

**Related:** [#heuristic-evaluation](#heuristic-evaluation), [02-cognitive-foundations.md](02-cognitive-foundations.md#mental-models).

**Sources:** [NN/g: Cognitive Walkthroughs](https://www.nngroup.com/articles/cognitive-walkthroughs/), [Wharton et al. 1994](https://dl.acm.org/doi/10.5555/189200.189214).

### A/B testing

**Definition:** Randomized controlled experiments on live traffic: variants A/B(/n) assigned randomly, a predefined metric compared for statistically significant difference — the strongest causal evidence for *which* performs better, and silent on *why*.

**Reasoning / Evidence:** Randomization eliminates confounds correlational analytics can't; at scale (Microsoft/Booking/Netflix experimentation literature) most shipped ideas *lose or tie* — humbling evidence for testing over opinion. The traps are statistical: peeking (checking until significance appears — inflates false positives massively), underpowering (small samples "prove" noise), metric myopia (winning clicks while losing retention), novelty effects (change itself moves metrics temporarily), and p-hacking (testing 20 metrics, reporting the significant one).

**When to use:**
- High-traffic surfaces + quantifiable proximate metrics (conversion, activation, completion); incremental optimizations; settling genuine team disagreements with reality.

**When NOT to use / exceptions:**
- Low traffic (weeks-to-significance = never; use qualitative methods); radical redesigns (novelty + too-many-variables — test the concept qualitatively first); anything unethical to randomize; strategy questions (A/B optimizes local hills, doesn't choose mountains).
- Never ship dark patterns because they "won" the metric ([18-patterns-antipatterns.md](18-patterns-antipatterns.md) — short-term metric wins routinely mask long-term trust losses).

**Pros:** Causal, scaled, opinion-proof.
**Cons:** Why-blind (pair with session replays/tests); guardrail-metric discipline required; experimentation infrastructure cost.

**Implementation guidance:**
- Pre-register: hypothesis, primary metric, guardrails (retention, support contacts, latency), sample size (power calculation), duration (≥1–2 full weekly cycles); no peeking (or use sequential-testing corrections).
- Segment results cautiously (subgroup fishing = p-hacking); replicate surprising wins before believing them.
- Long-term holdbacks for cumulative-effect visibility on engagement-adjacent changes.

**Related:** [#analytics-and-behavioral-data](#analytics-and-behavioral-data), [#usability-testing](#usability-testing), [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Sources:** [Kohavi et al.: Trustworthy Online Controlled Experiments](https://experimentguide.com/), [NN/g: A/B Testing 101](https://www.nngroup.com/articles/ab-testing/).

### Analytics and behavioral data

**Definition:** Instrumented product usage at scale: funnels (step-wise drop-off), feature adoption, paths, heatmaps/scroll maps, session recordings, error/rage-click tracking — the *what/where/how-many* layer that locates problems for qualitative methods to explain.

**Reasoning / Evidence:** Analytics see every user (no recruitment bias) and reveal magnitude (which drop-off costs most — [Pareto](01-core-principles.md#pareto-principle) targeting); they cannot see intent, confusion, or the users who never arrived. Rage-clicks/dead-clicks and error-hotspots are direct usability-defect detectors; funnel cliffs are prioritization gold.

**When to use:**
- Continuously: core-flow funnels, error rates per field/step ([07-forms-and-input.md](07-forms-and-input.md#error-message-design)), search-with-no-results logs ([08-navigation-ia.md](08-navigation-ia.md#no-results-recovery)), feature adoption before redesigns ([01-core-principles.md](01-core-principles.md#pareto-principle)).
- Session recordings: targeted (watch sessions that hit the funnel cliff), not voyeuristic browsing.

**When NOT to use / exceptions:**
- Privacy discipline is non-negotiable: consent compliance ([13-platform-web.md](13-platform-web.md#consent-and-cookie-ux)), PII masking in recordings (inputs masked by default), retention limits, and honesty that some jurisdictions/audiences preclude recording.
- Don't infer intent from behavior alone ("users love this page — 10min average!" …because they're stuck).

**Pros:** Scale, real behavior, continuous, prioritization-grade magnitudes.
**Cons:** Why-blind; instrumentation debt (broken events = confident lies); metric-fixation risk ([#ux-metrics-frameworks](#ux-metrics-frameworks) guardrails).

**Implementation guidance:**
- Instrument decisions, not everything: each dashboard answers a standing question with an owner.
- Pair quant→qual ritual: every funnel cliff gets session-replay review + a usability-test task probing it.
- Data hygiene: event naming schema, QA on instrumentation per release, annotation of releases/campaigns on time series.

**Related:** [#ab-testing](#ab-testing), [#ux-metrics-frameworks](#ux-metrics-frameworks), [#the-research-method-landscape](#the-research-method-landscape).

**Sources:** [NN/g: Analytics and User Experience](https://www.nngroup.com/articles/analytics-user-experience/).

---

## Measurement

### UX metrics frameworks

**Definition:** The standard instruments and frameworks: **SUS** (System Usability Scale: 10 items → 0–100 score; mean ≈ 68, so 68 = 50th percentile, 80+ ≈ top decile); **SEQ** (Single Ease Question, post-task 1–7; benchmark ≈ 5.5); **NPS** (loyalty proxy; methodologically criticized — coarse, volatile, gameable — but organizationally entrenched); **HEART** (Google's framework: Happiness, Engagement, Adoption, Retention, Task success — each with goals→signals→metrics); plus core behavioral trio: task success rate, time-on-task, error rate.

**Reasoning / Evidence:** Validated instruments (SUS: decades of benchmarks across thousands of studies; MeasuringU's rule of thumb is that SUS gives usable estimates around n≈12, though smaller samples widen confidence intervals sharply) beat homegrown "how easy was that 1–10" questions; HEART's contribution is the goals-signals-metrics discipline (metrics chosen to reflect decisions, not availability). NPS criticisms (Bain-marketing origins, statistical information loss from the 11-point trichotomy, question predicting *stated* intent not behavior) justify treating it as a coarse trend, never a design-quality metric.

**When to use:**
- SUS: perceived-usability benchmarking across releases/competitors; SEQ: per-task friction tracking in every usability test; HEART: choosing product-level UX metrics deliberately; behavioral trio: summative testing and monitoring.
- Operational extensions worth tracking alongside the trio: **task recovery rate** (of users who err, how many recover and complete — separates recoverable friction from dead ends), **zero-result search rate** ([08-navigation-ia.md](08-navigation-ia.md#no-results-recovery) — vocabulary/content gaps made measurable), **support contact rate** per flow (each contact is a UX failure with a price tag), and **accessibility defect counts** by severity ([#accessibility-auditing-in-process](#accessibility-auditing-in-process)).

**When NOT to use / exceptions:**
- Any single metric as *the* goal invites Goodhart's law (metric gamed, experience unimproved) — pair outcome metrics with guardrails; a metric that proxies business goals rather than user success can be "optimized" while harming users — pair quantitative signals with qualitative research.
- Small-N ritualism: SUS from 3 users is numerology.

**Pros:** Comparability (benchmarks), trend detection, executive-legible UX accountability.
**Cons:** Perceived ≠ actual usability ([aesthetic-usability](01-core-principles.md#aesthetic-usability-effect) inflates scores); framework theater (HEART posters, no decisions changed).

**Implementation guidance:**
- Pick ≤5 product UX metrics via HEART's goals→signals→metrics table; review in the same forum as business metrics.
- SUS administration: verbatim items, alternating polarity preserved, after real usage; report with confidence intervals.
- Track the trio per critical task: success ≥ target (e.g., 90%), time trend, error rate trend — the design-quality dashboard; add recovery rate where error rate alone hides whether users escape.

**Related:** [#analytics-and-behavioral-data](#analytics-and-behavioral-data), [#usability-testing](#usability-testing).

**Sources:** [MeasuringU: SUS](https://measuringu.com/sus/), [Google HEART framework (Rodden et al.)](https://research.google/pubs/pub36299/), [MeasuringU: SEQ](https://measuringu.com/seq10/).

### Accessibility auditing in process

**Definition:** Embedding accessibility verification into the delivery lifecycle — design-review checks, CI automation, release-gate manual passes, scheduled audits — rather than pre-launch heroics. (Methods detail: [05-accessibility.md](05-accessibility.md#testing-methodology).)

**Reasoning / Evidence:** Retrofit costs multiply; the ~30–40% automation ceiling means process must schedule *human* checks or they never happen; legal exposure ([05-accessibility.md](05-accessibility.md#why-accessibility)) makes "we'll audit someday" a liability posture.

**When to use:**
- Design reviews: contrast/target/focus-order annotations checked before build; component definition-of-done includes keyboard+SR ([10-design-systems.md](10-design-systems.md#accessibility-baked-into-components)).
- CI: axe-core per PR; release: keyboard-only pass of changed flows; quarterly-to-annual: full audit incl. SR matrix and disabled-user testing.

**Implementation guidance:**
- Track WCAG issues in the main bug tracker with severity = user-impact ([05-accessibility.md](05-accessibility.md#testing-methodology)); accessibility acceptance criteria in every UI story template.
- Trend accessibility defect counts as a product UX metric ([#ux-metrics-frameworks](#ux-metrics-frameworks)); a rising count under feature pressure is design debt with legal interest ([#design-debt](#design-debt)).

**Related:** [05-accessibility.md](05-accessibility.md#testing-methodology), [10-design-systems.md](10-design-systems.md#accessibility-baked-into-components), [#research-ethics-and-inclusive-research-operations](#research-ethics-and-inclusive-research-operations).

**Sources:** [W3C WAI: Involving Users](https://www.w3.org/WAI/planning/involving-users/), [GOV.UK: accessibility in agile](https://www.gov.uk/service-manual/helping-people-to-use-your-service/making-your-service-accessible-an-introduction).

---

## Ways of working

### Prototyping fidelity ladder

**Definition:** Matching prototype fidelity to the question asked: **paper/sketch** (concept/flow shape, minutes to make) → **wireframe/clickable low-fi** (IA, flow logic) → **hi-fi mockup/interactive prototype** (visual design, microcopy, detailed usability) → **coded prototype** (feel, motion, performance, real data). Rule: the cheapest artifact that can falsify your current hypothesis.

**Reasoning / Evidence:** Fidelity shapes feedback: rough artifacts invite structural criticism ("this flow is wrong"); polished artifacts suppress it and attract paint-color comments (participants assume it's nearly done — and the [aesthetic-usability effect](01-core-principles.md#aesthetic-usability-effect) inflates their ratings). Over-investing early also breeds sunk-cost attachment in the *team* ([02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux)).

**When to use:**
- Concept validity → paper/sketches with users or stakeholders; IA/flow → clickable wireframes ([tree tests](08-navigation-ia.md#tree-testing) for structure); usability of real designs → hi-fi interactive; motion/latency/real-content questions → coded (real data breaks layouts that lorem ipsum blessed).

**When NOT to use / exceptions:**
- Skipping fidelity rungs is right when conventions cover the skipped questions (a standard CRUD screen needs no paper phase).
- Hi-fi-first cultures (design-system-assembled screens are cheap now) are fine *if* structural questions were answered elsewhere.

**Pros:** Cost-matched learning; feedback-type steering; iteration speed at low rungs.
**Cons:** Prototype-reality gaps (canned data, no latency, happy paths) — flag what the prototype *can't* answer.

**Implementation guidance:**
- Label fidelity intent when sharing ("structure only — ignore visuals") to steer feedback.
- Test prototypes with realistic data extremes (longest names, empty states, error paths) before concluding.

**Related:** [#usability-testing](#usability-testing), [#iterative-design-and-the-double-diamond](#iterative-design-and-the-double-diamond).

**Sources:** [NN/g: Prototype fidelity](https://www.nngroup.com/articles/ux-prototype-hi-lo-fidelity/).

### Iterative design and the double diamond

**Definition:** The macro process models: **iteration** (design → test → learn → redesign, repeatedly — the empirical core of all UX process) framed by the **double diamond** (Design Council): diverge/converge on the *problem* (Discover→Define), then diverge/converge on the *solution* (Develop→Deliver) — "design the right thing, then design the thing right."

**Reasoning / Evidence:** Single-pass design fails because requirements are discovered through use (no analysis substitutes for contact with users — the entire evidence base of [#usability-testing](#usability-testing)); the double diamond's contribution is legitimizing problem-space work before solution generation (teams default to solutioneering the first problem framing).

**When to use:**
- Always iterate: budget ≥2–3 design-test cycles into any significant feature plan (a plan with one design phase is a plan to ship the first draft).
- Double diamond weighting by uncertainty: new-product/new-domain → heavy first diamond; well-understood incremental work → light-touch first diamond.

**When NOT to use / exceptions:**
- Process theater: diamonds as ceremony gates rather than thinking modes; iteration as excuse for shipping broken ("we'll fix it in iteration") — iteration is pre-release learning, not post-release apology.

**Pros:** Matches how design knowledge actually accrues; shared process vocabulary with stakeholders.
**Cons:** Time-boxing tension (converge deadlines vs learning) — resolve by scoping cycles, not skipping them.

**Implementation guidance:**
- Make iteration cheap structurally: design-system parts ([10-design-systems.md](10-design-systems.md)), standing research cadence, prototype infrastructure.
- Exit criteria per phase: Define exits with a problem statement + success metrics; Deliver exits with tested-against-those-metrics.

**Related:** [#prototyping-fidelity-ladder](#prototyping-fidelity-ladder), [#design-sprints](#design-sprints), [#the-gds-phase-model-discovery-alpha-beta-live](#the-gds-phase-model-discovery-alpha-beta-live).

**Sources:** [Design Council: Double Diamond](https://www.designcouncil.org.uk/our-resources/the-double-diamond/), [NN/g: Design Thinking 101](https://www.nngroup.com/articles/design-thinking/).

### The GDS phase model: Discovery, Alpha, Beta, Live

**Definition:** The UK Government Digital Service's delivery phasing — **Discovery** (understand users, needs, constraints, existing workarounds, the surrounding service ecosystem; decide whether to build at all) → **Alpha** (prototype and test the risky assumptions: concepts, IA, terminology, task flow) → **Beta** (build and test realistic end-to-end tasks: accessibility, edge cases, performance, content, at increasing scale) → **Live** (continuous monitoring and improvement: analytics, support contacts, search logs, feedback, accessibility defects, task success, service reliability). It is a different framing of the same idea as the [double diamond](#iterative-design-and-the-double-diamond): Discovery+Alpha ≈ the first diamond (right problem, tested direction); Beta+Live ≈ the second (right solution, then keep learning) — with the useful additions of explicit *kill points* between phases and a *Live* phase that says the process never ends at launch.

**Reasoning / Evidence:** GDS built the model for high-stakes services (government: everyone is a user, failure excludes citizens); its discipline — each phase gates the next, and Discovery can conclude "don't build this" — counters the near-universal failure of treating discovery research as a formality before a predetermined build. Mapping research methods to phases operationalizes [#the-research-method-landscape](#the-research-method-landscape): generative methods in Discovery, formative in Alpha/Beta, summative and monitoring in Live.

**When to use:**
- Service-scale work with real uncertainty about whether/what to build; regulated or public-sector contexts where the phase gates double as assurance points; teams needing shared phase vocabulary with non-design stakeholders.
- As a per-phase research-method default: Discovery → interviews/field studies/support-log analysis; Alpha → prototype tests, tree tests, terminology checks; Beta → end-to-end usability, accessibility testing with disabled participants, performance/content review; Live → analytics, search-log and support-contact monitoring, accessibility defect tracking.

**When NOT to use / exceptions:**
- Small incremental features don't need four ceremonial phases — the double diamond's "light-touch first diamond" advice applies equally here.
- Phase-model theater: phases as calendar milestones rather than evidence gates (an Alpha that cannot be cancelled is not an Alpha).

**Pros:** Explicit kill points before expensive build; Live phase institutionalizes post-launch learning; strong shared vocabulary across disciplines.
**Cons:** Government-scale ceremony can feel heavyweight for product teams; naming collision with agile "sprint/beta" vocabulary causes confusion — define terms.

**Implementation guidance:**
- Write each phase's exit question before entering it (Discovery: "is there a validated user need we should meet?"; Alpha: "does our riskiest assumption survive contact with users?").
- Don't skip Discovery because stakeholders "already know the problem" — that certainty is precisely what Discovery exists to test.
- In Live, wire the monitoring set into standing dashboards ([#analytics-and-behavioral-data](#analytics-and-behavioral-data), [#ux-metrics-frameworks](#ux-metrics-frameworks)) with owners, not a launch-week report.

**Related:** [#iterative-design-and-the-double-diamond](#iterative-design-and-the-double-diamond), [#the-research-method-landscape](#the-research-method-landscape), [#accessibility-auditing-in-process](#accessibility-auditing-in-process).

**Sources:** [GOV.UK Service Manual: Agile delivery phases](https://www.gov.uk/service-manual/agile-delivery), [GOV.UK Service Manual: user research](https://www.gov.uk/service-manual/user-research).

### Design sprints

**Definition:** The 5-day (or compressed) structured process (Knapp/GV): Monday map, Tuesday sketch, Wednesday decide, Thursday prototype, Friday test with 5 users — a time-boxed shortcut from big question to user-tested prototype without building anything.

**Reasoning / Evidence:** The sprint's real mechanisms: forced decision cadence (no committee drift), tangibility deadline (prototype by Thursday), and validation before investment (Friday's users). GV's portfolio track record popularized it; the failure literature is equally instructive — sprints on trivial questions, without decision-makers in the room, or whose findings die Monday.

**When to use:**
- Big, risky, contested questions at project inceptions/pivots: new product concepts, major redesign directions, cross-team deadlocks needing shared evidence.
- Decider present, team dedicated (truly full-time for the week), real users booked for Friday.

**When NOT to use / exceptions:**
- Incremental features (routine iteration is cheaper), questions already answerable by existing research, or as a recurring ritual (sprint fatigue devalues the format).
- No follow-through plan = expensive theater; schedule the post-sprint decision meeting *before* the sprint.

**Pros:** Week-boxed learning on quarter-sized risks; alignment through shared observation (stakeholders watching Friday's tests convert better than any report).
**Cons:** Calendar cost (5 people × 5 days); a single 5-user test is directional, not conclusive — treat the output as a validated *hypothesis*.

**Implementation guidance:**
- Compressed 3–4-day variants preserve the spine (map/decide/prototype/test) — protect the *test day* above all.
- Prototype at facade fidelity (Keynote/Figma smoke-and-mirrors) — Thursday builds a movie set, not a product.

**Related:** [#iterative-design-and-the-double-diamond](#iterative-design-and-the-double-diamond), [#usability-testing](#usability-testing).

**Sources:** [Knapp: Sprint](https://www.thesprintbook.com/), [NN/g: Design Sprints 101](https://www.nngroup.com/articles/design-sprints-101/).

### Design critique

**Definition:** Structured peer review of design work: presenter frames goals and questions, critics respond to *whether the design meets the stated goals* — not their taste — via facilitated formats (silent review, round-robin, "I like / I wish / what if") that separate observation from prescription.

**Reasoning / Evidence:** Unstructured feedback defaults to taste wars and HiPPO dominance; structure (goal-anchoring, question-scoping, facilitation) converts critique into design-quality infrastructure. Regular critique also spreads system knowledge and calibrates a team's shared quality bar ([10-design-systems.md](10-design-systems.md) culture layer).

**When to use:**
- Standing cadence (weekly) for work-in-progress; milestone reviews before build; onboarding vehicle for new designers.

**When NOT to use / exceptions:**
- Critique ≠ approval gate (decision rights stay with the owning designer/team unless explicitly a review board); ≠ user data (critique is expert inference — testing trumps it on disputes).

**Pros:** Cheap quality layer, knowledge diffusion, ego-safe conflict channel.
**Cons:** Facilitation-dependent (unfacilitated critique reverts to loudest-voice); time cost creep (scope tightly).

**Implementation guidance:**
- Presenter states: context, goals, constraints, specific questions, fidelity stage ([#prototyping-fidelity-ladder](#prototyping-fidelity-ladder) labeling).
- Critics: observations tied to goals/principles ("this violates [consistency](01-core-principles.md#4-consistency-and-standards) with our other pickers") before suggestions; questions before verdicts.
- Timebox per item (15–20 min); notes captured; unresolved disputes → route to evidence ([#usability-testing](#usability-testing), [#ab-testing](#ab-testing)).
- For repeatable milestone reviews, a shared scoring rubric reduces taste drift ([#design-review-rubrics-and-ai-assisted-review](#design-review-rubrics-and-ai-assisted-review)).

**Related:** [01-core-principles.md](01-core-principles.md), [#heuristic-evaluation](#heuristic-evaluation), [#design-review-rubrics-and-ai-assisted-review](#design-review-rubrics-and-ai-assisted-review).

**Sources:** [NN/g: Design Critiques](https://www.nngroup.com/articles/design-critiques/).

### Design debt

**Definition:** Accumulated UX inconsistency and deferred quality: pattern drift (five button styles), outdated flows beside new ones, abandoned half-migrations, unfixed known usability issues — the experience equivalent of technical debt, compounding user confusion and team velocity loss.

**Reasoning / Evidence:** Debt violates [consistency](01-core-principles.md#4-consistency-and-standards) at scale (users pay per inconsistency in re-learning); teams pay in per-feature design archaeology. Like tech debt, it's invisible in feature-by-feature review and obvious in cross-product audits.

**When to use (debt management):**
- Inventory: periodic UI audits (screenshot every pattern-instance; the wall of shame motivates); severity-rank by user-facing confusion × touch frequency.
- Paydown vehicles: design-system migrations ([10-design-systems.md](10-design-systems.md#versioning-and-deprecation)), boy-scout rule (leave touched screens consistent), dedicated debt budget (10–20% of design/eng capacity is a common working ratio).

**When NOT to use / exceptions:**
- Zero-debt absolutism blocks shipping; debt is a *loan* — legitimate when consciously taken with a paydown note, toxic when silent.

**Pros (managing it):** Compounding-cost interruption; velocity recovery; audit artifacts double as system-adoption evidence.
**Cons:** Invisible-work political problem (mitigate: quantify — support tickets, task-time deltas, migration-cost trends).

**Implementation guidance:**
- Debt register with the same rigor as bugs: item, user impact, spread (instances), paydown cost; review quarterly.
- Prevent at the source: system components as default parts ([10-design-systems.md](10-design-systems.md)), critique catching drift ([#design-critique](#design-critique)).

**Related:** [10-design-systems.md](10-design-systems.md), [01-core-principles.md](01-core-principles.md#4-consistency-and-standards), [#prioritizing-ux-work](#prioritizing-ux-work).

**Sources:** [NN/g: Design Debt](https://www.nngroup.com/articles/design-debt/) (concept treatments vary; also see design-systems literature).

### Prioritizing UX work

**Definition:** Frameworks for choosing among UX investments: **impact/effort matrix** (2×2: quick wins vs big bets vs fill-ins vs money pits), **RICE** (Reach × Impact × Confidence ÷ Effort), severity×frequency for usability-fix backlogs, and opportunity scoring (importance − satisfaction) for feature gaps.

**Reasoning / Evidence:** Unprioritized UX backlogs decay into squeaky-wheel ordering (loudest stakeholder, latest anecdote); scoring forces the [Pareto](01-core-principles.md#pareto-principle) question (reach × severity) and — crucially — the *confidence* term (RICE) makes evidence-quality explicit, rewarding validated problems over opinions.

**When to use:**
- Usability-fix backlogs: severity (blocker→cosmetic) × frequency (analytics-informed — [#analytics-and-behavioral-data](#analytics-and-behavioral-data)) as the sort key; severity scoring should weigh task blockage, error consequence, affected user groups, and accessibility exclusion explicitly, with fix complexity recorded separately as a scheduling input (cheap-fix-high-severity items are the quickest wins; expensive-fix-high-severity items need planning, not deprioritizing).
- Roadmap-level UX bets: RICE with honest confidence scoring (research-backed = high; hunch = 20%).

**When NOT to use / exceptions:**
- Scores are conversation-structurers, not oracles (garbage estimates in, confident garbage out); strategic bets legitimately override scores — visibly, not by fudging inputs.
- Accessibility blockers and data-loss bugs skip the queue (legal/ethical floors — [05-accessibility.md](05-accessibility.md#why-accessibility)).

**Pros:** Defensible ordering, evidence-quality visibility, stakeholder-negotiation medium.
**Cons:** False precision theater; effort estimates are chronically optimistic (use ranges).

**Implementation guidance:**
- Score in groups (calibration through argument); revisit quarterly; track prediction accuracy to calibrate future confidence terms.
- Publish the ranked list + rationale — the transparency is half the value.

**Related:** [#analytics-and-behavioral-data](#analytics-and-behavioral-data), [#design-debt](#design-debt), [01-core-principles.md](01-core-principles.md#pareto-principle).

**Sources:** [Intercom: RICE](https://www.intercom.com/blog/rice-simple-prioritization-for-product-managers/), [NN/g: Prioritization methods](https://www.nngroup.com/articles/prioritization-methods/).

### AI-assisted design and research

**Definition:** The current-generation toolkit: LLM-assisted research analysis (transcript summarization, theme extraction), generative UI drafting (screens/copy variants), synthetic-data prototyping, AI usability-evaluation claims, and "synthetic users" (LLM-simulated participants) — with sharply uneven reliability across these uses.

**Reasoning / Evidence:** Established value: transcription, first-pass coding of qualitative data (with human verification — LLMs miss and hallucinate themes), copy variant generation, boilerplate screen drafting inside a design system. Established danger: **synthetic users** — LLMs generate plausible-sounding *average* opinions reflecting training data, not your users' behavior; they cannot be surprised, don't struggle authentically, and systematically miss the tail behaviors testing exists to find (NN/g and academic evaluations converge: useful for pilot-question rehearsal, invalid as user evidence). AI evaluation tools (auto-heuristic-review) catch pattern-level issues, miss context.

**When to use:**
- Force-multiplier uses: analysis first drafts (human-verified against source), research-plan/discussion-guide drafting, copy exploration ([09-content-ux-writing.md](09-content-ux-writing.md)), design-system-constrained screen drafts, test-task piloting.

**When NOT to use / exceptions:**
- Never as substitute for real-user contact on real decisions; never present synthetic-user output as research findings; never auto-ship AI evaluation verdicts without expert review.
- Data governance: participant data into third-party models needs consent coverage and privacy review ([#research-ethics-and-inclusive-research-operations](#research-ethics-and-inclusive-research-operations)).

**Pros:** Analysis throughput, cheap divergence (variants), lowered documentation drudgery.
**Cons:** Fluent-wrongness (errors read as confident findings); homogenization pull (average-flavored designs); skill-atrophy risk in junior researchers skipping the raw data.

**Implementation guidance:**
- Verification ratio: any AI-coded qualitative claim gets source-quote spot-checks; keep humans on all severity judgments.
- Use AI *between* research touchpoints, never *instead of* them; disclose AI assistance in research reports.
- When AI performs structured design review, use an explicit rubric and human sign-off ([#design-review-rubrics-and-ai-assisted-review](#design-review-rubrics-and-ai-assisted-review)).

**Related:** [#usability-testing](#usability-testing), [#interviews-and-contextual-inquiry](#interviews-and-contextual-inquiry), [20-ai-product-ux.md](20-ai-product-ux.md).

**Sources:** [NN/g: AI for UX work](https://www.nngroup.com/articles/ai-tools-ux-work/), [NN/g: Synthetic Users](https://www.nngroup.com/articles/synthetic-users/).

### Design review rubrics and AI-assisted review

**Definition:** Structured 0–3 scoring rubrics for reviewing UI/UX designs and implementations — a universal UX rubric, a component checklist, a high-risk-flow rubric, and an accessibility rubric — plus a machine-readable output template so that reviews (human or AI-agent) produce comparable, actionable findings rather than free-form opinion.

**Caveat — status of these rubrics:** These rubrics are a working convention, not an industry standard; the criteria derive from established heuristics ([01-core-principles.md](01-core-principles.md#nielsens-10-usability-heuristics)) and WCAG ([05-accessibility.md](05-accessibility.md)), but the scoring bands and format are a practical convention for consistency, not validated instruments like SUS. A single-evaluator review — human or AI — catches pattern-level issues and misses context, exactly as the evaluator-effect evidence predicts for [heuristic evaluation](#heuristic-evaluation) (single evaluators find ~35% of problems); rubric scores are inputs to judgment, not verdicts. And per the warnings in [#ai-assisted-design-and-research](#ai-assisted-design-and-research): AI-generated review verdicts must be reviewed by a human expert before shipping — an agent's confident 3/3 is a claim to verify, not a certification.

**Reasoning / Evidence:** Unstructured review drifts to taste and to whatever the reviewer noticed last ([#design-critique](#design-critique) exists to counter exactly this); a rubric fixes the coverage (every review checks states, recovery, consistency, load, accessibility) and the scale (a "2" means the same thing across reviews and reviewers). For AI-agent review specifically, a rubric plus structured output makes the agent's reasoning auditable — a human can check *which* criterion failed and why, instead of accepting a fluent summary.

**Scoring scale:**

| Score | Meaning |
| --- | --- |
| 0 | Missing or harmful |
| 1 | Present but weak/incomplete |
| 2 | Adequate for normal use |
| 3 | Strong, accessible, resilient |

**Universal UX rubric:**

| Criterion | 0 | 1 | 2 | 3 |
| --- | --- | --- | --- | --- |
| Purpose clarity | No clear purpose | Purpose inferred | Purpose clear | Purpose clear with next step and scope |
| State visibility | State hidden | Some feedback | Key states visible | All states timely, accessible, local |
| User control | No exits/recovery | Some exits | Cancel/back/undo where needed | Robust recovery, drafts, undo, safeguards |
| Consistency | Inconsistent | Mostly consistent | Product/platform consistent | Consistent with documented deviations |
| Error prevention | Error-prone | Some constraints | Prevents common errors | Prevents high-risk errors and supports recovery |
| Cognitive load | Overwhelming | Some grouping | Clear hierarchy | Progressive, contextual, low memory burden |
| Accessibility | Blocks users | Partial | Meets baseline | Strong manual and AT-tested support |

**Component rubric** (pass/fail per item — detail in [11-components-and-overlays.md](11-components-and-overlays.md) and [06-aria-widget-reference.md](06-aria-widget-reference.md)):
- Correct component chosen for the job.
- Required states complete (default/hover/focus/active/disabled/loading/error/empty).
- Accessible name/role/value correct.
- Keyboard/touch behavior correct.
- Content labels clear.
- Error/empty/loading states handled.
- Platform conventions followed.

**High-risk flow rubric** (destructive, financial, irreversible, or high-visibility actions):

| Criterion | Required |
| --- | --- |
| Consequence clarity | User knows cost/effect/audience/data impact |
| Review | User can review before commit |
| Edit | User can correct before commit |
| Recovery | Undo/cancel/support/receipt exists where possible |
| Authentication | Step-up auth if needed, accessible |
| Audit | Confirmation record/log where appropriate |
| Accessibility | Manual keyboard/screen reader checks passed |

**Accessibility review rubric** (manual checks; criteria detail in [05-accessibility.md](05-accessibility.md)):
- Keyboard path complete.
- Focus visible and not obscured.
- Screen reader names/roles/states meaningful.
- Contrast and non-color cues pass.
- Text resize/reflow works.
- Motion reduced where requested.
- Target size and gesture alternatives pass.
- Errors associated and announced.
- Authentication accessible.

**Agent review output template** (machine-readable, one per reviewed artifact):

```md
Overall score: 0-3
Blockers: ...
High-risk issues: ...
Pattern mismatches: ...
Accessibility issues: ...
Recommended changes: ...
Acceptance criteria: ...
Sources/guidance: ...
```

**When to use:**
- Milestone design reviews and definition-of-done gates where consistency across reviewers matters; AI-agent first-pass reviews feeding a human reviewer; onboarding reviewers to a shared quality bar.
- High-risk flow rubric: mandatory for destructive/financial/irreversible actions before ship.

**When NOT to use / exceptions:**
- Not a replacement for [usability testing](#usability-testing) (rubrics inspect against principles; users reveal what principles missed) nor for multi-evaluator [heuristic evaluation](#heuristic-evaluation) on major surfaces.
- Not for early exploratory work — rubric-scoring a sketch punishes the fidelity choice ([#prototyping-fidelity-ladder](#prototyping-fidelity-ladder)); use goal-anchored [critique](#design-critique) instead.
- Don't average away blockers: a 0 on accessibility or high-risk recovery is a gate, whatever the overall mean.

**Pros:** Coverage consistency, comparable scores over time, auditable AI-review output, teaches the quality bar.
**Cons:** Checkbox theater risk (scoring without judgment); rubric criteria lag novel patterns; false authority of numeric scores over small-N expert inspection.

**Implementation guidance:**
- Route rubric use through the executable checklists in [21-agent-checklists.md](21-agent-checklists.md) for agent workflows; keep the rubric versioned alongside the design system ([10-design-systems.md](10-design-systems.md)).
- Require evidence per score ("2 on state visibility because loading state exists but error state missing at step 3") — unevidenced scores are opinions in a costume.
- Human sign-off rule: an accountable human expert reviews every AI-generated verdict before it gates or ships anything; log disagreements to improve the rubric.
- Treat blockers and accessibility-rubric failures as queue-skipping issues ([#prioritizing-ux-work](#prioritizing-ux-work)).

**Verification checklist:**
- [ ] Every criterion scored with cited evidence (screen, state, step).
- [ ] Blockers listed separately, not averaged into the overall score.
- [ ] Accessibility rubric run manually (keyboard + screen reader), not inferred from automated scans alone.
- [ ] AI-generated reviews carry a named human reviewer before acting on them.

**Related:** [#design-critique](#design-critique), [#heuristic-evaluation](#heuristic-evaluation), [#ai-assisted-design-and-research](#ai-assisted-design-and-research), [21-agent-checklists.md](21-agent-checklists.md), [05-accessibility.md](05-accessibility.md).

**Sources:** Working convention synthesized from [Nielsen's heuristics](https://www.nngroup.com/articles/ten-usability-heuristics/) and [WCAG 2.2](https://www.w3.org/TR/WCAG22/); see caveat above on status.

---

## Quick reference

| Method/Topic | One-line guidance | Key numbers |
|---|---|---|
| Method choice | Attitudinal/behavioral × qual/quant; question first | — |
| Research ethics | Consent, accommodations, data protection — every study | No exceptions; proportional apparatus |
| Interviews | Past specifics, not hypotheticals; observe in context | 5–8 per segment; 45–60 min |
| Surveys | Validated scales; pilot; never predict behavior | ≤5 min; balanced Likert |
| Personas | Evidence-based, behavioral, few | 3–5 max |
| JTBD | Situation-progress framing; switch interviews | — |
| Journey maps | Research-fed, pain-owned, data-overlaid | One persona × one goal |
| Usability testing | Watch task attempts; think-aloud; iterate | Success/time/errors + severity |
| 5-user rule | Rounds of ~5, per segment, formative only | L≈31%; ~85% at n=5; n≥20 for metrics |
| Card sort / tree test | Larger samples than usability rounds | Card sort ~15; tree test ~50 |
| Heuristic evaluation | 3–5 independent experts, then aggregate | Single expert ≈35% coverage |
| Cognitive walkthrough | Four questions per step, novice lens | 2–3 critical tasks |
| A/B testing | Pre-registered, powered, guardrailed, no peeking | ≥1–2 weekly cycles |
| Analytics | Funnels + rage-clicks locate; qual explains | Mask PII in replays |
| SUS | Perceived-usability benchmark | Mean 68; 80+ excellent; n≈12 rule of thumb |
| SEQ | Per-task ease pulse | 1–7; benchmark ≈5.5 |
| HEART | Goals→signals→metrics discipline | ≤5 product UX metrics |
| Operational metrics | Recovery, zero-result search, support contacts, a11y defects | Trend per critical flow |
| A11y in process | CI automation + scheduled human passes | Automation ≈30–40% ceiling |
| Prototyping | Cheapest artifact that falsifies the hypothesis | Fidelity steers feedback type |
| Double diamond | Right problem, then right solution | ≥2–3 test cycles budgeted |
| GDS phases | Discovery→Alpha→Beta→Live; evidence gates | Discovery can conclude "don't build" |
| Design sprints | Big contested questions; decider present | 5 days; 5 users Friday |
| Critique | Goal-anchored, facilitated, observation-first | 15–20 min per item |
| Review rubrics | Scored coverage; human signs off AI verdicts | 0–3 scale; blockers never averaged |
| Design debt | Register, budget, boy-scout rule | 10–20% capacity paydown |
| Prioritization | Severity×frequency; RICE with honest confidence | A11y blockers skip the queue |
| AI assistance | Drafts and analysis help; synthetic users ≠ users | Verify against source data |
