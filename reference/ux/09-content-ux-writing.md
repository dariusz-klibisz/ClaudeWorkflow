# Content & UX Writing

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers the words in interfaces: plain language, microcopy (buttons, links), error messages, empty states, confirmations, notifications, onboarding copy, voice and tone, capitalization, number/date formatting, truncation, content hierarchy and scannability, internationalization/localization, inclusive language, placeholder/helper text, and reading-level targets. Typography and visual text presentation are in [03-visual-design.md](03-visual-design.md); form validation mechanics in [07-forms-and-input.md](07-forms-and-input.md); navigation label strategy (information scent) in [08-navigation-ia.md](08-navigation-ia.md); accessible text requirements (contrast, alternatives) in [05-accessibility.md](05-accessibility.md); content documentation within design systems in [10-design-systems.md](10-design-systems.md).

## Contents

- [Language foundations](#language-foundations)
  - [Plain language](#plain-language)
  - [Reading level targets](#reading-level-targets)
  - [Voice and tone](#voice-and-tone)
  - [Terminology and naming consistency](#terminology-and-naming-consistency)
  - [Inclusive language](#inclusive-language)
- [Microcopy](#microcopy)
  - [Button labels](#button-labels)
  - [Link text](#link-text)
  - [Placeholder and helper text](#placeholder-and-helper-text)
  - [Confirmation and destructive action copy](#confirmation-and-destructive-action-copy)
  - [Notification and toast copy](#notification-and-toast-copy)
  - [Onboarding copy and tooltips](#onboarding-copy-and-tooltips)
- [Message design](#message-design)
  - [Error message writing](#error-message-writing)
  - [Empty states](#empty-states)
- [Mechanics and formatting](#mechanics-and-formatting)
  - [Capitalization conventions](#capitalization-conventions)
  - [Numbers, dates, and units formatting](#numbers-dates-and-units-formatting)
  - [Truncation and line length](#truncation-and-line-length)
- [Structure and scannability](#structure-and-scannability)
  - [Content hierarchy and front-loading](#content-hierarchy-and-front-loading)
  - [Scannability: headings, bullets, bold](#scannability-headings-bullets-bold)
- [Global content](#global-content)
  - [Internationalization and localization](#internationalization-and-localization)
    - [Data-format matrix](#data-format-matrix)
    - [Culturalization checklist](#culturalization-checklist)
    - [Layout and QA tests](#layout-and-qa-tests)
- [Quick reference](#quick-reference)

---

## Language foundations

### Plain language

**Definition:** Writing that readers can understand the first time they read it: common words, short sentences, active voice, direct address ("you"), and concrete instructions. Codified in the US Plain Writing Act (2010), plainlanguage.gov guidelines, and design-system content standards (GOV.UK, Carbon, Polaris).

**Reasoning / Evidence:** NN/g's studies with domain experts found that even highly educated specialists prefer plain language — experts scan, are time-pressed, and read above their reading level only when forced; simplified pages improved usability for both experts and lay users. Comprehension gains compound with stress: users in error or money-related contexts read worse than their baseline. GOV.UK's research showed common words ("help," "buy") outperform formal equivalents ("assistance," "purchase") in both comprehension and findability (people search using common words). Plain language is not "dumbing down": it preserves precision while removing decoding effort.

**When to use:**
- All interface copy, all audiences — including expert/technical products (Polaris: some merchant jargon is fine *only if users actually say it*).
- Especially: errors, legal/consent flows, money, health, anything read under stress.

**When NOT to use / exceptions:**
- Precise technical terms users know and search for (e.g., "commit," "invoice," "APR") should stay — replacing established domain vocabulary with folksy paraphrase harms scent.
- Regulated wording (legal disclosures) may be fixed; pair mandated text with a plain summary.

**Pros:**
- Faster comprehension, fewer support tickets, better search findability, better translation output (plain source text localizes cleanly).

**Cons:**
- Requires editorial effort and domain knowledge to compress without losing accuracy; stakeholders may perceive it as insufficiently "professional" (evidence says otherwise).

**Implementation guidance:**
- Sentence length ≤25 words; average ~15. One idea per sentence.
- Prefer active voice ("We sent your report") except where the actor is irrelevant or blame-avoidance matters ("Your session expired").
- Use "you/your" for the user; "we" for the product organization.
- Replace: utilize→use, assistance→help, in order to→to, prior to→before, terminate→end.
- Use contractions ("don't," "you're") — they test as more natural and are recommended by Polaris and GOV.UK.
- Front-load: key noun/verb in the first 2 words of headings, links, and buttons.
- Test with the read-aloud check (Polaris: "does it sound like something a human would say?") and a readability formula gate in CI/review.
- Avoid internal acronyms unless users demonstrably know them; define unavoidable technical or legal terms at first use.

**Related:** [#reading-level-targets](#reading-level-targets), [#content-hierarchy-and-front-loading](#content-hierarchy-and-front-loading), [08-navigation-ia.md](08-navigation-ia.md) (information scent).

**Sources:** [NN/g: Plain Language Is for Everyone, Even Experts](https://www.nngroup.com/articles/plain-language-experts/), [plainlanguage.gov: Federal Plain Language Guidelines](https://www.plainlanguage.gov/guidelines/), [GOV.UK: Writing for GOV.UK / Words to avoid](https://www.gov.uk/guidance/style-guide), [Shopify Polaris: Content fundamentals](https://polaris.shopify.com/content/fundamentals).

### Reading level targets

**Definition:** A measurable ceiling on text complexity, usually expressed as a school-grade score (Flesch-Kincaid, SMOG, Gunning Fog): target grade 8 or lower for general audiences; grade 6–7 for high-reach consumer products and critical flows.

**Reasoning / Evidence:** Roughly half of US adults read at or below 8th-grade level (PIAAC/NAAL literacy surveys). GOV.UK recommends writing for a 9-year-old reading age for general audiences — note this is a deliberate writing target, not a measured population average; the widely repeated claim that "the UK average reading age is 9" is a weakly sourced popularization of that guidance. Reading ability degrades further under stress, on mobile, for non-native speakers (a large share of any global product's users), and for users with dyslexia or cognitive disabilities. NN/g's expert-audience research removes the usual objection: even PhDs prefer and perform better with lower-grade text. Shopify Polaris targets 7th grade for merchant-facing copy.

**When to use:**
- Set an explicit grade target per content type: UI microcopy and errors grade ≤7; help content grade ≤8; technical reference may run higher but sentences stay short.
- Gate long-form and high-stakes content (billing, consent, medical) through a readability check.

**When NOT to use / exceptions:**
- Formulas measure syllables and sentence length, not clarity — a jargon-free short sentence of rare words scores well but reads badly, and vice versa. Use scores as smoke alarms, not proof.
- Code, identifiers, legal boilerplate, and domain terms inflate scores artificially; exclude them from measurement.

**Pros:**
- Objective, automatable quality gate; creates a shared bar across writers; catches drift in long documents.

**Cons:**
- Gaming the metric (chopping sentences arbitrarily) produces choppy, harder text; scores vary between tools; poorly suited to non-English.

**Implementation guidance:**
- Targets: general audience ≤ grade 8; critical flows (errors, money, health, legal summaries) ≤ grade 6–7; expert tools ≤ grade 10 with ≤20-word sentences.
- Tools: Hemingway Editor, readable.com, `textstat` (Python) in CI, Microsoft Word's built-in Flesch-Kincaid.
- Flesch Reading Ease ≥60 for general content (≈ plain conversational English).
- Reduce grade level by: shorter sentences first, then common word substitutions, then splitting clauses — never by deleting necessary information.
- For localized products, set equivalent targets per language (e.g., LIX for Scandinavian languages) rather than translating the English metric.

**Related:** [#plain-language](#plain-language), [05-accessibility.md](05-accessibility.md) (cognitive accessibility, WCAG 3.1.5 Reading Level AAA), [#internationalization-and-localization](#internationalization-and-localization).

**Sources:** [plainlanguage.gov: Test your content](https://www.plainlanguage.gov/resources/checklists/), [GOV.UK content design: writing for a 9-year-old reading age](https://www.gov.uk/guidance/content-design/writing-for-gov-uk), [NN/g: Legibility, Readability, and Comprehension](https://www.nngroup.com/articles/legibility-readability-comprehension/), [WCAG 2.2 SC 3.1.5 Reading Level](https://www.w3.org/WAI/WCAG22/Understanding/reading-level.html).

### Voice and tone

**Definition:** Voice is the product's stable personality across all copy (e.g., "confident, helpful, plain"); tone is how that voice flexes with context along spectra such as serious↔playful, formal↔casual, matter-of-fact↔enthusiastic, respectful↔irreverent (NN/g's four tone dimensions). One voice, many tones.

**Reasoning / Evidence:** NN/g's tone-of-voice studies found tone measurably affects perceived trustworthiness and brand attitude; MailChimp's "Voice & Tone" guide established the industry model of mapping user emotional state → appropriate tone. The core rule: tone must match the user's stakes and mood, not the brand's marketing preference. Errors, security, billing, and data loss demand serious, calm tone (a joke next to a failure reads as mockery — see NN/g error guidance: "avoid humor since it can become stale if users encounter the error frequently"); success and empty "fresh start" moments can afford lightness. Polaris explicitly de-prioritizes voice polish in favor of "sounding human."

**When to use:**
- Define voice once (3–4 adjectives with "this-not-that" contrasts: "confident, not arrogant") in the design system content section.
- Map tone per context: onboarding (warm, encouraging), success (light, brief), errors (serious, constructive), legal/security (precise, calm), marketing surfaces (brand-forward).

**When NOT to use / exceptions:**
- Never playful: error messages users may see repeatedly, destructive confirmations, payment failures, security alerts, accessibility-critical instructions.
- Humor risks: repeats badly (funny once, grating on the 40th encounter), translates poorly, and can misread severity; if used at all, restrict to rare, low-stakes, dead-end moments (NN/g's "novelty in dire situations" — e.g., total-outage pages) and never in place of status information.

**Pros:**
- Consistent voice builds recognition and trust; context-matched tone reduces frustration and increases perceived empathy.

**Cons:**
- Personality-heavy copy ages fast, localizes badly, and inflates word counts; enforcement across many writers needs governance.

**Implementation guidance:**
- Document: voice adjectives + tone map (context → tone setting per dimension) + before/after examples; store alongside component docs ([10-design-systems.md](10-design-systems.md)).
- Practical tone dial per NN/g's 4 dimensions: pick a default position and specify deviations by context; e.g., errors = maximally serious/formal-leaning/matter-of-fact/respectful.
- Success messages: short and warm is enough ("Report sent") — no exclamation-point inflation; ≤1 exclamation mark per screen, ideally zero in products used daily.
- Test tone with users: show alternative phrasings of the same message and probe perceived trust/competence.

**Related:** [#error-message-writing](#error-message-writing), [#notification-and-toast-copy](#notification-and-toast-copy), [10-design-systems.md](10-design-systems.md) (content guidelines as system artifact).

**Sources:** [NN/g: The Four Dimensions of Tone of Voice](https://www.nngroup.com/articles/tone-of-voice-dimensions/), [NN/g: Tone of Voice and User Impressions](https://www.nngroup.com/articles/tone-voice-users/), [Shopify Polaris: Content fundamentals](https://polaris.shopify.com/content/fundamentals), [Mailchimp Content Style Guide: Voice and Tone](https://styleguide.mailchimp.com/voice-and-tone/).

### Terminology and naming consistency

**Definition:** One concept, one term, everywhere: the same object, action, or feature is named identically across navigation, buttons, headings, errors, docs, and marketing — maintained through a published glossary/controlled vocabulary. Also covers naming new features (descriptive vs branded names).

**Reasoning / Evidence:** Consistency is usability heuristic #4: synonym drift ("workspace" in nav, "project" in the button, "site" in the error) forces users to guess whether two words mean one thing or two — each guess is load and error risk, and it poisons search (users search the word they saw). Translation costs multiply directly with synonym count (each variant is a separate string to translate and keep consistent in N languages). Polaris dedicates a "Naming" guideline to this: prefer descriptive names users can decode on sight ("Order status page") over invented brands ("OrderFlow™") unless the feature is genuinely marketed standalone; NN/g's tabs research showed branded labels ("Legit" for theater reviews) carry weaker scent and less engagement than plain ones.

**When to use:**
- Every noun and verb that names a domain object or action: define once in a glossary, enforce in review.
- Cross-functional surfaces: the term in the UI must match docs, support macros, API names (where user-visible), and release notes.

**When NOT to use / exceptions:**
- Deliberate register shifts are allowed between UI and marketing ("Buy" button vs "purchase" in legal terms) — but the *object* names stay fixed.
- When users' own vocabulary conflicts with internal naming, users win: rename internally rather than teaching users your jargon (validate with card sorts/search logs).

**Pros:**
- Lower learning cost, working in-product search, cheaper translation, coherent support conversations ("click Orders" always works).

**Cons:**
- Glossary maintenance and enforcement effort; renames after drift are expensive (UI + docs + translations + user relearning).

**Implementation guidance:**
- Maintain a glossary in the design system ([10-design-systems.md](10-design-systems.md)): term, definition, part of speech, approved verb pairings ("archive a project," not "store"), forbidden synonyms, translations.
- Pick one verb per action product-wide: delete vs remove (Polaris: "remove" = detach recoverable, "delete" = destroy), edit vs modify, create vs add vs new — define the distinction or standardize on one.
- Feature naming test: can a first-time user predict what it does from the name alone? Descriptive beats branded for anything inside the product.
- Enforce with a content linter (Vale rules from the glossary) in docs/UI-string CI; audit for synonym drift quarterly (grep UI strings for known synonym sets).
- When renaming, ship transitional copy ("Projects (previously Workspaces)") for 1–2 releases and update every surface in one release, not gradually.

**Related:** [#plain-language](#plain-language), [08-navigation-ia.md](08-navigation-ia.md) (labels, information scent), [10-design-systems.md](10-design-systems.md) (content standards in the system).

**Sources:** [Shopify Polaris: Naming](https://polaris.shopify.com/content/naming), [NN/g: Consistency and Standards (Heuristic #4)](https://www.nngroup.com/articles/consistency-and-standards/), [NN/g: Tabs, Used Right — descriptive labels](https://www.nngroup.com/articles/tabs-used-right/).

### Inclusive language

**Definition:** Word choices that avoid excluding, stereotyping, or othering users: gender-neutral phrasing, person-first or identity-first disability language as communities prefer, avoiding ableist/violent metaphors, culturally loaded idioms, and unnecessary demographic assumptions.

**Reasoning / Evidence:** Interface language reaches audiences of every gender, ability, culture, and background; exclusionary defaults ("he," "guys," "crazy," "grandmother-friendly") measurably alienate segments and create legal/brand risk. Industry standards have converged: Microsoft, Google, Apple style guides and design systems (Polaris's inclusive-language section, Carbon's content guidance) all specify neutral terms. Technical vocabulary is also shifting (allowlist/blocklist for whitelist/blacklist; primary/replica for master/slave) — adopted by Linux kernel, GitHub, Chromium.

**When to use:**
- All copy: UI, docs, error messages, sample data (names in examples should reflect diverse cultures), marketing.
- Forms asking about people: gender, titles, and name fields need inclusive options (see [07-forms-and-input.md](07-forms-and-input.md)).

**When NOT to use / exceptions:**
- Don't rewrite direct quotes or legally fixed text; annotate instead.
- Follow the specific community's stated preference where they conflict with generic rules (e.g., many Deaf and autistic communities prefer identity-first language: "Deaf person," not "person with deafness").

**Pros:**
- Widens usable audience; reduces harm and complaints; usually also simplifies copy (neutral phrasing is often shorter).

**Cons:**
- Terminology evolves; requires periodic audits and a maintained word list; over-policing tone can slow content work.

**Implementation guidance:**
- Use "they/their" for unknown persons; "everyone/folks/team" instead of "guys"; avoid gendering roles ("chairperson," "staffing" not "manpower").
- Replace ableist idioms: "crazy/insane" → "unexpected/huge"; "blind spot" → "gap"; "sanity check" → "quick check"; "cripple" → "impair."
- Technical terms: allowlist/blocklist, primary/replica, placeholder value (not "dummy"), main branch.
- Avoid idioms and pop-culture references entirely in UI copy — they exclude non-native speakers and break translation ([#internationalization-and-localization](#internationalization-and-localization)).
- Maintain a lintable word list (e.g., Vale, alex.js) wired into content review/CI.

**Related:** [#internationalization-and-localization](#internationalization-and-localization), [05-accessibility.md](05-accessibility.md), [07-forms-and-input.md](07-forms-and-input.md) (personal-data fields).

**Sources:** [Shopify Polaris: Inclusive language](https://polaris.shopify.com/content/inclusive-language), [Microsoft Style Guide: Bias-free communication](https://learn.microsoft.com/en-us/style-guide/bias-free-communication), [Google developer documentation style guide: Inclusive documentation](https://developers.google.com/style/inclusive-documentation), [IBM Carbon: Content overview](https://carbondesignsystem.com/guidelines/content/overview/).

---

## Microcopy

### Button labels

**Definition:** The text on action controls. Rule: verb-first, specific to the outcome ("Save changes," "Delete 3 files," "Send invoice"), never generic ("OK," "Submit," "Yes") when the action has consequences, and never "Click here."

**Reasoning / Evidence:** Users scan, they don't read; a button label is often the only text a user processes before acting. Specific verb labels let the button carry its full meaning standalone — critical because users frequently act without reading the preceding dialog text (dialog-blindness is well documented; NN/g and platform HIGs both instruct that button labels must make sense in isolation). Generic pairs ("OK/Cancel") force users to read and correctly map the question, doubling error opportunity on destructive dialogs. Polaris codifies: start with a verb, be direct ("Add apps," not "You can add apps").

**When to use:**
- All buttons and primary actions. Label = verb (+ object when ambiguous): "Create project," "Export CSV," "Delete account."
- Dialogs: buttons restate the specific action ("Delete draft" / "Keep editing"), never "Yes/No."

**When NOT to use / exceptions:**
- Universally-understood single-word conventions in tight chrome are fine: "Save," "Cancel," "Close," "Search."
- Icon-only buttons acceptable only for the few near-universal glyphs (search, close, play) and always with an accessible name and tooltip.
- Marketing CTAs may use value-phrases ("Start free trial") — still verb-first and specific.

**Pros:**
- Actions are self-describing, safer under scan-reading, and screen-reader friendly (labels announced out of visual context must stand alone).

**Cons:**
- Specific labels are longer — pressure in tight layouts and after 30–40% localization expansion; requires per-context writing rather than reusable generic strings.

**Implementation guidance:**
- Format: Verb or Verb+Noun, 1–3 words typical, ≤~25 characters before truncation risk in translation.
- Sentence case ("Save changes"), no ending punctuation; no "Please."
- Match the triggering intent: a button opening a further step can use "Continue"; the final commit must name the deed ("Place order").
- Destructive buttons: name the object and consequence ("Delete 12 photos"), styled as destructive, never the default focus ([#confirmation-and-destructive-action-copy](#confirmation-and-destructive-action-copy)).
- Avoid: "Submit" (name the outcome), "Click here," "Learn more" as a button label without context, ALL CAPS (legibility penalty; also reads as shouting).

**Verification checklist:**
- Every button/link label describes the *result* of activation ("Save draft," "Send invoice," "Download CSV") — not the mechanism ("OK," "Submit," "Click here").
- "Continue" appears only where the consequence of continuing is unambiguous.
- Each label makes sense read in isolation (screen-reader controls list test).
- Labels survive +30–40% translation expansion without truncation.

**Related:** [#link-text](#link-text), [#confirmation-and-destructive-action-copy](#confirmation-and-destructive-action-copy), [04-interaction-design.md](04-interaction-design.md) (button behavior), [07-forms-and-input.md](07-forms-and-input.md).

**Sources:** [Shopify Polaris: Content fundamentals — Inspire action](https://polaris.shopify.com/content/fundamentals), [NN/g: UX Writing — Study Guide](https://www.nngroup.com/articles/ux-writing-study-guide/), [GOV.UK Design System: Button component guidance](https://design-system.service.gov.uk/components/button/), [Apple HIG: Buttons](https://developer.apple.com/design/human-interface-guidelines/buttons).

### Link text

**Definition:** The clickable words of a hyperlink. Rule: the link text alone must describe the destination — meaningful out of context — and be front-loaded with the distinguishing words. "Click here," "here," "read more," and raw URLs are prohibited as link text.

**Reasoning / Evidence:** Three converging forces: (1) scanning — links are visual anchors in the F-pattern; users hop link-to-link and judge scent from link words only; (2) accessibility — screen-reader users pull a links list out of context, where five "Read more" links are indistinguishable (WCAG 2.4.4 Link Purpose); (3) SEO — anchor text is a relevance signal. "Click here" also fails device-neutrality (no clicking on touch/voice) and wastes the link's strongest words on an instruction users already know.

**When to use:**
- All inline links, card links, and "see also" references: link the noun phrase naming the destination ("view the billing FAQ," "download the annual report (PDF, 2.3 MB)").

**When NOT to use / exceptions:**
- Repeated card layouts may visually show a short "Read more" affordance, but the accessible name must be unique (aria-label or clipped full-title link).
- Don't over-link: more than ~1 link per sentence fragments reading; consolidate.

**Pros:**
- Higher scent, scannability, screen-reader usability, and search ranking from the same edit.

**Cons:**
- Requires writing sentences so a noun phrase can carry the link — occasional restructuring effort.

**Implementation guidance:**
- Link the 2–5 words that name the destination; front-load distinguishing terms ("2024 salary bands" not "information about salary bands for 2024").
- Sibling links on a page must be mutually distinguishable by text alone.
- Announce format/behavior surprises in the link text: "(PDF)", "(opens in new tab)" — and avoid new tabs by default.
- Style links distinctly (underline in body text; color alone fails WCAG 1.4.1); keep visited-state styling in link-dense content ([08-navigation-ia.md](08-navigation-ia.md) wayfinding).
- Match link text to the destination page's H1/title (scent confirmation on arrival).

**Related:** [#button-labels](#button-labels), [08-navigation-ia.md](08-navigation-ia.md) (information scent), [05-accessibility.md](05-accessibility.md) (WCAG 2.4.4).

**Sources:** [NN/g: Writing Hyperlinks — Salient, Descriptive, Start with Keyword](https://www.nngroup.com/articles/writing-links/), [NN/g: "Learn More" Links — You Can Do Better](https://www.nngroup.com/articles/learn-more-links/), [W3C WCAG 2.2 SC 2.4.4 Link Purpose (In Context)](https://www.w3.org/WAI/WCAG22/Understanding/link-purpose-in-context.html), [GOV.UK style guide: Links](https://www.gov.uk/guidance/style-guide/a-to-z#links).

### Placeholder and helper text

**Definition:** Placeholder = ghost text inside an empty input; helper text = persistent guidance below/beside a field (format hints, constraints, why-we-ask). Rule: placeholders never carry required information; anything the user needs while typing or reviewing goes in persistent helper text or the label.

**Reasoning / Evidence:** NN/g's placeholder research documents multiple failures: placeholder text disappears on focus/typing (users must clear the field to re-read instructions, straining short-term memory); fields with placeholders look pre-filled and get skipped in scanning; low-contrast gray fails low-vision users; placeholders-as-labels break entirely once the form is partially filled (no way to review which field is which) and are inconsistently announced by screen readers. Helper text avoids all of this at the cost of vertical space.

**When to use:**
- Helper text: format examples ("DD/MM/YYYY"), constraints ("8+ characters, one number"), sensitive-data justification ("We use this only for delivery updates"), consequences ("You can change this later").
- Placeholder: at most an *illustrative example* that duplicates nothing essential (e.g., "e.g., jane@company.com") — or simply omit.

**When NOT to use / exceptions:**
- Never use placeholder as the field label; never put validation rules only in placeholder.
- Search boxes are the tolerated exception (placeholder "Search products…" with an adjacent icon+label context).
- Don't write helper text for self-evident fields ("Enter your first name" under "First name") — noise erodes attention for hints that matter.

**Pros (helper text):**
- Persistent, reviewable, screen-reader associable (`aria-describedby`), visible during and after entry.

**Cons (helper text):**
- Consumes vertical space; long hints between fields slow form scanning — compress ruthlessly.

**Implementation guidance:**
- Position helper text below the label or below the field, consistently product-wide; keep it to 1 short line (~≤80 characters).
- Programmatically associate via `aria-describedby`; error text replaces or appends to helper text in the same location ([07-forms-and-input.md](07-forms-and-input.md)).
- Show password rules up front as a live checklist, not as a post-submit error.
- If a hint is only needed by some users, consider a progressive-disclosure "What's this?" toggle rather than always-on text.
- Placeholder styling, if used: WCAG 1.4.3 applies to visible placeholder text, so it must meet 4.5:1 contrast like any other text — but at that contrast the field looks pre-filled and gets skipped. This unresolvable tension is one more reason to avoid relying on placeholders (and never as labels): keep required information in the label and helper text instead.

**Related:** [07-forms-and-input.md](07-forms-and-input.md) (labels, validation), [#error-message-writing](#error-message-writing), [05-accessibility.md](05-accessibility.md).

**Sources:** [NN/g: Placeholders in Form Fields Are Harmful](https://www.nngroup.com/articles/form-design-placeholders/), [NN/g: Help and Hint Text in Forms](https://www.nngroup.com/articles/help-hint-text-forms/), [GOV.UK Design System: Text input — hint text](https://design-system.service.gov.uk/components/text-input/).

### Confirmation and destructive action copy

**Definition:** The wording of dialogs and flows that gate consequential or irreversible actions (delete, cancel subscription, overwrite, send to many people): a specific question, a statement of consequence, and buttons that name the action rather than "Yes/No/OK."

**Reasoning / Evidence:** Routine confirmations train click-through: users develop automatic dismissal within a handful of exposures, so confirmation copy only works when confirmations are rare and the copy forces re-engagement — specificity ("Delete 'Q3 Budget.xlsx'?") interrupts autopilot better than generic prompts ("Are you sure?"). Button-label research (platform HIGs, NN/g) shows verb-labeled buttons prevent the classic mis-map where users answer the wrong question. For truly destructive bulk actions, friction proportional to loss (type-to-confirm) is the industry standard (GitHub repo deletion). Undo beats confirmation for frequent reversible actions (see [04-interaction-design.md](04-interaction-design.md)).

**When to use:**
- Irreversible or costly actions: permanent deletion, sends to external recipients, financial commits, permission grants, overwrites.
- Scope escalation: acting on many items or someone else's data — state the count and scope in the copy.

**When NOT to use / exceptions:**
- Frequent, reversible actions: use undo (toast with "Undo," 5–10s window) instead of confirmation.
- Never stack confirmations (two dialogs) — a second dialog gets even less reading than the first.
- Don't use confirmations to shift blame for bad defaults ("You didn't save!") — fix with autosave.

**Pros:**
- Prevents high-cost slips when rare and specific; creates an explicit consent record for consequential operations.

**Cons:**
- Habituation destroys value if overused; every confirmation is friction charged to all users to protect a few.

**Implementation guidance:**
- Title = the question with specifics: "Delete 12 photos?" Body = consequence in one sentence: "They'll be removed from all shared albums. This can't be undone."
- Buttons: destructive verb + object ("Delete 12 photos") vs safe exit ("Keep photos" or "Cancel"); destructive button styled red and NOT the default/Enter-focused action.
- Include what will NOT happen when users plausibly fear it ("Your account and billing are unaffected").
- Type-to-confirm (retype resource name) only for catastrophic scope: production resources, whole-account deletion, repos.
- Cancel-flow copy must not guilt-trip ("confirmshaming" — "No thanks, I hate saving money" — is a dark pattern; see [18-patterns-antipatterns.md](18-patterns-antipatterns.md)).

**Verification checklist:**
- Consequences are explained *before* the irreversible action, in the dialog itself — not in distant documentation.
- The user could answer the dialog correctly from the buttons alone, without reading the body text.

**Related:** [#button-labels](#button-labels), [04-interaction-design.md](04-interaction-design.md) (undo, safety), [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Sources:** [NN/g: Confirmation Dialogs Can Prevent User Errors](https://www.nngroup.com/articles/confirmation-dialog/), [NN/g: User Errors — Slips and Mistakes](https://www.nngroup.com/articles/slips/), [Apple HIG: Alerts](https://developer.apple.com/design/human-interface-guidelines/alerts), [GOV.UK: Asking users to confirm](https://design-system.service.gov.uk/patterns/confirmation-pages/).

### Notification and toast copy

**Definition:** Short system-initiated messages: toasts/snackbars (transient, low-stakes feedback), banners (persistent conditions), badges, and push notifications. Copy must convey what happened, to what, and (if relevant) the next action — in the first few words.

**Reasoning / Evidence:** Toasts are glanced at for 1–3 seconds or missed entirely; comprehension depends on front-loading and brevity. NN/g's indicator/notification research: match persistence to importance (transient toasts must never carry information users need later or actions they must take — screen-reader and glance-away users will miss them); notification fatigue is real and dose-dependent — every low-value notification lowers attention to all future ones (cry-wolf effect). Push notifications compete on a lock screen: the first ~40 characters decide the swipe.

**When to use:**
- Toast: confirm a completed background action ("Invoice sent," "Link copied"), optionally with a single inline action ("Undo").
- Banner: ongoing states needing awareness (offline, degraded service, unpaid bill) — persist until resolved or dismissed.
- Push/email: user-relevant events they opted into; never for engagement bait.

**When NOT to use / exceptions:**
- Never put errors requiring action in auto-dismissing toasts; use inline errors or persistent alerts.
- Don't toast what the UI already shows (item visibly appears in list = no "Item added" toast needed).
- Silence is correct for micro-interactions with immediate visible results.

**Pros:**
- Lightweight status visibility (usability heuristic #1) without workflow interruption.

**Cons:**
- Transience excludes anyone not looking; stacking toasts becomes noise; overuse trains users to ignore all messaging.

**Implementation guidance:**
- Structure: [Object/outcome] + [optional next step]. "Report exported" / "Payment failed — update your card." Front-load the outcome word.
- Length: toast ≤~10 words / 1 line; push title ≤40 chars, body ≤120 chars visible on most lock screens.
- Duration: minimum 5s; a practical floor is (word count × 300ms) + 2s; persist on hover/focus; toasts with actions ≥8–10s; always dismissible.
- No "Success!"-only toasts — name what succeeded. Skip terminal punctuation on fragments; use it for full sentences.
- Announce via ARIA live region (`polite` for info, `assertive` only for urgent) — see [05-accessibility.md](05-accessibility.md).
- Batch and rate-limit push notifications; every category user-controllable; deep-link each notification to the exact relevant screen ([08-navigation-ia.md](08-navigation-ia.md) deep linking).

**Related:** [#voice-and-tone](#voice-and-tone), [#error-message-writing](#error-message-writing), [04-interaction-design.md](04-interaction-design.md) (feedback timing), [15-platform-mobile.md](15-platform-mobile.md) (push).

**Sources:** [NN/g: Indicators, Validations, and Notifications](https://www.nngroup.com/articles/indicators-validations-notifications/), [NN/g: Toasts — Definition and Guidelines](https://www.nngroup.com/articles/toast/), [Material Design 3: Snackbar](https://m3.material.io/components/snackbar/guidelines).

### Onboarding copy and tooltips

**Definition:** Text that teaches the product: first-run intros, coach marks, tooltips, feature-announcement popovers, and contextual "pull revelations" that surface when users engage with the relevant element.

**Reasoning / Evidence:** NN/g's onboarding and help research consistently finds front-loaded tutorials underperform contextual help: users skip carousels/tours to reach their task, and lessons taught before need have no anchor to bind to (little transfer to the actual interface). In-context cues at the moment of relevance ("pull" help, e.g., empty-state instructions, first-hover tooltips) are applied immediately and remembered. Interruptive multi-step product tours score worst: they block the goal, present information serially out of context, and get "X"-ed reflexively.

**When to use:**
- Contextual tooltips for unlabeled icon controls (hover/focus, 300–500ms delay) — a tooltip is a label of last resort, not a home for critical instructions.
- Just-in-time coach marks: single, dismissible, anchored to the element, at first encounter of a genuinely non-obvious feature.
- Empty states as onboarding ([#empty-states](#empty-states)); templates/demo data as "learning by example."

**When NOT to use / exceptions:**
- Don't tour every feature at first run; cap forced onboarding at 0–3 screens communicating only what's needed to start.
- Never gate core tasks behind completing a tutorial; always provide Skip.
- Tooltips can't be the only access to information on touch devices (no hover) — anything important needs a visible or tap path.

**Pros:**
- Context-bound learning sticks; users reach value faster; feature discovery without permanent UI weight.

**Cons:**
- Coach-mark fatigue: products that announce every release train reflexive dismissal; overlays block the UI they explain.

**Implementation guidance:**
- Tooltip copy: ≤~8 words for label tooltips; ≤2 short lines for descriptive ones; no interactive content inside plain tooltips (use popovers if interaction is needed); keyboard-focus triggers required.
- Coach marks: one at a time, max ~3 per release; verb-first benefit framing ("Filter by team to find work faster"), explicit dismiss, never re-show after dismissal.
- First-run: prefer a "get started" checklist (user-paced, persistent, 3–5 items) over a modal tour; mark completion.
- Write onboarding copy about the user's task, not the feature ("Track your first shipment" not "Introducing ShipTracker 2.0").
- Instrument dismiss-without-read rates; >80% instant dismissal means delete the message, not rewrite it.

**Related:** [#empty-states](#empty-states), [02-cognitive-foundations.md](02-cognitive-foundations.md) (learning, memory), [04-interaction-design.md](04-interaction-design.md) (progressive disclosure).

**Sources:** [NN/g: Onboarding Tutorials vs. Contextual Help](https://www.nngroup.com/articles/onboarding-tutorials/), [NN/g: Help and Documentation — pull revelations](https://www.nngroup.com/articles/help-and-documentation/), [NN/g: Tooltip Guidelines](https://www.nngroup.com/articles/tooltip-guidelines/), [NN/g: Mobile-App Onboarding](https://www.nngroup.com/articles/mobile-app-onboarding/).

---

## Message design

### Error message writing

**Definition:** Copy for system interruptions reporting incomplete, incompatible, or failed operations. Canonical structure: (1) what happened, (2) why / which input or condition caused it, (3) how to fix it — in plain, blame-free, jargon-free language, near the error's source.

**Reasoning / Evidence:** Nielsen's heuristic #9 ("Help users recognize, diagnose, and recover from errors") has 30+ years of validation. NN/g's 2023 error-message guidelines: generic messages ("An error occurred") strand users; technical codes serve diagnostics, not users; blame-words ("invalid," "illegal," "forbidden," "you failed") raise anxiety and abandonment; ALL CAPS reads as shouting and scans worse; premature validation (erroring on fields the user hasn't finished) is a hostile pattern; humor in errors goes stale on repetition and can obscure status (Disney's pun-instead-of-status example). Constructive advice with reduced correction effort (preserve input, suggest fixes, offer pick-one corrections) measurably improves recovery.

**When to use:**
- Every failure surface: form validation, network/server failures, permission denials, empty/failed searches, background job failures.
- The structure applies at every scale — a field-level error is a one-line compression of the same three parts.

**When NOT to use / exceptions:**
- Best error is none: prevent with constraints, good defaults, forgiving parsing (accept spaces in card numbers), and mistake-detection (missing-attachment check) before erroring.
- Security-sensitive messages must stay vague on purpose ("email or password incorrect" — don't reveal which); compensate with a strong recovery path.
- Total-outage/dead-end pages are the one place NN/g allows novelty/whimsy — status and apology first, personality second.

**Pros:**
- Directly converts failures into recoveries; reduces support volume; error copy quality is among the highest-ROI writing in a product.

**Cons:**
- Specific messages require enumerating failure modes with engineering (error taxonomy work); poorly mapped backend errors leak jargon.

**Implementation guidance:**
- Formula: [What happened] + [why/what condition] + [how to fix / next step]. "We couldn't save your changes — the file was modified by Ana at 14:32. Reload to see the latest version, or save a copy."
- Placement: adjacent to the source (field-level under the field; page-level summary linking to fields for long forms); visible without scrolling; icon + red + text (never color alone — ~300M people have color-vision deficiency).
- Words to ban: invalid, illegal, forbidden, fatal, bad, prohibited, "you failed to," oops/whoops (in serious contexts). Say "doesn't match the required format: DD/MM/YYYY" instead of "invalid date."
- No raw codes/stack traces in the primary message; if support needs a code, append it small ("Error ref: 4A7-002") with a copy affordance.
- Preserve everything the user entered; never clear a form on error. Offer computed fixes as tappable options where possible (ZIP↔city mismatch → suggest matching city).
- Timing: validate on blur-after-input or on submit, not on first keystroke or on focus-leave-while-empty; live checklist validation for known-hard fields (passwords).
- Tone: serious, calm, specific; the system takes responsibility ("We couldn't process…" not "You entered wrong data").

**Verification checklist:**
- Every error message states what happened, names the specific object/location, and gives a recovery action (formula: problem + specific location/object + recovery). Examples that pass: "Enter an email address in the format name@example.com." / "The file is larger than 10 MB. Choose a smaller file or compress it." / "We could not save your changes. Check your connection and try again."
- No message consists only of "An error occurred," a code, or "Invalid input."
- User input is preserved after every error path.

**Related:** [07-forms-and-input.md](07-forms-and-input.md) (validation mechanics), [#voice-and-tone](#voice-and-tone), [05-accessibility.md](05-accessibility.md) (error identification, `aria-live`), [01-core-principles.md](01-core-principles.md) (heuristic #9).

**Sources:** [NN/g: Error-Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/), [NN/g: Hostile Error Messages](https://www.nngroup.com/articles/hostile-error-messages/), [NN/g: 10 Design Guidelines for Reporting Errors in Forms](https://www.nngroup.com/articles/errors-forms-design-guidelines/), [Shopify Polaris: Error messages](https://polaris.shopify.com/content/error-messages).

### Empty states

**Definition:** The designed content of containers with nothing to show. Three distinct types requiring different copy: **first-use** (feature never used — onboarding opportunity), **cleared/completed** (user emptied it — status + positive closure), **no-results** (query/filter matched nothing — recovery; see also [08-navigation-ia.md#no-results-recovery](08-navigation-ia.md#no-results-recovery)).

**Reasoning / Evidence:** NN/g's empty-state research: a blank container is ambiguous — users can't distinguish "no data," "still loading," and "error," so they refresh, wait, or leave; a plain status line ("No records for the selected date range") alone restores confidence. Beyond status, empty states are prime onboarding real estate: contextual "pull revelations" ("Star your favorites to list them here" — DataDog; Power BI's recent-items explanation) teach features at exactly the moment of relevance, outperforming front-loaded tutorials. Direct pathways (a Create button, demo data — Loggly's dual path) convert dead space into activation. Worst pattern: showing a false "No records" while data is still loading — it teaches users to distrust every message in the product.

**When to use:**
- Every container that can be empty: lists, dashboards, inboxes, search results, favorites, notifications, filtered tables. Enumerate them during design — each needs a designed state.

**When NOT to use / exceptions:**
- While loading, show a loading indicator/skeleton — never the empty-state message; render the empty state only after the fetch resolves empty.
- Don't oversell in cleared states: "Inbox zero" celebration is fine once; heavy illustration + upsell in a routine emptied list is clutter.
- Illustration is optional garnish; never let it displace status + path-forward on small screens.

**Pros:**
- Converts confusion into orientation, onboarding, and activation; cheap to build relative to impact on first-run retention.

**Cons:**
- Each container × state variant multiplies design/QA surface; copy must stay in sync as features change.

**Implementation guidance:**
- Template: [Status: what this area shows and that it's empty] + [Why/how it fills] + [Primary action] + [optional Learn-more link]. Example (first-use alerts panel): "No alerts yet. Alerts notify you when metrics cross thresholds you set. [Create alert] · Learn more."
- First-use: benefit-first framing, one primary CTA in the empty area itself (not only in distant chrome); consider demo/sample data for complex tools with a clear "sample" label.
- Cleared: brief positive confirmation ("All caught up — no pending reviews") + what will repopulate it; no CTA needed.
- No-results: echo the query/filters, state no match, offer repair (edit query, remove the most restrictive filter, spelling suggestion) — details in [08-navigation-ia.md#no-results-recovery](08-navigation-ia.md#no-results-recovery).
- Keep copy ≤2 short sentences + CTA; headline ≤6 words; match tone to context (first-use may be warm; error-adjacent must be plain).

**Related:** [#onboarding-copy-and-tooltips](#onboarding-copy-and-tooltips), [08-navigation-ia.md](08-navigation-ia.md) (no-results), [01-core-principles.md](01-core-principles.md) (visibility of system status).

**Sources:** [NN/g: Designing Empty States in Complex Applications](https://www.nngroup.com/articles/empty-state-interface-design/), [NN/g: No Results Found SERP Design](https://www.nngroup.com/articles/search-no-results-serp/), [Material Design: Empty states](https://m2.material.io/design/communication/empty-states.html).

---

## Mechanics and formatting

### Capitalization conventions

**Definition:** The system-wide choice between sentence case ("Create new project") and title case ("Create New Project") for headings, buttons, labels, and navigation — plus the prohibition of ALL CAPS for running labels.

**Reasoning / Evidence:** Sentence case has become the dominant modern convention: GOV.UK mandates it everywhere (research: easier to read and scan; title case creates ambiguity with proper nouns — is "New Project Settings" a feature name?); Material Design, Polaris, Carbon, and most contemporary systems specify sentence case for all UI text. Apple platforms traditionally use title case for buttons/menus, so native-feel iOS/macOS apps may follow the platform. ALL CAPS slows reading of continuous text by roughly 13–20% in classic (Tinker-era) legibility studies and reads as shouting. On mechanism: the traditional "word shape" explanation (caps erase ascender/descender outlines) is contested in modern reading science; the better-supported account is reduced familiarity — readers simply have far less practice with all-caps text. The practical guidance is unchanged: avoid ALL CAPS for anything longer than a short label. NN/g's tabs guidance explicitly rejects all-caps labels.

**When to use:**
- Sentence case: default for everything — headings, buttons, menu items, labels, tabs, empty states (web products, Android, cross-platform).
- Title case: only when matching Apple-platform conventions or house editorial style for publication headlines.
- ALL CAPS: at most tiny meta-labels (badges "NEW", table column eyebrows) with letter-spacing, never for sentences or multi-word actions.

**When NOT to use / exceptions:**
- Proper nouns, trademarks, and acronyms keep their casing regardless of convention.
- Never mix systems: one convention product-wide, encoded in the design system ([10-design-systems.md](10-design-systems.md)).

**Pros (sentence case):**
- Faster to read/scan; no per-word capitalization judgment calls (which minor words to cap?); localizes trivially (most languages have no title case).

**Cons (sentence case):**
- Can look informal against Apple-native chrome; long-established brands with title-case identity face migration inconsistency.

**Implementation guidance:**
- Write the rule into the content style guide with examples per component (button, heading, tab, tooltip, menu).
- Enforce via component defaults where safe (never CSS `text-transform: capitalize` — it caps every word including "of/and" and breaks non-English).
- If using `text-transform: uppercase` for badge micro-labels, add `letter-spacing: 0.5–1px` and keep to 1–2 words.
- Don't capitalize to add emphasis mid-sentence; use bold sparingly instead ([#scannability-headings-bullets-bold](#scannability-headings-bullets-bold)).

**Related:** [03-visual-design.md](03-visual-design.md) (typography), [10-design-systems.md](10-design-systems.md) (content standards), [#internationalization-and-localization](#internationalization-and-localization).

**Sources:** [GOV.UK style guide: capitalisation](https://www.gov.uk/guidance/style-guide/a-to-z#capitalisation), [Material Design 3: Applying type — capitalization](https://m3.material.io/styles/typography/applying-type), [Shopify Polaris: Grammar and mechanics](https://polaris.shopify.com/content/grammar-and-mechanics), [NN/g: Tabs, Used Right — no ALL CAPS](https://www.nngroup.com/articles/tabs-used-right/).

### Numbers, dates, and units formatting

**Definition:** Consistent, locale-aware rendering of quantities, dates/times, currencies, file sizes, and durations — the highest-frequency "content" in data-heavy products.

**Reasoning / Evidence:** Ambiguous formats cause real errors: "03/04/2025" is March 4 (US) or April 3 (most of the world); a comma is a thousands separator in English and the decimal mark in most of Europe. Standards exist precisely to remove ambiguity (ISO 8601 for machine dates; CLDR locale data powering `Intl` APIs). GOV.UK's style research favors unambiguous human formats ("4 March 2025") over numeric ones. Relative times ("2 hours ago") match how people reason about recency but decay into vagueness beyond a few days and are useless for records.

**When to use:**
- Everywhere numbers/dates appear; centralize through locale-aware formatting utilities (`Intl.NumberFormat`, `Intl.DateTimeFormat`, ICU), never string concatenation.

**When NOT to use / exceptions:**
- Logs, filenames, APIs, exports: ISO 8601 (`2025-03-04T14:32Z`) regardless of UI locale.
- Relative timestamps: use only for <7-day recency, always with an absolute value on hover/tap or beside it in audit-grade contexts.

**Pros:**
- Removes a whole class of misreading errors; free correctness via platform i18n APIs.

**Cons:**
- Locale plumbing must exist early; retrofit is painful; mixed-locale teams need explicit test matrices.

**Implementation guidance:**
- Dates: spell the month where space allows ("4 Mar 2025" / "March 4, 2025" per locale); never all-numeric ambiguous formats in UI. Include the year unless current-year-and-obvious.
- Times: respect the user's 12/24-hour preference; show timezone when events cross zones ("14:00 CET"); store UTC, display local.
- Numbers: locale separators via APIs; right-align numeric table columns with tabular (monospaced) figures; consistent decimals per column.
- Large numbers: abbreviate in dashboards ("1.2M") with full value on hover; never abbreviate in financial/legal records.
- Currency: symbol + ISO code where multiple currencies coexist ("$1,200 USD"); position per locale rules (CLDR handles it).
- Units: non-breaking space between value and unit ("12 GB"); metric/imperial per locale with a user override; file sizes with 1 decimal max ("3.4 MB").
- Durations: "2h 34m" style for compact UI; avoid "0:02:34" ambiguity (hours or minutes leading?).
- Server-rendered output: never format via the *ambient* runtime locale (`toLocaleDateString()` with no locale argument, `Intl.Collator(undefined)`) — the server's locale and ICU build differ from the browser's, producing hydration mismatches and content that flickers on load. Pass the user's resolved locale explicitly through a shared formatting utility.

**Related:** [#internationalization-and-localization](#internationalization-and-localization), [03-visual-design.md](03-visual-design.md) (tabular figures), [12-data-tables-dashboards.md](12-data-tables-dashboards.md) (data tables).

**Sources:** [GOV.UK style guide: dates, times, numbers](https://www.gov.uk/guidance/style-guide/a-to-z#dates), [Unicode CLDR](https://cldr.unicode.org/), [MDN: Intl API](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl), [ISO 8601](https://www.iso.org/iso-8601-date-and-time-format.html).

### Truncation and line length

**Definition:** Policies for text that exceeds its container (ellipsis truncation, wrapping, fading, middle-truncation) and for optimal line length in reading contexts.

**Reasoning / Evidence:** Typographic research puts comfortable reading measure at ~45–75 characters per line (Bringhurst; web studies converge on 50–75); very long lines lose the return sweep, very short lines chop phrasing. Truncation trades space for information loss: end-ellipsis destroys exactly the distinguishing tail of similar items (versions, filenames differing at the end), which is why file systems middle-truncate; NN/g and Baymard document truncated product titles causing wrong selections. Data disappearing behind ellipsis without a hover/expand path is invisible to touch users.

**When to use:**
- Truncate only where space is genuinely fixed (table cells, list rows, breadcrumbs, tabs) AND full text is recoverable (tooltip on hover/focus, expand on tap, detail view).
- Wrap by default for body content, labels in cards, and anything users must fully read.

**When NOT to use / exceptions:**
- Never truncate: error messages, legal text, confirmation consequences, button labels (redesign instead), amounts.
- Don't clamp headings to 1 line at the cost of meaning; 2-line clamp usually preserves scent.
- Numbers must never be visually truncated ("1,000,00…") — resize or abbreviate deliberately instead.

**Pros:**
- Preserves layout density and alignment in tables/lists; predictable row heights aid scanning.

**Cons:**
- Hidden information, wrong-item selection risk, no hover on touch, screen readers read the full string while sighted users see less (verification mismatch).

**Implementation guidance:**
- Line length: 50–75 characters (≈ `max-width: 60–70ch`) for body text; 80–100ch acceptable for reference tables of short rows.
- End-truncate descriptive text; middle-truncate identifiers whose tails distinguish ("Quarterly…v3.xlsx"); show full text via `title`/tooltip and ensure a tap path on touch.
- Multi-line clamp: 2–3 lines for card titles/descriptions (`-webkit-line-clamp`), then ellipsis.
- Design with worst-case strings: German compound words, 30–40% expansion, long emails/URLs; test truncation at 320px width.
- Prefer layout fixes over truncation: smaller type tier, column priority (drop low-value columns first), wrapping with hanging indent.

**Related:** [03-visual-design.md](03-visual-design.md) (typography, measure), [#internationalization-and-localization](#internationalization-and-localization) (expansion), [08-navigation-ia.md](08-navigation-ia.md) (breadcrumbs), [12-data-tables-dashboards.md](12-data-tables-dashboards.md) (tables).

**Sources:** [Baymard: Truncation UX](https://baymard.com/blog/truncation-design), Bringhurst, *The Elements of Typographic Style* (measure: 45–75 cpl), [NN/g: Legibility, Readability, and Comprehension](https://www.nngroup.com/articles/legibility-readability-comprehension/).

---

## Structure and scannability

### Content hierarchy and front-loading

**Definition:** Ordering content by user importance: conclusions and key facts first, details after (inverted pyramid, from journalism); within sentences/headings/links, the information-carrying words come first (front-loading).

**Reasoning / Evidence:** NN/g eye-tracking: users read ~20–28% of words on an average page and scan in F-patterns — attention concentrates at the tops of pages and the first ~2 words of lines/links/headings. The inverted pyramid guarantees the majority-who-leave-early still got the point; front-loaded microcontent survives truncation, scanning, and screen-reader link lists. Journalism's century of practice plus NN/g's replications across decades make this among the most durable findings in web content.

**When to use:**
- Every page, article, help doc, release note, email, and notification: answer first, background later.
- Every heading, link, button, list item, subject line, and push notification: keyword-first word order.

**When NOT to use / exceptions:**
- Narrative/persuasive marketing may build tension deliberately — but even landing pages need the value proposition in the first viewport.
- Legal documents follow mandated structure; add a plain-language summary block up top instead.

**Pros:**
- Serves scanners and completionists simultaneously; improves SEO (engines weight early words); truncation-resilient.

**Cons:**
- Writers trained on academic build-up structure must invert habits; burying-the-lede is the most common content-review finding.

**Implementation guidance:**
- One-sentence summary (the "dek") under every long-form H1; first paragraph must stand alone as the complete short answer.
- Headings: front-load the topic noun ("Refund processing times," not "How long does it take to process refunds?" — unless matching users' literal search phrasing).
- Links/buttons: first 2 words carry the differentiator ([#link-text](#link-text)).
- Emails/notifications: outcome in the subject/title; details in body.
- Test: cover everything below the first paragraph — does a user still get the essential answer? If no, restructure.

**Related:** [#scannability-headings-bullets-bold](#scannability-headings-bullets-bold), [#plain-language](#plain-language), [02-cognitive-foundations.md](02-cognitive-foundations.md) (attention).

**Sources:** [NN/g: Inverted Pyramid — Writing for Comprehension](https://www.nngroup.com/articles/inverted-pyramid/), [NN/g: How Users Read on the Web (they don't)](https://www.nngroup.com/articles/how-users-read-on-the-web/), [NN/g: F-Shaped Pattern of Reading](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/), [NN/g: First 2 Words — A Signal for Scanning](https://www.nngroup.com/articles/first-2-words-a-signal-for-the-scanning-eye/).

### Scannability: headings, bullets, bold

**Definition:** Formatting long-form content for non-linear reading: descriptive headings every few paragraphs, bulleted lists for parallel items, bold on key phrases, short paragraphs, and generous whitespace.

**Reasoning / Evidence:** Nielsen's 1997 study (still replicated): "concise, scannable, objective" rewrites improved measured usability by 124% over promotional wall-of-text baselines; scannable formatting alone contributed ~47%. Mechanism: headings and bullets create entry points that let the F-pattern eye land on relevant content; bold keywords act as anchors for the scanning eye. Over-formatting destroys the effect — if everything is bold, nothing is.

**When to use:**
- Help docs, articles, settings descriptions, release notes, long emails — any content >~150 words.

**When NOT to use / exceptions:**
- Short microcopy needs no internal structure; a 2-sentence dialog with a bulleted list is over-engineered.
- Bullet abuse: sequential instructions belong in *numbered* lists; prose arguments with logical flow shouldn't be shredded into disconnected fragments.
- Bold ≤1 phrase per paragraph; never bold whole sentences/paragraphs.

**Pros:**
- Serves the scanning majority; helps screen-reader users too (heading navigation is the #1 SR browsing strategy — WebAIM surveys: ~68% navigate by headings first).

**Cons:**
- Fragmentation can hide argument structure; heading noise ("Introduction," "Overview") wastes the entry-point slots.

**Implementation guidance:**
- One idea per paragraph; paragraphs ≤3–4 sentences on the web.
- A meaningful heading every 2–4 paragraphs; headings must be descriptive standing alone (they're read out of context in SR lists and TOCs); proper nesting H1→H2→H3, no skips.
- Bullets: 3–7 items ideal; parallel grammatical structure; front-load each item; sub-bullets max 1 level.
- Numbered lists exclusively for sequences/rankings; name the count in the lead-in ("Three steps:").
- Bold sparingly for the scannable keywords a skimmer needs; avoid italics for emphasis on screens (weak salience, worse legibility at small sizes); never underline non-links.

**Related:** [#content-hierarchy-and-front-loading](#content-hierarchy-and-front-loading), [03-visual-design.md](03-visual-design.md) (type hierarchy), [05-accessibility.md](05-accessibility.md) (heading semantics).

**Sources:** [NN/g: How Users Read on the Web](https://www.nngroup.com/articles/how-users-read-on-the-web/), [Nielsen: Concise, Scannable, Objective (1997 study)](https://www.nngroup.com/articles/concise-scannable-and-objective-how-to-write-for-the-web/), [WebAIM: Screen Reader User Survey](https://webaim.org/projects/screenreadersurvey10/), [GOV.UK content design guidance](https://www.gov.uk/guidance/content-design/writing-for-gov-uk).

---

## Global content

### Internationalization and localization

**Definition:** Internationalization (i18n) = building the product so it *can* be adapted (no hardcoded strings, locale-aware formatting, flexible layouts, RTL support). Localization (l10n) = actually adapting content per locale (translation, formats, imagery, regulation). Writing-side implications: expansion-tolerant copy, no idioms, no string concatenation, locale-correct names/dates/colors.

**Reasoning / Evidence:** Hard numbers drive the layout rules: English→German/French/Portuguese routinely expands text 30–40% (short strings expand proportionally more — a 10-char English label can expand 100–300%, i.e., double or more); English→Finnish or →Tamil differ wildly; CJK contracts horizontally but needs taller line boxes; German unbroken compounds break word-wrap assumptions. RTL scripts (Arabic ~400M speakers, Hebrew, Farsi, Urdu) mirror not just text but layout direction, icons with directionality, and progression metaphors. Concatenated strings ("You have " + n + " new " + itemType) are untranslatable: word order, gender agreement, and plural systems vary (Slavic languages have 3–4 plural forms; Arabic has 6 — CLDR plural rules exist for this). W3C and every major localization guide (Material, Microsoft, Mozilla) codify these.

**When to use:**
- From day one on any product with plausible multi-locale future — retrofitting i18n is among the most expensive rewrites.
- Copywriting for translatable products: simple sentence structures, no idioms/wordplay/sports-or-military metaphors, consistent terminology (one term per concept — translation memory rewards it).

**When NOT to use / exceptions:**
- Truly single-locale internal tools may defer l10n but should still avoid concatenation and hardcoded formats (cheap insurance).
- Brand names and some technical terms stay untranslated; maintain a "do not translate" glossary.

**Pros:**
- Access to global markets; forced string externalization also improves copy governance and consistency.

**Cons:**
- 30–40% layout headroom constrains dense designs; translation cost scales with string count and churn; RTL mirroring doubles visual QA.

**Implementation guidance:**
- Layout: design at +40% width headroom for labels/buttons; test with pseudo-localization (e.g., "Àççôûñţ Šéţţîñĝš~~~") in CI screenshots; never fix button widths to English text.
- Strings: externalize all copy; use ICU MessageFormat for plurals/genders/placeholders; never concatenate fragments; give translators context notes and screenshots per string.
- Audit the strings a visible-copy review misses — hardcoded source-language text accumulates exactly where nobody looks: `aria-label`s, tooltips, toast/notification copy, input placeholders, validation/error messages, empty-state text, and legal/consent text. A lint rule banning bare string literals in templates and attributes catches these mechanically; a copy review does not.
- RTL: use logical CSS properties (`margin-inline-start`, not `margin-left`); mirror layout, chevrons, back arrows, and progress direction; do NOT mirror clocks, media-playback icons, checkmarks, or logos; test with `dir="rtl"` early.
- Formats: all dates/numbers/currency via CLDR-backed APIs ([#numbers-dates-and-units-formatting](#numbers-dates-and-units-formatting)); name fields as single "Full name" or locale-flexible given/family order without assuming Western order ([07-forms-and-input.md](07-forms-and-input.md)); addresses per-country schemas.
- Fonts: verify glyph coverage for target scripts; CJK needs larger minimum sizes (~12px+ becomes illegible; prefer 14px+); line-height +0.2 for diacritic-heavy scripts.
- Process: budget l10n QA per locale per release; watch string churn (every English edit costs × N languages).

#### Data-format matrix

Never assume first/last name, US addresses, Gregorian-only dates, or one phone format. Per data type:

| Data | Rule | Failure to avoid |
|---|---|---|
| Personal names | Flexible fields: accept single names (mononyms), no middle name, non-Latin scripts; avoid mandatory first/last split globally — prefer a single "Full name" field or locale-flexible given/family order | Rejecting valid names ("last name required," Latin-only validation) |
| Addresses | Country-specific formats or a flexible address component; postal code optional where countries lack them | US-only schema (mandatory state, 5-digit ZIP) rejecting valid international addresses |
| Dates | Locale-aware display with spelled month; explicit input hints or date pickers | Ambiguous all-numeric dates — "01/02/03" is three different dates depending on locale |
| Times | Time zone visible wherever events cross zones | Scheduling across zones without showing TZ |
| Currency | Locale-aware symbol *placement* (before/after amount varies by locale) plus ISO code where currencies coexist ("$1,200 USD") | Bare "$" ambiguity (USD? CAD? AUD? MXN?) |
| Numbers | Locale separators and decimal marks for both display *and parsing* | Parsing "1,5" as fifteen instead of one-and-a-half (decimal comma locales) |
| Units | Local units or conversion with user override | Forcing unfamiliar units |
| Phone | Country code and flexible formatting | Numeric-only, fixed-length local assumption |

#### Culturalization checklist

Beyond translation, audit per target market:
- **Color meanings:** white = mourning in parts of East Asia; red = danger in the West but prosperity in China — don't rely on one culture's color semantics for status.
- **Icons and metaphors:** hand gestures, animals, and body-part imagery carry different (sometimes offensive) meanings; avoid mailbox/piggy-bank-style culture-bound metaphors.
- **Imagery:** people, clothing, food, and settings in photos/illustrations should fit the market or stay neutral; avoid flags to represent languages (Spanish ≠ Spain's flag for Latin American users).
- **Calendars and holidays:** don't assume Gregorian-only (Hijri, Buddhist, Japanese era calendars are official in some locales); seasonal campaigns invert hemispheres.
- **Honorifics and formality:** formal/informal address (T-V distinction, Japanese honorifics) must be a deliberate, documented tone decision per locale.
- **Privacy, consent, and legal language:** expectations and required wording differ by jurisdiction.
- **Payment and delivery norms:** expected methods and address/delivery conventions vary by market.

#### Layout and QA tests

- Test with the *longest target locale's* strings (German/Finnish-like) and with compact CJK strings — both directions break layouts.
- RTL snapshot test: render key screens with `dir="rtl"` in visual-regression CI; verify navigation, icons with directionality, and progress metaphors mirror correctly (and that clocks/playback/checkmarks do not).
- Pseudo-localization pass in CI screenshots (expansion + accents + brackets to catch concatenation and clipping).
- No text baked into images; anything with words must be a translatable string.

**Verification checklist:**
- UI survives 30–50% text expansion (100–300% for short labels) without truncation or overlap.
- Dates, times, numbers, and currency are locale-aware end to end (display and parsing).
- RTL layout does not break navigation, icons, or reading order.
- Forms accept valid international data (names, addresses, phones) — see [#data-format-matrix](#data-format-matrix).
- No string is concatenated from translated fragments; plurals/gender use ICU MessageFormat or equivalent.
- No bare string literals in templates or attributes — all user-facing strings, including `aria-label`s, tooltips, placeholders, validation/error messages, empty states, and toasts, go through the i18n layer (lint-enforceable).
- Every key exists in every supported locale (completeness check in CI).

**Related:** [#numbers-dates-and-units-formatting](#numbers-dates-and-units-formatting), [#truncation-and-line-length](#truncation-and-line-length), [#inclusive-language](#inclusive-language), [03-visual-design.md](03-visual-design.md) (color meaning), [07-forms-and-input.md](07-forms-and-input.md).

**Sources:** [W3C Internationalization: Text size in translation](https://www.w3.org/International/articles/article-text-size), [Material Design: Bidirectionality](https://m2.material.io/design/usability/bidirectionality.html), [Unicode CLDR Plural Rules](https://cldr.unicode.org/index/cldr-spec/plural-rules), [Microsoft: Localization guidelines](https://learn.microsoft.com/en-us/globalization/localization/localization), [W3C: Personal names around the world](https://www.w3.org/International/questions/qa-personal-names).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Plain language | Common words, active voice, "you"; experts prefer it too | ≤25 words/sentence, ~15 avg |
| Reading level | Grade ≤8 general; ≤6–7 for critical flows; formulas as smoke alarms | ~50% of US adults read ≤ grade 8; Flesch ≥60; GOV.UK targets a 9-year-old reading age |
| Voice & tone | One voice, tone flexes to stakes: errors serious, success may be light | 4 tone dimensions; ≤1 exclamation/screen |
| Terminology | One concept = one term everywhere; glossary + linter; descriptive over branded names | 1 verb per action; rename in 1 release |
| Inclusive language | Neutral terms, no idioms/ableist metaphors; lint with a word list | they/their default; allowlist/blocklist |
| Button labels | Verb-first, outcome-specific; never Yes/No/OK for consequences | 1–3 words; ≤~25 chars pre-translation |
| Link text | Descriptive noun phrase, front-loaded, unique among siblings; never "click here" | Link 2–5 words; flag (PDF)/(new tab) |
| Placeholder/helper | Required info in persistent helper text, never in placeholders; visible placeholder text still owes 4.5:1 | Helper ≤1 line (~80 chars); `aria-describedby` |
| Confirmations | Rare, specific, verb-labeled buttons; undo for frequent reversible acts | Undo window 5–10s; type-to-confirm only for catastrophic |
| Toast copy | Outcome first, ≤1 line; never actionable-critical info in auto-dismiss | ≥5s duration; push title ≤40 chars |
| Onboarding/tooltips | Contextual beats front-loaded tours; skip always available | ≤3 coach marks/release; tooltip delay 300–500ms |
| Error messages | What happened + why + how to fix; no blame, jargon, caps, or humor | Adjacent to source; ban "invalid/illegal"; never color alone (~300M with CVD) |
| Empty states | Status + how it fills + direct CTA; three types (first-use/cleared/no-results) | ≤2 sentences + CTA; never show during load |
| Capitalization | Sentence case everywhere (except Apple-native chrome); no ALL CAPS labels | ALL CAPS ~13–20% slower for continuous text; caps only 1–2-word badges + letter-spacing |
| Numbers & dates | Locale APIs, spelled month, ISO 8601 for machines; relative time <7 days only; pass the user's locale explicitly in server-rendered output | "4 Mar 2025"; tabular figures in tables |
| Truncation & measure | Wrap by default; truncate only with recovery path; middle-truncate IDs | 50–75 chars/line; clamp 2–3 lines |
| Hierarchy | Inverted pyramid; answer in first paragraph; front-load first 2 words | Users read ~20–28% of words |
| Scannability | Descriptive headings, parallel bullets, sparse bold | Heading per 2–4 ¶; 3–7 bullets; ¶ ≤4 sentences |
| i18n/l10n | +40% expansion headroom, no concatenation, ICU plurals, logical CSS for RTL; flexible names/addresses; no bare string literals (incl. aria-labels, toasts, placeholders) | 30–40% expansion (100–300% short strings); 6 Arabic plural forms |
