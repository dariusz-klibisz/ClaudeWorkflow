package ops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func specPath(t *testing.T) string {
	t.Helper()
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	return p
}

func TestLessonsSuggestFromSignals(t *testing.T) {
	c, _ := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	// force + all-manual verdicts → two suggestions
	_ = c.Store.Append(&store.Event{Run: r.ID, Phase: "build", Kind: "phase", Actor: "engine", Data: map[string]any{"action": "force"}})
	_ = c.Store.Append(&store.Event{Run: r.ID, Phase: "build", Kind: "verdict", Actor: "agent", Data: map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0}})
	_ = c.Store.Append(&store.Event{Run: r.ID, Phase: "build", Kind: "verdict", Actor: "agent", Data: map[string]any{"agent": "adversary", "status": "clean", "criticals": 0, "majors": 0}})

	out, err := LessonsSuggest(c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "proposed 2 lesson(s)") {
		t.Fatalf("want 2 proposals:\n%s", out)
	}
	// idempotent: same signals, nothing new
	out, err = LessonsSuggest(c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "nothing to propose") {
		t.Fatalf("suggest must be idempotent:\n%s", out)
	}
}

// Artifact runs record manual doc-check test-runs by design — the
// runner-recognition suggestion must not fire there (the arch-design run's
// false positive).
func TestLessonsSuggestQuietForArtifactFamily(t *testing.T) {
	c, _ := newCtl(t)
	r, _ := c.RunStart("artifact", "arch-design")
	_ = c.Store.Append(&store.Event{Run: r.ID, Phase: "verify", Kind: "test-run", Actor: "agent",
		Data: map[string]any{"cmd": "grep -q x docs/d.md", "exit": 0, "grounded": true}})
	_ = c.Store.Append(&store.Event{Run: r.ID, Phase: "verify", Kind: "test-run", Actor: "agent",
		Data: map[string]any{"cmd": "grep -q y docs/d.md", "exit": 0, "grounded": true}})
	_ = c.Store.Append(&store.Event{Run: r.ID, Phase: "verify", Kind: "test-run", Actor: "agent",
		Data: map[string]any{"cmd": "grep -q z docs/d.md", "exit": 0, "grounded": true}})

	out, err := LessonsSuggest(c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "nothing to propose") {
		t.Fatalf("artifact family must not trigger the runner suggestion:\n%s", out)
	}
}

func TestLessonsAcceptProseAndReject(t *testing.T) {
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	ev, err := c.Record("lesson", map[string]any{"text": "run wf doctor at session start", "status": "proposed"}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}

	out, err := LessonsAccept(c, dir, specPath(t), ev.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "accepted") {
		t.Fatalf("accept output: %s", out)
	}
	rulePath := filepath.Join(dir, ".claude", "rules", "wf-lessons.md")
	raw, err := os.ReadFile(rulePath)
	if err != nil {
		t.Fatalf("prose rules file missing: %v", err)
	}
	if !strings.Contains(string(raw), "run wf doctor at session start") || !strings.Contains(string(raw), "<!-- wf:lessons:begin -->") {
		t.Fatalf("rules content wrong:\n%s", raw)
	}
	// approval + status recorded
	r, _ := c.Store.LoadRun()
	env, _ := c.Env(r)
	if got := len(env.Records("approval")); got != 1 {
		t.Fatalf("approval events = %d, want 1", got)
	}
	// reject flips it back out; engine-owned file goes away
	if _, err := LessonsReject(c, dir, specPath(t), ev.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(rulePath); !os.IsNotExist(err) {
		t.Fatal("rules file must be removed when no prose lessons remain")
	}
}

func TestLessonsCheckBecomesContractItem(t *testing.T) {
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	check := `{"phase":"frame","predicate":"record-exists","params":{"kind":"risk"},"remediation":"scan risks first: wf risk scan"}`
	ev, err := c.Record("lesson", map[string]any{"text": "Always scan risks in frame", "status": "proposed", "check": check}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LessonsAccept(c, dir, specPath(t), ev.ID); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(c.Store.ContractsDir(), "lessons.yaml")
	raw, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("lessons.yaml missing: %v", err)
	}
	if !strings.Contains(string(raw), "lesson.always-scan-risks-in-frame") {
		t.Fatalf("namespaced lesson item missing:\n%s", raw)
	}
	// one representation, one evaluator: the merged spec sees the item
	sp, err := spec.LoadStrict(specPath(t), c.Store.ContractsDir())
	if err != nil {
		t.Fatalf("merged spec must load strictly: %v", err)
	}
	found := false
	for _, ci := range sp.Contracts {
		if ci.ID == "lesson.always-scan-risks-in-frame" {
			found = true
			if ci.Phase != "frame" || ci.Predicate != "record-exists" {
				t.Fatalf("item content wrong: %+v", ci)
			}
			if ci.Source != "lessons.yaml" {
				t.Fatalf("item source: %s", ci.Source)
			}
		}
	}
	if !found {
		t.Fatal("lesson item not in merged contracts")
	}
	// apply is idempotent
	if _, err := LessonsApply(c, dir, specPath(t)); err != nil {
		t.Fatal(err)
	}
	raw2, _ := os.ReadFile(target)
	if string(raw2) != string(raw) {
		t.Fatal("apply must be idempotent")
	}
}

func TestLessonsInvalidCheckRefusedAtAccept(t *testing.T) {
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	bad := []string{
		`{"phase":"frame","predicate":"always-true","params":{},"remediation":"x"}`, // unknown predicate
		`{"phase":"nirvana","predicate":"record-exists","params":{"kind":"risk"},"remediation":"x"}`, // unknown phase
		`not: [valid, contract`, // unparseable
	}
	for _, check := range bad {
		ev, err := c.Record("lesson", map[string]any{"text": "bad " + check[:8], "status": "proposed", "check": check}, false, "agent")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := LessonsAccept(c, dir, specPath(t), ev.ID); err == nil {
			t.Fatalf("accept must refuse invalid check: %s", check)
		}
	}
	// nothing written, no acceptance recorded
	if _, err := os.Stat(filepath.Join(c.Store.ContractsDir(), "lessons.yaml")); !os.IsNotExist(err) {
		t.Fatal("no lessons.yaml may exist after refused accepts")
	}
	r, _ := c.Store.LoadRun()
	env, _ := c.Env(r)
	for _, l := range env.Records("lesson") {
		if s, _ := l.Data["status"].(string); s == "accepted" {
			t.Fatal("refused lesson must stay proposed")
		}
	}
}

func TestLessonsProseSurvivesUserContentOutsideMarkers(t *testing.T) {
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	rulePath := filepath.Join(dir, ".claude", "rules", "wf-lessons.md")
	if err := os.MkdirAll(filepath.Dir(rulePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(rulePath, []byte("# my own notes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ev, _ := c.Record("lesson", map[string]any{"text": "prose one", "status": "proposed"}, false, "agent")
	if _, err := LessonsAccept(c, dir, specPath(t), ev.ID); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(rulePath)
	if !strings.Contains(string(raw), "# my own notes") || !strings.Contains(string(raw), "prose one") {
		t.Fatalf("user content must survive regeneration:\n%s", raw)
	}
	// rejecting the only prose lesson keeps the file (user content remains)
	if _, err := LessonsReject(c, dir, specPath(t), ev.ID); err != nil {
		t.Fatal(err)
	}
	raw, _ = os.ReadFile(rulePath)
	if !strings.Contains(string(raw), "# my own notes") || strings.Contains(string(raw), "prose one") {
		t.Fatalf("reject must clear the block but keep user content:\n%s", raw)
	}
}
