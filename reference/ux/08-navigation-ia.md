# Information Architecture & Navigation

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers how to structure information (organization schemes, hierarchy depth vs breadth, validation methods) and how to expose that structure through navigation patterns: top nav, sidebars, tabs, hamburger menus, mega menus, breadcrumbs, bottom navigation, hub-and-spoke, wizards, command palettes, faceted filtering, sorting, search, pagination/infinite scroll, URL design, wayfinding, scope/context safety, and multi-level menu ergonomics. Interaction mechanics of individual controls (menus as widgets, gestures) are in [04-interaction-design.md](04-interaction-design.md); mobile-platform specifics in [15-platform-mobile.md](15-platform-mobile.md); form-flow structure in [07-forms-and-input.md](07-forms-and-input.md); search/nav microcopy in [09-content-ux-writing.md](09-content-ux-writing.md); accessibility of navigation (landmarks, skip links, focus order) in [05-accessibility.md](05-accessibility.md); data tables, data grids, and data-dense UI in [12-data-tables-dashboards.md](12-data-tables-dashboards.md).

## Contents

- [IA fundamentals](#ia-fundamentals)
  - [IA vs navigation](#ia-vs-navigation)
  - [Organization schemes](#organization-schemes)
  - [Hierarchy depth vs breadth](#hierarchy-depth-vs-breadth)
  - [Card sorting](#card-sorting)
  - [Tree testing](#tree-testing)
  - [Information scent](#information-scent)
- [Primary navigation patterns](#primary-navigation-patterns)
  - [Top navigation bar](#top-navigation-bar)
  - [Sidebar / left navigation](#sidebar--left-navigation)
  - [Tabs](#tabs)
  - [Hamburger menu](#hamburger-menu)
  - [Mega menus](#mega-menus)
  - [Bottom navigation (mobile)](#bottom-navigation-mobile)
  - [Hub-and-spoke](#hub-and-spoke)
  - [Sequential / wizard navigation](#sequential--wizard-navigation)
  - [Command palette](#command-palette)
- [Secondary and supporting navigation](#secondary-and-supporting-navigation)
  - [Breadcrumbs](#breadcrumbs)
  - [Faceted navigation and filtering](#faceted-navigation-and-filtering)
  - [Sorting](#sorting)
  - [Wayfinding and you-are-here signals](#wayfinding-and-you-are-here-signals)
  - [Scope and context safety](#scope-and-context-safety)
  - [Footer navigation](#footer-navigation)
- [Search](#search)
  - [Search vs browse](#search-vs-browse)
  - [Autocomplete and suggestions](#autocomplete-and-suggestions)
  - [No-results recovery](#no-results-recovery)
  - [Scoped search](#scoped-search)
- [Content traversal](#content-traversal)
  - [Pagination vs infinite scroll vs load more](#pagination-vs-infinite-scroll-vs-load-more)
  - [Deep linking and URL design](#deep-linking-and-url-design)
- [Navigation for data-heavy UIs](#navigation-for-data-heavy-uis)
- [Menu ergonomics](#menu-ergonomics)
  - [Multi-level menu ergonomics: hover intent and safe triangles](#multi-level-menu-ergonomics-hover-intent-and-safe-triangles)
- [Quick reference](#quick-reference)

---

## IA fundamentals

### IA vs navigation

**Definition:** Information architecture is the underlying organization, structure, and nomenclature defining relationships between a product's content and functionality. Navigation is the set of UI components (menus, breadcrumbs, filters, links) that let users move through that structure. IA informs UI; it is not itself visible on screen.

**Reasoning / Evidence:** NN/g describes navigation as "the tip of the iceberg" sitting atop the IA. Choosing navigation components before defining IA is risky: e.g., an inverted-L layout (top bar + left rail) accommodates at most ~4 tiers; discovering deeper content later forces either a redesign or cramming content into a structure that doesn't fit. Users never see the IA directly but feel its quality — content that is grouped and named to match their mental models feels findable; mismatches produce friction regardless of how polished the menus look. Structure around user tasks, not internal teams, backend services, or feature ownership: user task categories rarely match implementation boundaries, and org-chart IA is a classic failure.

**When to use:**
- At project start: run a content inventory and audit, define groupings, taxonomy, and metadata before wireframing navigation.
- During redesigns: re-validate IA before changing navigation chrome; most "navigation problems" are actually labeling/grouping problems.
- When analytics show high search usage from interior pages, pogo-sticking between siblings, or heavy back-button use — classic IA warning signs.

**When NOT to use / exceptions:**
- A first-pass IA is enough to start wireframing; do not block prototyping on a "final" IA.
- Single-purpose tools with <10 destinations need minimal formal IA work — a flat list validated with a quick label check suffices.

**Pros:**
- Prevents expensive late-stage navigation rework.
- Produces a shared vocabulary (taxonomy) across design, engineering, and content teams.

**Cons:**
- IA work (inventory, audit, taxonomy) is slow and invisible to stakeholders; hard to defend in fast-moving projects.
- Over-formalized IA can ossify; real products need polyhierarchy and cross-links, not a single pure tree.

**Implementation guidance:**
- Deliverables: content inventory spreadsheet, site map diagram (label tiers), controlled vocabulary, metadata schema for related-content generation.
- Sequence: inventory → audit → grouping (card sorting) → taxonomy → validate (tree testing) → then design navigation components.
- Triangulate inputs: card sorting, tree testing, search logs, support tickets, and user interviews all reveal the user's vocabulary and groupings.
- Allow polyhierarchy (an item appearing under multiple parents) when card sorts show users split ~40/60 on placement, but designate one canonical parent for breadcrumbs and URLs.
- Provide multiple paths to high-value content when different user groups approach it differently.
- Budget rule of thumb: for sites >200 pages, IA definition and validation should get at least 2–3 weeks before high-fidelity navigation design.

**Related:** [#card-sorting](#card-sorting), [#tree-testing](#tree-testing), [17-ux-process-research.md](17-ux-process-research.md), [#breadcrumbs](#breadcrumbs).

**Sources:** [NN/g: The Difference Between IA and Navigation](https://www.nngroup.com/articles/ia-vs-navigation/), [NN/g: IA Warning Signs in Analytics](https://www.nngroup.com/articles/ia-warning-signs-analytics/), [NN/g: Polyhierarchies Improve Findability](https://www.nngroup.com/articles/polyhierarchy/).

### Organization schemes

**Definition:** The principle by which content is grouped and ordered. Exact schemes: alphabetical, chronological, geographical. Ambiguous (subjective) schemes: topical (by subject), task-based (by user goal), audience-based (by user type), metaphor-driven. Richard Saul Wurman's LATCH framework summarizes the options: Location, Alphabet, Time, Category, Hierarchy (magnitude).

**Reasoning / Evidence:** Exact schemes support known-item lookup (user knows the name) but fail exploration (user must know the label before finding it). Ambiguous schemes support browsing and learning but require the designer's groupings to match user mental models — hence the need for card sorting. NN/g research specifically warns against audience-based navigation as a primary scheme: users often don't self-identify with the offered segments, content duplicated across audiences drifts out of sync, and users fear missing content placed under another audience's door.

**When to use:**
- Alphabetical: known-item lists >~25 entries with well-known names (countries, employee directories, glossaries, drug names).
- Chronological: time-anchored content (news, releases, activity feeds, changelogs, transaction histories). Default newest-first.
- Topical: general content sites and docs — the workhorse scheme; validate labels with users.
- Task-based: application navigation where users arrive with a goal ("Send money," "Create invoice"); works when tasks are few (<8) and distinct.
- Audience-based: only when audiences are truly disjoint with legally/functionally distinct content (e.g., patients vs clinicians, consumer vs institutional banking) and users reliably self-identify.

**When NOT to use / exceptions:**
- Alphabetical fails when users don't know the exact term ("Is it under 'Vacation' or 'Time off' or 'Leave'?") — it is a lazy default that hides poor categorization.
- Audience-based fails when one person belongs to multiple audiences, or content overlap exceeds ~30%.
- Chronological alone fails for evergreen content; pair with topical access.

**Pros:**
- Exact schemes: zero ambiguity, trivially maintainable, no research required.
- Ambiguous schemes: support exploration, teach the domain, group related tasks.

**Cons:**
- Exact schemes: no support for users who lack the vocabulary.
- Ambiguous schemes: expensive to get right; every new item raises a classification decision; internal org-chart structures leak in ("Conway's law IA").

**Implementation guidance:**
- Offer multiple access schemes to the same content where feasible: topical browse + search + A–Z index (for reference content) + chronological "what's new."
- Never mirror internal org structure; test whether users can predict where content lives (tree test success target ≥70–80% for top tasks).
- Keep categories mutually exclusive in meaning where possible; contested items are candidates for polyhierarchy or renaming.
- For task-based navigation, phrase labels as verb phrases ("Track a shipment"), 1–3 words for bars, up to ~5 words for link lists.
- If audience-based structure is mandated, deduplicate content by cross-linking to a single canonical page rather than copying.

**Related:** [#card-sorting](#card-sorting), [#tree-testing](#tree-testing), [09-content-ux-writing.md](09-content-ux-writing.md) (label writing), [02-cognitive-foundations.md](02-cognitive-foundations.md) (mental models).

**Sources:** [NN/g: Organize Information with LATCH](https://www.nngroup.com/videos/latch-framework/), [NN/g: Audience-Based Navigation — 5 Reasons to Avoid It](https://www.nngroup.com/articles/audience-based-navigation/), [NN/g: IA vs Navigation](https://www.nngroup.com/articles/ia-vs-navigation/).

### Hierarchy depth vs breadth

**Definition:** The trade-off between flat hierarchies (many options per level, few levels) and deep hierarchies (few options per level, many levels). A structure of 3 levels × 10 options each covers 1,000 items; 10 levels × 2 options each covers 1,024 with vastly more clicks.

**Reasoning / Evidence:** NN/g ("Flat vs. Deep Website Hierarchies") and classic HCI menu research (Miller 1981; Kiger 1984; Larson & Czerwinski 1998) converge: moderately broad and shallow beats narrow and deep. Each level traversed is a decision point where the user can misinterpret a label and go down the wrong branch; errors compound multiplicatively with depth. Scanning a longer list of well-labeled options is faster and less error-prone than committing to a sparse label and hoping. Deep structures also hide most content from the landing-page glance, weakening information scent, and increase backtracking cost when users guess wrong. The old "7±2 menu items" rule is a misapplication of memory research — navigation is recognition, not recall, so lists of 10–20 scannable options are fine.

**When to use:**
- Prefer breadth (2–3 levels total) for content sites up to tens of thousands of pages, exposing structure via mega menus or landing pages.
- Add depth only when categories are genuinely nested in users' mental models (e.g., Products → Category → Subcategory → Item in retail).

**When NOT to use / exceptions:**
- Extreme breadth (50+ undifferentiated links) shifts cost to visual scan; group with headings (a mega menu is structured breadth, not raw breadth).
- Very small screens tolerate less breadth per view; mobile may present the same IA with progressive disclosure without changing the underlying depth.
- Enterprise apps with true containment (Org → Project → Environment → Resource) rightly go deeper; compensate with breadcrumbs and search.

**Pros (broad/shallow):**
- Fewer decisions and clicks to any leaf; errors surface early and cheaply.
- Structure is visible; users build an accurate mental map faster.

**Cons (broad/shallow):**
- Long option lists need strong grouping, ordering, and typography to stay scannable.
- Category pages become large and need their own layout design.

**Implementation guidance:**
- Target ≤3 clicks structure for content sites, not as dogma ("3-click rule" is not empirically supported — users tolerate more clicks if scent stays strong) but as a depth-budget heuristic.
- Per level, 5–12 grouped options is a comfortable band for bars/menus; category landing pages can present 20–40 links grouped under headings.
- Flatten where categories are small; use landing pages for meaningful category groups; keep critical actions available near the content they affect.
- Diagnose with analytics: high exit rates on mid-tier category pages and long paths to conversion suggest excess depth.
- When deepening, add wayfinding: breadcrumbs at depth ≥3, persistent local nav or object headers, and "you are here" highlighting.

**Related:** [#information-scent](#information-scent), [#mega-menus](#mega-menus), [#breadcrumbs](#breadcrumbs), [02-cognitive-foundations.md](02-cognitive-foundations.md).

**Sources:** [NN/g: Flat vs. Deep Website Hierarchies](https://www.nngroup.com/articles/flat-vs-deep-hierarchy/), [NN/g: The 3-Click Rule is a Myth](https://www.nngroup.com/articles/3-click-rule/), Larson & Czerwinski, "Web Page Design: Implications of Memory, Structure and Scent for Information Retrieval," CHI 1998.

### Card sorting

**Definition:** A research method where participants group content items (cards) into categories and often name the groups. Open sort: participants create and name their own groups (generative — discovers mental models). Closed sort: participants place cards into predefined categories (evaluative — tests a proposed scheme). Hybrid allows both.

**Reasoning / Evidence:** Directly elicits users' mental models of content relationships instead of relying on internal assumptions. Quantitative similarity matrices and dendrograms across participants reveal consensus clusters and contested items. NN/g's long-standing guidance: ~15 participants for card sorting yields a correlation of ~0.90 with larger samples' groupings — diminishing returns beyond that.

**When to use:**
- Designing or restructuring the IA of a content-heavy site or app (≥30 distinct content items).
- Deciding between competing categorization schemes; identifying items with no natural home.
- Generating candidate category labels in users' own vocabulary.

**When NOT to use / exceptions:**
- Cannot validate findability — card sorting shows how users group, not whether they can find. Follow with tree testing.
- Useless for tiny IAs (<15 items) or purely task-flow structures (use task analysis instead).
- Sorting >60 cards per session exhausts participants; sample the content or split sessions.

**Pros:**
- Cheap, fast, runs remotely and unmoderated at scale (OptimalSort, Maze, UXtweak).
- Produces user-generated labels, gold for navigation copy.

**Cons:**
- Results are content-out, not task-in; users group in the abstract without a goal.
- Aggregation can mask distinct sub-populations with different models.

**Implementation guidance:**
- 30–60 cards per session; 15–20 participants for qualitative confidence, 30+ per segment for statistical cluster analysis.
- Write card labels without shared leading words (users will sort by surface wording, not meaning).
- Run open sort first to discover structure → design candidate IA → validate with closed sort or tree test.
- Analyze: similarity matrix (% of participants pairing two items), standardization grid for label frequency; flag items with <50% placement agreement as candidates for polyhierarchy or renaming.

**Related:** [#tree-testing](#tree-testing), [#organization-schemes](#organization-schemes), [17-ux-process-research.md](17-ux-process-research.md).

**Sources:** [NN/g: Card Sorting — Why & When](https://www.nngroup.com/videos/card-sorting-why-when/), [NN/g: Card Sorting: How Many Users to Test](https://www.nngroup.com/articles/card-sorting-how-many-users-to-test/), [NN/g: Card Sorting Definition](https://www.nngroup.com/articles/card-sorting-definition/).

### Tree testing

**Definition:** A findability evaluation ("reverse card sort") where participants locate items in a text-only representation of the hierarchy — no visual design, no search — by clicking through category labels. Metrics: task success, directness (first-click-correct path without backtracking), time.

**Reasoning / Evidence:** Isolates the IA variable: because there is no layout, color, or imagery, failures are attributable purely to structure and labeling. First-click research (Bailey et al.; MeasuringU summaries) shows users whose first click is correct are roughly 2× as likely to succeed overall (≈87% vs ≈46% task success), making directness a strong predictor of real navigation success.

**When to use:**
- Validating a candidate IA before investing in navigation UI design.
- A/B comparing two candidate trees (e.g., topical vs task-based top level).
- Benchmarking an existing site's structure before a redesign to prove/deny "the IA is fine."

**When NOT to use / exceptions:**
- Doesn't evaluate visual navigation design, mega menu layout, or scent from contextual links — pair with usability testing on prototypes.
- Doesn't test search-dominant behavior; if most real users search first, weight tree-test findings accordingly.

**Pros:**
- Fast, unmoderated, quantitative; 50 participants per tree is achievable in days.
- Directly comparable metrics across design iterations.

**Cons:**
- Artificially forces browsing; text-only trees lack the scent aids (icons, descriptions, promotions) of real pages.
- Task wording easily leaks answers ("Find the returns policy" when a node is named "Returns").

**Implementation guidance:**
- 8–12 tasks per session max; write tasks as scenarios avoiding words that appear in node labels.
- 50+ participants per tree for stable percentages; 15–20 acceptable for directional reads.
- Success benchmarks: ≥80% success and ≥70% directness on critical tasks is strong; <60% success on a top task means restructure or relabel before building.
- Iterate: fix worst labels, retest — tree tests are cheap enough to run 2–3 rounds.

**Related:** [#card-sorting](#card-sorting), [#hierarchy-depth-vs-breadth](#hierarchy-depth-vs-breadth), [17-ux-process-research.md](17-ux-process-research.md).

**Sources:** [NN/g: Tree Testing](https://www.nngroup.com/articles/tree-testing/), [NN/g: Navigation & IA Tests](https://www.nngroup.com/articles/navigation-ia-tests/), Bailey, Wolfson et al., first-click studies (usability.gov / MeasuringU summaries: [First Click Testing](https://measuringu.com/first-click/)).

### Information scent

**Definition:** The user's imperfect estimate, from proximal cues (link labels, snippets, icons, context), of whether a path leads to their goal. From Pirolli & Card's information foraging theory (Xerox PARC): users patrol for scent like animals foraging, following gradients and abandoning patches when scent drops.

**Reasoning / Evidence:** Information foraging is among the most robust theoretical frameworks in HCI. Practical corollaries: users don't click links they don't understand; they back out (pogo-stick) when a page's content mismatches the promise of its label; generic labels ("Products," "Solutions," "Resources") emit weak scent and depress click-through. Scent, not click count, determines perceived effort — users happily click 5 times when each click visibly converges on the goal.

**When to use:**
- Always: every link label, nav item, button, and card title is a scent carrier.
- Diagnosing navigation failures: weak scent explains most "users can't find X" findings.

**When NOT to use / exceptions:**
- Scent can be too strong toward the wrong thing: attractive but off-goal links (promos in nav) hijack tasks — police competing scent near critical paths.

**Pros:**
- Explains and predicts navigation behavior; converts vague "make it intuitive" into testable label quality.

**Cons:**
- Specific, descriptive labels are longer; tension with compact navigation bars.

**Implementation guidance:**
- Front-load labels with information-carrying words; make sibling links mutually exclusive in meaning (a user should never hesitate between two nav items for the same goal).
- Prefer specific over clever: "Pricing" beats "Explore"; "Returns & refunds" beats "Customer care."
- Enrich scent where space allows: mega menu items with 1-line descriptions, search results with keyword-highlighted snippets.
- Test with 5-second first-click tests: show the page, ask "where would you click to do X?" — target ≥80% correct first clicks on primary tasks.

**Related:** [#hierarchy-depth-vs-breadth](#hierarchy-depth-vs-breadth), [09-content-ux-writing.md](09-content-ux-writing.md) (link text), [02-cognitive-foundations.md](02-cognitive-foundations.md).

**Sources:** [NN/g: Information Scent](https://www.nngroup.com/articles/information-scent/), Pirolli & Card, "Information Foraging," Psychological Review 1999, [NN/g: Information Foraging — A Theory of How People Navigate the Web](https://www.nngroup.com/articles/information-foraging/).

---

## Primary navigation patterns

Choosing among the patterns below is a function of information structure, usage frequency, screen size, and platform convention — not aesthetics. Decision questions before committing: How many top-level destinations exist? Are sections peers or steps? Does the user need persistent context while working? Is the platform mobile, desktop, or web?

### Top navigation bar

**Definition:** A horizontal bar of primary destinations across the top of the viewport, often with dropdowns/mega menus for second-tier items, typically paired with a logo (home link) left and utility items (search, account, cart) right.

**Reasoning / Evidence:** The most familiar web pattern; users' F-pattern scanning starts top-left, and decades of convention (Jakob's Law — users spend most time on other sites) make top-bar placement the default expectation for marketing and content sites. Persistent visibility gives constant scent and orientation; NN/g's hidden-navigation research shows visible navigation roughly doubles usage vs hidden (48% vs 27% of tasks on desktop).

**When to use:**
- Content/marketing sites with 3–8 top-level categories with short labels.
- When horizontal space suffices and the IA is shallow (≤3 tiers, with dropdowns/mega menus for tier 2).
- Paired with left local nav (inverted-L) for sites needing 3–4 tiers.

**When NOT to use / exceptions:**
- More than ~8 top-level items or long labels (localization expands text 30–40%; see [09-content-ux-writing.md](09-content-ux-writing.md)) — items overflow or truncate.
- Complex applications with many tools/objects: sidebar scales better.
- Growing IAs: horizontal space is fixed; adding a category forces relabeling or restructuring.

**Pros:**
- Maximally conventional; zero learning cost.
- Leaves full content width below; plays well with hero layouts.
- Persistent scent for the whole site.

**Cons:**
- Hard cap on item count; labels must be 1–2 words.
- Dropdowns introduce hover-ergonomics problems ([#multi-level-menu-ergonomics-hover-intent-and-safe-triangles](#multi-level-menu-ergonomics-hover-intent-and-safe-triangles)).
- Collapses to hamburger on mobile, inheriting its discoverability costs.

**Implementation guidance:**
- 3–8 items; order by user priority left-to-right (RTL mirrored); highlight the current section persistently.
- Sticky top bars: keep height ≤64px when stuck; consider shrink-on-scroll. Don't let a sticky bar plus a cookie banner consume >20% of viewport height.
- Make top-level items clickable landing pages AND dropdown triggers with care — split-button navs test poorly (NN/g); prefer whole-item-opens-menu plus a "View all X" link inside the menu.
- Ensure keyboard operability: menus open on Enter/Space, arrow-key traversal, Escape closes (see [05-accessibility.md](05-accessibility.md)).

**Related:** [#mega-menus](#mega-menus), [#sidebar--left-navigation](#sidebar--left-navigation), [#hamburger-menu](#hamburger-menu), [13-platform-web.md](13-platform-web.md).

**Sources:** [NN/g: Hamburger Menus and Hidden Navigation Hurt UX Metrics](https://www.nngroup.com/articles/hamburger-menus/), [NN/g: Killing Off the Global Navigation — One Trend to Avoid](https://www.nngroup.com/articles/killing-global-navigation-one-trend-avoid/), [NN/g: Split Buttons for Navigation Menus](https://www.nngroup.com/articles/split-buttons-navigation/).

### Sidebar / left navigation

**Definition:** A vertical column of navigation links along the left (or leading) edge, optionally collapsible to icons ("rail") and supporting nested groups via expand/collapse.

**Reasoning / Evidence:** Vertical lists scale better than horizontal bars: adding items costs scroll, not restructure; labels can be longer; scanning a left-aligned vertical list is faster than scanning a horizontal one (consistent left edge for the eye). This is why nearly every complex web application (admin consoles, IDEs, dashboards — Gmail, GitHub, Jira, AWS) converges on left nav: apps have many tools of similar importance and need visible mode/section context at all times.

**When to use:**
- Applications with 8+ sections/tools of comparable importance.
- IAs that grow over time or have variable-length, descriptive labels.
- When persistent visibility of section structure aids orientation during long sessions.
- E-commerce category pages (left nav doubles as facet container).

**When NOT to use / exceptions:**
- Marketing/landing pages: a sidebar steals 200–280px of width and signals "tool," not "story."
- Very few destinations (<5): a top bar or tabs is lighter.
- Small screens: sidebar must collapse to a drawer — plan the collapsed state, don't bolt it on.

**Pros:**
- Scales to dozens of items with grouping and scroll.
- Longer labels fit; supports icons + text (best recognition).
- Collapsible variants recover width when users know the icons.

**Cons:**
- Permanent horizontal real-estate cost (~15–20% at 1280px).
- Icon-only collapsed rails lose scent for infrequent users; icon comprehension without labels is poor for abstract functions.
- Nested expand/collapse trees can hide location if state resets between sessions.

**Implementation guidance:**
- Expanded width 240–280px; collapsed rail 48–72px with tooltips on hover/focus showing full labels.
- Persist expand/collapse state per user; auto-expand the group containing the current page and highlight the current item (background fill + bold, not color alone).
- Limit nesting inside a sidebar to 2 visible levels; deeper structures should get a second pane or in-page navigation.
- Group with sticky section headers when list exceeds ~12 items; separate global items (Settings, Help) at the bottom.
- Collapse control: place at the bottom or top of the sidebar itself; use a chevron/rail icon, not a hamburger (different semantics).

**Related:** [#top-navigation-bar](#top-navigation-bar), [14-platform-desktop.md](14-platform-desktop.md), [12-data-tables-dashboards.md](12-data-tables-dashboards.md).

**Sources:** [NN/g: Vertical Navigation](https://www.nngroup.com/videos/vertical-navigation/), [NN/g: Left-Side Vertical Navigation on Desktop](https://www.nngroup.com/articles/vertical-nav/), [NN/g: Complex Application Design](https://www.nngroup.com/articles/complex-application-design/).

### Tabs

**Definition:** A control that lets users view one content panel at a time from a labeled set. Two distinct types that must not be mixed: **in-page (module) tabs** switch views within the current page (same layout, different data; instant); **navigation (page-level) tabs** move between pages (different content; loading expected).

**Reasoning / Evidence:** NN/g "Tabs, Used Right": tabs work by chunking content and using progressive disclosure, reducing cognitive load — but only when groupings are clean, labels are short and predictive, and users never need to compare content across tabs (cross-tab comparison forces memory-taxing switching). Mixing a navigation link into an in-page tab set (one tab that leaves the page) breaks learned expectations and disorients users. Stacked multi-row tabs (classic Amazon 2000 failure) destroy spatial memory because selecting a back-row tab forces reordering.

**When to use:**
- Lengthy content with 2–6 clear, mutually exclusive groupings of similar type.
- Content of unequal importance where one default view serves most visits (put the high-use tab first and default-selected).
- Alternate views of the same object (Details / Activity / Permissions on a record page) — peer views within the same context.

**When NOT to use / exceptions:**
- Users need to see or compare content from multiple tabs simultaneously → one long page with headings.
- Labels can't be kept to 1–2 words → groupings are too complex for tabs.
- More tabs than fit in one row → carousel/overflow tabs hide options; restructure instead.
- Sequential processes → wizard/stepper, not tabs (tabs imply free order).
- Mobile long-form content → accordions usually beat tabs on small screens.

**Pros:**
- Familiar, compact, supports progressive disclosure.
- Preserves context (in-page tabs keep the user "in place").

**Cons:**
- Content in non-default tabs gets dramatically less attention — never hide critical content there.
- Poor for printing/searching within page (hidden panels), and Ctrl+F misses hidden content unless rendered.

**Implementation guidance:**
- Single row only; 2–6 tabs; labels 1–2 words, sentence case, never ALL CAPS (legibility penalty).
- Tab list above the panel (not right/left/bottom); ≥2 selection indicators on the active tab (e.g., underline + bold), critical when only 2 tabs exist.
- Keep unselected tabs clearly visible/legible (≥4.5:1 contrast) so options aren't mistaken for dead text.
- In-page tab switch must be instant (<100ms perceived); navigation tabs may load, and should reflect URL changes for deep linking.
- Accessibility: WAI-ARIA tabs pattern — `role="tablist"`, arrow-key traversal, Enter/Space select, visible focus ring (see [06-aria-widget-reference.md](06-aria-widget-reference.md)).
- Don't mix in-page and navigation tabs in one control; if both exist on a page, style them differently.

**Related:** [#top-navigation-bar](#top-navigation-bar), [04-interaction-design.md](04-interaction-design.md) (accordions, progressive disclosure), [05-accessibility.md](05-accessibility.md).

**Sources:** [NN/g: Tabs, Used Right](https://www.nngroup.com/articles/tabs-used-right/), [NN/g: Accordions on Desktop](https://www.nngroup.com/articles/accordions-on-desktop/), [W3C WAI-ARIA APG: Tabs Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/tabs/).

### Hamburger menu

**Definition:** Primary navigation hidden behind a three-line icon (or "Menu" button), revealed as a drawer or overlay on demand. The canonical hidden-navigation pattern.

**Reasoning / Evidence:** NN/g's quantitative study (179 participants, 6 sites, desktop + mobile) measured the discoverability cost precisely: on desktop, hidden navigation was used in 27% of tasks vs 48–50% for visible/combo — usage nearly halved; task time was ≥39% slower and content discoverability dropped >20%. On mobile the penalty exists but is smaller: 57% (hidden) vs 86% (combo) usage; 15% slower. Causes: low salience of a small icon, near-zero information scent (the icon says nothing about what's inside), and extra interaction cost to peek. These three are structural and won't fade with familiarity.

**When to use:**
- Mobile screens where >4–5 top-level destinations must be accommodated and bottom navigation can't cover them.
- As secondary/overflow navigation alongside visible primary items ("combo" navigation — the best-measured mobile pattern: 86% usage).
- Rarely-needed global items (legal, settings, account switching) on constrained screens.

**When NOT to use / exceptions:**
- Desktop: essentially never for primary navigation — screen space exists; NN/g recommendation is unambiguous.
- Mobile apps/sites with ≤4–5 key destinations: use visible bottom navigation or top links instead.
- Anything revenue-critical (top categories of an e-commerce site) should not live only behind the icon.

**Pros:**
- Maximizes content space; accommodates unlimited items.
- Universally recognized icon on mobile (recognition, not comprehension of contents).

**Cons:**
- Halves navigation discoverability on desktop; measurable drops in success, speed, and perceived ease everywhere.
- Out of sight, out of mind: features hidden inside get less usage overall.
- Adds a tap + scan to every navigation act.

**Implementation guidance:**
- If used: label the icon with the word "Menu" (icon+text outperforms icon alone); place top-left or top-right consistently; keep the target ≥44×44px (48dp Android / 44pt iOS).
- Support the hamburger with in-page links to key content so primary tasks don't depend on it (NN/g "supporting the hamburger").
- Drawer contents: plain scannable list, current location highlighted, ≤2 levels of accordion nesting, visible Close affordance + tap-outside + Escape.
- On responsive sites, collapse to hamburger only below the breakpoint where full nav genuinely doesn't fit — don't hide nav at desktop widths for aesthetics.

**Related:** [#bottom-navigation-mobile](#bottom-navigation-mobile), [#top-navigation-bar](#top-navigation-bar), [15-platform-mobile.md](15-platform-mobile.md).

**Sources:** [NN/g: Hamburger Menus and Hidden Navigation Hurt UX Metrics](https://www.nngroup.com/articles/hamburger-menus/), [NN/g: Supporting Mobile Navigation in Spite of a Hamburger Menu](https://www.nngroup.com/articles/support-mobile-navigation/), [NN/g: Hidden-Navigation Study Methodology](https://www.nngroup.com/articles/hidden-navigation-methodology/).

### Mega menus

**Definition:** Large, rectangular dropdown panels attached to top-nav items, showing all (or most) second/third-tier options at once, organized into labeled groups in a 2D layout, optionally with icons or imagery.

**Reasoning / Evidence:** NN/g's usability testing concluded mega menus "work well" and outperform regular cascading dropdowns for large IAs: everything is visible at a glance (no scroll, no cascading hover chains), 2D grouping with headers communicates structure and improves scanning, and richer labels/imagery increase scent. Regular dropdowns punish users with sequential hover paths and hide relationships between options; mega menus replace navigation depth with structured breadth.

**When to use:**
- Sites with many second-level options per top category (e-commerce, universities, government, large B2B) — roughly >10 subitems under a top item.
- When showing tier-2 AND selected tier-3 items at once helps users jump deep directly.

**When NOT to use / exceptions:**
- Few subitems (<6 per category): a simple dropdown or direct landing pages are lighter.
- Touch-first contexts: hover-opened panels don't translate; use tap-to-open full-screen category sheets on mobile.
- Don't dump everything: mega menus curate the top of each branch; a 60-link wall without grouping is noise, not scent.

**Pros:**
- All options visible simultaneously; no memory or cascade cost.
- Grouping headers teach the IA; supports icons/product imagery for scent.
- Reduces required page-depth clicks (jump from home to tier 3).

**Cons:**
- Hover-trigger ergonomics: accidental opens on mouse-travel, accidental closes when cursor clips a gap ([#multi-level-menu-ergonomics-hover-intent-and-safe-triangles](#multi-level-menu-ergonomics-hover-intent-and-safe-triangles)).
- Covers page content; disorienting if it opens instantly on incidental hover.
- Heavy panels tempt marketing to insert promos that compete with navigation scent.
- Accessibility effort is significant (focus management, Escape, screen-reader group labeling).

**Implementation guidance:**
- Open on click, or on hover with an intent delay of ~300ms (and close delay ~400–500ms); never on 0ms hover.
- Group items under clear headers (headers themselves may be links to category landing pages); 3–5 columns max; keep panel within viewport (no internal scrolling).
- Indicate which top items open panels (chevron) vs navigate directly; keep behavior consistent across the bar.
- Keyboard: Enter/Space opens, arrows traverse, Escape closes and returns focus to the trigger; expose groups as labeled lists to assistive tech.
- Keep total links per panel ≤~40 and prioritize; link "View all [category]" for the long tail.

**Related:** [#top-navigation-bar](#top-navigation-bar), [#multi-level-menu-ergonomics-hover-intent-and-safe-triangles](#multi-level-menu-ergonomics-hover-intent-and-safe-triangles), [#hierarchy-depth-vs-breadth](#hierarchy-depth-vs-breadth).

**Sources:** [NN/g: Mega Menus Work Well for Site Navigation](https://www.nngroup.com/articles/mega-menus-work-well/), [NN/g: Mega Menus Gone Wrong](https://www.nngroup.com/articles/mega-menus-gone-wrong/), [W3C WAI-ARIA APG: Menu Button / Disclosure Navigation](https://www.w3.org/WAI/ARIA/apg/patterns/disclosure/examples/disclosure-navigation/).

### Bottom navigation (mobile)

**Definition:** A persistent bar of 3–5 top-level destinations at the bottom of a mobile screen, each with icon + label; the mobile-app equivalent of visible global navigation. iOS: tab bar; Android/Material: navigation bar.

**Reasoning / Evidence:** Puts primary destinations in the thumb's natural reach zone (bottom third of the screen — Hoober's research on one-handed grip: ~49% of users one-hand their phone) and keeps navigation permanently visible, avoiding the hamburger's discoverability penalty. Both Apple HIG and Material Design specify it as the default top-level navigation for apps, which makes it the learned convention.

**When to use:**
- Mobile apps (and app-like mobile web) with 3–5 top-level destinations of comparable importance, frequently switched between.
- When persistent one-tap access to sections measurably matters (social, commerce, media apps).

**When NOT to use / exceptions:**
- More than 5 destinations: don't add a "More" tab as a junk drawer if avoidable — restructure; a More tab is a mini-hamburger with the same scent problems.
- Single-flow or content-immersive experiences (readers, cameras, games) where chrome should recede.
- Desktop/tablet-wide layouts: convert to sidebar/rail (Material 3 explicitly maps bottom bar → navigation rail at ≥600dp).

**Pros:**
- Always visible; thumb-reachable; one tap to any section.
- Icon + label pairs give strong scent and support spatial memory.

**Cons:**
- Hard cap of 5 items; consumes ~56–83px of vertical space permanently.
- Tempts teams to promote marketing destinations into scarce slots.
- Behavior differences across platforms (iOS tab = stateful section stacks; Android historically varied) need deliberate spec.

**Implementation guidance:**
- 3–5 destinations; always show text labels (icon-only bars test poorly for anything but universally understood icons); 10–12sp/pt labels.
- Touch targets ≥48×48dp; active item indicated by color + filled icon + label emphasis (not color alone).
- Preserve per-tab state (scroll position, sub-screen) when switching; tapping the active tab scrolls to top or pops to section root.
- Don't hide the bar on scroll-down in navigation-heavy apps; if hiding for immersion, reveal on scroll-up.
- Never combine bottom nav + hamburger for two halves of primary navigation — one scheme must be primary.

**Related:** [#hamburger-menu](#hamburger-menu), [15-platform-mobile.md](15-platform-mobile.md), [#hub-and-spoke](#hub-and-spoke).

**Sources:** [Material Design 3: Navigation bar](https://m3.material.io/components/navigation-bar/guidelines), [Apple HIG: Tab bars](https://developer.apple.com/design/human-interface-guidelines/tab-bars), [NN/g: Basic Patterns for Mobile Navigation](https://www.nngroup.com/articles/mobile-navigation-patterns/).

### Hub-and-spoke

**Definition:** A central hub screen (dashboard, home, launcher) from which users go out to one section (spoke), complete a task, and return to the hub to start the next; spokes are weakly or not at all connected to each other.

**Reasoning / Evidence:** Matches workflows composed of discrete, independent tasks: the hub provides a single stable orientation point, minimizing the navigation model the user must learn (learn one screen, not a lattice). Classic in settings apps (iOS Settings), banking (account list → account detail → back), smart-home and device dashboards, and setup consoles. The cost is lateral movement: going spoke→spoke requires returning through the hub (2 steps instead of 1).

**When to use:**
- Independent task silos with little cross-navigation (settings, device management, one-account-at-a-time review).
- Novice-heavy or infrequent-use products where a single memorable home beats a rich but complex nav scheme.
- Mobile, where screen constraints reward one-thing-at-a-time flows.

**When NOT to use / exceptions:**
- Frequent lateral movement between sections (compare, cross-reference): the hub round-trip becomes a tax — add bottom nav/sidebar instead.
- Deep spokes (>3 levels) leave users far from the hub; add breadcrumbs or shortcuts.

**Pros:**
- Minimal navigation chrome; maximal focus within each spoke.
- Simple, teachable mental model; the hub can summarize status across spokes (dashboard value).

**Cons:**
- Spoke-to-spoke costs 2+ steps; hub becomes a bottleneck for power users.
- Hubs accrete clutter as sections multiply; requires curation.

**Implementation guidance:**
- Keep hub items scannable: 4–12 entries with icon + label + status glance (badge/summary), grouped if >8.
- Guarantee a one-tap return to hub from any spoke depth (persistent home affordance / logo / up-navigation).
- Add cross-spoke shortcuts only for measured frequent pairs (analytics-driven), not speculatively.
- On larger screens, consider hub-and-spoke degrading into list-detail (hub as persistent left pane).

**Related:** [#bottom-navigation-mobile](#bottom-navigation-mobile), [15-platform-mobile.md](15-platform-mobile.md), [#sequential--wizard-navigation](#sequential--wizard-navigation).

**Sources:** [NN/g: Basic Patterns for Mobile Navigation](https://www.nngroup.com/articles/mobile-navigation-patterns/), Tidwell, Brewer & Valencia, *Designing Interfaces* (3rd ed., O'Reilly, 2020) — Hub and Spoke pattern.

### Sequential / wizard navigation

**Definition:** A linear, step-by-step structure that guides users through an ordered process (checkout, onboarding, tax filing) with explicit Next/Back controls and a step indicator; free navigation is restricted to keep the sequence valid.

**Reasoning / Evidence:** Chunking a long process into steps reduces cognitive load per screen and enables validation at each stage; a visible step indicator sets expectations about length and progress (people abandon less when the end is visible — goal-gradient effect). Trade-off: wizards impose the system's order on the user, which frustrates experts and breaks down when users need to reference or edit earlier answers.

**When to use:**
- Infrequent, high-stakes, or genuinely order-dependent processes: checkout, account setup, filing, configuration with dependencies between steps.
- Novice users performing an unfamiliar task once.

**When NOT to use / exceptions:**
- Frequent expert tasks: a single form or dashboard is faster; wizards enrage repeat users.
- Steps with no real dependencies: let users complete sections in any order (checklist pattern) instead.
- ≤3 trivial fields: one screen, no wizard.

**Pros:**
- Low per-step load; early validation; measurable per-step drop-off for funnel optimization.
- Restricting navigation prevents invalid states.

**Cons:**
- Slower for experts; total time typically longer than a single well-designed form.
- Users can't preview effort or answers ahead; editing earlier steps risks losing later work if badly built.

**Implementation guidance:**
- 3–7 steps; label each step with a noun/verb phrase in the indicator; show "Step 2 of 5" numerically as well.
- Always allow Back without data loss; persist all entered data across steps and sessions (auto-save).
- Allow direct jump to completed steps via the indicator; lock only genuinely dependent future steps.
- Put a review/summary step before final commit for consequential flows; confirmation page after.
- One topic per step; if a step needs a scrollbar of dense fields, split it.

**Related:** [07-forms-and-input.md](07-forms-and-input.md), [#hub-and-spoke](#hub-and-spoke), [04-interaction-design.md](04-interaction-design.md).

**Sources:** [NN/g: Wizards — Definition and Design Recommendations](https://www.nngroup.com/articles/wizards/), [USWDS: Step indicator](https://designsystem.digital.gov/components/step-indicator/), [GOV.UK Design System: Question pages / one thing per page](https://design-system.service.gov.uk/patterns/question-pages/).

### Command palette

**Definition:** A keyboard-summoned overlay (typically Ctrl/Cmd+K or Ctrl/Cmd+Shift+P) that provides fuzzy search over commands, navigation destinations, and objects: typing filters the list, Enter executes or navigates. Popularized by code editors (Sublime Text, VS Code) and now standard in expert web tools (Slack, Linear, GitHub, Figma).

**Reasoning / Evidence:** A command palette flattens an arbitrarily deep IA into constant-time access: instead of recalling where a command lives in menus or a sidebar, the user recognizes it from a filtered list — recall converted to recognition, plus near-zero pointing cost for keyboard-centric users. It scales to hundreds of commands without adding chrome. Its structural weakness is discoverability: an invisible, shortcut-summoned surface has effectively zero information scent for users who don't know it exists, so it must supplement — never replace — visible navigation.

**When to use:**
- Expert productivity tools with large command/destination sets used frequently (editors, project trackers, consoles, design tools).
- Keyboard-centric workflows where round-trips to menus or the pointer are measurable friction.
- As a universal "jump to anything" layer over a deep object space (projects, documents, people).

**When NOT to use / exceptions:**
- Never as the only path to critical or novice-facing tasks; every palette command needs a visible alternative (menu item, button, nav link).
- Content-consumption sites and infrequent-use products: users won't invest in learning the invocation.
- Don't ship a palette as a substitute for fixing a broken IA — it rescues experts while novices stay lost.

**Pros:**
- Constant-time access regardless of IA depth; no permanent screen-space cost.
- Teaches command names and shortcuts (palettes that display the shortcut beside each command train users into faster paths).
- One learned interaction covers navigation, actions, and search.

**Cons:**
- Invisible to the uninitiated; adoption requires deliberate promotion (hints, onboarding, tooltips).
- Requires quality fuzzy matching, indexing, and ranking; a palette with poor matching is worse than menus.
- Mixed result types (commands vs pages vs objects) confuse without clear grouping.

**Implementation guidance:**
- Make the shortcut discoverable: show it in the UI (a visible search/command affordance labeled with the shortcut), in tooltips, and in the help menu.
- Group results by type (Commands / Pages / People / Files…), with type labels; rank by fuzzy-match quality blended with recency/frequency of use.
- Respect permissions: never list commands or destinations the current user cannot execute or access.
- Destructive commands invoked from the palette require the same confirmation as their visible counterparts.
- Support keyboard fully: ↓/↑ traverse, Enter executes, Escape dismisses and returns focus; show each command's dedicated shortcut inline where one exists.
- See [14-platform-desktop.md](14-platform-desktop.md#command-palettes) for platform-level implementation detail.

**Decision questions:**
- Do enough users work keyboard-first, frequently enough, to amortize the learning cost?
- Does every palette action have a visible alternative path?
- Is the command/object index complete and permission-aware?

**Related:** [#scoped-search](#scoped-search), [#sidebar--left-navigation](#sidebar--left-navigation), [14-platform-desktop.md](14-platform-desktop.md#command-palettes), [16-platform-cli-tui.md](16-platform-cli-tui.md).

**Sources:** Pattern conventions documented in product design guidance for VS Code, Slack, Linear, and GitHub command palettes; see [14-platform-desktop.md](14-platform-desktop.md#command-palettes).

---

## Secondary and supporting navigation

### Breadcrumbs

**Definition:** A trail of links at the top of the page showing the current page's ancestors in the site hierarchy back to home (location-based), with items separated by ">" or "/". Path-based breadcrumbs (showing session history) are a distinct — and rejected — variant.

**Reasoning / Evidence:** NN/g has recommended breadcrumbs since 1995: they support wayfinding at almost zero UI cost, and are most valuable for deep-linked arrivals (search engines, shared links) who land mid-hierarchy with no context — the breadcrumb tells them where they are and offers one-click access to broader parent categories. Location-based always beats path-based: path breadcrumbs duplicate the Back button, grow long and repetitive, and provide no value to external arrivals. Google also parses hierarchical breadcrumbs (structured data) for result display — an SEO side benefit. Breadcrumbs must reflect user-meaningful hierarchy, not technical/backend hierarchy.

**When to use:**
- Hierarchies ≥3 levels deep, especially with significant external/search entry traffic (e-commerce, docs, government, large content sites).
- As a supplement to global/local navigation — never a replacement.

**When NOT to use / exceptions:**
- Flat sites (1–2 levels): no wayfinding value; indicate current section in the main nav instead.
- Linear experiences (wizards, articles in a designed sequence): use step indicators or prev/next.
- Mobile: full trails wrap and crowd; truncate to parent-only or "‹ Parent" if space-constrained.

**Pros:**
- Tiny footprint; big orientation payoff; reduces back-button pogo-sticking; lets users "hop up" the tree in one click.

**Cons:**
- Ambiguous on polyhierarchical sites (item with multiple parents) — must pick one canonical path.
- Wasted space on flat sites; wrap/tap-target problems on small screens.

**Implementation guidance:**
- Position: top of content, below global nav. Start with Home (unless global nav already links home — don't duplicate); separate with ">" (or "/").
- Include the current page as the last item, unlinked, visually distinct from linked ancestors.
- Every node must be a real page ("click-and-go") — omit logical-only categories with no landing page.
- One canonical trail per page on polyhierarchical sites (align with canonical URL); never personalize the trail to the user's traversal.
- Mobile: don't wrap to multiple lines; show only the immediate parent if needed; keep tap targets ≥1×1cm.
- Mark up with `nav aria-label="Breadcrumb"` and Schema.org BreadcrumbList.

**Related:** [#wayfinding-and-you-are-here-signals](#wayfinding-and-you-are-here-signals), [#deep-linking-and-url-design](#deep-linking-and-url-design), [05-accessibility.md](05-accessibility.md).

**Sources:** [NN/g: Breadcrumbs — 11 Design Guidelines](https://www.nngroup.com/articles/breadcrumbs/), [Google Search Central: Breadcrumb structured data](https://developers.google.com/search/docs/appearance/structured-data/breadcrumb), [W3C WAI-ARIA APG: Breadcrumb pattern](https://www.w3.org/WAI/ARIA/apg/patterns/breadcrumb/).

### Faceted navigation and filtering

**Definition:** Navigation by successive refinement along independent attributes (facets): price, brand, size, color, date, status. Users compose a query by selecting facet values; results narrow accordingly. Distinct from a single-taxonomy hierarchy because facets combine freely.

**Reasoning / Evidence:** Marti Hearst's Flamenco project and subsequent e-commerce research established faceted search as the dominant pattern for large heterogeneous collections: users rarely know the exact item but can describe attributes. Showing result counts per facet value prevents dead ends (zero-result combinations) and provides scent. The batch-vs-instant apply question is device- and cost-dependent: instant apply gives immediate feedback but on slow connections/mobile causes jarring reflows mid-selection; batch apply lets users compose multi-facet queries but hides feedback until commit — Baymard's e-commerce testing recommends instant apply on desktop and batch ("Apply"/"View N results" button) inside mobile filter overlays.

**When to use:**
- Collections >~50 items with multiple meaningful attributes (products, listings, candidates, logs, tickets).
- When users' queries are attribute-shaped ("waterproof, under $100, size 42") rather than name-shaped.

**When NOT to use / exceptions:**
- Small collections (<~30 items): sorting alone suffices; facet UI is overhead.
- Attributes with inconsistent/missing data: a facet with 70% "unknown" misleads more than it helps — fix data first.
- Don't expose 20 facets; prioritize the 5–8 that drive real decisions (analytics + interviews).
- Filters narrow a set; they must not silently change the meaning or ranking of a search — keep filtering and relevance logic distinct.

**Pros:**
- Handles huge collections without deep hierarchy; supports partial-knowledge queries.
- Result counts teach the collection's shape; selections are reversible one at a time.

**Cons:**
- UI-dense; mobile requires an overlay/drawer that hides results while filtering.
- Engineering cost: counts, zero-result prevention, URL serialization, performance.
- Multi-select semantics (OR within facet, AND across facets) confuse if inconsistent.

**Implementation guidance:**
- Desktop: facets in left sidebar; instant apply with results updating <500ms; keep applied filters visible as removable chips above results with "Clear all."
- Mobile: filter button with active-count badge opens full-screen sheet; batch apply via sticky "Show 143 results" button (count updates live as facets are toggled); active filters must remain visible after the sheet closes.
- Show counts per value; hide or disable (with count 0) values that would empty the result set — prefer hiding for long lists, disabling for short stable lists.
- Within a facet, multi-select = OR; across facets = AND; label ranges clearly (e.g., "$50–$100").
- Collapse facets beyond the top 4–5 by default; show top ~7 values per facet with "Show more."
- Preserve filter state when users drill into a result and come back — losing composed filters on back-navigation is a top complaint.
- Design two distinct empty states: "no data exists at all" (onboarding/first-use framing, offer creation) vs "over-filtered" (show which filters caused it, one-click removal) — see [04-interaction-design.md](04-interaction-design.md#empty-states).
- Serialize the full filter state into the URL for shareability and back-button correctness ([#deep-linking-and-url-design](#deep-linking-and-url-design)).

**Verification checklist:**
- [ ] Active filters visible at all times as removable chips, with individual remove and "Clear all."
- [ ] Filter state survives detail navigation and Back.
- [ ] Counts shown where accurate; never misleading (stale or approximate counts labeled).
- [ ] Zero-result combinations prevented (hidden/disabled values) or recoverable ([#no-results-recovery](#no-results-recovery)).
- [ ] Accessibility labels describe scope and active constraints.

**Related:** [#search-vs-browse](#search-vs-browse), [#sorting](#sorting), [12-data-tables-dashboards.md](12-data-tables-dashboards.md), [#no-results-recovery](#no-results-recovery).

**Sources:** [NN/g: Filters vs. Facets](https://www.nngroup.com/articles/filters-vs-facets/), Hearst, *Search User Interfaces* (Cambridge, 2009), [Baymard Institute: Filtering UX research](https://baymard.com/blog/how-to-design-applied-filters), [NN/g: Applied Filters Should Be Removable](https://www.nngroup.com/articles/applied-filters/).

### Sorting

**Definition:** User control over the order of a list or result set by a chosen attribute (price, date, name, relevance, popularity) and direction (ascending/descending). Sorting reorders the whole set; it never removes items — conflating sort with filter confuses both.

**Reasoning / Evidence:** Sort order is the user's primary lever for surfacing "best for me" from a large set, and the current sort is part of the user's orientation state: if the order changes unpredictably (unstable sorts, silent re-ranking), users lose items they had visually anchored. "Relevance" and "Recommended" defaults are trust-sensitive: when hidden business logic (sponsorship, margin) shapes a default order without labeling, users who discover it recalibrate trust downward — see dark-pattern discussion in [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**When to use:**
- Any list or result set longer than one screen where multiple orderings are legitimately useful (price vs newest vs rating).
- Paired with filtering on browse/search results: filter narrows, sort ranks the remainder.

**When NOT to use / exceptions:**
- Tiny fixed lists (<~10 items) with one natural order: a sort control is noise.
- Don't use sort options to hide relevance logic — if the default is algorithmic, label it ("Recommended") rather than presenting it as neutral.
- Don't offer sorts on attributes with mostly missing data (sorting by a field that is 60% empty produces a meaningless order).

**Pros:**
- Cheap, familiar control; one dropdown serves orderings that would otherwise need separate views.
- Complements filters: together they compose precise "narrow then rank" queries.

**Cons:**
- Each additional sort option adds a decision; long sort menus (>~7 options) signal unclear priorities.
- Server-side sorting across paginated sets must be engineered; client-side sorting of a partial page misleads.

**Implementation guidance:**
- Make the current sort and its direction always visible (labeled control: "Sort: Price, low → high"), including the default ("Sorted by: Relevance") when the order isn't self-evident.
- Use stable sorting: items with equal keys keep their relative order across re-sorts, so the list doesn't shuffle unpredictably.
- Choose a meaningful, stated default per context: recency for activity streams, relevance for search, name for directories.
- Sort must apply to the entire result set (server-side for paginated data), not just the visible page.
- Serialize sort state (with filter state) into the URL so sorted views are shareable and Back-safe (`?sort=price-asc`) — [#deep-linking-and-url-design](#deep-linking-and-url-design).
- Keep sort and filter controls adjacent but visually distinct; a control that both narrows and reorders is unpredictable.
- For column-header sorting in data tables (click-to-cycle, secondary sort, indicators), see [12-data-tables-dashboards.md](12-data-tables-dashboards.md).

**Related:** [#faceted-navigation-and-filtering](#faceted-navigation-and-filtering), [12-data-tables-dashboards.md](12-data-tables-dashboards.md), [#deep-linking-and-url-design](#deep-linking-and-url-design), [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

**Sources:** Hearst, *Search User Interfaces* (Cambridge, 2009); [NN/g: Filters vs. Facets](https://www.nngroup.com/articles/filters-vs-facets/); e-commerce sorting conventions consolidated in [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md).

### Wayfinding and you-are-here signals

**Definition:** The set of cues that tell users where they are, where they can go, and how to get back: highlighted current nav item, breadcrumbs, page titles matching link labels, consistent chrome placement, and distinct visited-link styling.

**Reasoning / Evidence:** NN/g's "Navigation: You Are Here" applies physical wayfinding (Kevin Lynch's *The Image of the City*; mall you-are-here maps) to interfaces: users build cognitive maps from landmarks and paths; missing location signals force them to reconstruct context from content, which fails for deep-linked arrivals (frequently 50%+ of sessions arrive from search). Consistency of placement is itself a wayfinding device: chrome that moves between pages destroys the spatial memory users rely on.

**When to use:**
- Everywhere, always. Every page/screen must answer: What site/app is this? Where am I in it? What's here? Where can I go?

**When NOT to use / exceptions:**
- Deliberately minimal immersive states (fullscreen video, focus modes) may suppress chrome — but must restore orientation instantly on demand.

**Pros:**
- Reduces disorientation, back-button thrashing, and premature exits; supports deep-link traffic.

**Cons:**
- None material; costs are trivial. The only failure mode is inconsistency between signals (nav highlight says one section, breadcrumb another).

**Implementation guidance:**
- Highlight the current item in every visible nav level (global + local) with ≥2 cues (fill + weight); the current page label in nav should match the page's H1 and its `<title>`.
- Keep global chrome position, order, and labels identical across the whole product; place logo top-left, linked to home.
- Avoid identical page titles for different scopes or objects — the title should disambiguate which project/account/record the user is viewing ([#scope-and-context-safety](#scope-and-context-safety)).
- Preserve default visited-link coloring in content-heavy link lists — it is the user's own trail marker.
- Test with the "deep-link teleport" heuristic: land on any interior page cold; can a user state in 5 seconds what site, what section, what page?

**Related:** [#breadcrumbs](#breadcrumbs), [#scope-and-context-safety](#scope-and-context-safety), [03-visual-design.md](03-visual-design.md) (highlighting), [01-core-principles.md](01-core-principles.md) (consistency, visibility of status).

**Sources:** [NN/g: Navigation — You Are Here](https://www.nngroup.com/articles/navigation-you-are-here/), Lynch, *The Image of the City* (MIT Press, 1960), [NN/g: Visibility of System Status](https://www.nngroup.com/articles/visibility-system-status/).

### Scope and context safety

**Definition:** Making the operating scope — environment (production/staging/dev), tenant, organization, workspace, account, project, or file — persistently visible so users always know *where* an action will apply, not just *what page* they are on. Wayfinding answers "where am I in the IA"; scope safety answers "which data will this change touch."

**Reasoning / Evidence:** The canonical failure is a user editing production data while believing they are in a test environment: the screens are pixel-identical across scopes, the URL differs by one token, and nothing in the chrome disambiguates. The same class of error appears in multi-tenant admin tools (acting on the wrong customer), multi-account consoles (deploying to the wrong account), and multi-workspace apps (posting in the wrong workspace). These are high-severity, low-frequency errors: users don't develop vigilance for them precisely because they are rare, so the interface must carry the burden. This is visibility of system status ([01-core-principles.md](01-core-principles.md)) applied to the data dimension, and it is the navigation-side complement of destructive-action confirmation.

**When to use:**
- Any product where the same UI operates over multiple scopes with different consequences: environments, tenants, organizations, workspaces, accounts, projects, branches.
- Especially where mistakes are costly or irreversible: admin consoles, infrastructure tools, CMS/publishing, financial operations.

**When NOT to use / exceptions:**
- Single-scope products: a permanent scope indicator is noise; add it only when a second scope appears.
- Don't rely on theming alone (color-blind users, muscle memory); pair color with a text label.

**Pros:**
- Prevents a class of catastrophic, hard-to-undo errors that no amount of confirmation-dialog copy fully catches.
- Persistent scope labels double as orientation for users juggling several tenants/projects in multiple tabs.

**Cons:**
- Permanent chrome cost (a banner or badge); teams resist "wasting" the pixels until the first production incident.
- Environment theming must be maintained across every surface (including emails, CLIs, API consoles) or the inconsistency itself misleads.

**Implementation guidance:**
- Show the current scope persistently in the global chrome (header badge or sidebar label): workspace name, tenant, environment, account — whichever scopes carry consequence.
- Theme destructive-consequence environments distinctly: e.g., a red/amber banner or border for production, neutral for staging; keep the mapping consistent product-wide and pair color with text ("PRODUCTION").
- Name the scope inside destructive confirmations: "Delete database `orders-primary` in **Production / Acme Corp**?" — never a bare "Are you sure?"; require typed confirmation for the highest-stakes scopes.
- Differentiate page `<title>`s per scope so browser tabs are distinguishable when the same screen is open for two tenants.
- Make the scope switcher explicit and deliberate (a labeled control, not an easily mis-clicked dropdown next to unrelated controls); confirm or visibly announce scope changes.
- Reflect scope in the URL (subdomain or path segment) so links shared between colleagues carry the scope with them.
- Audit "same screen, different scope" pairs: any two scopes whose screens are visually identical need at least one always-visible differentiator.

**Decision questions:**
- What is the worst action a user could perform in the wrong scope, and would the current UI have warned them?
- Can a user with two scopes open in two tabs tell them apart in one glance (tab title, chrome, theme)?
- Does every destructive confirmation name the scope it applies to?

**Related:** [#wayfinding-and-you-are-here-signals](#wayfinding-and-you-are-here-signals), [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md) (multi-tenant SaaS patterns), [04-interaction-design.md](04-interaction-design.md) (destructive-action confirmation), [01-core-principles.md](01-core-principles.md).

**Sources:** [NN/g: Visibility of System Status](https://www.nngroup.com/articles/visibility-system-status/); environment-indicator conventions in infrastructure consoles (AWS, GitHub environments, Kubernetes dashboards); see [19-domain-ecommerce-saas.md](19-domain-ecommerce-saas.md).

### Footer navigation

**Definition:** A "fat footer" at page bottom carrying a structured sitemap excerpt, utility links (legal, contact, careers), and secondary wayfinding — the safety net users consult when in-page navigation has failed or when hunting for organizational metadata.

**Reasoning / Evidence:** Eye-tracking and NN/g observation show users deliberately scroll to footers for a predictable set of items (contact, legal, social, sitemap-style links); the footer is a convention-anchored "last resort" location, so burying those items elsewhere breaks expectations. Fat footers also recover value from dead ends: users who reach the bottom of a page are signaling completion or failure — both are good moments to offer paths onward.

**When to use:**
- All content/marketing sites; any product with legal/compliance links, contact routes, and secondary content that shouldn't consume primary nav slots.

**When NOT to use / exceptions:**
- Infinite-scroll pages make footers unreachable — a structural conflict; either paginate/load-more or move footer content elsewhere ([#pagination-vs-infinite-scroll-vs-load-more](#pagination-vs-infinite-scroll-vs-load-more)).
- Application UIs (dashboards) usually have no page bottom; put utility links in a help/account menu.

**Pros:**
- Cheap catch-all; predictable home for low-frequency, high-importance links; SEO-crawlable structure.

**Cons:**
- Genuinely important tasks placed only in the footer will underperform; it is a supplement, not primary navigation.

**Implementation guidance:**
- Group links under 3–5 headed columns; 5–10 links per column; sentence-case labels.
- Standard contents: contact, about, legal/privacy, careers, social, language/region selector, top categories.
- Keep footer identical on every page; place language/region switchers here or in the utility bar (users look in both).

**Related:** [#top-navigation-bar](#top-navigation-bar), [#wayfinding-and-you-are-here-signals](#wayfinding-and-you-are-here-signals), [13-platform-web.md](13-platform-web.md).

**Sources:** [NN/g: Footers 101](https://www.nngroup.com/articles/footers/), [NN/g: Footers Are Underrated](https://www.nngroup.com/videos/footers/).

---

## Search

### Search vs browse

**Definition:** The strategic choice of whether users primarily locate content by typing queries (search-first) or by traversing categories (browse-first), and how much UI prominence each gets.

**Reasoning / Evidence:** Users split by knowledge state: known-item seekers with precise vocabulary search; exploratory users who can't name the target browse (you can't search for what you can't name). NN/g's "Search Is Not Enough": search punishes poor spellers, vocabulary mismatches, and non-native speakers, and reveals nothing about the collection's structure; browse teaches the domain while navigating. Heavy search usage on a site is often a symptom of failed navigation, not a preference. The right posture is both, weighted by content type: huge/heterogeneous/known-name collections (Amazon, docs by error message, log search) lean search-first; curated/exploratory collections lean browse-first.

**When to use search-first (prominent box, search-centric home):**
- Collections too large or too flat to browse (>10⁴ items, user-generated content, logs).
- Users arrive knowing names/identifiers (SKU, error code, person, ticket #).
- Repeat expert users who have learned the vocabulary.

**When NOT to use / exceptions (browse should lead):**
- Users can't name what they want (gifts, inspiration, "what does this product do?").
- Small collections (<a few hundred items) where a good hierarchy answers faster than query formulation.
- Search quality is poor — a prominent bad search is worse than none; fix relevance before promoting it.

**Pros (search):** direct, scales infinitely, captures intent data (query logs are free research).
**Cons (search):** vocabulary-dependent; hides structure; requires real relevance investment (synonyms, typo tolerance, ranking).
**Pros (browse):** teaches structure; zero query formulation; works for the unnameable.
**Cons (browse):** breaks beyond a few thousand items; every added item raises classification cost.

**Implementation guidance:**
- Always provide both; never make search the only path.
- Prominence: persistent search box (not icon-only) top-right/top-center on content sites, ~27+ characters wide; icon-only search measurably reduces usage.
- Invest in the engine: typo tolerance, synonym maps built from query logs, stemming, and reasonable ranking before promoting search prominence.
- Define and state the searchable scope, including whether search covers archived, hidden, or private items — "I searched and it wasn't there" caused by a silently excluded archive is an invisible-content failure users cannot diagnose.
- When relevance is not obvious, results should indicate what matched (highlighted snippet, matched field label: "matched in comments") so users can judge results without opening each one.
- Preserve the query and any applied filters when users open a result and come back.
- Mine no-result and low-click queries monthly; they are a ranked to-do list for content and synonym gaps.

**Verification checklist:**
- [ ] Search input has a visible label or placeholder stating its scope.
- [ ] Loading, error, and no-results states all designed ([#no-results-recovery](#no-results-recovery)).
- [ ] Results show a count where useful; query is preserved on back-navigation.
- [ ] Coverage of archived/hidden/private content is deliberate and communicated.
- [ ] Results explain what matched when relevance is not obvious.

**Related:** [#autocomplete-and-suggestions](#autocomplete-and-suggestions), [#faceted-navigation-and-filtering](#faceted-navigation-and-filtering), [#organization-schemes](#organization-schemes), [#command-palette](#command-palette).

**Sources:** [NN/g: Search Is Not Enough](https://www.nngroup.com/articles/search-not-enough/), [NN/g: Site Search Suggestions](https://www.nngroup.com/articles/site-search-suggestions/), Hearst, *Search User Interfaces* (2009).

### Autocomplete and suggestions

**Definition:** As-you-type assistance in the search box: query suggestions (completed/popular queries), scoped suggestions ("levis in Jeans"), and direct-result previews (products, people, pages) with thumbnails.

**Reasoning / Evidence:** Suggestions convert recall into recognition: they repair spelling, teach the site's vocabulary, and shorten query formulation — especially valuable on mobile where typing is slow and error-prone. NN/g cautions that bad suggestions actively mislead (irrelevant or overly similar suggestions push users into poor queries); suggestion quality and ordering matter more than quantity.

**When to use:**
- Any site search with meaningful query volume; essential on mobile.
- Mixed suggestion types (queries + top entities) when direct hits are common (people search, product search).

**When NOT to use / exceptions:**
- Sparse content or query logs: suggestions built from thin data embarrass; start with content-title matching only.
- Sensitive contexts: don't suggest other users' queries where privacy leaks are possible.

**Pros:**
- Faster query entry (up to several seconds saved per query on mobile), fewer typos and zero-result searches; teaches available vocabulary.

**Cons:**
- Herds users into popular queries (feedback loop); poorly ranked suggestions worsen outcomes; engineering + latency budget required.

**Implementation guidance:**
- Trigger after 2–3 characters; respond <100ms perceived (debounce ~150–300ms); show ≤10 suggestions on desktop, 4–8 on mobile (keyboard occlusion).
- Highlight the completed portion (differentiate typed vs suggested text); support keyboard: ↓/↑ traverse, Enter selects, Escape dismisses; combobox ARIA pattern.
- Order: exact/prefix matches, then popularity-weighted; visually separate query suggestions from direct entity results with thumbnails.
- Never autocorrect silently at submit; show "showing results for X — search instead for Y."

**Related:** [#search-vs-browse](#search-vs-browse), [#no-results-recovery](#no-results-recovery), [05-accessibility.md](05-accessibility.md) (combobox pattern), [06-aria-widget-reference.md](06-aria-widget-reference.md).

**Sources:** [NN/g: Site Search Suggestions](https://www.nngroup.com/articles/site-search-suggestions/), [Baymard: Autocomplete design](https://baymard.com/blog/autocomplete-design), [W3C WAI-ARIA APG: Combobox pattern](https://www.w3.org/WAI/ARIA/apg/patterns/combobox/).

### No-results recovery

**Definition:** The designed response when a query or filter combination returns zero results: explanation, repair suggestions, relaxed alternatives, and paths onward — instead of a dead end.

**Reasoning / Evidence:** NN/g's SERP and empty-state research: a bare "No results found" strands the user and commonly ends the session; the empty state is a high-leverage recovery moment. The message must state clearly that nothing matched (system status), avoid blaming or joking away the failure, and do repair work for the user: spelling correction, partial-match relaxation, and category alternatives measurably rescue sessions.

**When to use:**
- Every search and every filterable list — zero-result states are guaranteed to occur and must be designed, not defaulted.

**When NOT to use / exceptions:**
- Don't silently substitute different results for the query without a visible notice — "we showed you something else" without saying so destroys trust.
- Prevent rather than recover where possible: facet counts and disabled zero-value options stop many dead ends upstream.

**Pros:**
- Converts abandonment into recovery; no-result query logs pinpoint content and synonym gaps.

**Cons:**
- Relaxation logic (dropping terms, fuzzy matching) risks confusing precision-seeking users if not clearly labeled.

**Implementation guidance:**
- State plainly what happened, echoing the query: "No results for 'wirless earbuds'."
- Offer in order: (1) spelling correction "Did you mean wireless?", (2) automatic relaxed results clearly labeled ("Showing results for wireless earbuds"), (3) tips (fewer/more general words), (4) category links or top products, (5) contact/support path — or a create/request path where the user may legitimately need to add the missing item.
- Distinguish the two zero-state causes and respond differently: "no data exists at all" (first-use/empty collection — explain and offer creation) vs "over-filtered" (data exists but the current query/filters exclude it — identify the culprit and offer relaxation); see [04-interaction-design.md](04-interaction-design.md#empty-states).
- For filtered lists: identify the culprit ("No results match Brand: Acme + Color: Red"), offer one-click removal of the most restrictive filter, keep all filters editable in place.
- Log every zero-result query with count; review top 50 monthly; fix via synonyms, content, or merchandising.
- No dead-end humor: puns instead of status ("Whoopsie-doodle!") obscure what happened — see [09-content-ux-writing.md](09-content-ux-writing.md).

**Related:** [#autocomplete-and-suggestions](#autocomplete-and-suggestions), [#faceted-navigation-and-filtering](#faceted-navigation-and-filtering), [04-interaction-design.md](04-interaction-design.md#empty-states), [09-content-ux-writing.md](09-content-ux-writing.md) (empty states, error copy).

**Sources:** [NN/g: No Results Found SERP Design](https://www.nngroup.com/articles/search-no-results-serp/), [NN/g: Designing Empty States in Complex Applications](https://www.nngroup.com/articles/empty-state-interface-design/), [Baymard: No-results page design](https://baymard.com/blog/no-results-page).

### Scoped search

**Definition:** Search constrained to a subset of the collection — by section (Docs vs Community), object type (People, Files, Messages), or current context (search within this project/folder) — via a scope selector or implicit contextual scoping.

**Reasoning / Evidence:** Scoping raises precision when the corpus is heterogeneous (a query like "billing" means different things in Docs vs Support vs Community). But NN/g found users frequently overlook scope selectors and then misinterpret results — searching an unintended silo and concluding content doesn't exist. The scope must be loudly visible at query time and on the results page, with a one-click widen.

**When to use:**
- Heterogeneous corpora with distinct result types users think of separately (mail vs drive vs people; code vs issues vs wiki).
- Contextual object search in apps ("find in this project") where the current container is the obvious default.

**When NOT to use / exceptions:**
- Homogeneous content: scoping adds a decision with no precision gain — one box, good ranking.
- Never default to a narrow scope silently on a global search box; default to all-content with type filters on the results page.

**Pros:**
- Higher precision; smaller index per query; matches expert mental models of the workspace.

**Cons:**
- Users miss the scope control and get invisible-content failures; scope selectors clutter compact headers.

**Implementation guidance:**
- Prefer scope-after-search: search everything, then filter results by type via tabs/facets on the SERP — safer than pre-scoping.
- If pre-scoping: show current scope inside the input as a removable chip/token ("in: Project Alpha ×"); echo scope prominently on results ("42 results in Documentation — search all content instead").
- Be explicit about whether the active scope includes archived, hidden, or restricted items; if permission filtering hides results, prefer a hint ("some results may be hidden by permissions") over silent omission where policy allows.
- When a scoped search returns few/zero results, automatically show hit counts from other scopes ("0 in Docs; 17 in Community").
- Keyboard-driven apps: scope operators (`from:`, `type:`, `in:`) with autocomplete for power users.

**Related:** [#search-vs-browse](#search-vs-browse), [#no-results-recovery](#no-results-recovery), [#command-palette](#command-palette), [16-platform-cli-tui.md](16-platform-cli-tui.md) (query operators).

**Sources:** [NN/g: Scoped Search — Dangerous, but Sometimes Useful](https://www.nngroup.com/articles/scoped-search/), Hearst, *Search User Interfaces* (2009).

---

## Content traversal

### Pagination vs infinite scroll vs load more

**Definition:** Three strategies for presenting result sets larger than one screen: pagination (discrete numbered pages), infinite scroll (auto-append on scroll), and load-more (manual append via button — a hybrid).

**Reasoning / Evidence:** NN/g and Baymard converge on task-fit: infinite scroll suits leisurely, discovery-oriented feeds where any item is equally acceptable and sequence position is meaningless (social feeds, image discovery); it fails goal-directed retrieval because it breaks location memory ("it was on page 3" has no equivalent), makes the footer unreachable, bloats memory, and cripples back-button return-to-position unless carefully engineered. Pagination supports orientation, memorability, jump-to-end, and comparison across a known extent — vital for ordered/searchable sets. Load-more keeps user control and footer reachability while allowing continuous flow, and is Baymard's recommended default for e-commerce product lists (with lazy-loaded imagery).

Decision table:

| Criterion | Pagination | Load more | Infinite scroll |
|---|---|---|---|
| User goal | Find a specific item; compare; resume later | Browse with intent; e-commerce lists | Passive discovery; entertainment feeds |
| Item order matters / recall position | Best | OK | Worst |
| Footer reachable | Yes | Yes | No |
| Back-button / return to position | Trivial | Good (with state) | Hard (needs scroll restoration) |
| Reach end of list / jump far | Trivial (last page) | Poor | Impossible |
| Perceived friction | Highest (click + reload) | Low | None |
| SEO / shareable location | Best (page URLs) | Medium | Poor |
| Typical batch size | 20–100/page | 20–50/click | 10–30/viewport |

**When to use:**
- Pagination: search results, tables, admin lists, ordered archives, anything users cite/resume ("results 40–60").
- Load more: e-commerce category lists, mixed browse/goal contexts, mobile lists generally.
- Infinite scroll: feeds where engagement-time is the goal and items are interchangeable.

**When NOT to use / exceptions:**
- Infinite scroll: never on task-oriented or comparison-driven lists, or any page with a needed footer; provide it only with scroll-position restoration and URL state.
- Pagination: avoid tiny page sizes (<20) that force excessive clicking through a modest set.

**Pros / Cons:** captured in the decision table above.

**Implementation guidance:**
- Pagination: show first/last, current ±2 pages, prev/next; current page non-linked and visually distinct; update URL per page (`?page=3`); keep per-page count user-selectable (25/50/100) on data tables.
- Load more: label with remaining context ("Show 24 more of 312"); after append, move focus to the first new item (screen-reader continuity); update URL/history so Back restores state.
- Infinite scroll: maintain scroll position on back-navigation; virtualize the DOM beyond ~200 items; append a sentinel ~1.5 viewports early for seamlessness; offer a "back to top" affordance after 2+ viewports.
- All three: communicate total count ("312 results") near the top — extent knowledge is orientation.

**Related:** [12-data-tables-dashboards.md](12-data-tables-dashboards.md), [#footer-navigation](#footer-navigation), [#deep-linking-and-url-design](#deep-linking-and-url-design).

**Sources:** [NN/g: Infinite Scrolling Is Not for Every Website](https://www.nngroup.com/articles/infinite-scrolling/), [NN/g: Infinite Scrolling — When to Use It, When to Avoid It](https://www.nngroup.com/articles/infinite-scrolling-tips/), [Baymard: Load More + Lazy Load recommendation](https://baymard.com/blog/pagination-load-more-infinite-scroll).

### Deep linking and URL design

**Definition:** Treating URLs as part of the UX: every meaningful state (page, record, filter set, tab, search) has a stable, readable, shareable address; links into deep states land users in a fully oriented view.

**Reasoning / Evidence:** Jakob Nielsen argued URI design is UI design ("URL as UI", 1999): users read, edit, share, and bookmark URLs; readable hierarchical URLs also serve as supplementary wayfinding and are correlated with search-result click confidence. In SPAs, failure to encode state in URLs breaks the web's core affordances — back button, refresh, share, bookmark — which users experience as data loss and disorientation. Deep-linked arrivals (search/social/chat) are the majority entry mode for many products; a deep link that lands on a redirect-to-home or a context-less fragment defeats the sharer's intent.

**When to use:**
- Everywhere on the web platform; in mobile apps via universal links/app links mapping URLs to in-app screens.
- Encode into the URL: entity IDs + slugs, pagination, sort, filter state, active tab, search query, selected item in list-detail.

**When NOT to use / exceptions:**
- Don't encode private/sensitive tokens or PII in URLs (they leak via logs, referrers, screenshots).
- Ephemeral UI state (open dropdown, scroll offset, hover) stays out of the URL; half-completed multi-step form state usually belongs server-side, not in the address bar.

**Pros:**
- Restores Back/refresh/share/bookmark; enables support ("send me the link"), QA repro, and analytics by state; SEO.

**Cons:**
- Engineering discipline: state serialization, permanence commitments (redirects when structures change), canonical-URL management for polyhierarchy.

**Implementation guidance:**
- Structure mirrors IA: `/category/subcategory/item-slug`; lowercase, hyphen-separated, human-readable; ≤~75 characters where practical.
- Query params for view state: `?page=3&sort=price-asc&color=red&q=boots`; stable param names; alphabetize serialization for cache/canonical consistency.
- Cool URIs don't change: 301 old URLs on restructures; never break inbound links (Tim Berners-Lee, 1998).
- Landing on a deep link must render full context: breadcrumb, highlighted nav, page title — test every template cold ([#wayfinding-and-you-are-here-signals](#wayfinding-and-you-are-here-signals)).
- SPAs: use the History API for every navigable state; Back must never dump the user out of the app or lose a modal's underlying page.

**Related:** [#breadcrumbs](#breadcrumbs), [#pagination-vs-infinite-scroll-vs-load-more](#pagination-vs-infinite-scroll-vs-load-more), [13-platform-web.md](13-platform-web.md).

**Sources:** [NN/g: URL as UI](https://www.nngroup.com/articles/url-as-ui/), [W3C/TimBL: Cool URIs don't change](https://www.w3.org/Provider/Style/URI), [Google: URL structure best practices](https://developers.google.com/search/docs/crawling-indexing/url-structure).

---

## Navigation for data-heavy UIs

In data-heavy enterprise UIs the table *is* the navigation surface: users locate records by composing sort + filter rather than by traversing menus, and table state (sort indicators, filter chips, "X of Y rows") functions as you-are-here signaling for the data space. The full treatment of large data tables — column sorting cycles, per-column and global filtering, column show/hide/reorder/pin, density, saved views, row-level navigation, and state persistence — lives with the rest of the data-dense-UI guidance. Data tables, data grids, and data-dense UI: see [12-data-tables-dashboards.md](12-data-tables-dashboards.md). The general principles in this file still apply there: preserve composed state across navigation ([#deep-linking-and-url-design](#deep-linking-and-url-design)), keep active constraints visible ([#faceted-navigation-and-filtering](#faceted-navigation-and-filtering), [#sorting](#sorting)), and design the over-filtered empty state ([#no-results-recovery](#no-results-recovery)).

---

## Menu ergonomics

### Multi-level menu ergonomics: hover intent and safe triangles

**Definition:** Techniques that make hover-opened hierarchical menus (dropdowns, flyouts, mega menus) tolerant of imprecise pointing: intent delays before open/close, and "safe triangle" geometry that keeps a submenu open while the cursor travels diagonally toward it — even when the diagonal path crosses other parent items.

**Reasoning / Evidence:** Fitts's Law and steering law: moving a cursor along a narrow corridor (across a parent item, then into a flyout) is slow and error-prone; naive implementations close the submenu the instant the cursor grazes an adjacent item, forcing users into unnatural right-angle paths. Amazon's category flyout famously solved this with triangle hit-detection (documented by Ben Kamens, 2013): compute the triangle between cursor and the submenu's top/bottom corners; while the cursor moves within it, suppress menu switching. Apple used a related timing approach in early Mac menus. Hover-intent delays (~150–300ms before open) prevent "menu flashing" as the cursor crosses a bar en route elsewhere.

**When to use:**
- Any hover-triggered submenu or mega menu on pointer-driven interfaces; cascading context menus in desktop apps.

**When NOT to use / exceptions:**
- Touch devices: no hover — every level must be tap-to-open (first tap opens, second navigates, or dedicated chevron targets).
- Prefer click-to-open menus outright where content covers the page (mega menus): click intent is unambiguous and needs no geometry tricks.
- Avoid >2 levels of cascading flyouts entirely; error rates compound per level — restructure the IA instead.

**Pros:**
- Cuts accidental menu switches to near zero; users can move diagonally at natural speed; perceived quality jumps.

**Cons:**
- Overlong delays make menus feel sluggish; triangle logic misfires if submenu geometry changes while open; more code to test.

**Implementation guidance:**
- Open delay ~150–300ms hover-intent; close delay ~300–500ms after pointer exit; cancel timers on re-entry.
- Safe triangle: on parent hover with open submenu, track cursor; if within the triangle (cursor origin → submenu near corners), delay item-switch ~300ms or until movement stalls; libraries: jQuery-menu-aim lineage, Floating UI's `safePolygon`.
- Also open on click as a redundant path (serves touch, trackpad, and motor-impaired users); Enter/Space and arrow keys must fully operate the menu (see [05-accessibility.md](05-accessibility.md)).
- Submenu targets ≥32px row height desktop; leave no dead gap between parent and flyout (visually adjacent or bridged hit area).
- Test at trackpad speeds and with a display scaled to 125–200% — geometry bugs surface there first.

**Related:** [#mega-menus](#mega-menus), [04-interaction-design.md](04-interaction-design.md) (Fitts's Law, hover states), [05-accessibility.md](05-accessibility.md), [11-components-and-overlays.md](11-components-and-overlays.md).

**Sources:** [Ben Kamens: Breaking down Amazon's mega dropdown](https://bjk5.com/post/44698559168/breaking-down-amazons-mega-dropdown), [NN/g: Mega Menus Work Well](https://www.nngroup.com/articles/mega-menus-work-well/), [Floating UI: safePolygon docs](https://floating-ui.com/docs/usehover#safepolygon), Accot & Zhai, "Beyond Fitts' Law: Models for Trajectory-Based HCI Tasks," CHI 1997.

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| IA before nav | Define/validate structure before choosing navigation chrome | Inverted-L caps at ~4 tiers |
| Organization schemes | Topical/task-based for browse; exact schemes for known-item; avoid audience-based | A–Z only for >~25 known-name items |
| Depth vs breadth | Prefer broad-shallow; each level is a compounding error point | ≤3 levels content sites; 5–12 items/level |
| Card sorting | Discover mental models and labels; can't prove findability | 30–60 cards; 15–20 participants |
| Tree testing | Validate findability text-only before UI design | ≥80% success, ≥70% directness targets; 50 users/tree; correct first click ≈2× success |
| Top nav bar | Default for content sites with few short-labeled categories | 3–8 items; sticky ≤64px |
| Sidebar nav | Default for apps; scales, allows long labels, collapsible | 240–280px expanded; 48–72px rail |
| Tabs | One row, short labels, no mixed types, default = high-use tab | 2–6 tabs; 1–2-word labels; ≥2 selection cues |
| Hamburger | Avoid on desktop; mobile only for >4–5 destinations, with in-page support links | Desktop usage 27% hidden vs 48–50% visible; 39% slower |
| Mega menus | Beat cascading dropdowns for large IAs; group under headers; click or delayed hover | ~300ms open delay; 3–5 columns; ≤~40 links |
| Bottom nav | Mobile default for 3–5 peer destinations; always label icons | 3–5 items; ≥48×48dp targets |
| Hub-and-spoke | Independent task silos; one-tap return to hub | 4–12 hub entries |
| Wizard | Order-dependent, infrequent flows; visible progress; free Back | 3–7 steps; persist data across steps |
| Command palette | Expert accelerator over a deep IA; never the only path to a task | Ctrl/Cmd+K convention; group results by type |
| Breadcrumbs | Location-based only, ≥3-level sites; current page last, unlinked | Truncate to parent on mobile; ≥1×1cm taps |
| Faceted filtering | Instant apply desktop; batch "Show N results" in mobile sheets | 5–8 facets; update <500ms; counts on values |
| Sorting | Stable sort; visible attribute + direction; label the default; state in URL | ≤~7 sort options; sort whole set server-side |
| Scope safety | Persistent scope indicator + environment theming; name scope in destructive confirms | 1-glance scope check across tabs |
| Search vs browse | Offer both; search-first only for huge/known-name corpora with good relevance | Visible box ≥27 chars wide |
| Autocomplete | Trigger early, respond fast, never silently rewrite queries | 2–3 chars; <100ms perceived; ≤10 suggestions |
| No results | State it, echo query, repair (spelling → relaxed → categories); over-filtered ≠ no data | Review top 50 zero-result queries monthly |
| Scoped search | Prefer filter-after-search; make active scope loud with 1-click widen | Show other-scope hit counts on 0 results |
| Paging strategy | Pagination = lookup; load-more = shopping; infinite = feeds only | 20–100/page; focus first new item on append |
| URL design | Every meaningful state addressable, readable, permanent | 301 all moves; slugs ≤~75 chars |
| Data tables | See [12-data-tables-dashboards.md](12-data-tables-dashboards.md) for sorting, filtering, column management, saved views | Persist state per user; "X of Y rows" |
| Hover menus | Intent delays + safe-triangle; click as redundant path; ≤2 flyout levels | 150–300ms open; 300–500ms close |
