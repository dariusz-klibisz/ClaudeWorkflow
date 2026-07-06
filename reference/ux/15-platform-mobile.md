# Platform: Mobile

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers mobile-native UX: the mobile usage context, touch targets and thumb zones, orientation and safe areas, Material Design 3 and Apple HIG conventions, cross-platform consistency strategy (with a compact Web/Desktop/Mobile comparison matrix), navigation and back behavior, gestures, mobile forms, notifications, onboarding, offline, performance perception, haptics, adaptive layouts (tablets/foldables) including the responsive-vs-adaptive distinction and localization across layouts, and mobile web vs native app trade-offs. Cross-platform interaction fundamentals are in [04-interaction-design.md](04-interaction-design.md); touch accessibility in [05-accessibility.md](05-accessibility.md#touchpointer-accessibility); mobile-web viewport mechanics in [13-platform-web.md](13-platform-web.md#viewport-quirks); desktop platform specifics in [14-platform-desktop.md](14-platform-desktop.md).

## Contents

- [Mobile context](#mobile-context)
  - [How mobile differs](#how-mobile-differs)
  - [Thumb zone](#thumb-zone)
  - [Touch targets](#touch-targets)
  - [Orientation and safe areas](#orientation-and-safe-areas)
- [Platform design languages](#platform-design-languages)
  - [Material Design 3 (Android)](#material-design-3-android)
  - [Apple HIG (iOS)](#apple-hig-ios)
  - [Cross-platform consistency vs platform-native](#cross-platform-consistency-vs-platform-native)
  - [Cross-platform comparison](#cross-platform-comparison)
- [Navigation](#navigation)
  - [Mobile navigation patterns](#mobile-navigation-patterns)
  - [Back behavior](#back-behavior)
  - [Deep links](#deep-links)
- [Input](#input)
  - [Gestures](#gestures)
  - [Mobile forms](#mobile-forms)
  - [Haptics](#haptics)
  - [Pull-to-refresh](#pull-to-refresh)
- [Communication](#communication)
  - [Notifications UX](#notifications-ux)
  - [Onboarding](#onboarding)
- [Resilience and performance](#resilience-and-performance)
  - [Offline and poor connectivity](#offline-and-poor-connectivity)
  - [Performance perception on mobile](#performance-perception-on-mobile)
- [Beyond the phone](#beyond-the-phone)
  - [Adaptive layouts: tablets and foldables](#adaptive-layouts-tablets-and-foldables)
  - [Dark mode on mobile](#dark-mode-on-mobile)
  - [Mobile web vs native app](#mobile-web-vs-native-app)
- [Quick reference](#quick-reference)

---

## Mobile context

### How mobile differs

**Definition:** Mobile usage is: one-handed much of the time (in Hoober's observational studies, ~49% of observed grips were one-handed, and ~75% of touch interactions were thumb-driven), interrupted (sessions average seconds-to-minutes, resumed constantly), glanceable (attention shared with the world), connectivity-variable, and sensor/identity-rich (camera, location, biometrics, always-signed-in).

**Reasoning / Evidence:** Steven Hoober's grip research and subsequent replications; session analytics across app categories show mobile sessions 3–10× shorter but more frequent than desktop. Design consequences: task flows must survive interruption ([state preservation](02-cognitive-foundations.md#interruption-cost)), primary actions must be thumb-reachable, and content must lead with the point ([02-cognitive-foundations.md](02-cognitive-foundations.md#scanning-patterns)).

**When to use:**
- Design for resumability: every flow survives backgrounding (OS kills apps freely) — state restored exactly on return.
- Front-load value per screen: one primary job per screen; progressive disclosure over dense layouts ([02-cognitive-foundations.md](02-cognitive-foundations.md#progressive-disclosure)).
- Exploit mobile superpowers: camera capture over file upload, biometric auth over passwords, location over address entry — each removes typing ([07-forms-and-input.md](07-forms-and-input.md#form-length-reduction)).

**When NOT to use / exceptions:**
- Tablet/landscape contexts relax one-handed constraints ([#adaptive-layouts-tablets-and-foldables](#adaptive-layouts-tablets-and-foldables)).
- Long-session mobile categories exist (games, reading, video) — glanceability rules apply to utility apps most.

**Pros (mobile-context design):** Fits actual usage; interruption-resilience improves *all* platforms.
**Cons:** One-job-per-screen lengthens flows (mitigate: smart defaults, fewer steps, not denser screens).

**Implementation guidance:**
- State restoration test: mid-flow → home → kill app → relaunch → exactly where you were, input preserved.
- Session design: assume 30–90s bursts; save continuously ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)).

**Related:** [#thumb-zone](#thumb-zone), [02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost).

**Sources:** [Hoober: How Do Users Really Hold Mobile Devices?](https://www.uxmatters.com/mt/archives/2013/02/how-do-users-really-hold-mobile-devices.php), [NN/g: Mobile UX](https://www.nngroup.com/articles/mobile-ux/).

### Thumb zone

**Definition:** The reachability map of one-handed phone use: the bottom-center third of the screen is easy; edges and especially the top corners (diagonal from the thumb's anchor) are hard, requiring grip shifts. Large phones worsen it.

**Reasoning / Evidence:** Hoober's grip studies; heat-map replications (Scott Hurff's thumb-zone maps). Platform evolution confirms: both OSes migrated key actions downward (iOS tab bars, Android navigation bar + FAB, bottom sheets everywhere, search bars moving bottom in Safari/Chrome experiments).

**When to use:**
- Primary actions and navigation: bottom third (tab bars, FABs, primary CTAs, bottom sheets).
- Destructive/rare actions: deliberately *out* of the easy zone (top corners) — inverse-Fitts protection ([01-core-principles.md](01-core-principles.md#fittss-law)).
- Reading content: top of screen is fine (eyes go everywhere; thumbs don't).

**When NOT to use / exceptions:**
- Two-handed contexts (typing, games, tablets) relax it.
- Don't cram *everything* bottom: 3–5 slots exist; hierarchy still governs ([#mobile-navigation-patterns](#mobile-navigation-patterns)).

**Pros:** Comfort, speed, fewer dropped-phone grip shifts.
**Cons:** Bottom real estate is contested (OS gesture bar, keyboard, tab bar, FAB, snackbars — choreograph stacking).

**Implementation guidance:**
- Primary CTA: bottom, full-width or FAB, above the OS gesture area (`safe-area-inset-bottom` — [13-platform-web.md](13-platform-web.md#viewport-quirks)).
- Screen-top actions get redundant paths (pull-down gestures, bottom-sheet equivalents) on large devices.

**Related:** [#touch-targets](#touch-targets), [01-core-principles.md](01-core-principles.md#fittss-law).

**Sources:** [Hoober: thumb zone research](https://www.uxmatters.com/mt/archives/2017/03/design-for-fingers-touch-and-people-part-1.php), [Hurff: thumb zones](https://scotthurff.com/posts/how-to-design-for-thumbs-in-the-era-of-huge-screens/).

### Touch targets

**Definition:** Minimum interactive sizes: Apple HIG **44×44pt**; Material **48×48dp** (with ≥8dp spacing); WCAG 2.5.8 legal floor **24×24 CSS px** ([05-accessibility.md](05-accessibility.md#target-size)). Visual element may be smaller; the *hit area* may not.

**Reasoning / Evidence:** Fingertip contact patches average ~7–10mm (~45–57px at typical densities); MIT Touch Lab and platform-vendor research converge on ~7mm/44–48 units as the error-rate knee. NN/g: target size failures are a top mobile usability defect, hitting motor-impaired and large-fingered users first, everyone in motion second.

**When to use:** Everything tappable — including icon buttons (biggest offender: 24px icons with 24px hit areas), list-row chevrons, close ×s, checkbox rows (make the whole row tappable), inline links (avoid tiny adjacent links).

**When NOT to use / exceptions:** Dense text links mid-paragraph are the WCAG-exempt case — still avoid stacking them.

**Pros:** Directly measurable error-rate reduction. **Cons:** Layout pressure in dense toolbars (solve with padding, not shrinkage).

**Implementation guidance:**
- Hit area = visual + padding to ≥44pt/48dp; adjacent targets ≥8dp apart.
- Audit tool: enable layout-bounds debugging (Android) / hit-area overlays; tap-error analytics (rage taps, mis-tap corrections) flag violators.

**Verification checklist:**
- [ ] Every tappable element has a hit area ≥44×44pt (iOS) / 48×48dp (Android) — measured, not eyeballed (layout-bounds / hit-area overlay on).
- [ ] Adjacent targets have ≥8dp spacing between hit areas.
- [ ] Icon buttons, list-row chevrons, and close buttons audited specifically (the most common violators).
- [ ] Tap-error analytics (rage taps, immediate corrections) reviewed for hotspots.

**Related:** [05-accessibility.md](05-accessibility.md#target-size), [#thumb-zone](#thumb-zone).

**Sources:** [Apple HIG: Layout](https://developer.apple.com/design/human-interface-guidelines/layout), [Material: Touch targets](https://m3.material.io/foundations/designing/structure), [NN/g: Touch Targets](https://www.nngroup.com/articles/touch-target-size/).

### Orientation and safe areas

**Definition:** Two device-variation obligations: (1) **orientation** — do not lock to portrait or landscape unless a single orientation is essential to the task (WCAG 1.3.4 Orientation makes forced orientation an accessibility failure: users with mounted devices — wheelchairs, bed mounts — cannot rotate); (2) **safe areas** — layout must respect notches, camera cutouts, rounded corners, the home-indicator/gesture bar, and system status bars via platform safe-area insets.

**Reasoning / Evidence:** Orientation locks assume a hand-held, rotatable device; mounted-device users and users who prefer landscape (larger text, external keyboards) are locked out — hence WCAG 1.3.4's Level AA requirement that content not restrict operation to a single orientation unless essential. Safe-area violations (content under notches, buttons under the gesture bar) are among the most visible "unfinished app" signals and cause real mis-taps at the screen edges.

**When to use:**
- Support both orientations by default; let layouts reflow (not just rotate-and-stretch) — landscape often wants two panes or wider content measure.
- Apply safe-area insets to all edge-adjacent UI: bottom bars/CTAs above the gesture area, top content below status bar/notch, horizontal insets in landscape (notch moves to the side).

**When NOT to use / exceptions:**
- Essential-orientation cases exist (a piano app, a landscape-only video game level, a check-scanning camera view) — WCAG's "essential" exception; document the justification in accessibility review.
- If an orientation lock is genuinely essential, still verify the locked experience works with large text and external keyboards.

**Pros:** Works for mounted-device users; survives foldable/tablet postures. **Cons:** Two-orientation QA per screen (mitigate: size-class-driven layouts make orientation just another width — [#adaptive-layouts-tablets-and-foldables](#adaptive-layouts-tablets-and-foldables)).

**Implementation guidance:**
- Rotation preserves state and scroll position — rotating mid-form must never lose input.
- Use platform inset APIs (`safeAreaInsets`, `WindowInsets`, CSS `env(safe-area-inset-*)` for web — [13-platform-web.md](13-platform-web.md#viewport-quirks)); never hardcode status-bar heights.

**Verification checklist:**
- [ ] App usable in both orientations, or the lock is documented as essential (WCAG 1.3.4).
- [ ] Rotation mid-flow preserves state, input, and scroll position.
- [ ] No content or controls under the notch, status bar, or home-indicator area on notched devices — both orientations tested.
- [ ] Landscape tested on a notched device (side cutout) and with the keyboard open.

**Related:** [05-accessibility.md](05-accessibility.md), [#adaptive-layouts-tablets-and-foldables](#adaptive-layouts-tablets-and-foldables), [#thumb-zone](#thumb-zone).

**Sources:** [WCAG 2.1: Orientation (1.3.4)](https://www.w3.org/WAI/WCAG21/Understanding/orientation.html), [Apple HIG: Layout](https://developer.apple.com/design/human-interface-guidelines/layout), [Android: display cutouts](https://developer.android.com/develop/ui/views/layout/display-cutout).

---

## Platform design languages

### Material Design 3 (Android)

**Definition:** Google's current design system for Android: **dynamic color** (palettes generated from user wallpaper, applied via tonal roles), tonal elevation (surfaces tinted, not just shadowed), updated components (navigation bar/rail/drawer, FAB variants, chips, cards), large-corner shape language, emphasized-easing motion, and window size classes for adaptive layout — extended by "M3 Expressive" refinements.

**Reasoning / Evidence:** M3 is the OS-native visual language: users' system apps, settings, and expectations render in it. Dynamic color personalizes the whole device coherently — apps opting in feel *of the device*; apps ignoring it feel foreign (a legitimate brand choice, consciously made).

**When to use:**
- Android-first/Android-only apps: adopt M3 components and theming; map brand colors through M3's tonal-role system (don't fight it with hardcoded hexes).
- Dynamic color: support it for system-harmony (utility apps especially); brand-critical surfaces may pin brand seeds instead.
- Standard components first: navigation bar (bottom, 3–5), top app bar, FAB for the *single* signature creative action, bottom sheets, M3 date/time pickers.

**When NOT to use / exceptions:**
- Heavy-brand consumer apps (games, media identity apps) legitimately override visuals — but keep *behavioral* Material (back handling, touch feedback/ripple, navigation patterns).
- FAB: skip when there's no single dominant action (a misused FAB is prime screen real estate wasted).

**Pros:** Native feel, free theming/dark-mode machinery, accessibility-tested components, user personalization harmony.
**Cons:** Brand dilution risk; M2→M3 migration churn; dynamic color complicates brand QA (test extreme wallpaper palettes).

**Implementation guidance:**
- Theme via tonal roles (primary/secondary/tertiary + surface tones), not raw colors — dark mode then ships nearly free ([03-visual-design.md](03-visual-design.md#dark-mode-design), [10-design-systems.md](10-design-systems.md#theming-dark-mode-multi-brand-density)).
- Type: Roboto/Google Sans-adjacent scale, respect user font-size settings ([05-accessibility.md](05-accessibility.md#zoom-and-text-resize)).
- Motion: M3 duration/easing tokens (emphasized ~500ms for large, standard 200–300ms) with reduced-motion respect.

**Related:** [#cross-platform-consistency-vs-platform-native](#cross-platform-consistency-vs-platform-native), [10-design-systems.md](10-design-systems.md#reference-systems-compared).

**Sources:** [Material Design 3](https://m3.material.io/), [Android design for developers](https://developer.android.com/design/ui).

### Apple HIG (iOS)

**Definition:** Apple's Human Interface Guidelines for iOS: SF Pro type with Dynamic Type scaling, SF Symbols icon system, tab bars (bottom, 3–5), navigation stacks with large titles and universal swipe-back, sheets (detents: medium/large) as the modal workhorse, context menus via long-press, subtle depth/materials (blur/translucency — and the "Liquid Glass" direction in current iOS), restrained color with tint-color semantics.

**Reasoning / Evidence:** iOS users share strong learned behaviors: edge-swipe back is *reflexive*; sheets are dismissed by drag; long-press previews are expected. HIG compliance isn't aesthetic conformity — it's muscle-memory compatibility ([Jakob's Law](01-core-principles.md#jakobs-law) at OS level). App Review also enforces slices of it.

**When to use:**
- Structure: tab bar for top-level (3–5), navigation stack per tab (state preserved per tab), sheets for scoped tasks, full-screen covers sparingly.
- System components first: SF Symbols (thousands of glyphs, weight-matched to text), system pickers/share sheet/context menus, Dynamic Type mandatory support.
- Swipe-back: never block the leading-edge back gesture; never attach conflicting edge gestures.

**When NOT to use / exceptions:**
- Cross-platform brand surfaces may diverge visually — but behavioral iOS-isms (swipe back, sheet dismissal, tab persistence) must survive the divergence.
- Don't import Android idioms (system-back assumptions, FAB-as-primary, top drawers) wholesale.

**Pros:** Muscle-memory fit, free accessibility (Dynamic Type, VoiceOver-tested components), review-friction avoidance.
**Cons:** Distinctiveness constraints; HIG shifts with iOS versions (design refreshes require keeping pace).

**Implementation guidance:**
- Dynamic Type: test at largest accessibility sizes — layouts must reflow, not truncate ([05-accessibility.md](05-accessibility.md#zoom-and-text-resize)).
- Navigation: large title collapsing on scroll; back button labeled with the previous screen where space allows.
- Haptics via system feedback generators tied to meaning ([#haptics](#haptics)); respect Reduce Motion/Transparency.

**Related:** [#cross-platform-consistency-vs-platform-native](#cross-platform-consistency-vs-platform-native), [#back-behavior](#back-behavior).

**Sources:** [Apple HIG](https://developer.apple.com/design/human-interface-guidelines/), [Apple HIG: Designing for iOS](https://developer.apple.com/design/human-interface-guidelines/designing-for-ios).

### Cross-platform consistency vs platform-native

**Definition:** The strategic decision for dual-platform apps: one unified design everywhere (brand consistency, single design/eng effort — the Flutter/React Native default posture) vs per-platform nativeness (two adapted designs) vs the pragmatic hybrid: **shared content/brand layer, platform-adapted behavioral layer**.

**Reasoning / Evidence:** Users don't compare your iOS and Android builds — they compare your app to *other apps on their device* ([01-core-principles.md](01-core-principles.md#4-consistency-and-standards)). Behavioral violations (wrong back behavior, non-native controls with wrong physics) test worse than visual unification; visual brand sameness with correct per-platform behavior is the documented sweet spot for most products.

**When to use (hybrid, the default recommendation):**
- Unify: brand color/type/illustration, content structure, feature set, iconography style (mapped to SF Symbols/Material Symbols equivalents), terminology and object model.
- Adapt per platform: navigation chrome (tab bar semantics differ subtly), back behavior ([#back-behavior](#back-behavior)), pickers/dialogs/share sheets (always system), text editing behaviors (selection loupes/handles), scroll physics, haptics, settings-screen conventions.

**When NOT to use / exceptions:**
- Games/immersive media: fully custom worlds are the norm — platform conventions apply only at system touchpoints.
- Tiny teams may ship one design first — accepting measured friction on the non-primary platform as debt.

**Pros (hybrid):** Brand + belonging; one design system with platform tokens ([10-design-systems.md](10-design-systems.md#theming-dark-mode-multi-brand-density)).
**Cons:** Requires platform-fluent designers/QA on both; component library must support behavioral forks (Flutter/RN can — with deliberate effort).

**Implementation guidance:**
- Fork list (minimum): back handling, system pickers/share, context menus (long-press vs varying), notification patterns, typography scale mapping (SF↔Roboto metrics differ), edge gestures.
- Document every intentional platform deviation — an undocumented fork becomes an accidental inconsistency at the next redesign.
- QA with native-platform users per build; the "feels off" reports are behavioral 90% of the time.

**Related:** [01-core-principles.md](01-core-principles.md#jakobs-law), [10-design-systems.md](10-design-systems.md), [14-platform-desktop.md](14-platform-desktop.md#cross-platform-convention-mapping), [#cross-platform-comparison](#cross-platform-comparison).

**Sources:** [Material: platform adaptation](https://m3.material.io/foundations/adaptive-design), [React Native: platform-specific code](https://reactnative.dev/docs/platform-specific-code).

### Cross-platform comparison

**Definition:** A compact decision matrix for what changes across web, desktop, and mobile/touch — the "adapt" column checklist for the hybrid strategy above. Rule of thumb: **shared** = concepts, terminology, object model, core IA, design tokens; **adapted** = navigation, controls, density, gestures, shortcuts, permissions, notifications, file handling.

| Topic | Web | Desktop | Mobile/Touch |
|---|---|---|---|
| Navigation | Links, URLs, browser history, landmarks | Menus, sidebars, tabs, windows, documents | Bottom nav, tabs, drill-down, gestures with alternatives |
| Primary input | Keyboard, pointer, touch, AT | Keyboard, pointer, shortcuts, menus | Touch, virtual keyboard, sensors, voice, biometrics |
| Back behavior | Browser back / route history | App/window/document conventions | System/app back, navigation stack |
| Density | Responsive, user zoom | Higher density acceptable; density controls useful | Lower density, larger targets |
| Shortcuts | Avoid browser/OS conflicts | Expected for frequent commands | Rare; external keyboard optional |
| Files | Downloads/uploads/browser sandbox | Native file dialogs, recent files | Pickers, share sheets, cloud providers |
| Notifications | Browser permission; use sparingly | OS notifications/status/tray | OS notifications with interruption sensitivity |
| Permissions | Browser permission prompts | OS/app permissions | Contextual runtime permissions |
| Accessibility | WCAG / semantic HTML / ARIA | OS accessibility APIs | Platform accessibility APIs, dynamic type, touch targets |
| Performance | Core Web Vitals, network variability | Startup, memory, responsiveness | Battery, network, startup, offline/interruption |

**Implementation guidance:**
- Detect capabilities (pointer precision, hover, viewport, input methods), not device categories — a tablet with keyboard and mouse is effectively a laptop; a touch-screen laptop needs touch targets ([#adaptive-layouts-tablets-and-foldables](#adaptive-layouts-tablets-and-foldables)).
- Web specifics (semantic HTML, history, Core Web Vitals): [13-platform-web.md](13-platform-web.md). Desktop specifics — including Windows/Fluent, GNOME ("design for people, reduce user effort"), and KDE ("simple by default, powerful when needed") platform notes: [14-platform-desktop.md](14-platform-desktop.md). CLI/TUI: [16-platform-cli-tui.md](16-platform-cli-tui.md).

**Related:** [#cross-platform-consistency-vs-platform-native](#cross-platform-consistency-vs-platform-native), [13-platform-web.md](13-platform-web.md), [14-platform-desktop.md](14-platform-desktop.md), [16-platform-cli-tui.md](16-platform-cli-tui.md).

**Sources:** [Material: platform adaptation](https://m3.material.io/foundations/adaptive-design), [Apple HIG](https://developer.apple.com/design/human-interface-guidelines/).

---

## Navigation

### Mobile navigation patterns

**Definition:** The mobile top-level navigation repertoire: **bottom tab bar** (3–5 destinations, the modern default), **navigation drawer** (hamburger side panel — demoted from default status), **hub-and-spoke** (home menu → dedicated flows), **top tabs** (swipeable siblings within a section), and **bottom sheets** for contextual action sets.

**Reasoning / Evidence:** Hidden-navigation research (NN/g): drawer-hidden destinations get roughly half the engagement of visible tabs; both platforms converged on bottom tabs for reachability ([#thumb-zone](#thumb-zone)) + visibility ([01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall)). Drawers survive legitimately for long-tail destination lists *behind* a tabbed core.

**When to use:**
- Bottom tabs: 3–5 true top-level destinations, always-visible, with icon + **label** (icon-only tabs measurably fail — [03-visual-design.md](03-visual-design.md#iconography)), per-tab state preservation.
- Drawer: secondary/long-tail destinations (settings, help, account switcher) supplementing tabs — not replacing them.
- Hub-and-spoke: task-launcher products (banking: many discrete flows) — see [08-navigation-ia.md](08-navigation-ia.md#hub-and-spoke).
- Top tabs (swipeable): content siblings within one destination (categories in a feed).

**When NOT to use / exceptions:**
- >5 top-level candidates = IA problem, not a bigger-tab-bar problem ([08-navigation-ia.md](08-navigation-ia.md#hierarchy-depth-vs-breadth)).
- Don't combine drawer + tabs *as peers* (two competing top levels confuse); drawer subordinates.
- "More" tab as a dumping ground: acceptable pressure valve, but everything important lives in the visible 4.
- Don't bury primary tasks in overflow menus — if a task is core, it earns visible chrome.

**Pros (tabs):** Visibility, reach, state persistence, orientation ([wayfinding](08-navigation-ia.md#wayfinding-and-you-are-here-signals)).
**Cons (tabs):** Hard cap of 5; contested bottom space (keyboards, sheets, FABs — choreograph).

**Implementation guidance:**
- Tab bar: 49pt (iOS) / 80dp (M3 nav bar) heights per platform; active state clear (color + filled-icon variant); never hide the bar on inner screens of a tab (users lose the map) except full-screen media.
- Re-tap active tab = scroll-to-top / pop-to-root (both platforms' convention).
- Badge counts on tabs sparingly ([04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)).

**Related:** [08-navigation-ia.md](08-navigation-ia.md#bottom-navigation-mobile), [08-navigation-ia.md](08-navigation-ia.md#hamburger-menu), [#thumb-zone](#thumb-zone).

**Sources:** [NN/g: Mobile Navigation Patterns](https://www.nngroup.com/articles/mobile-navigation-patterns/), [Material: Navigation bar](https://m3.material.io/components/navigation-bar/overview), [Apple HIG: Tab bars](https://developer.apple.com/design/human-interface-guidelines/tab-bars).

### Back behavior

**Definition:** The platform-divergent back model: **Android** — universal system back (gesture/button) traversing a back stack the app must curate predictably (with predictive-back previews in current Android); **iOS** — per-screen navigation-bar back button + leading-edge swipe-back gesture, no system-wide back.

**Reasoning / Evidence:** Back mishandling is a top Android complaint class (back exiting the app unexpectedly, back loops, back losing input); Android's predictive back (showing where back will land before committing) exists precisely because apps were unpredictable. iOS swipe-back is so ingrained that blocking it (custom gesture conflicts, full-screen takeovers without exits) reads as broken.

**When to use:**
- Android: back = reverse chronological within the app's sensible hierarchy; from a deep-linked screen, back synthesizes the logical parent stack (not exit); double-back-to-exit pattern with toast on the root only.
- iOS: every pushed screen swipe-backable; back button labeled with parent title when space allows; modally-presented sheets dismiss by swipe-down, not back.
- Both: back never silently discards meaningful input — unsaved-changes guard ([07-forms-and-input.md](07-forms-and-input.md#destructive-action-confirmation)).

**When NOT to use / exceptions:**
- Blocking back is legitimate mid-payment-processing (with visible state); otherwise never trap.
- In-app browsers and cross-app flows: back must return to the app sanely (test the full journey).

**Pros (correct back):** Trust in exploration (back = safe exit — [01-core-principles.md](01-core-principles.md#3-user-control-and-freedom)).
**Cons:** Back-stack curation for deep links and multi-entry flows is genuinely intricate — design the stack map explicitly.

**Implementation guidance:**
- Map every screen: "back from here goes exactly where?" — write it down, test deep-link entries.
- Android: implement predictive-back callbacks; never override system back for custom drawers/sheets without closing them first (back closes topmost layer, then navigates).
- Cross-platform frameworks: this is fork item #1 ([#cross-platform-consistency-vs-platform-native](#cross-platform-consistency-vs-platform-native)).

**Related:** [13-platform-web.md](13-platform-web.md#browser-conventions), [#deep-links](#deep-links).

**Sources:** [Android: back navigation](https://developer.android.com/guide/navigation/custom-back), [Android: predictive back](https://developer.android.com/guide/navigation/custom-back/predictive-back-gesture), [Apple HIG: Navigation bars](https://developer.apple.com/design/human-interface-guidelines/navigation-bars).

### Deep links

**Definition:** URLs opening specific in-app locations (universal links/app links), with web fallbacks, correct back-stack synthesis, and deferred deep linking (install → land on the intended content).

**Reasoning / Evidence:** Links are how content travels (messages, email, search, QR); apps that open to home instead of the linked content break the fundamental contract ([13-platform-web.md](13-platform-web.md#urls-as-ux) applies to apps too). Broken deep links silently divert to bad app-store interstitials ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#app-interstitials)).

**When to use:**
- Every shareable entity (product, post, profile, document) gets a canonical URL that opens natively when the app is installed, web otherwise.
- Back-stack synthesis: deep-link entry builds the logical parent chain so back ascends sensibly ([#back-behavior](#back-behavior)).
- Notification taps = deep links to the *specific* referenced content, never the app home.

**When NOT to use / exceptions:**
- Auth-gated content: deep link → auth → *continue to the original target* (return-URL discipline — [13-platform-web.md](13-platform-web.md#session-expiry-handling)).

**Pros:** Share loops, notification efficacy, search/app-indexing traffic.
**Cons:** Route-mapping maintenance across app versions; test-matrix (installed/not, logged-in/not, cold/warm).

**Implementation guidance:**
- Universal Links (iOS) / App Links (Android) over custom schemes (which lack web fallback and are hijackable).
- Test matrix per link type: app installed cold/warm/background × authed/not × valid/deleted target (deleted target → graceful in-app explanation, not crash/home).

**Related:** [13-platform-web.md](13-platform-web.md#urls-as-ux), [#notifications-ux](#notifications-ux).

**Sources:** [Android App Links](https://developer.android.com/training/app-links), [Apple Universal Links](https://developer.apple.com/documentation/xcode/allowing-apps-and-websites-to-link-to-your-content).

---

## Input

### Gestures

**Definition:** Mobile's gesture vocabulary as platform infrastructure: system-owned edges (iOS back-swipe left edge, home indicator bottom; Android back both edges, home bottom) constrain app gestures; standard app gestures (swipe actions on rows, swipe-dismiss sheets/images, pinch-zoom media, long-press context menus) ride on convention. Doctrine: [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures) — redundancy rule and discoverability limits apply doubly on mobile.

**Reasoning / Evidence:** Edge-gesture conflicts are the modern minefield: an app's edge-drawer fights Android back gestures; custom bottom swipes fight home indicators. Platform HIGs explicitly reserve the edges.

**When to use:**
- Standard set with visible equivalents: swipe row → action buttons also in long-press menu/detail; swipe-dismiss → close button too; pinch → zoom controls available (a11y — [05-accessibility.md](05-accessibility.md#touchpointer-accessibility)).
- Long-press: context menus/previews (both platforms now converge).

**When NOT to use / exceptions:**
- No functional gestures at screen edges; no custom multi-finger requirements without alternatives (WCAG 2.5.1); no gesture-only features period.

**Implementation guidance:**
- Swipe rows: reveal-actions style (partial swipe shows buttons; full swipe = default action) with color+icon coding; destructive full-swipe requires undo ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- First-use hint animation once; respect system gesture insets APIs.

**Related:** [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures), [#back-behavior](#back-behavior).

**Sources:** [Apple HIG: Gestures](https://developer.apple.com/design/human-interface-guidelines/gestures), [Android: gesture navigation](https://developer.android.com/develop/ui/views/touch-and-input/gestures/gesturenav).

### Mobile forms

**Definition:** Forms optimized for the worst input environment in computing: correct keyboard per field, minimal typing via platform superpowers (autofill, OTP auto-read, camera scan, contact pickers, biometrics), and thumb-reachable single-column flow. (Foundation: [07-forms-and-input.md](07-forms-and-input.md).)

**Reasoning / Evidence:** Mobile typing is ~2–3× slower and more error-prone than desktop; every removed keystroke compounds. Platform affordances remove entire fields: OS autofill (addresses, payments, passwords/passkeys), SMS OTP auto-fill (`one-time-code` / SMS Retriever), card-scan via camera, biometric re-auth replacing password re-entry.

**When to use:**
- Keyboard types: email/number/phone/URL keyboards via input types; correct capitalization/autocorrect flags per field (names: caps-words, no autocorrect; codes: no autocorrect).
- Autofill: full attribute/hint coverage (this is also WCAG 1.3.5 / 3.3.8 — [05-accessibility.md](05-accessibility.md#accessible-authentication)); passkeys/biometric prompts over password fields wherever possible.
- Field order: keyboard "next" advances logically; submit reachable from final field ("go"/"done" action).

**When NOT to use / exceptions:**
- Number keyboard for numeric *codes* with possible leading zeros/letters: use text + `inputmode` ([07-forms-and-input.md](07-forms-and-input.md#input-type-selection)).
- Don't auto-advance between fields ([07-forms-and-input.md](07-forms-and-input.md#input-masking-and-formatting)).
- Camera scan, contact pickers, location: offer only when helpful and permission-appropriate — a permission prompt to save ten keystrokes is a bad trade on first use.

**Implementation guidance:**
- Keyboard-avoidance: focused field scrolls above keyboard, submit visible or reachable; test small devices + large text.
- One column always; labels above and always visible ([07-forms-and-input.md](07-forms-and-input.md#label-placement)); 16px+ input font (also prevents iOS zoom — [13-platform-web.md](13-platform-web.md#viewport-quirks)).

**Related:** [07-forms-and-input.md](07-forms-and-input.md), [05-accessibility.md](05-accessibility.md#accessible-authentication).

**Sources:** [Android: autofill](https://developer.android.com/guide/topics/text/autofill), [Apple: Password AutoFill](https://developer.apple.com/documentation/security/password_autofill), [web.dev: mobile forms](https://web.dev/learn/forms/mobile/).

### Haptics

**Definition:** Tactile feedback via device vibration motors, semantically tiered (iOS: impact light/medium/heavy, notification success/warning/error, selection ticks; Android: equivalent haptic constants) — a feedback channel that works eyes-free.

**Reasoning / Evidence:** Haptics confirm actions without visual attention (glanceable context — [#how-mobile-differs](#how-mobile-differs)) and add perceived quality (the iOS toggle/picker feel). Overuse numbs: constant buzzing habituates and drains battery/goodwill ([02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)).

**When to use:**
- Meaningful moments: action confirmation (send, complete), state boundaries (toggle, picker detents, pull-to-refresh threshold), errors (distinct pattern), drag interactions (pickup/drop ticks).

**When NOT to use / exceptions:**
- Not for every tap (system keyboards handle their own); not for passive events (scroll); respect system haptics-off settings; never as the *sole* feedback channel (pair with visual — [05-accessibility.md](05-accessibility.md)).

**Pros:** Eyes-free confirmation, perceived polish. **Cons:** Habituation on overuse; hardware variance (Android motor quality varies wildly — test on mid-range devices).

**Implementation guidance:**
- Use semantic system APIs (UIFeedbackGenerator / HapticFeedbackConstants), not raw vibration durations; map consistently: one meaning = one pattern product-wide.

**Related:** [04-interaction-design.md](04-interaction-design.md#microinteractions), [01-core-principles.md](01-core-principles.md#feedback).

**Sources:** [Apple HIG: Playing haptics](https://developer.apple.com/design/human-interface-guidelines/playing-haptics), [Android: haptics](https://developer.android.com/develop/ui/views/haptics).

### Pull-to-refresh

**Definition:** The drag-down-past-top gesture triggering content refresh — a solved convention (spinner reveal, threshold haptic, release-to-refresh) for reverse-chronological/feed content.

**Reasoning / Evidence:** Invented by Loren Brichter in Tweetie 2 (2009), now universal for feeds; users *attempt it speculatively* on any scrollable list — absence where expected reads as staleness.

**When to use:** Feeds, inboxes, timelines, dashboards — any "what's new?" content; alongside (not replacing) auto-refresh logic.

**When NOT to use / exceptions:** Static content, paginated tables mid-state, screens where top ≠ newest; don't conflict with collapsing-header gestures.

**Implementation guidance:** Platform components (they handle physics/thresholds); threshold haptic tick; show refresh result ("3 new posts" pill) rather than silent completion ([01-core-principles.md](01-core-principles.md#1-visibility-of-system-status)).

**Related:** [#gestures](#gestures), [04-interaction-design.md](04-interaction-design.md#microinteractions).

**Sources:** [Apple HIG: Refresh content controls](https://developer.apple.com/design/human-interface-guidelines/refresh-content-controls), [Material: pull to refresh](https://m3.material.io/components/loading-indicator/overview).

---

## Communication

### Notifications UX

**Definition:** Mobile push done sustainably: permission earned in context (pre-permission priming), content that's user-relevant and actionable, Android channels / iOS interruption-levels used honestly, deep-linked taps, and volume discipline against fatigue.

**Reasoning / Evidence:** Opt-in rates collapse for on-launch asks (single-digit-to-low-double %) and recover with contextual two-step priming ([13-platform-web.md](13-platform-web.md#notification-permission-ux) — same math natively). Notification fatigue drives app *uninstalls*, not just mutes; per-channel controls (Android) mean users prune categories — apps without honest categories get all-or-nothing muted.

**When to use:**
- Permission: after demonstrated value, tied to a concrete benefit ("Alert me when the price drops"), via in-app pre-prompt → OS prompt on yes; iOS provisional (quiet) delivery as a soft entry.
- Content: user-triggered events (reply, shipment, mention) > business-triggered (promos); actionable inline (reply, mark done, snooze); deep-linked to the exact object ([#deep-links](#deep-links)).
- Channels/categories: honest granularity (orders / messages / promotions separately) so pruning ≠ total loss.

**When NOT to use / exceptions:**
- Never: on-launch permission ask, promo storms through transactional channels (users channel-mute you forever, and platforms police it), notifications as engagement-decay Hail Marys ("we miss you" — uninstall bait — [18-patterns-antipatterns.md](18-patterns-antipatterns.md#nagging)).
- Quiet hours respect; batch low-priority digests.

**Pros (disciplined push):** Retention channel with real utility; re-engagement that's welcomed.
**Cons:** Organizational pressure to over-send is relentless — set and defend volume budgets with uninstall-correlation data.

**Implementation guidance:**
- Budget: cap per-day/week by category; measure mute/uninstall rates per campaign, not just opens.
- Rich content where valuable (images, progress); update/replace stale notifications (order status evolves in one notification, not five).
- In-app settings mirror OS controls (per-category toggles) — discoverable, honest.

**Related:** [04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design), [02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost), [13-platform-web.md](13-platform-web.md#notification-permission-ux).

**Sources:** [NN/g: Push notification guidelines](https://www.nngroup.com/articles/push-notification/), [Android: notification channels](https://developer.android.com/develop/ui/views/notifications/channels), [Apple HIG: Managing notifications](https://developer.apple.com/design/human-interface-guidelines/managing-notifications).

### Onboarding

**Definition:** The first-launch path to first value: skippable by default, value-first (show the product, not slides about it), permissions primed in context (not front-loaded), and progressive teaching at moments of relevance.

**Reasoning / Evidence:** Slide-carousel tutorials are skipped by most users and retained by few (recall of front-loaded instructions ≈ nil — [02-cognitive-foundations.md](02-cognitive-foundations.md#working-memory-and-chunking)); do-first with contextual hints outperforms in activation metrics. Front-loaded permission barrages (notifications+location+camera on launch) tank both grant rates and trust.

**When to use:**
- Minimal path: launch → (auth if truly required — social/passkey fast paths) → core surface with guided empty state ([04-interaction-design.md](04-interaction-design.md#empty-states)) → first-success moment ASAP.
- Contextual permission priming: ask camera when they tap scan; location when they tap near-me ([13-platform-web.md](13-platform-web.md#notification-permission-ux) doctrine).
- Personalization quizzes: only if they *visibly* shape the immediate experience; always skippable.

**When NOT to use / exceptions:**
- Complex/regulated products (banking KYC) have mandated steps — design them as one clear checklist with progress and save-resume ([07-forms-and-input.md](07-forms-and-input.md#multi-step-vs-single-page-forms)).
- A 1–3 screen value proposition *for genuinely novel products* can earn its place — test skip rates honestly.

**Pros (value-first):** Activation, permission grant rates, day-1 retention.
**Cons:** Guided-empty-state and contextual-hint engineering exceeds slide-deck effort — worth it.

**Implementation guidance:**
- Measure: time-to-first-value, step-level drop-off, D1/D7 retention per onboarding variant ([17-ux-process-research.md](17-ux-process-research.md#analytics-and-behavioral-data)).
- Defer sign-up where possible (browse first, auth at save/purchase) — every pre-value screen sheds users.

**Related:** [14-platform-desktop.md](14-platform-desktop.md#first-run-experience), [09-content-ux-writing.md](09-content-ux-writing.md#onboarding-copy-and-tooltips), [#notifications-ux](#notifications-ux).

**Sources:** [NN/g: Mobile-app onboarding](https://www.nngroup.com/articles/mobile-app-onboarding/).

---

## Resilience and performance

### Offline and poor connectivity

**Definition:** Designing for the network reality of mobile: offline states that inform without blocking, cached content available always, queued mutations with sync indication, and graceful degradation across the connectivity spectrum (offline → 2G-ish → flaky → fast).

**Reasoning / Evidence:** Mobile connectivity is *intermittent by physics* (elevators, transit, roaming, congestion); apps treating offline as an error state (blocking spinners, data loss on submit) fail users at predictable daily moments. Offline-first architectures (local store + background sync) turn the failure mode into a non-event ([14-platform-desktop.md](14-platform-desktop.md#offline-first-and-autosave) shares doctrine).

**When to use:**
- Read paths: cache-first display with freshness indication ("Updated 2h ago"); stale beats blank.
- Write paths: queue locally, optimistic UI with pending markers ([04-interaction-design.md](04-interaction-design.md#optimistic-ui)), background sync on reconnect, conflict strategy defined.
- Status: subtle persistent offline banner; per-item sync badges where items queue; never modal "No connection" walls over usable cached content.

**When NOT to use / exceptions:**
- Genuinely realtime features (calls, live bidding) degrade to honest unavailable states with retry.
- Sensitive cached data: respect security requirements (lock cached financial data behind biometric re-auth).

**Pros:** Reliability perception ("just works on the subway") is a retention differentiator.
**Cons:** Sync/conflict engineering; cache-invalidation correctness; testing discipline (network-condition simulation must be routine, not launch-week).

**Implementation guidance:**
- Test matrix: airplane mode mid-flow, flaky (packet loss) profiles, slow (400ms RTT) profiles — on every core flow.
- Retry: automatic with exponential backoff + manual retry affordance; failed queued items surfaced actionably, never silently dropped.

**Related:** [04-interaction-design.md](04-interaction-design.md#optimistic-ui), [13-platform-web.md](13-platform-web.md#pwa-capabilities), [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery).

**Sources:** [Google: offline UX design guidelines](https://developer.android.com/design/ui/mobile/guides/patterns/offline), [Ink & Switch: local-first](https://www.inkandswitch.com/local-first/).

### Performance perception on mobile

**Definition:** Mobile-specific performance UX: cold-start time (target ≤2s to interactive content; OS-flagged if ≥5s), 60fps scroll (dropped frames are *felt* in touch-tracking immediately), touch-response latency (<100ms), and battery/data respect as UX dimensions.

**Reasoning / Evidence:** App-start abandonment mirrors web-load abandonment; scroll jank breaks the direct-manipulation illusion ([04-interaction-design.md](04-interaction-design.md#direct-manipulation)) more harshly than on desktop because the finger *is* the scrollbar. Android vitals/App Store both surface slow-start and jank metrics to stores — performance is discoverability too.

**When to use:**
- Cold start: splash-to-content ≤2s; skeleton the first screen ([04-interaction-design.md](04-interaction-design.md#loading-strategies)); defer non-critical init.
- Lists: virtualization, image thumbnailing/lazy decode, stable item heights (mobile CLS equivalent).
- Test on mid-range Android hardware (the global median device, not your flagship) under battery saver.

**When NOT to use / exceptions:**
- Splash screens beyond the OS-provided launch frame: brand moments >1s are self-indulgence ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#splash-screens)).

**Implementation guidance:**
- Budgets in CI: startup time, frame-drop rate, APK/IPA size (download friction), memory ceilings.
- Data respect: image quality tiers by connection class; user-facing data-saver options for media-heavy apps.

**Related:** [04-interaction-design.md](04-interaction-design.md#feedback-and-response-time-thresholds), [13-platform-web.md](13-platform-web.md#core-web-vitals).

**Sources:** [Android vitals](https://developer.android.com/topic/performance/vitals), [Apple: reducing launch time](https://developer.apple.com/documentation/xcode/reducing-your-app-s-launch-time).

---

## Beyond the phone

### Adaptive layouts: tablets and foldables

**Definition:** Scaling phone UIs up responsibly: Material **window size classes** (compact <600dp / medium 600–840dp / expanded >840dp width) and iPadOS multitasking (Split View, Stage Manager, resizable scenes) demand layout adaptation — list-detail panes, navigation rail replacing bottom bar, multi-column grids — not stretched phone layouts. Terminology: **responsive** design adapts *fluidly and continuously* to available space (flexible grids, reflowing content); **adaptive** design switches between *discrete* layouts or behaviors at known thresholds (size classes, posture changes). Mobile platform practice is primarily adaptive (size-class switches) built on responsive content within each layout.

**Reasoning / Evidence:** Stretched-phone-UI on tablets (one column of 900px-wide list rows) is the canonical adaptive failure — wasted space, absurd reach distances. Foldables add postures (folded/unfolded transitions mid-session, hinge-aware layouts); iPad multitasking means your "tablet" layout may live at phone width — size classes, not device detection, are the only sane model. The same logic generalizes: detect capabilities and available space, never assume viewport width equals device type.

**When to use:**
- Expanded widths: list-detail (master-detail) two-pane, navigation rail (left) instead of bottom bar, dialogs as centered sheets not full-screen, content max-widths for readability ([03-visual-design.md](03-visual-design.md#line-height-and-line-length)).
- Continuity: fold/unfold and window-resize preserve state and scroll position seamlessly.
- Use responsive techniques (fluid grids, flexible content) *within* each size class; use adaptive switches (pane count, navigation morph) *between* them.

**When NOT to use / exceptions:**
- Phone-only products can defer — but honest window-size-class plumbing from day one is cheap; retrofits aren't.

**Pros:** Tablet/foldable users are high-engagement segments; store featuring favors adaptive quality.
**Cons:** Layout matrix growth; per-posture QA.

**Implementation guidance:**
- Breakpoints by size class, never by device model; test: phone, unfolded foldable (~673dp), tablet split-view (~half-width), full tablet.
- Navigation morphs: bottom bar (compact) ↔ rail (medium+) ↔ drawer (expanded, optionally) — same destinations, same order.
- Localization must survive breakpoint changes: RTL locales mirror layouts (navigation rail flips sides, back gestures/icons mirror), translated strings expand (German/Finnish +30–40% is routine) — test that adaptive switches don't truncate or overflow at expanded text, and locale-aware date/number formats persist across every layout variant ([09-content-ux-writing.md](09-content-ux-writing.md)).

**Verification checklist:**
- [ ] Layouts keyed to window size classes (600/840dp), not device model or "is tablet" checks.
- [ ] Fold/unfold, rotation, and split-screen resize preserve state and scroll position.
- [ ] Navigation destinations identical (same set, same order) across bar/rail/drawer morphs.
- [ ] RTL mirroring and +40% text expansion tested at each size class.

**Related:** [13-platform-web.md](13-platform-web.md#responsive-design), [08-navigation-ia.md](08-navigation-ia.md), [14-platform-desktop.md](14-platform-desktop.md#density-and-information-rich-layouts), [#orientation-and-safe-areas](#orientation-and-safe-areas), [09-content-ux-writing.md](09-content-ux-writing.md).

**Sources:** [Material: window size classes](https://m3.material.io/foundations/layout/applying-layout/window-size-classes), [Android: large screens](https://developer.android.com/guide/topics/large-screens/get-started-with-large-screens), [Apple HIG: iPadOS](https://developer.apple.com/design/human-interface-guidelines/designing-for-ipados).

### Dark mode on mobile

**Definition:** Mobile dark-theme specifics on top of the general doctrine ([03-visual-design.md](03-visual-design.md#dark-mode-design)): follow the system setting by default (both OSes schedule it), offer in-app override (light/dark/system triple), OLED battery reality, and map correctly through platform theming (M3 tonal roles / iOS semantic colors).

**Reasoning / Evidence:** System-scheduled dark mode means apps ignoring the setting glare at users nightly; platform semantic color systems (iOS dynamic colors, M3 roles) make compliance nearly free *if adopted* — hardcoded hexes are the failure source. OLED savings are real with true-dark surfaces at high brightness.

**When to use:** All apps; semantic tokens throughout; test both themes per screen (including images/illustrations/map styles).

**Implementation guidance:**
- Default: system-follow; persisted in-app override for the minority who split preferences per app.
- Use platform semantic colors/elevation systems; audit every hardcoded color; dark-variant assets where needed ([03-visual-design.md](03-visual-design.md#dark-mode-design) for surface/contrast numbers).

**Related:** [03-visual-design.md](03-visual-design.md#dark-mode-design), [10-design-systems.md](10-design-systems.md#theming-dark-mode-multi-brand-density).

**Sources:** [Material: dark theme](https://m2.material.io/design/color/dark-theme.html), [Apple HIG: Dark Mode](https://developer.apple.com/design/human-interface-guidelines/dark-mode).

### Mobile web vs native app

**Definition:** The channel decision: mobile web/PWA (zero install friction, linkable, one codebase) vs native app (capability depth, home-screen presence, store distribution, notification reliability) vs both (web for acquisition/long-tail, app for retention/power) — the dominant mature-product answer.

**Reasoning / Evidence:** Install friction is enormous (each step from tap to first-open sheds users) — casual/first-time intent belongs to web; engaged repeat use rewards the app's speed/offline/notification advantages. Forcing app installs on web visitors (interstitial walls) measurably destroys acquisition and is search-penalized ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#app-interstitials)).

**When to use (web/PWA-first):**
- Acquisition surfaces, content consumption, infrequent transactions (annual bookings), emerging-market data-sensitive audiences ([13-platform-web.md](13-platform-web.md#pwa-capabilities)).
**When to use (native):**
- Daily-use products, deep hardware (camera pipelines, sensors, background location), notification-critical services, offline-heavy workflows.

**When NOT to use / exceptions:**
- Never gate web content behind install walls; smart banners (dismissible, small) are the acceptable cross-promotion ceiling.
- Don't ship a webview-wrapper "app" that's just the site minus browser features — it earns 1-star honesty.

**Pros (both-channel):** Full-funnel coverage; each channel at its strength.
**Cons:** Two surfaces to maintain (mitigate: shared design system — [10-design-systems.md](10-design-systems.md), shared APIs; feature-parity discipline for core tasks).

**Implementation guidance:**
- Parity rule: any task a user *started* on one channel must be completable there (no "finish in the app" cliffs for core flows).
- Hand-off: web→app via deep links preserving context ([#deep-links](#deep-links)); app→web share URLs canonical ([13-platform-web.md](13-platform-web.md#urls-as-ux)).

**Related:** [13-platform-web.md](13-platform-web.md#pwa-capabilities), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#app-interstitials).

**Sources:** [Google: avoid intrusive interstitials](https://developers.google.com/search/docs/appearance/avoid-intrusive-interstitials), [web.dev: PWA vs native considerations](https://web.dev/learn/pwa/).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Mobile context | Interruption-proof, one job per screen | Sessions 30–90s bursts |
| Thumb zone | Primary actions bottom third | ~49% grips one-handed; ~75% of touches thumb-driven (Hoober) |
| Touch targets | Hit areas beat visuals | 44pt / 48dp / 24px floor; 8dp gaps |
| Orientation & safe areas | Don't lock unless essential; inset from notches/gesture bar | WCAG 1.3.4 |
| Material 3 | Tonal roles, dynamic color, standard components | Nav bar 3–5, 80dp |
| Apple HIG | Tab bars, swipe-back sacred, Dynamic Type | Tab bar 3–5, 49pt |
| Cross-platform | Unify brand, fork behavior | Back handling = fork #1 |
| Navigation | Bottom tabs with labels; drawer for long tail | Re-tap = top/root |
| Back (Android) | Predictable stack; synthesize from deep links | Predictive back support |
| Back (iOS) | Never block edge swipe-back | — |
| Deep links | Every entity; back stack synthesized | Universal/App Links, not schemes |
| Gestures | Standard set + visible equivalents; edges are OS-owned | — |
| Forms | Right keyboard, autofill everything, biometrics | Inputs ≥16px |
| Haptics | Semantic APIs, meaningful moments only | One meaning = one pattern |
| Notifications | Contextual permission, honest channels, budgets | Never ask on launch |
| Onboarding | Value first, ask in context, skippable | Time-to-first-value metric |
| Offline | Cache-first reads, queued writes, stale > blank | Test airplane mid-flow |
| Performance | Fast start, 60fps scroll, mid-range test devices | Cold start ≤2s; touch <100ms |
| Adaptive | Size classes, list-detail, rail | 600/840dp breakpoints |
| Responsive vs adaptive | Responsive = fluid within layouts; adaptive = discrete switches between them | Never infer device from width |
| Dark mode | System-follow + override; semantic tokens | — |
| Web vs native | Web acquires, app retains; never install-wall | Parity for core tasks |
