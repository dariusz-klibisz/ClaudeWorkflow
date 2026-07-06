# Visual Design

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers the visual layer: hierarchy, typography, color, spacing/layout, density, elevation, shape, iconography, imagery, dark mode, and trend analysis (flat, neumorphism, glassmorphism, brutalism). Perceptual grouping laws behind these rules are in [01-core-principles.md](01-core-principles.md#gestalt-principles); contrast and non-visual requirements in [05-accessibility.md](05-accessibility.md); motion in [04-interaction-design.md](04-interaction-design.md#motion-and-animation); tokenizing these decisions in [10-design-systems.md](10-design-systems.md).

## Contents

- [Hierarchy](#hierarchy)
  - [Visual hierarchy](#visual-hierarchy)
- [Typography](#typography)
  - [Type scale](#type-scale)
  - [Line height and line length](#line-height-and-line-length)
  - [Font selection and pairing](#font-selection-and-pairing)
  - [Font weight and emphasis](#font-weight-and-emphasis)
- [Color](#color)
  - [Building a UI palette](#building-a-ui-palette)
  - [Semantic color](#semantic-color)
  - [Contrast](#contrast)
  - [Color blindness and color independence](#color-blindness-and-color-independence)
  - [Modern color spaces (OKLCH)](#modern-color-spaces-oklch)
- [Spacing and layout](#spacing-and-layout)
  - [The 8pt spacing grid](#the-8pt-spacing-grid)
  - [Whitespace](#whitespace)
  - [Alignment and grids](#alignment-and-grids)
- [Density](#density)
  - [Density modes](#density-modes)
- [Surfaces and shape](#surfaces-and-shape)
  - [Elevation, shadows, depth](#elevation-shadows-depth)
  - [Border radius and shape language](#border-radius-and-shape-language)
- [Iconography and imagery](#iconography-and-imagery)
  - [Iconography](#iconography)
  - [Imagery and illustration](#imagery-and-illustration)
- [Dark mode](#dark-mode)
  - [Dark mode design](#dark-mode-design)
- [Trends and styles](#trends-and-styles)
  - [Flat design and minimalism](#flat-design-and-minimalism)
  - [Neumorphism](#neumorphism)
  - [Glassmorphism](#glassmorphism)
  - [Brutalism](#brutalism)
- [Quick reference](#quick-reference)

---

## Hierarchy

### Visual hierarchy

**Definition:** The deliberate ranking of elements by visual prominence — via size, weight, color, contrast, position, and whitespace — so users perceive importance and reading order without effort.

**Reasoning / Evidence:** Pre-attentive processing: size/color/contrast differences register in <200ms, before conscious reading. Hierarchy is how a screen answers "what is this, what matters, what do I do" in the first glance (users form impressions in ~50ms). Eyetracking shows users follow prominence gradients, not document order.

**When to use:**
- Every screen: define 3 levels max in most views — primary (the point of the screen), secondary (supporting), tertiary (metadata/chrome).
- Landing/marketing: one dominant message + one dominant action.
- Data views: the decision-driving value visually largest; context smaller.

**When NOT to use / exceptions:**
- Reference/table content where all cells are peers — impose hierarchy on structure (headers, grouping), not on data values.
- Don't fake hierarchy that misrepresents importance to the user (business-preferred option shouting over user-relevant info — [18-patterns-antipatterns.md](18-patterns-antipatterns.md#visual-interference-misdirection)).

**Pros:**
- Faster comprehension, correct first clicks, lower perceived complexity.

**Cons:**
- Requires real prioritization decisions; committee design ("make everything prominent") flattens it back to noise.

**Implementation guidance:**
- Squint test / 5-second test: blur the screen; the primary element should still dominate.
- Combine at least two cues for each hierarchy jump (size + weight, color + position) — single-cue differences are fragile.
- Type hierarchy: each level ≥1.2× size step or a weight step from the next; body 400, emphasis 500–600, headings 600–700.
- One primary button per view; secondary as outline/ghost; destructive visually distinct and separated.
- Position: top-left region (LTR) carries the entry point ([02-cognitive-foundations.md](02-cognitive-foundations.md#scanning-patterns)).

**Verification checklist:**
- Each screen has one clear primary purpose, and the primary action is visually prominent without overwhelming the content.
- Headings expose the structure of the page; a heading-only skim conveys the outline.
- Related items are aligned and grouped; unrelated groups are visibly separated.
- No cluster of elements is equally dominant — the squint test yields a single winner.

**Related:** [01-core-principles.md](01-core-principles.md#von-restorff-effect), [01-core-principles.md](01-core-principles.md#8-aesthetic-and-minimalist-design), [#type-scale](#type-scale).

**Sources:** [NN/g: Visual Hierarchy](https://www.nngroup.com/articles/visual-hierarchy-ux-definition/), [NN/g: Principles of Visual Design](https://www.nngroup.com/articles/principles-visual-design/).

---

## Typography

### Type scale

**Definition:** A constrained, ratio-based set of font sizes (e.g., 12, 14, 16, 20, 25, 31, 39) derived from a base size and a modular ratio (1.2 minor third, 1.25 major third, 1.333 perfect fourth), replacing ad-hoc size choices. Note: in the 1.25 example the small steps (12→14→16) are not true 1.25 steps — small sizes are conventionally rounded and hand-tuned for legibility rather than strictly ratio-derived; the ratio governs the scale from the base upward.

**Reasoning / Evidence:** Ratio-derived scales produce harmonious, clearly distinguishable steps (from centuries of print typography). Constraint enforces consistency and makes hierarchy levels recognizable across screens.

**When to use:**
- Every project: define the scale as tokens before designing screens ([10-design-systems.md](10-design-systems.md#design-tokens)).
- Content-heavy sites: larger ratios (1.25–1.414) for dramatic hierarchy.
- Dense applications: smaller ratios (1.125–1.2) — big jumps waste space in UI chrome.

**When NOT to use / exceptions:**
- Slavish ratio purity is unnecessary — round to whole/even pixels; adjust individual steps optically.
- Marketing display sizes can break scale for impact.

**Pros:** Consistency for free; fewer decisions; recognizable hierarchy levels.
**Cons:** Pure ratios yield awkward fractions (must round); one scale rarely fits both marketing and app density (define two, related).

**Implementation guidance:**
- Base: 16px body for web (never smaller than 16px for primary reading; 14px acceptable for dense app UI, 12px floor for captions/labels).
- App scale example (1.2): 12 / 14 / 16 / 20 / 24 / 29 / 35. Limit to 5–7 sizes total.
- Use `rem` units on web so user font-size preferences scale everything ([13-platform-web.md](13-platform-web.md#web-typography-specifics)).
- Pair each size with a set line height token (see next entry).

**Verification checklist:**
- All text is real text, not images of text.
- Layout survives user text resizing (browser zoom, OS dynamic type) without loss of content or function.
- No meaning is conveyed solely by capitalization or color of text.

**Related:** [#visual-hierarchy](#visual-hierarchy), [10-design-systems.md](10-design-systems.md#design-tokens), [05-accessibility.md](05-accessibility.md#zoom-and-text-resize).

**Sources:** [Material Design: Type scale](https://m3.material.io/styles/typography/type-scale-tokens), [type-scale.com](https://typescale.com/).

### Line height and line length

**Definition:** Line height (leading): vertical space per text line, expressed as a multiplier. Line length (measure): characters per line. Both directly govern reading comfort and speed.

**Reasoning / Evidence:** Legibility research (print + screen): body text reads best at 45–75 characters/line (66 the classic ideal); longer lines cause return-sweep errors (eye loses the next line), shorter cause choppy saccades. Line height 1.4–1.6 for body; WCAG 1.4.8 (AAA) specifies ≥1.5 line spacing and ≤80 characters.

**When to use:**
- All body/paragraph text: enforce measure via container widths (`max-width: 65ch`).
- Long-form reading: top of the ranges (1.6 leading, ~66ch).
- UI labels/single lines: tighter (1.2–1.3) is fine; headings 1.1–1.3.

**When NOT to use / exceptions:**
- Data tables and code: line-length rules don't apply (structure ≠ prose); code uses ~1.4–1.5 leading with monospace.
- Very large display text needs negative letter-spacing and tight leading (1.0–1.15) for cohesion.

**Pros:** Cheap, high-impact readability wins; concrete numbers to lint against.
**Cons:** Fluid layouts fight fixed measures — solve with max-width, not fixed width.

**Implementation guidance:**
- Body: `line-height: 1.5` default, `max-width: 60–75ch`.
- Never justify text on the web without hyphenation (rivers); left-align LTR ragged-right.
- Paragraph spacing: 0.75–1× font size between paragraphs; no first-line indents in UI contexts.
- Respect WCAG 1.4.12 Text Spacing: layout must survive user overrides (line-height 1.5×, paragraph 2×, letter 0.12×, word 0.16× font size) without loss ([05-accessibility.md](05-accessibility.md#text-spacing)).

**Related:** [#type-scale](#type-scale), [05-accessibility.md](05-accessibility.md#text-spacing), [09-content-ux-writing.md](09-content-ux-writing.md#truncation-and-line-length).

**Sources:** [WCAG 1.4.8 Visual Presentation](https://www.w3.org/WAI/WCAG22/Understanding/visual-presentation.html), [Butterick's Practical Typography: line length](https://practicaltypography.com/line-length.html).

### Font selection and pairing

**Definition:** Choosing typefaces for UI text, reading text, and display, and combining them (typically ≤2 families) coherently.

**Reasoning / Evidence:** Screen UI defaults to sans-serifs for small-size clarity (simpler letterforms survive low sizes/hinting); serif readability disadvantages on modern high-DPI screens are largely gone for body sizes ≥16px, so the choice is now tonal. System font stacks render instantly (no FOUT), match platform feel, and cost zero bytes.

**When to use:**
- App UIs: system stack (`-apple-system, "Segoe UI", Roboto, ...`) or one webfont family with broad weights; consider variable fonts (one file, all weights).
- Brand-heavy surfaces: display font for headings + workhorse for body/UI.
- Multilingual products: verify glyph coverage (CJK, Cyrillic, Arabic) before committing.

**When NOT to use / exceptions:**
- >2 families is rarely justified; use weight/size/case for variety instead.
- Decorative fonts never for body or UI controls; thin weights (<300) fail on low-DPI and dark backgrounds.
- Don't use fonts lacking distinguishable Il1 0O for anything showing codes/IDs (use tabular, distinguishable faces; enable `font-variant-numeric: tabular-nums` for aligned numbers in tables).

**Pros (system fonts):** Zero load time, platform-native feel, automatic OS optimization.
**Cons (system fonts):** No brand distinctiveness; metrics differ across platforms (layout shifts in testing).

**Implementation guidance:**
- Pairing rule: pair by contrast (serif display + sans body, or one superfamily), match x-heights, never pair near-identical faces.
- Load ≤2 families × ≤4 weights; `font-display: swap` with metric-matched fallbacks to kill CLS ([13-platform-web.md](13-platform-web.md#web-typography-specifics)).
- Tabular numerals for any column of numbers; proportional for prose.

**Related:** [#type-scale](#type-scale), [13-platform-web.md](13-platform-web.md#web-typography-specifics), [10-design-systems.md](10-design-systems.md).

**Sources:** [Butterick's Practical Typography](https://practicaltypography.com/), [web.dev: font best practices](https://web.dev/articles/font-best-practices).

### Font weight and emphasis

**Definition:** Using weight (400/500/600/700), style (italic), case, and color — not size alone — to create emphasis and structure within text.

**Reasoning / Evidence:** Weight differences are pre-attentively detected and cheaper spatially than size jumps; scanning studies show bold keywords aid the "spotted" scanning pattern. Overuse destroys the signal ([Von Restorff](01-core-principles.md#von-restorff-effect)).

**When to use:**
- Keyword emphasis in scannable text (1–3 words, sparingly).
- Label vs value pairs (label 500 muted, value 400 strong — or inverse, consistently).
- Hierarchy within same-size text (list item title 600, meta 400).

**When NOT to use / exceptions:**
- ALL CAPS for body/sentences (slower reading, shouting tone) — acceptable for short labels/overlines with letter-spacing +0.05–0.1em.
- Italic on screens at small sizes renders poorly; prefer weight or color.
- Underline is reserved for links on the web — never underline non-links ([01-core-principles.md](01-core-principles.md#similarity)).
- Color-only emphasis fails color-blind users; pair with weight ([05-accessibility.md](05-accessibility.md#color-independence)).

**Pros:** Space-free hierarchy; supports scanning.
**Cons:** Each additional emphasized element halves the value of all of them.

**Implementation guidance:**
- Jump ≥2 weight steps for visible difference (400→600, not 400→500 in many faces).
- Muted text: reduce contrast, but keep ≥4.5:1 (typical: 60–70% opacity of text color on the same background).
- Budget: ≤10% of visible words emphasized.

**Related:** [#visual-hierarchy](#visual-hierarchy), [09-content-ux-writing.md](09-content-ux-writing.md#scannability-headings-bullets-bold).

**Sources:** [NN/g: Text Scanning Patterns](https://www.nngroup.com/articles/text-scanning-patterns-eyetracking/).

---

## Color

### Building a UI palette

**Definition:** A structured color system: neutral ramp (greys for text/surfaces/borders, ~8–12 steps), one primary/brand hue ramp, 1–2 secondary/accent ramps, and semantic ramps (success/warning/error/info) — each ramp spanning light→dark for state and theming needs.

**Reasoning / Evidence:** Ramps (not single colors) are required because every hue needs multiple lightness steps for hover/active states, backgrounds vs text usage, and dark-mode inversion. The 60-30-10 rule (60% neutral, 30% secondary/surfaces, 10% accent) keeps accent colors meaningful — accents work by scarcity ([Von Restorff](01-core-principles.md#von-restorff-effect)).

**When to use:**
- Project start: define ramps as primitive tokens; map to semantic tokens (text-primary, surface-raised, border-subtle) before any screen design ([10-design-systems.md](10-design-systems.md#design-tokens)).
- Most UI should be neutral: color carries meaning (interaction, status, brand moments), not decoration.

**When NOT to use / exceptions:**
- Data visualization needs its own categorical/sequential/diverging palettes (12+ distinguishable categorical colors is the practical max; use shape/pattern beyond).
- Marketing surfaces legitimately run higher accent ratios than app UI.

**Pros:** Consistency, theming-readiness, restrained accent = clear interactive signaling.
**Cons:** Up-front system work; retrofitting ramps onto ad-hoc colors is painful.

**Implementation guidance:**
- Neutral ramp: 8–12 steps from near-white to near-black; avoid pure black text on pure white (harsh halation) — e.g., #1a1a1a on #ffffff or off-white.
- Generate ramps in a perceptual space (OKLCH — see below) so steps are visually even and contrast is predictable.
- Reserve the accent hue exclusively for interactive elements; never decorative use of the link color.
- Check every text/background token pair against WCAG at definition time, not per-screen ([#contrast](#contrast)).
- Avoid oversaturated palettes for dense work surfaces — sustained-use tools need muted, low-chroma surfaces.

**Related:** [#semantic-color](#semantic-color), [#modern-color-spaces-oklch](#modern-color-spaces-oklch), [10-design-systems.md](10-design-systems.md#design-tokens).

**Sources:** [Carbon Design System: Color](https://carbondesignsystem.com/elements/color/usage/), [Refactoring UI: color](https://www.refactoringui.com/).

### Semantic color

**Definition:** Colors carrying conventional meaning: green=success/positive, red=error/destructive/negative, yellow/amber=warning, blue=information/interactive. Applied via semantic tokens so meaning, not raw hue, drives usage.

**Reasoning / Evidence:** These mappings are strong learned conventions in Western/global software (traffic-light metaphor); violating them (red=good) causes real errors. But they are culture-variable (red=prosperity in China; finance apps there use red=up) and invisible to some color-blind users — hence redundancy requirements.

**When to use:**
- Status communication: form validation, alerts, badges, system state.
- Destructive actions: red styling on the button *and* separation/labeling.
- Financial/quantitative deltas: green/red with sign and icon redundancy.

**When NOT to use / exceptions:**
- Never color alone: pair with icon + text (WCAG 1.4.1) ([05-accessibility.md](05-accessibility.md#color-independence)).
- Localize semantics where stakes are high (China-market finance: configurable/inverted delta colors).
- Don't dilute: red only for errors/destruction — using red for emphasis trains habituation ([02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)).

**Pros:** Instant, learned recognition; cross-product convention transfer.
**Cons:** Color-blind invisibility without redundancy; cultural variance; semantic hues must also hit contrast targets (yellow on white is a chronic failure — darken to amber for text/icons).

**Implementation guidance:**
- Four semantic ramps (success/warning/error/info), each with: strong (text/icon on light), subtle (background tint), border variants.
- Warning text/icons: use amber ≥ #b45309-region on white, not pure yellow.
- Every semantic signal = color + icon + label ("Error:" prefix or equivalent).

**Related:** [#color-blindness-and-color-independence](#color-blindness-and-color-independence), [09-content-ux-writing.md](09-content-ux-writing.md#error-message-writing).

**Sources:** [Polaris: Colors](https://polaris.shopify.com/design/colors), [WCAG 1.4.1 Use of Color](https://www.w3.org/WAI/WCAG22/Understanding/use-of-color.html).

### Contrast

**Definition:** Luminance difference between foreground and background, expressed as a ratio 1:1–21:1. WCAG minimums: 4.5:1 normal text, 3:1 large text (≥24px or ≥18.7px bold), 3:1 UI components and meaningful graphics (AA); 7:1 / 4.5:1 for AAA.

**Reasoning / Evidence:** WCAG 2.x contrast math models low-vision legibility; ~1 in 12 men has some color vision deficiency and low-vision prevalence rises steeply with age. Contrast failures are the single most common accessibility defect in audits (WebAIM Million: ~80% of home pages).

**When to use:**
- Every text/background pair, icon, input border, focus indicator, chart element carrying meaning.
- Check at token-definition time and in both themes; automate in CI.

**When NOT to use / exceptions:**
- Exempt: disabled controls, decorative elements, logos (per WCAG) — but visibly-too-faint disabled states still harm; aim ≥3:1 even there when feasible.
- Placeholder text is NOT exempt in practice — chronic failure; style ≥4.5:1 or don't rely on placeholders ([07-forms-and-input.md](07-forms-and-input.md#labels-vs-placeholders)).
- WCAG ratio math penalizes light-on-dark unevenly; APCA (candidate for WCAG 3) models perception better — useful as secondary check, not yet the legal standard.

**Pros:** Objective, automatable, legally load-bearing.
**Cons:** Passing ratios can still look bad (vibrating saturated pairs); ratio ignores font weight/size nuances beyond the large-text cutoff (APCA addresses this).

**Implementation guidance:**
- Targets: body text ≥4.5:1 (aim 7:1), large headings ≥3:1, borders/icons/focus rings ≥3:1 against adjacent colors.
- Grey text floor on white: #767676 (4.54:1). On dark #121212: light grey ≥ #a0a0a0-region — verify per pair.
- Text over images: scrim/gradient, measure worst point.
- Tools: axe/Lighthouse in CI; contrast tokens documented in the design system.

**Verification checklist:**
- All text meets WCAG minimums (4.5:1 normal, 3:1 large) against its actual rendered background.
- Non-text indicators (input borders, focus rings, icons, chart marks) meet 3:1 against adjacent colors.
- Dark mode and OS high-contrast/forced-colors modes have been tested — every token pair re-checked, not assumed.
- Text over imagery is measured at its worst-case point (behind the lightest/darkest region).

**Related:** [05-accessibility.md](05-accessibility.md#contrast-minimum), [#dark-mode-design](#dark-mode-design).

**Sources:** [WCAG 1.4.3](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html), [WebAIM Million](https://webaim.org/projects/million/), [APCA](https://git.apcacontrast.com/).

### Color blindness and color independence

**Definition:** Designing so information survives color vision deficiency (CVD): deuteranomaly/deuteranopia (green-weak/blind, most common), protanopia (red), tritanopia (blue, rare). ~8% of men, ~0.5% of women have some CVD.

**Reasoning / Evidence:** WCAG 1.4.1: color must not be the only visual means of conveying information. Red/green distinctions — the most common semantic pairing — are precisely what the most common CVD collapses.

**When to use:**
- All status/semantic signaling; charts (lines/series); diffs; heatmaps; form validation.
- Map/geo visualization palettes.

**When NOT to use / exceptions:**
- None — this is a hard requirement, not a trade-off.

**Pros:** Redundant encoding (icon+text+position) helps *everyone* in bad viewing conditions (sunlight, projectors, night mode).
**Cons:** Adds design constraints to viz palettes; some dense visualizations need rework (patterns/labels).

**Implementation guidance:**
- Always pair color with a second channel: icon shape, text label, position, pattern.
- CVD-safe categorical palettes: Okabe-Ito (8 colors), viridis family for sequential.
- Red/green deltas: add +/− signs and arrows.
- Test with simulators (Chrome DevTools vision deficiency emulation, Sim Daltonism) as routine QA.

**Verification checklist:**
- Every color-coded state is also distinguishable by text, icon, pattern, or shape with color removed (greyscale test).
- Charts and diffs remain readable under deuteranopia/protanopia simulation.

**Related:** [#semantic-color](#semantic-color), [05-accessibility.md](05-accessibility.md#color-independence).

**Sources:** [WCAG 1.4.1](https://www.w3.org/WAI/WCAG22/Understanding/use-of-color.html), [Okabe & Ito palette](https://jfly.uni-koeln.de/color/).

### Modern color spaces (OKLCH)

**Definition:** Perceptually uniform color spaces (OKLCH/OKLab) where equal numeric steps look equally different, unlike HSL where lightness lies (HSL yellow and blue at "50%" differ wildly in perceived brightness).

**Reasoning / Evidence:** Björn Ottosson's OKLab (2020); CSS Color 4 ships `oklch()` with wide-gamut (P3) support. Ramps generated in OKLCH keep consistent perceived lightness/chroma across hues — predictable contrast and even steps, which HSL cannot deliver.

**When to use:**
- Generating palette ramps programmatically; theming systems computing hover/active shades; dark-mode palette derivation.
- Wide-gamut (P3) brand colors on modern displays.

**When NOT to use / exceptions:**
- Legacy browser targets need sRGB fallbacks (trivial with CSS `@supports` or build-time conversion).
- Hand-picked small palettes don't require it — value is in systematic generation.

**Pros:** Even ramps, predictable contrast, hue-stable lightness manipulation, P3 access.
**Cons:** Team familiarity; tooling still catching up; numbers less guessable than hex.

**Implementation guidance:**
- Store primitives as OKLCH; export hex/sRGB fallbacks at build.
- Ramp recipe: fix hue & chroma, step lightness evenly (e.g., L 0.97→0.15 in ~10 steps); reduce chroma at extremes to stay in gamut.
- Derive dark theme by remapping the lightness axis, not inverting hex values.

**Related:** [#building-a-ui-palette](#building-a-ui-palette), [#dark-mode-design](#dark-mode-design), [10-design-systems.md](10-design-systems.md#design-tokens).

**Sources:** [Ottosson: OKLab](https://bottosson.github.io/posts/oklab/), [MDN: oklch()](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value/oklch).

---

## Spacing and layout

### The 8pt spacing grid

**Definition:** All spacing (padding, margins, gaps) and most component dimensions use multiples of 8 (with 4 for fine adjustments): 4, 8, 12, 16, 24, 32, 48, 64...

**Reasoning / Evidence:** Constraint system: removes per-decision bikeshedding, guarantees vertical rhythm, divides cleanly across 1×/2×/3× DPI screens. Adopted by Material, Fluent, most mature systems.

**When to use:**
- Universally, as spacing tokens; component internals (button padding 8×16, input height 40/48) and layout gaps alike.

**When NOT to use / exceptions:**
- Optical adjustments override math (icon visually centered ≠ mathematically centered; text baseline tweaks of 1–2px are fine).
- Typography line-heights may land off-grid; prioritize readable leading over grid purity.

**Pros:** Consistency, speed, DPI-clean rendering, easy handoff ("space-4" not "23px").
**Cons:** None significant; the constraint occasionally forces 8 where 6 looked better — accept it.

**Implementation guidance:**
- Token scale: 2 (hairline exception), 4, 8, 12, 16, 24, 32, 48, 64, 96.
- Spacing communicates grouping: within-group ≤8–12, between-group ≥16–24, between-section ≥32–48 ([Proximity](01-core-principles.md#proximity)).
- Touch target minimum 44–48px tall regardless of visual size (padding extends hit area).

**Verification checklist:**
- Labels sit visibly nearer to their own controls than to neighboring controls.
- Field sets and sections are grouped under headings, with larger gaps between groups than within them.
- Action buttons sit adjacent to the content they affect, not stranded elsewhere on the screen.

**Related:** [01-core-principles.md](01-core-principles.md#proximity), [#alignment-and-grids](#alignment-and-grids), [10-design-systems.md](10-design-systems.md#design-tokens).

**Sources:** [Material Design: Layout spacing](https://m3.material.io/foundations/layout/understanding-layout/spacing).

### Whitespace

**Definition:** Empty space around and between elements — an active design element that groups, separates, emphasizes, and paces content; not "wasted" space.

**Reasoning / Evidence:** Proximity grouping ([Gestalt](01-core-principles.md#proximity)) works through whitespace; studies link generous spacing to improved comprehension and perceived quality/luxury. Crowding degrades reading and target acquisition (Fitts penalty for adjacent targets).

**When to use:**
- To group without boxes; to isolate the primary CTA; to pace long pages into digestible sections; around headings (more above than below binds heading to its content).

**When NOT to use / exceptions:**
- Dense professional tools: whitespace competes with information density — compress deliberately there ([#density-modes](#density-modes)).
- Excessive spacing on long forms/pages inflates scrolling; balance against interaction cost.

**Pros:** Free clarity and perceived quality; the cheapest grouping tool.
**Cons:** Screen real estate cost; stakeholder pressure to fill it is constant.

**Implementation guidance:**
- Heading spacing ratio: above ≥1.5–2× the space below.
- Padding scales with container: cards 16–24px, modals 24–32px, page sections 48–96px vertical.
- Empty is a statement: an isolated button in whitespace is the strongest CTA treatment.
- Prefer whitespace over borders/boxes for grouping where it suffices — reach for dividers only when spacing alone can't carry the structure.

**Related:** [01-core-principles.md](01-core-principles.md#proximity), [#the-8pt-spacing-grid](#the-8pt-spacing-grid).

**Sources:** [NN/g: The Power of White Space](https://www.nngroup.com/articles/whitespace/).

### Alignment and grids

**Definition:** Layout structured on a column grid (typically 12 columns with fixed gutters) and consistent alignment edges, so elements share invisible lines.

**Reasoning / Evidence:** Continuity/alignment grouping ([Gestalt](01-core-principles.md#continuity)); misalignment registers pre-attentively as disorder. 12-column grids dominate because 12 divides by 2/3/4/6 — flexible splits without fractional columns.

**When to use:**
- Page-level layout: 12-col desktop (gutters 16–24px), 8-col tablet, 4-col mobile.
- Within components: shared left edges for labels, values, controls.
- Cross-screen consistency: same grid = same perceived family.

**When NOT to use / exceptions:**
- Canvas/creative tools and data-viz interiors don't grid.
- CSS Grid/Flexbox make content-driven layout easy — the 12-col mental model matters more than literal framework columns.
- Breaking the grid deliberately draws attention (hero moments) — once per page at most.

**Pros:** Order, rhythm, responsive rationality, design-dev shared vocabulary.
**Cons:** Grid-first thinking can produce boxy sameness; content should drive exceptions.

**Implementation guidance:**
- One primary left alignment edge per region; every extra alignment axis costs scanning effort.
- Center-align only short content blocks (heroes, empty states); never center-align form labels or long text.
- Max content width: 1200–1440px typical; text columns capped at 60–75ch regardless ([#line-height-and-line-length](#line-height-and-line-length)).

**Related:** [01-core-principles.md](01-core-principles.md#continuity), [13-platform-web.md](13-platform-web.md#responsive-design).

**Sources:** [Material Design: Layout](https://m3.material.io/foundations/layout/understanding-layout/overview).

---

## Density

### Density modes

**Definition:** The information-per-pixel dial: comfortable (consumer default), compact (pro/data-heavy), spacious (touch-first/accessibility). Mature systems ship density as a token-level mode.

**Reasoning / Evidence:** Optimal density is audience- and task-dependent: NN/g's complex-application research shows expert users of data tools *prefer and perform better with* higher density (fewer scrolls/clicks, more context per glance), while consumer audiences drop off with crowding. One density cannot serve both.

**When to use:**
- Data tables, dashboards, admin consoles: compact default or user-switchable density (row heights 32/40/48).
- Consumer/marketing: comfortable/spacious.
- Touch contexts: enforce ≥44–48px targets regardless of visual density.
- High density helps expert comparison and monitoring; low density helps reading, onboarding, touch, and low-stress decisions — choose per task and audience.

**When NOT to use / exceptions:**
- Density switching adds testing surface; small products can pick one density and skip the toggle.
- Never densify by shrinking text below 12px or targets below WCAG 24px minimum ([05-accessibility.md](05-accessibility.md#target-size)).

**Pros:** Serves experts without alienating novices; user control over a real preference axis.
**Cons:** Multiplies component states/QA; density-dependent bugs (truncation, wrapping).

**Implementation guidance:**
- Implement via spacing/size tokens with density multipliers (e.g., compact = 0.75× paddings, row 32px; comfortable = 40px; spacious = 48px).
- Persist per-user, per-view (a user may want compact tables but comfortable settings).
- Keep type size constant across densities where possible; compress padding, not text.

**Related:** [#whitespace](#whitespace), [08-navigation-ia.md](12-data-tables-dashboards.md#sorting-filtering-and-column-management), [10-design-systems.md](10-design-systems.md#theming-dark-mode-multi-brand-density).

**Sources:** [NN/g: Complex Application Design](https://www.nngroup.com/articles/complex-application-design/).

---

## Surfaces and shape

### Elevation, shadows, depth

**Definition:** A layering model expressing which surfaces float above others, via shadows (light themes), lighter surface tones (dark themes), and z-order — communicating hierarchy, interactivity, and modality.

**Reasoning / Evidence:** Figure/ground perception ([Gestalt](01-core-principles.md#figureground)); shadows are a physical-world depth cue processed automatically. Material Design formalized elevation levels; consistent elevation = consistent meaning (menu always closer than page).

**When to use:**
- Temporary surfaces above content: menus, popovers, dialogs, toasts (highest).
- Interactive lift: cards/buttons raising on hover to signal affordance.
- Sticky chrome (app bars) gaining shadow only when content scrolls beneath (scroll-state signifier).

**When NOT to use / exceptions:**
- Decorative shadows on everything = muddy noise; elevation must mean something.
- Dark mode: shadows nearly invisible — express elevation with surface lightness steps instead ([#dark-mode-design](#dark-mode-design)).
- Flat-styled brands can encode layering with borders/scrims instead; the *model* matters, not the shadow.

**Pros:** Instant layer comprehension; modality signaling.
**Cons:** Shadow rendering cost (blur) on low-end devices; easy to overdo.

**Implementation guidance:**
- 4–6 elevation tokens max: base(0), raised card(1), sticky bar(2), dropdown/popover(3), dialog(4), toast(5).
- Realistic shadows: low offset + low blur + low opacity, layered (e.g., `0 1px 2px rgba(0,0,0,.08), 0 4px 12px rgba(0,0,0,.10)`); avoid single harsh spreads.
- One light source: shadows offset downward consistently.
- Dark theme: surface lightness +4–8% per level in the 1–5 range on a #121212 base.

**Related:** [01-core-principles.md](01-core-principles.md#figureground), [#dark-mode-design](#dark-mode-design).

**Sources:** [Material Design: Elevation](https://m3.material.io/styles/elevation/overview).

### Border radius and shape language

**Definition:** The corner-rounding scale and overall geometry vocabulary (sharp/rounded/pill/circular) applied consistently — shape as brand tone: sharp reads precise/serious/technical; rounded reads friendly/approachable; pill reads playful/consumer.

**Reasoning / Evidence:** Shape-emotion associations replicated in design research (curved contours preferred, sharp associated with threat/precision); consistency matters more than the chosen value — mixed radii read as sloppy.

**When to use:**
- Define a radius token scale at project start (e.g., 4 / 8 / 12 / 16 / full).
- Nested radii: inner radius = outer radius − padding (or inner elements look wrong).
- Match brand: fintech/pro tools smaller radii (2–6px); consumer/social larger (8–16px, pills).

**When NOT to use / exceptions:**
- Pill buttons truncate text sooner and read as tags/chips in some contexts — check for confusion with your chip components.
- Full-round avatars vs squared content thumbnails: keep object-type shape rules consistent.

**Pros:** Cheap brand cohesion; instant recognizability of interactive object classes.
**Cons:** Trend churn (radius inflation/deflation cycles); nested-radius math often botched.

**Implementation guidance:**
- 3–4 radius tokens: sm(4) inputs/chips, md(8) cards/buttons, lg(12–16) modals/sheets, full(9999) pills/avatars.
- Same-class = same radius everywhere (all buttons identical).
- Nested: `inner = max(outer − gap, 0)`.

**Related:** [#visual-hierarchy](#visual-hierarchy), [10-design-systems.md](10-design-systems.md#design-tokens).

**Sources:** [Bar & Neta 2006 (curvature preference)](https://doi.org/10.1111/j.1467-9280.2006.01759.x).

---

## Iconography and imagery

### Iconography

**Definition:** Pictographic symbols for actions/objects/states. Core truth: very few icons are universally understood (home, search magnifier, cart, play, gear-settings roughly qualify); everything else needs a label.

**Reasoning / Evidence:** NN/g icon-usability research: icon-only interfaces consistently fail recognition tests; icon+label outperforms both icon-only and label-only for speed and accuracy. Icons still add value: faster visual scanning anchors, target enlargement, aesthetic rhythm.

**When to use:**
- Icon + visible label as default for navigation and primary actions.
- Icon-only permitted where space is severe AND the icon is standard AND a tooltip + accessible name exist (toolbar density cases).
- Status icons paired with text ([#semantic-color](#semantic-color)).

**When NOT to use / exceptions:**
- Never icon-only for ambiguous/abstract concepts (what's "sync" vs "refresh" vs "repeat"?).
- Don't invent novel metaphors for established actions — use the platform's standard glyphs ([Jakob's Law](01-core-principles.md#jakobs-law)).
- Hamburger/kebab/meatball overflow menus: recognized now, but burying primary nav in them has measured discoverability cost ([08-navigation-ia.md](08-navigation-ia.md#hamburger-menu)).

**Pros:** Scan anchors, space efficiency, language-independence (for the truly universal few).
**Cons:** Recognition failure risk; maintenance of a consistent custom set is costly — use an established library (Material Symbols, SF Symbols, Lucide, Phosphor).

**Implementation guidance:**
- One family, one stroke weight, one corner style across the product; sizes 16/20/24px on the pixel grid (unaligned icons blur).
- Icon-only buttons: `aria-label` + tooltip (300–500ms delay); hit area ≥44px even if glyph is 20px.
- Test recognition: label-free icon test with 5 users; <80% recognition → needs its label.
- Decorative icons: `aria-hidden="true"`.

**Related:** [05-accessibility.md](05-accessibility.md#screen-readers), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#icon-only-interfaces-without-labels).

**Sources:** [NN/g: Icon Usability](https://www.nngroup.com/articles/icon-usability/).

### Imagery and illustration

**Definition:** Photography, illustration, and decorative graphics in UI — brand tone carriers that must serve, not obstruct, content.

**Reasoning / Evidence:** Eyetracking: users ignore stock photography ("filler" detection is fast and merciless) but study relevant images (real products, real people relevant to the task, diagrams). Faces attract fixation strongly — and can pull attention *away* from content; gaze direction of pictured people steers viewer attention.

**When to use:**
- Product photos (multiple angles, zoom — commerce conversions depend on them).
- Diagrams/screenshots in docs (show, don't tell).
- Illustration for empty states/onboarding — with restraint and a consistent style.

**When NOT to use / exceptions:**
- Generic stock (handshakes, headset-women) actively signals inauthenticity — worse than nothing.
- Background images behind text without heavy scrims ([#contrast](#contrast)).
- Large hero media that delays LCP ([13-platform-web.md](13-platform-web.md#core-web-vitals)).

**Pros:** Emotional tone, conversion lift (relevant imagery), comprehension (diagrams).
**Cons:** Payload weight; localization/diversity considerations; decoration crowding function.

**Implementation guidance:**
- Every content image: meaningful `alt`; decorative: empty `alt=""` ([05-accessibility.md](05-accessibility.md#screen-readers)).
- Reserve layout space (aspect-ratio) to prevent CLS; responsive `srcset`; lazy-load below-fold only.
- Pictured people should face/gaze toward the key content (attention steering).

**Related:** [13-platform-web.md](13-platform-web.md#responsive-design), [05-accessibility.md](05-accessibility.md).

**Sources:** [NN/g: Photos as Web Content](https://www.nngroup.com/articles/photos-as-web-content/).

---

## Dark mode

### Dark mode design

**Definition:** A full alternate theme with dark surfaces and light text — not a color inversion. Requires its own surface ramp, adjusted accents, and elevation model.

**Reasoning / Evidence:** User demand is strong (comfort in low light, OLED battery savings, preference); readability research is nuanced: light-on-dark slightly reduces reading performance for most sighted users at length (halation makes thin light text glow), while users with cataracts/photophobia benefit. Conclusion: offer both, respect the OS setting, default to it.

**When to use:**
- Media/entertainment, developer tools, night-context apps: dark-first is reasonable.
- Everything else: support both via `prefers-color-scheme`, with an in-app override persisted per user.

**When NOT to use / exceptions:**
- Long-form reading products: keep light as default; dark optional.
- Brand-critical color surfaces may not translate — define dark-specific brand tokens rather than skipping dark mode.
- Don't ship dark mode by inverting: unadjusted saturated colors vibrate, shadows vanish, images/illustrations may need variants.

**Pros:** User comfort/preference, OLED power savings (only with true dark surfaces), pro-tool aesthetics.
**Cons:** Doubles the theme QA surface; every asset/color decision made twice; done badly it's worse than absent.

**Implementation guidance:**
- Base surface: dark grey `#121212`-region, not pure black `#000` (pure black maximizes halation contrast and makes elevation impossible; exception: OLED "true black" modes as an extra option).
- Elevation = lighter surfaces: overlay white at ~5/8/11/12/14% for levels 1–5 (Material convention) — shadows are secondary in dark.
- Desaturate accents: mid-tone versions of brand hues (saturated colors on dark vibrate); text not pure white — `#dedede`–`#e0e0e0` (~87% white) for body (`#f5f5f5` is ~96% white, brighter than the Material high-emphasis convention).
- Contrast still applies fully: 4.5:1 text — check every token pair in dark separately.
- Elevated surfaces get *lighter*, never darker; text emphasis levels (Material convention): high-emphasis ~87% white, medium-emphasis 60%, disabled 38%. Note: disabled text is exempt from WCAG contrast minimums (the WCAG 1.4.3 exemption for inactive components), but medium- and high-emphasis text is not — both must still hit 4.5:1.
- Media: slightly dim images (`filter: brightness(.9)`) optional; provide dark variants of illustrations with white-filled shapes.
- Implement as token remap ([10-design-systems.md](10-design-systems.md#theming-dark-mode-multi-brand-density)): semantic tokens flip, components untouched; respect `prefers-color-scheme`, offer light/dark/system triple toggle.

**Related:** [#elevation-shadows-depth](#elevation-shadows-depth), [#contrast](#contrast), [10-design-systems.md](10-design-systems.md#theming-dark-mode-multi-brand-density).

**Sources:** [NN/g: Dark Mode](https://www.nngroup.com/articles/dark-mode/), [Material: Dark theme](https://m2.material.io/design/color/dark-theme.html).

---

## Trends and styles

### Flat design and minimalism

**Definition:** The post-2013 aesthetic removing skeuomorphic texture, gradients, and shadows in favor of flat surfaces, simple shapes, and typography-led hierarchy. "Flat 2.0" restored subtle shadows/layering after usability losses.

**Reasoning / Evidence:** NN/g measured the cost: flat designs with weak clickability signifiers caused click uncertainty and increased time-on-task (users hover-hunting for what's interactive). Aesthetic benefits are real (clean, fast, scalable); signifier loss is the failure mode, not flatness itself.

**When to use:**
- Default modern baseline — with sufficient signifiers retained (Flat 2.0: subtle elevation, filled buttons, clear link styling).

**When NOT to use / exceptions:**
- Audit any design where buttons are borderless text, links are color-only at low contrast, or inputs lack visible boundaries — restore signifiers.

**Pros:** Performance (no texture assets), scalability, timelessness relative to heavy skeuomorphism, content focus.
**Cons:** Documented signifier weakness; sameness across products (brand differentiation burden shifts to color/type/motion).

**Implementation guidance:**
- Non-negotiable signifiers: buttons visibly bounded (fill or border), inputs visibly bounded, links distinguishable in body text (color + underline), interactive hover/focus states.
- Clickability test: >90% of users identify interactive elements from a static screenshot.

**Related:** [01-core-principles.md](01-core-principles.md#signifiers), [04-interaction-design.md](04-interaction-design.md#affordances-and-signifiers-in-practice).

**Sources:** [NN/g: Flat Design](https://www.nngroup.com/articles/flat-design/), [NN/g: Flat-Design Best Practices](https://www.nngroup.com/articles/flat-design-best-practices/).

### Neumorphism

**Definition:** "Soft UI" (2020 trend): elements extruded from the background via paired light/dark shadows, same-color surfaces, minimal contrast.

**Reasoning / Evidence:** Aesthetically distinctive; usability analysis is broadly negative — the style's defining trait (low contrast between element and background) directly violates WCAG 1.4.11 (3:1 non-text contrast) and cripples state signaling (pressed vs raised is nearly indistinguishable at a glance).

**When to use:**
- Decorative accents, portfolio/branding contexts, isolated non-critical widgets (music player skin) where an accessible alternative exists.

**When NOT to use / exceptions:**
- Any functional UI at scale; forms; anything requiring accessibility conformance — the style cannot pass 1.4.11 without abandoning its core look.

**Pros:** Distinctive softness; tactile character.
**Cons:** Inaccessible by construction; disabled/pressed/focus states barely expressible; fails in dark mode and low-quality screens.

**Implementation guidance:**
- If used decoratively: keep all text/controls on conventional high-contrast styling; treat neumorphic surfaces as backgrounds only.

**Related:** [#contrast](#contrast), [05-accessibility.md](05-accessibility.md#contrast-minimum).

**Sources:** [WCAG 1.4.11](https://www.w3.org/WAI/WCAG22/Understanding/non-text-contrast.html).

### Glassmorphism

**Definition:** Frosted-glass surfaces: background blur + translucency + thin light border (macOS/iOS materials, Windows acrylic/mica, Apple's Liquid Glass direction).

**Reasoning / Evidence:** Platform vendors ship it because blur preserves context (you sense what's behind a sidebar) while separating layers. Risks: text over translucency has *unpredictable* contrast (depends on what scrolls beneath) and blur is GPU-expensive.

**When to use:**
- Chrome layers (sidebars, bars, overlays) on platforms where it's native convention ([14-platform-desktop.md](14-platform-desktop.md#platform-conventions)).
- Large surfaces with controlled/blurred-enough backgrounds.

**When NOT to use / exceptions:**
- Body text containers over arbitrary content; small text on thin translucent chips.
- Low-end/battery-sensitive targets: provide reduced-transparency fallback; honor OS "reduce transparency" accessibility setting.

**Pros:** Spatial layering with context; native feel on Apple/Microsoft platforms.
**Cons:** Contrast unpredictability (worst-case must still pass 4.5:1); GPU cost; trend-dependency.

**Implementation guidance:**
- Blur radius high enough (≥20px) + tint layer (e.g., 60–80% surface opacity) so worst-case background can't break contrast.
- Test over pathological backgrounds (white text areas, busy imagery).
- CSS: `backdrop-filter: blur(20px)` + fallback solid; respect `prefers-reduced-transparency`.

**Related:** [#contrast](#contrast), [14-platform-desktop.md](14-platform-desktop.md#platform-conventions).

**Sources:** [MDN: backdrop-filter](https://developer.mozilla.org/en-US/docs/Web/CSS/backdrop-filter), [Microsoft: Materials](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/materials).

### Brutalism

**Definition:** Deliberately raw web aesthetics: system fonts, harsh contrast, exposed structure, anti-polish — from actual brutalism (unstyled rawness) to "neo-brutalism" (thick borders, hard shadows, saturated blocks).

**Reasoning / Evidence:** Differentiation strategy in a sea of same-looking minimal UIs; attention economics favor distinctiveness ([Von Restorff](01-core-principles.md#von-restorff-effect) at brand scale). Usability is orthogonal: brutalist styling can sit atop perfectly conventional structure — or excuse chaos.

**When to use:**
- Brand-forward contexts: portfolios, culture/media, campaign sites, dev-tool branding where memorability outweighs mainstream polish.

**When NOT to use / exceptions:**
- Broad-audience transactional products (commerce checkouts, government, health): unfamiliarity taxes exactly the users who can least afford it.
- Never let the aesthetic break conventions of function: navigation, links, buttons still must be findable ([Jakob's Law](01-core-principles.md#jakobs-law)).

**Pros:** Memorable; fast (no asset weight); high contrast often *helps* accessibility.
**Cons:** Polarizing; easily slides from styled-raw to actually-unusable; dated risk as the trend cycles.

**Implementation guidance:**
- Keep structure conventional (standard nav placement, visible labels); apply brutalism to surfaces, type, and color only.
- Neo-brutalist hard shadows/borders actually restore strong signifiers — usable foundation if contrast is managed.

**Related:** [#flat-design-and-minimalism](#flat-design-and-minimalism), [01-core-principles.md](01-core-principles.md#jakobs-law).

**Sources:** [NN/g: Brutalism and Antidesign](https://www.nngroup.com/articles/brutalism-antidesign/).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Visual hierarchy | 3 levels; two cues per jump; squint test | ~50ms first impression |
| Type scale | Ratio-based token scale, rem units | Base 16px; ratios 1.2–1.333; 5–7 sizes |
| Line height/length | Comfortable leading, capped measure | 1.5 body; 45–75ch (66 ideal) |
| Fonts | ≤2 families; system stack for apps | ≤4 weights; tabular nums in tables |
| Emphasis | Weight over caps/italic; budget it | ≥2 weight steps; ≤10% of words |
| Palette | Neutral-dominant ramps + scarce accent | 60-30-10; 8–12 neutral steps |
| Semantic color | Color + icon + text, always | Red=error only; amber not yellow for text |
| Contrast | Check every token pair, both themes | 4.5:1 text; 3:1 large text/UI; 7:1 AAA |
| Color blindness | Second channel always | ~8% of men; Okabe-Ito palette |
| OKLCH | Generate ramps perceptually | CSS Color 4; sRGB fallback |
| 8pt grid | All spacing in 4/8 multiples | Within ≤8–12, between ≥16–24 |
| Whitespace | Active grouping tool | Heading: above ≥1.5–2× below |
| Grids | 12-col desktop; one left edge per region | Gutters 16–24px; max 1200–1440px |
| Density | Audience-driven; token-switchable | Rows 32/40/48; targets ≥24px always |
| Elevation | 4–6 meaningful levels; one light source | Dark mode: +4–8% lightness per level |
| Radius | 3–4 tokens; nested = outer − gap | sm 4 / md 8 / lg 12–16 / full |
| Icons | Icon + label default; standard glyphs | 16/20/24px; ≥80% recognition or label it |
| Imagery | Relevant or nothing; never stock filler | alt text always; reserve aspect-ratio |
| Dark mode | Own theme, not inversion; respect OS | #121212 base; ~87% white text; 4.5:1 holds |
| Flat design | Keep signifiers (Flat 2.0) | >90% clickability recognition |
| Neumorphism | Decorative only — fails contrast | Violates 3:1 (1.4.11) |
| Glassmorphism | Chrome layers, tint for worst case | blur ≥20px; 60–80% tint |
| Brutalism | Style surfaces, keep structure conventional | — |
