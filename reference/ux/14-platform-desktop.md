# Platform: Desktop

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers UX for desktop applications (Windows, macOS, Linux): what makes desktop distinct, per-platform conventions and their comparison, keyboard shortcuts, menus, window management, information density, native vs cross-platform frameworks, system integration, autosave/offline expectations, update UX, and first-run experience. Generic interaction rules are in [04-interaction-design.md](04-interaction-design.md); density and visual language in [03-visual-design.md](03-visual-design.md); web-app-in-desktop concerns overlap with [13-platform-web.md](13-platform-web.md).

## Contents

- [The desktop context](#the-desktop-context)
  - [How desktop differs](#how-desktop-differs)
- [Platform conventions](#platform-conventions)
  - [Windows conventions](#windows-conventions)
  - [macOS conventions](#macos-conventions)
  - [Linux conventions](#linux-conventions)
  - [Cross-platform convention mapping](#cross-platform-convention-mapping)
- [Input and commands](#input-and-commands)
  - [Keyboard shortcuts design](#keyboard-shortcuts-design)
  - [Menus](#menus)
  - [Context menus](#context-menus)
  - [Command palettes](#command-palettes)
- [Windows and workspace](#windows-and-workspace)
  - [Window management and state restoration](#window-management-and-state-restoration)
  - [Modality](#modality)
  - [Density and information-rich layouts](#density-and-information-rich-layouts)
  - [Desktop drag and drop](#desktop-drag-and-drop)
- [Platform integration](#platform-integration)
  - [Native vs custom chrome and cross-platform frameworks](#native-vs-custom-chrome-and-cross-platform-frameworks)
  - [System integration](#system-integration)
  - [High-DPI and multi-monitor](#high-dpi-and-multi-monitor)
  - [Touch on hybrid devices](#touch-on-hybrid-devices)
- [Lifecycle](#lifecycle)
  - [Offline-first and autosave](#offline-first-and-autosave)
  - [Update UX](#update-ux)
  - [First-run experience](#first-run-experience)
- [Quick reference](#quick-reference)

---

## The desktop context

### How desktop differs

**Definition:** Desktop usage patterns diverge from web/mobile: long focused sessions, keyboard-primary expert input, multiple simultaneous windows/monitors, dense professional data, files and local resources, persistent installation (users chose you deliberately), and multi-year retention.

**Reasoning / Evidence:** Desktop users are disproportionately *committed* users: they installed, they return daily, they build muscle memory and workflows around the tool. This flips several web priorities: learnability still matters, but *efficiency ceiling* matters more ([01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use)); discoverability chrome can be lighter because sessions are long enough to learn; interruption cost is higher because sessions are deep ([02-cognitive-foundations.md](02-cognitive-foundations.md#flow-state)).

**When to use (desktop-first thinking):**
- Professional/creation tools: optimize for the 100th session, not just the 1st — shortcuts, density options, customization, bulk operations.
- Respect the OS: users chose a platform and expect apps to belong to it (conventions below).

**When NOT to use / exceptions:**
- Utility apps used rarely (installers, settings panels) should optimize first-use clarity like web pages.
- Don't ship a stretched phone UI to desktop (sparse, tap-target-huge, single-column) — it wastes the platform.

**Pros (desktop-native strengths):** Keyboard mastery, window workflows, file integration, offline capability, performance headroom.
**Cons:** Higher user expectations (stability, data safety, OS fit); install friction means fewer casual trials; platform fragmentation costs.

**Implementation guidance:**
- Session-continuity budget: restore exactly where the user left off (windows, tabs, scroll, selection) — desktop users end sessions by OS shutdown, not by "finishing."
- Every frequent action: mouse path + keyboard path + (where sensible) menu path.
- Test with real workflows: 8-hour-session concerns (memory leaks, notification fatigue, focus stealing) don't show in demos.

**Verification checklist (simple by default, powerful when needed):**
- [ ] Default UI exposes common tasks clearly, without requiring configuration first.
- [ ] Advanced capabilities remain reachable through menus, shortcuts, command palette, preferences, or progressive disclosure — layered, not removed.
- [ ] Simplifying the novice view has not stripped expert efficiency (shortcuts, bulk actions, density options still exist).
- [ ] Power features that would overwhelm default tasks are not all exposed at top level.

**Related:** [01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use), [#window-management-and-state-restoration](#window-management-and-state-restoration), [03-visual-design.md](03-visual-design.md#density-modes).

**Sources:** [NN/g: Complex Application Design](https://www.nngroup.com/articles/complex-application-design/), [macOS HIG](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos), [Windows design basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/).

---

## Platform conventions

### Windows conventions

**Definition:** Fluent 2-era Windows 11 conventions: title bar with right-aligned minimize/maximize/close, optional Mica/Acrylic materials, Segoe UI Variable type, taskbar presence with jump lists, system tray for background agents, snap layouts participation, in-window menu bar or command bar (ribbon for document-heavy suites), Alt-key mnemonics, per-monitor DPI awareness.

**Reasoning / Evidence:** Microsoft's Fluent 2 design system defines current guidance; users' muscle memory spans decades (Alt+F4, right-click depth, window snapping). Windows tolerates both menu-bar and command-bar/ribbon apps; consistency with *category* leaders matters as much as OS guidance.

**When to use:**
- Follow always: caption buttons at top-right (never restyled beyond recognition), Ctrl-based shortcuts, right-click context menus everywhere, snap-layout compatibility (don't block Win+arrow / hover-maximize layouts), jump lists for recent files/actions, dark/light theme + accent color respect.
- System tray: only for genuinely background-resident apps (sync clients, comms) — with quit truly quitting when asked.

**When NOT to use / exceptions:**
- Ribbons: only for command-rich document apps (Office-class); overkill elsewhere — Windows 11 direction favors simpler command bars.
- Don't clone macOS idioms (global top menu bar, Cmd metaphors, traffic lights left) on Windows.

**Pros:** Belonging, muscle-memory transfer, OS feature participation (snap, jump lists) for free engagement.
**Cons:** Fluent churn (UWP→WinUI→…) makes "native" a moving target; Win10/Win11 visual split.

**Implementation guidance:**
- Standard keys: Ctrl+C/V/X/Z/Y/S/F/N/W, F2 rename, F5 refresh, Alt+F4 quit, Ctrl+Tab view-switch; Alt-mnemonics in menus (underlined access keys).
- Dialogs: historically OK-left/Cancel-right; Windows 11 accepts trailing-primary — pick the platform-current pattern, keep it consistent, and never swap Cancel/confirm positions between dialogs.
- Title bar customization: keep drag regions, double-click-maximize, and system menu (Alt+Space) working.

**Related:** [#cross-platform-convention-mapping](#cross-platform-convention-mapping), [#keyboard-shortcuts-design](#keyboard-shortcuts-design).

**Sources:** [Windows design principles](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/design-principles), [Fluent 2](https://fluent2.microsoft.design/), [Windows app design basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/).

### macOS conventions

**Definition:** Apple HIG conventions for Mac: the **global menu bar at screen top** (every app must populate it fully), traffic-light window controls top-left, Cmd-based shortcuts, Dock integration, sheets (window-attached dialogs), popovers, unified/inline toolbars, sidebar source lists, SF Pro type and SF Symbols icons, Settings window (not dialog) with standard placement, system-wide dark mode/accent/transparency compliance.

**Reasoning / Evidence:** Mac users are convention-sensitive to a degree that surprises cross-platform teams: apps without a real menu bar, with wrong Cmd shortcuts, or with Windows-style in-window menus get review-bombed as "not a Mac app." The menu bar doubles as documentation — users search it (Help menu search searches *menu items*), so every command belongs there even if it also lives in a toolbar.

**When to use:**
- Always: full menu bar (App/File/Edit/View/Window/Help minimum) with every command + shortcut listed; Cmd+, for Settings; Cmd+Q quits, Cmd+W closes window (app may stay running — the Mac lifecycle difference); traffic lights untouched; sheets for document-scoped dialogs.
- Standard idioms: sidebar navigation with translucency, toolbar in the title-bar region, drag-and-drop richness (files in/out, text, images), Services/Share menu participation where relevant.

**When NOT to use / exceptions:**
- Don't replicate a ribbon or an in-window hamburger menu bar; don't quit-on-last-window-close unless the app is genuinely single-window utility class.

**Pros:** Menu bar = free discoverability + shortcut teaching; strong idioms mean strong defaults.
**Cons:** Idiom richness = more to implement correctly; cross-platform frameworks approximate it poorly by default (Electron apps must hand-build a proper menu).

**Implementation guidance:**
- Standard keys: Cmd+C/V/X/Z (Shift+Cmd+Z redo)/S/F/N/W/Q/T, Cmd+, settings, Space quick-look-style preview where apt; full keyboard navigation compliance (Tab behavior per system setting).
- Dialog buttons: primary at bottom-**right**, Cancel to its left; destructive actions separated; sheet attachment for document-modal.
- Respect: system appearance (light/dark/auto), accent color, reduced transparency/motion accessibility settings ([05-accessibility.md](05-accessibility.md#reduced-motion)).

**Related:** [#cross-platform-convention-mapping](#cross-platform-convention-mapping), [#menus](#menus).

**Sources:** [Apple HIG: macOS](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos), [Apple HIG: Menus](https://developer.apple.com/design/human-interface-guidelines/menus).

### Linux conventions

**Definition:** The fragmented-but-navigable Linux desktop: GNOME HIG (header bars merging title+toolbar, minimal chrome, adaptive layouts, no menu bar — primary menu button instead) vs KDE/Plasma (Windows-adjacent: menu bars, dense options, high configurability), plus the packaging layer (Flatpak/Snap/AppImage/native) shaping integration.

**Reasoning / Evidence:** No single HIG governs; the pragmatic strategy is: follow **GNOME HIG** for GNOME-targeting apps (largest default-desktop share via Ubuntu/Fedora), respect **freedesktop.org standards** (XDG dirs, desktop entries, notifications, tray protocols) for cross-desktop correctness, and let toolkit choice (GTK vs Qt) align with the primary target.

**When to use:**
- GNOME-target: header bar with window controls per user setting, primary menu (hamburger) top-right, `Ctrl+,`-style shortcuts increasingly adopted, adaptive width tolerance (libadwaita patterns).
- KDE-target/cross-desktop: traditional menu bar acceptable; respect global themes where feasible.
- Always: XDG base directories (config in `~/.config`, not dotfile sprawl), `.desktop` entries with actions, standard notification portal, honor system dark mode preference (portal), Flatpak-friendly sandboxing.

**When NOT to use / exceptions:**
- Don't hard-code one DE's assumptions (tray icon existence varies; test without).
- Electron/cross-platform apps: at minimum theme-follow (dark mode), use portals for file dialogs (native pickers), and package as Flatpak for the mainstream path.

**Pros:** Technical users forgiving of chrome differences but demanding about *behavioral* correctness (shortcuts, files, terminal-friendliness).
**Cons:** Fragmentation tax: two design languages (GNOME/KDE), N packaging formats; small market share for consumer apps (but overweight in developer tools).

**Implementation guidance:**
- Shortcuts: Ctrl-based, Windows-adjacent set; middle-click paste tolerance; respect `$XDG_CONFIG_HOME` ([16-platform-cli-tui.md](16-platform-cli-tui.md#configuration-precedence) shares this).
- Test matrix: GNOME (Wayland), KDE, one tiling WM (developer audiences love them — don't fight unusual window sizes).

**Related:** [#cross-platform-convention-mapping](#cross-platform-convention-mapping), [16-platform-cli-tui.md](16-platform-cli-tui.md#configuration-precedence).

**Sources:** [GNOME HIG](https://developer.gnome.org/hig/), [KDE HIG](https://develop.kde.org/hig/), [freedesktop.org standards](https://www.freedesktop.org/wiki/Specifications/).

### Cross-platform convention mapping

**Definition:** The translation table a cross-platform desktop app must implement per-OS rather than shipping one platform's habits everywhere.

**Reasoning / Evidence:** [Consistency with the platform beats consistency with yourself](01-core-principles.md#4-consistency-and-standards) — users compare your app to their OS's other apps, not to your other builds.

| Concern | Windows | macOS | Linux (GNOME) |
|---|---|---|---|
| Primary modifier | Ctrl | Cmd (⌘) | Ctrl |
| Quit app | Alt+F4 / Ctrl+W closes | Cmd+Q (Cmd+W = window only) | Ctrl+Q |
| Settings location | Hamburger/File menu → Settings | App menu → Settings… (Cmd+,) | Primary menu → Preferences |
| Menu style | In-window bar / command bar | Global top menu bar (mandatory) | Header-bar primary menu |
| Window controls | Right: min/max/close | Left: traffic lights | Right (user-configurable) |
| Dialog primary button | Trailing (right) accepted; legacy left | Bottom-right | Trailing (right) |
| Redo | Ctrl+Y (also Ctrl+Shift+Z) | Shift+Cmd+Z | Ctrl+Shift+Z |
| Find/Find-next | Ctrl+F / F3 | Cmd+F / Cmd+G | Ctrl+F / Ctrl+G |
| App lifecycle | Last window close = exit (typical) | App outlives windows | Last window close = exit (typical) |
| File dialogs | System picker | System sheet/panel | XDG portal picker |
| Rename (file lists) | F2 | Return | F2 |

**When to use:** Every cross-platform build: per-OS keymaps, menu structures, dialog button orders, and lifecycle behavior — not a compile flag afterthought.

**When NOT to use / exceptions:** In-content interactions (canvas tools, custom widgets) can be identical everywhere; it's the *chrome and system-touching* layer that must localize.

**Pros:** Each audience gets a native-feeling app. **Cons:** Test matrix ×3; per-platform QA of every shortcut/dialog.

**Implementation guidance:**
- Abstraction layer for: accelerator display (⌘ vs Ctrl+ rendering), menu templates, dialog button order, settings-window style.
- Document per-platform deltas in the design system ([10-design-systems.md](10-design-systems.md)); test with platform-native users, who spot violations instantly.

**Related:** [#windows-conventions](#windows-conventions), [#macos-conventions](#macos-conventions), [#linux-conventions](#linux-conventions).

**Sources:** [Apple HIG](https://developer.apple.com/design/human-interface-guidelines/), [Windows design](https://learn.microsoft.com/en-us/windows/apps/design/basics/), [GNOME HIG](https://developer.gnome.org/hig/).

---

## Input and commands

### Keyboard shortcuts design

**Definition:** The layered shortcut system of a serious desktop app: universal standards (copy/save/find), category conventions (your competitors' de-facto set), app-specific accelerators, and user customization — all discoverable through menus, tooltips, and a shortcut reference.

**Reasoning / Evidence:** Shortcuts are the highest-leverage expert accelerator ([01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use)); discoverability is the bottleneck — menus displaying shortcuts beside commands are the classic teaching mechanism (each menu use is a lesson). Muscle memory transfers across apps, so violating standards (Ctrl+S doing something weird) actively causes errors.

**When to use:**
- Standard set first, never repurposed: Ctrl/Cmd + C/V/X/Z/S/F/N/W/A/P; platform redo variant; Esc cancels; Enter confirms; F2/Return rename; Delete deletes (with undo — [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)).
- Category conventions: study the tools your users come from (video editors expect J/K/L; IDEs expect Ctrl+P/Cmd+P quick-open) — adopt, don't innovate, on the basics.
- Every menu item shows its shortcut; tooltips on toolbar buttons show them; a searchable shortcut sheet (often `?` or Ctrl+/) lists all.
- Customization for pro tools: remappable keys with conflict detection, preset schemes ("VS Code keys", "Vim mode") for migrating users.

**When NOT to use / exceptions:**
- Don't consume single letters globally while text inputs exist (see WCAG 2.1.4); don't override OS-reserved combos (Cmd+Tab, Win+L, Alt+Tab).
- Chorded/multi-stroke shortcuts (Ctrl+K Ctrl+S) only in developer-class tools where the audience accepts them.

**Pros:** Order-of-magnitude speed for experts; retention (muscle memory = switching cost in your favor).
**Cons:** Conflict management complexity; international keyboard layouts break punctuation-based combos (test AZERTY/QWERTZ — Ctrl+/ may not exist).

**Implementation guidance:**
- Display: platform-correct rendering (⌘⇧Z vs Ctrl+Shift+Z); in menus right-aligned.
- Scope layers: global (app-wide) < view < focused-widget; innermost wins; document the model.
- Log shortcut usage: unused shortcuts on frequent commands = discoverability failure → surface hints ("Tip: Ctrl+K opens search") sparingly.

**Verification checklist:**
- [ ] Standard shortcuts (copy/paste/cut/undo/redo/save/find/new/close/select-all/print) work with their platform-conventional meanings — desktop users expect them to just work.
- [ ] Every shortcut action has a visible menu path (nothing is shortcut-only).
- [ ] Frequent commands have keyboard accelerators, shown beside their menu items and in toolbar tooltips.
- [ ] No app shortcut overrides an OS-reserved combination.
- [ ] Shortcuts verified per-OS (Ctrl vs Cmd mapping, redo variant) and on non-US keyboard layouts.

**Related:** [04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns), [#menus](#menus), [#command-palettes](#command-palettes).

**Sources:** [Apple HIG: Keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards), [Windows: Keyboard accelerators](https://learn.microsoft.com/en-us/windows/apps/design/input/keyboard-accelerators).

### Menus

**Definition:** The command-organization backbone: menu bar (global on macOS, in-window on Windows/Linux) or primary-menu button, organized File/Edit/View/…-style, holding *every* command with its shortcut — the app's complete, searchable command inventory.

**Reasoning / Evidence:** Menus serve recognition ([01-core-principles.md](01-core-principles.md#6-recognition-rather-than-recall)): rarely-used commands don't need toolbar slots because the menu holds them findable. macOS Help-menu search searches menus — instant command findability. Apps that hid menus behind minimal chrome (post-2010 trend) traded discoverability for looks; command palettes ([#command-palettes](#command-palettes)) are the modern compensation, not a replacement.

**When to use:**
- Any app with >~15 commands: full menu structure, standard top-level names (File, Edit, View, Window, Help), standard item placement (Quit/Exit at File's bottom or app menu; Undo/Redo top of Edit).
- Toolbar = frequency-based *subset* of menu commands, user-customizable in pro tools.

**When NOT to use / exceptions:**
- Single-purpose utilities (<10 commands) can live with a primary-menu button (GNOME style) or toolbar-only.
- Never menu-only for frequent actions (menus are slow paths — Fitts-distant, two-level) — frequent actions also get toolbar/shortcut/context-menu presence.

**Pros:** Complete inventory, shortcut teaching, standardized findability.
**Cons:** Menu real estate invites bloat — curate; deep nesting (>2 levels) hides ([08-navigation-ia.md](08-navigation-ia.md#multi-level-menu-ergonomics-hover-intent-and-safe-triangles)).

**Implementation guidance:**
- Group with separators by task cluster; ≤~9 items per group; ellipsis (…) on items opening dialogs; checkmarks/radio marks for state items; disable-with-reason over hiding (menu items that vanish confuse; disabled items teach "exists, not now").
- Label items verb-first ("Export as PDF…"), consistent with buttons elsewhere ([09-content-ux-writing.md](09-content-ux-writing.md#button-labels)).

**Related:** [#keyboard-shortcuts-design](#keyboard-shortcuts-design), [#context-menus](#context-menus), [08-navigation-ia.md](08-navigation-ia.md).

**Sources:** [Apple HIG: Menus](https://developer.apple.com/design/human-interface-guidelines/menus), [GNOME HIG: Menus](https://developer.gnome.org/hig/patterns/controls/menus.html).

### Context menus

**Definition:** Right-click (or long-press/menu-key) menus offering actions relevant to the clicked object — the desktop's ubiquitous object-action accelerator.

**Reasoning / Evidence:** Context menus minimize Fitts distance (menu at cursor — [01-core-principles.md](01-core-principles.md#fittss-law)) and match object-oriented thinking ("act on this thing"). Desktop users right-click *speculatively* to learn what's possible — an empty or missing context menu reads as a broken app. Iron rule: context menus are accelerators, **never the sole access path** (invisible, undiscoverable by definition; inaccessible to some inputs).

**When to use:**
- Every meaningful object (files, rows, tabs, cards, canvas elements, text selections): 3–12 relevant actions, most-used first, destructive last and separated.
- Match selection: right-click on unselected item acts on it (select-then-menu); on selection acts on all selected.

**When NOT to use / exceptions:**
- Never web-page-hijack right-click without content justification (canvas/editor apps justify it; content sites don't).
- Keyboard path required: Menu key / Shift+F10 must open it on the focused item ([05-accessibility.md](05-accessibility.md#keyboard-access-and-no-keyboard-trap)).

**Pros:** Speed, spatial efficiency, expert affordance.
**Cons:** Zero discoverability for novices (hence: duplicate everything visible somewhere); per-object-type menu maintenance.

**Implementation guidance:**
- Order template: primary action (bold/first, = double-click action), object verbs, clipboard verbs, "Rename/Duplicate", separator, "Delete", separator, "Properties/Details".
- Icons sparingly (destructive + a few anchors); shortcuts displayed for items that have them.

**Related:** [#menus](#menus), [01-core-principles.md](01-core-principles.md#fittss-law), [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures).

**Sources:** [Apple HIG: Context menus](https://developer.apple.com/design/human-interface-guidelines/context-menus), [Windows: Context menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus-and-context-menus).

### Command palettes

**Definition:** The keyboard-summoned fuzzy-search command launcher (Ctrl/Cmd+K or Ctrl/Cmd+Shift+P): type to filter all commands/destinations/objects, Enter to execute — recognition-assisted recall for expert speed.

**Reasoning / Evidence:** Popularized by Sublime/VS Code, now near-universal in productivity software (Figma, Linear, Slack, GitHub). It resolves the menus-vs-shortcuts tension: faster than menu mousing, no memorization like shortcuts, self-documenting (results show the shortcut), and it scales to hundreds of commands where toolbars can't ([02-cognitive-foundations.md](02-cognitive-foundations.md#recognition-vs-recall) hybrid).

**When to use:**
- Any command-rich app (>~30 commands) with a keyboard-comfortable audience; also as navigation (jump to file/project/person).
- As the discoverability layer for minimal-chrome designs — but paired with menus, not replacing them (palettes require knowing they exist).

**When NOT to use / exceptions:**
- Casual/consumer apps with few commands: unnecessary machinery.
- Not an excuse to skip information architecture — browsing must still work ([08-navigation-ia.md](08-navigation-ia.md)).

**Pros:** Expert throughput, zero-real-estate command access, shortcut teaching, unified action+navigation search.
**Cons:** Invisible to novices (advertise it: hint in empty states/Help); ranking quality is the whole game (bad fuzzy match = abandoned feature).

**Implementation guidance:**
- Binding: Ctrl/Cmd+K (consumer convention) or Ctrl/Cmd+Shift+P (dev convention); Esc closes; arrows navigate; recents first on empty query.
- Fuzzy matching with camelCase/initialism support ("ofp" → "Open File Preview"); show category + shortcut per result; support parameterized commands (type-ahead into arguments).
- Full ARIA combobox semantics ([05-accessibility.md](05-accessibility.md#aria)).

**Related:** [#keyboard-shortcuts-design](#keyboard-shortcuts-design), [02-cognitive-foundations.md](02-cognitive-foundations.md#recognition-vs-recall).

**Sources:** [VS Code command palette](https://code.visualstudio.com/docs/getstarted/userinterface#_command-palette), [NN/g: UI accelerators](https://www.nngroup.com/articles/ui-accelerators/).

---

## Windows and workspace

### Window management and state restoration

**Definition:** Correct citizenship in the OS windowing system: remembering size/position/monitor per window, restoring full session state on relaunch, choosing among multi-window vs tabs vs split panes, and cooperating with OS snapping/tiling/full-screen.

**Reasoning / Evidence:** Desktop users build spatial workspaces (this app on that monitor at that size); apps that forget placement or open centered-default every launch destroy invested arrangement. Session restoration is table stakes post-Docs/VS Code: reopening = exactly where you were, including unsaved content ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)).

**When to use:**
- Persist per-window: size, position, monitor, maximized/full-screen state; validate against current monitor topology on restore (never restore off-screen — the classic disconnected-laptop bug).
- Session restore: open documents/tabs/views, selection, scroll, panel layout; instant-resume feel.
- Multi-window when: users compare/reference across displays (documents, chats, inspectors); tabs when: many peer documents in one workspace; split panes when: intra-window comparison — pro tools often need all three.

**When NOT to use / exceptions:**
- Single-instance enforcement only with strong reason (and then: focus existing window on re-launch, never silent nothing).
- Kiosk/utility dialogs may fix size — resizable is otherwise the default; content must reflow at small sizes ([#density-and-information-rich-layouts](#density-and-information-rich-layouts)).

**Pros:** Workspace trust; zero re-setup cost per session.
**Cons:** Multi-monitor edge cases (DPI change mid-drag, topology changes); state-file versioning across app updates.

**Implementation guidance:**
- Restore rules: intersect saved bounds with current displays; fall back to primary-display center at default size.
- Support OS features: Windows Snap (don't custom-frame it away), macOS full-screen spaces + Stage Manager tolerance, Linux tiling WMs (unusual aspect ratios must not break layout).
- New-window defaults: cascade or platform default; remember "last closed" as the template.

**Verification checklist (window state and multi-window):**
- [ ] Window size, position, monitor, and maximized state are restored on relaunch; saved bounds are validated against current topology (nothing restores off-screen).
- [ ] Previous workspace state (open documents/tabs, panel layout, selection, scroll) is restored where useful.
- [ ] In multi-window use, the active document/window scope is always clear — commands act on the window the user thinks they act on.
- [ ] Accidental edits in the wrong window are prevented (clear focus indication; per-window undo/edit scope).
- [ ] Unsaved changes are handled explicitly on window close and app quit (prompt or reliable autosave — never silent loss).
- [ ] OS window management (snap layouts, full-screen spaces, tiling WMs) works without fighting the app's chrome.

**Related:** [#how-desktop-differs](#how-desktop-differs), [#high-dpi-and-multi-monitor](#high-dpi-and-multi-monitor), [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery).

**Sources:** [Apple HIG: Windows](https://developer.apple.com/design/human-interface-guidelines/windows), [Windows: window sizing](https://learn.microsoft.com/en-us/windows/apps/design/layout/screen-sizes-and-breakpoints-for-responsive-design).

### Modality

**Definition:** The dialog-blocking spectrum: modeless (palettes, inspectors — never block), window-modal (sheets: block one document), app-modal (block the whole app), system-modal (rare, OS-level). Doctrine: use the *least* modality that guarantees the needed decision.

**Reasoning / Evidence:** App-modal dialogs freeze every window — hostile in multi-document apps (can't reference doc B to answer a question about doc A). macOS sheets exist precisely to scope modality to one window. Modeless-with-undo beats modal-confirm for most flows ([04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs), [02-cognitive-foundations.md](02-cognitive-foundations.md#habituation)).

**When to use:**
- Modeless: tool palettes, find bars, inspectors — let work continue.
- Window-modal (sheet): document-scoped decisions (save changes? export options).
- App-modal: genuinely app-wide decisions only (quit with unsaved changes across windows, licensing).

**When NOT to use / exceptions:**
- Never app-modal progress dialogs for backgroundable work (export can run while the user continues — progress in a status area).
- Focus-stealing across apps: never grab foreground for non-critical events ([02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost)) — badge/bounce/notify instead.

**Pros (minimal modality):** Multi-window workflows preserved; interruption cost contained.
**Cons:** Modeless state management is harder (dialog outcomes while the world changes).

**Implementation guidance:**
- Every modal: Esc = cancel, Enter = default (non-destructive) action, focus-trapped, returns focus on close ([04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns)).
- Sheets/dialog buttons per-platform order ([#cross-platform-convention-mapping](#cross-platform-convention-mapping)); verbs not Yes/No ([09-content-ux-writing.md](09-content-ux-writing.md#confirmation-and-destructive-action-copy)).

**Related:** [04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design), [18-patterns-antipatterns.md](18-patterns-antipatterns.md#modal-overuse).

**Sources:** [Apple HIG: Sheets](https://developer.apple.com/design/human-interface-guidelines/sheets), [NN/g: Modal & Nonmodal](https://www.nngroup.com/articles/modal-nonmodal-dialog/).

### Density and information-rich layouts

**Definition:** Desktop's data-dense workhorse patterns: master-detail splits, data grids, tree views, dockable/resizable panels, inspector sidebars — engineered for professionals scanning and manipulating volumes of information.

**Reasoning / Evidence:** NN/g complex-app research: experts *prefer* density (context per glance beats clicks per datum); the sin is not density but density *without hierarchy* ([03-visual-design.md](03-visual-design.md#density-modes)). Desktop's pointer precision and screen area make 24–32px rows, hover reveals, and multi-panel layouts viable in ways mobile cannot.

**When to use:**
- Master-detail: list/table left, detail right — the default record-work layout (select from a collection, inspect/edit one item); three-pane (nav/list/detail) for mail-class apps.
- Data grids: sortable columns, resize/reorder/hide columns (persisted per user), inline editing, multi-select + bulk bar, virtualized scrolling for large sets ([08-navigation-ia.md](12-data-tables-dashboards.md#sorting-filtering-and-column-management)).
- Resizable panes wherever workflows require simultaneous reference/comparison; dockable panels for pro tools (IDEs, DAWs, editors): draggable, collapsible, layout-persisted, reset-to-default escape hatch.

**When NOT to use / exceptions:**
- Consumer-facing desktop surfaces (onboarding, settings) stay comfortable-density.
- Panel-docking freedom without a reset button = users lost in their own layout — always ship "Reset layout."

**Pros:** Expert throughput; screen investment repaid.
**Cons:** Complexity budget: every panel/column option multiplies states; small-window degradation needs design (what collapses first?).

**Implementation guidance:**
- Rows 32–40px comfortable / 24–28px compact; density toggle persisted ([03-visual-design.md](03-visual-design.md#density-modes)).
- Resizable splits: min-widths per pane, double-click-divider to auto-fit, persisted positions.
- Keyboard: full grid navigation (arrows/Home/End/PageUp/Down, type-ahead), selection model (Shift/Ctrl-click + keyboard equivalents) ([04-interaction-design.md](04-interaction-design.md#keyboard-interaction-patterns)); multi-select and bulk actions designed carefully (clear selection count, undoable bulk operations); context preserved during multi-step editing.

**Related:** [03-visual-design.md](03-visual-design.md#density-modes), [08-navigation-ia.md](08-navigation-ia.md#navigation-for-data-heavy-uis), [12-data-tables-dashboards.md](12-data-tables-dashboards.md).

**Sources:** [NN/g: Complex Application Design](https://www.nngroup.com/articles/complex-application-design/).

### Desktop drag and drop

**Definition:** Desktop-grade drag-drop expectations beyond in-app reordering ([04-interaction-design.md](04-interaction-design.md#drag-and-drop)): files in from the OS (and out, to Finder/Explorer), content between apps (text, images, rich data), spring-loaded targets (hover-to-open folders/tabs), and modifier semantics (move vs copy vs link).

**Reasoning / Evidence:** OS-level drag-drop is core desktop grammar (Finder/Explorer taught it for decades); apps that accept files only via dialog feel amputated. Drag *out* (export by dragging a canvas object to the desktop) is a delight marker of platform-serious apps.

**When to use:**
- Accept file drops anywhere conceptually sensible (import zones = whole window with overlay indicator, not a tiny target).
- Provide drag-out for user-owned content (attachments, images, exports).
- Spring-loading: hovering a drag over a folder/tab/collapsed section opens it after ~500–800ms.

**When NOT to use / exceptions:**
- Every drop path needs a dialog/menu/keyboard equivalent ([05-accessibility.md](05-accessibility.md#dragging-alternatives)).
- Cross-app drag on Linux/Windows sandboxing (Flatpak, store apps) has real limits — degrade gracefully.

**Pros:** Workflow fluency with the whole OS. **Cons:** Per-platform data-format plumbing (pasteboard types, DataTransfer, MIME negotiation).

**Implementation guidance:**
- Whole-window drop overlay: "Drop to import" with accepted-types note; per-zone highlighting when multiple targets exist.
- Modifiers: platform defaults (Windows: Ctrl=copy, Shift=move; macOS: Option=copy) — show cursor badges.
- Also support clipboard equivalence: whatever drags should copy/paste too.

**Related:** [04-interaction-design.md](04-interaction-design.md#drag-and-drop), [#system-integration](#system-integration).

**Sources:** [Apple HIG: Drag and drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop), [MDN: HTML Drag and Drop](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API).

---

## Platform integration

### Native vs custom chrome and cross-platform frameworks

**Definition:** The strategic choice among native toolkits (WinUI, AppKit/SwiftUI, GTK/Qt), cross-platform native-rendering frameworks (Qt, Flutter, Compose Multiplatform), and web-wrapper frameworks (Electron, Tauri) — and how far to customize window chrome — with direct UX consequences: fidelity, performance, memory, and the uncanny valley of almost-native.

**Reasoning / Evidence:** Trade-off is real and measured: Electron ships Chromium per app (memory/disk footprint, but full web-stack velocity and consistency); Tauri uses system webviews (lighter, more variance); native maximizes fidelity and system integration at multiplied engineering cost. The *uncanny valley* problem: an app that looks 90% native but scrolls wrong, breaks shortcuts, or mis-renders text triggers stronger rejection than an honestly custom-designed one — users forgive "its own thing" more than "broken native."

**When to use (web-wrapper):**
- Web-parallel products (the app IS the web app), fast-iterating startups, developer tools where the audience tolerates footprint for features (VS Code proves the ceiling is high).
**When to use (native/Qt-class):**
- Performance-critical (media, CAD), deep OS integration needs, platform-flagship experiences, resource-constrained targets.

**When NOT to use / exceptions:**
- Don't fake platform chrome badly: either use real native chrome or design confidently custom (consistent everywhere) — half-imitations age worst.
- Whatever the stack, the *behavioral* layer must be per-platform correct ([#cross-platform-convention-mapping](#cross-platform-convention-mapping)): menus, shortcuts, dialogs, lifecycle — Electron does not excuse a missing macOS menu bar.

**Pros (custom chrome):** Brand consistency, cross-platform sameness, design freedom.
**Cons (custom chrome):** You now own: window dragging, snap integration, caption buttons, accessibility of the title bar, per-OS quirks — a permanent tax.

**Implementation guidance:**
- Checklist for any stack: platform menus ✓, standard shortcuts per-OS ✓, native file dialogs ✓, OS notifications ✓, dark-mode follow ✓, screen-reader tested ✓ ([05-accessibility.md](05-accessibility.md)), 60fps scroll ✓, cold start ≤2–3s.
- If custom title bar: keep double-click-maximize, drag regions generous, snap/traffic-light behaviors intact.

**Verification checklist (desktop accessibility, any stack):**
- [ ] Controls expose native accessibility APIs (UIA on Windows, NSAccessibility on macOS, AT-SPI on Linux) — prefer platform controls where possible; custom controls implement the semantics.
- [ ] OS accessibility settings respected: high contrast, reduced motion, text scaling, system keyboard navigation ([05-accessibility.md](05-accessibility.md)).
- [ ] Screen reader tested on each shipping platform.
- [ ] Visible focus indication never removed for aesthetic reasons — including in custom title bars and chrome.

**Related:** [#cross-platform-convention-mapping](#cross-platform-convention-mapping), [01-core-principles.md](01-core-principles.md#4-consistency-and-standards).

**Sources:** [Electron docs](https://www.electronjs.org/docs/latest/), [Tauri](https://tauri.app/), [Qt](https://doc.qt.io/).

### System integration

**Definition:** The integration surface distinguishing a desktop *citizen* from a ported webpage: file-type associations, recent-files (OS jump lists / Dock menus / GNOME recents), OS notifications, global shortcuts, rich clipboard, share targets, spotlight/search indexing, protocol handlers (`yourapp://`), startup-at-login options.

**Reasoning / Evidence:** Integration features are retention infrastructure: file association makes the app the *default verb* for its documents; recents shortcut the most common task (resume); protocol handlers connect web→app flows. Each is also an expectation — a document editor without file association feels broken.

**When to use:**
- Document apps: register types (with proper icons), populate OS recents, restore-last-session default.
- Communication/monitoring apps: OS notification center (with actions), tray/menu-bar agent mode, startup-at-login (opt-IN, never silent).
- Any app pairing with web: protocol handler for deep links (auth callbacks, "open in app").

**When NOT to use / exceptions:**
- Startup-at-login default-on without asking = resentment + startup-cleanup uninstalls; always opt-in with visible value.
- Global shortcuts: register few, show conflicts, make remappable (they collide across apps).
- Don't grab file associations already owned without asking (the classic PDF-reader war).

**Pros:** Workflow embedding = retention; OS surfaces (recents, search) become your free UI.
**Cons:** Per-platform APIs ×3; sandboxing limits (store-distributed builds lose some powers — document the delta).

**Implementation guidance:**
- Notifications: through OS APIs (respect Focus/DND), actionable where supported, never re-implemented as always-on-top windows ([04-interaction-design.md](04-interaction-design.md#notifications-and-interruption-design)).
- Clipboard: write rich + plaintext fallbacks; paste-special where formats matter.
- Uninstall respect: leave user documents; offer settings cleanup choice.

**Related:** [#windows-conventions](#windows-conventions), [#macos-conventions](#macos-conventions), [15-platform-mobile.md](15-platform-mobile.md#notifications-ux).

**Sources:** [Windows: App integration](https://learn.microsoft.com/en-us/windows/apps/develop/), [Apple: App integration technologies](https://developer.apple.com/documentation/).

### High-DPI and multi-monitor

**Definition:** Rendering correctly across display scale factors (100–300%, fractional like 125/150%) and monitor topologies (mixed-DPI setups, hotplug, different color profiles) — per-monitor DPI awareness as the standard.

**Reasoning / Evidence:** Mixed-DPI is now the norm (laptop 150% + external 100%); apps without per-monitor awareness render blurry (bitmap-stretched) or wrong-sized when dragged across displays. Fractional scaling breaks px-hinted assets and off-by-one borders.

**When to use:**
- Always: vector/DPI-independent assets (SVG, icon fonts, layered icon sizes), layout in DIPs/points not raw pixels, per-monitor-v2 awareness (Windows), test at 125/150/200%.
- Window drag across monitors: live rescale without layout explosion; remember per-monitor placement ([#window-management-and-state-restoration](#window-management-and-state-restoration)).

**When NOT to use / exceptions:** None — this is correctness; legacy-mode bitmap scaling is the fallback of shame.

**Pros:** Crisp everywhere; future display tolerance. **Cons:** Testing matrix; icon pipelines need multi-resolution discipline.

**Implementation guidance:**
- Icons: SVG or @1x/@1.5x/@2x/@3x sets; hairlines as 1 *physical* px where intended.
- Test: 100%+200% dual setup, drag every window across, check text sharpness, cursor sizes, and popup positioning at boundaries.

**Related:** [#window-management-and-state-restoration](#window-management-and-state-restoration), [03-visual-design.md](03-visual-design.md#iconography).

**Sources:** [Windows: DPI awareness](https://learn.microsoft.com/en-us/windows/win32/hidpi/high-dpi-desktop-application-development-on-windows), [Apple: Hi-DPI guidelines](https://developer.apple.com/design/human-interface-guidelines/images).

### Touch on hybrid devices

**Definition:** Accommodating touch/pen on convertible and touchscreen desktops (Surface-class, touch laptops): target sizes and gesture support layered onto pointer-first design without degrading the mouse experience.

**Reasoning / Evidence:** Windows convertibles are mainstream; users mix modalities mid-session (scroll by touch, precise-click by trackpad). Pointer-only assumptions (tiny targets, hover-gated actions) fail the touch moments ([04-interaction-design.md](04-interaction-design.md#hover-dependence)).

**When to use:**
- Baseline hybrid readiness: interactive targets ≥ ~32–40px effective (padding-extended), touch scrolling/momentum native, no hover-only functions, pinch-zoom where zooming exists.
- Pen (creation tools): pressure/tilt support, palm rejection, hover preview where the hardware offers it.

**When NOT to use / exceptions:**
- Don't inflate the whole UI to tablet scale on desktop — density stays; hit areas grow invisibly ([05-accessibility.md](05-accessibility.md#target-size)).
- Pointer-only pro tools (CAD fine-work) may reasonably limit touch to pan/zoom.

**Pros:** One app serves the device's whole posture range. **Cons:** More input QA; gesture/shortcut conflicts to arbitrate.

**Implementation guidance:**
- Detect input type per-interaction (pointer events), not per-device; optionally offer a "touch mode" spacing toggle (Office pattern).
- Test: run your critical flows entirely by touch on a real convertible.

**Related:** [15-platform-mobile.md](15-platform-mobile.md#touch-targets), [04-interaction-design.md](04-interaction-design.md#touch-and-pointer-gestures).

**Sources:** [Windows: Touch interactions](https://learn.microsoft.com/en-us/windows/apps/design/input/touch-interactions).

---

## Lifecycle

### Offline-first and autosave

**Definition:** Desktop-grade data safety: local-first storage with background sync (where cloud exists), continuous autosave, crash recovery restoring everything, and honest sync-state indication.

**Reasoning / Evidence:** Desktop users' baseline was set by native apps that never lost work offline; cloud-era desktop apps that die without connectivity betray the form factor. Crash recovery (restore all windows + unsaved content) is the difference between an incident and a catastrophe ([02-cognitive-foundations.md](02-cognitive-foundations.md#learned-helplessness-and-error-tolerant-design)).

**When to use:**
- Local-first: reads/writes hit local store; sync reconciles in background; full functionality offline for core tasks.
- Autosave continuous ([07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery)) + version history for documents; crash recovery journal.
- Sync status: quiet indicator (✓ synced / ↻ syncing / ⚠ conflict), never modal sync errors.

**When NOT to use / exceptions:**
- Genuinely connection-defined apps (trading terminals, comms) degrade to clearly-labeled read-only/reconnecting states rather than pretending.
- Conflict UX must exist before shipping sync: last-write-wins silently corrupts trust; offer both-versions resolution for real conflicts.

**Pros:** Resilience, speed (local reads), trust. **Cons:** Sync engineering is deep (conflicts, migrations, partial states) — budget accordingly.

**Implementation guidance:**
- Save-state visibility: document apps show "Edited/Saved" quietly by the title; never a Save button that lies (grayed while dirty).
- Unsaved changes handled explicitly: quit/close with dirty state prompts (or autosaves reliably) — never silent loss.
- Recovery drill: kill -9 during editing → relaunch must restore session + content; test it in CI if possible.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#autosave-and-draft-recovery), [15-platform-mobile.md](15-platform-mobile.md#offline-and-poor-connectivity).

**Sources:** [Local-first software (Ink & Switch)](https://www.inkandswitch.com/local-first/).

### Update UX

**Definition:** How the app ships new versions: silent background auto-update (browser model), prompted update (notify → user chooses when), or manual — plus restart choreography and change communication.

**Reasoning / Evidence:** Update prompts are recurring interruptions ([02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost)); the browser model (download silently, apply on natural restart) won because it converts updates from events into non-events. But silent updates that *change UX* without warning violate user trust — behavior changes need communication even when installation is silent.

**When to use:**
- Default: background download + apply-on-next-launch + subtle "restart to update" affordance (never forced mid-session restart).
- Enterprise contexts: admin-controllable channels/deferral (IT requirement).
- Change communication: in-app "what's new" *available but not modal* (badge on menu, first-run-after-update panel dismissible instantly); highlight breaking changes honestly.

**When NOT to use / exceptions:**
- Forced-restart countdowns over user work: never (the OS-update pattern users hate most).
- Nag-until-updated dialogs: only for security-critical/EOL versions, with a real "remind me later."
- Don't move users' cheese silently: major UI reorganizations deserve opt-in transitions or at minimum orientation aids ([02-cognitive-foundations.md](02-cognitive-foundations.md#mental-models)).

**Pros (silent+notify):** Security currency, no version-fragmentation, minimal interruption.
**Cons:** Session-state preservation across restart required ([#window-management-and-state-restoration](#window-management-and-state-restoration)); rollback path needed for bad releases.

**Implementation guidance:**
- Restart-to-update: one click, full session restoration, <10s total perceived cost.
- Downgrade/rollback: at minimum, settings-file forward-compatibility or migration warnings.
- Release notes: user-language benefits, not commit logs ([09-content-ux-writing.md](09-content-ux-writing.md)).

**Related:** [#offline-first-and-autosave](#offline-first-and-autosave), [02-cognitive-foundations.md](02-cognitive-foundations.md#interruption-cost).

**Sources:** [Chrome update model](https://www.chromium.org/developers/design-documents/software-updates-courgette/), [Sparkle (macOS)](https://sparkle-project.org/).

### First-run experience

**Definition:** The install→first-value path on desktop: minimal-permission startup, immediate usable state (sample content/empty-state guidance over tour modals), progressive feature introduction, and import/migration from predecessor tools.

**Reasoning / Evidence:** Desktop installs signal intent — users arrive with a task; blocking it with account walls, permission barrages, and tour carousels burns the strongest moment ([Peak-End](01-core-principles.md#peak-end-rule): beginnings anchor too). Migration/import (from the competitor they're leaving) is the highest-leverage first-run feature for switchers.

**When to use:**
- Get to the working surface in seconds: skip-able sign-in where possible, sample project or guided empty state ([04-interaction-design.md](04-interaction-design.md#empty-states)), first-task tooltip *at the moment of relevance* not before.
- Import: detect predecessor data (bookmarks, settings, projects) and offer one-click migration.
- Permission asks (notifications, startup-at-login): deferred until the feature that needs them is invoked ([13-platform-web.md](13-platform-web.md#notification-permission-ux) — same doctrine).

**When NOT to use / exceptions:**
- Genuinely account-required products (sync-defined) keep the wall but make it fast (SSO/passkeys — [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)).
- Complex pro tools may offer an optional interactive tutorial project — optional being the operative word.

**Pros:** Activation rate; word-of-mouth first impressions.
**Cons:** Sample content/import engineering; the org's every team wants a first-run slide (defend the path).

**Implementation guidance:**
- Measure time-to-first-value (install→first meaningful action) as a core metric; target minutes, cut everything that doesn't serve it.
- Checklist pattern for multi-step setup (connect account, import, invite) — visible, optional, dismissible ([01-core-principles.md](01-core-principles.md#zeigarnik-effect)).

**Related:** [04-interaction-design.md](04-interaction-design.md#empty-states), [15-platform-mobile.md](15-platform-mobile.md#onboarding), [09-content-ux-writing.md](09-content-ux-writing.md#onboarding-copy-and-tooltips).

**Sources:** [NN/g: Onboarding tutorials](https://www.nngroup.com/articles/onboarding-tutorials/).

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Desktop mindset | Optimize the 100th session; restore everything | — |
| Windows | Fluent 2, Ctrl keys, snap-friendly, jump lists | Alt+F4; F2 rename |
| macOS | Full menu bar mandatory; Cmd keys; sheets | Cmd+, settings; Cmd+Q vs Cmd+W |
| Linux | GNOME HIG + freedesktop standards; Flatpak | XDG dirs |
| Cross-platform | Localize chrome/shortcuts/dialogs per OS | See mapping table |
| Shortcuts | Standards sacred; menus teach; remappable (pro) | Never repurpose Ctrl/Cmd+S/Z/C/V |
| Menus | Complete command inventory + shortcuts shown | ≤9 items/group; … for dialogs |
| Context menus | Accelerator, never sole path | 3–12 items; Shift+F10 opens |
| Command palette | Ctrl/Cmd+K fuzzy everything | Recents on empty query |
| Windows/state | Restore size/position/session; never off-screen | Validate against topology |
| Modality | Least modality that works; sheets over app-modal | Esc cancels, Enter defaults |
| Density | Master-detail, grids, panels; reset-layout escape | Rows 24–40px; density toggle |
| Drag & drop | OS files in/out; spring-load; modifier badges | Spring 500–800ms |
| Framework choice | Behavior must be native even if chrome isn't | Cold start ≤2–3s |
| Integration | Associations, recents, notifications, protocols | Startup-at-login opt-in only |
| Hi-DPI | Per-monitor aware, vector assets | Test 100/125/150/200% |
| Hybrid touch | Bigger hit areas, no hover-gating | Effective targets ≥32–40px |
| Offline/autosave | Local-first, crash recovery, quiet sync state | kill -9 drill must pass |
| Updates | Silent download, user-timed restart, real notes | Restart cost <10s |
| First run | Task in seconds; import rivals; defer asks | Time-to-first-value metric |
