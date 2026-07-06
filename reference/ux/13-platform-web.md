# Platform: Web

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers UX concerns specific to the web platform: responsive design, browser conventions, perceived performance and Core Web Vitals, URL design, progressive enhancement, SPA vs MPA, PWA, scroll UX, consent and permission prompts, typography and viewport quirks, multi-tab behavior, session expiry, and agent-friendly content. Platform-agnostic interaction rules are in [04-interaction-design.md](04-interaction-design.md); accessibility (largely web-native in vocabulary) in [05-accessibility.md](05-accessibility.md); desktop-class web apps also draw on [14-platform-desktop.md](14-platform-desktop.md); mobile-web overlap in [15-platform-mobile.md](15-platform-mobile.md).

## Contents

- [Layout across devices](#layout-across-devices)
  - [Responsive design](#responsive-design)
  - [Viewport quirks](#viewport-quirks)
- [Respecting the browser](#respecting-the-browser)
  - [Browser conventions](#browser-conventions)
  - [Links vs buttons](#links-vs-buttons)
  - [Multi-tab behavior and state sync](#multi-tab-behavior-and-state-sync)
  - [Session expiry handling](#session-expiry-handling)
- [Performance as UX](#performance-as-ux)
  - [Core Web Vitals](#core-web-vitals)
  - [Perceived performance techniques](#perceived-performance-techniques)
  - [Web typography specifics](#web-typography-specifics)
- [Architecture choices](#architecture-choices)
  - [URLs as UX](#urls-as-ux)
  - [SPA vs MPA](#spa-vs-mpa)
  - [Progressive enhancement](#progressive-enhancement)
  - [PWA capabilities](#pwa-capabilities)
- [Page-level behaviors](#page-level-behaviors)
  - [Scroll UX](#scroll-ux)
  - [Consent and cookie UX](#consent-and-cookie-ux)
  - [Notification permission UX](#notification-permission-ux)
  - [SEO-UX overlap](#seo-ux-overlap)
  - [Agent-friendly web content](#agent-friendly-web-content)
- [Quick reference](#quick-reference)

---

## Layout across devices

### Responsive design

**Definition:** One codebase adapting layout to any viewport via fluid grids, flexible media, and breakpoints — mobile-first (styling the small screen as the base, layering complexity up) as the standard methodology, now augmented by container queries (components responding to their *container*, not the viewport).

**Reasoning / Evidence:** Device diversity is unbounded (320px phones to 5K desktops, foldables, embedded webviews); separate m-dot sites lost to responsive years ago (URL fragmentation, maintenance). Mobile-first forces content priority decisions early ([01-core-principles.md](01-core-principles.md#8-aesthetic-and-minimalist-design)) and matches the majority-mobile traffic reality of most consumer products. WCAG reflow (320px) makes responsiveness a legal floor ([05-accessibility.md](05-accessibility.md#reflow)).

**When to use:**
- Everything on the web. Breakpoints from content (where the layout breaks), not device names — typical working set: ~360 (mobile), ~768 (tablet/narrow), ~1024 (desktop), ~1440 (wide).
- Container queries for reusable components living in variable slots (card in sidebar vs main column).
- Fluid type/spacing (`clamp()`) to reduce breakpoint count.

**When NOT to use / exceptions:**
- Genuinely desktop-only professional tools (complex editors) may set a floor ("best on larger screens") — but chrome/auth/marketing pages still must reflow.
- Don't hide *content* on mobile as a layout escape hatch — reprioritize, don't amputate (mobile users have the same goals).

**Pros:** One URL space, one codebase, legal compliance, future-device tolerance.
**Cons:** Testing matrix; complex data-dense views need real per-breakpoint design (tables → cards decisions), not just reflow.

**Implementation guidance:**
- Base: fluid single column; enhance up (`min-width` queries). Max content width 1200–1440px, text columns capped at 60–75ch ([03-visual-design.md](03-visual-design.md#alignment-and-grids)).
- Tables on small screens: horizontal scroll *with sticky first column* for true tables; card-per-row transformation for record lists; column-priority hiding with user control as third option ([12-data-tables-dashboards.md](12-data-tables-dashboards.md)).
- Images: `srcset`/`sizes`, `aspect-ratio` reserved to prevent CLS.
- Account for coarse pointers and no-hover environments — don't gate functionality on hover alone ([04-interaction-design.md](04-interaction-design.md)).
- Test at: 320, 375, 768, 1024, 1440, plus 200% zoom ([05-accessibility.md](05-accessibility.md#zoom-and-text-resize)), and both landscape/portrait on touch devices.

**Related:** [05-accessibility.md](05-accessibility.md#reflow), [15-platform-mobile.md](15-platform-mobile.md), [03-visual-design.md](03-visual-design.md#alignment-and-grids).

**Sources:** [MDN: Responsive design](https://developer.mozilla.org/en-US/docs/Learn/CSS/CSS_layout/Responsive_Design), [web.dev: Learn Responsive Design](https://web.dev/learn/design/), [MDN: Container queries](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_containment/Container_queries).

### Viewport quirks

**Definition:** The mobile-web viewport's sharp edges: dynamic browser chrome making `100vh` overflow (address bar collapse), notches/home-indicators requiring safe-area insets, keyboard-driven viewport resizing, and pinch-zoom interplay.

**Reasoning / Evidence:** Classic bugs: full-screen heroes cut off by the address bar (100vh > visible height on iOS Safari); fixed bottom bars hidden behind home indicators; inputs obscured by the on-screen keyboard. CSS shipped fixes: dynamic viewport units (`dvh`/`svh`/`lvh`), `env(safe-area-inset-*)`, `interactive-widget` viewport meta.

**When to use:**
- Any full-height layout: `100dvh` (dynamic) or `100svh` (small, conservative) instead of `100vh`.
- Fixed bottom UI: `padding-bottom: env(safe-area-inset-bottom)` + `viewport-fit=cover` when drawing edge-to-edge.
- Forms: ensure focused inputs scroll above the keyboard (test iOS Safari specifically).

**When NOT to use / exceptions:**
- Never `user-scalable=no` / `maximum-scale=1` — zoom-blocking is an accessibility violation ([05-accessibility.md](05-accessibility.md#zoom-and-text-resize)); prevent iOS auto-zoom-on-focus by giving inputs ≥16px font instead.

**Pros:** Correct rendering on the devices most users hold.
**Cons:** Cross-browser unit support nuances; requires on-device testing (simulators lie about chrome behavior).

**Implementation guidance:**
- `<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">`.
- Bottom-fixed CTAs: `bottom: max(16px, env(safe-area-inset-bottom))`.
- Test: iOS Safari (scroll up/down chrome states), Android Chrome, in-app webviews (Instagram/Gmail browsers have their own chrome).

**Related:** [#responsive-design](#responsive-design), [15-platform-mobile.md](15-platform-mobile.md).

**Sources:** [web.dev: viewport units](https://web.dev/blog/viewport-units), [MDN: env()](https://developer.mozilla.org/en-US/docs/Web/CSS/env).

---

## Respecting the browser

### Browser conventions

**Definition:** The browser-provided behaviors users own and expect: working Back/Forward, scroll restoration, find-in-page (Ctrl/Cmd+F), zoom, copy/select, open-in-new-tab (middle-click/Cmd-click), bookmarkable URLs, refresh safety, autofill and password-manager compatibility. Breaking any of them is breaking the user's tools.

**Reasoning / Evidence:** [Jakob's Law](01-core-principles.md#jakobs-law) at maximum strength: Back is the most-used browser action; its breakage (SPA hijacks, redirect traps, state loss) tests as one of the most trust-destroying failures. Find-in-page fails on virtualized/lazy content; copy-blocking and right-click-blocking punish legitimate users while stopping no one. Custom text controls that break IME input, paste, or password managers are recurring failure modes.

**When to use:**
- Back button: every distinct view/state a user perceives as "a place" gets a history entry; modals/lightboxes on mobile-web should close on Back (pushState discipline); never trap with redirect stacks.
- Scroll restoration: browser default for MPA; SPAs must reimplement (save/restore per history entry — [04-interaction-design.md](04-interaction-design.md#scroll-behaviors)).
- Refresh: no data loss (drafts autosaved — [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)); no double-submit (POST-redirect-GET or idempotency keys).
- New-tab: real links (`<a href>`) so middle-click/Cmd-click work everywhere navigation happens.

**When NOT to use / exceptions:**
- Deliberate history flattening is right for wizard *sub*-steps that shouldn't be re-enterable (payment processing) — use `replaceState`, and keep Back meaning "leave the flow safely."
- Games/canvas apps legitimately capture some inputs — with escape hatches.

**Pros:** Trust, muscle-memory speed, shareability, power-user retention.
**Cons:** SPA implementations must hand-build what MPAs get free — real engineering cost (the argument for MPA/hybrid where app-ness isn't needed — [#spa-vs-mpa](#spa-vs-mpa)).

**Implementation guidance:**
- Audit ritual: navigate 5 screens deep, hit Back 5 times — every step must land where expected with state/scroll intact.
- Never `history.pushState` on filter twiddles that users wouldn't call "a place" (use `replaceState`), but DO put filters in the URL ([#urls-as-ux](#urls-as-ux)).
- Don't block: text selection, right-click, paste ([07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)), zoom, find (avoid full virtualization on text-searchable content where feasible).

**Verification checklist:**
- [ ] Back button returns to the previous meaningful view with state and scroll intact — never exits the app unexpectedly or fires an action.
- [ ] Middle-click/Cmd-click opens every navigational control in a new tab.
- [ ] Refresh mid-task loses no user input and never double-submits.
- [ ] Zoom, text selection, copy/paste, context menus, find-in-page, autofill, and password managers all work.
- [ ] Custom inputs handle IME composition and paste correctly.

**Related:** [#spa-vs-mpa](#spa-vs-mpa), [#urls-as-ux](#urls-as-ux), [01-core-principles.md](01-core-principles.md#jakobs-law).

**Sources:** [NN/g: Browser back button](https://www.nngroup.com/articles/recalling-vs-recognizing-browser-back/), [web.dev: History API](https://developer.mozilla.org/en-US/docs/Web/API/History_API).

### Links vs buttons

**Definition:** The semantic contract: **links** (`<a href>`) navigate to a URL; **buttons** (`<button>`) perform actions. Visual style may vary; the underlying element must match the behavior.

**Reasoning / Evidence:** Links carry free behavior users depend on: new-tab open, copy-address, hover-preview in status bar, visited styling, SR "links list" navigation. A `<div onclick>` or a button-that-navigates loses all of it; a link-that-mutates breaks middle-click horribly (action fires? or nothing?). AT users navigate by element type — miscategorized controls scramble their map.

**When to use:**
- Navigation (URL changes, even in-SPA): `<a href>` — styled as a button if design wants, but semantically an anchor.
- Actions (submit, open modal, toggle, delete): `<button type=…>`.

**When NOT to use / exceptions:**
- `href="#"` + onclick is never right; `javascript:` URLs never.
- Ambiguous cases (open detail *panel* that also has a URL): prefer link + progressive enhancement.

**Pros:** Free platform behavior, AT correctness, SEO crawlability.
**Cons:** None; this is a correctness rule.

**Implementation guidance:**
- Lint for click-handlers on divs/spans; design systems expose `as`/`href` props so one visual Button renders either element ([10-design-systems.md](10-design-systems.md#component-api-design)).
- Buttons in forms: explicit `type="button"` unless submitting (default is submit — classic bug).

**Related:** [05-accessibility.md](05-accessibility.md#semantic-html-first), [#browser-conventions](#browser-conventions).

**Sources:** [MDN: `<a>`](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/a), [W3C: links vs buttons (APG)](https://www.w3.org/WAI/ARIA/apg/patterns/button/).

### Multi-tab behavior and state sync

**Definition:** Handling users' natural multi-tab workflows: comparison across tabs, long-lived stale tabs, concurrent edits, auth state changes in one tab affecting others.

**Reasoning / Evidence:** Tabs are user-owned working memory ([02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking)) — comparison shopping, reference-while-writing, queued reading. Failures: cart edited in tab A silently wrong in tab B; logout in one tab leaving others in undead sessions; "your session is invalid" on the tab a user returns to after an hour.

**When to use:**
- Auth sync: login/logout broadcast across tabs (BroadcastChannel/storage events) — other tabs show a gentle "signed out — sign in again" state, never data-destroying redirects.
- Shared-mutable data (carts, documents): sync on visibilitychange/focus; conflict strategy for concurrent edits (last-write warning or real merge).
- Stale-tab refresh: on tab focus after inactivity, silently revalidate data; show "updated" indicators rather than jarring reloads.

**When NOT to use / exceptions:**
- Don't force single-tab ("app already open in another tab" blocks) except genuinely stateful exclusive sessions (some trading/telephony) — and then explain, offer takeover.

**Pros:** Matches real behavior; prevents silent data conflicts.
**Cons:** Sync engineering; edge cases multiply (mitigate: revalidate-on-focus covers most value cheaply).

**Implementation guidance:**
- `document.visibilitychange` → revalidate stale queries (stale-while-revalidate libraries do this by default).
- BroadcastChannel for auth events; drafts keyed per-entity so two tabs editing different records never collide.
- Preserve deep state per tab (each tab = its own history/scroll) — don't globalize UI state that should be per-view.

**Verification checklist:**
- [ ] Open the app in two tabs, log out in one — the other degrades gracefully (clear signed-out state, no data-destroying redirect).
- [ ] Edit shared data (cart, document) in one tab — the other reflects it on focus, or warns on conflicting save.
- [ ] Return to a tab left idle for an hour — data revalidates without a jarring reload or state loss.

**Related:** [#session-expiry-handling](#session-expiry-handling), [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery).

**Sources:** [MDN: BroadcastChannel](https://developer.mozilla.org/en-US/docs/Web/API/Broadcast_Channel_API), [MDN: Page Visibility API](https://developer.mozilla.org/en-US/docs/Web/API/Page_Visibility_API).

### Session expiry handling

**Definition:** What happens when auth/session lifetime ends while the user has state in flight: warnings before expiry, state preservation through re-auth, and graceful expired-state UI.

**Reasoning / Evidence:** WCAG 2.2.1 (timing adjustable) requires warning + extension for session limits; data lost to silent expiry is a catastrophic-severity failure ([02-cognitive-foundations.md](02-cognitive-foundations.md#learned-helplessness-and-error-tolerant-design)). Banking-style short sessions collide with long-form tasks constantly.

**When to use:**
- Timed sessions: modal warning ≥60s before expiry ("Still there? Extend session") with one-click extension; screen-reader announced.
- Expiry with unsaved work: preserve locally, re-auth in place (modal/redirect-return), restore seamlessly ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)).
- API 401s mid-session: intercept globally → re-auth flow → replay/restore, never dump to login losing context.

**When NOT to use / exceptions:**
- High-security contexts may forbid local persistence of sensitive fields — preserve the non-sensitive scaffold, explain what must be re-entered.

**Pros:** Trust; compliance (2.2.1); support-ticket reduction.
**Cons:** Re-auth-and-restore flows are fiddly (deep links back to exact state).

**Implementation guidance:**
- Return-URL discipline: login redirects carry the full deep link (path + query) and restore it.
- Idle warning: 60s countdown, Extend + Save-and-logout options, `role="alertdialog"`.
- Test the brutal path: fill a long form, expire the session server-side, submit — the user must not lose work.

**Verification checklist:**
- [ ] Session expiry while a form is half-filled: warning appears before expiry, and work survives re-authentication.
- [ ] Post-login redirect lands on the exact deep link (path + query) the user was on.
- [ ] Expiry in a background tab does not destroy state when the user returns.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery), [05-accessibility.md](05-accessibility.md#cognitive-accessibility).

**Sources:** [Understanding WCAG 2.2.1](https://www.w3.org/WAI/WCAG22/Understanding/timing-adjustable.html).

---

## Performance as UX

### Core Web Vitals

**Definition:** Google's user-experience performance metrics with "good" thresholds at the 75th percentile of real users: **LCP** (Largest Contentful Paint — main content visible) ≤ **2.5s**; **INP** (Interaction to Next Paint — responsiveness of all interactions) ≤ **200ms**; **CLS** (Cumulative Layout Shift — visual stability) ≤ **0.1**.

**Reasoning / Evidence:** Each metric operationalizes a classic UX limit: LCP ≈ "is it loading?" (attention/abandonment window), INP ≈ the 0.1–0.4s response band ([04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds)), CLS ≈ change-blindness/mis-click harm from shifting layouts (rage-click generator: the button moves as you tap). Field studies across large sites correlate vitals improvements with bounce/conversion gains; also a Google Search ranking signal.

**Field vs lab:** Field data (RUM, CrUX) reflects real users on real devices and networks and should drive production priorities; lab data (Lighthouse, WebPageTest) is for pre-release regression catching and diagnosis, not verdicts. Evaluate at the 75th percentile and segment mobile vs desktop where possible — averages hide the users who suffer. Note that Lighthouse cannot directly measure real INP without user interaction; TBT (Total Blocking Time) serves as a lab proxy only.

**When to use:**
- Continuous field monitoring (CrUX/RUM) as UX health metrics, budgeted per template; lab tools (Lighthouse) for diagnosis.
- Design-stage prevention: reserve media space (CLS), question heavy heroes (LCP), audit long-task-causing interactions (INP).

**When NOT to use / exceptions:**
- Lab scores ≠ user reality (field data rules); vitals don't measure task success — a fast, unusable page still fails ([17-ux-process-research.md](17-ux-process-research.md#ux-metrics-frameworks)).
- Internal tools without search/SEO stakes still deserve INP-level responsiveness; thresholds are human, not Google-specific.

**Pros:** Shared UX-eng vocabulary with hard numbers; field-measurable; forcing function against regression.
**Cons:** Metric-gaming risk (skeleton-everything to fake LCP); INP debugging is genuinely hard (long tasks, hydration).

**Implementation guidance:**
- LCP: server-render the hero, preload the LCP image, no lazy-loading above the fold, compress aggressively (AVIF/WebP).
- CLS: `aspect-ratio`/dimensions on all media/embeds/ads; `font-display: swap` with metric-matched fallbacks; never inject banners above content post-load.
- INP: break long JS tasks (<50ms chunks), debounce expensive handlers, optimistic UI for slow mutations ([04-interaction-design.md](04-interaction-design.md#optimistic-ui)), virtualize huge lists.
- Set performance budgets per template for critical pages; verify layout shifts never move controls during interaction, and test under realistic (throttled) device/network conditions.

**Related:** [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds), [#perceived-performance-techniques](#perceived-performance-techniques).

**Sources:** [web.dev: Core Web Vitals](https://web.dev/articles/vitals), [web.dev: INP](https://web.dev/articles/inp), [web.dev: Optimize CLS](https://web.dev/articles/optimize-cls).

### Perceived performance techniques

**Definition:** Making waits feel shorter independent of actual duration: progressive rendering, skeletons, optimistic UI, prefetching, streaming, and honest-but-managed progress. (Core taxonomy: [04-interaction-design.md](04-interaction-design.md#loading-strategies).)

**Reasoning / Evidence:** Perception research: occupied time feels shorter than idle; determinate progress beats unknown; early meaningful content anchors "it's working." Prefetch-on-intent (hover/viewport-approach) hides hundreds of ms for near-zero cost.

**When to use:**
- Streaming/progressive render: shell first, content as ready (SSR streaming, out-of-order chunks).
- Prefetch: next-likely routes/data on link hover (~200ms head start) or viewport entry; paginated next-page prefetch.
- Instant transitions: navigate immediately to a skeletoned detail view using already-known data (title/thumb from the list) while the rest loads.

**When NOT to use / exceptions:**
- Prefetch discipline on metered/mobile connections (respect Data Saver; prefetch only high-probability targets).
- Don't fake progress bars disconnected from work; don't skeleton content whose shape you can't predict ([04-interaction-design.md](04-interaction-design.md#loading-strategies)).

**Pros:** Large felt-speed wins without infrastructure rewrites.
**Cons:** Complexity budget; over-prefetching wastes bandwidth/server.

**Implementation guidance:**
- Carry list-item data into detail routes for instant partial paint.
- `<link rel="prefetch">`/router prefetch on hover+focus; `IntersectionObserver` for scroll-approach.
- Measure with real-user LCP/INP before/after — perceived tricks must show up in behavior (bounce/engagement), not just demos.

**Related:** [04-interaction-design.md](04-interaction-design.md#loading-strategies), [#core-web-vitals](#core-web-vitals).

**Sources:** [web.dev: perceived performance](https://web.dev/learn/performance/why-speed-matters), [NN/g: Progress Indicators](https://www.nngroup.com/articles/progress-indicators/).

### Web typography specifics

**Definition:** Web-platform type mechanics: `rem`-based sizing (respecting user browser font settings), webfont loading strategy (FOUT/FOIT/CLS management), and system-font fallback engineering.

**Reasoning / Evidence:** Users set browser base font sizes for a reason; `px`-locked text ignores them (accessibility failure in spirit, 1.4.4 in letter — [05-accessibility.md](05-accessibility.md#zoom-and-text-resize)). Webfont loading historically caused invisible text (FOIT) or layout-jumping swaps (FOUT→CLS); modern stack (font-display, size-adjust metrics matching, preload) tames it.

**When to use:**
- All font sizes in `rem` (spacing may be rem or px by system convention); base = user default (16px nominal), never `html { font-size: 62.5% }` hacks that break assumptions.
- Webfonts: `font-display: swap` + metric-matched fallback (`size-adjust`, `ascent-override`) to make the swap visually near-invisible; preload the 1–2 critical files; WOFF2 only; subset.

**When NOT to use / exceptions:**
- System-font-stack products skip the whole problem ([03-visual-design.md](03-visual-design.md#font-selection-and-pairing)) — right answer for many apps.
- `font-display: optional` for truly cosmetic fonts (skip on slow loads).

**Pros:** User-preference respect; zero-jank loading achievable.
**Cons:** Metric-matching setup per font (tooling exists: Fontaine, Capsize).

**Implementation guidance:**
- Budget: ≤2 families, ≤4 files, subset to used scripts; self-host (privacy + speed vs third-party CDNs).
- Test: slow-3G load for FOUT behavior; browser font-size changed to 20px — layout must respect it.

**Related:** [03-visual-design.md](03-visual-design.md#font-selection-and-pairing), [#core-web-vitals](#core-web-vitals), [05-accessibility.md](05-accessibility.md#zoom-and-text-resize).

**Sources:** [web.dev: font best practices](https://web.dev/articles/font-best-practices), [MDN: font-display](https://developer.mozilla.org/en-US/docs/Web/CSS/@font-face/font-display).

---

## Architecture choices

### URLs as UX

**Definition:** Treating the URL as user interface: readable, guessable, shareable, stateful. Application state a user would want to share, bookmark, or return to (search query, filters, page, selected tab/item) belongs in the URL.

**Reasoning / Evidence:** NN/g ("URL as UI"): users read, edit, and infer from URLs; shared links are a primary growth/collaboration channel. State-in-URL is what makes Back work meaningfully, refresh non-destructive, and "send me the link" possible. Cool URIs don't change (W3C): link rot punishes users and SEO.

**When to use:**
- Filters/search/sort/pagination → query params; entity identity → path (`/projects/atlas/settings`); UI panel state worth sharing → param or path.
- Human-readable slugs over opaque IDs where feasible (`/blog/dark-mode-guide` not `/p/8842`); lowercase, hyphen-separated.
- Permanence: redirects (301) forever when structures change.

**When NOT to use / exceptions:**
- Ephemeral micro-state (hover, scroll, half-open accordions) stays out; sensitive data (tokens, PII) never in URLs (logs, referrers, shoulder-surfing).
- Long filter states may need shortening (saved-view IDs) beyond ~2000 chars.

**Pros:** Shareability, bookmarkability, Back/refresh integrity, support ("send me the URL you see"), SEO.
**Cons:** State-serialization discipline; migration/redirect maintenance.

**Implementation guidance:**
- Test: copy URL from any meaningful state → paste in incognito → same view (post-auth).
- `replaceState` for keystroke-level changes, `pushState` at "place" granularity ([#browser-conventions](#browser-conventions)).
- Structure: `/{section}/{entity-or-slug}/{sub-view}`; params for cross-cutting state (`?q=…&status=open&page=2`).

**Related:** [08-navigation-ia.md](08-navigation-ia.md#deep-linking-and-url-design), [#browser-conventions](#browser-conventions).

**Sources:** [NN/g: URL as UI](https://www.nngroup.com/articles/url-as-ui/), [W3C: Cool URIs don't change](https://www.w3.org/Provider/Style/URI).

### SPA vs MPA

**Definition:** Single-Page Application (client-side routing/rendering, app-like statefulness) vs Multi-Page Application (server-rendered documents per navigation) — plus the modern hybrid middle (server-rendered frameworks with client hydration/islands, and MPA + View Transitions).

**Reasoning / Evidence:** SPAs won app-like domains (editors, dashboards, chat: persistent state, instant in-app transitions, offline) but at recurring UX tax: hand-reimplementing Back/scroll/focus/announcements ([#browser-conventions](#browser-conventions), [05-accessibility.md](05-accessibility.md#keyboard-navigation)), hydration cost hurting INP/LCP, and memory-leaking long sessions. MPAs get browser correctness free and now get app-feel cheaply (speculation rules prefetch, cross-document View Transitions). The dichotomy is dissolving — choose per-surface.

**When to use (SPA/client-heavy):**
- True applications: persistent cross-view state (playing media, live collaboration, complex editing), offline requirements, sub-100ms in-app interactions.

**When to use (MPA/server-first):**
- Content, commerce, marketing, docs, dashboards-that-are-mostly-reads: navigation-centric surfaces where URL/back/SEO correctness is the experience.

**When NOT to use / exceptions:**
- Never SPA-by-default for content sites (common failure: blog as SPA — slower, buggier back, worse SEO for zero benefit).
- Hybrid reality: app shell SPA + server-rendered content routes is a legitimate architecture.

**Pros (SPA):** Statefulness, transition control, offline. **Cons (SPA):** Browser-behavior reimplementation burden, JS weight, a11y regressions by default.
**Pros (MPA):** Free correctness, resilience, simpler mental model. **Cons (MPA):** Cross-page state needs server/storage; historically clunky transitions (now largely solved).

**Implementation guidance:**
- If SPA: router must handle scroll restoration, focus-to-heading on navigation + live-region announcement, title updates, real `<a href>` links, error boundaries per route.
- If MPA: speculation-rules prefetch for instant-feel; View Transitions API for continuity ([04-interaction-design.md](04-interaction-design.md#view-transitions)).
- If server-rendering with hydration: format dates/numbers with the user's resolved locale passed explicitly, never the server's ambient locale — mismatched output causes hydration errors and text that flickers on load ([09 › Numbers, dates, and units formatting](09-content-ux-writing.md#numbers-dates-and-units-formatting)).
- Decide per route-group, not per company religion.

**Verification checklist (SPA routing):**
- [ ] Route changes update the document title.
- [ ] Route changes move focus to a logical heading/container (or preserve focus intentionally).
- [ ] Back/forward restore meaningful state, including scroll for list → detail → back flows.
- [ ] Unsaved user work is never lost silently on navigation.

**Related:** [#browser-conventions](#browser-conventions), [#core-web-vitals](#core-web-vitals), [04-interaction-design.md](04-interaction-design.md#view-transitions).

**Sources:** [NN/g: SPAs and seamless navigation](https://www.nngroup.com/articles/spa-seamless-navigation/), [MDN: View Transitions](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API).

### Progressive enhancement

**Definition:** Building in layers: functional HTML core → CSS presentation → JS enhancement, so the experience degrades gracefully when any layer fails (slow networks, old devices, JS errors, aggressive extensions, bots/readers).

**Reasoning / Evidence:** JS fails more than teams think (flaky networks mid-load, CDN hiccups, extension conflicts, corporate proxies) — GDS measured ~1% of *capable* browsers failing to run JS per visit; at scale that's constant breakage. Server-rendered HTML cores also align with LCP, SEO, and AT reliability. Full purity is impractical for true apps; the *posture* (core tasks survive degradation) scales to everything.

**When to use:**
- Public, broad-audience, or critical services (government, health, commerce checkout): forms that POST without JS, links that navigate, content in HTML.
- Everything: server-render first paint; treat JS as enhancement for speed/richness rather than existence-requirement where feasible.

**When NOT to use / exceptions:**
- Inherently-JS products (canvas editors, video calls, games): provide honest `<noscript>` messaging and accessible surrounding chrome instead.

**Pros:** Resilience, reach (old/cheap devices globally), performance floor, SEO/AT alignment.
**Cons:** Dual-path testing for enhanced/basic; some modern stacks fight it (mitigated by SSR frameworks with progressive form actions).

**Implementation guidance:**
- Core-task audit with JS disabled: can users read, navigate, submit the money path? Rank gaps by criticality.
- Native elements first ([05-accessibility.md](05-accessibility.md#semantic-html-first)) — they're the zero-JS baseline.
- Feature-detect (`@supports`, capability checks), never browser-sniff.
- Handle offline, poor-network, and failed-API states explicitly — degradation should be designed, not accidental.

**Related:** [05-accessibility.md](05-accessibility.md#semantic-html-first), [#spa-vs-mpa](#spa-vs-mpa).

**Sources:** [GOV.UK: Progressive enhancement](https://www.gov.uk/service-manual/technology/using-progressive-enhancement), [GDS: how many people disable JS](https://gds.blog.gov.uk/2013/10/21/how-many-people-are-missing-out-on-javascript-enhancement/).

### PWA capabilities

**Definition:** Progressive Web App features closing native gaps: installability (home screen/dock, standalone window), offline via service workers, push notifications, background sync, share targets, file handling.

**Reasoning / Evidence:** For repeat-use products, install + offline measurably lift engagement (persistent entry point, resilience); PWAs avoid app-store gatekeeping/fees and unify codebases. Limits: iOS support lags selectively (push arrived 2023, install less discoverable), discovery lacks store surface, deep hardware access is uneven — set expectations per capability, per platform.

**When to use:**
- Repeat-use web products with mobile traffic: offline shell + install prompt-after-engagement is cheap value.
- Emerging-market/flaky-network audiences: offline-first is transformative.
- Internal tools: install for app-like windows without distribution overhead.

**When NOT to use / exceptions:**
- One-shot visit sites (campaigns, articles): skip install prompts entirely.
- Never browser-prompt-on-load for install/notifications ([#notification-permission-ux](#notification-permission-ux)) — in-context, post-value custom prompts only.
- Needs deep native integration (advanced camera control, background geolocation, platform widgets)? Native remains right ([15-platform-mobile.md](15-platform-mobile.md#mobile-web-vs-native-app)).

**Pros:** One codebase, linkability retained, no store tax, instant updates.
**Cons:** iOS second-class realities; offline sync complexity (conflicts — [15-platform-mobile.md](15-platform-mobile.md#offline-and-poor-connectivity)); service-worker cache bugs are miserable (stale-app traps).

**Implementation guidance:**
- Offline: cache app shell + last-viewed content; explicit offline UI states (never silent stale data); background sync for queued mutations.
- Install: custom prompt after demonstrated engagement (2nd+ visit / task completion), never nag; `beforeinstallprompt` capture.
- Update flow: detect new SW → unobtrusive "Update available — Refresh" (never auto-reload over user state).

**Related:** [15-platform-mobile.md](15-platform-mobile.md#offline-and-poor-connectivity), [#notification-permission-ux](#notification-permission-ux).

**Sources:** [web.dev: Learn PWA](https://web.dev/learn/pwa/), [MDN: Progressive web apps](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps).

---

## Page-level behaviors

### Scroll UX

**Definition:** Web-specific scrolling discipline: no scrolljacking, restrained sticky elements, scroll restoration, reading-progress affordances, and anchor/fragment correctness. (Cross-platform doctrine: [04-interaction-design.md](04-interaction-design.md#scroll-behaviors).)

**Reasoning / Evidence:** NN/g: scrolljacking (hijacked speed/snap/direction) produces some of the strongest negative user reactions on record. Sticky-chrome over-accumulation (header + subnav + cookie bar + chat bubble) can consume a large fraction of a phone viewport — audit the stack as a whole, not element by element.

**When to use:**
- Sticky header kept short (rule of thumb: ~64px or less, or shrink-on-scroll), optional hide-down/show-up on reading surfaces; one sticky layer per edge, with total sticky chrome budgeted to roughly a quarter of the viewport or less (rule of thumb, not a researched threshold).
- Anchor links: `scroll-margin-top` compensating sticky headers (fragment jumps must not hide the target under chrome); smooth scroll behind `prefers-reduced-motion`.
- Back-to-top button after ~2–3 viewport-heights on long pages.

**When NOT to use / exceptions:**
- Full-page section-snap "storytelling" scroll: reserve for marketing one-offs, keep interruptible, honor reduced-motion — never on functional surfaces.
- Chat/CTA floating widgets covering content on mobile: collapse or remove.

**Implementation guidance:**
- CSS `scroll-behavior: smooth` (media-query gated), `scroll-margin-top: var(--header-height)`; passive scroll listeners; test with keyboard scrolling (Space/PageDown must not be hijacked).

**Related:** [04-interaction-design.md](04-interaction-design.md#scroll-behaviors), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#custom-scrollbars-and-scrolljacking).

**Sources:** [NN/g: Scrolljacking 101](https://www.nngroup.com/articles/scrolljacking-101/).

### Consent and cookie UX

**Definition:** Privacy-consent interfaces (GDPR/ePrivacy and successors) done lawfully and respectfully: equal-prominence Accept/Reject, no dark patterns, minimal interruption, honest scope.

**Reasoning / Evidence:** Legal floor is explicit: consent must be freely given, specific, informed, unambiguous; EU regulators (EDPB, national DPAs) have ruled that Reject must be as easy as Accept (first layer), pre-ticked boxes invalid, and "legitimate interest" abuse sanctionable. UX floor: consent walls are the first impression — hostile banners open the relationship with a fight ([Peak-End](01-core-principles.md#peak-end-rule) — bad start).

**When to use:**
- Only when actually needed: no banner at all if you use only essential cookies/storage (the genuinely best UX — and increasingly a differentiator).
- Required consent: compact banner (not full-screen wall where law permits), Accept-all and Reject-all equally prominent on layer one, granular controls one tap deeper, decision remembered and revisitable (persistent footer link).

**When NOT to use / exceptions:**
- Never: buried reject (2+ clicks vs 1), color-coded coercion (bright Accept, ghost Reject — regulators have fined this), fake "necessary" categorization, re-prompting after rejection ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#nagging)).

**Pros (honest consent):** Legal safety, trust signal, cleaner data (consented analytics are defensible).
**Cons:** Real analytics coverage loss with easy reject — the lawful cost of the business model; mitigate with cookieless/aggregate measurement.

**Implementation guidance:**
- Bottom banner, compact (well under a quarter of the viewport as a rule of thumb), focus-reachable but not focus-trapped-on-load, dismiss = no consent (not implied accept).
- Copy: plain language, purpose-first ("We use cookies for analytics and ads personalization"), no guilt trips ([09-content-ux-writing.md](09-content-ux-writing.md)).
- Honor Global Privacy Control signals where jurisdictionally relevant; log consent states auditable.

**Verification checklist (consent wiring integrity):**
- [ ] "Accept all" actually enables every category — an accept-all that silently leaves a category off corrupts the consent record.
- [ ] Rejecting all optional categories is persisted as an explicit rejection event, not as absence of a record.
- [ ] The stored consent decision is derived from the actual category choices, not from which button was pressed.

**Related:** [18-patterns-antipatterns.md](18-patterns-antipatterns.md#dark-patterns-catalog), [09-content-ux-writing.md](09-content-ux-writing.md).

**Sources:** [EDPB cookie banner taskforce report](https://www.edpb.europa.eu/our-work-tools/our-documents/other/report-work-undertaken-cookie-banner-taskforce_en), [GDPR Art. 7](https://gdpr-info.eu/art-7-gdpr/).

### Notification permission UX

**Definition:** Requesting browser push (and similar powerful permissions: location, camera) at the right moment: never on page load, always after context and expressed intent, via a two-step soft-prompt pattern.

**Reasoning / Evidence:** Permission-on-load has catastrophic accept rates (single-digit %) and worse: browser-level *permanent* denial (Chrome/Firefox now auto-suppress abusive prompt patterns) — one bad ask can lock you out forever. Contextual asks ("Get notified when your order ships" at order completion) invert the numbers.

**When to use:**
- Two-step: in-page custom prompt explaining concrete value ("Enable shipping alerts?") → only on yes, trigger the browser prompt (high accept, since pre-qualified).
- At the moment of relevant intent: after subscribing/ordering/enabling a feature that needs it — never speculative.

**When NOT to use / exceptions:**
- Products with no genuinely user-valuable push have no business asking; email/in-app cover most cases.
- Same doctrine for geolocation/camera/mic: ask in-context of the feature invocation, explain first.

**Pros:** Accept rates, preserved future askability, trust.
**Cons:** Fewer subscribers than spray-and-pray *appears* to yield — but denied users were never subscribers anyway.

**Implementation guidance:**
- Track soft-prompt dismissals; cool-down (weeks) before re-offering; never re-ask after hard denial (show settings-path instructions instead if the user later wants it).
- Pair with notification *content* discipline ([15-platform-mobile.md](15-platform-mobile.md#notifications-ux), [04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)).

**Related:** [15-platform-mobile.md](15-platform-mobile.md#notifications-ux), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#notification-permission-on-page-load).

**Sources:** [web.dev: permissions UX](https://web.dev/articles/push-notifications-permissions-ux), [NN/g: permission requests](https://www.nngroup.com/articles/permission-requests/).

### SEO-UX overlap

**Definition:** The shared substrate of search optimization and UX: semantic structure (headings, landmarks), descriptive titles/links, fast loading, mobile usability, structured data, and honest content — modern SEO is mostly UX made crawlable.

**Reasoning / Evidence:** Search engines rank what serves users: Core Web Vitals are a ranking input, heading structure serves both SR users and crawlers, descriptive link text serves both scanning and PageRank context, `<title>` is both the SERP headline and the tab/bookmark/history label users depend on.

**When to use:**
- Page titles: unique, front-loaded, `Primary thing – Section – Site` (~50–60 visible chars) — this is tab-management UX as much as SEO.
- One `h1` = page topic; heading hierarchy real ([05-accessibility.md](05-accessibility.md#keyboard-navigation)); descriptive link text ([09-content-ux-writing.md](09-content-ux-writing.md#link-text)).
- Structured data (Product, Article, FAQ, Breadcrumb) where content types match — richer SERPs are better first-impressions.
- Favicon set + social share metadata (OG/Twitter cards): links shared anywhere are your UI there.

**When NOT to use / exceptions:**
- SEO-driven content pollution (keyword-stuffed headings, FAQ-spam) degrades the UX it borrowed authority from — and algorithms increasingly punish it.

**Pros:** One investment, two returns; structure work compounds with accessibility.
**Cons:** None inherent; conflicts signal doing one of them wrong.

**Implementation guidance:**
- Title template enforced app-wide (SPAs: update on route change — also the SR announcement hook).
- Canonicals + 301 discipline ([#urls-as-ux](#urls-as-ux)); image alt serves both ([05-accessibility.md](05-accessibility.md#screen-readers)).

**Related:** [#urls-as-ux](#urls-as-ux), [05-accessibility.md](05-accessibility.md#semantic-html-first), [09-content-ux-writing.md](09-content-ux-writing.md), [#agent-friendly-web-content](#agent-friendly-web-content).

**Sources:** [Google: SEO Starter Guide](https://developers.google.com/search/docs/fundamentals/seo-starter-guide), [web.dev: vitals & ranking](https://developers.google.com/search/docs/appearance/core-web-vitals).

### Agent-friendly web content

**Definition:** Building pages that AI agents and assistive technology can parse and operate alike: meaningful content as real text (not text baked into images), semantic HTML structure, programmatically determinable form fields and labels, stable URLs for content, and structured data where content types match. (UX-WEB-008.)

**Reasoning / Evidence:** An emerging concern that is largely a restatement of existing robustness principles: the same properties that make a page work for screen readers (POUR's "robust" — [05-accessibility.md](05-accessibility.md)) and for search crawlers ([#seo-ux-overlap](#seo-ux-overlap)) make it work for AI agents acting on a user's behalf — summarizers, shopping/booking agents, research assistants. Agents, like AT, consume the accessibility tree and DOM semantics, not pixels: a `<div onclick>` "button," an image-of-text price, or an unlabeled input is invisible or ambiguous to all three audiences. There is no separate "agent optimization" discipline yet worth inventing — accessibility and SEO done properly cover most of it.

**When to use:**
- Any content users might want machines to read, compare, or act on: product/pricing pages, documentation, articles, schedules, forms.
- Core flows (search, checkout, booking): keep actions programmatically understandable — real `<form>`s, labeled fields, native or APG-correct widgets.

**When NOT to use / exceptions:**
- Content you deliberately don't want machine-extracted (paywalled bodies, anti-scraping surfaces) — a business decision, but be aware you're also degrading AT and search in most implementations.
- Don't add speculative agent-specific markup beyond established standards (semantic HTML, ARIA, schema.org); unproven conventions churn.

**Pros:** One investment serving accessibility, SEO, and agent interoperability; future-proofing as agent traffic grows; no conflict with human UX.
**Cons:** None inherent for honest content; genuine tension only where scraping-resistance is a requirement.

**Implementation guidance:**
- Important content as text: no text-in-images for prices, headings, policies; provide alt text that carries the information when images are unavoidable ([05-accessibility.md](05-accessibility.md#screen-readers)).
- Semantic structure: real headings in order, landmarks, lists, tables-with-headers; descriptive link text (no "click here").
- Forms: `<label>` associations, `autocomplete` attributes, programmatically determinable error/description text ([07-forms-and-input.md](07-forms-and-input.md)).
- Stable, meaningful URLs for content ([#urls-as-ux](#urls-as-ux)); structured data (schema.org) where applicable ([#seo-ux-overlap](#seo-ux-overlap)).
- Avoid requiring visual-only interaction (drag-only, canvas-only, hover-only) for core content or actions.
- Verification: if the page passes an accessibility audit and its content is readable with CSS/images disabled, it is largely agent-ready.

**Related:** [05-accessibility.md](05-accessibility.md#semantic-html-first), [#seo-ux-overlap](#seo-ux-overlap), [#links-vs-buttons](#links-vs-buttons), [20-ai-product-ux.md](20-ai-product-ux.md).

**Sources:** [W3C: WCAG POUR principles](https://www.w3.org/WAI/WCAG22/Understanding/intro), [Google: SEO Starter Guide](https://developers.google.com/search/docs/fundamentals/seo-starter-guide).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Responsive | Mobile-first, content-driven breakpoints | 320/768/1024/1440; reflow at 320px |
| Viewport | dvh not vh; safe-area insets; never block zoom | Inputs ≥16px (iOS zoom) |
| Browser conventions | Back/scroll/refresh/find must work | 5-deep back-button audit |
| Links vs buttons | href navigates, button acts — always | — |
| Multi-tab | Revalidate on focus; broadcast auth | — |
| Session expiry | Warn, extend, preserve, restore | ≥60s warning |
| Core Web Vitals | Field-measured UX floor (75th percentile) | LCP ≤2.5s; INP ≤200ms; CLS ≤0.1 |
| Perceived perf | Stream, prefetch on intent, carry data forward | Hover prefetch ≈200ms head start |
| Typography | rem sizing; metric-matched font swap | ≤2 families; WOFF2 subset |
| URLs | Shareable state; readable paths; 301 forever | Incognito-paste test |
| SPA vs MPA | App-ness earns SPA; content stays server-first | Decide per route-group |
| Progressive enhancement | Core tasks survive JS failure | ~1% JS-failure baseline |
| PWA | Offline + install after engagement | Custom prompt, never on load |
| Scroll | Never hijack; budget sticky chrome | Rules of thumb: ≤64px header; ≤25% viewport |
| Consent | Reject as easy as accept; or need no banner; stored decision mirrors actual category choices | Equal prominence, layer one |
| Permissions | Two-step soft prompt, in context | Never on load |
| SEO-UX | Structure once, benefit twice | Titles ~50–60 chars, front-loaded |
| Agent-friendly content | Real text, semantics, labeled forms — a11y/SEO done right covers agents | Readable with CSS/images off |
