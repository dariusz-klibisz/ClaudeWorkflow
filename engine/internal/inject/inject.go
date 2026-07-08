// Package inject renders the context payloads (05 §5): the SessionStart
// status block (≤60 lines), the per-turn UserPromptSubmit anchor (≤10 lines),
// and the SubagentStart reviewer briefing. All content is regenerated from
// disk state — the conversation is a cache; these repopulate it.
// Formatting rule: factual, declarative statements only (no imperative
// "SYSTEM:" framing — prompt-injection defenses, 01 §1).
package inject

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/doctor"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

// Session renders the full re-anchoring block.
func Session(c *runctl.Ctl) (string, error) {
	r, err := c.Store.LoadRun()
	if err != nil {
		return "", err
	}
	if r == nil {
		return "[wf] no active run. Start one with /wf:dev (wf run start --family diff|artifact|assessment --intent …). State lives in .workflow/; this block is regenerated from disk.", nil
	}
	var b strings.Builder
	phasePos := phaseIndex(c, r.Family, r.Phase)
	fmt.Fprintf(&b, "[wf] run %s · %s/%s · phase: %s (%s) · started %s · status: %s\n",
		r.ID, r.Family, orDash(r.Intent), orDash(r.Phase), phasePos, ago(r.Started), r.Status)

	if r.Status == "parked" {
		b.WriteString("run is parked — resume with `wf run resume`, or start a branch (/wf:dev)\n")
		return b.String(), nil
	}
	if r.Phase == "" {
		b.WriteString("all phases complete — close with `wf run close`\n")
		return b.String(), nil
	}

	env, err := c.Env(r)
	if err != nil {
		return "", err
	}
	findings, err := contracts.EvaluatePhase(env, r.Phase)
	if err != nil {
		return b.String() + fmt.Sprintf("contract evaluation unavailable (%v) — run `wf doctor`\n", err), nil
	}

	// missing phase inputs (adopt/resume/force landings) — early so the line
	// survives Turn's 10-line window
	if entry, err := contracts.EvaluateEntry(env, r.Phase); err == nil && len(entry) > 0 {
		fmt.Fprintf(&b, "⚠ phase inputs missing: %s → %s (deliberate skip: wf contract waive %s --reason …)\n",
			entry[0].ID, entry[0].Remediation, entry[0].ID)
	}

	agentItems, userItems := split(findings)
	switch {
	case len(findings) == 0:
		fmt.Fprintf(&b, "phase contract met — exit with `wf phase exit`\n")
	case len(agentItems) == 0:
		fmt.Fprintf(&b, "waiting-on: the user — %s\n", userItems[0].Remediation)
	default:
		fmt.Fprintf(&b, "waiting-on: nothing — %d contract item(s) open\n", len(findings))
	}

	if len(findings) > 0 {
		b.WriteString("next actions:\n")
		n := 0
		for _, f := range append(agentItems, userItems...) {
			n++
			if n > 5 {
				fmt.Fprintf(&b, "  … and %d more (wf status)\n", len(findings)-5)
				break
			}
			line := f.Remediation
			if f.Detail != "" {
				line += " [" + f.Detail + "]"
			}
			fmt.Fprintf(&b, "  %d. %s → %s\n", n, f.ID, line)
		}
	}

	open, total := taskCounts(env)
	fmt.Fprintf(&b, "open tasks: %d/%d · loops: %d · forces: %d\n", open, total, r.Loops, r.Forces)
	if warn := doctor.HookLiveness(c, r); warn != "" {
		fmt.Fprintf(&b, "⚠ %s\n", warn)
	}
	fmt.Fprintf(&b, "resume procedure: /wf:%s   escapes: /wf:park /wf:force\n", skillFor(c, r.Phase))
	b.WriteString("state lives in .workflow/ — after compaction or resume, this block (not memory) is authoritative\n")
	return b.String(), nil
}

// Turn renders the ≤10-line per-prompt anchor: the head of the session block.
func Turn(c *runctl.Ctl) (string, error) {
	full, err := Session(c)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(full, "\n"), "\n")
	max := 10
	if len(lines) < max {
		max = len(lines)
	}
	return strings.Join(lines[:max], "\n"), nil
}

