package views

// wf report — the health-signals view (08 §4): the honest tiers kept from
// v0.36. Loops per run, escapes, self-attested counts, ungrounded ACs,
// lesson efficacy, deliver-reached. Derived, never source; log replay is
// allowed here (07 §5) — this feeds humans and the run archive, never gates.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

// RunSignals is one run's health snapshot. It is computed purely from the
// run's events, so it works identically for the active run and archives.
type RunSignals struct {
	Run    string `json:"run"`
	Family string `json:"family"`
	Intent string `json:"intent"`
	Status string `json:"status"`          // active | parked | closed
	Phase  string `json:"phase,omitempty"` // open runs only

	// escapes (04 §7: recorded, reported, escalating)
	Loops           int            `json:"loops"`
	Forces          int            `json:"forces"`
	Parks           int            `json:"parks"`
	EscapesByAction map[string]int `json:"escapes_by_action,omitempty"`
	Waivers         int            `json:"waivers"` // contract-item waivers
	WaivedPhases    []string       `json:"waived_phases,omitempty"`

	// self-attestation (04 §8: manual records are auto:false and reported)
	Verdicts     int `json:"verdicts"`
	AutoVerdicts int `json:"auto_verdicts"`
	TestRuns     int `json:"test_runs"`
	AutoTestRuns int `json:"auto_test_runs"`

	// grounding
	UngroundedTestRuns int      `json:"ungrounded_test_runs"`
	ACPasses           int      `json:"ac_passes"`
	UngroundedACs      []string `json:"ungrounded_acs,omitempty"` // pass without a grounded green tagged to the AC

	// lesson efficacy (03 §4.7)
	LessonsProposed   int `json:"lessons_proposed"`
	LessonsAccepted   int `json:"lessons_accepted"`
	LessonsRejected   int `json:"lessons_rejected"`
	LessonItemWaivers int `json:"lesson_item_waivers"` // waived lesson.* contract items — accepted but dodged

	DeliverReached bool `json:"deliver_reached"`
}

