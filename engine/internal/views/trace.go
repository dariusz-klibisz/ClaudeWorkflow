// Package views renders derived, never-source views over run state (08 §4):
// the Ship trace report with its idempotent trace-finding records.
package views

import (
	"fmt"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

// Trace computes the Ship close-out findings from state (03 §4.7), writes
// missing trace-finding records (idempotent via a `key` field), and returns
// the human report. The agent resolves findings via
// `wf record trace-finding updates=<id> status=resolved|dispositioned`.
func Trace(c *runctl.Ctl) (string, error) {
	r, err := c.MustRun()
	if err != nil {
		return "", err
	}
	env, err := c.Env(r)
	if err != nil {
		return "", err
	}

	type want struct {
		key      string
		text     string
		severity string
	}
	var wants []want

	// forced phase exits — always high: the gate was bypassed
	for _, ev := range env.Events {
		if ev.Kind == "phase" && ev.Str("action") == "force" {
			wants = append(wants, want{
				key:      "force:" + ev.ID,
				text:     fmt.Sprintf("phase %s was force-exited: %s", ev.Phase, ev.Str("reason")),
				severity: "high",
			})
		}
	}
	// open followups — must become tasks now or be carried
	for _, f := range env.Records("followup") {
		if s, _ := f.Data["status"].(string); s == "open" {
			wants = append(wants, want{
				key:      "followup:" + f.ID,
				text:     fmt.Sprintf("open followup: %v", f.Data["text"]),
				severity: "high",
			})
		}
	}
	// pending deviations (belt — build gate should have caught them)
	for _, d := range env.Records("deviation") {
		if s, _ := d.Data["status"].(string); s == "pending" {
			wants = append(wants, want{
				key:      "deviation:" + d.ID,
				text:     fmt.Sprintf("unacked deviation: %v", d.Data["text"]),
				severity: "high",
			})
		}
	}
	// ambiguities still deferred/open at ship
	for _, a := range env.Records("ambiguity") {
		if s, _ := a.Data["disposition"].(string); s == "deferred" || s == "open" {
			wants = append(wants, want{
				key:      "ambiguity:" + a.ID,
				text:     fmt.Sprintf("ambiguity never resolved (%v lens): %v", a.Data["lens"], a.Data["text"]),
				severity: "medium",
			})
		}
	}

	// idempotent write: only keys not yet recorded
	existing := map[string]bool{}
	for _, tf := range env.Records("trace-finding") {
		if k, _ := tf.Data["key"].(string); k != "" {
			existing[k] = true
		}
	}
	created := 0
	for _, w := range wants {
		if existing[w.key] {
			continue
		}
		if _, err := c.Record("trace-finding", map[string]any{
			"key": w.key, "text": w.text, "severity": w.severity, "status": "open",
		}, true, "engine"); err != nil {
			return "", err
		}
		created++
	}

	// re-read for the report
	env, err = c.Env(r)
	if err != nil {
		return "", err
	}
	return renderTrace(c, r, env, created), nil
}

func renderTrace(c *runctl.Ctl, r *store.Run, env *contracts.Env, created int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "[wf trace] run %s (%s/%s)\n", r.ID, r.Family, r.Intent)

	// phase coverage
	b.WriteString("phase coverage:\n")
	for _, p := range c.Spec.PhasesFor(r.Family) {
		state := "pending"
		switch {
		case has(r.ExitedPh, p.ID):
			state = "exited"
		case has(r.WaivedPh, p.ID):
			state = "WAIVED"
		case p.ID == r.Phase:
			state = "current"
		}
		fmt.Fprintf(&b, "  %-8s %s\n", p.ID, state)
	}

	// escapes and waivers (informational — the auditor judges them)
	var waivers, escapes []string
	for _, w := range env.Records("waiver") {
		waivers = append(waivers, fmt.Sprintf("%v (%v)", w.Data["item"], w.Data["reason"]))
	}
	for _, e := range env.Records("escape") {
		escapes = append(escapes, fmt.Sprintf("%v: %v", e.Data["action"], e.Data["reason"]))
	}
	if len(waivers) > 0 {
		fmt.Fprintf(&b, "waivers (%d): %s\n", len(waivers), strings.Join(waivers, " · "))
	}
	if len(escapes) > 0 {
		fmt.Fprintf(&b, "escapes (%d): %s\n", len(escapes), strings.Join(escapes, " · "))
	}
	fmt.Fprintf(&b, "loops: %d · forces: %d\n", r.Loops, r.Forces)

	// findings
	open := 0
	var lines []string
	for _, tf := range env.Records("trace-finding") {
		status, _ := tf.Data["status"].(string)
		mark := "·"
		if status == "open" {
			open++
			mark = "✗"
		}
		lines = append(lines, fmt.Sprintf("  %s [%v] %v (%s) — id %s", mark, tf.Data["severity"], tf.Data["text"], status, tf.ID))
	}
	if len(lines) == 0 {
		b.WriteString("findings: none — clean close-out\n")
	} else {
		fmt.Fprintf(&b, "findings (%d new, %d open):\n%s\n", created, open, strings.Join(lines, "\n"))
		if open > 0 {
			b.WriteString("resolve each: wf record trace-finding updates=<id> status=resolved|dispositioned note=\"…\"\n")
		}
	}
	return b.String()
}

func has(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