// Agent renders the SubagentStart briefing: scope, corpus routing, and —
// for gating reviewers — the verdict-block contract (04 §4). Author-side
// agents with a roster corpus (designer, ux-designer, implementer) get the
// scope + corpus half; without it their corpus routing is dead weight.
func Agent(c *runctl.Ctl, agentName string) (string, error) {
	r, err := c.Store.LoadRun()
	if err != nil || r == nil {
		return "", err
	}
	ag, ok := c.Spec.AgentByName(agentName)
	if !ok || (!ag.Gating && len(ag.Corpus) == 0) {
		return "", nil
	}
	var b strings.Builder
	role := "work"
	if ag.Gating {
		role = "review"
	}
	fmt.Fprintf(&b, "[wf] %s scope: run %s (%s/%s), phase %s.\n", role, r.ID, r.Family, orDash(r.Intent), r.Phase)
	if ag.Gating {
		if mode := agentMode(agentName, r.Phase); mode != "" {
			fmt.Fprintf(&b, "assigned mode/scope for this spawn: %s — include it as the `scope:` line of the verdict block.\n", mode)
		}
		if agentName == "compliance-reviewer" {
			complianceBriefing(c, &b)
		}
	} else if agentName == "implementer" {
		implementerBriefing(c, r, &b)
	} else if stage, reentry := designStage(c, r, agentName); stage != "" {
		fmt.Fprintf(&b, "assigned design stage for this spawn: %s — name it in the `stage` field of the returned option-set.\n", stage)
		if reentry {
			priorDesignFindings(c, r, &b)
		}
	}
	if len(ag.Corpus) > 0 {
		root := pluginRoot(c)
		verb := "authoring"
		if ag.Gating {
			verb = "judging"
		}
		fmt.Fprintf(&b, "reference corpus for this %s (read before %s; cite file+rule):\n", role, verb)
		// Briefing paths render /-separated on EVERY platform: Windows
		// tools accept them, and the agent-visible text stays deterministic
		// (filepath.Join's backslashes broke the corpus assertions on
		// Windows CI).
		for _, p := range ag.Corpus {
			fmt.Fprintf(&b, "  - %s\n", filepath.ToSlash(filepath.Join(root, filepath.FromSlash(p))))
		}
		b.WriteString("corpus absent/unreadable ⇒ use your own knowledge and say so in the output.\n")
	}
	if ag.Gating {
		b.WriteString("The verdict is machine-parsed. The final message must end with exactly this fenced block:\n")
		b.WriteString("```verdict\nstatus: <clean|changes-required|safe|risky|unsafe|n/a>\ncriticals: <int>\nmajors: <int>\nscope: <your assigned mode/lens, if any>\nreason: <required for n/a — one line: why this review does not apply>\n```\n")
		b.WriteString("clean/safe require criticals=0 and majors=0; risky requires each concern listed for disposition; n/a needs one line of reason above the block.\n")
	}
	return b.String(), nil
}

// complianceBriefing names the standards in force (from installed pack
// items) and routes the pack-shipped checklist documents — these live
// project-side under .workflow/packs/, not in the plugin corpus.
func complianceBriefing(c *runctl.Ctl, b *strings.Builder) {
	stds := c.Spec.ComplianceStandards()
	if len(stds) == 0 {
		b.WriteString("no regulated standards are in force (no installed pack references you) — verdict n/a with that reason.\n")
		return
	}
	fmt.Fprintf(b, "standards in force for this project: %s", strings.Join(stds, ", "))
	if len(stds) > 1 {
		b.WriteString(" — one standard per spawn; declare yours in the verdict `scope:` line")
	}
	b.WriteString(".\n")
	packsDir := filepath.Join(c.Store.Root, "packs")
	var docs []string
	_ = filepath.Walk(packsDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".md") {
			docs = append(docs, path)
		}
		return nil
	})
	if len(docs) > 0 {
		b.WriteString("standard checklists (installed with the packs — read before judging; cite clause IDs):\n")
		for _, d := range docs {
			fmt.Fprintf(b, "  - %s\n", filepath.ToSlash(d)) // /-separated on every platform, like the corpus routes
		}
	} else {
		b.WriteString("no pack checklists found under .workflow/packs — review from your own knowledge, say so, and cap status at changes-required.\n")
	}
}