// Report computes signals for every archived run plus the active one —
// "loops per run", "lesson efficacy" and "deliver-reached" only mean
// anything across runs.
func Report(c *runctl.Ctl) ([]RunSignals, error) {
	ids, err := c.Store.ListArchivedRuns()
	if err != nil {
		return nil, err
	}
	var out []RunSignals
	for _, id := range ids {
		s, err := runSignals(c, id, true)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	if r, err := c.Store.LoadRun(); err == nil && r != nil {
		s, err := runSignals(c, r.ID, false)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return out, nil
}

// ReportRun computes one run's signals. runID "" or "current" = the active
// run; otherwise an archive ID.
func ReportRun(c *runctl.Ctl, runID string) (*RunSignals, error) {
	if runID == "" || runID == "current" {
		r, err := c.MustRun()
		if err != nil {
			return nil, err
		}
		return runSignals(c, r.ID, false)
	}
	// prefer the archive; fall back to the live log (a just-closed or
	// still-open run addressed by ID)
	if _, err := c.Store.ArchivedRunEvents(runID); err == nil {
		return runSignals(c, runID, true)
	}
	return runSignals(c, runID, false)
}

// WriteRunSignals freezes the snapshot into the run archive
// (runs/<id>/signals.json) — the Ship close-out step of 03 §4.7.
func WriteRunSignals(c *runctl.Ctl, runID string) (string, error) {
	s, err := runSignals(c, runID, true)
	if err != nil {
		return "", err
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(c.Store.RunsDir(), runID, "signals.json")
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// runSignals folds one run's events into signals. Lesson events stay in the
// live log across close (keepLive, 08 §6), so archived runs merge their
// archive slice with any live events still carrying their run ID.
func runSignals(c *runctl.Ctl, runID string, archived bool) (*RunSignals, error) {
	var evs []store.Event
	if archived {
		arch, err := c.Store.ArchivedRunEvents(runID)
		if err != nil {
			return nil, err
		}
		live, err := c.Store.RunEvents(runID)
		if err != nil {
			return nil, err
		}
		evs = append(arch, live...)
	} else {
		var err error
		evs, err = c.Store.RunEvents(runID)
		if err != nil {
			return nil, err
		}
		if len(evs) == 0 {
			return nil, fmt.Errorf("no events for run %s", runID)
		}
	}

	s := &RunSignals{Run: runID, Status: "active", EscapesByAction: map[string]int{}}

	// engine transitions (raw events — run/phase kinds are not records)
	for _, e := range evs {
		switch e.Kind {
		case "run":
			switch e.Str("action") {
			case "start", "branch", "adopt":
				s.Family, s.Intent = e.Str("family"), e.Str("intent")
			case "close":
				s.Status = "closed"
			}
		case "phase":
			switch e.Str("action") {
			case "enter":
				s.Phase = e.Str("target")
			case "waive":
				s.WaivedPhases = append(s.WaivedPhases, e.Str("target"))
			case "loop":
				s.Loops++
			case "force":
				s.Forces++
			case "park":
				if s.Status == "active" {
					s.Status = "parked"
				}
			case "resume":
				s.Status = "active"
			case "exit":
				if e.Phase == "ship" {
					s.DeliverReached = true
				}
			}
		case "escape":
			s.EscapesByAction[e.Str("action")]++
			switch e.Str("action") {
			case "park":
				s.Parks++
			}
		}
	}
	if s.Status == "closed" {
		s.Phase = ""
	}

	// records (updates-folded — lesson status flips must resolve)
	env := &contracts.Env{Events: evs}
	groundedGreens := map[string]bool{}
	for _, tr := range env.Records("test-run") {
		s.TestRuns++
		if tr.Auto {
			s.AutoTestRuns++
		}
		grounded, _ := tr.Data["grounded"].(bool)
		if !grounded {
			s.UngroundedTestRuns++
			continue
		}
		if exit, ok := tr.Data["exit"].(float64); ok && exit == 0 {
			if ac, _ := tr.Data["ac"].(string); ac != "" {
				groundedGreens[ac] = true
			}
		}
	}
	for _, v := range env.Records("verdict") {
		s.Verdicts++
		if v.Auto {
			s.AutoVerdicts++
		}
	}
	for _, av := range env.Records("ac-verdict") {
		if status, _ := av.Data["status"].(string); status != "pass" {
			continue
		}
		s.ACPasses++
		if ac, _ := av.Data["ac"].(string); ac != "" && !groundedGreens[ac] {
			s.UngroundedACs = append(s.UngroundedACs, ac)
		}
	}
	for _, w := range env.Records("waiver") {
		s.Waivers++
		if item, _ := w.Data["item"].(string); strings.HasPrefix(item, "lesson.") {
			s.LessonItemWaivers++
		}
	}
	for _, l := range env.Records("lesson") {
		switch status, _ := l.Data["status"].(string); status {
		case "proposed":
			s.LessonsProposed++
		case "accepted":
			s.LessonsAccepted++
		case "rejected":
			s.LessonsRejected++
		}
	}
	if len(s.EscapesByAction) == 0 {
		s.EscapesByAction = nil
	}
	return s, nil
}

// RenderReport renders the aggregate text view.
func RenderReport(sigs []RunSignals) string {
	if len(sigs) == 0 {
		return "[wf report] no runs recorded yet\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[wf report] %d run(s)\n", len(sigs))
	b.WriteString("run · family/intent · status · loops · forces · parks · waivers · verdicts(auto) · tests(auto) · delivered\n")
	var loops, forces, parks, delivered int
	for _, s := range sigs {
		mark := "✗"
		if s.DeliverReached {
			mark = "✓"
			delivered++
		}
		loops += s.Loops
		forces += s.Forces
		parks += s.Parks
		fmt.Fprintf(&b, "  %s · %s/%s · %s · %d · %d · %d · %d · %d(%d) · %d(%d) · %s\n",
			s.Run, s.Family, s.Intent, s.Status, s.Loops, s.Forces, s.Parks, s.Waivers,
			s.Verdicts, s.AutoVerdicts, s.TestRuns, s.AutoTestRuns, mark)
	}
	fmt.Fprintf(&b, "totals: loops %d · forces %d · parks %d · delivered %d/%d\n", loops, forces, parks, delivered, len(sigs))
	for _, s := range sigs {
		b.WriteString(renderConcerns(s))
	}
	return b.String()
}

// RenderRunSignals renders one run's full snapshot.
func RenderRunSignals(s *RunSignals) string {
	var b strings.Builder
	fmt.Fprintf(&b, "[wf report] run %s (%s/%s) — %s\n", s.Run, s.Family, s.Intent, s.Status)
	fmt.Fprintf(&b, "  loops %d · forces %d · parks %d · waivers %d (phases: %s)\n",
		s.Loops, s.Forces, s.Parks, s.Waivers, orDash(s.WaivedPhases))
	fmt.Fprintf(&b, "  verdicts %d (%d auto) · test-runs %d (%d auto, %d ungrounded)\n",
		s.Verdicts, s.AutoVerdicts, s.TestRuns, s.AutoTestRuns, s.UngroundedTestRuns)
	fmt.Fprintf(&b, "  AC passes %d (ungrounded: %s)\n", s.ACPasses, orDash(s.UngroundedACs))
	fmt.Fprintf(&b, "  lessons: %d proposed · %d accepted · %d rejected · %d lesson-item waivers\n",
		s.LessonsProposed, s.LessonsAccepted, s.LessonsRejected, s.LessonItemWaivers)
	fmt.Fprintf(&b, "  deliver-reached: %v\n", s.DeliverReached)
	b.WriteString(renderConcerns(*s))
	return b.String()
}

// renderConcerns flags the dishonesty signatures worth a human look.
func renderConcerns(s RunSignals) string {
	var c []string
	if len(s.UngroundedACs) > 0 {
		c = append(c, fmt.Sprintf("AC passes without grounded greens: %s", strings.Join(s.UngroundedACs, ", ")))
	}
	if s.Verdicts >= 2 && s.AutoVerdicts == 0 {
		c = append(c, "all verdicts self-attested (SubagentStop capture dead?)")
	}
	if s.TestRuns >= 3 && s.AutoTestRuns == 0 {
		c = append(c, "all test-runs self-attested (Bash capture dead or runner unrecognized?)")
	}
	if s.LessonItemWaivers > 0 {
		c = append(c, fmt.Sprintf("%d accepted lesson item(s) waived — lessons dodged", s.LessonItemWaivers))
	}
	if len(c) == 0 {
		return ""
	}
	return fmt.Sprintf("  ⚠ %s: %s\n", s.Run, strings.Join(c, " · "))
}

func orDash(xs []string) string {
	if len(xs) == 0 {
		return "—"
	}
	return strings.Join(xs, ", ")
}
