package doctor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func newCtl(t *testing.T) *runctl.Ctl {
	t.Helper()
	s, err := store.Open(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	sp, err := spec.Load(p, "")
	if err != nil {
		t.Fatal(err)
	}
	return &runctl.Ctl{Store: s, Spec: sp, Config: &store.Config{}}
}

// The dead-hooks incident: a run past Frame with a rich ledger but zero
// hook-captured events must be flagged loudly.
func TestHookLivenessDetection(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")

	// still in frame: no judgment yet
	if msg := HookLiveness(c, r); msg != "" {
		t.Fatalf("no warning before frame exit: %s", msg)
	}

	// simulate the power5 pattern: many manual events, frame exited, no hooks
	r.ExitedPh = []string{"frame"}
	r.Phase = "build"
	_ = c.Store.SaveRun(r)
	for i := 0; i < 16; i++ {
		if _, err := c.Record("assumption", map[string]any{"text": "x"}, false, "agent"); err != nil {
			t.Fatal(err)
		}
	}
	msg := HookLiveness(c, r)
	if !strings.Contains(msg, "DEAD") {
		t.Fatalf("all-manual ledger past frame must warn: %q", msg)
	}

	// one hook-captured event clears signal 2
	if _, err := c.Record("test-run", map[string]any{"cmd": "go test", "exit": 0, "grounded": true}, true, "hook"); err != nil {
		t.Fatal(err)
	}
	if msg := HookLiveness(c, r); msg != "" {
		t.Fatalf("hook activity present — no warning expected: %s", msg)
	}
}

// The power5 signature precisely: manual captures exist (masking signal 2)
// but every reviewer verdict is manual — signal 1 must still fire.
func TestHookLivenessManualVerdicts(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	r.ExitedPh = []string{"frame"}
	r.Phase = "build"
	_ = c.Store.SaveRun(r)
	// agent-piped captures carry actor=hook (indistinguishable by design)
	_, _ = c.Record("test-run", map[string]any{"cmd": "go test", "exit": 0, "grounded": true}, true, "hook")
	// manual verdicts only
	_, _ = c.Record("verdict", map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0}, false, "agent")
	_, _ = c.Record("verdict", map[string]any{"agent": "adversary", "status": "clean", "criticals": 0, "majors": 0, "scope": "red-team"}, false, "agent")
	msg := HookLiveness(c, r)
	if !strings.Contains(msg, "SubagentStop") {
		t.Fatalf("all-manual verdicts must trigger the verdict-gate warning: %q", msg)
	}
	// an auto-captured verdict clears it
	_, _ = c.Record("verdict", map[string]any{"agent": "auditor", "status": "clean", "criticals": 0, "majors": 0}, true, "hook")
	if msg := HookLiveness(c, r); msg != "" {
		t.Fatalf("auto verdict present — expected clear: %s", msg)
	}
}