// agentMode mirrors the verdict gate's phase-derived default scopes.
func agentMode(agent, phase string) string {
	if agent == "adversary" {
		switch phase {
		case "frame":
			return "abuse-case"
		case "design":
			return "attack-tree"
		default:
			return "red-team"
		}
	}
	return ""
}

// designStage derives the stage a designer spawn is for from the recorded
// option-sets — deterministic staging: system is fixed before software
// (mirrors agentMode's phase-derived scopes; there is no per-spawn
// addressing channel). ux-designer always works the ux stage. The second
// return marks a loop re-entry (both stages already recorded).
func designStage(c *runctl.Ctl, r *store.Run, agent string) (string, bool) {
	switch agent {
	case "ux-designer":
		return "ux", false
	case "designer":
	default:
		return "", false
	}
	env, err := c.Env(r)
	if err != nil {
		return "", false
	}
	stages := map[string]bool{}
	for _, rec := range env.Records("option-set") {
		if s, _ := rec.Data["stage"].(string); s != "" {
			stages[s] = true
		}
	}
	switch {
	case !stages["system"]:
		return "system", false
	case !stages["software"]:
		return "software", false
	default:
		return "system or software (loop re-entry — state which; carry the recorded rejected option IDs forward, a rejected option may not reappear)", true
	}
}

// priorDesignFindings surfaces the finding CONTENT of the latest failing
// design-stage verdicts on a loop re-entry — the rework must see WHAT
// failed, not just that something did (transcripts die at compaction; the
// verdict record's hook-captured `findings` lines are the durable copy).
func priorDesignFindings(c *runctl.Ctl, r *store.Run, b *strings.Builder) {
	env, err := c.Env(r)
	if err != nil {
		return
	}
	latest := map[string][]any{}
	for _, v := range env.Records("verdict") {
		agent, _ := v.Data["agent"].(string)
		switch agent {
		case "design-reviewer", "critic", "adversary":
		default:
			continue
		}
		switch s, _ := v.Data["status"].(string); s {
		case "changes-required", "unsafe", "risky":
			fl, _ := v.Data["findings"].([]any)
			latest[agent] = fl // stream order: last failing verdict wins
		default:
			delete(latest, agent) // a later pass supersedes the failure
		}
	}
	printed := false
	for _, agent := range []string{"design-reviewer", "critic", "adversary"} {
		fl := latest[agent]
		if len(fl) == 0 {
			continue
		}
		if !printed {
			b.WriteString("prior failing design verdicts — findings the rework must answer:\n")
			printed = true
		}
		fmt.Fprintf(b, "  %s:\n", agent)
		for i, f := range fl {
			if i == 5 {
				fmt.Fprintf(b, "    … and %d more (see the verdict record)\n", len(fl)-5)
				break
			}
			fmt.Fprintf(b, "    - %v\n", f)
		}
	}
}

// implementerBriefing scopes an implementer spawn to the single active
// task: tid, DoD, its ACs (text + the verification command whose exact
// invocation the capture hook recognizes), the approved design refs the
// conformance reviewer will hold the diff to, and the recorded out-of-scope
// boundary.
func implementerBriefing(c *runctl.Ctl, r *store.Run, b *strings.Builder) {
	env, err := c.Env(r)
	if err != nil {
		return
	}
	task, ok := activeTaskRecord(env)
	if !ok {
		b.WriteString("no single active task — the main thread marks exactly one task in_progress (wf record task updates=<id> status=in_progress) before spawning; return and say so.\n")
		return
	}
	tid, _ := task.Data["tid"].(string)
	subject, _ := task.Data["subject"].(string)
	fmt.Fprintf(b, "assigned task for this spawn: %s — %s\n", tid, subject)
	if dod, _ := task.Data["dod"].([]any); len(dod) > 0 {
		b.WriteString("definition of done:\n")
		for _, d := range dod {
			fmt.Fprintf(b, "  - %v\n", d)
		}
	}
	if acs := strList(task.Data["ac_links"]); len(acs) > 0 {
		acText, strat := acIndex(env)
		b.WriteString("acceptance criteria this task carries (test-first: the red test encodes the AC):\n")
		for _, ac := range acs {
			line := "  - " + ac
			if t := acText[ac]; t != "" {
				line += ": " + t
			}
			b.WriteString(line + "\n")
			if s := strat[ac]; s != "" {
				fmt.Fprintf(b, "    verification: %s — run EXACTLY this invocation (it is what the capture hook recognizes)\n", s)
			}
		}
	}
	var sels []string
	for _, os := range env.Records("option-set") {
		if s, _ := os.Data["selected"].(string); s != "" {
			st, _ := os.Data["stage"].(string)
			sels = append(sels, st+":"+s)
		}
	}
	if len(sels) > 0 {
		fmt.Fprintf(b, "approved design selections (conformance-reviewed; departures need wf record deviation, never improvised): %s\n", strings.Join(sels, ", "))
	}
	for _, a := range env.Records("artifact") {
		if t, _ := a.Data["template"].(string); t == "adr" {
			if p, _ := a.Data["path"].(string); p != "" {
				fmt.Fprintf(b, "ADR: %s\n", p)
			}
		}
	}
	for _, sb := range env.Records("scope-boundary") {
		if oos := strList(sb.Data["out_of_scope"]); len(oos) > 0 {
			fmt.Fprintf(b, "out of scope (do not touch — discoveries become wf record followup): %s\n", strings.Join(oos, "; "))
		}
	}
}

