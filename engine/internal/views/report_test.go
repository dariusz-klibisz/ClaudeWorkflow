package views

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

// raw appends a synthetic event — report is a view over events, so tests
// build ledgers directly (write-time validation is covered elsewhere).
func raw(t *testing.T, c *runctl.Ctl, run, phase, kind string, auto bool, actor string, data map[string]any) {
	t.Helper()
	if err := c.Store.Append(&store.Event{Run: run, Phase: phase, Kind: kind, Auto: auto, Actor: actor, Data: data}); err != nil {
		t.Fatal(err)
	}
}

func TestReportRunSignals(t *testing.T) {
	c := newCtl(t)
	r, err := c.RunStart("diff", "fix")
	if err != nil {
		t.Fatal(err)
	}
	id := r.ID

	// escapes + transitions
	raw(t, c, id, "build", "phase", false, "engine", map[string]any{"action": "loop", "target": "build"})
	raw(t, c, id, "build", "phase", false, "engine", map[string]any{"action": "force"})
	raw(t, c, id, "build", "escape", false, "user", map[string]any{"action": "force", "reason": "r", "level": 1})
	raw(t, c, id, "build", "escape", false, "user", map[string]any{"action": "park", "reason": "blocked"})
	raw(t, c, id, "design", "phase", false, "engine", map[string]any{"action": "waive", "target": "design"})
	// self-attestation mix
	raw(t, c, id, "build", "verdict", true, "hook", map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0})
	raw(t, c, id, "build", "verdict", false, "agent", map[string]any{"agent": "adversary", "status": "clean", "criticals": 0, "majors": 0})
	// grounding: AC-1 has a grounded green; AC-2 passes without one
	raw(t, c, id, "build", "test-run", true, "hook", map[string]any{"cmd": "go test", "exit": 0.0, "grounded": true, "ac": "AC-1"})
	raw(t, c, id, "build", "test-run", false, "agent", map[string]any{"cmd": "go test | tail", "exit": 0.0, "grounded": false})
	raw(t, c, id, "verify", "ac-verdict", false, "agent", map[string]any{"ac": "AC-1", "status": "pass"})
	raw(t, c, id, "verify", "ac-verdict", false, "agent", map[string]any{"ac": "AC-2", "status": "pass"})
	// waivers incl. a dodged lesson item; lesson lifecycle with an update
	raw(t, c, id, "build", "waiver", false, "user", map[string]any{"item": "T-3", "reason": "testless"})
	raw(t, c, id, "build", "waiver", false, "user", map[string]any{"item": "lesson.check-hooks", "reason": "n/a"})
	lev := &store.Event{Run: id, Phase: "ship", Kind: "lesson", Actor: "agent", Data: map[string]any{"text": "l1", "status": "proposed"}}
	if err := c.Store.Append(lev); err != nil {
		t.Fatal(err)
	}
	raw(t, c, id, "ship", "lesson", false, "user", map[string]any{"updates": lev.ID, "status": "accepted"})

	s, err := ReportRun(c, "current")
	if err != nil {
		t.Fatal(err)
	}
	check := func(name string, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}
	check("loops", s.Loops, 1)
	check("forces", s.Forces, 1)
	check("parks", s.Parks, 1)
	check("verdicts", s.Verdicts, 2)
	check("autoVerdicts", s.AutoVerdicts, 1)
	check("testRuns", s.TestRuns, 2)
	check("autoTestRuns", s.AutoTestRuns, 1)
	check("ungroundedTestRuns", s.UngroundedTestRuns, 1)
	check("acPasses", s.ACPasses, 2)
	check("waivers", s.Waivers, 2)
	check("lessonItemWaivers", s.LessonItemWaivers, 1)
	// the update folded: accepted, not proposed
	check("lessonsProposed", s.LessonsProposed, 0)
	check("lessonsAccepted", s.LessonsAccepted, 1)
	if len(s.UngroundedACs) != 1 || s.UngroundedACs[0] != "AC-2" {
		t.Errorf("ungrounded ACs = %v, want [AC-2]", s.UngroundedACs)
	}
	if !strings.Contains(s.WaivedPhases[0], "design") {
		t.Errorf("waived phases = %v", s.WaivedPhases)
	}
	if s.DeliverReached {
		t.Error("run not shipped — deliver-reached must be false")
	}
	// concerns rendered
	out := RenderRunSignals(s)
	if !strings.Contains(out, "AC passes without grounded greens: AC-2") || !strings.Contains(out, "lessons dodged") {
		t.Errorf("concern lines missing:\n%s", out)
	}
}

func TestReportArchivedAndAggregate(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "new")
	id := r.ID
	raw(t, c, id, "ship", "phase", false, "engine", map[string]any{"action": "exit"})
	// a lesson (keepLive: survives archival in the live log)
	raw(t, c, id, "ship", "lesson", false, "agent", map[string]any{"text": "l", "status": "proposed"})
	raw(t, c, id, "", "run", false, "engine", map[string]any{"action": "close"})
	if err := c.Store.ArchiveRun(id); err != nil {
		t.Fatal(err)
	}

	// snapshot freeze (the run close CLI step)
	path, err := WriteRunSignals(c, id)
	if err != nil {
		t.Fatal(err)
	}
	var snap RunSignals
	rawJSON, _ := os.ReadFile(path)
	if err := json.Unmarshal(rawJSON, &snap); err != nil {
		t.Fatal(err)
	}
	if snap.Run != id || snap.Status != "closed" || !snap.DeliverReached {
		t.Fatalf("snapshot wrong: %+v", snap)
	}
	if snap.LessonsProposed != 1 {
		t.Fatalf("keep-live lesson must be merged into archived signals: %+v", snap)
	}
	if filepath.Base(path) != "signals.json" {
		t.Fatalf("snapshot path: %s", path)
	}

	// second, active run — aggregate sees both
	r2, err := c.RunStart("assessment", "investigate")
	if err != nil {
		t.Fatal(err)
	}
	sigs, err := Report(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 2 || sigs[0].Run != id || sigs[1].Run != r2.ID {
		t.Fatalf("aggregate = %+v", sigs)
	}
	if !sigs[0].DeliverReached || sigs[1].DeliverReached {
		t.Fatalf("deliver flags wrong: %+v", sigs)
	}
	out := RenderReport(sigs)
	if !strings.Contains(out, "2 run(s)") || !strings.Contains(out, "delivered 1/2") {
		t.Errorf("aggregate render:\n%s", out)
	}

	// --run against the archive ID
	s, err := ReportRun(c, id)
	if err != nil {
		t.Fatal(err)
	}
	if s.Status != "closed" {
		t.Fatalf("archived run status: %s", s.Status)
	}
}
