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
	} else if stage := designStage(c, r, agentName); stage != "" {
		fmt.Fprintf(&b, "assigned design stage for this spawn: %s — name it in the `stage` field of the returned option-set.\n", stage)
	}
	if len(ag.Corpus) > 0 {
		root := pluginRoot(c)
		verb := "authoring"
		if ag.Gating {
			verb = "judging"
		}
		fmt.Fprintf(&b, "reference corpus for this %s (read before %s; cite file+rule):\n", role, verb)
		for _, p := range ag.Corpus {
			fmt.Fprintf(&b, "  - %s\n", filepath.Join(root, filepath.FromSlash(p)))
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
			fmt.Fprintf(b, "  - %s\n", d)
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
// addressing channel). ux-designer always works the ux stage.
func designStage(c *runctl.Ctl, r *store.Run, agent string) string {
	switch agent {
	case "ux-designer":
		return "ux"
	case "designer":
	default:
		return ""
	}
	env, err := c.Env(r)
	if err != nil {
		return ""
	}
	stages := map[string]bool{}
	for _, rec := range env.Records("option-set") {
		if s, _ := rec.Data["stage"].(string); s != "" {
			stages[s] = true
		}
	}
	switch {
	case !stages["system"]:
		return "system"
	case !stages["software"]:
		return "software"
	default:
		return "system or software (loop re-entry — state which; carry the recorded rejected option IDs forward, a rejected option may not reappear)"
	}
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