// activeTaskRecord mirrors the task gates' selection: the single
// in_progress task, falling back to a single open one.
func activeTaskRecord(env *contracts.Env) (contracts.Record, bool) {
	for _, statuses := range [][]string{{"in_progress"}, {"in_progress", "open"}} {
		var found contracts.Record
		count := 0
		for _, tr := range env.Records("task") {
			s, _ := tr.Data["status"].(string)
			for _, want := range statuses {
				if s == want {
					count++
					found = tr
					break
				}
			}
		}
		if count == 1 {
			return found, true
		}
	}
	return contracts.Record{}, false
}

// acIndex maps AC id → requirement AC text and AC id → verification command
// (falling back to the method when no command is recorded).
func acIndex(env *contracts.Env) (map[string]string, map[string]string) {
	text := map[string]string{}
	for _, req := range env.Records("requirement") {
		if acs, _ := req.Data["acs"].([]any); acs != nil {
			for _, raw := range acs {
				if m, ok := raw.(map[string]any); ok {
					id, _ := m["id"].(string)
					t, _ := m["text"].(string)
					if id != "" {
						text[id] = t
					}
				}
			}
		}
	}
	strat := map[string]string{}
	for _, vs := range env.Records("verification-strategy") {
		ac, _ := vs.Data["ac"].(string)
		if ac == "" {
			continue
		}
		if cmd, _ := vs.Data["command"].(string); cmd != "" {
			strat[ac] = cmd
		} else if m, _ := vs.Data["method"].(string); m != "" {
			strat[ac] = m
		}
	}
	return text, strat
}

func strList(v any) []string {
	raw, _ := v.([]any)
	var out []string
	for _, el := range raw {
		if s, ok := el.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

func pluginRoot(c *runctl.Ctl) string {
	if r := os.Getenv("CLAUDE_PLUGIN_ROOT"); r != "" {
		return r
	}
	return c.Spec.PluginRoot()
}

// ---------------------------------------------------------------------------

func split(fs []contracts.Finding) (agent, user []contracts.Finding) {
	for _, f := range fs {
		if f.UserBlocked {
			user = append(user, f)
		} else {
			agent = append(agent, f)
		}
	}
	return
}

func taskCounts(env *contracts.Env) (open, total int) {
	for _, tr := range env.Records("task") {
		total++
		switch fmt.Sprintf("%v", tr.Data["status"]) {
		case "open", "in_progress":
			open++
		}
	}
	return
}

func phaseIndex(c *runctl.Ctl, family, phase string) string {
	ph := c.Spec.PhasesFor(family)
	for i, p := range ph {
		if p.ID == phase {
			return fmt.Sprintf("%d/%d", i+1, len(ph))
		}
	}
	return "-"
}

func skillFor(c *runctl.Ctl, phase string) string {
	if p, ok := c.Spec.PhaseByID(phase); ok {
		return p.Skill
	}
	return "dev"
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func ago(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return "?"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
