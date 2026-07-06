# Patterns to Avoid & Controversial Patterns

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file catalogs what *not* to do and what to do only carefully: deceptive/dark patterns (with legal status, detection heuristics for automated/agent review, and ethical alternatives), the deceptive–neutral–persuasive spectrum, common non-deceptive anti-patterns, trust/privacy/security notes, and controversial patterns with balanced trade-off analysis. Positive doctrine lives in the rest of the reference; entries here link back to the rules they violate. Legal context: the EU DSA (Art. 25 bans dark patterns on platforms), GDPR (consent manipulation), EU consumer-protection directives, and US FTC enforcement (e.g., the Epic/Fortnite settlement) have converted much of this catalog from "bad practice" into legal liability. Pattern names follow the current deceptive.design taxonomy (revised 2023), with Brignull's original 2010 names in parentheses for recognizability.

## Contents

- [The deceptive–neutral–persuasive spectrum](#the-deceptive-neutral-persuasive-spectrum)
- [Dark patterns catalog](#dark-patterns-catalog)
  - [Confirmshaming](#confirmshaming)
  - [Hard to cancel (roach motel)](#hard-to-cancel-roach-motel)
  - [Hidden subscription (forced continuity)](#hidden-subscription-forced-continuity)
  - [Hidden costs and drip pricing](#hidden-costs-and-drip-pricing)
  - [Bait and switch / trick wording](#bait-and-switch--trick-wording)
  - [Disguised ads](#disguised-ads)
  - [Fake scarcity, fake urgency, and fake social proof](#fake-scarcity-fake-urgency-and-fake-social-proof)
  - [Nagging](#nagging)
  - [Obstruction](#obstruction)
  - [Preselection](#preselection)
  - [Forced action (privacy Zuckering)](#forced-action-privacy-zuckering)
  - [Visual interference (misdirection)](#visual-interference-misdirection)
  - [Sneaking (sneak into basket)](#sneaking-sneak-into-basket)
  - [Friend spam](#friend-spam)
  - [Comparison prevention (hard-to-compare pricing)](#comparison-prevention-hard-to-compare-pricing)
- [Common anti-patterns](#common-anti-patterns)
  - [Mystery meat navigation](#mystery-meat-navigation)
  - [Auto-rotating carousels](#auto-rotating-carousels)
  - [Splash screens](#splash-screens)
  - [Modal overuse](#modal-overuse)
  - [Forced tutorial tours](#forced-tutorial-tours)
  - [Placeholder-as-label](#placeholder-as-label)
  - [Disabled submit without explanation](#disabled-submit-without-explanation)
  - [Custom scrollbars and scrolljacking](#custom-scrollbars-and-scrolljacking)
  - [Breaking browser back](#breaking-browser-back)
  - [Unsolicited chat popups](#unsolicited-chat-popups)
  - [Notification permission on page load](#notification-permission-on-page-load)
  - [App interstitials](#app-interstitials)
  - [Form data loss](#form-data-loss)
  - [Aggressive session timeout](#aggressive-session-timeout)
  - [Low-contrast text](#low-contrast-text)
  - [Icon-only interfaces without labels](#icon-only-interfaces-without-labels)
  - [Gesture-only features](#gesture-only-features)
  - [Fake loading](#fake-loading)
  - [Hover-only disclosure on touch](#hover-only-disclosure-on-touch)
  - [PDF for web content](#pdf-for-web-content)
- [Trust, privacy, and security notes](#trust-privacy-and-security-notes)
  - [Security UX](#security-ux)
  - [Privacy by design](#privacy-by-design)
  - [Public-service and institutional trust](#public-service-and-institutional-trust)
- [Controversial patterns](#controversial-patterns)
  - [Infinite scroll](#infinite-scroll)
  - [Hamburger menu](#hamburger-menu)
  - [Carousels (user-controlled)](#carousels-user-controlled)
  - [Gamification and streaks](#gamification-and-streaks)
  - [Personalization and filter bubbles](#personalization-and-filter-bubbles)
  - [Autoplay video](#autoplay-video)
  - [Floating action buttons](#floating-action-buttons)
  - [AI chatbots as primary support](#ai-chatbots-as-primary-support)
- [Quick reference](#quick-reference)

---

## The deceptive-neutral-persuasive spectrum

**Definition:** The ethical scale for influence in design: **persuasive** (transparent influence serving informed user choice: clear value propositions, honest defaults, real social proof) → **neutral** (no engineered influence) → **manipulative/deceptive** (influence exploiting cognitive bias against the user's interest: hidden information, engineered confusion, coerced consent). The line: would the design survive full user awareness of its mechanism — and does it serve the user's goals or only extract from them?

**Reasoning / Evidence:** All design influences behavior ([defaults](02-cognitive-foundations.md#cognitive-biases-in-ux) prove neutrality is a myth); the ethical question is direction and transparency. Two workable tests: the **disclosure test** (explain the mechanism to the user — do they feel served or played?) and the **regret test** (do users, fully informed later, regret the choice the design steered?). Regulators have operationalized versions of this: FTC's "unfair or deceptive," DSA's "materially distorting the ability to make free and informed decisions."

**When to use (persuasion legitimately):**
- Honest recommendation (best-for-most defaults, "most popular" when true), real urgency (actual stock, actual deadlines), friction *for* the user (typed confirmation before irreversible deletion — [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).

**Detection heuristic (for design review):** For each influence element ask: (1) Is the claim verifiable and true? (2) Is the user's alternative equally visible and easy? (3) Who benefits if the user acts without thinking — the user or the business? Two-or-more business-only answers = investigate as dark pattern.

**Four verification questions (apply to any contested design):**
1. Would the user consent if fully informed of all consequences?
2. Does the design benefit the user, or only the business?
3. Is declining as easy as accepting where the choice should be symmetrical?
4. Would you defend this design publicly — to a journalist, a regulator, or the user's face?

Any "no" is a signal to redesign; two or more is a dark pattern in progress.

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux), every entry below.

**Sources:** [Deceptive Design (Brignull)](https://www.deceptive.design/), [FTC: Bringing Dark Patterns to Light](https://www.ftc.gov/reports/bringing-dark-patterns-light), [EU DSA Art. 25](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32022R2065).

---

## Dark patterns catalog

Format note: dark-pattern entries use **Why it's harmful**, **Detection heuristic** (how an agent reviewing a design spots it), **Ethical alternative**, and **Legal status** in place of the standard pros/cons — there is no legitimate "when to use." Names follow the current deceptive.design taxonomy; legacy Brignull names appear in parentheses.

### Confirmshaming

**Definition:** Guilt-tripping decline options: "No thanks, I hate saving money."

**Why it's harmful:** Weaponizes emotion against free choice; degrades brand voice; regulators treat emotionally-loaded consent steering as manipulation (GDPR consent must be freely given).

**Detection heuristic:** Decline/opt-out copy contains self-deprecation, loss-language, or judgment ("I prefer paying full price") instead of neutral "No thanks."

**Ethical alternative:** Symmetric, neutral choices: "Yes, subscribe" / "No thanks." State the benefit in the *offer*, not shame in the refusal ([09-content-ux-writing.md](09-content-ux-writing.md#voice-and-tone)).

**Legal status:** Emotionally-steered consent risks invalidity under GDPR's freely-given standard; EDPB deceptive-design guidelines and FTC dark-pattern reports both name it.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

### Hard to cancel (roach motel)

**Definition:** Easy to get in, hard to get out: one-click subscribe, cancellation buried behind phone calls, retention mazes, and chat-agent gauntlets.

**Why it's harmful:** Converts consent into capture; regulators on both sides of the Atlantic increasingly demand cancellation parity (cancel as easily as sign up); Amazon Prime's cancellation flow (internally code-named "Iliad") drew FTC action on exactly this.

**Detection heuristic:** Count steps/modalities: sign-up steps vs cancellation steps; any channel switch (web signup but phone-only cancel) is a red flag; count retention interstitials before the cancel confirm.

**Ethical alternative:** Cancellation parity: same channel, ≤ same steps; one clearly-labeled path (Settings → Subscription → Cancel); one *dismissible* save-offer maximum; confirm + immediate email receipt. Easy exits build return-willingness ([Peak-End](01-core-principles.md#peak-end-rule) — the end is what they remember).

**Legal status:** In flux but hostile to the pattern. The FTC's Click-to-Cancel rule (announced October 2024) was **vacated by the Eighth Circuit Court of Appeals in July 2025** before taking effect — but FTC enforcement against subscription traps continues under FTC Act §5 (unfair/deceptive practices) and ROSCA, as in the Amazon Prime "Iliad flow" case. EU consumer law and several US state auto-renewal statutes independently require easy cancellation. Do not design to the vacated rule's absence: the enforcement basis survives it.

**Sources:** [FTC Click-to-Cancel announcement (2024; rule later vacated)](https://www.ftc.gov/news-events/news/press-releases/2024/10/federal-trade-commission-announces-final-click-cancel-rule-making-it-easier-consumers-end-recurring), [deceptive.design taxonomy](https://www.deceptive.design/types).

### Hidden subscription (forced continuity)

**Definition:** Silent trial-to-paid conversion: card required upfront, no expiry reminder, billing begins invisibly — plus its sibling, subscription auto-renewal without notice.

**Why it's harmful:** Monetizes forgetting ([post-completion error](02-cognitive-foundations.md#post-completion-errors) exploitation); chargebacks and support burden offset the captured revenue; trust damage compounds when the charge is discovered.

**Detection heuristic:** Trial flow demands payment credentials without a pre-charge notification commitment; no reminder email exists in the flow map; cancellation not available *during* trial.

**Ethical alternative:** Trial without card where viable; if card required: reminder 3–7 days pre-charge (email + in-app), one-tap cancel from the reminder, and easy in-trial cancellation. Netflix-style "your trial ends on X" transparency measurably preserves trust.

**Legal status:** ROSCA and FTC negative-option enforcement apply in the US; EU consumer law and multiple US states require renewal reminders and explicit consent to recurring charges.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types), [FTC negative-option enforcement](https://www.ftc.gov/reports/bringing-dark-patterns-light).

### Hidden costs and drip pricing

**Definition:** Advertising a price that inflates through the funnel: fees, "service charges," and shipping revealed at the last step.

**Why it's harmful:** Sunk-cost exploitation (users deep in checkout tolerate the reveal — [02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux)); measured trust destruction and cart abandonment (top abandonment cause in Baymard's checkout research).

**Detection heuristic:** Compare first-shown price vs final total across the funnel; any mandatory fee appearing after step 1 flags; "+ taxes and fees" without amounts flags.

**Ethical alternative:** All-in pricing upfront (or full fee disclosure at first price display); shipping estimators early; itemized-but-complete totals from the cart onward ([07-forms-and-input.md](07-forms-and-input.md), [02-cognitive-foundations.md](02-cognitive-foundations.md#progressive-disclosure) — never hide *decision* information).

**Legal status:** Banned/regulated in growing jurisdictions: US junk-fee rules and EU price-indication law both target mandatory fees revealed late.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types), [Baymard: checkout abandonment](https://baymard.com/lists/cart-abandonment-rate).

### Bait and switch / trick wording

**Definition:** Advertising X, delivering Y: the promoted deal that's "unavailable" (upsell offered), the button that does something other than its label, the classic Windows-10-era close-box-means-consent dialog. The current taxonomy covers much of this under **trick wording** (language engineered to mislead) and the **sneaking** family.

**Why it's harmful:** Direct deception; destroys the label-action contract that all interface trust rests on ([affordances/signifiers](01-core-principles.md#signifiers) weaponized).

**Detection heuristic:** Does every control do exactly what its label says? Does advertised inventory/pricing exist and remain available through the flow? Does dismissing/closing anything ever equal consent? Read every consent sentence for double negatives and scope-switching mid-sentence.

**Ethical alternative:** Labels = actions, always; out-of-stock promotions clearly marked before click-through; close means close; consent copy in plain language ([09-content-ux-writing.md](09-content-ux-writing.md#plain-language)).

**Legal status:** Textbook FTC "deceptive practice"; misleading commercial communication violates the EU Unfair Commercial Practices Directive.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

### Disguised ads

**Definition:** Advertising camouflaged as content, navigation, or UI: fake download buttons on file-hosting sites, "sponsored" results styled identically to organic, ad-as-article without disclosure.

**Why it's harmful:** Exploits learned trust in UI conventions ([Jakob's Law](01-core-principles.md#jakobs-law) inverted); fake-button ads on download flows cause real malware harm.

**Detection heuristic:** Screenshot test: can a user distinguish every paid unit from content/chrome in 2 seconds? Are there buttons styled like the primary action that lead off-task? Is "Sponsored/Ad" labeling present, legible (not 8px grey-on-grey), and adjacent?

**Ethical alternative:** Visually distinct ad containers, clear conventional labels, and never ad units mimicking the page's primary action ([02-cognitive-foundations.md](02-cognitive-foundations.md#selective-attention-and-banner-blindness) — honest styling costs clicks and keeps trust).

**Legal status:** FTC native-advertising rules and EU transparency obligations (including the DSA) require clear, proximate ad labeling.

**Sources:** [FTC: Native Advertising guidance](https://www.ftc.gov/business-guidance/resources/native-advertising-guide-businesses), [deceptive.design taxonomy](https://www.deceptive.design/types).

### Fake scarcity, fake urgency, and fake social proof

**Definition:** Manufactured pressure — three sibling types in the current taxonomy: countdown timers that reset (fake urgency), "3 left!" without inventory basis (fake scarcity), "12 people are viewing this" fabrications (fake social proof).

**Why it's harmful:** Exploits scarcity bias and time pressure to suppress deliberation ([System 2](02-cognitive-foundations.md#system-1-and-system-2) bypass); discovered fakery (timer resets on refresh) permanently marks the brand as a liar.

**Detection heuristic:** Refresh the page — does the timer reset? Is the scarcity number generated client-side or from inventory? Do "expiring" deals reappear? Can the claim be audited against a data source?

**Ethical alternative:** Real deadlines and real inventory only, with basis shown where possible ("sale ends Sunday"; "2 in stock" from actual inventory); no urgency theater on evergreen offers ([Parkinson's Law](01-core-principles.md#parkinsons-law) honest-expectation framing).

**Legal status:** *False* claims are illegal in the EU (Unfair Commercial Practices blacklist) and FTC-actionable in the US as deceptive practices.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types), [EU UCPD blacklist](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32005L0029).

### Nagging

**Definition:** Repeated interruption to extract a yes: rating prompts every session, notification re-asks after refusal, upsell modals on every launch, "we miss you" pushes.

**Why it's harmful:** Attrition-warfare against stated preference; interruption cost compounds ([02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost)); habituated dismissal then blinds users to *legitimate* dialogs ([habituation](02-cognitive-foundations.md#habituation)).

**Detection heuristic:** Is every prompt's refusal *persisted*? What's the re-ask interval per prompt type? Does any prompt lack a "don't ask again"? Count interrupts per session for a returning user who has said no to everything.

**Ethical alternative:** Ask once, in context, after positive signal ([15-platform-mobile.md](15-platform-mobile.md#notifications-ux) permission priming); persist refusals (re-ask ≥ months, if ever); passive availability (a settings toggle, a banner-free "rate us" menu item) over interruption.

**Legal status:** Named in EDPB deceptive-design guidelines; platforms police re-prompt frequency directly (iOS rating API caps at 3 prompts/year).

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

### Obstruction

**Definition:** Deliberate friction on disfavored choices: privacy settings six menus deep, decline flows with extra confirmations, data-export that "takes up to 30 days," comparison-blocking.

**Why it's harmful:** Weaponized [Tesler's Law](01-core-principles.md#teslers-law) — complexity pushed *onto* users selectively to steer them.

**Detection heuristic:** Click-count asymmetry audit: steps to accept vs decline, subscribe vs cancel, share-data vs refuse; any deliberate delay, channel switch, or repeated confirmation appearing *only* on the business-disfavored path.

**Ethical alternative:** Symmetric friction: equal steps for opposite choices; friction budgeted by *user* risk (irreversibility), never by business preference ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).

**Legal status:** DSA Art. 25 names it directly for platforms; GDPR Art. 7(3) requires withdrawal of consent to be as easy as giving it.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types), [EDPB dark-pattern guidelines](https://www.edpb.europa.eu/our-work-tools/our-documents/guidelines/guidelines-032022-deceptive-design-patterns-social-media_en).

### Preselection

**Definition:** Business-favoring options pre-checked: newsletter boxes, add-on insurance, larger plans, data-sharing toggles defaulted on.

**Why it's harmful:** Exploits default bias ([02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux)) — defaults are among the most powerful steering mechanisms in design, though the magnitude of the default effect varies enormously by domain and stakes; the pattern works precisely because most users don't revisit pre-made choices.

**Detection heuristic:** Inventory every default: does any pre-checked/pre-selected state cost money, share data, or create obligations? Would the business accept the opposite default?

**Ethical alternative:** User-interest defaults only ([07-forms-and-input.md](07-forms-and-input.md#defaults-and-prefill)); consent and purchases always start unchecked; recommend visibly instead ("Most popular" badge on an *unselected* option).

**Legal status:** Legally void for GDPR consent — pre-ticked boxes explicitly invalid per CJEU *Planet49* (C-673/17); pre-selected paid extras banned in EU consumer law.

**Sources:** [CJEU Planet49 (C-673/17)](https://curia.europa.eu/juris/document/document.jsf?docid=218462), [deceptive.design taxonomy](https://www.deceptive.design/types).

### Forced action (privacy Zuckering)

**Definition:** Coercing users into disclosure or unwanted actions as the price of proceeding: bundled consent, confusing privacy controls, "personalization" framing for surveillance, defaults maximizing collection. The legacy name "privacy Zuckering" covers the data-oversharing variant; the current taxonomy files coerced disclosure under **forced action**.

**Why it's harmful:** Consent theater over informed choice; EDPB guidelines catalog the sub-patterns (overloading, skipping, stirring, hindering).

**Detection heuristic:** Can a user grant service-essential permissions *without* marketing/tracking ones (granularity)? Is the data use explained at the point of collection in plain language ([09-content-ux-writing.md](09-content-ux-writing.md#plain-language))? Are privacy-protective choices equally prominent?

**Ethical alternative:** Granular, plain-language, point-of-need consent; privacy-protective defaults; data-use honesty ("for delivery updates only" — [07-forms-and-input.md](07-forms-and-input.md#form-length-reduction)); collect less ([Tesler](01-core-principles.md#teslers-law): absorb the cost of not knowing).

**Legal status:** GDPR purpose-limitation and granular-consent requirements target it directly; bundled take-it-or-leave-it consent is presumptively not freely given.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types), [EDPB deceptive-design guidelines](https://www.edpb.europa.eu/our-work-tools/our-documents/guidelines/guidelines-032022-deceptive-design-patterns-social-media_en).

### Visual interference (misdirection)

**Definition:** Visual hierarchy weaponized: the business-preferred button huge and green, the user's alternative a grey whisper; attention steered away from unfavorable information.

**Why it's harmful:** Corrupts the [visual hierarchy](03-visual-design.md#visual-hierarchy) contract (prominence = importance *to the user*).

**Detection heuristic:** For each decision point: are the options' visual weights proportionate to *user* relevance or to business preference? Squint test on dialogs — does the safe/decline option survive? Equal-prominence check on all consent surfaces.

**Ethical alternative:** Prominence follows user importance; equal choices get equal weight (consent — [13-platform-web.md](13-platform-web.md#consent-and-cookie-ux)); primary-button status for the *likely-intended* action, not the extractive one ([Von Restorff](01-core-principles.md#von-restorff-effect) used for the user).

**Legal status:** DSA/EDPB name asymmetric choice-presentation explicitly; cookie-banner button asymmetry has drawn national-DPA fines.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

### Sneaking (sneak into basket)

**Definition:** Items/services added to the cart without explicit user action — bundled extras, "protection plans," donation add-ons appearing pre-added. The **sneaking** family also covers information smuggled past the user (undisclosed material terms).

**Why it's harmful:** Bypasses consent entirely (beyond nudging into theft-adjacent territory); charge-dispute and trust costs.

**Detection heuristic:** Diff the cart against explicit add-to-cart actions; anything present without a user add = violation. Audit bundles: is the base product purchasable alone at the shown price?

**Ethical alternative:** Opt-in add-ons at decision points, unselected ([#preselection](#preselection)); order review step itemizing everything with removal one click ([07-forms-and-input.md](07-forms-and-input.md#multi-step-vs-single-page-forms) review pattern).

**Legal status:** Explicitly blacklisted in EU consumer law (inertia selling / unrequested additional payments); FTC-actionable in the US.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

### Friend spam

**Definition:** Harvesting contacts under one pretext ("find your friends") then messaging them as the user ("X invited you!") — LinkedIn's historic version drove a class-action settlement. (Legacy Brignull name; the current taxonomy treats it under forced action and address-book abuse.)

**Why it's harmful:** Impersonation + consent-scope violation; burns the user's *social* capital, not just their trust; contact-access permissions now stricter on both mobile platforms partly in response.

**Detection heuristic:** Trace every contact-permission grant to its messaging consequences: does anything send in the user's name without per-message (or clearly per-batch) consent? Is the permission's stated purpose the actual use?

**Ethical alternative:** Explicit preview-and-confirm before anything sends; scoped asks ("find who's already here" ≠ "invite everyone"); user-controlled invite selection.

**Legal status:** Litigation-tested (LinkedIn *Add Connections* settlement); purpose-mismatched contact harvesting also raises GDPR purpose-limitation issues.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

### Comparison prevention (hard-to-compare pricing)

**Definition:** Deliberately incommensurable plans/prices: different units across tiers, feature matrices with strategic gaps, bundle-only pricing obscuring per-item cost.

**Why it's harmful:** Comparison-blocking exploits [choice overload](02-cognitive-foundations.md#decision-fatigue-and-choice-overload) to push users toward satisficing on the promoted option.

**Detection heuristic:** Can plans be compared on a single consistent unit (per user/month, per GB)? Does the comparison table omit rows where a cheaper tier wins? Is the promoted tier's advantage real on the user's likely dimensions?

**Ethical alternative:** Consistent units, complete comparison tables, honest recommendation logic ("best if you…" mapping to needs — [02-cognitive-foundations.md](02-cognitive-foundations.md#decision-fatigue-and-choice-overload) curation done *for* the user).

**Legal status:** Unit-price transparency is mandated in retail contexts for exactly this reason; deliberate comparison prevention feeds unfair-practice analysis in the EU.

**Sources:** [deceptive.design taxonomy](https://www.deceptive.design/types).

---

## Common anti-patterns

Non-deceptive but harmful: these fail users through negligence or fashion rather than intent. Standard template resumes (compressed).

### Mystery meat navigation

**Definition:** Navigation requiring interaction to reveal meaning: unlabeled icons, hover-to-discover menus, cryptic labels. **Why harmful:** Violates [recognition](01-core-principles.md#6-recognition-rather-than-recall) and [signifiers](01-core-principles.md#signifiers); every destination becomes a guess. **Exceptions:** Truly universal icons ([03-visual-design.md](03-visual-design.md#iconography)). **Fix:** Icon+label; test-verified icon recognition ≥80% or label it. **Related:** [08-navigation-ia.md](08-navigation-ia.md). **Sources:** [NN/g: Icon Usability](https://www.nngroup.com/articles/icon-usability/).

### Auto-rotating carousels

**Definition:** Hero banners cycling automatically. **Why harmful:** Erik Runyon's Notre Dame measurements: only ~1% of visitors clicked a carousel slide at all, and ~84–89% of those clicks were on slide 1 — content beyond the first slide is effectively invisible; auto-motion triggers [banner blindness](02-cognitive-foundations.md#selective-attention-and-banner-blindness), frustrates reading-in-progress (WCAG 2.2.2 requires pause), and exists mainly to defer stakeholder priority fights. **Exceptions:** [User-controlled carousels](#carousels-user-controlled) for browsing content. **Fix:** One prioritized hero message; static grid for multiple promotions. **Sources:** [Erik Runyon: Carousel Stats](https://erikrunyon.com/2013/01/carousel-stats/), [NN/g: Auto-Forwarding Carousels](https://www.nngroup.com/articles/auto-forwarding/).

### Splash screens

**Definition:** Branded intro screens/animations delaying content beyond technical necessity. **Why harmful:** Pure wait-cost ([response limits](04-interaction-design.md#feedback-and-response-time-thresholds)); brand goodwill burns at ~1s+. **Exceptions:** OS-mandated launch frames (make them content-skeletons); unavoidable heavy init (show progress, not logos). **Fix:** Launch to content ≤2s ([15-platform-mobile.md](15-platform-mobile.md#performance-perception-on-mobile)); skeleton over logo theater. **Sources:** [NN/g: splash screens](https://www.nngroup.com/articles/mobile-first-impressions/).

### Modal overuse

**Definition:** Blocking dialogs for non-blocking content: newsletter modals on arrival, rating asks mid-task, announcements as interstitials. **Why harmful:** Interruption cost ([02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost)), habituated dismissal ([habituation](02-cognitive-foundations.md#habituation)), and entry-modal bounce; modality is for decisions that must block ([14-platform-desktop.md](14-platform-desktop.md#modality)). **Exceptions:** Genuine blocking decisions (irreversible confirms, required auth). **Fix:** Escalation ladder — badge → banner → toast before any modal ([04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)); never modal-on-load. **Sources:** [NN/g: Modal & Nonmodal](https://www.nngroup.com/articles/modal-nonmodal-dialog/).

### Forced tutorial tours

**Definition:** Mandatory multi-step walkthroughs before product use. **Why harmful:** Front-loaded instructions exceed working memory ([02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking)), get skipped (measured skip rates dominate), and delay time-to-value ([15-platform-mobile.md](15-platform-mobile.md#onboarding)). **Exceptions:** Genuinely novel interaction paradigms (one *interactive* first task, not slides). **Fix:** Do-first with contextual hints at moment of need; optional help center; guided empty states ([04-interaction-design.md](04-interaction-design.md#empty-states)). **Sources:** [NN/g: Onboarding Tutorials](https://www.nngroup.com/articles/onboarding-tutorials/).

### Placeholder-as-label

**Definition:** Placeholder text as the only field label. **Why harmful:** Fully documented failure set — see [07-forms-and-input.md](07-forms-and-input.md#labels-vs-placeholders). **Fix:** Persistent visible labels; placeholders only as supplementary examples. **Sources:** [NN/g: Placeholders](https://www.nngroup.com/articles/form-design-placeholders/).

### Disabled submit without explanation

**Definition:** Greyed-out submit buttons with no indication of what's missing. **Why harmful:** Dead-end guessing game ([error prevention](01-core-principles.md#5-error-prevention) inverted — the "prevention" hides the requirement); worst for AT and cognitive-load-sensitive users. **Exceptions:** Momentary in-flight disabling (with spinner) is fine ([04-interaction-design.md](04-interaction-design.md#component-states)). **Fix:** Keep submit enabled → validate on click → error summary + field messages ([07-forms-and-input.md](07-forms-and-input.md#error-message-design)); or disabled *with* adjacent explanation of what remains. **Sources:** [GOV.UK: never disable submit](https://design-system.service.gov.uk/components/button/).

### Custom scrollbars and scrolljacking

**Definition:** Overridden scroll speed/direction/physics; hidden or restyled-beyond-recognition scrollbars; forced section snapping. **Why harmful:** Hijacks the most practiced interaction in computing ([04-interaction-design.md](04-interaction-design.md#scroll-behaviors)); breaks keyboard scrolling, find-in-page positioning, and motion tolerance ([05-accessibility.md](05-accessibility.md#reduced-motion)). **Exceptions:** Subtle scrollbar *styling* (thin, still visible/grabbable) is acceptable; media scrub areas own their gestures. **Fix:** Native scroll physics always; CSS scroll-snap only on carousels/galleries, interruptible. **Sources:** [NN/g: Scrolljacking 101](https://www.nngroup.com/articles/scrolljacking-101/).

### Breaking browser back

**Definition:** Back button trapping (redirect stacks), state loss on back, or back exiting instead of un-stepping. **Why harmful:** Breaks the web's most-used control ([13-platform-web.md](13-platform-web.md#browser-conventions)); historically a deliberate ad-page trap, now mostly SPA negligence. **Fix:** History discipline per perceived "place"; the 5-deep back audit; Android back-stack curation ([15-platform-mobile.md](15-platform-mobile.md#back-behavior)). **Sources:** [NN/g: back button expectations](https://www.nngroup.com/articles/recalling-vs-recognizing-browser-back/).

### Unsolicited chat popups

**Definition:** Auto-opening chat widgets ("Hi! Need help?") with sound/motion/bounce, often minutes into reading. **Why harmful:** Interruption without intent signal; fakes human attention (bait-adjacent when it's a bot pretending); covers content on mobile. **Exceptions:** Contextual proactive chat on *struggle signals* (repeated errors, long dwell on help pages) tests defensibly. **Fix:** Passive availability (visible button, no auto-open, no sound); open on user intent. **Sources:** [NN/g: chat UX](https://www.nngroup.com/articles/chat-ux/).

### Notification permission on page load

**Definition:** Browser/OS permission prompts fired on arrival. **Why harmful:** Single-digit grant rates, permanent-denial lockout, browser-level suppression of abusers — fully covered in [13-platform-web.md](13-platform-web.md#notification-permission-ux). **Fix:** Two-step contextual priming after demonstrated value. **Sources:** [web.dev: permission UX](https://web.dev/articles/push-notifications-permissions-ux).

### App interstitials

**Definition:** Full-screen "open in app" walls over mobile-web content. **Why harmful:** Blocks the content the user came for; measured engagement loss (Google's own case study on its +1 interstitial); search-ranking penalty for intrusive interstitials; the app-vs-web decision belongs to the user ([15-platform-mobile.md](15-platform-mobile.md#mobile-web-vs-native-app)). **Exceptions:** Legally required gates (age verification); small dismissible smart banners. **Fix:** Smart banner (platform-native, small, dismissible); feature-parity web so the push is unnecessary. **Sources:** [Google: avoid intrusive interstitials](https://developers.google.com/search/docs/appearance/avoid-intrusive-interstitials).

### Form data loss

**Definition:** User input destroyed by validation errors, session expiry, navigation, or crashes. **Why harmful:** The highest-severity form failure; teaches lasting distrust ([learned helplessness](02-cognitive-foundations.md#learned-helplessness-and-error-tolerant-design)). **Fix:** Preserve input through every error round-trip ([07-forms-and-input.md](07-forms-and-input.md#error-message-design)); autosave ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)); session-expiry state restoration ([13-platform-web.md](13-platform-web.md#session-expiry-handling)). **Sources:** [GOV.UK validation pattern](https://design-system.service.gov.uk/patterns/validation/).

### Aggressive session timeout

**Definition:** Short expiry without warning, discarding state. **Why harmful:** WCAG 2.2.1 violation posture + data loss; punishes exactly the deliberate/slow users who need more time ([05-accessibility.md](05-accessibility.md#cognitive-accessibility)). **Exceptions:** High-security contexts justify short sessions — never silent, never state-destroying (see [Security UX](#security-ux)). **Fix:** ≥60s warning + one-click extend + draft survival ([13-platform-web.md](13-platform-web.md#session-expiry-handling)). **Sources:** [WCAG 2.2.1](https://www.w3.org/WAI/WCAG22/Understanding/timing-adjustable.html).

### Low-contrast text

**Definition:** Fashionable grey-on-grey, thin-weight, sub-4.5:1 text. **Why harmful:** The most common accessibility failure on the web (WebAIM Million ~80% of pages); excludes low-vision users and degrades everyone in glare/night conditions ([03-visual-design.md](03-visual-design.md#contrast)). **Fix:** Token-level contrast enforcement; #767676 floor on white; de-emphasize with weight/size, not just lightness. **Sources:** [WebAIM Million](https://webaim.org/projects/million/).

### Icon-only interfaces without labels

**Definition:** Toolbars/navigation of unlabeled glyphs. **Why harmful:** Recognition failure for all but ~5 universal icons ([03-visual-design.md](03-visual-design.md#iconography)); voice-control and SR failures when accessible names are missing too. **Exceptions:** Space-critical toolbars with tooltips + aria-labels + universal glyphs. **Fix:** Icon+label default; measured recognition before any label removal. **Sources:** [NN/g: Icon Usability](https://www.nngroup.com/articles/icon-usability/).

### Gesture-only features

**Definition:** Functionality reachable only by swipe/long-press/multi-touch with no visible path. **Why harmful:** Invisible ([signifiers](01-core-principles.md#signifiers) absent), motor-exclusionary (WCAG 2.5.1/2.5.7 — [05-accessibility.md](05-accessibility.md#touchpointer-accessibility)), and undiscoverable-by-design. **Fix:** Gesture-as-accelerator rule: every gesture has a visible equivalent ([04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures)). **Sources:** [WCAG 2.5.1](https://www.w3.org/WAI/WCAG22/Understanding/pointer-gestures.html).

### Fake loading

**Definition:** Decorative delays and fabricated progress: spinners over instant results, progress theater unconnected to work. **Why harmful:** Wastes user time to manufacture perception — deception by latency. **The documented exception:** *labor illusion* — brief deliberate pacing with honest work-description ("checking 40 insurers…") measurably increases trust in results *when work is real*; security contexts sometimes pace to imply thoroughness. Keep any such pacing ≤1–2s, truthful about the work, and never on repeated operations. **Fix:** Show real speed; if pacing for comprehension (flash-of-change blindness — [02-cognitive-foundations.md](02-cognitive-foundations.md#change-blindness)), use motion, not delay. **Sources:** [Buell & Norton: labor illusion](https://www.hbs.edu/ris/Publication%20Files/11-055.pdf).

### Hover-only disclosure on touch

**Definition:** Hover-revealed actions/info shipped to touch devices where hover doesn't exist. **Why harmful & fix:** Fully covered in [04-interaction-design.md](04-interaction-design.md#hover-dependence) — hover accelerates, never gates. **Sources:** [WCAG 1.4.13](https://www.w3.org/WAI/WCAG22/Understanding/content-on-hover-or-focus.html).

### PDF for web content

**Definition:** Publishing routine web content (forms, guides, policies) as PDFs. **Why harmful:** NN/g's long-running finding: PDFs break reading flow (download, viewer context switch), responsiveness (fixed layout on phones), accessibility (untagged PDFs are SR-hostile), search UX, and updateability. **Exceptions:** Genuinely print-destined artifacts (tickets, certificates, regulated forms) — offered *alongside* HTML. **Fix:** HTML-first content; PDF as secondary export. **Sources:** [NN/g: PDF — Unfit for Human Consumption](https://www.nngroup.com/articles/pdf-unfit-for-human-consumption/).

---

## Trust, privacy, and security notes

The positive counterpart to the catalog above: trust-preserving practices that frequently get sacrificed to the same pressures that produce dark patterns.

### Security UX

- **Friction proportional to risk:** budget security friction by consequence, not uniformly — step-up authentication (re-auth, second factor) for sensitive actions (payout changes, data export, deletion) rather than blanket hostility on every login. See [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux) for the full auth doctrine and [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md) for checkout/account contexts.
- **Make secure behavior the easiest path:** support password managers and paste; offer accessible alternatives to cognitive-function tests (WCAG 3.3.8).
- **Don't train users to click through warnings:** every unnecessary or crying-wolf security dialog erodes response to real ones ([habituation](02-cognitive-foundations.md#habituation)). Reserve interruptive warnings for genuine risk, explain security events in plain language, and never make "ignore" the practiced path.

### Privacy by design

- **Data minimization, visibly:** ask for the minimum data needed and let users *see* that you did — explain why each sensitive field exists at the point of collection ([07-forms-and-input.md](07-forms-and-input.md#form-length-reduction)).
- **Clear retention:** state how long data is kept in plain language, not policy-PDF language ([09-content-ux-writing.md](09-content-ux-writing.md#plain-language)).
- **Easy export and deletion:** self-service, same-channel, without [obstruction](#obstruction)-style delays; account-level controls belong where users expect them ([19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md)).
- **Just-in-time permission prompts** with granular choices; never nudge toward oversharing ([forced action](#forced-action-privacy-zuckering)).

### Public-service and institutional trust

Government and public-institution services (GOV.UK, USWDS doctrine) operate under stricter rules than commercial products:

- **No persuasion mechanics:** citizens cannot take their business elsewhere, so urgency, gamification, upsells, and engagement optimization have no legitimate place; success is measured by user outcomes, not completion volume or engagement.
- **Plain language throughout** ([09-content-ux-writing.md](09-content-ux-writing.md#plain-language)); clearly identify the service provider; use secure domains and avoid deceptive trust signals.
- **Continuity of service:** keep content accurate and reviewed; provide support and feedback channels; degrade gracefully — an unavailable government service can't be "churned" away from, only failed.

**Sources:** [GOV.UK Service Manual](https://www.gov.uk/service-manual), [USWDS](https://designsystem.digital.gov/).

---

## Controversial patterns

Legitimate uses and real risks; the entries weigh both.

### Infinite scroll

**Definition:** Content loading continuously as the user scrolls, without pagination boundaries.

**Reasoning / Evidence:** Fits open-ended *discovery* feeds (social, inspiration) where the goal is grazing; hostile to goal-directed tasks (can't return to position, can't estimate scope, footer unreachable, back-navigation position loss) and implicated in compulsive-use critiques ("time well spent," bottomless-bowl effect — the variable-reward loop is the [flow](02-cognitive-foundations.md#flow-state) exploit).

**When to use:** Discovery feeds with no task-completion semantics; paired with position restoration and "new since" markers.
**When NOT to use:** Search results, product listings, tables, anything users compare/return-to ([08-navigation-ia.md](08-navigation-ia.md#pagination-vs-infinite-scroll-vs-load-more) decision table — load-more is usually the better hybrid); anywhere a footer must be reachable.

**Pros:** Frictionless grazing; mobile-ergonomic. **Cons:** Position loss, scope-blindness, footer-kill, compulsion mechanics, memory bloat.

**Implementation guidance:** Load-more button beats pure infinite for semi-goal contexts; virtualize; restore scroll on back ([13-platform-web.md](13-platform-web.md#browser-conventions)); provide natural stopping cues ("You're all caught up") — the ethical-design marker.

**Related:** [08-navigation-ia.md](08-navigation-ia.md#pagination-vs-infinite-scroll-vs-load-more), [#gamification-and-streaks](#gamification-and-streaks).
**Sources:** [NN/g: Infinite Scrolling](https://www.nngroup.com/articles/infinite-scrolling-tips/).

### Hamburger menu

**Definition:** Primary navigation collapsed behind the ☰ icon.

**Reasoning / Evidence:** NN/g's hidden-navigation studies: hidden nav halves discoverability and slows tasks (desktop especially); but screen economics on mobile are real, and the icon itself is now universally recognized — the cost is *hiding*, not the glyph.

**When to use:** Mobile secondary/long-tail navigation behind visible bottom tabs ([15-platform-mobile.md](15-platform-mobile.md#mobile-navigation-patterns)); genuinely rarely-needed destinations.
**When NOT to use:** Desktop (space exists — show the nav); as the *sole* top-level navigation on mobile for engagement-critical destinations.

**Pros:** Space recovery; familiar icon. **Cons:** Measured discoverability/engagement loss for everything inside.

**Implementation guidance:** Priority+ hybrid (show top 3–5, overflow the rest); label it "Menu" where space allows (labeled beats bare icon in tests).

**Related:** [08-navigation-ia.md](08-navigation-ia.md#hamburger-menu), [15-platform-mobile.md](15-platform-mobile.md#mobile-navigation-patterns).
**Sources:** [NN/g: Hamburger Menus](https://www.nngroup.com/articles/hamburger-menus/).

### Carousels (user-controlled)

**Definition:** Swipeable/arrow-stepped content sequences *without* auto-rotation.

**Reasoning / Evidence:** Distinct from [auto-rotating banners](#auto-rotating-carousels): user-controlled carousels work where browsing-a-set is the actual task (product image galleries, media rows — Netflix's entire IA) with position clarity and peek affordances ([closure](01-core-principles.md#closure)).

**When to use:** Image galleries; homogeneous content rows in browse contexts; mobile-native swipe sets.
**When NOT to use:** Heterogeneous priority messages (stakeholder parking lot); content users must *compare* (carousels serialize what should be parallel).

**Pros:** Space-efficient set browsing; natural touch idiom. **Cons:** Items beyond slide 1 get steeply less exposure; arrow-stepping on desktop is high-interaction-cost.

**Implementation guidance:** Peek next item (10–20% visible); position indicators; swipe + arrows + keyboard; lazy-load with reserved space (CLS — [13-platform-web.md](13-platform-web.md#core-web-vitals)).

**Related:** [#auto-rotating-carousels](#auto-rotating-carousels), [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures).
**Sources:** [NN/g: Carousel Usability](https://www.nngroup.com/articles/designing-effective-carousels/).

### Gamification and streaks

**Definition:** Game mechanics on non-game tasks: points, badges, levels, leaderboards, and streaks (consecutive-day counters with loss-on-break).

**Reasoning / Evidence:** Real efficacy for genuinely-user-valued goals (Duolingo's streak drives measurable learning retention; fitness rings drive activity) via [goal-gradient](01-core-principles.md#goal-gradient-effect) and [Zeigarnik](01-core-principles.md#zeigarnik-effect) mechanics. The dark twin: loss-aversion weaponized ([02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux)) — streak anxiety, engagement-for-engagement's-sake, guilt mechanics on entertainment products where "commitment" serves only retention metrics.

**When to use:** User-chosen goals with real progress (learning, health, savings); celebration over punishment framing.
**When NOT to use:** Consumption products (streaks on watching/scrolling = compulsion engineering); anywhere breaking the streak punishes disproportionately to the user's actual goal.

**Pros:** Motivation scaffolding for hard-to-sustain user goals. **Cons:** Anxiety, crowding out intrinsic motivation (overjustification effect), metric-gaming behavior, ethical decay gradient.

**Implementation guidance:** Streak freezes/repair (mercy mechanics) as default; celebrate milestones, don't mourn breaks; audit with the [spectrum test](#the-deceptive-neutral-persuasive-spectrum): does the mechanic serve the user's stated goal or the retention dashboard?

**Related:** [01-core-principles.md](01-core-principles.md#goal-gradient-effect), [#infinite-scroll](#infinite-scroll).
**Sources:** [NN/g: Gamification](https://www.nngroup.com/articles/gamification/).

### Personalization and filter bubbles

**Definition:** Algorithmic tailoring of content/products/rankings to inferred individual preference.

**Reasoning / Evidence:** Genuine relevance value (recommendation quality drives satisfaction in commerce/media measurably) vs documented costs: filter bubbles (narrowing exposure), opacity ("why am I seeing this?"), cold-start clumsiness, privacy input costs, and error persistence (one gift purchase haunts recommendations for months). DSA requires recommender-system transparency and a non-profiling option for large platforms.

**When to use:** High-volume heterogeneous catalogs (commerce, media) with explicit user benefit; always with visible controls.
**When NOT to use:** Critical/civic information (news diversity), pricing (personalized pricing reads as discrimination and is regulated), safety information.

**Pros:** Relevance, discovery, reduced choice overload ([02-cognitive-foundations.md](02-cognitive-foundations.md#decision-fatigue-and-choice-overload)). **Cons:** Bubble effects, creepiness threshold, error lock-in, regulatory exposure.

**Implementation guidance:** Explainability ("Because you watched…"), correction controls ("not interested," category resets), non-personalized fallback view, and personalization *of chrome* (reordering a user's own tools) before personalization *of truth* (what exists).

**Related:** [02-cognitive-foundations.md](02-cognitive-foundations.md#cognitive-biases-in-ux), [#the-deceptive-neutral-persuasive-spectrum](#the-deceptive-neutral-persuasive-spectrum).
**Sources:** [NN/g: Customization vs Personalization](https://www.nngroup.com/articles/customization-personalization/), [DSA Art. 27/38](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32022R2065).

### Autoplay video

**Definition:** Video starting without user action — muted-inline (feeds, previews) through full-sound (the indefensible end).

**Reasoning / Evidence:** Muted autoplay previews measurably aid content selection (Netflix previews, feed videos) and browsers have forced the armistice (sound-on autoplay blocked by Chrome/Safari policies). Costs: data consumption, attention capture ([banner blindness](02-cognitive-foundations.md#selective-attention-and-banner-blindness)'s motion-grab inverse), motion-sensitivity ([05-accessibility.md](05-accessibility.md#reduced-motion)), and WCAG 2.2.2 (pause/stop/hide for >5s auto-content, 1.4.2 audio control).

**When to use:** Muted, short, in-viewport-only previews in media-browsing contexts, with visible pause and reduced-motion/data-saver respect.
**When NOT to use:** Sound-on ever; content pages (article videos); metered-data contexts without opt-in; background hero videos that fight text legibility ([03-visual-design.md](03-visual-design.md#imagery-and-illustration)).

**Pros:** Preview efficiency; engagement in browse contexts. **Cons:** Data, distraction, accessibility, battery.

**Implementation guidance:** Muted + captions default ([05-accessibility.md](05-accessibility.md)); pause when out of viewport; honor `prefers-reduced-motion` and data-saver; global "never autoplay" setting.

**Related:** [05-accessibility.md](05-accessibility.md#motion-and-flashing-limits), [04-interaction-design.md](04-interaction-design.md#motion-and-animation).
**Sources:** [WCAG 2.2.2](https://www.w3.org/WAI/WCAG22/Understanding/pause-stop-hide.html), [Chrome autoplay policy](https://developer.chrome.com/blog/autoplay/).

### Floating action buttons

**Definition:** Material's circular floating primary-action button, docked above content.

**Reasoning / Evidence:** Works when the screen has one dominant creative action (compose, add) — thumb-optimal ([thumb zone](15-platform-mobile.md#thumb-zone)), always available, strong [Von Restorff](01-core-principles.md#von-restorff-effect) isolate. Fails when: no single dominant action exists (FAB becomes arbitrary), it occludes content/list actions (covering the last row's controls), or it hides functionality behind an ambiguous "+" ([mystery meat](#mystery-meat-navigation) adjacency).

**When to use:** Clear single-primary-action screens (compose in mail, new-item in lists); Android-first products (platform-native — [15-platform-mobile.md](15-platform-mobile.md#material-design-3-android)).
**When NOT to use:** Multi-action screens (speed-dial FABs are a smell), reading surfaces (occlusion), iOS-first designs (less native — toolbar actions fit better).

**Pros:** Reachable, prominent, conventional on Android. **Cons:** Occlusion, one-action limit, ambiguity when misassigned.

**Implementation guidance:** Content bottom-padding clearing the FAB; hide-on-scroll-down/return-on-up for reading; label extended-FAB variant where ambiguity exists.

**Related:** [15-platform-mobile.md](15-platform-mobile.md#material-design-3-android), [01-core-principles.md](01-core-principles.md#von-restorff-effect).
**Sources:** [Material: FAB](https://m3.material.io/components/floating-action-button/overview).

### AI chatbots as primary support

**Definition:** LLM/scripted chatbots as the front line (or entirety) of customer support.

**Reasoning / Evidence:** Legitimate wins: instant 24/7 answers for the FAQ head (deflection of genuinely simple queries), language coverage, and triage. Documented failures: hallucinated policies (Air Canada was held *legally liable* for its chatbot's invented refund policy, 2024), obstruction-by-bot (the hard-to-cancel variant where the bot is the maze — [#obstruction](#obstruction)), user rage at escalation-blocking, and trust damage when bots impersonate humans (DSA/FTC transparency direction: disclose the bot).

**When to use:** FAQ-head triage with *instant, visible* human escalation; drafted-by-AI human-reviewed replies; internal support tooling.
**When NOT to use:** Sole channel for consequential issues (billing disputes, safety, cancellations — regulators increasingly mandate human paths); unreviewed authoritative policy statements.

**Pros:** Latency, scale, cost, coverage. **Cons:** Hallucination liability, escalation rage, trust erosion, complex-case failure.

**Implementation guidance:** Disclose bot status upfront; escalation to human within ≤2 exchanges on request or detected frustration; ground responses in versioned policy sources (RAG with citations); log-and-audit answer accuracy; measure *resolution*, not deflection ([17-ux-process-research.md](17-ux-process-research.md#ux-metrics-frameworks) — deflection is a cost metric wearing a savings costume). Full AI-product doctrine: [20-ai-product-ux.md](20-ai-product-ux.md).

**Related:** [#obstruction](#obstruction), [09-content-ux-writing.md](09-content-ux-writing.md#voice-and-tone), [20-ai-product-ux.md](20-ai-product-ux.md).
**Sources:** [Air Canada chatbot ruling (CBC)](https://www.cbc.ca/news/canada/british-columbia/air-canada-chatbot-lawsuit-1.7116416), [NN/g: AI chat UX](https://www.nngroup.com/articles/ai-chat-ux/).

---

## Quick reference

| Pattern | Verdict | Key rule |
|---|---|---|
| Confirmshaming | Never | Neutral decline copy |
| Hard to cancel (roach motel) | Never (FTC §5/ROSCA enforced; Click-to-Cancel rule vacated 2025 but enforcement continues) | Cancellation parity |
| Hidden subscription (forced continuity) | Never | Pre-charge reminder 3–7 days |
| Hidden costs / drip pricing | Never (regulated) | All-in price at first display |
| Bait and switch / trick wording | Never | Label = action |
| Disguised ads | Never (FTC rules) | 2-second ad-distinction test |
| Fake scarcity/urgency/social proof | Never (EU blacklist) | Refresh test; real inventory only |
| Nagging | Never | Persist refusals; ask once in context |
| Obstruction | Never (DSA Art. 25) | Click-count symmetry audit |
| Preselection | Never (Planet49) | Consent/purchases start unchecked |
| Forced action (privacy Zuckering) | Never (GDPR) | Granular, plain, point-of-need |
| Visual interference (misdirection) | Never | Prominence = user importance |
| Sneaking (sneak into basket) | Never (blacklisted) | Cart diff vs explicit adds |
| Friend spam | Never | Preview-and-confirm sends |
| Comparison prevention | Never | One consistent unit |
| Mystery meat nav | Avoid | Icon+label; ≥80% recognition |
| Auto-carousels | Avoid | One hero; static grid |
| Splash screens | Avoid | Content ≤2s |
| Modal overuse | Avoid | Badge→banner→toast→modal ladder |
| Forced tours | Avoid | Do-first, contextual hints |
| Placeholder-as-label | Avoid | Persistent labels |
| Disabled submit | Avoid | Enabled + validate on click |
| Scrolljacking | Avoid | Native physics always |
| Breaking back | Avoid | 5-deep back audit |
| Chat popups | Avoid | Passive availability |
| Permission on load | Avoid | Two-step contextual ask |
| App interstitials | Avoid (SEO penalty) | Smart banner max |
| Form data loss | Avoid (severity max) | Preserve through everything |
| Aggressive timeout | Avoid (WCAG 2.2.1) | Warn 60s + extend + drafts |
| Low contrast | Avoid (top WCAG failure) | 4.5:1 token enforcement |
| Icon-only UI | Avoid | Labels or proven recognition |
| Gesture-only | Avoid (WCAG 2.5.x) | Visible equivalent always |
| Fake loading | Avoid (labor-illusion exception) | ≤1–2s, honest, real work |
| Hover-only on touch | Avoid | Hover accelerates, never gates |
| PDF content | Avoid | HTML first |
| Infinite scroll | Conditional | Feeds yes; goals no; stopping cues |
| Hamburger | Conditional | Mobile long-tail only; label it |
| Carousels (manual) | Conditional | Galleries yes; priority lists no |
| Gamification | Conditional | User goals only; mercy mechanics |
| Personalization | Conditional | Explain, control, fallback |
| Autoplay video | Conditional | Muted, short, pausable, in-viewport |
| FAB | Conditional | One dominant action, Android-first |
| AI support bots | Conditional | Disclose, ground, escalate ≤2 turns |
