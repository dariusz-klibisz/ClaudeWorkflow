# AI Product UX & Human-AI Interaction

> Part of the [UI/UX Reference](00-index.md). See index for the full documentation map.

## Scope

This file gives decision guidance for products that include AI-generated, predictive, automated, or agentic behavior: which interaction paradigm fits which task, how to manage the human-AI relationship across its lifecycle (initial expectations, during interaction, when the AI is wrong, over time), how to tier AI features by risk and escalate safeguards accordingly, and the concrete UX of prompts, outputs, agentic actions, accessibility, and support chatbots. Two published frameworks anchor this guidance: the [Microsoft HAX Guidelines for Human-AI Interaction](https://www.microsoft.com/en-us/haxtoolkit/) (18 guidelines organized by lifecycle phase) and the [Google PAIR People + AI Guidebook](https://pair.withgoogle.com/guidebook/) — the lifecycle structure below follows HAX's phase organization without reproducing its numbered guidelines. Generic loading/feedback mechanics live in [04-interaction-design.md](04-interaction-design.md#loading-and-waiting); accessibility fundamentals in [05-accessibility.md](05-accessibility.md); AI as a *design-process* tool (not a product feature) in [17-ux-process-research.md](17-ux-process-research.md#ai-assisted-design-and-research).

## Contents

- [Core AI UX rules](#core-ai-ux-rules)
- [Interaction paradigms](#interaction-paradigms)
- [AI risk tiers](#ai-risk-tiers)
- [Human-AI interaction lifecycle](#human-ai-interaction-lifecycle)
  - [Initially: set expectations](#initially-set-expectations)
  - [During interaction](#during-interaction)
  - [When the AI is wrong](#when-the-ai-is-wrong)
  - [Over time](#over-time)
- [Prompt UX](#prompt-ux)
- [Output states and presentation](#output-states-and-presentation)
- [Agentic actions](#agentic-actions)
- [AI accessibility](#ai-accessibility)
- [AI chatbots for support](#ai-chatbots-for-support)
- [AI-assisted design and research](#ai-assisted-design-and-research)
- [Acceptance criteria](#acceptance-criteria)
- [Quick reference](#quick-reference)

---

## Core AI UX rules

These apply to every AI feature regardless of paradigm or risk tier:

- Communicate what the AI can and cannot do.
- Set expectations before users rely on output.
- Show uncertainty, provenance, or confidence when trust matters.
- Keep users in control of high-risk outcomes.
- Make AI outputs editable, rejectable, and reversible where feasible.
- Provide fallback paths when AI fails.
- Do not blame the user for model failure.
- Evaluate AI behavior continuously with real tasks and failure cases.

**Sources:** [Microsoft HAX Toolkit](https://www.microsoft.com/en-us/haxtoolkit/), [Google PAIR Guidebook](https://pair.withgoogle.com/guidebook/).

---

## Interaction paradigms

**Definition:** The four ways users direct AI-capable systems, from full manual control to full delegation. Choosing the wrong paradigm for the task is the root cause of most AI-feature failures — the matrix decides before any screen is designed.

| Paradigm | Use when | Risks | Required safeguards |
| --- | --- | --- | --- |
| **Command-based GUI** (no AI, or AI invisible) | User needs precision, control, and predictable results | Slow for complex or generative tasks | Direct manipulation, undo, clear state |
| **Intent-based prompting** | User can describe the desired outcome but not the steps | Ambiguous prompts, hallucination, blank-page paralysis | Examples, constraints, previews, refinement loop |
| **Hybrid GUI + AI** | User needs both generation speed and fine control | Interface complexity; unclear which layer owns state | Editable outputs, structured controls, explainable state |
| **Autonomous agent** | System acts over time on the user's behalf | Trust, safety, runaway actions | Permissions, logs, confirmations, limits, stop control |

**Decision questions:**

- Can the user specify the exact result faster than they could describe it? → Command GUI.
- Is the outcome describable but the production laborious? → Intent-based or hybrid.
- Does the user need to adjust the output after generation? → Hybrid (never prompt-only for editable artifacts).
- Does value come from acting *without* the user present? → Agent — apply [Agentic actions](#agentic-actions) rules and the [risk-tier](#ai-risk-tiers) safeguards.

**When not to use AI at all:** deterministic tasks with one correct answer (calculation, lookup, format conversion) where a conventional control is faster, more reliable, and cheaper; tasks where an error is unacceptable and unverifiable by the user.

**Failure modes:** prompt-only interfaces for tasks users could click through in seconds; "AI-washing" a working GUI into a chatbot; agents launched without stop controls; hybrid interfaces where AI edits silently overwrite user edits.

**Verification checklist:**
- [ ] Paradigm choice documented against the task, not the technology.
- [ ] A non-AI path exists for critical tasks.
- [ ] Editable artifacts can be edited directly, not only re-prompted.
- [ ] Agent-paradigm features meet every rule in [Agentic actions](#agentic-actions).

---

## AI risk tiers

**Rule:** Classify every AI feature into a risk tier before designing it; required safeguards escalate with tier and are cumulative (each tier includes everything from the tiers below).

| Tier | Examples | Required UX safeguards |
| --- | --- | --- |
| **1 — Low** (suggestions, easily dismissed) | Summarize notes, rewrite tone, autocomplete | Disclose AI generation; edit/reject controls |
| **2 — Medium** (influences decisions) | Recommend a workflow, classify a support ticket, rank results | + Confidence/uncertainty indication, correction mechanism, feedback channel, **preview before apply** |
| **3 — High** (consequential output) | Financial/medical/legal drafting, public posting, messages sent as the user | + Source/provenance display, **explicit confirmation**, human review step, hallucination disclaimer |
| **4 — Critical** (autonomous high-impact actions) | Autonomous transactions, account deletion, health triage, irreversible external actions | + **Human approval before execution**, scoped permissions, **audit log**, human escalation path, kill switch |

**Decision questions:**

- What is the worst plausible outcome of a wrong output? Who bears it?
- Can the user detect the error before harm occurs? If not, move up a tier.
- Is the action reversible? Irreversible + external = tier 4 treatment.
- Would the organization defend this output in a dispute? (See [AI chatbots for support](#ai-chatbots-for-support) — courts have held companies liable for AI output.)

**Failure modes:** tier-4 actions shipped with tier-1 safeguards ("the model is usually right"); disclaimers used as a substitute for safeguards at high tiers; confirmation fatigue from applying tier-3 friction to tier-1 suggestions (which trains users to click through the confirmations that matter).

**Verification checklist:**
- [ ] Every AI feature has a documented tier.
- [ ] Safeguards match or exceed the tier's requirements.
- [ ] No tier-3/4 action executes without explicit user confirmation or approval.
- [ ] Tier-4 features have audit logs and a tested human-escalation path.

---

## Human-AI interaction lifecycle

The HAX guidelines organize human-AI interaction into four phases: initially, during interaction, when wrong, and over time ([HAX Toolkit](https://www.microsoft.com/en-us/haxtoolkit/)). Design for all four — most teams design only the happy path of phase two.

### Initially: set expectations

**Rule:** Before users rely on output, make clear what the system can do, how well it does it, and what it does with their data.

**Implementation guidance:**
- Explain the AI's role in plain language ("Drafts a reply you review and send" — not "AI-powered magic").
- State capability *and* error expectation: what it's good at, where it's unreliable.
- Provide examples, starter prompts, or templates so the first interaction succeeds.
- State limitations and data use (is input used for training? retained? shared?).
- Offer a non-AI path where the task is critical.

**Failure modes:** marketing copy promising accuracy the model can't deliver (destroys trust on first error); zero-state blank prompt boxes with no examples; burying data-use terms.

**Verification checklist:**
- [ ] A first-time user can state what the feature does and how reliable it is.
- [ ] Examples or templates are visible at the empty state.
- [ ] Data use is disclosed before first input.
- [ ] Critical tasks have a visible non-AI alternative.

### During interaction

**Rule:** Show what the system is doing, why, and with what confidence; keep the user able to steer, stop, and refine.

**Implementation guidance:**
- Show progress for long generation ([04-interaction-design.md](04-interaction-design.md#loading-and-waiting) for thresholds and skeleton/streaming mechanics).
- Allow cancel/stop at any time.
- Support refinement without starting over — carry the conversation/context forward.
- Preserve prior context visibly when it affects output (which document, which selection, which history the model saw).
- Show confidence or uncertainty when it affects the user's decision; make clear *why* the system produced what it did (sources used, rules applied, context considered) at tier 2+.
- Match displayed confidence to actual reliability — fabricated confidence scores are worse than none.

**Failure modes:** unexplained multi-second silences; outputs that ignore visible context without saying so; uniform authoritative tone regardless of actual certainty; no way to interrupt generation.

**Verification checklist:**
- [ ] Generation shows progress and can be cancelled.
- [ ] Refinement builds on prior output instead of restarting.
- [ ] Context affecting the output is visible or inspectable.
- [ ] Uncertainty is surfaced where the risk tier requires it.

### When the AI is wrong

**Rule:** Wrong output is a designed-for state, not an edge case. Support efficient correction, dismissal, and reporting — and never blame the user.

**Implementation guidance:**
- Allow correction and reporting in place (thumbs-down + structured reason, inline edit).
- Offer regenerate, refine, and manual-edit paths — escalating user control.
- Show the fallback path (manual mode, human contact, non-AI feature).
- Error copy owns the failure: "I couldn't find that" — never "your prompt was unclear."
- Make dismissal cheap: one action to reject a suggestion, no justification demanded.

**Failure modes:** regenerate as the *only* recovery (gambling, not correcting); feedback buttons that go nowhere; failures that dead-end without a fallback; retry loops that lose the user's input ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#form-data-loss)).

**Verification checklist:**
- [ ] Every output can be corrected, edited, or rejected.
- [ ] A fallback path exists and is discoverable from the failure state.
- [ ] Error copy does not blame the user.
- [ ] Feedback is collected and routed to evaluation.

### Over time

**Rule:** If the system learns from behavior, it must do so transparently and under user control.

**Implementation guidance:**
- Let users adjust personalization — including turning it off.
- Show what the system learned when it affects results ("Because you usually…" with an edit control).
- Allow reset/clear history where appropriate.
- Notify users of significant capability changes; don't silently change behavior they've calibrated trust against.
- Monitor drift and recurring failure patterns; recurring tier-2+ failures are a product incident, not a model quirk.

**Failure modes:** invisible personalization users can't inspect or correct; learned behavior that persists after the user's context changes; model updates that silently regress workflows.

**Verification checklist:**
- [ ] Personalization is visible, adjustable, and resettable.
- [ ] Behavior changes from learning are attributable ("why am I seeing this?").
- [ ] Drift and failure recurrence are monitored with owners assigned.

---

## Prompt UX

**Definition:** The input surface for intent-based interaction: affordances, examples, constraints, history, and re-run mechanics.

**Rules:**
- Provide structured inputs where possible instead of a blank prompt alone — dropdowns, chips, and toggles for tone, length, audience, format, and sources beat free-text for constrained parameters.
- Use examples to teach expected prompt shape; placeholder examples and starter-prompt galleries reduce blank-page failure.
- Let users constrain output explicitly: tone, length, audience, format, sources.
- Keep **prompt history**: previous prompts retrievable, searchable where volume warrants.
- Make prompts **editable and re-runnable**: users iterate by editing the prior prompt, not retyping it; an edited re-run should show its relationship to the original.
- Provide previews before commit for anything at tier 2+.
- Never require articulate prose — structured alternatives serve non-native speakers, motor-impaired users dictating input, and everyone in a hurry.

**When not to use free-text prompting:** parameter choices with known valid values; tasks where users can't evaluate whether the prompt was "good"; high-frequency repeat tasks (promote them to buttons).

**Failure modes:** blank textarea as the entire interface; prompts lost on navigation or error ([18-patterns-antipatterns.md](18-patterns-antipatterns.md#form-data-loss) applies fully to prompts); no way to see what prompt produced a saved output; forcing prompt-craft expertise the product should encode.

**Verification checklist:**
- [ ] Empty state shows examples or templates.
- [ ] Common constraints have structured controls.
- [ ] Prompt history exists; prompts survive errors and navigation.
- [ ] Any prior prompt can be edited and re-run.

---

## Output states and presentation

**Rule:** Design every output state, not just the successful draft.

**Required states:**

| State | Requirements |
| --- | --- |
| Loading/generating | Progress indication, cancel control ([04-interaction-design.md](04-interaction-design.md#loading-and-waiting)) |
| **Streaming** | Progressive text render; stop button; layout that doesn't jump; see [AI accessibility](#ai-accessibility) for screen-reader handling |
| **Partial result** | Clearly marked incomplete; continue/retry affordance |
| Generated draft | Marked as AI-generated; edit, copy, regenerate, refine controls |
| Low-confidence result | Uncertainty flagged; alternatives or verification prompt offered |
| **Citations/sources** | Provenance links for factual claims at tier 2+; "source unavailable" as an honest state, never a fabricated citation |
| Failed generation | Owns the failure; retry + fallback path; input preserved |
| User-edited result | Distinguishable from raw AI output; user edits never silently overwritten |

**Hallucination disclaimers — proportional to risk:** a blanket "AI can make mistakes" line is table stakes at tier 1 and *insufficient* at tier 3+, where the UI must actively surface provenance and require verification steps (e.g., "check the cited policy before sending"). Disclaimers do not substitute for safeguards, and courts have not accepted them as liability shields (see [AI chatbots for support](#ai-chatbots-for-support)).

**Regenerate/refine controls:** regenerate (new attempt), refine (directed change: "shorter," "more formal"), and manual edit form an escalating-control ladder; offer all three for editable artifacts. Show which version the user is looking at when multiple generations exist.

**Failure modes:** citations that don't support the claims they decorate; streaming output that reflows and shifts already-read text; regenerate discarding a user's manual edits without warning; confident tone on low-confidence output.

**Verification checklist:**
- [ ] All eight states above are designed and reachable.
- [ ] Factual outputs at tier 2+ carry provenance or an honest "no source" state.
- [ ] Disclaimer strength scales with the feature's risk tier.
- [ ] Regenerate warns before destroying user edits.

---

## Agentic actions

**Rule:** An agent acting on the user's behalf requires more transparency and control than a tool the user drives — never less.

**Implementation guidance:**
- **Plan preview before execution:** for medium/high-risk actions, show the planned steps and let the user approve, edit, or reject the plan — not just the final outcome.
- **Explicit scopes and permissions:** ask permission before acting outside the current context; scope grants narrowly (this task, this account, this time window) rather than blanket authority.
- **Visible activity log:** every action the agent takes is recorded and reviewable — what, when, why, and on whose authorization. Tier 4 requires a durable audit log.
- **Stop/pause control:** always available, always effective mid-run.
- **Undo where possible:** prefer reversible actions; where irreversibility is unavoidable (external sends, transactions, deletions), require explicit confirmation per the [risk tiers](#ai-risk-tiers).
- **Limits:** rate, spend, and blast-radius caps enforced by the system, not by prompt instructions.

**Decision questions:**

- Would the user be surprised by any single action in this plan? Then the plan needs preview.
- Can the user reconstruct what the agent did a week later? Then the log is sufficient.
- What happens if the agent runs with a wrong premise for an hour? Then limits and stop controls are sized correctly — or not.

**Failure modes:** agents that report success without an inspectable record; "confirm plan" screens too vague to evaluate ("I'll organize your files"); stop buttons that only stop *new* actions while queued ones fire; permission asked once, exercised forever.

**Verification checklist:**
- [ ] Medium/high-risk plans preview before execution.
- [ ] Activity log covers every externally visible action.
- [ ] Stop/pause halts the run, verifiably, mid-execution.
- [ ] Permissions are scoped, inspectable, and revocable.
- [ ] Irreversible/external actions confirm explicitly; reversible ones offer undo.

---

## AI accessibility

**Rule:** AI features and AI-generated content must meet the same accessibility bar as everything else — see [05-accessibility.md](05-accessibility.md) for the foundation.

**Implementation guidance:**
- **Streaming vs screen readers:** do not pipe token-by-token streams into an `aria-live` region — it spams the queue and renders output unintelligible. Announce meaningful chunks (paragraph or message boundaries) or announce *completion* ("Response ready, 3 paragraphs") and let the user read at their own pace. See [05-accessibility.md](05-accessibility.md#aria) for live-region politeness rules.
- **Keyboard control of generation:** starting, stopping, and regenerating must all be keyboard-operable; focus lands predictably on completed output or its controls, never lost during streaming ([05-accessibility.md](05-accessibility.md#keyboard-navigation)).
- Chat interfaces need full keyboard and screen-reader support: message history navigable, roles (user vs AI) programmatically distinguishable, input reachable.
- Generated content must itself be accessible: real headings, alt text for generated images, captions for generated audio/video, semantic structure — generate it accessible, don't ask users to repair it.
- Don't require articulate prose from users; structured alternatives ([Prompt UX](#prompt-ux)) are an accessibility feature for cognitive, motor, and language needs.

**Failure modes:** live regions announcing every token; focus stolen to the output area mid-typing; stop buttons that are pointer-only; generated images shipped with empty alt.

**Verification checklist:**
- [ ] Streaming output tested with a screen reader; announcements are chunked or completion-based.
- [ ] Generate/stop/regenerate all work by keyboard with sane focus order.
- [ ] Generated artifacts pass the same accessibility checks as authored content.
- [ ] A structured (non-prose) input path exists.

---

## AI chatbots for support

**Definition:** LLM or scripted chatbots as the front line — or the entirety — of customer support. A controversial pattern with legitimate wins and documented, legally consequential failures; also cataloged in [18-patterns-antipatterns.md](18-patterns-antipatterns.md#ai-chatbots-as-primary-support).

**Reasoning / Evidence:** Legitimate wins: instant 24/7 answers for the FAQ head (deflection of genuinely simple queries), language coverage, and tier-1 triage. Documented failures: hallucinated policies — Air Canada was held **legally liable** for its chatbot's invented bereavement-refund policy (*Moffatt v. Air Canada*, BC Civil Resolution Tribunal, 2024; the tribunal rejected the argument that the chatbot was a "separate legal entity" responsible for its own statements) — obstruction-by-bot (the hard-to-cancel variant where the bot is the maze — [18-patterns-antipatterns.md](18-patterns-antipatterns.md#obstruction)), user rage at escalation-blocking, and trust damage when bots impersonate humans (DSA/FTC transparency direction: disclose the bot).

**Verdict:** Good for tier-1 triage **with clear escape hatches to humans**. Never as the only support channel.

**When to use:** FAQ-head triage with *instant, visible* human escalation; drafted-by-AI, human-reviewed replies; internal support tooling.
**When not to use:** Sole channel for consequential issues (billing disputes, safety, cancellations — regulators increasingly mandate human paths); unreviewed authoritative policy statements.

**Pros:** Latency, scale, cost, language coverage. **Cons:** Hallucination liability, escalation rage, trust erosion, complex-case failure.

**Implementation guidance:**
- Disclose bot status upfront; never impersonate a human.
- Escalate to a human within ≤2 exchanges on request or detected frustration; the escape hatch is visible from the first message.
- Ground responses in versioned policy sources (RAG with citations); the bot states policy only from documents the company stands behind — *the company is liable for what the bot says*.
- Log and audit answer accuracy ([risk tier](#ai-risk-tiers) 3: the bot makes claims users act on).
- Measure *resolution*, not deflection — deflection is a cost metric wearing a savings costume ([17-ux-process-research.md](17-ux-process-research.md)).
- Tone: helpful, honest about being a bot ([09-content-ux-writing.md](09-content-ux-writing.md#voice-and-tone)).

**Failure modes:** bot-as-maze blocking the human path; chatbot confidently inventing policy (now established legal liability, not just UX debt); measuring success by contact deflection while resolution craters.

**Verification checklist:**
- [ ] Bot discloses itself in the first message.
- [ ] Human escalation reachable within two exchanges, and always on request.
- [ ] Policy answers cite versioned source documents.
- [ ] Accuracy audited; resolution (not deflection) is the success metric.
- [ ] Consequential issues (billing, safety, cancellation) have non-bot channels.

**Sources:** [Air Canada chatbot ruling (CBC)](https://www.cbc.ca/news/canada/british-columbia/air-canada-chatbot-lawsuit-1.7116416), [NN/g: AI chat UX](https://www.nngroup.com/articles/ai-chat-ux/).

---

## AI-assisted design and research

This file covers AI *in the product*. For AI in the design **process** — LLM-assisted research analysis, generative UI drafting, AI evaluation tools, and the synthetic-users caution (LLM-simulated participants are invalid as user evidence; they produce plausible averages, not your users' behavior) — see [17-ux-process-research.md](17-ux-process-research.md#ai-assisted-design-and-research).

---

## Acceptance criteria

- User can understand the AI's role and limitations before relying on it.
- User can inspect and control AI output before high-risk use.
- User can recover from wrong output via correction, regeneration, or fallback.
- Data use and privacy implications are clear at the point of input.
- Agentic actions are previewed, logged, stoppable, and scoped.
- AI features have measurable evaluation and a working feedback loop.

## Quick reference

| Topic | Key rule | Escalation trigger |
| --- | --- | --- |
| Paradigm choice | Match to task: precision→GUI, describable outcome→prompt, editable artifact→hybrid, delegation→agent | Never AI-wash a working GUI |
| Risk tiers | Tier the feature; safeguards are cumulative: disclose→preview/confidence→confirm+provenance→approve+audit | Irreversible + external = tier 4 |
| Initially | Set capability *and* error expectations; examples at empty state | Critical task → non-AI path required |
| During | Progress, cancel, refinement without restart, visible context | Tier 2+ → show confidence and why |
| When wrong | Correct/edit/reject in place; fallback path; never blame the user | Regenerate-only recovery is a failure |
| Over time | Personalization visible, adjustable, resettable | Silent behavior change breaks calibrated trust |
| Prompt UX | Structured inputs + examples + history + editable re-runs | Blank textarea alone is a failure |
| Output states | Design all states incl. streaming, partial, low-confidence, failed | Disclaimer strength ∝ risk tier |
| Citations | Provenance at tier 2+; honest "no source" state | Fabricated citations = never |
| Agentic actions | Plan preview, activity log, stop control, scoped permissions, undo | Irreversible action → explicit confirm |
| Accessibility | Chunk/completion announcements, keyboard generation control | Token-by-token live regions = never |
| Support chatbots | Tier-1 triage only; human within ≤2 exchanges; RAG-grounded policy | Company is liable for bot claims (Air Canada, 2024) |
| Design-process AI | See [17-ux-process-research.md](17-ux-process-research.md#ai-assisted-design-and-research) | Synthetic users ≠ user evidence |
