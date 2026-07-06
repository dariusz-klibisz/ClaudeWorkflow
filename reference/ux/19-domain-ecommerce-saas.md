# Domain Patterns: E-commerce, SaaS & Accounts

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file covers four domains where generic UX rules get their highest-stakes application: e-commerce and checkout (listing → PDP → cart → checkout, plus the conversion-ethics line), SaaS/enterprise/admin products (scope visibility, permission-aware UI, bulk operations, expert accelerators, destructive admin actions), settings and preferences (IA, save models, honest defaults), and accounts/permissions/authentication (login patterns, signup, permission prompts, roles, deletion/export). Each entry keeps the decision questions and failure modes that make it usable for review. The deep authentication doctrine — NIST SP 800-63B password rules, passkeys, magic links, OTP mechanics — lives in [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux) and is cross-linked, not duplicated. Dark-pattern definitions and legal status live in [18-patterns-antipatterns.md](18-patterns-antipatterns.md).

## Contents

- [E-commerce & checkout](#e-commerce--checkout)
  - [Product listing and category pages](#product-listing-and-category-pages)
  - [Product detail page (PDP)](#product-detail-page-pdp)
  - [Cart](#cart)
  - [Checkout flow](#checkout-flow)
  - [Conversion ethics](#conversion-ethics)
- [SaaS, enterprise & admin](#saas-enterprise--admin)
  - [Enterprise UX priorities](#enterprise-ux-priorities)
  - [Scope visibility](#scope-visibility)
  - [Permission-aware UI](#permission-aware-ui)
  - [Bulk operations](#bulk-operations)
  - [Expert accelerators](#expert-accelerators)
  - [Admin destructive actions](#admin-destructive-actions)
- [Settings & preferences](#settings--preferences)
  - [Settings information architecture](#settings-information-architecture)
  - [Setting types and save behavior](#setting-types-and-save-behavior)
  - [Honest defaults](#honest-defaults)
  - [Reset, undo, and settings search](#reset-undo-and-settings-search)
- [Accounts, permissions & authentication UX](#accounts-permissions--authentication-ux)
  - [Login pattern matrix](#login-pattern-matrix)
  - [Signup minimization](#signup-minimization)
  - [Permission-prompt timing](#permission-prompt-timing)
  - [Roles and permission models](#roles-and-permission-models)
  - [Account deletion and data export](#account-deletion-and-data-export)
  - [Authentication failure modes](#authentication-failure-modes)
- [Quick reference](#quick-reference)

---

## E-commerce & checkout

The governing principle: **reduce uncertainty before purchase**. Total cost, delivery, returns, and availability must be clear as early as possible; cart and progress must be preserved; conversion must never be optimized through deception. The Baymard Institute's long-running checkout research consistently finds that **unexpected extra costs (shipping, taxes, fees) are the most common reason shoppers abandon checkout** — which makes early total-cost disclosure a conversion tactic, not just an ethical one.

### Product listing and category pages

**Use when:** Users browse and compare product sets — category pages, search results, collection pages.

**Avoid when:** A single-product or single-plan business; don't force a listing layer where there's nothing to compare.

**Required content / decisions:**
- Product card: image, name, price, availability, rating, and key variants — the differentiators a user needs to decide *without* opening every PDP.
- Filters: visible sidebar vs drawer — decide by screen size and filter count.
- Sort: relevance, price, rating, newest, availability — offer what the catalog makes meaningful.
- Pagination vs infinite scroll vs load-more: decide by return-to-position needs and footer reachability (decision table in [08-navigation-ia.md](08-navigation-ia.md)).
- Sale pricing: original price, discounted price, and the *total* price story must be unambiguous.

**Failure modes:**
- Filters produce zero results with no recovery path (no "remove filter" chips, no relaxation suggestions).
- Product cards omit the key differentiator, forcing pogo-sticking into every PDP.
- Sale price vs unit price vs total price unclear — a comparison-prevention smell ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#comparison-prevention-hard-to-compare-pricing)).

**Verification checklist:**
- [ ] Cards show enough to compare (price, availability, key variant info) without opening the PDP.
- [ ] Zero-result filter states offer one-click recovery.
- [ ] Applied filters are visible and individually removable.
- [ ] Returning from a PDP restores scroll/list position.

### Product detail page (PDP)

**Use when:** The decision point for a single product — everything needed to commit belongs here.

**Required content:**
- Product name; images/media with accessible alternatives.
- Price **and fees** — no drip pricing deferred to checkout.
- Variants/options with clear selection state; invalid combinations disabled or explained *before* add-to-cart, never as a post-click error.
- Availability (in stock, backorder, delivery estimate).
- Delivery and returns information — accessible from the PDP, not buried in footer legal.
- One primary CTA (Add to cart / Buy now).
- Reviews or trust signals where genuinely available; specifications for comparison-heavy products.

**Decision question:** If users must choose variants, can they tell what's selected, what's unavailable, and why — before pressing the CTA?

**Failure modes:**
- Variant selection resets images/price silently, or add-to-cart fails because a variant wasn't chosen and nothing indicates which.
- Fees revealed only at checkout (see [Conversion ethics](#conversion-ethics)).
- Fabricated review counts or urgency banners ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#fake-scarcity-fake-urgency-and-fake-social-proof)).

**Verification checklist:**
- [ ] Price shown includes or clearly states all mandatory fees.
- [ ] Variant selection state is visible; invalid combos are explained before the CTA.
- [ ] Delivery cost/date estimate and return policy reachable from the PDP.
- [ ] Media has text alternatives ([05-accessibility.md](05-accessibility.md)).

### Cart

**Use when:** Holding state between browsing and checkout — the cart is a *workspace*, not a toll booth.

**Behavior contract:**
- Shows items, quantities, variants, and a full price breakdown.
- Quantity editing, removal, and save-for-later work in place.
- Delivery/tax estimate shown, or an explicit "calculated at checkout" statement — never silence.
- Items that became unavailable or changed price are flagged with recovery options, not silently dropped or silently repriced.
- Cart persists across sessions where the platform permits.

**Failure modes:**
- Silent quantity caps or removals discovered only at payment.
- "Sneak into basket" additions — items or services the user never chose ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#sneaking-sneak-into-basket)).
- Cart lost on login, tab close, or session expiry.

**Verification checklist:**
- [ ] Every line item is editable/removable in the cart.
- [ ] Price breakdown distinguishes items, delivery, taxes, discounts.
- [ ] Price/availability changes since adding are surfaced with an explanation.
- [ ] Nothing appears in the cart the user didn't explicitly add.

### Checkout flow

**Use when:** Purchase completion — the highest-stakes form sequence in commercial UX. Form mechanics (layout, validation, autofill, payment fields) are covered in [07-forms-and-input.md](07-forms-and-input.md#high-risk-and-regulated-forms); this entry covers the flow contract.

**Required steps:** Contact/shipping → delivery method → payment → review → confirmation/receipt. Collapse steps where data is minimal, but never remove the review before commit.

**Behavior contract:**
- **Guest checkout first.** Do not require account creation before purchase unless legally or operationally essential; offer optional account creation *after* the order (one extra field: password — or none, with a passkey/magic-link invite).
- **Preserve data across steps** — back navigation, validation errors, and payment failures must never wipe entered data (payment errors preserve all non-sensitive input and explain the next step).
- **No redundant asking:** billing address defaults to **same as shipping** (pre-checked, single checkbox to change).
- **Total cost as early as possible:** shipping, taxes, and fees appear no later than the step where they're determinable — and always before the final commit (Baymard: surprise costs at the end are the top abandonment trigger).
- **Final button names the outcome:** `Pay $42.00`, `Place order`, `Confirm booking` — never `Continue` or `Submit` ([09-content-ux-writing.md](09-content-ux-writing.md#button-labels)).
- Delivery date/cost, and access to the return/cancellation policy, visible before final commit; secure-payment context communicated honestly (real processor identity, not decorative trust badges).

**Accessibility requirements:** autofill-compatible fields (`autocomplete` attributes), paste allowed in card/password/OTP fields ([05-accessibility.md](05-accessibility.md#accessible-authentication), WCAG 3.3.8), inline errors plus an error summary, and mobile-adequate target sizes.

**Failure modes:**
- Forced registration wall before checkout (a classic abandonment driver).
- Costs revealed only on the review step ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#hidden-costs-and-drip-pricing)).
- Payment decline dumps the user back to an empty form.
- Address validation rejects real addresses with no override.

**Verification checklist:**
- [ ] Guest checkout available; account creation offered after purchase, not before.
- [ ] Billing = shipping is the default with a one-click change.
- [ ] Full total (items + delivery + taxes + fees) visible before the final button.
- [ ] Final button states the consequence and amount.
- [ ] Errors at any step preserve all entered non-sensitive data.
- [ ] Confirmation includes receipt, order reference, and what happens next.

**Related:** [07-forms-and-input.md](07-forms-and-input.md#high-risk-and-regulated-forms), [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs).

**Sources:** [Baymard Institute checkout research](https://baymard.com/lists/cart-abandonment-rate).

### Conversion ethics

Conversion optimization has a bright line: **help users decide; never decide for them or deceive them.** Full pattern catalog with legal status in [18-patterns-antipatterns.md](18-patterns-antipatterns.md#dark-patterns-catalog).

| Allowed | Not allowed |
|---|---|
| Clarify value (better copy, comparisons, honest social proof) | Hidden fees / drip pricing ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#hidden-costs-and-drip-pricing)) |
| Reduce unnecessary friction (fewer fields, autofill, guest checkout) | Preselected add-ons, insurance, donations ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#preselection)) |
| Improve error recovery (preserved data, clear next steps) | Obstructed cancellation ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#hard-to-cancel-roach-motel)) |
| Make choices understandable (plain pricing, honest defaults) | Fake scarcity/urgency/social proof ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#fake-scarcity-fake-urgency-and-fake-social-proof)) |
| Genuine, verifiable scarcity ("2 left" that is true) | Confusing opt-outs, confirmshaming ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#confirmshaming)) |
| Remembering user context to shorten repeat purchase | Visual manipulation that hides cost or consequence ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#visual-interference-misdirection)) |

**Decision question:** Would this tactic still work if the user fully understood it? If yes, it's persuasion; if no, it's deception — and increasingly, legal liability.

---

## SaaS, enterprise & admin

### Enterprise UX priorities

Enterprise, B2B, internal-tool, and admin-console UX reorders the usual priorities:

- **Efficiency** for frequent expert workflows (users live in the tool 8 hours a day; every extra click compounds).
- **Learnability** for occasional users (the quarterly-report user must not be locked out by expert-only design).
- **Error prevention at scale** — one admin mistake can affect thousands of users or records; safeguards scale with blast radius.
- **Clear scope** across tenants, workspaces, environments, roles, and permissions.
- **Auditability and recovery** — who did what, when, and how to undo it.
- **Configurability without chaos** — flexibility bounded by sane defaults.

**Decision question for any admin feature:** what is the blast radius of the worst plausible mistake here, and does the safeguard match it?

### Scope visibility

**Use when:** Multi-tenant, multi-project, multi-environment, or role-switching products — anywhere "which context am I acting in?" has more than one answer.

**Avoid when:** Single-tenant, single-context tools — don't add scope chrome with nothing to disambiguate.

**Required indicators (persistently visible, not hover-only):**
- Current organization/account (tenant).
- Workspace/project.
- **Environment: production / staging / test** — with strong visual differentiation (color, label) for production.
- Role/permission mode (e.g., "viewing as admin", impersonation banners).
- The specific object being edited.

**Failure modes (the two canonical enterprise disasters):**
- Admin deletes a production item believing they're in staging.
- Support agent edits another customer's workspace while impersonating.

**Verification checklist:**
- [ ] Tenant, workspace, and environment are readable at a glance from any screen.
- [ ] Production is visually distinct from non-production (not just a small label).
- [ ] Impersonation/role-switch states show a persistent banner with one-click exit.
- [ ] Destructive confirmations restate the scope ("Delete *prod-db-1* in **Production / Acme Corp**").

**Related:** [08-navigation-ia.md](08-navigation-ia.md) (wayfinding and scope-safety), [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs).

### Permission-aware UI

**Use when:** Any product where users have differing permissions and could encounter actions they can't perform.

**Behavior contract — the hide vs disable decision rule:**
- **Disable with explanation** when discoverability matters: show the action, disabled, with the reason ("Requires Billing Admin role") and an escalation path ("Request access" / "Contact your admin"). This teaches users what the product can do and how to get there.
- **Hide** only when the action is genuinely irrelevant to the role, or security policy requires non-disclosure of the capability's existence.
- **Never fail silently** because of permissions — an action that appears to work but does nothing (or errors generically) is the worst of all options.

**Decision question:** If this user got the missing permission tomorrow, would they know this feature exists? If they should, disable-with-reason; if it would only confuse, hide.

**Failure modes:**
- Grayed-out buttons with no tooltip or reason ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#disabled-submit-without-explanation) is the form-level cousin).
- Actions that 403 only after the user filled in a form.
- Permission checks inconsistently applied — visible in list view, forbidden in detail view.

**Verification checklist:**
- [ ] Every disabled-by-permission control explains what role/permission is needed.
- [ ] An escalation/request path exists where organizationally appropriate.
- [ ] No permission failure is silent; no permission error arrives only after data entry.

### Bulk operations

**Use when:** Tables/lists where users act on many records at once — the efficiency signature of admin tools, and their biggest error amplifier.

**Required safeguards:**
- **Selection count and scope always visible** ("34 selected").
- **Explicit distinction** between *selected on this page*, *all selected*, and *all matching the filter* — the "all 12,401 matching" expansion must be a separate, deliberate click.
- **Preview affected records/count before commit** — especially for destructive or state-changing operations.
- **Undo window or confirmation** for destructive changes (prefer undo where technically feasible — [04-interaction-design.md](04-interaction-design.md#undoredo-vs-confirmation-dialogs)); confirmation restates count and scope.
- **Progress indication and a partial-failure report** ("212 updated, 3 failed — view failures") with per-item reasons and retry.
- **Audit log entry** recording actor, action, scope, and count.

**Failure modes:**
- "Select all" silently means "all matching" and the user deletes 10× what they saw.
- Bulk job fails halfway with no report of which records changed.
- No way to reverse a bulk mistake short of a support ticket.

**Verification checklist:**
- [ ] Count shown at selection time and restated at confirmation.
- [ ] Page-selection vs all-matching is an explicit, separate action.
- [ ] Destructive bulk ops have undo or typed/structured confirmation proportional to blast radius.
- [ ] Partial failures reported per item, with retry.
- [ ] Bulk actions appear in the audit log.

### Expert accelerators

**Use when:** Frequent expert workflows in productivity/admin tools — the payoff for daily users.

**Patterns:** keyboard shortcuts and command palettes ([14-platform-desktop.md](14-platform-desktop.md#keyboard-shortcuts-design), [14-platform-desktop.md](14-platform-desktop.md#command-palettes)), bulk edit, saved filters/views, templates, recent items, density controls.

**The rule that keeps them honest:** expert accelerators must **never replace the visible novice path** for critical tasks. A shortcut is a second way, not the only way ([01-core-principles.md](01-core-principles.md#7-flexibility-and-efficiency-of-use)).

**Failure modes:**
- Features reachable only by shortcut or palette (invisible to occasional users).
- Saved views that silently change shared defaults for teammates.
- Density controls that break touch targets or accessibility minimums.

**Verification checklist:**
- [ ] Every accelerated action has a discoverable menu/UI path.
- [ ] Shortcuts follow platform conventions and are displayed in menus/tooltips.
- [ ] Saved views distinguish personal vs shared scope.

### Admin destructive actions

**Use when:** Deleting/disabling tenants, users, projects, data, or configuration — actions whose consequences land on *other people*.

**Required content in the confirmation moment:**
- **Object name and scope** — restate exactly what, and in which tenant/environment.
- **Consequence** — what happens, to whom ("14 users lose access immediately").
- **Reversibility** — reversible (grace period, soft delete) or permanent, stated plainly.
- **Dependencies/affected users** — downstream breakage enumerated, not implied.
- **Audit log** entry; for the highest-impact actions, consider step-up authentication ([#login-pattern-matrix](#login-pattern-matrix)) and typed confirmation of the object name.

Copy guidance for these dialogs: [09-content-ux-writing.md](09-content-ux-writing.md#confirmation-and-destructive-action-copy).

**Failure modes:**
- Generic "Are you sure?" that names neither object nor consequence (trains reflexive confirmation).
- Permanent deletion presented identically to reversible archiving.
- No record of who deleted what.

**Verification checklist:**
- [ ] Confirmation names the object, scope, consequence, and reversibility.
- [ ] Friction is proportional: reversible → undo; irreversible + high-impact → typed confirmation and/or step-up auth.
- [ ] Action is audit-logged with actor and timestamp.

---

## Settings & preferences

### Settings information architecture

**Use when:** Organizing any settings surface beyond a handful of toggles.

**Behavior contract:**
- **Group by user mental model, not by implementation service.** "Notifications", "Privacy", "Appearance" — not "Sync Engine" or the names of your microservices. Users search by task language.
- **Make scope explicit** on every settings page: does this apply to the device, the app, the account, the organization, the project, the document, or this session? Ambiguous scope is the number-one settings failure.
- Card-sort or mine support queries to name categories ([17-ux-process-research.md](17-ux-process-research.md)).

**Failure modes:**
- User can't find a setting because the category mirrors engineering structure.
- Same setting appears in two places with different values (scope collision).
- Org-level policy silently overrides a user setting with no indication of why the control is locked.

**Verification checklist:**
- [ ] A user can find any setting using task language, not internal names.
- [ ] Every settings page states its scope.
- [ ] Admin-locked settings show that they're locked, by whom, and why.

### Setting types and save behavior

Match save behavior to setting type — mixing models on one screen is a reliable source of "did it save?" confusion.

| Setting type | Pattern | Save behavior | Notes |
|---|---|---|---|
| Immediate preference | Switch/select | **Instant apply** | Show a saved/synced indicator if persistence matters; switch = instant is the contract ([07-forms-and-input.md](07-forms-and-input.md)) |
| Complex configuration | Form section | **Explicit save** | Multiple related fields change together; support cancel/reset and validation before commit |
| Disruptive/previewable change | Preview | **Apply/Cancel** | Display, layout, theme — let users see the result and back out |
| Dangerous setting | Review + confirmation | **Explicit save/confirm** | Explain consequence and scope before commit |
| Personalization | Preferences panel | Immediate or saved | Respect accessibility preferences and user autonomy |
| Admin policy | Admin settings page | **Explicit save + audit** | Show affected users/scope; log the change |

**Save-model decision questions:**
- **Autosave/instant:** is the change low-risk, reversible, and its feedback immediately visible? Then a save button is pure friction.
- **Explicit save:** are multiple related changes made together, is the change high-impact, or does it need review? Then instant-apply causes half-configured states.
- **Apply/Cancel:** is the change disruptive or previewable (resolution, theme, layout) where users may want to abandon? Provide both, plus revert-on-timeout for changes that can lock users out (display settings).

**Failure modes:**
- Toggle looks instant but actually requires a hidden Save — user leaves, change lost.
- Changes save silently but fail server-side with no error surfaced.
- One screen mixes instant toggles and explicit-save fields with no visual distinction.

**Verification checklist:**
- [ ] Each settings screen uses one save model, or visibly separates sections that differ.
- [ ] Save state (saved / saving / failed) is communicated, not assumed.
- [ ] Dangerous settings state consequence and scope before commit.

### Honest defaults

**Behavior contract:**
- Defaults optimize for **user benefit, safety, and common success** — the default is what most users would choose if fully informed.
- Privacy- and security-relevant defaults are **conservative** (share less, expose less).
- **Never use defaults to trick users** into sharing, buying, or subscribing — preselected consent, marketing opt-ins, and add-ons are the [preselection dark pattern](18-patterns-antipatterns.md#preselection), and preselected consent fails GDPR's active-consent standard.
- **Explain non-obvious defaults** inline ("On by default so teammates can find you").

**Failure modes:**
- Privacy-impacting defaults preselected without clarity.
- Defaults chosen for business metrics that most users would reject if asked.

**Verification checklist:**
- [ ] No consent, purchase, or data-sharing option is preselected in the business's favor.
- [ ] Non-obvious defaults carry a one-line explanation.
- [ ] Security/privacy defaults are the conservative option.

### Reset, undo, and settings search

**Behavior contract:**
- **Reset to default** for complex settings and per-section groups — and state exactly what reset affects before doing it ("Resets appearance settings only; account settings unchanged").
- **Undo** for recent reversible changes where feasible.
- **Search within settings** for large sets (roughly a screen's worth of categories or more): match on synonyms and task language ("dark mode", "theme", "appearance" all hit), and jump straight to the control, highlighted.

**Failure modes:**
- Reset-all destroys hours of configuration with no scoped alternative and no warning about extent.
- Settings search matches only exact internal labels.

**Verification checklist:**
- [ ] Reset is available, scoped, and explains its extent before executing.
- [ ] Products with large settings surfaces have settings search with synonym matching.

---

## Accounts, permissions & authentication UX

Deep authentication mechanics — NIST SP 800-63B password policy, passkeys/WebAuthn, magic links, OTP field behavior, password managers — are covered in [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux); WCAG 3.3.8 accessible-authentication requirements in [05-accessibility.md](05-accessibility.md#accessible-authentication). This section covers the account-level decisions: which pattern for which audience, when to prompt, and how accounts end.

### Login pattern matrix

**Use when:** Choosing the authentication pattern for a product or a specific action class.

| Need | Recommended pattern | Avoid |
|---|---|---|
| Consumer login | Email + password / passkey / magic link | Username-only identity with ambiguous recovery |
| Enterprise login | **SSO (SAML/OIDC) + fallback** | Local password-only where buyers expect SSO |
| High-security action (payout change, data export, deletion, role grants) | **Step-up auth** — re-authenticate or second factor for that action only | Re-auth for every minor action (blanket friction trains bypass habits) |
| Low-friction return | Remembered session with clear sign-out | Hidden persistent sessions on shared devices |
| Account recovery | Email/phone/passkey recovery + human support path | Memory-based security questions (deprecated by NIST — see [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)) |

**Decision questions:**
- Who is the buyer? Enterprise buyers require SSO; shipping without it blocks deals and forces shadow credentials.
- Which actions are sensitive enough for step-up? Budget security friction by consequence, not uniformly ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#security-ux)).
- Is the device likely shared? If yes, session persistence needs explicit consent and visible sign-out.

**Verification checklist:**
- [ ] Authentication works with password managers, paste, keyboard, and screen readers.
- [ ] Step-up auth preserves the in-progress task state.
- [ ] Recovery paths exist that don't depend on memorized answers.

### Signup minimization

**Use when:** Any registration flow.

**Behavior contract:**
- **Ask only for what's needed to create value now.** Every field must justify itself; defer optional profile data to post-signup, in context.
- Explain password/passkey requirements **before** the user errs (live checklist — [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)).
- Verify email/phone at the moment risk requires it — not necessarily as a wall before first value.
- Offer passkeys or social/SSO where the audience supports them; keep an email fallback.

**Failure modes:**
- Long profile-building forms before the user has seen any value.
- Email-verification wall that blocks first use when nothing sensitive is at stake yet.
- Password rules revealed only via rejection errors.

**Verification checklist:**
- [ ] Signup collects the minimum for first value; everything else is deferred.
- [ ] Requirements shown up front; validation is live, not punitive.
- [ ] On registration with an existing email, say so clearly with a sign-in link (registration is the safe place for this — see [Authentication failure modes](#authentication-failure-modes) for why login is not).

### Permission-prompt timing

**Use when:** Requesting OS/browser permissions — notifications, location, camera, contacts.

**Behavior contract:**
- **Ask in context, when value is clear** — never on load. The prompt should arrive attached to a user action whose benefit needs the permission ("Alert me when the price drops" → notification prompt).
- **Pre-permission priming:** explain why before the system prompt when it improves informed consent; the OS prompt is often your only shot (platform mechanics: [15-platform-mobile.md](15-platform-mobile.md#notifications-ux) for mobile, [13-platform-web.md](13-platform-web.md#notification-permission-ux) for web).
- **Provide a fallback if denied** — the feature degrades, it doesn't dead-end.
- **Do not nag after denial** ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#nagging)); offer a settings path for users who change their minds.

**Failure modes:**
- Permission prompt on page/app load ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#notification-permission-on-page-load)) — reflexive denial, often permanent.
- Denial leaves the feature broken with no explanation or alternative.
- Re-prompting loops after explicit denial.

**Verification checklist:**
- [ ] No permission is requested before the user takes an action that needs it.
- [ ] Each permission has a denial fallback path.
- [ ] Denied state is remembered; recovery is via settings, not repeated prompts.

### Roles and permission models

**Use when:** Designing how roles and grants are presented to the humans who assign and hold them.

**Behavior contract:**
- **Role names understandable** by the people assigning them — "Billing Admin", not "role_fin_l2".
- **Permissions summarized by consequence**, not by API scope list: "Can invite and remove members; cannot see billing."
- **Dangerous grants require review** — granting admin/owner shows what the grantee will be able to do before confirm.
- **Visibility of access:** where appropriate, users can see who has access to a resource and at what level.
- Escalation requests flow through the product (see [Permission-aware UI](#permission-aware-ui)), not through out-of-band folklore.

**Failure modes:**
- Role-picker with names nobody can distinguish, leading to over-granting (everyone becomes admin).
- Permission changes take effect silently with no notice to affected users.
- No answer to "who can see this?" for sensitive resources.

**Verification checklist:**
- [ ] A non-technical admin can predict what each role allows from its name + summary.
- [ ] High-privilege grants show consequences at grant time and are audit-logged.
- [ ] Access to a resource is inspectable where the domain expects it.

### Account deletion and data export

**Use when:** Every product with accounts. Under GDPR, erasure (Art. 17) and data portability (Art. 20) are user rights — and the UX bar is: **deletion should be as easy as signup**. Obstruction here is the [roach-motel pattern](18-patterns-antipatterns.md#hard-to-cancel-roach-motel) with legal exposure.

**Behavior contract:**
- Self-service deletion, findable in account settings, same channel as signup — no mandatory phone calls or chat gauntlets.
- **Explain consequences, retention, and reversibility** before commit: what is deleted, what is retained (and the legal basis), whether there's a grace period.
- **Offer export before deletion** — machine-readable, complete, and not a dark archive of unusable blobs.
- **Confirm identity for high-risk deletion** (step-up auth) — this is legitimate friction; retention-begging screens are not.
- **Provide receipt/status:** confirmation of the request, timeline, and completion notice.

**Failure modes:**
- Deletion hidden, gated behind support tickets, or padded with repeated retention offers.
- "Deleted" accounts that silently persist.
- Export that omits key data or arrives in proprietary formats.

**Verification checklist:**
- [ ] Deletion and export are self-service and findable within account settings.
- [ ] Pre-deletion screen states scope, retention, grace period, and offers export.
- [ ] Identity confirmation is proportional (step-up), not obstructive (phone-only).
- [ ] User receives confirmation and completion status.

**Sources:** [GDPR Art. 17 (right to erasure)](https://gdpr-info.eu/art-17-gdpr/), [GDPR Art. 20 (data portability)](https://gdpr-info.eu/art-20-gdpr/).

### Authentication failure modes

The recurring implementation bugs that break real users — each is a review checklist item:

- **Blocking password managers** (autofill-hostile fields, non-standard inputs) — violates WCAG 3.3.8 ([05-accessibility.md](05-accessibility.md#accessible-authentication)) and NIST guidance ([07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)).
- **Disabling paste** in password/OTP/card fields — same violations, punishes exactly the users doing security right.
- **Split OTP digit boxes that break autofill** — one-input-per-digit defeats `autocomplete="one-time-code"` OS autofill and paste; prefer a single field ([07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux)).
- **Permission prompts before value is understood** (see [Permission-prompt timing](#permission-prompt-timing)).
- **Account enumeration via error messages:** login errors that reveal whether an account exists ("no account with this email") let attackers harvest valid accounts — on *login*, say "email or password didn't match"; on *registration*, do say the email is taken (with a sign-in link), since that channel is rate-limitable and the honest message prevents duplicate-account confusion.
- **Cognitive-function tests without accessible alternatives** (puzzle CAPTCHAs, memorized security questions) — WCAG 3.3.8/3.3.9 territory ([05-accessibility.md](05-accessibility.md#accessible-authentication)).

**Verification checklist:**
- [ ] Paste and autofill work in every credential field (password, OTP, card).
- [ ] OTP is a single field with `one-time-code` autofill, chunked visually if needed.
- [ ] Login errors don't confirm account existence; registration handles duplicates honestly.
- [ ] No auth step requires a cognitive test without an accessible alternative.

---

## Quick reference

| Topic | One-line rule | Details |
|---|---|---|
| Listing pages | Cards carry the comparison-critical facts; zero-result filters offer recovery | [#product-listing-and-category-pages](#product-listing-and-category-pages) |
| PDP | Price + all fees, variants, availability, delivery/returns, one primary CTA | [#product-detail-page-pdp](#product-detail-page-pdp) |
| Cart | Editable workspace; flag changes; never sneak items in | [#cart](#cart) |
| Guest checkout | Default path; account creation offered *after* purchase | [#checkout-flow](#checkout-flow) |
| Total cost | Shipping/taxes/fees as early as determinable, always before commit (Baymard: surprise costs = top abandonment reason) | [#checkout-flow](#checkout-flow) |
| Final button | Names the outcome: "Pay $42.00", never "Continue" | [#checkout-flow](#checkout-flow) |
| Billing address | Same-as-shipping pre-checked, one click to change | [#checkout-flow](#checkout-flow) |
| Conversion ethics | Persuasion survives full user understanding; deception doesn't | [#conversion-ethics](#conversion-ethics) |
| Scope visibility | Tenant + workspace + environment always visible; production visually loud | [#scope-visibility](#scope-visibility) |
| Hide vs disable | Disable-with-reason when discoverability matters; hide when irrelevant or policy demands | [#permission-aware-ui](#permission-aware-ui) |
| Bulk operations | Count + preview + undo/confirm + partial-failure report + audit log | [#bulk-operations](#bulk-operations) |
| Expert accelerators | Second way, never the only way | [#expert-accelerators](#expert-accelerators) |
| Admin deletes | Name object/scope/consequence/reversibility; audit; step-up for worst cases | [#admin-destructive-actions](#admin-destructive-actions) |
| Settings IA | Group by user mental model; state scope on every page | [#settings-information-architecture](#settings-information-architecture) |
| Save behavior | Toggles apply instantly; forms save explicitly; disruptive changes get apply/cancel | [#setting-types-and-save-behavior](#setting-types-and-save-behavior) |
| Defaults | User-benefit only; privacy defaults conservative; never preselected consent/add-ons | [#honest-defaults](#honest-defaults) |
| Reset & search | Scoped reset with stated extent; settings search for large sets | [#reset-undo-and-settings-search](#reset-undo-and-settings-search) |
| Login pattern | Consumer: email/passkey/magic link; enterprise: SSO+fallback; sensitive actions: step-up | [#login-pattern-matrix](#login-pattern-matrix) |
| Signup | Minimum fields for first value; defer the rest | [#signup-minimization](#signup-minimization) |
| Permission prompts | In context, primed, with denial fallback; never on load | [#permission-prompt-timing](#permission-prompt-timing) |
| Roles | Named for humans, summarized by consequence; dangerous grants reviewed | [#roles-and-permission-models](#roles-and-permission-models) |
| Deletion/export | As easy as signup; explain scope/retention; export first; receipt (GDPR Art. 17/20) | [#account-deletion-and-data-export](#account-deletion-and-data-export) |
| Auth bugs | No paste-blocking, no split-OTP autofill breakage, no enumeration via login errors | [#authentication-failure-modes](#authentication-failure-modes) |
| Auth doctrine | Full NIST/passkey/OTP coverage lives in file 07 | [07-forms-and-input.md](07-forms-and-input.md#password-and-auth-ux) |
