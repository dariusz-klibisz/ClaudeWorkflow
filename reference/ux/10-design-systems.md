# Design Systems

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers design systems as products: what they are and when they pay off, design tokens (tiers, naming, theming, standards), component library architecture (atomic design and alternatives, API design, composition vs configuration, variants/states), documentation, theming (dark mode, multi-brand, density), versioning and governance (contribution models, deprecation, adoption metrics), accessibility baked into components, notable public systems and what each is good for, design-dev handoff and Figma-code sync, and the adopt-vs-build decision. Visual style rules themselves (color, type, spacing scales) are in [03-visual-design.md](03-visual-design.md); component interaction behavior in [04-interaction-design.md](04-interaction-design.md); accessibility requirements in [05-accessibility.md](05-accessibility.md); content standards that live inside a system in [09-content-ux-writing.md](09-content-ux-writing.md); team process in [17-ux-process-research.md](17-ux-process-research.md).

## Contents

- [Foundations](#foundations)
  - [What a design system is and why](#what-a-design-system-is-and-why)
  - [Costs and when NOT to build one](#costs-and-when-not-to-build-one)
- [Design tokens](#design-tokens)
  - [Token tiers: primitive, semantic, component](#token-tiers-primitive-semantic-component)
  - [Token naming conventions](#token-naming-conventions)
  - [Theming via tokens](#theming-via-tokens)
  - [Token tooling and standards](#token-tooling-and-standards)
- [Component architecture](#component-architecture)
  - [Atomic design: critique and alternatives](#atomic-design-critique-and-alternatives)
  - [Component API design](#component-api-design)
  - [Composition vs configuration](#composition-vs-configuration)
  - [Variants and states](#variants-and-states)
  - [Accessibility baked into components](#accessibility-baked-into-components)
- [Documentation and operations](#documentation-and-operations)
  - [Documentation practices](#documentation-practices)
  - [Stable test hooks](#stable-test-hooks)
  - [Theming: dark mode, multi-brand, density](#theming-dark-mode-multi-brand-density)
  - [Versioning and deprecation](#versioning-and-deprecation)
  - [Governance and contribution models](#governance-and-contribution-models)
  - [Adoption metrics](#adoption-metrics)
- [Ecosystem and strategy](#ecosystem-and-strategy)
  - [Reference systems compared](#reference-systems-compared)
  - [Design-dev handoff](#design-dev-handoff)
  - [Figma-code sync trends](#figma-code-sync-trends)
  - [Adopt vs adapt vs build](#adopt-vs-adapt-vs-build)
- [Quick reference](#quick-reference)

---

## Foundations

### What a design system is and why

**Definition:** A design system is a complete set of standards intended to manage design at scale using reusable components and patterns (NN/g). Concretely: a token/style foundation, a component library (design + code), pattern guidance, content guidelines, and — critically — the team and processes that maintain them. It is a product serving other products, not a static deliverable. A design system encodes reusable *decisions* — about behavior, accessibility, content, and usage — not just visual assets.

**Reasoning / Evidence:** NN/g's Design Systems 101 identifies the value drivers: (1) **velocity** — replication of solved UI at near-zero marginal cost frees design/engineering for novel problems; (2) **consistency** — one source of truth prevents visual/behavioral drift across teams and channels, and enables org-wide rebrands at scale; (3) **quality** — accessibility, i18n, and interaction details solved once and inherited everywhere; (4) **shared language** — a controlled vocabulary ("dropdown" means one thing) removes a whole class of cross-functional miscommunication; (5) **onboarding** — explicit standards teach juniors and contractors. Industry case studies report 30–50% faster UI build time for screens composed from mature systems — treat this figure as practitioner folklore / vendor-published estimates without independent verification, though the direction (mature systems speed up composition) is consistent across reports.

**When to use:**
- Multiple teams shipping UI to the same product family or brand.
- Foreseeable years of replicable design work (NN/g's stated ROI condition).
- Products with hundreds+ of screens or several platforms needing coherence.
- Regulated/accessibility-committed orgs: centralizing compliance in components is the only strategy that scales.

**When NOT to use / exceptions:** see [#costs-and-when-not-to-build-one](#costs-and-when-not-to-build-one).

**Pros:**
- Compounding velocity; consistent multi-product UX; centralized quality (a11y, i18n, performance); rebrand leverage; education artifact.

**Cons:**
- Standing maintenance cost; can suppress justified divergence and experimentation if governance is rigid; a bad component propagates its badness everywhere ("consistency of error").

**Implementation guidance:**
- Treat it as a product: roadmap, versioned releases, support channel, user research where "users" are your product teams. Components are product infrastructure — fund and staff them accordingly.
- Minimum viable team (NN/g): 1 interaction designer + 1 visual designer + 1 developer; add part-time researcher, content designer; secure an executive sponsor.
- Extract, don't invent: harvest patterns from shipped product UI (proven use ≥3 places) rather than designing components speculatively.
- Start with the highest-frequency primitives: color/type/spacing tokens, button, input, select, checkbox/radio, modal, table, toast — these cover the bulk of any UI.
- Measure and report value ([#adoption-metrics](#adoption-metrics)) from day one — systems die when their funding case is vibes.

**Related:** [#costs-and-when-not-to-build-one](#costs-and-when-not-to-build-one), [#adopt-vs-adapt-vs-build](#adopt-vs-adapt-vs-build), [17-ux-process-research.md](17-ux-process-research.md), [01-core-principles.md](01-core-principles.md) (consistency).

**Sources:** [NN/g: Design Systems 101](https://www.nngroup.com/articles/design-systems-101/), Kholmatova, *Design Systems* (Smashing Magazine, 2017), [Sparkbox: Design Systems Survey](https://designsystemssurvey.seesparkbox.com/).

### Costs and when NOT to build one

**Definition:** The standing liabilities of owning a design system — maintenance, governance overhead, adoption friction — and the situations where building one destroys value.

**Reasoning / Evidence:** NN/g is explicit: design systems are "not a one-and-done solution"; they require a dedicated team, continuous evolution, and active teaching, or they decay into "unwieldy collections of components and code." The failure modes are well documented across the industry: zombie systems (built once, unfunded, drift from product reality, teams route around them); governance bottlenecks (every UI change queued behind a small core team); premature systematization (tokens and a 40-component library for a pre-product-market-fit startup whose UI will be thrown away); and adoption debt (a system nobody migrates to doubles the surface: old + new UI both live forever).

**When to use (i.e., when NOT to build):**
- Single small team (<~5 designers+engineers on UI): a lightweight style guide + a good open-source component library covers 90% of the benefit at 5% of the cost.
- Prototypes, MVPs, proofs of concept: NN/g — replicability value is in the future; a pivoting product has no future screens to amortize against.
- One-off projects (campaign sites, short-lived tools).
- No budget for a standing team: an unmaintained system is worse than none — it becomes misinformation with authority.

**When NOT to use / exceptions:**
- Even tiny teams should adopt token thinking (named scales instead of magic numbers) and one component library — the "no system" option still means conventions.
- A design-system *lite* (tokens + 10 core components + usage notes) is a legitimate middle path and a sound foundation to grow later.

**Pros (of deferring/not building):**
- Speed now; freedom to explore divergent UI; no governance tax; avoids betting standards on unvalidated product decisions.

**Cons (of deferring):**
- Inconsistency accretes; the eventual consolidation cost grows with every screen shipped; accessibility fixes stay per-screen.

**Implementation guidance:**
- Decision gate: build only when you can name (a) ≥2 consuming teams, (b) a funded owner team, (c) a multi-year product horizon. Fewer than that → adopt/adapt an existing system ([#adopt-vs-adapt-vs-build](#adopt-vs-adapt-vs-build)).
- Budget realistically: mature systems consume roughly 3–8 full-time people (design+eng+content) for a mid-size org — a practitioner rule of thumb without independent verification, but a useful planning anchor; the system team's cost must be < the duplication it removes.
- Sunset criteria: if adoption (share of product UI on system components) stays <30% after 18–24 months, restructure or fold the effort into a consuming team.
- Avoid the "big bang": don't freeze product work for a 12-month system build; extract incrementally alongside feature delivery.

**Related:** [#what-a-design-system-is-and-why](#what-a-design-system-is-and-why), [#governance-and-contribution-models](#governance-and-contribution-models), [#adoption-metrics](#adoption-metrics).

**Sources:** [NN/g: Design Systems 101 — "Why Not Use a Design System?"](https://www.nngroup.com/articles/design-systems-101/), [Sparkbox Design Systems Survey (maintenance/adoption findings)](https://designsystemssurvey.seesparkbox.com/), Kholmatova, *Design Systems* (2017).

---

## Design tokens

### Token tiers: primitive, semantic, component

**Definition:** Design tokens are named, platform-agnostic values for visual decisions (color, spacing, type, radius, shadow, motion, z-index). Mature systems layer them in three tiers: **primitive/global** (raw palette: `blue-600: #2563EB`, `space-4: 16px`), **semantic/alias** (meaning-bearing references: `color-action-primary → blue-600`, `color-text-danger → red-700`), and **component** (scoped decisions: `button-primary-bg → color-action-primary`).

**Reasoning / Evidence:** Tokens originated at Salesforce — Jina Anne and Jon Levine coined and popularized "design tokens" (~2014, Lightning Design System) — to solve multi-platform style drift. USWDS articulates the core rationale: CSS accepts ~16.7M colors and continuous values; unbounded choice slows design and makes designer-developer communication "unnecessarily granular" — tokens are "the presets on a car radio," a curated discrete palette. The tier separation is what makes theming and rebranding cheap: products reference semantic names; changing what `color-action-primary` points at re-skins everything without touching component code. Systems skipping the semantic tier (components referencing `blue-600` directly) rediscover this the hard way at their first dark mode or rebrand.

**When to use:**
- Any product beyond a prototype; the discipline costs little and prevents magic-number sprawl.
- Multi-platform (web/iOS/Android) or multi-theme products: tokens are the only sane single source of truth.

**When NOT to use / exceptions:**
- Don't tokenize everything: one-off decorative values in a marketing page don't need names; token bloat (thousands of unused names) is its own maintenance failure.
- Component tier is optional at small scale — many healthy systems run two tiers (primitive + semantic) and add component tokens only where per-component theming is truly needed.

**Pros:**
- Single source of truth across design tools and codebases; theming/rebrand leverage; enforced visual consistency; reviewable diffs for style changes.

**Cons:**
- Indirection cost: three hops from `button-primary-bg` to a hex value; naming governance is real work; tooling pipeline (build, sync, versioning) must be owned.

**Implementation guidance:**
- Keep primitives *meaning-free* (scale names: `gray-50…900`, `space-1…12`) and semantics *usage-named* (`color-bg-surface`, `color-text-secondary`, `color-border-focus`).
- Rule for consumers: product/component code references semantic (or component) tokens only — lint against direct primitive/hex usage.
- Scale sizes: color primitives typically 8–12 steps/hue; spacing on a 4px or 8px base scale (USWDS, Material); type in a 6–10 step modular scale.
- Cover the non-obvious categories early: focus ring, elevation/shadow, motion durations/easings, z-index — these drift worst without tokens.
- Version tokens as a package; changing a semantic token's *reference* is minor, renaming/removing a token is breaking ([#versioning-and-deprecation](#versioning-and-deprecation)).

**Related:** [#token-naming-conventions](#token-naming-conventions), [#theming-via-tokens](#theming-via-tokens), [03-visual-design.md](03-visual-design.md) (scales), [#token-tooling-and-standards](#token-tooling-and-standards).

**Sources:** [USWDS: Design tokens](https://designsystem.digital.gov/design-tokens/), [Salesforce Lightning: Design tokens](https://www.lightningdesignsystem.com/design-tokens/), [Material 3: Design tokens](https://m3.material.io/foundations/design-tokens/overview), [Nathan Curtis: Naming Tokens in Design Systems](https://medium.com/eightshapes-llc/naming-tokens-in-design-systems-9e86c7444676).

### Token naming conventions

**Definition:** A predictable grammar for token names, commonly ordered as: [namespace] – category/property – concept/role – variant – state – scale. Examples: `color-text-danger`, `button-primary-bg-hover`, `space-inset-lg`, `esds-color-action-primary-disabled`.

**Reasoning / Evidence:** Nathan Curtis's naming framework (EightShapes) is the de facto industry reference: names are the API of the token system — they must let a consumer *guess* the right token without a lookup, and let maintainers add tokens without collisions or synonym sprawl ("primary" vs "brand" vs "main" chaos). Consistent ordering enables tooling (sorting, grouping, find-and-replace across themes) and diff-able theme files.

**When to use:**
- Define the grammar before minting more than ~30 tokens; retro-renaming is a breaking change across every consumer.

**When NOT to use / exceptions:**
- Don't encode *values* into semantic names (`color-red-error` — what happens when error becomes orange?); value-ish names belong only in the primitive tier.
- Avoid encoding platform specifics in shared names; platform transforms handle unit differences.

**Pros:**
- Guessable, collision-free, tool-friendly; onboarding devs can navigate by pattern.

**Cons:**
- Long names (`color-background-interactive-primary-hover`); grammar debates consume real committee time — timebox them.

**Implementation guidance:**
- Fix an order and publish it: e.g., `[system]-[category]-[concept]-[property]-[variant]-[state]-[scale]`; omit levels that don't apply, never reorder.
- Controlled vocabularies per level: states = {hover, active, focus, disabled, selected, error}; emphasis = {primary, secondary, tertiary}; surfaces = {canvas, surface, overlay}. Publish the enum lists.
- Use T-shirt sizes (`xs…3xl`) or numeric scales (`100…900`) consistently — pick one per category, don't mix.
- kebab-case for CSS/web, transform to camelCase/snake_case per platform automatically (Style Dictionary does this).
- Reserve a namespace prefix (`acme-`) for tokens shipped to third parties to avoid collisions.

**Related:** [#token-tiers-primitive-semantic-component](#token-tiers-primitive-semantic-component), [#token-tooling-and-standards](#token-tooling-and-standards).

**Sources:** [Nathan Curtis: Naming Tokens in Design Systems](https://medium.com/eightshapes-llc/naming-tokens-in-design-systems-9e86c7444676), [USWDS: Design tokens — keys and values](https://designsystem.digital.gov/design-tokens/), [W3C Design Tokens Community Group format](https://design-tokens.github.io/community-group/format/).

### Theming via tokens

**Definition:** Implementing visual themes (light/dark, brands, densities, high-contrast) as alternative *value sets* bound to the same semantic token names — components never change; only the token resolution does.

**Reasoning / Evidence:** This is the payoff of the semantic tier: because `color-bg-surface` means "surface background" rather than "white," dark mode is a re-mapping (`→ gray-900`) rather than a rewrite. Systems that themed by forking components (a `<ButtonDark>`) or sprinkling conditionals drowned in combinatorics — every new theme multiplied maintenance. Material 3's dynamic color and Fluent's brand ramps both operate by generating primitive palettes and re-binding semantic roles.

**When to use:**
- Dark mode (now user-expected on all platforms), OS `prefers-color-scheme` respect.
- Multi-brand/white-label products; high-contrast accessibility themes; density modes.

**When NOT to use / exceptions:**
- Dark mode is not color inversion: shadows fail on dark (use lighter surfaces to signal elevation), saturated colors vibrate (desaturate ~1–2 steps), pure black backgrounds cause smearing on OLED and harsh contrast (prefer `#121212`-ish grays — Material guidance).
- Don't theme what shouldn't vary: focus indicators, error semantics, and brand-legal marks may be intentionally theme-invariant.

**Pros:**
- New theme ≈ one values file; theme QA reduces to token-mapping review + spot checks; user preference support nearly free after setup.

**Cons:**
- Every semantic token must have a defined value in every theme (coverage discipline); imagery/illustrations don't theme via tokens and need per-theme assets; contrast must be re-verified per theme.

**Implementation guidance:**
- Web mechanism: CSS custom properties scoped by `[data-theme="dark"]` or media query; never hardcode colors in component CSS.
- Author dark palettes deliberately: elevation = lighter surface (Material: overlay white at 5–16% by level); maintain WCAG contrast (≥4.5:1 text) in *both* themes — check every text/background token pair per theme in CI.
- Respect OS preference by default (`prefers-color-scheme`, `prefers-contrast`, `prefers-reduced-motion`), offer in-app override with a "match system" option, persist per user.
- Multi-brand: brands swap primitive palettes + type + radius tokens; semantic layer stays identical — this is the test of your tier hygiene.
- Test themes with real content screenshots (visual regression per theme), not just palette swatches.

**Related:** [#token-tiers-primitive-semantic-component](#token-tiers-primitive-semantic-component), [#theming-dark-mode-multi-brand-density](#theming-dark-mode-multi-brand-density), [03-visual-design.md](03-visual-design.md), [05-accessibility.md](05-accessibility.md) (contrast).

**Sources:** [Material Design: Dark theme](https://m2.material.io/design/color/dark-theme.html), [Material 3: Color system and dynamic color](https://m3.material.io/styles/color/system/overview), [web.dev: prefers-color-scheme](https://web.dev/articles/prefers-color-scheme), [USWDS: Design tokens](https://designsystem.digital.gov/design-tokens/).

### Token tooling and standards

**Definition:** The pipeline that stores tokens as data and distributes them to every platform: a source format (increasingly the W3C Design Tokens Community Group JSON spec), a transformer (Style Dictionary, Theo, Terrazzo), design-tool bindings (Figma Variables, Tokens Studio), and versioned distribution (npm packages, platform artifacts).

**Reasoning / Evidence:** Without a pipeline, "tokens" are just a spreadsheet that drifts. The W3C DTCG (Design Tokens Community Group, with participants from Figma, Adobe, Google, Microsoft, Salesforce) has been standardizing a JSON interchange format (`$value`, `$type`, aliases/references, composite types) so tools interoperate instead of each inventing a schema; Style Dictionary (Amazon) is the dominant open-source transformer, emitting CSS custom properties, Sass, iOS Swift, Android XML/Compose from one source. Figma Variables (2023+) brought token-like aliasing and modes (theming) natively into the design tool, closing part of the design-code gap.

**When to use:**
- The moment tokens serve ≥2 consumers (e.g., web + Figma, or two codebases): single-source + transform beats manual sync immediately.

**When NOT to use / exceptions:**
- A single web app can start with plain CSS custom properties in one file; add pipeline when a second consumer appears.
- The DTCG spec has long been in draft; treat it as the interchange direction, but expect tool-specific extensions.

**Pros:**
- One source of truth, N outputs; tokens become reviewable, diffable, CI-testable data; design-tool and code stay in sync mechanically.

**Cons:**
- Another build system to own; Figma Variables ↔ token-file round-tripping still needs plugins/glue (Tokens Studio, API scripts); composite tokens (typography, shadows) have patchier tool support.

**Implementation guidance:**
- Store tokens as JSON in the design-system repo (DTCG format or Style Dictionary's, with a migration path); PR review = design review for style changes.
- Transform via Style Dictionary to: CSS custom properties + TS constants (web), Swift (iOS), Kotlin/XML (Android), and a Figma-sync artifact.
- Bind Figma Variables to the same names; use Variable modes for themes; automate sync in CI rather than trusting manual updates.
- Add CI checks: every semantic token resolves in every theme; contrast pairs pass WCAG; no orphaned/unused tokens beyond a threshold.
- Publish as a versioned package (`@acme/tokens`) with a changelog; consumers pin and upgrade deliberately.

**Related:** [#token-naming-conventions](#token-naming-conventions), [#figma-code-sync-trends](#figma-code-sync-trends), [#versioning-and-deprecation](#versioning-and-deprecation).

**Sources:** [W3C Design Tokens Community Group: Format module](https://design-tokens.github.io/community-group/format/), [Style Dictionary](https://styledictionary.com/), [Figma: Variables documentation](https://help.figma.com/hc/en-us/articles/15339657135383-Guide-to-variables), [Tokens Studio](https://tokens.studio/).

---

## Component architecture

> **Cross-reference:** per-component behavioral contracts and overlay patterns (modal, drawer, popover, toast, tooltip…) live in [11-components-and-overlays.md](11-components-and-overlays.md); per-widget ARIA/keyboard specs in [06-aria-widget-reference.md](06-aria-widget-reference.md). This section covers the *architecture* of a component library, not individual component behavior.

### Atomic design: critique and alternatives

**Definition:** Brad Frost's atomic design (2013) decomposes UI into atoms (button, label) → molecules (search form) → organisms (header) → templates → pages, by chemistry metaphor. Alternatives organize by *usage semantics* instead: flat component lists with categories (inputs, navigation, feedback — what Material, Carbon, Polaris actually ship), or tiered "core/composite/pattern" models. The component/pattern distinction matters: **components** are self-contained, reusable parts (Button, Select, Modal); **patterns** are documented recipes that assemble components to do a job (a filter bar, a settings page, an error summary) and carry usage guidance and trade-offs rather than code alone.

**Reasoning / Evidence:** Atomic design's lasting contribution is compositional thinking: big UI is assembled from small tested parts. The critique from practice: the chemistry taxonomy is ambiguous at the boundaries (is a labeled input a molecule or organism? teams burn hours on classification with zero user value); the levels don't map to how consumers *search* for components (nobody looks for "molecules" — they look for "date picker"); and five levels overweight structure for systems that mostly need two (elements and patterns). Frost himself notes the mental model matters more than the vocabulary. Every major public system abandoned the chemistry labels in favor of functional categories, which is revealed preference at industry scale.

**When to use:**
- Use the *principle* (compose small→large, test small parts hard) universally.
- The named taxonomy can help teams new to systems thinking as a teaching scaffold.

**When NOT to use / exceptions:**
- Don't let taxonomy debates block shipping; if classification takes >5 minutes, the taxonomy is the problem.
- Don't expose chemistry terms to component-library consumers; organize docs by function and alphabet.

**Pros (atomic thinking):**
- Reuse, isolated testing, clear dependency direction (atoms never import organisms).

**Cons (atomic taxonomy):**
- Classification ambiguity; consumer-hostile navigation; encourages premature abstraction of "atoms" that have only one real use.

**Implementation guidance:**
- Organize the library in 2 tiers: **components** (self-contained, general: Button, Select, Modal, Table) and **patterns/compositions** (recipes assembling components for a job: filter bar, settings page, empty state). Add a third "primitives/utilities" tier (Box, Stack, Text) if using a styled-primitive approach.
- Pattern-tier coverage checklist: forms, tables, navigation, dialogs, empty states, errors, search, filtering, and onboarding are the patterns every product needs — document them even before componentizing them.
- Attach research notes and known trade-offs to each pattern (GOV.UK-style "when this was tested…"), and mark experimental patterns clearly so consumers know what's stable ([#governance-and-contribution-models](#governance-and-contribution-models) incubator tier).
- Docs navigation: alphabetical index + functional categories (actions, inputs, navigation, feedback, layout, data display) — mirror Carbon/Polaris.
- Enforce dependency direction with lint: components may not import from patterns; patterns may not import from product code.
- Promote a product-built composition into the system only after ≥3 independent usages (the "rule of three").

**Related:** [#component-api-design](#component-api-design), [#composition-vs-configuration](#composition-vs-configuration), [#documentation-practices](#documentation-practices), [11-components-and-overlays.md](11-components-and-overlays.md), [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Sources:** [Brad Frost: Atomic Design](https://atomicdesign.bradfrost.com/chapter-2/), [Brad Frost: Atomic Design and Storybook / extending atomic design](https://bradfrost.com/blog/post/extending-atomic-design/), [Carbon Design System: Components](https://carbondesignsystem.com/components/overview/components/), [Shopify Polaris: Components](https://polaris.shopify.com/components).

### Component API design

**Definition:** The public interface of a component — props/parameters, slots/children, events, and defaults — treated with the same rigor as any developer-facing API: predictable, minimal, consistent across the library, hard to misuse. The API is a *UX contract*: it should make correct UX easy and incorrect UX difficult.

**Reasoning / Evidence:** Component APIs are used thousands of times by developers who read docs once; inconsistency (one component's `isDisabled`, another's `disabled`, a third's `state="disabled"`) taxes every usage and breeds bugs. The "pit of success" principle: defaults should produce the correct, accessible result with zero configuration (a bare `<Button>` must be keyboard-operable and correctly styled), and dangerous options should be loud (`dangerouslySetInnerHTML`-style naming). Boolean-prop explosion (`primary`, `large`, `quiet`, each pairwise-conflicting) creates unrepresentable-state bugs; enum variants make invalid combinations unconstructible.

**When to use:**
- Every system component; write the API sketch (props table + usage examples) *before* implementation and review it like an RFC.

**When NOT to use / exceptions:**
- Internal sub-components need not carry full API polish; mark them private and exclude from docs.
- Escape hatches (`className`/`style` passthrough, `unstyled` mode) are pragmatic — include them deliberately with guidance, or teams fork components instead.

**Pros:**
- Predictability compounds: learning one component teaches the library; misuse-resistant APIs reduce a11y and state bugs at the source.

**Cons:**
- API rigor slows initial component delivery; changing a shipped API is a breaking change, so early mistakes are expensive ([#versioning-and-deprecation](#versioning-and-deprecation)).

**Implementation guidance:**
- Library-wide prop vocabulary (publish it): `variant` (visual style enum), `size` (`sm|md|lg`), `disabled`, `loading`, `error`, `onChange/onClick` event naming, controlled+uncontrolled input support (`value`/`defaultValue`).
- Prefer enums over booleans for mutually exclusive options: `variant="primary" | "secondary" | "danger"`, not three booleans.
- Defaults = the most common accessible case; require explicit opt-in for risky behavior (e.g., `<Modal dismissible={false}>` must be deliberate).
- Forward refs, spread ARIA/data attributes to the root element, and document which element receives them.
- Type strictly (TypeScript): make invalid prop combinations type errors where possible (discriminated unions for variant-dependent props).
- Guard runtime inputs in presentational components — compile-time types don't protect against runtime data. Clamp numeric props that reach inline styles, SVG geometry, or width calculations (percentages clamped to 0–100; counts checked for finiteness — `NaN`/negative values render broken bars and invalid markup), and render a safe fallback (or nothing) for unknown enum values rather than throwing during (server-side) render. Rule of thumb: for bad data, presentational components degrade — they don't crash the page.
- Keep required props ≤2; a component demanding 6 props to render will be wrapped or avoided.

**Verification checklist:**
- [ ] Bare/default usage renders a correct, accessible result (safe defaults).
- [ ] Accessible name is required where needed — omitting it (no visible label, no `ariaLabel`) is a type error, not a runtime surprise.
- [ ] Variants are an intentional, documented enum — not an open styling escape hatch.
- [ ] No style-only components that omit semantics (a "Button" rendering a bare `<div>` fails review).
- [ ] Props follow the published library-wide vocabulary; events use the standard naming.
- [ ] Numeric props feeding styles/geometry are clamped and finiteness-guarded; unknown enum values degrade to a documented fallback instead of throwing.

**Related:** [#composition-vs-configuration](#composition-vs-configuration), [#variants-and-states](#variants-and-states), [#accessibility-baked-into-components](#accessibility-baked-into-components).

**Sources:** [Shopify Polaris: Contributing — component API guidance](https://polaris.shopify.com/contributing), [Carbon Design System: Developing — component conventions](https://carbondesignsystem.com/developing/frameworks/react/), [React docs: Thinking in React / component design](https://react.dev/learn/thinking-in-react), Nathan Curtis, EightShapes: [Component QA in Design Systems](https://medium.com/eightshapes-llc/component-qa-in-design-systems-b18cb4decb9c).

### Composition vs configuration

**Definition:** Two API strategies for component flexibility. **Configuration**: one component with many props (`<Card title icon actions footerButtons […]>`). **Composition**: small parts assembled by the consumer (`<Card><Card.Header><Card.Title/></Card.Header><Card.Body/></Card>`), including compound-component and headless (logic-only, e.g., Radix, Headless UI, React Aria) approaches.

**Reasoning / Evidence:** Configuration optimizes the common case (one line, consistent output, easy to enforce standards) but hits the "prop apocalypse": every layout request adds a prop, the component becomes an unmaintainable God-object, and unsupported arrangements force forks. Composition keeps each part simple, makes arrangement the consumer's power, and matches how the underlying platforms work — the industry's headless-component wave (Radix/React Aria providing behavior+a11y, styling per-system) is composition taken to its logical end. Cost: composition permits *wrong* assemblies, so consistency depends on docs, lint rules, and recipes instead of API constraints.

**When to use (configuration):**
- Small, closed components with genuinely finite variation: Button, Badge, Checkbox, Spinner.
- Places where enforcing uniformity outweighs flexibility (e.g., a system-mandated page header).

**When to use (composition):**
- Container/layout components with unpredictable content: Card, Modal, Toolbar, Menu, Table.
- Systems serving many teams with divergent needs — composition absorbs variation without core changes.

**When NOT to use / exceptions:**
- Don't compose primitives so far that every button takes five elements to write; provide *both*: composable parts + preassembled convenience components for the common case ("recipes over configuration").
- Headless-only strategies still need a styled default layer, or every team restyles from scratch and drifts.

**Pros / Cons:**
- Configuration — Pros: terse, consistent, discoverable via props table. Cons: prop explosion, rigid layouts, fork pressure.
- Composition — Pros: flexible, small maintainable parts, no God-components. Cons: verbose, misassembly risk, higher consumer skill demand.

**Implementation guidance:**
- Default heuristic: configuration for leaf/closed components; compound composition for containers; headless layer only if you must serve multiple styling regimes.
- With compound components, validate assembly in dev mode (warn when `<Card.Header>` is used outside `<Card>`).
- Ship "recipes": documented, copy-pasteable assemblies for the top 5 use cases of each composable component — this recovers configuration's consistency.
- Watch the fork signal: when a team copies a component's source to change layout, the API chose configuration where composition was needed.

**Related:** [#component-api-design](#component-api-design), [#atomic-design-critique-and-alternatives](#atomic-design-critique-and-alternatives), [#accessibility-baked-into-components](#accessibility-baked-into-components).

**Sources:** [Radix Primitives: Composition philosophy](https://www.radix-ui.com/primitives/docs/overview/introduction), [React Aria: Architecture](https://react-spectrum.adobe.com/react-aria/why.html), Kent C. Dodds: [Inversion of Control / compound components](https://kentcdodds.com/blog/inversion-of-control), [Shopify Polaris: Components](https://polaris.shopify.com/components).

### Variants and states

**Definition:** **Variants** are design-time alternatives of a component (visual emphasis: primary/secondary/tertiary/danger; size: sm/md/lg; density). **States** are runtime conditions (default, hover, focus, active/pressed, disabled, loading, error, selected, readonly). Every variant must define every applicable state.

**Reasoning / Evidence:** The variant×state matrix is where systems silently fail: a danger-variant button without a defined focus style, a compact table row without a selected state — gaps surface as inconsistent ad-hoc fixes in product code. Figma's component-variant model and code enum props both exist to make the matrix explicit and enumerable; visual regression tools (Chromatic, Percy) test exactly this grid. Constraining variant count is governance: every added variant multiplies the matrix (variants × sizes × states × themes), and unconstrained "just add a variant" requests are how 14-variant buttons happen.

**When to use:**
- Define the canonical state list system-wide and require each new component/variant PR to fill the full matrix (design + code + docs).

**When NOT to use / exceptions:**
- Not every component has every state (a Badge has no hover); document *not applicable* explicitly rather than leaving holes.
- Reject variant requests that encode one team's context ("marketing-blue button") or lack a demonstrated product need; solve with tokens/theming or say no — see [#governance-and-contribution-models](#governance-and-contribution-models).

**Pros:**
- Complete matrices make quality testable (snapshot the grid) and behavior predictable everywhere.

**Cons:**
- Matrix size grows multiplicatively; each theme doubles visual QA surface.

**Implementation guidance:**
- Canonical interactive-state set: default, hover, focus-visible, active, disabled, loading; add error/selected/readonly for inputs and selectables. Focus style must meet WCAG 2.4.7/2.4.11 (visible, ≥3:1 contrast, not color-only).
- Emphasis variants: 3 levels (primary/secondary/tertiary) + danger covers nearly all needs; one primary button per view is a *usage* rule to document.
- Sizes: 3 (sm/md/lg) with touch-target floors — interactive targets ≥24×24px minimum (WCAG 2.5.8), 44–48px on touch ([15-platform-mobile.md](15-platform-mobile.md)).
- Disabled: prefer explaining *why* (tooltip/helper) or hiding over bare disabled controls; disabled elements are exempt from contrast rules but near-invisible grays still harm discoverability.
- Loading state must preserve layout (no width jump) and communicate via `aria-busy`/live region, not spinner alone.
- Build a "kitchen sink" story per component rendering the full variant×state×theme grid; wire to visual regression in CI.

**Verification checklist:**
- [ ] Every applicable state is defined in design, code, and docs; non-applicable states are explicitly marked N/A, not left blank.
- [ ] Focus-visible style present on every interactive variant and meets WCAG contrast.
- [ ] Loading preserves layout and announces (`aria-busy`/live region).
- [ ] Full variant×state×theme grid is rendered in a story and covered by visual regression.

**Related:** [#component-api-design](#component-api-design), [04-interaction-design.md](04-interaction-design.md) (state feedback), [05-accessibility.md](05-accessibility.md) (focus, contrast).

**Sources:** [Material 3: Interaction states](https://m3.material.io/foundations/interaction/states/overview), [Carbon: Button usage (variants/states)](https://carbondesignsystem.com/components/button/usage/), [WCAG 2.2: Focus Appearance & Target Size](https://www.w3.org/WAI/WCAG22/Understanding/), [NN/g: Button States 101](https://www.nngroup.com/videos/button-states-101/).

### Accessibility baked into components

**Definition:** Making the component library the primary delivery vehicle for accessibility: semantics, keyboard behavior, focus management, ARIA wiring, contrast-safe tokens, and reduced-motion support implemented once in components so product teams inherit compliance by default.

**Reasoning / Evidence:** Per-screen accessibility retrofits don't scale — the same widget bugs recur across hundreds of screens. Centralizing in components inverts the economics: fix a Modal's focus trap once, fix it everywhere. This is the explicit strategy of GOV.UK (components researched and tested with assistive technology, WCAG 2.2-compliant out of the box), USWDS (Section 508 mandate), and Carbon (per-component accessibility documentation). The W3C ARIA Authoring Practices Guide supplies the canonical keyboard/ARIA contract per widget type (tabs, dialog, combobox, menu…), and headless libraries (React Aria, Radix) exist mainly to package these contracts. Caveat proven repeatedly in audits: using an accessible library does not make a product accessible — composition-level failures (missing labels, bad heading order, focus order across components, custom overrides) remain the product team's job.

**When to use:**
- Every component, from day one — retrofitting a11y into a shipped library is a breaking-change minefield.
- Regulated contexts (government, EU EN 301 549, ADA-exposed commerce): a compliant library is table stakes.

**When NOT to use / exceptions:**
- Don't over-ARIA: no role is better than a wrong role; native HTML elements (`<button>`, `<select>`, `<dialog>`) carry semantics free — first rule of ARIA is don't use ARIA if native works.
- Component-level a11y cannot cover page-level obligations: landmarks, heading structure, focus order across widgets, alt text, error-summary patterns remain per-product.

**Pros:**
- Compliance economics inverted; a11y expertise concentrated where it exists; every consuming product's baseline rises with each library release.

**Cons:**
- Correct focus/ARIA behavior is genuinely hard (combobox, grid, date picker are months of edge cases — a reason to build on headless primitives); escape hatches let teams break inherited guarantees.

**Implementation guidance:**
- Implement to WAI-ARIA APG per widget: dialogs trap and restore focus; menus/tabs use roving tabindex + arrow keys; Escape closes layers; comboboxes follow the APG combobox contract — per-widget contracts in [06-aria-widget-reference.md](06-aria-widget-reference.md).
- Native elements first; ARIA only to fill gaps.
- Token-level guarantees: text/background semantic pairs ≥4.5:1 (3:1 for large text), UI component boundaries ≥3:1, focus ring token ≥3:1 against adjacents — verified in CI per theme.
- Respect `prefers-reduced-motion` inside components (disable non-essential animation), `prefers-contrast` where supported.
- Per-component a11y docs (Carbon-style): keyboard interaction table, screen-reader behavior, required labels/props (`ariaLabel` required when no visible label — make it a type error to omit).
- CI: axe-core automated checks on every story (catches ~30–40% of issue types), plus manual screen-reader passes (NVDA+Chrome, VoiceOver+Safari) per component release; document AT support matrix.

**Verification checklist (release gate per component):**
- [ ] Ships with all states defined ([#variants-and-states](#variants-and-states)).
- [ ] Full keyboard support per the APG contract for its widget type, documented in a keyboard interaction table.
- [ ] Documented usage rules: when to use, when not to use, required labels/props.
- [ ] axe-core passes on every story; manual screen-reader pass recorded (NVDA, VoiceOver).
- [ ] Contrast token pairs verified per theme; `prefers-reduced-motion` respected.

**Related:** [05-accessibility.md](05-accessibility.md) (full requirements), [06-aria-widget-reference.md](06-aria-widget-reference.md) (per-widget specs), [#component-api-design](#component-api-design), [#variants-and-states](#variants-and-states), [#reference-systems-compared](#reference-systems-compared).

**Sources:** [W3C WAI-ARIA Authoring Practices Guide](https://www.w3.org/WAI/ARIA/apg/patterns/), [GOV.UK Design System: Accessibility](https://design-system.service.gov.uk/accessibility/), [Carbon: Accessibility guidance](https://carbondesignsystem.com/guidelines/accessibility/overview/), [USWDS: Accessibility](https://designsystem.digital.gov/documentation/accessibility/), [Deque: axe-core](https://github.com/dequelabs/axe-core).

---

## Documentation and operations

### Documentation practices

**Definition:** The written layer that makes components usable: what each component is for, when to use it (and what to use instead), do/don't examples, content guidelines, props/API reference, live examples, and accessibility notes — with design guidance and code docs unified per component.

**Reasoning / Evidence:** NN/g: a design system undocumented "risks being applied inconsistently or incorrectly" — the artifact without the instructions fails its purpose. Best-in-class public systems converge on the same per-component template (evidence of what works): Polaris and Carbon pair usage guidance + content guidelines + code on one page; GOV.UK adds research findings ("when this component was tested…"); Material adds do/don't imagery. Do/don't pairs are disproportionately effective because they answer the question consumers actually have ("is *my* case right?") faster than prose. The most-cited documentation failure: design docs and code docs in different places drifting apart — consumers then trust neither.

**When to use:**
- Every published component/pattern/token set; documentation ships *with* the release, not after.

**When NOT to use / exceptions:**
- Don't document speculative components before they're stable — docs create adoption you'll break. Mark experimental components/patterns clearly so consumers can't mistake them for stable.
- Internal-only utilities can carry lighter docs (props table + one example).

**Pros:**
- Self-service adoption (the system team stops being a help desk); consistent usage; onboarding material for free (NN/g's education benefit).

**Cons:**
- Documentation is commonly estimated at 30–50% of the cost of shipping a component (practitioner folklore, not an independently verified figure — but directionally right: docs are a major cost center) and rots fastest; needs an owner and a review cadence.

**Implementation guidance:**
- Per-component template: one-line purpose → when to use / when NOT to use (link alternatives) → anatomy diagram → variants/states demo → do/don't pairs (≥2, with images) → content guidelines (label writing per [09-content-ux-writing.md](09-content-ux-writing.md)) → accessibility (keyboard table, ARIA) → props/API reference (generated from types, never hand-maintained) → live editable example.
- Single source: docs site renders the *actual coded components* (Storybook/Docusaurus + live playground); design specs link to the same page — one URL per component for both audiences.
- Auto-generate what can be generated: props tables from TypeScript, token tables from token JSON, changelog from commits.
- Write "use instead" cross-links aggressively (Modal vs Drawer vs Popover — see [11-components-and-overlays.md](11-components-and-overlays.md)) — component *choice* errors outnumber component *usage* errors.
- Review docs on every breaking release; run link checks and screenshot freshness checks in CI.

**Verification checklist (docs done per component):**
- [ ] When to use / when NOT to use, with linked alternatives.
- [ ] ≥2 do/don't pairs; content guidelines; anatomy diagram.
- [ ] Keyboard/ARIA section present; states demo covers the full matrix.
- [ ] API reference generated from types (not hand-maintained).
- [ ] Version and migration guidance linked from the component page.

**Related:** [#component-api-design](#component-api-design), [09-content-ux-writing.md](09-content-ux-writing.md), [#design-dev-handoff](#design-dev-handoff).

**Sources:** [NN/g: Design Systems 101](https://www.nngroup.com/articles/design-systems-101/), [Shopify Polaris component docs (template exemplar)](https://polaris.shopify.com/components), [Carbon: Content overview + component pages](https://carbondesignsystem.com/guidelines/content/overview/), [GOV.UK Design System (research notes per component)](https://design-system.service.gov.uk/components/).

### Stable test hooks

**Definition:** Treating test selectors as part of the component contract: interactive components expose stable, centrally documented test-hook attributes (commonly `data-testid`, following a published naming convention such as kebab-case `feature-element[-modifier]`), and the hooks are declared alongside the component (exported constants and/or the component's docs page) so tests reference declared hooks rather than guessing selectors. This is a design-system governance concern — the hook contract ships with the component like props and tokens do — not a build-tooling detail.

**Reasoning / Evidence:** Convention-without-enforcement fails in both directions — a recurring finding in large-scale code review of real products: components ship without hooks, so tests fall back to brittle text or CSS-class selectors that break on copy edits or restyling; and tests assert fabricated hook names that match nothing, silently skipping (or vacuously passing) forever. Cypress and Playwright both recommend dedicated test attributes decoupled from text and styling; Testing Library ranks `data-testid` as the last-resort query after role/label queries. Lint/CI enforcement (interactive components must carry hooks; test selectors must resolve to declared hooks) beats prose convention.

**When to use:**
- Systems with E2E/integration suites maintained across teams; any component whose visible text or class names change independently of behavior.

**When NOT to use / exceptions:**
- Prefer accessibility-first queries (role + accessible name) where unambiguous — they double as accessibility verification; test IDs are for elements without stable roles/names and for repeated instances.
- Don't add hooks to decorative or static content (headings, body text, icons) — query those by role or text.

**Pros:**
- Selectors survive copy and styling changes; hook names become shared vocabulary between component docs and test suites; missing or fabricated hooks are machine-detectable.

**Cons:**
- Extra attribute noise in markup (can be mitigated by stripping test IDs in production builds if desired); one more contract to version — renaming a hook is a breaking change for consumers' tests.

**Implementation guidance:**
- Publish one naming convention library-wide, e.g. kebab-case `feature-element[-modifier]` with the modifier distinguishing repeated instances (`result-row-${id}`).
- Declare hooks next to the component — exported constants or a documented table on the component page; tests import or reference declared names, never ad-hoc strings.
- Enforce both directions in CI: a lint rule requiring hooks on interactive elements, and test tooling that fails (rather than skips) when a selector resolves to nothing.
- Decide and document the production-stripping policy explicitly (keep hooks for supportability vs strip for markup hygiene) so teams don't decide per project.

**Related:** [#component-api-design](#component-api-design), [#documentation-practices](#documentation-practices), [#governance-and-contribution-models](#governance-and-contribution-models).

**Sources:** [Cypress: Best Practices — Selecting Elements](https://docs.cypress.io/app/core-concepts/best-practices#Selecting-Elements), [Playwright: Locators — locate by test id](https://playwright.dev/docs/locators#locate-by-test-id), [Testing Library: About Queries — Priority](https://testing-library.com/docs/queries/about/#priority), [MDN: `data-*` global attributes](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/data-*).

### Theming: dark mode, multi-brand, density

**Definition:** The user- and business-facing theme axes a system must support: color schemes (light/dark/high-contrast), brand skins (multi-tenant/white-label), and density modes (comfortable/compact spacing for information-dense contexts).

**Reasoning / Evidence:** Dark mode is now a platform-level user expectation (OS toggles, `prefers-color-scheme`) with real benefits (OLED power savings, low-light comfort, user preference ~evenly split in surveys) — products without it read as unmaintained. Multi-brand is an architecture question that decides acquisitions and white-label deals: systems with clean token tiers onboard a brand in days (semantic re-mapping), token-less systems in months. Density: enterprise users at data-heavy screens measurably benefit from compact modes (more rows in view = less scrolling/paging); Material 3, Fluent 2, and Carbon all ship density scales as first-class.

**When to use:**
- Dark mode: any consumer-facing product; developer tools (dominant preference); media/reading apps.
- Multi-brand: white-label products, post-acquisition portfolios, regional brand families.
- Density: data tables, admin consoles, trading/ops dashboards — offer compact; keep comfortable default.

**When NOT to use / exceptions:**
- Content-branding conflicts: photo/commerce imagery may look wrong on dark surfaces; brands may legally require specific chrome.
- Don't build multi-brand machinery speculatively for a single-brand product — it's among the priciest YAGNI violations.
- Density mode ≠ shrinking touch targets below minimums on touch devices; compact is a pointer-first mode.

**Pros:**
- User preference respect; enterprise-sales enablement (white-label, compact); accessibility wins (high-contrast themes, light-sensitivity).

**Cons:**
- Each axis multiplies the QA matrix (themes × densities × states × breakpoints); imagery/illustration needs per-theme variants; teams must design in ≥2 themes forever after.

**Implementation guidance:**
- All three axes via tokens ([#theming-via-tokens](#theming-via-tokens)): color themes re-map color tokens; density re-maps spacing/size tokens (e.g., row height 52→40→32px, control height 40→32px); brands re-map primitives + type + radius.
- Dark specifics: elevation as lighter surfaces (not shadow); desaturate accents 1–2 steps; avoid pure black (`#121212`+ grays); re-verify all contrast pairs; provide dark-variant illustrations/logos.
- Defaults: follow OS scheme; in-app override with "match system"; persist per account, sync across devices.
- Density: global user setting for workhorse apps plus per-table override; never mix densities within one view.
- Multi-brand: one brand = one token package; forbid brand conditionals inside component code (the moment `if (brand === 'x')` appears in a component, the architecture has failed).

**Related:** [#theming-via-tokens](#theming-via-tokens), [#token-tiers-primitive-semantic-component](#token-tiers-primitive-semantic-component), [03-visual-design.md](03-visual-design.md), [08-navigation-ia.md](08-navigation-ia.md) (table density), [12-data-tables-dashboards.md](12-data-tables-dashboards.md).

**Sources:** [Material Design: Dark theme](https://m2.material.io/design/color/dark-theme.html), [Fluent 2: Design tokens & theming](https://fluent2.microsoft.design/design-tokens), [Carbon: Themes](https://carbondesignsystem.com/guidelines/themes/overview/), [web.dev: prefers-color-scheme](https://web.dev/articles/prefers-color-scheme).

### Versioning and deprecation

**Definition:** Managing change in a system that other teams build on: semantic versioning of packages, deprecation policy (announce → warn → migrate → remove), codemods/migration tooling, and communication rhythms that let consumers upgrade on their own schedule without being stranded.

**Reasoning / Evidence:** A design system is a dependency; dependency-management economics apply. Breaking changes without migration paths teach teams to pin old versions forever — the system then supports N versions in the wild and consistency (its core value) is gone. Mature systems (Polaris, Carbon, MUI) converge on: semver, long deprecation windows with runtime/console warnings, codemods for mechanical migrations, and published support windows. The empirical failure mode is well documented in design-systems surveys: upgrade pain is a top-3 reason teams abandon systems.

**When to use:**
- From the first external consumer; versioning discipline is cheapest when started early.

**When NOT to use / exceptions:**
- Pre-1.0 internal phases can move fast with a declared "unstable" contract — but declare it explicitly.
- Visual-only token value tweaks (same name, new value) can ship as minor even though pixels change; define and publish what *you* count as breaking (API? pixels? DOM structure? a11y behavior?).

**Pros:**
- Consumers upgrade predictably; the system team can evolve without freezing; trust compounds.

**Cons:**
- Deprecation windows mean maintaining old+new in parallel (real cost); codemod authoring is engineering work per breaking change.

**Implementation guidance:**
- Semver strictly: breaking = major (removed/renamed props, tokens, DOM/ARIA structure changes that break tests/styling overrides); publish the breaking-change definition.
- Deprecation pipeline: mark in types (`@deprecated` + replacement) → console warning in dev builds → docs banner with migration recipe → remove after ≥2 minors or ≥6 months, whichever is longer.
- Deprecate outdated *patterns*, not just components: when guidance changes, mark the old pattern page superseded and link the replacement.
- Ship codemods (jscodeshift-style) for every mechanical rename/restructure; a breaking change without a codemod or written migration guide shouldn't ship.
- Release cadence: minors on a fixed rhythm (e.g., monthly), majors ≤1–2/year batched; changelog per release with upgrade-effort estimates.
- Support window: current major + previous major (security/critical fixes) is the common workable policy.
- Version design assets in lockstep: Figma library versions labeled to match npm majors; publish both together.

**Related:** [#governance-and-contribution-models](#governance-and-contribution-models), [#token-tooling-and-standards](#token-tooling-and-standards), [#adoption-metrics](#adoption-metrics).

**Sources:** [Semantic Versioning](https://semver.org/), [Shopify Polaris: Versioning and changelog practices](https://polaris.shopify.com/version-guides), [Carbon: Migration guides](https://carbondesignsystem.com/migrating/guide/overview/), [Sparkbox Design Systems Survey](https://designsystemssurvey.seesparkbox.com/).

### Governance and contribution models

**Definition:** Who decides what enters, changes, and leaves the system. Three canonical models: **centralized** (a dedicated team owns everything; consumers request), **federated/distributed** (practitioners from product teams govern jointly; no single owning team), and **hybrid** (core team owns foundations + pipeline; product teams contribute components through a defined process) — the dominant mature model.

**Reasoning / Evidence:** Nathan Curtis's team-model analysis (EightShapes) frames the trade: centralized teams produce coherent quality but become bottlenecks and drift from product reality ("ivory tower" — components that don't fit real use cases because the core team doesn't ship product); federated models stay grounded but struggle with consistency, ownership of maintenance, and "who fixes the bug" ambiguity. Industry surveys (Sparkbox) show most surviving systems run hybrid: a small core owning tokens/infrastructure/quality bars, with a documented contribution path that converts product teams from requesters into co-owners — which also drives adoption (people use what they helped build).

**When to use:**
- Centralized: early-stage systems (first 1–2 years), small orgs, or when quality/a11y baselines must be established fast.
- Federated: strong design-ops culture, senior practitioners, org allergy to central teams — rare as a pure model.
- Hybrid: default for orgs with ≥3 product teams once the foundation exists.

**When NOT to use / exceptions:**
- Contribution without a quality gate imports inconsistency into the source of truth — a contribution process is a *review* process, not an open door.
- Don't accept single-team-specific components into the core; park them in a documented "local/incubator" tier until a second team needs them. Same rule for variants: no one-off variants without a demonstrated multi-team product need.

**Pros / Cons:**
- Centralized — Pros: coherence, speed of decision, clear accountability. Cons: bottleneck, ivory-tower risk, scaling ceiling.
- Federated — Pros: grounded in product, shared ownership, no central headcount. Cons: slow consensus, maintenance orphans, inconsistency.
- Hybrid — Pros: grounded + coherent, adoption via participation. Cons: process overhead, contribution review capacity becomes the new bottleneck.

**Implementation guidance:**
- Publish the contribution ladder: propose (issue + evidence of ≥2–3 use cases) → design review → build to the quality bar (API review, a11y checklist, docs, tests, full state matrix) → core team merge + co-maintenance agreement.
- Decision forum: a small standing council (core team + rotating product reps) meeting biweekly; decisions logged publicly (ADR-style) — silence and private Slack decisions kill trust.
- Quality bar as checklist in the PR template; the bar is the same for core and contributed components.
- Incubator tier: let teams publish experimental patterns under a clearly labeled namespace; promote by the rule of three.
- Audit periodically: usage and accessibility of shipped components (not just new ones); feed regressions and unused components into the deprecation/redesign queue ([#adoption-metrics](#adoption-metrics)).
- Fund maintenance explicitly: every accepted contribution names an owning team for bugs until core adopts it.

**Related:** [#versioning-and-deprecation](#versioning-and-deprecation), [#adoption-metrics](#adoption-metrics), [17-ux-process-research.md](17-ux-process-research.md).

**Sources:** [Nathan Curtis: Team Models for Scaling a Design System](https://medium.com/eightshapes-llc/team-models-for-scaling-a-design-system-2cf9d03be6a0), [Shopify Polaris: Contributing](https://polaris.shopify.com/contributing), [Carbon: How to contribute](https://carbondesignsystem.com/contributing/overview/), [Sparkbox Design Systems Survey](https://designsystemssurvey.seesparkbox.com/).

### Adoption metrics

**Definition:** Quantifying whether the system is actually used and paying off: component coverage (share of product UI built from system components), token coverage (share of style values from tokens vs hardcoded), version currency (consumers on supported versions), detachment/override rates, and outcome proxies (build time, a11y defect rates, design QA findings).

**Reasoning / Evidence:** Systems compete for funding against features; without numbers, they lose. Adoption is also the leading indicator of the value thesis: consistency and velocity benefits only materialize on UI that *uses* the system. Instrumentation is tractable: static analysis counts component imports and raw hex/px values in code; Figma library analytics report component insertions and detach rates; telemetry can tag rendered components in production.

**When to use:**
- From first release; baseline before you evangelize, so improvement is attributable.
- Before/after case studies on pilot teams are the strongest funding argument.

**When NOT to use / exceptions:**
- Don't weaponize metrics against product teams (mandated 100% adoption breeds hostile compliance and hacky wrapping); metrics diagnose, incentives persuade.
- Raw component-count vanity metrics ("we shipped 60 components") measure supply, not value — measure *usage*.

**Pros:**
- Funding defense; prioritization signal (most-used components deserve the most investment; unused ones deserve deprecation review); early warning (rising detach/override rates = the component doesn't fit needs).

**Cons:**
- Measurement infrastructure is real work; cross-codebase static analysis needs per-stack tooling; numbers invite gaming.

**Implementation guidance:**
- Core dashboard: % screens/views using system components; count of hardcoded colors/spacing values per codebase (lint rule + trend); % consumers on current or previous major; Figma detach rate per component; top-20 component usage counts.
- Outcome studies: time-to-build a standard screen (pilot vs control team); accessibility audit findings per screen (system vs legacy screens); design-review rework rounds.
- Commonly cited practitioner targets (folklore, not independently verified benchmarks): >60% component coverage in active codebases by end of year 2; hardcoded-value count trending to near-zero in new code (enforced by lint); >80% of consumers within one major of current.
- Review quarterly with the council; feed unused/high-detach components into the deprecation or redesign queue.

**Related:** [#governance-and-contribution-models](#governance-and-contribution-models), [#costs-and-when-not-to-build-one](#costs-and-when-not-to-build-one), [#versioning-and-deprecation](#versioning-and-deprecation).

**Sources:** [Figma: Library analytics](https://help.figma.com/hc/en-us/articles/360039238353-View-and-explore-library-analytics), [Sparkbox Design Systems Survey](https://designsystemssurvey.seesparkbox.com/), Nathan Curtis, EightShapes: [Measuring Design System Success](https://medium.com/eightshapes-llc/and-you-thought-buttons-were-easy-26eb5b5c1871).

---

## Ecosystem and strategy

### Reference systems compared

**Definition:** The major public design systems, each embodying different priorities — useful as adoption candidates, benchmarks, and pattern references even when building your own.

**Reasoning / Evidence:** Each system encodes years of research and production hardening in a particular context; knowing what each optimizes for prevents both blind adoption and reinvention. Comparison:

| System | Owner | Optimized for | Adopt when | Watch out for |
|---|---|---|---|---|
| Material 3 | Google | Android-native feel; theming (dynamic color); huge component range; extensive guidance prose | Android apps; cross-platform consumer apps; teams wanting maximal free guidance | Strong opinionated aesthetic reads "Google"; web implementations lag the spec; M2→M3 churn |
| Fluent 2 | Microsoft | Windows + Microsoft 365 ecosystem coherence; desktop app conventions; React implementation | Windows/desktop apps; M365 integrations/add-ins | Web-outside-Microsoft feels branded; multiple historical Fluent versions confuse |
| Carbon | IBM | Data-dense enterprise products; strong data-table/dashboard patterns; per-component a11y docs; multi-framework (React/Angular/Vue/Web Components) | Enterprise/admin/analytics tools; IBM-adjacent stacks | Utilitarian aesthetic; heavier learning curve |
| Polaris | Shopify | Commerce admin UX; the best content/voice guidelines in the industry; opinionated merchant-task patterns | Shopify apps (mandatory-ish); commerce/admin tools; teams that need content standards to copy | Tightly coupled to Shopify's domain; less general-purpose |
| Atlassian DS | Atlassian | Collaboration/project-tool patterns; pragmatic tokens + governance docs | Jira/Confluence-ecosystem apps; team-tool products | Ecosystem-specific interaction idioms |
| GOV.UK Design System | UK Government | Accessibility + plain language above all; research-backed per-component findings; forms/services excellence | Government/civic services; form-heavy transactional services; any team wanting the a11y gold standard | Deliberately austere; not a branding vehicle; page-based (not SPA-first) |
| USWDS | US GSA | US federal compliance (Section 508); token architecture (well-documented tiers); templates for gov sites | US government and contractors (often required); token-system learners | US-gov visual identity; Sass-centric tooling |

**When to use:**
- As adoption base ([#adopt-vs-adapt-vs-build](#adopt-vs-adapt-vs-build)); as pattern reference when designing your own component (read 3–4 systems' take before deciding); as content-guideline source (Polaris, GOV.UK) even if you use none of their code.

**When NOT to use / exceptions:**
- Don't mix component libraries from multiple systems in one product (inconsistent behavior, doubled bundle, clashing tokens); reference many, ship one.
- A system's *guidance* can be adopted without its *code* and vice versa — they're separable layers.

**Pros:** Battle-tested, accessible, documented, free; instant baseline quality.
**Cons:** Someone else's brand and opinions; upstream churn on their schedule; deep customization can exceed build-your-own cost ([#adopt-vs-adapt-vs-build](#adopt-vs-adapt-vs-build)).

**Implementation guidance:**
- Benchmark ritual: before building any component, read its page in Material, Carbon, Polaris, and GOV.UK; note where they agree (near-standard — follow it) and disagree (context-dependent — decide for your context).
- Steal content guidelines from Polaris/GOV.UK, a11y documentation format from Carbon, token architecture from USWDS, theming from Material — provenance-tagged in your own docs.

**Related:** [#adopt-vs-adapt-vs-build](#adopt-vs-adapt-vs-build), [#accessibility-baked-into-components](#accessibility-baked-into-components), [09-content-ux-writing.md](09-content-ux-writing.md).

**Sources:** [Material 3](https://m3.material.io/), [Fluent 2](https://fluent2.microsoft.design/), [Carbon Design System](https://carbondesignsystem.com/), [Shopify Polaris](https://polaris.shopify.com/), [Atlassian Design System](https://atlassian.design/), [GOV.UK Design System](https://design-system.service.gov.uk/), [USWDS](https://designsystem.digital.gov/).

### Design-dev handoff

**Definition:** The transfer of design intent to implementation: specs, redlines, tokens, interaction notes, and the collaboration process around them. In system-mature teams, handoff shifts from "annotate every pixel" to "reference the system": designs point at existing components/tokens, and only deltas need specification.

**Reasoning / Evidence:** Handoff failure is a top source of UI defects and rework: unspecified states (error, loading, empty, overflow), unspecified breakpoints, and values eyeballed from mocks. Design systems attack the root cause — when a mock is assembled from library components bound to tokens, most of the spec *is* the reference, and tools (Figma Dev Mode, Zeplin) can surface component names, token names, and code snippets instead of raw px/hex. Remaining failures are process-shaped: designing only the happy path, and treating handoff as a throw-over-the-wall event instead of a conversation (pairing during build catches intent drift cheaply).

**When to use:**
- Every feature; formalize more for distributed/async teams and less-experienced implementers.

**When NOT to use / exceptions:**
- Over-specification of things the system already defines wastes time and invites divergence ("the mock says 15px" when the token says 16) — spec deltas, not the whole world.
- Prototypes/spikes can skip formal handoff; label them non-production explicitly.

**Pros:**
- Fewer defects and review rounds; faster builds; specs stay honest because they reference living components.

**Cons:**
- Requires designers to actually use the library (library hygiene, no detached copies); Dev-Mode-style tooling has per-seat costs.

**Implementation guidance:**
- Definition of done for a design: all states (default/hover/focus/disabled/loading/error/empty), responsive behavior at agreed breakpoints, content overflow behavior (long strings, [09-content-ux-writing.md](09-content-ux-writing.md) truncation), interaction notes (focus order, keyboard, motion), a11y annotations (labels, headings, alt) — checklist in the design review template.
- Name the system parts: annotate "Button / primary / md" and token names, not hex/px; treat any needed value *without* a token name as a signal to extend tokens or question the design.
- Kickoff + mid-build pairing beats end-stage inspection; budget designer time during implementation (~10–20% of design effort) for QA against intent.
- Design QA gate before release: designer reviews the built UI against the spec on real devices/themes; log deltas as bugs, not vibes.

**Related:** [#figma-code-sync-trends](#figma-code-sync-trends), [#documentation-practices](#documentation-practices), [17-ux-process-research.md](17-ux-process-research.md).

**Sources:** [Figma: Dev Mode](https://help.figma.com/hc/en-us/articles/15023124644247-Guide-to-Dev-Mode), [NN/g: DesignOps 101](https://www.nngroup.com/articles/design-operations-101/), [Zeplin: Design delivery](https://zeplin.io/).

### Figma-code sync trends

**Definition:** The tightening loop between design tools and production code: Figma Variables/modes mirroring token systems, Code Connect mapping Figma components to real code components (so Dev Mode shows your actual API snippets), token pipelines syncing both directions, Storybook-Figma embeds, and AI-assisted design-to-code generation.

**Reasoning / Evidence:** The historical root cause of design-code drift is dual sources of truth maintained by hand. Each capability removes a manual sync: Variables + modes let a single token JSON drive both Figma and CSS (via plugins/REST API); Code Connect (2024+) replaces autogenerated CSS-ish snippets with your real component code in Dev Mode, so developers copy correct API usage; visual regression + Storybook embeds keep docs honest. Direction of travel: code as the source of truth for components, design tools as a *view* over it — several teams now generate Figma libraries *from* coded components rather than vice versa. AI design-to-code currently produces plausible one-off markup but not system-consistent code without being constrained to your component APIs and tokens — treat as scaffolding, not delivery.

**When to use:**
- Token sync + Variables: any team with both a Figma library and a token package — highest ROI, lowest risk.
- Code Connect / snippet mapping: once the coded library is stable and adopted.

**When NOT to use / exceptions:**
- Don't auto-round-trip edits both directions without review gates — silent overwrites in either direction destroy trust; pick one canonical source per artifact type (commonly: tokens in repo, visuals in Figma, component truth in code).
- Don't ship AI-generated component code into the system without the same review bar as human contributions.

**Pros:**
- Drift reduced mechanically rather than by discipline; developers see correct usage at inspection time; theming testable in design and code from the same data.

**Cons:**
- Vendor lock-in around one design tool's API; plugin/pipeline maintenance; enterprise-tier feature gating (Variables modes, Code Connect require paid plans).

**Implementation guidance:**
- Phase 1: tokens as JSON in repo → CI syncs to Figma Variables (REST API/Tokens Studio); lint code against token usage.
- Phase 2: map library components to code with Code Connect (or Storybook links in component descriptions); enforce "no detached instances" in design review.
- Phase 3: visual regression (Chromatic/Percy) on the coded library as the drift alarm; audit Figma-vs-code state matrices quarterly.
- Keep the canonical-source table documented: who edits what, where, and how it propagates.

**Related:** [#token-tooling-and-standards](#token-tooling-and-standards), [#design-dev-handoff](#design-dev-handoff), [#documentation-practices](#documentation-practices).

**Sources:** [Figma: Guide to Variables](https://help.figma.com/hc/en-us/articles/15339657135383-Guide-to-variables), [Figma: Code Connect](https://www.figma.com/code-connect-docs/), [Style Dictionary](https://styledictionary.com/), [Storybook: Figma integration](https://storybook.js.org/docs/sharing/design-integrations).

### Adopt vs adapt vs build

**Definition:** The strategic choice among three approaches (NN/g's framing): **adopt** an existing open system wholesale, **adapt** one (theme/extend an existing system or headless base under your brand), or **create** a proprietary system from scratch. Control, cost, and brand differentiation all rise together along that spectrum.

**Reasoning / Evidence:** NN/g: adopting is lowest cost/fastest; custom is justified only when needs can't be met by existing systems; and — the key trap — as customizations accumulate, "the cost savings you may have gained from using the existing design system will diminish, and in the long run you may be better off creating your own anyway." Heavily overriding a themable-but-opinionated library (fighting its DOM, specificity, and update churn) is routinely more expensive than owning components built on headless primitives. The middle path has matured: headless behavior layers (Radix, React Aria) + your tokens/styling gives custom brand with borrowed correctness (a11y, focus, keyboard) — most of "build" at a fraction of the cost.

**When to use (decision framework):**
- **Adopt** when: internal tools/MVPs; small team, no design-system headcount; brand differentiation irrelevant; platform-ecosystem apps where matching the host wins (Shopify app → Polaris; Android → Material; US gov → USWDS).
- **Adapt** when: brand matters but budget is moderate; an existing system covers ≥80% of component needs; you can absorb upstream churn. Sub-decision: theme an opinionated system (small deltas) vs style a headless base (large visual deltas).
- **Build** when: multi-product org with dedicated funded team; distinctive interaction patterns existing systems don't cover (novel domains, canvas/editor UIs); design as competitive moat; strict long-horizon control needs. Even then: build *on* headless primitives and the DTCG token stack, not from bare DOM.

**When NOT to use / exceptions:**
- Don't build to satisfy team preference for greenfield work — the maintenance tail is the real purchase.
- Don't adopt a system whose interaction idioms fight your platform (a desktop-dense trading app on a touch-first library).
- Reassess at inflection points (new platform, acquisition, rebrand): the right answer changes as the org changes.

**Pros / Cons:**
- Adopt — Pros: days-to-value, free maintenance/a11y, hiring familiarity. Cons: generic look, upstream churn, customization ceiling.
- Adapt — Pros: brand + borrowed correctness, moderate cost. Cons: override fragility, two-master problem (upstream + local layers).
- Build — Pros: total fit and control, differentiation. Cons: standing cost of roughly 3–8 FTE (practitioner estimate, unverified), years to maturity, you own every a11y bug.

**Implementation guidance:**
- Score the decision on: team size, product horizon (years of replicable screens — NN/g's ROI test), brand-differentiation need, component-coverage fit (audit your UI inventory vs the candidate's catalog: adopt if ≥80% covered, adapt 60–80%, consider build <60%), and appetite for standing maintenance.
- Whatever the path, own your *token layer* and *content guidelines* from day one — they're portable across all three strategies and make later migration survivable.
- If adapting, wrap third-party components behind your own package/namespace so a future swap is a re-implementation behind a stable API, not a product-wide rewrite.
- Revisit the decision explicitly every 12–18 months against adoption and override metrics ([#adoption-metrics](#adoption-metrics)).

**Related:** [#costs-and-when-not-to-build-one](#costs-and-when-not-to-build-one), [#reference-systems-compared](#reference-systems-compared), [#composition-vs-configuration](#composition-vs-configuration).

**Sources:** [NN/g: Design Systems 101 — How to Approach Design-System Adoption](https://www.nngroup.com/articles/design-systems-101/), [Radix Primitives](https://www.radix-ui.com/primitives), [React Aria](https://react-spectrum.adobe.com/react-aria/), [USWDS: How to use USWDS](https://designsystem.digital.gov/how-to-use-uswds/).

---

## Quick reference

*Note: numbers marked with † below (build-speed gains, FTE counts, docs-cost share, coverage targets) are practitioner folklore / vendor-published estimates without independent verification — useful anchors, not benchmarks.*

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Build rationale | System = product for products; justified by years of replicable screens across teams | Min team: 1 IxD + 1 visual + 1 dev; mature: 3–8 FTE† |
| When not to build | Skip for MVPs, single small teams, one-offs; adopt instead | <5 UI builders → adopt; <30% adoption at 18–24mo → restructure |
| Token tiers | Primitive (meaning-free scales) → semantic (usage names) → component; consumers use semantic only | 8–12 steps/hue; 4/8px spacing base |
| Token naming | Fixed grammar: category–concept–property–variant–state–scale; controlled vocab per slot | Define before ~30 tokens exist |
| Theming via tokens | Themes = value re-mappings of semantic names; components never change | Dark: no pure black (#121212+), desaturate 1–2 steps, ≥4.5:1 per theme |
| Token tooling | JSON source (DTCG) → Style Dictionary → CSS/Swift/Kotlin + Figma Variables sync | 1 source, N outputs; CI contrast + coverage checks |
| Atomic design | Keep compositional thinking, drop chemistry taxonomy; organize by function | 2 tiers: components + patterns; promote after 3 uses |
| Component APIs | Enum variants over booleans; accessible defaults; shared prop vocabulary; clamp/guard numeric props at runtime — degrade, don't crash | ≤2 required props; sizes sm/md/lg; percentages clamped 0–100 |
| Composition vs config | Config for leaf components, compound composition for containers, recipes for consistency | Fork/copy events = composition signal |
| Variants & states | Full variant×state matrix required per component; snapshot-test the grid | 6 core states; 3 emphasis levels + danger; targets ≥24px (44–48px touch) |
| A11y in components | APG-compliant behavior baked in; native HTML first; tokens guarantee contrast | axe catches ~30–40%; 4.5:1 text, 3:1 UI/focus |
| Documentation | Per-component: when-to-use, do/don't pairs, states, a11y, generated API docs — one page for design+code | Docs ≈ 30–50% of component cost†; ≥2 do/don't pairs |
| Test hooks | Declared, lint-enforced `data-testid` hooks on interactive elements; tests use declared names, role/label queries first | kebab-case `feature-element[-modifier]`; prod stripping optional |
| Theming axes | Dark/high-contrast, multi-brand, density — all as token re-maps, never component forks | Rows 52/40/32px; no `if (brand)` in components |
| Versioning | Semver + deprecation pipeline + codemods; publish what counts as breaking | Remove after ≥2 minors or 6mo; support current+previous major |
| Governance | Hybrid model: core owns foundations, teams contribute through a quality-gated ladder | Council biweekly; rule of three for promotion |
| Adoption metrics | Measure usage (coverage, detach rate, version currency), not component count | >60% coverage by yr 2†; >80% within 1 major |
| Reference systems | Material=theming/Android; Carbon=enterprise data; Polaris=content; GOV.UK/USWDS=a11y+forms; Fluent=Windows | Read 3–4 systems before designing any component |
| Handoff | Spec deltas against the system; all states + overflow + a11y in definition of done | Designer QA ~10–20% of design effort |
| Figma-code sync | One canonical source per artifact: tokens in repo, component truth in code, visuals in Figma | Sync via CI, not manually |
| Adopt vs build | Coverage audit: ≥80% adopt, 60–80% adapt, <60% consider build (on headless primitives) | Own tokens + content guidelines regardless; review every 12–18mo |
