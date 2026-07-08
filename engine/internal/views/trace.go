// Package views renders derived, never-source views over run state (08 §4):
// the Ship trace report with its idempotent trace-finding records.
package views

import (
	"fmt"
	"path"
	"path/filepath"
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
	// approval drift — records added AFTER the gate's last approval bound
	// its refs: the user approved a baseline that has since grown without
	// re-approval (re-approving the gate binds the new refs; the finding
	// then wants a disposition naming that re-approval)
	for _, gate := range []string{"scope", "design", "plan"} {
		approved := latestApprovedRefs(env, gate)
		if approved == nil {
			continue // gate never approved with bound refs — nothing to drift from
		}
		for _, ref := range contracts.ApprovalRefs(env, gate) {
			if !approved[ref] {
				wants = append(wants, want{
					key:      "drift:" + gate + ":" + ref,
					text:     fmt.Sprintf("approval drift: %s appeared after the last %s approval — re-approve the %s gate or disposition why it needs none", ref, gate, gate),
					severity: "medium",
				})
			}
		}
	}

	// out-of-scope edits — the recorded scope boundary, mechanically
	// checked: path-like out_of_scope entries (globs, /-paths, filenames)
	// are matched against the auto-captured edit paths; prose entries stay
	// the auditor's to judge. The boundary used to be write-only.
	seenEdits := map[string]bool{}
	for _, sb := range env.Records("scope-boundary") {
		oos, _ := sb.Data["out_of_scope"].([]any)
		for _, raw := range oos {
			pat, ok := raw.(string)
			if !ok || !pathLike(pat) {
				continue
			}
			for _, ed := range env.Records("edit") {
				p, _ := ed.Data["path"].(string)
				if p == "" || seenEdits[ed.ID] || !pathMatches(pat, p) {
					continue
				}
				seenEdits[ed.ID] = true
				wants = append(wants, want{
					key:      "scope:" + ed.ID,
					text:     fmt.Sprintf("edit touched declared out-of-scope path %s (boundary entry %q)", p, pat),
					severity: "high",
				})
			}
		}
	}
	// gating-verdict staleness — edits recorded AFTER the latest gating
	// verdict mean the reviews ran on a diff that has since changed.
	// Informational, not a gate: the auditor judges whether the late change
	// needed re-review (verdict-in deliberately carries verdicts forward).
	if lastV := lastGatingVerdictOrder(env); lastV >= 0 {
		late := 0
		for _, ed := range env.Records("edit") {
			if ed.Order > lastV {
				late++
			}
		}
		if late > 0 {
			wants = append(wants, want{
				key:      "stale-verdicts",
				text:     fmt.Sprintf("%d edit(s) recorded after the latest gating verdict — the reviewed diff changed since review; re-run the affected reviewers or disposition why no re-review was needed", late),
				severity: "medium",
			})
		}
	}
	// a forced Frame exit without a risk record silently vacates every
	// when.signals-conditioned security item (threat model, attack tree,
	// attack paths): false-by-absence must be visible, never silent.
	if len(env.Records("risk")) == 0 {
		for _, ev := range env.Events {
			if ev.Kind == "phase" && ev.Str("action") == "force" && ev.Phase == "frame" {
				wants = append(wants, want{
					key:      "vacated-signals:" + ev.ID,
					text:     "frame was forced without a risk scan — signal-conditioned security items (threat model, attack tree, attack paths) never applied; wf risk scan and re-check, or disposition why the change is signal-free",
					severity: "high",
				})
				break
			}
		}
	}
	// invalidated high-risk assumptions — the work was built on something
	// that turned out false; the close-out needs an explicit disposition
	// naming why the result still stands.
	for _, a := range env.Records("assumption") {
		hr, _ := a.Data["high_risk"].(bool)
		if s, _ := a.Data["status"].(string); hr && s == "invalidated" {
			wants = append(wants, want{
				key:      "assumption-invalidated:" + a.ID,
				text:     fmt.Sprintf("high-risk assumption invalidated: %v — disposition why the delivered result still stands", a.Data["text"]),
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

// latestApprovedRefs returns the ref set bound by the gate's newest
// approval, or nil when the gate was never approved with bound refs
// (approvals predating the refs feature carry none — no baseline, no drift).
func latestApprovedRefs(env *contracts.Env, gate string) map[string]bool {
	var refs map[string]bool
	for _, a := range env.Records("approval") {
		if g, _ := a.Data["gate"].(string); g != gate {
			continue
		}
		raw, ok := a.Data["approved_refs"].([]any)
		if !ok {
			continue
		}
		refs = map[string]bool{}
		for _, r := range raw {
			if s, ok := r.(string); ok {
				refs[s] = true
			}
		}
	}
	return refs
}

func has(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// pathLike decides whether an out_of_scope entry is mechanically checkable:
// a glob, a /-separated path, or a filename with an extension. Prose entries
// ("authentication") are skipped — they are the auditor's to judge.
func pathLike(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || strings.ContainsAny(s, " \t") {
		return false // prose
	}
	return strings.ContainsAny(s, "/*?[") || strings.Contains(s, ".")
}

// pathMatches matches a boundary entry against a normalized project-relative
// edit path: trailing-/ entries and bare directories match by prefix, glob
// entries via path.Match against the full path and the basename.
func pathMatches(pat, p string) bool {
	pat = strings.TrimPrefix(strings.TrimSpace(filepath.ToSlash(pat)), "./")
	p = strings.TrimPrefix(filepath.ToSlash(p), "./")
	if pat == "" {
		return false
	}
	if strings.HasSuffix(pat, "/") {
		return strings.HasPrefix(p, pat)
	}
	if strings.ContainsAny(pat, "*?[") {
		if ok, _ := path.Match(pat, p); ok {
			return true
		}
		ok, _ := path.Match(pat, path.Base(p))
		return ok
	}
	return p == pat || strings.HasPrefix(p, pat+"/")
}

// lastGatingVerdictOrder returns the stream position of the newest verdict
// by a gating roster agent, or -1 when none exists.
func lastGatingVerdictOrder(env *contracts.Env) int {
	last := -1
	for _, v := range env.Records("verdict") {
		agent, _ := v.Data["agent"].(string)
		ag, ok := env.Spec.AgentByName(agent)
		if !ok || !ag.Gating {
			continue
		}
		if v.Order > last {
			last = v.Order
		}
	}
	return last
}
