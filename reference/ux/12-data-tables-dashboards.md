# Data Tables, Grids & Dashboards

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers presenting and manipulating structured data at scale: choosing between tables, cards, grids, and trees; static-table semantics and accessibility; sorting, filtering, and column management in large tables; interactive (editable) data grids; bulk actions; density decisions; and dashboards with their metric cards, chart selection, and visualization accessibility. Faceted filtering, search, and pagination mechanics live in [08-navigation-ia.md](08-navigation-ia.md); grid/table/treegrid keyboard patterns in [06-aria-widget-reference.md](06-aria-widget-reference.md); progress, loading, and empty states in [11-components-and-overlays.md](11-components-and-overlays.md); density tokens in [03-visual-design.md](03-visual-design.md#density-modes); desktop master-detail and data-rich layouts in [14-platform-desktop.md](14-platform-desktop.md).

## Contents

- [Part 1 — Tables and data grids](#part-1--tables-and-data-grids)
  - [Pattern selection](#pattern-selection)
  - [Static table](#static-table)
  - [Sorting, filtering, and column management](#sorting-filtering-and-column-management)
  - [Interactive data grid](#interactive-data-grid)
  - [Bulk actions](#bulk-actions)
  - [Data density decision rules](#data-density-decision-rules)
- [Part 2 — Dashboards and data visualization](#part-2--dashboards-and-data-visualization)
  - [Dashboard fit check](#dashboard-fit-check)
  - [Metric card](#metric-card)
  - [Chart selection matrix](#chart-selection-matrix)
  - [Visualization accessibility](#visualization-accessibility)
  - [Dashboard layout rules](#dashboard-layout-rules)
- [Part 3 — Acceptance criteria](#part-3--acceptance-criteria)
- [Quick reference](#quick-reference)

---

## Part 1 — Tables and data grids

### Pattern selection

Choose the presentation pattern from the data need, not from what the component library makes easy:

| Data need | Use | Avoid | Why |
| --- | --- | --- | --- |
| Compare many records by the same attributes | Table | Cards | Tables align attributes in columns for vertical scanning |
| Edit many cells rapidly | Data grid | Form per row | Grid supports expert efficiency (spreadsheet mental model) |
| Browse rich visual items | Cards/grid | Table | Visual content needs flexible layout; imagery is the decision input |
| Show hierarchy | Tree / tree table | Flat table | Hierarchy must be navigable (expand/collapse, level context) |
| Monitor metrics over time | Dashboard / table hybrid | Static table only | Monitoring needs status, trends, and drilldown |
| Mobile summary of records | Responsive cards/list | Squashed table | Narrow screens need attribute prioritization, not compression |

Decision questions (answer before building):
- Do users compare records attribute-by-attribute? → table.
- Do users *edit* in place, at volume? → data grid (a much larger commitment — see [Interactive data grid](#interactive-data-grid)).
- Is imagery the primary decision input with few attributes? → cards beat tables ([NN/g data-table research](https://www.nngroup.com/articles/data-tables/) distinguishes lookup/compare tasks, which favor tables, from browse tasks, which don't).
- Will this render on touch/narrow viewports? → design the mobile list variant explicitly; never rely on horizontal squashing.

---

### Static table

ID: DATA-TABLE

**Use when:** noninteractive (or lightly interactive: sort, link-per-row) tabular data — records sharing the same attributes, read and compared more than edited.

**Avoid when:** using a table for page layout; complex interactive editing (use a data grid); imagery-driven browsing (use cards).

**Behavior contract:**
- Semantic table element (or platform-native table API) — never divs styled to look tabular.
- Every data cell programmatically associated with its header(s): `<th scope="col">` / `<th scope="row">`; `headers`/`id` for complex multi-level headers. This is what lets screen readers announce "Status: Active, row 12" instead of a context-free word ([W3C APG table pattern](https://www.w3.org/WAI/ARIA/apg/patterns/table/)).
- A `<caption>` or programmatically associated heading identifies what the table contains.
- First column carries the identifying information (name/ID) — it's where scanning starts ([02-cognitive-foundations.md](02-cognitive-foundations.md#scanning-patterns)).
- If sortable: sort state exposed both visually (arrow on the sorted column) and programmatically (`aria-sort`).
- Numeric columns right-aligned with tabular numerals; text columns left-aligned ([03-visual-design.md](03-visual-design.md)).

**Failure modes:**
- Table used for layout (screen readers announce meaningless row/column positions).
- Header row visually styled as a header but marked up as `<td>` — no cell association exists.
- Responsive collapse (stacking cells vertically on mobile) breaks the header↔cell relationship: a bare "42" with no visible or programmatic label. If you linearize a table on small screens, repeat the label per value (visually or via generated content with an accessible equivalent), or switch to a designed list/card variant.
- Truncated cell values with no way to see the full value (no tooltip, no expand, no detail view).

**Verification checklist:**
- [ ] Screen reader announces column and row context when navigating cells.
- [ ] Caption/heading states table purpose and scope.
- [ ] Sorted column exposes `aria-sort`; sort control is a real button with an accessible name.
- [ ] Mobile/narrow rendering keeps every value labeled.
- [ ] Every truncated value has a full-value path (tooltip on hover *and* focus, expandable row, or detail view).

---

### Sorting, filtering, and column management

**Definition:** Navigation and orientation mechanics inside tables of hundreds to millions of rows: column sorting, per-column and global filtering, column show/hide/reorder/resize/pin, density control, and row-level navigation to detail views.

**Reasoning / Evidence:** In data-heavy enterprise UIs the table *is* the navigation surface: users locate records by composing sort + filter, not by menus. Research-backed patterns from enterprise design systems (Carbon, Fluent, Ant Design) and NN/g's complex-application work: visible sort indicators and persistent filter chips serve as you-are-here signals for the data space; losing table state on navigation (drill into a row, come back, state reset) is one of the most-reported enterprise UX failures because it destroys the user's constructed context.

**Use when:**
- Any list >~25 rows used for lookup, triage, or comparison; anything called "admin," "console," or "ops."

**Avoid when:**
- Card grids beat tables when imagery is the primary decision input and attributes are few.
- Don't add column management to simple 4–5-column tables that never scroll horizontally — configuration UI without need is noise.

**Behavior contract:**
- **Sorting:** click header cycles asc → desc (a third click may clear); show arrow on the sorted column and gray sort affordances on hover of others; support secondary sort (Shift+click) for power users; default sort must be meaningful (recency for activity, name for directories) and stated. Sort labels should be unambiguous: "Sort by Due date, ascending." Preserve sort across navigation.
- **Filtering:** global quick-filter box above the table (<300ms response) plus structured per-column filters; render active filters as removable chips with "Clear all"; show "143 of 5,204 rows" so users know they're viewing a subset — the single most missed status signal. The no-results state must distinguish *no data exists* from *over-filtered* ([08-navigation-ia.md](08-navigation-ia.md#no-results-recovery)).
- **Column management:** a "Columns" control listing all columns with checkboxes and drag-reorder; pin/freeze the identity column(s) left; persist widths, order, visibility, sort, and filters per user per table; provide "Reset to default."
- **Sticky headers:** keep the header row visible during vertical scroll — column meaning must never scroll away from the data it labels.
- **Density:** comfortable ~52px rows default, compact ~32px option for experts; persist choice ([03-visual-design.md](03-visual-design.md#density-modes)).
- **Row navigation:** entire row clickable to detail (with visible hover state), or explicit link on identity column; preserve table state on back-navigation (URL-serialized or session-restored); support open-in-new-tab (real links, not JS handlers).
- **Saved views:** named sort+filter+column bundles shareable by URL — the highest-leverage feature for teams above ~10⁴ rows.
- **Pagination integration:** virtualize rendering beyond ~200 visible rows; pair with pagination or load-more semantics per [08-navigation-ia.md](08-navigation-ia.md#pagination-vs-infinite-scroll-vs-load-more). Use pagination when users need a stable location, comparison across pages, or footer access; infinite scroll only for exploratory feeds with robust state restoration.

**Failure modes:**
- Feature-dense headers become cluttered; column management, saved views, and virtualization are significant engineering — budget for them or cut scope.
- Horizontal scroll hides identity columns unless pinned.
- Table state (sort/filter/scroll) reset on drill-in-and-back.
- Filter applied silently with no row-count feedback — users conclude data is missing.

**Verification checklist:**
- [ ] Current sort column and direction visible and exposed via `aria-sort`.
- [ ] Active filters shown as chips with per-chip removal and "Clear all."
- [ ] Row count shows filtered subset vs total ("143 of 5,204").
- [ ] Header stays visible while scrolling; identity column pinned when the table scrolls horizontally.
- [ ] Drill into a row and navigate back: sort, filters, and scroll position survive.

**Related:** [08-navigation-ia.md](08-navigation-ia.md#faceted-navigation-and-filtering) (faceted filtering), [08-navigation-ia.md](08-navigation-ia.md#pagination-vs-infinite-scroll-vs-load-more), [14-platform-desktop.md](14-platform-desktop.md#density-and-information-rich-layouts), [10-design-systems.md](10-design-systems.md) (component APIs).

**Sources:** [Carbon Design System: Data table](https://carbondesignsystem.com/components/data-table/usage/), [NN/g: Data Tables — Four Major User Tasks](https://www.nngroup.com/articles/data-tables/), [Pencil & Paper: Data table UX patterns](https://www.pencilandpaper.io/articles/ux-pattern-analysis-enterprise-data-tables).

---

### Interactive data grid

ID: DATA-GRID

**Use when:** spreadsheet-like navigation/editing, virtualized data, selectable rows/cells — users manipulate many values quickly and expect Excel-adjacent keyboard behavior.

**Avoid when:** data is static or lightly interactive (a plain table is cheaper and more accessible); each record needs rich per-field validation and guidance (a form or side panel per record is safer).

**Required decisions** (make each explicitly; defaults chosen by accident are the main source of grid bugs):
- **Row vs cell focus model.** Row focus suits triage/selection workflows; cell focus suits editing. In either case the grid is a *single tab stop*: Tab enters/leaves the grid, arrow keys move within it ([W3C APG grid pattern](https://www.w3.org/WAI/ARIA/apg/patterns/grid/)).
- **Selection model:** row, cell, range, multi-select — and how selection is shown and announced.
- **Edit model:** inline (fast, low-context — best for single-field corrections), modal (full validation, loses table context), side panel (edit while seeing the row in context — often the best middle ground; see [14-platform-desktop.md](14-platform-desktop.md) master-detail), or batch (edit many, commit once).
- **Keyboard model:** arrow keys move cell focus; Tab behavior predictable (leave grid, or move cell-to-cell *only* inside an active edit row); Home/End/PageUp/PageDown where users expect them; **Enter starts or commits an edit; Escape cancels the edit and restores the prior value** — this pair is the non-negotiable core of edit semantics.
- **Virtualization/accessibility model:** virtualized rows must still expose accurate row/column counts and positions (`aria-rowcount`, `aria-rowindex`) so assistive tech isn't told the table is 40 rows when it's 40,000.
- **Undo and validation model:** inline edits need cell-level validation shown at the cell (not a distant banner), a clear saving/failed state per cell or row, and undo for committed edits — inline editing without undo converts slips into data corruption.

**Required states** — design all of them, not just the happy path: loading, empty, no results, sorted, filtered, selected, focused, editing, invalid, saving, failed, **stale data** (the row changed on the server since it was loaded — show it, and never silently overwrite either version). Empty/loading/progress patterns: [11-components-and-overlays.md](11-components-and-overlays.md).

**Failure modes:**
- Every cell is a tab stop — keyboard users need hundreds of Tab presses to cross the grid.
- Edit committed on blur with no feedback, or lost on scroll (virtualization discards uncommitted state).
- Validation error shown only after a bulk save, with no pointer to the offending cell.
- Concurrent edit conflict resolved by last-write-wins with no warning.

**Verification checklist:**
- [ ] Tab enters the grid once; arrows navigate; Tab exits.
- [ ] Enter opens/commits an edit; Escape cancels and restores the previous value.
- [ ] Every state in the required-states list has a designed rendering.
- [ ] Committed edits are undoable; failed saves are visibly attributed to the row/cell.
- [ ] Screen reader announces row/column position and edit state (see [06-aria-widget-reference.md](06-aria-widget-reference.md) grid/treegrid patterns).

---

### Bulk actions

**Use when:** users repeat the same operation across many records (assign, tag, archive, delete, export).

**Avoid when:** the action is rarely applied to more than one record, or per-record confirmation is genuinely required (e.g., irreversible legal actions) — then bulk selection creates risk without saving time.

**Behavior contract:**
- **Selection count always visible** while anything is selected ("12 selected"), typically in a toolbar that appears with selection — the count is the user's only guard against acting on the wrong scope.
- **"Select all visible" vs "select all matching filter" must be explicit.** Checking the header checkbox selects the visible page; offer a distinct, labeled step to extend to all N matching rows ("Select all 5,204 matching"). Silently expanding scope is how mass-deletion incidents happen.
- Bulk action scope stated in the action itself or its confirmation ("Archive 12 items?").
- **Destructive bulk actions require review/confirmation or undo** — prefer undo with a persistent window for reversible operations, confirmation with an explicit count for irreversible ones ([04-interaction-design.md](04-interaction-design.md)).
- Partial failure reported per item ("9 archived, 3 failed: …"), not as a single opaque error.

**Failure modes:**
- Selection cleared by sorting, filtering, or paging without warning.
- Header checkbox ambiguity: user believes "all" means all matching rows, system means visible page (or vice versa).
- Bulk destructive action fires with no count in the confirmation.

**Verification checklist:**
- [ ] Selection count visible at all times while selection is active.
- [ ] All-visible vs all-matching are separate, labeled steps.
- [ ] Destructive bulk action has undo or an explicit count-bearing confirmation.
- [ ] Selection behavior across sort/filter/page changes is defined and, when destructive, conservative.

---

### Data density decision rules

Use **high density** (compact rows, more columns, hover-revealed secondary actions) when:
- Users are experts with long sessions.
- Comparison and monitoring are the core tasks.
- Pointer/keyboard precision is expected (desktop).

Use **low density** (taller rows, fewer attributes, visible actions) when:
- Users are novices or infrequent.
- Touch input dominates — keep targets ≥44–48px regardless of visual density ([15-platform-mobile.md](15-platform-mobile.md)).
- Reading comprehension matters more than scanning throughput.
- Decisions are high-risk.

Agent rule: for complex enterprise/productivity tools, offer a user-switchable density control (comfortable/compact) and persist it per user per view. Token-level implementation and row-height numbers (32/40/48–52px): [03-visual-design.md](03-visual-design.md#density-modes); desktop-specific guidance: [14-platform-desktop.md](14-platform-desktop.md#density-and-information-rich-layouts). NN/g's complex-application research finds expert users prefer and perform better with higher density; the sin is density *without hierarchy*, not density itself.

---

## Part 2 — Dashboards and data visualization

### Dashboard fit check

Use a dashboard when users need to **monitor, compare, diagnose, or act on changing information**. Do not use a dashboard as a dumping ground for every available metric — Stephen Few's dashboard-design work defines a dashboard as a single-screen display of the information needed to achieve one or more objectives, monitorable at a glance; anything that doesn't serve a decision dilutes what does.

Decision questions:
- What question does each widget answer, and what would the user *do* differently based on the answer? No action path → cut it (vanity metric).
- Is the data changing often enough to monitor? If not, a report beats a dashboard.

### Metric card

ID: DATA-VIZ-METRIC-CARD

**Use when:** a single decision-relevant KPI needs at-a-glance visibility with context.

**Avoid when:** the metric has no decision relevance, no owner, or no meaningful baseline.

**Behavior contract — required content, all of it:**
- Metric name (in the user's vocabulary, not the data warehouse's).
- Value, with unit.
- Time range the value covers ("last 7 days") — a number without a period is unfalsifiable.
- Comparison/baseline or trend where relevant (vs previous period, vs target), with direction *and* whether that direction is good (color + arrow + text, not color alone).
- Freshness timestamp ("as of 09:42") and update cadence.
- Status/threshold meaning if the card carries a status (what makes it red?).

**Failure modes:**
- Metric lacks denominator or time period ("Errors: 34" — of how many? since when?).
- Red/green-only status excludes colorblind users (WCAG 1.4.1 Use of Color — [05-accessibility.md](05-accessibility.md#use-of-color)).
- Delta shown without baseline ("+12%" against what?).

**Verification checklist:**
- [ ] Card answers: what is it, how much, over what period, compared to what, how fresh.
- [ ] Trend/status readable without color (icon/text redundancy).
- [ ] Tapping/clicking the card drills into detail or the card explicitly doesn't drill.

### Chart selection matrix

Choose the chart from the analytical question:

| Need | Use | Avoid |
| --- | --- | --- |
| Trend over time | Line chart | Pie chart |
| Compare categories | Bar chart | 3D chart (depth distorts value perception) |
| Part-to-whole, few categories | Stacked bar; donut cautiously | Many-slice pie |
| Distribution | Histogram / box plot | Average-only metric (hides shape) |
| Relationship between variables | Scatter plot | Dual-axis chart unless the audience is expert (dual axes invite spurious correlation reading) |
| Geographic pattern | Map | Map for non-geographic data |
| Exact values needed | Table | Chart-only |
| Part-to-whole across a hierarchy | Treemap | Nested pies |

General principles (Tufte): maximize data-ink, remove chartjunk — decoration that obscures comparison is a failure mode, not a style choice. Start quantitative axes for bar charts at zero; label axes, units, and series directly (direct labeling beats a distant legend when series count is small).

### Visualization accessibility

**Behavior contract:**
- **Never rely on color alone** to distinguish series or encode status — add pattern, shape, direct labels, or position (WCAG 1.4.1; [05-accessibility.md](05-accessibility.md#use-of-color)).
- Provide a **text summary of the chart's insight** ("Signups rose 18% after the March release") — the takeaway, not a pixel description.
- Provide a **data table or accessible data alternative** for every important chart (toggle, link, or adjacent table); charts rendered as canvas/SVG without one are invisible to screen readers.
- Label axes, units, and series; ensure contrast for lines, bars, and text against the chart background (3:1 minimum for graphical objects — [05-accessibility.md](05-accessibility.md)).
- Don't require animation to understand a change; respect reduced-motion preferences ([04-interaction-design.md](04-interaction-design.md)).
- Colorblind-safe palettes for categorical series ([03-visual-design.md](03-visual-design.md#color-blindness-and-color-independence)).

**Verification checklist:**
- [ ] Chart understandable in grayscale.
- [ ] Every important chart has a text insight + data-table alternative.
- [ ] Axes, units, and series labeled; no unlabeled dual axis.

### Dashboard layout rules

**Behavior contract:**
- **Group widgets by user question, not by data source** — "Is checkout healthy?" collects metrics from three systems; users don't care which.
- **Most decision-critical content top-left** (scanning starts there in LTR locales — [02-cognitive-foundations.md](02-cognitive-foundations.md#scanning-patterns)); supporting detail below and right.
- Filters and scope (date range, environment, team) global, visible, and persistent — a dashboard whose scope is ambiguous produces confidently wrong readings.
- **Data freshness visible** per dashboard and, where cadences differ, per widget; distinguish real-time, delayed, sampled, estimated, and forecasted data explicitly.
- **Every widget needs its own empty, loading, and error states** ([11-components-and-overlays.md](11-components-and-overlays.md)) — one failed query must not render as a zero.
- Don't auto-refresh content out from under an interacting user; refresh on an announced cadence with a pause, or update passively with a "new data" affordance.

**Failure modes:**
- Dashboard organized by database schema.
- A widget silently showing 0 because its query failed.
- Mixed freshness (one card real-time, its neighbor daily-batch) with no indication.

---

## Part 3 — Acceptance criteria

A tables-and-dashboards surface ships when all of the following hold:

**Comprehension**
- [ ] User can identify data scope, active sort, active filters, selection, and freshness at a glance.
- [ ] Every chart has title, scope, time range, unit, and an accessible alternative.
- [ ] Dashboard user can answer: what changed, why it matters, what to do next.
- [ ] Empty, no-results, loading, and error states are distinct and designed ([11-components-and-overlays.md](11-components-and-overlays.md)).

**Accessibility**
- [ ] Screen reader announces row and column context for any cell; header association is programmatic, not visual (verify against [06-aria-widget-reference.md](06-aria-widget-reference.md) table/grid/treegrid patterns).
- [ ] No value is truncated without an accessible path to the full value.
- [ ] No status or series distinction depends on color alone (WCAG 1.4.1).

**Keyboard**
- [ ] Every mouse interaction — sort, filter, select, edit, bulk action, drilldown, chart data access — has a keyboard path ([06-aria-widget-reference.md](06-aria-widget-reference.md)).
- [ ] Grids are a single tab stop with arrow-key navigation; Enter/Escape edit semantics work as specified.

**State integrity**
- [ ] Table state (sort/filter/columns/scroll) survives drill-in and back-navigation.
- [ ] Bulk selection scope is explicit; destructive bulk actions are undoable or count-confirmed.
- [ ] Stale/conflicting data is surfaced, never silently overwritten.

---

## Quick reference

| Topic | One-line guidance | Key numbers |
|---|---|---|
| Pattern choice | Table to compare, grid to edit, cards to browse, tree for hierarchy, list on mobile | — |
| Static table semantics | Real `<th scope>` + caption; responsive collapse must keep labels with values | — |
| Sorting | Visible arrow + `aria-sort`; meaningful stated default; Shift+click secondary sort | Click cycles asc → desc (→ clear) |
| Filtering | Chips + "Clear all" + filtered-count status | "143 of 5,204"; quick filter <300ms |
| Column management | Show/hide/reorder/pin; persist per user; reset to default | Pin identity column left |
| Table state | Preserve sort/filter/scroll across drill-in and back; URL-serialized | Saved views pay off above ~10⁴ rows |
| Density | Expert+pointer → compact; novice/touch/high-risk → comfortable; user toggle, persisted | Rows ~52px comfortable / ~32px compact; virtualize >~200 rows |
| Grid keyboard | One tab stop; arrows move; Enter edits/commits; Escape cancels | `aria-rowcount`/`aria-rowindex` when virtualized |
| Grid editing | Choose inline/modal/side-panel/batch deliberately; cell-level validation; undo for commits | Design all states incl. stale/failed |
| Bulk actions | Count always visible; all-visible ≠ all-matching; undo or count-bearing confirm for destructive | — |
| Metric card | Name, value, unit, period, comparison, freshness, threshold meaning | — |
| Chart choice | Trend→line, compare→bar, part-of-whole→stacked bar/treemap; no many-slice pie, 3D, casual dual-axis | Bar axes start at zero |
| Viz accessibility | Never color alone; text insight + data table per important chart | 3:1 graphical contrast; WCAG 1.4.1 |
| Dashboard layout | Group by user question; critical top-left; visible scope + freshness; per-widget error states | — |
